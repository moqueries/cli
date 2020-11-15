package testmocks_test

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/pkg/generator"
	"github.com/myshkin5/moqueries/pkg/generator/testmocks/exported"
	"github.com/myshkin5/moqueries/pkg/moq"
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
	newRecorder(sParams []string, bParam bool) interface{}
	any(rec interface{}, sParams, bParam bool) interface{}
	results(rec interface{}, results ...results) interface{}
	invokeMockAndExpectResults(sParams []string, bParam bool, res results)
	bundleParams(sParams []string, bParam bool) interface{}
	sceneMock() moq.Mock
}

func expectCall(a adaptor, sParams []string, bParam bool, results ...results) interface{} {
	rec := a.newRecorder(sParams, bParam)
	return a.results(rec, results...)
}

type tableEntry struct {
	description string
	a           adaptor
}

var (
	scene    *moq.Scene
	moqTMock *moq.MockMoqT
)

var _ = BeforeSuite(func() {
	scene = moq.NewScene(GinkgoT())
	moqTMock = moq.NewMockMoqT(scene, nil)
})

func entries(moqTMock *moq.MockMoqT, c moq.MockConfig) ([]tableEntry, *moq.Scene) {
	mockScene := moq.NewScene(moqTMock.Mock())

	var entries []tableEntry
	entries = append(entries, tableEntry{"usualFn", &usualFnAdaptor{
		m: newMockUsualFn(mockScene, &c)}})
	entries = append(entries, tableEntry{"exportedUsualFn", &exportedUsualFnAdaptor{
		m: exported.NewMockUsualFn(mockScene, &c)}})
	entries = append(entries, tableEntry{"noNamesFn", &noNamesFnAdaptor{
		m: newMockNoNamesFn(mockScene, &c)}})
	entries = append(entries, tableEntry{"exportedNoNamesFn", &exportedNoNamesFnAdaptor{
		m: exported.NewMockNoNamesFn(mockScene, &c)}})
	entries = append(entries, tableEntry{"noResultsFn", &noResultsFnAdaptor{
		m: newMockNoResultsFn(mockScene, &c)}})
	entries = append(entries, tableEntry{"exportedNoResultsFn", &exportedNoResultsFnAdaptor{
		m: exported.NewMockNoResultsFn(mockScene, &c)}})
	entries = append(entries, tableEntry{"noParamsFn", &noParamsFnAdaptor{
		m: newMockNoParamsFn(mockScene, &c)}})
	entries = append(entries, tableEntry{"exportedNoParamsFn", &exportedNoParamsFnAdaptor{
		m: exported.NewMockNoParamsFn(mockScene, &c)}})
	entries = append(entries, tableEntry{"nothingFn", &nothingFnAdaptor{
		m: newMockNothingFn(mockScene, &c)}})
	entries = append(entries, tableEntry{"exportedNothingFn", &exportedNothingFnAdaptor{
		m: exported.NewMockNothingFn(mockScene, &c)}})
	entries = append(entries, tableEntry{"variadicFn", &variadicFnAdaptor{
		m: newMockVariadicFn(mockScene, &c)}})
	entries = append(entries, tableEntry{"exportedVariadicFn", &exportedVariadicFnAdaptor{
		m: exported.NewMockVariadicFn(mockScene, &c)}})
	entries = append(entries, tableEntry{"repeatedIdsFn", &repeatedIdsFnAdaptor{
		m: newMockRepeatedIdsFn(mockScene, &c)}})
	entries = append(entries, tableEntry{"exportedRepeatedIdsFn", &exportedRepeatedIdsFnAdaptor{
		m: exported.NewMockRepeatedIdsFn(mockScene, &c)}})

	usualMock := newMockUsual(mockScene, &c)
	exportUsualMock := exported.NewMockUsual(mockScene, &c)
	entries = append(entries, tableEntry{"usual", &usualAdaptor{
		m: usualMock}})
	entries = append(entries, tableEntry{"exportedUsual", &exportedUsualAdaptor{
		m: exportUsualMock}})
	entries = append(entries, tableEntry{"noNames", &noNamesAdaptor{
		m: usualMock}})
	entries = append(entries, tableEntry{"exportedNoNames", &exportedNoNamesAdaptor{
		m: exportUsualMock}})
	entries = append(entries, tableEntry{"noResults", &noResultsAdaptor{
		m: usualMock}})
	entries = append(entries, tableEntry{"exportedNoResults", &exportedNoResultsAdaptor{
		m: exportUsualMock}})
	entries = append(entries, tableEntry{"noParams", &noParamsAdaptor{
		m: usualMock}})
	entries = append(entries, tableEntry{"exportedNoParams", &exportedNoParamsAdaptor{
		m: exportUsualMock}})
	entries = append(entries, tableEntry{"nothing", &nothingAdaptor{
		m: usualMock}})
	entries = append(entries, tableEntry{"exportedNothing", &exportedNothingAdaptor{
		m: exportUsualMock}})
	entries = append(entries, tableEntry{"variadic", &variadicAdaptor{
		m: usualMock}})
	entries = append(entries, tableEntry{"exportedVariadic", &exportedVariadicAdaptor{
		m: exportUsualMock}})
	entries = append(entries, tableEntry{"repeatedIds", &repeatedIdsAdaptor{
		m: usualMock}})
	entries = append(entries, tableEntry{"exportedRepeatedIds", &exportedRepeatedIdsAdaptor{
		m: exportUsualMock}})

	return entries, mockScene
}

