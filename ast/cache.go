// Package ast provides utilities for working with Go's Abstract Syntax Tree
package ast

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"golang.org/x/tools/go/packages"

	"github.com/myshkin5/moqueries/logs"
	"github.com/myshkin5/moqueries/metrics"
)

//go:generate moqueries LoadFn

const (
	builtinPkg = "builtin"

	genTypeSuffix     = "_genType"
	starGenTypeSuffix = "_starGenType"
	testPkgSuffix     = "_test"
)

// LoadFn is the function type of packages.Load
type LoadFn func(cfg *packages.Config, patterns ...string) ([]*packages.Package, error)

var (
	// ErrTypeNotFound is returned when a type isn't in the cache and can't be
	// loaded
	ErrTypeNotFound = errors.New("type not found")
	// ErrInvalidType is returned when the type doesn't have the expected
	// structure
	ErrInvalidType = errors.New("type did not have expected format")
)

// Cache loads packages from the AST and caches the results
type Cache struct {
	load    LoadFn
	metrics metrics.Metrics

	typesByIdent       map[string]*typInfo
	funcDeclsByIdent   map[string]*funcDeclInfo
	methodDeclsByIdent map[string]*methodDeclInfo
	loadedPkgs         map[string]*pkgInfo
}

type typInfo struct {
	id  dst.Ident
	typ *dst.TypeSpec
}

type methodDeclInfo struct {
	id    dst.Ident
	recv  *dst.Expr
	funcs []*funcDeclInfo
}

type funcDeclInfo struct {
	id  dst.Ident
	typ *dst.FuncType
}

type pkgInfo struct {
	directLoaded bool
	loadTestPkgs bool
	pkg          *decorator.Package
}

// NewCache returns a new empty Caches
func NewCache(load LoadFn, metrics metrics.Metrics) *Cache {
	return &Cache{
		load:    load,
		metrics: metrics,

		typesByIdent:       make(map[string]*typInfo),
		funcDeclsByIdent:   make(map[string]*funcDeclInfo),
		methodDeclsByIdent: make(map[string]*methodDeclInfo),
		loadedPkgs:         make(map[string]*pkgInfo),
	}
}

// TypeInfo returns all the information the cache holds for a type
type TypeInfo struct {
	Type       *dst.TypeSpec
	PkgPath    string
	Exported   bool
	Fabricated bool
}

// Type returns the requested TypeSpec or an error if the type can't be found
func (c *Cache) Type(id dst.Ident, contextPkg string, testImport bool) (TypeInfo, error) {
	loadPkg := id.Path
	if loadPkg == "" {
		loadPkg = contextPkg
	}
	if strings.HasSuffix(loadPkg, testPkgSuffix) {
		// Strip the _test suffix when loading a package but set testImport so
		// we know to add it back later
		loadPkg = strings.TrimSuffix(loadPkg, testPkgSuffix)
		testImport = true
	}

	pkgPath, err := c.loadPackage(loadPkg, testImport)
	if err != nil {
		return TypeInfo{}, err
	}

	if testImport && !strings.HasSuffix(pkgPath, testPkgSuffix) {
		pkgPath += testPkgSuffix
	}

	realId := IdPath(id.Name, pkgPath).String()
	if typ, ok := c.typesByIdent[realId]; ok {
		return TypeInfo{
			Type:       typ.typ,
			PkgPath:    pkgPath,
			Exported:   isExported(typ.typ.Name.Name, pkgPath),
			Fabricated: false,
		}, nil
	}
	if funcDecl, ok := c.funcDeclsByIdent[realId]; ok {
		return TypeInfo{
			Type: &dst.TypeSpec{
				Name: &funcDecl.id,
				Type: funcDecl.typ,
			},
			PkgPath:    pkgPath,
			Exported:   isExported(funcDecl.id.Name, pkgPath),
			Fabricated: true,
		}, nil
	}
	if methodDecl, ok := c.methodDeclsByIdent[realId]; ok {
		return TypeInfo{
			Type: &dst.TypeSpec{
				Name: &methodDecl.id,
				Type: c.fabricateInterfaceType(methodDecl.funcs),
			},
			PkgPath:    pkgPath,
			Exported:   isExported(methodDecl.id.Name, pkgPath),
			Fabricated: true,
		}, nil
	}

	if id.Path != "" {
		return TypeInfo{}, fmt.Errorf(
			"%w: %q (original package %q)", ErrTypeNotFound, realId, id.Path)
	}

	pkgPath, err = c.loadPackage(builtinPkg, false)
	if err != nil {
		return TypeInfo{}, err
	}

	realId = IdPath(id.Name, pkgPath).String()
	typ, ok := c.typesByIdent[realId]
	if !ok {
		return TypeInfo{}, fmt.Errorf(
			"%w: %q (original package %q)", ErrTypeNotFound, realId, id.Path)
	}

	return TypeInfo{
		Type: typ.typ,
		// PkgPath for builtin types is always empty
		PkgPath:    "",
		Exported:   isExported(typ.typ.Name.Name, builtinPkg),
		Fabricated: false,
	}, nil
}

