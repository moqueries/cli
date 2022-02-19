package ast_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"

	"github.com/myshkin5/moqueries/ast"
	"github.com/myshkin5/moqueries/moq"
)

var (
	scene     *moq.Scene
	loadFnMoq *moqLoadTypesFn

	cache *ast.Cache

	type1 *dst.TypeSpec
	type2 *dst.TypeSpec

	ioTypes []*dst.TypeSpec
)

func TestMain(m *testing.M) {
	var err error
	ioTypes, _, err = ast.LoadTypes("io", false)
	if err != nil {
		panic(fmt.Sprintf("Could not load io types: %#v", err))
	}
	os.Exit(m.Run())
}

func beforeEach(t *testing.T) {
	t.Helper()
	if scene != nil {
		t.Fatal("afterEach not called")
	}
	scene = moq.NewScene(t)
	loadFnMoq = newMoqLoadTypesFn(scene, &moq.Config{Sequence: moq.SeqDefaultOn})

	cache = ast.NewCache(loadFnMoq.mock())

	type1 = &dst.TypeSpec{Name: ast.Id("type1")}
	type2 = &dst.TypeSpec{Name: ast.Id("type2")}
}

func afterEach() {
	scene.AssertExpectationsMet()
	scene = nil
}

func TestTypeSimpleLoad(t *testing.T) {
	// ASSEMBLE
	beforeEach(t)

	loadFnMoq.onCall(".", true).
		returnResults([]*dst.TypeSpec{type1, type2}, "the-pkg", nil)

	// ACT
	actualType, actualPkg, actualErr := cache.Type(*ast.IdPath("type1", "."), true)

	// ASSERT
	if actualErr != nil {
		t.Errorf("Expected no error, was %#v", actualErr)
	}

	if actualType != type1 {
		t.Errorf("Expected type to be %#v, was %#v", type1, actualType)
	}

	if actualPkg != "the-pkg" {
		t.Errorf("Expected pkg to be the-pkg, was %s", actualPkg)
	}

	afterEach()
}

func TestTypeLoadError(t *testing.T) {
	// ASSEMBLE
	beforeEach(t)

	err := errors.New("load error")
	loadFnMoq.onCall(".", true).
		returnResults(nil, "", err)

	// ACT
	actualType, actualPkg, actualErr := cache.Type(*ast.IdPath("type1", "."), true)

	// ASSERT
	if actualErr != err {
		t.Errorf("Expected error %#v, was %#v", err, actualErr)
	}

	if actualType != nil {
		t.Errorf("Expected nil type, was %#v", actualType)
	}

	if actualPkg != "" {
		t.Errorf("Expected empty pkg, was %s", actualPkg)
	}

	afterEach()
}

func TestTypeNotFound(t *testing.T) {
	// ASSEMBLE
	beforeEach(t)

	loadFnMoq.onCall("the-pkg", false).
		returnResults([]*dst.TypeSpec{type1, type2}, "the-pkg", nil)

	// ACT
	actualType, actualPkg, actualErr := cache.Type(*ast.IdPath("type3", "the-pkg"), false)

	// ASSERT
	if actualErr == nil || !strings.Contains(actualErr.Error(), "not found") {
		t.Errorf("Expected error to contain 'not found', was %#v", actualErr)
	}

	if actualType != nil {
		t.Errorf("Expected nil type, was %#v", actualType)
	}

	if actualPkg != "" {
		t.Errorf("Expected empty pkg, was %s", actualPkg)
	}

	afterEach()
}

func TestTypeOnlyLoadPackageOnce(t *testing.T) {
	// ASSEMBLE
	beforeEach(t)

	loadFnMoq.onCall("the-pkg", true).
		returnResults([]*dst.TypeSpec{type1, type2}, "the-pkg", nil)
	_, _, _ = cache.Type(*ast.IdPath("type1", "the-pkg"), true)

	// ACT
	actualType, actualPkg, actualErr := cache.Type(*ast.IdPath("type2", "the-pkg"), false)

	// ASSERT
	if actualErr != nil {
		t.Errorf("Expected no error, was %#v", actualErr)
	}

	if actualType != type2 {
		t.Errorf("Expected type to be %#v, was %#v", type2, actualType)
	}

	if actualPkg != "the-pkg" {
		t.Errorf("Expected pkg to be the-pkg, was %s", actualPkg)
	}

	afterEach()
}

