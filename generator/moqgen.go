package generator

import (
	"errors"
	"fmt"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/goast"
	"moqueries.org/runtime/logs"

	"moqueries.org/cli/ast"
)

const (
	headerComment = "// Code generated by Moqueries - https://moqueries.org - DO NOT EDIT."
	testPkgSuffix = "_test"
)

var (
	// ErrInvalidConfig is returned when configuration values are invalid or
	// conflict with each other
	ErrInvalidConfig = errors.New("invalid configuration")
	// ErrUnknownFieldType is returned when a param or result list has an
	// unknown expression
	ErrUnknownFieldType = errors.New("unknown field type")
)

// Type describes the type being converted into a mock
type Type struct {
	TypeInfo   ast.TypeInfo
	Funcs      []Func
	OutPkgPath string
	Reduced    bool
}

//go:generate moqueries GetwdFunc

// GetwdFunc is the signature of os.Getwd
type GetwdFunc func() (string, error)

//go:generate moqueries NewConverterFunc

// NewConverterFunc creates a new converter for a specific type
type NewConverterFunc func(typ Type, export bool) Converterer

//go:generate moqueries Converterer

// Converterer is the interface used by MoqGenerator to invoke a Converter
type Converterer interface {
	BaseDecls() (baseDecls []dst.Decl, err error)
	MockStructs() (structDecls []dst.Decl, err error)
	MethodStructs(fn Func) (structDecls []dst.Decl, err error)
	NewFunc() (funcDecl *dst.FuncDecl, err error)
	IsolationAccessor(suffix, fnName string) (funcDecl *dst.FuncDecl, err error)
	FuncClosure(fn Func) (funcDecl *dst.FuncDecl, err error)
	MockMethod(fn Func) (funcDecl *dst.FuncDecl, err error)
	RecorderMethods(fn Func) (funcDecls []dst.Decl, err error)
	ResetMethod() (funcDecl *dst.FuncDecl, err error)
	AssertMethod() (funcDecl *dst.FuncDecl, err error)
}

// MoqGenerator generates moqs
type MoqGenerator struct {
	typeCache      TypeCache
	getwdFn        GetwdFunc
	newConverterFn NewConverterFunc
}

// New returns a new MoqGenerator
func New(
	typeCache TypeCache,
	getwdFn GetwdFunc,
	newConverterFn NewConverterFunc,
) *MoqGenerator {
	return &MoqGenerator{
		typeCache:      typeCache,
		getwdFn:        getwdFn,
		newConverterFn: newConverterFn,
	}
}

// MoqResponse contains the results of generating a mock
type MoqResponse struct {
	File       *dst.File
	DestPath   string
	OutPkgPath string
}

// Generate generates moqs for the given types in the given destination package
func (g *MoqGenerator) Generate(req GenerateRequest) (MoqResponse, error) {
	relPath, err := g.relativePath(req.WorkingDir)
	if err != nil {
		return MoqResponse{}, err
	}

	outPkgPath, err := g.outPackagePath(req, relPath)
	if err != nil {
		return MoqResponse{}, err
	}
	file := initializeFile(outPkgPath)

	imp := importPath(req.Import, relPath)

	destPath, err := destinationPath(req, relPath)
	if err != nil {
		return MoqResponse{}, err
	}

	var decls []dst.Decl
	for _, inType := range req.Types {
		id := ast.IdPath(inType, imp)
		typeInfo, err := g.typeCache.Type(*id, imp, req.TestImport)
		if err != nil {
			return MoqResponse{}, err
		}

		if req.ExcludeNonExported && !typeInfo.Exported {
			return MoqResponse{}, fmt.Errorf("%w: %s mocked type is not exported", ErrNonExported, id.String())
		}

		fInfo := &funcInfo{excludeNonExported: req.ExcludeNonExported, fabricated: typeInfo.Fabricated}
		// Clone the type because findFuncs may reduce the interface
		//nolint:forcetypeassert // if dst.Clone returns a different type, panic
		typeInfo.Type = dst.Clone(typeInfo.Type).(*dst.TypeSpec)
		tErr := g.findFuncs(typeInfo, fInfo)
		if tErr != nil {
			return MoqResponse{}, tErr
		}

		if req.ExcludeNonExported && fInfo.reduced {
			if len(fInfo.funcs) == 0 {
				return MoqResponse{}, fmt.Errorf("%w: type %s only contains non-exported types",
					ErrNonExported, id.String())
			}

			logs.Warnf("Mock implementation removed non-exported"+
				" methods for %s because ExcludeNonExported was set",
				id.String())
		}

		typ := Type{
			TypeInfo:   typeInfo,
			Funcs:      fInfo.funcs,
			OutPkgPath: outPkgPath,
			Reduced:    fInfo.reduced,
		}
		converter := g.newConverterFn(typ, req.Export)

		structs, err := g.structs(converter, typ)
		if err != nil {
			return MoqResponse{}, err
		}
		decls = append(decls, structs...)

		decl, err := converter.NewFunc()
		if err != nil {
			return MoqResponse{}, err
		}
		decls = append(decls, decl)

		meths, err := g.methods(converter, typ, fInfo.funcs)
		if err != nil {
			return MoqResponse{}, err
		}
		decls = append(decls, meths...)

		decl, err = converter.ResetMethod()
		if err != nil {
			return MoqResponse{}, err
		}
		decls = append(decls, decl)

		decl, err = converter.AssertMethod()
		if err != nil {
			return MoqResponse{}, err
		}
		decls = append(decls, decl)
	}
	file.Decls = decls

	return MoqResponse{
		File:       file,
		DestPath:   destPath,
		OutPkgPath: outPkgPath,
	}, nil
}