// IsComparable determines if an expression is comparable
func (c *Cache) IsComparable(expr dst.Expr) (bool, error) {
	return c.isDefaultComparable(expr, true)
}

// IsDefaultComparable determines if an expression is comparable. Returns the
// same results as IsComparable but pointers and interfaces are not comparable
// by default (interface implementations that are not comparable and put into a
// map key will panic at runtime and by default pointers use a deep hash to be
// comparable).
func (c *Cache) IsDefaultComparable(expr dst.Expr) (bool, error) {
	return c.isDefaultComparable(expr, false)
}

// FindPackage finds the package for a given directory
func (c *Cache) FindPackage(dir string) (string, error) {
	if !filepath.IsAbs(dir) && !strings.HasPrefix(dir, ".") {
		// go list (which is called by packages.Load) requires that relative
		// paths start with a ./
		dir = "./" + dir
	}
	pkgPath, err := c.loadPackage(dir, false)
	if err != nil {
		return "", err
	}
	return pkgPath, nil
}

// LoadPackage loads the specified pattern of package(s) and returns a list of
// mock-able types
func (c *Cache) LoadPackage(pkgPattern string) error {
	_, err := c.loadPackage(pkgPattern, false)
	if err != nil {
		return err
	}

	return nil
}

// MockableTypes returns all the mockable types loaded so far
func (c *Cache) MockableTypes(onlyExported bool) []dst.Ident {
	filterOut := func(id dst.Ident) bool {
		dir := id.Path
		var file string
		internal := false
		vendor := false
		const (
			internalPkg = "internal"
			vendorPkg   = "vendor"
		)
		for {
			dir, file = filepath.Split(dir)
			dir = filepath.Clean(dir)
			if dir == internalPkg || file == internalPkg {
				internal = true
				break
			}
			if dir == vendorPkg || file == vendorPkg {
				vendor = true
				break
			}
			if dir == "." {
				break
			}
		}
		if internal || vendor {
			return true
		}

		return false
	}

	var typs []dst.Ident
	for _, typ := range c.typesByIdent {
		_, okInterface := typ.typ.Type.(*dst.InterfaceType)
		_, okFunc := typ.typ.Type.(*dst.FuncType)
		if !okInterface && !okFunc {
			continue
		}

		if onlyExported && !typ.id.IsExported() {
			continue
		}

		if typ.id.Path == builtinPkg {
			continue
		}

		if filterOut(typ.id) {
			continue
		}

		typs = append(typs, typ.id)
	}

	for _, funcDecl := range c.funcDeclsByIdent {
		if funcDecl.id.Path != builtinPkg && onlyExported && !funcDecl.id.IsExported() {
			continue
		}

		if filterOut(funcDecl.id) {
			continue
		}

		typs = append(typs, funcDecl.id)
	}

	for _, methodDecl := range c.methodDeclsByIdent {
		if methodDecl.id.Path != builtinPkg && onlyExported && !methodDecl.id.IsExported() {
			continue
		}

		if filterOut(methodDecl.id) {
			continue
		}

		typs = append(typs, methodDecl.id)
	}

	sort.Slice(typs, func(i, j int) bool {
		if typs[i].Path == typs[j].Path {
			return typs[i].Name < typs[j].Name
		}
		return typs[i].Path < typs[j].Path
	})

	return typs
}

