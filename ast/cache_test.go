package ast_test

import (
	"errors"
	"fmt"
	goAst "go/ast"
	"go/token"
	"go/types"
	"os"
	"path"
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

		pkgs []*packages.Package
	)

	beforeEach := func(t *testing.T, loadTestPkgs bool) {
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
			Tests: loadTestPkgs,
		}

		fs := token.NewFileSet()
		fs.AddFile("file1", 1, 0)
		fs.AddFile("file2", 2, 0)
		pkgs = []*packages.Package{
			{
				Syntax: []*goAst.File{
					{
						Package: 1,
						Decls: []goAst.Decl{
							&goAst.GenDecl{
								Specs: []goAst.Spec{
									&goAst.TypeSpec{
										Name: goAst.NewIdent("type1"),
										Type: &goAst.StructType{
											Fields: &goAst.FieldList{},
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
										Name: goAst.NewIdent("type2"),
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
				PkgPath:   "the-pkg",
			},
		}
	}

	afterEach := func() {
		scene.AssertExpectationsMet()
		scene = nil
	}

	t.Run("Type", func(t *testing.T) {
		t.Run("simple load", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, true)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, ".").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()

			id := ast.IdPath("type1", ".")

			// ACT
			actualType, actualPkg, actualErr := cache.Type(*id, true)

			// ASSERT
			if actualErr != nil {
				t.Fatalf("got %#v, want no error", actualErr)
			}

			if actualType.Name.Name != "type1" {
				t.Errorf("got %#v, want type1", actualType.Name.Name)
			}

			if actualPkg != "the-pkg" {
				t.Errorf("got %s, want the-pkg", actualPkg)
			}
		})

		t.Run("load error", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, true)
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
			beforeEach(t, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "the-pkg").
				returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()

			id := ast.IdPath("type3", "the-pkg")

			// ACT
			actualType, actualPkg, actualErr := cache.Type(*id, false)

			// ASSERT
			if actualErr == nil || !strings.Contains(actualErr.Error(), "not found") {
				t.Errorf("got %#v, want to contain 'not found'", actualErr)
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
				"any package":     "the-pkg",
				"default package": ".",
			}

			for name, pkg := range testCases {
				t.Run(name, func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t, true)
					defer afterEach()

					metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
					metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
					loadFnMoq.onCall(loadCfg, pkg).returnResults(pkgs, nil)
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
					metricsMoq.OnCall().ASTPkgCacheHitsInc().ReturnResults()
					_, _, _ = cache.Type(*ast.IdPath("type1", pkg), true)

					id := ast.IdPath("type2", "the-pkg")

					// ACT
					actualType, actualPkg, actualErr := cache.Type(*id, false)

					// ASSERT
					if actualErr != nil {
						t.Errorf("got %#v, want no error", actualErr)
					}

					if actualType.Name.Name != "type2" {
						t.Errorf("got %#v, want type2", actualType.Name.Name)
					}

					if actualPkg != "the-pkg" {
						t.Errorf("got %s, want 'the-pkg'", actualPkg)
					}
				})
			}
		})

		t.Run("reload test package", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach()

			fs := token.NewFileSet()
			fs.AddFile("file1", 1, 0)
			nonTestPkgs := []*packages.Package{
				{
					Syntax: []*goAst.File{
						{
							Package: 1,
							Decls: []goAst.Decl{
								&goAst.GenDecl{
									Specs: []goAst.Spec{
										&goAst.TypeSpec{
											Name: goAst.NewIdent("type1"),
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
					GoFiles:   []string{"file1"},
					PkgPath:   "the-pkg",
				},
			}

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "the-pkg").returnResults(nonTestPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			_, _, _ = cache.Type(*ast.IdPath("type1", "the-pkg"), false)
			loadCfg.Tests = true
			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "the-pkg").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()

			id := ast.IdPath("type2", "the-pkg")

			// ACT
			actualType, actualPkg, actualErr := cache.Type(*id, true)

			// ASSERT
			if actualErr != nil {
				t.Fatalf("got %#v, want no error", actualErr)
			}

			if actualType.Name.Name != "type2" {
				t.Errorf("got %#v, want type2", actualType.Name.Name)
			}

			if actualPkg != "the-pkg" {
				t.Errorf("got %s, want 'the-pkg'", actualPkg)
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
			file := path.Join(dir, "code.go")
			err = os.WriteFile(file, []byte(code), 0o600)
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			pkgs, err := packages.Load(&packages.Config{
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
			pkgs[0].PkgPath = pkgPath

			return pkgs
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
							beforeEach(t, false)
							defer afterEach()

							metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.Optional())
							metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults().Repeat(moq.Optional())
							loadFnMoq.onCall(loadCfg, "io").
								returnResults(ioTypes, nil).repeat(moq.AnyTimes())
							metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults().Repeat(moq.Optional())

							expr := simpleExpr(t, paramType)

							// ACT
							isComparable, err := stc.compFn(cache, expr)
							// ASSERT
							if err != nil {
								t.Errorf("got %#v, want no error", err)
							}
							if isComparable != stc.comparable {
								t.Errorf("got %t, want %t", isComparable, stc.comparable)
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

							beforeEach(t, false)
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

					beforeEach(t, false)
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
		beforeEach(t, false)
		defer afterEach()

		cache := ast.NewCache(packages.Load, metricsMoq.Mock())

		metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.AnyTimes())
		metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults().Repeat(moq.AnyTimes())
		metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults().Repeat(moq.AnyTimes())

		typ, _, err := cache.Type(
			*ast.IdPath("TypeCache", "github.com/myshkin5/moqueries/generator"), true)
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
		isComparable, err := cache.IsComparable(expr)
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
			beforeEach(t, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "./the-pkg").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()

			// ACT
			pkgPath, err := cache.FindPackage("the-pkg")
			// ASSERT
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			if pkgPath != "the-pkg" {
				t.Errorf("got %s, want the-pkg", pkgPath)
			}
		})

		t.Run("abs dir", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "/this-dir").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()

			// ACT
			pkgPath, err := cache.FindPackage("/this-dir")
			// ASSERT
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			if pkgPath != "the-pkg" {
				t.Errorf("got %s, want the-pkg", pkgPath)
			}
		})

		t.Run("current dir", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, ".").returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()

			// ACT
			pkgPath, err := cache.FindPackage(".")
			// ASSERT
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			if pkgPath != "the-pkg" {
				t.Errorf("got %s, want the-pkg", pkgPath)
			}
		})
	})
}
