package generator_test

import (
	"github.com/dave/dst"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/generator"
)

var _ = Describe("Converter", func() {
	var (
		converter *generator.Converter

		func1Params  *dst.FieldList
		func1Results *dst.FieldList

		iSpec      *dst.TypeSpec
		iSpecFuncs []generator.Func

		fnSpec      *dst.TypeSpec
		fnSpecFuncs []generator.Func
	)

	BeforeEach(func() {
		converter = generator.NewConverter(false)

		func1Params = &dst.FieldList{
			List: []*dst.Field{
				{
					Names: []*dst.Ident{dst.NewIdent("firstParam")},
					Type: &dst.StarExpr{
						X: &dst.SelectorExpr{
							X:   dst.NewIdent("cobra"),
							Sel: dst.NewIdent("Command"),
						},
					},
				},
				{
					Names: []*dst.Ident{dst.NewIdent("secondParam")},
					Type:  dst.NewIdent("string"),
				},
				{
					Names: []*dst.Ident{dst.NewIdent("thirdParam")},
					Type: &dst.StarExpr{
						X: &dst.SelectorExpr{
							X:   dst.NewIdent("dst"),
							Sel: dst.NewIdent("InterfaceType"),
						},
					},
				},
				{
					Names: []*dst.Ident{dst.NewIdent("fourthParam")},
					Type: &dst.MapType{
						Key:   dst.NewIdent("string"),
						Value: dst.NewIdent("string"),
					},
				},
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
	})

	Describe("BaseStruct", func() {
		It("creates a base mock for an interface", func() {
			// ASSEMBLE

			// ACT
			decl := converter.BaseStruct(iSpec, iSpecFuncs)

			// ASSERT
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// mockPublicInterface holds the state of a mock of the PublicInterface type"))
		})

		It("creates a base mock for a function", func() {
			// ASSEMBLE

			// ACT
			decl := converter.BaseStruct(fnSpec, fnSpecFuncs)

			// ASSERT
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// mockPublicFunction holds the state of a mock of the PublicFunction type"))
		})
	})

	Describe("IsolationStruct", func() {
		It("creates a struct", func() {
			// ASSEMBLE

			// ACT
			decl := converter.IsolationStruct("MyInterface", "mock")

			// ASSERT
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// mockMyInterface_mock isolates the mock interface of the MyInterface type"))
		})
	})

	Describe("MethodStructs", func() {
		It("creates structs for a function", func() {
			// ASSEMBLE
			expectedParams := dst.Clone(func1Params).(*dst.FieldList)
			// Map params are represented as a deep hash when the struct is comparable
			expectedParams.List[3].Type = &dst.Ident{
				Path: "github.com/myshkin5/moqueries/hash",
				Name: "Hash",
			}
			fn := generator.Func{
				Name:    "Func1",
				Params:  func1Params,
				Results: func1Results,
			}

			// ACT
			decls := converter.MethodStructs(iSpec, fn)

			// ASSERT
			Expect(decls).To(HaveLen(7))
			decl, ok := decls[0].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// mockPublicInterface_Func1_params holds the params of the PublicInterface type"))

			decl, ok = decls[1].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// mockPublicInterface_Func1_paramsKey holds the map key params of the PublicInterface type"))

			decl, ok = decls[2].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal("// mockPublicInterface_Func1_resultsByParams " +
				"contains the results for a given set of parameters for the PublicInterface type"))

			decl, ok = decls[3].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal("// mockPublicInterface_Func1_doFn " +
				"defines the type of function needed when calling andDo for the PublicInterface type"))

			decl, ok = decls[4].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal("// mockPublicInterface_Func1_doReturnFn " +
				"defines the type of function needed when calling doReturnResults for the PublicInterface type"))

			decl, ok = decls[5].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal("// mockPublicInterface_Func1_results " +
				"holds the results of the PublicInterface type"))

			decl, ok = decls[6].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal("// mockPublicInterface_Func1_fnRecorder " +
				"routes recorded function calls to the mockPublicInterface mock"))
		})
	})

	Describe("NewFunc", func() {
		It("creates a new mock function for an interface", func() {
			// ASSEMBLE

			// ACT
			decl := converter.NewFunc(iSpec)

			// ASSERT
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// newMockPublicInterface creates a new mock of the PublicInterface type"))
		})

		It("creates a new mock function for a function", func() {
			// ASSEMBLE

			// ACT
			decl := converter.NewFunc(fnSpec)

			// ASSERT
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// newMockPublicFunction creates a new mock of the PublicFunction type"))
		})
	})

	Describe("IsolationAccessor", func() {
		It("creates a func", func() {
			// ASSEMBLE

			// ACT
			decl := converter.IsolationAccessor("MyInterface", "recorder", "onCall")

			// ASSERT
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// onCall returns the recorder implementation of the MyInterface type"))
		})
	})

	Describe("FuncClosure", func() {
		It("creates a closure function for a function type", func() {
			// ASSEMBLE

			// ACT
			decl := converter.FuncClosure("MyFn", "github.com/myshkin5/moqueries/generator", fnSpecFuncs[0])

			// ASSERT
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// mock returns the mock implementation of the MyFn type"))
		})
	})
})
