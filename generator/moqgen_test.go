package generator_test

import (
	"errors"
	"testing"

	"github.com/dave/dst"

	"github.com/myshkin5/moqueries/ast"
	"github.com/myshkin5/moqueries/generator"
	"github.com/myshkin5/moqueries/moq"
)

func TestMoqGenerator(t *testing.T) {
	var (
		scene             *moq.Scene
		typeCacheMoq      *moqTypeCache
		getwdFnMoq        *moqGetwdFunc
		newConverterFnMoq *moqNewConverterFunc
		converter1Moq     *moqConverterer
		converter2Moq     *moqConverterer

		gen *generator.MoqGenerator

		ifaceSpec1    *dst.TypeSpec
		ifaceSpec2    *dst.TypeSpec
		ifaceMethods1 *dst.FieldList
		ifaceMethods2 *dst.FieldList
		func1         *dst.Field
		func1Params   *dst.FieldList

		readFnType *dst.FuncType
		readerSpec *dst.TypeSpec

		fnSpec *dst.TypeSpec
	)

	beforeEach := func(t *testing.T) {
		t.Helper()

		if scene != nil {
			t.Fatal("afterEach not called")
		}
		scene = moq.NewScene(t)
		typeCacheMoq = newMoqTypeCache(scene, nil)
		getwdFnMoq = newMoqGetwdFunc(scene, nil)
		newConverterFnMoq = newMoqNewConverterFunc(scene, nil)
		converter1Moq = newMoqConverterer(scene, nil)
		converter2Moq = newMoqConverterer(scene, nil)

		gen = generator.New(
			typeCacheMoq.mock(),
			getwdFnMoq.mock(),
			newConverterFnMoq.mock())

		func1Params = &dst.FieldList{List: []*dst.Field{
			{
				Names: []*dst.Ident{dst.NewIdent("firstParm")},
				Type: &dst.StarExpr{X: &dst.SelectorExpr{
					X:   dst.NewIdent("cobra"),
					Sel: dst.NewIdent("Command"),
				}},
			},
			{Type: dst.NewIdent("string")},
			{
				Type: &dst.StarExpr{X: &dst.SelectorExpr{
					X:   dst.NewIdent("dst"),
					Sel: dst.NewIdent("InterfaceType"),
				}},
			},
		}}
		func1 = &dst.Field{
			Names: []*dst.Ident{dst.NewIdent("Func1")},
			Type:  &dst.FuncType{Params: func1Params, Results: nil},
		}
		ifaceMethods1 = &dst.FieldList{List: []*dst.Field{func1}}
		ifaceSpec1 = &dst.TypeSpec{
			Name: dst.NewIdent("PublicInterface"),
			Type: &dst.InterfaceType{Methods: ifaceMethods1},
		}
		ifaceMethods2 = &dst.FieldList{}
		ifaceSpec2 = &dst.TypeSpec{
			Name: dst.NewIdent("privateInterface"),
			Type: &dst.InterfaceType{Methods: ifaceMethods2},
		}

		readFnType = &dst.FuncType{
			Params: &dst.FieldList{List: []*dst.Field{{
				Names: []*dst.Ident{dst.NewIdent("p")},
				Type:  &dst.ArrayType{Elt: dst.NewIdent("byte")},
			}}},
			Results: &dst.FieldList{List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent("n")},
					Type:  dst.NewIdent("int"),
				},
				{
					Names: []*dst.Ident{dst.NewIdent("err")},
					Type:  dst.NewIdent("error"),
				},
			}},
		}
		readerSpec = &dst.TypeSpec{
			Name: &dst.Ident{Name: "Reader", Path: "io"},
			Type: &dst.InterfaceType{Methods: &dst.FieldList{
				List: []*dst.Field{{
					Names: []*dst.Ident{dst.NewIdent("Read")},
					Type:  readFnType,
				}},
			}},
		}

		fnSpec = &dst.TypeSpec{
			Name: dst.NewIdent("PublicFn"),
			Type: &dst.FuncType{Params: func1Params, Results: nil},
		}
	}

	afterEach := func(t *testing.T) {
		t.Helper()
		scene.AssertExpectationsMet()
		scene = nil
	}

	t.Run("always returns a header comment", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach(t)
		req := generator.GenerateRequest{
			Types:       nil,
			Export:      false,
			Destination: "dir/file_test.go",
			Package:     "",
			Import:      ".",
			TestImport:  false,
			WorkingDir:  "/some-nice-path",
		}

		typeCacheMoq.onCall().FindPackage("dir").returnResults("myrepo.com/dir", nil)
		getwdFnMoq.onCall().returnResults("/some-nice-path", nil)

		// ACT
		resp, err := gen.Generate(req)
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted nil err", err)
		}
		if resp.File == nil {
			t.Fatalf("got nil file, wanted not nil")
		}
		if len(resp.File.Decs.Start) < 1 {
			t.Errorf("got %d, wanted > 0 len start", len(resp.File.Decs.Start))
		}
		expectedStart := "// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!"
		if resp.File.Decs.Start[0] != expectedStart {
			t.Errorf("got %s, wanted %s", resp.File.Decs.Start[0], expectedStart)
		}
		if resp.DestPath != "dir/file_test.go" {
			t.Errorf("got %s, wanted dir/file_test.go", resp.DestPath)
		}
		if resp.OutPkgPath != "myrepo.com/dir_test" {
			t.Errorf("got %s, wanted myrepo.com/dir_test", resp.OutPkgPath)
		}
	})

	t.Run("can put mocks in parent packages when given a relative destination", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach(t)
		req := generator.GenerateRequest{
			Types:       nil,
			Export:      true,
			Destination: "../file_test.go",
			Package:     "",
			Import:      ".",
			TestImport:  false,
			WorkingDir:  "/some-nice-path",
		}

		typeCacheMoq.onCall().FindPackage("..").returnResults("myrepo.com/otherpkg", nil)
		getwdFnMoq.onCall().returnResults("/some-nice-path", nil)

		// ACT
		resp, err := gen.Generate(req)
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted nil err", err)
		}
		if resp.File.Name.Name != "otherpkg" {
			t.Errorf("got %s, wanted otherpkg", resp.File.Name.Name)
		}
		if resp.DestPath != "../file_test.go" {
			t.Errorf("got %s, wanted ../file_test.go", resp.DestPath)
		}
		if resp.OutPkgPath != "myrepo.com/otherpkg" {
			t.Errorf("got %s, wanted myrepo.com/otherpkg", resp.OutPkgPath)
		}
	})

	t.Run("package naming based on current directory and export flag", func(t *testing.T) {
		type testCase struct {
			export bool
			pkg    string
		}
		testCases := map[string]testCase{
			"test package when not exported": {export: false, pkg: "thispkg_test"},
			"non-test package when exported": {export: true, pkg: "thispkg"},
		}

		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				beforeEach(t)
				defer afterEach(t)
				req := generator.GenerateRequest{
					Types:       nil,
					Export:      tc.export,
					Destination: "file_test.go",
					Package:     "",
					Import:      ".",
					TestImport:  false,
					WorkingDir:  "/some-nice-path",
				}

				typeCacheMoq.onCall().FindPackage(".").returnResults("myrepo.com/thispkg", nil)
				getwdFnMoq.onCall().returnResults("/some-nice-path", nil)

				// ACT
				resp, err := gen.Generate(req)
				// ASSERT
				if err != nil {
					t.Errorf("got %#v, wanted nil err", err)
				}
				if resp.File.Name.Name != tc.pkg {
					t.Errorf("got %s, wanted %s", resp.File.Name.Name, tc.pkg)
				}
				if resp.DestPath != "file_test.go" {
					t.Errorf("got %s, wanted file_test.go", resp.DestPath)
				}
				outPkgPath := "myrepo.com/" + tc.pkg
				if resp.OutPkgPath != outPkgPath {
					t.Errorf("got %s, wanted %s", resp.OutPkgPath, outPkgPath)
				}
			})
		}
	})

	t.Run("recursively looks up nested interfaces", func(t *testing.T) {
		types := []string{"PublicInterface", "privateInterface"}
		type testCase struct {
			request    generator.GenerateRequest
			findPkgDir string
			findPkgOut string
			outPkgPath string
			getwdDir   string
			typePath   string
			destPath   string
		}
		testCases := map[string]testCase{
			"current working dir same as req working dir": {
				request: generator.GenerateRequest{
					Types:       types,
					Export:      false,
					Destination: "file_test.go",
					Package:     "",
					Import:      ".",
					TestImport:  false,
					WorkingDir:  "/some-nice-path",
				},
				findPkgDir: ".",
				findPkgOut: "thispkg",
				outPkgPath: "thispkg_test",
				getwdDir:   "/some-nice-path",
				typePath:   ".",
				destPath:   "file_test.go",
			},
			"gives a sane path when destination is an absolute path": {
				request: generator.GenerateRequest{
					Types:       types,
					Export:      false,
					Destination: "/exactly-here/file_test.go",
					Package:     "",
					Import:      ".",
					TestImport:  false,
					WorkingDir:  "/some-nice-path",
				},
				findPkgDir: "exactly-here",
				findPkgOut: "thispkg",
				outPkgPath: "thispkg_test",
				getwdDir:   "/some-nice-path",
				typePath:   ".",
				destPath:   "/exactly-here/file_test.go",
			},
			"gives a sane path when destination dir is an absolute path": {
				request: generator.GenerateRequest{
					Types:          types,
					Export:         false,
					DestinationDir: "/exactly-here",
					Package:        "",
					Import:         ".",
					TestImport:     false,
					WorkingDir:     "/some-nice-path",
				},
				findPkgDir: "exactly-here",
				findPkgOut: "thispkg",
				outPkgPath: "thispkg_test",
				getwdDir:   "/some-nice-path",
				typePath:   ".",
				destPath:   "/exactly-here/moq_publicinterface_privateinterface_test.go",
			},
			"req working dir not set": {
				request: generator.GenerateRequest{
					Types:       types,
					Export:      false,
					Destination: "file_test.go",
					Package:     "",
					Import:      ".",
					TestImport:  false,
				},
				findPkgDir: ".",
				findPkgOut: "thispkg",
				outPkgPath: "thispkg_test",
				getwdDir:   "/some-nice-path",
				typePath:   ".",
				destPath:   "file_test.go",
			},
			"current working dir parent of req working dir": {
				request: generator.GenerateRequest{
					Types:       types,
					Export:      false,
					Destination: "file_test.go",
					Package:     "",
					Import:      ".",
					TestImport:  false,
					WorkingDir:  "/some-nice-path/some-child-dir",
				},
				findPkgDir: "some-child-dir",
				findPkgOut: "thispkg",
				outPkgPath: "thispkg_test",
				getwdDir:   "/some-nice-path",
				typePath:   "./some-child-dir",
				destPath:   "some-child-dir/file_test.go",
			},
			"specific absolute import": {
				request: generator.GenerateRequest{
					Types:       types,
					Export:      false,
					Destination: "file_test.go",
					Package:     "",
					Import:      "io",
					TestImport:  false,
					WorkingDir:  "/some-nice-path",
				},
				findPkgDir: ".",
				findPkgOut: "thispkg",
				outPkgPath: "thispkg_test",
				getwdDir:   "/some-nice-path",
				typePath:   "io",
				destPath:   "file_test.go",
			},
			"specific relative import": {
				request: generator.GenerateRequest{
					Types:       types,
					Export:      false,
					Destination: "file_test.go",
					Package:     "",
					Import:      "./somechilddir",
					TestImport:  false,
					WorkingDir:  "/some-nice-path",
				},
				findPkgDir: ".",
				findPkgOut: "thispkg",
				outPkgPath: "thispkg_test",
				getwdDir:   "/some-nice-path",
				typePath:   "./somechilddir",
				destPath:   "file_test.go",
			},
			"current working dir parent of req working dir w/ relative import": {
				request: generator.GenerateRequest{
					Types:       types,
					Export:      false,
					Destination: "file_test.go",
					Package:     "",
					Import:      "./secondchilddir",
					TestImport:  false,
					WorkingDir:  "/some-nice-path/firstchilddir",
				},
				findPkgDir: "firstchilddir",
				findPkgOut: "thispkg",
				outPkgPath: "thispkg_test",
				getwdDir:   "/some-nice-path",
				typePath:   "./firstchilddir/secondchilddir",
				destPath:   "firstchilddir/file_test.go",
			},
			"warns (?!) when trying to export a test file": {
				request: generator.GenerateRequest{
					Types:       types,
					Export:      true,
					Destination: "file_test.go",
					Package:     "",
					Import:      ".",
					TestImport:  false,
					WorkingDir:  "/some-nice-path",
				},
				findPkgDir: ".",
				findPkgOut: "thispkg",
				outPkgPath: "thispkg",
				getwdDir:   "/some-nice-path",
				typePath:   ".",
				destPath:   "file_test.go",
			},
			"generated filename for a test package": {
				request: generator.GenerateRequest{
					Types:          types,
					Export:         false,
					DestinationDir: "subpkg",
					Package:        "",
					Import:         ".",
					TestImport:     false,
					WorkingDir:     "/some-nice-path",
				},
				findPkgDir: "subpkg",
				findPkgOut: "thispkg",
				outPkgPath: "thispkg_test",
				getwdDir:   "/some-nice-path",
				typePath:   ".",
				destPath:   "subpkg/moq_publicinterface_privateinterface_test.go",
			},
			"generated filename for an exported package": {
				request: generator.GenerateRequest{
					Types:          types,
					Export:         true,
					DestinationDir: "subpkg",
					Package:        "",
					Import:         ".",
					TestImport:     false,
					WorkingDir:     "/some-nice-path",
				},
				findPkgDir: "subpkg",
				findPkgOut: "thispkg",
				outPkgPath: "thispkg",
				getwdDir:   "/some-nice-path",
				typePath:   ".",
				destPath:   "subpkg/moq_publicinterface_privateinterface.go",
			},
			"generated filename w/ a relative path": {
				request: generator.GenerateRequest{
					Types:          types,
					Export:         true,
					DestinationDir: "destpkg3",
					Package:        "subpkg2",
					Import:         ".",
					TestImport:     false,
					WorkingDir:     "/some-nice-path/subdir1",
				},
				findPkgDir: "subdir1/destpkg3",
				findPkgOut: "thispkg",
				outPkgPath: "subpkg2",
				getwdDir:   "/some-nice-path",
				typePath:   "./subdir1",
				destPath:   "subdir1/destpkg3/moq_publicinterface_privateinterface.go",
			},
		}

		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				beforeEach(t)
				defer afterEach(t)

				typeCacheMoq.onCall().FindPackage(tc.findPkgDir).
					returnResults(tc.findPkgOut, nil)
				getwdFnMoq.onCall().returnResults(tc.getwdDir, nil)

				// PublicInterface embeds privateInterface which embeds io.Reader
				ifaceMethods1.List = append(ifaceMethods1.List, &dst.Field{
					Type: &dst.Ident{
						Name: "privateInterface",
						Path: "github.com/myshkin5/moqueries/generator",
					},
				})
				ifaceMethods2.List = append(ifaceMethods2.List, &dst.Field{
					Type: &dst.Ident{
						Name: "Reader",
						Path: "io",
					},
				})
				typeCacheMoq.onCall().Type(*ast.IdPath("PublicInterface", tc.typePath), false).
					returnResults(ifaceSpec1, "github.com/myshkin5/moqueries/generator", nil)
				typeCacheMoq.onCall().Type(*ast.IdPath("privateInterface", "github.com/myshkin5/moqueries/generator"), false).
					returnResults(ifaceSpec2, "github.com/myshkin5/moqueries/generator", nil)
				typeCacheMoq.onCall().Type(*ast.IdPath("Reader", "io"), false).
					returnResults(readerSpec, "", nil)
				typeCacheMoq.onCall().Type(*ast.IdPath("privateInterface", tc.typePath), false).
					returnResults(ifaceSpec2, "github.com/myshkin5/moqueries/generator", nil)
				typeCacheMoq.onCall().Type(*ast.IdPath("Reader", "io"), false).
					returnResults(readerSpec, "", nil)
				ifaceFuncs := []generator.Func{
					{Name: "Func1", Params: func1Params},
					{
						Name:    "Read",
						Params:  readFnType.Params,
						Results: readFnType.Results,
					},
				}
				newConverterFnMoq.onCall(generator.Type{
					TypeSpec:   ifaceSpec1,
					Funcs:      ifaceFuncs,
					InPkgPath:  "github.com/myshkin5/moqueries/generator",
					OutPkgPath: tc.outPkgPath,
				}, tc.request.Export).returnResults(converter1Moq.mock())
				converter1Moq.onCall().BaseStruct().returnResults(&dst.GenDecl{Specs: []dst.Spec{&dst.TypeSpec{
					Name: dst.NewIdent("pub-decl"),
				}}})
				converter1Moq.onCall().IsolationStruct("mock").
					returnResults(nil)
				converter1Moq.onCall().IsolationStruct("recorder").
					returnResults(nil)
				converter1Moq.onCall().MethodStructs(ifaceFuncs[0]).
					returnResults(nil, nil)
				converter1Moq.onCall().MethodStructs(ifaceFuncs[1]).
					returnResults(nil, nil)
				converter1Moq.onCall().NewFunc().
					returnResults(nil)
				converter1Moq.onCall().IsolationAccessor("mock", "mock").
					returnResults(nil)
				converter1Moq.onCall().MockMethod(ifaceFuncs[0]).
					returnResults(nil)
				converter1Moq.onCall().MockMethod(ifaceFuncs[1]).
					returnResults(nil)
				converter1Moq.onCall().IsolationAccessor("recorder", "onCall").
					returnResults(nil)
				converter1Moq.onCall().RecorderMethods(ifaceFuncs[0]).
					returnResults(nil)
				converter1Moq.onCall().RecorderMethods(ifaceFuncs[1]).
					returnResults(nil)
				converter1Moq.onCall().ResetMethod().
					returnResults(nil)
				converter1Moq.onCall().AssertMethod().
					returnResults(nil)
				iface2Funcs := []generator.Func{{
					Name:    "Read",
					Params:  readFnType.Params,
					Results: readFnType.Results,
				}}
				newConverterFnMoq.onCall(generator.Type{
					TypeSpec:   ifaceSpec2,
					Funcs:      iface2Funcs,
					InPkgPath:  "github.com/myshkin5/moqueries/generator",
					OutPkgPath: tc.outPkgPath,
				}, tc.request.Export).returnResults(converter2Moq.mock())
				converter2Moq.onCall().BaseStruct().returnResults(nil)
				converter2Moq.onCall().IsolationStruct("mock").
					returnResults(nil)
				converter2Moq.onCall().IsolationStruct("recorder").
					returnResults(nil)
				converter2Moq.onCall().MethodStructs(iface2Funcs[0]).
					returnResults(nil, nil)
				converter2Moq.onCall().NewFunc().
					returnResults(nil)
				converter2Moq.onCall().IsolationAccessor("mock", "mock").
					returnResults(nil)
				converter2Moq.onCall().MockMethod(iface2Funcs[0]).
					returnResults(nil)
				converter2Moq.onCall().IsolationAccessor("recorder", "onCall").
					returnResults(nil)
				converter2Moq.onCall().RecorderMethods(iface2Funcs[0]).
					returnResults(nil)
				converter2Moq.onCall().ResetMethod().
					returnResults(nil)
				converter2Moq.onCall().AssertMethod().
					returnResults(nil)

				// ACT
				resp, err := gen.Generate(tc.request)
				// ASSERT
				if err != nil {
					t.Errorf("got %#v, wanted nil err", err)
				}
				if resp.DestPath != tc.destPath {
					t.Errorf("got %s, wanted %s", resp.DestPath, tc.destPath)
				}
				if resp.OutPkgPath != tc.outPkgPath {
					t.Errorf("got %s, wanted %s", resp.OutPkgPath, tc.outPkgPath)
				}
			})
		}
	})

	t.Run("ErrNonExported processing", func(t *testing.T) {
		t.Run("unexported top-level type", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)
			defer afterEach(t)

			typeCacheMoq.onCall().FindPackage(".").
				returnResults("destdir", nil)
			getwdFnMoq.onCall().returnResults("/some-nice-path", nil)

			req := generator.GenerateRequest{
				Types:              []string{"privateInterface"},
				Export:             true,
				Destination:        "file_test.go",
				Package:            "",
				Import:             ".",
				TestImport:         false,
				WorkingDir:         "/some-nice-path",
				ErrorOnNonExported: true,
			}

			// ACT
			resp, err := gen.Generate(req)

			// ASSERT
			if err == nil {
				t.Fatal("got no error, wanted error")
			}
			expectedMsg := "non-exported types: privateInterface mocked type is not exported"
			if err.Error() != expectedMsg {
				t.Errorf("got %s, wanted %s", err.Error(), expectedMsg)
			}
			if resp.File != nil {
				t.Errorf("got %#v, wanted nil", resp.File)
			}
			if resp.DestPath != "" {
				t.Errorf("got %s, wanted nothing", resp.DestPath)
			}
			if resp.OutPkgPath != "" {
				t.Errorf("got %s, wanted nothing", resp.OutPkgPath)
			}
		})

		t.Run("unexported method", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)
			defer afterEach(t)

			typeCacheMoq.onCall().FindPackage(".").
				returnResults("destdir", nil)
			getwdFnMoq.onCall().returnResults("/some-nice-path", nil)

			func1.Names[0].Name = "func1"
			typeCacheMoq.onCall().Type(*ast.IdPath("MyInterface", "."), false).
				returnResults(ifaceSpec1, "github.com/myshkin5/moqueries/generator", nil)
			req := generator.GenerateRequest{
				Types:              []string{"MyInterface"},
				Export:             true,
				Destination:        "file_test.go",
				Package:            "",
				Import:             ".",
				TestImport:         false,
				WorkingDir:         "/some-nice-path",
				ErrorOnNonExported: true,
			}

			// ACT
			resp, err := gen.Generate(req)

			// ASSERT
			if err == nil {
				t.Fatal("got no error, wanted error")
			}
			expectedMsg := "non-exported types: func1 method is not exported"
			if err.Error() != expectedMsg {
				t.Errorf("got %s, wanted %s", err.Error(), expectedMsg)
			}
			if resp.File != nil {
				t.Errorf("got %#v, wanted nil", resp.File)
			}
			if resp.DestPath != "" {
				t.Errorf("got %s, wanted nothing", resp.DestPath)
			}
			if resp.OutPkgPath != "" {
				t.Errorf("got %s, wanted nothing", resp.OutPkgPath)
			}
		})

		t.Run("unexported embedded interface", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)
			defer afterEach(t)

			typeCacheMoq.onCall().FindPackage(".").
				returnResults("destdir", nil)
			getwdFnMoq.onCall().returnResults("/some-nice-path", nil)

			// PublicInterface embeds privateInterface
			ifaceMethods1.List = append(ifaceMethods1.List, &dst.Field{
				Type: &dst.Ident{
					Name: "privateInterface",
					Path: "github.com/myshkin5/moqueries/generator",
				},
			})
			typeCacheMoq.onCall().Type(*ast.IdPath("PublicInterface", "."), false).
				returnResults(ifaceSpec1, "github.com/myshkin5/moqueries/generator", nil)
			privateId := ast.IdPath("privateInterface", "github.com/myshkin5/moqueries/generator")
			typeCacheMoq.onCall().Type(*privateId, false).
				returnResults(ifaceSpec2, "github.com/myshkin5/moqueries/generator", nil)
			req := generator.GenerateRequest{
				Types:              []string{"PublicInterface"},
				Export:             true,
				Destination:        "file_test.go",
				Package:            "",
				Import:             ".",
				TestImport:         false,
				WorkingDir:         "/some-nice-path",
				ErrorOnNonExported: true,
			}

			// ACT
			resp, err := gen.Generate(req)

			// ASSERT
			if err == nil {
				t.Fatal("got no error, wanted error")
			}
			expectedMsg := "non-exported types: privateInterface embedded type is not exported"
			if err.Error() != expectedMsg {
				t.Errorf("got %s, wanted %s", err.Error(), expectedMsg)
			}
			if resp.File != nil {
				t.Errorf("got %#v, wanted nil", resp.File)
			}
			if resp.DestPath != "" {
				t.Errorf("got %s, wanted nothing", resp.DestPath)
			}
			if resp.OutPkgPath != "" {
				t.Errorf("got %s, wanted nothing", resp.OutPkgPath)
			}
		})
	})

	t.Run("successfully navigates type aliases", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach(t)

		typeCacheMoq.onCall().FindPackage(".").returnResults("thispkg", nil)
		getwdFnMoq.onCall().returnResults("/some-nice-path", nil)

		ifaceSpec := &dst.TypeSpec{
			Name: dst.NewIdent("AliasType"),
			Type: ast.IdPath("Reader", "io"),
		}

		typeCacheMoq.onCall().Type(*ast.IdPath("AliasType", "."), false).
			returnResults(ifaceSpec, "github.com/myshkin5/moqueries/generator", nil)
		typeCacheMoq.onCall().Type(*ast.IdPath("Reader", "io"), false).
			returnResults(readerSpec, "", nil)

		ifaceFuncs := []generator.Func{
			{
				Name:    "Read",
				Params:  readFnType.Params,
				Results: readFnType.Results,
			},
		}
		newConverterFnMoq.onCall(generator.Type{
			TypeSpec:   ifaceSpec,
			Funcs:      ifaceFuncs,
			InPkgPath:  "github.com/myshkin5/moqueries/generator",
			OutPkgPath: "thispkg_test",
		}, false).returnResults(converter1Moq.mock())
		converter1Moq.onCall().BaseStruct().returnResults(nil)
		converter1Moq.onCall().IsolationStruct("mock").
			returnResults(nil)
		converter1Moq.onCall().IsolationStruct("recorder").
			returnResults(nil)
		converter1Moq.onCall().MethodStructs(ifaceFuncs[0]).
			returnResults(nil, nil)
		converter1Moq.onCall().NewFunc().
			returnResults(nil)
		converter1Moq.onCall().IsolationAccessor("mock", "mock").
			returnResults(nil)
		converter1Moq.onCall().MockMethod(ifaceFuncs[0]).
			returnResults(nil)
		converter1Moq.onCall().IsolationAccessor("recorder", "onCall").
			returnResults(nil)
		converter1Moq.onCall().RecorderMethods(ifaceFuncs[0]).
			returnResults(nil)
		converter1Moq.onCall().ResetMethod().
			returnResults(nil)
		converter1Moq.onCall().AssertMethod().
			returnResults(nil)

		req := generator.GenerateRequest{
			Types:       []string{"AliasType"},
			Export:      false,
			Destination: "file_test.go",
			Package:     "",
			Import:      ".",
			TestImport:  false,
			WorkingDir:  "/some-nice-path",
		}

		// ACT
		_, err := gen.Generate(req)
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted nil err", err)
		}
	})

	t.Run("returns an os.Getwd error", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach(t)

		getwdFnMoq.onCall().returnResults("/some-nice-path", errors.New("os.Getwd-error"))

		req := generator.GenerateRequest{
			Types:       []string{"PublicInterface"},
			Export:      false,
			Destination: "file_test.go",
			Package:     "",
			Import:      ".",
			TestImport:  false,
			WorkingDir:  "/some-nice-path",
		}

		// ACT
		resp, err := gen.Generate(req)

		// ASSERT
		if err == nil {
			t.Fatal("got no error, wanted error")
		}
		expectedMsg := "error getting current working directory: os.Getwd-error"
		if err.Error() != expectedMsg {
			t.Errorf("got %s, wanted %s", err.Error(), expectedMsg)
		}
		if resp.File != nil {
			t.Errorf("got %#v, wanted nil", resp.File)
		}
		if resp.DestPath != "" {
			t.Errorf("got %s, wanted nothing", resp.DestPath)
		}
		if resp.OutPkgPath != "" {
			t.Errorf("got %s, wanted nothing", resp.OutPkgPath)
		}
	})

	t.Run("returns an ErrInvalidConfig when defining both destination"+
		" and destination directory", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach(t)

		typeCacheMoq.onCall().FindPackage("moqdir").returnResults("thispkg", nil)
		getwdFnMoq.onCall().returnResults("/some-nice-path", nil)

		req := generator.GenerateRequest{
			Types:          []string{"PublicInterface"},
			Export:         false,
			Destination:    "file_test.go",
			DestinationDir: "./moqdir",
			Package:        "",
			Import:         ".",
			TestImport:     false,
			WorkingDir:     "/some-nice-path",
		}

		// ACT
		resp, err := gen.Generate(req)

		// ASSERT
		if err == nil {
			t.Fatal("got no error, wanted error")
		}
		if !errors.Is(err, generator.ErrInvalidConfig) {
			t.Errorf("got %#v, wanted %#v", err, generator.ErrInvalidConfig)
		}
		expectedMsg := "invalid configuration: both --destination and" +
			" --destination-dir flags must not be present together"
		if err.Error() != expectedMsg {
			t.Errorf("got %s, wanted %s", err.Error(), expectedMsg)
		}
		if resp.File != nil {
			t.Errorf("got %#v, wanted nil", resp.File)
		}
		if resp.DestPath != "" {
			t.Errorf("got %s, wanted nothing", resp.DestPath)
		}
		if resp.OutPkgPath != "" {
			t.Errorf("got %s, wanted nothing", resp.OutPkgPath)
		}
	})

	t.Run("returns a convertor error", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach(t)

		typeCacheMoq.onCall().FindPackage(".").returnResults("thispkg", nil)
		getwdFnMoq.onCall().returnResults("/some-nice-path", nil)
		typeCacheMoq.onCall().Type(*ast.IdPath("PublicInterface", "."), false).
			returnResults(ifaceSpec1, "github.com/myshkin5/moqueries/generator", nil)
		ifaceFuncs := []generator.Func{{Name: "Func1", Params: func1Params}}
		newConverterFnMoq.onCall(generator.Type{
			TypeSpec:   ifaceSpec1,
			Funcs:      ifaceFuncs,
			InPkgPath:  "github.com/myshkin5/moqueries/generator",
			OutPkgPath: "thispkg_test",
		}, false).returnResults(converter1Moq.mock())
		converter1Moq.onCall().BaseStruct().returnResults(&dst.GenDecl{Specs: []dst.Spec{&dst.TypeSpec{
			Name: dst.NewIdent("pub-decl"),
		}}})
		converter1Moq.onCall().IsolationStruct("mock").
			returnResults(nil)
		converter1Moq.onCall().IsolationStruct("recorder").
			returnResults(nil)
		expectedErr := errors.New("bad convertor")
		converter1Moq.onCall().MethodStructs(ifaceFuncs[0]).
			returnResults(nil, expectedErr)

		req := generator.GenerateRequest{
			Types:       []string{"PublicInterface"},
			Export:      false,
			Destination: "file_test.go",
			Package:     "",
			Import:      ".",
			TestImport:  false,
			WorkingDir:  "/some-nice-path",
		}

		// ACT
		resp, err := gen.Generate(req)

		// ASSERT
		if err != expectedErr {
			t.Errorf("got %#v, wanted %#v", err, expectedErr)
		}
		if resp.File != nil {
			t.Errorf("got %#v, wanted nil", resp.File)
		}
		if resp.DestPath != "" {
			t.Errorf("got %s, wanted nothing", resp.DestPath)
		}
		if resp.OutPkgPath != "" {
			t.Errorf("got %s, wanted nothing", resp.OutPkgPath)
		}
	})

	t.Run("loads tests types when requested", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach(t)

		typeCacheMoq.onCall().FindPackage(".").returnResults("thispkg", nil)
		getwdFnMoq.onCall().returnResults("/some-nice-path", nil)
		typeCacheMoq.onCall().Type(*ast.IdPath("PublicInterface", "."), true).
			returnResults(ifaceSpec1, "github.com/myshkin5/moqueries/generator", nil)
		typeCacheMoq.onCall().Type(*ast.IdPath("privateInterface", "."), true).
			returnResults(ifaceSpec2, "github.com/myshkin5/moqueries/generator", nil)

		iface1Funcs := []generator.Func{{Name: "Func1", Params: func1Params}}
		newConverterFnMoq.onCall(generator.Type{
			TypeSpec:   ifaceSpec1,
			Funcs:      iface1Funcs,
			InPkgPath:  "github.com/myshkin5/moqueries/generator",
			OutPkgPath: "thispkg_test",
		}, false).returnResults(converter1Moq.mock())
		converter1Moq.onCall().BaseStruct().returnResults(nil)
		converter1Moq.onCall().IsolationStruct("mock").
			returnResults(nil)
		converter1Moq.onCall().IsolationStruct("recorder").
			returnResults(nil)
		converter1Moq.onCall().MethodStructs(iface1Funcs[0]).
			returnResults(nil, nil)
		converter1Moq.onCall().NewFunc().
			returnResults(nil)
		converter1Moq.onCall().IsolationAccessor("mock", "mock").
			returnResults(nil)
		converter1Moq.onCall().MockMethod(iface1Funcs[0]).
			returnResults(nil)
		converter1Moq.onCall().IsolationAccessor("recorder", "onCall").
			returnResults(nil)
		converter1Moq.onCall().RecorderMethods(iface1Funcs[0]).
			returnResults(nil)
		converter1Moq.onCall().ResetMethod().
			returnResults(nil)
		converter1Moq.onCall().AssertMethod().
			returnResults(nil)

		var iface2Funcs []generator.Func
		newConverterFnMoq.onCall(generator.Type{
			TypeSpec:   ifaceSpec2,
			Funcs:      iface2Funcs,
			InPkgPath:  "github.com/myshkin5/moqueries/generator",
			OutPkgPath: "thispkg_test",
		}, false).returnResults(converter2Moq.mock())
		converter2Moq.onCall().BaseStruct().returnResults(nil)
		converter2Moq.onCall().IsolationStruct("mock").
			returnResults(nil)
		converter2Moq.onCall().IsolationStruct("recorder").
			returnResults(nil)
		converter2Moq.onCall().NewFunc().
			returnResults(nil)
		converter2Moq.onCall().IsolationAccessor("mock", "mock").
			returnResults(nil)
		converter2Moq.onCall().IsolationAccessor("recorder", "onCall").
			returnResults(nil)
		converter2Moq.onCall().ResetMethod().
			returnResults(nil)
		converter2Moq.onCall().AssertMethod().
			returnResults(nil)

		req := generator.GenerateRequest{
			Types:       []string{"PublicInterface", "privateInterface"},
			Export:      false,
			Destination: "file_test.go",
			Package:     "",
			Import:      ".",
			TestImport:  true,
			WorkingDir:  "/some-nice-path",
		}

		// ACT
		_, err := gen.Generate(req)
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted nil err", err)
		}
	})

	t.Run("returns an error when the type cache returns an error", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach(t)

		typeCacheMoq.onCall().FindPackage(".").returnResults("thispkg", nil)
		getwdFnMoq.onCall().returnResults("/some-nice-path", nil)
		expectedErr := errors.New("bad cache")
		typeCacheMoq.onCall().Type(*ast.IdPath("BadInterface", "."), false).
			returnResults(nil, "", expectedErr)

		req := generator.GenerateRequest{
			Types:       []string{"BadInterface"},
			Export:      false,
			Destination: "file_test.go",
			Package:     "",
			Import:      ".",
			TestImport:  false,
			WorkingDir:  "/some-nice-path",
		}

		// ACT
		resp, err := gen.Generate(req)

		// ASSERT
		if err != expectedErr {
			t.Errorf("got %#v, wanted %#v", err, expectedErr)
		}
		if resp.File != nil {
			t.Errorf("got %#v, wanted nil", resp.File)
		}
		if resp.DestPath != "" {
			t.Errorf("got %s, wanted nothing", resp.DestPath)
		}
		if resp.OutPkgPath != "" {
			t.Errorf("got %s, wanted nothing", resp.OutPkgPath)
		}
	})

	t.Run("returns an error when recursive lookups return an error", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach(t)

		typeCacheMoq.onCall().FindPackage(".").returnResults("thispkg", nil)
		getwdFnMoq.onCall().returnResults("/some-nice-path", nil)
		// PublicInterface embeds privateInterface
		ifaceMethods1.List = append(ifaceMethods1.List, &dst.Field{
			Type: &dst.Ident{
				Name: "privateInterface",
				Path: "github.com/myshkin5/moqueries/generator",
			},
		})
		typeCacheMoq.onCall().Type(*ast.IdPath("PublicInterface", "."), false).
			returnResults(ifaceSpec1, "github.com/myshkin5/moqueries/generator", nil)
		expectedErr := errors.New("bad cache")
		typeCacheMoq.onCall().Type(*ast.IdPath("privateInterface", "github.com/myshkin5/moqueries/generator"), false).
			returnResults(nil, "", expectedErr)

		req := generator.GenerateRequest{
			Types:       []string{"PublicInterface"},
			Export:      false,
			Destination: "file_test.go",
			Package:     "",
			Import:      ".",
			TestImport:  false,
			WorkingDir:  "/some-nice-path",
		}

		// ACT
		resp, err := gen.Generate(req)

		// ASSERT
		if err != expectedErr {
			t.Errorf("got %#v, wanted %#v", err, expectedErr)
		}
		if resp.File != nil {
			t.Errorf("got %#v, wanted nil", resp.File)
		}
		if resp.DestPath != "" {
			t.Errorf("got %s, wanted nothing", resp.DestPath)
		}
		if resp.OutPkgPath != "" {
			t.Errorf("got %s, wanted nothing", resp.OutPkgPath)
		}
	})

	t.Run("returns an error when recursive function lookups return an error", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach(t)

		typeCacheMoq.onCall().FindPackage(".").returnResults("thispkg", nil)
		getwdFnMoq.onCall().returnResults("/some-nice-path", nil)

		// PublicInterface embeds privateInterface which embeds io.Reader
		ifaceMethods1.List = append(ifaceMethods1.List, &dst.Field{
			Type: &dst.Ident{
				Name: "privateInterface",
				Path: "github.com/myshkin5/moqueries/generator",
			},
		})
		ifaceMethods2.List = append(ifaceMethods2.List, &dst.Field{
			Type: &dst.Ident{
				Name: "Reader",
				Path: "io",
			},
		})
		typeCacheMoq.onCall().Type(*ast.IdPath("PublicInterface", "."), false).
			returnResults(ifaceSpec1, "github.com/myshkin5/moqueries/generator", nil)
		typeCacheMoq.onCall().Type(*ast.IdPath("privateInterface", "github.com/myshkin5/moqueries/generator"), false).
			returnResults(ifaceSpec2, "github.com/myshkin5/moqueries/generator", nil)
		expectedErr := errors.New("bad cache")
		typeCacheMoq.onCall().Type(*ast.IdPath("Reader", "io"), false).
			returnResults(nil, "", expectedErr)

		req := generator.GenerateRequest{
			Types:       []string{"PublicInterface"},
			Export:      false,
			Destination: "file_test.go",
			Package:     "",
			Import:      ".",
			TestImport:  false,
			WorkingDir:  "/some-nice-path",
		}

		// ACT
		resp, err := gen.Generate(req)

		// ASSERT
		if err != expectedErr {
			t.Errorf("got %#v, wanted %#v", err, expectedErr)
		}
		if resp.File != nil {
			t.Errorf("got %#v, wanted nil", resp.File)
		}
		if resp.DestPath != "" {
			t.Errorf("got %s, wanted nothing", resp.DestPath)
		}
		if resp.OutPkgPath != "" {
			t.Errorf("got %s, wanted nothing", resp.OutPkgPath)
		}
	})

	t.Run("handles function types", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach(t)

		typeCacheMoq.onCall().FindPackage(".").returnResults("thispkg", nil)
		getwdFnMoq.onCall().returnResults("/some-nice-path", nil)
		typeCacheMoq.onCall().Type(*ast.IdPath("PublicFn", "."), false).
			returnResults(fnSpec, "github.com/myshkin5/moqueries/generator", nil)
		fnFuncs := []generator.Func{{Params: func1Params}}
		newConverterFnMoq.onCall(generator.Type{
			TypeSpec:   fnSpec,
			Funcs:      fnFuncs,
			InPkgPath:  "github.com/myshkin5/moqueries/generator",
			OutPkgPath: "thispkg_test",
		}, false).returnResults(converter1Moq.mock())
		converter1Moq.onCall().BaseStruct().returnResults(&dst.GenDecl{Specs: []dst.Spec{&dst.TypeSpec{
			Name: dst.NewIdent("pub-decl"),
		}}})
		converter1Moq.onCall().IsolationStruct("mock").
			returnResults(nil)
		converter1Moq.onCall().MethodStructs(fnFuncs[0]).
			returnResults(nil, nil)
		converter1Moq.onCall().NewFunc().
			returnResults(nil)
		converter1Moq.onCall().FuncClosure(fnFuncs[0]).
			returnResults(nil)
		converter1Moq.onCall().MockMethod(fnFuncs[0]).
			returnResults(nil)
		converter1Moq.onCall().RecorderMethods(fnFuncs[0]).
			returnResults(nil)
		converter1Moq.onCall().ResetMethod().
			returnResults(nil)
		converter1Moq.onCall().AssertMethod().
			returnResults(nil)

		req := generator.GenerateRequest{
			Types:       []string{"PublicFn"},
			Export:      false,
			Destination: "file_test.go",
			Package:     "",
			Import:      ".",
			TestImport:  false,
			WorkingDir:  "/some-nice-path",
		}

		// ACT
		_, err := gen.Generate(req)
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted nil err", err)
		}
	})
}
