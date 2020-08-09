package generator

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/dave/dst"

	"github.com/myshkin5/moqueries/pkg/logs"
)

const (
	mockReceiverName    = "m"
	okName              = "ok"
	paramPrefix         = "param"
	paramsName          = "params"
	resultPrefix        = "result"
	resultsByParamsName = "resultsByParams"
	resultsName         = "results"
)

// Converter converts various interface and function types to AST structs to
// build a mock
type Converter struct{}

// NewConverter creates a new Converter
func NewConverter() *Converter {
	return &Converter{}
}

// Func holds on to function related data
type Func struct {
	Name    string
	Params  *dst.FieldList
	Results *dst.FieldList
}

// BaseStruct generates the base structure used to store the mock's state
func (c *Converter) BaseStruct(typeSpec *dst.TypeSpec, funcs []Func) *dst.GenDecl {
	mName := mockName(typeSpec.Name.Name)
	mock := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: dst.NewIdent(mName),
				Type: &dst.StructType{Fields: &dst.FieldList{
					List: baseMockFieldList(typeSpec, funcs),
				}},
			},
		},
	}
	mock.Decs.Before = dst.NewLine
	mock.Decs.Start.Append(fmt.Sprintf("// %s holds the state of a mock of the %s type", mName, typeSpec.Name.Name))

	return mock
}

// ParamResultStruct generates a structure for storing a set of parameters or
// a set of results for a method invocation of a mock
func (c *Converter) ParamResultStruct(typeName, prefix, label string, fieldList *dst.FieldList, comparable bool) *dst.GenDecl {
	var unnamedPrefix string
	switch label {
	case "params":
		unnamedPrefix = paramPrefix
	case "results":
		unnamedPrefix = resultPrefix
	default:
		logs.Panicf("Unknown label: %s", label)
	}
	fieldList = dst.Clone(fieldList).(*dst.FieldList)

	if fieldList.List != nil {
		for n, f := range fieldList.List {
			if len(f.Names) == 0 {
				f.Names = []*dst.Ident{dst.NewIdent(fmt.Sprintf("%s%d", unnamedPrefix, n+1))}
			}

			// Map params are represented as a deep hash
			if comparable && !isComparable(f.Type) {
				f.Type = &dst.Ident{
					Path: "github.com/myshkin5/moqueries/pkg/hash",
					Name: "Hash",
				}
			}
		}
	}

	structName := fmt.Sprintf("%s_%s", prefix, label)
	sType := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: dst.NewIdent(structName),
				Type: &dst.StructType{Fields: fieldList},
			},
		},
	}
	sType.Decs.Before = dst.NewLine
	sType.Decs.Start.Append(fmt.Sprintf("// %s holds the %s of the %s type",
		structName,
		label,
		typeName))

	return sType
}

// NewMockFn generates a function for constructing a mock
func (c *Converter) NewMockFn(typeSpec *dst.TypeSpec, funcs []Func) (funcDecl *dst.FuncDecl) {
	fnName := "newMock" + typeSpec.Name.Name
	mName := mockName(typeSpec.Name.Name)
	mockFn := &dst.FuncDecl{
		Name: dst.NewIdent(fnName),
		Type: &dst.FuncType{
			Params: &dst.FieldList{},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{
						Type: &dst.StarExpr{X: dst.NewIdent(mName)},
					},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{
						&dst.UnaryExpr{
							Op: token.AND,
							X: &dst.CompositeLit{
								Type: dst.NewIdent(mName),
								Elts: newMockElements(typeSpec, funcs),
								Decs: dst.CompositeLitDecorations{Lbrace: []string{"\n"}},
							},
						},
					},
				},
			},
		},
	}

	mockFn.Decs.Before = dst.NewLine
	mockFn.Decs.Start.Append(fmt.Sprintf("// %s creates a new mock of the %s type", fnName, typeSpec.Name.Name))

	return mockFn
}

// Method generates a mock implementation of a method
func (c *Converter) Method(typeName string, fn Func) *dst.FuncDecl {
	mName := mockName(typeName)
	typePrefix := fmt.Sprintf("%s_%s", mName, fn.Name)
	return &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(mockReceiverName)},
					Type:  &dst.StarExpr{X: dst.NewIdent(mName)},
				},
			},
		},
		Name: dst.NewIdent(fn.Name),
		Type: &dst.FuncType{
			Params:  dst.Clone(fn.Params).(*dst.FieldList),
			Results: dst.Clone(fn.Results).(*dst.FieldList),
		},
		Body: mockFunc(typePrefix, "_"+fn.Name, fn),
		Decs: dst.FuncDeclDecorations{
			NodeDecs: dst.NodeDecs{
				Before: dst.EmptyLine,
				After:  dst.EmptyLine,
			},
		},
	}
}

