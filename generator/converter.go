package generator

import (
	"errors"
	"fmt"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/dave/dst"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"moqueries.org/runtime/logs"

	. "moqueries.org/cli/ast"
)

const (
	moqueriesPkg = "moqueries.org/runtime"
	hashPkg      = moqueriesPkg + "/hash"
	implPkg      = moqueriesPkg + "/impl"
	moqPkg       = moqueriesPkg + "/moq"

	intType           = "int"
	configType        = "Config"
	hashType          = "Hash"
	moqType           = "Moq"
	paramIndexingType = "ParamIndexing"
	recorderType      = "Recorder"
	repeaterType      = "Repeater"
	sceneType         = "Scene"
	tType             = "T"

	adaptorIdent           = "adaptor"
	adaptorReceiverIdent   = "a"
	anyParamsIdent         = "anyParams"
	anyParamsReceiverIdent = "a"
	blankIdent             = "_"
	configIdent            = "config"
	doFnIdent              = "doFn"
	doReturnFnIdent        = "doReturnFn"
	hashesIdent            = "hashes"
	iIdent                 = "i"
	mockIdent              = "mock"
	moqIdent               = "moq"
	moqReceiverIdent       = "m"
	nilIdent               = "nil"
	paramIndexingIdent     = "paramIndexing"
	parameterIndexingIdent = "parameterIndexing"
	paramIndexByValueIdent = "ParamIndexByValue"
	paramIndexByHashIdent  = "ParamIndexByHash"
	paramsIdent            = "params"
	paramsKeyIdent         = "paramsKey"
	recorderIdent          = "recorder"
	recorderReceiverIdent  = "r"
	repeatersIdent         = "repeaters"
	resultIdent            = "result"
	resultsIdent           = "results"
	runtimeIdent           = "runtime"
	sceneIdent             = "scene"
	sequenceIdent          = "sequence"

	andDoFnName            = "andDo"
	anyParamFnName         = "AnyParam"
	assertFnName           = "AssertExpectationsMet"
	doReturnResultsFnName  = "doReturnResults"
	fnFnName               = "fn"
	functionFnName         = "Function"
	hashOnlyParamKeyFnName = "HashOnlyParamKey"
	helperFnName           = "Helper"
	isAnyPermittedFnName   = "IsAnyPermitted"
	mockFnName             = "mock"
	newMoqFnName           = "NewMoq"
	onCallFnName           = "onCall"
	paramKeyFnName         = "ParamKey"
	paramsKeyFnName        = "paramsKey"
	prettyParamsFnName     = "prettyParams"
	repeatFnName           = "repeat"
	resetFnName            = "Reset"
	returnResultsFnName    = "returnResults"
	seqFnName              = "Seq"

	recorderSuffix = "recorder"
	paramPrefix    = "param"
	resultPrefix   = "result"
	usedSuffix     = "Used"
	usedHashSuffix = "UsedHash"

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
	blankIdent:            {},
	iIdent:                {},
	intType:               {},
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

	err error
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
	Name       string
	ParentType TypeInfo
	FuncType   *dst.FuncType
}

// BaseDecls generates the base declarations used to store the moq's state and
// establish the type
func (c *Converter) BaseDecls() ([]dst.Decl, error) {
	mName := c.moqName()
	moqName := fmt.Sprintf(double, mName, mockIdent)

	var fields []*dst.Field
	if c.isInterface() {
		fields = append(fields, Field(Star(c.genericExpr(Id(moqName), clone)).Obj).Names(c.exportId(moqIdent)).
			Decs(FieldDecs(dst.None, dst.EmptyLine).Obj).Obj)
	}

	for _, fn := range c.typ.Funcs {
		fnPrefix := c.export(moqIdent)
		if fn.Name != "" {
			fnPrefix += sep + fn.Name
		}
		fields = append(fields,
			Field(Star(Index(IdPath(moqType, implPkg)).
				Sub(c.baseTypeParams(fn)...).Obj).Obj).Names(c.exportId(fnPrefix)).Obj)
	}

	fields = append(fields, Field(c.exportId(fmt.Sprintf(double, mName, runtimeIdent))).
		Names(c.exportId(runtimeIdent)).
		Decs(FieldDecs(dst.EmptyLine, dst.None).Obj).Obj)

	var decls []dst.Decl

	typeName := c.typeName()
	id := IdPath(typeName, c.typ.TypeInfo.PkgPath)
	idStr := fmt.Sprintf("%s.%s", filepath.Base(id.Path), id.Name)
	if c.typ.Reduced && !c.typ.TypeInfo.Fabricated {
		typeName = fmt.Sprintf(double, typeName, "reduced")
		id = Id(typeName)
		idStr = typeName
	}

	if c.isInterface() {
		if c.typ.TypeInfo.Fabricated {
			id = Id(id.Name)
		}
		decls = append(decls, VarDecl(
			Value(c.genericExpr(id, typeAssertSafe)).Names(Id(blankIdent)).
				Values(Call(Paren(Star(c.genericExpr(Id(moqName), typeAssertSafe)).Obj)).
					Args(Id(nilIdent)).Obj).Obj).
			Decs(genDeclDecf("// The following type assertion assures"+
				" that %s is mocked completely", idStr)).Obj,
		)
	}

	if c.typ.TypeInfo.Fabricated || c.typ.Reduced {
		typ := c.resolveExpr(cloneExpr(c.typ.TypeInfo.Type.Type), c.typ.TypeInfo)
		msg := "emitted when mocking functions directly and not from a function type"
		if c.isInterface() {
			msg = "emitted when mocking a collections of methods directly and not from an interface type"
		}
		if !c.typ.TypeInfo.Fabricated {
			msg = "emitted when the original interface contains non-exported methods"
		}
		decls = append(decls, TypeDecl(TypeSpec(typeName).Type(typ).TypeParams(c.typeParams()).Obj).
			Decs(genDeclDecf("// %s is the fabricated implementation type of this mock (%s)",
				typeName, msg)).Obj)
	}

	decls = append(decls, TypeDecl(
		TypeSpec(mName).Type(Struct(fields...)).TypeParams(c.typeParams()).Obj).
		Decs(genDeclDecf("// %s holds the state of a moq of the %s type",
			mName, typeName)).Obj)

	if c.err != nil {
		return nil, c.err
	}

	return decls, nil
}

// MockStructs generates several structs used to isolate functionality
// for the moq
func (c *Converter) MockStructs() ([]dst.Decl, error) {
	mName := c.moqName()

	gen := func(suffix string) *dst.GenDecl {
		iName := fmt.Sprintf(double, mName, suffix)

		return TypeDecl(TypeSpec(iName).Type(Struct(Field(Star(c.genericExpr(Id(mName), clone)).Obj).
			Names(c.exportId(moqIdent)).Obj)).TypeParams(c.typeParams()).Obj).
			Decs(genDeclDecf("// %s isolates the %s interface of the %s type",
				iName, suffix, c.typeName())).Obj
	}

	var iStructs []dst.Decl
	if c.isInterface() {
		iStructs = append(iStructs, gen(mockIdent))
	}

	_, isId := c.typ.TypeInfo.Type.Type.(*dst.Ident)
	if c.isInterface() || isId {
		iStructs = append(iStructs, gen(recorderIdent))
	}

	rName := fmt.Sprintf(double, mName, runtimeIdent)
	iStructs = append(iStructs, TypeDecl(TypeSpec(rName).Type(c.runtimeStruct()).Obj).
		Decs(genDeclDecf("// %s holds runtime configuration for the %s type",
			rName, c.typeName())).Obj)

	if c.err != nil {
		return nil, c.err
	}

	return iStructs, nil
}

// MethodStructs generates a structure for storing a set of parameters or
// a set of results for a method invocation of a mock
func (c *Converter) MethodStructs(fn Func) ([]dst.Decl, error) {
	prefix := c.typePrefix(fn)

	decls := []dst.Decl{
		c.adaptorStructDecl(prefix),
		c.paramsStructDecl(prefix, false, fn.FuncType.Params, fn.ParentType),
		c.paramsStructDecl(prefix, true, fn.FuncType.Params, fn.ParentType),
		c.resultsStruct(prefix, fn.FuncType.Results, fn.ParentType),
		c.fnParamIndexingStruct(prefix, fn.FuncType.Params, fn.ParentType),
		c.doFuncType(prefix, fn.FuncType.Params, fn.ParentType),
		c.doReturnFuncType(prefix, fn),
		c.recorderStruct(prefix, fn),
		c.anyParamsStruct(prefix),
	}

	if c.err != nil {
		return nil, c.err
	}

	return decls, nil
}

