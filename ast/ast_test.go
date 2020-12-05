package ast_test

import (
	"strings"

	"github.com/dave/dst"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/ast"
)

var _ = Describe("Ast", func() {
	Describe("FindPackageDir", func() {
		It("finds the package dir of the current directory", func() {
			// ASSEMBLE

			// ACT
			dir, err := ast.FindPackageDir(".")

			// ASSERT
			Expect(err).NotTo(HaveOccurred())
			Expect(dir).To(HaveSuffix("github.com/myshkin5/moqueries/ast"))
		})

		It("finds the package dir of an external module", func() {
			// ASSEMBLE

			// ACT
			dir, err := ast.FindPackageDir("github.com/onsi/ginkgo")

			// ASSERT
			Expect(err).NotTo(HaveOccurred())
			splits := strings.Split(dir, "@")
			Expect(splits).To(HaveLen(2))
			Expect(splits[0]).To(HaveSuffix("pkg/mod/github.com/onsi/ginkgo"))
		})

		It("returns an error when there are no files in the package", func() {
			// ASSEMBLE

			// ACT
			_, err := ast.FindPackageDir("randomandnonexistent")

			// ASSERT
			Expect(err).To(MatchError("no files found for package randomandnonexistent"))
		})
	})

	Describe("LoadTypes", func() {
		It("loads the expected interfaces", func() {
			// ASSEMBLE

			// ACT
			typs, pkgPath, err := ast.LoadTypes("builtin", false)

			// ASSERT
			Expect(err).NotTo(HaveOccurred())
			var iTypes []*dst.TypeSpec
			for _, typ := range typs {
				if _, ok := typ.Type.(*dst.InterfaceType); ok {
					iTypes = append(iTypes, typ)
				}
			}
			Expect(iTypes).To(HaveLen(1))
			Expect(iTypes[0].Name.Name).To(Equal("error"))
			Expect(pkgPath).To(Equal("builtin"))
		})

		It("loads the expected functions", func() {
			// ASSEMBLE

			// ACT
			typs, pkgPath, err := ast.LoadTypes("bufio", false)

			// ASSERT
			Expect(err).NotTo(HaveOccurred())
			var fTypes []*dst.TypeSpec
			for _, typ := range typs {
				if _, ok := typ.Type.(*dst.FuncType); ok {
					fTypes = append(fTypes, typ)
				}
			}
			Expect(fTypes).To(HaveLen(1))
			Expect(fTypes[0].Name.Name).To(Equal("SplitFunc"))
			Expect(pkgPath).To(Equal("bufio"))
		})

		It("resolves local paths", func() {
			// ASSEMBLE

			// ACT
			typs, pkgPath, err := ast.LoadTypes("github.com/myshkin5/moqueries/generator", false)

			// ASSERT
			Expect(err).NotTo(HaveOccurred())

			var iTypes []*dst.TypeSpec
			for _, typ := range typs {
				if _, ok := typ.Type.(*dst.InterfaceType); ok {
					iTypes = append(iTypes, typ)
				}
			}
			Expect(iTypes).To(HaveLen(1))
			Expect(iTypes[0].Name.Name).To(Equal("Converterer"))
			Expect(pkgPath).To(Equal("github.com/myshkin5/moqueries/generator"))

			var baseStruct *dst.FuncType
			for _, field := range iTypes[0].Type.(*dst.InterfaceType).Methods.List {
				if field.Names[0].Name == "BaseStruct" {
					baseStruct = field.Type.(*dst.FuncType)
				}
			}
			Expect(baseStruct).NotTo(BeNil())
			Expect(baseStruct.Params.List[1].Names[0].Name).To(Equal("funcs"))
			funcIdent := baseStruct.Params.List[1].Type.(*dst.ArrayType).Elt.(*dst.Ident)
			Expect(funcIdent.Path).To(Equal("github.com/myshkin5/moqueries/generator"))
		})

		It("loads test types", func() {
			// ASSEMBLE

			// ACT
			typs, pkgPath, err := ast.LoadTypes("github.com/myshkin5/moqueries/ast", true)

			// ASSERT
			Expect(err).NotTo(HaveOccurred())

			var fTypes []*dst.TypeSpec
			for _, typ := range typs {
				if _, ok := typ.Type.(*dst.FuncType); ok {
					fTypes = append(fTypes, typ)
				}
			}
			Expect(fTypes).To(HaveLen(1))
			Expect(fTypes[0].Name.Name).To(Equal("TestFn"))
			Expect(pkgPath).To(Equal("github.com/myshkin5/moqueries/ast.test"))
		})
	})
})

type TestFn func()
