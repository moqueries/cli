package testmoqs_test

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"

	"github.com/myshkin5/moqueries/generator"
	"github.com/myshkin5/moqueries/generator/testmoqs"
	"github.com/myshkin5/moqueries/generator/testmoqs/exported"
	"github.com/myshkin5/moqueries/moq"
)

type results struct {
	sResults []string
	err      error
}

type adaptorConfig struct {
	// exported is true when the generated test code is exported
	exported bool
	// noParams is true when the function being tested has no parameters
	noParams bool
	// opaqueParams is true when the function being tested doesn't show
	// parameter content when an error occurs (params might be contained in
	// a separate struct)
	opaqueParams bool
}

type adaptor interface {
	config() adaptorConfig
	mock() interface{}
	newRecorder(sParams []string, bParam bool) recorder
	invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results)
	bundleParams(sParams []string, bParam bool) interface{}
	sceneMoq() moq.Moq
}

type recorder interface {
	anySParam()
	anyBParam()
	seq()
	noSeq()
	returnResults(sResults []string, err error)
	andDo(t moq.T, fn func(), sParams []string, bParam bool)
	doReturnResults(t moq.T, fn func(), sParams []string, bParam bool, sResults []string, err error)
	repeat(repeaters ...moq.Repeater)
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
	tMoq     *moq.MoqT
	moqScene *moq.Scene
)

func testCases(t *testing.T, c moq.Config) map[string]adaptor {
	t.Helper()
	scene = moq.NewScene(t)
	tMoq = moq.NewMoqT(scene, nil)
	moqScene = moq.NewScene(tMoq.Mock())

	usualMoq := newMoqUsual(moqScene, &c)
	exportUsualMoq := exported.NewMoqUsual(moqScene, &c)
	//nolint:lll // chopped down entries are reverted by gofumpt
	entries := map[string]adaptor{
		"usual fn":                           &usualFnAdaptor{m: newMoqUsualFn(moqScene, &c)},
		"exported usual fn":                  &exportedUsualFnAdaptor{m: exported.NewMoqUsualFn(moqScene, &c)},
		"no names fn":                        &noNamesFnAdaptor{m: newMoqNoNamesFn(moqScene, &c)},
		"exported no names fn":               &exportedNoNamesFnAdaptor{m: exported.NewMoqNoNamesFn(moqScene, &c)},
		"no results fn":                      &noResultsFnAdaptor{m: newMoqNoResultsFn(moqScene, &c)},
		"exported no results fn":             &exportedNoResultsFnAdaptor{m: exported.NewMoqNoResultsFn(moqScene, &c)},
		"no params fn":                       &noParamsFnAdaptor{m: newMoqNoParamsFn(moqScene, &c)},
		"exported no params fn":              &exportedNoParamsFnAdaptor{m: exported.NewMoqNoParamsFn(moqScene, &c)},
		"nothing fn":                         &nothingFnAdaptor{m: newMoqNothingFn(moqScene, &c)},
		"exported nothing fn":                &exportedNothingFnAdaptor{m: exported.NewMoqNothingFn(moqScene, &c)},
		"variadic fn":                        &variadicFnAdaptor{m: newMoqVariadicFn(moqScene, &c)},
		"exported variadic fn":               &exportedVariadicFnAdaptor{m: exported.NewMoqVariadicFn(moqScene, &c)},
		"repeated ids fn":                    &repeatedIdsFnAdaptor{m: newMoqRepeatedIdsFn(moqScene, &c)},
		"exported repeated ids fn":           &exportedRepeatedIdsFnAdaptor{m: exported.NewMoqRepeatedIdsFn(moqScene, &c)},
		"times fn":                           &timesFnAdaptor{m: newMoqTimesFn(moqScene, &c)},
		"exported times fn":                  &exportedTimesFnAdaptor{m: exported.NewMoqTimesFn(moqScene, &c)},
		"difficult param names fn":           &difficultParamNamesFnAdaptor{m: newMoqDifficultParamNamesFn(moqScene, &c)},
		"exported difficult param names fn":  &exportedDifficultParamNamesFnAdaptor{m: exported.NewMoqDifficultParamNamesFn(moqScene, &c)},
		"difficult result names fn":          &difficultResultNamesFnAdaptor{m: newMoqDifficultResultNamesFn(moqScene, &c)},
		"exported difficult result names fn": &exportedDifficultResultNamesFnAdaptor{m: exported.NewMoqDifficultResultNamesFn(moqScene, &c)},
		"pass by ref fn":                     &passByReferenceFnAdaptor{m: newMoqPassByReferenceFn(moqScene, &c)},
		"exported pass by ref fn":            &exportedPassByReferenceFnAdaptor{m: exported.NewMoqPassByReferenceFn(moqScene, &c)},
		"interface param fn":                 &interfaceParamFnAdaptor{m: newMoqInterfaceParamFn(moqScene, &c)},
		"exported interface param fn":        &exportedInterfaceParamFnAdaptor{m: exported.NewMoqInterfaceParamFn(moqScene, &c)},
		"interface result fn":                &interfaceResultFnAdaptor{m: newMoqInterfaceResultFn(moqScene, &c)},
		"exported interface result param fn": &exportedInterfaceResultFnAdaptor{m: exported.NewMoqInterfaceResultFn(moqScene, &c)},

		"usual":                           &usualAdaptor{m: usualMoq},
		"exported usual":                  &exportedUsualAdaptor{m: exportUsualMoq},
		"no names":                        &noNamesAdaptor{m: usualMoq},
		"exported no names":               &exportedNoNamesAdaptor{m: exportUsualMoq},
		"no results":                      &noResultsAdaptor{m: usualMoq},
		"exported no results":             &exportedNoResultsAdaptor{m: exportUsualMoq},
		"no params":                       &noParamsAdaptor{m: usualMoq},
		"exported no params":              &exportedNoParamsAdaptor{m: exportUsualMoq},
		"nothing":                         &nothingAdaptor{m: usualMoq},
		"exported nothing":                &exportedNothingAdaptor{m: exportUsualMoq},
		"variadic":                        &variadicAdaptor{m: usualMoq},
		"exported variadic":               &exportedVariadicAdaptor{m: exportUsualMoq},
		"repeated ids":                    &repeatedIdsAdaptor{m: usualMoq},
		"exported repeat ids":             &exportedRepeatedIdsAdaptor{m: exportUsualMoq},
		"times":                           &timesAdaptor{m: usualMoq},
		"exported times":                  &exportedTimesAdaptor{m: exportUsualMoq},
		"difficult param names":           &difficultParamNamesAdaptor{m: usualMoq},
		"exported difficult param names":  &exportedDifficultParamNamesAdaptor{m: exportUsualMoq},
		"difficult result names":          &difficultResultNamesAdaptor{m: usualMoq},
		"exported difficult result names": &exportedDifficultResultNamesAdaptor{m: exportUsualMoq},
		"pass by ref":                     &passByReferenceAdaptor{m: usualMoq},
		"exported pass by ref":            &exportedPassByReferenceAdaptor{m: exportUsualMoq},
		"interface param":                 &interfaceParamAdaptor{m: usualMoq},
		"exported interface param":        &exportedInterfaceParamAdaptor{m: exportUsualMoq},
		"interface result":                &interfaceResultAdaptor{m: usualMoq},
		"exported interface result param": &exportedInterfaceResultAdaptor{m: exportUsualMoq},
	}

	return entries
}

