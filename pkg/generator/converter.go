package generator

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/dave/dst"

	. "github.com/myshkin5/moqueries/pkg/ast"
	"github.com/myshkin5/moqueries/pkg/logs"
)

const (
	moqueriesPkg  = "github.com/myshkin5/moqueries"
	hashPkg       = moqueriesPkg + "/pkg/hash"
	moqPkg        = moqueriesPkg + "/pkg/moq"
	syncAtomicPkg = "sync/atomic"

	intType        = "int"
	mockConfigType = "MockConfig"
	moqTType       = "MoqT"
	sceneType      = "Scene"

	anyCountIdent         = "anyCount"
	anyParamsIdent        = "anyParams"
	anyTimesIdent         = "anyTimes"
	configIdent           = "config"
	countIdent            = "count"
	expectationIdent      = "Expectation"
	iIdent                = "i"
	indexIdent            = "index"
	insertAtIdent         = "insertAt"
	lastIdent             = "last"
	missingIdent          = "missing"
	mockIdent             = "mock"
	moqIdent              = "moq"
	mockReceiverIdent     = "m"
	nIdent                = "n"
	nilIdent              = "nil"
	okIdent               = "ok"
	paramsIdent           = "params"
	paramsKeyIdent        = "paramsKey"
	recorderIdent         = "recorder"
	recorderReceiverIdent = "r"
	resultsByParamsIdent  = "resultsByParams"
	resIdent              = "res"
	resultIdent           = "result"
	resultsIdent          = "results"
	sceneIdent            = "scene"
	sequenceIdent         = "sequence"
	strictIdent           = "Strict"

	anyTimesFnName = "anyTimes"
	assertFnName   = "AssertExpectationsMet"
	errorfFnName   = "Errorf"
	fatalfFnName   = "Fatalf"
	fnFnName       = "fn"
	lenFnName      = "len"
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
	isExported bool
}

