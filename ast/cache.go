package ast

import (
	"errors"
	"fmt"

	"github.com/dave/dst"
)

var (
	// ErrTypeNotFound is returned when a type isn't in the cache and can't be
	// loaded
	ErrTypeNotFound = errors.New("type not found")
	// ErrInvalidType is returned when the type doesn't have the expected
	// structure
	ErrInvalidType = errors.New("type did not have expected format")
)

// Cache loads types from the AST and caches results
type Cache struct {
	loadTypes LoadTypesFn

	typesByIdent map[string]*dst.TypeSpec
	loadedPkgs   map[string]pkg
}

type pkg struct {
	loadTestTypes bool
	pkgPath       string
}

// NewCache returns a new empty Cache
func NewCache(loadTypesFn LoadTypesFn) *Cache {
	return &Cache{
		loadTypes: loadTypesFn,

		typesByIdent: make(map[string]*dst.TypeSpec),
		loadedPkgs:   make(map[string]pkg),
	}
}

// Type returns the requested TypeSpec or an error if the type can't be found
func (c *Cache) Type(id dst.Ident, loadTestTypes bool) (*dst.TypeSpec, string, error) {
	pkgPath, err := c.loadPackage(id.Path, loadTestTypes)
	if err != nil {
		return nil, "", err
	}

	realId := IdPath(id.Name, pkgPath).String()
	typ, ok := c.typesByIdent[realId]
	if !ok {
		return nil, "", fmt.Errorf(
			"%q (original package %q): %w", realId, id.Path, ErrTypeNotFound)
	}

	return typ, pkgPath, nil
}

// IsComparable determines if an expression is comparable
func (c *Cache) IsComparable(expr dst.Expr) (bool, error) {
	return c.isDefaultComparable(expr, true)
}

// IsDefaultComparable determines if an expression is comparable. Returns the
// same results as IsComparable but pointers and interfaces are not comparable
// by default (interfaces that are not comparable and put into a map key will
// panic at runtime and by default pointers use a deep hash to be comparable).
func (c *Cache) IsDefaultComparable(expr dst.Expr) (bool, error) {
	return c.isDefaultComparable(expr, false)
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

func (c *Cache) loadPackage(path string, loadTestTypes bool) (string, error) {
	loadedPkg, ok := c.loadedPkgs[path]
	if ok {
		// If we already loaded the test types or if the test types aren't
		// requested, we're done
		if loadedPkg.loadTestTypes || !loadTestTypes {
			return loadedPkg.pkgPath, nil
		}
	}

	typeSpecs, pkgPath, err := c.loadTypes(path, loadTestTypes)
	if err != nil {
		return "", err
	}

	for _, typeSpec := range typeSpecs {
		// We have to make a new ident as key because loaded types don't have
		// Path set?
		ident := dst.Ident{
			Name: typeSpec.Name.Name,
			Path: pkgPath,
		}
		c.typesByIdent[ident.String()] = typeSpec
	}

	c.loadedPkgs[path] = pkg{
		loadTestTypes: loadTestTypes,
		pkgPath:       pkgPath,
	}
	if path != pkgPath {
		c.loadedPkgs[pkgPath] = c.loadedPkgs[path]
	}

	return pkgPath, nil
}
