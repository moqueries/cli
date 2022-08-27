package generator

import (
	"fmt"
	"go/token"

	"github.com/dave/dst"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	. "github.com/myshkin5/moqueries/ast"
	"github.com/myshkin5/moqueries/logs"
)

const (
	moqueriesPkg  = "github.com/myshkin5/moqueries"
	hashPkg       = moqueriesPkg + "/hash"
	moqPkg        = moqueriesPkg + "/moq"
	syncAtomicPkg = "sync/atomic"

	intType           = "int"
	configType        = "Config"
	hashType          = "Hash"
	paramIndexingType = "ParamIndexing"
	repeaterType      = "Repeater"
	repeatValType     = "RepeatVal"
	sceneType         = "Scene"
	tType             = "T"

	anyCountIdent          = "anyCount"
	anyParamsIdent         = "anyParams"
	anyParamsReceiverIdent = "a"
	anyTimesIdent          = "AnyTimes"
	configIdent            = "config"
	doFnIdent              = "doFn"
	doReturnFnIdent        = "doReturnFn"
	expectationIdent       = "Expectation"
	hashesIdent            = "hashes"
	iIdent                 = "i"
	indexIdent             = "index"
	insertAtIdent          = "insertAt"
	lastIdent              = "last"
	minTimesIdent          = "MinTimes"
	missingIdent           = "missing"
	mockIdent              = "mock"
	moqIdent               = "moq"
	moqReceiverIdent       = "m"
	nIdent                 = "n"
	nilIdent               = "nil"
	okIdent                = "ok"
	parameterIndexingIdent = "parameterIndexing"
	paramIndexByValueIdent = "ParamIndexByValue"
	paramIndexByHashIdent  = "ParamIndexByHash"
	paramsIdent            = "params"
	paramsKeyIdent         = "paramsKey"
	recorderIdent          = "recorder"
	recorderReceiverIdent  = "r"
	repeatIdent            = "repeat"
	repeatersIdent         = "repeaters"
	resultsByParamsIdent   = "resultsByParams"
	resIdent               = "res"
	resultIdent            = "result"
	resultCountIdent       = "ResultCount"
	resultsIdent           = "results"
	runtimeIdent           = "runtime"
	sceneIdent             = "scene"
	sequenceIdent          = "sequence"
	strictIdent            = "Strict"
	valuesIdent            = "values"

	andDoFnName           = "andDo"
	assertFnName          = "AssertExpectationsMet"
	doReturnResultsFnName = "doReturnResults"
	errorfFnName          = "Errorf"
	fatalfFnName          = "Fatalf"
	findResultsFnName     = "findResults"
	fnFnName              = "fn"
	helperFnName          = "Helper"
	incrementFnName       = "Increment"
	lenFnName             = "len"
	mockFnName            = "mock"
	onCallFnName          = "onCall"
	paramsKeyFnName       = "paramsKey"
	prettyParamsFnName    = "prettyParams"
	repeatFnName          = "repeat"
	resetFnName           = "Reset"
	returnFnName          = "returnResults"

	fnRecorderSuffix = "fnRecorder"
	paramPrefix      = "param"
	resultPrefix     = "result"
	usedSuffix       = "Used"
	usedHashSuffix   = "UsedHash"

	sep     = "_"
	double  = "%s" + sep + "%s"
	triple  = "%s" + sep + "%s" + sep + "%s"
	unnamed = "%s%d"
)

// invalidNames is a list of names that if seen in field lists could cause a
// bad mock to be generated. When these names are seen in a parameter list or a
// result list, they are treated as though the name doesn't exist and a generic
// name is given (e.g.: param1).
var invalidNames = map[string]struct{}{
	moqReceiverIdent:      {},
	recorderReceiverIdent: {},
	sequenceIdent:         {},
	paramsIdent:           {},
	resultsIdent:          {},
	resultIdent:           {},
}

var titler = cases.Title(language.Und, cases.NoLower)

// Converter converts various interface and function types to AST structs to
// build a moq
type Converter struct {
	typ        Type
	isExported bool
	typeCache  TypeCache
}

// NewConverter creates a new Converter
func NewConverter(typ Type, isExported bool, typeCache TypeCache) *Converter {
	return &Converter{
		typ:        typ,
		isExported: isExported,
		typeCache:  typeCache,
	}
}

// Func holds on to function related data
type Func struct {
	Name    string
	Params  *dst.FieldList
	Results *dst.FieldList
}

// BaseStruct generates the base structure used to store the moq's state
func (c *Converter) BaseStruct() *dst.GenDecl {
	mName := c.moqName()
	moqName := fmt.Sprintf(double, mName, mockIdent)

	fields := []*dst.Field{
		Field(Star(c.idPath(sceneType, moqPkg))).Names(c.exportId(sceneIdent)).Obj,
		Field(c.idPath(configType, moqPkg)).Names(c.exportId(configIdent)).Obj,
		Field(Star(Id(moqName))).Names(c.exportId(moqIdent)).
			Decs(FieldDecs(dst.None, dst.EmptyLine).Obj).Obj,
	}

	for _, fn := range c.typ.Funcs {
		typePrefix := c.typePrefix(fn)
		fieldSuffix := ""
		if _, ok := c.typ.TypeSpec.Type.(*dst.InterfaceType); ok {
			fieldSuffix = sep + fn.Name
		}
		fields = append(fields,
			Field(SliceType(Idf(double, typePrefix, resultsByParamsIdent))).
				Names(c.exportId(resultsByParamsIdent+fieldSuffix)).Obj)
	}

	fields = append(fields, Field(c.runtimeStruct()).
		Names(c.exportId(runtimeIdent)).
		Decs(FieldDecs(dst.EmptyLine, dst.None).Obj).Obj)

	return TypeDecl(TypeSpec(mName).Type(Struct(fields...)).Obj).
		Decs(genDeclDec("// %s holds the state of a moq of the %s type",
			mName, c.typ.TypeSpec.Name.Name)).Obj
}

// IsolationStruct generates a struct used to isolate an interface for the moq
func (c *Converter) IsolationStruct(suffix string) *dst.GenDecl {
	mName := c.moqName()
	iName := fmt.Sprintf(double, mName, suffix)

	return TypeDecl(TypeSpec(iName).Type(Struct(Field(Star(Id(mName))).
		Names(c.exportId(moqIdent)).Obj)).Obj).
		Decs(genDeclDec("// %s isolates the %s interface of the %s type",
			iName, suffix, c.typ.TypeSpec.Name.Name)).Obj
}

// MethodStructs generates a structure for storing a set of parameters or
// a set of results for a method invocation of a mock
func (c *Converter) MethodStructs(fn Func) ([]dst.Decl, error) {
	prefix := c.typePrefix(fn)

	paramsStruct, err := c.paramsStructDecl(prefix, false, fn.Params)
	if err != nil {
		logs.Panic("Creating params struct should never generate errors", err)
	}
	paramsKeyStruct, err := c.paramsStructDecl(prefix, true, fn.Params)
	if err != nil {
		return nil, err
	}

	return []dst.Decl{
		paramsStruct,
		paramsKeyStruct,
		c.resultByParamsStruct(prefix),
		c.doFuncType(prefix, fn.Params),
		c.doReturnFuncType(prefix, fn),
		c.resultsStruct(prefix, fn.Results),
		c.fnRecorderStruct(prefix),
		c.anyParamsStruct(prefix),
	}, nil
}

// NewFunc generates a function for constructing a moq
func (c *Converter) NewFunc() *dst.FuncDecl {
	fnName := c.export("newMoq" + c.typ.TypeSpec.Name.Name)
	mName := c.moqName()
	moqName := fmt.Sprintf(double, mName, mockIdent)
	return Fn(fnName).
		Params(
			Field(Star(c.idPath(sceneType, moqPkg))).Names(Id(sceneIdent)).Obj,
			Field(Star(c.idPath(configType, moqPkg))).Names(Id(configIdent)).Obj,
		).
		Results(Field(Star(Id(mName))).Obj).
		Body(
			If(Bin(Id(configIdent)).Op(token.EQL).Y(Id(nilIdent)).Obj).Body(
				Assign(Id(configIdent)).Tok(token.ASSIGN).Rhs(Un(token.AND,
					Comp(c.idPath(configType, moqPkg)).Obj)).Obj).Obj,
			Assign(Id(moqReceiverIdent)).Tok(token.DEFINE).
				Rhs(Un(token.AND, Comp(Id(mName)).Elts(
					Key(c.exportId(sceneIdent)).
						Value(Id(sceneIdent)).Decs(kvExprDec(dst.None)).Obj,
					Key(c.exportId(configIdent)).
						Value(Star(Id(configIdent))).Decs(kvExprDec(dst.None)).Obj,
					Key(c.exportId(moqIdent)).
						Value(Un(token.AND, Comp(Id(moqName)).Obj)).Obj,
					Key(c.exportId(runtimeIdent)).
						Value(Comp(c.runtimeStruct()).
							Elts(c.runtimeValues()...).Obj).Decs(kvExprDec(dst.None)).
						Decs(kvExprDec(dst.EmptyLine)).Obj,
				).Decs(litDec()).Obj)).Obj,
			Assign(Sel(Sel(Id(moqReceiverIdent)).
				Dot(c.exportId(moqIdent)).Obj).
				Dot(c.exportId(moqIdent)).Obj).
				Tok(token.ASSIGN).
				Rhs(Id(moqReceiverIdent)).
				Decs(AssignDecs(dst.None).After(dst.EmptyLine).Obj).Obj,
			Expr(Call(Sel(Id(sceneIdent)).Dot(Id("AddMoq")).Obj).
				Args(Id(moqReceiverIdent)).Obj).Obj,
			Return(Id(moqReceiverIdent)),
		).
		Decs(fnDeclDec("// %s creates a new moq of the %s type",
			fnName, c.typ.TypeSpec.Name.Name)).Obj
}

