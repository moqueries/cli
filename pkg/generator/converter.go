package generator

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/dave/dst"

	"github.com/myshkin5/moqueries/pkg/logs"
)

const (
	moqueriesPkg = "github.com/myshkin5/moqueries"
	hashPkg      = moqueriesPkg + "/pkg/hash"
	testingPkg   = moqueriesPkg + "/pkg/testing"

	moqTType = "MoqT"

	anyTimesIdent         = "anyTimes"
	countIdent            = "count"
	iIdent                = "i"
	indexIdent            = "index"
	lastIdent             = "last"
	mockIdent             = "mock"
	mockReceiverIdent     = "m"
	okIdent               = "ok"
	paramsIdent           = "params"
	recorderIdent         = "recorder"
	recorderReceiverIdent = "r"
	resultsByParamsIdent  = "resultsByParams"
	resultIdent           = "result"
	resultsIdent          = "results"
	testingIdent          = "t"

	anyTimesFnName = "anyTimes"
	fnFnName       = "fn"
	mockFnName     = "mock"
	onCallFnName   = "onCall"
	returnFnName   = "returnResults"
	timesFnName    = "times"

	fnRecorderSuffix = "fnRecorder"
	paramPrefix      = "param"
	resultMgrSuffix  = "resultMgr"
	resultPrefix     = "result"
)

// Converter converts various interface and function types to AST structs to
// build a mock
type Converter struct {
	export bool
}

// NewConverter creates a new Converter
func NewConverter(export bool) *Converter {
	return &Converter{
		export: export,
	}
}

// Func holds on to function related data
type Func struct {
	Name    string
	Params  *dst.FieldList
	Results *dst.FieldList
}

// BaseStruct generates the base structure used to store the mock's state
func (c *Converter) BaseStruct(typeSpec *dst.TypeSpec, funcs []Func) *dst.GenDecl {
	mName := c.mockName(typeSpec.Name.Name)
	mock := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: dst.NewIdent(mName),
				Type: &dst.StructType{Fields: &dst.FieldList{
					List: c.baseMockFieldList(typeSpec, funcs),
				}},
			},
		},
	}
	mock.Decs.Before = dst.NewLine
	mock.Decs.Start.Append(fmt.Sprintf(
		"// %s holds the state of a mock of the %s type",
		mName,
		typeSpec.Name.Name))

	return mock
}

// IsolationStruct generates a struct used to isolate an interface for the mock
func (c *Converter) IsolationStruct(typeName, suffix string) (structDecl *dst.GenDecl) {
	mName := c.mockName(typeName)
	iName := fmt.Sprintf("%s_%s", mName, suffix)
	isolate := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: dst.NewIdent(iName),
				Type: &dst.StructType{Fields: &dst.FieldList{
					List: []*dst.Field{{
						Names: []*dst.Ident{
							dst.NewIdent(c.exportName(mockIdent)),
						},
						Type: &dst.StarExpr{X: dst.NewIdent(mName)},
					}},
				}},
			},
		},
	}
	isolate.Decs.Before = dst.NewLine
	isolate.Decs.Start.Append(fmt.Sprintf(
		"// %s isolates the %s interface of the %s type",
		iName,
		suffix,
		typeName))

	return isolate
}

// MethodStructs generates a structure for storing a set of parameters or
// a set of results for a method invocation of a mock
func (c *Converter) MethodStructs(typeSpec *dst.TypeSpec, fn Func) []dst.Decl {
	prefix := c.mockName(typeSpec.Name.Name)
	if _, ok := typeSpec.Type.(*dst.InterfaceType); ok {
		prefix = fmt.Sprintf("%s_%s", prefix, fn.Name)
	}

	var decls []dst.Decl
	decls = append(decls,
		c.methodStruct(typeSpec.Name.Name, prefix, paramsIdent, fn.Params))
	decls = append(decls, c.resultMgrStruct(typeSpec.Name.Name, prefix))
	decls = append(decls,
		c.methodStruct(typeSpec.Name.Name, prefix, resultsIdent, fn.Results))
	decls = append(decls, c.fnRecorderStruct(typeSpec.Name.Name, prefix))
	return decls
}

// NewFunc generates a function for constructing a mock
func (c *Converter) NewFunc(typeSpec *dst.TypeSpec, funcs []Func) (
	funcDecl *dst.FuncDecl,
) {
	fnName := c.exportName("newMock" + typeSpec.Name.Name)
	mName := c.mockName(typeSpec.Name.Name)
	mockFn := &dst.FuncDecl{
		Name: dst.NewIdent(fnName),
		Type: &dst.FuncType{
			Params: &dst.FieldList{
				List: []*dst.Field{{
					Names: []*dst.Ident{dst.NewIdent(testingIdent)},
					Type:  &dst.Ident{Name: moqTType, Path: testingPkg},
				}},
			},
			Results: &dst.FieldList{List: []*dst.Field{{
				Type: &dst.StarExpr{X: dst.NewIdent(mName)},
			}}},
		},
		Body: &dst.BlockStmt{
			List: []dst.Stmt{
				&dst.ReturnStmt{
					Results: []dst.Expr{
						&dst.UnaryExpr{
							Op: token.AND,
							X: &dst.CompositeLit{
								Type: dst.NewIdent(mName),
								Elts: c.newMockElements(typeSpec, funcs),
								Decs: dst.CompositeLitDecorations{
									Lbrace: []string{"\n"},
								},
							},
						},
					},
				},
			},
		},
	}

	mockFn.Decs.Before = dst.NewLine
	mockFn.Decs.Start.Append(fmt.Sprintf(
		"// %s creates a new mock of the %s type", fnName, typeSpec.Name.Name))

	return mockFn
}

