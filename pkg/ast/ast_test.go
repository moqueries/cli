package ast_test

import (
	"strings"

	"github.com/dave/dst"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/pkg/ast"
)

var _ = Describe("Ast", func() {
	Describe("FindPackageDir", func() {
		It("finds the package dir of the current directory", func() {
			// ASSEMBLE

			// ACT
			dir, err := ast.FindPackageDir(".")

			// ASSERT
			Expect(err).NotTo(HaveOccurred())
			Expect(dir).To(HaveSuffix("github.com/myshkin5/moqueries/pkg/ast"))
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
			typs, pkgPath, err := ast.LoadTypes("builtin")

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
			typs, pkgPath, err := ast.LoadTypes("bufio")

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
	})
})
