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
		scene      *moq.Scene
		loadFnMoq  *moqLoadFn
		metricsMoq *metrics.MoqMetrics

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
					metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, exportPkg).returnResults(exportPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
					metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, noExportPkg).returnResults(noExportPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().
				ReturnResults().Repeat(moq.MinTimes(1))
			type1Id := ast.IdPath("type1", noExportPkg)
			_, _ = cache.Type(*type1Id, noExportPkg, false)

			loadCfg.Tests = true

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
			metricsMoq.OnCall().ASTTypeCacheMissesInc().
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
					metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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

					subTestCases := map[string]struct {
						compFn     func(c *ast.Cache, expr dst.Expr) (bool, error)
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
					}

					for name, stc := range subTestCases {
						t.Run(name, func(t *testing.T) {
							// ASSEMBLE
							beforeEach(t, false)
							defer afterEach(t)

							metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.Optional())
							metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults().Repeat(moq.Optional())
							loadFnMoq.onCall(loadCfg, "io").
								returnResults(ioPkgs, nil).repeat(moq.AnyTimes())
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

		comparableStructCases := map[string]struct {
			code    string
			declIdx int
		}{
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
							defer afterEach(t)

							f := parse(t, fmt.Sprintf(stc.code, paramType))
							fn, ok := f.Decls[stc.declIdx].(*dst.FuncDecl)
							if !ok {
								t.Fatalf("got %#v, want a function declaration", f.Decls[stc.declIdx])
							}
							expr := fn.Type.Params.List[0].Type

							metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.Optional())
							metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults().Repeat(moq.Optional())
							loadFnMoq.onCall(loadCfg, "io").
								returnResults(ioPkgs, nil).repeat(moq.AnyTimes())
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
					metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults().Repeat(moq.Optional())
					loadFnMoq.onCall(loadCfg, "b").returnResults(bPkg, nil)
					metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults().Repeat(moq.Optional())

					metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.Optional())
					metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults().Repeat(moq.Optional())
					loadFnMoq.onCall(loadCfg, "io").
						returnResults(ioPkgs, nil).repeat(moq.AnyTimes())
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
		beforeEach(t, false)
		defer afterEach(t)

		c := ast.NewCache(packages.Load, metricsMoq.Mock())

		metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults().Repeat(moq.AnyTimes())
		metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults().Repeat(moq.AnyTimes())
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
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, "./testpkgs/noexport").returnResults(noExportPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

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

		t.Run("abs dir", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			dir, err := os.Getwd()
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}
			absDir := filepath.Join(dir, "testpkgs/noexport")
			loadFnMoq.onCall(loadCfg, absDir).returnResults(noExportPkgs, nil)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			// ACT
			pkgPath, err := cache.FindPackage(absDir)
			// ASSERT
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			if pkgPath != noExportPkg {
				t.Errorf("got %s, want %s", pkgPath, noExportPkg)
			}
		})

		t.Run("current dir", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
			loadFnMoq.onCall(loadCfg, ".").doReturnResults(packages.Load)
			metricsMoq.OnCall().ASTTotalLoadTimeInc(0).Any().D().ReturnResults()
			metricsMoq.OnCall().ASTTotalDecorationTimeInc(0).Any().D().ReturnResults()

			// ACT
			pkgPath, err := cache.FindPackage(".")
			// ASSERT
			if err != nil {
				t.Fatalf("got %#v, want no error", err)
			}

			if pkgPath != "moqueries.org/cli/ast" {
				t.Errorf("got %s, want moqueries.org/cli/ast", pkgPath)
			}
		})
	})

	t.Run("LoadPackage", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
					metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
					if len(typs) != 5 {
						t.Fatalf("got %d types, want 5", len(typs))
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
				})
			}
		})

		t.Run("filtering to nothing", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t, false)
			defer afterEach(t)

			metricsMoq.OnCall().ASTPkgCacheMissesInc().ReturnResults()
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
			if len(typs) != 5 {
				t.Fatalf("got %d types, want 5", len(typs))
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
			metricsMoq.OnCall().ASTTypeCacheMissesInc().ReturnResults()
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
			if len(typs) != 5 {
				t.Fatalf("got %d types, want 5", len(typs))
			}
		})
	})
}