// IsolationAccessor generates a function to access an isolation interface
func (c *Converter) IsolationAccessor(suffix, fnName string) *dst.FuncDecl {
	mName := c.moqName()
	iName := fmt.Sprintf(double, mName, suffix)

	var retVal dst.Expr
	retVal = Sel(Id(moqReceiverIdent)).Dot(c.exportId(moqIdent)).Obj
	if fnName != mockFnName {
		retVal = Un(token.AND, Comp(Id(iName)).Elts(
			Key(c.exportId(moqIdent)).Value(Id(moqReceiverIdent)).
				Decs(kvExprDec(dst.None)).Obj).Decs(litDec()).Obj,
		)
	}

	fnName = c.export(fnName)
	return Fn(fnName).
		Recv(Field(Star(Id(mName))).Names(Id(moqReceiverIdent)).Obj).
		Results(Field(Star(Id(iName))).Obj).
		Body(Return(retVal)).
		Decs(fnDeclDec("// %s returns the %s implementation of the %s type",
			fnName, suffix, c.typ.TypeSpec.Name.Name)).Obj
}

// FuncClosure generates a mock implementation of function type wrapped in a
// closure
func (c *Converter) FuncClosure(fn Func) *dst.FuncDecl {
	mName := c.moqName()
	fnLitCall := Call(Sel(Id(moqIdent)).Dot(c.exportId(fnFnName)).Obj).
		Args(passthroughFields(paramPrefix, fn.Params)...).
		Ellipsis(isVariadic(fn.Params)).Obj
	var fnLitRetStmt dst.Stmt
	fnLitRetStmt = Return(fnLitCall)
	if fn.Results == nil {
		fnLitRetStmt = Expr(fnLitCall).Obj
	}

	return Fn(c.export(mockFnName)).
		Recv(Field(Star(Id(mName))).Names(Id(moqReceiverIdent)).Obj).
		Results(Field(c.idPath(c.typ.TypeSpec.Name.Name, c.typ.TypeSpec.Name.Path)).Obj).
		Body(Return(FnLit(FnType(cloneAndNameUnnamed(paramPrefix, fn.Params)).
			Results(cloneFieldList(fn.Results, true)).Obj).
			Body(Assign(Id(moqIdent)).
				Tok(token.DEFINE).
				Rhs(Un(
					token.AND,
					Comp(Idf(double, mName, mockIdent)).
						Elts(Key(c.exportId(moqIdent)).
							Value(Id(moqReceiverIdent)).Obj).Obj,
				)).Obj,
				fnLitRetStmt,
			).Obj)).
		Decs(fnDeclDec("// %s returns the %s implementation of the %s type",
			c.export(mockFnName), moqIdent, c.typ.TypeSpec.Name.Name)).Obj
}

// MockMethod generates a mock implementation of a method
func (c *Converter) MockMethod(fn Func) *dst.FuncDecl {
	mName := c.moqName()
	recv := fmt.Sprintf(double, mName, mockIdent)

	fnName := fn.Name
	fieldSuffix := sep + fn.Name
	typePrefix := c.typePrefix(fn)
	if fnName == "" {
		fnName = c.export(fnFnName)
		fieldSuffix = ""
	}

	return Fn(fnName).
		Recv(Field(Star(Id(recv))).Names(Id(moqReceiverIdent)).Obj).
		ParamList(cloneAndNameUnnamed(paramPrefix, fn.Params)).
		ResultList(cloneAndNameUnnamed(resultPrefix, fn.Results)).
		Body(c.mockFunc(typePrefix, fieldSuffix, fn)...).
		Decs(stdFuncDec()).Obj
}

// RecorderMethods generates a recorder implementation of a method and
// associated return method
func (c *Converter) RecorderMethods(fn Func) []dst.Decl {
	decls := []dst.Decl{
		c.recorderFn(fn),
	}

	decls = append(decls, c.anyParamFns(fn)...)
	decls = append(decls, c.recorderSeqFns(fn)...)
	decls = append(decls,
		c.returnResultsFn(fn),
		c.andDoFn(fn),
		c.doReturnResultsFn(fn),
		c.findResultsFn(fn),
		c.recorderRepeatFn(fn),
		c.prettyParamsFn(fn),
		c.paramsKeyFn(fn),
	)

	return decls
}

// ResetMethod generates a method to reset the moq's state
func (c *Converter) ResetMethod() *dst.FuncDecl {
	mName := c.moqName()

	var stmts []dst.Stmt
	for _, fn := range c.typ.Funcs {
		fieldSuffix := ""
		if _, ok := c.typ.TypeSpec.Type.(*dst.InterfaceType); ok {
			fieldSuffix = sep + fn.Name
		}

		stmts = append(stmts, Assign(Sel(Id(moqReceiverIdent)).
			Dot(c.exportId(resultsByParamsIdent+fieldSuffix)).Obj).
			Tok(token.ASSIGN).
			Rhs(Id(nilIdent)).Obj)
	}

	return Fn(resetFnName).
		Recv(Field(Star(Id(mName))).Names(Id(moqReceiverIdent)).Obj).
		Body(stmts...).
		Decs(fnDeclDec("// %s resets the state of the moq", resetFnName)).Obj
}

// AssertMethod generates a method to assert all expectations are met
func (c *Converter) AssertMethod() *dst.FuncDecl {
	mName := c.moqName()

	stmts := []dst.Stmt{
		c.helperCallExpr(Id(moqReceiverIdent)),
	}
	for _, fn := range c.typ.Funcs {
		fieldSuffix := ""
		if _, ok := c.typ.TypeSpec.Type.(*dst.InterfaceType); ok {
			fieldSuffix = sep + fn.Name
		}

		stmts = append(stmts, Range(Sel(Id(moqReceiverIdent)).
			Dot(c.exportId(resultsByParamsIdent+fieldSuffix)).Obj).
			Key(Id(sep)).Value(Id(resIdent)).Tok(token.DEFINE).
			Body(Range(Sel(Id(resIdent)).Dot(c.exportId(resultsIdent)).Obj).
				Key(Id(sep)).Value(Id(resultsIdent)).Tok(token.DEFINE).
				Body(
					Assign(Id(missingIdent)).
						Tok(token.DEFINE).
						Rhs(Bin(Sel(Sel(Id(resultsIdent)).
							Dot(c.exportId(repeatIdent)).Obj).
							Dot(Id(minTimesIdent)).Obj).
							Op(token.SUB).
							Y(Call(Id(intType)).Args(
								Call(c.idPath("LoadUint32", syncAtomicPkg)).Args(Un(
									token.AND,
									Sel(Id(resultsIdent)).
										Dot(c.exportId(indexIdent)).Obj)).Obj).Obj).Obj).Obj,
					If(Bin(Id(missingIdent)).Op(token.GTR).
						Y(LitInt(0)).Obj).
						Body(
							Expr(Call(Sel(Sel(Sel(Id(moqReceiverIdent)).
								Dot(c.exportId(sceneIdent)).Obj).
								Dot(Id(tType)).Obj).Dot(Id(errorfFnName)).Obj).
								Args(
									LitString("Expected %d additional call(s) to %s"),
									Id(missingIdent),
									c.callPrettyParams(fn,
										Id(moqReceiverIdent),
										Id(resultsIdent))).Obj).Obj,
						).Obj,
				).Obj,
			).Obj)
	}

	return Fn(assertFnName).
		Recv(Field(Star(Id(mName))).Names(Id(moqReceiverIdent)).Obj).
		Body(stmts...).
		Decs(fnDeclDec("// %s asserts that all expectations have been met",
			assertFnName)).Obj
}

func (c *Converter) typePrefix(fn Func) string {
	mName := c.moqName()
	typePrefix := fmt.Sprintf(double, mName, fn.Name)
	if fn.Name == "" {
		typePrefix = mName
	}
	return typePrefix
}

func (c *Converter) runtimeStruct() *dst.StructType {
	return Struct(Field(c.paramIndexingStruct()).
		Names(c.exportId(parameterIndexingIdent)).Obj)
}

func (c *Converter) paramIndexingStruct() *dst.StructType {
	var piFields []*dst.Field
	for _, fn := range c.typ.Funcs {
		if fn.Name == "" {
			piFields = append(piFields, c.paramIndexingFnStruct(fn).Fields.List...)
		} else {
			piFields = append(piFields,
				Field(c.paramIndexingFnStruct(fn)).Names(Id(fn.Name)).Obj)
		}
	}

	return Struct(piFields...)
}

func (c *Converter) paramIndexingFnStruct(fn Func) *dst.StructType {
	var piParamFields []*dst.Field
	count := 0
	for _, f := range fn.Params.List {
		if len(f.Names) == 0 {
			piParamFields = append(piParamFields,
				c.paramIndexingField(fmt.Sprintf(unnamed, paramPrefix, count+1)))
			count++
		}

		for _, name := range f.Names {
			piParamFields = append(piParamFields,
				c.paramIndexingField(validName(name.Name, paramPrefix, count)))
			count++
		}
	}

	return Struct(piParamFields...)
}

func (c *Converter) paramIndexingField(name string) *dst.Field {
	return Field(c.idPath(paramIndexingType, moqPkg)).Names(c.exportId(name)).Obj
}

func (c *Converter) runtimeValues() []dst.Expr {
	var vals []dst.Expr
	kvDec := dst.NewLine
	for _, fn := range c.typ.Funcs {
		if fn.Name == "" {
			vals = append(vals, c.paramIndexingFnValues(fn.Params.List)...)
		} else {
			vals = append(vals, Key(Id(fn.Name)).
				Value(Comp(c.paramIndexingFnStruct(fn)).
					Elts(c.paramIndexingFnValues(fn.Params.List)...).Obj).Decs(kvExprDec(kvDec)).Obj)
		}
		kvDec = dst.None
	}

	return []dst.Expr{Key(c.exportId(parameterIndexingIdent)).
		Value(Comp(c.paramIndexingStruct()).Elts(vals...).Obj).Obj}
}