// NewFunc generates a function for constructing a moq
func (c *Converter) NewFunc() (*dst.FuncDecl, error) {
	fnName := c.export("newMoq" + c.typeName())
	mName := c.moqName()
	body := c.defineAdaptors()
	var elts []dst.Expr
	var moqStmt dst.Stmt
	if c.isInterface() {
		elts = append(elts, Key(c.exportId(moqIdent)).Value(Un(token.AND, Comp(
			c.genericExpr(Id(fmt.Sprintf(double, mName, mockIdent)), clone)).Obj)).
			Decs(KeyValueDecs(dst.None).After(dst.EmptyLine).Obj).Obj)
		moqStmt = Assign(Sel(Sel(Id(moqReceiverIdent)).
			Dot(c.exportId(moqIdent)).Obj).
			Dot(c.exportId(moqIdent)).Obj).
			Tok(token.ASSIGN).
			Rhs(Id(moqReceiverIdent)).
			Decs(AssignDecs(dst.None).After(dst.EmptyLine).Obj).Obj
	}
	elts = append(elts, c.newMoqs()...)
	elts = append(elts, Key(c.exportId(runtimeIdent)).Value(Comp(
		c.exportId(fmt.Sprintf(double, mName, runtimeIdent))).
		Elts(c.runtimeValues()...).Obj).Decs(kvExprDec(dst.None)).
		Decs(kvExprDec(dst.EmptyLine)).Obj)
	body = append(body,
		Assign(Id(moqReceiverIdent)).Tok(token.DEFINE).Rhs(Un(token.AND,
			Comp(c.genericExpr(Id(mName), clone)).Elts(elts...).Decs(litDec()).Obj)).Obj,
		moqStmt,
	)
	body = append(body, c.updateAdaptors()...)
	body = append(body,
		Expr(Call(Sel(Id(sceneIdent)).Dot(Id("AddMoq")).Obj).
			Args(Id(moqReceiverIdent)).Decs(CallDecs(dst.EmptyLine, dst.None)).Obj).Obj,
		Return(Id(moqReceiverIdent)),
	)

	decl := Fn(fnName).
		TypeParams(c.typeParams()).
		Params(
			Field(Star(c.idPath(sceneType, moqPkg)).Obj).Names(Id(sceneIdent)).Obj,
			Field(Star(c.idPath(configType, moqPkg)).Obj).Names(Id(configIdent)).Obj,
		).
		Results(Field(Star(c.genericExpr(Id(mName), clone)).Obj).Obj).
		Body(body...).
		Decs(fnDeclDecf("// %s creates a new moq of the %s type",
			fnName, c.typeName())).Obj

	if c.err != nil {
		return nil, c.err
	}

	return decl, nil
}

// IsolationAccessor generates a function to access an isolation interface
func (c *Converter) IsolationAccessor(suffix, fnName string) (*dst.FuncDecl, error) {
	mName := c.moqName()
	iName := c.genericExpr(Id(fmt.Sprintf(double, mName, suffix)), clone)

	var retVal dst.Expr
	retVal = Sel(Id(moqReceiverIdent)).Dot(c.exportId(moqIdent)).Obj
	if fnName != mockFnName {
		retVal = Un(token.AND, Comp(cloneExpr(iName)).Elts(
			Key(c.exportId(moqIdent)).Value(Id(moqReceiverIdent)).
				Decs(kvExprDec(dst.None)).Obj).Decs(litDec()).Obj,
		)
	}

	fnName = c.export(fnName)
	decl := Fn(fnName).
		Recv(Field(Star(c.genericExpr(Id(mName), clone)).Obj).Names(Id(moqReceiverIdent)).Obj).
		Results(Field(Star(iName).Obj).Obj).
		Body(Return(retVal)).
		Decs(fnDeclDecf("// %s returns the %s implementation of the %s type",
			fnName, suffix, c.typeName())).Obj

	if c.err != nil {
		return nil, c.err
	}

	return decl, nil
}

// FuncClosure generates a mock implementation of function type wrapped in a
// closure
func (c *Converter) FuncClosure(fn Func) (*dst.FuncDecl, error) {
	mName := c.moqName()
	resType := c.idPath(c.typeName(), c.typ.TypeInfo.Type.Name.Path)
	if c.typ.TypeInfo.Fabricated {
		resType = Id(c.typeName())
	}

	decl := Fn(c.export(mockFnName)).
		Recv(Field(Star(c.genericExpr(Id(mName), clone)).Obj).Names(Id(moqReceiverIdent)).Obj).
		Results(Field(c.genericExpr(resType, clone)).Obj).
		Body(Return(FnLit(FnType(c.cloneAndNameUnnamed(paramPrefix, fn.FuncType.Params, fn.ParentType)).
			Results(c.cloneFieldList(fn.FuncType.Results, true, fn.ParentType)).Obj).
			Body(c.mockFunc(c.typePrefix(fn), fn)...).Obj)).
		Decs(fnDeclDecf("// %s returns the %s implementation of the %s type",
			c.export(mockFnName), moqIdent, c.typeName())).Obj

	if c.err != nil {
		return nil, c.err
	}

	return decl, nil
}

// MockMethod generates a mock implementation of a method
func (c *Converter) MockMethod(fn Func) (*dst.FuncDecl, error) {
	mName := c.moqName()

	fnName := fn.Name
	typePrefix := c.typePrefix(fn)
	if fnName == "" {
		fnName = c.export(fnFnName)
	}

	decl := Fn(fnName).
		Recv(Field(Star(c.genericExpr(Id(fmt.Sprintf(double, mName, mockIdent)), clone)).Obj).
			Names(Id(moqReceiverIdent)).Obj).
		ParamList(c.cloneAndNameUnnamed(paramPrefix, fn.FuncType.Params, fn.ParentType)).
		ResultList(c.cloneFieldList(fn.FuncType.Results, true, fn.ParentType)).
		Body(c.mockFunc(typePrefix, fn)...).
		Decs(stdFuncDec()).Obj

	if c.err != nil {
		return nil, c.err
	}

	return decl, nil
}

// RecorderMethods generates a recorder implementation of a method and
// associated return method
func (c *Converter) RecorderMethods(fn Func) ([]dst.Decl, error) {
	decls := []dst.Decl{c.recorderFn(fn)}

	decls = append(decls, c.anyParamFns(fn)...)
	decls = append(decls, c.recorderSeqFns(fn)...)
	decls = append(decls,
		c.returnResultsFn(fn),
		c.andDoFn(fn),
		c.doReturnResultsFn(fn),
		c.recorderRepeatFn(fn),
		c.prettyParamsFn(fn),
		c.paramsKeyFn(fn),
	)

	if c.err != nil {
		return nil, c.err
	}

	return decls, nil
}

// ResetMethod generates a method to reset the moq's state
func (c *Converter) ResetMethod() (*dst.FuncDecl, error) {
	var stmts []dst.Stmt
	for _, fn := range c.typ.Funcs {
		fieldSuffix := ""
		if c.isInterface() {
			fieldSuffix = sep + fn.Name
		}

		stmts = append(stmts,
			Expr(Call(Sel(Sel(Id(moqReceiverIdent)).
				Dot(c.exportId(moqIdent+fieldSuffix)).Obj).Dot(Id(resetFnName)).Obj).Obj).
				Decs(ExprDecs(dst.NewLine).Obj).Obj)
	}

	decl := Fn(resetFnName).
		Recv(Field(Star(c.genericExpr(Id(c.moqName()), clone)).Obj).Names(Id(moqReceiverIdent)).Obj).
		Body(stmts...).
		Decs(fnDeclDecf("// %s resets the state of the moq", resetFnName)).Obj

	if c.err != nil {
		return nil, c.err
	}

	return decl, nil
}

