package ast

import (
	"fmt"

	"github.com/dave/dst"
)

// Cache loads types from the AST and caches results
type Cache struct {
	shallowPointerCompare   bool
	shallowInterfaceCompare bool
	loadTypes               LoadTypesFn

	typesByIdent map[string]*dst.TypeSpec
	loadedPkgs   map[string]pkg
}

type pkg struct {
	loadTestTypes bool
	pkgPath       string
}

// NewCache returns a new empty Cache
func NewCache(shallowPointerCompare bool, shallowInterfaceCompare bool, loadTypesFn LoadTypesFn) *Cache {
	return &Cache{
		shallowPointerCompare:   shallowPointerCompare,
		shallowInterfaceCompare: shallowInterfaceCompare,
		loadTypes:               loadTypesFn,

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
		return nil, "", fmt.Errorf("type %s not found, original package %s", realId, id.Path)
	}

	return typ, pkgPath, nil
}

// IsComparable determines if an expression is comparable
func (c *Cache) IsComparable(expr dst.Expr) (bool, error) {
	switch e := expr.(type) {
	case *dst.ArrayType:
		if e.Len == nil {
			return false, nil
		}
		return c.IsComparable(e.Elt)
	case *dst.MapType, *dst.Ellipsis:
		return false, nil
	case *dst.StarExpr:
		return c.shallowPointerCompare, nil
	case *dst.InterfaceType:
		return c.shallowInterfaceCompare, nil
	case *dst.Ident:
		if e.Obj != nil {
			return c.IsComparable(e.Obj.Decl.(*dst.TypeSpec).Type)
		}
		typ, ok := c.typesByIdent[e.String()]
		if ok {
			return c.IsComparable(typ.Type)
		}

		// Builtin type?
		if e.Path == "" {
			// error is the one builtin type that isn't comparable (it's an interface)
			if e.Name == "error" {
				return false, nil
			}

			return true, nil
		}

		_, err := c.loadPackage(e.Path, false)
		if err != nil {
			return false, err
		}

		typ, ok = c.typesByIdent[e.String()]
		if ok {
			return c.IsComparable(typ.Type)
		}

		return true, nil
	case *dst.SelectorExpr:
		path := e.X.(*dst.Ident).Name
		_, err := c.loadPackage(path, false)
		if err != nil {
			return false, err
		}

		typ, ok := c.typesByIdent[IdPath(e.Sel.Name, path).String()]
		if ok {
			return c.IsComparable(typ.Type)
		}

		// Builtin type?
		return true, nil
	case *dst.StructType:
		for _, f := range e.Fields.List {
			comp, err := c.IsComparable(f.Type)
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