func (c *Converter) paramIndexingFnValues(params []*dst.Field) []dst.Expr {
	var vals []dst.Expr
	kvDec := dst.NewLine
	count := 0
	for _, f := range params {
		if len(f.Names) == 0 {
			val := c.paramIndexingValue(f.Type,
				fmt.Sprintf(unnamed, paramPrefix, count+1), kvDec)

			vals = append(vals, val)
			count++
			kvDec = dst.None
		}

		for _, name := range f.Names {
			val := c.paramIndexingValue(f.Type,
				validName(name.Name, paramPrefix, count), kvDec)

			vals = append(vals, val)
			count++
			kvDec = dst.None
		}
	}

	return vals
}

func (c *Converter) paramIndexingValue(typ dst.Expr, name string, kvDec dst.SpaceType) *dst.KeyValueExpr {
	comp, err := c.typeCache.IsDefaultComparable(typ)
	if err != nil {
		logs.Panic("Call MethodStructs first to get a meaningful error", err)
	}

	val := paramIndexByValueIdent
	if !comp {
		val = paramIndexByHashIdent
	}

	return Key(c.exportId(name)).Value(c.idPath(val, moqPkg)).Decs(kvExprDec(kvDec)).Obj
}

func (c *Converter) paramsStructDecl(
	prefix string, paramsKey bool, fieldList *dst.FieldList,
) (*dst.GenDecl, error) {
	var mStruct *dst.StructType
	var label, goDocDesc string
	if paramsKey {
		label = paramsKeyIdent
		goDocDesc = "map key params"

		pStruct, err := c.methodStruct(paramsKeyIdent, fieldList)
		if err != nil {
			return nil, err
		}
		pkStruct, err := c.methodStruct(hashesIdent, fieldList)
		if err != nil {
			return nil, err
		}

		mStruct = Struct(Field(pStruct).Names(c.exportId(paramsIdent)).Obj,
			Field(pkStruct).Names(c.exportId(hashesIdent)).Obj)
	} else {
		label = paramsIdent
		goDocDesc = label

		var err error
		mStruct, err = c.methodStruct(paramsIdent, fieldList)
		if err != nil {
			return nil, err
		}
	}

	structName := fmt.Sprintf(double, prefix, label)
	return TypeDecl(TypeSpec(structName).Type(mStruct).Obj).
		Decs(genDeclDec("// %s holds the %s of the %s type",
			structName, goDocDesc, c.typ.TypeSpec.Name.Name)).Obj, nil
}

func (c *Converter) methodStruct(label string, fieldList *dst.FieldList) (*dst.StructType, error) {
	unnamedPrefix, _ := labelDirection(label)
	fieldList = cloneFieldList(fieldList, false)

	if fieldList == nil {
		return StructFromList(nil), nil
	}

	count := 0
	var fList []*dst.Field
	for _, f := range fieldList.List {
		if len(f.Names) == 0 {
			f.Names = []*dst.Ident{Idf(unnamed, unnamedPrefix, count+1)}
		}

		for n, name := range f.Names {
			f.Names[n] = Id(c.export(validName(name.Name, unnamedPrefix, count)))
			count++
		}

		typ, err := c.comparableType(label, f.Type)
		if err != nil {
			return nil, err
		}
		if typ != nil {
			f.Type = typ
			fList = append(fList, f)
		}
	}

	if len(fList) != 0 {
		fieldList.List = fList
	} else {
		fieldList = nil
	}
	return StructFromList(fieldList), nil
}

func (c *Converter) comparableType(label string, typ dst.Expr) (dst.Expr, error) {
	switch label {
	case paramsIdent:
	case resultsIdent:
	case paramsKeyIdent:
		comp, err := c.typeCache.IsComparable(typ)
		if err != nil {
			return nil, err
		}
		if !comp {
			// Non-comparable params are not represented in the params section
			// of the paramsKey
			return nil, nil
		}
	case hashesIdent:
		// Everything is represented as a hash in the hashes section of the
		// paramsKey
		return c.idPath(hashType, hashPkg), nil
	default:
		logs.Panicf("Unknown label: %s", label)
	}

	if ellipsis, ok := typ.(*dst.Ellipsis); ok {
		// Ellipsis params are represented as a slice (when not comparable)
		return SliceType(ellipsis.Elt), nil
	}

	return typ, nil
}

func (c *Converter) resultByParamsStruct(prefix string) *dst.GenDecl {
	structName := fmt.Sprintf(double, prefix, resultsByParamsIdent)

	return TypeDecl(TypeSpec(structName).Type(Struct(
		Field(Id(intType)).Names(c.exportId(anyCountIdent)).Obj,
		Field(Id("uint64")).Names(c.exportId(anyParamsIdent)).Obj,
		Field(MapType(Id(c.export(fmt.Sprintf(double, prefix, paramsKeyIdent)))).
			Value(Star(Idf(double, prefix, resultsIdent))).Obj,
		).Names(c.exportId(resultsIdent)).Obj,
	)).Obj).Decs(genDeclDec(
		"// %s contains the results for a given set of parameters for the %s type",
		structName,
		c.typ.TypeSpec.Name.Name)).Obj
}

func (c *Converter) doFuncType(prefix string, params *dst.FieldList) *dst.GenDecl {
	fnName := fmt.Sprintf(double, prefix, doFnIdent)
	return TypeDecl(TypeSpec(fnName).
		Type(FnType(cloneFieldList(params, false)).Obj).Obj).
		Decs(genDeclDec(
			"// %s defines the type of function needed when calling %s for the %s type",
			fnName,
			c.export(andDoFnName),
			c.typ.TypeSpec.Name.Name)).Obj
}

func (c *Converter) doReturnFuncType(prefix string, fn Func) *dst.GenDecl {
	fnName := fmt.Sprintf(double, prefix, doReturnFnIdent)
	return TypeDecl(TypeSpec(fnName).
		Type(FnType(cloneFieldList(fn.Params, false)).
			Results(cloneFieldList(fn.Results, false)).Obj).Obj).
		Decs(genDeclDec(
			"// %s defines the type of function needed when calling %s for the %s type",
			fnName,
			c.export(doReturnResultsFnName),
			c.typ.TypeSpec.Name.Name)).Obj
}

func (c *Converter) resultsStruct(prefix string, results *dst.FieldList) *dst.GenDecl {
	structName := fmt.Sprintf(double, prefix, resultsIdent)

	return TypeDecl(TypeSpec(structName).Type(Struct(
		Field(Idf(double, prefix, paramsIdent)).
			Names(c.exportId(paramsIdent)).Obj,
		Field(SliceType(c.innerResultsStruct(prefix, results))).
			Names(c.exportId(resultsIdent)).Obj,
		Field(Id("uint32")).Names(c.exportId(indexIdent)).Obj,
		Field(Star(c.idPath(repeatValType, moqPkg))).Names(c.exportId(repeatIdent)).Obj,
	)).Obj).Decs(genDeclDec("// %s holds the results of the %s type",
		structName,
		c.typ.TypeSpec.Name.Name)).Obj
}

func (c *Converter) innerResultsStruct(prefix string, results *dst.FieldList) *dst.StructType {
	mStruct, err := c.methodStruct(resultsIdent, results)
	if err != nil {
		logs.Panic("Creating results struct should never generate errors", err)
	}

	return Struct(
		Field(Star(mStruct)).Names(c.exportId(valuesIdent)).Obj,
		Field(Id("uint32")).Names(c.exportId(sequenceIdent)).Obj,
		Field(Idf(double, prefix, doFnIdent)).Names(c.exportId(doFnIdent)).Obj,
		Field(Idf(double, prefix, doReturnFnIdent)).Names(c.exportId(doReturnFnIdent)).Obj,
	)
}

func (c *Converter) fnRecorderStruct(prefix string) *dst.GenDecl {
	mName := c.moqName()
	structName := fmt.Sprintf(double, prefix, fnRecorderSuffix)
	return TypeDecl(TypeSpec(structName).Type(Struct(
		Field(Idf(double, prefix, paramsIdent)).
			Names(c.exportId(paramsIdent)).Obj,
		Field(Id("uint64")).
			Names(c.exportId(anyParamsIdent)).Obj,
		Field(Id("bool")).
			Names(c.exportId(sequenceIdent)).Obj,
		Field(Star(Idf(double, prefix, resultsIdent))).
			Names(c.exportId(resultsIdent)).Obj,
		Field(Star(Id(mName))).
			Names(c.exportId(moqIdent)).Obj,
	)).Obj).Decs(genDeclDec("// %s routes recorded function calls to the %s moq",
		structName, mName)).Obj
}

func (c *Converter) anyParamsStruct(prefix string) *dst.GenDecl {
	structName := fmt.Sprintf(double, prefix, anyParamsIdent)
	recStructName := fmt.Sprintf(double, prefix, fnRecorderSuffix)
	return TypeDecl(TypeSpec(structName).Type(Struct(Field(Star(Id(recStructName))).
		Names(c.exportId(recorderIdent)).Obj)).Obj).
		Decs(genDeclDec("// %s isolates the any params functions of the %s type",
			structName, c.typ.TypeSpec.Name.Name)).Obj
}

