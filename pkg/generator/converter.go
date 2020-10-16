package generator

import (
	"fmt"
	"go/token"
	"strings"

	. "github.com/dave/dst"

	"github.com/myshkin5/moqueries/pkg/logs"
)

const (
	moqueriesPkg = "github.com/myshkin5/moqueries"
	hashPkg      = moqueriesPkg + "/pkg/hash"
	moqPkg       = moqueriesPkg + "/pkg/moq"

	mockConfigType = "MockConfig"
	moqTType       = "MoqT"
	sceneType      = "Scene"

	anyTimesIdent         = "anyTimes"
	configIdent           = "config"
	countIdent            = "count"
	expectationIdent      = "Expectation"
	iIdent                = "i"
	indexIdent            = "index"
	lastIdent             = "last"
	missingIdent          = "missing"
	mockIdent             = "mock"
	mockReceiverIdent     = "m"
	okIdent               = "ok"
	paramsIdent           = "params"
	recorderIdent         = "recorder"
	recorderReceiverIdent = "r"
	resultsByParamsIdent  = "resultsByParams"
	resultIdent           = "result"
	resultsIdent          = "results"
	sceneIdent            = "scene"
	strictIdent           = "Strict"

	anyTimesFnName = "anyTimes"
	assertFnName   = "AssertExpectationsMet"
	errorfFnName   = "Errorf"
	fatalfFnName   = "Fatalf"
	fnFnName       = "fn"
	mockFnName     = "mock"
	onCallFnName   = "onCall"
	resetFnName    = "Reset"
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
	Params  *FieldList
	Results *FieldList
}

// BaseStruct generates the base structure used to store the mock's state
func (c *Converter) BaseStruct(typeSpec *TypeSpec, funcs []Func) *GenDecl {
	mName := c.mockName(typeSpec.Name.Name)
	mock := &GenDecl{
		Tok: token.TYPE,
		Specs: []Spec{&TypeSpec{
			Name: NewIdent(mName),
			Type: &StructType{Fields: &FieldList{
				List: c.baseMockFieldList(typeSpec, funcs),
			}},
		}},
	}
	mock.Decs.Before = NewLine
	mock.Decs.Start.Append(fmt.Sprintf(
		"// %s holds the state of a mock of the %s type",
		mName,
		typeSpec.Name.Name))

	return mock
}

// IsolationStruct generates a struct used to isolate an interface for the mock
func (c *Converter) IsolationStruct(typeName, suffix string) (structDecl *GenDecl) {
	mName := c.mockName(typeName)
	iName := fmt.Sprintf("%s_%s", mName, suffix)
	isolate := &GenDecl{
		Tok: token.TYPE,
		Specs: []Spec{&TypeSpec{
			Name: NewIdent(iName),
			Type: &StructType{Fields: &FieldList{
				List: []*Field{{
					Names: []*Ident{NewIdent(c.exportName(mockIdent))},
					Type:  &StarExpr{X: NewIdent(mName)},
				}},
			}},
		}},
	}
	isolate.Decs.Before = NewLine
	isolate.Decs.Start.Append(fmt.Sprintf(
		"// %s isolates the %s interface of the %s type",
		iName,
		suffix,
		typeName))

	return isolate
}

// MethodStructs generates a structure for storing a set of parameters or
// a set of results for a method invocation of a mock
func (c *Converter) MethodStructs(typeSpec *TypeSpec, fn Func) []Decl {
	prefix := c.mockName(typeSpec.Name.Name)
	if _, ok := typeSpec.Type.(*InterfaceType); ok {
		prefix = fmt.Sprintf("%s_%s", prefix, fn.Name)
	}

	var decls []Decl
	decls = append(decls,
		c.methodStruct(typeSpec.Name.Name, prefix, paramsIdent, fn.Params))
	decls = append(decls, c.resultMgrStruct(typeSpec.Name.Name, prefix))
	decls = append(decls,
		c.methodStruct(typeSpec.Name.Name, prefix, resultsIdent, fn.Results))
	decls = append(decls, c.fnRecorderStruct(typeSpec.Name.Name, prefix))
	return decls
}

// NewFunc generates a function for constructing a mock
func (c *Converter) NewFunc(typeSpec *TypeSpec, funcs []Func) (funcDecl *FuncDecl) {
	fnName := c.exportName("newMock" + typeSpec.Name.Name)
	mName := c.mockName(typeSpec.Name.Name)
	mockFn := &FuncDecl{
		Name: NewIdent(fnName),
		Type: &FuncType{
			Params: &FieldList{List: []*Field{
				{
					Names: []*Ident{NewIdent(sceneIdent)},
					Type:  &StarExpr{X: &Ident{Name: sceneType, Path: moqPkg}},
				},
				{
					Names: []*Ident{NewIdent(configIdent)},
					Type:  &StarExpr{X: &Ident{Name: mockConfigType, Path: moqPkg}},
				},
			}},
			Results: &FieldList{List: []*Field{{Type: &StarExpr{X: NewIdent(mName)}}}},
		},
		Body: &BlockStmt{List: []Stmt{
			&IfStmt{
				Cond: &BinaryExpr{
					X:  NewIdent(configIdent),
					Op: token.EQL,
					Y:  NewIdent("nil"),
				},
				Body: &BlockStmt{List: []Stmt{&AssignStmt{
					Lhs: []Expr{NewIdent(configIdent)},
					Tok: token.ASSIGN,
					Rhs: []Expr{&UnaryExpr{
						Op: token.AND,
						X: &CompositeLit{Type: &Ident{
							Name: mockConfigType,
							Path: moqPkg,
						}},
					}},
				}}},
			},
			&AssignStmt{
				Lhs: []Expr{NewIdent(mockReceiverIdent)},
				Tok: token.DEFINE,
				Rhs: []Expr{&UnaryExpr{
					Op: token.AND,
					X: &CompositeLit{
						Type: NewIdent(mName),
						Elts: c.newMockElements(typeSpec, funcs),
						Decs: CompositeLitDecorations{Lbrace: []string{"\n"}},
					},
				}},
			},
			&ExprStmt{X: &CallExpr{Fun: &SelectorExpr{
				X:   NewIdent(mockReceiverIdent),
				Sel: NewIdent(resetFnName),
			}}},
			&ExprStmt{X: &CallExpr{
				Fun: &SelectorExpr{
					X:   NewIdent(sceneIdent),
					Sel: NewIdent("AddMock"),
				},
				Args: []Expr{NewIdent(mockReceiverIdent)},
			}},
			&ReturnStmt{Results: []Expr{NewIdent(mockReceiverIdent)}},
		}},
	}

	mockFn.Decs.Before = NewLine
	mockFn.Decs.Start.Append(fmt.Sprintf(
		"// %s creates a new mock of the %s type", fnName, typeSpec.Name.Name))

	return mockFn
}

