package generator

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/dave/dst"

	"github.com/myshkin5/moqueries/pkg/logs"
)

const (
	fnRecorderSuffix      = "fnRecorder"
	fnFnName              = "fn"
	mockFnName            = "mock"
	mockIdent             = "mock"
	mockReceiverIdent     = "m"
	okIdent               = "ok"
	onCallIdent           = "onCall"
	paramPrefix           = "param"
	paramsIdent           = "params"
	recorderIdent         = "recorder"
	recorderReceiverIdent = "r"
	returnFnName          = "ret"
	resultPrefix          = "result"
	resultsByParamsIdent  = "resultsByParams"
	resultsIdent          = "results"
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

// IsolationStruct generates a struct used to isolate an interface for the mock
func (c *Converter) IsolationStruct(typeName, suffix string) (structDecl *dst.GenDecl) {
	mName := mockName(typeName)
	iName := fmt.Sprintf("%s_%s", mName, suffix)
	isolate := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: dst.NewIdent(iName),
				Type: &dst.StructType{Fields: &dst.FieldList{
					List: []*dst.Field{
						{
							Names: []*dst.Ident{dst.NewIdent(mockIdent)},
							Type: &dst.StarExpr{
								X: dst.NewIdent(mName),
							},
						},
					},
				}},
			},
		},
	}
	isolate.Decs.Before = dst.NewLine
	isolate.Decs.Start.Append(fmt.Sprintf("// %s isolates the %s interface of the %s type", iName, suffix, typeName))

	return isolate
}

// MethodStructs generates a structure for storing a set of parameters or
// a set of results for a method invocation of a mock
func (c *Converter) MethodStructs(typeName, prefix string, fn Func) []dst.Decl {
	decls := make([]dst.Decl, 3)
	decls[0] = methodStruct(typeName, prefix, paramsIdent, fn.Params)
	decls[1] = methodStruct(typeName, prefix, resultsIdent, fn.Results)
	decls[2] = fnRecorderStruct(typeName, prefix)
	return decls
}

// NewFunc generates a function for constructing a mock
func (c *Converter) NewFunc(typeSpec *dst.TypeSpec, funcs []Func) (funcDecl *dst.FuncDecl) {
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

// IsolationAccessor generates a function to access an isolation interface
func (c *Converter) IsolationAccessor(typeName, suffix, fnName string) (funcDecl *dst.FuncDecl) {
	mName := mockName(typeName)
	iName := fmt.Sprintf("%s_%s", mName, suffix)
	mockFn := &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(mockReceiverIdent)},
					Type:  &dst.StarExpr{X: dst.NewIdent(mName)},
				},
			},
		},
		Name: dst.NewIdent(fnName),
		Type: &dst.FuncType{
			Params: &dst.FieldList{},
			Results: &dst.FieldList{
				List: []*dst.Field{
					{
						Type: &dst.StarExpr{X: dst.NewIdent(iName)},
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
								Type: dst.NewIdent(iName),
								Elts: []dst.Expr{
									&dst.KeyValueExpr{
										Key:   dst.NewIdent(mockIdent),
										Value: dst.NewIdent(mockReceiverIdent),
										Decs: dst.KeyValueExprDecorations{
											NodeDecs: dst.NodeDecs{After: dst.NewLine},
										},
									},
								},
								Decs: dst.CompositeLitDecorations{Lbrace: []string{"\n"}},
							},
						},
					},
				},
			},
		},
	}

	mockFn.Decs.Before = dst.NewLine
	mockFn.Decs.Start.Append(fmt.Sprintf("// %s returns the %s implementation of the %s type",
		fnName, suffix, typeName))

	return mockFn
}