func (c *Converter) mockFunc(typePrefix, fieldSuffix string, fn Func) []dst.Stmt {
	stateSelector := Sel(Id(moqReceiverIdent)).Dot(c.exportId(moqIdent)).Obj

	paramsKeyFn := c.export(paramsKeyFnName)
	if fn.Name != "" {
		paramsKeyFn = c.export(fmt.Sprintf(double, paramsKeyFnName, fn.Name))
	}
	moqSel := Sel(Id(moqReceiverIdent)).Dot(c.exportId(moqIdent)).Obj

	stmts := []dst.Stmt{
		c.helperCallExpr(Sel(Id(moqReceiverIdent)).Dot(c.exportId(moqIdent)).Obj),
		Assign(Id(paramsIdent)).
			Tok(token.DEFINE).
			Rhs(Comp(Idf(double, typePrefix, paramsIdent)).
				Elts(c.passthroughElements(fn.Params, paramsIdent, "")...).Obj).Obj,
		Var(Value(Star(Idf(double, typePrefix, resultsIdent))).
			Names(Id(resultsIdent)).Obj),
		Range(Sel(cloneExpr(stateSelector)).
			Dot(c.exportId(resultsByParamsIdent+fieldSuffix)).Obj).
			Key(Id(sep)).
			Value(Id(resultsByParamsIdent)).
			Tok(token.DEFINE).
			Body(
				Assign(Id(paramsKeyIdent)).
					Tok(token.DEFINE).
					Rhs(Call(Sel(Sel(Id(moqReceiverIdent)).
						Dot(c.exportId(moqIdent)).Obj).
						Dot(Id(paramsKeyFn)).Obj).
						Args(Id(paramsIdent), Sel(Id(resultsByParamsIdent)).
							Dot(c.exportId(anyParamsIdent)).Obj).Obj).Obj,
				Var(Value(Id("bool")).Names(Id(okIdent)).Obj),
				Assign(Id(resultsIdent), Id(okIdent)).
					Tok(token.ASSIGN).
					Rhs(Index(Sel(Id(resultsByParamsIdent)).
						Dot(c.exportId(resultsIdent)).Obj).
						Sub(Id(paramsKeyIdent)).Obj).Obj,
				If(Id(okIdent)).Body(Break()).Obj,
			).Obj,
	}

	stmts = append(stmts,
		If(Bin(Id(resultsIdent)).Op(token.EQL).Y(Id(nilIdent)).Obj).Body(
			If(Bin(Sel(Sel(cloneExpr(stateSelector)).
				Dot(c.exportId(configIdent)).Obj).
				Dot(Id(expectationIdent)).Obj).
				Op(token.EQL).
				Y(c.idPath(strictIdent, moqPkg)).Obj).
				Body(
					Expr(Call(Sel(Sel(Sel(cloneExpr(stateSelector)).
						Dot(c.exportId(sceneIdent)).Obj).
						Dot(Id(tType)).Obj).
						Dot(Id(fatalfFnName)).Obj).
						Args(LitString("Unexpected call to %s"),
							c.callPrettyParams(fn, moqSel, nil)).Obj).Obj).Obj,
			Return(),
		).Obj)

	stmts = append(stmts, Assign(Id(iIdent)).
		Tok(token.DEFINE).
		Rhs(Bin(Call(Id(intType)).
			Args(Call(c.idPath("AddUint32", syncAtomicPkg)).Args(Un(
				token.AND,
				Sel(Id(resultsIdent)).Dot(c.exportId(indexIdent)).Obj),
				LitInt(1)).Obj).Obj).
			Op(token.SUB).
			Y(LitInt(1)).Obj).
		Decs(AssignDecs(dst.EmptyLine).Obj).Obj)
	stmts = append(stmts,
		If(Bin(Id(iIdent)).Op(token.GEQ).Y(
			Sel(Sel(Id(resultsIdent)).Dot(c.exportId(repeatIdent)).Obj).
				Dot(Id(resultCountIdent)).Obj).Obj).
			Body(
				If(Un(token.NOT, Sel(Sel(Id(resultsIdent)).
					Dot(c.exportId(repeatIdent)).Obj).Dot(Id(anyTimesIdent)).Obj)).
					Body(
						If(Bin(Sel(Sel(cloneExpr(stateSelector)).
							Dot(c.exportId(configIdent)).Obj).
							Dot(Id(expectationIdent)).Obj).
							Op(token.EQL).
							Y(c.idPath(strictIdent, moqPkg)).Obj).
							Body(Expr(Call(Sel(Sel(Sel(cloneExpr(stateSelector)).
								Dot(c.exportId(sceneIdent)).Obj).
								Dot(Id(tType)).Obj).
								Dot(Id(fatalfFnName)).Obj).
								Args(
									LitString("Too many calls to %s"),
									c.callPrettyParams(fn, moqSel, nil),
								).Obj).Obj).Obj,
						Return(),
					).Obj,
				Assign(Id(iIdent)).
					Tok(token.ASSIGN).
					Rhs(Bin(Sel(Sel(Id(resultsIdent)).
						Dot(c.exportId(repeatIdent)).Obj).Dot(Id(resultCountIdent)).Obj).
						Op(token.SUB).
						Y(LitInt(1)).Obj).Obj,
			).Decs(IfDecs(dst.EmptyLine).Obj).Obj)

	stmts = append(stmts, Assign(Id(resultIdent)).
		Tok(token.DEFINE).
		Rhs(Index(Sel(Id(resultsIdent)).Dot(c.exportId(resultsIdent)).Obj).
			Sub(Id(iIdent)).Obj).Obj)
	stmts = append(stmts, If(Bin(
		Sel(Id(resultIdent)).Dot(c.exportId(sequenceIdent)).Obj).
		Op(token.NEQ).
		Y(LitInt(0)).Obj).Body(
		Assign(Id(sequenceIdent)).Tok(token.DEFINE).Rhs(Call(
			Sel(Sel(Sel(Id(moqReceiverIdent)).
				Dot(c.exportId(moqIdent)).Obj).Dot(c.exportId(sceneIdent)).Obj).
				Dot(Id("NextMockSequence")).Obj).Obj).Obj,
		If(Bin(Paren(Bin(Un(token.NOT, Sel(Sel(Id(resultsIdent)).
			Dot(c.exportId(repeatIdent)).Obj).Dot(Id(anyTimesIdent)).Obj)).Op(token.LAND).
			Y(Bin(Sel(Id(resultIdent)).Dot(c.exportId(sequenceIdent)).Obj).
				Op(token.NEQ).Y(Id(sequenceIdent)).Obj).Obj)).Op(token.LOR).
			Y(Bin(Sel(Id(resultIdent)).Dot(c.exportId(sequenceIdent)).Obj).
				Op(token.GTR).Y(Id(sequenceIdent)).Obj).Obj).Body(
			Expr(Call(Sel(Sel(Sel(cloneExpr(stateSelector)).
				Dot(c.exportId(sceneIdent)).Obj).
				Dot(Id(tType)).Obj).
				Dot(Id(fatalfFnName)).Obj).
				Args(LitString("Call sequence does not match call to %s"),
					c.callPrettyParams(fn, moqSel, nil)).Obj).Obj,
		).Obj,
	).Decs(IfDecs(dst.EmptyLine).Obj).Obj)

	variadic := isVariadic(fn.Params)
	stmts = append(stmts, If(Bin(Sel(Id(resultIdent)).
		Dot(c.exportId(doFnIdent)).Obj).Op(token.NEQ).Y(Id(nilIdent)).Obj).Body(
		Expr(Call(Sel(Id(resultIdent)).Dot(c.exportId(doFnIdent)).Obj).
			Args(passthroughFields(paramPrefix, fn.Params)...).Ellipsis(variadic).Obj).Obj,
	).Decs(IfDecs(dst.EmptyLine).Obj).Obj)

	doReturnCall := Call(Sel(Id(resultIdent)).Dot(c.exportId(doReturnFnIdent)).Obj).
		Args(passthroughFields(paramPrefix, fn.Params)...).Ellipsis(variadic).Obj
	var doReturnStmt dst.Stmt = Expr(doReturnCall).Obj
	if fn.Results != nil {
		stmts = append(stmts, If(Bin(Sel(Id(resultIdent)).
			Dot(c.exportId(valuesIdent)).Obj).Op(token.NEQ).Y(Id(nilIdent)).Obj).Body(
			c.assignResult(fn.Results)...).Obj)
		doReturnStmt = Assign(passthroughFields(resultPrefix, fn.Results)...).
			Tok(token.ASSIGN).Rhs(doReturnCall).Obj
	}

	stmts = append(stmts, If(Bin(Sel(Id(resultIdent)).
		Dot(c.exportId(doReturnFnIdent)).Obj).Op(token.NEQ).Y(Id(nilIdent)).Obj).
		Body(doReturnStmt).Obj)
	stmts = append(stmts, Return())

	return stmts
}

func (c *Converter) recorderFn(fn Func) *dst.FuncDecl {
	mName := c.moqName()

	recvType := fmt.Sprintf(double, mName, recorderIdent)
	fnName := fn.Name
	fnRecName := fmt.Sprintf(triple, mName, fn.Name, fnRecorderSuffix)
	typePrefix := c.typePrefix(fn)
	var moqVal dst.Expr = Sel(Id(moqReceiverIdent)).
		Dot(c.exportId(moqIdent)).Obj
	if fn.Name == "" {
		recvType = mName
		fnName = c.export(onCallFnName)
		fnRecName = fmt.Sprintf(double, mName, fnRecorderSuffix)
		moqVal = Id(moqReceiverIdent)
	}

	return Fn(fnName).
		Recv(Field(Star(Id(recvType))).Names(Id(moqReceiverIdent)).Obj).
		ParamList(cloneAndNameUnnamed(paramPrefix, fn.Params)).
		Results(Field(Star(Id(fnRecName))).Obj).
		Body(c.recorderFnInterfaceBody(fnRecName, typePrefix, moqVal, fn)...).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) recorderFnInterfaceBody(
	fnRecName, typePrefix string, moqVal dst.Expr, fn Func,
) []dst.Stmt {
	return []dst.Stmt{Return(Un(
		token.AND,
		Comp(Id(fnRecName)).
			Elts(
				Key(c.exportId(paramsIdent)).
					Value(Comp(Idf(double, typePrefix, paramsIdent)).
						Elts(c.passthroughElements(fn.Params, paramsIdent, "")...).Obj,
					).Decs(kvExprDec(dst.None)).Obj,
				Key(c.exportId(sequenceIdent)).
					Value(Bin(Sel(Sel(moqVal).
						Dot(c.exportId(configIdent)).Obj).
						Dot(Id(titler.String(sequenceIdent))).Obj).
						Op(token.EQL).
						Y(c.idPath("SeqDefaultOn", moqPkg)).Obj).
					Decs(kvExprDec(dst.None)).Obj,
				Key(c.exportId(moqIdent)).
					Value(cloneExpr(moqVal)).Decs(kvExprDec(dst.None)).Obj,
			).Decs(litDec()).Obj,
	))}
}