// IsolationAccessor generates a function to access an isolation interface
func (c *Converter) IsolationAccessor(typeName, suffix, fnName string) (funcDecl *FuncDecl) {
	fnName = c.exportName(fnName)
	mName := c.mockName(typeName)
	iName := fmt.Sprintf("%s_%s", mName, suffix)
	mockFn := &FuncDecl{
		Recv: &FieldList{List: []*Field{{
			Names: []*Ident{NewIdent(mockReceiverIdent)},
			Type:  &StarExpr{X: NewIdent(mName)},
		}}},
		Name: NewIdent(fnName),
		Type: &FuncType{
			Params:  &FieldList{},
			Results: &FieldList{List: []*Field{{Type: &StarExpr{X: NewIdent(iName)}}}},
		},
		Body: &BlockStmt{List: []Stmt{&ReturnStmt{
			Results: []Expr{&UnaryExpr{
				Op: token.AND,
				X: &CompositeLit{
					Type: NewIdent(iName),
					Elts: []Expr{&KeyValueExpr{
						Key:   NewIdent(c.exportName(mockIdent)),
						Value: NewIdent(mockReceiverIdent),
						Decs: KeyValueExprDecorations{
							NodeDecs: NodeDecs{After: NewLine},
						},
					}},
					Decs: CompositeLitDecorations{Lbrace: []string{"\n"}},
				},
			}},
		}}},
	}

	mockFn.Decs.Before = NewLine
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
	funcDecl *FuncDecl,
) {
	mName := c.mockName(typeName)
	ellipsis := false
	if len(fn.Params.List) > 0 {
		if _, ok := fn.Params.List[len(fn.Params.List)-1].Type.(*Ellipsis); ok {
			ellipsis = true
		}
	}
	fnLitCall := &CallExpr{
		Fun: &SelectorExpr{
			X:   NewIdent(mockIdent),
			Sel: NewIdent(c.exportName(fnFnName)),
		},
		Args:     passthroughFields(fn.Params),
		Ellipsis: ellipsis,
	}
	var fnLitRetStmt Stmt
	fnLitRetStmt = &ReturnStmt{Results: []Expr{fnLitCall}}
	if fn.Results == nil {
		fnLitRetStmt = &ExprStmt{X: fnLitCall}
	}

	mockFn := &FuncDecl{
		Recv: &FieldList{List: []*Field{{
			Names: []*Ident{NewIdent(mockReceiverIdent)},
			Type:  &StarExpr{X: NewIdent(mName)},
		}}},
		Name: NewIdent(c.exportName(mockFnName)),
		Type: &FuncType{
			Params: &FieldList{},
			Results: &FieldList{List: []*Field{{
				Type: &Ident{Path: pkgPath, Name: typeName},
			}}},
		},
		Body: &BlockStmt{List: []Stmt{&ReturnStmt{
			Results: []Expr{&FuncLit{
				Type: &FuncType{
					Params:  Clone(fn.Params).(*FieldList),
					Results: cloneNilableFieldList(fn.Results),
				},
				Body: &BlockStmt{List: []Stmt{&AssignStmt{
					Lhs: []Expr{NewIdent(mockIdent)},
					Tok: token.DEFINE,
					Rhs: []Expr{&UnaryExpr{
						Op: token.AND,
						X: &CompositeLit{
							Type: NewIdent(fmt.Sprintf("%s_%s", mName, mockIdent)),
							Elts: []Expr{&KeyValueExpr{
								Key:   NewIdent(c.exportName(mockIdent)),
								Value: NewIdent(mockReceiverIdent),
							}},
						},
					}}},
					fnLitRetStmt,
				}},
			}},
		}}},
		Decs: stdFuncDec(),
	}

	mockFn.Decs.Before = NewLine
	mockFn.Decs.Start.Append(fmt.Sprintf(
		"// %s returns the %s implementation of the %s type",
		c.exportName(mockFnName),
		mockIdent,
		typeName))

	return mockFn
}

// MockMethod generates a mock implementation of a method
func (c *Converter) MockMethod(typeName string, fn Func) *FuncDecl {
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

	return &FuncDecl{
		Recv: &FieldList{List: []*Field{{
			Names: []*Ident{NewIdent(mockReceiverIdent)},
			Type:  &StarExpr{X: NewIdent(recv)},
		}}},
		Name: NewIdent(fnName),
		Type: &FuncType{
			Params:  cloneAndNameUnnamed(paramPrefix, fn.Params),
			Results: cloneAndNameUnnamed(resultPrefix, fn.Results),
		},
		Body: c.mockFunc(typePrefix, fieldSuffix, fn),
		Decs: stdFuncDec(),
	}
}

// RecorderMethods generates a recorder implementation of a method and
// associated return method
func (c *Converter) RecorderMethods(typeName string, fn Func) (funcDecls []Decl) {
	return []Decl{
		c.recorderFn(typeName, fn),
		c.recorderReturnFn(typeName, fn),
		c.recorderTimesFn(typeName, fn),
		c.recorderAnyTimesFn(typeName, fn),
	}
}

// ResetMethod generates a method to reset the mock's state
func (c *Converter) ResetMethod(typeSpec *TypeSpec, funcs []Func) (funcDecl *FuncDecl) {
	mName := c.mockName(typeSpec.Name.Name)

	var stmts []Stmt
	for _, fn := range funcs {
		typePrefix := mName
		fieldSuffix := ""
		if _, ok := typeSpec.Type.(*InterfaceType); ok {
			typePrefix = fmt.Sprintf("%s_%s", mName, fn.Name)
			fieldSuffix = "_" + fn.Name
		}

		pName := fmt.Sprintf("%s_%s", typePrefix, paramsIdent)

		stmts = append(stmts, &AssignStmt{
			Lhs: []Expr{&SelectorExpr{
				X:   NewIdent(mockReceiverIdent),
				Sel: NewIdent(c.exportName(resultsByParamsIdent + fieldSuffix)),
			}},
			Tok: token.ASSIGN,
			Rhs: []Expr{&CompositeLit{Type: &MapType{
				Key: NewIdent(pName),
				Value: &StarExpr{
					X: NewIdent(fmt.Sprintf("%s_%s", typePrefix, resultMgrSuffix)),
				},
			}}},
		})
	}

	fn := &FuncDecl{
		Recv: &FieldList{List: []*Field{{
			Names: []*Ident{NewIdent(mockReceiverIdent)},
			Type:  &StarExpr{X: NewIdent(mName)},
		}}},
		Name: NewIdent(resetFnName),
		Type: &FuncType{},
		Body: &BlockStmt{List: stmts},
	}

	fn.Decs.Before = NewLine
	fn.Decs.Start.Append(fmt.Sprintf(
		"// %s resets the state of the mock", resetFnName))

	return fn
}