// NewConverter creates a new Converter
func NewConverter(isExported bool) *Converter {
	return &Converter{
		isExported: isExported,
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
	fields := []*dst.Field{
		Field(Star(IdPath(sceneType, moqPkg))).Names(Id(c.export(sceneIdent))).Obj,
		Field(IdPath(mockConfigType, moqPkg)).Names(Id(c.export(configIdent))).Obj,
	}

	mName := c.mockName(typeSpec.Name.Name)
	for _, fn := range funcs {
		typePrefix := mName
		fieldSuffix := ""
		if _, ok := typeSpec.Type.(*dst.InterfaceType); ok {
			typePrefix = fmt.Sprintf("%s_%s", mName, fn.Name)
			fieldSuffix = "_" + fn.Name
		}
		fields = append(fields,
			Field(SliceType(Idf("%s_%s", typePrefix, resultsByParamsIdent))).
				Names(Id(c.export(resultsByParamsIdent+fieldSuffix))).Obj)
	}

	return Struct(mName).Fields(fields...).Decs(genDeclDec(
		"// %s holds the state of a mock of the %s type",
		mName, typeSpec.Name.Name)).Obj
}

// IsolationStruct generates a struct used to isolate an interface for the mock
func (c *Converter) IsolationStruct(typeName, suffix string) (structDecl *dst.GenDecl) {
	mName := c.mockName(typeName)
	iName := fmt.Sprintf("%s_%s", mName, suffix)

	return Struct(iName).
		Fields(Field(Star(Id(mName))).Names(Id(c.export(mockIdent))).Obj).
		Decs(genDeclDec("// %s isolates the %s interface of the %s type",
			iName, suffix, typeName)).Obj
}

// MethodStructs generates a structure for storing a set of parameters or
// a set of results for a method invocation of a mock
func (c *Converter) MethodStructs(typeSpec *dst.TypeSpec, fn Func) []dst.Decl {
	prefix := c.mockName(typeSpec.Name.Name)
	if _, ok := typeSpec.Type.(*dst.InterfaceType); ok {
		prefix = fmt.Sprintf("%s_%s", prefix, fn.Name)
	}

	return []dst.Decl{
		c.methodStruct(typeSpec.Name.Name, prefix, paramsIdent, fn.Params),
		c.methodStruct(typeSpec.Name.Name, prefix, paramsKeyIdent, fn.Params),
		c.resultByParamsStruct(typeSpec.Name.Name, prefix),
		c.resultMgrStruct(typeSpec.Name.Name, prefix),
		c.methodStruct(typeSpec.Name.Name, prefix, resultsIdent, fn.Results),
		c.fnRecorderStruct(typeSpec.Name.Name, prefix),
	}
}

// NewFunc generates a function for constructing a mock
func (c *Converter) NewFunc(typeSpec *dst.TypeSpec) (funcDecl *dst.FuncDecl) {
	fnName := c.export("newMock" + typeSpec.Name.Name)
	mName := c.mockName(typeSpec.Name.Name)
	return Fn(fnName).
		Params(
			Field(Star(IdPath(sceneType, moqPkg))).Names(Id(sceneIdent)).Obj,
			Field(Star(IdPath(mockConfigType, moqPkg))).Names(Id(configIdent)).Obj,
		).
		Results(Field(Star(Id(mName))).Obj).
		Body(
			If(Bin(Id(configIdent)).Op(token.EQL).Y(Id(nilIdent)).Obj).Body(
				Assign(Id(configIdent)).Tok(token.ASSIGN).Rhs(Un(token.AND,
					Comp(IdPath(mockConfigType, moqPkg)).Obj)).Obj).Obj,
			Assign(Id(mockReceiverIdent)).Tok(token.DEFINE).
				Rhs(Un(token.AND, Comp(Id(mName)).Elts(
					Key(Id(c.export(sceneIdent))).
						Value(Id(sceneIdent)).Decs(kvExprDec(dst.None)).Obj,
					Key(Id(c.export(configIdent))).
						Value(Star(Id(configIdent))).Decs(kvExprDec(dst.None)).Obj,
				).Decs(litDec()).Obj)).Obj,
			Expr(Call(Sel(Id(sceneIdent)).Dot(Id("AddMock")).Obj).
				Args(Id(mockReceiverIdent)).Obj),
			Return(Id(mockReceiverIdent)),
		).
		Decs(fnDeclDec("// %s creates a new mock of the %s type",
			fnName, typeSpec.Name.Name)).Obj
}

// IsolationAccessor generates a function to access an isolation interface
func (c *Converter) IsolationAccessor(typeName, suffix, fnName string) (funcDecl *dst.FuncDecl) {
	fnName = c.export(fnName)
	mName := c.mockName(typeName)
	iName := fmt.Sprintf("%s_%s", mName, suffix)
	return Fn(fnName).
		Recv(Field(Star(Id(mName))).Names(Id(mockReceiverIdent)).Obj).
		Results(Field(Star(Id(iName))).Obj).
		Body(Return(Un(token.AND, Comp(Id(iName)).Elts(
			Key(Id(c.export(mockIdent))).Value(Id(mockReceiverIdent)).
				Decs(kvExprDec(dst.None)).Obj).Decs(litDec()).Obj,
		))).
		Decs(fnDeclDec("// %s returns the %s implementation of the %s type",
			fnName, suffix, typeName)).Obj
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
	fnLitCall := Call(Sel(Id(mockIdent)).Dot(Id(c.export(fnFnName))).Obj).
		Args(passthroughFields(paramPrefix, fn.Params)...).
		Ellipsis(ellipsis).Obj
	var fnLitRetStmt dst.Stmt
	fnLitRetStmt = Return(fnLitCall)
	if fn.Results == nil {
		fnLitRetStmt = Expr(fnLitCall)
	}

	return Fn(c.export(mockFnName)).
		Recv(Field(Star(Id(mName))).Names(Id(mockReceiverIdent)).Obj).
		Results(Field(IdPath(typeName, pkgPath)).Obj).
		Body(Return(FnLit(FnType(cloneAndNameUnnamed(paramPrefix, fn.Params)).
			Results(cloneNilableFieldList(fn.Results)).Obj).
			Body(Assign(Id(mockIdent)).
				Tok(token.DEFINE).
				Rhs(Un(
					token.AND,
					Comp(Idf("%s_%s", mName, mockIdent)).
						Elts(Key(Id(c.export(mockIdent))).
							Value(Id(mockReceiverIdent)).Obj).Obj,
				)).Obj,
				fnLitRetStmt,
			).Obj)).
		Decs(fnDeclDec("// %s returns the %s implementation of the %s type",
			c.export(mockFnName), mockIdent, typeName)).Obj
}

// MockMethod generates a mock implementation of a method
func (c *Converter) MockMethod(typeName string, fn Func) *dst.FuncDecl {
	mName := c.mockName(typeName)
	recv := fmt.Sprintf("%s_%s", mName, mockIdent)

	fnName := fn.Name
	fieldSuffix := "_" + fn.Name
	typePrefix := fmt.Sprintf("%s_%s", mName, fn.Name)
	if fnName == "" {
		fnName = c.export(fnFnName)
		fieldSuffix = ""
		typePrefix = mName
	}

	return Fn(fnName).
		Recv(Field(Star(Id(recv))).Names(Id(mockReceiverIdent)).Obj).
		ParamFieldList(cloneAndNameUnnamed(paramPrefix, fn.Params)).
		ResultFieldList(cloneAndNameUnnamed(resultPrefix, fn.Results)).
		Body(c.mockFunc(typePrefix, fieldSuffix, fn)...).
		Decs(stdFuncDec()).Obj
}

// RecorderMethods generates a recorder implementation of a method and
// associated return method
func (c *Converter) RecorderMethods(typeName string, fn Func) (funcDecls []dst.Decl) {
	decls := []dst.Decl{
		c.recorderFn(typeName, fn),
	}

	decls = append(decls, c.anyParamFns(typeName, fn)...)
	decls = append(decls, c.recorderSeqFns(typeName, fn)...)
	decls = append(decls,
		c.recorderReturnFn(typeName, fn),
		c.recorderTimesFn(typeName, fn),
		c.recorderAnyTimesFn(typeName, fn),
	)

	return decls
}

// ResetMethod generates a method to reset the mock's state
func (c *Converter) ResetMethod(typeSpec *dst.TypeSpec, funcs []Func) (funcDecl *dst.FuncDecl) {
	mName := c.mockName(typeSpec.Name.Name)

	var stmts []dst.Stmt
	for _, fn := range funcs {
		fieldSuffix := ""
		if _, ok := typeSpec.Type.(*dst.InterfaceType); ok {
			fieldSuffix = "_" + fn.Name
		}

		stmts = append(stmts, Assign(Sel(Id(mockReceiverIdent)).
			Dot(Id(c.export(resultsByParamsIdent+fieldSuffix))).Obj).
			Tok(token.ASSIGN).
			Rhs(Id(nilIdent)).Obj)
	}

	return Fn(resetFnName).
		Recv(Field(Star(Id(mName))).Names(Id(mockReceiverIdent)).Obj).
		Body(stmts...).
		Decs(fnDeclDec("// %s resets the state of the mock", resetFnName)).Obj
}

// AssertMethod generates a method to assert all expectations are met
func (c *Converter) AssertMethod(typeSpec *dst.TypeSpec, funcs []Func) (funcDecl *dst.FuncDecl) {
	mName := c.mockName(typeSpec.Name.Name)

	var stmts []dst.Stmt
	for _, fn := range funcs {
		fieldSuffix := ""
		if _, ok := typeSpec.Type.(*dst.InterfaceType); ok {
			fieldSuffix = "_" + fn.Name
		}

		stmts = append(stmts, Range(Sel(Id(mockReceiverIdent)).
			Dot(Id(c.export(resultsByParamsIdent+fieldSuffix))).Obj).
			Key(Id("_")).Value(Id(resIdent)).Tok(token.DEFINE).
			Body(Range(Sel(Id(resIdent)).Dot(Id(c.export(resultsIdent))).Obj).
				Key(Id("_")).Value(Id(resultsIdent)).Tok(token.DEFINE).
				Body(
					Assign(Id(missingIdent)).
						Tok(token.DEFINE).
						Rhs(Bin(Call(Id(lenFnName)).Args(Sel(Id(resultsIdent)).
							Dot(Id(c.export(resultsIdent))).Obj).Obj).
							Op(token.SUB).
							Y(Call(Id(intType)).Args(
								Call(IdPath("LoadUint32", syncAtomicPkg)).Args(Un(
									token.AND,
									Sel(Id(resultsIdent)).
										Dot(Id(c.export(indexIdent))).Obj)).Obj).Obj).Obj).Obj,
					If(Bin(Bin(Id(missingIdent)).Op(token.EQL).
						Y(LitInt(1)).Obj).
						Op(token.LAND).
						Y(Bin(Sel(Id(resultsIdent)).
							Dot(Id(c.export(anyTimesIdent))).Obj).
							Op(token.EQL).
							Y(Id("true")).Obj).Obj).
						Body(Continue()).Obj,
					If(Bin(Id(missingIdent)).Op(token.GTR).Y(LitInt(0)).Obj).
						Body(
							Expr(Call(Sel(Sel(Sel(Id(mockReceiverIdent)).
								Dot(Id(c.export(sceneIdent))).Obj).
								Dot(Id(moqTType)).Obj).Dot(Id(errorfFnName)).Obj).
								Args(
									LitString("Expected %d additional call(s) with parameters %#v"),
									Id(missingIdent),
									Sel(Id(resultsIdent)).Dot(Id(c.export(paramsIdent))).Obj).Obj),
						).Obj,
				).Obj,
			).Obj)
	}

	return Fn(assertFnName).
		Recv(Field(Star(Id(mName))).Names(Id(mockReceiverIdent)).Obj).
		Body(stmts...).
		Decs(fnDeclDec("// %s asserts that all expectations have been met",
			assertFnName)).Obj
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
		fieldList = FieldList()
	} else {
		for n, f := range fieldList.List {
			if len(f.Names) == 0 {
				f.Names = []*dst.Ident{Idf("%s%d", unnamedPrefix, n+1)}
			}

			for nn := range f.Names {
				f.Names[nn] = Id(c.export(f.Names[nn].Name))
			}

			f.Type = comparableType(comparable, f.Type)
		}
	}

	if label == resultsIdent {
		fieldList.List = append(fieldList.List,
			Field(Id("uint32")).
				Names(Idf("%s_%s", c.export(moqIdent), c.export(sequenceIdent))).Obj)
	}

	goDocDesc := label
	if label == paramsKeyIdent {
		goDocDesc = "map key params"
	}

	structName := fmt.Sprintf("%s_%s", prefix, label)
	return Struct(structName).FieldList(fieldList).Decs(genDeclDec(
		"// %s holds the %s of the %s type",
		structName, goDocDesc, typeName)).Obj
}