// AssertMethod generates a method to assert all expectations are met
func (c *Converter) AssertMethod() (*dst.FuncDecl, error) {
	sel := func(fn Func) dst.Expr {
		dotId := moqIdent
		if c.isInterface() {
			dotId = fmt.Sprintf(double, moqIdent, fn.Name)
		}
		return Sel(Id(moqReceiverIdent)).Dot(c.exportId(dotId)).Obj
	}
	var stmts []dst.Stmt
	if len(c.typ.Funcs) != 0 {
		stmts = append(stmts, c.helperCallExpr(sel(c.typ.Funcs[0])))
	}
	for _, fn := range c.typ.Funcs {
		stmts = append(stmts, Expr(Call(Sel(sel(fn)).Dot(Id(titler.String(assertFnName))).Obj).Obj).Obj)
	}

	decl := Fn(assertFnName).
		Recv(Field(Star(c.genericExpr(Id(c.moqName()), clone)).Obj).Names(Id(moqReceiverIdent)).Obj).
		Body(stmts...).
		Decs(fnDeclDecf("// %s asserts that all expectations have been met",
			assertFnName)).Obj

	if c.err != nil {
		return nil, c.err
	}

	return decl, nil
}

func (c *Converter) typePrefix(fn Func) string {
	mName := c.moqName()
	typePrefix := mName
	if c.isInterface() {
		typePrefix = fmt.Sprintf(double, mName, fn.Name)
	}
	return typePrefix
}

func (c *Converter) baseTypeParams(fn Func) []dst.Expr {
	typePrefix := c.typePrefix(fn)
	return []dst.Expr{
		Star(c.genericExpr(Idf(double, typePrefix, adaptorIdent), clone)).
			Decs(StarDecs(dst.NewLine, dst.NewLine)).Obj,
		c.genericExpr(IdDecs(Idf(double, typePrefix, paramsIdent), IdentDecs(dst.NewLine, dst.NewLine)), clone),
		c.genericExpr(IdDecs(Idf(double, typePrefix, paramsKeyIdent), IdentDecs(dst.NewLine, dst.NewLine)), clone),
		c.genericExpr(IdDecs(Idf(double, typePrefix, resultsIdent), IdentDecs(dst.NewLine, dst.NewLine)), clone),
	}
}

func (c *Converter) runtimeStruct() *dst.StructType {
	mName := c.moqName()
	var piExpr dst.Expr
	var piFields []*dst.Field
	for _, fn := range c.typ.Funcs {
		if c.isInterface() {
			piFields = append(piFields,
				Field(Idf(triple, mName, fn.Name, paramIndexingIdent)).Names(Id(fn.Name)).Obj)
		} else {
			piExpr = Idf(double, mName, paramIndexingIdent)
		}
	}
	if c.isInterface() {
		piExpr = Struct(piFields...)
	}
	return Struct(Field(piExpr).Names(c.exportId(parameterIndexingIdent)).Obj)
}

func (c *Converter) adaptorStructDecl(prefix string) *dst.GenDecl {
	aName := fmt.Sprintf(double, prefix, adaptorIdent)
	return TypeDecl(TypeSpec(aName).
		Type(Struct(Field(Star(c.genericExpr(Id(c.moqName()), clone)).Obj).
			Names(c.exportId(moqIdent)).Obj)).
		TypeParams(c.typeParams()).Obj).
		Decs(genDeclDecf("// %s adapts %s as needed by the runtime", aName, c.moqName())).Obj
}

func (c *Converter) paramsStructDecl(
	prefix string, paramsKey bool, fieldList *dst.FieldList, parentType TypeInfo,
) *dst.GenDecl {
	var mStruct *dst.StructType
	var label, goDocDesc string
	if paramsKey {
		label = paramsKeyIdent
		goDocDesc = "map key params"

		mStruct = Struct(Field(c.methodStruct(paramsKeyIdent, fieldList, parentType)).
			Names(c.exportId(paramsIdent)).Obj,
			Field(c.methodStruct(hashesIdent, fieldList, parentType)).
				Names(c.exportId(hashesIdent)).Obj)
	} else {
		label = paramsIdent
		goDocDesc = label

		mStruct = c.methodStruct(paramsIdent, fieldList, parentType)
	}

	structName := fmt.Sprintf(double, prefix, label)
	return TypeDecl(TypeSpec(structName).Type(mStruct).TypeParams(c.typeParams()).Obj).
		Decs(genDeclDecf("// %s holds the %s of the %s type",
			structName, goDocDesc, c.typeName())).Obj
}

func (c *Converter) methodStruct(label string, fieldList *dst.FieldList, parentType TypeInfo) *dst.StructType {
	unnamedPrefix, _ := labelDirection(label)
	fieldList = c.cloneFieldList(fieldList, false, parentType)

	var fList []*dst.Field
	if fieldList != nil {
		count := 1
		for _, f := range fieldList.List {
			for n, name := range f.Names {
				f.Names[n] = Id(c.export(validName(name.Name, unnamedPrefix, count)))
				count++
			}

			if len(f.Names) == 0 {
				f.Names = []*dst.Ident{c.exportId(fmt.Sprintf(unnamed, unnamedPrefix, count))}
				count++
			}

			typ := c.comparableType(label, f.Type)
			if typ != nil {
				f.Type = typ
				fList = append(fList, f)
			}
		}
	}

	if len(fList) != 0 {
		fieldList.List = fList
	} else {
		fieldList = nil
	}
	return StructFromList(fieldList)
}

func (c *Converter) comparableType(label string, typ dst.Expr) dst.Expr {
	switch label {
	case paramsIdent:
	case resultsIdent:
	case paramsKeyIdent:
		comp, err := c.typeCache.IsComparable(typ, c.typ.TypeInfo)
		if err != nil {
			if c.err == nil {
				c.err = err
			}
			return nil
		}
		if !comp {
			// Non-comparable params are not represented in the params section
			// of the paramsKey
			return nil
		}
	case hashesIdent:
		// Everything is represented as a hash in the hashes section of the
		// paramsKey
		return c.idPath(hashType, hashPkg)
	default:
		logs.Panicf("Unknown label: %s", label)
	}

	if ellipsis, ok := typ.(*dst.Ellipsis); ok {
		// Ellipsis params are represented as a slice (when not comparable)
		return SliceType(ellipsis.Elt)
	}

	return typ
}

func (c *Converter) resultsStruct(prefix string, results *dst.FieldList, parentType TypeInfo) *dst.GenDecl {
	structName := fmt.Sprintf(double, prefix, resultsIdent)

	return TypeDecl(TypeSpec(structName).
		Type(c.methodStruct(resultsIdent, results, parentType)).TypeParams(c.typeParams()).Obj).
		Decs(genDeclDecf("// %s holds the results of the %s type", structName, c.typeName())).Obj
}

func (c *Converter) fnParamIndexingStruct(prefix string, params *dst.FieldList, parentType TypeInfo) *dst.GenDecl {
	structName := fmt.Sprintf(double, prefix, paramIndexingIdent)
	unnamedPrefix, _ := labelDirection(paramsIdent)
	params = c.cloneFieldList(params, false, parentType)

	var fList []*dst.Field
	if params != nil {
		count := 1
		for _, f := range params.List {
			for n, name := range f.Names {
				f.Names[n] = Id(c.export(validName(name.Name, unnamedPrefix, count)))
				count++
			}

			if len(f.Names) == 0 {
				f.Names = []*dst.Ident{c.exportId(fmt.Sprintf(unnamed, unnamedPrefix, count))}
				count++
			}

			f.Type = IdPath(paramIndexingType, moqPkg)
			fList = append(fList, f)
		}
	}

	return TypeDecl(TypeSpec(structName).Type(Struct(fList...)).Obj).
		Decs(genDeclDecf("// %s holds the parameter indexing runtime configuration for the %s type",
			structName, c.typeName())).Obj
}

func (c *Converter) doFuncType(prefix string, params *dst.FieldList, parentType TypeInfo) *dst.GenDecl {
	fnName := fmt.Sprintf(double, prefix, doFnIdent)
	return TypeDecl(TypeSpec(fnName).
		Type(FnType(c.cloneFieldList(params, false, parentType)).Obj).
		TypeParams(c.typeParams()).Obj).
		Decs(genDeclDecf(
			"// %s defines the type of function needed when calling %s for the %s type",
			fnName,
			c.export(andDoFnName),
			c.typeName())).Obj
}