func (c *Converter) AssertMethod(typeSpec *TypeSpec, funcs []Func) (funcDecl *FuncDecl) {
	mName := c.mockName(typeSpec.Name.Name)

	var stmts []Stmt
	for _, fn := range funcs {
		fieldSuffix := ""
		if _, ok := typeSpec.Type.(*InterfaceType); ok {
			fieldSuffix = "_" + fn.Name
		}

		stmts = append(stmts, &RangeStmt{
			Key:   NewIdent(paramsIdent),
			Value: NewIdent(resultsIdent),
			Tok:   token.DEFINE,
			X: &SelectorExpr{
				X:   NewIdent(mockReceiverIdent),
				Sel: NewIdent(c.exportName(resultsByParamsIdent + fieldSuffix))},
			Body: &BlockStmt{List: []Stmt{
				&AssignStmt{
					Lhs: []Expr{NewIdent(missingIdent)},
					Tok: token.DEFINE,
					Rhs: []Expr{&BinaryExpr{
						X: &CallExpr{
							Fun: NewIdent("len"),
							Args: []Expr{&SelectorExpr{
								X:   NewIdent(resultsIdent),
								Sel: NewIdent(c.exportName(resultsIdent)),
							}},
						},
						Op: token.SUB,
						Y: &CallExpr{
							Fun: NewIdent("int"),
							Args: []Expr{&CallExpr{
								Fun: &Ident{
									Name: "LoadUint32",
									Path: "sync/atomic",
								},
								Args: []Expr{&UnaryExpr{
									Op: token.AND,
									X: &SelectorExpr{
										X:   NewIdent(resultsIdent),
										Sel: NewIdent(c.exportName(indexIdent)),
									},
								}},
							}},
						},
					}},
				},
				&IfStmt{
					Cond: &BinaryExpr{
						X: &BinaryExpr{
							X:  NewIdent(missingIdent),
							Op: token.EQL,
							Y:  &BasicLit{Kind: token.INT, Value: "1"},
						},
						Op: token.LAND,
						Y: &BinaryExpr{
							X: &SelectorExpr{
								X:   NewIdent(resultsIdent),
								Sel: NewIdent(c.exportName(anyTimesIdent)),
							},
							Op: token.EQL,
							Y:  NewIdent("true"),
						},
					},
					Body: &BlockStmt{List: []Stmt{&BranchStmt{Tok: token.CONTINUE}}},
				},
				&IfStmt{
					Cond: &BinaryExpr{
						X:  NewIdent(missingIdent),
						Op: token.GTR,
						Y:  &BasicLit{Kind: token.INT, Value: "0"},
					},
					Body: &BlockStmt{List: []Stmt{
						&ExprStmt{X: &CallExpr{
							Fun: &SelectorExpr{
								X: &SelectorExpr{
									X: &SelectorExpr{
										X:   NewIdent(mockReceiverIdent),
										Sel: NewIdent(c.exportName(sceneIdent)),
									},
									Sel: NewIdent(moqTType),
								},
								Sel: NewIdent(errorfFnName),
							},
							Args: []Expr{
								&BasicLit{
									Kind:  token.STRING,
									Value: "\"Expected %d additional call(s) with parameters %#v\"",
								},
								NewIdent(missingIdent),
								NewIdent(paramsIdent),
							},
						}},
					}},
				},
			}},
		})
	}

	fn := &FuncDecl{
		Recv: &FieldList{List: []*Field{{
			Names: []*Ident{NewIdent(mockReceiverIdent)},
			Type:  &StarExpr{X: NewIdent(mName)},
		}}},
		Name: NewIdent(assertFnName),
		Type: &FuncType{},
		Body: &BlockStmt{List: stmts},
	}

	fn.Decs.Before = NewLine
	fn.Decs.Start.Append(fmt.Sprintf(
		"// %s asserts that all expectations have been met", assertFnName))

	return fn
}

func (c *Converter) baseMockFieldList(typeSpec *TypeSpec, funcs []Func) []*Field {
	var fields []*Field

	fields = append(fields, &Field{
		Names: []*Ident{NewIdent(c.exportName(sceneIdent))},
		Type:  &StarExpr{X: &Ident{Name: sceneType, Path: moqPkg}},
	})
	fields = append(fields, &Field{
		Names: []*Ident{NewIdent(c.exportName(configIdent))},
		Type:  &Ident{Name: mockConfigType, Path: moqPkg},
	})

	mName := c.mockName(typeSpec.Name.Name)
	for _, fn := range funcs {
		typePrefix := mName
		fieldSuffix := ""
		if _, ok := typeSpec.Type.(*InterfaceType); ok {
			typePrefix = fmt.Sprintf("%s_%s", mName, fn.Name)
			fieldSuffix = "_" + fn.Name
		}
		fields = c.baseMockFieldsPerFn(fields, typePrefix, fieldSuffix)
	}

	return fields
}

func (c *Converter) baseMockFieldsPerFn(
	fields []*Field, typePrefix, fieldSuffix string,
) []*Field {
	pName := c.exportName(fmt.Sprintf("%s_%s", typePrefix, paramsIdent))
	fields = append(fields, &Field{
		Names: []*Ident{NewIdent(c.exportName(resultsByParamsIdent + fieldSuffix))},
		Type: &MapType{
			Key: NewIdent(pName),
			Value: &StarExpr{
				X: NewIdent(fmt.Sprintf("%s_%s", typePrefix, resultMgrSuffix)),
			},
		},
	})

	// TODO: param chans
	//fields = append(fields, &Field{
	//	Names: []*Ident{NewIdent(c.exportName(paramsIdent + fieldSuffix))},
	//	Type:  &ChanType{Dir: SEND | RECV, Value: NewIdent(pName)},
	//})

	return fields
}

