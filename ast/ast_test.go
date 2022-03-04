package ast_test

import (
	"strings"
	"testing"

	"github.com/dave/dst"

	"github.com/myshkin5/moqueries/ast"
)

func TestFindPackage(t *testing.T) {
	t.Run("finds the package of the current directory", func(t *testing.T) {
		// ASSEMBLE

		// ACT
		pkg, err := ast.FindPackage(".")
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}
		expectedPkg := "ast"
		if !strings.HasSuffix(pkg, expectedPkg) {
			t.Errorf("got %s, wanted suffix %s", pkg, expectedPkg)
		}
	})

	t.Run("finds a package relative to the current package", func(t *testing.T) {
		// ASSEMBLE

		// ACT
		pkg, err := ast.FindPackage("./testpkg")
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}
		expectedPkg := "testpkg"
		if !strings.HasSuffix(pkg, expectedPkg) {
			t.Errorf("got %s, wanted suffix %s", pkg, expectedPkg)
		}
	})

	t.Run("finds a package relative to the current package (without initial `.`)", func(t *testing.T) {
		// ASSEMBLE

		// ACT
		pkg, err := ast.FindPackage("testpkg")
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}
		expectedPkg := "testpkg"
		if !strings.HasSuffix(pkg, expectedPkg) {
			t.Errorf("got %s, wanted suffix %s", pkg, expectedPkg)
		}
	})

	t.Run("returns an error when there are no files in the package", func(t *testing.T) {
		// ASSEMBLE

		// ACT
		_, err := ast.FindPackage("randomandnonexistent")

		// ASSERT
		if err == nil {
			t.Errorf("got nil, wanted not nil")
		}
		expectedErr := "ast/randomandnonexistent: directory not found"
		if err != nil && !strings.Contains(err.Error(), expectedErr) {
			t.Errorf("got %s, wanted to contain %s", err.Error(), expectedErr)
		}
	})
}

func TestLoadTypes(t *testing.T) {
	t.Run("loads the expected interfaces", func(t *testing.T) {
		// ASSEMBLE

		// ACT
		typs, pkgPath, err := ast.LoadTypes("builtin", false)
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}
		var iTypes []*dst.TypeSpec
		for _, typ := range typs {
			if _, ok := typ.Type.(*dst.InterfaceType); ok {
				iTypes = append(iTypes, typ)
			}
		}
		if len(iTypes) != 1 {
			t.Errorf("got %d, wanted 1", len(iTypes))
		}
		if len(iTypes) > 0 && iTypes[0].Name.Name != "error" {
			t.Errorf("got %s, wanted error", iTypes[0].Name.Name)
		}
		if pkgPath != "builtin" {
			t.Errorf("got %s, wanted builtin", pkgPath)
		}
	})

	t.Run("loads the expected functions", func(t *testing.T) {
		// ASSEMBLE

		// ACT
		typs, pkgPath, err := ast.LoadTypes("bufio", false)
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}
		var fTypes []*dst.TypeSpec
		for _, typ := range typs {
			if _, ok := typ.Type.(*dst.FuncType); ok {
				fTypes = append(fTypes, typ)
			}
		}
		if len(fTypes) != 1 {
			t.Errorf("got %d, wanted 1", len(fTypes))
		}
		if len(fTypes) > 0 && fTypes[0].Name.Name != "SplitFunc" {
			t.Errorf("got %s, wanted error", fTypes[0].Name.Name)
		}
		if pkgPath != "bufio" {
			t.Errorf("got %s, wanted bufio", pkgPath)
		}
	})

	t.Run("resolves local paths", func(t *testing.T) {
		// ASSEMBLE

		// ACT
		typs, pkgPath, err := ast.LoadTypes("github.com/myshkin5/moqueries/generator", false)
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		var iTypes []*dst.TypeSpec
		for _, typ := range typs {
			if _, ok := typ.Type.(*dst.InterfaceType); ok {
				if typ.Name.Name == "Converterer" {
					iTypes = append(iTypes, typ)
				}
			}
		}
		if len(iTypes) != 1 {
			t.Fatalf("got %d, wanted 1", len(iTypes))
		}
		expectedPkgPath := "github.com/myshkin5/moqueries/generator"
		if pkgPath != expectedPkgPath {
			t.Errorf("got %s, wanted %s", pkgPath, expectedPkgPath)
		}

		var fnType *dst.FuncType
		iType, ok := iTypes[0].Type.(*dst.InterfaceType)
		if !ok {
			t.Fatalf("wanted an interface type, got %#v", iTypes[0].Type)
		}
		for _, field := range iType.Methods.List {
			if field.Names[0].Name == "MockMethod" {
				fnType, ok = field.Type.(*dst.FuncType)
				if !ok {
					t.Fatalf("wanted a function type, got %#v", field.Type)
				}
			}
		}
		if fnType == nil {
			t.Fatalf("got nil, wanted not nil")
		}

		paramIdent := fnType.Params.List[0].Names[0].Name
		if paramIdent != "fn" {
			t.Errorf("got %s, wanted typ", paramIdent)
		}
		idType, ok := fnType.Params.List[0].Type.(*dst.Ident)
		if !ok {
			t.Fatalf("wanted an ident type, got %#v", fnType.Params.List[0].Type)
		}
		if idType.Name != "Func" {
			t.Errorf("got %s, wanted Type", idType.Name)
		}
		if idType.Path != expectedPkgPath {
			t.Errorf("got %s, wanted %s", idType.Path, expectedPkgPath)
		}
	})

	t.Run("resolves relative paths", func(t *testing.T) {
		// ASSEMBLE

		// ACT
		_, pkgPath, err := ast.LoadTypes(".", false)
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		expectedPkgPath := "github.com/myshkin5/moqueries/ast"
		if pkgPath != expectedPkgPath {
			t.Errorf("got %s, wanted %s", pkgPath, expectedPkgPath)
		}
	})

	t.Run("loads test types", func(t *testing.T) {
		// ASSEMBLE
		var _ TestFn = func() {}

		// ACT
		typs, pkgPath, err := ast.LoadTypes("github.com/myshkin5/moqueries/ast", true)
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		var fTypes []*dst.TypeSpec
		for _, typ := range typs {
			if _, ok := typ.Type.(*dst.FuncType); ok {
				fTypes = append(fTypes, typ)
			}
		}
		found := false
		for _, fType := range fTypes {
			if fType.Name.Name == "TestFn" {
				found = true
			}
		}
		if !found {
			t.Errorf("got not found, wanted found")
		}
		expectedPkgPath := "github.com/myshkin5/moqueries/ast.test"
		if pkgPath != expectedPkgPath {
			t.Errorf("got %s, wanted %s", pkgPath, expectedPkgPath)
		}
	})
}

type TestFn func()
