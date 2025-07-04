package generator_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dave/dst"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"moqueries.org/runtime/moq"

	"moqueries.org/cli/ast"
	"moqueries.org/cli/generator"
)

func TestConverter(t *testing.T) {
	var (
		scene        *moq.Scene
		typeCacheMoq *moqTypeCache

		func1Param1  *dst.Ident
		func1Param2  *dst.Ident
		func1Param3  *dst.Ident
		func1Param4  *dst.Ident
		func1Params  *dst.FieldList
		func1Results *dst.FieldList

		iSpec      *dst.TypeSpec
		iSpecFuncs []generator.Func

		fnSpec      *dst.TypeSpec
		fnSpecFuncs []generator.Func

		titler = cases.Title(language.Und, cases.NoLower)
	)

	beforeEach := func(t *testing.T) {
		t.Helper()
		if scene != nil {
			t.Fatal("afterEach not called")
		}
		scene = moq.NewScene(t)
		typeCacheMoq = newMoqTypeCache(scene, nil)

		func1Param1 = dst.NewIdent("pType1")
		func1Param2 = dst.NewIdent("pType2")
		func1Param3 = dst.NewIdent("pType3")
		func1Param4 = dst.NewIdent("pType4")
		func1Params = &dst.FieldList{
			List: []*dst.Field{
				{Names: []*dst.Ident{dst.NewIdent("firstParam")}, Type: func1Param1},
				{Names: []*dst.Ident{dst.NewIdent("secondParam")}, Type: func1Param2},
				{Names: []*dst.Ident{dst.NewIdent("thirdParam")}, Type: func1Param3},
				{Names: []*dst.Ident{dst.NewIdent("fourthParam")}, Type: func1Param4},
			},
		}
		func1Results = &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent("result")},
					Type:  dst.NewIdent("string"),
				},
				{
					Names: []*dst.Ident{dst.NewIdent("err")},
					Type:  dst.NewIdent("error"),
				},
			},
		}

		iSpec = &dst.TypeSpec{
			Name: dst.NewIdent("PublicInterface"),
			Type: &dst.InterfaceType{
				Methods: &dst.FieldList{
					List: []*dst.Field{
						{
							Names: []*dst.Ident{dst.NewIdent("Func1")},
							Type: &dst.FuncType{
								Params:  func1Params,
								Results: func1Results,
							},
						},
						{
							Names: []*dst.Ident{dst.NewIdent("func2")},
							Type:  &dst.FuncType{},
						},
					},
				},
			},
		}
		iSpecFuncs = []generator.Func{
			{
				Name: "Func1",
				FuncType: &dst.FuncType{
					Params:  func1Params,
					Results: func1Results,
				},
			},
			{
				Name: "func2",
				FuncType: &dst.FuncType{
					Params: &dst.FieldList{List: nil},
				},
			},
		}

		fnSpec = &dst.TypeSpec{
			Name: dst.NewIdent("PublicFunction"),
			Type: &dst.FuncType{
				Params:  func1Params,
				Results: func1Results,
			},
		}
		fnSpecFuncs = []generator.Func{
			{
				FuncType: &dst.FuncType{
					Params:  func1Params,
					Results: func1Results,
				},
			},
		}
	}

	afterEach := func(t *testing.T) {
		t.Helper()
		scene.AssertExpectationsMet()
		scene = nil
	}

	entries := map[string]bool{
		"exported": true,
		"private":  false,
	}

	for name, isExported := range entries {
		exported := func(name string) string {
			if isExported {
				name = titler.String(name)
			}
			return name
		}

		t.Run(name, func(t *testing.T) {
			t.Run("BaseDecls", func(t *testing.T) {
				t.Run("creates a base moq for an interface", func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t)
					defer afterEach(t)

					typ := generator.Type{
						TypeInfo: ast.TypeInfo{
							Type:    iSpec,
							PkgPath: "thatmodule/pkg",
						},
						Funcs: iSpecFuncs,
					}
					converter := generator.NewConverter(typ, isExported, typeCacheMoq.mock())

					// ACT
					decls, err := converter.BaseDecls()
					// ASSERT
					if err != nil {
						t.Fatalf("got %#v, want no error", err)
					}
					if len(decls) != 2 {
						t.Fatalf("got %#v, want 1 decl", decls)
					}
					decl := decls[0].(*dst.GenDecl)
					if len(decl.Decs.Start) != 2 {
						t.Errorf("got len %d, wanted 2", len(decl.Decs.Start))
					}
					expectedStart := "// The following type assertion assures that pkg.PublicInterface is mocked"
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
					expectedStart = "// completely"
					if decl.Decs.Start[1] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[1], expectedStart)
					}
					decl = decls[1].(*dst.GenDecl)
					expectedStart = fmt.Sprintf("// %s holds the state of a moq of the PublicInterface type",
						exported("moqPublicInterface"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
				})

				t.Run("creates a base moq for a function", func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t)
					defer afterEach(t)

					typ := generator.Type{
						TypeInfo: ast.TypeInfo{Type: fnSpec},
						Funcs:    fnSpecFuncs,
					}
					converter := generator.NewConverter(typ, isExported, typeCacheMoq.mock())

					// ACT
					decls, err := converter.BaseDecls()
					// ASSERT
					if err != nil {
						t.Fatalf("got %#v, want no error", err)
					}
					if len(decls) != 1 {
						t.Fatalf("got %#v, want 1 decl", decls)
					}
					decl := decls[0].(*dst.GenDecl)
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart := fmt.Sprintf("// %s holds the state of a moq of the PublicFunction type",
						exported("moqPublicFunction"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
				})

				t.Run("includes the interface when mocking a fabricated interface", func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t)
					defer afterEach(t)

					typ := generator.Type{
						TypeInfo: ast.TypeInfo{
							Type:       iSpec,
							PkgPath:    "thatmodule/pkg",
							Fabricated: true,
						},
						Funcs: iSpecFuncs,
					}
					typeCacheMoq.onCall().Type(dst.Ident{}, "thatmodule/pkg", false).
						any().id().returnResults(ast.TypeInfo{
						Type: &dst.TypeSpec{Name: dst.NewIdent("fab")},
					}, nil).repeat(moq.AnyTimes())
					converter := generator.NewConverter(typ, isExported, typeCacheMoq.mock())

					// ACT
					decls, err := converter.BaseDecls()
					// ASSERT
					if err != nil {
						t.Fatalf("got %#v, want no error", err)
					}
					if len(decls) != 3 {
						t.Fatalf("got %#v, want 3 decl", decls)
					}
					decl := decls[0].(*dst.GenDecl)
					if len(decl.Decs.Start) != 2 {
						t.Errorf("got len %d, wanted 2", len(decl.Decs.Start))
					}
					expectedStart := "// The following type assertion assures that pkg.PublicInterface is mocked"
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
					expectedStart = "// completely"
					if decl.Decs.Start[1] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[1], expectedStart)
					}
					decl = decls[1].(*dst.GenDecl)
					if len(decl.Decs.Start) != 3 {
						t.Errorf("got len %d, wanted 3", len(decl.Decs.Start))
					}
					expectedStart = "// PublicInterface is the fabricated implementation type of this mock (emitted"
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
					expectedStart = "// when mocking a collections of methods directly and not from an interface"
					if decl.Decs.Start[1] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[1], expectedStart)
					}
					expectedStart = "// type)"
					if decl.Decs.Start[2] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[2], expectedStart)
					}
				})

				t.Run("includes the func type when mocking a fabricated func type", func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t)
					defer afterEach(t)

					typ := generator.Type{
						TypeInfo: ast.TypeInfo{Type: fnSpec, Fabricated: true},
						Funcs:    fnSpecFuncs,
					}
					typeCacheMoq.onCall().Type(dst.Ident{}, "", false).
						any().id().returnResults(ast.TypeInfo{
						Type: &dst.TypeSpec{Name: dst.NewIdent("fab")},
					}, nil).repeat(moq.AnyTimes())
					converter := generator.NewConverter(typ, isExported, typeCacheMoq.mock())

					// ACT
					decls, err := converter.BaseDecls()
					// ASSERT
					if err != nil {
						t.Fatalf("got %#v, want no error", err)
					}
					if len(decls) != 2 {
						t.Fatalf("got %#v, want 1 decl", decls)
					}
					decl := decls[0].(*dst.GenDecl)
					if len(decl.Decs.Start) != 2 {
						t.Errorf("got len %d, wanted 2", len(decl.Decs.Start))
					}
					expectedStart := "// PublicFunction is the fabricated implementation type of this mock (emitted"
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
					expectedStart = "// when mocking functions directly and not from a function type)"
					if decl.Decs.Start[1] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[1], expectedStart)
					}
				})
			})

			t.Run("MockStructs", func(t *testing.T) {
				t.Run("create structs", func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t)
					defer afterEach(t)

					typ := generator.Type{
						TypeInfo: ast.TypeInfo{Type: &dst.TypeSpec{
							Name: dst.NewIdent("MyInterface"),
							Type: &dst.InterfaceType{},
						}},
					}
					converter := generator.NewConverter(typ, isExported, typeCacheMoq.mock())

					// ACT
					decls, err := converter.MockStructs()
					// ASSERT
					if err != nil {
						t.Fatalf("got %#v, want no error", err)
					}
					if len(decls) != 3 {
						t.Fatalf("got %#v, want 1 decl", decls)
					}

					decl := decls[0].(*dst.GenDecl)
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart := fmt.Sprintf("// %s isolates the mock interface of the MyInterface type",
						exported("moqMyInterface_mock"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}

					decl = decls[1].(*dst.GenDecl)
					if len(decl.Decs.Start) < 2 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart = fmt.Sprintf("// %s isolates the recorder interface of the MyInterface",
						exported("moqMyInterface_recorder"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
					expectedStart = "// type"
					if decl.Decs.Start[1] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[1], expectedStart)
					}

					decl = decls[2].(*dst.GenDecl)
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart = fmt.Sprintf("// %s holds runtime configuration for the MyInterface type",
						exported("moqMyInterface_runtime"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
				})
			})

			t.Run("MethodStructs", func(t *testing.T) {
				t.Run("creates structs for a function", func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t)
					defer afterEach(t)

					expectedParams, ok := dst.Clone(func1Params).(*dst.FieldList)
					if !ok {
						t.Fatalf("type assertion failed")
					}
					// Non-comparables are represented as a deep hash
					expectedParams.List[0].Type = &dst.Ident{
						Path: "moqueries.org/runtime/hash",
						Name: "Hash",
					}
					expectedParams.List[3].Type = &dst.Ident{
						Path: "moqueries.org/runtime/hash",
						Name: "Hash",
					}
					fn := generator.Func{
						Name: "Func1",
						FuncType: &dst.FuncType{
							Params:  func1Params,
							Results: func1Results,
						},
					}
					typeCacheMoq.onCall().Type(*func1Param1, "", false).
						returnResults(ast.TypeInfo{
							Type: &dst.TypeSpec{Name: func1Param1},
						}, nil).repeat(moq.AnyTimes())
					typeCacheMoq.onCall().Type(*func1Param2, "", false).
						returnResults(ast.TypeInfo{
							Type: &dst.TypeSpec{Name: func1Param2},
						}, nil).repeat(moq.AnyTimes())
					typeCacheMoq.onCall().Type(*func1Param3, "", false).
						returnResults(ast.TypeInfo{
							Type: &dst.TypeSpec{Name: func1Param3},
						}, nil).repeat(moq.AnyTimes())
					typeCacheMoq.onCall().Type(*func1Param4, "", false).
						returnResults(ast.TypeInfo{
							Type: &dst.TypeSpec{Name: func1Param4},
						}, nil).repeat(moq.AnyTimes())
					tInfo := ast.TypeInfo{Type: iSpec}
					typeCacheMoq.onCall().IsComparable(func1Param1, tInfo).
						returnResults(false, nil)
					typeCacheMoq.onCall().IsComparable(func1Param2, tInfo).
						returnResults(true, nil)
					typeCacheMoq.onCall().IsComparable(func1Param3, tInfo).
						returnResults(false, nil)
					typeCacheMoq.onCall().IsComparable(func1Param4, tInfo).
						returnResults(true, nil)
					id := ast.Id("string")
					typeCacheMoq.onCall().Type(*id, "", false).
						returnResults(ast.TypeInfo{
							Type: &dst.TypeSpec{Name: id},
						}, nil).repeat(moq.AnyTimes())
					id = ast.Id("error")
					typeCacheMoq.onCall().Type(*id, "", false).
						returnResults(ast.TypeInfo{
							Type: &dst.TypeSpec{Name: id},
						}, nil).repeat(moq.AnyTimes())

					typ := generator.Type{
						TypeInfo: ast.TypeInfo{Type: iSpec},
					}
					converter := generator.NewConverter(typ, isExported, typeCacheMoq.mock())

					// ACT
					decls, err := converter.MethodStructs(fn)
					// ASSERT
					if err != nil {
						t.Errorf("got %#v, wanted nil err", err)
					}
					if len(decls) != 9 {
						t.Errorf("got len %d, wanted 9", len(decls))
					}

					decl, ok := decls[0].(*dst.GenDecl)
					if !ok {
						t.Errorf("got %#v, wanted GenDecl type", decls[0])
					}
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart := fmt.Sprintf("// %s adapts %s as needed by the",
						exported("moqPublicInterface_Func1_adaptor"),
						exported("moqPublicInterface"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}

					decl, ok = decls[1].(*dst.GenDecl)
					if !ok {
						t.Errorf("got %#v, wanted GenDecl type", decls[0])
					}
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart = fmt.Sprintf("// %s holds the params of the PublicInterface type",
						exported("moqPublicInterface_Func1_params"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}

					decl, ok = decls[2].(*dst.GenDecl)
					if !ok {
						t.Errorf("got %#v, wanted GenDecl type", decls[1])
					}
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart = fmt.Sprintf("// %s holds the map key params of the",
						exported("moqPublicInterface_Func1_paramsKey"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
					expectedStart = "// PublicInterface type"
					if decl.Decs.Start[1] != expectedStart {
						t.Errorf("got %s, want %s", decl.Decs.Start[0], expectedStart)
					}

					decl, ok = decls[3].(*dst.GenDecl)
					if !ok {
						t.Errorf("got %#v, wanted GenDecl type", decls[2])
					}
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart = fmt.Sprintf("// %s holds the results of the PublicInterface",
						exported("moqPublicInterface_Func1_results"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
					expectedStart = "// type"
					if decl.Decs.Start[1] != expectedStart {
						t.Errorf("got %s, want %s", decl.Decs.Start[1], expectedStart)
					}

					decl, ok = decls[4].(*dst.GenDecl)
					if !ok {
						t.Errorf("got %#v, wanted GenDecl type", decls[3])
					}
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart = fmt.Sprintf("// %s holds the parameter indexing runtime",
						exported("moqPublicInterface_Func1_paramIndexing"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
					expectedStart = "// configuration for the PublicInterface type"
					if decl.Decs.Start[1] != expectedStart {
						t.Errorf("got %s, want %s", decl.Decs.Start[1], expectedStart)
					}

					decl, ok = decls[5].(*dst.GenDecl)
					if !ok {
						t.Errorf("got %#v, wanted GenDecl type", decls[4])
					}
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart = fmt.Sprintf("// %s defines the type of function needed when",
						exported("moqPublicInterface_Func1_doFn"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
					expectedStart = fmt.Sprintf("// calling %s for the PublicInterface type",
						exported("andDo"))
					if decl.Decs.Start[1] != expectedStart {
						t.Errorf("got %s, want %s", decl.Decs.Start[1], expectedStart)
					}

					decl, ok = decls[6].(*dst.GenDecl)
					if !ok {
						t.Errorf("got %#v, wanted GenDecl type", decls[5])
					}
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart = fmt.Sprintf("// %s defines the type of function needed when",
						exported("moqPublicInterface_Func1_doReturnFn"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
					expectedStart = fmt.Sprintf("// calling %s for the PublicInterface type",
						exported("doReturnResults"))
					if decl.Decs.Start[1] != expectedStart {
						t.Errorf("got %s, want %s", decl.Decs.Start[1], expectedStart)
					}

					decl, ok = decls[7].(*dst.GenDecl)
					if !ok {
						t.Errorf("got %#v, wanted GenDecl type", decls[6])
					}
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart = fmt.Sprintf("// %s routes recorded function calls to the",
						exported("moqPublicInterface_Func1_recorder"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
					expectedStart = fmt.Sprintf("// %s moq",
						exported("moqPublicInterface"))
					if decl.Decs.Start[1] != expectedStart {
						t.Errorf("got %s, want %s", decl.Decs.Start[1], expectedStart)
					}

					decl, ok = decls[8].(*dst.GenDecl)
					if !ok {
						t.Errorf("got %#v, wanted GenDecl type", decls[7])
					}
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart = fmt.Sprintf("// %s isolates the any params functions of the",
						exported("moqPublicInterface_Func1_anyParams"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
					expectedStart = "// PublicInterface type"
					if decl.Decs.Start[1] != expectedStart {
						t.Errorf("got %s, want %s", decl.Decs.Start[1], expectedStart)
					}
				})

				t.Run("returns errors", func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t)
					defer afterEach(t)

					fn := generator.Func{
						Name: "Func1",
						FuncType: &dst.FuncType{
							Params:  func1Params,
							Results: func1Results,
						},
					}
					expectedErr := errors.New("type error")
					typeCacheMoq.onCall().Type(*func1Param1, "", false).
						returnResults(ast.TypeInfo{}, expectedErr).repeat(moq.AnyTimes())
					typeCacheMoq.onCall().Type(dst.Ident{}, "", false).any().id().
						returnResults(ast.TypeInfo{
							Type: &dst.TypeSpec{Name: func1Param1},
						}, nil).repeat(moq.AnyTimes())
					typeCacheMoq.onCall().IsComparable(nil, ast.TypeInfo{}).
						any().expr().any().parentType().
						returnResults(false, nil).repeat(moq.AnyTimes())

					typ := generator.Type{
						TypeInfo: ast.TypeInfo{Type: iSpec},
					}
					converter := generator.NewConverter(typ, isExported, typeCacheMoq.mock())

					// ACT
					decls, err := converter.MethodStructs(fn)

					// ASSERT
					if !errors.Is(err, expectedErr) {
						t.Errorf("got %#v, wanted %#v", err, expectedErr)
					}
					if decls != nil {
						t.Errorf("got %#v, wanted nil", decls)
					}
				})
			})

			t.Run("NewFunc", func(t *testing.T) {
				t.Run("creates a new moq function for an interface", func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t)
					defer afterEach(t)

					typ := generator.Type{
						TypeInfo: ast.TypeInfo{Type: iSpec},
					}
					converter := generator.NewConverter(typ, isExported, typeCacheMoq.mock())

					// ACT
					decl, err := converter.NewFunc()
					// ASSERT
					if err != nil {
						t.Errorf("got error %#v, want no error", err)
					}
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart := fmt.Sprintf("// %s creates a new moq of the PublicInterface type",
						exported("newMoqPublicInterface"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
				})

				t.Run("creates a new moq function for a function", func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t)
					defer afterEach(t)

					typ := generator.Type{
						TypeInfo: ast.TypeInfo{Type: fnSpec},
					}
					converter := generator.NewConverter(typ, isExported, typeCacheMoq.mock())

					// ACT
					decl, err := converter.NewFunc()
					// ASSERT
					if err != nil {
						t.Errorf("got error %#v, want no error", err)
					}
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
					}
					expectedStart := fmt.Sprintf("// %s creates a new moq of the PublicFunction type",
						exported("newMoqPublicFunction"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
				})

				t.Run("returns errors", func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t)
					defer afterEach(t)

					typ := generator.Type{
						TypeInfo: ast.TypeInfo{Type: iSpec},
						Funcs:    iSpecFuncs,
					}
					converter := generator.NewConverter(typ, isExported, typeCacheMoq.mock())

					expectedErr := errors.New("type error")
					typeCacheMoq.onCall().Type(*func1Param1, "", false).
						returnResults(ast.TypeInfo{}, expectedErr)
					typeCacheMoq.onCall().Type(dst.Ident{}, "", false).any().id().
						returnResults(ast.TypeInfo{}, expectedErr).repeat(moq.AnyTimes())
					typeCacheMoq.onCall().IsDefaultComparable(nil, ast.TypeInfo{}).
						any().expr().any().parentType().
						returnResults(false, nil).repeat(moq.AnyTimes())

					// ACT
					decl, err := converter.NewFunc()

					// ASSERT
					if !errors.Is(err, expectedErr) {
						t.Errorf("got %#v, wanted %#v", err, expectedErr)
					}
					if decl != nil {
						t.Errorf("got %#v, wanted nil", decl)
					}
				})
			})

			t.Run("IsolationAccessor", func(t *testing.T) {
				t.Run("creates a func", func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t)
					defer afterEach(t)

					typ := generator.Type{
						TypeInfo: ast.TypeInfo{Type: &dst.TypeSpec{
							Name: dst.NewIdent("MyInterface"),
						}},
					}
					converter := generator.NewConverter(typ, isExported, typeCacheMoq.mock())

					// ACT
					decl, err := converter.IsolationAccessor("recorder", "onCall")
					// ASSERT
					if err != nil {
						t.Fatalf("got error %#v, want no error", err)
					}
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted > 0", len(decl.Decs.Start))
					}
					expectedStart := fmt.Sprintf("// %s returns the recorder implementation of the MyInterface type",
						exported("onCall"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
				})
			})

			t.Run("FuncClosure", func(t *testing.T) {
				t.Run("creates a closure function for a function type", func(t *testing.T) {
					// ASSEMBLE
					beforeEach(t)
					defer afterEach(t)

					typ := generator.Type{
						TypeInfo: ast.TypeInfo{Type: &dst.TypeSpec{
							Name: dst.NewIdent("MyFn"),
						}},
					}
					converter := generator.NewConverter(typ, isExported, typeCacheMoq.mock())

					typeCacheMoq.onCall().Type(*func1Param1, "", false).
						returnResults(ast.TypeInfo{
							Type: &dst.TypeSpec{Name: func1Param1},
						}, nil).repeat(moq.AnyTimes())
					typeCacheMoq.onCall().Type(*func1Param2, "", false).
						returnResults(ast.TypeInfo{
							Type: &dst.TypeSpec{Name: func1Param2},
						}, nil).repeat(moq.AnyTimes())
					typeCacheMoq.onCall().Type(*func1Param3, "", false).
						returnResults(ast.TypeInfo{
							Type: &dst.TypeSpec{Name: func1Param3},
						}, nil).repeat(moq.AnyTimes())
					typeCacheMoq.onCall().Type(*func1Param4, "", false).
						returnResults(ast.TypeInfo{
							Type: &dst.TypeSpec{Name: func1Param4},
						}, nil).repeat(moq.AnyTimes())
					id := ast.Id("string")
					typeCacheMoq.onCall().Type(*id, "", false).
						returnResults(ast.TypeInfo{
							Type: &dst.TypeSpec{Name: id},
						}, nil).repeat(moq.AnyTimes())
					id = ast.Id("error")
					typeCacheMoq.onCall().Type(*id, "", false).
						returnResults(ast.TypeInfo{
							Type: &dst.TypeSpec{Name: id},
						}, nil).repeat(moq.AnyTimes())
					typeCacheMoq.onCall().IsComparable(func1Param1, typ.TypeInfo).
						returnResults(true, nil)
					typeCacheMoq.onCall().IsComparable(func1Param2, typ.TypeInfo).
						returnResults(true, nil)
					typeCacheMoq.onCall().IsComparable(func1Param3, typ.TypeInfo).
						returnResults(true, nil)
					typeCacheMoq.onCall().IsComparable(func1Param4, typ.TypeInfo).
						returnResults(true, nil)

					// ACT
					decl, err := converter.FuncClosure(fnSpecFuncs[0])
					// ASSERT
					if err != nil {
						t.Fatalf("got %#v, want no error", err)
					}
					if len(decl.Decs.Start) < 1 {
						t.Errorf("got len %d, wanted > 0", len(decl.Decs.Start))
					}
					expectedStart := fmt.Sprintf("// %s returns the moq implementation of the MyFn type",
						exported("mock"))
					if decl.Decs.Start[0] != expectedStart {
						t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
					}
				})
			})
		})
	}
}