func TestMoqs(t *testing.T) {
	t.Run("can return different values when configured to do so", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				config := entry.config()
				if !config.noParams {
					expectCall(entry, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				expectCall(entry, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil},
					results{sResults: []string{"green", "purple"}, err: nil})
				if !config.noParams {
					expectCall(entry, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})

				// ASSERT
				moqScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			})
		}
	})

	t.Run("adds additional results when an expectation has already been set", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				expectCall(entry, []string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})

				expectCall(entry, []string{"Hi", "you"}, true,
					results{sResults: []string{"red", "blue"}, err: nil})

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"red", "blue"}, err: nil})

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})
}

func TestExpectations(t *testing.T) {
	t.Run("fails if an expectation is not set in strict mode", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				msg := "Unexpected call with parameters %#v"
				params := entry.bundleParams([]string{"Hi", "you"}, true)
				tMoq.OnCall().Fatalf(msg, params).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, params)

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"", ""}, err: nil})

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
				config := entry.config()
				if !config.noParams && !config.opaqueParams {
					if !strings.Contains(fmtMsg, "Hi") {
						t.Errorf("got: %s, want to contain Hi", fmtMsg)
					}
				}
			})
		}
	})

	t.Run("returns zero values if an expectation is not set in nice mode", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{Expectation: moq.Nice}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"", ""}, err: nil})

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})
}

func TestRepeaters(t *testing.T) {
	t.Run("can return the same values multiple times", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				config := entry.config()
				if !config.noParams {
					expectCall(entry, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.repeat(moq.Times(4))
				rec.returnResults([]string{"green", "purple"}, nil)

				if !config.noParams {
					expectCall(entry, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				for n := 0; n < 4; n++ {
					entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
						results{sResults: []string{"blue", "orange"}, err: nil})
				}
				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})

	t.Run("returns the same value any number of times", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				config := entry.config()
				if !config.noParams {
					expectCall(entry, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.returnResults([]string{"green", "purple"}, nil)
				rec.repeat(moq.AnyTimes())

				if !config.noParams {
					expectCall(entry, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				// ACT
				for n := 0; n < 20; n++ {
					entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
						results{sResults: []string{"green", "purple"}, err: nil})
				}

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})

	t.Run("fails if Repeat is called without a preceding Return call", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				msg := fmt.Sprintf("%s or %s must be called before calling %s",
					export("returnResults", entry),
					export("doReturnResults", entry),
					export("repeat", entry))

				tMoq.OnCall().Fatalf(msg).
					ReturnResults()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)

				// ACT
				rec.repeat(moq.Times(4))

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})

	t.Run("fails if the function is called too many times in strict mode", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				expectCall(entry, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil},
					results{sResults: []string{"green", "purple"}, err: io.EOF})

				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: io.EOF})

				tMoq.OnCall().Fatalf(
					"Too many calls to mock with parameters %#v",
					entry.bundleParams([]string{"Hi", "you"}, true)).ReturnResults()

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"", ""}, err: nil})

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})

	t.Run("returns zero values if the function is called too many times in nice mode", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{Expectation: moq.Nice}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				expectCall(entry, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil},
					results{sResults: []string{"green", "purple"}, err: io.EOF})

				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: io.EOF})

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"", ""}, err: nil})

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})
}