// FuncClosure generates a mock implementation of function type wrapped in a
// closure
func (c *Converter) FuncClosure(typeName, pkgPath string, fn Func) (funcDecl *dst.FuncDecl) {
	mName := mockName(typeName)
	return &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(mockReceiverName)},
					Type:  &dst.StarExpr{X: dst.NewIdent(mName)},
				},
			},
		},
		Name: dst.NewIdent("fn"),
		Type: &dst.FuncType{
			Params: &dst.FieldList{},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{
						Type: &dst.Ident{
							Path: pkgPath,
							Name: typeName,
						},
					},
				},
			},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{
						&dst.FuncLit{
							Type: &dst.FuncType{
								Params:  dst.Clone(fn.Params).(*dst.FieldList),
								Results: dst.Clone(fn.Results).(*dst.FieldList),
							},
							Body: mockFunc(mName, "", fn),
						},
					},
				},
			},
		},
		Decs: dst.FuncDeclDecorations{
			NodeDecs: dst.NodeDecs{
				Before: dst.EmptyLine,
				After:  dst.EmptyLine,
			},
		},
	}
}

func baseMockFieldList(typeSpec *dst.TypeSpec, funcs []Func) []*dst.Field {
	var fields []*dst.Field

	mName := mockName(typeSpec.Name.Name)
	for _, fn := range funcs {
		typePrefix := mName
		fieldSuffix := ""
		if _, ok := typeSpec.Type.(*dst.InterfaceType); ok {
			typePrefix = fmt.Sprintf("%s_%s", mName, fn.Name)
			fieldSuffix = "_" + fn.Name
		}
		fields = baseMockFieldsPerFn(fields, typePrefix, fieldSuffix)
	}

	return fields
}

func baseMockFieldsPerFn(fields []*dst.Field, typePrefix, fieldSuffix string) []*dst.Field {
	pName := fmt.Sprintf("%s_%s", typePrefix, paramsName)
	fields = append(fields, &dst.Field{
		Names: []*dst.Ident{
			dst.NewIdent(resultsByParamsName + fieldSuffix),
		},
		Type: &dst.MapType{
			Key:   &dst.Ident{Name: pName},
			Value: &dst.Ident{Name: fmt.Sprintf("%s_%s", typePrefix, resultsName)},
		},
	})

	fields = append(fields, &dst.Field{
		Names: []*dst.Ident{dst.NewIdent(paramsName + fieldSuffix)},
		Type: &dst.ChanType{
			Dir:   dst.SEND | dst.RECV,
			Value: &dst.Ident{Name: pName},
		},
	})

	return fields
}

func isComparable(expr dst.Expr) bool {
	// TODO this logic needs to be expanded -- also should check structs recursively
	switch expr.(type) {
	case *dst.ArrayType, *dst.MapType:
		return false
	}

	return true
}

func newMockElements(typeSpec *dst.TypeSpec, funcs []Func) []dst.Expr {
	var elems []dst.Expr

	mName := mockName(typeSpec.Name.Name)
	for _, fn := range funcs {
		typePrefix := mName
		fieldSuffix := ""
		if _, ok := typeSpec.Type.(*dst.InterfaceType); ok {
			typePrefix = fmt.Sprintf("%s_%s", mName, fn.Name)
			fieldSuffix = "_" + fn.Name
		}
		elems = newMockElement(elems, typePrefix, fieldSuffix)
	}

	return elems
}

func newMockElement(elems []dst.Expr, typePrefix, fieldSuffix string) []dst.Expr {
	pName := fmt.Sprintf("%s_%s", typePrefix, paramsName)
	elems = append(elems, &dst.KeyValueExpr{
		Key: dst.NewIdent(resultsByParamsName + fieldSuffix),
		Value: &dst.CompositeLit{
			Type: &dst.MapType{
				Key:   dst.NewIdent(pName),
				Value: dst.NewIdent(fmt.Sprintf("%s_%s", typePrefix, resultsName)),
			},
		},
		Decs: dst.KeyValueExprDecorations{
			NodeDecs: dst.NodeDecs{After: dst.NewLine},
		},
	})
	elems = append(elems, &dst.KeyValueExpr{
		Key: dst.NewIdent(paramsName + fieldSuffix),
		Value: &dst.CallExpr{
			Fun: dst.NewIdent("make"),
			Args: []dst.Expr{
				&dst.ChanType{
					Dir:   dst.SEND | dst.RECV,
					Value: dst.NewIdent(pName),
				},
				&dst.BasicLit{
					Kind:  token.INT,
					Value: "100",
				},
			},
		},
		Decs: dst.KeyValueExprDecorations{
			NodeDecs: dst.NodeDecs{After: dst.NewLine},
		},
	})
	return elems
}

