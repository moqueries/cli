package generator_test

import (
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"

	"moqueries.org/cli/generator"
	"moqueries.org/cli/pkg"
)

//go:generate moqueries --test-import testInterface

// testInterface verifies that the generator can access types in the test
// package
type testInterface interface {
	something()
}

func TestGenerating(t *testing.T) {
	const imp = "moqueries.org/cli/generator/testmoqs"

	t.Run("generates lots of different types of moqs which are then tested by testmoqs", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping generate test in short mode.")
		}

		// NB: Keep in sync with testmoqs/types.go go:generate directives

		// These lines generate the same moqs listed in types.go go:generate
		// directives. Remove the "pending" flag on this test to verify code
		// coverage.

		types := []string{
			"UsualFn",
			"NoNamesFn",
			"NoResultsFn",
			"NoParamsFn",
			"NothingFn",
			"VariadicFn",
			"RepeatedIdsFn",
			"TimesFn",
			"DifficultParamNamesFn",
			"DifficultResultNamesFn",
			"PassByArrayFn",
			"PassByChanFn",
			"PassByEllipsisFn",
			"PassByMapFn",
			"PassByReferenceFn",
			"PassBySliceFn",
			"PassByValueFn",
			"InterfaceParamFn",
			"InterfaceResultFn",
			"Usual",
		}

		err := generator.Generate(
			generator.GenerateRequest{
				Destination: "testmoqs/moq_testmoqs_test.go",
				Types:       types,
				Import:      imp,
			},
			generator.GenerateRequest{
				Destination: "testmoqs/exported/moq_exported_testmoqs.go",
				Export:      true,
				Types:       types,
				Import:      imp,
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}
	})

	t.Run("manual package test", func(t *testing.T) {
		t.Skip("manual test")
		err := pkg.Generate(pkg.PackageGenerateRequest{
			DestinationDir: "../../std", PkgPatterns: []string{"std"},
		})
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}
	})

	t.Run("dumps the DST of a moq", func(t *testing.T) {
		t.SkipNow()
		filePath := "./moq_usual_test.go"
		outPath := "./dst.txt"

		fSet := token.NewFileSet()
		// inFile, err := parser.ParseDir(fSet, filePath, nil, parser.AllErrors)
		inFile, err := parser.ParseFile(fSet, filePath, nil, parser.ParseComments)
		if err != nil {
			t.Errorf("got %#v, want no err", err)
		}

		// pkgs, err := packages.Load(&packages.Config{
		// 	Mode: packages.NeedName |
		// 		packages.NeedFiles |
		// 		packages.NeedCompiledGoFiles |
		// 		packages.NeedImports |
		// 		packages.NeedTypes |
		// 		packages.NeedSyntax |
		// 		packages.NeedTypesInfo |
		// 		packages.NeedTypesSizes,
		// 	Tests: false,
		// }, filePath)
		dstFile, err := decorator.DecorateFile(fSet, inFile)
		if err != nil {
			t.Errorf("got %#v, want no err", err)
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			t.Errorf("got %#v, want no err", err)
		}

		err = dst.Fprint(outFile, dstFile, dst.NotNilFilter)
		if err != nil {
			t.Errorf("got %#v, want no err", err)
		}
	})
}