// IsolationAccessor generates a function to access an isolation interface
func (c *Converter) IsolationAccessor(typeName, suffix, fnName string) (
	funcDecl *dst.FuncDecl,
) {
	fnName = c.exportName(fnName)
	mName := c.mockName(typeName)
	iName := fmt.Sprintf("%s_%s", mName, suffix)
	mockFn := &dst.FuncDecl{
		Recv: &dst.FieldList{List: []*dst.Field{{
			Names: []*dst.Ident{dst.NewIdent(mockReceiverIdent)},
			Type:  &dst.StarExpr{X: dst.NewIdent(mName)},
		}}},
		Name: dst.NewIdent(fnName),
		Type: &dst.FuncType{
			Params: &dst.FieldList{},
			Results: &dst.FieldList{List: []*dst.Field{{
				Type: &dst.StarExpr{X: dst.NewIdent(iName)},
			}}},
		},
		Body: &dst.BlockStmt{List: []dst.Stmt{&dst.ReturnStmt{
			Results: []dst.Expr{&dst.UnaryExpr{
				Op: token.AND,
				X: &dst.CompositeLit{
					Type: dst.NewIdent(iName),
					Elts: []dst.Expr{&dst.KeyValueExpr{
						Key:   dst.NewIdent(c.exportName(mockIdent)),
						Value: dst.NewIdent(mockReceiverIdent),
						Decs: dst.KeyValueExprDecorations{
							NodeDecs: dst.NodeDecs{After: dst.NewLine},
						},
					}},
					Decs: dst.CompositeLitDecorations{Lbrace: []string{"\n"}},
				},
			}},
		}}},
	}

	mockFn.Decs.Before = dst.NewLine
	mockFn.Decs.Start.Append(fmt.Sprintf(
		"// %s returns the %s implementation of the %s type",
		fnName,
		suffix,
		typeName))

	return mockFn
}

// FuncClosure generates a mock implementation of function type wrapped in a
// closure
func (c *Converter) FuncClosure(typeName, pkgPath string, fn Func) (
	funcDecl *dst.FuncDecl,
) {
	mName := c.mockName(typeName)
	ellipsis := false
	if len(fn.Params.List) > 0 {
		if _, ok := fn.Params.List[len(fn.Params.List)-1].Type.(*dst.Ellipsis); ok {
			ellipsis = true
		}
	}
	fnLitCall := &dst.CallExpr{
		Fun: &dst.SelectorExpr{
			X:   dst.NewIdent(mockIdent),
			Sel: dst.NewIdent(c.exportName(fnFnName)),
		},
		Args:     passthroughFields(fn.Params),
		Ellipsis: ellipsis,
	}
	var fnLitRetStmt dst.Stmt
	fnLitRetStmt = &dst.ReturnStmt{Results: []dst.Expr{fnLitCall}}
	if fn.Results == nil {
		fnLitRetStmt = &dst.ExprStmt{X: fnLitCall}
	}

	mockFn := &dst.FuncDecl{
		Recv: &dst.FieldList{List: []*dst.Field{{
			Names: []*dst.Ident{dst.NewIdent(mockReceiverIdent)},
			Type:  &dst.StarExpr{X: dst.NewIdent(mName)},
		}}},
		Name: dst.NewIdent(c.exportName(mockFnName)),
		Type: &dst.FuncType{
			Params: &dst.FieldList{},
			Results: &dst.FieldList{List: []*dst.Field{{
				Type: &dst.Ident{Path: pkgPath, Name: typeName},
			}}},
		},
		Body: &dst.BlockStmt{List: []dst.Stmt{&dst.ReturnStmt{
			Results: []dst.Expr{&dst.FuncLit{
				Type: &dst.FuncType{
					Params:  dst.Clone(fn.Params).(*dst.FieldList),
					Results: cloneNilableFieldList(fn.Results),
				},
				Body: &dst.BlockStmt{List: []dst.Stmt{&dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent(mockIdent)},
					Tok: token.DEFINE,
					Rhs: []dst.Expr{&dst.UnaryExpr{
						Op: token.AND,
						X: &dst.CompositeLit{
							Type: dst.NewIdent(fmt.Sprintf(
								"%s_%s", mName, mockIdent)),
							Elts: []dst.Expr{&dst.KeyValueExpr{
								Key:   dst.NewIdent(c.exportName(mockIdent)),
								Value: dst.NewIdent(mockReceiverIdent),
							}},
						},
					}}},
					fnLitRetStmt,
				}},
			}},
		}}},
		Decs: stdFuncDec(),
	}

	mockFn.Decs.Before = dst.NewLine
	mockFn.Decs.Start.Append(fmt.Sprintf(
		"// %s returns the %s implementation of the %s type",
		c.exportName(mockFnName),
		mockIdent,
		typeName))

	return mockFn
}

// MockMethod generates a mock implementation of a method
func (c *Converter) MockMethod(typeName string, fn Func) *dst.FuncDecl {
	mName := c.mockName(typeName)
	recv := fmt.Sprintf("%s_%s", mName, mockIdent)

	fnName := fn.Name
	fieldSuffix := "_" + fn.Name
	typePrefix := fmt.Sprintf("%s_%s", mName, fn.Name)
	if fnName == "" {
		fnName = c.exportName(fnFnName)
		fieldSuffix = ""
		typePrefix = mName
	}

	return &dst.FuncDecl{
		Recv: &dst.FieldList{List: []*dst.Field{{
			Names: []*dst.Ident{dst.NewIdent(mockReceiverIdent)},
			Type:  &dst.StarExpr{X: dst.NewIdent(recv)},
		}}},
		Name: dst.NewIdent(fnName),
		Type: &dst.FuncType{
			Params:  cloneAndNameUnnamed(paramPrefix, fn.Params),
			Results: cloneAndNameUnnamed(resultPrefix, fn.Results),
		},
		Body: c.mockFunc(typePrefix, fieldSuffix, fn),
		Decs: stdFuncDec(),
	}
}

// RecorderMethods generates a recorder implementation of a method and
// associated return method
func (c *Converter) RecorderMethods(typeName string, fn Func) (
	funcDecls []dst.Decl,
) {
	return []dst.Decl{
		c.recorderFn(typeName, fn),
		c.recorderReturnFn(typeName, fn),
		c.recorderTimesFn(typeName, fn),
		c.recorderAnyTimesFn(typeName, fn),
	}
}

