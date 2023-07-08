package ast_test

import (
	"errors"
	"fmt"
	goAst "go/ast"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"golang.org/x/tools/go/packages"
	"moqueries.org/runtime/moq"

	"moqueries.org/cli/ast"
	"moqueries.org/cli/metrics"
)

const (
	testPkgs          = "moqueries.org/cli/ast/testpkgs/"
	noExportPkg       = testPkgs + "noexport"
	noExportTestPkg   = noExportPkg + "_test"
	exportPkg         = testPkgs + "export"
	replacebuiltinPkg = testPkgs + "replacebuiltin"
)

var (
	builtinPkgs        []*packages.Package
	ioPkgs             []*packages.Package
	noExportPkgs       []*packages.Package
	exportPkgs         []*packages.Package
	replacebuiltinPkgs []*packages.Package

	noExportWTestsPkgs []*packages.Package
)

func TestMain(m *testing.M) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes,
		Tests: false,
	}
	var err error
	builtinPkgs, err = packages.Load(cfg, "builtin")
	if err != nil {
		panic(fmt.Sprintf("Could not load io package: %#v", err))
	}
	ioPkgs, err = packages.Load(cfg, "io")
	if err != nil {
		panic(fmt.Sprintf("Could not load io package: %#v", err))
	}
	noExportPkgs, err = packages.Load(cfg, noExportPkg)
	if err != nil {
		panic(fmt.Sprintf("Could not load noexport package: %#v", err))
	}
	exportPkgs, err = packages.Load(cfg, exportPkg)
	if err != nil {
		panic(fmt.Sprintf("Could not load export package: %#v", err))
	}
	replacebuiltinPkgs, err = packages.Load(cfg, replacebuiltinPkg)
	if err != nil {
		panic(fmt.Sprintf("Could not load export package: %#v", err))
	}

	cfg.Tests = true
	noExportWTestsPkgs, err = packages.Load(cfg, noExportPkg)
	if err != nil {
		panic(fmt.Sprintf("Could not load noexport package with test types: %#v", err))
	}

	os.Exit(m.Run())
}