func TestReset(t *testing.T) {
	for name, entry := range testCases(t, moq.Config{}) {
		t.Run(name, func(t *testing.T) {
			// ASSEMBLE
			scene.Reset()
			moqScene.Reset()

			expectCall(entry, []string{"Hi", "you"}, true,
				results{sResults: []string{"green", "purple"}, err: nil})

			// ACT
			entry.sceneMoq().Reset()

			// ASSERT
			expectCall(entry, []string{"Hi", "you"}, true,
				results{sResults: []string{"grey", "indigo"}, err: nil})
			entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
				results{sResults: []string{"grey", "indigo"}, err: nil})

			scene.AssertExpectationsMet()
			moqScene.AssertExpectationsMet()
		})
	}
}

func TestAnyValues(t *testing.T) {
	t.Run("ignores a param when instructed to do so", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				config := entry.config()
				if !config.noParams {
					expectCall(entry, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.anySParam()
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.returnResults([]string{"green", "purple"}, nil)

				if !config.noParams {
					expectCall(entry, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Goodbye", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})

	t.Run("returns an error if any functions are called after returning results", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				config := entry.config()
				if config.noParams {
					t.Skip("With no params to track, there will be no `any` functions")
				}
				scene.Reset()
				moqScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				rrFn := "returnResults"
				drrFn := "doReturnResults"
				if config.exported {
					rrFn = strings.Title(rrFn)
					drrFn = strings.Title(drrFn)
				}

				msg := fmt.Sprintf(
					"Any functions must be called before %s or %s calls, parameters: %%#v",
					rrFn, drrFn)
				bParams := entry.bundleParams([]string{"Hi", "you"}, true)
				tMoq.OnCall().Fatalf(msg, bParams).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, bParams)

				// ACT
				rec.anySParam()

				// ASSERT
				if !rec.isNil() {
					t.Errorf("got: %t, want true (nil)", rec.isNil())
				}
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
				if !config.noParams && !config.opaqueParams {
					if !strings.Contains(fmtMsg, "Hi") {
						t.Errorf("got: %s, want to contain Hi", fmtMsg)
					}
				}
			})
		}
	})
}

func TestAssertExpectationsMet(t *testing.T) {
	t.Run("reports success when there ae no expectations", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				// ACT
				entry.sceneMoq().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})

	t.Run("reports success when all expectations are met", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				config := entry.config()
				if !config.noParams {
					expectCall(entry, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				expectCall(entry, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if !config.noParams {
					expectCall(entry, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				// ACT
				entry.sceneMoq().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})

	t.Run("fails when one expectation is not met", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				// ASSEMBLE
				config := entry.config()
				if !config.noParams {
					expectCall(entry, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				expectCall(entry, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if !config.noParams {
					expectCall(entry, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				msg := "Expected %d additional call(s) with parameters %#v"
				params := entry.bundleParams([]string{"Hi", "you"}, true)
				tMoq.OnCall().Errorf(msg, 1, params).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, 1, params)

				// ACT
				entry.sceneMoq().AssertExpectationsMet()

				// ASSERT
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
				if !config.noParams && !config.opaqueParams {
					if !strings.Contains(fmtMsg, "Hi") {
						t.Errorf("got: %s, want to contain Hi", fmtMsg)
					}
				}
			})
		}
	})

	t.Run("succeeds when an anyTimes expectation is not called at all", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.repeat(moq.AnyTimes())

				// ACT
				entry.sceneMoq().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})

	t.Run("succeeds when an anyTimes expectation is called once", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.repeat(moq.AnyTimes())

				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ACT
				entry.sceneMoq().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})

	t.Run("succeeds when an anyTimes expectation is called many times", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				// ASSEMBLE
				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.repeat(moq.AnyTimes())

				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ACT
				entry.sceneMoq().AssertExpectationsMet()

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})
}