func mockFunc(typePrefix, fieldSuffix string, fn Func) *dst.BlockStmt {
	return &dst.BlockStmt{
		List: []dst.Stmt{
			&dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent(paramsName)},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{
					&dst.CompositeLit{
						Type: dst.NewIdent(fmt.Sprintf("%s_%s", typePrefix, paramsName)),
						Elts: passthroughElements(fn.Params.List),
					},
				},
			},
			&dst.SendStmt{
				Chan: &dst.SelectorExpr{
					X:   dst.NewIdent(mockReceiverName),
					Sel: dst.NewIdent(paramsName + fieldSuffix),
				},
				Value: dst.NewIdent(paramsName),
			},
			&dst.AssignStmt{
				Lhs: []dst.Expr{
					dst.NewIdent(resultsName),
					dst.NewIdent(okName),
				},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{
					&dst.IndexExpr{
						X: &dst.SelectorExpr{
							X:   dst.NewIdent(mockReceiverName),
							Sel: dst.NewIdent(resultsByParamsName + fieldSuffix),
						},
						Index: dst.NewIdent(paramsName),
					},
				},
			},
			&dst.IfStmt{
				Cond: dst.NewIdent(okName),
				Body: &dst.BlockStmt{List: assignResults(fn.Results.List)},
			},
			returnStmt(fn.Results.List),
		},
	}
}

func passthroughElements(fields []*dst.Field) []dst.Expr {
	var elts []dst.Expr
	beforeDec := dst.NewLine
	for n, field := range fields {
		if len(field.Names) == 0 {
			pName := fmt.Sprintf("%s%d", paramPrefix, n+1)
			elts = append(elts, &dst.KeyValueExpr{
				Key:   dst.NewIdent(pName),
				Value: passthroughValue(pName, field),
				Decs: dst.KeyValueExprDecorations{
					NodeDecs: dst.NodeDecs{
						Before: beforeDec,
						After:  dst.NewLine,
					},
				},
			})
			beforeDec = dst.None
		}

		for _, name := range field.Names {
			elts = append(elts, &dst.KeyValueExpr{
				Key:   dst.NewIdent(name.Name),
				Value: passthroughValue(name.Name, field),
				Decs: dst.KeyValueExprDecorations{
					NodeDecs: dst.NodeDecs{
						Before: beforeDec,
						After:  dst.NewLine,
					},
				},
			})
			beforeDec = dst.None
		}
	}

	return elts
}

func passthroughValue(name string, field *dst.Field) dst.Expr {
	var val dst.Expr
	val = dst.NewIdent(name)
	if !isComparable(field.Type) {
		val = &dst.CallExpr{
			Fun: &dst.Ident{
				Path: "github.com/myshkin5/moqueries/pkg/hash",
				Name: "DeepHash",
			},
			Args: []dst.Expr{val},
		}
	}
	return val
}

func assignResults(results []*dst.Field) []dst.Stmt {
	var assigns []dst.Stmt
	for n, result := range results {
		if len(result.Names) == 0 {
			rName := fmt.Sprintf("%s%d", resultPrefix, n+1)
			assigns = append(assigns, &dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent(rName)},
				Tok: token.ASSIGN,
				Rhs: []dst.Expr{
					&dst.SelectorExpr{
						X:   dst.NewIdent(resultsName),
						Sel: dst.NewIdent(rName),
					},
				},
			})
		}

		for _, name := range result.Names {
			assigns = append(assigns, &dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent(name.Name)},
				Tok: token.ASSIGN,
				Rhs: []dst.Expr{
					&dst.SelectorExpr{
						X:   dst.NewIdent(resultsName),
						Sel: dst.NewIdent(name.Name),
					},
				},
			})
		}
	}
	return assigns
}

func returnStmt(results []*dst.Field) dst.Stmt {
	var exprs []dst.Expr
	for n, result := range results {
		if len(result.Names) == 0 {
			rName := fmt.Sprintf("%s%d", resultPrefix, n+1)
			exprs = append(exprs, dst.NewIdent(rName))
		}
		for _, name := range result.Names {
			exprs = append(exprs, dst.NewIdent(name.Name))
		}
	}
	return &dst.ReturnStmt{Results: exprs}
}

func mockName(typeName string) string {
	return "mock" + strings.Title(typeName)
}