func (g *MoqGenerator) relativePath(workingDir string) (string, error) {
	wd, err := g.getwdFn()
	if err != nil {
		return "", fmt.Errorf("error getting current working directory: %w", err)
	}

	if workingDir == wd || workingDir == "" {
		return ".", nil
	}

	relPath, err := filepath.Rel(wd, workingDir)
	if err != nil {
		return "", fmt.Errorf("error getting relative import path from %s to %s: %w",
			wd, workingDir, err)
	}

	return relPath, err
}

func (g *MoqGenerator) outPackagePath(req GenerateRequest, relPath string) (string, error) {
	destDir := req.DestinationDir
	if !filepath.IsAbs(destDir) {
		destDir = filepath.Join(relPath, destDir)
	}
	if filepath.IsAbs(req.Destination) {
		destDir = req.Destination
	} else {
		destDir = filepath.Join(destDir, req.Destination)
	}
	if strings.HasSuffix(destDir, ".go") {
		destDir = filepath.Dir(destDir)
	}
	outPkgPath, err := g.typeCache.FindPackage(destDir)
	if err != nil {
		return "", err
	}
	if req.Package == "" || req.Package == "." {
		if !req.Export {
			outPkgPath += testPkgSuffix
		}
	} else {
		outPkgPath = filepath.Join(filepath.Dir(outPkgPath), req.Package)
	}
	logs.Debugf("Output package: %s", outPkgPath)
	return outPkgPath, nil
}

func initializeFile(pkg string) *dst.File {
	fSet := token.NewFileSet()

	base := filepath.Base(pkg)
	src := fmt.Sprintf("%s\n\npackage %s\n", headerComment, base)
	file, err := decorator.NewDecoratorWithImports(fSet, pkg, goast.New()).Parse(src)
	if err != nil {
		logs.Panic("Could not create decorator", err)
	}

	return file
}

func importPath(imp, relPath string) string {
	if imp != "." && !strings.HasPrefix(imp, "./") {
		return imp
	}

	imp = filepath.Join(relPath, imp)
	if strings.HasPrefix(imp, ".") {
		return imp
	}

	// Relative imports must always start with a `.` and filepath.Join will
	// remove a prefixed `./` of a relative path
	return "./" + imp
}