func comparableType(needComparable bool, typ dst.Expr) dst.Expr {
	if needComparable && !isComparable(typ) {
		// Non-comparable params are represented as a deep hash
		return IdPath("Hash", hashPkg)
	} else if ellipsis, ok := typ.(*dst.Ellipsis); ok {
		// Ellipsis params are represented as a slice (when not comparable)
		return SliceType(ellipsis.Elt)
	}

	return typ
}

func (c *Converter) resultByParamsStruct(typeName, prefix string) *dst.GenDecl {
	structName := fmt.Sprintf("%s_%s", prefix, resultsByParamsIdent)

	return Struct(structName).Fields(
		Field(Id(intType)).Names(Id(c.export(anyCountIdent))).Obj,
		Field(Id("uint64")).Names(Id(c.export(anyParamsIdent))).Obj,
		Field(MapType(Id(c.export(fmt.Sprintf("%s_%s", prefix, paramsKeyIdent)))).
			Value(Star(Idf("%s_%s", prefix, resultMgrSuffix))).Obj,
		).Names(Id(c.export(resultsIdent))).Obj,
	).Decs(genDeclDec(
		"// %s contains the results for a given set of parameters for the %s type",
		structName,
		typeName)).Obj
}

func (c *Converter) resultMgrStruct(typeName, prefix string) *dst.GenDecl {
	structName := fmt.Sprintf("%s_%s", prefix, resultMgrSuffix)

	return Struct(structName).Fields(
		Field(Idf("%s_%s", prefix, paramsIdent)).
			Names(Id(c.export(paramsIdent))).Obj,
		Field(SliceType(Star(Id(c.export(fmt.Sprintf(
			"%s_%s", prefix, resultsIdent)))))).
			Names(Id(c.export(resultsIdent))).Obj,
		Field(Id("uint32")).Names(Id(c.export(indexIdent))).Obj,
		Field(Id("bool")).Names(Id(c.export(anyTimesIdent))).Obj,
	).Decs(genDeclDec(
		"// %s manages multiple results and the state of the %s type",
		structName,
		typeName)).Obj
}

func (c *Converter) fnRecorderStruct(typeName string, prefix string) *dst.GenDecl {
	mName := c.mockName(typeName)
	structName := fmt.Sprintf("%s_%s", prefix, fnRecorderSuffix)
	return Struct(structName).Fields(
		Field(Idf("%s_%s", prefix, paramsIdent)).
			Names(Id(c.export(paramsIdent))).Obj,
		Field(Idf("%s_%s", prefix, paramsKeyIdent)).
			Names(Id(c.export(paramsKeyIdent))).Obj,
		Field(Id("uint64")).
			Names(Id(c.export(anyParamsIdent))).Obj,
		Field(Id("bool")).
			Names(Id(c.export(sequenceIdent))).Obj,
		Field(Star(Idf("%s_%s", prefix, resultMgrSuffix))).
			Names(Id(c.export(resultsIdent))).Obj,
		Field(Star(Id(mName))).
			Names(Id(c.export(mockIdent))).Obj,
	).Decs(genDeclDec("// %s routes recorded function calls to the %s mock",
		structName, mName)).Obj
}