func isComparable(expr Expr) bool {
	// TODO this logic needs to be expanded -- also should check structs recursively
	switch expr.(type) {
	case *ArrayType, *MapType, *Ellipsis:
		return false
	}

	return true
}

func (c *Converter) methodStruct(
	typeName, prefix, label string, fieldList *FieldList,
) *GenDecl {
	unnamedPrefix, comparable := labelDirection(label)
	fieldList = cloneNilableFieldList(fieldList)

	if fieldList == nil {
		// Result field lists can be nil (rather than containing an empty
		// list). Struct field lists cannot be nil.
		fieldList = &FieldList{}
	} else {
		for n, f := range fieldList.List {
			if len(f.Names) == 0 {
				f.Names = []*Ident{NewIdent(fmt.Sprintf("%s%d", unnamedPrefix, n+1))}
			}

			for nn := range f.Names {
				f.Names[nn] = NewIdent(c.exportName(f.Names[nn].Name))
			}

			// Map params are represented as a deep hash
			if comparable && !isComparable(f.Type) {
				f.Type = &Ident{Path: hashPkg, Name: "Hash"}
			}
		}
	}

	structName := fmt.Sprintf("%s_%s", prefix, label)
	sType := &GenDecl{
		Tok: token.TYPE,
		Specs: []Spec{&TypeSpec{
			Name: NewIdent(structName),
			Type: &StructType{Fields: fieldList},
		}},
	}
	sType.Decs.Before = NewLine
	sType.Decs.Start.Append(fmt.Sprintf("// %s holds the %s of the %s type",
		structName,
		label,
		typeName))

	return sType
}

func (c *Converter) resultMgrStruct(typeName, prefix string) *GenDecl {
	structName := fmt.Sprintf("%s_%s", prefix, resultMgrSuffix)

	sType := &GenDecl{
		Tok: token.TYPE,
		Specs: []Spec{&TypeSpec{
			Name: NewIdent(structName),
			Type: &StructType{Fields: &FieldList{List: []*Field{
				{
					Names: []*Ident{NewIdent(c.exportName(resultsIdent))},
					Type: &ArrayType{Elt: &StarExpr{
						X: NewIdent(c.exportName(fmt.Sprintf(
							"%s_%s", prefix, resultsIdent))),
					}},
				},
				{
					Names: []*Ident{NewIdent(c.exportName(indexIdent))},
					Type:  NewIdent("uint32"),
				},
				{
					Names: []*Ident{NewIdent(c.exportName(anyTimesIdent))},
					Type:  NewIdent("bool"),
				},
			}}},
		}},
	}
	sType.Decs.Before = NewLine
	sType.Decs.Start.Append(fmt.Sprintf(
		"// %s manages multiple results and the state of the %s type",
		structName,
		typeName))

	return sType
}

func (c *Converter) fnRecorderStruct(typeName string, prefix string) *GenDecl {
	mName := c.mockName(typeName)
	structName := fmt.Sprintf("%s_%s", prefix, fnRecorderSuffix)
	sType := &GenDecl{
		Tok: token.TYPE,
		Specs: []Spec{&TypeSpec{
			Name: NewIdent(structName),
			Type: &StructType{Fields: &FieldList{List: []*Field{
				{
					Names: []*Ident{NewIdent(c.exportName(paramsIdent))},
					Type:  NewIdent(fmt.Sprintf("%s_%s", prefix, paramsIdent)),
				},
				{
					Names: []*Ident{NewIdent(c.exportName(resultsIdent))},
					Type: &StarExpr{
						X: NewIdent(fmt.Sprintf("%s_%s", prefix, resultMgrSuffix)),
					},
				},
				{
					Names: []*Ident{NewIdent(c.exportName(mockIdent))},
					Type:  &StarExpr{X: NewIdent(mName)},
				},
			}}},
		}},
	}
	sType.Decs.Before = NewLine
	sType.Decs.Start.Append(fmt.Sprintf(
		"// %s routes recorded function calls to the %s mock",
		structName,
		mName))

	return sType
}

func (c *Converter) newMockElements(typeSpec *TypeSpec, funcs []Func) []Expr {
	var elems []Expr

	elems = append(elems, &KeyValueExpr{
		Key:   NewIdent(c.exportName(sceneIdent)),
		Value: NewIdent(sceneIdent),
		Decs:  KeyValueExprDecorations{NodeDecs: NodeDecs{After: NewLine}},
	})
	elems = append(elems, &KeyValueExpr{
		Key:   NewIdent(c.exportName(configIdent)),
		Value: &StarExpr{X: NewIdent(configIdent)},
		Decs:  KeyValueExprDecorations{NodeDecs: NodeDecs{After: NewLine}},
	})

	// TODO: param chans
	//mName := c.mockName(typeSpec.Name.Name)
	//for _, fn := range funcs {
	//	typePrefix := mName
	//	fieldSuffix := ""
	//	if _, ok := typeSpec.Type.(*InterfaceType); ok {
	//		typePrefix = fmt.Sprintf("%s_%s", mName, fn.Name)
	//		fieldSuffix = "_" + fn.Name
	//	}
	//	elems = append(elems, &KeyValueExpr{
	//		Key: NewIdent(c.exportName(paramsIdent + fieldSuffix)),
	//		Value: &CallExpr{
	//			Fun: NewIdent("make"),
	//			Args: []Expr{
	//				&ChanType{
	//					Dir: SEND | RECV,
	//					Value: NewIdent(
	//						fmt.Sprintf("%s_%s", typePrefix, paramsIdent)),
	//				},
	//				&BasicLit{Kind: token.INT, Value: "100"},
	//			},
	//		},
	//		Decs: KeyValueExprDecorations{NodeDecs: NodeDecs{After: NewLine}},
	//	})
	//}

	return elems
}