func isExported(name, pkgPath string) bool {
	if dst.IsExported(name) {
		return true
	}
	if pkgPath == builtinPkg {
		// builtin types are always exported
		return true
	}
	return false
}

func (c *Cache) isDefaultComparable(expr dst.Expr, interfacePointerDefault bool) (bool, error) {
	switch e := expr.(type) {
	case *dst.ArrayType:
		if e.Len == nil {
			return false, nil
		}
		return c.isDefaultComparable(e.Elt, interfacePointerDefault)
	case *dst.MapType, *dst.Ellipsis, *dst.FuncType:
		return false, nil
	case *dst.StarExpr:
		return interfacePointerDefault, nil
	case *dst.InterfaceType:
		return interfacePointerDefault, nil
	case *dst.Ident:
		if e.Obj != nil {
			typ, ok := e.Obj.Decl.(*dst.TypeSpec)
			if !ok {
				return false, fmt.Errorf("%q: %w", e.String(), ErrInvalidType)
			}

			if typ.Name.Name == "string" && typ.Name.Path == "" {
				return true, nil
			}

			return c.isDefaultComparable(typ.Type, interfacePointerDefault)
		}
		typ, ok := c.typesByIdent[e.String()]
		if ok {
			return c.isDefaultComparable(typ.typ.Type, interfacePointerDefault)
		}

		// Builtin type?
		if e.Path == "" {
			// error is the one builtin type that may not be comparable (it's
			// an interface so return the same result as an interface)
			if e.Name == "error" {
				return interfacePointerDefault, nil
			}

			return true, nil
		}

		_, err := c.loadPackage(e.Path, false)
		if err != nil {
			return false, err
		}

		typ, ok = c.typesByIdent[e.String()]
		if ok {
			return c.isDefaultComparable(typ.typ.Type, interfacePointerDefault)
		}

		return true, nil
	case *dst.SelectorExpr:
		ex, ok := e.X.(*dst.Ident)
		if !ok {
			return false, fmt.Errorf("%q: %w", e.X, ErrInvalidType)
		}
		path := ex.Name
		_, err := c.loadPackage(path, false)
		if err != nil {
			return false, err
		}

		typ, ok := c.typesByIdent[IdPath(e.Sel.Name, path).String()]
		if ok {
			return c.isDefaultComparable(typ.typ.Type, interfacePointerDefault)
		}

		// Builtin type?
		return true, nil
	case *dst.StructType:
		for _, f := range e.Fields.List {
			comp, err := c.isDefaultComparable(f.Type, interfacePointerDefault)
			if err != nil || !comp {
				return false, err
			}
		}
	}

	return true, nil
}

func (c *Cache) loadPackage(path string, testImport bool) (string, error) {
	indexPath := path
	if strings.HasPrefix(path, ".") {
		var err error
		indexPath, err = filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("%w: error getting absolute path for %s", err, path)
		}
	}
	loadedPkg, ok := c.loadedPkgs[indexPath]
	if ok {
		// If we already loaded the test packages or if the test packages
		// aren't requested, we're done
		if loadedPkg.loadTestPkgs || !testImport {
			// If we direct loaded and indexed the types, we're done
			if loadedPkg.directLoaded {
				c.metrics.ASTPkgCacheHitsInc()
				return loadedPkg.pkg.PkgPath, nil
			}
		}
	}
	c.metrics.ASTPkgCacheMissesInc()

	pkgPath, err := c.loadTypes(path, testImport)
	if err != nil {
		return "", err
	}

	if path != pkgPath {
		logs.Debugf("Requested package %s loaded as %s", path, pkgPath)
	}

	return pkgPath, nil
}