func (c *Converter) mockFunc(typePrefix, fieldSuffix string, fn Func) []dst.Stmt {
	stateSelector := Sel(Id(mockReceiverIdent)).Dot(Id(c.export(mockIdent))).Obj

	stmts := []dst.Stmt{
		Assign(Id(paramsIdent)).
			Tok(token.DEFINE).
			Rhs(Comp(Idf("%s_%s", typePrefix, paramsIdent)).
				Elts(c.passthroughElements(fn.Params, paramsIdent, "", nil)...).Obj).Obj,
		Var(Value(Star(Idf("%s_%s", typePrefix, resultMgrSuffix))).
			Names(Id(resultsIdent)).Obj),
		Range(Sel(dst.Clone(stateSelector).(dst.Expr)).
			Dot(Id(c.export(resultsByParamsIdent + fieldSuffix))).Obj).
			Key(Id("_")).
			Value(Id(resultsByParamsIdent)).
			Tok(token.DEFINE).
			Body(c.mockFuncFindResults(typePrefix, fn)...).Obj,
	}

	stmts = append(stmts,
		If(Bin(Id(resultsIdent)).Op(token.EQL).Y(Id(nilIdent)).Obj).Body(
			If(Bin(Sel(Sel(dst.Clone(stateSelector).(dst.Expr)).
				Dot(Id(c.export(configIdent))).Obj).
				Dot(Id(expectationIdent)).Obj).
				Op(token.EQL).
				Y(IdPath(strictIdent, moqPkg)).Obj).
				Body(
					Expr(Call(Sel(Sel(Sel(dst.Clone(stateSelector).(dst.Expr)).
						Dot(Id(c.export(sceneIdent))).Obj).
						Dot(Id(moqTType)).Obj).
						Dot(Id(fatalfFnName)).Obj).
						Args(LitString("Unexpected call with parameters %#v"),
							Id(paramsIdent)).Obj)).Obj,
			Return(),
		).Obj)

	stmts = append(stmts, Assign(Id(iIdent)).
		Tok(token.DEFINE).
		Rhs(Bin(Call(Id(intType)).
			Args(Call(IdPath("AddUint32", syncAtomicPkg)).Args(Un(
				token.AND,
				Sel(Id(resultsIdent)).Dot(Id(c.export(indexIdent))).Obj),
				LitInt(1)).Obj).Obj).
			Op(token.SUB).
			Y(LitInt(1)).Obj).
		Decs(AssignDecs(dst.EmptyLine).Obj).Obj)
	stmts = append(stmts,
		If(Bin(Id(iIdent)).Op(token.GEQ).Y(Call(Id(lenFnName)).
			Args(Sel(Id(resultsIdent)).
				Dot(Id(c.export(resultsIdent))).Obj).Obj).Obj).
			Body(
				If(Un(token.NOT, Sel(Id(resultsIdent)).
					Dot(Id(c.export(anyTimesIdent))).Obj)).
					Body(
						If(Bin(Sel(Sel(dst.Clone(stateSelector).(dst.Expr)).
							Dot(Id(c.export(configIdent))).Obj).
							Dot(Id(expectationIdent)).Obj).
							Op(token.EQL).
							Y(IdPath(strictIdent, moqPkg)).Obj).
							Body(Expr(Call(Sel(Sel(Sel(dst.Clone(stateSelector).(dst.Expr)).
								Dot(Id(c.export(sceneIdent))).Obj).
								Dot(Id(moqTType)).Obj).
								Dot(Id(fatalfFnName)).Obj).
								Args(
									LitString("Too many calls to mock with parameters %#v"),
									Id(paramsIdent),
								).Obj)).Obj,
						Return(),
					).Obj,
				Assign(Id(iIdent)).
					Tok(token.ASSIGN).
					Rhs(Bin(Call(Id(lenFnName)).
						Args(Sel(Id(resultsIdent)).
							Dot(Id(c.export(resultsIdent))).Obj).Obj).
						Op(token.SUB).
						Y(LitInt(1)).Obj).Obj,
			).Decs(IfDecs(dst.EmptyLine).Obj).Obj)

	stmts = append(stmts, Assign(Id(resultIdent)).
		Tok(token.DEFINE).
		Rhs(Index(Sel(Id(resultsIdent)).
			Dot(Id(c.export(resultsIdent))).Obj).Sub(Id(iIdent)).Obj).Obj)
	stmts = append(stmts, If(Bin(
		Sel(Id(resultIdent)).Dot(
			Idf("%s_%s", c.export(moqIdent), c.export(sequenceIdent))).Obj).
		Op(token.NEQ).
		Y(LitInt(0)).Obj).Body(
		Assign(Id(sequenceIdent)).Tok(token.DEFINE).Rhs(Call(
			Sel(Sel(Sel(Id(mockReceiverIdent)).
				Dot(Id(c.export(mockIdent))).Obj).Dot(Id(c.export(sceneIdent))).Obj).
				Dot(Id("NextMockSequence")).Obj).Obj).Obj,
		If(Bin(Paren(Bin(Un(token.NOT, Sel(Id(resultsIdent)).
			Dot(Id(c.export(anyTimesIdent))).Obj)).Op(token.LAND).
			Y(Bin(Sel(Id(resultIdent)).Dot(
				Idf("%s_%s", c.export(moqIdent), c.export(sequenceIdent))).Obj).
				Op(token.NEQ).Y(Id(sequenceIdent)).Obj).Obj)).Op(token.LOR).
			Y(Bin(Sel(Id(resultIdent)).
				Dot(Idf("%s_%s", c.export(moqIdent), c.export(sequenceIdent))).Obj).
				Op(token.GTR).Y(Id(sequenceIdent)).Obj).Obj).Body(
			Expr(Call(Sel(Sel(Sel(dst.Clone(stateSelector).(dst.Expr)).
				Dot(Id(c.export(sceneIdent))).Obj).
				Dot(Id(moqTType)).Obj).
				Dot(Id(fatalfFnName)).Obj).
				Args(LitString("Call sequence does not match %#v"), Id(paramsIdent)).Obj),
		).Obj,
	).Decs(IfDecs(dst.EmptyLine).Obj).Obj)

	if fn.Results != nil {
		stmts = append(stmts, c.assignResult(fn.Results)...)
	}

	stmts = append(stmts, Return())

	return stmts
}

func (c *Converter) mockFuncFindResults(typePrefix string, fn Func) []dst.Stmt {
	var stmts []dst.Stmt
	var paramPos int
	for n, param := range fn.Params.List {
		if len(param.Names) == 0 {
			pName := fmt.Sprintf("%s%d", paramPrefix, n+1)
			stmts = append(stmts, c.mockFuncFindResultsParam(
				resultsByParamsIdent, nil, pName, paramPos, param.Type)...)
			paramPos++
		}

		for _, name := range param.Names {
			stmts = append(stmts, c.mockFuncFindResultsParam(
				resultsByParamsIdent, nil, name.Name, paramPos, param.Type)...)
			paramPos++
		}
	}

	stmts = append(stmts, Assign(Id(paramsKeyIdent)).
		Tok(token.DEFINE).
		Rhs(Comp(Idf("%s_%s", typePrefix, paramsKeyIdent)).
			// Passing through as params not paramsKey as hashing is already done
			Elts(c.passthroughElements(fn.Params, paramsIdent, "Used", nil)...).Obj).Obj)
	stmts = append(stmts, Var(Value(Id("bool")).Names(Id(okIdent)).Obj))
	stmts = append(stmts, Assign(Id(resultsIdent), Id(okIdent)).
		Tok(token.ASSIGN).
		Rhs(Index(Sel(Id(resultsByParamsIdent)).Dot(Id(c.export(resultsIdent))).Obj).
			Sub(Id(paramsKeyIdent)).Obj).Obj)
	stmts = append(stmts, If(Id(okIdent)).Body(Break()).Obj)

	return stmts
}

func (c *Converter) mockFuncFindResultsParam(
	sel string, pKeySel *dst.SelectorExpr, pName string, paramPos int, typ dst.Expr,
) []dst.Stmt {
	comparable := true
	if pKeySel != nil {
		comparable = false
	}
	pUsed := fmt.Sprintf("%sUsed", pName)
	return []dst.Stmt{
		Var(Value(comparableType(true, dst.Clone(typ).(dst.Expr))).
			Names(Id(pUsed)).Obj),
		If(Bin(Bin(Sel(Id(sel)).Dot(Id(c.export(anyParamsIdent))).Obj).
			Op(token.AND).
			Y(Paren(Bin(LitInt(1)).Op(token.SHL).Y(LitInt(paramPos)).Obj)).Obj).
			Op(token.EQL).
			Y(LitInt(0)).Obj).
			Body(Assign(Id(pUsed)).
				Tok(token.ASSIGN).
				Rhs(c.passthroughValue(Id(pName), typ, comparable, pKeySel)).Obj).Obj,
	}
}