// FuncClosure generates a mock implementation of function type wrapped in a
// closure
func (c *Converter) FuncClosure(typeName, pkgPath string, fn Func) (funcDecl *dst.FuncDecl) {
	mName := mockName(typeName)
	mockFn := &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(mockReceiverIdent)},
					Type:  &dst.StarExpr{X: dst.NewIdent(mName)},
				},
			},
		},
		Name: dst.NewIdent(mockFnName),
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
							Body: &dst.BlockStmt{
								List: []dst.Stmt{
									&dst.AssignStmt{
										Lhs: []dst.Expr{dst.NewIdent(mockIdent)},
										Tok: token.DEFINE,
										Rhs: []dst.Expr{
											&dst.UnaryExpr{
												Op: token.AND,
												X: &dst.CompositeLit{
													Type: dst.NewIdent(fmt.Sprintf("%s_%s", mName, mockIdent)),
													Elts: []dst.Expr{
														&dst.KeyValueExpr{
															Key:   dst.NewIdent(mockIdent),
															Value: dst.NewIdent(mockReceiverIdent),
														},
													},
												},
											},
										},
									},
									&dst.ReturnStmt{
										Results: []dst.Expr{
											&dst.CallExpr{
												Fun: &dst.SelectorExpr{
													X:   dst.NewIdent(mockIdent),
													Sel: dst.NewIdent(fnFnName),
												},
												Args: passthroughFields(fn.Params),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Decs: stdFuncDec(),
	}

	mockFn.Decs.Before = dst.NewLine
	mockFn.Decs.Start.Append(fmt.Sprintf("// %s returns the %s implementation of the %s type",
		mockFnName, mockIdent, typeName))

	return mockFn
}

// MockMethod generates a mock implementation of a method
func (c *Converter) MockMethod(typeName string, fn Func) *dst.FuncDecl {
	mName := mockName(typeName)
	recv := fmt.Sprintf("%s_%s", mName, mockIdent)

	fnName := fn.Name
	fieldSuffix := "_" + fn.Name
	typePrefix := fmt.Sprintf("%s_%s", mName, fn.Name)
	if fnName == "" {
		fnName = "fn"
		fieldSuffix = ""
		typePrefix = mName
	}

	return &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(mockReceiverIdent)},
					Type:  &dst.StarExpr{X: dst.NewIdent(recv)},
				},
			},
		},
		Name: dst.NewIdent(fnName),
		Type: &dst.FuncType{
			Params:  cloneAndNameUnnamed(paramPrefix, fn.Params),
			Results: cloneAndNameUnnamed(resultPrefix, fn.Results),
		},
		Body: mockFunc(typePrefix, fieldSuffix, fn),
		Decs: stdFuncDec(),
	}
}

// RecorderMethods generates a recorder implementation of a method and associated
// return method
func (c *Converter) RecorderMethods(typeName string, fn Func) (funcDecls []dst.Decl) {
	return []dst.Decl{
		recorderFn(typeName, fn),
		recorderReturnFn(typeName, fn),
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
	pName := fmt.Sprintf("%s_%s", typePrefix, paramsIdent)
	fields = append(fields, &dst.Field{
		Names: []*dst.Ident{
			dst.NewIdent(resultsByParamsIdent + fieldSuffix),
		},
		Type: &dst.MapType{
			Key:   &dst.Ident{Name: pName},
			Value: &dst.Ident{Name: fmt.Sprintf("%s_%s", typePrefix, resultsIdent)},
		},
	})

	fields = append(fields, &dst.Field{
		Names: []*dst.Ident{dst.NewIdent(paramsIdent + fieldSuffix)},
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

func methodStruct(typeName, prefix, label string, fieldList *dst.FieldList) *dst.GenDecl {
	unnamedPrefix, comparable := labelDirection(label)
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

func fnRecorderStruct(typeName string, prefix string) *dst.GenDecl {
	mName := mockName(typeName)
	structName := fmt.Sprintf("%s_%s", prefix, fnRecorderSuffix)
	sType := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: dst.NewIdent(structName),
				Type: &dst.StructType{Fields: &dst.FieldList{List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent(paramsIdent)},
						Type:  dst.NewIdent(fmt.Sprintf("%s_%s", prefix, paramsIdent)),
					},
					{
						Names: []*dst.Ident{dst.NewIdent(mockIdent)},
						Type:  &dst.StarExpr{X: dst.NewIdent(mName)},
					},
				}}},
			},
		},
	}
	sType.Decs.Before = dst.NewLine
	sType.Decs.Start.Append(fmt.Sprintf("// %s routes recorded function calls to the %s mock",
		structName,
		mName))

	return sType
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
	pName := fmt.Sprintf("%s_%s", typePrefix, paramsIdent)
	elems = append(elems, &dst.KeyValueExpr{
		Key: dst.NewIdent(resultsByParamsIdent + fieldSuffix),
		Value: &dst.CompositeLit{
			Type: &dst.MapType{
				Key:   dst.NewIdent(pName),
				Value: dst.NewIdent(fmt.Sprintf("%s_%s", typePrefix, resultsIdent)),
			},
		},
		Decs: dst.KeyValueExprDecorations{
			NodeDecs: dst.NodeDecs{After: dst.NewLine},
		},
	})
	elems = append(elems, &dst.KeyValueExpr{
		Key: dst.NewIdent(paramsIdent + fieldSuffix),
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
	stateSelector := &dst.SelectorExpr{
		X:   dst.NewIdent(mockReceiverIdent),
		Sel: dst.NewIdent(mockIdent),
	}
	return &dst.BlockStmt{
		List: []dst.Stmt{
			&dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent(paramsIdent)},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{
					&dst.CompositeLit{
						Type: dst.NewIdent(fmt.Sprintf("%s_%s", typePrefix, paramsIdent)),
						Elts: passthroughElements(fn.Params.List, paramsIdent),
					},
				},
			},
			&dst.SendStmt{
				Chan: &dst.SelectorExpr{
					X:   dst.Clone(stateSelector).(dst.Expr),
					Sel: dst.NewIdent(paramsIdent + fieldSuffix),
				},
				Value: dst.NewIdent(paramsIdent),
			},
			&dst.AssignStmt{
				Lhs: []dst.Expr{
					dst.NewIdent(resultsIdent),
					dst.NewIdent(okIdent),
				},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{
					&dst.IndexExpr{
						X: &dst.SelectorExpr{
							X:   dst.Clone(stateSelector).(dst.Expr),
							Sel: dst.NewIdent(resultsByParamsIdent + fieldSuffix),
						},
						Index: dst.NewIdent(paramsIdent),
					},
				},
			},
			&dst.IfStmt{
				Cond: dst.NewIdent(okIdent),
				Body: &dst.BlockStmt{List: assignResults(fn.Results.List)},
			},
			mockFnReturn(fn.Results.List),
		},
	}
}