func TestSequences(t *testing.T) {
	t.Run("passes when sequences are correct", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				config := entry.config()
				if !config.noParams {
					expectCall(entry, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}
				expectCall(entry, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil},
					results{sResults: []string{"green", "purple"}, err: nil})
				if !config.noParams {
					expectCall(entry, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				// ACT
				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})

				if !config.noParams {
					entry.invokeMockAndExpectResults(t, []string{"Bye", "now"}, true,
						results{sResults: []string{"red", "yellow"}, err: io.EOF})
				}

				// ASSERT
				moqScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			})
		}
	})

	t.Run("fails when sequences are incorrect", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{Sequence: moq.SeqDefaultOn}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				config := entry.config()
				if config.noParams {
					t.Skip("With no params to track, hard to make conflicting calls")
				}
				scene.Reset()
				moqScene.Reset()

				expectCall(entry, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
				expectCall(entry, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				msg := "Call sequence does not match %#v"
				params1 := entry.bundleParams([]string{"Hi", "you"}, true)
				tMoq.OnCall().Fatalf(msg, params1).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, params1)
				params2 := entry.bundleParams([]string{"Hello", "there"}, false)
				tMoq.OnCall().Errorf("Expected %d additional call(s) with parameters %#v", 1, params2).
					ReturnResults()

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ASSERT
				if !config.opaqueParams {
					if !strings.Contains(fmtMsg, "Hi") {
						t.Errorf("got: %s, want to contain Hi", fmtMsg)
					}
				}

				moqScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			})
		}
	})

	t.Run("can have sequences turned on for select calls", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{Sequence: moq.SeqDefaultOff}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				if entry.config().noParams {
					t.Skip("With no params to track, hard to make conflicting calls")
				}
				scene.Reset()
				moqScene.Reset()

				expectCall(entry, []string{"Eh", "you"}, true, results{
					sResults: []string{"grey", "brown"}, err: nil,
				})

				rec := entry.newRecorder([]string{"Hello", "there"}, false)
				rec.seq()
				rec.returnResults([]string{"red", "yellow"}, io.EOF)

				rec = entry.newRecorder([]string{"Hi", "you"}, true)
				rec.seq()
				rec.returnResults([]string{"blue", "orange"}, nil)

				expectCall(entry, []string{"Bye", "now"}, true, results{
					sResults: []string{"silver", "black"}, err: nil,
				})

				entry.invokeMockAndExpectResults(t, []string{"Bye", "now"}, true,
					results{sResults: []string{"silver", "black"}, err: nil})

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ASSERT
				entry.invokeMockAndExpectResults(t, []string{"Eh", "you"}, true,
					results{sResults: []string{"grey", "brown"}, err: nil})

				moqScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			})
		}
	})

	t.Run("returns an error if sequences are altered after returnResults is called", func(t *testing.T) {
		for _, seqNoSeq := range []string{"seq", "noSeq"} {
			t.Run(seqNoSeq, func(t *testing.T) {
				for name, entry := range testCases(t, moq.Config{}) {
					t.Run(name, func(t *testing.T) {
						// ASSEMBLE
						scene.Reset()
						moqScene.Reset()

						result := results{sResults: []string{"red", "yellow"}, err: io.EOF}
						rec := entry.newRecorder([]string{"Hello", "there"}, false)
						rec.returnResults(result.sResults, result.err)

						msg := fmt.Sprintf("%s must be called before %s or %s calls, parameters: %%#v",
							export(seqNoSeq, entry),
							export("returnResults", entry),
							export("doReturnResults", entry))
						bParams := entry.bundleParams([]string{"Hello", "there"}, false)
						tMoq.OnCall().Fatalf(msg, bParams).
							ReturnResults()
						fmtMsg := fmt.Sprintf(msg, bParams)

						// ACT
						switch seqNoSeq {
						case "seq":
							rec.seq()
						case "noSeq":
							rec.noSeq()
						default:
							t.Fatalf(seqNoSeq)
						}

						// ASSERT
						entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false, result)
						if !rec.isNil() {
							t.Errorf("got: %t, want true (nil)", rec.isNil())
						}
						config := entry.config()
						if !config.noParams && !config.opaqueParams {
							if !strings.Contains(fmtMsg, "Hello") {
								t.Errorf("got: %s, want to contain Hello", fmtMsg)
							}
						}

						moqScene.AssertExpectationsMet()
						scene.AssertExpectationsMet()
					})
				}
			})
		}
	})

	t.Run("can have sequences turned off for select calls", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{Sequence: moq.SeqDefaultOn}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				if entry.config().noParams {
					t.Skip("With no params to track, hard to make conflicting calls")
				}
				scene.Reset()
				moqScene.Reset()

				rec := entry.newRecorder([]string{"Eh", "you"}, true)
				rec.noSeq()
				rec.returnResults([]string{"grey", "brown"}, nil)

				expectCall(entry, []string{"Hello", "there"}, false, results{
					sResults: []string{"red", "yellow"}, err: io.EOF,
				})

				expectCall(entry, []string{"Hi", "you"}, true, results{
					sResults: []string{"blue", "orange"}, err: nil,
				})

				rec = entry.newRecorder([]string{"Bye", "now"}, true)
				rec.noSeq()
				rec.returnResults([]string{"silver", "black"}, nil)

				entry.invokeMockAndExpectResults(t, []string{"Bye", "now"}, true,
					results{sResults: []string{"silver", "black"}, err: nil})

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ASSERT
				entry.invokeMockAndExpectResults(t, []string{"Eh", "you"}, true,
					results{sResults: []string{"grey", "brown"}, err: nil})

				moqScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			})
		}
	})

	t.Run("reserves multiple sequences when times is called", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{Sequence: moq.SeqDefaultOn}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				if entry.config().noParams {
					t.Skip("With no params to track, hard to make conflicting calls")
				}
				scene.Reset()
				moqScene.Reset()

				rec := entry.newRecorder([]string{"Hello", "there"}, false)
				rec.returnResults([]string{"red", "yellow"}, io.EOF)
				rec.repeat(moq.Times(2))

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
				entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})

				// ASSERT
				moqScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			})
		}
	})

	t.Run("allows sequences to work with anyTimes", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{Sequence: moq.SeqDefaultOn}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				rec := entry.newRecorder([]string{"Hello", "there"}, false)
				rec.returnResults([]string{"red", "yellow"}, io.EOF)
				rec.repeat(moq.AnyTimes())

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
				entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})
				entry.invokeMockAndExpectResults(t, []string{"Hello", "there"}, false,
					results{sResults: []string{"red", "yellow"}, err: io.EOF})

				// ASSERT
				moqScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			})
		}
	})
}