func (c *Converter) doReturnFuncType(prefix string, fn Func) *dst.GenDecl {
	fnName := fmt.Sprintf(double, prefix, doReturnFnIdent)
	return TypeDecl(TypeSpec(fnName).
		Type(FnType(c.cloneFieldList(fn.FuncType.Params, false, fn.ParentType)).
			Results(c.cloneFieldList(fn.FuncType.Results, false, fn.ParentType)).Obj).
		TypeParams(c.typeParams()).Obj).
		Decs(genDeclDecf(
			"// %s defines the type of function needed when calling %s for the %s type",
			fnName,
			c.export(doReturnResultsFnName),
			c.typeName())).Obj
}

func (c *Converter) recorderStruct(prefix string, fn Func) *dst.GenDecl {
	mName := c.moqName()
	structName := fmt.Sprintf(double, prefix, recorderSuffix)
	return TypeDecl(TypeSpec(structName).Type(Struct(Field(Star(Index(IdPath(recorderType, implPkg)).
		Sub(c.baseTypeParams(fn)...).Obj).Obj).Names(c.exportId(recorderIdent)).Obj,
	)).TypeParams(c.typeParams()).Obj).Decs(genDeclDecf(
		"// %s routes recorded function calls to the %s moq",
		structName, mName)).Obj
}

func (c *Converter) anyParamsStruct(prefix string) *dst.GenDecl {
	structName := fmt.Sprintf(double, prefix, anyParamsIdent)
	return TypeDecl(TypeSpec(structName).Type(Struct(
		Field(Star(c.genericExpr(Id(fmt.Sprintf(double, prefix, recorderSuffix)), clone)).Obj).
			Names(c.exportId(recorderIdent)).Obj)).TypeParams(c.typeParams()).Obj).
		Decs(genDeclDecf("// %s isolates the any params functions of the %s type",
			structName, c.typeName())).Obj
}

func (c *Converter) defineAdaptors() []dst.Stmt {
	var stmts []dst.Stmt
	for n, fn := range c.typ.Funcs {
		stmts = append(stmts, Assign(Idf(unnamed, adaptorIdent, n+1)).Tok(token.DEFINE).
			Rhs(Un(token.AND, Comp(c.genericExpr(Idf(double, c.typePrefix(fn), adaptorIdent), clone)).Obj)).Obj)
	}
	return stmts
}

func (c *Converter) updateAdaptors() []dst.Stmt {
	var stmts []dst.Stmt
	for n := range c.typ.Funcs {
		stmts = append(stmts, Assign(Sel(Idf(unnamed, adaptorIdent, n+1)).
			Dot(c.exportId(moqIdent)).Obj).Tok(token.ASSIGN).Rhs(Id(moqReceiverIdent)).Obj)
	}
	return stmts
}

func (c *Converter) newMoqs() []dst.Expr {
	var elts []dst.Expr
	for n, fn := range c.typ.Funcs {
		name := c.exportId(moqIdent)
		if c.isInterface() {
			name = c.exportId(fmt.Sprintf(double, moqIdent, fn.Name))
		}
		elts = append(elts, Key(name).Value(Call(
			Index(IdPath(newMoqFnName, implPkg)).Sub(c.baseTypeParams(fn)...).Obj).
			Args(Id(sceneIdent), Id(fmt.Sprintf(unnamed, adaptorIdent, n+1)), Id(configIdent)).Obj).
			Decs(kvExprDec(dst.NewLine)).Obj)
	}
	return elts
}

func (c *Converter) runtimeValues() []dst.Expr {
	var vals []dst.Expr
	kvDec := dst.NewLine
	for _, fn := range c.typ.Funcs {
		subVals := c.paramIndexingFnValues(fn)
		if c.isInterface() {
			subVals = []dst.Expr{Key(Id(fn.Name)).Value(Comp(c.exportId(
				fmt.Sprintf(double, c.typePrefix(fn), paramIndexingIdent))).
				Elts(subVals...).Obj).Decs(kvExprDec(kvDec)).Obj}
		}
		vals = append(vals, subVals...)
		kvDec = dst.None
	}

	return []dst.Expr{Key(c.exportId(parameterIndexingIdent)).
		Value(Comp(c.paramIndexingStruct()).Elts(vals...).Obj).Obj}
}

func (c *Converter) paramIndexingFnValues(fn Func) []dst.Expr {
	var vals []dst.Expr
	kvDec := dst.NewLine
	count := 1
	for _, f := range fn.FuncType.Params.List {
		typ := c.resolveExpr(f.Type, fn.ParentType)

		for _, name := range f.Names {
			vals = append(vals, c.paramIndexingValue(
				typ, validName(name.Name, paramPrefix, count), kvDec))
			count++
			kvDec = dst.None
		}

		if len(f.Names) == 0 {
			vals = append(vals, c.paramIndexingValue(
				typ, fmt.Sprintf(unnamed, paramPrefix, count), kvDec))
			count++
			kvDec = dst.None
		}
	}

	return vals
}

func (c *Converter) paramIndexingValue(typ dst.Expr, name string, kvDec dst.SpaceType) *dst.KeyValueExpr {
	comp, err := c.typeCache.IsDefaultComparable(typ, c.typ.TypeInfo)
	if err != nil {
		if c.err == nil {
			c.err = err
		}
		return nil
	}

	val := paramIndexByValueIdent
	if !comp {
		val = paramIndexByHashIdent
	}

	return Key(c.exportId(name)).Value(c.idPath(val, moqPkg)).Decs(kvExprDec(kvDec)).Obj
}

func (c *Converter) paramIndexingStruct() dst.Expr {
	if !c.isInterface() {
		return c.exportId(fmt.Sprintf(double, c.moqName(), paramIndexingIdent))
	}
	var piFields []*dst.Field
	for _, fn := range c.typ.Funcs {
		piFields = append(piFields, Field(Idf(
			double, c.typePrefix(fn), paramIndexingIdent)).Names(Id(fn.Name)).Obj)
	}

	return Struct(piFields...)
}

func (c *Converter) mockFunc(typePrefix string, fn Func) []dst.Stmt {
	moqSel := Sel(Id(moqReceiverIdent)).Dot(c.exportId(moqIdent)).Obj
	if c.isInterface() {
		moqSel = Sel(Sel(Id(moqReceiverIdent)).Dot(c.exportId(moqIdent)).Obj).
			Dot(c.exportId(fmt.Sprintf(double, moqIdent, fn.Name))).Obj
	}

	stmts := []dst.Stmt{
		c.helperCallExpr(moqSel),
		Assign(Id(paramsIdent)).
			Tok(token.DEFINE).
			Rhs(Comp(c.genericExpr(Idf(double, typePrefix, paramsIdent), clone)).
				Elts(c.passthroughKeyValues(
					fn.FuncType.Params, paramsIdent, "", fn.ParentType)...).Obj).
			Decs(AssignDecs(dst.None).After(dst.EmptyLine).Obj).Obj,
	}

	count := 1
	if fn.FuncType.Results != nil {
		for _, res := range fn.FuncType.Results.List {
			for range res.Names {
				stmts = append(stmts, Var(Value(c.resolveExpr(cloneExpr(res.Type), fn.ParentType)).
					Names(Idf(unnamed, resultPrefix, count)).Obj))
				count++
			}
			if len(res.Names) == 0 {
				stmts = append(stmts, Var(Value(c.resolveExpr(cloneExpr(res.Type), fn.ParentType)).
					Names(Idf(unnamed, resultPrefix, count)).Obj))
				count++
			}
		}
	}
	fnCall := Call(Sel(cloneExpr(moqSel)).Dot(Id(functionFnName)).Obj).Args(Id(paramsIdent)).Obj
	if fn.FuncType.Results == nil {
		stmts = append(stmts, Expr(fnCall).Obj)
		return stmts
	}

	stmts = append(stmts, If(Bin(Id(resultIdent)).Op(token.NEQ).Y(Id(nilIdent)).Obj).
		Init(Assign(Id(resultIdent)).Tok(token.DEFINE).Rhs(fnCall).Obj).
		Body(c.assignResult(fn.FuncType.Results)...).Obj)

	var rets []dst.Expr
	for n := 1; n < count; n++ {
		rets = append(rets, Idf(unnamed, resultPrefix, n))
	}
	stmts = append(stmts, Return(rets...))

	return stmts
}