func TestTypeReloadTestPackage(t *testing.T) {
	// ASSEMBLE
	beforeEach(t)

	loadFnMoq.onCall("the-pkg", false).
		returnResults([]*dst.TypeSpec{type1}, "the-pkg", nil)
	_, _, _ = cache.Type(*ast.IdPath("type1", "the-pkg"), false)
	loadFnMoq.onCall("the-pkg", true).
		returnResults([]*dst.TypeSpec{type1, type2}, "the-pkg", nil)

	// ACT
	actualType, actualPkg, actualErr := cache.Type(*ast.IdPath("type2", "the-pkg"), true)

	// ASSERT
	if actualErr != nil {
		t.Errorf("Expected no error, was %#v", actualErr)
	}

	if actualType != type2 {
		t.Errorf("Expected type to be %#v, was %#v", type2, actualType)
	}

	if actualPkg != "the-pkg" {
		t.Errorf("Expected pkg to be the-pkg, was %s", actualPkg)
	}

	afterEach()
}

func TestTypeOnlyLoadDefaultPackageOnce(t *testing.T) {
	// ASSEMBLE
	beforeEach(t)

	loadFnMoq.onCall(".", true).
		returnResults([]*dst.TypeSpec{type1, type2}, "the-pkg", nil)
	_, _, _ = cache.Type(*ast.IdPath("type1", "."), true)

	// ACT
	actualType, actualPkg, actualErr := cache.Type(*ast.IdPath("type2", "the-pkg"), false)

	// ASSERT
	if actualErr != nil {
		t.Errorf("Expected no error, was %#v", actualErr)
	}

	if actualType != type2 {
		t.Errorf("Expected type to be %#v, was %#v", type2, actualType)
	}

	if actualPkg != "the-pkg" {
		t.Errorf("Expected pkg to be the-pkg, was %s", actualPkg)
	}

	afterEach()
}

type tableEntry struct {
	paramType         string
	comparable        bool
	defaultComparable bool
	structable        bool
}

var entries = []tableEntry{
	{
		paramType:         "string",
		comparable:        true,
		defaultComparable: true,
		structable:        true,
	},
	{
		paramType:         "[]string",
		comparable:        false,
		defaultComparable: false,
		structable:        true,
	},
	{
		paramType:         "[3]string",
		comparable:        true,
		defaultComparable: true,
		structable:        true,
	},
	{
		paramType:         "map[string]string",
		comparable:        false,
		defaultComparable: false,
		structable:        true,
	},
	{
		paramType:         "...string",
		comparable:        false,
		defaultComparable: false,
		structable:        false,
	},
	{
		paramType:         "*string",
		comparable:        true,
		defaultComparable: false,
		structable:        true,
	},
	{
		paramType:         "error",
		comparable:        true,
		defaultComparable: false,
		structable:        true,
	},
	{
		paramType:         "[3]error",
		comparable:        true,
		defaultComparable: false,
		structable:        true,
	},
	{
		paramType:         "io.Reader",
		comparable:        true,
		defaultComparable: false,
		structable:        true,
	},
	{
		paramType:         "[3]io.Reader",
		comparable:        true,
		defaultComparable: false,
		structable:        true,
	},
}

func simpleExpr(t *testing.T, paramType string) dst.Expr {
	t.Helper()
	code := `package a

import _ "io"

func b(c %s) {}
`
	f := parse(t, fmt.Sprintf(code, paramType))
	fn, ok := f.Decls[1].(*dst.FuncDecl)
	if !ok {
		t.Fatalf("wanted a function declaration, got %#v", f.Decls[1])
	}
	return fn.Type.Params.List[0].Type
}

func TestIsComparableSimpleExprs(t *testing.T) {
	for _, entry := range entries {
		t.Run(entry.paramType, func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)

			loadFnMoq.onCall("io", false).
				returnResults(ioTypes, "io", nil).repeat(moq.AnyTimes())

			// ACT
			comparable, err := cache.IsComparable(simpleExpr(t, entry.paramType))
			// ASSERT
			if err != nil {
				t.Errorf("Unexpected error checking comparable: %s, err: %#v", entry.paramType, err)
			}
			if comparable != entry.comparable {
				t.Errorf("IsComparable(%s) = %t; want %t",
					entry.paramType, comparable, entry.comparable)
			}

			afterEach()
		})
	}
}

func TestIsDefaultComparableSimpleExprs(t *testing.T) {
	for _, entry := range entries {
		t.Run(entry.paramType, func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)

			loadFnMoq.onCall("io", false).
				returnResults(ioTypes, "io", nil).repeat(moq.AnyTimes())

			// ACT
			comparable, err := cache.IsDefaultComparable(simpleExpr(t, entry.paramType))
			// ASSERT
			if err != nil {
				t.Errorf("Unexpected error checking comparable: %s, err: %#v", entry.paramType, err)
			}
			if comparable != entry.defaultComparable {
				t.Errorf("IsComparable(%s) = %t; want %t",
					entry.paramType, comparable, entry.defaultComparable)
			}

			afterEach()
		})
	}
}

type comparableStructEntry struct {
	code    string
	declIdx int
}