func TestDoFuncs(t *testing.T) {
	t.Run("can call different do functions when configured to do so", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				var firstCall bool
				rec.andDo(t, func() {
					firstCall = true
				}, []string{"Hi", "you"}, true)
				rec.returnResults([]string{"green", "purple"}, nil)
				var secondCall bool
				rec.andDo(t, func() {
					secondCall = true
				}, []string{"Hi", "you"}, true)

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if !firstCall {
					t.Errorf("got no call, want first call")
				}
				if secondCall {
					t.Errorf("got call, want no second call")
				}
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})

				// ASSERT
				if !secondCall {
					t.Errorf("got no call, want second call")
				}
				moqScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			})
		}
	})

	t.Run("fails if andDo is called without calling returnResults first", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				msg := fmt.Sprintf("%s must be called before calling %s",
					export("returnResults", entry), export("andDo", entry))
				tMoq.OnCall().Fatalf(msg).
					ReturnResults()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)

				// ACT
				rec.andDo(t, func() {}, []string{"Hi", "you"}, true)

				// ASSERT
				if !rec.isNil() {
					t.Errorf("got: %t, want true (nil)", rec.isNil())
				}
				moqScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			})
		}
	})

	t.Run("can derive return values from doReturnResults functions when configured to do so", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				var firstCall bool
				rec.doReturnResults(t, func() {
					firstCall = true
				}, []string{"Hi", "you"}, true, []string{"blue", "orange"}, nil)
				var secondCall bool
				rec.doReturnResults(t, func() {
					secondCall = true
				}, []string{"Hi", "you"}, true, []string{"green", "purple"}, nil)

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})
				if !firstCall {
					t.Errorf("got no call, want first call")
				}
				if secondCall {
					t.Errorf("got call, want no second call")
				}
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"green", "purple"}, err: nil})

				// ASSERT
				if !secondCall {
					t.Errorf("got no call, want second call")
				}
				moqScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			})
		}
	})

	t.Run("fails when sequences are incorrect with a doReturnResults function", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{Sequence: moq.SeqDefaultOn}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				config := entry.config()
				if config.noParams {
					t.Skip("With no params to track, hard to make conflicting calls")
				}
				scene.Reset()
				moqScene.Reset()

				rec := entry.newRecorder([]string{"Hello", "there"}, false)
				rec.doReturnResults(
					t, func() {}, []string{"Hello", "there"}, false, []string{"red", "yellow"}, io.EOF)
				rec = entry.newRecorder([]string{"Hi", "you"}, true)
				rec.doReturnResults(
					t, func() {}, []string{"Hi", "you"}, true, []string{"blue", "orange"}, nil)

				msg := "Call sequence does not match %#v"
				params1 := entry.bundleParams([]string{"Hi", "you"}, true)
				tMoq.OnCall().Fatalf(msg, params1).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, params1)
				params2 := entry.bundleParams([]string{"Hello", "there"}, false)
				tMoq.OnCall().Errorf("Expected %d additional call(s) with parameters %#v", 1, params2).
					ReturnResults()

				// ACT
				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				// ASSERT
				if !config.opaqueParams {
					if !strings.Contains(fmtMsg, "Hi") {
						t.Errorf("got: %s, want to contain Hi", fmtMsg)
					}
				}

				moqScene.AssertExpectationsMet()
				scene.AssertExpectationsMet()
			})
		}
	})
}

