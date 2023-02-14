package metrics_test

import (
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"moqueries.org/runtime/moq"

	"moqueries.org/cli/metrics"
)

func TestSummarize(t *testing.T) {
	type testCase struct {
		lines         []string
		expectedRuns  int
		expectedState map[string]interface{}
	}
	testCases := map[string]testCase{
		"nothing to summarize": {
			lines:        []string{"Hello 1", "Hello 2", "Hello 3"},
			expectedRuns: 0,
			expectedState: map[string]interface{}{
				"ast-pkg-cache-hits":            0.0,
				"ast-pkg-cache-misses":          0.0,
				"ast-type-cache-hits":           0.0,
				"ast-type-cache-misses":         0.0,
				"ast-total-load-time":           0.0,
				"ast-total-load-time-str":       "0s",
				"ast-total-decoration-time":     0.0,
				"ast-total-decoration-time-str": "0s",
				"total-processing-time":         0.0,
				"total-processing-time-str":     "0s",
			},
		},
		"sparse": {
			lines: []string{
				"Hello 1",
				"Type cache metrics bad line",
				"Hello 2",
				"DEBUG: 2022/06/16 14:39:10 Type cache metrics {\"ast-pkg-cache-hits\":1,\"ast-total-load-time\":145394479," +
					"\"ast-total-load-time-str\":\"999.999ms\"}",
				"Hello 3",
				"Type cache metrics {\"ast-pkg-cache-misses\":2,\"ast-total-load-time\":42,\"ast-total-load-time-str\":\"1234ms\"}",
				"hello 4",
				"Type cache metrics {\"ast-type-cache-misses\":1,\"total-processing-time\":394392," +
					"\"total-processing-time-str\":\"ignored\",\"ast-total-decoration-time\":12345}",
				"hello 5",
				"Type cache metrics {\"ast-type-cache-misses\":5,\"total-processing-time\":8374583}",
				"hello 6",
				"Type cache metrics {\"ast-type-cache-misses\":0,\"ast-total-decoration-time\":54321," +
					"\"ast-total-decoration-time-str\":\"ignored\"}",
			},
			expectedRuns: 5,
			expectedState: map[string]interface{}{
				"ast-pkg-cache-hits":            1.0,
				"ast-pkg-cache-misses":          2.0,
				"ast-type-cache-hits":           0.0,
				"ast-type-cache-misses":         6.0,
				"ast-total-load-time":           145394521.0,
				"ast-total-load-time-str":       "145.394521ms",
				"ast-total-decoration-time":     66666.0,
				"ast-total-decoration-time-str": "66.666µs",
				"total-processing-time":         8768975.0,
				"total-processing-time-str":     "8.768975ms",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// ASSEMBLE
			scene := moq.NewScene(t)
			config := &moq.Config{Sequence: moq.SeqDefaultOn}
			fMoq := newMoqReader(scene, config)
			loggingfFnMoq := newMoqLoggingfFn(scene, config)

			for _, line := range tc.lines {
				fn := func(line string) func(p []byte) (int, error) {
					return func(p []byte) (int, error) {
						line += "\n"
						copy(p, line)
						return len(line), nil
					}
				}
				fMoq.onCall().Read(nil).any().p().
					doReturnResults(fn(line))
			}
			fMoq.onCall().Read(nil).any().p().
				doReturnResults(func(p []byte) (int, error) {
					return 0, io.EOF
				})
			var actualArgs []interface{}
			loggingfFnMoq.onCall("Scanned %d runs, summary: %s", nil).
				any().args().
				doReturnResults(func(format string, args ...interface{}) {
					actualArgs = args
				})

			// ACT
			err := metrics.Summarize(fMoq.mock(), loggingfFnMoq.mock())
			// ASSERT
			if err != nil {
				t.Errorf("got %#v, wanted no error", err)
			}
			if len(actualArgs) != 2 {
				t.Fatalf("got %d args, wanted 2", len(actualArgs))
			}
			runs, ok := actualArgs[0].(int)
			if !ok {
				t.Errorf("got %#v, wanted int", runs)
			}
			if runs != tc.expectedRuns {
				t.Errorf("got %d runs, wanted %d", runs, tc.expectedRuns)
			}
			buf, ok := actualArgs[1].([]byte)
			if !ok {
				t.Errorf("got %#v, wanted []byte", buf)
			}
			msi := map[string]interface{}{}
			err = json.Unmarshal(buf, &msi)
			if err != nil {
				t.Fatalf("got %#v, wanted no err", err)
			}
			if !reflect.DeepEqual(msi, tc.expectedState) {
				t.Errorf("got %#v, wanted %#v", msi, tc.expectedState)
			}
		})
	}
}
