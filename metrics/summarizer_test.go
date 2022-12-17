package metrics_test

import (
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"moqueries.org/cli/metrics"
	"moqueries.org/cli/moq"
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
				"ASTPkgCacheHits":           0.0,
				"ASTPkgCacheMisses":         0.0,
				"ASTTypeCacheHits":          0.0,
				"ASTTypeCacheMisses":        0.0,
				"ASTTotalLoadTime":          0.0,
				"ASTTotalLoadTimeStr":       "0s",
				"ASTTotalDecorationTime":    0.0,
				"ASTTotalDecorationTimeStr": "0s",
				"TotalProcessingTime":       0.0,
				"TotalProcessingTimeStr":    "0s",
			},
		},
		"sparse": {
			lines: []string{
				"Hello 1",
				"Type cache metrics bad line",
				"Hello 2",
				"DEBUG: 2022/06/16 14:39:10 Type cache metrics {\"ASTPkgCacheHits\":1,\"ASTTotalLoadTime\":145394479," +
					"\"ASTTotalLoadTimeStr\":\"999.999ms\"}",
				"Hello 3",
				"Type cache metrics {\"ASTPkgCacheMisses\":2,\"ASTTotalLoadTime\":42,\"ASTTotalLoadTimeStr\":\"1234ms\"}",
				"hello 4",
				"Type cache metrics {\"ASTTypeCacheMisses\":1,\"TotalProcessingTime\":394392," +
					"\"TotalProcessingTimeStr\":\"ignored\",\"ASTTotalDecorationTime\":12345}",
				"hello 5",
				"Type cache metrics {\"ASTTypeCacheMisses\":5,\"TotalProcessingTime\":8374583}",
				"hello 6",
				"Type cache metrics {\"ASTTypeCacheMisses\":0,\"ASTTotalDecorationTime\":54321," +
					"\"ASTTotalDecorationTimeStr\":\"ignored\"}",
			},
			expectedRuns: 5,
			expectedState: map[string]interface{}{
				"ASTPkgCacheHits":           1.0,
				"ASTPkgCacheMisses":         2.0,
				"ASTTypeCacheHits":          0.0,
				"ASTTypeCacheMisses":        6.0,
				"ASTTotalLoadTime":          145394521.0,
				"ASTTotalLoadTimeStr":       "145.394521ms",
				"ASTTotalDecorationTime":    66666.0,
				"ASTTotalDecorationTimeStr": "66.666µs",
				"TotalProcessingTime":       8768975.0,
				"TotalProcessingTimeStr":    "8.768975ms",
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
