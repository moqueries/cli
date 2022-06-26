package metrics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/myshkin5/moqueries/logs"
)

//go:generate moqueries --import io Reader

// Summarize reads a log file, selects multiple metrics from multiple generate
// calls, and summarizes the metrics
func Summarize(f io.Reader, loggingfFn LoggingfFn) error {
	var totals metricsState
	scanner := bufio.NewScanner(f)
	runs := 0
	for scanner.Scan() {
		txt := scanner.Text()
		idx := strings.Index(txt, metricsLogKey)
		if idx == -1 {
			continue
		}

		jTxt := txt[idx+len(metricsLogKey)+1:]
		var metrics metricsState
		err := json.Unmarshal([]byte(jTxt), &metrics)
		if err != nil {
			logs.Errorf("Ignoring line, error JSON parsing %s: %#v", jTxt, err)
			continue
		}
		totals = add(totals, metrics)
		runs++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %w", err)
	}

	loggingfFn("Scanned %d runs, summary: %s", runs, serializeState(totals))

	return nil
}