func TestOptionalInvocations(t *testing.T) {
	t.Run("success when optional invocations are not made", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.repeat(moq.MinTimes(2), moq.MaxTimes(4))

				// ACT
				for n := 0; n < 2; n++ {
					entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
						results{sResults: []string{"blue", "orange"}, err: nil})
				}

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})

	t.Run("success when optional invocations are made", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.repeat(moq.MinTimes(2), moq.MaxTimes(4))

				// ACT
				for n := 0; n < 4; n++ {
					entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
						results{sResults: []string{"blue", "orange"}, err: nil})
				}

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})

	t.Run("failure when fewer than min invocations are not made", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.repeat(moq.MinTimes(2), moq.MaxTimes(4))

				entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
					results{sResults: []string{"blue", "orange"}, err: nil})

				msg := "Expected %d additional call(s) with parameters %#v"
				params := entry.bundleParams([]string{"Hi", "you"}, true)
				tMoq.OnCall().Errorf(msg, 1, params).
					ReturnResults()
				fmtMsg := fmt.Sprintf(msg, 1, params)

				// ACT
				moqScene.AssertExpectationsMet()

				// ASSERT
				config := entry.config()
				if !config.noParams && !config.opaqueParams {
					if !strings.Contains(fmtMsg, "Hi") {
						t.Errorf("got: %s, want to contain Hi", fmtMsg)
					}
				}

				scene.AssertExpectationsMet()
			})
		}
	})

	t.Run("success with multiple independent identical expectations and just optional invocations", func(t *testing.T) {
		for name, entry := range testCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				scene.Reset()
				moqScene.Reset()

				rec := entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.repeat(moq.MinTimes(2), moq.MaxTimes(4))

				rec = entry.newRecorder([]string{"Hi", "you"}, true)
				rec.returnResults([]string{"blue", "orange"}, nil)
				rec.repeat(moq.MinTimes(2), moq.MaxTimes(4))

				// ACT
				for n := 0; n < 4; n++ {
					entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
						results{sResults: []string{"blue", "orange"}, err: nil})
				}

				// ASSERT
				scene.AssertExpectationsMet()
				moqScene.AssertExpectationsMet()
			})
		}
	})

	t.Run("failure with multiple independent identical expectations and less than min optional invocations",
		func(t *testing.T) {
			for name, entry := range testCases(t, moq.Config{}) {
				t.Run(name, func(t *testing.T) {
					// ASSEMBLE
					scene.Reset()
					moqScene.Reset()

					rec := entry.newRecorder([]string{"Hi", "you"}, true)
					rec.returnResults([]string{"blue", "orange"}, nil)
					rec.repeat(moq.MinTimes(2), moq.MaxTimes(4))

					rec = entry.newRecorder([]string{"Hi", "you"}, true)
					rec.returnResults([]string{"blue", "orange"}, nil)
					rec.repeat(moq.MinTimes(2), moq.MaxTimes(4))

					for n := 0; n < 3; n++ {
						entry.invokeMockAndExpectResults(t, []string{"Hi", "you"}, true,
							results{sResults: []string{"blue", "orange"}, err: nil})
					}

					msg := "Expected %d additional call(s) with parameters %#v"
					params := entry.bundleParams([]string{"Hi", "you"}, true)
					tMoq.OnCall().Errorf(msg, 1, params).
						ReturnResults()
					fmtMsg := fmt.Sprintf(msg, 1, params)

					// ACT
					moqScene.AssertExpectationsMet()

					// ASSERT
					config := entry.config()
					if !config.noParams && !config.opaqueParams {
						if !strings.Contains(fmtMsg, "Hi") {
							t.Errorf("got: %s, want to contain Hi", fmtMsg)
						}
					}
					scene.AssertExpectationsMet()
				})
			}
		})
}

func TestConsistentMockInstance(t *testing.T) {
	for name, entry := range testCases(t, moq.Config{}) {
		t.Run(name, func(t *testing.T) {
			// ASSEMBLE
			scene.Reset()
			moqScene.Reset()
			mockAddr1 := fmt.Sprintf("%v", entry.mock())

			// ACT
			mockAddr2 := fmt.Sprintf("%v", entry.mock())

			// ASSERT
			if mockAddr1 != mockAddr2 {
				t.Errorf("wanted same instance, got %s != %s", mockAddr1, mockAddr2)
			}
		})
	}
}

type paramIndexingAdaptor interface {
	setParamIndexing(pi moq.ParamIndexing)
	onCall(params *testmoqs.PassByReferenceParams, sResult string, err error)
	invokeMockAndExpectResults(t moq.T, params *testmoqs.PassByReferenceParams, sResult string, err error)
	bundleParams(params *testmoqs.PassByReferenceParams) interface{}
}

type passByReferenceFnParamIndexingAdaptor struct {
	m *moqPassByReferenceFn
}

func (a *passByReferenceFnParamIndexingAdaptor) setParamIndexing(pi moq.ParamIndexing) {
	a.m.runtime.parameterIndexing.p = pi
}

func (a *passByReferenceFnParamIndexingAdaptor) onCall(
	params *testmoqs.PassByReferenceParams, sResult string, err error) {
	a.m.onCall(params).returnResults(sResult, err)
}

func (a *passByReferenceFnParamIndexingAdaptor) invokeMockAndExpectResults(
	t moq.T, params *testmoqs.PassByReferenceParams, sResult string, err error) {
	actualSResult, actualErr := a.m.mock()(params)
	if sResult != actualSResult {
		t.Errorf("got %s, wanted %s", actualSResult, sResult)
	}
	if err != actualErr {
		t.Errorf("got %#v, wanted %#v", actualErr, err)
	}
}

func (a *passByReferenceFnParamIndexingAdaptor) bundleParams(params *testmoqs.PassByReferenceParams) interface{} {
	return moqPassByReferenceFn_params{p: params}
}

type exportedPassByReferenceParamIndexingFnAdaptor struct {
	m *exported.MoqPassByReferenceFn
}

