package generator_test

import (
	"errors"
	"testing"

	"github.com/dave/dst"

	"github.com/myshkin5/moqueries/generator"
	"github.com/myshkin5/moqueries/moq"
)

func TestConverterPublic(t *testing.T) {
	var (
		scene        *moq.Scene
		typeCacheMoq *moqTypeCache

		converter *generator.Converter

		func1Param1  dst.Expr
		func1Param2  dst.Expr
		func1Param3  dst.Expr
		func1Param4  dst.Expr
		func1Params  *dst.FieldList
		func1Results *dst.FieldList

		iSpec      *dst.TypeSpec
		iSpecFuncs []generator.Func

		fnSpec      *dst.TypeSpec
		fnSpecFuncs []generator.Func
	)

	beforeEach := func(t *testing.T) {
		if scene != nil {
			t.Fatal("afterEach not called")
		}
		scene = moq.NewScene(t)
		typeCacheMoq = newMoqTypeCache(scene, nil)

		converter = generator.NewConverter(true, typeCacheMoq.mock())

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
				Name:    "Func1",
				Params:  func1Params,
				Results: func1Results,
			},
			{
				Name: "func2",
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
				Params:  func1Params,
				Results: func1Results,
			},
		}
	}

	afterEach := func() {
		scene.AssertExpectationsMet()
		scene = nil
	}

	t.Run("BaseStruct", func(t *testing.T) {
		t.Run("creates a base moq for an interface", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)

			// ACT
			decl := converter.BaseStruct(iSpec, iSpecFuncs)

			// ASSERT
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
			}
			expectedStart := "// MoqPublicInterface holds the state of a moq of the PublicInterface type"
			if decl.Decs.Start[0] != expectedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
			}
			afterEach()
		})

		t.Run("creates a base moq for a function", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)

			// ACT
			decl := converter.BaseStruct(fnSpec, fnSpecFuncs)

			// ASSERT
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
			}
			expwctedStart := "// MoqPublicFunction holds the state of a moq of the PublicFunction type"
			if decl.Decs.Start[0] != expwctedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expwctedStart)
			}
			afterEach()
		})
	})

	t.Run("IsolationStruct", func(t *testing.T) {
		t.Run("creates a struct", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)

			// ACT
			decl := converter.IsolationStruct("MyInterface", "mock")

			// ASSERT
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
			}
			expwctedStart := "// MoqMyInterface_mock isolates the mock interface of the MyInterface type"
			if decl.Decs.Start[0] != expwctedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expwctedStart)
			}
			afterEach()
		})
	})

	t.Run("MethodStructs", func(t *testing.T) {
		t.Run("creates structs for a function", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)
			expectedParams := dst.Clone(func1Params).(*dst.FieldList)
			// Non-comparables are represented as a deep hash
			expectedParams.List[0].Type = &dst.Ident{
				Path: "github.com/myshkin5/moqueries/hash",
				Name: "Hash",
			}
			expectedParams.List[3].Type = &dst.Ident{
				Path: "github.com/myshkin5/moqueries/hash",
				Name: "Hash",
			}
			fn := generator.Func{
				Name:    "Func1",
				Params:  func1Params,
				Results: func1Results,
			}
			typeCacheMoq.onCall().IsComparable(func1Param1).
				returnResults(false, nil)
			typeCacheMoq.onCall().IsComparable(func1Param2).
				returnResults(true, nil)
			typeCacheMoq.onCall().IsComparable(func1Param3).
				returnResults(false, nil)
			typeCacheMoq.onCall().IsComparable(func1Param4).
				returnResults(true, nil)

			// ACT
			decls, err := converter.MethodStructs(iSpec, fn)
			// ASSERT
			if err != nil {
				t.Errorf("got %#v, wanted nil err", err)
			}
			if len(decls) != 8 {
				t.Errorf("got len %d, wanted 8", len(decls))
			}
			decl, ok := decls[0].(*dst.GenDecl)
			if !ok {
				t.Errorf("got %#v, wanted GenDecl type", decls[0])
			}
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
			}
			expectedStart := "// MoqPublicInterface_Func1_params holds the params of the PublicInterface type"
			if decl.Decs.Start[0] != expectedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
			}

			decl, ok = decls[1].(*dst.GenDecl)
			if !ok {
				t.Errorf("got %#v, wanted GenDecl type", decls[1])
			}
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
			}
			expectedStart = "// MoqPublicInterface_Func1_paramsKey holds the map key params of the PublicInterface type"
			if decl.Decs.Start[0] != expectedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
			}

			decl, ok = decls[2].(*dst.GenDecl)
			if !ok {
				t.Errorf("got %#v, wanted GenDecl type", decls[2])
			}
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
			}
			expectedStart = "// MoqPublicInterface_Func1_resultsByParams " +
				"contains the results for a given set of parameters for the PublicInterface type"
			if decl.Decs.Start[0] != expectedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
			}

			decl, ok = decls[3].(*dst.GenDecl)
			if !ok {
				t.Errorf("got %#v, wanted GenDecl type", decls[3])
			}
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
			}
			expectedStart = "// MoqPublicInterface_Func1_doFn " +
				"defines the type of function needed when calling AndDo for the PublicInterface type"
			if decl.Decs.Start[0] != expectedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
			}

			decl, ok = decls[4].(*dst.GenDecl)
			if !ok {
				t.Errorf("got %#v, wanted GenDecl type", decls[4])
			}
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
			}
			expectedStart = "// MoqPublicInterface_Func1_doReturnFn " +
				"defines the type of function needed when calling DoReturnResults for the PublicInterface type"
			if decl.Decs.Start[0] != expectedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
			}

			decl, ok = decls[5].(*dst.GenDecl)
			if !ok {
				t.Errorf("got %#v, wanted GenDecl type", decls[5])
			}
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
			}
			expectedStart = "// MoqPublicInterface_Func1_results " +
				"holds the results of the PublicInterface type"
			if decl.Decs.Start[0] != expectedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
			}

			decl, ok = decls[6].(*dst.GenDecl)
			if !ok {
				t.Errorf("got %#v, wanted GenDecl type", decls[6])
			}
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
			}
			expectedStart = "// MoqPublicInterface_Func1_fnRecorder " +
				"routes recorded function calls to the MoqPublicInterface moq"
			if decl.Decs.Start[0] != expectedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
			}

			decl, ok = decls[7].(*dst.GenDecl)
			if !ok {
				t.Errorf("got %#v, wanted GenDecl type", decls[7])
			}
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted < 1", len(decl.Decs.Start))
			}
			expectedStart = "// MoqPublicInterface_Func1_anyParams " +
				"isolates the any params functions of the PublicInterface type"
			if decl.Decs.Start[0] != expectedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
			}
			afterEach()
		})

		t.Run("returns a type cache error", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)
			fn := generator.Func{
				Name:    "Func1",
				Params:  func1Params,
				Results: func1Results,
			}
			expectedErr := errors.New("type error")
			typeCacheMoq.onCall().IsComparable(func1Param1).
				returnResults(false, expectedErr)

			// ACT
			decls, err := converter.MethodStructs(iSpec, fn)

			// ASSERT
			if err != expectedErr {
				t.Errorf("got %#v, wanted %#v", err, expectedErr)
			}
			if decls != nil {
				t.Errorf("got %#v, wanted nil", decls)
			}
			afterEach()
		})
	})

	t.Run("NewFunc", func(t *testing.T) {
		t.Run("creates a new moq function for an interface", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)

			// ACT
			decl := converter.NewFunc(iSpec)

			// ASSERT
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted > 0", len(decl.Decs.Start))
			}
			expectedStart := "// NewMoqPublicInterface creates a new moq of the PublicInterface type"
			if decl.Decs.Start[0] != expectedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
			}
			afterEach()
		})

		t.Run("creates a new moq function for a function", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)

			// ACT
			decl := converter.NewFunc(fnSpec)

			// ASSERT
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted > 0", len(decl.Decs.Start))
			}
			expectedStart := "// NewMoqPublicFunction creates a new moq of the PublicFunction type"
			if decl.Decs.Start[0] != expectedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
			}
			afterEach()
		})
	})

	t.Run("IsolationAccessor", func(t *testing.T) {
		t.Run("creates a func", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)

			// ACT
			decl := converter.IsolationAccessor("MyInterface", "recorder", "onCall")

			// ASSERT
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted > 0", len(decl.Decs.Start))
			}
			expectedStart := "// OnCall returns the recorder implementation of the MyInterface type"
			if decl.Decs.Start[0] != expectedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
			}
			afterEach()
		})

		t.Run("creates a closure function for a function type", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)

			// ACT
			decl := converter.FuncClosure(
				"MyFn", "github.com/myshkin5/moqueries/generator", fnSpecFuncs[0])

			// ASSERT
			if len(decl.Decs.Start) < 1 {
				t.Errorf("got len %d, wanted > 0", len(decl.Decs.Start))
			}
			expectedStart := "// Mock returns the moq implementation of the MyFn type"
			if decl.Decs.Start[0] != expectedStart {
				t.Errorf("got %s, wanted %s", decl.Decs.Start[0], expectedStart)
			}
			afterEach()
		})
	})
}