func (c *Converter) anyParamFns(fn Func) []dst.Decl {
	mName := c.moqName()

	fnRecName := fmt.Sprintf(triple, mName, fn.Name, fnRecorderSuffix)
	anyParamsName := fmt.Sprintf(triple, mName, fn.Name, anyParamsIdent)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf(double, mName, fnRecorderSuffix)
		anyParamsName = fmt.Sprintf(double, mName, anyParamsIdent)
	}

	decls := []dst.Decl{c.anyParamAnyFn(fn, anyParamsName, fnRecName)}
	count := 0
	for _, param := range fn.Params.List {
		if len(param.Names) == 0 {
			pName := fmt.Sprintf(unnamed, paramPrefix, count+1)
			decls = append(decls, c.anyParamFn(anyParamsName, fnRecName, pName, count))
			count++
		}

		for _, name := range param.Names {
			decls = append(decls, c.anyParamFn(anyParamsName, fnRecName,
				validName(name.Name, paramPrefix, count), count))
			count++
		}
	}
	return decls
}

func (c *Converter) anyParamAnyFn(fn Func, anyParamsName, fnRecName string) *dst.FuncDecl {
	moqSel := Sel(Id(recorderReceiverIdent)).Dot(c.exportId(moqIdent)).Obj

	return Fn(c.export("any")).
		Recv(Field(Star(Id(fnRecName))).Names(Id(recorderReceiverIdent)).Obj).
		Results(Field(Star(Id(anyParamsName))).Obj).
		Body(
			c.helperCallExpr(Sel(Id(recorderReceiverIdent)).Dot(c.exportId(moqIdent)).Obj),
			If(Bin(Sel(Id(recorderReceiverIdent)).Dot(c.exportId(resultsIdent)).Obj).
				Op(token.NEQ).
				Y(Id(nilIdent)).Obj).
				Body(
					Expr(Call(Sel(Sel(c.selExport(moqSel, sceneIdent)).
						Dot(Id(tType)).Obj).
						Dot(Id(fatalfFnName)).Obj).
						Args(LitStringf(
							"Any functions must be called before %s or %s calls, recording %%s",
							c.export(returnFnName), c.export(doReturnResultsFnName)),
							c.callPrettyParams(fn, moqSel, Id(recorderReceiverIdent)),
						).Obj).Obj,
					Return(Id(nilIdent))).Obj,
			Return(Un(
				token.AND,
				Comp(Id(anyParamsName)).Elts(Key(c.exportId(recorderIdent)).
					Value(Id(recorderReceiverIdent)).Obj).Obj)),
		).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) anyParamFn(anyParamsName, fnRecName, pName string, paramPos int) *dst.FuncDecl {
	return Fn(c.export(pName)).
		Recv(Field(Star(Id(anyParamsName))).Names(Id(anyParamsReceiverIdent)).Obj).
		Results(Field(Star(Id(fnRecName))).Obj).
		Body(
			Assign(Sel(Sel(Id(anyParamsReceiverIdent)).
				Dot(c.exportId(recorderIdent)).Obj).
				Dot(c.exportId(anyParamsIdent)).Obj).
				Tok(token.OR_ASSIGN).
				Rhs(Bin(LitInt(1)).Op(token.SHL).Y(LitInt(paramPos)).Obj).Obj,
			Return(Sel(Id(anyParamsReceiverIdent)).Dot(c.exportId(recorderIdent)).Obj),
		).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) returnResultsFn(fn Func) *dst.FuncDecl {
	mStruct, err := c.methodStruct(resultsIdent, fn.Results)
	if err != nil {
		logs.Panic("Creating results struct should never generate errors", err)
	}

	params := cloneAndNameUnnamed(resultPrefix, fn.Results)
	resExprs := []dst.Expr{
		Key(c.exportId(valuesIdent)).
			Value(Un(token.AND, Comp(mStruct).
				Elts(c.passthroughElements(fn.Results, resultsIdent, "")...).Obj)).
			Decs(kvExprDec(dst.NewLine)).Obj,
		Key(c.exportId(sequenceIdent)).
			Value(Id(sequenceIdent)).Decs(kvExprDec(dst.None)).Obj,
	}

	return c.returnFn(returnFnName, fn, params, resExprs)
}

func (c *Converter) returnFn(
	fnName string,
	fn Func,
	params *dst.FieldList,
	resultExprs []dst.Expr,
) *dst.FuncDecl {
	mName := c.moqName()

	fnRecName := fmt.Sprintf(triple, mName, fn.Name, fnRecorderSuffix)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf(double, mName, fnRecorderSuffix)
	}

	resStruct := c.innerResultsStruct(c.typePrefix(fn), fn.Results)

	return Fn(c.export(fnName)).
		Recv(Field(Star(Id(fnRecName))).Names(Id(recorderReceiverIdent)).Obj).
		ParamList(params).
		Results(Field(Star(Id(fnRecName))).Obj).
		Body(
			c.helperCallExpr(Sel(Id(recorderReceiverIdent)).Dot(c.exportId(moqIdent)).Obj),
			Expr(Call(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(findResultsFnName)).Obj).Obj).
				Decs(ExprDecs(dst.EmptyLine).Obj).Obj,
			Var(Value(Id("uint32")).Names(Id(sequenceIdent)).Obj),
			If(Sel(Id(recorderReceiverIdent)).Dot(c.exportId(sequenceIdent)).Obj).
				Body(Assign(Id(sequenceIdent)).Tok(token.ASSIGN).Rhs(
					Call(Sel(Sel(Sel(Id(recorderReceiverIdent)).
						Dot(c.exportId(moqIdent)).Obj).
						Dot(c.exportId(sceneIdent)).Obj).
						Dot(Id("NextRecorderSequence")).Obj).Obj,
				).Obj).
				Decs(IfDecs(dst.EmptyLine).Obj).Obj,
			Assign(
				Sel(Sel(Id(recorderReceiverIdent)).
					Dot(c.exportId(resultsIdent)).Obj).
					Dot(c.exportId(resultsIdent)).Obj).
				Tok(token.ASSIGN).
				Rhs(Call(Id("append")).Args(Sel(Sel(Id(recorderReceiverIdent)).
					Dot(c.exportId(resultsIdent)).Obj).
					Dot(c.exportId(resultsIdent)).Obj,
					Comp(resStruct).Elts(resultExprs...).Obj).Obj).Obj,
			Return(Id(recorderReceiverIdent)),
		).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) andDoFn(fn Func) *dst.FuncDecl {
	mName := c.moqName()

	fnRecName := fmt.Sprintf(triple, mName, fn.Name, fnRecorderSuffix)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf(double, mName, fnRecorderSuffix)
	}
	typePrefix := c.typePrefix(fn)
	fnName := fmt.Sprintf(double, typePrefix, doFnIdent)

	return Fn(c.export(andDoFnName)).
		Recv(Field(Star(Id(fnRecName))).Names(Id(recorderReceiverIdent)).Obj).
		Params(Field(Id(fnName)).Names(Id(fnFnName)).Obj).
		Results(Field(Star(Id(fnRecName))).Obj).Body(
		c.helperCallExpr(Sel(Id(recorderReceiverIdent)).Dot(c.exportId(moqIdent)).Obj),
		If(Bin(Sel(Id(recorderReceiverIdent)).
			Dot(c.exportId(resultsIdent)).Obj).
			Op(token.EQL).
			Y(Id(nilIdent)).Obj).
			Body(Expr(Call(Sel(Sel(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(moqIdent)).Obj).
				Dot(c.exportId(sceneIdent)).Obj).
				Dot(Id(tType)).Obj).Dot(Id(fatalfFnName)).Obj).
				Args(LitStringf("%s must be called before calling %s",
					c.export(returnFnName), c.export(andDoFnName))).Obj).Obj,
				Return(Id(nilIdent))).Obj,
		c.lastResult(true),
		Assign(Sel(Id(lastIdent)).Dot(c.exportId(doFnIdent)).Obj).
			Tok(token.ASSIGN).Rhs(Id(fnFnName)).Obj,
		Return(Id(recorderReceiverIdent)),
	).Decs(stdFuncDec()).Obj
}

func (c *Converter) doReturnResultsFn(fn Func) *dst.FuncDecl {
	typePrefix := c.typePrefix(fn)
	fnName := fmt.Sprintf(double, typePrefix, doReturnFnIdent)
	params := FieldList(Field(Id(fnName)).Names(Id(fnFnName)).Obj)
	resExprs := []dst.Expr{
		Key(c.exportId(sequenceIdent)).Value(Id(sequenceIdent)).Obj,
		Key(c.exportId(doReturnFnIdent)).Value(Id(fnFnName)).Obj,
	}

	return c.returnFn(doReturnResultsFnName, fn, params, resExprs)
}

func (c *Converter) findResultsFn(fn Func) *dst.FuncDecl {
	mName := c.moqName()

	fnRecName := fmt.Sprintf(triple, mName, fn.Name, fnRecorderSuffix)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf(double, mName, fnRecorderSuffix)
	}

	incrRepeat := Expr(Call(Sel(Sel(Sel(Id(recorderReceiverIdent)).
		Dot(c.exportId(resultsIdent)).Obj).
		Dot(c.exportId(repeatIdent)).Obj).
		Dot(Id(incrementFnName)).Obj).
		Args(Sel(Sel(Sel(Id(recorderReceiverIdent)).
			Dot(c.exportId(moqIdent)).Obj).
			Dot(c.exportId(sceneIdent)).Obj).
			Dot(Id(tType)).Obj).Obj).Obj
	body := []dst.Stmt{
		If(Bin(Sel(Id(recorderReceiverIdent)).
			Dot(c.exportId(resultsIdent)).Obj).
			Op(token.NEQ).
			Y(Id(nilIdent)).Obj).
			Body(
				cloneStmt(incrRepeat),
				Return(),
			).Decs(IfDecs(dst.EmptyLine).Obj).Obj,
	}
	body = append(body, c.findRecorderResults(fn)...)
	body = append(body, cloneStmt(incrRepeat))

	return Fn(c.export(findResultsFnName)).
		Recv(Field(Star(Id(fnRecName))).Names(Id(recorderReceiverIdent)).Obj).
		Body(body...).Decs(stdFuncDec()).Obj
}

