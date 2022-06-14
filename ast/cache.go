// Package ast provides utilities for working with Go's Abstract Syntax Tree
package ast

import (
	"errors"
	"fmt"
	"path/filepath"
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
	testPkgSuffix = "_test"
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

// Cache loads types from the AST and caches the results
type Cache struct {
	load    LoadFn
	metrics metrics.Metrics

	typesByIdent map[string]*dst.TypeSpec
	loadedPkgs   map[string]*pkgInfo
}

type pkgInfo struct {
	directLoaded bool
	loadTestPkgs bool
	typesIndexed bool
	pkg          *decorator.Package
}

// NewCache returns a new empty Caches
func NewCache(load LoadFn, metrics metrics.Metrics) *Cache {
	return &Cache{
		load:    load,
		metrics: metrics,

		typesByIdent: make(map[string]*dst.TypeSpec),
		loadedPkgs:   make(map[string]*pkgInfo),
	}
}

// Type returns the requested TypeSpec or an error if the type can't be found
func (c *Cache) Type(id dst.Ident, testImport bool) (*dst.TypeSpec, string, error) {
	if strings.HasSuffix(id.Path, testPkgSuffix) {
		// Strip the _test suffix when loading a package but set testImport so
		// we know to add it back later
		id.Path = strings.TrimSuffix(id.Path, testPkgSuffix)
		testImport = true
	}

	pkgPath, err := c.loadPackage(id.Path, testImport)
	if err != nil {
		return nil, "", err
	}

	if testImport && !strings.HasSuffix(pkgPath, testPkgSuffix) {
		pkgPath += testPkgSuffix
	}

	realId := IdPath(id.Name, pkgPath).String()
	typ, ok := c.typesByIdent[realId]
	if !ok {
		return nil, "", fmt.Errorf(
			"%w: %q (original package %q)", ErrTypeNotFound, realId, id.Path)
	}

	return typ, pkgPath, nil
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

func (c *Cache) isDefaultComparable(expr dst.Expr, interfacePointerDefault bool) (bool, error) {
	switch e := expr.(type) {
	case *dst.ArrayType:
		if e.Len == nil {
			return false, nil
		}
		return c.isDefaultComparable(e.Elt, interfacePointerDefault)
	case *dst.MapType, *dst.Ellipsis:
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
			return c.isDefaultComparable(typ.Type, interfacePointerDefault)
		}
		typ, ok := c.typesByIdent[e.String()]
		if ok {
			return c.isDefaultComparable(typ.Type, interfacePointerDefault)
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
			return c.isDefaultComparable(typ.Type, interfacePointerDefault)
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
			return c.isDefaultComparable(typ.Type, interfacePointerDefault)
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
	loadedPkg, ok := c.loadedPkgs[path]
	if ok {
		// If we already loaded the test packages or if the test packages
		// aren't requested, we're done
		if loadedPkg.loadTestPkgs || !testImport {
			// If we direct loaded and indexed the types, we're done
			if loadedPkg.directLoaded && loadedPkg.typesIndexed {
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
		c.loadedPkgs[path] = c.loadedPkgs[pkgPath]
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
		}
		for _, file := range pkg.pkg.Syntax {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*dst.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, okT := spec.(*dst.TypeSpec); okT {
							ident := dst.Ident{
								Name: typeSpec.Name.Name,
								Path: pkg.pkg.PkgPath,
							}
							c.typesByIdent[ident.String()] = typeSpec
						}
					}
				}
			}
		}
		pkg.typesIndexed = true
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
	logs.Debugf("Loading package %s (test packages: %t) took %s",
		loadPkg, testImport, loadTime.String())
	if err != nil {
		return nil, err
	}

	var out []*pkgInfo
	for _, pkg := range pkgs {
		p, cErr := c.convert(pkg, testImport, true)
		if cErr != nil {
			return nil, cErr
		}
		out = append(out, p)
	}

	return out, nil
}

// convert was copied from github.com/dave/dst/decorator.Load with minor modifications
func (c *Cache) convert(pkg *packages.Package, testImport, directLoaded bool) (*pkgInfo, error) {
	p := &pkgInfo{
		directLoaded: directLoaded,
		loadTestPkgs: testImport,
		pkg: &decorator.Package{
			Package: pkg,
			Imports: map[string]*decorator.Package{},
		},
	}
	c.loadedPkgs[pkg.PkgPath] = p
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
		logs.Debugf("Decorating package %s took %s",
			pkg.PkgPath, time.Since(start).String())

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