func recorderFn(typeName string, fn Func) *dst.FuncDecl {
	mName := mockName(typeName)

	recvType := fmt.Sprintf("%s_%s", mName, recorderIdent)
	fnName := fn.Name
	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	typePrefix := fmt.Sprintf("%s_%s", mName, fn.Name)
	var mockVal dst.Expr = &dst.SelectorExpr{
		X:   dst.NewIdent(mockReceiverIdent),
		Sel: dst.NewIdent(mockIdent),
	}
	if fn.Name == "" {
		recvType = mName
		fnName = onCallIdent
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
		typePrefix = mName
		mockVal = dst.NewIdent(mockReceiverIdent)
	}

	return &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(mockReceiverIdent)},
					Type:  &dst.StarExpr{X: dst.NewIdent(recvType)},
				},
			},
		},
		Name: dst.NewIdent(fnName),
		Type: &dst.FuncType{
			Params: dst.Clone(fn.Params).(*dst.FieldList),
			Results: &dst.FieldList{
				List: []*dst.Field{
					{
						Type: &dst.StarExpr{
							X: dst.NewIdent(fnRecName),
						},
					},
				},
			},
		},
		Body: recorderFnInterfaceBody(fnRecName, typePrefix, mockVal, fn),
		Decs: stdFuncDec(),
	}
}

func recorderFnInterfaceBody(fnRecName, typePrefix string, mockValue dst.Expr, fn Func) *dst.BlockStmt {
	return &dst.BlockStmt{
		List: []dst.Stmt{
			&dst.ReturnStmt{
				Results: []dst.Expr{
					&dst.UnaryExpr{
						Op: token.AND,
						X: &dst.CompositeLit{
							Type: dst.NewIdent(fnRecName),
							Elts: []dst.Expr{
								&dst.KeyValueExpr{
									Key: dst.NewIdent(paramsIdent),
									Value: &dst.CompositeLit{
										Type: dst.NewIdent(fmt.Sprintf("%s_%s", typePrefix, paramsIdent)),
										Elts: passthroughElements(fn.Params.List, paramsIdent),
									},
									Decs: dst.KeyValueExprDecorations{
										NodeDecs: dst.NodeDecs{After: dst.NewLine},
									},
								},
								&dst.KeyValueExpr{
									Key:   dst.NewIdent(mockIdent),
									Value: mockValue,
									Decs: dst.KeyValueExprDecorations{
										NodeDecs: dst.NodeDecs{After: dst.NewLine},
									},
								},
							},
							Decs: dst.CompositeLitDecorations{Lbrace: []string{"\n"}},
						},
					},
				},
			},
		},
	}
}