func (c *Converter) baseMockFieldList(
	typeSpec *dst.TypeSpec, funcs []Func,
) []*dst.Field {
	var fields []*dst.Field

	fields = append(fields, &dst.Field{
		Names: []*dst.Ident{dst.NewIdent(c.exportName(testingIdent))},
		Type:  &dst.Ident{Name: moqTType, Path: testingPkg},
	})

	mName := c.mockName(typeSpec.Name.Name)
	for _, fn := range funcs {
		typePrefix := mName
		fieldSuffix := ""
		if _, ok := typeSpec.Type.(*dst.InterfaceType); ok {
			typePrefix = fmt.Sprintf("%s_%s", mName, fn.Name)
			fieldSuffix = "_" + fn.Name
		}
		fields = c.baseMockFieldsPerFn(fields, typePrefix, fieldSuffix)
	}

	return fields
}

func (c *Converter) baseMockFieldsPerFn(
	fields []*dst.Field, typePrefix, fieldSuffix string,
) []*dst.Field {
	pName := c.exportName(fmt.Sprintf("%s_%s", typePrefix, paramsIdent))
	fields = append(fields, &dst.Field{
		Names: []*dst.Ident{
			dst.NewIdent(c.exportName(resultsByParamsIdent + fieldSuffix)),
		},
		Type: &dst.MapType{
			Key: dst.NewIdent(pName),
			Value: &dst.StarExpr{
				X: dst.NewIdent(fmt.Sprintf(
					"%s_%s", typePrefix, resultMgrSuffix)),
			},
		},
	})

	fields = append(fields, &dst.Field{
		Names: []*dst.Ident{
			dst.NewIdent(c.exportName(paramsIdent + fieldSuffix)),
		},
		Type: &dst.ChanType{
			Dir:   dst.SEND | dst.RECV,
			Value: dst.NewIdent(pName),
		},
	})

	return fields
}

func isComparable(expr dst.Expr) bool {
	// TODO this logic needs to be expanded -- also should check structs recursively
	switch expr.(type) {
	case *dst.ArrayType, *dst.MapType, *dst.Ellipsis:
		return false
	}

	return true
}

func (c *Converter) methodStruct(
	typeName, prefix, label string, fieldList *dst.FieldList,
) *dst.GenDecl {
	unnamedPrefix, comparable := labelDirection(label)
	fieldList = cloneNilableFieldList(fieldList)

	if fieldList == nil {
		// Result field lists can be nil (rather than containing an empty
		// list). Struct field lists cannot be nil.
		fieldList = &dst.FieldList{}
	} else {
		for n, f := range fieldList.List {
			if len(f.Names) == 0 {
				f.Names = []*dst.Ident{
					dst.NewIdent(fmt.Sprintf("%s%d", unnamedPrefix, n+1)),
				}
			}

			for nn := range f.Names {
				f.Names[nn] = dst.NewIdent(c.exportName(f.Names[nn].Name))
			}

			// Map params are represented as a deep hash
			if comparable && !isComparable(f.Type) {
				f.Type = &dst.Ident{Path: hashPkg, Name: "Hash"}
			}
		}
	}

	structName := fmt.Sprintf("%s_%s", prefix, label)
	sType := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{&dst.TypeSpec{
			Name: dst.NewIdent(structName),
			Type: &dst.StructType{Fields: fieldList},
		}},
	}
	sType.Decs.Before = dst.NewLine
	sType.Decs.Start.Append(fmt.Sprintf("// %s holds the %s of the %s type",
		structName,
		label,
		typeName))

	return sType
}

func (c *Converter) resultMgrStruct(typeName, prefix string) *dst.GenDecl {
	structName := fmt.Sprintf("%s_%s", prefix, resultMgrSuffix)

	sType := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{&dst.TypeSpec{
			Name: dst.NewIdent(structName),
			Type: &dst.StructType{Fields: &dst.FieldList{List: []*dst.Field{
				{
					Names: []*dst.Ident{
						dst.NewIdent(c.exportName(resultsIdent)),
					},
					Type: &dst.ArrayType{Elt: &dst.StarExpr{
						X: dst.NewIdent(c.exportName(fmt.Sprintf(
							"%s_%s", prefix, resultsIdent))),
					}},
				},
				{
					Names: []*dst.Ident{
						dst.NewIdent(c.exportName(indexIdent)),
					},
					Type: dst.NewIdent("uint32"),
				},
				{
					Names: []*dst.Ident{
						dst.NewIdent(c.exportName(anyTimesIdent)),
					},
					Type: dst.NewIdent("bool"),
				},
			}}},
		}},
	}
	sType.Decs.Before = dst.NewLine
	sType.Decs.Start.Append(fmt.Sprintf(
		"// %s manages multiple results and the state of the %s type",
		structName,
		typeName))

	return sType
}

func (c *Converter) fnRecorderStruct(
	typeName string, prefix string,
) *dst.GenDecl {
	mName := c.mockName(typeName)
	structName := fmt.Sprintf("%s_%s", prefix, fnRecorderSuffix)
	sType := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{&dst.TypeSpec{
			Name: dst.NewIdent(structName),
			Type: &dst.StructType{Fields: &dst.FieldList{List: []*dst.Field{
				{
					Names: []*dst.Ident{
						dst.NewIdent(c.exportName(paramsIdent)),
					},
					Type: dst.NewIdent(fmt.Sprintf(
						"%s_%s", prefix, paramsIdent)),
				},
				{
					Names: []*dst.Ident{
						dst.NewIdent(c.exportName(resultsIdent)),
					},
					Type: &dst.StarExpr{X: dst.NewIdent(fmt.Sprintf(
						"%s_%s", prefix, resultMgrSuffix))},
				},
				{
					Names: []*dst.Ident{dst.NewIdent(c.exportName(mockIdent))},
					Type:  &dst.StarExpr{X: dst.NewIdent(mName)},
				},
			}}},
		}},
	}
	sType.Decs.Before = dst.NewLine
	sType.Decs.Start.Append(fmt.Sprintf(
		"// %s routes recorded function calls to the %s mock",
		structName,
		mName))

	return sType
}