func (a *exportedPassByReferenceParamIndexingFnAdaptor) setParamIndexing(pi moq.ParamIndexing) {
	a.m.Runtime.ParameterIndexing.P = pi
}

func (a *exportedPassByReferenceParamIndexingFnAdaptor) onCall(
	params *testmoqs.PassByReferenceParams, sResult string, err error) {
	a.m.OnCall(params).ReturnResults(sResult, err)
}

func (a *exportedPassByReferenceParamIndexingFnAdaptor) invokeMockAndExpectResults(
	t moq.T, params *testmoqs.PassByReferenceParams, sResult string, err error) {
	actualSResult, actualErr := a.m.Mock()(params)
	if sResult != actualSResult {
		t.Errorf("got %s, wanted %s", actualSResult, sResult)
	}
	if err != actualErr {
		t.Errorf("got %#v, wanted %#v", actualErr, err)
	}
}

func (a *exportedPassByReferenceParamIndexingFnAdaptor) bundleParams(
	params *testmoqs.PassByReferenceParams) interface{} {
	return exported.MoqPassByReferenceFn_params{P: params}
}

type passByReferenceParamIndexingAdaptor struct {
	m *moqUsual
}

func (a *passByReferenceParamIndexingAdaptor) setParamIndexing(pi moq.ParamIndexing) {
	a.m.runtime.parameterIndexing.PassByReference.p = pi
}

func (a *passByReferenceParamIndexingAdaptor) onCall(
	params *testmoqs.PassByReferenceParams, sResult string, err error) {
	a.m.onCall().PassByReference(params).returnResults(sResult, err)
}

func (a *passByReferenceParamIndexingAdaptor) invokeMockAndExpectResults(
	t moq.T, params *testmoqs.PassByReferenceParams, sResult string, err error) {
	actualSResult, actualErr := a.m.mock().PassByReference(params)
	if sResult != actualSResult {
		t.Errorf("got %s, wanted %s", actualSResult, sResult)
	}
	if err != actualErr {
		t.Errorf("got %#v, wanted %#v", actualErr, err)
	}
}

func (a *passByReferenceParamIndexingAdaptor) bundleParams(params *testmoqs.PassByReferenceParams) interface{} {
	return moqUsual_PassByReference_params{p: params}
}

type exportedPassByReferenceParamIndexingAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedPassByReferenceParamIndexingAdaptor) setParamIndexing(pi moq.ParamIndexing) {
	a.m.Runtime.ParameterIndexing.PassByReference.P = pi
}

func (a *exportedPassByReferenceParamIndexingAdaptor) onCall(
	params *testmoqs.PassByReferenceParams, sResult string, err error) {
	a.m.OnCall().PassByReference(params).ReturnResults(sResult, err)
}

func (a *exportedPassByReferenceParamIndexingAdaptor) invokeMockAndExpectResults(
	t moq.T, params *testmoqs.PassByReferenceParams, sResult string, err error) {
	actualSResult, actualErr := a.m.Mock().PassByReference(params)
	if sResult != actualSResult {
		t.Errorf("got %s, wanted %s", actualSResult, sResult)
	}
	if err != actualErr {
		t.Errorf("got %#v, wanted %#v", actualErr, err)
	}
}

func (a *exportedPassByReferenceParamIndexingAdaptor) bundleParams(params *testmoqs.PassByReferenceParams) interface{} {
	return exported.MoqUsual_PassByReference_params{P: params}
}

func paramIndexingTestCases(t *testing.T, c moq.Config) map[string]paramIndexingAdaptor {
	t.Helper()
	scene = moq.NewScene(t)
	tMoq = moq.NewMoqT(scene, nil)
	moqScene = moq.NewScene(tMoq.Mock())

	usualMoq := newMoqUsual(moqScene, &c)
	exportUsualMoq := exported.NewMoqUsual(moqScene, &c)
	entries := map[string]paramIndexingAdaptor{
		"pass by ref fn": &passByReferenceFnParamIndexingAdaptor{
			m: newMoqPassByReferenceFn(moqScene, &c),
		},
		"exported pass by ref fn": &exportedPassByReferenceParamIndexingFnAdaptor{
			m: exported.NewMoqPassByReferenceFn(moqScene, &c),
		},

		"pass by ref":          &passByReferenceParamIndexingAdaptor{m: usualMoq},
		"exported pass by ref": &exportedPassByReferenceParamIndexingAdaptor{m: exportUsualMoq},
	}

	return entries
}

type sliceWriter []int

func (w sliceWriter) Write([]byte) (n int, err error) {
	return 0, err
}