func (c *Converter) recorderFn(typeName string, fn Func) *dst.FuncDecl {
	mName := c.mockName(typeName)

	recvType := fmt.Sprintf("%s_%s", mName, recorderIdent)
	fnName := fn.Name
	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	typePrefix := fmt.Sprintf("%s_%s", mName, fn.Name)
	var mockVal dst.Expr = Sel(Id(mockReceiverIdent)).
		Dot(Id(c.export(mockIdent))).Obj
	if fn.Name == "" {
		recvType = mName
		fnName = c.export(onCallFnName)
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
		typePrefix = mName
		mockVal = Id(mockReceiverIdent)
	}

	return Fn(fnName).
		Recv(Field(Star(Id(recvType))).Names(Id(mockReceiverIdent)).Obj).
		ParamFieldList(cloneAndNameUnnamed(paramPrefix, fn.Params)).
		Results(Field(Star(Id(fnRecName))).Obj).
		Body(c.recorderFnInterfaceBody(fnRecName, typePrefix, mockVal, fn)...).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) recorderFnInterfaceBody(
	fnRecName, typePrefix string, mockVal dst.Expr, fn Func,
) []dst.Stmt {
	return []dst.Stmt{Return(Un(
		token.AND,
		Comp(Id(fnRecName)).
			Elts(
				Key(Id(c.export(paramsIdent))).
					Value(Comp(Idf("%s_%s", typePrefix, paramsIdent)).
						Elts(c.passthroughElements(fn.Params, paramsIdent, "", nil)...).Obj,
					).Decs(kvExprDec(dst.None)).Obj,
				Key(Id(c.export(paramsKeyIdent))).
					Value(Comp(Idf("%s_%s", typePrefix, paramsKeyIdent)).
						Elts(c.passthroughElements(fn.Params, paramsKeyIdent, "", nil)...).Obj,
					).Decs(kvExprDec(dst.None)).Obj,
				Key(Id(c.export(sequenceIdent))).
					Value(Bin(Sel(Sel(mockVal).
						Dot(Id(c.export(configIdent))).Obj).
						Dot(Id(strings.Title(sequenceIdent))).Obj).
						Op(token.EQL).
						Y(IdPath("SeqDefaultOn", moqPkg)).Obj).
					Decs(kvExprDec(dst.None)).Obj,
				Key(Id(c.export(mockIdent))).
					Value(dst.Clone(mockVal).(dst.Expr)).Decs(kvExprDec(dst.None)).Obj,
			).Decs(litDec()).Obj,
	))}
}

func (c *Converter) anyParamFns(typeName string, fn Func) []dst.Decl {
	mName := c.mockName(typeName)

	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
	}

	var decls []dst.Decl
	var paramPos int
	for n, param := range fn.Params.List {
		if len(param.Names) == 0 {
			pName := fmt.Sprintf("%s%d", paramPrefix, n+1)
			decls = append(decls, c.anyParamFn(fnRecName, pName, paramPos))
			paramPos++
		}

		for _, name := range param.Names {
			decls = append(decls, c.anyParamFn(fnRecName, name.Name, paramPos))
			paramPos++
		}
	}
	return decls
}