func TestCache(t *testing.T) {
	var (
		scene         *moq.Scene
		loadFnMoq     *moqLoadFn
		statFnMoq     *moqStatFn
		readFileFnMoq *moqReadFileFn
		metricsMoq    *metrics.MoqMetrics

		cache *ast.Cache

		loadCfg *packages.Config
	)

	beforeEach := func(t *testing.T, testImport bool) {
		t.Helper()
		if scene != nil {
			t.Fatal("afterEach not called")
		}
		scene = moq.NewScene(t)
		loadFnMoq = newMoqLoadFn(scene, nil)
		statFnMoq = newMoqStatFn(scene, nil)
		readFileFnMoq = newMoqReadFileFn(scene, nil)
		metricsMoq = metrics.NewMoqMetrics(scene, nil)

		cache = ast.NewCache(loadFnMoq.mock(), statFnMoq.mock(), readFileFnMoq.mock(), metricsMoq.Mock())

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
	}
	afterEach := func(t *testing.T) {
		t.Helper()
		scene.AssertExpectationsMet()
		scene = nil
	}

	t.Run("Type", func(t *testing.T) {
		t.Run("simple load", func(t *testing.T) {
			validateFuncFn := func(t *testing.T, fnType *dst.FuncType) {
				t.Helper()
				if fnType.Params == nil {
					t.Fatalf("got nil, want not nil")
				}
				if len(fnType.Params.List) != 0 {
					t.Errorf("got %#v, want empty param list",
						fnType.Params.List)
				}
			}
			validateInterfaceFn := func(t *testing.T, iType *dst.InterfaceType, name1, name2 string) {
				t.Helper()
				if iType.Methods == nil {
					t.Fatalf("got nil, want not nil")
				}

				if len(iType.Methods.List) != 2 {
					t.Fatalf("got %#v, want length 2", iType.Methods.List)
				}

				if iType.Methods.List[0].Names[0].String() != name1 {
					t.Errorf("got %s, want %s",
						iType.Methods.List[0].Names[0].String(), name1)
				}

				if iType.Methods.List[1].Names[0].String() != name2 {
					t.Errorf("got %s, want %s",
						iType.Methods.List[1].Names[0].String(), name2)
				}
			}
			testCases := map[string]struct {
				pkgs               []*packages.Package
				typeToLoad         string
				testImport         bool
				expectedInterface  bool
				expectedPkg        string
				expectedExported   bool
				expectedFabricated bool
				validateInterface  func(t *testing.T, iType *dst.InterfaceType)
				validateFunc       func(t *testing.T, fnType *dst.FuncType)
			}{
				"regular non-test": {
					pkgs:               noExportPkgs,
					typeToLoad:         "type1",
					testImport:         false,
					expectedInterface:  true,
					expectedPkg:        noExportPkg,
					expectedExported:   false,
					expectedFabricated: false,
				},
				"exported regular non-test": {
					pkgs:               exportPkgs,
					typeToLoad:         "Type1",
					testImport:         false,
					expectedInterface:  true,
					expectedPkg:        exportPkg,
					expectedExported:   true,
					expectedFabricated: false,
				},
				"regular test": {
					pkgs:               noExportWTestsPkgs,
					typeToLoad:         "test_type1",
					testImport:         true,
					expectedInterface:  true,
					expectedPkg:        noExportTestPkg,
					expectedExported:   false,
					expectedFabricated: false,
				},
				"func non-test": {
					pkgs:               noExportPkgs,
					typeToLoad:         "type4_genType",
					testImport:         false,
					expectedInterface:  false,
					expectedPkg:        noExportPkg,
					expectedExported:   false,
					expectedFabricated: true,
					validateFunc:       validateFuncFn,
				},
				"exported func non-test": {
					pkgs:               exportPkgs,
					typeToLoad:         "Type4_genType",
					testImport:         false,
					expectedInterface:  false,
					expectedPkg:        exportPkg,
					expectedExported:   true,
					expectedFabricated: true,
					validateFunc:       validateFuncFn,
				},
				"func test": {
					pkgs:               noExportWTestsPkgs,
					typeToLoad:         "test_type4_genType",
					testImport:         true,
					expectedInterface:  false,
					expectedPkg:        noExportTestPkg,
					expectedExported:   false,
					expectedFabricated: true,
					validateFunc:       validateFuncFn,
				},
				"method non-test": {
					pkgs:               noExportPkgs,
					typeToLoad:         "widget_genType",
					testImport:         false,
					expectedInterface:  true,
					expectedPkg:        noExportPkg,
					expectedExported:   false,
					expectedFabricated: true,
					validateInterface: func(t *testing.T, iType *dst.InterfaceType) {
						t.Helper()
						validateInterfaceFn(t, iType, "method1", "method2")
					},
				},
				"exported method non-test": {
					pkgs:               exportPkgs,
					typeToLoad:         "Widget_genType",
					testImport:         false,
					expectedInterface:  true,
					expectedPkg:        exportPkg,
					expectedExported:   true,
					expectedFabricated: true,
					validateInterface: func(t *testing.T, iType *dst.InterfaceType) {
						t.Helper()
						validateInterfaceFn(t, iType, "Type5", "Type6")
					},
				},
				"method test": {
					pkgs:               noExportWTestsPkgs,
					typeToLoad:         "test_widget_genType",
					testImport:         true,
					expectedInterface:  true,
					expectedPkg:        noExportTestPkg,
					expectedExported:   false,
					expectedFabricated: true,
					validateInterface: func(t *testing.T, iType *dst.InterfaceType) {
						t.Helper()
						validateInterfaceFn(t, iType, "test_method1", "test_method2")
					},
				},
				"star method non-test": {
					pkgs:               noExportPkgs,
					typeToLoad:         "widget_starGenType",
					testImport:         false,
					expectedInterface:  true,
					expectedPkg:        noExportPkg,
					expectedExported:   false,
					expectedFabricated: true,
					validateInterface: func(t *testing.T, iType *dst.InterfaceType) {
						t.Helper()
						validateInterfaceFn(t, iType, "method3", "method4")
					},
				},
				"exported star method non-test": {
					pkgs:               exportPkgs,
					typeToLoad:         "Widget_starGenType",
					testImport:         false,
					expectedInterface:  true,
					expectedPkg:        exportPkg,
					expectedExported:   true,
					expectedFabricated: true,
					validateInterface: func(t *testing.T, iType *dst.InterfaceType) {
						t.Helper()
						validateInterfaceFn(t, iType, "Type7", "Type8")
					},
				},
				"star method test": {
					pkgs:               noExportWTestsPkgs,
					typeToLoad:         "test_widget_starGenType",
					testImport:         true,
					expectedInterface:  true,
					expectedPkg:        noExportTestPkg,
					expectedExported:   false,
					expectedFabricated: true,
					validateInterface: func(t *testing.T, iType *dst.InterfaceType) {
						t.Helper()
						validateInterfaceFn(t, iType, "test_method3", "test_method4")
					},
				},
			}

			for name, tc := range testCases {
				t.Run(name, func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t, tc.testImport)
					defer afterEach(t)

					metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
					loadFnMoq.onCall(loadCfg, noExportPkg).returnResults(tc.pkgs, nil)
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
					metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
						ReturnResults().Repeat(moq.MinTimes(1))

					id := ast.IdPath(tc.typeToLoad, noExportPkg)

					// ACT
					actualType, actualErr := cache.Type(*id, noExportPkg, tc.testImport)

					// ASSERT
					if actualErr != nil {
						t.Fatalf("got %#v, want no error", actualErr)
					}

					if actualType.Type.Name.Name != tc.typeToLoad {
						t.Errorf("got %#v, want %s", actualType.Type.Name.Name, tc.typeToLoad)
					}

					if actualType.PkgPath != tc.expectedPkg {
						t.Errorf("got %s, want %s", actualType.PkgPath, tc.expectedPkg)
					}

					if actualType.Exported != tc.expectedExported {
						t.Errorf("got %t, want %t", actualType.Exported, tc.expectedExported)
					}

					if actualType.Fabricated != tc.expectedFabricated {
						t.Errorf("got %t, want %t", actualType.Fabricated, tc.expectedFabricated)
					}

					if tc.expectedInterface {
						iType, ok := actualType.Type.Type.(*dst.InterfaceType)
						if !ok {
							t.Fatalf("got %#v, want *dst.InterfaceType", actualType.Type.Type)
						}
						if tc.validateInterface != nil {
							tc.validateInterface(t, iType)
						}
					} else {
						fnType, ok := actualType.Type.Type.(*dst.FuncType)
						if !ok {
							t.Fatalf("got %#v, want *dst.FuncType", actualType.Type.Type)
						}
						if tc.validateFunc != nil {
							tc.validateFunc(t, fnType)
						}
					}
				})
			}
		})

		t.Run("can load builtins", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, exportPkg).returnResults(exportPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "builtin").returnResults(builtinPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			id := ast.Id("string")

			// ACT
			typ, err := cache.Type(*id, exportPkg, false)
			// ASSERT
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			if typ.Type.Name.Name != "string" {
				t.Errorf("got %s, want string", typ.Type.Name.Name)
			}
			if typ.PkgPath != "" {
				t.Errorf("got %s, want empty package", typ.PkgPath)
			}
			if !typ.Exported {
				t.Errorf("got not exported, want exported")
			}
			if typ.Fabricated {
				t.Errorf("got fabricated, want not fabricated")
			}
		})

		t.Run("can load replaced builtins", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, replacebuiltinPkg).returnResults(replacebuiltinPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			id := ast.Id("string")

			// ACT
			typ, err := cache.Type(*id, replacebuiltinPkg, false)
			// ASSERT
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			if typ.Type.Name.Name != "string" {
				t.Errorf("got %s, want string", typ.Type.Name.Name)
			}
			if typ.PkgPath != replacebuiltinPkg {
				t.Errorf("got %s, want %s", typ.PkgPath, replacebuiltinPkg)
			}
			if typ.Exported {
				t.Errorf("got exported, want not exported")
			}
			if typ.Fabricated {
				t.Errorf("got fabricated, want not fabricated")
			}
		})

		t.Run("loads test package when given a test package", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, true)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, noExportPkg).returnResults(noExportWTestsPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
				ReturnResults().Repeat(moq.MinTimes(1))

			id := ast.IdPath("test_type1", noExportTestPkg)

			// ACT
			actualType, actualErr := cache.Type(*id, noExportTestPkg, false)

			// ASSERT
			if actualErr != nil {
				t.Fatalf("got %#v, want no error", actualErr)
			}

			if actualType.Type.Name.Name != "test_type1" {
				t.Errorf("got %#v, want test_type1", actualType.Type.Name.Name)
			}

			if actualType.PkgPath != noExportTestPkg {
				t.Errorf("got %s, want %s", actualType.PkgPath, noExportTestPkg)
			}

			if actualType.Fabricated {
				t.Errorf("got %t, want false", actualType.Fabricated)
			}
		})

		t.Run("load error", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, true)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			err := errors.New("load error")
			loadFnMoq.onCall(loadCfg, ".").returnResults(nil, err)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()

			id := ast.IdPath("type1", ".")

			// ACT
			actualType, actualErr := cache.Type(*id, ".", true)

			// ASSERT
			if actualErr != err {
				t.Errorf("got %#v, want %#v", actualErr, err)
			}

			if actualType.Type != nil {
				t.Errorf("got %#v, want nil type", actualType.Type)
			}

			if actualType.PkgPath != "" {
				t.Errorf("got %s, want empty pkg", actualType.PkgPath)
			}
		})

		t.Run("not found", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, noExportPkg).returnResults(noExportPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			id := ast.IdPath("notthere1", noExportPkg)

			// ACT
			actualType, actualErr := cache.Type(*id, noExportPkg, false)

			// ASSERT
			if actualErr == nil || !strings.Contains(actualErr.Error(), "not found") {
				t.Errorf("got %#v, want to contain 'not found'", actualErr)
			}
			if !errors.Is(actualErr, ast.ErrTypeNotFound) {
				t.Errorf("got %#v, want ast.ErrTypeNotFound", actualErr)
			}

			if actualType.Type != nil {
				t.Errorf("got %#v, want nil type", actualType.Type)
			}

			if actualType.PkgPath != "" {
				t.Errorf("got %s, want empty pkg", actualType.PkgPath)
			}
		})

		t.Run("load packages only once", func(t *testing.T) {
			testCases := map[string]struct {
				typeToLoad        string
				pkgName           string
				initialTestImport bool
				actTestImport     bool
				expectedPkg       string
			}{
				"test load then no test load": {
					typeToLoad:        "type2",
					pkgName:           noExportPkg,
					initialTestImport: true,
					actTestImport:     false,
					expectedPkg:       noExportPkg,
				},
				"relative package": {
					typeToLoad:        "type2",
					pkgName:           "./testpkgs/noexport",
					initialTestImport: true,
					actTestImport:     false,
					expectedPkg:       noExportPkg,
				},
				"no test load, twice": {
					typeToLoad:        "type2",
					pkgName:           noExportPkg,
					initialTestImport: false,
					actTestImport:     false,
					expectedPkg:       noExportPkg,
				},
				"test load, twice": {
					typeToLoad:        "test_type2",
					pkgName:           noExportPkg,
					initialTestImport: true,
					actTestImport:     true,
					expectedPkg:       noExportTestPkg,
				},
			}

			for name, tc := range testCases {
				t.Run(name, func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t, true)
					defer afterEach(t)

					metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
					loadCfg.Tests = tc.initialTestImport
					loadFnMoq.onCall(loadCfg, tc.pkgName).returnResults(noExportWTestsPkgs, nil)
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
					metricsMoq.OnCall().ASTPkgCacheHitsInc().ReturnResults()
					metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
						ReturnResults().Repeat(moq.MinTimes(1))
					mehId := ast.IdPath("doesnotmatter", tc.pkgName)
					_, _ = cache.Type(*mehId, tc.pkgName, tc.initialTestImport)

					id := ast.IdPath(tc.typeToLoad, noExportPkg)

					// ACT
					actualType, actualErr := cache.Type(*id, noExportPkg, tc.actTestImport)

					// ASSERT
					if actualErr != nil {
						t.Fatalf("got %#v, want no error", actualErr)
					}

					if actualType.Type.Name.Name != tc.typeToLoad {
						t.Errorf("got %#v, want %s", actualType.Type.Name.Name, tc.typeToLoad)
					}

					if actualType.PkgPath != tc.expectedPkg {
						t.Errorf("got %s, want %s", actualType.PkgPath, tc.expectedPkg)
					}
				})
			}
		})

		t.Run("reload test package", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, noExportPkg).returnResults(noExportPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
				ReturnResults().Repeat(moq.MinTimes(1))
			type1Id := ast.IdPath("type1", noExportPkg)
			_, _ = cache.Type(*type1Id, noExportPkg, false)

			loadCfg.Tests = true

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, noExportPkg).returnResults(noExportWTestsPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()

			id := ast.IdPath("test_type2", noExportPkg)

			// ACT
			actualType, actualErr := cache.Type(*id, noExportPkg, true)

			// ASSERT
			if actualErr != nil {
				t.Fatalf("got %#v, want no error", actualErr)
			}

			if actualType.Type.Name.Name != "test_type2" {
				t.Errorf("got %#v, want test_type2", actualType.Type.Name.Name)
			}

			if actualType.PkgPath != noExportTestPkg {
				t.Errorf("got %s, want %s", actualType.PkgPath, noExportTestPkg)
			}
		})

		t.Run("knows how to load builtin error type", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().
				ReturnResults().Repeat(moq.Times(2))
			loadFnMoq.onCall(loadCfg, noExportPkg).returnResults(noExportPkgs, nil)
			loadFnMoq.onCall(loadCfg, "builtin").returnResults(builtinPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().
				ReturnResults().Repeat(moq.Times(2))
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
				ReturnResults().Repeat(moq.Times(2))

			id := ast.IdPath("error", "")

			// ACT
			actualType, actualErr := cache.Type(*id, noExportPkg, false)

			// ASSERT
			if actualErr != nil {
				t.Fatalf("got %#v, want no error", actualErr)
			}

			if actualType.Type.Name.Name != "error" {
				t.Errorf("got %#v, want error", actualType.Type.Name.Name)
			}

			if actualType.PkgPath != "" {
				t.Errorf("got %s, want empty package name", actualType.PkgPath)
			}
		})

		t.Run("rel/abs only once", func(t *testing.T) {
			relNoExportPkg := "./testpkgs/noexport"
			testCases := map[string]struct {
				firstPath  string
				secondPath string
			}{
				"relative directory then full package": {
					firstPath:  relNoExportPkg,
					secondPath: noExportPkg,
				},
				"full package then relative directory": {
					firstPath:  noExportPkg,
					secondPath: relNoExportPkg,
				},
			}

			for name, tc := range testCases {
				t.Run(name, func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t, false)
					defer afterEach(t)

					metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
					loadFnMoq.onCall(loadCfg, tc.firstPath).returnResults(noExportPkgs, nil)
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
					metricsMoq.OnCall().ASTPkgCacheHitsInc().ReturnResults()
					metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
						ReturnResults().Repeat(moq.MinTimes(1))

					id := ast.IdPath("type2", tc.firstPath)
					actualRelType, actualErr := cache.Type(*id, tc.firstPath, false)
					if actualErr != nil {
						t.Fatalf("got %#v, want no error", actualErr)
					}
					id = ast.IdPath("type2", tc.secondPath)

					// ACT
					actualPkgType, actualErr := cache.Type(*id, tc.secondPath, false)

					// ASSERT
					if actualErr != nil {
						t.Fatalf("got %#v, want no error", actualErr)
					}

					if actualPkgType.Type.Name.Name != "type2" {
						t.Errorf("got %#v, want %s", actualPkgType.Type.Name.Name, "type2")
					}

					if actualPkgType.PkgPath != noExportPkg {
						t.Errorf("got %s, want %s", actualPkgType.PkgPath, noExportPkg)
					}

					if actualPkgType != actualRelType {
						t.Errorf("got %#v, want %#v", actualPkgType, actualRelType)
					}
				})
			}
		})

		t.Run("loads a directory then the same package only once", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, ".").returnResults(noExportWTestsPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTPkgCacheHitsInc().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
				ReturnResults().Repeat(moq.MinTimes(1))

			id := ast.IdPath("type2", ".")
			actualRelType, actualErr := cache.Type(*id, ".", false)
			if actualErr != nil {
				t.Fatalf("got %#v, want no error", actualErr)
			}
			id = ast.IdPath("type2", noExportPkg)

			// ACT
			actualPkgType, actualErr := cache.Type(*id, noExportPkg, false)

			// ASSERT
			if actualErr != nil {
				t.Fatalf("got %#v, want no error", actualErr)
			}

			if actualPkgType.Type.Name.Name != "type2" {
				t.Errorf("got %#v, want %s", actualPkgType.Type.Name.Name, "type2")
			}

			if actualPkgType.PkgPath != noExportPkg {
				t.Errorf("got %s, want %s", actualPkgType.PkgPath, noExportPkg)
			}

			if actualPkgType != actualRelType {
				t.Errorf("got %#v, want %#v", actualPkgType, actualRelType)
			}
		})
	})

	t.Run("IsComparable/IsDefaultComparable", func(t *testing.T) {
		testCases := map[string]struct {
			comparable        bool
			defaultComparable bool
			structable        bool
		}{
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

		loadPackages := func(t *testing.T, code map[string]string) []*packages.Package {
			t.Helper()
			dir, err := os.MkdirTemp("", "cache-test-*")
			if err != nil {
				t.Fatalf("got os.MkdirTemp err: %#v, want no error", err)
			}
			defer func() {
				err := os.RemoveAll(dir)
				if err != nil {
					t.Errorf("got os.RemoveAll err: %#v, want no error", err)
				}
			}()

			for srcPath, src := range code {
				fPath := filepath.Join(dir, srcPath)
				fDir, _ := filepath.Split(fPath)
				if err := os.Mkdir(fDir, fs.ModePerm); err != nil && !errors.Is(err, fs.ErrExist) {
					t.Fatalf("got os.Mkdir err: %#v, want no error", err)
				}
				err := os.WriteFile(fPath, []byte(src), fs.ModePerm)
				if err != nil {
					t.Fatalf("got os.WriteFile err: %#v, want no error", err)
				}
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
				Dir:   dir,
				Tests: false,
			}, "a")
			if err != nil {
				t.Fatalf("got packages.Load err: %#v, want no error", err)
			}
			return pkgs
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

		isComparable := func(c *ast.Cache, expr dst.Expr, parentType ast.TypeInfo) (bool, error) {
			return c.IsComparable(expr, parentType)
		}
		isDefaultComparable := func(c *ast.Cache, expr dst.Expr, parentType ast.TypeInfo) (bool, error) {
			return c.IsDefaultComparable(expr, parentType)
		}
		type compFn func(c *ast.Cache, expr dst.Expr, parentType ast.TypeInfo) (bool, error)

		t.Run("simple exprs", func(t *testing.T) {
			for paramType, tc := range testCases {
				t.Run(paramType, func(t *testing.T) {
					for name, stc := range map[string]struct {
						compFn     compFn
						comparable bool
					}{
						"IsComparable": {
							compFn:     isComparable,
							comparable: tc.comparable,
						},
						"IsDefaultComparable": {
							compFn:     isDefaultComparable,
							comparable: tc.defaultComparable,
						},
					} {
						t.Run(name, func(t *testing.T) {
							// ASSEMBLE
							beforeEach(t, false)
							defer afterEach(t)

							metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.Optional())
							loadFnMoq.onCall(loadCfg, "io").
								returnResults(ioPkgs, nil).repeat(moq.AnyTimes())
							metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults().Repeat(moq.Optional())
							metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
								ReturnResults().Repeat(moq.AnyTimes())

							expr := simpleExpr(t, paramType)

							// ACT
							is, err := stc.compFn(cache, expr, ast.TypeInfo{})
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

		t.Run("IsComparable - struct exprs", func(t *testing.T) {
			code := `package a

import _ "io"

type b struct {
	c %s
}

func d(e b) {}
`

			for paramType, tc := range testCases {
				t.Run(paramType, func(t *testing.T) {
					// ASSEMBLE
					if !tc.structable {
						t.Skipf("%s can't be put into a struct, skipping", paramType)
					}

					beforeEach(t, false)
					defer afterEach(t)

					pkgs := loadPackages(t, map[string]string{
						"code.go": fmt.Sprintf(code, paramType),
						"go.mod":  "module a",
					})

					metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().
						Repeat(moq.MinTimes(1), moq.MaxTimes(2))
					loadFnMoq.onCall(loadCfg, "a").
						returnResults(pkgs, nil)
					loadFnMoq.onCall(loadCfg, "io").
						returnResults(ioPkgs, nil).repeat(moq.Optional())
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults().
						Repeat(moq.MinTimes(1), moq.MaxTimes(2))
					metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
						ReturnResults().Repeat(moq.AnyTimes())

					// ACT
					isComparable, err := cache.IsComparable(ast.IdPath("b", "a"), ast.TypeInfo{})
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

		t.Run("IsComparable - imported", func(t *testing.T) {
			for paramType, tc := range testCases {
				t.Run(paramType, func(t *testing.T) {
					// ASSEMBLE
					if !tc.structable {
						t.Skipf("%s can't be put into a struct, skipping", paramType)
					}

					beforeEach(t, false)
					defer afterEach(t)

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
					loadFnMoq.onCall(loadCfg, "b").returnResults(bPkg, nil)
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults().Repeat(moq.Optional())

					metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.Optional())
					loadFnMoq.onCall(loadCfg, "io").
						returnResults(ioPkgs, nil).repeat(moq.AnyTimes())
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults().Repeat(moq.Optional())
					metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
						ReturnResults().Repeat(moq.AnyTimes())

					// ACT
					isComparable, err := cache.IsComparable(expr, ast.TypeInfo{})
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

		t.Run("generics", func(t *testing.T) {
			for name, tc := range map[string]struct {
				comparable        bool
				defaultComparable bool
				typeConstraints   string
				structContents    string
				typeParams        string
				errorContains     string
				alterMethod       func(*dst.FuncDecl)
				skipNonMethodTest bool
			}{
				"[T any, U comparable, V any]": {
					comparable:        true,
					defaultComparable: true,
					typeConstraints:   "[T any, U comparable, V any]",
					structContents:    "t T; u U; v V",
					typeParams:        "T, U, V",
				},
				"[T any, U any, V any]": {
					comparable:        false,
					defaultComparable: false,
					typeConstraints:   "[T any, U any, V any]",
					structContents:    "t T; u U; v V",
					typeParams:        "T, U, V",
				},
				"[U comparable]": {
					comparable:        true,
					defaultComparable: true,
					typeConstraints:   "[U comparable]",
				},
				"[U any]": {
					comparable:        false,
					defaultComparable: false,
					typeConstraints:   "[U any]",
				},
				"[U interface{}]": {
					comparable:        false,
					defaultComparable: false,
					typeConstraints:   "[U interface{}]",
				},
				"[U []int]": {
					comparable:        false,
					defaultComparable: false,
					typeConstraints:   "[U []int]",
				},
				"[U int]": {
					comparable:        true,
					defaultComparable: true,
					typeConstraints:   "[U int]",
				},
				"[U interface{ []int }]": {
					comparable:        false,
					defaultComparable: false,
					typeConstraints:   "[U interface{ []int }]",
				},
				"[U interface{ int }]": {
					comparable:        true,
					defaultComparable: true,
					typeConstraints:   "[U interface{ int }]",
				},
				"[U notComparable]": {
					comparable:        false,
					defaultComparable: false,
					typeConstraints:   "[U notComparable]",
				},
				"[U isComparable]": {
					comparable:        true,
					defaultComparable: true,
					typeConstraints:   "[U isComparable]",
				},
				"[U ~notComparable]": {
					comparable:        false,
					defaultComparable: false,
					typeConstraints:   "[U ~notComparable]",
				},
				"[U ~isComparable]": {
					comparable:        true,
					defaultComparable: true,
					typeConstraints:   "[U ~isComparable]",
				},
				"[U interface{ ~int | ~string }]": {
					comparable:        true,
					defaultComparable: true,
					typeConstraints:   "[U interface{ ~int | ~string }]",
				},
				"[U interface{ ~int | ~string | ~notComparable }]": {
					comparable:        false,
					defaultComparable: false,
					typeConstraints:   "[U interface{ ~int | ~string | ~notComparable }]",
				},
				"[U interface{ ~notComparable | ~int | ~string }]": {
					comparable:        false,
					defaultComparable: false,
					typeConstraints:   "[U interface{ ~notComparable | ~int | ~string }]",
				},
				"[U interface{ ~int | ~isComparable | ~string }]": {
					comparable:        true,
					defaultComparable: true,
					typeConstraints:   "[U interface{ ~int | ~isComparable | ~string }]",
				},
				"[U interface{ comparable; m() }]": {
					comparable:        true,
					defaultComparable: true,
					typeConstraints:   "[U interface{ comparable; m() }]",
				},
				"no type constraints on struct": {
					typeConstraints: "",
					errorContains:   "base type to method type param mismatch",
					// A generic function (not a method) doesn't depend on the
					// type params on a struct. Skip
					skipNonMethodTest: true,
				},
				"non-id IndexListExpr.X": {
					typeConstraints: "[T any, U any, V any]",
					structContents:  "t T; u U; v V",
					typeParams:      "T, U, V",
					alterMethod: func(decl *dst.FuncDecl) {
						decl.Recv.List[0].Type.(*dst.IndexListExpr).X = ast.Un(token.AND, ast.Id("hi"))
					},
					// There is no index list for a generic function
					skipNonMethodTest: true,
					errorContains:     "expecting *dst.Ident in IndexListExpr.X",
				},
				"non-id IndexExpr.X": {
					alterMethod: func(decl *dst.FuncDecl) {
						decl.Recv.List[0].Type.(*dst.IndexExpr).X = ast.Un(token.AND, ast.Id("hi"))
					},
					// There is no index list for a generic function
					skipNonMethodTest: true,
					errorContains:     "expecting *dst.Ident in IndexExpr.X",
				},
				"unexpected index type": {
					alterMethod: func(decl *dst.FuncDecl) {
						decl.Recv.List[0].Type = ast.Un(token.AND, ast.Id("hi"))
					},
					// There is no index list for a generic function
					skipNonMethodTest: true,
					errorContains:     "unexpected index type",
				},
				"too many method type params": {
					typeParams:    "A, B, C, U",
					errorContains: "base type to method type param mismatch",
					// A generic function (not a method) doesn't depend on
					// the type params on a struct. Skip
					skipNonMethodTest: true,
				},
				"no type spec": {
					alterMethod: func(decl *dst.FuncDecl) {
						decl.Recv.List[0].Type.(*dst.IndexExpr).X.(*dst.Ident).Obj = nil
					},
					// There is no index list for a generic function
					skipNonMethodTest: true,
					errorContains:     "expecting Obj",
				},
				"non-type spec decl": {
					alterMethod: func(decl *dst.FuncDecl) {
						decl.Recv.List[0].Type.(*dst.IndexExpr).X.(*dst.Ident).Obj.Decl = "hi"
					},
					// There is no index list for a generic function
					skipNonMethodTest: true,
					errorContains:     "expecting *dst.TypeSpec",
				},
				"bad underlying type operator": {
					typeConstraints: "[U ~int]",
					alterMethod: func(decl *dst.FuncDecl) {
						decl.Recv.List[0].
							Type.(*dst.IndexExpr).
							X.(*dst.Ident).Obj.
							Decl.(*dst.TypeSpec).TypeParams.List[0].
							Type.(*dst.UnaryExpr).Op = token.AND
					},
					// There is no index list for a generic function
					skipNonMethodTest: true,
					errorContains:     "unexpected unary operator &",
				},
			} {
				t.Run(name, func(t *testing.T) {
					for name, stc := range map[string]struct {
						compFn     compFn
						comparable bool
						structable bool
					}{
						"IsComparable": {
							compFn:     isComparable,
							comparable: tc.comparable,
							structable: false,
						},
						"IsDefaultComparable": {
							compFn:     isDefaultComparable,
							comparable: tc.defaultComparable,
							structable: true,
						},
					} {
						t.Run(name, func(t *testing.T) {
							for name, fn := range map[string]func(*testing.T) ast.TypeInfo{
								// "struct context": func(t *testing.T, f *dst.File) (dst.Expr, *dst.TypeSpec) {
								// 	if tc.skipNonMethodTest {
								// 		t.Skip()
								// 	}
								//
								// 	gen, ok := f.Decls[2].(*dst.GenDecl)
								// 	if !ok {
								// 		t.Fatalf("got %#v, want a generic declaration", f.Decls[1])
								// 	}
								// 	fields := gen.Specs[0].(*dst.TypeSpec).Type.(*dst.StructType).Fields.List
								// 	var idx int
								// 	// A bit brittle but at least one of the
								// 	// tests puts U in the middle of the list
								// 	if len(fields) == 1 {
								// 		idx = 0
								// 	} else {
								// 		idx = 1
								// 	}
								// 	expr := fields[idx].Type
								// 	return expr, nil
								// },
								// "method context": func(t *testing.T, f *dst.File) (dst.Expr, *dst.FuncDecl) {
								// 	fn, ok := f.Decls[3].(*dst.FuncDecl)
								// 	if !ok {
								// 		t.Fatalf("got %#v, want a function declaration", f.Decls[1])
								// 	}
								// 	expr := fn.Type.Params.List[0].Type
								// 	if tc.alterMethod != nil {
								// 		tc.alterMethod(fn)
								// 	}
								// 	return expr, fn
								// },
								// "function context": func(t *testing.T, f *dst.File) (dst.Expr, *dst.FuncDecl) {
								// 	if tc.skipNonMethodTest {
								// 		t.Skip()
								// 	}
								//
								// 	fn, ok := f.Decls[4].(*dst.FuncDecl)
								// 	if !ok {
								// 		t.Fatalf("got %#v, want a function declaration", f.Decls[1])
								// 	}
								// 	expr := fn.Type.Params.List[0].Type
								//
								// 	return expr, fn
								// },
								"type context": func(t *testing.T) ast.TypeInfo {
									t.Helper()
									tSpec, err := cache.Type(*ast.IdPath("b", "a"), "", false)
									if err != nil {
										t.Fatalf("got Type error: %#v, want no error", err)
									}

									return tSpec
								},
							} {
								t.Run(name, func(t *testing.T) {
									if tc.skipNonMethodTest {
										t.Skip()
									}

									// ASSEMBLE
									beforeEach(t, false)
									defer afterEach(t)

									code := `package a

type notComparable []int

type isComparable int

type b%s struct{%s}

func (b[%s]) c(U) {}

func d%s(U) {}
`
									structContents := tc.structContents
									if structContents == "" {
										structContents = "u U"
									}
									typeParams := tc.typeParams
									if typeParams == "" {
										typeParams = "U"
									}
									code = fmt.Sprintf(code, tc.typeConstraints,
										structContents, typeParams,
										tc.typeConstraints)
									pkgs := loadPackages(t, map[string]string{
										"code.go": code,
										"go.mod":  "module a",
									})

									metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
									loadFnMoq.onCall(loadCfg, "a").
										returnResults(pkgs, nil)
									metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
									metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

									tSpec := fn(t)

									// ACT
									isComp, err := stc.compFn(cache, ast.Id("U"), tSpec)
									// ASSERT
									if tc.errorContains == "" {
										if err != nil {
											t.Errorf("got %#v, want no error", err)
										}
										if isComp != stc.comparable {
											t.Errorf("got %t, want %t", isComp, stc.comparable)
										}
									} else {
										if err == nil {
											t.Fatalf("got no error, want err")
										}
										if !errors.Is(err, ast.ErrInvalidType) {
											t.Errorf("got %#v, want ast.ErrInvalidType", err)
										}
										if !strings.Contains(err.Error(), tc.errorContains) {
											t.Errorf("got %s, want to contain %s", err.Error(), tc.errorContains)
										}
									}
								})
							}
						})
					}
				})
			}
		})
	})

	t.Run("DST ident not comparable", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t, false)
		defer afterEach(t)

		c := ast.NewCache(packages.Load, os.Stat, os.ReadFile, metricsMoq.Mock())

		metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.AnyTimes())
		metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().
			ReturnResults().Repeat(moq.AnyTimes())
		metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
			ReturnResults().Repeat(moq.AnyTimes())

		pkg := "moqueries.org/cli/generator"
		typ, err := c.Type(*ast.IdPath("TypeCache", pkg), pkg, false)
		if err != nil {
			t.Fatalf("got %#v, want no error", err)
		}
		iface, ok := typ.Type.Type.(*dst.InterfaceType)
		if !ok {
			t.Fatalf("got %#v, want an interface", typ.Type)
		}
		fType, ok := iface.Methods.List[0].Type.(*dst.FuncType)
		if !ok {
			t.Fatalf("got %#v, want a function type", iface.Methods.List[0].Type)
		}
		expr := fType.Params.List[0].Type

		// ACT
		isComparable, err := c.IsComparable(expr, ast.TypeInfo{})
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, want no error", err)
		}
		if isComparable {
			t.Errorf("got %t, want false", isComparable)
		}
	})

	t.Run("FindPackage", func(t *testing.T) {
		t.Run("uses cached package if it exists", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "./testpkgs/noexport").returnResults(noExportPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			err := cache.LoadPackage("./testpkgs/noexport")
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			// ACT
			pkgPath, err := cache.FindPackage("testpkgs/noexport")
			// ASSERT
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			if pkgPath != noExportPkg {
				t.Errorf("got %s, want %s", pkgPath, noExportPkg)
			}
		})

		t.Run("reads the go.mod file from a parent directory", func(t *testing.T) {
			goodModFile := "module moqueries.org/cli"
			statError := errors.New("bad stat call")
			readError := errors.New("bad read file call")
			for name, tc := range map[string]struct {
				modFile   string
				statError error
				readError error
				isError   error
				errString string
			}{
				"happy path": {modFile: goodModFile},
				"stat error": {
					statError: statError,
					isError:   statError,
					errString: "bad stat call: error stat-ing %s",
				},
				"read error": {
					readError: readError,
					isError:   readError,
					errString: "bad read file call: error reading %s",
				},
				"bad mod file": {
					modFile:   "module",
					errString: "%s:1: usage: module module/path: error parsing",
				},
				"missing module directive": {
					modFile:   "something not sure what but not a mod file",
					isError:   ast.ErrMissingModuleDirective,
					errString: "missing module directive: error parsing %s",
				},
			} {
				t.Run(name, func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t, false)
					defer afterEach(t)

					dir, err := os.Getwd()
					if err != nil {
						t.Fatalf("got %#v, want no error", err)
					}
					rootDir := filepath.Dir(dir)
					absDir := filepath.Join(dir, "testpkgs/noexport")
					var modPath string
					for {
						statFnMoq.onCall(filepath.Join(absDir, "go.mod")).
							returnResults(nil, fs.ErrNotExist)
						absDir = filepath.Dir(absDir)
						if absDir == rootDir {
							modPath = filepath.Join(absDir, "go.mod")
							statFnMoq.onCall(modPath).returnResults(nil, tc.statError)
							break
						}
					}
					if tc.statError == nil {
						readFileFnMoq.onCall(modPath).
							returnResults([]byte(tc.modFile), tc.readError)
					}

					// ACT
					pkgPath, err := cache.FindPackage("testpkgs/noexport")
					// ASSERT
					if tc.errString == "" {
						if err != nil {
							t.Fatalf("got %#v, want no error", err)
						}

						if pkgPath != noExportPkg {
							t.Errorf("got %s, want %s", pkgPath, noExportPkg)
						}
					} else {
						if err == nil {
							t.Fatalf("got no error, want error")
						}
						if tc.isError != nil && !errors.Is(err, tc.isError) {
							t.Errorf("got %#v, want is %#v", err, tc.isError)
						}
						errString := fmt.Sprintf(tc.errString, modPath)
						if err.Error() != errString {
							t.Errorf("got %s, want %s", err.Error(), errString)
						}

						if pkgPath != "" {
							t.Errorf("got %s, want empty path", pkgPath)
						}
					}
				})
			}
		})
	})

	t.Run("LoadPackage", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, ".").doReturnResults(packages.Load)
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
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
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
		pkg := func(pkgPath string) *packages.Package {
			fs := token.NewFileSet()
			fs.AddFile("file1", 1, 0)
			fs.AddFile("file2", 2, 0)

			return &packages.Package{
				Syntax: []*goAst.File{
					{
						Package: 1,
						Decls: []goAst.Decl{&goAst.GenDecl{
							Specs: []goAst.Spec{&goAst.TypeSpec{
								Name: goAst.NewIdent("Type1"),
								Type: &goAst.InterfaceType{
									Methods: &goAst.FieldList{},
								},
							}},
						}},
					},
					{
						Package: 2,
						Decls: []goAst.Decl{&goAst.GenDecl{
							Specs: []goAst.Spec{&goAst.TypeSpec{
								Name: goAst.NewIdent("Type2"),
								Type: &goAst.FuncType{
									Params: &goAst.FieldList{},
								},
							}},
						}},
					},
					{
						Package: 2,
						Decls: []goAst.Decl{&goAst.GenDecl{
							Specs: []goAst.Spec{&goAst.TypeSpec{
								Name: goAst.NewIdent("Type3"),
								Type: &goAst.StructType{
									Fields: &goAst.FieldList{},
								},
							}},
						}},
					},
				},
				TypesInfo: &types.Info{},
				Fset:      fs,
				GoFiles:   []string{"file1", "file2"},
				PkgPath:   pkgPath,
			}
		}

		t.Run("everything returned", func(t *testing.T) {
			testCases := map[string]struct {
				pkg          string
				pkgs         []*packages.Package
				onlyExported bool
				prefix       string
				widgetPrefix string
			}{
				"no filtering": {
					pkg:          noExportPkg,
					pkgs:         noExportPkgs,
					onlyExported: false,
					prefix:       "type",
					widgetPrefix: "widget",
				},
				"filtering": {
					pkg:          exportPkg,
					pkgs:         exportPkgs,
					onlyExported: true,
					prefix:       "Type",
					widgetPrefix: "Widget",
				},
			}

			for name, tc := range testCases {
				t.Run(name, func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t, false)
					defer afterEach(t)

					metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
					loadFnMoq.onCall(loadCfg, tc.pkg).returnResults(tc.pkgs, nil)
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
					metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

					err := cache.LoadPackage(tc.pkg)
					if err != nil {
						t.Fatalf("got %#v, want no error", err)
					}

					// ACT
					typs := cache.MockableTypes(tc.onlyExported)

					// ASSERT
					if len(typs) != 7 {
						t.Fatalf("got %d types, want 7", len(typs))
					}

					typsByName := map[string]dst.Ident{}
					for _, typ := range typs {
						id := typ
						typsByName[id.Name] = id
					}

					if _, ok := typsByName[tc.prefix+"1"]; !ok {
						t.Errorf("got nothing, want %s1", tc.prefix)
					}
					if _, ok := typsByName[tc.prefix+"2"]; !ok {
						t.Errorf("got nothing, want %s2", tc.prefix)
					}
					if _, ok := typsByName[tc.prefix+"4_genType"]; !ok {
						t.Errorf("got nothing, want %s4_genType", tc.prefix)
					}
					if _, ok := typsByName[tc.widgetPrefix+"_genType"]; !ok {
						t.Errorf("got nothing, want %s_genType", tc.widgetPrefix)
					}
					if _, ok := typsByName[tc.widgetPrefix+"_starGenType"]; !ok {
						t.Errorf("got nothing, want %s_starGenType", tc.widgetPrefix)
					}
					// TODO: Check new extensions
				})
			}
		})

		t.Run("filtering to nothing", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, noExportPkg).returnResults(noExportPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			err := cache.LoadPackage(noExportPkg)
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
			beforeEach(t, false)
			defer afterEach(t)

			pkgs := append(exportPkgs, pkg("vendor/other-pkgs"))

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, exportPkg).returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults().Repeat(moq.AnyTimes())

			err := cache.LoadPackage(exportPkg)
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			// ACT
			typs := cache.MockableTypes(true)

			// ASSERT
			if len(typs) != 7 {
				t.Fatalf("got %d types, want 7", len(typs))
			}
		})

		t.Run("always filters internal packages", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			pkgs := append(exportPkgs,
				pkg("elsewhere/internal"),
				pkg("elsewhere/internal/again"),
				pkg("internal"))

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, exportPkg).returnResults(pkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults().Repeat(moq.AnyTimes())

			err := cache.LoadPackage(exportPkg)
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			// ACT
			typs := cache.MockableTypes(true)

			// ASSERT
			if len(typs) != 7 {
				t.Fatalf("got %d types, want 7", len(typs))
			}
		})
	})
}
