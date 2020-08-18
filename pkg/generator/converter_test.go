package generator_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/pkg/generator"
)

var _ = Describe("Converter", func() {
	var (
		converter *generator.Converter

		func1Params  *dst.FieldList
		func1Results *dst.FieldList

		func1ParamsPassthrough *dst.CompositeLit

		iSpec      *dst.TypeSpec
		iSpecFuncs []generator.Func

		fnSpec      *dst.TypeSpec
		fnSpecFuncs []generator.Func
	)

	BeforeEach(func() {
		converter = generator.NewConverter()

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

		func1ParamsPassthrough = &dst.CompositeLit{
			Type: dst.NewIdent("mockMyInterface_Func1_params"),
			Elts: []dst.Expr{
				&dst.KeyValueExpr{
					Key:   dst.NewIdent("firstParam"),
					Value: dst.NewIdent("firstParam"),
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{
							// first element has a new line before and after
							Before: dst.NewLine,
							After:  dst.NewLine,
						},
					},
				},
				&dst.KeyValueExpr{
					Key:   dst.NewIdent("secondParam"),
					Value: dst.NewIdent("secondParam"),
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{After: dst.NewLine},
					},
				},
				&dst.KeyValueExpr{
					Key:   dst.NewIdent("thirdParam"),
					Value: dst.NewIdent("thirdParam"),
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{After: dst.NewLine},
					},
				},
				&dst.KeyValueExpr{
					Key: dst.NewIdent("fourthParam"),
					Value: &dst.CallExpr{
						Fun: &dst.Ident{
							Path: "github.com/myshkin5/moqueries/pkg/hash",
							Name: "DeepHash",
						},
						Args: []dst.Expr{dst.NewIdent("fourthParam")},
					},
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{After: dst.NewLine},
					},
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
			Expect(decl.Tok).To(Equal(token.TYPE))
			Expect(decl.Specs).To(Equal([]dst.Spec{
				&dst.TypeSpec{
					Name: dst.NewIdent("mockPublicInterface"),
					Type: &dst.StructType{Fields: &dst.FieldList{
						List: []*dst.Field{
							{
								Names: []*dst.Ident{dst.NewIdent("resultsByParams_Func1")},
								Type: &dst.MapType{
									Key:   dst.NewIdent("mockPublicInterface_Func1_params"),
									Value: dst.NewIdent("mockPublicInterface_Func1_results"),
								},
							},
							{
								Names: []*dst.Ident{dst.NewIdent("params_Func1")},
								Type: &dst.ChanType{
									Dir:   dst.SEND | dst.RECV,
									Value: dst.NewIdent("mockPublicInterface_Func1_params"),
								},
							},
							{
								Names: []*dst.Ident{dst.NewIdent("resultsByParams_func2")},
								Type: &dst.MapType{
									Key:   dst.NewIdent("mockPublicInterface_func2_params"),
									Value: dst.NewIdent("mockPublicInterface_func2_results"),
								},
							},
							{
								Names: []*dst.Ident{dst.NewIdent("params_func2")},
								Type: &dst.ChanType{
									Dir:   dst.SEND | dst.RECV,
									Value: dst.NewIdent("mockPublicInterface_func2_params"),
								},
							},
						},
					}},
				},
			}))
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// mockPublicInterface holds the state of a mock of the PublicInterface type"))
		})

		It("creates a base mock for a function", func() {
			// ASSEMBLE

			// ACT
			decl := converter.BaseStruct(fnSpec, fnSpecFuncs)

			// ASSERT
			Expect(decl.Tok).To(Equal(token.TYPE))
			Expect(decl.Specs).To(Equal([]dst.Spec{
				&dst.TypeSpec{
					Name: dst.NewIdent("mockPublicFunction"),
					Type: &dst.StructType{Fields: &dst.FieldList{
						List: []*dst.Field{
							{
								Names: []*dst.Ident{dst.NewIdent("resultsByParams")},
								Type: &dst.MapType{
									Key:   dst.NewIdent("mockPublicFunction_params"),
									Value: dst.NewIdent("mockPublicFunction_results"),
								},
							},
							{
								Names: []*dst.Ident{dst.NewIdent("params")},
								Type: &dst.ChanType{
									Dir:   dst.SEND | dst.RECV,
									Value: dst.NewIdent("mockPublicFunction_params"),
								},
							},
						},
					}},
				},
			}))
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
			Expect(decl.Tok).To(Equal(token.TYPE))
			Expect(decl.Specs).To(Equal([]dst.Spec{
				&dst.TypeSpec{
					Name: dst.NewIdent("mockMyInterface_mock"),
					Type: &dst.StructType{Fields: &dst.FieldList{
						List: []*dst.Field{
							{
								Names: []*dst.Ident{dst.NewIdent("mock")},
								Type: &dst.StarExpr{
									X: dst.NewIdent("mockMyInterface"),
								},
							},
						},
					}},
				},
			}))
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
				Path: "github.com/myshkin5/moqueries/pkg/hash",
				Name: "Hash",
			}
			fn := generator.Func{
				Name:    "Func1",
				Params:  func1Params,
				Results: func1Results,
			}

			// ACT
			decls := converter.MethodStructs("MyInterface", "mockMyInterface_Func1", fn)

			// ASSERT
			Expect(decls).To(HaveLen(3))
			decl, ok := decls[0].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(decl.Tok).To(Equal(token.TYPE))
			Expect(decl.Specs).To(Equal([]dst.Spec{
				&dst.TypeSpec{
					Name: dst.NewIdent("mockMyInterface_Func1_params"),
					Type: &dst.StructType{Fields: expectedParams},
				},
			}))
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// mockMyInterface_Func1_params holds the params of the MyInterface type"))

			decl, ok = decls[1].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(decl.Tok).To(Equal(token.TYPE))
			Expect(decl.Specs).To(Equal([]dst.Spec{
				&dst.TypeSpec{
					Name: dst.NewIdent("mockMyInterface_Func1_results"),
					Type: &dst.StructType{Fields: func1Results},
				},
			}))
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// mockMyInterface_Func1_results holds the results of the MyInterface type"))

			decl, ok = decls[2].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(decl.Tok).To(Equal(token.TYPE))
			Expect(decl.Specs).To(Equal([]dst.Spec{
				&dst.TypeSpec{
					Name: dst.NewIdent("mockMyInterface_Func1_fnRecorder"),
					Type: &dst.StructType{Fields: &dst.FieldList{List: []*dst.Field{
						{
							Names: []*dst.Ident{dst.NewIdent("params")},
							Type:  dst.NewIdent("mockMyInterface_Func1_params"),
						},
						{
							Names: []*dst.Ident{dst.NewIdent("mock")},
							Type:  &dst.StarExpr{X: dst.NewIdent("mockMyInterface")},
						},
					}}},
				},
			}))
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// mockMyInterface_Func1_fnRecorder routes recorded function calls to the mockMyInterface mock"))
		})

		It("doesn't convert non-comparables to hashes when making a non-comparable struct", func() {
			// ASSEMBLE
			fn := generator.Func{
				Name: "Func1",
				// Params are comparable (used in a map key) so converted to a hash
				Params: func1Params,
				// Same list passed as a result does not get converted to a hash
				Results: func1Params,
			}

			expectedParams := dst.Clone(func1Params).(*dst.FieldList)
			// Map params are represented as a deep hash
			expectedParams.List[3].Type = &dst.Ident{
				Path: "github.com/myshkin5/moqueries/pkg/hash",
				Name: "Hash",
			}

			// ACT
			decls := converter.MethodStructs("MyInterface", "mockMyInterface_Func1", fn)

			// ASSERT
			Expect(len(decls)).To(BeNumerically(">", 2))
			decl, ok := decls[0].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(decl.Specs).To(Equal([]dst.Spec{
				&dst.TypeSpec{
					Name: dst.NewIdent("mockMyInterface_Func1_params"),
					Type: &dst.StructType{Fields: expectedParams},
				},
			}))
			decl, ok = decls[1].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(decl.Specs).To(Equal([]dst.Spec{
				&dst.TypeSpec{
					Name: dst.NewIdent("mockMyInterface_Func1_results"),
					// Types still match as nothing was converted to a hash
					Type: &dst.StructType{Fields: func1Params},
				},
			}))
		})

		It("adds names to unnamed parameters and results", func() {
			// ASSEMBLE
			for _, f := range func1Params.List {
				f.Names = nil
			}
			expectedParams := dst.Clone(func1Params).(*dst.FieldList)
			for n, f := range expectedParams.List {
				f.Names = []*dst.Ident{dst.NewIdent(fmt.Sprintf("param%d", n+1))}
			}
			// Map params are represented as a deep hash
			expectedParams.List[3].Type = &dst.Ident{
				Path: "github.com/myshkin5/moqueries/pkg/hash",
				Name: "Hash",
			}

			for _, f := range func1Results.List {
				f.Names = nil
			}
			expectedResults := dst.Clone(func1Results).(*dst.FieldList)
			for n, f := range expectedResults.List {
				f.Names = []*dst.Ident{dst.NewIdent(fmt.Sprintf("result%d", n+1))}
			}

			fn := generator.Func{
				Name:    "Func1",
				Params:  func1Params,
				Results: func1Results,
			}

			// ACT
			decls := converter.MethodStructs("MyInterface", "mockMyInterface_Func1", fn)

			// ASSERT
			Expect(len(decls)).To(BeNumerically(">", 2))
			decl, ok := decls[0].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(decl.Specs).To(Equal([]dst.Spec{
				&dst.TypeSpec{
					Name: dst.NewIdent("mockMyInterface_Func1_params"),
					Type: &dst.StructType{Fields: expectedParams},
				},
			}))
			// Verify the source AST wasn't modified
			Expect(func1Params.List[0].Names).To(BeNil())

			decl, ok = decls[1].(*dst.GenDecl)
			Expect(ok).To(BeTrue())
			Expect(decl.Specs).To(Equal([]dst.Spec{
				&dst.TypeSpec{
					Name: dst.NewIdent("mockMyInterface_Func1_results"),
					Type: &dst.StructType{Fields: expectedResults},
				},
			}))
			// Verify the source AST wasn't modified
			Expect(func1Results.List[0].Names).To(BeNil())
		})
	})

	Describe("NewFunc", func() {
		It("creates a new mock function for an interface", func() {
			// ASSEMBLE

			// ACT
			decl := converter.NewFunc(iSpec, iSpecFuncs)

			// ASSERT
			Expect(decl.Name).To(Equal(dst.NewIdent("newMockPublicInterface")))
			Expect(decl.Type).To(Equal(&dst.FuncType{
				Params: &dst.FieldList{},
				Results: &dst.FieldList{
					List: []*dst.Field{
						{
							Type: &dst.StarExpr{X: dst.NewIdent("mockPublicInterface")},
						},
					},
				},
			}))
			Expect(decl.Body.List).To(HaveLen(1))
			returnStmt, ok := decl.Body.List[0].(*dst.ReturnStmt)
			Expect(ok).To(BeTrue())
			Expect(returnStmt.Results).To(HaveLen(1))
			unaryExpr, ok := returnStmt.Results[0].(*dst.UnaryExpr)
			Expect(ok).To(BeTrue())
			Expect(unaryExpr.Op).To(Equal(token.AND))
			compositeLit, ok := unaryExpr.X.(*dst.CompositeLit)
			Expect(ok).To(BeTrue())
			Expect(compositeLit.Type).To(Equal(dst.NewIdent("mockPublicInterface")))
			Expect(compositeLit.Elts).To(HaveLen(4))
			Expect(compositeLit.Elts[0]).To(Equal(&dst.KeyValueExpr{
				Key: dst.NewIdent("resultsByParams_Func1"),
				Value: &dst.CompositeLit{
					Type: &dst.MapType{
						Key:   dst.NewIdent("mockPublicInterface_Func1_params"),
						Value: dst.NewIdent("mockPublicInterface_Func1_results"),
					},
				},
				Decs: dst.KeyValueExprDecorations{
					NodeDecs: dst.NodeDecs{After: dst.NewLine},
				},
			}))
			Expect(compositeLit.Elts[1]).To(Equal(&dst.KeyValueExpr{
				Key: dst.NewIdent("params_Func1"),
				Value: &dst.CallExpr{
					Fun: dst.NewIdent("make"),
					Args: []dst.Expr{
						&dst.ChanType{
							Dir:   dst.SEND | dst.RECV,
							Value: dst.NewIdent("mockPublicInterface_Func1_params"),
						},
						&dst.BasicLit{
							Kind:  token.INT,
							Value: "100",
						},
					},
				},
				Decs: dst.KeyValueExprDecorations{
					NodeDecs: dst.NodeDecs{After: dst.NewLine},
				},
			}))
			Expect(compositeLit.Elts[2]).To(Equal(&dst.KeyValueExpr{
				Key: dst.NewIdent("resultsByParams_func2"),
				Value: &dst.CompositeLit{
					Type: &dst.MapType{
						Key:   dst.NewIdent("mockPublicInterface_func2_params"),
						Value: dst.NewIdent("mockPublicInterface_func2_results"),
					},
				},
				Decs: dst.KeyValueExprDecorations{
					NodeDecs: dst.NodeDecs{After: dst.NewLine},
				},
			}))
			Expect(compositeLit.Elts[3]).To(Equal(&dst.KeyValueExpr{
				Key: dst.NewIdent("params_func2"),
				Value: &dst.CallExpr{
					Fun: dst.NewIdent("make"),
					Args: []dst.Expr{
						&dst.ChanType{
							Dir:   dst.SEND | dst.RECV,
							Value: dst.NewIdent("mockPublicInterface_func2_params"),
						},
						&dst.BasicLit{
							Kind:  token.INT,
							Value: "100",
						},
					},
				},
				Decs: dst.KeyValueExprDecorations{
					NodeDecs: dst.NodeDecs{After: dst.NewLine},
				},
			}))
			Expect(compositeLit.Decs).To(Equal(dst.CompositeLitDecorations{Lbrace: []string{"\n"}}))
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// newMockPublicInterface creates a new mock of the PublicInterface type"))
		})

		It("creates a new mock function for a function", func() {
			// ASSEMBLE

			// ACT
			decl := converter.NewFunc(fnSpec, fnSpecFuncs)

			// ASSERT
			Expect(decl.Name).To(Equal(dst.NewIdent("newMockPublicFunction")))
			Expect(decl.Type).To(Equal(&dst.FuncType{
				Params: &dst.FieldList{},
				Results: &dst.FieldList{
					List: []*dst.Field{
						{
							Type: &dst.StarExpr{X: dst.NewIdent("mockPublicFunction")},
						},
					},
				},
			}))
			Expect(decl.Body.List).To(HaveLen(1))
			returnStmt, ok := decl.Body.List[0].(*dst.ReturnStmt)
			Expect(ok).To(BeTrue())
			Expect(returnStmt.Results).To(HaveLen(1))
			unaryExpr, ok := returnStmt.Results[0].(*dst.UnaryExpr)
			Expect(ok).To(BeTrue())
			Expect(unaryExpr.Op).To(Equal(token.AND))
			compositeLit, ok := unaryExpr.X.(*dst.CompositeLit)
			Expect(ok).To(BeTrue())
			Expect(compositeLit.Type).To(Equal(dst.NewIdent("mockPublicFunction")))
			Expect(compositeLit.Elts).To(HaveLen(2))
			Expect(compositeLit.Elts[0]).To(Equal(&dst.KeyValueExpr{
				Key: dst.NewIdent("resultsByParams"),
				Value: &dst.CompositeLit{
					Type: &dst.MapType{
						Key:   dst.NewIdent("mockPublicFunction_params"),
						Value: dst.NewIdent("mockPublicFunction_results"),
					},
				},
				Decs: dst.KeyValueExprDecorations{
					NodeDecs: dst.NodeDecs{After: dst.NewLine},
				},
			}))
			Expect(compositeLit.Elts[1]).To(Equal(&dst.KeyValueExpr{
				Key: dst.NewIdent("params"),
				Value: &dst.CallExpr{
					Fun: dst.NewIdent("make"),
					Args: []dst.Expr{
						&dst.ChanType{
							Dir:   dst.SEND | dst.RECV,
							Value: dst.NewIdent("mockPublicFunction_params"),
						},
						&dst.BasicLit{
							Kind:  token.INT,
							Value: "100",
						},
					},
				},
				Decs: dst.KeyValueExprDecorations{
					NodeDecs: dst.NodeDecs{After: dst.NewLine},
				},
			}))
			Expect(compositeLit.Decs).To(Equal(dst.CompositeLitDecorations{Lbrace: []string{"\n"}}))
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
			Expect(decl.Recv).To(Equal(&dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("m")},
						Type: &dst.StarExpr{
							X: dst.NewIdent("mockMyInterface"),
						},
					},
				},
			}))
			Expect(decl.Name).To(Equal(dst.NewIdent("onCall")))
			Expect(decl.Type).To(Equal(&dst.FuncType{
				Params: &dst.FieldList{},
				Results: &dst.FieldList{
					List: []*dst.Field{
						{
							Type: &dst.StarExpr{X: dst.NewIdent("mockMyInterface_recorder")},
						},
					},
				},
			}))
			Expect(decl.Body.List).To(HaveLen(1))
			returnStmt, ok := decl.Body.List[0].(*dst.ReturnStmt)
			Expect(ok).To(BeTrue())
			Expect(returnStmt.Results).To(HaveLen(1))
			unaryExpr, ok := returnStmt.Results[0].(*dst.UnaryExpr)
			Expect(ok).To(BeTrue())
			Expect(unaryExpr.Op).To(Equal(token.AND))
			compositeLit, ok := unaryExpr.X.(*dst.CompositeLit)
			Expect(ok).To(BeTrue())
			Expect(compositeLit.Type).To(Equal(dst.NewIdent("mockMyInterface_recorder")))
			Expect(compositeLit.Elts).To(HaveLen(1))
			Expect(compositeLit.Elts[0]).To(Equal(&dst.KeyValueExpr{
				Key:   dst.NewIdent("mock"),
				Value: dst.NewIdent("m"),
				Decs: dst.KeyValueExprDecorations{
					NodeDecs: dst.NodeDecs{After: dst.NewLine},
				},
			}))
			Expect(compositeLit.Decs).To(Equal(dst.CompositeLitDecorations{Lbrace: []string{"\n"}}))
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// onCall returns the recorder implementation of the MyInterface type"))
		})
	})

	Describe("MockMethod", func() {
		It("creates a mock function for a method of the interface", func() {
			// ASSEMBLE

			// ACT
			decl := converter.MockMethod("MyInterface", iSpecFuncs[0])

			// ASSERT
			Expect(decl.Recv).To(Equal(&dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("m")},
						Type: &dst.StarExpr{
							X: dst.NewIdent("mockMyInterface_mock"),
						},
					},
				},
			}))
			Expect(decl.Name).To(Equal(dst.NewIdent("Func1")))
			Expect(decl.Type).To(Equal(&dst.FuncType{
				Params:  func1Params,
				Results: func1Results,
			}))
			Expect(decl.Type.Params.List[0]).NotTo(BeIdenticalTo(func1Params.List[0]), "should be cloned")
			Expect(decl.Type.Results.List[0]).NotTo(BeIdenticalTo(func1Results.List[0]), "should be cloned")
			Expect(decl.Body.List).To(HaveLen(5))
			Expect(decl.Body.List[0]).To(Equal(&dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent("params")},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{func1ParamsPassthrough},
			}))
			Expect(decl.Body.List[1]).To(Equal(&dst.SendStmt{
				Chan: &dst.SelectorExpr{
					X: &dst.SelectorExpr{
						X:   dst.NewIdent("m"),
						Sel: dst.NewIdent("mock"),
					},
					Sel: dst.NewIdent("params_Func1"),
				},
				Value: dst.NewIdent("params"),
			}))
			Expect(decl.Body.List[2]).To(Equal(&dst.AssignStmt{
				Lhs: []dst.Expr{
					dst.NewIdent("results"),
					dst.NewIdent("ok"),
				},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{
					&dst.IndexExpr{
						X: &dst.SelectorExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("m"),
								Sel: dst.NewIdent("mock"),
							},
							Sel: dst.NewIdent("resultsByParams_Func1"),
						},
						Index: dst.NewIdent("params"),
					},
				},
			}))
			Expect(decl.Body.List[3]).To(Equal(&dst.IfStmt{
				Cond: dst.NewIdent("ok"),
				Body: &dst.BlockStmt{
					List: []dst.Stmt{
						&dst.AssignStmt{
							Lhs: []dst.Expr{dst.NewIdent("result")},
							Tok: token.ASSIGN,
							Rhs: []dst.Expr{
								&dst.SelectorExpr{
									X:   dst.NewIdent("results"),
									Sel: dst.NewIdent("result"),
								},
							},
						},
						&dst.AssignStmt{
							Lhs: []dst.Expr{dst.NewIdent("err")},
							Tok: token.ASSIGN,
							Rhs: []dst.Expr{
								&dst.SelectorExpr{
									X:   dst.NewIdent("results"),
									Sel: dst.NewIdent("err"),
								},
							},
						},
					},
				},
			}))
			Expect(decl.Body.List[4]).To(Equal(&dst.ReturnStmt{
				Results: []dst.Expr{
					dst.NewIdent("result"),
					dst.NewIdent("err"),
				},
			}))
			Expect(decl.Decs).To(Equal(dst.FuncDeclDecorations{
				NodeDecs: dst.NodeDecs{
					Before: dst.EmptyLine,
					After:  dst.EmptyLine,
				},
			}))
		})

		It("creates a mock function for a function", func() {
			// ASSEMBLE
			func1ParamsPassthrough.Type = dst.NewIdent("mockMyFn_params")

			// ACT
			decl := converter.MockMethod("MyFn", fnSpecFuncs[0])

			// ASSERT
			Expect(decl.Recv).To(Equal(&dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("m")},
						Type: &dst.StarExpr{
							X: dst.NewIdent("mockMyFn_mock"),
						},
					},
				},
			}))
			Expect(decl.Name).To(Equal(dst.NewIdent("fn")))
			Expect(decl.Type).To(Equal(&dst.FuncType{
				Params:  func1Params,
				Results: func1Results,
			}))
			Expect(decl.Type.Params.List[0]).NotTo(BeIdenticalTo(func1Params.List[0]), "should be cloned")
			Expect(decl.Type.Results.List[0]).NotTo(BeIdenticalTo(func1Results.List[0]), "should be cloned")
			Expect(decl.Body.List).To(HaveLen(5))
			Expect(decl.Body.List[0]).To(Equal(&dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent("params")},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{func1ParamsPassthrough},
			}))
			Expect(decl.Body.List[1]).To(Equal(&dst.SendStmt{
				Chan: &dst.SelectorExpr{
					X: &dst.SelectorExpr{
						X:   dst.NewIdent("m"),
						Sel: dst.NewIdent("mock"),
					},
					Sel: dst.NewIdent("params"),
				},
				Value: dst.NewIdent("params"),
			}))
			Expect(decl.Body.List[2]).To(Equal(&dst.AssignStmt{
				Lhs: []dst.Expr{
					dst.NewIdent("results"),
					dst.NewIdent("ok"),
				},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{
					&dst.IndexExpr{
						X: &dst.SelectorExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("m"),
								Sel: dst.NewIdent("mock"),
							},
							Sel: dst.NewIdent("resultsByParams"),
						},
						Index: dst.NewIdent("params"),
					},
				},
			}))
			Expect(decl.Body.List[3]).To(Equal(&dst.IfStmt{
				Cond: dst.NewIdent("ok"),
				Body: &dst.BlockStmt{
					List: []dst.Stmt{
						&dst.AssignStmt{
							Lhs: []dst.Expr{dst.NewIdent("result")},
							Tok: token.ASSIGN,
							Rhs: []dst.Expr{
								&dst.SelectorExpr{
									X:   dst.NewIdent("results"),
									Sel: dst.NewIdent("result"),
								},
							},
						},
						&dst.AssignStmt{
							Lhs: []dst.Expr{dst.NewIdent("err")},
							Tok: token.ASSIGN,
							Rhs: []dst.Expr{
								&dst.SelectorExpr{
									X:   dst.NewIdent("results"),
									Sel: dst.NewIdent("err"),
								},
							},
						},
					},
				},
			}))
			Expect(decl.Body.List[4]).To(Equal(&dst.ReturnStmt{
				Results: []dst.Expr{
					dst.NewIdent("result"),
					dst.NewIdent("err"),
				},
			}))
			Expect(decl.Decs).To(Equal(dst.FuncDeclDecorations{
				NodeDecs: dst.NodeDecs{
					Before: dst.EmptyLine,
					After:  dst.EmptyLine,
				},
			}))
		})

		It("copies all params when multiple params all use the same type", func() {
			// ASSEMBLE
			func1Params.List = func1Params.List[0:1]
			func1Params.List[0].Names = append(func1Params.List[0].Names, dst.NewIdent("secondParam"))

			// ACT
			decl := converter.MockMethod("MyInterface", iSpecFuncs[0])

			// ASSERT
			elements := decl.Body.List[0].(*dst.AssignStmt).Rhs[0].(*dst.CompositeLit).Elts
			Expect(elements).To(Equal([]dst.Expr{
				&dst.KeyValueExpr{
					Key:   dst.NewIdent("firstParam"),
					Value: dst.NewIdent("firstParam"),
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{
							// first element has a new line before and after
							Before: dst.NewLine,
							After:  dst.NewLine,
						},
					},
				},
				&dst.KeyValueExpr{
					Key:   dst.NewIdent("secondParam"),
					Value: dst.NewIdent("secondParam"),
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{
							After: dst.NewLine,
						},
					},
				},
			}))
		})

		It("copies unnamed params", func() {
			// ASSEMBLE
			for _, f := range func1Params.List {
				f.Names = nil
			}

			// ACT
			decl := converter.MockMethod("MyInterface", iSpecFuncs[0])

			// ASSERT
			elements := decl.Body.List[0].(*dst.AssignStmt).Rhs[0].(*dst.CompositeLit).Elts
			Expect(elements).To(Equal([]dst.Expr{
				&dst.KeyValueExpr{
					Key:   dst.NewIdent("param1"),
					Value: dst.NewIdent("param1"),
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{
							// first element has a new line before and after
							Before: dst.NewLine,
							After:  dst.NewLine,
						},
					},
				},
				&dst.KeyValueExpr{
					Key:   dst.NewIdent("param2"),
					Value: dst.NewIdent("param2"),
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{After: dst.NewLine},
					},
				},
				&dst.KeyValueExpr{
					Key:   dst.NewIdent("param3"),
					Value: dst.NewIdent("param3"),
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{After: dst.NewLine},
					},
				},
				&dst.KeyValueExpr{
					Key: dst.NewIdent("param4"),
					Value: &dst.CallExpr{
						Fun: &dst.Ident{
							Path: "github.com/myshkin5/moqueries/pkg/hash",
							Name: "DeepHash",
						},
						Args: []dst.Expr{dst.NewIdent("param4")},
					},
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{After: dst.NewLine},
					},
				},
			}))
		})

		It("copies all results when multiple results all use the same type", func() {
			// ASSEMBLE
			func1Results.List = func1Results.List[0:1]
			func1Results.List[0].Names = append(func1Results.List[0].Names, dst.NewIdent("secondResult"))

			// ACT
			decl := converter.MockMethod("MyInterface", iSpecFuncs[0])

			// ASSERT
			resultBody := decl.Body.List[3].(*dst.IfStmt).Body.List
			Expect(resultBody).To(Equal([]dst.Stmt{
				&dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent("result")},
					Tok: token.ASSIGN,
					Rhs: []dst.Expr{
						&dst.SelectorExpr{
							X:   dst.NewIdent("results"),
							Sel: dst.NewIdent("result"),
						},
					},
				},
				&dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent("secondResult")},
					Tok: token.ASSIGN,
					Rhs: []dst.Expr{
						&dst.SelectorExpr{
							X:   dst.NewIdent("results"),
							Sel: dst.NewIdent("secondResult"),
						},
					},
				},
			}))
		})

		It("copies unnamed results", func() {
			// ASSEMBLE
			for _, f := range func1Results.List {
				f.Names = nil
			}

			// ACT
			decl := converter.MockMethod("MyInterface", iSpecFuncs[0])

			// ASSERT
			resultBody := decl.Body.List[3].(*dst.IfStmt).Body.List
			Expect(resultBody).To(Equal([]dst.Stmt{
				&dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent("result1")},
					Tok: token.ASSIGN,
					Rhs: []dst.Expr{
						&dst.SelectorExpr{
							X:   dst.NewIdent("results"),
							Sel: dst.NewIdent("result1"),
						},
					},
				},
				&dst.AssignStmt{
					Lhs: []dst.Expr{dst.NewIdent("result2")},
					Tok: token.ASSIGN,
					Rhs: []dst.Expr{
						&dst.SelectorExpr{
							X:   dst.NewIdent("results"),
							Sel: dst.NewIdent("result2"),
						},
					},
				},
			}))
		})

		It("returns all results when multiple results all use the same type", func() {
			// ASSEMBLE
			func1Results.List = func1Results.List[0:1]
			func1Results.List[0].Names = append(func1Results.List[0].Names, dst.NewIdent("secondResult"))

			// ACT
			decl := converter.MockMethod("MyInterface", iSpecFuncs[0])

			// ASSERT
			Expect(decl.Body.List[4]).To(Equal(&dst.ReturnStmt{
				Results: []dst.Expr{
					dst.NewIdent("result"),
					dst.NewIdent("secondResult"),
				},
			}))
		})

		It("returns unnamed results", func() {
			// ASSEMBLE
			for _, f := range func1Results.List {
				f.Names = nil
			}

			// ACT
			decl := converter.MockMethod("MyInterface", iSpecFuncs[0])

			// ASSERT
			Expect(decl.Body.List[4]).To(Equal(&dst.ReturnStmt{
				Results: []dst.Expr{
					dst.NewIdent("result1"),
					dst.NewIdent("result2"),
				},
			}))
		})
	})

	Describe("FuncClosure", func() {
		It("creates a closure function for a function type", func() {
			// ASSEMBLE

			// ACT
			decl := converter.FuncClosure("MyFn", "github.com/myshkin5/moqueries/pkg/generator", fnSpecFuncs[0])

			// ASSERT
			Expect(decl.Recv).To(Equal(&dst.FieldList{
				List: []*dst.Field{
					{
						Names: []*dst.Ident{dst.NewIdent("m")},
						Type:  &dst.StarExpr{X: dst.NewIdent("mockMyFn")},
					},
				},
			}))
			Expect(decl.Name).To(Equal(dst.NewIdent("mock")))
			Expect(decl.Type).To(Equal(&dst.FuncType{
				Params: &dst.FieldList{},
				Results: &dst.FieldList{
					List: []*dst.Field{
						{
							Type: &dst.Ident{
								Path: "github.com/myshkin5/moqueries/pkg/generator",
								Name: "MyFn",
							},
						},
					},
				},
			}))
			Expect(decl.Body.List).To(HaveLen(1))
			rStmt, ok := decl.Body.List[0].(*dst.ReturnStmt)
			Expect(ok).To(BeTrue())
			Expect(rStmt.Results).To(HaveLen(1))
			fLit, ok := rStmt.Results[0].(*dst.FuncLit)
			Expect(ok).To(BeTrue())
			Expect(fLit.Type.Params).To(Equal(func1Params))
			Expect(fLit.Type.Params).NotTo(BeIdenticalTo(func1Params))
			Expect(fLit.Type.Results).To(Equal(func1Results))
			Expect(fLit.Type.Results).NotTo(BeIdenticalTo(func1Results))
			Expect(fLit.Body.List).To(HaveLen(2))
			Expect(fLit.Body.List[0]).To(Equal(&dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent("mock")},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{
					&dst.UnaryExpr{
						Op: token.AND,
						X: &dst.CompositeLit{
							Type: dst.NewIdent("mockMyFn_mock"),
							Elts: []dst.Expr{
								&dst.KeyValueExpr{
									Key:   dst.NewIdent("mock"),
									Value: dst.NewIdent("m"),
								},
							},
						},
					},
				},
			}))
			Expect(fLit.Body.List[1]).To(Equal(&dst.ReturnStmt{
				Results: []dst.Expr{
					&dst.CallExpr{
						Fun: &dst.SelectorExpr{
							X:   dst.NewIdent("mock"),
							Sel: dst.NewIdent("fn"),
						},
						Args: []dst.Expr{
							dst.NewIdent("firstParam"),
							dst.NewIdent("secondParam"),
							dst.NewIdent("thirdParam"),
							dst.NewIdent("fourthParam"),
						},
					},
				},
			}))
			Expect(len(decl.Decs.Start)).To(BeNumerically(">", 0))
			Expect(decl.Decs.Start[0]).To(Equal(
				"// mock returns the mock implementation of the MyFn type"))
		})
	})

	Describe("RecorderMethods", func() {
		Context("Interfaces", func() {
			It("generates a recorder method", func() {
				// ASSEMBLE

				// ACT
				decls := converter.RecorderMethods("MyInterface", iSpecFuncs[0])

				// ASSERT
				Expect(len(decls)).To(BeNumerically(">", 0))
				decl, ok := decls[0].(*dst.FuncDecl)
				Expect(ok).To(BeTrue())
				Expect(decl.Recv).To(Equal(&dst.FieldList{
					List: []*dst.Field{
						{
							Names: []*dst.Ident{dst.NewIdent("m")},
							Type: &dst.StarExpr{
								X: dst.NewIdent("mockMyInterface_recorder"),
							},
						},
					},
				}))
				Expect(decl.Name).To(Equal(dst.NewIdent("Func1")))
				Expect(decl.Type).To(Equal(&dst.FuncType{
					Params: func1Params,
					Results: &dst.FieldList{List: []*dst.Field{{
						Type: &dst.StarExpr{X: dst.NewIdent("mockMyInterface_Func1_fnRecorder")},
					}}},
				}))
				Expect(decl.Type.Params.List[0]).NotTo(BeIdenticalTo(func1Params.List[0]), "should be cloned")
				Expect(decl.Body.List).To(HaveLen(1))
				returnStmt, ok := decl.Body.List[0].(*dst.ReturnStmt)
				Expect(ok).To(BeTrue())
				Expect(returnStmt.Results).To(HaveLen(1))
				unaryExpr, ok := returnStmt.Results[0].(*dst.UnaryExpr)
				Expect(ok).To(BeTrue())
				Expect(unaryExpr.Op).To(Equal(token.AND))
				compositeLit, ok := unaryExpr.X.(*dst.CompositeLit)
				Expect(ok).To(BeTrue())
				Expect(compositeLit.Type).To(Equal(dst.NewIdent("mockMyInterface_Func1_fnRecorder")))
				Expect(compositeLit.Elts).To(HaveLen(2))
				Expect(compositeLit.Elts[0]).To(Equal(&dst.KeyValueExpr{
					Key:   dst.NewIdent("params"),
					Value: func1ParamsPassthrough,
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{After: dst.NewLine},
					},
				}))
				Expect(compositeLit.Elts[1]).To(Equal(&dst.KeyValueExpr{
					Key: dst.NewIdent("mock"),
					Value: &dst.SelectorExpr{
						X:   dst.NewIdent("m"),
						Sel: dst.NewIdent("mock"),
					},
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{After: dst.NewLine},
					},
				}))
				Expect(compositeLit.Decs).To(Equal(dst.CompositeLitDecorations{Lbrace: []string{"\n"}}))
				Expect(decl.Decs).To(Equal(dst.FuncDeclDecorations{
					NodeDecs: dst.NodeDecs{
						Before: dst.EmptyLine,
						After:  dst.EmptyLine,
					},
				}))
			})

			It("generates a return method", func() {
				// ASSEMBLE

				// ACT
				decls := converter.RecorderMethods("MyInterface", iSpecFuncs[0])

				// ASSERT
				Expect(len(decls)).To(BeNumerically(">", 1))
				decl, ok := decls[1].(*dst.FuncDecl)
				Expect(ok).To(BeTrue())
				Expect(decl.Recv).To(Equal(&dst.FieldList{
					List: []*dst.Field{
						{
							Names: []*dst.Ident{dst.NewIdent("r")},
							Type: &dst.StarExpr{
								X: dst.NewIdent("mockMyInterface_Func1_fnRecorder"),
							},
						},
					},
				}))
				Expect(decl.Name).To(Equal(dst.NewIdent("ret")))
				Expect(decl.Type).To(Equal(&dst.FuncType{
					Params: func1Results,
				}))
				Expect(decl.Type.Params.List[0]).NotTo(BeIdenticalTo(func1Results.List[0]), "should be cloned")
				Expect(decl.Body.List).To(HaveLen(1))
				assignStmt, ok := decl.Body.List[0].(*dst.AssignStmt)
				Expect(ok).To(BeTrue())
				Expect(assignStmt.Lhs).To(Equal([]dst.Expr{
					&dst.IndexExpr{
						X: &dst.SelectorExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("r"),
								Sel: dst.NewIdent("mock"),
							},
							Sel: dst.NewIdent("resultsByParams_Func1"),
						},
						Index: &dst.SelectorExpr{
							X:   dst.NewIdent("r"),
							Sel: dst.NewIdent("params"),
						},
					},
				}))
				Expect(assignStmt.Tok).To(Equal(token.ASSIGN))
				Expect(assignStmt.Rhs).To(Equal([]dst.Expr{
					&dst.CompositeLit{
						Type: dst.NewIdent("mockMyInterface_Func1_results"),
						Elts: []dst.Expr{
							&dst.KeyValueExpr{
								Key:   dst.NewIdent("result"),
								Value: dst.NewIdent("result"),
								Decs: dst.KeyValueExprDecorations{
									NodeDecs: dst.NodeDecs{
										// first element has a new line before and after
										Before: dst.NewLine,
										After:  dst.NewLine,
									},
								},
							},
							&dst.KeyValueExpr{
								Key:   dst.NewIdent("err"),
								Value: dst.NewIdent("err"),
								Decs: dst.KeyValueExprDecorations{
									NodeDecs: dst.NodeDecs{After: dst.NewLine},
								},
							},
						},
					},
				}))
			})

			It("doesn't hash return values in the return method", func() {
				// ASSEMBLE
				func1Results.List[0] = &dst.Field{
					Names: []*dst.Ident{dst.NewIdent("headerMap")},
					Type: &dst.MapType{
						Key:   dst.NewIdent("string"),
						Value: dst.NewIdent("string"),
					},
				}

				// ACT
				decls := converter.RecorderMethods("MyInterface", iSpecFuncs[0])

				// ASSERT
				Expect(len(decls)).To(BeNumerically(">", 1))
				decl, ok := decls[1].(*dst.FuncDecl)
				Expect(ok).To(BeTrue())
				Expect(decl.Body.List).To(HaveLen(1))
				assignStmt, ok := decl.Body.List[0].(*dst.AssignStmt)
				Expect(ok).To(BeTrue())
				Expect(assignStmt.Rhs[0].(*dst.CompositeLit).Elts[0]).To(Equal(&dst.KeyValueExpr{
					Key:   dst.NewIdent("headerMap"),
					Value: dst.NewIdent("headerMap"),
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{
							// first element has a new line before and after
							Before: dst.NewLine,
							After:  dst.NewLine,
						},
					},
				}))
			})

			It("names unnamed return values in the return method", func() {
				// ASSEMBLE
				for _, f := range func1Results.List {
					f.Names = nil
				}

				// ACT
				decls := converter.RecorderMethods("MyInterface", iSpecFuncs[0])

				// ASSERT
				Expect(len(decls)).To(BeNumerically(">", 1))
				decl, ok := decls[1].(*dst.FuncDecl)
				Expect(ok).To(BeTrue())
				Expect(decl.Body.List).To(HaveLen(1))
				assignStmt, ok := decl.Body.List[0].(*dst.AssignStmt)
				Expect(ok).To(BeTrue())
				Expect(assignStmt.Rhs[0].(*dst.CompositeLit).Elts).To(Equal([]dst.Expr{
					&dst.KeyValueExpr{
						Key:   dst.NewIdent("result1"),
						Value: dst.NewIdent("result1"),
						Decs: dst.KeyValueExprDecorations{
							NodeDecs: dst.NodeDecs{
								// first element has a new line before and after
								Before: dst.NewLine,
								After:  dst.NewLine,
							},
						},
					},
					&dst.KeyValueExpr{
						Key:   dst.NewIdent("result2"),
						Value: dst.NewIdent("result2"),
						Decs: dst.KeyValueExprDecorations{
							NodeDecs: dst.NodeDecs{After: dst.NewLine},
						},
					},
				}))
			})
		})

		Context("Functions", func() {
			It("generates a recorder method", func() {
				// ASSEMBLE
				func1ParamsPassthrough.Type = dst.NewIdent("mockMyFunc_params")

				// ACT
				decls := converter.RecorderMethods("MyFunc", fnSpecFuncs[0])

				// ASSERT
				Expect(len(decls)).To(BeNumerically(">", 0))
				decl, ok := decls[0].(*dst.FuncDecl)
				Expect(ok).To(BeTrue())
				Expect(decl.Recv).To(Equal(&dst.FieldList{
					List: []*dst.Field{
						{
							Names: []*dst.Ident{dst.NewIdent("m")},
							Type: &dst.StarExpr{
								X: dst.NewIdent("mockMyFunc"),
							},
						},
					},
				}))
				Expect(decl.Name).To(Equal(dst.NewIdent("onCall")))
				Expect(decl.Type).To(Equal(&dst.FuncType{
					Params: func1Params,
					Results: &dst.FieldList{
						List: []*dst.Field{
							{
								Type: &dst.StarExpr{
									X: dst.NewIdent("mockMyFunc_fnRecorder"),
								},
							},
						},
					},
				}))
				Expect(decl.Body.List).To(HaveLen(1))
				returnStmt, ok := decl.Body.List[0].(*dst.ReturnStmt)
				Expect(ok).To(BeTrue())
				Expect(returnStmt.Results).To(HaveLen(1))
				Expect(ok).To(BeTrue())
				unaryExpr, ok := returnStmt.Results[0].(*dst.UnaryExpr)
				Expect(ok).To(BeTrue())
				Expect(unaryExpr.Op).To(Equal(token.AND))
				compositeLit, ok := unaryExpr.X.(*dst.CompositeLit)
				Expect(ok).To(BeTrue())
				Expect(compositeLit.Type).To(Equal(dst.NewIdent("mockMyFunc_fnRecorder")))
				Expect(compositeLit.Elts).To(HaveLen(2))
				Expect(compositeLit.Elts[0]).To(Equal(&dst.KeyValueExpr{
					Key:   dst.NewIdent("params"),
					Value: func1ParamsPassthrough,
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{After: dst.NewLine},
					},
				}))
				Expect(compositeLit.Elts[1]).To(Equal(&dst.KeyValueExpr{
					Key:   dst.NewIdent("mock"),
					Value: dst.NewIdent("m"),
					Decs: dst.KeyValueExprDecorations{
						NodeDecs: dst.NodeDecs{After: dst.NewLine},
					},
				}))
				Expect(compositeLit.Decs).To(Equal(dst.CompositeLitDecorations{Lbrace: []string{"\n"}}))
				Expect(decl.Decs).To(Equal(dst.FuncDeclDecorations{
					NodeDecs: dst.NodeDecs{
						Before: dst.EmptyLine,
						After:  dst.EmptyLine,
					},
				}))
			})

			It("generates a return method", func() {
				// ASSEMBLE

				// ACT
				decls := converter.RecorderMethods("MyFunc", fnSpecFuncs[0])

				// ASSERT
				Expect(len(decls)).To(BeNumerically(">", 1))
				decl, ok := decls[1].(*dst.FuncDecl)
				Expect(ok).To(BeTrue())
				Expect(decl.Recv).To(Equal(&dst.FieldList{
					List: []*dst.Field{
						{
							Names: []*dst.Ident{dst.NewIdent("r")},
							Type: &dst.StarExpr{
								X: dst.NewIdent("mockMyFunc_fnRecorder"),
							},
						},
					},
				}))
				Expect(decl.Name).To(Equal(dst.NewIdent("ret")))
				Expect(decl.Type).To(Equal(&dst.FuncType{
					Params: func1Results,
				}))
				Expect(decl.Type.Params.List[0]).NotTo(BeIdenticalTo(func1Results.List[0]), "should be cloned")
				Expect(decl.Body.List).To(HaveLen(1))
				assignStmt, ok := decl.Body.List[0].(*dst.AssignStmt)
				Expect(ok).To(BeTrue())
				Expect(assignStmt.Lhs).To(Equal([]dst.Expr{
					&dst.IndexExpr{
						X: &dst.SelectorExpr{
							X: &dst.SelectorExpr{
								X:   dst.NewIdent("r"),
								Sel: dst.NewIdent("mock"),
							},
							Sel: dst.NewIdent("resultsByParams"),
						},
						Index: &dst.SelectorExpr{
							X:   dst.NewIdent("r"),
							Sel: dst.NewIdent("params"),
						},
					},
				}))
				Expect(assignStmt.Tok).To(Equal(token.ASSIGN))
				Expect(assignStmt.Rhs).To(Equal([]dst.Expr{
					&dst.CompositeLit{
						Type: dst.NewIdent("mockMyFunc_results"),
						Elts: []dst.Expr{
							&dst.KeyValueExpr{
								Key:   dst.NewIdent("result"),
								Value: dst.NewIdent("result"),
								Decs: dst.KeyValueExprDecorations{
									NodeDecs: dst.NodeDecs{
										// first element has a new line before and after
										Before: dst.NewLine,
										After:  dst.NewLine,
									},
								},
							},
							&dst.KeyValueExpr{
								Key:   dst.NewIdent("err"),
								Value: dst.NewIdent("err"),
								Decs: dst.KeyValueExprDecorations{
									NodeDecs: dst.NodeDecs{After: dst.NewLine},
								},
							},
						},
					},
				}))
			})
		})
	})

	It("dumps the AST of an interface mock", func() {
		filePath := "./mock_converterer_test.go"
		outPath := "./mock_converterer_ast.txt"

		fSet := token.NewFileSet()
		inFile, err := parser.ParseFile(fSet, filePath, nil, parser.ParseComments)
		Expect(err).NotTo(HaveOccurred())

		outFile, err := os.Create(outPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(ast.Fprint(outFile, fSet, inFile, ast.NotNilFilter)).To(Succeed())
	})

	It("dumps the DST of an interface mock", func() {
		filePath := "./mock_converterer_test.go"
		outPath := "./mock_converterer_dst.txt"

		fSet := token.NewFileSet()
		inFile, err := parser.ParseFile(fSet, filePath, nil, parser.ParseComments)
		Expect(err).NotTo(HaveOccurred())

		dstFile, err := decorator.DecorateFile(fSet, inFile)
		Expect(err).NotTo(HaveOccurred())

		outFile, err := os.Create(outPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(dst.Fprint(outFile, dstFile, dst.NotNilFilter)).To(Succeed())
	})

	It("dumps the AST of a function mock", func() {
		filePath := "./mock_loadtypesfn_test.go"
		outPath := "./mock_loadtypesfn_ast.txt"

		fSet := token.NewFileSet()
		inFile, err := parser.ParseFile(fSet, filePath, nil, parser.ParseComments)
		Expect(err).NotTo(HaveOccurred())

		outFile, err := os.Create(outPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(ast.Fprint(outFile, fSet, inFile, ast.NotNilFilter)).To(Succeed())
	})

	It("dumps the DST of a function mock", func() {
		filePath := "./mock_loadtypesfn_test.go"
		outPath := "./mock_loadtypesfn_dst.txt"

		fSet := token.NewFileSet()
		inFile, err := parser.ParseFile(fSet, filePath, nil, parser.ParseComments)
		Expect(err).NotTo(HaveOccurred())

		dstFile, err := decorator.DecorateFile(fSet, inFile)
		Expect(err).NotTo(HaveOccurred())

		outFile, err := os.Create(outPath)
		Expect(err).NotTo(HaveOccurred())

		Expect(dst.Fprint(outFile, dstFile, dst.NotNilFilter)).To(Succeed())
	})
})