func (c *Converter) findRecorderResults(fn Func) []dst.Stmt {
	mName := c.moqName()

	results := fmt.Sprintf(triple, mName, fn.Name, resultsIdent)
	resultsByParamsType := fmt.Sprintf(triple, mName, fn.Name, resultsByParamsIdent)
	paramsKey := fmt.Sprintf(triple, mName, fn.Name, paramsKeyIdent)
	paramsKeyFn := c.export(fmt.Sprintf(double, paramsKeyFnName, fn.Name))
	resultsByParams := fmt.Sprintf(double, resultsByParamsIdent, fn.Name)
	if fn.Name == "" {
		results = fmt.Sprintf(double, mName, resultsIdent)
		resultsByParamsType = fmt.Sprintf(double, mName, resultsByParamsIdent)
		paramsKey = fmt.Sprintf(double, mName, paramsKeyIdent)
		paramsKeyFn = c.export(paramsKeyFnName)
		resultsByParams = resultsByParamsIdent
	}

	moqSel := Sel(Id(recorderReceiverIdent)).Dot(c.exportId(moqIdent)).Obj

	return []dst.Stmt{
		Assign(Id(anyCountIdent)).
			Tok(token.DEFINE).
			Rhs(Call(c.idPath("OnesCount64", "math/bits")).Args(
				Sel(Id(recorderReceiverIdent)).
					Dot(c.exportId(anyParamsIdent)).Obj).Obj).Obj,
		Assign(Id(insertAtIdent)).Tok(token.DEFINE).Rhs(LitInt(-1)).Obj,
		Var(Value(Star(Id(resultsByParamsType))).
			Names(Id(resultsIdent)).Obj),
		Range(c.selExport(moqSel, resultsByParams)).
			Key(Id(nIdent)).
			Value(Id(resIdent)).
			Tok(token.DEFINE).
			Body(
				If(Bin(Sel(Id(resIdent)).Dot(c.exportId(anyParamsIdent)).Obj).
					Op(token.EQL).
					Y(Sel(Id(recorderReceiverIdent)).
						Dot(c.exportId(anyParamsIdent)).Obj).Obj).Body(
					Assign(Id(resultsIdent)).
						Tok(token.ASSIGN).
						Rhs(Un(token.AND, Id(resIdent))).Obj,
					Break(),
				).Obj,
				If(Bin(Sel(Id(resIdent)).Dot(c.exportId(anyCountIdent)).Obj).
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
					Key(c.exportId(anyCountIdent)).Value(Id(anyCountIdent)).
						Decs(kvExprDec(dst.NewLine)).Obj,
					Key(c.exportId(anyParamsIdent)).
						Value(Sel(Id(recorderReceiverIdent)).
							Dot(c.exportId(anyParamsIdent)).Obj).
						Decs(kvExprDec(dst.None)).Obj,
					Key(c.exportId(resultsIdent)).Value(Comp(MapType(Id(paramsKey)).
						Value(Star(Id(results))).Obj).Obj).
						Decs(kvExprDec(dst.NewLine)).Obj,
				).Obj)).Obj,
			Assign(c.selExport(moqSel, resultsByParams)).
				Tok(token.ASSIGN).Rhs(Call(Id("append")).Args(
				c.selExport(moqSel, resultsByParams),
				Star(Id(resultsIdent))).Obj).Obj,
			If(Bin(Bin(Id(insertAtIdent)).Op(token.NEQ).
				Y(LitInt(-1)).Obj).Op(token.LAND).
				Y(Bin(Bin(Id(insertAtIdent)).Op(token.ADD).
					Y(LitInt(1)).Obj).Op(token.LSS).Y(Call(Id(lenFnName)).
					Args(c.selExport(moqSel, resultsByParams)).Obj).Obj).Obj).Body(
				Expr(Call(Id("copy")).Args(
					SliceExpr(c.selExport(moqSel, resultsByParams)).
						Low(Bin(Id(insertAtIdent)).Op(token.ADD).Y(LitInt(1)).Obj).Obj,
					SliceExpr(c.selExport(moqSel, resultsByParams)).
						Low(Id(insertAtIdent)).High(LitInt(0)).Obj).Obj).Obj,
				Assign(Index(c.selExport(moqSel, resultsByParams)).
					Sub(Id(insertAtIdent)).Obj).
					Tok(token.ASSIGN).
					Rhs(Star(Id(resultsIdent))).Obj,
			).Obj,
		).Decs(IfDecs(dst.EmptyLine).Obj).Obj,
		Assign(Id(paramsKeyIdent)).
			Tok(token.DEFINE).
			Rhs(Call(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(moqIdent)).Obj).Dot(Id(paramsKeyFn)).Obj).
				Args(Sel(Id(recorderReceiverIdent)).Dot(c.exportId(paramsIdent)).Obj,
					Sel(Id(recorderReceiverIdent)).Dot(c.exportId(anyParamsIdent)).Obj).Obj).
			Decs(AssignDecs(dst.None).After(dst.EmptyLine).Obj).Obj,
		Var(Value(Id("bool")).Names(Id(okIdent)).Obj),
		Assign(Sel(Id(recorderReceiverIdent)).Dot(c.exportId(resultsIdent)).Obj, Id(okIdent)).
			Tok(token.ASSIGN).
			Rhs(Index(Sel(Id(resultsIdent)).Dot(c.exportId(resultsIdent)).Obj).
				Sub(Id(paramsKeyIdent)).Obj).Obj,
		If(Un(token.NOT, Id(okIdent))).
			Body(
				Assign(Sel(Id(recorderReceiverIdent)).
					Dot(c.exportId(resultsIdent)).Obj).
					Tok(token.ASSIGN).
					Rhs(Un(
						token.AND,
						Comp(c.exportId(results)).
							Elts(
								Key(c.exportId(paramsIdent)).
									Value(Sel(Id(recorderReceiverIdent)).
										Dot(c.exportId(paramsIdent)).Obj).
									Decs(kvExprDec(dst.NewLine)).Obj,
								Key(c.exportId(resultsIdent)).Value(Id(nilIdent)).
									Decs(kvExprDec(dst.None)).Obj,
								Key(c.exportId(indexIdent)).Value(
									LitInt(0)).Decs(kvExprDec(dst.None)).Obj,
								Key(c.exportId(repeatIdent)).Value(
									Un(token.AND, Comp(c.idPath(repeatValType, moqPkg)).Obj)).
									Decs(kvExprDec(dst.None)).Obj,
							).Obj,
					)).Obj,
				Assign(Index(Sel(Id(resultsIdent)).Dot(c.exportId(resultsIdent)).Obj).
					Sub(Id(paramsKeyIdent)).Obj).
					Tok(token.ASSIGN).
					Rhs(Sel(Id(recorderReceiverIdent)).
						Dot(c.exportId(resultsIdent)).Obj).Obj,
			).Decs(IfDecs(dst.EmptyLine).Obj).Obj,
	}
}