func (c *Converter) assignResult(resFL *dst.FieldList) []dst.Stmt {
	var assigns []dst.Stmt
	if resFL != nil {
		count := 1
		for _, result := range resFL.List {
			for _, name := range result.Names {
				rName := fmt.Sprintf(unnamed, resultPrefix, count)
				vName := validName(name.Name, resultPrefix, count)
				assigns = append(assigns, Assign(Id(rName)).
					Tok(token.ASSIGN).
					Rhs(Sel(Id(resultIdent)).Dot(c.exportId(vName)).Obj).Obj)
				count++
			}

			if len(result.Names) == 0 {
				rName := fmt.Sprintf(unnamed, resultPrefix, count)
				assigns = append(assigns, Assign(Id(rName)).
					Tok(token.ASSIGN).
					Rhs(Sel(Id(resultIdent)).Dot(c.exportId(rName)).Obj).Obj)
				count++
			}
		}
	}
	return assigns
}

func (c *Converter) recorderFn(fn Func) *dst.FuncDecl {
	mName := c.moqName()

	recvType := mName
	fnName := onCallFnName
	fnRecName := fmt.Sprintf(double, mName, recorderSuffix)
	if c.isInterface() {
		recvType = fmt.Sprintf(double, mName, recorderIdent)
		fnName = fn.Name
		fnRecName = fmt.Sprintf(triple, mName, fn.Name, recorderSuffix)
	}

	return Fn(c.export(fnName)).
		Recv(Field(Star(c.genericExpr(Id(recvType), clone)).Obj).Names(Id(moqReceiverIdent)).Obj).
		ParamList(c.cloneAndNameUnnamed(paramPrefix, fn.FuncType.Params, fn.ParentType)).
		Results(Field(Star(c.genericExpr(Id(fnRecName), clone)).Obj).Obj).
		Body(c.recorderFnInterfaceBody(fnRecName, c.typePrefix(fn), fn)...).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) recorderFnInterfaceBody(fnRecName, typePrefix string, fn Func) []dst.Stmt {
	var onCallSel dst.Expr = Sel(Id(moqReceiverIdent)).Dot(c.exportId(moqIdent)).Obj
	if c.isInterface() {
		onCallSel = Sel(onCallSel).Dot(c.exportId(fmt.Sprintf(double, moqIdent, fn.Name))).Obj
	}
	return []dst.Stmt{Return(Un(token.AND, Comp(c.genericExpr(Id(fnRecName), clone)).
		Elts(Key(c.exportId(recorderIdent)).Value(Call(Sel(onCallSel).
			Dot(Id(titler.String(onCallFnName))).Obj).
			Args(Comp(c.genericExpr(Idf(double, typePrefix, paramsIdent), clone)).Elts(
				c.passthroughKeyValues(fn.FuncType.Params, paramsIdent, "", fn.ParentType)...).Obj,
			).Obj).Decs(kvExprDec(dst.None)).Obj).Decs(litDec()).Obj,
	))}
}

func (c *Converter) anyParamFns(fn Func) []dst.Decl {
	mName := c.moqName()

	fnRecName := fmt.Sprintf(double, mName, recorderSuffix)
	anyParamsName := fmt.Sprintf(double, mName, anyParamsIdent)
	if c.isInterface() {
		fnRecName = fmt.Sprintf(triple, mName, fn.Name, recorderSuffix)
		anyParamsName = fmt.Sprintf(triple, mName, fn.Name, anyParamsIdent)
	}

	decls := []dst.Decl{c.anyRecorderFn(anyParamsName, fnRecName)}
	count := 1
	for _, param := range fn.FuncType.Params.List {
		for _, name := range param.Names {
			decls = append(decls, c.anyParamFn(anyParamsName, fnRecName,
				validName(name.Name, paramPrefix, count), count))
			count++
		}

		if len(param.Names) == 0 {
			pName := fmt.Sprintf(unnamed, paramPrefix, count)
			decls = append(decls, c.anyParamFn(anyParamsName, fnRecName, pName, count))
			count++
		}
	}
	return decls
}

func (c *Converter) anyRecorderFn(anyParamsName, fnRecName string) *dst.FuncDecl {
	return Fn(c.export("any")).
		Recv(Field(Star(c.genericExpr(Id(fnRecName), clone)).Obj).
			Names(Id(recorderReceiverIdent)).Obj).
		Results(Field(Star(c.genericExpr(Id(anyParamsName), clone)).Obj).Obj).
		Body(
			c.helperCallExpr(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(recorderIdent)).Obj).Dot(Id(titler.String(moqIdent))).Obj),
			If(Un(token.NOT, Call(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(recorderIdent)).Obj).Dot(Id(isAnyPermittedFnName)).Obj).
				Args(LitBool(c.isExported)).Obj)).Body(
				Return(Id(nilIdent)),
			).Obj,
			Return(Un(token.AND, Comp(cloneExpr(c.genericExpr(Id(anyParamsName), clone))).
				Elts(Key(c.exportId(recorderIdent)).Value(Id(recorderReceiverIdent)).Obj).Obj)),
		).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) anyParamFn(anyParamsName, fnRecName, pName string, paramPos int) *dst.FuncDecl {
	return Fn(c.export(pName)).
		Recv(Field(Star(c.genericExpr(Id(anyParamsName), clone)).Obj).Names(Id(anyParamsReceiverIdent)).Obj).
		Results(Field(Star(c.genericExpr(Id(fnRecName), clone)).Obj).Obj).
		Body(
			Expr(Call(Sel(Sel(Sel(Id(anyParamsReceiverIdent)).Dot(c.exportId(recorderIdent)).Obj).
				Dot(c.exportId(recorderIdent)).Obj).Dot(Id(anyParamFnName)).Obj).
				Args(LitInt(paramPos)).Obj).Obj,
			Return(Sel(Id(anyParamsReceiverIdent)).Dot(c.exportId(recorderIdent)).Obj),
		).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) recorderSeqFns(fn Func) []dst.Decl {
	return []dst.Decl{
		c.recorderSeqFn(fn, "seq", true),
		c.recorderSeqFn(fn, "noSeq", false),
	}
}

func (c *Converter) recorderSeqFn(fn Func, fnName string, seq bool) *dst.FuncDecl {
	mName := c.moqName()

	fnRecName := fmt.Sprintf(double, mName, recorderSuffix)
	if c.isInterface() {
		fnRecName = fmt.Sprintf(triple, mName, fn.Name, recorderSuffix)
	}

	fnName = c.export(fnName)
	return Fn(fnName).
		Recv(Field(Star(c.genericExpr(Id(fnRecName), clone)).Obj).Names(Id(recorderReceiverIdent)).Obj).
		Results(Field(Star(c.genericExpr(Id(fnRecName), clone)).Obj).Obj).
		Body(
			c.helperCallExpr(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(recorderIdent)).Obj).Dot(Id(titler.String(moqIdent))).Obj),
			If(Un(token.NOT, Call(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(recorderIdent)).Obj).Dot(Id(seqFnName)).Obj).
				Args(LitBool(seq), LitString(fnName), LitBool(c.isExported)).Obj)).Body(
				Return(Id(nilIdent)),
			).Obj,
			Return(Id(recorderReceiverIdent)),
		).Decs(stdFuncDec()).Obj
}

func (c *Converter) returnResultsFn(fn Func) *dst.FuncDecl {
	mName := c.moqName()

	fnRecName := fmt.Sprintf(double, mName, recorderSuffix)
	resType := fmt.Sprintf(double, mName, resultsIdent)
	if c.isInterface() {
		fnRecName = fmt.Sprintf(triple, mName, fn.Name, recorderSuffix)
		resType = fmt.Sprintf(triple, mName, fn.Name, resultsIdent)
	}

	return Fn(c.export(returnResultsFnName)).
		Recv(Field(Star(c.genericExpr(Id(fnRecName), clone)).Obj).
			Names(Id(recorderReceiverIdent)).Obj).
		ParamList(c.cloneAndNameUnnamed(resultPrefix, fn.FuncType.Results, fn.ParentType)).
		Results(Field(Star(c.genericExpr(Id(fnRecName), clone)).Obj).Obj).
		Body(
			c.helperCallExpr(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(recorderIdent)).Obj).Dot(Id(titler.String(moqIdent))).Obj),
			Expr(Call(Sel(Sel(Id(recorderReceiverIdent)).Dot(c.exportId(recorderIdent)).Obj).
				Dot(Id(titler.String(returnResultsFnName))).Obj).
				Args(Comp(c.genericExpr(Id(resType), clone)).
					Elts(c.passthroughKeyValues(fn.FuncType.Results,
						resultsIdent, "", fn.ParentType)...).Obj).Obj).Obj,
			Return(Id(recorderReceiverIdent)),
		).Obj
}