func TestParamIndexing(t *testing.T) {
	t.Run("can index a param by value", func(t *testing.T) {
		for name, entry := range paramIndexingTestCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				p := testmoqs.PassByReferenceParams{
					SParam: "Hi",
					BParam: true,
				}

				entry.setParamIndexing(moq.ParamIndexByValue)
				entry.onCall(&p, "Hello", nil)

				// If we are indexing by hash, this change would make the expectation
				// not match
				p.SParam = "Changed my mind, good bye!"
				p.BParam = false

				// ACT
				entry.invokeMockAndExpectResults(t, &p, "Hello", nil)

				// ASSERT
				scene.AssertExpectationsMet()
			})
		}
	})

	t.Run("errors when indexing params by value and the values are equal but different instances", func(t *testing.T) {
		for name, entry := range paramIndexingTestCases(t, moq.Config{}) {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				p1 := testmoqs.PassByReferenceParams{
					SParam: "Hi",
					BParam: true,
				}
				p2 := testmoqs.PassByReferenceParams{
					SParam: "Hi",
					BParam: true,
				}

				entry.setParamIndexing(moq.ParamIndexByValue)
				entry.onCall(&p1, "Hello", nil)

				params := entry.bundleParams(&p1)
				tMoq.OnCall().Fatalf("Unexpected call with parameters %#v", params).ReturnResults()

				// ACT
				entry.invokeMockAndExpectResults(t, &p2, "", nil)

				// ASSERT
				scene.AssertExpectationsMet()
			})
		}
	})

	t.Run("panics if an interface parameter that can't be indexed by value gets indexed by value", func(t *testing.T) {
		// ASSEMBLE
		var r interface{}
		{
			defer func() {
				r = recover()
			}()

			scene = moq.NewScene(t)
			usualMoq := newMoqUsual(scene, nil)
			usualMoq.runtime.parameterIndexing.InterfaceParam.w = moq.ParamIndexByValue
			w := sliceWriter{}

			// ACT
			usualMoq.onCall().InterfaceParam(w).returnResults("", nil)
		}

		// ASSERT
		if r == nil {
			t.Fatalf("wanted panic, got none")
		}

		err, ok := r.(error)
		if !ok {
			t.Fatalf("wanted error, got none")
		}

		msg := "hash of unhashable type testmoqs_test.sliceWriter"
		if err.Error() != msg {
			t.Errorf("wanted %s, got %s", msg, err.Error())
		}
	})
}

func TestGenerating(t *testing.T) {
	t.Run("generates moqs", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping generate test in short mode.")
		}

		// NB: Keep in sync with types.go go:generate directives

		// These lines generate the same moqs listed in types.go go:generate
		// directives. Remove the "pending" flag on this test to verify code
		// coverage.

		err := generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_usualfn_test.go", Types: []string{"UsualFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_usualfn.go", Export: true, Types: []string{"UsualFn"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_nonamesfn_test.go", Types: []string{"NoNamesFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_nonamesfn.go", Export: true, Types: []string{"NoNamesFn"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_noresultsfn_test.go", Types: []string{"NoResultsFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_noresultsfn.go", Export: true, Types: []string{"NoResultsFn"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_noparamsfn_test.go", Types: []string{"NoParamsFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_noparamsfn.go", Export: true, Types: []string{"NoParamsFn"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_nothingfn_test.go", Types: []string{"NothingFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_nothingfn.go", Export: true, Types: []string{"NothingFn"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_variadicfn_test.go", Types: []string{"VariadicFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_variadicfn.go", Export: true, Types: []string{"VariadicFn"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_repeatedidsfn_test.go", Types: []string{"RepeatedIdsFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_repeatedidsfn.go", Export: true, Types: []string{"RepeatedIdsFn"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_timesfn_test.go", Types: []string{"TimesFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_timesfn.go", Export: true, Types: []string{"TimesFn"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_difficultparamnamesfn_test.go", Types: []string{"DifficultParamNamesFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_difficultparamnamesfn.go", Export: true, Types: []string{"DifficultParamNamesFn"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_difficultresultnamesfn_test.go", Types: []string{"DifficultResultNamesFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_difficultresultnamesfn.go", Export: true, Types: []string{"DifficultResultNamesFn"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_passbyreferencefn_test.go", Types: []string{"PassByReferenceFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_passbyreferencefn.go", Export: true, Types: []string{"PassByReferenceFn"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_interfaceparamfn_test.go", Types: []string{"InterfaceParamFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_interfaceparamfn.go", Export: true, Types: []string{"InterfaceParamFn"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_interfaceresultfn_test.go", Types: []string{"InterfaceResultFn"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_interfaceresultfn.go", Export: true, Types: []string{"InterfaceResultFn"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = generator.Generate(
			generator.GenerateRequest{
				Destination: "moq_usual_test.go", Types: []string{"Usual"},
			},
			generator.GenerateRequest{
				Destination: "exported/moq_usual.go", Export: true, Types: []string{"Usual"},
			},
		)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}
	})

	t.Run("dumps the DST of a moq", func(t *testing.T) {
		filePath := "./moq_usual_test.go"
		outPath := "./moq_usual_test_dst.txt"

		fSet := token.NewFileSet()
		inFile, err := parser.ParseFile(fSet, filePath, nil, parser.ParseComments)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		dstFile, err := decorator.DecorateFile(fSet, inFile)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}

		err = dst.Fprint(outFile, dstFile, dst.NotNilFilter)
		if err != nil {
			t.Errorf("got %#v, wanted no err", err)
		}
	})
}

func export(id string, a adaptor) string {
	if a.config().exported {
		id = strings.Title(id)
	}
	return id
}
