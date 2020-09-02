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
	mock.Decs.Start.Append(fmt.Sprintf("// %s holds the state of a mock of the %s type", mName, typeSpec.Name.Name))

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
					List: []*dst.Field{
						{
							Names: []*dst.Ident{dst.NewIdent(c.exportName(mockIdent))},
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
func (c *Converter) MethodStructs(typeSpec *dst.TypeSpec, fn Func) []dst.Decl {
	prefix := c.mockName(typeSpec.Name.Name)
	if _, ok := typeSpec.Type.(*dst.InterfaceType); ok {
		prefix = fmt.Sprintf("%s_%s", prefix, fn.Name)
	}

	decls := make([]dst.Decl, 3)
	decls[0] = c.methodStruct(typeSpec.Name.Name, prefix, paramsIdent, fn.Params)
	decls[1] = c.methodStruct(typeSpec.Name.Name, prefix, resultsIdent, fn.Results)
	decls[2] = c.fnRecorderStruct(typeSpec.Name.Name, prefix)
	return decls
}

// NewFunc generates a function for constructing a mock
func (c *Converter) NewFunc(typeSpec *dst.TypeSpec, funcs []Func) (funcDecl *dst.FuncDecl) {
	fnName := c.exportName("newMock" + typeSpec.Name.Name)
	mName := c.mockName(typeSpec.Name.Name)
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
								Elts: c.newMockElements(typeSpec, funcs),
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
	fnName = c.exportName(fnName)
	mName := c.mockName(typeName)
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
										Key:   dst.NewIdent(c.exportName(mockIdent)),
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
	mName := c.mockName(typeName)
	mockFn := &dst.FuncDecl{
		Recv: &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent(mockReceiverIdent)},
					Type:  &dst.StarExpr{X: dst.NewIdent(mName)},
				},
			},
		},
		Name: dst.NewIdent(c.exportName(mockFnName)),
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
															Key:   dst.NewIdent(c.exportName(mockIdent)),
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
													Sel: dst.NewIdent(c.exportName(fnFnName)),
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
		c.exportName(mockFnName), mockIdent, typeName))

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
		Body: c.mockFunc(typePrefix, fieldSuffix, fn),
		Decs: stdFuncDec(),
	}
}

// RecorderMethods generates a recorder implementation of a method and associated
// return method
func (c *Converter) RecorderMethods(typeName string, fn Func) (funcDecls []dst.Decl) {
	return []dst.Decl{
		c.recorderFn(typeName, fn),
		c.recorderReturnFn(typeName, fn),
	}
}