func (c *Converter) newMockElements(
	typeSpec *dst.TypeSpec, funcs []Func,
) []dst.Expr {
	var elems []dst.Expr

	elems = append(elems, &dst.KeyValueExpr{
		Key:   dst.NewIdent(c.exportName(testingIdent)),
		Value: dst.NewIdent(testingIdent),
		Decs: dst.KeyValueExprDecorations{
			NodeDecs: dst.NodeDecs{After: dst.NewLine},
		},
	})

	mName := c.mockName(typeSpec.Name.Name)
	for _, fn := range funcs {
		typePrefix := mName
		fieldSuffix := ""
		if _, ok := typeSpec.Type.(*dst.InterfaceType); ok {
			typePrefix = fmt.Sprintf("%s_%s", mName, fn.Name)
			fieldSuffix = "_" + fn.Name
		}
		elems = c.newMockElement(elems, typePrefix, fieldSuffix)
	}

	return elems
}

func (c *Converter) newMockElement(
	elems []dst.Expr, typePrefix, fieldSuffix string,
) []dst.Expr {
	pName := fmt.Sprintf("%s_%s", typePrefix, paramsIdent)
	elems = append(elems, &dst.KeyValueExpr{
		Key: dst.NewIdent(c.exportName(resultsByParamsIdent + fieldSuffix)),
		Value: &dst.CompositeLit{Type: &dst.MapType{
			Key: dst.NewIdent(pName),
			Value: &dst.StarExpr{
				X: dst.NewIdent(fmt.Sprintf(
					"%s_%s", typePrefix, resultMgrSuffix)),
			},
		}},
		Decs: dst.KeyValueExprDecorations{
			NodeDecs: dst.NodeDecs{After: dst.NewLine},
		},
	})
	elems = append(elems, &dst.KeyValueExpr{
		Key: dst.NewIdent(c.exportName(paramsIdent + fieldSuffix)),
		Value: &dst.CallExpr{
			Fun: dst.NewIdent("make"),
			Args: []dst.Expr{
				&dst.ChanType{
					Dir:   dst.SEND | dst.RECV,
					Value: dst.NewIdent(pName),
				},
				&dst.BasicLit{Kind: token.INT, Value: "100"},
			},
		},
		Decs: dst.KeyValueExprDecorations{
			NodeDecs: dst.NodeDecs{After: dst.NewLine},
		},
	})
	return elems
}

func (c *Converter) mockFunc(typePrefix, fieldSuffix string, fn Func) *dst.BlockStmt {
	stateSelector := &dst.SelectorExpr{
		X:   dst.NewIdent(mockReceiverIdent),
		Sel: dst.NewIdent(c.exportName(mockIdent)),
	}
	stmts := []dst.Stmt{
		&dst.AssignStmt{
			Lhs: []dst.Expr{dst.NewIdent(paramsIdent)},
			Tok: token.DEFINE,
			Rhs: []dst.Expr{&dst.CompositeLit{
				Type: dst.NewIdent(fmt.Sprintf(
					"%s_%s", typePrefix, paramsIdent)),
				Elts: c.passthroughElements(fn.Params, paramsIdent),
			}},
		},
		&dst.SendStmt{
			Chan: &dst.SelectorExpr{
				X:   dst.Clone(stateSelector).(dst.Expr),
				Sel: dst.NewIdent(c.exportName(paramsIdent + fieldSuffix)),
			},
			Value: dst.NewIdent(paramsIdent),
		},
	}

	stmts = append(stmts, &dst.AssignStmt{
		Lhs: []dst.Expr{
			dst.NewIdent(resultsIdent),
			dst.NewIdent(okIdent),
		},
		Tok: token.DEFINE,
		Rhs: []dst.Expr{&dst.IndexExpr{
			X: &dst.SelectorExpr{
				X: dst.Clone(stateSelector).(dst.Expr),
				Sel: dst.NewIdent(c.exportName(
					resultsByParamsIdent + fieldSuffix)),
			},
			Index: dst.NewIdent(paramsIdent),
		}},
	})
	stmts = append(stmts, &dst.IfStmt{
		Cond: dst.NewIdent(okIdent),
		Body: c.findResultsBody(fn),
	})
	stmts = append(stmts, mockFnReturn(fn.Results))

	return &dst.BlockStmt{List: stmts}
}