func (c *Converter) mockFunc(typePrefix, fieldSuffix string, fn Func) *BlockStmt {
	stateSelector := &SelectorExpr{
		X:   NewIdent(mockReceiverIdent),
		Sel: NewIdent(c.exportName(mockIdent)),
	}
	stmts := []Stmt{
		&AssignStmt{
			Lhs: []Expr{NewIdent(paramsIdent)},
			Tok: token.DEFINE,
			Rhs: []Expr{&CompositeLit{
				Type: NewIdent(fmt.Sprintf("%s_%s", typePrefix, paramsIdent)),
				Elts: c.passthroughElements(fn.Params, paramsIdent),
			}},
		},
		// TODO: param chans
		//&SendStmt{
		//	Chan: &SelectorExpr{
		//		X:   Clone(stateSelector).(Expr),
		//		Sel: NewIdent(c.exportName(paramsIdent + fieldSuffix)),
		//	},
		//	Value: NewIdent(paramsIdent),
		//},
	}

	stmts = append(stmts, &AssignStmt{
		Lhs: []Expr{NewIdent(resultsIdent), NewIdent(okIdent)},
		Tok: token.DEFINE,
		Rhs: []Expr{&IndexExpr{
			X: &SelectorExpr{
				X:   Clone(stateSelector).(Expr),
				Sel: NewIdent(c.exportName(resultsByParamsIdent + fieldSuffix)),
			},
			Index: NewIdent(paramsIdent),
		}},
	})
	stmts = append(stmts, &IfStmt{
		Cond: &UnaryExpr{Op: token.NOT, X: NewIdent(okIdent)},
		Body: &BlockStmt{List: []Stmt{
			&IfStmt{
				Cond: &BinaryExpr{
					X: &SelectorExpr{
						X: &SelectorExpr{
							X: &SelectorExpr{
								X:   NewIdent(mockReceiverIdent),
								Sel: NewIdent(c.exportName(mockIdent)),
							},
							Sel: NewIdent(c.exportName(configIdent)),
						},
						Sel: NewIdent(expectationIdent),
					},
					Op: token.EQL,
					Y:  &Ident{Name: strictIdent, Path: moqPkg},
				},
				Body: &BlockStmt{List: []Stmt{
					&ExprStmt{X: &CallExpr{
						Fun: &SelectorExpr{
							X: &SelectorExpr{
								X: &SelectorExpr{
									X: &SelectorExpr{
										X:   NewIdent(mockReceiverIdent),
										Sel: NewIdent(c.exportName(mockIdent)),
									},
									Sel: NewIdent(c.exportName(sceneIdent)),
								},
								Sel: NewIdent(moqTType),
							},
							Sel: NewIdent(fatalfFnName),
						},
						Args: []Expr{
							&BasicLit{
								Kind:  token.STRING,
								Value: "\"Unexpected call with parameters %#v\"",
							},
							NewIdent(paramsIdent),
						},
					}},
				}},
			},
			&ReturnStmt{},
		}},
	})

	stmts = append(stmts, &AssignStmt{
		Lhs: []Expr{NewIdent(iIdent)},
		Tok: token.DEFINE,
		Rhs: []Expr{&BinaryExpr{
			X: &CallExpr{
				Fun: NewIdent("int"),
				Args: []Expr{&CallExpr{
					Fun: &Ident{Name: "AddUint32", Path: "sync/atomic"},
					Args: []Expr{
						&UnaryExpr{
							Op: token.AND,
							X: &SelectorExpr{
								X:   NewIdent(resultsIdent),
								Sel: NewIdent(c.exportName(indexIdent)),
							},
						},
						&BasicLit{Kind: token.INT, Value: "1"},
					},
				}},
			},
			Op: token.SUB,
			Y:  &BasicLit{Kind: token.INT, Value: "1"},
		}},
		Decs: AssignStmtDecorations{NodeDecs: NodeDecs{Before: EmptyLine}},
	})
	stmts = append(stmts, &IfStmt{
		Cond: &BinaryExpr{
			X:  NewIdent(iIdent),
			Op: token.GEQ,
			Y: &CallExpr{
				Fun: NewIdent("len"),
				Args: []Expr{&SelectorExpr{
					X:   NewIdent(resultsIdent),
					Sel: NewIdent(c.exportName(resultsIdent)),
				}},
			},
		},
		Body: &BlockStmt{List: []Stmt{
			&IfStmt{
				Cond: &UnaryExpr{
					Op: token.NOT,
					X: &SelectorExpr{
						X:   NewIdent(resultsIdent),
						Sel: NewIdent(c.exportName(anyTimesIdent)),
					},
				},
				Body: &BlockStmt{List: []Stmt{
					&IfStmt{
						Cond: &BinaryExpr{
							X: &SelectorExpr{
								X: &SelectorExpr{
									X: &SelectorExpr{
										X:   NewIdent(mockReceiverIdent),
										Sel: NewIdent(c.exportName(mockIdent)),
									},
									Sel: NewIdent(c.exportName(configIdent)),
								},
								Sel: NewIdent(expectationIdent),
							},
							Op: token.EQL,
							Y:  &Ident{Name: strictIdent, Path: moqPkg},
						},
						Body: &BlockStmt{List: []Stmt{&ExprStmt{X: &CallExpr{
							Fun: &SelectorExpr{
								X: &SelectorExpr{
									X: &SelectorExpr{
										X: &SelectorExpr{
											X:   NewIdent(mockReceiverIdent),
											Sel: NewIdent(c.exportName(mockIdent)),
										},
										Sel: NewIdent(c.exportName(sceneIdent)),
									},
									Sel: NewIdent(moqTType),
								},
								Sel: NewIdent(fatalfFnName),
							},
							Args: []Expr{
								&BasicLit{
									Kind:  token.STRING,
									Value: "\"Too many calls to mock with parameters %#v\"",
								},
								NewIdent(paramsIdent),
							},
						}}}},
					},
					&ReturnStmt{},
				}},
			},
			&AssignStmt{
				Lhs: []Expr{NewIdent(iIdent)},
				Tok: token.ASSIGN,
				Rhs: []Expr{&BinaryExpr{
					X: &CallExpr{
						Fun: NewIdent("len"),
						Args: []Expr{&SelectorExpr{
							X:   NewIdent(resultsIdent),
							Sel: NewIdent(c.exportName(resultsIdent)),
						}},
					},
					Op: token.SUB,
					Y:  &BasicLit{Kind: token.INT, Value: "1"},
				}},
			},
		}},
	})

	if fn.Results != nil {
		stmts = append(stmts, &AssignStmt{
			Lhs: []Expr{NewIdent(resultIdent)},
			Tok: token.DEFINE,
			Rhs: []Expr{&IndexExpr{
				X: &SelectorExpr{
					X:   NewIdent(resultsIdent),
					Sel: NewIdent(c.exportName(resultsIdent)),
				},
				Index: NewIdent(iIdent),
			}},
		})
		stmts = append(stmts, c.assignResult(fn.Results)...)
	}

	stmts = append(stmts, &ReturnStmt{})

	return &BlockStmt{List: stmts}
}