func (c *Converter) recorderRepeatFn(fn Func) *dst.FuncDecl {
	mName := c.moqName()

	fnRecName := fmt.Sprintf(triple, mName, fn.Name, fnRecorderSuffix)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf(double, mName, fnRecorderSuffix)
	}

	lastVal := Comp(c.innerResultsStruct(c.typePrefix(fn), fn.Results)).Elts(
		Key(c.exportId(valuesIdent)).
			Value(Sel(Id(lastIdent)).Dot(c.exportId(valuesIdent)).Obj).
			Decs(kvExprDec(dst.NewLine)).Obj,
		Key(c.exportId(sequenceIdent)).
			Value(Call(Sel(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(moqIdent)).Obj).
				Dot(c.exportId(sceneIdent)).Obj).
				Dot(Id("NextRecorderSequence")).Obj).Obj).
			Decs(kvExprDec(dst.None)).Obj,
	).Obj

	return Fn(c.export(repeatFnName)).
		Recv(Field(Star(Id(fnRecName))).Names(Id(recorderReceiverIdent)).Obj).
		Params(Field(Ellipsis(c.idPath(repeaterType, moqPkg))).Names(Id(repeatersIdent)).Obj).
		Results(Field(Star(Id(fnRecName))).Obj).
		Body(
			c.helperCallExpr(Sel(Id(recorderReceiverIdent)).Dot(c.exportId(moqIdent)).Obj),
			If(Bin(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(resultsIdent)).Obj).
				Op(token.EQL).
				Y(Id(nilIdent)).Obj).
				Body(
					Expr(Call(Sel(Sel(Sel(Sel(Id(recorderReceiverIdent)).
						Dot(c.exportId(moqIdent)).Obj).
						Dot(c.exportId(sceneIdent)).Obj).
						Dot(Id(tType)).Obj).Dot(Id(fatalfFnName)).Obj).
						Args(LitStringf("%s or %s must be called before calling %s",
							c.export(returnFnName),
							c.export(doReturnResultsFnName),
							c.export(repeatFnName))).Obj).Obj,
					Return(Id(nilIdent)),
				).Obj,
			Expr(Call(Sel(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(resultsIdent)).Obj).
				Dot(c.exportId(repeatIdent)).Obj).
				Dot(Id(titler.String(repeatFnName))).Obj).
				Args(Sel(Sel(Sel(Id(recorderReceiverIdent)).
					Dot(c.exportId(moqIdent)).Obj).
					Dot(c.exportId(sceneIdent)).Obj).
					Dot(Id(tType)).Obj,
					Id(repeatersIdent)).Obj).Obj,
			c.lastResult(false),
			For(Assign(Id(nIdent)).Tok(token.DEFINE).Rhs(LitInt(0)).Obj).
				Cond(Bin(Id(nIdent)).Op(token.LSS).
					Y(Bin(Sel(Sel(Sel(Id(recorderReceiverIdent)).
						Dot(c.exportId(resultsIdent)).Obj).
						Dot(c.exportId(repeatIdent)).Obj).
						Dot(Id(resultCountIdent)).Obj).
						Op(token.SUB).Y(LitInt(1)).Obj).Obj).
				Post(IncStmt(Id(nIdent))).Body(
				If(Sel(Id(recorderReceiverIdent)).Dot(c.exportId(sequenceIdent)).Obj).Body(
					Assign(Id(lastIdent)).Tok(token.ASSIGN).Rhs(lastVal).Obj).Obj,
				Assign(Sel(Sel(Id(recorderReceiverIdent)).
					Dot(c.exportId(resultsIdent)).Obj).
					Dot(c.exportId(resultsIdent)).Obj).
					Tok(token.ASSIGN).
					Rhs(Call(Id("append")).Args(Sel(Sel(Id(recorderReceiverIdent)).
						Dot(c.exportId(resultsIdent)).Obj).
						Dot(c.exportId(resultsIdent)).Obj,
						Id(lastIdent)).Obj).Obj,
			).Obj,
			Return(Id(recorderReceiverIdent)),
		).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) prettyParamsFn(fn Func) *dst.FuncDecl {
	mName := c.moqName()
	params := fmt.Sprintf(triple, mName, fn.Name, paramsIdent)
	fnName := fmt.Sprintf(double, prettyParamsFnName, fn.Name)
	sfmt := fn.Name + "("
	if fn.Name == "" {
		sfmt = c.typ.TypeSpec.Name.Name + "("
		params = fmt.Sprintf(double, mName, paramsIdent)
		fnName = prettyParamsFnName
	}
	var pExprs []dst.Expr
	count := 0
	for _, param := range fn.Params.List {
		if len(param.Names) == 0 {
			sfmt += "%#v, "
			pExprs = append(pExprs, Sel(Id(paramsIdent)).
				Dot(c.exportId(fmt.Sprintf(unnamed, paramPrefix, count+1))).Obj)
			count++
		}

		for _, name := range param.Names {
			sfmt += "%#v, "
			vName := validName(name.Name, paramPrefix, count)
			pExprs = append(pExprs, Sel(Id(paramsIdent)).Dot(c.exportId(vName)).Obj)
			count++
		}
	}
	if count > 0 {
		sfmt = sfmt[0 : len(sfmt)-2]
	}
	sfmt += ")"
	pExprs = append([]dst.Expr{LitString(sfmt)}, pExprs...)
	return Fn(c.export(fnName)).
		Recv(Field(Star(Id(mName))).Names(Id(moqReceiverIdent)).Obj).
		Params(Field(Id(params)).Names(Id(paramsIdent)).Obj).
		Results(Field(Id("string")).Obj).
		Body(Return(
			Call(IdPath("Sprintf", "fmt")).
				Args(pExprs...).Obj)).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) paramsKeyFn(fn Func) *dst.FuncDecl {
	var stmts []dst.Stmt
	count := 0
	for _, param := range fn.Params.List {
		if len(param.Names) == 0 {
			stmts = append(stmts, c.mockFuncFindResultsParam(
				fn, fmt.Sprintf(unnamed, paramPrefix, count+1), count, param.Type)...)
			count++
		}

		for _, name := range param.Names {
			stmts = append(stmts,
				c.mockFuncFindResultsParam(fn, name.Name, count, param.Type)...)
			count++
		}
	}

	mName := c.moqName()
	params := fmt.Sprintf(triple, mName, fn.Name, paramsIdent)
	paramsKey := fmt.Sprintf(triple, mName, fn.Name, paramsKeyIdent)
	fnName := fmt.Sprintf(double, paramsKeyFnName, fn.Name)
	if fn.Name == "" {
		params = fmt.Sprintf(double, mName, paramsIdent)
		paramsKey = fmt.Sprintf(double, mName, paramsKeyIdent)
		fnName = paramsKeyFnName
	}

	pStruct, err := c.methodStruct(paramsKeyIdent, fn.Params)
	if err != nil {
		logs.Panic("Call MethodStructs first to get a meaningful error", err)
	}
	pkStruct, err := c.methodStruct(hashesIdent, fn.Params)
	if err != nil {
		logs.Panic("Call MethodStructs first to get a meaningful error", err)
	}
	stmts = append(stmts, Return(Comp(Id(paramsKey)).Elts(
		Key(c.exportId(paramsIdent)).Value(Comp(pStruct).
			Elts(c.passthroughElements(fn.Params, paramsKeyIdent, usedSuffix)...).Obj).
			Decs(kvExprDec(dst.NewLine)).Obj,
		Key(c.exportId(hashesIdent)).Value(Comp(pkStruct).
			Elts(c.passthroughElements(fn.Params, hashesIdent, usedHashSuffix)...).Obj).
			Decs(kvExprDec(dst.NewLine)).Obj).Obj))

	return Fn(c.export(fnName)).
		Recv(Field(Star(Id(mName))).Names(Id(moqReceiverIdent)).Obj).
		Params(Field(Id(params)).Names(Id(paramsIdent)).Obj,
			Field(Id("uint64")).Names(Id(anyParamsIdent)).Obj).
		Results(Field(Id(paramsKey)).Obj).
		Body(stmts...).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) mockFuncFindResultsParam(
	fn Func, pName string, paramPos int, typ dst.Expr,
) []dst.Stmt {
	comp, err := c.typeCache.IsComparable(typ)
	if err != nil {
		logs.Panic("Call MethodStructs first to get a meaningful error", err)
	}

	vName := validName(pName, paramPrefix, paramPos)

	var stmts []dst.Stmt
	pUsed := fmt.Sprintf("%s%s", vName, usedSuffix)
	if comp {
		stmts = append(stmts, Var(Value(cloneExpr(typ)).Names(Id(pUsed)).Obj))
	}
	hashUsed := fmt.Sprintf("%s%s", vName, usedHashSuffix)
	stmts = append(stmts, Var(Value(c.idPath(hashType, hashPkg)).Names(Id(hashUsed)).Obj))

	ifSel := Sel(Sel(Sel(Sel(Id(moqReceiverIdent)).
		Dot(c.exportId(runtimeIdent)).Obj).
		Dot(c.exportId(parameterIndexingIdent)).Obj).
		Dot(Id(fn.Name)).Obj).
		Dot(c.exportId(vName)).Obj
	fatalMsg := LitStringf("The %s parameter of the %s function can't be indexed by value",
		vName, fn.Name)
	if fn.Name == "" {
		ifSel = Sel(Sel(Sel(Id(moqReceiverIdent)).
			Dot(c.exportId(runtimeIdent)).Obj).
			Dot(c.exportId(parameterIndexingIdent)).Obj).
			Dot(c.exportId(vName)).Obj
		fatalMsg = LitStringf("The %s parameter can't be indexed by value", vName)
	}

	ifCond := If(Bin(ifSel).
		Op(token.EQL).
		Y(c.idPath(paramIndexByValueIdent, moqPkg)).Obj)
	pKeySel := Id(paramsIdent)
	hashAssign := Assign(Id(hashUsed)).
		Tok(token.ASSIGN).
		Rhs(c.passthroughValue(Id(vName), true, pKeySel)).Obj
	var usedBody []dst.Stmt
	if comp {
		usedBody = append(usedBody, ifCond.Body(
			Assign(Id(pUsed)).
				Tok(token.ASSIGN).
				Rhs(c.passthroughValue(Id(vName), false, pKeySel)).Obj).
			Else(hashAssign).Obj)
	} else {
		usedBody = append(usedBody, ifCond.Body(Expr(Call(Sel(Sel(Sel(Id(moqReceiverIdent)).
			Dot(c.exportId(sceneIdent)).Obj).
			Dot(Id(tType)).Obj).
			Dot(Id(fatalfFnName)).Obj).
			Args(fatalMsg).Obj).Obj).Obj)
		usedBody = append(usedBody, hashAssign)
	}

	stmts = append(stmts, If(Bin(Bin(Id(anyParamsIdent)).
		Op(token.AND).
		Y(Paren(Bin(LitInt(1)).Op(token.SHL).Y(LitInt(paramPos)).Obj)).Obj).
		Op(token.EQL).
		Y(LitInt(0)).Obj).
		Body(usedBody...).Obj)
	return stmts
}

func (c *Converter) lastResult(forUpdate bool) *dst.AssignStmt {
	var rhs dst.Expr = Index(Sel(Sel(Id(recorderReceiverIdent)).
		Dot(c.exportId(resultsIdent)).Obj).
		Dot(c.exportId(resultsIdent)).Obj).
		Sub(Bin(Call(Id(lenFnName)).
			Args(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(resultsIdent)).Obj).
				Dot(c.exportId(resultsIdent)).Obj).Obj).
			Op(token.SUB).
			Y(LitInt(1)).Obj).Obj
	if forUpdate {
		rhs = Un(token.AND, rhs)
	}

	return Assign(Id(lastIdent)).
		Tok(token.DEFINE).
		Rhs(rhs).Obj
}

func (c *Converter) recorderSeqFns(fn Func) []dst.Decl {
	return []dst.Decl{
		c.recorderSeqFn(fn, "seq", "true"),
		c.recorderSeqFn(fn, "noSeq", "false"),
	}
}