func (c *Converter) findResultsBody(fn Func) *dst.BlockStmt {
	stmts := []dst.Stmt{
		&dst.AssignStmt{
			Lhs: []dst.Expr{dst.NewIdent(iIdent)},
			Tok: token.DEFINE,
			Rhs: []dst.Expr{&dst.BinaryExpr{
				X: &dst.CallExpr{
					Fun: dst.NewIdent("int"),
					Args: []dst.Expr{&dst.CallExpr{
						Fun: &dst.Ident{
							Name: "AddUint32",
							Path: "sync/atomic",
						},
						Args: []dst.Expr{
							&dst.UnaryExpr{
								Op: token.AND,
								X: &dst.SelectorExpr{
									X:   dst.NewIdent(resultsIdent),
									Sel: dst.NewIdent(c.exportName(indexIdent)),
								},
							},
							&dst.BasicLit{Kind: token.INT, Value: "1"},
						},
					}},
				},
				Op: token.SUB,
				Y:  &dst.BasicLit{Kind: token.INT, Value: "1"},
			}},
		},
		&dst.IfStmt{
			Cond: &dst.BinaryExpr{
				X:  dst.NewIdent(iIdent),
				Op: token.GEQ,
				Y: &dst.CallExpr{
					Fun: dst.NewIdent("len"),
					Args: []dst.Expr{&dst.SelectorExpr{
						X:   dst.NewIdent(resultsIdent),
						Sel: dst.NewIdent(c.exportName(resultsIdent)),
					}},
				},
			},
			Body: &dst.BlockStmt{List: []dst.Stmt{
				&dst.IfStmt{
					Cond: &dst.UnaryExpr{
						Op: token.NOT,
						X: &dst.SelectorExpr{
							X:   dst.NewIdent(resultsIdent),
							Sel: dst.NewIdent(c.exportName(anyTimesIdent)),
						},
					},
					Body: &dst.BlockStmt{List: []dst.Stmt{
						&dst.ExprStmt{X: &dst.CallExpr{
							Fun: &dst.SelectorExpr{
								X: &dst.SelectorExpr{
									X: &dst.SelectorExpr{
										X: dst.NewIdent(mockReceiverIdent),
										Sel: dst.NewIdent(
											c.exportName(mockIdent)),
									},
									Sel: dst.NewIdent(
										c.exportName(testingIdent)),
								},
								Sel: dst.NewIdent("Fatalf"),
							},
							Args: []dst.Expr{
								&dst.BasicLit{
									Kind: token.STRING,
									Value: "\"Too many calls to" +
										" mock with parameters %#v\"",
								},
								dst.NewIdent(paramsIdent),
							},
						}},
						&dst.ReturnStmt{},
					}},
				},
				&dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent(iIdent)},
					Tok: token.ASSIGN,
					Rhs: []dst.Expr{&dst.BinaryExpr{
						X: &dst.CallExpr{
							Fun: dst.NewIdent("len"),
							Args: []dst.Expr{&dst.SelectorExpr{
								X:   dst.NewIdent(resultsIdent),
								Sel: dst.NewIdent(c.exportName(resultsIdent)),
							}},
						},
						Op: token.SUB,
						Y:  &dst.BasicLit{Kind: token.INT, Value: "1"},
					}},
				},
			}},
		},
	}

	if fn.Results != nil {
		stmts = append(stmts, &dst.AssignStmt{
			Lhs: []dst.Expr{dst.NewIdent(resultIdent)},
			Tok: token.DEFINE,
			Rhs: []dst.Expr{&dst.IndexExpr{
				X: &dst.SelectorExpr{
					X:   dst.NewIdent(resultsIdent),
					Sel: dst.NewIdent(c.exportName(resultsIdent)),
				},
				Index: dst.NewIdent(iIdent),
			}},
		})
		stmts = append(stmts, c.assignResult(fn.Results)...)
	}

	return &dst.BlockStmt{List: stmts}
}

func (c *Converter) recorderFn(typeName string, fn Func) *dst.FuncDecl {
	mName := c.mockName(typeName)

	recvType := fmt.Sprintf("%s_%s", mName, recorderIdent)
	fnName := fn.Name
	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	typePrefix := fmt.Sprintf("%s_%s", mName, fn.Name)
	var mockVal dst.Expr = &dst.SelectorExpr{
		X:   dst.NewIdent(mockReceiverIdent),
		Sel: dst.NewIdent(c.exportName(mockIdent)),
	}
	if fn.Name == "" {
		recvType = mName
		fnName = c.exportName(onCallFnName)
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
		typePrefix = mName
		mockVal = dst.NewIdent(mockReceiverIdent)
	}

	return &dst.FuncDecl{
		Recv: &dst.FieldList{List: []*dst.Field{{
			Names: []*dst.Ident{dst.NewIdent(mockReceiverIdent)},
			Type:  &dst.StarExpr{X: dst.NewIdent(recvType)},
		}}},
		Name: dst.NewIdent(fnName),
		Type: &dst.FuncType{
			Params: cloneAndNameUnnamed(paramPrefix, fn.Params),
			Results: &dst.FieldList{List: []*dst.Field{{
				Type: &dst.StarExpr{X: dst.NewIdent(fnRecName)},
			}}},
		},
		Body: c.recorderFnInterfaceBody(fnRecName, typePrefix, mockVal, fn),
		Decs: stdFuncDec(),
	}
}

func (c *Converter) recorderFnInterfaceBody(
	fnRecName, typePrefix string, mockValue dst.Expr, fn Func,
) *dst.BlockStmt {
	return &dst.BlockStmt{List: []dst.Stmt{&dst.ReturnStmt{
		Results: []dst.Expr{&dst.UnaryExpr{
			Op: token.AND,
			X: &dst.CompositeLit{
				Type: dst.NewIdent(fnRecName),
				Elts: []dst.Expr{
					&dst.KeyValueExpr{
						Key: dst.NewIdent(c.exportName(paramsIdent)),
						Value: &dst.CompositeLit{
							Type: dst.NewIdent(fmt.Sprintf(
								"%s_%s", typePrefix, paramsIdent)),
							Elts: c.passthroughElements(
								fn.Params, paramsIdent),
						},
						Decs: dst.KeyValueExprDecorations{
							NodeDecs: dst.NodeDecs{After: dst.NewLine},
						},
					},
					&dst.KeyValueExpr{
						Key:   dst.NewIdent(c.exportName(mockIdent)),
						Value: mockValue,
						Decs: dst.KeyValueExprDecorations{
							NodeDecs: dst.NodeDecs{After: dst.NewLine},
						},
					},
				},
				Decs: dst.CompositeLitDecorations{Lbrace: []string{"\n"}},
			},
		}},
	}}}
}

