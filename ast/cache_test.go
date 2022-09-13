package ast_test

import (
	"errors"
	"fmt"
	goAst "go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"golang.org/x/tools/go/packages"

	"github.com/myshkin5/moqueries/ast"
	"github.com/myshkin5/moqueries/metrics"
	"github.com/myshkin5/moqueries/moq"
)

var ioTypes []*packages.Package

func TestMain(m *testing.M) {
	var err error
	ioTypes, err = packages.Load(&packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes,
		Tests: false,
	}, "io")
	if err != nil {
		panic(fmt.Sprintf("Could not load io types: %#v", err))
	}
	os.Exit(m.Run())
}

func TestCache(t *testing.T) {
	var (
		scene      *moq.Scene
		loadFnMoq  *moqLoadFn
		metricsMoq *metrics.MoqMetrics

		cache *ast.Cache

		loadCfg *packages.Config

		pkgs       []*packages.Package
		builtinPkg *packages.Package
	)

	pkg := func(pkgPath string, testPkg, export bool, fs *token.FileSet) *packages.Package {
		typePrefix := ""
		if testPkg {
			pkgPath += "_test"
			typePrefix = "test_"
			if export {
				typePrefix = "Test_"
			}
		}
		name1 := "type1"
		name2 := "type2"
		name3 := "type3"
		if export {
			name1 = "Type1"
			name2 = "Type2"
			name3 = "Type3"
		}

		return &packages.Package{
			Syntax: []*goAst.File{
				{
					Package: 1,
					Decls: []goAst.Decl{
						&goAst.GenDecl{
							Specs: []goAst.Spec{
								&goAst.TypeSpec{
									Name: goAst.NewIdent(typePrefix + name1),
									Type: &goAst.InterfaceType{
										Methods: &goAst.FieldList{},
									},
								},
							},
						},
					},
				},
				{
					Package: 2,
					Decls: []goAst.Decl{
						&goAst.GenDecl{
							Specs: []goAst.Spec{
								&goAst.TypeSpec{
									Name: goAst.NewIdent(typePrefix + name2),
									Type: &goAst.FuncType{
										Params: &goAst.FieldList{},
									},
								},
							},
						},
					},
				},
				{
					Package: 2,
					Decls: []goAst.Decl{
						&goAst.GenDecl{
							Specs: []goAst.Spec{
								&goAst.TypeSpec{
									Name: goAst.NewIdent(typePrefix + name3),
									Type: &goAst.StructType{
										Fields: &goAst.FieldList{},
									},
								},
							},
						},
					},
				},
			},
			TypesInfo: &types.Info{},
			Fset:      fs,
			GoFiles:   []string{"file1", "file2"},
			PkgPath:   pkgPath,
		}
	}
	beforeEach := func(t *testing.T, testImport, export bool) {
		t.Helper()
		if scene != nil {
			t.Fatal("afterEach not called")
		}
		scene = moq.NewScene(t)
		loadFnMoq = newMoqLoadFn(scene, nil)
		metricsMoq = metrics.NewMoqMetrics(scene, nil)

		cache = ast.NewCache(loadFnMoq.mock(), metricsMoq.Mock())

		loadCfg = &packages.Config{
			Mode: packages.NeedName |
				packages.NeedFiles |
				packages.NeedCompiledGoFiles |
				packages.NeedImports |
				packages.NeedTypes |
				packages.NeedSyntax |
				packages.NeedTypesInfo |
				packages.NeedTypesSizes,
			Tests: testImport,
		}

		fs := token.NewFileSet()
		fs.AddFile("file1", 1, 0)
		fs.AddFile("file2", 2, 0)
		pkgs = []*packages.Package{
			pkg("the_pkg", false, export, fs),
		}
		if testImport {
			// When loading test packages, the regular package typically
			// precedes the test package
			pkgs = append(pkgs, pkg("the_pkg", true, export, fs))
		}

		builtinPkg = &packages.Package{
			Syntax: []*goAst.File{
				{
					Package: 1,
					Decls: []goAst.Decl{
						&goAst.GenDecl{
							Specs: []goAst.Spec{
								&goAst.TypeSpec{
									Name: goAst.NewIdent("error"),
									Type: &goAst.InterfaceType{
										Methods: &goAst.FieldList{},
									},
								},
							},
						},
					},
				},
			},
			TypesInfo: &types.Info{},
			Fset:      fs,
			GoFiles:   []string{"file1"},
			PkgPath:   "builtin",
		}
	}
	afterEach := func() {
		scene.AssertExpectationsMet()
		scene = nil
	}

	t.Run("Type", func(t *testing.T) {
		t.Run("simple load", func(t *testing.T) {
			type testCase struct {
				typeToLoad, expectedPkg string
			}
			testCases := map[bool]testCase{
				false: {typeToLoad: "type1", expectedPkg: "the_pkg"},
				true:  {typeToLoad: "test_type1", expectedPkg: "the_pkg_test"},
			}

			for testImport, tc := range testCases {
				t.Run(fmt.Sprintf("testImport: %t", testImport), func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t, testImport, false)
					defer afterEach()

					metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
					metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
					loadFnMoq.onCall(loadCfg, ".").returnResults(pkgs, nil)
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
					metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
						ReturnResults().Repeat(moq.MinTimes(1))

					id := ast.IdPath(tc.typeToLoad, ".")

					// ACT
					actualType, actualPkg, actualErr := cache.Type(*id, testImport)

					// ASSERT
					if actualErr != nil {
						t.Fatalf("got %#v, want no error", actualErr)
					}

					if actualType.Name.Name != tc.typeToLoad {
						t.Errorf("got %#v, want %s", actualType.Name.Name, tc.typeToLoad)
					}

					if actualPkg != tc.expectedPkg {
						t.Errorf("got %s, want %s", actualPkg, tc.expectedPkg)
					}
				})
			}
		})

		t.Run("loads test package when given a test package", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, true, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "the_pkg").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
				ReturnResults().Repeat(moq.MinTimes(1))

			id := ast.IdPath("test_type1", "the_pkg_test")

			// ACT
			actualType, actualPkg, actualErr := cache.Type(*id, false)

			// ASSERT
			if actualErr != nil {
				t.Fatalf("got %#v, want no error", actualErr)
			}

			if actualType.Name.Name != "test_type1" {
				t.Errorf("got %#v, want test_type1", actualType.Name.Name)
			}

			if actualPkg != "the_pkg_test" {
				t.Errorf("got %s, want the_pkg_test", actualPkg)
			}
		})

		t.Run("load error", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, true, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			err := errors.New("load error")
			loadFnMoq.onCall(loadCfg, ".").returnResults(nil, err)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()

			id := ast.IdPath("type1", ".")

			// ACT
			actualType, actualPkg, actualErr := cache.Type(*id, true)

			// ASSERT
			if actualErr != err {
				t.Errorf("got %#v, want %#v", actualErr, err)
			}

			if actualType != nil {
				t.Errorf("got %#v, want nil type", actualType)
			}

			if actualPkg != "" {
				t.Errorf("got %s, want empty pkg", actualPkg)
			}
		})

		t.Run("not found", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "the_pkg").
				returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			id := ast.IdPath("type4", "the_pkg")

			// ACT
			actualType, actualPkg, actualErr := cache.Type(*id, false)

			// ASSERT
			if actualErr == nil || !strings.Contains(actualErr.Error(), "not found") {
				t.Errorf("got %#v, want to contain 'not found'", actualErr)
			}
			if !errors.Is(actualErr, ast.ErrTypeNotFound) {
				t.Errorf("got %#v, want ast.ErrTypeNotFound", actualErr)
			}

			if actualType != nil {
				t.Errorf("got %#v, want nil type", actualType)
			}

			if actualPkg != "" {
				t.Errorf("got %s, want empty pkg", actualPkg)
			}
		})

		t.Run("load packages only once", func(t *testing.T) {
			testCases := map[string]string{
				"any package":     "the_pkg",
				"default package": ".",
			}

			for name, pkgName := range testCases {
				t.Run(name, func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t, true, false)
					defer afterEach()

					metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
					metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
					loadFnMoq.onCall(loadCfg, pkgName).returnResults(pkgs, nil)
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
					metricsMoq.OnCall().ASTPkgCacheHitsInc().ReturnResults()
					metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
						ReturnResults().Repeat(moq.MinTimes(1))
					_, _, _ = cache.Type(*ast.IdPath("type1", pkgName), true)

					id := ast.IdPath("type2", "the_pkg")

					// ACT
					actualType, actualPkg, actualErr := cache.Type(*id, false)

					// ASSERT
					if actualErr != nil {
						t.Errorf("got %#v, want no error", actualErr)
					}

					if actualType.Name.Name != "type2" {
						t.Errorf("got %#v, want type2", actualType.Name.Name)
					}

					if actualPkg != "the_pkg" {
						t.Errorf("got %s, want 'the_pkg'", actualPkg)
					}
				})
			}
		})

		t.Run("reload test package", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "the_pkg").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
				ReturnResults().Repeat(moq.MinTimes(1))
			_, _, _ = cache.Type(*ast.IdPath("type1", "the_pkg"), false)

			loadCfg.Tests = true

			fs := token.NewFileSet()
			fs.AddFile("file1", 1, 0)
			fs.AddFile("file2", 2, 0)
			pkgs = append(pkgs, pkg("the_pkg", true, false, fs))

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "the_pkg").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()

			id := ast.IdPath("test_type2", "the_pkg")

			// ACT
			actualType, actualPkg, actualErr := cache.Type(*id, true)

			// ASSERT
			if actualErr != nil {
				t.Fatalf("got %#v, want no error", actualErr)
			}

			if actualType.Name.Name != "test_type2" {
				t.Errorf("got %#v, want test_type2", actualType.Name.Name)
			}

			if actualPkg != "the_pkg_test" {
				t.Errorf("got %s, want 'the_pkg_test'", actualPkg)
			}
		})

		t.Run("knows how to load builtin error type", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().
				ReturnResults().Repeat(moq.Times(2))
			metricsMoq.OnCall().ASTTypeCacheMissesInc().
				ReturnResults().Repeat(moq.Times(2))
			loadFnMoq.onCall(loadCfg, ".").returnResults(pkgs, nil)
			loadFnMoq.onCall(loadCfg, "builtin").
				returnResults([]*packages.Package{builtinPkg}, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().
				ReturnResults().Repeat(moq.Times(2))
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
				ReturnResults().Repeat(moq.Times(2))

			id := ast.IdPath("error", ".")

			// ACT
			actualType, actualPkg, actualErr := cache.Type(*id, false)

			// ASSERT
			if actualErr != nil {
				t.Fatalf("got %#v, want no error", actualErr)
			}

			if actualType.Name.Name != "error" {
				t.Errorf("got %#v, want error", actualType.Name.Name)
			}

			if actualPkg != "" {
				t.Errorf("got %s, want empty package name", actualPkg)
			}
		})
	})

	t.Run("IsComparable/IsDefaultComparable", func(t *testing.T) {
		type testCase struct {
			comparable        bool
			defaultComparable bool
			structable        bool
		}

		testCases := map[string]testCase{
			"string": {
				comparable:        true,
				defaultComparable: true,
				structable:        true,
			},
			"[]string": {
				comparable:        false,
				defaultComparable: false,
				structable:        true,
			},
			"[3]string": {
				comparable:        true,
				defaultComparable: true,
				structable:        true,
			},
			"map[string]string": {
				comparable:        false,
				defaultComparable: false,
				structable:        true,
			},
			"...string": {
				comparable:        false,
				defaultComparable: false,
				structable:        false,
			},
			"*string": {
				comparable:        true,
				defaultComparable: false,
				structable:        true,
			},
			"error": {
				comparable:        true,
				defaultComparable: false,
				structable:        true,
			},
			"[3]error": {
				comparable:        true,
				defaultComparable: false,
				structable:        true,
			},
			"io.Reader": {
				comparable:        true,
				defaultComparable: false,
				structable:        true,
			},
			"[3]io.Reader": {
				comparable:        true,
				defaultComparable: false,
				structable:        true,
			},
			"func()": {
				comparable:        false,
				defaultComparable: false,
				structable:        false,
			},
		}

		parse := func(t *testing.T, code string) *dst.File {
			t.Helper()
			f, err := decorator.Parse(code)
			if err != nil {
				t.Errorf("got err: %#v, want no error", err)
			}
			return f
		}

		simpleExpr := func(t *testing.T, paramType string) dst.Expr {
			t.Helper()
			code := `package a

import _ "io"

func b(c %s) {}
`
			f := parse(t, fmt.Sprintf(code, paramType))
			fn, ok := f.Decls[1].(*dst.FuncDecl)
			if !ok {
				t.Fatalf("got %#v, want a function declaration", f.Decls[1])
			}
			return fn.Type.Params.List[0].Type
		}

		parseASTPackage := func(t *testing.T, code, pkgPath string) []*packages.Package {
			t.Helper()

			dir, err := os.MkdirTemp("", "moq-ast-cache-*")
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}
			defer func() {
				err = os.RemoveAll(dir)
				if err != nil {
					t.Fatalf("got %#v, want no error", err)
				}
			}()
			file := filepath.Join(dir, "code.go")
			err = os.WriteFile(file, []byte(code), 0o600)
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			p, err := packages.Load(&packages.Config{
				Mode: packages.NeedName |
					packages.NeedFiles |
					packages.NeedCompiledGoFiles |
					packages.NeedImports |
					packages.NeedTypes |
					packages.NeedSyntax |
					packages.NeedTypesInfo |
					packages.NeedTypesSizes,
				Tests: false,
			}, file)
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}
			p[0].PkgPath = pkgPath

			return p
		}

		t.Run("simple exprs", func(t *testing.T) {
			for paramType, tc := range testCases {
				t.Run(paramType, func(t *testing.T) {
					isComparable := func(c *ast.Cache, expr dst.Expr) (bool, error) {
						return c.IsComparable(expr)
					}
					isDefaultComparable := func(c *ast.Cache, expr dst.Expr) (bool, error) {
						return c.IsDefaultComparable(expr)
					}

					type subTestCase struct {
						compFn     func(c *ast.Cache, expr dst.Expr) (bool, error)
						comparable bool
					}
					subTestCases := map[string]subTestCase{
						"IsComparable": {
							compFn:     isComparable,
							comparable: tc.comparable,
						},
						"IsDefaultComparable": {
							compFn:     isDefaultComparable,
							comparable: tc.defaultComparable,
						},
					}

					for name, stc := range subTestCases {
						t.Run(name, func(t *testing.T) {
							// ASSEMBLE
							beforeEach(t, false, false)
							defer afterEach()

							metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.Optional())
							metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults().Repeat(moq.Optional())
							loadFnMoq.onCall(loadCfg, "io").
								returnResults(ioTypes, nil).repeat(moq.AnyTimes())
							metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults().Repeat(moq.Optional())
							metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
								ReturnResults().Repeat(moq.AnyTimes())

							expr := simpleExpr(t, paramType)

							// ACT
							is, err := stc.compFn(cache, expr)
							// ASSERT
							if err != nil {
								t.Errorf("got %#v, want no error", err)
							}
							if is != stc.comparable {
								t.Errorf("got %t, want %t", is, stc.comparable)
							}
						})
					}
				})
			}
		})

		type comparableStructCase struct {
			code    string
			declIdx int
		}

		comparableStructCases := map[string]comparableStructCase{
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

		t.Run("IsComparable - struct exprs", func(t *testing.T) {
			for name, stc := range comparableStructCases {
				t.Run(name, func(t *testing.T) {
					for paramType, tc := range testCases {
						t.Run(paramType, func(t *testing.T) {
							// ASSEMBLE
							if !tc.structable {
								t.Skipf("%s can't be put into a struct, skipping", paramType)
							}

							beforeEach(t, false, false)
							defer afterEach()

							f := parse(t, fmt.Sprintf(stc.code, paramType))
							fn, ok := f.Decls[stc.declIdx].(*dst.FuncDecl)
							if !ok {
								t.Fatalf("got %#v, want a function declaration", f.Decls[stc.declIdx])
							}
							expr := fn.Type.Params.List[0].Type

							metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.Optional())
							metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults().Repeat(moq.Optional())
							loadFnMoq.onCall(loadCfg, "io").
								returnResults(ioTypes, nil).repeat(moq.AnyTimes())
							metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults().Repeat(moq.Optional())
							metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
								ReturnResults().Repeat(moq.AnyTimes())

							// ACT
							isComparable, err := cache.IsComparable(expr)
							// ASSERT
							if err != nil {
								t.Errorf("got %#v, want no error", err)
							}
							if isComparable != tc.comparable {
								t.Errorf("got %t, want %t", isComparable, tc.comparable)
							}
						})
					}
				})
			}
		})

		t.Run("IsComparable - imported", func(t *testing.T) {
			for paramType, tc := range testCases {
				t.Run(paramType, func(t *testing.T) {
					// ASSEMBLE
					if !tc.structable {
						t.Skipf("%s can't be put into a struct, skipping", paramType)
					}

					beforeEach(t, false, false)
					defer afterEach()

					code1 := `package a

import "b"

func c(d b.e) {}
`
					f1 := parse(t, code1)
					fn, ok := f1.Decls[1].(*dst.FuncDecl)
					if !ok {
						t.Fatalf("got %#v, want a function declaration", f1.Decls[1])
					}
					expr := fn.Type.Params.List[0].Type

					code2 := `package b

import _ "io"

type e struct {
	f %s
}
`
					bPkg := parseASTPackage(t, fmt.Sprintf(code2, paramType), "b")
					metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.Optional())
					metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults().Repeat(moq.Optional())
					loadFnMoq.onCall(loadCfg, "b").returnResults(bPkg, nil)
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults().Repeat(moq.Optional())

					metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.Optional())
					metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults().Repeat(moq.Optional())
					loadFnMoq.onCall(loadCfg, "io").
						returnResults(ioTypes, nil).repeat(moq.AnyTimes())
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults().Repeat(moq.Optional())
					metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
						ReturnResults().Repeat(moq.AnyTimes())

					// ACT
					isComparable, err := cache.IsComparable(expr)
					// ASSERT
					if err != nil {
						t.Errorf("got %#v, want no error", err)
					}
					if isComparable != tc.comparable {
						t.Errorf("got %t, want %t", isComparable, tc.comparable)
					}
				})
			}
		})
	})

	t.Run("DST ident not comparable", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t, false, false)
		defer afterEach()

		c := ast.NewCache(packages.Load, metricsMoq.Mock())

		metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.AnyTimes())
		metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults().Repeat(moq.AnyTimes())
		metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().
			ReturnResults().Repeat(moq.AnyTimes())
		metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
			ReturnResults().Repeat(moq.AnyTimes())

		typ, _, err := c.Type(
			*ast.IdPath("TypeCache", "github.com/myshkin5/moqueries/generator"), false)
		if err != nil {
			t.Fatalf("got %#v, want no error", err)
		}
		iface, ok := typ.Type.(*dst.InterfaceType)
		if !ok {
			t.Fatalf("got %#v, want an interface", typ.Type)
		}
		fType, ok := iface.Methods.List[0].Type.(*dst.FuncType)
		if !ok {
			t.Fatalf("got %#v, want a function type", iface.Methods.List[0].Type)
		}
		expr := fType.Params.List[0].Type

		// ACT
		isComparable, err := c.IsComparable(expr)
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, want no error", err)
		}
		if isComparable {
			t.Errorf("got %t, want false", isComparable)
		}
	})

	t.Run("FindPackage", func(t *testing.T) {
		t.Run("relative dir", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "./the_pkg").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			// ACT
			pkgPath, err := cache.FindPackage("the_pkg")
			// ASSERT
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			if pkgPath != "the_pkg" {
				t.Errorf("got %s, want the_pkg", pkgPath)
			}
		})

		t.Run("abs dir", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "/this-dir/the_pkg").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			// ACT
			pkgPath, err := cache.FindPackage("/this-dir/the_pkg")
			// ASSERT
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			if pkgPath != "the_pkg" {
				t.Errorf("got %s, want the_pkg", pkgPath)
			}
		})

		t.Run("current dir", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, ".").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			// ACT
			pkgPath, err := cache.FindPackage(".")
			// ASSERT
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			if pkgPath != "the_pkg" {
				t.Errorf("got %s, want the_pkg", pkgPath)
			}
		})
	})

	t.Run("LoadPackage", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, ".").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			// ACT
			err := cache.LoadPackage(".")
			// ASSERT
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}
		})

		t.Run("error", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			err := errors.New("load error")
			loadFnMoq.onCall(loadCfg, ".").returnResults(nil, err)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()

			// ACT
			actualErr := cache.LoadPackage(".")
			// ASSERT
			if actualErr != err {
				t.Errorf("got %#v, want %#v", actualErr, err)
			}
		})
	})

	t.Run("MockableTypes", func(t *testing.T) {
		t.Run("everything returned", func(t *testing.T) {
			type testCase struct {
				export       bool
				onlyExported bool
				prefix       string
			}
			testCases := map[string]testCase{
				"no filtering": {export: false, onlyExported: false, prefix: "type"},
				"filtering":    {export: true, onlyExported: true, prefix: "Type"},
			}

			for name, tc := range testCases {
				t.Run(name, func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t, false, tc.export)
					defer afterEach()

					metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
					metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
					loadFnMoq.onCall(loadCfg, ".").returnResults(pkgs, nil)
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
					metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

					err := cache.LoadPackage(".")
					if err != nil {
						t.Fatalf("got %#v, want no error", err)
					}

					// ACT
					typs := cache.MockableTypes(tc.onlyExported)

					// ASSERT
					if len(typs) != 2 {
						t.Fatalf("got %d types, want 2", len(typs))
					}

					typ1 := typs[0]
					typ2 := typs[1]
					if typ1.Name != tc.prefix+"1" {
						typ1 = typs[1]
						typ2 = typs[0]
					}
					if typ1.Name != tc.prefix+"1" {
						t.Errorf("got %s, want type1", typ1.Name)
					}
					if typ2.Name != tc.prefix+"2" {
						t.Errorf("got %s, want type2", typ2.Name)
					}
				})
			}
		})

		t.Run("filtering to nothing", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, ".").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			err := cache.LoadPackage(".")
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			// ACT
			typs := cache.MockableTypes(true)

			// ASSERT
			if len(typs) != 0 {
				t.Fatalf("got %d types, want none", len(typs))
			}
		})

		t.Run("always filters vendor packages", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false, false)
			defer afterEach()

			fs := token.NewFileSet()
			fs.AddFile("file1", 1, 0)
			fs.AddFile("file2", 2, 0)
			pkgs = append(pkgs, pkg("vendor/other-pkgs", true, true, fs))

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, ".").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults().Repeat(moq.AnyTimes())

			err := cache.LoadPackage(".")
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			// ACT
			typs := cache.MockableTypes(true)

			// ASSERT
			if len(typs) != 0 {
				t.Fatalf("got %d types, want none", len(typs))
			}
		})

		t.Run("always filters internal packages", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false, false)
			defer afterEach()

			fs := token.NewFileSet()
			fs.AddFile("file1", 1, 0)
			fs.AddFile("file2", 2, 0)
			pkgs = append(pkgs, pkg("elsewhere/internal", false, true, fs))
			pkgs = append(pkgs, pkg("elsewhere/internal/again", false, true, fs))
			pkgs = append(pkgs, pkg("internal", false, true, fs))

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, ".").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults().Repeat(moq.AnyTimes())

			err := cache.LoadPackage(".")
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			// ACT
			typs := cache.MockableTypes(true)

			// ASSERT
			if len(typs) != 0 {
				t.Fatalf("got %d types, want none", len(typs))
			}
		})
	})
}