func (c *Converter) baseMockFieldList(typeSpec *dst.TypeSpec, funcs []Func) []*dst.Field {
	var fields []*dst.Field

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

func (c *Converter) baseMockFieldsPerFn(fields []*dst.Field, typePrefix, fieldSuffix string) []*dst.Field {
	pName := c.exportName(fmt.Sprintf("%s_%s", typePrefix, paramsIdent))
	fields = append(fields, &dst.Field{
		Names: []*dst.Ident{
			dst.NewIdent(c.exportName(resultsByParamsIdent + fieldSuffix)),
		},
		Type: &dst.MapType{
			Key:   &dst.Ident{Name: pName},
			Value: &dst.Ident{Name: fmt.Sprintf("%s_%s", typePrefix, resultsIdent)},
		},
	})

	fields = append(fields, &dst.Field{
		Names: []*dst.Ident{dst.NewIdent(c.exportName(paramsIdent + fieldSuffix))},
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

func (c *Converter) methodStruct(typeName, prefix, label string, fieldList *dst.FieldList) *dst.GenDecl {
	unnamedPrefix, comparable := labelDirection(label)
	fieldList = dst.Clone(fieldList).(*dst.FieldList)

	if fieldList.List != nil {
		for n, f := range fieldList.List {
			if len(f.Names) == 0 {
				f.Names = []*dst.Ident{dst.NewIdent(fmt.Sprintf("%s%d", unnamedPrefix, n+1))}
			}

			for nn := range f.Names {
				f.Names[nn] = dst.NewIdent(c.exportName(f.Names[nn].Name))
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

func (c *Converter) fnRecorderStruct(typeName string, prefix string) *dst.GenDecl {
	mName := c.mockName(typeName)
	structName := fmt.Sprintf("%s_%s", prefix, fnRecorderSuffix)
	sType := &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: dst.NewIdent(structName),
				Type: &dst.StructType{Fields: &dst.FieldList{List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent(c.exportName(paramsIdent))},
						Type:  dst.NewIdent(fmt.Sprintf("%s_%s", prefix, paramsIdent)),
					},
					{
						Names: []*dst.Ident{dst.NewIdent(c.exportName(mockIdent))},
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

func (c *Converter) newMockElements(typeSpec *dst.TypeSpec, funcs []Func) []dst.Expr {
	var elems []dst.Expr

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

func (c *Converter) newMockElement(elems []dst.Expr, typePrefix, fieldSuffix string) []dst.Expr {
	pName := fmt.Sprintf("%s_%s", typePrefix, paramsIdent)
	elems = append(elems, &dst.KeyValueExpr{
		Key: dst.NewIdent(c.exportName(resultsByParamsIdent + fieldSuffix)),
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
		Key: dst.NewIdent(c.exportName(paramsIdent + fieldSuffix)),
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

func (c *Converter) mockFunc(typePrefix, fieldSuffix string, fn Func) *dst.BlockStmt {
	stateSelector := &dst.SelectorExpr{
		X:   dst.NewIdent(mockReceiverIdent),
		Sel: dst.NewIdent(c.exportName(mockIdent)),
	}
	return &dst.BlockStmt{
		List: []dst.Stmt{
			&dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent(paramsIdent)},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{
					&dst.CompositeLit{
						Type: dst.NewIdent(fmt.Sprintf("%s_%s", typePrefix, paramsIdent)),
						Elts: c.passthroughElements(fn.Params.List, paramsIdent),
					},
				},
			},
			&dst.SendStmt{
				Chan: &dst.SelectorExpr{
					X:   dst.Clone(stateSelector).(dst.Expr),
					Sel: dst.NewIdent(c.exportName(paramsIdent + fieldSuffix)),
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
							Sel: dst.NewIdent(c.exportName(resultsByParamsIdent + fieldSuffix)),
						},
						Index: dst.NewIdent(paramsIdent),
					},
				},
			},
			&dst.IfStmt{
				Cond: dst.NewIdent(okIdent),
				Body: &dst.BlockStmt{List: c.assignResults(fn.Results.List)},
			},
			mockFnReturn(fn.Results.List),
		},
	}
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
		fnName = c.exportName(onCallIdent)
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
		Body: c.recorderFnInterfaceBody(fnRecName, typePrefix, mockVal, fn),
		Decs: stdFuncDec(),
	}
}

func (c *Converter) recorderFnInterfaceBody(fnRecName, typePrefix string, mockValue dst.Expr, fn Func) *dst.BlockStmt {
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
									Key: dst.NewIdent(c.exportName(paramsIdent)),
									Value: &dst.CompositeLit{
										Type: dst.NewIdent(fmt.Sprintf("%s_%s", typePrefix, paramsIdent)),
										Elts: c.passthroughElements(fn.Params.List, paramsIdent),
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
					},
				},
			},
		},
	}
}

func (c *Converter) recorderReturnFn(typeName string, fn Func) *dst.FuncDecl {
	mName := c.mockName(typeName)

	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	resultsByParams := c.exportName(fmt.Sprintf("%s_%s", resultsByParamsIdent, fn.Name))
	results := fmt.Sprintf("%s_%s_%s", mName, fn.Name, resultsIdent)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
		resultsByParams = c.exportName(resultsByParamsIdent)
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
		Name: dst.NewIdent(c.exportName(returnFnName)),
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
									Sel: dst.NewIdent(c.exportName(mockIdent)),
								},
								Sel: dst.NewIdent(c.exportName(resultsByParams)),
							},
							Index: &dst.SelectorExpr{
								X:   dst.NewIdent(recorderReceiverIdent),
								Sel: dst.NewIdent(c.exportName(paramsIdent)),
							},
						},
					},
					Tok: token.ASSIGN,
					Rhs: []dst.Expr{
						&dst.CompositeLit{
							Type: dst.NewIdent(results),
							Elts: c.passthroughElements(fn.Results.List, resultsIdent),
						},
					},
				},
			},
		},
		Decs: stdFuncDec(),
	}
}

func (c *Converter) passthroughElements(fields []*dst.Field, label string) []dst.Expr {
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
				Key:   dst.NewIdent(c.exportName(name.Name)),
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

func (c *Converter) assignResults(results []*dst.Field) []dst.Stmt {
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
						Sel: dst.NewIdent(c.exportName(name.Name)),
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