var comparableStructEntries = map[string]comparableStructEntry{
	"inline": {
		code: `package a

import _ "io"

func a(b struct{ c %s }) {}
`,
		declIdx: 1,
	},
	"standard": {
		code: `package a

import _ "io"

type b struct {
	c %s
}

func d(e b) {}
`,
		declIdx: 2,
	},
}

func TestIsComparableStructExprs(t *testing.T) {
	for name, sEntry := range comparableStructEntries {
		t.Run(name, func(t *testing.T) {
			for _, entry := range entries {
				t.Run(entry.paramType, func(t *testing.T) {
					// ASSEMBLE
					if !entry.structable {
						t.Skipf("%s can't be put into a struct, skipping", entry.paramType)
					}

					beforeEach(t)

					f := parse(t, fmt.Sprintf(sEntry.code, entry.paramType))
					fn, ok := f.Decls[sEntry.declIdx].(*dst.FuncDecl)
					if !ok {
						t.Fatalf("wanted a function declaration, got %#v", f.Decls[sEntry.declIdx])
					}
					expr := fn.Type.Params.List[0].Type

					loadFnMoq.onCall("io", false).
						returnResults(ioTypes, "io", nil).repeat(moq.AnyTimes())

					// ACT
					comparable, err := cache.IsComparable(expr)
					// ASSERT
					if err != nil {
						t.Errorf("Unexpected error checking comparable: %s, err: %#v", sEntry.code, err)
					}
					if comparable != entry.comparable {
						t.Errorf("IsComparable(%s) = %t; want %t",
							entry.paramType, comparable, entry.comparable)
					}

					afterEach()
				})
			}
		})
	}
}

func TestIsComparableImported(t *testing.T) {
	for _, entry := range entries {
		t.Run(entry.paramType, func(t *testing.T) {
			// ASSEMBLE
			if !entry.structable {
				t.Skipf("%s can't be put into a struct, skipping", entry.paramType)
			}

			beforeEach(t)

			code1 := `package a

import "b"

func c(d b.e) {}
`
			f1 := parse(t, code1)
			fn, ok := f1.Decls[1].(*dst.FuncDecl)
			if !ok {
				t.Fatalf("wanted a function declaration, got %#v", f1.Decls[1])
			}
			expr := fn.Type.Params.List[0].Type

			code2 := `package b

import _ "io"

type e struct {
	f %s
}
`
			f2 := parse(t, fmt.Sprintf(code2, entry.paramType))
			gen, ok := f2.Decls[1].(*dst.GenDecl)
			if !ok {
				t.Fatalf("wanted a generic declaration, got %#v", f2.Decls[1])
			}
			typ, ok := gen.Specs[0].(*dst.TypeSpec)
			if !ok {
				t.Fatalf("wanted a type, got %#v", gen.Specs[0])
			}
			loadFnMoq.onCall("b", false).
				returnResults(
					[]*dst.TypeSpec{typ},
					"b",
					nil)

			loadFnMoq.onCall("io", false).
				returnResults(ioTypes, "io", nil).repeat(moq.AnyTimes())

			// ACT
			comparable, err := cache.IsComparable(expr)
			// ASSERT
			if err != nil {
				t.Errorf("Unexpected error checking comparable: %s, err: %#v", code2, err)
			}
			if comparable != entry.comparable {
				t.Errorf("IsComparable(%s) = %t; want %t",
					entry.paramType, comparable, entry.comparable)
			}

			afterEach()
		})
	}
}

func TestDSTIdentNotComparable(t *testing.T) {
	// ASSEMBLE
	cache := ast.NewCache(ast.LoadTypes)
	typ, _, err := cache.Type(
		*ast.IdPath("TypeCache", "github.com/myshkin5/moqueries/generator"), false)
	if err != nil {
		t.Errorf("Unexpected error loading type, err: %#v", err)
	}
	iface, ok := typ.Type.(*dst.InterfaceType)
	if !ok {
		t.Fatalf("wanted an interface, got %#v", typ.Type)
	}
	fType, ok := iface.Methods.List[0].Type.(*dst.FuncType)
	if !ok {
		t.Fatalf("wanted a function type, got %#v", iface.Methods.List[0].Type)
	}
	expr := fType.Params.List[0].Type

	// ACT
	comparable, err := cache.IsComparable(expr)
	// ASSERT
	if err != nil {
		t.Errorf("Unexpected error checking comparable, err: %#v", err)
	}
	if comparable {
		t.Errorf("IsComparable = %t; want false", comparable)
	}
}

func parse(t *testing.T, code string) *dst.File {
	t.Helper()
	f, err := decorator.Parse(code)
	if err != nil {
		t.Errorf("Unexpected error parsing code: %s, err: %#v", code, err)
	}
	return f
}
