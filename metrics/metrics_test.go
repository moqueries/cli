package metrics_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/myshkin5/moqueries/metrics"
	"github.com/myshkin5/moqueries/moq"
)

func TestMetrics(t *testing.T) {
	var (
		scene          *moq.Scene
		isLoggingFnMoq *moqIsLoggingFn
		loggingfFnMoq  *moqLoggingfFn

		m *metrics.Processor
	)

	beforeEach := func(t *testing.T) {
		t.Helper()
		if scene != nil {
			t.Fatalf("afterEach not called")
		}
		scene = moq.NewScene(t)
		config := &moq.Config{Sequence: moq.SeqDefaultOn}
		isLoggingFnMoq = newMoqIsLoggingFn(scene, config)
		loggingfFnMoq = newMoqLoggingfFn(scene, config)

		m = metrics.NewMetrics(isLoggingFnMoq.mock(), loggingfFnMoq.mock())
	}

	afterEach := func() {
		scene.AssertExpectationsMet()
		scene = nil
	}

	t.Run("no-op", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()

		m.ASTPkgCacheHitsInc()
		m.ASTPkgCacheMissesInc()
		m.ASTTypeCacheHitsInc()
		m.ASTTypeCacheMissesInc()
		m.ASTTotalLoadTimeInc(1234 * time.Millisecond)
		m.ASTTotalDecorationTimeInc(9999 * time.Millisecond)
		m.TotalProcessingTimeInc(4321 * time.Millisecond)

		isLoggingFnMoq.onCall().returnResults(false)

		// ACT
		m.Finalize()

		// ASSERT
	})

	t.Run("single int increment", func(t *testing.T) {
		incFuncs := map[string]func(m *metrics.Processor){
			"ASTPkgCacheHits":    func(m *metrics.Processor) { m.ASTPkgCacheHitsInc() },
			"ASTPkgCacheMisses":  func(m *metrics.Processor) { m.ASTPkgCacheMissesInc() },
			"ASTTypeCacheHits":   func(m *metrics.Processor) { m.ASTTypeCacheHitsInc() },
			"ASTTypeCacheMisses": func(m *metrics.Processor) { m.ASTTypeCacheMissesInc() },
		}

		for name, incFunc := range incFuncs {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				beforeEach(t)
				defer afterEach()

				isLoggingFnMoq.onCall().returnResults(true)
				var data interface{}
				loggingfFnMoq.onCall("Type cache metrics %s", nil).
					any().args().
					doReturnResults(func(_ string, args ...interface{}) {
						if len(args) != 1 {
							t.Fatalf("got %d args, wanted 1", len(args))
						}
						data = args[0]
					})

				// ACT
				incFunc(m)
				m.Finalize()

				// ASSERT
				buf, ok := data.([]byte)
				if !ok {
					t.Fatalf("got %#v, wanted []byte", data)
				}
				msi := map[string]interface{}{}
				err := json.Unmarshal(buf, &msi)
				if err != nil {
					t.Fatalf("got %#v, wanted no err", err)
				}

				found := false
				for k, v := range msi {
					var i float64
					if k == name {
						found = true
						i = 1
					} else {
						i = 0
					}
					switch typ := v.(type) {
					case float64:
						if typ != i {
							t.Errorf("got %f, wanted %f", typ, i)
						}
					case string:
						if typ != "0s" {
							t.Errorf("got %s, wanted 0s", typ)
						}
					default:
						t.Errorf("got %#v, wanted float64 or string", typ)
					}
				}

				if !found {
					t.Errorf("got not found, wanted found")
				}
			})
		}
	})

	t.Run("single duration increment", func(t *testing.T) {
		incFuncs := map[string]func(m *metrics.Processor, d time.Duration){
			"ASTTotalLoadTime": func(m *metrics.Processor, d time.Duration) {
				m.ASTTotalLoadTimeInc(d)
			},
			"ASTTotalDecorationTime": func(m *metrics.Processor, d time.Duration) {
				m.ASTTotalDecorationTimeInc(d)
			},
			"TotalProcessingTime": func(m *metrics.Processor, d time.Duration) {
				m.TotalProcessingTimeInc(d)
			},
		}

		for name, incFunc := range incFuncs {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				beforeEach(t)
				defer afterEach()

				isLoggingFnMoq.onCall().returnResults(true)
				var data interface{}
				loggingfFnMoq.onCall("Type cache metrics %s", nil).
					any().args().
					doReturnResults(func(_ string, args ...interface{}) {
						if len(args) != 1 {
							t.Fatalf("got %d args, wanted 1", len(args))
						}
						data = args[0]
					})

				// ACT
				incFunc(m, 123*time.Millisecond)
				m.Finalize()

				// ASSERT
				buf, ok := data.([]byte)
				if !ok {
					t.Fatalf("got %#v, wanted []byte", data)
				}
				msi := map[string]interface{}{}
				err := json.Unmarshal(buf, &msi)
				if err != nil {
					t.Fatalf("got %#v, wanted no err", err)
				}

				found := false
				foundStr := false
				for k, v := range msi {
					if k == name {
						found = true
						lTime := float64(123 * time.Millisecond)
						if v != lTime {
							t.Errorf("got %#v, wanted %f", v, lTime)
						}
					}
					if k == name+"Str" {
						foundStr = true
						lTime := "123ms"
						if v != lTime {
							t.Errorf("got %#v, wanted %s", v, lTime)
						}
					}
				}
				if !found {
					t.Errorf("got not found, wanted found")
				}
				if !foundStr {
					t.Errorf("got not foundStr, wanted foundStr")
				}
			})
		}
	})
}