func (c *Converter) recorderReturnFn(typeName string, fn Func) *dst.FuncDecl {
	mName := c.mockName(typeName)

	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	results := fmt.Sprintf("%s_%s_%s", mName, fn.Name, resultsIdent)
	resultMgr := fmt.Sprintf("%s_%s_%s", mName, fn.Name, resultMgrSuffix)
	resultsByParams := fmt.Sprintf("%s_%s", resultsByParamsIdent, fn.Name)
	mockSel := &dst.SelectorExpr{X: &dst.SelectorExpr{
		X:   dst.NewIdent(recorderReceiverIdent),
		Sel: dst.NewIdent(c.exportName(mockIdent)),
	}}
	if fn.Name == "" {
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
		results = fmt.Sprintf("%s_%s", mName, resultsIdent)
		resultMgr = fmt.Sprintf("%s_%s", mName, resultMgrSuffix)
		resultsByParams = resultsByParamsIdent
		mockSel = &dst.SelectorExpr{X: &dst.SelectorExpr{
			X:   dst.NewIdent(recorderReceiverIdent),
			Sel: dst.NewIdent(c.exportName(mockIdent)),
		}}
	}

	return &dst.FuncDecl{
		Recv: &dst.FieldList{List: []*dst.Field{{
			Names: []*dst.Ident{dst.NewIdent(recorderReceiverIdent)},
			Type:  &dst.StarExpr{X: dst.NewIdent(fnRecName)},
		}}},
		Name: dst.NewIdent(c.exportName(returnFnName)),
		Type: &dst.FuncType{
			Params: cloneAndNameUnnamed(resultPrefix, fn.Results),
			Results: &dst.FieldList{List: []*dst.Field{{
				Type: &dst.StarExpr{X: dst.NewIdent(fnRecName)},
			}}},
		},
		Body: &dst.BlockStmt{List: []dst.Stmt{
			&dst.IfStmt{
				Cond: &dst.BinaryExpr{
					X: &dst.SelectorExpr{
						X:   dst.NewIdent(recorderReceiverIdent),
						Sel: dst.NewIdent(c.exportName(resultsIdent)),
					},
					Op: token.EQL,
					Y:  dst.NewIdent("nil"),
				},
				Body: &dst.BlockStmt{List: []dst.Stmt{
					&dst.IfStmt{
						Init: &dst.AssignStmt{
							Lhs: []dst.Expr{
								dst.NewIdent("_"),
								dst.NewIdent(okIdent),
							},
							Tok: token.DEFINE,
							Rhs: []dst.Expr{&dst.IndexExpr{
								X: cloneSelect(
									mockSel, c.exportName(resultsByParams)),
								Index: &dst.SelectorExpr{
									X: dst.NewIdent(recorderReceiverIdent),
									Sel: dst.NewIdent(
										c.exportName(paramsIdent)),
								},
							}},
						},
						Cond: dst.NewIdent(okIdent),
						Body: &dst.BlockStmt{List: []dst.Stmt{
							&dst.ExprStmt{X: &dst.CallExpr{
								Fun: &dst.SelectorExpr{
									X: cloneSelect(
										mockSel, c.exportName(testingIdent)),
									Sel: dst.NewIdent("Fatalf"),
								},
								Args: []dst.Expr{
									&dst.BasicLit{
										Kind: token.STRING,
										Value: "\"Expectations already" +
											" recorded for mock with" +
											" parameters %#v\"",
									},
									&dst.SelectorExpr{
										X: dst.NewIdent(
											recorderReceiverIdent),
										Sel: dst.NewIdent(c.exportName(
											paramsIdent)),
									},
								},
							}},
							&dst.ReturnStmt{Results: []dst.Expr{
								dst.NewIdent("nil")},
							},
						}},
						Decs: dst.IfStmtDecorations{NodeDecs: dst.NodeDecs{
							After: dst.EmptyLine,
						}},
					},
					&dst.AssignStmt{
						Lhs: []dst.Expr{&dst.SelectorExpr{
							X:   dst.NewIdent(recorderReceiverIdent),
							Sel: dst.NewIdent(c.exportName(resultsIdent)),
						}},
						Tok: token.ASSIGN,
						Rhs: []dst.Expr{&dst.UnaryExpr{
							Op: token.AND,
							X: &dst.CompositeLit{
								Type: dst.NewIdent(c.exportName(resultMgr)),
								Elts: []dst.Expr{
									&dst.KeyValueExpr{
										Key: dst.NewIdent(c.exportName(
											resultsIdent)),
										Value: &dst.CompositeLit{
											Type: &dst.ArrayType{
												Elt: &dst.StarExpr{
													X: dst.NewIdent(
														c.exportName(results)),
												},
											},
										},
									},
									&dst.KeyValueExpr{
										Key: dst.NewIdent(c.exportName(
											indexIdent)),
										Value: &dst.BasicLit{
											Kind:  token.INT,
											Value: "0",
										},
									},
									&dst.KeyValueExpr{
										Key: dst.NewIdent(c.exportName(
											anyTimesIdent)),
										Value: dst.NewIdent("false"),
									},
								},
							},
						}},
					},
					&dst.AssignStmt{
						Lhs: []dst.Expr{&dst.IndexExpr{
							X: &dst.SelectorExpr{
								X: &dst.SelectorExpr{
									X:   dst.NewIdent(recorderReceiverIdent),
									Sel: dst.NewIdent(c.exportName(mockIdent)),
								},
								Sel: dst.NewIdent(c.exportName(resultsByParams)),
							},
							Index: &dst.SelectorExpr{
								X:   dst.NewIdent(recorderReceiverIdent),
								Sel: dst.NewIdent(c.exportName(paramsIdent)),
							},
						}},
						Tok: token.ASSIGN,
						Rhs: []dst.Expr{&dst.SelectorExpr{
							X:   dst.NewIdent(recorderReceiverIdent),
							Sel: dst.NewIdent(c.exportName(resultsIdent)),
						}},
					},
				}},
			},
			&dst.AssignStmt{
				Lhs: []dst.Expr{
					&dst.SelectorExpr{
						X: &dst.SelectorExpr{
							X:   dst.NewIdent(recorderReceiverIdent),
							Sel: dst.NewIdent(c.exportName(resultsIdent)),
						},
						Sel: dst.NewIdent(c.exportName(resultsIdent)),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []dst.Expr{&dst.CallExpr{
					Fun: dst.NewIdent("append"),
					Args: []dst.Expr{
						&dst.SelectorExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent(recorderReceiverIdent),
								Sel: dst.NewIdent(c.exportName(resultsIdent)),
							},
							Sel: dst.NewIdent(c.exportName(resultsIdent)),
						},
						&dst.UnaryExpr{
							Op: token.AND,
							X: &dst.CompositeLit{
								Type: dst.NewIdent(c.exportName(results)),
								Elts: c.passthroughElements(
									fn.Results, resultsIdent),
							},
						},
					},
				}},
			},
			&dst.ReturnStmt{Results: []dst.Expr{
				dst.NewIdent(recorderReceiverIdent)},
			},
		}},
		Decs: stdFuncDec(),
	}
}