func (c *Converter) andDoFn(fn Func) *dst.FuncDecl {
	mName := c.moqName()

	fnRecName := fmt.Sprintf(double, mName, recorderSuffix)
	paramsType := fmt.Sprintf(double, mName, paramsIdent)
	if c.isInterface() {
		fnRecName = fmt.Sprintf(triple, mName, fn.Name, recorderSuffix)
		paramsType = fmt.Sprintf(triple, mName, fn.Name, paramsIdent)
	}
	typePrefix := c.typePrefix(fn)
	fnName := fmt.Sprintf(double, typePrefix, doFnIdent)

	ell := false
	if fn.FuncType.Params != nil && len(fn.FuncType.Params.List) > 0 {
		last := fn.FuncType.Params.List[len(fn.FuncType.Params.List)-1]
		_, ell = last.Type.(*dst.Ellipsis)
	}

	return Fn(c.export(andDoFnName)).
		Recv(Field(Star(c.genericExpr(Id(fnRecName), clone)).Obj).
			Names(Id(recorderReceiverIdent)).Obj).
		Params(Field(c.genericExpr(Id(fnName), clone)).Names(Id(fnFnName)).Obj).
		Results(Field(Star(c.genericExpr(Id(fnRecName), clone)).Obj).Obj).Body(
		c.helperCallExpr(Sel(Sel(Id(recorderReceiverIdent)).
			Dot(c.exportId(recorderIdent)).Obj).Dot(Id(titler.String(moqIdent))).Obj),
		If(Un(token.NOT, Call(Sel(Sel(Id(recorderReceiverIdent)).
			Dot(c.exportId(recorderIdent)).Obj).Dot(Id(titler.String(andDoFnName))).Obj).
			Args(FnLit(FnType(FieldList(Field(c.genericExpr(Id(paramsType), clone)).
				Names(Id(paramsIdent)).Obj)).Obj).Body(
				Expr(Call(Id(fnFnName)).Args(c.passthroughValues(
					fn.FuncType.Params, paramPrefix, Id(paramsIdent))...).Ellipsis(ell).Obj).
					Decs(ExprDecs(dst.NewLine).Obj).Obj,
			).Obj, LitBool(c.isExported)).Obj)).
			Body(Return(Id(nilIdent))).Obj,
		Return(Id(recorderReceiverIdent)),
	).Decs(stdFuncDec()).Obj
}

func (c *Converter) doReturnResultsFn(fn Func) *dst.FuncDecl {
	mName := c.moqName()

	fnRecName := fmt.Sprintf(double, mName, recorderSuffix)
	paramsType := fmt.Sprintf(double, mName, paramsIdent)
	resType := fmt.Sprintf(double, mName, resultsIdent)
	fnType := fmt.Sprintf(double, mName, doReturnFnIdent)
	if c.isInterface() {
		fnRecName = fmt.Sprintf(triple, mName, fn.Name, recorderSuffix)
		paramsType = fmt.Sprintf(triple, mName, fn.Name, paramsIdent)
		resType = fmt.Sprintf(triple, mName, fn.Name, resultsIdent)
		fnType = fmt.Sprintf(triple, mName, fn.Name, doReturnFnIdent)
	}

	return Fn(c.export(doReturnResultsFnName)).
		Recv(Field(Star(c.genericExpr(Id(fnRecName), clone)).Obj).
			Names(Id(recorderReceiverIdent)).Obj).
		Params(Field(c.genericExpr(Id(fnType), clone)).Names(Id(fnFnName)).Obj).
		Results(Field(Star(c.genericExpr(Id(fnRecName), clone)).Obj).Obj).
		Body(
			c.helperCallExpr(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(recorderIdent)).Obj).Dot(Id(titler.String(moqIdent))).Obj),
			Expr(Call(Sel(Sel(Id(recorderReceiverIdent)).
				Dot(c.exportId(recorderIdent)).Obj).
				Dot(Id(titler.String(doReturnResultsFnName))).Obj).
				Args(FnLit(FnType(FieldList(Field(c.genericExpr(Id(paramsType), clone)).
					Names(Id(paramsIdent)).Obj)).
					Results(FieldList(Field(Star(c.genericExpr(Id(resType), clone)).Obj).Obj)).Obj).Body(
					c.doReturnResultsFnLitBody(fn)...,
				).Obj).Obj).Obj,
			Return(Id(recorderReceiverIdent)),
		).Obj
}

func (c *Converter) doReturnResultsFnLitBody(fn Func) []dst.Stmt {
	mName := c.moqName()
	resType := fmt.Sprintf(double, mName, resultsIdent)
	if c.isInterface() {
		resType = fmt.Sprintf(triple, mName, fn.Name, resultsIdent)
	}

	ell := false
	if fn.FuncType.Params != nil && len(fn.FuncType.Params.List) > 0 {
		last := fn.FuncType.Params.List[len(fn.FuncType.Params.List)-1]
		_, ell = last.Type.(*dst.Ellipsis)
	}

	call := Call(Id(fnFnName)).
		Args(c.passthroughValues(fn.FuncType.Params, paramPrefix, Id(paramsIdent))...).Ellipsis(ell).Obj

	var fnCall dst.Stmt = Assign(c.passthroughValues(fn.FuncType.Results, resultPrefix, nil)...).
		Tok(token.DEFINE).Rhs(call).Decs(AssignDecs(dst.NewLine).Obj).Obj
	if fn.FuncType.Results == nil {
		fnCall = Expr(call).Decs(ExprDecs(dst.NewLine).Obj).Obj
	}

	return []dst.Stmt{
		fnCall,
		Return(Un(token.AND, Comp(c.genericExpr(Id(resType), clone)).
			Elts(c.passthroughKeyValues(
				fn.FuncType.Results, resultsIdent, "", fn.ParentType)...).Obj)),
	}
}

func (c *Converter) recorderRepeatFn(fn Func) *dst.FuncDecl {
	mName := c.moqName()

	fnRecName := fmt.Sprintf(double, mName, recorderSuffix)
	if c.isInterface() {
		fnRecName = fmt.Sprintf(triple, mName, fn.Name, recorderSuffix)
	}

	return Fn(c.export(repeatFnName)).
		Recv(Field(Star(c.genericExpr(Id(fnRecName), clone)).Obj).
			Names(Id(recorderReceiverIdent)).Obj).
		Params(Field(Ellipsis(c.idPath(repeaterType, moqPkg))).Names(Id(repeatersIdent)).Obj).
		Results(Field(Star(c.genericExpr(Id(fnRecName), clone)).Obj).Obj).Body(
		c.helperCallExpr(Sel(Sel(Id(recorderReceiverIdent)).
			Dot(c.exportId(recorderIdent)).Obj).Dot(Id(titler.String(moqIdent))).Obj),
		If(Un(token.NOT, Call(Sel(Sel(Id(recorderReceiverIdent)).
			Dot(c.exportId(recorderIdent)).Obj).Dot(Id(titler.String(repeatFnName))).Obj).
			Args(Id(repeatersIdent), LitBool(c.isExported)).Obj)).
			Body(Return(Id(nilIdent))).Obj,
		Return(Id(recorderReceiverIdent)),
	).Decs(stdFuncDec()).Obj
}