func (c *Converter) anyParamFn(fnRecName, pName string, paramPos int) dst.Decl {
	mockSel := Sel(Sel(Id(recorderReceiverIdent)).
		Dot(Id(c.export(mockIdent))).Obj).Obj

	return Fn(c.export("any"+strings.Title(pName))).
		Recv(Field(Star(Id(fnRecName))).Names(Id(recorderReceiverIdent)).Obj).
		Results(Field(Star(Id(fnRecName))).Obj).
		Body(
			If(Bin(Sel(Id(recorderReceiverIdent)).Dot(Id(c.export(resultsIdent))).Obj).
				Op(token.NEQ).
				Y(Id(nilIdent)).Obj).
				Body(
					Expr(Call(Sel(Sel(cloneSelect(mockSel, c.export(sceneIdent))).
						Dot(Id(moqTType)).Obj).
						Dot(Id(fatalfFnName)).Obj).
						Args(LitString(
							"Any functions must be called prior to returning results, parameters: %#v"),
							Sel(Id(recorderReceiverIdent)).
								Dot(Id(c.export(paramsIdent))).Obj,
						).Obj),
					Return(Id(nilIdent))).Obj,
			Assign(Sel(Id(recorderReceiverIdent)).
				Dot(Id(c.export(anyParamsIdent))).Obj).
				Tok(token.OR_ASSIGN).
				Rhs(Bin(LitInt(1)).Op(token.SHL).Y(LitInt(paramPos)).Obj).Obj,
			Return(Id(recorderReceiverIdent)),
		).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) recorderReturnFn(typeName string, fn Func) *dst.FuncDecl {
	mName := c.mockName(typeName)

	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	results := fmt.Sprintf("%s_%s_%s", mName, fn.Name, resultsIdent)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
		results = fmt.Sprintf("%s_%s", mName, resultsIdent)
	}

	passThrough := c.passthroughElements(fn.Results, resultsIdent, "", nil)
	passThrough = append(passThrough,
		Key(Idf("%s_%s", c.export(moqIdent), c.export(sequenceIdent))).
			Value(Id(sequenceIdent)).
			Decs(kvExprDec(dst.None)).Obj)

	return Fn(c.export(returnFnName)).
		Recv(Field(Star(Id(fnRecName))).Names(Id(recorderReceiverIdent)).Obj).
		ParamFieldList(cloneAndNameUnnamed(resultPrefix, fn.Results)).
		Results(Field(Star(Id(fnRecName))).Obj).
		Body(
			If(Bin(Sel(Id(recorderReceiverIdent)).
				Dot(Id(c.export(resultsIdent))).Obj).
				Op(token.EQL).
				Y(Id(nilIdent)).Obj).
				Body(c.findRecorderResults(typeName, fn)...).
				Decs(IfDecs(dst.EmptyLine).Obj).Obj,
			Var(Value(Id("uint32")).Names(Id(sequenceIdent)).Obj),
			If(Sel(Id(recorderReceiverIdent)).Dot(Id(c.export(sequenceIdent))).Obj).
				Body(Assign(Id(sequenceIdent)).Tok(token.ASSIGN).Rhs(
					Call(Sel(Sel(Sel(Id(recorderReceiverIdent)).
						Dot(Id(c.export(mockIdent))).Obj).
						Dot(Id(c.export(sceneIdent))).Obj).
						Dot(Id("NextRecorderSequence")).Obj).Obj,
				).Obj).
				Decs(IfDecs(dst.EmptyLine).Obj).Obj,
			Assign(
				Sel(Sel(Id(recorderReceiverIdent)).
					Dot(Id(c.export(resultsIdent))).Obj).
					Dot(Id(c.export(resultsIdent))).Obj).
				Tok(token.ASSIGN).
				Rhs(Call(Id("append")).Args(Sel(Sel(Id(recorderReceiverIdent)).
					Dot(Id(c.export(resultsIdent))).Obj).
					Dot(Id(c.export(resultsIdent))).Obj,
					Un(token.AND, Comp(Id(c.export(results))).
						Elts(passThrough...).Obj)).Obj).Obj,
			Return(Id(recorderReceiverIdent)),
		).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) findRecorderResults(typeName string, fn Func) []dst.Stmt {
	mName := c.mockName(typeName)

	results := fmt.Sprintf("%s_%s_%s", mName, fn.Name, resultsIdent)
	resultsByParamsType := fmt.Sprintf("%s_%s_%s", mName, fn.Name, resultsByParamsIdent)
	resultMgr := fmt.Sprintf("%s_%s_%s", mName, fn.Name, resultMgrSuffix)
	paramsKey := fmt.Sprintf("%s_%s_%s", mName, fn.Name, paramsKeyIdent)
	resultsByParams := fmt.Sprintf("%s_%s", resultsByParamsIdent, fn.Name)
	if fn.Name == "" {
		results = fmt.Sprintf("%s_%s", mName, resultsIdent)
		resultsByParamsType = fmt.Sprintf("%s_%s", mName, resultsByParamsIdent)
		resultMgr = fmt.Sprintf("%s_%s", mName, resultMgrSuffix)
		paramsKey = fmt.Sprintf("%s_%s", mName, paramsKeyIdent)
		resultsByParams = resultsByParamsIdent
	}

	mockSel := Sel(Sel(Id(recorderReceiverIdent)).
		Dot(Id(c.export(mockIdent))).Obj).Obj

	stmts := []dst.Stmt{
		Assign(Id(anyCountIdent)).
			Tok(token.DEFINE).
			Rhs(Call(IdPath("OnesCount64", "math/bits")).Args(
				Sel(Id(recorderReceiverIdent)).
					Dot(Id(c.export(anyParamsIdent))).Obj).Obj).Obj,
		Assign(Id(insertAtIdent)).Tok(token.DEFINE).Rhs(LitInt(-1)).Obj,
		Var(Value(Star(Id(resultsByParamsType))).
			Names(Id(resultsIdent)).Obj),
		Range(cloneSelect(mockSel, c.export(resultsByParams))).
			Key(Id(nIdent)).
			Value(Id(resIdent)).
			Tok(token.DEFINE).
			Body(
				If(Bin(Sel(Id(resIdent)).Dot(Id(c.export(anyParamsIdent))).Obj).
					Op(token.EQL).
					Y(Sel(Id(recorderReceiverIdent)).
						Dot(Id(c.export(anyParamsIdent))).Obj).Obj).Body(
					Assign(Id(resultsIdent)).
						Tok(token.ASSIGN).
						Rhs(Un(token.AND, Id(resIdent))).Obj,
					Break(),
				).Obj,
				If(Bin(Sel(Id(resIdent)).Dot(Id(c.export(anyCountIdent))).Obj).
					Op(token.GTR).
					Y(Id(anyCountIdent)).Obj).Body(
					Assign(Id(insertAtIdent)).
						Tok(token.ASSIGN).
						Rhs(Id(nIdent)).Obj,
				).Obj,
			).Obj,
		If(Bin(Id(resultsIdent)).Op(token.EQL).Y(Id(nilIdent)).Obj).Body(
			Assign(Id(resultsIdent)).Tok(token.ASSIGN).Rhs(Un(
				token.AND, Comp(Id(resultsByParamsType)).Elts(
					Key(Id(c.export(anyCountIdent))).Value(Id(anyCountIdent)).
						Decs(kvExprDec(dst.NewLine)).Obj,
					Key(Id(c.export(anyParamsIdent))).
						Value(Sel(Id(recorderReceiverIdent)).
							Dot(Id(c.export(anyParamsIdent))).Obj).
						Decs(kvExprDec(dst.None)).Obj,
					Key(Id(c.export(resultsIdent))).Value(Comp(MapType(Id(paramsKey)).
						Value(Star(Id(resultMgr))).Obj).Obj).
						Decs(kvExprDec(dst.NewLine)).Obj,
				).Obj)).Obj,
			Assign(cloneSelect(mockSel, c.export(resultsByParams))).
				Tok(token.ASSIGN).Rhs(Call(Id("append")).Args(
				cloneSelect(mockSel, c.export(resultsByParams)),
				Star(Id(resultsIdent))).Obj).Obj,
			If(Bin(Bin(Id(insertAtIdent)).Op(token.NEQ).
				Y(LitInt(-1)).Obj).Op(token.LAND).
				Y(Bin(Bin(Id(insertAtIdent)).Op(token.ADD).
					Y(LitInt(1)).Obj).Op(token.LSS).Y(Call(Id(lenFnName)).
					Args(cloneSelect(mockSel,
						c.export(resultsByParams))).Obj).Obj).Obj).Body(
				Expr(Call(Id("copy")).Args(
					SliceExpr(cloneSelect(mockSel, c.export(resultsByParams))).
						Low(Bin(Id(insertAtIdent)).Op(token.ADD).Y(LitInt(1)).Obj).Obj,
					SliceExpr(cloneSelect(mockSel, c.export(resultsByParams))).
						Low(Id(insertAtIdent)).High(LitInt(0)).Obj).Obj),
				Assign(Index(cloneSelect(mockSel, c.export(resultsByParams))).
					Sub(Id(insertAtIdent)).Obj).
					Tok(token.ASSIGN).
					Rhs(Star(Id(resultsIdent))).Obj,
			).Obj,
		).Decs(IfDecs(dst.EmptyLine).Obj).Obj,
	}

	pKeySel := Sel(Sel(Id(recorderReceiverIdent)).
		Dot(Id(c.export(paramsKeyIdent))).Obj).Obj
	var paramPos int
	for n, param := range fn.Params.List {
		if len(param.Names) == 0 {
			pName := fmt.Sprintf("%s%d", paramPrefix, n+1)
			stmts = append(stmts, c.mockFuncFindResultsParam(
				recorderReceiverIdent, pKeySel, pName, paramPos, param.Type)...)
			paramPos++
		}

		for _, name := range param.Names {
			stmts = append(stmts, c.mockFuncFindResultsParam(
				recorderReceiverIdent, pKeySel, name.Name, paramPos, param.Type)...)
			paramPos++
		}
	}

	stmts = append(stmts,
		Assign(Id(paramsKeyIdent)).
			Tok(token.DEFINE).
			Rhs(Comp(Id(paramsKey)).
				// Passing through as params not paramsKey as hashing is already done
				Elts(c.passthroughElements(fn.Params, paramsIdent, "Used", nil)...).Obj).
			Decs(AssignDecs(dst.None).After(dst.EmptyLine).Obj).Obj,
		IfInit(Assign(Id("_"), Id(okIdent)).
			Tok(token.DEFINE).
			Rhs(Index(Sel(Id(resultsIdent)).Dot(Id(c.export(resultsIdent))).Obj).
				Sub(Id(paramsKeyIdent)).Obj).Obj).
			Cond(Id(okIdent)).
			Body(
				Expr(Call(Sel(Sel(cloneSelect(mockSel, c.export(sceneIdent))).
					Dot(Id(moqTType)).Obj).
					Dot(Id(fatalfFnName)).Obj).
					Args(LitString(
						"Expectations already recorded for mock with parameters %#v"),
						Sel(Id(recorderReceiverIdent)).
							Dot(Id(c.export(paramsIdent))).Obj,
					).Obj),
				Return(Id(nilIdent)),
			).Decs(IfDecs(dst.EmptyLine).Obj).Obj,
		Assign(Sel(Id(recorderReceiverIdent)).
			Dot(Id(c.export(resultsIdent))).Obj).
			Tok(token.ASSIGN).
			Rhs(Un(
				token.AND,
				Comp(Id(c.export(resultMgr))).
					Elts(
						Key(Id(c.export(paramsIdent))).
							Value(Sel(Id(recorderReceiverIdent)).
								Dot(Id(c.export(paramsIdent))).Obj).
							Decs(kvExprDec(dst.NewLine)).Obj,
						Key(Id(c.export(resultsIdent))).Value(
							Comp(SliceType(Star(Id(c.export(results))))).Obj).
							Decs(kvExprDec(dst.None)).Obj,
						Key(Id(c.export(indexIdent))).Value(
							LitInt(0)).Decs(kvExprDec(dst.None)).Obj,
						Key(Id(c.export(anyTimesIdent))).Value(
							Id("false")).Decs(kvExprDec(dst.None)).Obj,
					).Obj,
			)).Obj,
		Assign(Index(Sel(Id(resultsIdent)).Dot(Id(c.export(resultsIdent))).Obj).
			Sub(Id(paramsKeyIdent)).Obj).
			Tok(token.ASSIGN).
			Rhs(Sel(Id(recorderReceiverIdent)).
				Dot(Id(c.export(resultsIdent))).Obj).Obj,
	)

	return stmts
}

