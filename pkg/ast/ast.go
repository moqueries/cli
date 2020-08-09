package ast

import (
	"fmt"
	"path/filepath"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"golang.org/x/tools/go/packages"
)

// FindPackageDir find the directory containing a package where pkg is either
// the current directory (".") or an external module
func FindPackageDir(pkg string) (string, error) {
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedFiles}, pkg)
	if err != nil {
		return "", err
	}
	if len(pkgs) == 0 {
		return "", fmt.Errorf("no packages found for package %s", pkg)
	}
	if len(pkgs) > 1 {
		return "", fmt.Errorf("more than one packages found for package %s", pkg)
	}
	if len(pkgs[0].GoFiles) == 0 {
		return "", fmt.Errorf("no files found for package %s", pkg)
	}

	return filepath.Dir(pkgs[0].GoFiles[0]), nil
}

// LoadTypes loads all of the types in the specified package
func LoadTypes(loadPkg string) ([]*dst.TypeSpec, string, error) {
	pkgs, err := decorator.Load(&packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes,
		Env: []string{"DST_INCLUDE_LOCAL_PKG=true"},
	}, loadPkg)
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