func (c *Converter) recorderFn(typeName string, fn Func) *FuncDecl {
	mName := c.mockName(typeName)

	recvType := fmt.Sprintf("%s_%s", mName, recorderIdent)
	fnName := fn.Name
	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	typePrefix := fmt.Sprintf("%s_%s", mName, fn.Name)
	var mockVal Expr = &SelectorExpr{
		X:   NewIdent(mockReceiverIdent),
		Sel: NewIdent(c.exportName(mockIdent)),
	}
	if fn.Name == "" {
		recvType = mName
		fnName = c.exportName(onCallFnName)
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
		typePrefix = mName
		mockVal = NewIdent(mockReceiverIdent)
	}

	return &FuncDecl{
		Recv: &FieldList{List: []*Field{{
			Names: []*Ident{NewIdent(mockReceiverIdent)},
			Type:  &StarExpr{X: NewIdent(recvType)},
		}}},
		Name: NewIdent(fnName),
		Type: &FuncType{
			Params: cloneAndNameUnnamed(paramPrefix, fn.Params),
			Results: &FieldList{List: []*Field{{
				Type: &StarExpr{X: NewIdent(fnRecName)},
			}}},
		},
		Body: c.recorderFnInterfaceBody(fnRecName, typePrefix, mockVal, fn),
		Decs: stdFuncDec(),
	}
}

func (c *Converter) recorderFnInterfaceBody(
	fnRecName, typePrefix string, mockValue Expr, fn Func,
) *BlockStmt {
	return &BlockStmt{List: []Stmt{&ReturnStmt{
		Results: []Expr{&UnaryExpr{
			Op: token.AND,
			X: &CompositeLit{
				Type: NewIdent(fnRecName),
				Elts: []Expr{
					&KeyValueExpr{
						Key: NewIdent(c.exportName(paramsIdent)),
						Value: &CompositeLit{
							Type: NewIdent(fmt.Sprintf(
								"%s_%s", typePrefix, paramsIdent)),
							Elts: c.passthroughElements(fn.Params, paramsIdent),
						},
						Decs: KeyValueExprDecorations{
							NodeDecs: NodeDecs{After: NewLine},
						},
					},
					&KeyValueExpr{
						Key:   NewIdent(c.exportName(mockIdent)),
						Value: mockValue,
						Decs: KeyValueExprDecorations{
							NodeDecs: NodeDecs{After: NewLine},
						},
					},
				},
				Decs: CompositeLitDecorations{Lbrace: []string{"\n"}},
			},
		}},
	}}}
}

func (c *Converter) recorderReturnFn(typeName string, fn Func) *FuncDecl {
	mName := c.mockName(typeName)

	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	results := fmt.Sprintf("%s_%s_%s", mName, fn.Name, resultsIdent)
	resultMgr := fmt.Sprintf("%s_%s_%s", mName, fn.Name, resultMgrSuffix)
	resultsByParams := fmt.Sprintf("%s_%s", resultsByParamsIdent, fn.Name)
	mockSel := &SelectorExpr{X: &SelectorExpr{
		X:   NewIdent(recorderReceiverIdent),
		Sel: NewIdent(c.exportName(mockIdent)),
	}}
	if fn.Name == "" {
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
		results = fmt.Sprintf("%s_%s", mName, resultsIdent)
		resultMgr = fmt.Sprintf("%s_%s", mName, resultMgrSuffix)
		resultsByParams = resultsByParamsIdent
		mockSel = &SelectorExpr{X: &SelectorExpr{
			X:   NewIdent(recorderReceiverIdent),
			Sel: NewIdent(c.exportName(mockIdent)),
		}}
	}

	return &FuncDecl{
		Recv: &FieldList{List: []*Field{{
			Names: []*Ident{NewIdent(recorderReceiverIdent)},
			Type:  &StarExpr{X: NewIdent(fnRecName)},
		}}},
		Name: NewIdent(c.exportName(returnFnName)),
		Type: &FuncType{
			Params: cloneAndNameUnnamed(resultPrefix, fn.Results),
			Results: &FieldList{List: []*Field{{
				Type: &StarExpr{X: NewIdent(fnRecName)},
			}}},
		},
		Body: &BlockStmt{List: []Stmt{
			&IfStmt{
				Cond: &BinaryExpr{
					X: &SelectorExpr{
						X:   NewIdent(recorderReceiverIdent),
						Sel: NewIdent(c.exportName(resultsIdent)),
					},
					Op: token.EQL,
					Y:  NewIdent("nil"),
				},
				Body: &BlockStmt{List: []Stmt{
					&IfStmt{
						Init: &AssignStmt{
							Lhs: []Expr{NewIdent("_"), NewIdent(okIdent)},
							Tok: token.DEFINE,
							Rhs: []Expr{&IndexExpr{
								X: cloneSelect(mockSel, c.exportName(resultsByParams)),
								Index: &SelectorExpr{
									X:   NewIdent(recorderReceiverIdent),
									Sel: NewIdent(c.exportName(paramsIdent)),
								},
							}},
						},
						Cond: NewIdent(okIdent),
						Body: &BlockStmt{List: []Stmt{
							&ExprStmt{X: &CallExpr{
								Fun: &SelectorExpr{
									X: &SelectorExpr{
										X: cloneSelect(
											mockSel, c.exportName(sceneIdent)),
										Sel: NewIdent(moqTType),
									},
									Sel: NewIdent(fatalfFnName),
								},
								Args: []Expr{
									&BasicLit{
										Kind:  token.STRING,
										Value: "\"Expectations already recorded for mock with parameters %#v\"",
									},
									&SelectorExpr{
										X:   NewIdent(recorderReceiverIdent),
										Sel: NewIdent(c.exportName(paramsIdent)),
									},
								},
							}},
							&ReturnStmt{Results: []Expr{NewIdent("nil")}},
						}},
						Decs: IfStmtDecorations{NodeDecs: NodeDecs{
							After: EmptyLine,
						}},
					},
					&AssignStmt{
						Lhs: []Expr{&SelectorExpr{
							X:   NewIdent(recorderReceiverIdent),
							Sel: NewIdent(c.exportName(resultsIdent)),
						}},
						Tok: token.ASSIGN,
						Rhs: []Expr{&UnaryExpr{
							Op: token.AND,
							X: &CompositeLit{
								Type: NewIdent(c.exportName(resultMgr)),
								Elts: []Expr{
									&KeyValueExpr{
										Key: NewIdent(c.exportName(resultsIdent)),
										Value: &CompositeLit{Type: &ArrayType{
											Elt: &StarExpr{
												X: NewIdent(c.exportName(results)),
											},
										}},
									},
									&KeyValueExpr{
										Key:   NewIdent(c.exportName(indexIdent)),
										Value: &BasicLit{Kind: token.INT, Value: "0"},
									},
									&KeyValueExpr{
										Key:   NewIdent(c.exportName(anyTimesIdent)),
										Value: NewIdent("false"),
									},
								},
							},
						}},
					},
					&AssignStmt{
						Lhs: []Expr{&IndexExpr{
							X: &SelectorExpr{
								X: &SelectorExpr{
									X:   NewIdent(recorderReceiverIdent),
									Sel: NewIdent(c.exportName(mockIdent)),
								},
								Sel: NewIdent(c.exportName(resultsByParams)),
							},
							Index: &SelectorExpr{
								X:   NewIdent(recorderReceiverIdent),
								Sel: NewIdent(c.exportName(paramsIdent)),
							},
						}},
						Tok: token.ASSIGN,
						Rhs: []Expr{&SelectorExpr{
							X:   NewIdent(recorderReceiverIdent),
							Sel: NewIdent(c.exportName(resultsIdent)),
						}},
					},
				}},
			},
			&AssignStmt{
				Lhs: []Expr{
					&SelectorExpr{
						X: &SelectorExpr{
							X:   NewIdent(recorderReceiverIdent),
							Sel: NewIdent(c.exportName(resultsIdent)),
						},
						Sel: NewIdent(c.exportName(resultsIdent)),
					},
				},
				Tok: token.ASSIGN,
				Rhs: []Expr{&CallExpr{
					Fun: NewIdent("append"),
					Args: []Expr{
						&SelectorExpr{
							X: &SelectorExpr{
								X:   NewIdent(recorderReceiverIdent),
								Sel: NewIdent(c.exportName(resultsIdent)),
							},
							Sel: NewIdent(c.exportName(resultsIdent)),
						},
						&UnaryExpr{
							Op: token.AND,
							X: &CompositeLit{
								Type: NewIdent(c.exportName(results)),
								Elts: c.passthroughElements(fn.Results, resultsIdent),
							},
						},
					},
				}},
			},
			&ReturnStmt{Results: []Expr{NewIdent(recorderReceiverIdent)}},
		}},
		Decs: stdFuncDec(),
	}
}