func (c *Converter) recorderTimesFn(typeName string, fn Func) *dst.FuncDecl {
	mName := c.mockName(typeName)

	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	results := fmt.Sprintf("%s_%s_%s", mName, fn.Name, resultsIdent)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
		results = fmt.Sprintf("%s_%s", mName, resultsIdent)
	}

	passThrough := c.passthroughElements(fn.Results, resultsIdent, "", Sel(Id(lastIdent)).Obj)
	passThrough = append(passThrough,
		Key(Idf("%s_%s", c.export(moqIdent), c.export(sequenceIdent))).
			Value(Call(Sel(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(Id(c.export(mockIdent))).Obj).
				Dot(Id(c.export(sceneIdent))).Obj).
				Dot(Id("NextRecorderSequence")).Obj).Obj).
			Decs(kvExprDec(dst.None)).Obj)

	return Fn(c.export(timesFnName)).
		Recv(Field(Star(Id(fnRecName))).Names(Id(recorderReceiverIdent)).Obj).
		Params(Field(Id(intType)).Names(Id(countIdent)).Obj).
		Results(Field(Star(Id(fnRecName))).Obj).
		Body(
			If(Bin(Sel(Id(recorderReceiverIdent)).
				Dot(Id(c.export(resultsIdent))).Obj).
				Op(token.EQL).
				Y(Id(nilIdent)).Obj).
				Body(
					Expr(Call(Sel(Sel(Sel(Sel(Id(recorderReceiverIdent)).
						Dot(Id(c.export(mockIdent))).Obj).
						Dot(Id(c.export(sceneIdent))).Obj).
						Dot(Id(moqTType)).Obj).Dot(Id(fatalfFnName)).Obj).
						Args(LitString(
							"Return must be called before calling Times")).Obj),
					Return(Id(nilIdent)),
				).Obj,
			Assign(Id(lastIdent)).
				Tok(token.DEFINE).
				Rhs(Index(Sel(Sel(Id(recorderReceiverIdent)).
					Dot(Id(c.export(resultsIdent))).Obj).
					Dot(Id(c.export(resultsIdent))).Obj).
					Sub(Bin(Call(Id(lenFnName)).
						Args(Sel(Sel(Id(recorderReceiverIdent)).
							Dot(Id(c.export(resultsIdent))).Obj).
							Dot(Id(c.export(resultsIdent))).Obj).Obj).
						Op(token.SUB).
						Y(LitInt(1)).Obj).Obj).Obj,
			For(Assign(Id("n")).Tok(token.DEFINE).Rhs(LitInt(0)).Obj).
				Cond(Bin(Id("n")).Op(token.LSS).
					Y(Bin(Id(countIdent)).Op(token.SUB).Y(LitInt(1)).Obj).Obj).
				Post(IncStmt(Id("n"))).Body(
				If(Bin(Sel(Id(lastIdent)).
					Dot(Idf("%s_%s", c.export(moqIdent), c.export(sequenceIdent))).Obj).
					Op(token.NEQ).Y(LitInt(0)).Obj).Body(
					Assign(Id(lastIdent)).Tok(token.ASSIGN).Rhs(
						Un(token.AND, Comp(Id(c.export(results))).
							Elts(passThrough...).Obj)).Obj).Obj,
				Assign(Sel(Sel(Id(recorderReceiverIdent)).
					Dot(Id(c.export(resultsIdent))).Obj).
					Dot(Id(c.export(resultsIdent))).Obj).
					Tok(token.ASSIGN).
					Rhs(Call(Id("append")).Args(Sel(Sel(Id(recorderReceiverIdent)).
						Dot(Id(c.export(resultsIdent))).Obj).
						Dot(Id(c.export(resultsIdent))).Obj,
						Id(lastIdent)).Obj).Obj,
			).Obj,
			Return(Id(recorderReceiverIdent)),
		).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) recorderAnyTimesFn(typeName string, fn Func) *dst.FuncDecl {
	mName := c.mockName(typeName)

	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
	}

	return Fn(c.export(anyTimesFnName)).
		Recv(Field(Star(Id(fnRecName))).Names(Id(recorderReceiverIdent)).Obj).
		Body(
			If(Bin(Sel(Id(recorderReceiverIdent)).
				Dot(Id(c.export(resultsIdent))).Obj).
				Op(token.EQL).
				Y(Id(nilIdent)).Obj).
				Body(
					Expr(Call(Sel(Sel(Sel(Sel(Id(recorderReceiverIdent)).
						Dot(Id(c.export(mockIdent))).Obj).
						Dot(Id(c.export(sceneIdent))).Obj).
						Dot(Id(moqTType)).Obj).
						Dot(Id(fatalfFnName)).Obj).
						Args(LitString(
							"Return must be called before calling AnyTimes")).Obj),
					Return(),
				).Obj,
			Assign(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(Id(c.export(resultsIdent))).Obj).
				Dot(Id(c.export(anyTimesIdent))).Obj).
				Tok(token.ASSIGN).
				Rhs(Id("true")).Obj,
		).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) recorderSeqFns(typeName string, fn Func) []dst.Decl {
	return []dst.Decl{
		c.recorderSeqFn("seq", "true", typeName, fn),
		c.recorderSeqFn("noSeq", "false", typeName, fn),
	}
}