var _ = Describe("TestMocks", func() {
	It("can return different values when configured to do so", func() {
		entries, mockScene := entries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			if entry.a.tracksParams() {
				expectCall(entry.a, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}
			expectCall(entry.a, []string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil},
				results{sResults: []string{"green", "purple"}, err: nil})
			if entry.a.tracksParams() {
				expectCall(entry.a, []string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			if entry.a.tracksParams() {
				entry.a.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}
			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil})
			if entry.a.tracksParams() {
				entry.a.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			// ACT
			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})

			// ASSERT
			mockScene.AssertExpectationsMet()
			scene.AssertExpectationsMet()
		}
	})

	It("fails if an expectation is not set in strict mode", func() {
		entries, mockScene := entries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			msg := "Unexpected call with parameters %#v"
			params := entry.a.bundleParams([]string{"Hi", "you"}, true)
			moqTMock.OnCall().Fatalf(msg, params).
				ReturnResults()
			fmtMsg := fmt.Sprintf(msg, params)

			// ACT
			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"", ""}, err: nil})

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
			if entry.a.tracksParams() {
				Expect(fmtMsg).To(ContainSubstring("Hi"))
			}
		}
	})

	It("returns zero values if an expectation is not set in nice mode", func() {
		entries, mockScene := entries(
			moqTMock, moq.MockConfig{Expectation: moq.Nice})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			// ACT
			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"", ""}, err: nil})

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("can return the same values multiple times", func() {
		entries, mockScene := entries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			if entry.a.tracksParams() {
				expectCall(entry.a, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}
			expectCall(entry.a, []string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil, times: 4},
				results{sResults: []string{"green", "purple"}, err: nil})
			if entry.a.tracksParams() {
				expectCall(entry.a, []string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			if entry.a.tracksParams() {
				entry.a.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}
			for n := 0; n < 4; n++ {
				entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
			}
			if entry.a.tracksParams() {
				entry.a.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			// ACT
			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("returns the same value any number of times", func() {
		entries, mockScene := entries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			if entry.a.tracksParams() {
				expectCall(entry.a, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}
			expectCall(entry.a, []string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil},
				results{sResults: []string{"green", "purple"}, err: nil, anyTimes: true})
			if entry.a.tracksParams() {
				expectCall(entry.a, []string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			if entry.a.tracksParams() {
				entry.a.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}
			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil})
			if entry.a.tracksParams() {
				entry.a.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			// ACT
			// NB: Don't use values larger than the results channel will hold
			for n := 0; n < 20; n++ {
				entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})
			}

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("fails if Times is called without a preceding Return call", func() {
		entries, mockScene := entries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			// ASSEMBLE
			moqTMock.OnCall().Fatalf("Return must be called before calling Times").
				ReturnResults()

			// ACT
			expectCall(entry.a, []string{"Hi", "you"}, true,
				results{noReturnResults: true, times: 4})

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("fails if AnyTimes is called without a preceding Return call", func() {
		entries, mockScene := entries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			// ASSEMBLE
			moqTMock.OnCall().Fatalf("Return must be called before calling AnyTimes").
				ReturnResults()

			// ACT
			expectCall(entry.a, []string{"Hi", "you"}, true,
				results{noReturnResults: true, anyTimes: true})

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("fails if the function is called too many times in strict mode", func() {
		entries, mockScene := entries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			expectCall(entry.a, []string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil},
				results{sResults: []string{"green", "purple"}, err: io.EOF})

			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil})

			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: io.EOF})

			moqTMock.OnCall().Fatalf(
				"Too many calls to mock with parameters %#v",
				entry.a.bundleParams([]string{"Hi", "you"}, true)).ReturnResults()

			// ACT
			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"", ""}, err: nil})

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("returns zero values if the function is called too many times in nice mode", func() {
		entries, mockScene := entries(moqTMock, moq.MockConfig{Expectation: moq.Nice})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			expectCall(entry.a, []string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil},
				results{sResults: []string{"green", "purple"}, err: io.EOF})

			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil})

			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: io.EOF})

			// ACT
			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"", ""}, err: nil})

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("fails if expectations are set more than once for the same parameter set", func() {
		entries, mockScene := entries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			expectCall(entry.a, []string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})
			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})

			msg := "Expectations already recorded for mock with parameters %#v"
			bParams := entry.a.bundleParams([]string{"Hi", "you"}, true)
			moqTMock.OnCall().Fatalf(msg, bParams).
				ReturnResults()
			fmtMsg := fmt.Sprintf(msg, bParams)

			// ACT
			fnRec := expectCall(entry.a, []string{"Hi", "you"}, true,
				results{sResults: []string{"red", "purple"}, err: io.EOF})

			// ASSERT
			Expect(fnRec).To(BeNil())
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
			if entry.a.tracksParams() {
				Expect(fmtMsg).To(ContainSubstring("Hi"))
			}
		}
	})

	It("resets its state", func() {
		entries, mockScene := entries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			expectCall(entry.a, []string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})

			// ACT
			entry.a.sceneMock().Reset()

			// ASSERT
			expectCall(entry.a, []string{"Hi", "you"}, true,
				results{sResults: []string{"grey", "indigo"}, err: nil})
			entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"grey", "indigo"}, err: nil})

			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	Describe("any values", func() {
		It("ignores a param when instructed to do so", func() {
			entries, mockScene := entries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				if entry.a.tracksParams() {
					expectCall(entry.a, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				rec := entry.a.newRecorder([]string{"Hi", "you"}, true)
				rec = entry.a.any(rec, true, false)
				entry.a.results(rec,
					results{sResults: []string{"blue", "orange"}, err: nil},
					results{sResults: []string{"green", "purple"}, err: nil})
				if entry.a.tracksParams() {
					expectCall(entry.a, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				if entry.a.tracksParams() {
					entry.a.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if entry.a.tracksParams() {
					entry.a.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				// ACT
				entry.a.invokeMockAndExpectResults([]string{"Goodbye", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})

				// ASSERT
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
			}
		})

		It("returns an error if any functions are called after returning results", func() {
			entries, mockScene := entries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				if !entry.a.tracksParams() {
					// With no params to track, there will be no `any` functions
					continue
				}
				scene.Reset()
				mockScene.Reset()

				rec := entry.a.newRecorder([]string{"Hi", "you"}, true)
				rec = entry.a.results(rec,
					results{sResults: []string{"blue", "orange"}, err: nil})
				entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				msg := "Any functions must be called prior to returning results, parameters: %#v"
				bParams := entry.a.bundleParams([]string{"Hi", "you"}, true)
				moqTMock.OnCall().Fatalf(msg, bParams).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, bParams)

				// ACT
				rec = entry.a.any(rec, true, false)

				// ASSERT
				Expect(rec).To(BeNil())
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
				if entry.a.tracksParams() {
					Expect(fmtMsg).To(ContainSubstring("Hi"))
				}
			}
		})
	})

	Describe("AssertExpectationsMet", func() {
		It("reports success when there ae no expectations", func() {
			entries, mockScene := entries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				// ACT
				entry.a.sceneMock().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
			}
		})

		It("reports success when all expectations are met", func() {
			entries, mockScene := entries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				if entry.a.tracksParams() {
					expectCall(entry.a, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				expectCall(entry.a, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if entry.a.tracksParams() {
					expectCall(entry.a, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				if entry.a.tracksParams() {
					entry.a.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if entry.a.tracksParams() {
					entry.a.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				// ACT
				entry.a.sceneMock().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
			}
		})

		It("fails when one expectation is not met", func() {
			entries, mockScene := entries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				// ASSEMBLE
				if entry.a.tracksParams() {
					expectCall(entry.a, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				expectCall(entry.a, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if entry.a.tracksParams() {
					expectCall(entry.a, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				if entry.a.tracksParams() {
					entry.a.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				if entry.a.tracksParams() {
					entry.a.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				msg := "Expected %d additional call(s) with parameters %#v"
				params := entry.a.bundleParams([]string{"Hi", "you"}, true)
				moqTMock.OnCall().Errorf(msg, 1, params).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, 1, params)

				// ACT
				entry.a.sceneMock().AssertExpectationsMet()

				// ASSERT
				entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
				if entry.a.tracksParams() {
					Expect(fmtMsg).To(ContainSubstring("Hi"))
				}
			}
		})

		It("succeeds when an anyTimes expectation is not called at all", func() {
			entries, mockScene := entries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				expectCall(entry.a, []string{"Hi", "you"}, true,
					results{
						sResults: []string{"blue", "orange"},
						err:      nil,
						anyTimes: true,
					})

				// ACT
				entry.a.sceneMock().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
			}
		})

		It("succeeds when an anyTimes expectation is called once", func() {
			entries, mockScene := entries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				// ASSEMBLE
				expectCall(entry.a, []string{"Hi", "you"}, true,
					results{
						sResults: []string{"blue", "orange"},
						err:      nil,
						anyTimes: true,
					})

				entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ACT
				entry.a.sceneMock().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
			}
		})

		It("succeeds when an anyTimes expectation is called many times", func() {
			entries, mockScene := entries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				// ASSEMBLE
				expectCall(entry.a, []string{"Hi", "you"}, true,
					results{
						sResults: []string{"blue", "orange"},
						err:      nil,
						anyTimes: true,
					})

				entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				entry.a.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ACT
				entry.a.sceneMock().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
			}
		})
	})

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

	It("dumps the DST of an interface mock", func() {
		filePath := "./moq_usual_test.go"
		outPath := "./moq_usual_test_dst.txt"

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
