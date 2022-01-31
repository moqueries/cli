package generator

import (
	"fmt"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/goast"

	"github.com/myshkin5/moqueries/ast"
	"github.com/myshkin5/moqueries/logs"
)

const (
	headerComment = "// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!"
)

//go:generate moqueries --destination moq_converterer_test.go Converterer

// Converterer is the interface used by MoqGenerator to invoke a Converter
type Converterer interface {
	BaseStruct(typeSpec *dst.TypeSpec, funcs []Func) (structDecl *dst.GenDecl)
	IsolationStruct(typeName, suffix string) (structDecl *dst.GenDecl)
	MethodStructs(typeSpec *dst.TypeSpec, fn Func) (structDecls []dst.Decl, err error)
	NewFunc(typeSpec *dst.TypeSpec, funcs []Func) (funcDecl *dst.FuncDecl)
	IsolationAccessor(typeName, suffix, fnName string) (funcDecl *dst.FuncDecl)
	FuncClosure(typeName, pkgPath string, fn Func) (funcDecl *dst.FuncDecl)
	MockMethod(typeName string, fn Func) (funcDecl *dst.FuncDecl)
	RecorderMethods(typeName string, fn Func) (funcDecls []dst.Decl)
	ResetMethod(typeSpec *dst.TypeSpec, funcs []Func) (funcDecl *dst.FuncDecl)
	AssertMethod(typeSpec *dst.TypeSpec, funcs []Func) (funcDecl *dst.FuncDecl)
}

// MoqGenerator generates moqs
type MoqGenerator struct {
	export    bool
	pkg       string
	dest      string
	typeCache TypeCache
	converter Converterer
}

//go:generate moqueries --destination moq_typecache_test.go TypeCache

// TypeCache defines the interface to the Cache type
type TypeCache interface {
	Type(id dst.Ident, loadTestTypes bool) (*dst.TypeSpec, string, error)
	IsComparable(expr dst.Expr) (bool, error)
	IsDefaultComparable(expr dst.Expr) (bool, error)
}

// New returns a new MoqGenerator
func New(
	export bool,
	pkg, dest string,
	typeCache TypeCache,
	converter Converterer,
) *MoqGenerator {
	return &MoqGenerator{
		export:    export,
		pkg:       pkg,
		dest:      dest,
		typeCache: typeCache,
		converter: converter,
	}
}

// Generate generates moqs for the given types in the given destination package
func (g *MoqGenerator) Generate(inTypes []string, imp string, loadTestTypes bool) (
	*token.FileSet,
	*dst.File,
	error,
) {
	fSet, file := initializeFile(g.defaultPackage())

	var decls []dst.Decl
	for _, inType := range inTypes {
		typeSpec, pkgPath, err := g.loadInType(inType, imp, loadTestTypes)
		if err != nil {
			return nil, nil, err
		}

		funcs, tErr := g.findFuncs(typeSpec)
		if tErr != nil {
			return nil, nil, tErr
		}

		structs, err := g.structs(typeSpec, funcs)
		if err != nil {
			return nil, nil, err
		}
		decls = append(decls, structs...)

		decls = append(decls, g.converter.NewFunc(typeSpec, funcs))

		decls = append(decls, g.methods(typeSpec, pkgPath, funcs)...)

		decls = append(decls, g.converter.ResetMethod(typeSpec, funcs))

		decls = append(decls, g.converter.AssertMethod(typeSpec, funcs))
	}
	file.Decls = decls

	return fSet, file, nil
}

func (g *MoqGenerator) defaultPackage() string {
	pkg := g.pkg
	if pkg == "" {
		abs, err := filepath.Abs(g.dest)
		if err != nil {
			logs.Panicf("Could not get absolute path to destination %s: %#v", g.dest, err)
		}
		dirName := filepath.Base(filepath.Dir(abs))
		pkg = dirName
		if !g.export {
			pkg = pkg + "_test"
		}
	}
	logs.Debugf("Output package: %s", pkg)
	return pkg
}

func initializeFile(pkg string) (*token.FileSet, *dst.File) {
	fSet := token.NewFileSet()

	src := fmt.Sprintf("%s\n\npackage %s\n", headerComment, pkg)
	file, err := decorator.NewDecoratorWithImports(fSet, pkg, goast.New()).Parse(src)
	if err != nil {
		logs.Panic("Could not create decorator", err)
	}

	return fSet, file
}