func (c *Converter) recorderSeqFn(fnName, assign, typeName string, fn Func) *dst.FuncDecl {
	mName := c.mockName(typeName)

	fnRecName := fmt.Sprintf("%s_%s_%s", mName, fn.Name, fnRecorderSuffix)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf("%s_%s", mName, fnRecorderSuffix)
	}

	fnName = c.export(fnName)
	return Fn(fnName).
		Results(Field(Star(Id(fnRecName))).Obj).
		Recv(Field(Star(Id(fnRecName))).Names(Id(recorderReceiverIdent)).Obj).
		Body(
			If(Bin(Sel(Id(recorderReceiverIdent)).
				Dot(Id(c.export(resultsIdent))).Obj).
				Op(token.NEQ).
				Y(Id(nilIdent)).Obj).
				Body(
					Expr(Call(Sel(Sel(Sel(Sel(Id(recorderReceiverIdent)).
						Dot(Id(c.export(mockIdent))).Obj).
						Dot(Id(c.export(sceneIdent))).Obj).
						Dot(Id(moqTType)).Obj).
						Dot(Id(fatalfFnName)).Obj).
						Args(LitString(fmt.Sprintf(
							"%s must be called prior to returning results, parameters: %%#v", fnName)),
							Sel(Id(recorderReceiverIdent)).
								Dot(Id(c.export(paramsIdent))).Obj).Obj),
					Return(Id(nilIdent)),
				).Obj,
			Assign(Sel(Id(recorderReceiverIdent)).
				Dot(Id(c.export(sequenceIdent))).Obj).
				Tok(token.ASSIGN).
				Rhs(Id(assign)).
				Decs(AssignDecs(dst.NewLine).Obj).Obj,
			Return(Id(recorderReceiverIdent)),
		).Decs(stdFuncDec()).Obj
}

func (c *Converter) passthroughElements(
	fl *dst.FieldList, label, valSuffix string, sel *dst.SelectorExpr,
) []dst.Expr {
	unnamedPrefix, comparable := labelDirection(label)
	var elts []dst.Expr
	if fl != nil {
		beforeDec := dst.NewLine
		fields := fl.List
		for n, field := range fields {
			if len(field.Names) == 0 {
				pName := fmt.Sprintf("%s%d", unnamedPrefix, n+1)
				elts = append(elts, Key(Id(c.export(pName))).Value(
					c.passthroughValue(Id(pName+valSuffix), field.Type, comparable, sel)).
					Decs(kvExprDec(beforeDec)).Obj)
				beforeDec = dst.None
			}

			for _, name := range field.Names {
				elts = append(elts, Key(Id(c.export(name.Name))).Value(
					c.passthroughValue(Id(name.Name+valSuffix), field.Type, comparable, sel)).
					Decs(kvExprDec(beforeDec)).Obj)
				beforeDec = dst.None
			}
		}
	}

	return elts
}

func (c *Converter) passthroughValue(
	src *dst.Ident, typ dst.Expr, comparable bool, sel *dst.SelectorExpr,
) dst.Expr {
	var val dst.Expr = src
	if sel != nil {
		val = cloneSelect(sel, c.export(src.Name))
	}
	if comparable && !isComparable(typ) {
		val = Call(IdPath("DeepHash", hashPkg)).Args(val).Obj
	}
	return val
}

func passthroughFields(prefix string, fields *dst.FieldList) []dst.Expr {
	var exprs []dst.Expr
	for n, f := range fields.List {
		if len(f.Names) == 0 {
			exprs = append(exprs, Idf("%s%d", prefix, n+1))
		}

		for _, name := range f.Names {
			exprs = append(exprs, Id(name.Name))
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
				assigns = append(assigns, Assign(Id(rName)).
					Tok(token.ASSIGN).
					Rhs(Sel(Id(resultIdent)).
						Dot(Id(c.export(rName))).Obj).Obj)
			}

			for _, name := range result.Names {
				assigns = append(assigns, Assign(Id(name.Name)).
					Tok(token.ASSIGN).
					Rhs(Sel(Id(resultIdent)).
						Dot(Id(c.export(name.Name))).Obj).Obj)
			}
		}
	}
	return assigns
}

func cloneAndNameUnnamed(prefix string, fieldList *dst.FieldList) *dst.FieldList {
	fieldList = cloneNilableFieldList(fieldList)
	if fieldList != nil {
		for n, f := range fieldList.List {
			if len(f.Names) == 0 {
				f.Names = []*dst.Ident{Idf("%s%d", prefix, n+1)}
			}
		}
	}
	return fieldList
}

func (c *Converter) mockName(typeName string) string {
	return c.export(mockIdent + strings.Title(typeName))
}

func (c *Converter) export(name string) string {
	if c.isExported {
		name = strings.Title(name)
	}
	return name
}

func stdFuncDec() dst.FuncDeclDecorations {
	return dst.FuncDeclDecorations{
		NodeDecs: dst.NodeDecs{Before: dst.EmptyLine, After: dst.EmptyLine},
	}
}

func labelDirection(label string) (unnamedPrefix string, comparable bool) {
	switch label {
	case paramsIdent:
		unnamedPrefix = paramPrefix
		comparable = false
	case paramsKeyIdent:
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
	x.Sel = Id(sel)
	return x
}

func genDeclDec(format string, a ...interface{}) dst.GenDeclDecorations {
	return dst.GenDeclDecorations{
		NodeDecs: nodeDec(format, a...),
	}
}

func fnDeclDec(format string, a ...interface{}) dst.FuncDeclDecorations {
	return dst.FuncDeclDecorations{
		NodeDecs: nodeDec(format, a...),
	}
}

func nodeDec(format string, a ...interface{}) dst.NodeDecs {
	return dst.NodeDecs{
		Before: dst.NewLine,
		Start:  dst.Decorations{fmt.Sprintf(format, a...)},
	}
}

func litDec() dst.CompositeLitDecorations {
	return dst.CompositeLitDecorations{Lbrace: []string{"\n"}}
}

func kvExprDec(before dst.SpaceType) dst.KeyValueExprDecorations {
	return KeyValueDecs(before).After(dst.NewLine).Obj
}