func (c *Converter) recorderSeqFn(fn Func, fnName, assign string) *dst.FuncDecl {
	mName := c.moqName()

	fnRecName := fmt.Sprintf(triple, mName, fn.Name, fnRecorderSuffix)
	if fn.Name == "" {
		fnRecName = fmt.Sprintf(double, mName, fnRecorderSuffix)
	}

	moqSel := Sel(Id(recorderReceiverIdent)).Dot(c.exportId(moqIdent)).Obj

	fnName = c.export(fnName)
	return Fn(fnName).
		Results(Field(Star(Id(fnRecName))).Obj).
		Recv(Field(Star(Id(fnRecName))).Names(Id(recorderReceiverIdent)).Obj).
		Body(
			c.helperCallExpr(Sel(Id(recorderReceiverIdent)).Dot(c.exportId(moqIdent)).Obj),
			If(Bin(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(resultsIdent)).Obj).
				Op(token.NEQ).
				Y(Id(nilIdent)).Obj).
				Body(
					Expr(Call(Sel(Sel(Sel(Sel(Id(recorderReceiverIdent)).
						Dot(c.exportId(moqIdent)).Obj).
						Dot(c.exportId(sceneIdent)).Obj).
						Dot(Id(tType)).Obj).
						Dot(Id(fatalfFnName)).Obj).
						Args(LitStringf(
							"%s must be called before %s or %s calls, recording %%s",
							fnName, c.export(returnFnName), c.export(doReturnResultsFnName)),
							c.callPrettyParams(fn, moqSel,
								Id(recorderReceiverIdent))).Obj).Obj,
					Return(Id(nilIdent)),
				).Obj,
			Assign(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(sequenceIdent)).Obj).
				Tok(token.ASSIGN).
				Rhs(Id(assign)).
				Decs(AssignDecs(dst.NewLine).Obj).Obj,
			Return(Id(recorderReceiverIdent)),
		).Decs(stdFuncDec()).Obj
}

func (c *Converter) callPrettyParams(fn Func, moqExpr, paramsExpr dst.Expr) *dst.CallExpr {
	prettyParamsFn := prettyParamsFnName
	if fn.Name != "" {
		prettyParamsFn = fmt.Sprintf(double, prettyParamsFnName, fn.Name)
	}

	return Call(c.selExport(moqExpr, prettyParamsFn)).
		Args(c.selExport(paramsExpr, paramsIdent)).Obj
}

func (c *Converter) passthroughElements(fl *dst.FieldList, label, valSuffix string) []dst.Expr {
	if fl == nil {
		return nil
	}

	unnamedPrefix, dropNonComparable := labelDirection(label)
	var elts []dst.Expr
	beforeDec := dst.NewLine
	count := 0
	for _, field := range fl.List {
		comp, err := c.typeCache.IsComparable(field.Type)
		if err != nil {
			logs.Panic("Call MethodStructs first to get a meaningful error", err)
		}

		if len(field.Names) == 0 {
			if comp || !dropNonComparable {
				pName := fmt.Sprintf(unnamed, unnamedPrefix, count+1)
				elts = append(elts, Key(c.exportId(pName)).Value(
					c.passthroughValue(Id(pName+valSuffix), false, nil)).
					Decs(kvExprDec(beforeDec)).Obj)
				beforeDec = dst.None
			}
			count++
		}

		for _, name := range field.Names {
			if comp || !dropNonComparable {
				vName := validName(name.Name, unnamedPrefix, count)
				elts = append(elts, Key(c.exportId(vName)).Value(
					c.passthroughValue(Id(vName+valSuffix), false, nil)).
					Decs(kvExprDec(beforeDec)).Obj)
				beforeDec = dst.None
			}
			count++
		}
	}

	return elts
}

func (c *Converter) passthroughValue(
	src *dst.Ident, needComparable bool, sel dst.Expr,
) dst.Expr {
	var val dst.Expr = src
	if sel != nil {
		val = c.selExport(sel, src.Name)
	}
	if needComparable {
		val = Call(c.idPath("DeepHash", hashPkg)).Args(val).Obj
	}
	return val
}

func passthroughFields(prefix string, fields *dst.FieldList) []dst.Expr {
	var exprs []dst.Expr
	count := 0
	for _, f := range fields.List {
		if len(f.Names) == 0 {
			exprs = append(exprs, Idf(unnamed, prefix, count+1))
			count++
		}

		for _, name := range f.Names {
			exprs = append(exprs, Id(validName(name.Name, prefix, count)))
			count++
		}
	}
	return exprs
}

func (c *Converter) assignResult(resFL *dst.FieldList) []dst.Stmt {
	var assigns []dst.Stmt
	if resFL != nil {
		count := 0
		for _, result := range resFL.List {
			if len(result.Names) == 0 {
				rName := fmt.Sprintf(unnamed, resultPrefix, count+1)
				assigns = append(assigns, Assign(Id(rName)).
					Tok(token.ASSIGN).
					Rhs(Sel(Sel(Id(resultIdent)).
						Dot(c.exportId(valuesIdent)).Obj).
						Dot(c.exportId(rName)).Obj).Obj)
				count++
			}

			for _, name := range result.Names {
				vName := validName(name.Name, resultPrefix, count)
				assigns = append(assigns, Assign(Id(vName)).
					Tok(token.ASSIGN).
					Rhs(Sel(Sel(Id(resultIdent)).
						Dot(c.exportId(valuesIdent)).Obj).
						Dot(c.exportId(vName)).Obj).Obj)
				count++
			}
		}
	}
	return assigns
}

func cloneAndNameUnnamed(prefix string, fieldList *dst.FieldList) *dst.FieldList {
	fieldList = cloneFieldList(fieldList, false)
	if fieldList != nil {
		count := 0
		for _, f := range fieldList.List {
			if len(f.Names) == 0 {
				f.Names = []*dst.Ident{Idf(unnamed, prefix, count+1)}
			}
			for n, name := range f.Names {
				f.Names[n] = Idf(validName(name.Name, prefix, count))
				count++
			}
		}
	}
	return fieldList
}

func validName(name, prefix string, count int) string {
	if _, ok := invalidNames[name]; ok {
		name = fmt.Sprintf(unnamed, prefix, count+1)
	}
	return name
}

func (c *Converter) moqName() string {
	return c.export(moqIdent + titler.String(c.typ.TypeSpec.Name.Name))
}

func (c *Converter) export(name string) string {
	if c.isExported {
		name = titler.String(name)
	}
	return name
}

func (c *Converter) exportId(name string) *dst.Ident {
	return Id(c.export(name))
}

func (c *Converter) idPath(name, path string) *dst.Ident {
	switch path {
	case "":
		return IdPath(name, c.typ.InPkgPath)
	case c.typ.OutPkgPath:
		return Id(name)
	default:
		return IdPath(name, path)
	}
}

func (c *Converter) helperCallExpr(selector dst.Expr) dst.Stmt {
	return Expr(Call(Sel(Sel(Sel(selector).
		Dot(c.exportId(sceneIdent)).Obj).
		Dot(Id(tType)).Obj).
		Dot(Id(helperFnName)).Obj).Obj).Obj
}

func (c *Converter) selExport(x dst.Expr, sel string) dst.Expr {
	if x == nil {
		return Id(sel)
	}

	switch v := x.(type) {
	case *dst.SelectorExpr:
		return Sel(cloneSelect(v)).Dot(c.exportId(sel)).Obj
	case *dst.Ident:
		return Sel(cloneIdent(v)).Dot(c.exportId(sel)).Obj
	default:
		logs.Panicf("unsupported selector type: %#v", v)
		return nil
	}
}

func stdFuncDec() dst.FuncDeclDecorations {
	return dst.FuncDeclDecorations{
		NodeDecs: dst.NodeDecs{Before: dst.EmptyLine, After: dst.EmptyLine},
	}
}

func labelDirection(label string) (string, bool) {
	var unnamedPrefix string
	var dropNonComparable bool
	switch label {
	case paramsIdent:
		unnamedPrefix = paramPrefix
		dropNonComparable = false
	case paramsKeyIdent:
		unnamedPrefix = paramPrefix
		dropNonComparable = true
	case hashesIdent:
		unnamedPrefix = paramPrefix
		dropNonComparable = false
	case resultsIdent:
		unnamedPrefix = resultPrefix
		dropNonComparable = false
	default:
		logs.Panicf("Unknown label: %s", label)
	}

	return unnamedPrefix, dropNonComparable
}

func isVariadic(fl *dst.FieldList) bool {
	if len(fl.List) > 0 {
		if _, ok := fl.List[len(fl.List)-1].Type.(*dst.Ellipsis); ok {
			return true
		}
	}
	return false
}

func cloneIdent(i *dst.Ident) dst.Expr {
	//nolint:forcetypeassert // if dst.Clone returns a different type, panic
	return dst.Clone(i).(*dst.Ident)
}

func cloneFieldList(fl *dst.FieldList, removeNames bool) *dst.FieldList {
	if fl != nil {
		//nolint:forcetypeassert // if dst.Clone returns a different type, panic
		fl = dst.Clone(fl).(*dst.FieldList)
		for _, field := range fl.List {
			if id, ok := field.Type.(*dst.Ident); ok {
				if id.Path == "builtin" {
					// When mocking errors, the string return value is reported
					// in the builtin package
					id.Path = ""
				}
			}
		}
		if removeNames {
			for _, field := range fl.List {
				for n := range field.Names {
					field.Names[n] = Id(sep)
				}
			}
		}
	}
	return fl
}

func cloneSelect(sel *dst.SelectorExpr) dst.Expr {
	//nolint:forcetypeassert // if dst.Clone returns a different type, panic
	return dst.Clone(sel).(*dst.SelectorExpr)
}

func cloneExpr(expr dst.Expr) dst.Expr {
	//nolint:forcetypeassert // if dst.Clone returns a different type, panic
	return dst.Clone(expr).(dst.Expr)
}

func cloneStmt(stmt dst.Stmt) dst.Stmt {
	//nolint:forcetypeassert // if dst.Clone returns a different type, panic
	return dst.Clone(stmt).(dst.Stmt)
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