const testPkgSuffix = "_test"

func (g *MoqGenerator) loadInType(inType, imp string, loadTestTypes bool) (
	*dst.TypeSpec, string, error,
) {
	if strings.HasSuffix(imp, testPkgSuffix) {
		imp = strings.TrimSuffix(imp, testPkgSuffix)
		loadTestTypes = true
	}
	return g.typeCache.Type(*ast.IdPath(inType, imp), loadTestTypes)
}

func (g *MoqGenerator) findFuncs(typeSpec *dst.TypeSpec) ([]Func, error) {
	switch typ := typeSpec.Type.(type) {
	case *dst.InterfaceType:
		return g.loadNestedInterfaces(typ)
	case *dst.FuncType:
		return []Func{
			{
				Params:  typ.Params,
				Results: typ.Results,
			},
		}, nil
	case *dst.Ident:
		var funcs []Func
		return g.loadTypeEquivalent(funcs, typ)
	default:
		logs.Panicf("Unknown type: %v", typeSpec.Type)
		panic("unreachable")
	}
}

func (g *MoqGenerator) loadNestedInterfaces(iType *dst.InterfaceType) ([]Func, error) {
	var funcs []Func

	for _, method := range iType.Methods.List {
		switch typ := method.Type.(type) {
		case *dst.FuncType:
			funcs = append(funcs, Func{
				Name:    method.Names[0].Name,
				Params:  typ.Params,
				Results: typ.Results,
			})
		case *dst.Ident:
			var err error
			funcs, err = g.loadTypeEquivalent(funcs, typ)
			if err != nil {
				return nil, err
			}
		default:
			logs.Panicf("Unknown type in interface method list: %v", method.Type)
		}
	}

	return funcs, nil
}

func (g *MoqGenerator) loadTypeEquivalent(funcs []Func, id *dst.Ident) ([]Func, error) {
	nestedType, _, err := g.typeCache.Type(*id, false)
	if err != nil {
		return nil, err
	}

	newFuncs, err := g.findFuncs(nestedType)
	if err != nil {
		return nil, err
	}
	funcs = append(funcs, newFuncs...)

	return funcs, nil
}

func (g *MoqGenerator) structs(typeSpec *dst.TypeSpec, funcs []Func) ([]dst.Decl, error) {
	decls := []dst.Decl{
		g.converter.BaseStruct(typeSpec, funcs),
		g.converter.IsolationStruct(typeSpec.Name.Name, mockIdent),
	}

	_, iOk := typeSpec.Type.(*dst.InterfaceType)
	_, aOk := typeSpec.Type.(*dst.Ident)
	if iOk || aOk {
		decls = append(decls,
			g.converter.IsolationStruct(typeSpec.Name.Name, recorderIdent))
	}

	for _, fn := range funcs {
		structs, err := g.converter.MethodStructs(typeSpec, fn)
		if err != nil {
			return nil, err
		}
		decls = append(decls, structs...)
	}

	return decls, nil
}

func (g *MoqGenerator) methods(
	typeSpec *dst.TypeSpec, pkgPath string, funcs []Func,
) []dst.Decl {
	var decls []dst.Decl

	switch typeSpec.Type.(type) {
	case *dst.InterfaceType, *dst.Ident:
		decls = append(
			decls, g.converter.IsolationAccessor(
				typeSpec.Name.Name, mockIdent, mockFnName))

		for _, fn := range funcs {
			decls = append(
				decls, g.converter.MockMethod(typeSpec.Name.Name, fn))
		}

		decls = append(
			decls, g.converter.IsolationAccessor(
				typeSpec.Name.Name, recorderIdent, onCallFnName))

		for _, fn := range funcs {
			decls = append(
				decls, g.converter.RecorderMethods(typeSpec.Name.Name, fn)...)
		}
	case *dst.FuncType:
		if len(funcs) != 1 {
			logs.Panicf("Function moqs should have just one function, found: %d",
				len(funcs))
		}

		decls = append(
			decls, g.converter.FuncClosure(
				typeSpec.Name.Name, pkgPath, funcs[0]))

		decls = append(
			decls, g.converter.MockMethod(
				typeSpec.Name.Name, funcs[0]))

		decls = append(
			decls, g.converter.RecorderMethods(
				typeSpec.Name.Name, funcs[0])...)
	default:
		logs.Panicf("Unknown type: %v", typeSpec.Type)
	}

	return decls
}