func (c *Converter) prettyParamsFn(fn Func) *dst.FuncDecl {
	mName := c.moqName()
	rType := fmt.Sprintf(double, mName, adaptorIdent)
	params := fmt.Sprintf(double, mName, paramsIdent)
	sfmt := c.typeName() + "("
	if c.isInterface() {
		rType = fmt.Sprintf(triple, mName, fn.Name, adaptorIdent)
		params = fmt.Sprintf(triple, mName, fn.Name, paramsIdent)
		sfmt = fn.Name + "("
	}
	var pExprs []dst.Expr
	count := 1
	for _, param := range fn.FuncType.Params.List {
		for _, name := range param.Names {
			sfmt += "%#v, "
			vName := validName(name.Name, paramPrefix, count)
			pExpr := Sel(Id(paramsIdent)).Dot(c.exportId(vName)).Obj
			pExprs = append(pExprs, c.prettyParam(param, pExpr))
			count++
		}

		if len(param.Names) == 0 {
			sfmt += "%#v, "
			pExpr := Sel(Id(paramsIdent)).
				Dot(c.exportId(fmt.Sprintf(unnamed, paramPrefix, count))).Obj
			pExprs = append(pExprs, c.prettyParam(param, pExpr))
			count++
		}
	}
	if count > 1 {
		sfmt = sfmt[0 : len(sfmt)-2]
	}
	sfmt += ")"
	pExprs = append([]dst.Expr{LitString(sfmt)}, pExprs...)
	return Fn(titler.String(prettyParamsFnName)).
		Recv(Field(Star(c.genericExpr(Id(rType), clone)).Obj).Obj).
		Params(Field(c.genericExpr(Id(params), clone)).Names(Id(paramsIdent)).Obj).
		Results(Field(Id("string")).Obj).
		Body(Return(
			Call(IdPath("Sprintf", "fmt")).Args(pExprs...).Obj)).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) prettyParam(param *dst.Field, expr *dst.SelectorExpr) dst.Expr {
	if _, ok := param.Type.(*dst.FuncType); !ok {
		return expr
	}

	return Call(IdPath("FnString", moqPkg)).Args(expr).Obj
}

func (c *Converter) paramsKeyFn(fn Func) *dst.FuncDecl {
	mName := c.moqName()
	rType := fmt.Sprintf(double, mName, adaptorIdent)
	params := fmt.Sprintf(double, mName, paramsIdent)
	paramsKey := fmt.Sprintf(double, mName, paramsKeyIdent)
	hSel := moqIdent
	if c.isInterface() {
		rType = fmt.Sprintf(triple, mName, fn.Name, adaptorIdent)
		params = fmt.Sprintf(triple, mName, fn.Name, paramsIdent)
		paramsKey = fmt.Sprintf(triple, mName, fn.Name, paramsKeyIdent)
		hSel = fmt.Sprintf(double, moqIdent, fn.Name)
	}

	stmts := []dst.Stmt{
		c.helperCallExpr(Sel(Sel(Id(adaptorReceiverIdent)).
			Dot(c.exportId(moqIdent)).Obj).Dot(c.exportId(hSel)).Obj),
	}
	count := 1
	for _, param := range fn.FuncType.Params.List {
		typ := c.resolveExpr(param.Type, fn.ParentType)

		for _, name := range param.Names {
			stmts = append(stmts,
				c.mockFuncFindResultsParam(fn, name.Name, count, typ)...)
			count++
		}

		if len(param.Names) == 0 {
			stmts = append(stmts, c.mockFuncFindResultsParam(
				fn, fmt.Sprintf(unnamed, paramPrefix, count), count, typ)...)
			count++
		}
	}

	stmts = append(stmts, Return(Comp(c.genericExpr(Id(paramsKey), clone)).Elts(
		Key(c.exportId(paramsIdent)).Value(Comp(
			c.methodStruct(paramsKeyIdent, fn.FuncType.Params, fn.ParentType)).
			Elts(c.passthroughKeyValues(
				fn.FuncType.Params, paramsKeyIdent, usedSuffix, fn.ParentType)...).Obj).
			Decs(kvExprDec(dst.NewLine)).Obj,
		Key(c.exportId(hashesIdent)).Value(Comp(
			c.methodStruct(hashesIdent, fn.FuncType.Params, fn.ParentType)).
			Elts(c.passthroughKeyValues(
				fn.FuncType.Params, hashesIdent, usedHashSuffix, fn.ParentType)...).Obj).
			Decs(kvExprDec(dst.NewLine)).Obj).Obj))

	return Fn(titler.String(paramsKeyFnName)).
		Recv(Field(Star(c.genericExpr(Id(rType), clone)).Obj).Names(Id(adaptorReceiverIdent)).Obj).
		Params(Field(c.genericExpr(Id(params), clone)).Names(Id(paramsIdent)).Obj,
			Field(Id("uint64")).Names(Id(anyParamsIdent)).Obj).
		Results(Field(c.genericExpr(Id(paramsKey), clone)).Obj).
		Body(stmts...).
		Decs(stdFuncDec()).Obj
}

func (c *Converter) mockFuncFindResultsParam(
	fn Func, pName string, paramPos int, typ dst.Expr,
) []dst.Stmt {
	pSel := Sel(Sel(Sel(Sel(Id(adaptorReceiverIdent)).
		Dot(c.exportId(moqIdent)).Obj).
		Dot(c.exportId(runtimeIdent)).Obj).
		Dot(c.exportId(parameterIndexingIdent)).Obj)
	sSel := moqIdent
	if c.isInterface() {
		pSel = Sel(pSel.Dot(Id(fn.Name)).Obj)
		sSel = fmt.Sprintf(double, moqIdent, fn.Name)
	}
	comp, err := c.typeCache.IsComparable(typ, c.typ.TypeInfo)
	if err != nil {
		if c.err == nil {
			c.err = err
		}
		return nil
	}

	vName := validName(pName, paramPrefix, paramPos)
	pUsed := fmt.Sprintf("%s%s", vName, usedSuffix)
	hashUsed := fmt.Sprintf("%s%s", vName, usedHashSuffix)

	lhs := []dst.Expr{Id(hashUsed)}
	fnName := hashOnlyParamKeyFnName
	var tExpr dst.Expr = Sel(Sel(Sel(Sel(Id(adaptorReceiverIdent)).Dot(c.exportId(moqIdent)).Obj).
		Dot(c.exportId(sSel)).Obj).Dot(Id(titler.String(sceneIdent))).Obj).Dot(Id(tType)).Obj
	var pNameLit dst.Expr = LitString(vName)
	if comp {
		lhs = append([]dst.Expr{Id(pUsed)}, lhs...)
		fnName = paramKeyFnName
		tExpr = nil
		pNameLit = nil
	}

	return []dst.Stmt{
		Assign(lhs...).Tok(token.DEFINE).Rhs(Call(IdPath(fnName, implPkg)).Args(
			tExpr,
			Sel(Id(paramsIdent)).Dot(c.exportId(vName)).Decs(SelDecs(dst.NewLine, dst.None)).Obj,
			pNameLit,
			LitInt(paramPos),
			pSel.Dot(c.exportId(vName)).Obj,
			Id(anyParamsIdent)).Obj).Obj,
	}
}

func (c *Converter) typeParams() *dst.FieldList {
	return c.cloneFieldList(c.typ.TypeInfo.Type.TypeParams, false, c.typ.TypeInfo)
}

type genericExprMode int

const (
	// clone clones the generic expression as-is
	clone genericExprMode = iota
	// typeAssertSafe makes the expression safe for use in a type assertion
	// (omitNames is implied)
	typeAssertSafe
)

func (c *Converter) genericExpr(in dst.Expr, mode genericExprMode) dst.Expr {
	if c.typ.TypeInfo.Type.TypeParams == nil || len(c.typ.TypeInfo.Type.TypeParams.List) == 0 {
		return in
	}

	var subs []dst.Expr
	for _, p := range c.typ.TypeInfo.Type.TypeParams.List {
		for _, n := range p.Names {
			var sub dst.Expr
			if mode == clone {
				sub = Id(n.Name)
			} else {
				sub = cloneExpr(p.Type)
				if mode == typeAssertSafe {
					if un, ok := sub.(*dst.UnaryExpr); ok && un.Op == token.TILDE {
						sub = un.X
					}
				}
			}
			subs = append(subs, c.resolveExpr(sub, c.typ.TypeInfo))
		}
	}

	return Index(in).Sub(subs...).Obj
}