func (c *Converter) recorderTimesFn(typeName string, fn Func) *FuncDecl {
	mName := c.mockName(typeName)

	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
	}

	return &FuncDecl{
		Recv: &FieldList{List: []*Field{{
			Names: []*Ident{NewIdent(recorderReceiverIdent)},
			Type:  &StarExpr{X: NewIdent(fnRecName)},
		}}},
		Name: NewIdent(c.exportName(timesFnName)),
		Type: &FuncType{
			Params: &FieldList{List: []*Field{{
				Names: []*Ident{NewIdent(countIdent)},
				Type:  NewIdent("int"),
			}}},
			Results: &FieldList{List: []*Field{{
				Type: &StarExpr{X: NewIdent(fnRecName)},
			}}},
		},
		Body: &BlockStmt{List: []Stmt{
			&IfStmt{
				Cond: &BinaryExpr{
					X: &SelectorExpr{
						X:   NewIdent(recorderReceiverIdent),
						Sel: NewIdent(c.exportName(resultsIdent)),
					},
					Op: token.EQL,
					Y:  NewIdent("nil"),
				},
				Body: &BlockStmt{List: []Stmt{
					&ExprStmt{X: &CallExpr{
						Fun: &SelectorExpr{
							X: &SelectorExpr{
								X: &SelectorExpr{
									X: &SelectorExpr{
										X:   NewIdent(recorderReceiverIdent),
										Sel: NewIdent(c.exportName(mockIdent)),
									},
									Sel: NewIdent(c.exportName(sceneIdent)),
								},
								Sel: NewIdent(moqTType),
							},
							Sel: NewIdent(fatalfFnName),
						},
						Args: []Expr{&BasicLit{
							Kind:  token.STRING,
							Value: "\"Return must be called before calling Times\"",
						}},
					}},
					&ReturnStmt{Results: []Expr{NewIdent("nil")}},
				}},
			},
			&AssignStmt{
				Lhs: []Expr{NewIdent(lastIdent)},
				Tok: token.DEFINE,
				Rhs: []Expr{&IndexExpr{
					X: &SelectorExpr{
						X: &SelectorExpr{
							X:   NewIdent(recorderReceiverIdent),
							Sel: NewIdent(c.exportName(resultsIdent)),
						},
						Sel: NewIdent(c.exportName(resultsIdent)),
					},
					Index: &BinaryExpr{
						X: &CallExpr{
							Fun: NewIdent("len"),
							Args: []Expr{&SelectorExpr{
								X: &SelectorExpr{
									X:   NewIdent(recorderReceiverIdent),
									Sel: NewIdent(c.exportName(resultsIdent)),
								},
								Sel: NewIdent(c.exportName(resultsIdent)),
							}},
						},
						Op: token.SUB,
						Y:  &BasicLit{Kind: token.INT, Value: "1"},
					},
				}},
			},
			&ForStmt{
				Init: &AssignStmt{
					Lhs: []Expr{NewIdent("n")},
					Tok: token.DEFINE,
					Rhs: []Expr{&BasicLit{Kind: token.INT, Value: "0"}},
				},
				Cond: &BinaryExpr{
					X:  NewIdent("n"),
					Op: token.LSS,
					Y: &BinaryExpr{
						X:  NewIdent(countIdent),
						Op: token.SUB,
						Y:  &BasicLit{Kind: token.INT, Value: "1"},
					},
				},
				Post: &IncDecStmt{X: NewIdent("n"), Tok: token.INC},
				Body: &BlockStmt{List: []Stmt{&AssignStmt{
					Lhs: []Expr{&SelectorExpr{
						X: &SelectorExpr{
							X:   NewIdent(recorderReceiverIdent),
							Sel: NewIdent(c.exportName(resultsIdent)),
						},
						Sel: NewIdent(c.exportName(resultsIdent)),
					}},
					Tok: token.ASSIGN,
					Rhs: []Expr{&CallExpr{
						Fun: NewIdent("append"),
						Args: []Expr{
							&SelectorExpr{
								X: &SelectorExpr{
									X:   NewIdent(recorderReceiverIdent),
									Sel: NewIdent(c.exportName(resultsIdent)),
								},
								Sel: NewIdent(c.exportName(resultsIdent)),
							},
							NewIdent(lastIdent),
						},
					}},
				}}},
			},
			&ReturnStmt{Results: []Expr{NewIdent(recorderReceiverIdent)}},
		}},
		Decs: stdFuncDec(),
	}
}

