package testmocks_test

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/pkg/generator"
	"github.com/myshkin5/moqueries/pkg/generator/testmocks/exported"
	"github.com/myshkin5/moqueries/pkg/moq"
)

type results struct {
	sResults []string
	err      error
}

type adaptor interface {
	exported() bool
	tracksParams() bool
	newRecorder(sParams []string, bParam bool) recorder
	invokeMockAndExpectResults(sParams []string, bParam bool, res results)
	bundleParams(sParams []string, bParam bool) interface{}
	sceneMock() moq.Mock
}

type recorder interface {
	anySParam()
	anyBParam()
	seq()
	noSeq()
	returnResults(sResults []string, err error)
	andDo(fn func(), sParams []string, bParam bool)
	doReturnResults(fn func(), sParams []string, bParam bool, sResults []string, err error)
	times(count int)
	anyTimes()
	isNil() bool
}

func expectCall(a adaptor, sParams []string, bParam bool, results ...results) {
	rec := a.newRecorder(sParams, bParam)
	for _, res := range results {
		rec.returnResults(res.sResults, res.err)
	}
}

var (
	scene    *moq.Scene
	moqTMock *moq.MockMoqT
)

var _ = BeforeSuite(func() {
	scene = moq.NewScene(GinkgoT())
	moqTMock = moq.NewMockMoqT(scene, nil)
})

func tableEntries(moqTMock *moq.MockMoqT, c moq.MockConfig) ([]adaptor, *moq.Scene) {
	mockScene := moq.NewScene(moqTMock.Mock())

	var entries []adaptor
	entries = append(entries, &usualFnAdaptor{m: newMockUsualFn(mockScene, &c)})
	entries = append(entries, &exportedUsualFnAdaptor{m: exported.NewMockUsualFn(mockScene, &c)})
	entries = append(entries, &noNamesFnAdaptor{m: newMockNoNamesFn(mockScene, &c)})
	entries = append(entries, &exportedNoNamesFnAdaptor{m: exported.NewMockNoNamesFn(mockScene, &c)})
	entries = append(entries, &noResultsFnAdaptor{m: newMockNoResultsFn(mockScene, &c)})
	entries = append(entries, &exportedNoResultsFnAdaptor{m: exported.NewMockNoResultsFn(mockScene, &c)})
	entries = append(entries, &noParamsFnAdaptor{m: newMockNoParamsFn(mockScene, &c)})
	entries = append(entries, &exportedNoParamsFnAdaptor{m: exported.NewMockNoParamsFn(mockScene, &c)})
	entries = append(entries, &nothingFnAdaptor{m: newMockNothingFn(mockScene, &c)})
	entries = append(entries, &exportedNothingFnAdaptor{m: exported.NewMockNothingFn(mockScene, &c)})
	entries = append(entries, &variadicFnAdaptor{m: newMockVariadicFn(mockScene, &c)})
	entries = append(entries, &exportedVariadicFnAdaptor{m: exported.NewMockVariadicFn(mockScene, &c)})
	entries = append(entries, &repeatedIdsFnAdaptor{m: newMockRepeatedIdsFn(mockScene, &c)})
	entries = append(entries, &exportedRepeatedIdsFnAdaptor{m: exported.NewMockRepeatedIdsFn(mockScene, &c)})

	usualMock := newMockUsual(mockScene, &c)
	exportUsualMock := exported.NewMockUsual(mockScene, &c)
	entries = append(entries, &usualAdaptor{m: usualMock})
	entries = append(entries, &exportedUsualAdaptor{m: exportUsualMock})
	entries = append(entries, &noNamesAdaptor{m: usualMock})
	entries = append(entries, &exportedNoNamesAdaptor{m: exportUsualMock})
	entries = append(entries, &noResultsAdaptor{m: usualMock})
	entries = append(entries, &exportedNoResultsAdaptor{m: exportUsualMock})
	entries = append(entries, &noParamsAdaptor{m: usualMock})
	entries = append(entries, &exportedNoParamsAdaptor{m: exportUsualMock})
	entries = append(entries, &nothingAdaptor{m: usualMock})
	entries = append(entries, &exportedNothingAdaptor{m: exportUsualMock})
	entries = append(entries, &variadicAdaptor{m: usualMock})
	entries = append(entries, &exportedVariadicAdaptor{m: exportUsualMock})
	entries = append(entries, &repeatedIdsAdaptor{m: usualMock})
	entries = append(entries, &exportedRepeatedIdsAdaptor{m: exportUsualMock})

	return entries, mockScene
}