func (c *Converter) recorderTimesFn(typeName string, fn Func) *dst.FuncDecl {
	mName := c.mockName(typeName)

	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
	}

	return &dst.FuncDecl{
		Recv: &dst.FieldList{List: []*dst.Field{{
			Names: []*dst.Ident{dst.NewIdent(recorderReceiverIdent)},
			Type:  &dst.StarExpr{X: dst.NewIdent(fnRecName)},
		}}},
		Name: dst.NewIdent(c.exportName(timesFnName)),
		Type: &dst.FuncType{
			Params: &dst.FieldList{List: []*dst.Field{{
				Names: []*dst.Ident{dst.NewIdent(countIdent)},
				Type:  dst.NewIdent("int"),
			}}},
			Results: &dst.FieldList{List: []*dst.Field{{
				Type: &dst.StarExpr{X: dst.NewIdent(fnRecName)},
			}}},
		},
		Body: &dst.BlockStmt{List: []dst.Stmt{
			&dst.IfStmt{
				Cond: &dst.BinaryExpr{
					X: &dst.SelectorExpr{
						X:   dst.NewIdent(recorderReceiverIdent),
						Sel: dst.NewIdent(c.exportName(resultsIdent)),
					},
					Op: token.EQL,
					Y:  dst.NewIdent("nil"),
				},
				Body: &dst.BlockStmt{List: []dst.Stmt{
					&dst.ExprStmt{X: &dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X: &dst.SelectorExpr{
								X: &dst.SelectorExpr{
									X:   dst.NewIdent(recorderReceiverIdent),
									Sel: dst.NewIdent(c.exportName(mockIdent)),
								},
								Sel: dst.NewIdent(c.exportName(testingIdent)),
							},
							Sel: dst.NewIdent("Fatalf"),
						},
						Args: []dst.Expr{
							&dst.BasicLit{
								Kind: token.STRING,
								Value: "\"Return must be called" +
									" before calling Times\"",
							},
						},
					}},
					&dst.ReturnStmt{Results: []dst.Expr{dst.NewIdent("nil")}},
				}},
			},
			&dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent(lastIdent)},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{&dst.IndexExpr{
					X: &dst.SelectorExpr{
						X: &dst.SelectorExpr{
							X:   dst.NewIdent(recorderReceiverIdent),
							Sel: dst.NewIdent(c.exportName(resultsIdent)),
						},
						Sel: dst.NewIdent(c.exportName(resultsIdent)),
					},
					Index: &dst.BinaryExpr{
						X: &dst.CallExpr{
							Fun: dst.NewIdent("len"),
							Args: []dst.Expr{
								&dst.SelectorExpr{
									X: &dst.SelectorExpr{
										X: dst.NewIdent(recorderReceiverIdent),
										Sel: dst.NewIdent(c.exportName(
											resultsIdent)),
									},
									Sel: dst.NewIdent(c.exportName(
										resultsIdent)),
								},
							},
						},
						Op: token.SUB,
						Y:  &dst.BasicLit{Kind: token.INT, Value: "1"},
					},
				}},
			},
			&dst.ForStmt{
				Init: &dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent("n")},
					Tok: token.DEFINE,
					Rhs: []dst.Expr{&dst.BasicLit{Kind: token.INT, Value: "0"}},
				},
				Cond: &dst.BinaryExpr{
					X:  dst.NewIdent("n"),
					Op: token.LSS,
					Y: &dst.BinaryExpr{
						X:  dst.NewIdent(countIdent),
						Op: token.SUB,
						Y:  &dst.BasicLit{Kind: token.INT, Value: "1"},
					},
				},
				Post: &dst.IncDecStmt{X: dst.NewIdent("n"), Tok: token.INC},
				Body: &dst.BlockStmt{List: []dst.Stmt{&dst.AssignStmt{
					Lhs: []dst.Expr{&dst.SelectorExpr{
						X: &dst.SelectorExpr{
							X:   dst.NewIdent(recorderReceiverIdent),
							Sel: dst.NewIdent(c.exportName(resultsIdent)),
						},
						Sel: dst.NewIdent(c.exportName(resultsIdent)),
					}},
					Tok: token.ASSIGN,
					Rhs: []dst.Expr{&dst.CallExpr{
						Fun: dst.NewIdent("append"),
						Args: []dst.Expr{
							&dst.SelectorExpr{
								X: &dst.SelectorExpr{
									X: dst.NewIdent(recorderReceiverIdent),
									Sel: dst.NewIdent(c.exportName(
										resultsIdent)),
								},
								Sel: dst.NewIdent(c.exportName(resultsIdent)),
							},
							dst.NewIdent(lastIdent),
						},
					}},
				}}},
			},
			&dst.ReturnStmt{Results: []dst.Expr{
				dst.NewIdent(recorderReceiverIdent)},
			},
		}},
		Decs: stdFuncDec(),
	}
}

