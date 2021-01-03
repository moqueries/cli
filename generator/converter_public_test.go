package generator_test

import (
	"github.com/dave/dst"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/generator"
	"github.com/myshkin5/moqueries/moq"
)

var _ = Describe("Converter", func() {
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

	BeforeEach(func() {
		scene = moq.NewScene(GinkgoT())
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
	})

	Describe("BaseStruct", func() {
		It("creates a base moq for an interface", func() {
			// ASSEMBLE

			// ACT
			decl := converter.BaseStruct(iSpec, iSpecFuncs)

			// ASSERT
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// MoqPublicInterface holds the state of a moq of the PublicInterface type"))
		})

		It("creates a base moq for a function", func() {
			// ASSEMBLE

			// ACT
			decl := converter.BaseStruct(fnSpec, fnSpecFuncs)

			// ASSERT
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// MoqPublicFunction holds the state of a moq of the PublicFunction type"))
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
				"// MoqMyInterface_mock isolates the mock interface of the MyInterface type"))
		})
	})

	Describe("MethodStructs", func() {
		It("creates structs for a function", func() {
			// ASSEMBLE
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
			Expect(err).NotTo(HaveOccurred())
			Expect(decls).To(HaveLen(7))
			decl, ok := decls[0].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// MoqPublicInterface_Func1_params holds the params of the PublicInterface type"))

			decl, ok = decls[1].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// MoqPublicInterface_Func1_paramsKey holds the map key params of the PublicInterface type"))

			decl, ok = decls[2].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal("// MoqPublicInterface_Func1_resultsByParams " +
				"contains the results for a given set of parameters for the PublicInterface type"))

			decl, ok = decls[3].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal("// MoqPublicInterface_Func1_doFn " +
				"defines the type of function needed when calling AndDo for the PublicInterface type"))

			decl, ok = decls[4].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal("// MoqPublicInterface_Func1_doReturnFn " +
				"defines the type of function needed when calling DoReturnResults for the PublicInterface type"))

			decl, ok = decls[5].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal("// MoqPublicInterface_Func1_results " +
				"holds the results of the PublicInterface type"))

			decl, ok = decls[6].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal("// MoqPublicInterface_Func1_fnRecorder " +
				"routes recorded function calls to the MoqPublicInterface moq"))
		})
	})

	Describe("NewFunc", func() {
		It("creates a new moq function for an interface", func() {
			// ASSEMBLE

			// ACT
			decl := converter.NewFunc(iSpec)

			// ASSERT
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// NewMoqPublicInterface creates a new moq of the PublicInterface type"))
		})

		It("creates a new moq function for a function", func() {
			// ASSEMBLE

			// ACT
			decl := converter.NewFunc(fnSpec)

			// ASSERT
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// NewMoqPublicFunction creates a new moq of the PublicFunction type"))
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
				"// OnCall returns the recorder implementation of the MyInterface type"))
		})
	})

	Describe("FuncClosure", func() {
		It("creates a closure function for a function type", func() {
			// ASSEMBLE

			// ACT
			decl := converter.FuncClosure(
				"MyFn", "github.com/myshkin5/moqueries/generator", fnSpecFuncs[0])

			// ASSERT
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// Mock returns the moq implementation of the MyFn type"))
		})
	})
})