func (c *Converter) recorderAnyTimesFn(typeName string, fn Func) *FuncDecl {
	mName := c.mockName(typeName)

	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
	}

	return &FuncDecl{
		Recv: &FieldList{List: []*Field{{
			Names: []*Ident{NewIdent(recorderReceiverIdent)},
			Type:  &StarExpr{X: NewIdent(fnRecName)},
		}}},
		Name: NewIdent(c.exportName(anyTimesFnName)),
		Type: &FuncType{Params: &FieldList{}},
		Body: &BlockStmt{List: []Stmt{
			&IfStmt{
				Cond: &BinaryExpr{
					X: &SelectorExpr{
						X:   NewIdent(recorderReceiverIdent),
						Sel: NewIdent(c.exportName(resultsIdent)),
					},
					Op: token.EQL,
					Y:  NewIdent("nil"),
				},
				Body: &BlockStmt{List: []Stmt{
					&ExprStmt{X: &CallExpr{
						Fun: &SelectorExpr{
							X: &SelectorExpr{
								X: &SelectorExpr{
									X: &SelectorExpr{
										X:   NewIdent(recorderReceiverIdent),
										Sel: NewIdent(c.exportName(mockIdent)),
									},
									Sel: NewIdent(c.exportName(sceneIdent)),
								},
								Sel: NewIdent(moqTType),
							},
							Sel: NewIdent(fatalfFnName),
						},
						Args: []Expr{&BasicLit{
							Kind:  token.STRING,
							Value: "\"Return must be called before calling AnyTimes\"",
						}},
					}},
					&ReturnStmt{},
				}},
			},
			&AssignStmt{
				Lhs: []Expr{&SelectorExpr{
					X: &SelectorExpr{
						X:   NewIdent(recorderReceiverIdent),
						Sel: NewIdent(c.exportName(resultsIdent)),
					},
					Sel: NewIdent(c.exportName(anyTimesIdent)),
				}},
				Tok: token.ASSIGN,
				Rhs: []Expr{NewIdent("true")},
			},
		}},
		Decs: stdFuncDec(),
	}
}

func (c *Converter) passthroughElements(fl *FieldList, label string) []Expr {
	unnamedPrefix, comparable := labelDirection(label)
	var elts []Expr
	if fl != nil {
		beforeDec := NewLine
		fields := fl.List
		for n, field := range fields {
			if len(field.Names) == 0 {
				pName := fmt.Sprintf("%s%d", unnamedPrefix, n+1)
				elts = append(elts, &KeyValueExpr{
					Key:   NewIdent(c.exportName(pName)),
					Value: passthroughValue(pName, field, comparable),
					Decs: KeyValueExprDecorations{
						NodeDecs: NodeDecs{Before: beforeDec, After: NewLine},
					},
				})
				beforeDec = None
			}

			for _, name := range field.Names {
				elts = append(elts, &KeyValueExpr{
					Key:   NewIdent(c.exportName(name.Name)),
					Value: passthroughValue(name.Name, field, comparable),
					Decs: KeyValueExprDecorations{
						NodeDecs: NodeDecs{Before: beforeDec, After: NewLine},
					},
				})
				beforeDec = None
			}
		}
	}

	return elts
}

func passthroughValue(name string, field *Field, comparable bool) Expr {
	var val Expr
	val = NewIdent(name)
	if comparable && !isComparable(field.Type) {
		val = &CallExpr{
			Fun:  &Ident{Path: hashPkg, Name: "DeepHash"},
			Args: []Expr{val},
		}
	}
	return val
}

func passthroughFields(fields *FieldList) []Expr {
	var exprs []Expr
	for _, f := range fields.List {
		for _, n := range f.Names {
			exprs = append(exprs, NewIdent(n.Name))
		}
	}
	return exprs
}

func (c *Converter) assignResult(resFL *FieldList) []Stmt {
	var assigns []Stmt
	if resFL != nil {
		results := resFL.List
		for n, result := range results {
			if len(result.Names) == 0 {
				rName := fmt.Sprintf("%s%d", resultPrefix, n+1)
				assigns = append(assigns, &AssignStmt{
					Lhs: []Expr{NewIdent(rName)},
					Tok: token.ASSIGN,
					Rhs: []Expr{&SelectorExpr{
						X:   NewIdent(resultIdent),
						Sel: NewIdent(c.exportName(rName)),
					}},
				})
			}

			for _, name := range result.Names {
				assigns = append(assigns, &AssignStmt{
					Lhs: []Expr{NewIdent(name.Name)},
					Tok: token.ASSIGN,
					Rhs: []Expr{&SelectorExpr{
						X:   NewIdent(resultIdent),
						Sel: NewIdent(c.exportName(name.Name)),
					}},
				})
			}
		}
	}
	return assigns
}

func cloneAndNameUnnamed(prefix string, fieldList *FieldList) *FieldList {
	fieldList = cloneNilableFieldList(fieldList)
	if fieldList != nil {
		for n, f := range fieldList.List {
			if len(f.Names) == 0 {
				f.Names = []*Ident{NewIdent(fmt.Sprintf("%s%d", prefix, n+1))}
			}
		}
	}
	return fieldList
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

func stdFuncDec() FuncDeclDecorations {
	return FuncDeclDecorations{
		NodeDecs: NodeDecs{Before: EmptyLine, After: EmptyLine},
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

func cloneNilableFieldList(fl *FieldList) *FieldList {
	if fl != nil {
		fl = Clone(fl).(*FieldList)
	}
	return fl
}

func cloneSelect(x *SelectorExpr, sel string) *SelectorExpr {
	x = Clone(x).(*SelectorExpr)
	x.Sel = NewIdent(sel)
	return x
}