func (c *Converter) recorderAnyTimesFn(
	typeName string, fn Func,
) *dst.FuncDecl {
	mName := c.mockName(typeName)

	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
	}

	return &dst.FuncDecl{
		Recv: &dst.FieldList{List: []*dst.Field{{
			Names: []*dst.Ident{dst.NewIdent(recorderReceiverIdent)},
			Type:  &dst.StarExpr{X: dst.NewIdent(fnRecName)},
		}}},
		Name: dst.NewIdent(c.exportName(anyTimesFnName)),
		Type: &dst.FuncType{Params: &dst.FieldList{}},
		Body: &dst.BlockStmt{List: []dst.Stmt{
			&dst.IfStmt{
				Cond: &dst.BinaryExpr{
					X: &dst.SelectorExpr{
						X:   dst.NewIdent(recorderReceiverIdent),
						Sel: dst.NewIdent(c.exportName(resultsIdent)),
					},
					Op: token.EQL,
					Y:  dst.NewIdent("nil"),
				},
				Body: &dst.BlockStmt{List: []dst.Stmt{
					&dst.ExprStmt{X: &dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X: &dst.SelectorExpr{
								X: &dst.SelectorExpr{
									X:   dst.NewIdent(recorderReceiverIdent),
									Sel: dst.NewIdent(c.exportName(mockIdent)),
								},
								Sel: dst.NewIdent(c.exportName(testingIdent)),
							},
							Sel: dst.NewIdent("Fatalf"),
						},
						Args: []dst.Expr{&dst.BasicLit{
							Kind: token.STRING,
							Value: "\"Return must be called before" +
								" calling AnyTimes\"",
						}},
					}},
					&dst.ReturnStmt{},
				}},
			},
			&dst.AssignStmt{
				Lhs: []dst.Expr{&dst.SelectorExpr{
					X: &dst.SelectorExpr{
						X:   dst.NewIdent(recorderReceiverIdent),
						Sel: dst.NewIdent(c.exportName(resultsIdent)),
					},
					Sel: dst.NewIdent(c.exportName(anyTimesIdent)),
				}},
				Tok: token.ASSIGN,
				Rhs: []dst.Expr{dst.NewIdent("true")},
			},
		}},
		Decs: stdFuncDec(),
	}
}

func (c *Converter) passthroughElements(
	fl *dst.FieldList, label string,
) []dst.Expr {
	unnamedPrefix, comparable := labelDirection(label)
	var elts []dst.Expr
	if fl != nil {
		beforeDec := dst.NewLine
		fields := fl.List
		for n, field := range fields {
			if len(field.Names) == 0 {
				pName := fmt.Sprintf("%s%d", unnamedPrefix, n+1)
				elts = append(elts, &dst.KeyValueExpr{
					Key:   dst.NewIdent(c.exportName(pName)),
					Value: passthroughValue(pName, field, comparable),
					Decs: dst.KeyValueExprDecorations{NodeDecs: dst.NodeDecs{
						Before: beforeDec,
						After:  dst.NewLine,
					}},
				})
				beforeDec = dst.None
			}

			for _, name := range field.Names {
				elts = append(elts, &dst.KeyValueExpr{
					Key:   dst.NewIdent(c.exportName(name.Name)),
					Value: passthroughValue(name.Name, field, comparable),
					Decs: dst.KeyValueExprDecorations{NodeDecs: dst.NodeDecs{
						Before: beforeDec,
						After:  dst.NewLine,
					}},
				})
				beforeDec = dst.None
			}
		}
	}

	return elts
}

func passthroughValue(
	name string, field *dst.Field, comparable bool,
) dst.Expr {
	var val dst.Expr
	val = dst.NewIdent(name)
	if comparable && !isComparable(field.Type) {
		val = &dst.CallExpr{
			Fun:  &dst.Ident{Path: hashPkg, Name: "DeepHash"},
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

func (c *Converter) assignResult(resFL *dst.FieldList) []dst.Stmt {
	var assigns []dst.Stmt
	if resFL != nil {
		results := resFL.List
		for n, result := range results {
			if len(result.Names) == 0 {
				rName := fmt.Sprintf("%s%d", resultPrefix, n+1)
				assigns = append(assigns, &dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent(rName)},
					Tok: token.ASSIGN,
					Rhs: []dst.Expr{&dst.SelectorExpr{
						X:   dst.NewIdent(resultIdent),
						Sel: dst.NewIdent(c.exportName(rName)),
					}},
				})
			}

			for _, name := range result.Names {
				assigns = append(assigns, &dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent(name.Name)},
					Tok: token.ASSIGN,
					Rhs: []dst.Expr{&dst.SelectorExpr{
						X:   dst.NewIdent(resultIdent),
						Sel: dst.NewIdent(c.exportName(name.Name)),
					}},
				})
			}
		}
	}
	return assigns
}

func cloneAndNameUnnamed(
	prefix string, fieldList *dst.FieldList,
) *dst.FieldList {
	fieldList = cloneNilableFieldList(fieldList)
	if fieldList != nil {
		for n, f := range fieldList.List {
			if len(f.Names) == 0 {
				f.Names = []*dst.Ident{dst.NewIdent(fmt.Sprintf(
					"%s%d", prefix, n+1))}
			}
		}
	}
	return fieldList
}

func mockFnReturn(resFL *dst.FieldList) dst.Stmt {
	var exprs []dst.Expr
	if resFL != nil {
		results := resFL.List
		for n, result := range results {
			if len(result.Names) == 0 {
				rName := fmt.Sprintf("%s%d", resultPrefix, n+1)
				exprs = append(exprs, dst.NewIdent(rName))
			}
			for _, name := range result.Names {
				exprs = append(exprs, dst.NewIdent(name.Name))
			}
		}
	}
	return &dst.ReturnStmt{Results: exprs}
}

func (c *Converter) mockName(typeName string) string {
	return c.exportName(mockIdent + strings.Title(typeName))
}

func (c *Converter) exportName(name string) string {
	if c.export {
		name = strings.Title(name)
	}
	return name
}

func stdFuncDec() dst.FuncDeclDecorations {
	return dst.FuncDeclDecorations{NodeDecs: dst.NodeDecs{
		Before: dst.EmptyLine,
		After:  dst.EmptyLine,
	}}
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

func cloneNilableFieldList(fl *dst.FieldList) *dst.FieldList {
	if fl != nil {
		fl = dst.Clone(fl).(*dst.FieldList)
	}
	return fl
}

func cloneSelect(x *dst.SelectorExpr, sel string) *dst.SelectorExpr {
	x = dst.Clone(x).(*dst.SelectorExpr)
	x.Sel = dst.NewIdent(sel)
	return x
}