func (c *Converter) passthroughKeyValues(fl *dst.FieldList, label, valSuffix string, parentType TypeInfo) []dst.Expr {
	if fl == nil {
		return nil
	}

	unnamedPrefix, dropNonComparable := labelDirection(label)
	var elts []dst.Expr
	beforeDec := dst.NewLine
	count := 1
	for _, field := range fl.List {
		comp, err := c.typeCache.IsComparable(c.resolveExpr(field.Type, parentType), c.typ.TypeInfo)
		if err != nil {
			if c.err == nil {
				c.err = err
			}
			return nil
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

		if len(field.Names) == 0 {
			if comp || !dropNonComparable {
				pName := fmt.Sprintf(unnamed, unnamedPrefix, count)
				elts = append(elts, Key(c.exportId(pName)).Value(
					c.passthroughValue(Id(pName+valSuffix), false, nil)).
					Decs(kvExprDec(beforeDec)).Obj)
				beforeDec = dst.None
			}
			count++
		}
	}

	return elts
}

func (c *Converter) passthroughValues(fl *dst.FieldList, prefix string, sel dst.Expr) []dst.Expr {
	if fl == nil {
		return nil
	}

	pass := func(name string) dst.Expr {
		if sel != nil {
			// If we are selecting off of something, we need to use an
			// exported name
			name = c.export(name)
		}
		return c.passthroughValue(Id(name), false, sel)
	}

	var elts []dst.Expr
	count := 1
	for _, field := range fl.List {
		for _, name := range field.Names {
			elts = append(elts, pass(validName(name.Name, prefix, count)))
			count++
		}

		if len(field.Names) == 0 {
			elts = append(elts, pass(fmt.Sprintf(unnamed, prefix, count)))
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

func (c *Converter) cloneAndNameUnnamed(prefix string, fieldList *dst.FieldList, parentType TypeInfo) *dst.FieldList {
	fieldList = c.cloneFieldList(fieldList, false, parentType)
	if fieldList != nil {
		count := 1
		for _, f := range fieldList.List {
			for n, name := range f.Names {
				f.Names[n] = Id(validName(name.Name, prefix, count))
				count++
			}
			if len(f.Names) == 0 {
				f.Names = []*dst.Ident{Idf(unnamed, prefix, count)}
				count++
			}
		}
	}
	return fieldList
}

func validName(name, prefix string, count int) string {
	if _, ok := invalidNames[name]; ok {
		name = fmt.Sprintf(unnamed, prefix, count)
	}
	return name
}

func (c *Converter) moqName() string {
	return c.export(moqIdent + titler.String(c.typeName()))
}

func (c *Converter) typeName() string {
	typ := c.typ.TypeInfo.Type.Name.Name
	if strings.HasPrefix(typ, GenTypeSuffix) {
		typ = typ[:len(typ)-len(GenTypeSuffix)]
	}

	return typ
}

func (c *Converter) isInterface() bool {
	if _, ok := c.typ.TypeInfo.Type.Type.(*dst.InterfaceType); ok {
		return true
	}
	return false
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
		return IdPath(name, c.typ.TypeInfo.PkgPath)
	case c.typ.OutPkgPath:
		return Id(name)
	default:
		return IdPath(name, path)
	}
}

func (c *Converter) helperCallExpr(selector dst.Expr) dst.Stmt {
	return Expr(Call(Sel(Sel(Sel(selector).
		Dot(Id(titler.String(sceneIdent))).Obj).
		Dot(Id(tType)).Obj).
		Dot(Id(helperFnName)).Obj).Obj).
		Decs(ExprDecs(dst.NewLine).Obj).Obj
}

func (c *Converter) selExport(x dst.Expr, sel string) dst.Expr {
	if x == nil {
		return Id(sel)
	}

	switch v := x.(type) {
	case *dst.SelectorExpr:
		return Sel(cloneSelect(v)).Dot(c.exportId(sel)).Obj
	case *dst.Ident:
		return Sel(cloneExpr(v)).Dot(c.exportId(sel)).Obj
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

func (c *Converter) cloneFieldList(fl *dst.FieldList, removeNames bool, parentType TypeInfo) *dst.FieldList {
	if fl != nil {
		//nolint:forcetypeassert // if dst.Clone returns a different type, panic
		fl = dst.Clone(fl).(*dst.FieldList)
		c.resolveFieldList(fl, parentType)
		if removeNames {
			var fields []*dst.Field
			for _, field := range fl.List {
				for range field.Names {
					fields = append(fields, Field(cloneExpr(field.Type)).Obj)
				}
				if len(field.Names) == 0 {
					fields = append(fields, field)
				}
				field.Names = nil
			}
			fl.List = fields
		}
	}
	return fl
}

func (c *Converter) resolveFieldList(fl *dst.FieldList, parentType TypeInfo) {
	if fl == nil {
		return
	}
	for _, field := range fl.List {
		if id, ok := field.Type.(*dst.Ident); ok {
			if id.Path == "builtin" {
				// When mocking errors, the string return value is reported
				// in the builtin package
				id.Path = ""
			}
		}

		field.Type = c.resolveExpr(field.Type, parentType)
	}
}

func (c *Converter) resolveExpr(expr dst.Expr, parentType TypeInfo) dst.Expr {
	switch e := expr.(type) {
	case *dst.ArrayType:
		e.Elt = c.resolveExpr(e.Elt, parentType)
		e.Len = c.resolveExpr(e.Len, parentType)
	case *dst.ChanType:
		e.Value = c.resolveExpr(e.Value, parentType)
	case *dst.Ellipsis:
		e.Elt = c.resolveExpr(e.Elt, parentType)
	case *dst.Ident:
		if e.Path == "" && c.typ.TypeInfo.Type.TypeParams != nil {
			// Priority is always highest for type params
			for _, f := range c.typ.TypeInfo.Type.TypeParams.List {
				for _, n := range f.Names {
					if n.Name == e.Name {
						// So just use the name as is with no path
						return e
					}
				}
			}
		}
		typ, err := c.typeCache.Type(*e, parentType.PkgPath, false)
		if err != nil {
			if errors.Is(err, ErrTypeNotFound) {
				ppkg := e.Path
				if ppkg == "" {
					ppkg = parentType.PkgPath
				}
				id := IdPath(e.Name, ppkg)
				logs.Warnf("Passing through unknown type %s as %s", e.String(), id.String())
				return id
			}
			if c.err == nil {
				c.err = err
			}
			return nil
		}
		return IdPath(typ.Type.Name.Name, typ.PkgPath)
	case *dst.IndexExpr:
		e.X = c.resolveExpr(e.X, parentType)
		e.Index = c.resolveExpr(e.Index, parentType)
	case *dst.InterfaceType:
		for _, method := range e.Methods.List {
			method.Type = c.resolveExpr(method.Type, parentType)
		}
	case *dst.FuncType:
		c.resolveFieldList(e.Params, parentType)
		c.resolveFieldList(e.Results, parentType)
	case *dst.MapType:
		e.Key = c.resolveExpr(e.Key, parentType)
		e.Value = c.resolveExpr(e.Value, parentType)
	case *dst.StarExpr:
		e.X = c.resolveExpr(e.X, parentType)
	case *dst.StructType:
		for n, f := range e.Fields.List {
			e.Fields.List[n].Type = c.resolveExpr(f.Type, parentType)
		}
	}

	return expr
}

func cloneSelect(sel *dst.SelectorExpr) dst.Expr {
	//nolint:forcetypeassert // if dst.Clone returns a different type, panic
	return dst.Clone(sel).(*dst.SelectorExpr)
}

func cloneExpr(expr dst.Expr) dst.Expr {
	//nolint:forcetypeassert // if dst.Clone returns a different type, panic
	return dst.Clone(expr).(dst.Expr)
}

func genDeclDecf(format string, a ...interface{}) dst.GenDeclDecorations {
	return dst.GenDeclDecorations{
		NodeDecs: NodeDecsf(format, a...),
	}
}

func fnDeclDecf(format string, a ...interface{}) dst.FuncDeclDecorations {
	return dst.FuncDeclDecorations{
		NodeDecs: NodeDecsf(format, a...),
	}
}

func litDec() dst.CompositeLitDecorations {
	return dst.CompositeLitDecorations{Lbrace: []string{"\n"}}
}

func kvExprDec(before dst.SpaceType) dst.KeyValueExprDecorations {
	return KeyValueDecs(before).After(dst.NewLine).Obj
}