func destinationPath(req GenerateRequest, relPath string) (string, error) {
	if req.Export && strings.HasSuffix(req.Destination, "_test.go") {
		logs.Warn("Exported moq in a test file will not be accessible in" +
			" other packages. Remove --export option or set the --destination" +
			" to a non-test file.")
	}

	if req.Destination != "" && req.DestinationDir != "" {
		return "", fmt.Errorf("%w: both --destination and"+
			" --destination-dir flags must not be present together", ErrInvalidConfig)
	}

	destPath := req.Destination
	if destPath == "" {
		destPath = "moq_"
		for n, typ := range req.Types {
			typ = strings.TrimSuffix(typ, ast.GenTypeSuffix)
			destPath += strings.ToLower(typ)
			if n+1 < len(req.Types) {
				destPath += "_"
			}
		}
		if !req.Export || (req.Package != "" && strings.HasSuffix(req.Package, testPkgSuffix)) {
			destPath += testPkgSuffix
		}
		destPath += ".go"
	}

	if filepath.IsAbs(destPath) {
		return destPath, nil
	}

	if filepath.IsAbs(req.DestinationDir) {
		return filepath.Join(req.DestinationDir, destPath), nil
	}

	return filepath.Join(relPath, req.DestinationDir, destPath), nil
}

type funcInfo struct {
	excludeNonExported bool
	funcs              []Func
	reduced            bool
	fabricated         bool
}

func (g *MoqGenerator) findFuncs(tInfo ast.TypeInfo, fInfo *funcInfo) error {
	switch typ := tInfo.Type.Type.(type) {
	case *dst.InterfaceType:
		return g.loadNestedInterfaces(typ, tInfo, fInfo)
	case *dst.FuncType:
		fn := Func{FuncType: typ, ParentType: tInfo}
		fully, err := g.isFnFullyExported(fn)
		if err != nil {
			return err
		}

		if fInfo.excludeNonExported && !fully {
			return fmt.Errorf("%w: %s (%s) mocked type is not exported",
				ErrNonExported, tInfo.Type.Name.String(), tInfo.PkgPath)
		}

		fInfo.funcs = append(fInfo.funcs, fn)
		return nil
	case *dst.Ident:
		return g.loadTypeEquivalent(typ, tInfo, fInfo)
	default:
		logs.Panicf("Unknown type: %#v", tInfo.Type.Type)
		panic("unreachable")
	}
}

func (g *MoqGenerator) loadNestedInterfaces(iType *dst.InterfaceType, tInfo ast.TypeInfo, fInfo *funcInfo) error {
	var finalFuncs []*dst.Field
	for _, method := range iType.Methods.List {
		switch typ := method.Type.(type) {
		case *dst.FuncType:
			name := method.Names[0].Name
			if !dst.IsExported(name) && fInfo.excludeNonExported {
				fInfo.reduced = true
				continue
			}

			fn := Func{
				Name:       name,
				ParentType: tInfo,
				FuncType:   typ,
			}

			if fInfo.excludeNonExported {
				fully, err := g.isFnFullyExported(fn)
				if err != nil {
					return err
				}

				if !fully {
					fInfo.reduced = true
					continue
				}
			}

			fInfo.funcs = append(fInfo.funcs, fn)
		case *dst.Ident:
			var err error
			err = g.loadTypeEquivalent(typ, tInfo, fInfo)
			if err != nil {
				return err
			}
		default:
			logs.Panicf("Unknown type in interface method list: %#v", method.Type)
		}
		finalFuncs = append(finalFuncs, method)
	}
	if fInfo.reduced && (tInfo.Fabricated || fInfo.excludeNonExported) {
		// Reduces fabricated interface if any methods were removed
		iType.Methods.List = finalFuncs
	}

	return nil
}

func (g *MoqGenerator) loadTypeEquivalent(id *dst.Ident, tInfo ast.TypeInfo, fInfo *funcInfo) error {
	nestedType, err := g.typeCache.Type(*id, tInfo.PkgPath, false)
	if err != nil {
		return err
	}

	if fInfo.excludeNonExported && !nestedType.Exported {
		fInfo.reduced = true
		return nil
	}

	err = g.findFuncs(nestedType, fInfo)
	if err != nil {
		return err
	}

	return nil
}

func (g *MoqGenerator) isFnFullyExported(fn Func) (bool, error) {
	if fn.Name != "" && !dst.IsExported(fn.Name) {
		return false, nil
	}

	ex, err := g.isFieldListFullyExported(fn.FuncType.Params, fn.ParentType.PkgPath)
	if err != nil || !ex {
		return ex, err
	}
	ex, err = g.isFieldListFullyExported(fn.FuncType.Results, fn.ParentType.PkgPath)
	if err != nil || !ex {
		return ex, err
	}

	return true, nil
}