var _ = Describe("TestMocks", func() {
	It("can return different values when configured to do so", func() {
		entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			if entry.tracksParams() {
				expectCall(entry, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}
			expectCall(entry, []string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil},
				results{sResults: []string{"green", "purple"}, err: nil})
			if entry.tracksParams() {
				expectCall(entry, []string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			if entry.tracksParams() {
				entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}
			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil})
			if entry.tracksParams() {
				entry.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			// ACT
			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})

			// ASSERT
			mockScene.AssertExpectationsMet()
			scene.AssertExpectationsMet()
		}
	})

	It("fails if an expectation is not set in strict mode", func() {
		entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			msg := "Unexpected call with parameters %#v"
			params := entry.bundleParams([]string{"Hi", "you"}, true)
			moqTMock.OnCall().Fatalf(msg, params).
				ReturnResults()
			fmtMsg := fmt.Sprintf(msg, params)

			// ACT
			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"", ""}, err: nil})

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
			if entry.tracksParams() {
				Expect(fmtMsg).To(ContainSubstring("Hi"))
			}
		}
	})

	It("returns zero values if an expectation is not set in nice mode", func() {
		entries, mockScene := tableEntries(
			moqTMock, moq.MockConfig{Expectation: moq.Nice})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			// ACT
			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"", ""}, err: nil})

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("can return the same values multiple times", func() {
		entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			if entry.tracksParams() {
				expectCall(entry, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			rec := entry.newRecorder([]string{"Hi", "you"}, true)
			rec.returnResults([]string{"blue", "orange"}, nil)
			rec.times(4)
			rec.returnResults([]string{"green", "purple"}, nil)

			if entry.tracksParams() {
				expectCall(entry, []string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			if entry.tracksParams() {
				entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}
			for n := 0; n < 4; n++ {
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
			}
			if entry.tracksParams() {
				entry.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			// ACT
			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("returns the same value any number of times", func() {
		entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			if entry.tracksParams() {
				expectCall(entry, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			rec := entry.newRecorder([]string{"Hi", "you"}, true)
			rec.returnResults([]string{"blue", "orange"}, nil)
			rec.returnResults([]string{"green", "purple"}, nil)
			rec.anyTimes()

			if entry.tracksParams() {
				expectCall(entry, []string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			if entry.tracksParams() {
				entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}
			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil})
			if entry.tracksParams() {
				entry.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
			}

			// ACT
			// NB: Don't use values larger than the results channel will hold
			for n := 0; n < 20; n++ {
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})
			}

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("fails if Times is called without a preceding Return call", func() {
		entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			msg := fmt.Sprintf("%s or %s must be called before calling %s",
				export("returnResults", entry),
				export("doReturnResults", entry),
				export("times", entry))

			moqTMock.OnCall().Fatalf(msg).
				ReturnResults()

			rec := entry.newRecorder([]string{"Hi", "you"}, true)

			// ACT
			rec.times(4)

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("fails if AnyTimes is called without a preceding Return call", func() {
		entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			msg := fmt.Sprintf("%s or %s must be called before calling %s",
				export("returnResults", entry),
				export("doReturnResults", entry),
				export("anyTimes", entry))

			moqTMock.OnCall().Fatalf(msg).
				ReturnResults()

			rec := entry.newRecorder([]string{"Hi", "you"}, true)

			// ACT
			rec.anyTimes()

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("fails if the function is called too many times in strict mode", func() {
		entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			expectCall(entry, []string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil},
				results{sResults: []string{"green", "purple"}, err: io.EOF})

			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil})

			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: io.EOF})

			moqTMock.OnCall().Fatalf(
				"Too many calls to mock with parameters %#v",
				entry.bundleParams([]string{"Hi", "you"}, true)).ReturnResults()

			// ACT
			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"", ""}, err: nil})

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("returns zero values if the function is called too many times in nice mode", func() {
		entries, mockScene := tableEntries(moqTMock, moq.MockConfig{Expectation: moq.Nice})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			expectCall(entry, []string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil},
				results{sResults: []string{"green", "purple"}, err: io.EOF})

			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"blue", "orange"}, err: nil})

			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: io.EOF})

			// ACT
			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"", ""}, err: nil})

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("adds additional results when an expectation has already been set", func() {
		entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			expectCall(entry, []string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})

			expectCall(entry, []string{"Hi", "you"}, true,
				results{sResults: []string{"red", "blue"}, err: nil})

			// ACT
			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})
			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"red", "blue"}, err: nil})

			// ASSERT
			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	It("resets its state", func() {
		entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
		for _, entry := range entries {
			// ASSEMBLE
			scene.Reset()
			mockScene.Reset()

			expectCall(entry, []string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})

			// ACT
			entry.sceneMock().Reset()

			// ASSERT
			expectCall(entry, []string{"Hi", "you"}, true,
				results{sResults: []string{"grey", "indigo"}, err: nil})
			entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
				results{sResults: []string{"grey", "indigo"}, err: nil})

			scene.AssertExpectationsMet()
			mockScene.AssertExpectationsMet()
		}
	})

	Describe("any values", func() {
		It("ignores a param when instructed to do so", func() {
			entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				if entry.tracksParams() {
					expectCall(entry, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.anySParam()
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.returnResults([]string{"green", "purple"}, nil)

				if entry.tracksParams() {
					expectCall(entry, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				if entry.tracksParams() {
					entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if entry.tracksParams() {
					entry.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				// ACT
				entry.invokeMockAndExpectResults([]string{"Goodbye", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})

				// ASSERT
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
			}
		})

		It("returns an error if any functions are called after returning results", func() {
			entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				if !entry.tracksParams() {
					// With no params to track, there will be no `any` functions
					continue
				}
				scene.Reset()
				mockScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				rrFn := "returnResults"
				drrFn := "doReturnResults"
				if entry.exported() {
					rrFn = strings.Title(rrFn)
					drrFn = strings.Title(drrFn)
				}

				msg := fmt.Sprintf(
					"Any functions must be called before %s or %s calls, parameters: %%#v",
					rrFn, drrFn)
				bParams := entry.bundleParams([]string{"Hi", "you"}, true)
				moqTMock.OnCall().Fatalf(msg, bParams).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, bParams)

				// ACT
				rec.anySParam()

				// ASSERT
				Expect(rec.isNil()).To(BeTrue())
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
				if entry.tracksParams() {
					Expect(fmtMsg).To(ContainSubstring("Hi"))
				}
			}
		})
	})

	Describe("AssertExpectationsMet", func() {
		It("reports success when there ae no expectations", func() {
			entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				// ACT
				entry.sceneMock().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
			}
		})

		It("reports success when all expectations are met", func() {
			entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				if entry.tracksParams() {
					expectCall(entry, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				expectCall(entry, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if entry.tracksParams() {
					expectCall(entry, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				if entry.tracksParams() {
					entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if entry.tracksParams() {
					entry.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				// ACT
				entry.sceneMock().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
			}
		})

		It("fails when one expectation is not met", func() {
			entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				// ASSEMBLE
				if entry.tracksParams() {
					expectCall(entry, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				expectCall(entry, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if entry.tracksParams() {
					expectCall(entry, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				if entry.tracksParams() {
					entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				if entry.tracksParams() {
					entry.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				msg := "Expected %d additional call(s) with parameters %#v"
				params := entry.bundleParams([]string{"Hi", "you"}, true)
				moqTMock.OnCall().Errorf(msg, 1, params).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, 1, params)

				// ACT
				entry.sceneMock().AssertExpectationsMet()

				// ASSERT
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
				if entry.tracksParams() {
					Expect(fmtMsg).To(ContainSubstring("Hi"))
				}
			}
		})

		It("succeeds when an anyTimes expectation is not called at all", func() {
			entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.anyTimes()

				// ACT
				entry.sceneMock().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
			}
		})

		It("succeeds when an anyTimes expectation is called once", func() {
			entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.anyTimes()

				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ACT
				entry.sceneMock().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
			}
		})

		It("succeeds when an anyTimes expectation is called many times", func() {
			entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				// ASSEMBLE
				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.anyTimes()

				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ACT
				entry.sceneMock().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				mockScene.AssertExpectationsMet()
			}
		})
	})

	Describe("sequences", func() {
		It("passes when sequences are correct", func() {
			entries, mockScene := tableEntries(
				moqTMock, moq.MockConfig{Sequence: moq.SeqDefaultOn})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				if entry.tracksParams() {
					expectCall(entry, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				expectCall(entry, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil},
					results{sResults: []string{"green", "purple"}, err: nil})
				if entry.tracksParams() {
					expectCall(entry, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				// ACT
				if entry.tracksParams() {
					entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})

				if entry.tracksParams() {
					entry.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				// ASSERT
				mockScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			}
		})

		It("fails when sequences are incorrect", func() {
			entries, mockScene := tableEntries(
				moqTMock, moq.MockConfig{Sequence: moq.SeqDefaultOn})
			for _, entry := range entries {
				// ASSEMBLE
				if !entry.tracksParams() {
					// With no params to track, hard to make conflicting calls
					continue
				}
				scene.Reset()
				mockScene.Reset()

				expectCall(entry, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
				expectCall(entry, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				msg := "Call sequence does not match %#v"
				params1 := entry.bundleParams([]string{"Hi", "you"}, true)
				moqTMock.OnCall().Fatalf(msg, params1).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, params1)
				params2 := entry.bundleParams([]string{"Hello", "there"}, false)
				moqTMock.OnCall().Errorf("Expected %d additional call(s) with parameters %#v", 1, params2).
					ReturnResults()

				// ACT
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ASSERT
				Expect(fmtMsg).To(ContainSubstring("Hi"))

				mockScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			}
		})

		It("can have sequences turned on for select calls", func() {
			entries, mockScene := tableEntries(
				moqTMock, moq.MockConfig{Sequence: moq.SeqDefaultOff})
			for _, entry := range entries {
				// ASSEMBLE
				if !entry.tracksParams() {
					// With no params to track, hard to make conflicting calls
					continue
				}
				scene.Reset()
				mockScene.Reset()

				expectCall(entry, []string{"Eh", "you"}, true, results{
					sResults: []string{"grey", "brown"}, err: nil})

				rec := entry.newRecorder([]string{"Hello", "there"}, false)
				rec.seq()
				rec.returnResults([]string{"red", "yellow"}, io.EOF)

				rec = entry.newRecorder([]string{"Hi", "you"}, true)
				rec.seq()
				rec.returnResults([]string{"blue", "orange"}, nil)

				expectCall(entry, []string{"Bye", "now"}, true, results{
					sResults: []string{"silver", "black"}, err: nil})

				entry.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
					results{sResults: []string{"silver", "black"}, err: nil})

				// ACT
				entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ASSERT
				entry.invokeMockAndExpectResults([]string{"Eh", "you"}, true,
					results{sResults: []string{"grey", "brown"}, err: nil})

				mockScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			}
		})

		It("returns an error if sequences are added after returnResults is called", func() {
			entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				result := results{sResults: []string{"red", "yellow"}, err: io.EOF}
				rec := entry.newRecorder([]string{"Hello", "there"}, false)
				rec.returnResults(result.sResults, result.err)

				msg := fmt.Sprintf("%s must be called before %s or %s calls, parameters: %%#v",
					export("seq", entry),
					export("returnResults", entry),
					export("doReturnResults", entry))
				bParams := entry.bundleParams([]string{"Hello", "there"}, false)
				moqTMock.OnCall().Fatalf(msg, bParams).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, bParams)

				// ACT
				rec.seq()

				// ASSERT
				entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false, result)
				Expect(rec.isNil()).To(BeTrue())
				if entry.tracksParams() {
					Expect(fmtMsg).To(ContainSubstring("Hello"))
				}

				mockScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			}
		})

		It("returns an error if sequences are removed after returnResults is called", func() {
			entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				result := results{sResults: []string{"red", "yellow"}, err: io.EOF}
				rec := entry.newRecorder([]string{"Hello", "there"}, false)
				rec.returnResults(result.sResults, result.err)

				msg := fmt.Sprintf("%s must be called before %s or %s calls, parameters: %%#v",
					export("noSeq", entry),
					export("returnResults", entry),
					export("doReturnResults", entry))
				bParams := entry.bundleParams([]string{"Hello", "there"}, false)
				moqTMock.OnCall().Fatalf(msg, bParams).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, bParams)

				// ACT
				rec.noSeq()

				// ASSERT
				entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false, result)
				Expect(rec.isNil()).To(BeTrue())
				if entry.tracksParams() {
					Expect(fmtMsg).To(ContainSubstring("Hello"))
				}

				mockScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			}
		})

		It("can have sequences turned off for select calls", func() {
			entries, mockScene := tableEntries(
				moqTMock, moq.MockConfig{Sequence: moq.SeqDefaultOn})
			for _, entry := range entries {
				// ASSEMBLE
				if !entry.tracksParams() {
					// With no params to track, hard to make conflicting calls
					continue
				}
				scene.Reset()
				mockScene.Reset()

				rec := entry.newRecorder([]string{"Eh", "you"}, true)
				rec.noSeq()
				rec.returnResults([]string{"grey", "brown"}, nil)

				expectCall(entry, []string{"Hello", "there"}, false, results{
					sResults: []string{"red", "yellow"}, err: io.EOF})

				expectCall(entry, []string{"Hi", "you"}, true, results{
					sResults: []string{"blue", "orange"}, err: nil})

				rec = entry.newRecorder([]string{"Bye", "now"}, true)
				rec.noSeq()
				rec.returnResults([]string{"silver", "black"}, nil)

				entry.invokeMockAndExpectResults([]string{"Bye", "now"}, true,
					results{sResults: []string{"silver", "black"}, err: nil})

				// ACT
				entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ASSERT
				entry.invokeMockAndExpectResults([]string{"Eh", "you"}, true,
					results{sResults: []string{"grey", "brown"}, err: nil})

				mockScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			}
		})

		It("reserves multiple sequences when times is called", func() {
			entries, mockScene := tableEntries(
				moqTMock, moq.MockConfig{Sequence: moq.SeqDefaultOn})
			for _, entry := range entries {
				// ASSEMBLE
				if !entry.tracksParams() {
					// With no params to track, hard to make conflicting calls
					continue
				}
				scene.Reset()
				mockScene.Reset()

				rec := entry.newRecorder([]string{"Hello", "there"}, false)
				rec.returnResults([]string{"red", "yellow"}, io.EOF)
				rec.times(2)

				// ACT
				entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
				entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})

				// ASSERT
				mockScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			}
		})

		It("allows sequences to work with anyTimes", func() {
			entries, mockScene := tableEntries(
				moqTMock, moq.MockConfig{Sequence: moq.SeqDefaultOn})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				rec := entry.newRecorder([]string{"Hello", "there"}, false)
				rec.returnResults([]string{"red", "yellow"}, io.EOF)
				rec.anyTimes()

				// ACT
				entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
				entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
				entry.invokeMockAndExpectResults([]string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})

				// ASSERT
				mockScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			}
		})
	})

	Describe("do functions", func() {
		It("can call different do functions when configured to do so", func() {
			entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				var firstCall bool
				rec.andDo(func() {
					firstCall = true
				}, []string{"Hi", "you"}, true)
				rec.returnResults([]string{"green", "purple"}, nil)
				var secondCall bool
				rec.andDo(func() {
					secondCall = true
				}, []string{"Hi", "you"}, true)

				// ACT
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				Expect(firstCall).To(BeTrue())
				Expect(secondCall).To(BeFalse())
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})

				// ASSERT
				Expect(secondCall).To(BeTrue())
				mockScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			}
		})

		It("fails if andDo is called without calling returnResults first", func() {
			entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				msg := fmt.Sprintf("%s must be called before calling %s",
					export("returnResults", entry), export("andDo", entry))
				moqTMock.OnCall().Fatalf(msg).
					ReturnResults()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)

				// ACT
				rec.andDo(func() {}, []string{"Hi", "you"}, true)

				// ASSERT
				Expect(rec.isNil()).To(BeTrue())
				mockScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			}
		})

		It("can derive return values from doReturnResults functions when configured to do so", func() {
			entries, mockScene := tableEntries(moqTMock, moq.MockConfig{})
			for _, entry := range entries {
				// ASSEMBLE
				scene.Reset()
				mockScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				var firstCall bool
				rec.doReturnResults(func() {
					firstCall = true
				}, []string{"Hi", "you"}, true, []string{"blue", "orange"}, nil)
				var secondCall bool
				rec.doReturnResults(func() {
					secondCall = true
				}, []string{"Hi", "you"}, true, []string{"green", "purple"}, nil)

				// ACT
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				Expect(firstCall).To(BeTrue())
				Expect(secondCall).To(BeFalse())
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})

				// ASSERT
				Expect(secondCall).To(BeTrue())
				mockScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			}
		})

		It("fails when sequences are incorrect with a doReturnResults function", func() {
			entries, mockScene := tableEntries(
				moqTMock, moq.MockConfig{Sequence: moq.SeqDefaultOn})
			for _, entry := range entries {
				// ASSEMBLE
				if !entry.tracksParams() {
					// With no params to track, hard to make conflicting calls
					continue
				}
				scene.Reset()
				mockScene.Reset()

				rec := entry.newRecorder([]string{"Hello", "there"}, false)
				rec.doReturnResults(
					func() {}, []string{"Hello", "there"}, false, []string{"red", "yellow"}, io.EOF)
				rec = entry.newRecorder([]string{"Hi", "you"}, true)
				rec.doReturnResults(
					func() {}, []string{"Hi", "you"}, true, []string{"blue", "orange"}, nil)

				msg := "Call sequence does not match %#v"
				params1 := entry.bundleParams([]string{"Hi", "you"}, true)
				moqTMock.OnCall().Fatalf(msg, params1).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, params1)
				params2 := entry.bundleParams([]string{"Hello", "there"}, false)
				moqTMock.OnCall().Errorf("Expected %d additional call(s) with parameters %#v", 1, params2).
					ReturnResults()

				// ACT
				entry.invokeMockAndExpectResults([]string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ASSERT
				Expect(fmtMsg).To(ContainSubstring("Hi"))

				mockScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
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
		filePath := "./exported/moq_usualfn.go"
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

func export(id string, a adaptor) string {
	if a.exported() {
		id = strings.Title(id)
	}
	return id
}
