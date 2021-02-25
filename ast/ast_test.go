package ast_test

import (
	"strings"
	"testing"

	"github.com/dave/dst"

	"github.com/myshkin5/moqueries/ast"
)

func TestFindPackageDir(t *testing.T) {
	t.Run("finds the package dir of the current directory", func(t *testing.T) {
		// ASSEMBLE

		// ACT
		dir, err := ast.FindPackageDir(".")
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}
		expectedSuffix := "github.com/myshkin5/moqueries/ast"
		if !strings.HasSuffix(dir, expectedSuffix) {
			t.Errorf("got %s, wanted suffix %s", dir, expectedSuffix)
		}
	})

	t.Run("finds the package dir of an external module", func(t *testing.T) {
		// ASSEMBLE

		// ACT
		dir, err := ast.FindPackageDir("github.com/dave/dst")
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}
		splits := strings.Split(dir, "@")
		if len(splits) != 2 {
			t.Errorf("got %d, wanted 2", len(splits))
		}
		expectedSuffix := "pkg/mod/github.com/dave/dst"
		if !strings.HasSuffix(splits[0], expectedSuffix) {
			t.Errorf("got %s, wanted suffix %s", splits[0], expectedSuffix)
		}
	})

	t.Run("returns an error when there are no files in the package", func(t *testing.T) {
		// ASSEMBLE

		// ACT
		_, err := ast.FindPackageDir("randomandnonexistent")

		// ASSERT
		if err == nil {
			t.Errorf("got nil, wanted not nil")
		}
		expectedErr := "no files found for package randomandnonexistent"
		if err != nil && err.Error() != expectedErr {
			t.Errorf("got %s, wanted %s", err.Error(), expectedErr)
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
			t.Errorf("got %d, wanted 1", len(iTypes))
		}
		expectedPkgPath := "github.com/myshkin5/moqueries/generator"
		if pkgPath != expectedPkgPath {
			t.Errorf("got %s, wanted %s", pkgPath, expectedPkgPath)
		}

		var baseStruct *dst.FuncType
		if len(iTypes) > 0 {
			for _, field := range iTypes[0].Type.(*dst.InterfaceType).Methods.List {
				if field.Names[0].Name == "BaseStruct" {
					baseStruct = field.Type.(*dst.FuncType)
				}
			}
		}
		if baseStruct == nil {
			t.Errorf("got nil, wanted not nil")
		} else {
			paramIdent := baseStruct.Params.List[1].Names[0].Name
			if paramIdent != "funcs" {
				t.Errorf("got %s, wanted funcs", paramIdent)
			}
			funcIdent := baseStruct.Params.List[1].Type.(*dst.ArrayType).Elt.(*dst.Ident)
			expectedPath := "github.com/myshkin5/moqueries/generator"
			if funcIdent.Path != expectedPath {
				t.Errorf("got %s, wanted %s", funcIdent, expectedPath)
			}
		}
	})

	t.Run("loads test types", func(t *testing.T) {
		// ASSEMBLE

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