func recorderReturnFn(typeName string, fn Func) *dst.FuncDecl {
	mName := mockName(typeName)

	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	resultsByParams := fmt.Sprintf("%s_%s", resultsByParamsIdent, fn.Name)
	results := fmt.Sprintf("%s_%s_%s", mName, fn.Name, resultsIdent)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
		resultsByParams = resultsByParamsIdent
		results = fmt.Sprintf("%s_%s", mName, resultsIdent)
	}

	return &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(recorderReceiverIdent)},
					Type:  &dst.StarExpr{X: dst.NewIdent(fnRecName)},
				},
			},
		},
		Name: dst.NewIdent(returnFnName),
		Type: &dst.FuncType{
			Params: cloneAndNameUnnamed(resultPrefix, fn.Results),
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.AssignStmt{
					Lhs: []dst.Expr{
						&dst.IndexExpr{
							X: &dst.SelectorExpr{
								X: &dst.SelectorExpr{
									X:   dst.NewIdent(recorderReceiverIdent),
									Sel: dst.NewIdent(mockIdent),
								},
								Sel: dst.NewIdent(resultsByParams),
							},
							Index: &dst.SelectorExpr{
								X:   dst.NewIdent(recorderReceiverIdent),
								Sel: dst.NewIdent(paramsIdent),
							},
						},
					},
					Tok: token.ASSIGN,
					Rhs: []dst.Expr{
						&dst.CompositeLit{
							Type: dst.NewIdent(results),
							Elts: passthroughElements(fn.Results.List, resultsIdent),
						},
					},
				},
			},
		},
		Decs: stdFuncDec(),
	}
}

func passthroughElements(fields []*dst.Field, label string) []dst.Expr {
	unnamedPrefix, comparable := labelDirection(label)
	var elts []dst.Expr
	beforeDec := dst.NewLine
	for n, field := range fields {
		if len(field.Names) == 0 {
			pName := fmt.Sprintf("%s%d", unnamedPrefix, n+1)
			elts = append(elts, &dst.KeyValueExpr{
				Key:   dst.NewIdent(pName),
				Value: passthroughValue(pName, field, comparable),
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
				Value: passthroughValue(name.Name, field, comparable),
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

func passthroughValue(name string, field *dst.Field, comparable bool) dst.Expr {
	var val dst.Expr
	val = dst.NewIdent(name)
	if comparable && !isComparable(field.Type) {
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

func passthroughFields(fields *dst.FieldList) []dst.Expr {
	var exprs []dst.Expr
	for _, f := range fields.List {
		for _, n := range f.Names {
			exprs = append(exprs, dst.NewIdent(n.Name))
		}
	}
	return exprs
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
						X:   dst.NewIdent(resultsIdent),
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
						X:   dst.NewIdent(resultsIdent),
						Sel: dst.NewIdent(name.Name),
					},
				},
			})
		}
	}
	return assigns
}

func cloneAndNameUnnamed(prefix string, fieldList *dst.FieldList) *dst.FieldList {
	fieldList = dst.Clone(fieldList).(*dst.FieldList)
	for n, f := range fieldList.List {
		if len(f.Names) == 0 {
			f.Names = []*dst.Ident{dst.NewIdent(fmt.Sprintf("%s%d", prefix, n+1))}
		}
	}
	return fieldList
}

func mockFnReturn(results []*dst.Field) dst.Stmt {
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
	return mockIdent + strings.Title(typeName)
}

func stdFuncDec() dst.FuncDeclDecorations {
	return dst.FuncDeclDecorations{
		NodeDecs: dst.NodeDecs{
			Before: dst.EmptyLine,
			After:  dst.EmptyLine,
		},
	}
}

func labelDirection(label string) (unnamedPrefix string, comparable bool) {
	switch label {
	case paramsIdent:
		unnamedPrefix = paramPrefix
		comparable = true
	case resultsIdent:
		unnamedPrefix = resultPrefix
		comparable = false
	default:
		logs.Panicf("Unknown label: %s", label)
	}

	return unnamedPrefix, comparable
}