func (g *MoqGenerator) isFieldListFullyExported(fl *dst.FieldList, contextPkg string) (bool, error) {
	if fl == nil {
		return true, nil
	}

	for _, f := range fl.List {
		exported, err := g.isExprFullyExported(f.Type, contextPkg)
		if err != nil || !exported {
			return false, err
		}
	}

	return true, nil
}

func (g *MoqGenerator) isExprFullyExported(expr dst.Expr, contextPkg string) (bool, error) {
	switch typ := expr.(type) {
	case *dst.ArrayType:
		return g.isExprFullyExported(typ.Elt, contextPkg)
	case *dst.ChanType:
		return g.isExprFullyExported(typ.Value, contextPkg)
	case *dst.Ellipsis:
		return g.isExprFullyExported(typ.Elt, contextPkg)
	case *dst.FuncType:
		exported, err := g.isFieldListFullyExported(typ.Params, contextPkg)
		if err != nil || !exported {
			return false, err
		}
		return g.isFieldListFullyExported(typ.Results, contextPkg)
	case *dst.Ident:
		if typ.IsExported() {
			// Quick check if exported. If not exported, need to get the
			// actual type from the type cache. For instance, IsExported
			// returns false for primitive types but the cache works around
			// this.
			return true, nil
		}

		subType, err := g.typeCache.Type(*typ, contextPkg, false)
		if err != nil {
			return false, err
		}

		return subType.Exported, nil
	case *dst.IndexExpr:
		return g.isExprFullyExported(typ.X, contextPkg)
	case *dst.InterfaceType:
		return g.isFieldListFullyExported(typ.Methods, contextPkg)
	case *dst.MapType:
		exported, err := g.isExprFullyExported(typ.Key, contextPkg)
		if err != nil || !exported {
			return false, err
		}
		return g.isExprFullyExported(typ.Value, contextPkg)
	case *dst.StarExpr:
		return g.isExprFullyExported(typ.X, contextPkg)
	case *dst.StructType:
		return g.isFieldListFullyExported(typ.Fields, contextPkg)
	default:
		return false, fmt.Errorf("%w: function contains a %#v",
			ErrUnknownFieldType, typ)
	}
}

func (g *MoqGenerator) structs(converter Converterer, typ Type) ([]dst.Decl, error) {
	decls, err := converter.BaseDecls()
	if err != nil {
		return nil, err
	}
	mockStructs, err := converter.MockStructs()
	if err != nil {
		return nil, err
	}
	decls = append(decls, mockStructs...)

	for _, fn := range typ.Funcs {
		structs, err := converter.MethodStructs(fn)
		if err != nil {
			return nil, err
		}
		decls = append(decls, structs...)
	}

	return decls, nil
}

func (g *MoqGenerator) methods(
	converter Converterer, typ Type, funcs []Func,
) ([]dst.Decl, error) {
	var decls []dst.Decl

	switch typ.TypeInfo.Type.Type.(type) {
	case *dst.InterfaceType, *dst.Ident:
		decl, err := converter.IsolationAccessor(mockIdent, mockFnName)
		if err != nil {
			return nil, err
		}
		decls = append(decls, decl)

		for _, fn := range funcs {
			meth, err := converter.MockMethod(fn)
			if err != nil {
				return nil, err
			}
			decls = append(decls, meth)
		}

		decl, err = converter.IsolationAccessor(recorderIdent, onCallFnName)
		if err != nil {
			return nil, err
		}
		decls = append(decls, decl)

		for _, fn := range funcs {
			meths, err := converter.RecorderMethods(fn)
			if err != nil {
				return nil, err
			}
			decls = append(decls, meths...)
		}
	case *dst.FuncType:
		if len(funcs) != 1 {
			logs.Panicf("Function moqs should have just one function, found: %d",
				len(funcs))
		}

		fnClos, err := converter.FuncClosure(funcs[0])
		if err != nil {
			return nil, err
		}
		decls = append(decls, fnClos)

		meths, err := converter.RecorderMethods(funcs[0])
		if err != nil {
			return nil, err
		}
		decls = append(decls, meths...)
	default:
		logs.Panicf("Unknown type: %#v", typ.TypeInfo.Type.Type)
	}

	return decls, nil
}