func (c *Cache) loadTypes(loadPkg string, testImport bool) (string, error) {
	pkgs, err := c.loadAST(loadPkg, testImport)
	if err != nil {
		return "", err
	}

	foundPkg := ""
	for _, pkg := range pkgs {
		if foundPkg == "" {
			// The first package is the main package (any subsequent packages
			// are most likely test packages)
			foundPkg = pkg.pkg.PkgPath
			break
		}
	}

	return foundPkg, nil
}

func (c *Cache) loadAST(loadPkg string, testImport bool) ([]*pkgInfo, error) {
	if dp, ok := c.loadedPkgs[loadPkg]; ok {
		// If we already loaded the test types or if the test types aren't
		// requested, we're done
		if dp.loadTestPkgs || !testImport {
			// If we direct loaded, we're done
			if dp.directLoaded {
				c.metrics.ASTTypeCacheHitsInc()
				return []*pkgInfo{dp}, nil
			}
		}
	}
	c.metrics.ASTTypeCacheMissesInc()

	start := time.Now()
	pkgs, err := c.load(&packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes,
		Tests: testImport,
	}, loadPkg)
	loadTime := time.Since(start)
	c.metrics.ASTTotalLoadTimeInc(loadTime)
	logs.Debugf("Loading package %s (test packages: %t) took %s, found %d packages",
		loadPkg, testImport, loadTime.String(), len(pkgs))
	if err != nil {
		return nil, err
	}

	var out []*pkgInfo
	for _, pkg := range pkgs {
		logs.Debugf("Converting %s", pkg.PkgPath)
		p, cErr := c.convert(pkg, testImport, true)
		if cErr != nil {
			return nil, cErr
		}

		for _, file := range p.pkg.Syntax {
			for _, d := range file.Decls {
				switch decl := d.(type) {
				case *dst.GenDecl:
					c.storeTypeSpecs(decl, p)
				case *dst.FuncDecl:
					c.storeFuncDecl(decl, p)
				}
			}
		}

		out = append(out, p)
	}

	return out, nil
}

// convert was copied from github.com/dave/dst/decorator.Load with minor modifications
func (c *Cache) convert(pkg *packages.Package, testImport, directLoaded bool) (*pkgInfo, error) {
	p, ok := c.loadedPkgs[pkg.PkgPath]
	if ok && p.directLoaded {
		return p, nil
	}
	p = &pkgInfo{
		directLoaded: directLoaded,
		loadTestPkgs: testImport,
		pkg: &decorator.Package{
			Package: pkg,
			Imports: map[string]*decorator.Package{},
		},
	}
	c.loadedPkgs[pkg.PkgPath] = p
	absPath, err := findAbsPath(pkg)
	if err != nil {
		// If we can't find the absolute path, we swallow the error and don't
		// add it to the map. This is less efficient but not horrible. In
		// practice, this seems to only happen with standard library packages
		// when generating std and unsafe whenever it is referenced.
		if pkg.PkgPath != "unsafe" {
			logs.Warnf("Could not find the absolute path of the %s package", pkg.PkgPath)
		}
	} else {
		c.loadedPkgs[absPath] = p
	}
	if len(pkg.Syntax) > 0 {
		// Only decorate files in the GoFiles list. Syntax also has preprocessed cgo files which
		// break things.
		goFiles := make(map[string]bool, len(pkg.GoFiles))
		for _, fpath := range pkg.GoFiles {
			goFiles[fpath] = true
		}

		start := time.Now()
		p.pkg.Decorator = decorator.NewDecoratorFromPackage(pkg)
		p.pkg.Decorator.ResolveLocalPath = true
		for _, f := range pkg.Syntax {
			fpath := pkg.Fset.File(f.Pos()).Name()
			if !goFiles[fpath] {
				continue
			}
			file, err := p.pkg.Decorator.DecorateFile(f)
			if err != nil {
				return nil, err
			}
			p.pkg.Syntax = append(p.pkg.Syntax, file)
		}
		dTime := time.Since(start)
		c.metrics.ASTTotalDecorationTimeInc(dTime)
		logs.Debugf("Decorating package id %s, name %s, pkgPath %s took %s",
			pkg.ID, pkg.Name, pkg.PkgPath, dTime.String())

		dir, _ := filepath.Split(pkg.Fset.File(pkg.Syntax[0].Pos()).Name())
		p.pkg.Dir = dir

		for path, imp := range pkg.Imports {
			dimp, err := c.convert(imp, false, false)
			if err != nil {
				return nil, err
			}
			p.pkg.Imports[path] = dimp.pkg
		}
	}
	return p, nil
}

