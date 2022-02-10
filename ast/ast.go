package ast

import (
	"fmt"
	"path/filepath"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"golang.org/x/tools/go/packages"
)

// FindPackage finds the package for a given directory
func FindPackage(dir string) (string, error) {
	if !filepath.IsAbs(dir) {
		// go list (which is called by packages.Load) requires that relative
		// paths start with a ./
		dir = "./" + dir
	}
	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.NeedName,
		Tests: false,
	}, dir)
	if err != nil {
		return "", err
	}
	if len(pkgs) == 0 {
		return "", fmt.Errorf("no packages found in %s", dir)
	}
	if len(pkgs) > 1 {
		return "", fmt.Errorf("too many packages found in %s (%d)", dir, len(pkgs))
	}
	if len(pkgs[0].Errors) > 0 {
		return "", pkgs[0].Errors[0]
	}

	return pkgs[0].Name, nil
}

//go:generate moqueries LoadTypesFn

// LoadTypesFn is the function type of LoadTypes
type LoadTypesFn func(pkg string, loadTestTypes bool) (
	typeSpecs []*dst.TypeSpec, pkgPath string, err error)

// LoadTypes loads all the types in the specified package
func LoadTypes(loadPkg string, loadTestTypes bool) ([]*dst.TypeSpec, string, error) {
	pkgs, err := load(loadPkg, loadTestTypes)
	if err != nil {
		return nil, "", err
	}

	var pkgPath string

	var typs []*dst.TypeSpec
	for _, pkg := range pkgs {
		pkgPath = pkg.PkgPath
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*dst.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, okT := spec.(*dst.TypeSpec); okT {
							typs = append(typs, typeSpec)
						}
					}
				}
			}
		}
	}

	return typs, pkgPath, nil
}

// load was copied from github.com/dave/dst/decorator.Load with minor modifications
func load(loadPkg string, loadTestTypes bool) ([]*decorator.Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes,
		Tests: loadTestTypes,
	}, loadPkg)
	if err != nil {
		return nil, err
	}

	dpkgs := map[*packages.Package]*decorator.Package{}

	var out []*decorator.Package
	for _, pkg := range pkgs {
		p, cErr := convert(pkg, dpkgs)
		if cErr != nil {
			return nil, cErr
		}
		out = append(out, p)
	}

	return out, nil
}

// convert was copied from github.com/dave/dst/decorator.Load with minor modifications
func convert(pkg *packages.Package, dpkgs map[*packages.Package]*decorator.Package) (*decorator.Package, error) {
	if dp, ok := dpkgs[pkg]; ok {
		return dp, nil
	}
	p := &decorator.Package{
		Package: pkg,
		Imports: map[string]*decorator.Package{},
	}
	dpkgs[pkg] = p
	if len(pkg.Syntax) > 0 {

		// Only decorate files in the GoFiles list. Syntax also has preprocessed cgo files which
		// break things.
		goFiles := make(map[string]bool, len(pkg.GoFiles))
		for _, fpath := range pkg.GoFiles {
			goFiles[fpath] = true
		}

		p.Decorator = decorator.NewDecoratorFromPackage(pkg)
		p.Decorator.ResolveLocalPath = true
		for _, f := range pkg.Syntax {
			fpath := pkg.Fset.File(f.Pos()).Name()
			if !goFiles[fpath] {
				continue
			}
			file, err := p.Decorator.DecorateFile(f)
			if err != nil {
				return nil, err
			}
			p.Syntax = append(p.Syntax, file)
		}

		dir, _ := filepath.Split(pkg.Fset.File(pkg.Syntax[0].Pos()).Name())
		p.Dir = dir

		for path, imp := range pkg.Imports {
			dimp, err := convert(imp, dpkgs)
			if err != nil {
				return nil, err
			}
			p.Imports[path] = dimp
		}
	}
	return p, nil
}
