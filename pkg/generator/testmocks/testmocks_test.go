package testmocks_test

import (
	"go/parser"
	"go/token"
	"io"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/pkg/generator"
	"github.com/myshkin5/moqueries/pkg/generator/testmocks/exported"
	"github.com/myshkin5/moqueries/pkg/hash"
	"github.com/myshkin5/moqueries/pkg/testing"
)

type results struct {
	sResults        []string
	err             error
	noReturnResults bool
	times           int
	anyTimes        bool
}

type adaptor interface {
	tracksParams() bool
	expectCall(sParams []string, bParam bool, results ...results) interface{}
	invokeMockAndExpectResults(sParams []string, bParam bool, res results)
	bundleParams(sParams []string, bParam bool) interface{}
}

var _ = Describe("TestMocks", func() {
	entries := func() []TableEntry {
		tMoq := testing.NewMockMoqT(GinkgoT())

		var entries []TableEntry
		entries = append(entries, Entry("usualFn", &usualFnAdaptor{
			m: newMockUsualFn(tMoq.Mock())}, tMoq))
		entries = append(entries, Entry("exportedUsualFn", &exportedUsualFnAdaptor{
			m: exported.NewMockUsualFn(tMoq.Mock())}, tMoq))
		entries = append(entries, Entry("noNamesFn", &noNamesFnAdaptor{
			m: newMockNoNamesFn(tMoq.Mock())}, tMoq))
		entries = append(entries, Entry("exportedNoNamesFn", &exportedNoNamesFnAdaptor{
			m: exported.NewMockNoNamesFn(tMoq.Mock())}, tMoq))
		entries = append(entries, Entry("noResultsFn", &noResultsFnAdaptor{
			m: newMockNoResultsFn(tMoq.Mock())}, tMoq))
		entries = append(entries, Entry("exportedNoResultsFn", &exportedNoResultsFnAdaptor{
			m: exported.NewMockNoResultsFn(tMoq.Mock())}, tMoq))
		entries = append(entries, Entry("noParamsFn", &noParamsFnAdaptor{
			m: newMockNoParamsFn(tMoq.Mock())}, tMoq))
		entries = append(entries, Entry("exportedNoParamsFn", &exportedNoParamsFnAdaptor{
			m: exported.NewMockNoParamsFn(tMoq.Mock())}, tMoq))
		entries = append(entries, Entry("nothingFn", &nothingFnAdaptor{
			m: newMockNothingFn(tMoq.Mock())}, tMoq))
		entries = append(entries, Entry("exportedNothingFn", &exportedNothingFnAdaptor{
			m: exported.NewMockNothingFn(tMoq.Mock())}, tMoq))
		entries = append(entries, Entry("variadicFn", &variadicFnAdaptor{
			m: newMockVariadicFn(tMoq.Mock())}, tMoq))
		entries = append(entries, Entry("exportedVariadicFn", &exportedVariadicFnAdaptor{
			m: exported.NewMockVariadicFn(tMoq.Mock())}, tMoq))
		entries = append(entries, Entry("repeatedIdsFn", &repeatedIdsFnAdaptor{
			m: newMockRepeatedIdsFn(tMoq.Mock())}, tMoq))
		entries = append(entries, Entry("exportedRepeatedIdsFn", &exportedRepeatedIdsFnAdaptor{
			m: exported.NewMockRepeatedIdsFn(tMoq.Mock())}, tMoq))

		usualMock := newMockUsual(tMoq.Mock())
		exportUsualMock := exported.NewMockUsual(tMoq.Mock())
		entries = append(entries, Entry("usual", &usualAdaptor{
			m: usualMock}, tMoq))
		entries = append(entries, Entry("exportedUsual", &exportedUsualAdaptor{
			m: exportUsualMock}, tMoq))
		entries = append(entries, Entry("noNames", &noNamesAdaptor{
			m: usualMock}, tMoq))
		entries = append(entries, Entry("exportedNoNames", &exportedNoNamesAdaptor{
			m: exportUsualMock}, tMoq))
		entries = append(entries, Entry("noResults", &noResultsAdaptor{
			m: usualMock}, tMoq))
		entries = append(entries, Entry("exportedNoResults", &exportedNoResultsAdaptor{
			m: exportUsualMock}, tMoq))
		entries = append(entries, Entry("noParams", &noParamsAdaptor{
			m: usualMock}, tMoq))
		entries = append(entries, Entry("exportedNoParams", &exportedNoParamsAdaptor{
			m: exportUsualMock}, tMoq))
		entries = append(entries, Entry("nothing", &nothingAdaptor{
			m: usualMock}, tMoq))
		entries = append(entries, Entry("exportedNothing", &exportedNothingAdaptor{
			m: exportUsualMock}, tMoq))
		entries = append(entries, Entry("variadic", &variadicAdaptor{
			m: usualMock}, tMoq))
		entries = append(entries, Entry("exportedVariadic", &exportedVariadicAdaptor{
			m: exportUsualMock}, tMoq))
		entries = append(entries, Entry("repeatedIds", &repeatedIdsAdaptor{
			m: usualMock}, tMoq))
		entries = append(entries, Entry("exportedRepeatedIds", &exportedRepeatedIdsAdaptor{
			m: exportUsualMock}, tMoq))

		return entries
	}

	DescribeTable("can return different values when configured to do so",
		func(a adaptor, tMoq *testing.MockMoqT) {
			// ASSEMBLE
			if a.tracksParams() {
				a.expectCall([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}
			a.expectCall([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil},
				results{sResults: []string{"green", "purple"}, err: nil})
			if a.tracksParams() {
				a.expectCall([]string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil})

			// ACT
			a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})

			// ASSERT
			Expect(tMoq.Params_Fatalf).To(BeEmpty())
		},
		entries()...,
	)

	DescribeTable("can return the same values multiple times",
		func(a adaptor, tMoq *testing.MockMoqT) {
			// ASSEMBLE
			if a.tracksParams() {
				a.expectCall([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}
			a.expectCall([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil, times: 4},
				results{sResults: []string{"green", "purple"}, err: nil})
			if a.tracksParams() {
				a.expectCall([]string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			for n := 0; n < 4; n++ {
				a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
			}

			// ACT
			a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})

			// ASSERT
			Expect(tMoq.Params_Fatalf).To(BeEmpty())
		},
		entries()...,
	)

	DescribeTable("returns the same value any number of times",
		func(a adaptor, tMoq *testing.MockMoqT) {
			// ASSEMBLE
			if a.tracksParams() {
				a.expectCall([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}
			a.expectCall([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil},
				results{sResults: []string{"green", "purple"}, err: nil, anyTimes: true})
			if a.tracksParams() {
				a.expectCall([]string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil})

			// ACT
			// NB: Don't use values larger than the results channel will hold
			for n := 0; n < 99; n++ {
				a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})
			}
			// ASSERT
			Expect(tMoq.Params_Fatalf).To(BeEmpty())
		},
		entries()...,
	)

	DescribeTable("fails if Times is called without a preceding Return call",
		func(a adaptor, tMoq *testing.MockMoqT) {
			// ASSEMBLE

			// ACT
			a.expectCall([]string{"Hi", "you"}, true,
				results{noReturnResults: true, times: 4})

			// ASSERT
			Expect(tMoq.Params_Fatalf).To(Receive(Equal(testing.MockMoqT_Fatalf_params{
				Format: "Return must be called before calling Times",
				Args:   hash.DeepHash([]interface{}{}),
			})))
		},
		entries()...,
	)

	DescribeTable("fails if AnyTimes is called without a preceding Return call",
		func(a adaptor, tMoq *testing.MockMoqT) {
			// ASSEMBLE

			// ACT
			a.expectCall([]string{"Hi", "you"}, true,
				results{noReturnResults: true, anyTimes: true})

			// ASSERT
			Expect(tMoq.Params_Fatalf).To(Receive(Equal(testing.MockMoqT_Fatalf_params{
				Format: "Return must be called before calling AnyTimes",
				Args:   hash.DeepHash([]interface{}{}),
			})))
		},
		entries()...,
	)

	DescribeTable("fails if the function is called too many times",
		func(a adaptor, tMoq *testing.MockMoqT) {
			// ASSEMBLE
			a.expectCall([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil},
				results{sResults: []string{"green", "purple"}, err: io.EOF})

			a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil})

			a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: io.EOF})

			// ACT
			a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"", ""}, err: nil})

			// ASSERT
			Expect(tMoq.Params_Fatalf).To(Receive(Equal(testing.MockMoqT_Fatalf_params{
				Format: "Too many calls to mock with parameters %#v",
				Args: hash.DeepHash([]interface{}{
					a.bundleParams([]string{"Hi", "you"}, true),
				}),
			})))
		},
		entries()...,
	)

	DescribeTable("fails if expectations are set more than once for the same parameter set",
		func(a adaptor, tMoq *testing.MockMoqT) {
			// ASSEMBLE
			a.expectCall([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})

			// ACT
			fnRec := a.expectCall([]string{"Hi", "you"}, true,
				results{sResults: []string{"red", "purple"}, err: io.EOF})

			// ASSERT
			Expect(fnRec).To(BeNil())
			Expect(tMoq.Params_Fatalf).To(Receive(Equal(testing.MockMoqT_Fatalf_params{
				Format: "Expectations already recorded for mock with parameters %#v",
				Args: hash.DeepHash([]interface{}{
					a.bundleParams([]string{"Hi", "you"}, true),
				}),
			})))
		},
		entries()...,
	)

	PIt("generates mocks", func() {
		// NB: Keep in sync with types.go go:generate directives

		// These lines generate the same mocks listed in types.go go:generate
		// directives. Remove the "pending" flag on this test to verify code
		// coverage.

		Expect(generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_usualfn_test.go", Types: []string{"UsualFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_usualfn.go", Export: true, Types: []string{"UsualFn"},
			},
		)).To(Succeed())

		Expect(generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_nonamesfn_test.go", Types: []string{"NoNamesFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_nonamesfn.go", Export: true, Types: []string{"NoNamesFn"},
			},
		)).To(Succeed())

		Expect(generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_noresultsfn_test.go", Types: []string{"NoResultsFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_noresultsfn.go", Export: true, Types: []string{"NoResultsFn"},
			},
		)).To(Succeed())

		Expect(generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_noparamsfn_test.go", Types: []string{"NoParamsFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_noparamsfn.go", Export: true, Types: []string{"NoParamsFn"},
			},
		)).To(Succeed())

		Expect(generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_nothingfn_test.go", Types: []string{"NothingFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_nothingfn.go", Export: true, Types: []string{"NothingFn"},
			},
		)).To(Succeed())

		Expect(generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_variadicfn_test.go", Types: []string{"VariadicFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_variadicfn.go", Export: true, Types: []string{"VariadicFn"},
			},
		)).To(Succeed())

		Expect(generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_usual_test.go", Types: []string{"Usual"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_usual.go", Export: true, Types: []string{"Usual"},
			},
		)).To(Succeed())
	})

	PIt("dumps the DST of an interface mock", func() {
		filePath := "./moq_usualfn_test.go"
		outPath := "./moq_usualfn_test_dst.txt"

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