func findAbsPath(pkg *packages.Package) (string, error) {
	checkList := func(fpaths []string) (string, bool) {
		if len(fpaths) > 0 {
			return filepath.Dir(fpaths[0]), true
		}

		return "", false
	}

	if absPath, ok := checkList(pkg.GoFiles); ok {
		return absPath, nil
	}
	if absPath, ok := checkList(pkg.CompiledGoFiles); ok {
		return absPath, nil
	}
	if absPath, ok := checkList(pkg.OtherFiles); ok {
		return absPath, nil
	}

	return "", fmt.Errorf("could not find absolute path for %s package", pkg.PkgPath)
}

func (c *Cache) storeTypeSpecs(decl *dst.GenDecl, pkg *pkgInfo) {
	for _, spec := range decl.Specs {
		if typeSpec, okT := spec.(*dst.TypeSpec); okT {
			ident := dst.Ident{
				Name: typeSpec.Name.Name,
				Path: pkg.pkg.PkgPath,
			}
			c.typesByIdent[ident.String()] = &typInfo{
				id:  ident,
				typ: typeSpec,
			}
		}
	}
}

func (c *Cache) storeFuncDecl(decl *dst.FuncDecl, pkg *pkgInfo) {
	ident := dst.Ident{
		Name: decl.Name.Name,
		Path: pkg.pkg.PkgPath,
	}
	fnInfo := &funcDeclInfo{
		id:  ident,
		typ: decl.Type,
	}
	if decl.Recv == nil {
		fnInfo.id.Name += genTypeSuffix
		// TODO: error if func seen more than once?
		c.funcDeclsByIdent[fnInfo.id.String()] = fnInfo
		return
	}

	if len(decl.Recv.List) != 1 {
		logs.Panicf("%s has a receiver list of length %d, expected length of 1",
			ident.String(), len(decl.Recv.List))
	}
	recv := decl.Recv.List[0].Type
	suffix := genTypeSuffix
	expr := recv
	if sExpr, ok := expr.(*dst.StarExpr); ok {
		suffix = starGenTypeSuffix
		expr = sExpr.X
	}
	exprId, ok := expr.(*dst.Ident)
	if !ok {
		logs.Panicf("%s has a non-Ident (or StarExpr/Ident) receiver: %#v",
			ident.String(), expr)
	}
	keyId := dst.Ident{
		Name: exprId.Name + suffix,
		Path: pkg.pkg.PkgPath,
	}
	declInfo, ok := c.methodDeclsByIdent[keyId.String()]
	if !ok {
		declInfo = &methodDeclInfo{id: keyId, recv: &recv}
		c.methodDeclsByIdent[keyId.String()] = declInfo
	}
	declInfo.funcs = append(declInfo.funcs, fnInfo)
}

func (c *Cache) fabricateInterfaceType(funcs []*funcDeclInfo) *dst.InterfaceType {
	fl := &dst.FieldList{}
	for _, fn := range funcs {
		fl.List = append(fl.List, &dst.Field{
			Names: []*dst.Ident{dst.NewIdent(fn.id.Name)},
			Type:  fn.typ,
		})
	}
	return &dst.InterfaceType{Methods: fl}
}
