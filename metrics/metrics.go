// Package metrics implements a simple metrics mechanism for gathering,
// reporting and aggregating metrics
package metrics

import (
	"bytes"
	"encoding/json"
	"time"

	"moqueries.org/runtime/logs"
)

const (
	metricsLogKey = "Type cache metrics"
)

//go:generate moqueries --export Metrics

// Metrics is the interface to the metrics system
type Metrics interface {
	ASTPkgCacheHitsInc()
	ASTPkgCacheMissesInc()
	ASTTotalLoadTimeInc(d time.Duration)
	ASTTotalDecorationTimeInc(d time.Duration)

	TotalProcessingTimeInc(d time.Duration)

	Finalize()
}

type metricsState struct {
	ASTPkgCacheHits    int `json:"ast-pkg-cache-hits"`
	ASTPkgCacheMisses  int `json:"ast-pkg-cache-misses"`
	ASTTypeCacheHits   int `json:"ast-type-cache-hits"`
	ASTTypeCacheMisses int `json:"ast-type-cache-misses"`

	ASTTotalLoadTime    time.Duration `json:"ast-total-load-time"`
	ASTTotalLoadTimeStr string        `json:"ast-total-load-time-str"`

	ASTTotalDecorationTime    time.Duration `json:"ast-total-decoration-time"`
	ASTTotalDecorationTimeStr string        `json:"ast-total-decoration-time-str"`

	TotalProcessingTime    time.Duration `json:"total-processing-time"`
	TotalProcessingTimeStr string        `json:"total-processing-time-str"`
}

func add(m1, m2 metricsState) metricsState {
	// NOTE: Intentionally close to struct definition as a reminder to update
	//   this func when the struct is updated
	return metricsState{
		ASTPkgCacheHits:        m1.ASTPkgCacheHits + m2.ASTPkgCacheHits,
		ASTPkgCacheMisses:      m1.ASTPkgCacheMisses + m2.ASTPkgCacheMisses,
		ASTTypeCacheHits:       m1.ASTTypeCacheHits + m2.ASTTypeCacheHits,
		ASTTypeCacheMisses:     m1.ASTTypeCacheMisses + m2.ASTTypeCacheMisses,
		ASTTotalLoadTime:       m1.ASTTotalLoadTime + m2.ASTTotalLoadTime,
		ASTTotalDecorationTime: m1.ASTTotalDecorationTime + m2.ASTTotalDecorationTime,
		TotalProcessingTime:    m1.TotalProcessingTime + m2.TotalProcessingTime,
	}
}

// Processor maintains the state of the metrics system
type Processor struct {
	isLoggingFn IsLoggingFn
	loggingfFn  LoggingfFn

	state metricsState
}

//go:generate moqueries IsLoggingFn

// IsLoggingFn is the function that determines if logging is on
type IsLoggingFn func() bool

//go:generate moqueries LoggingfFn

// LoggingfFn is the logging function to output finalized metrics
type LoggingfFn func(format string, args ...interface{})

// NewMetrics returns a new system for gathering metrics
func NewMetrics(isLoggingFn IsLoggingFn, loggingfFn LoggingfFn) *Processor {
	return &Processor{
		isLoggingFn: isLoggingFn,
		loggingfFn:  loggingfFn,
	}
}

// Finalize is called after generating mocks to log metrics
func (m *Processor) Finalize() {
	if m.isLoggingFn() {
		m.loggingfFn(metricsLogKey+" %s", serializeState(m.state))
	}
}

// ASTPkgCacheHitsInc increments the ASTPkgCacheHits metric
func (m *Processor) ASTPkgCacheHitsInc() {
	m.state.ASTPkgCacheHits++
}

// ASTPkgCacheMissesInc increments the ASTPkgCacheMisses metric
func (m *Processor) ASTPkgCacheMissesInc() {
	m.state.ASTPkgCacheMisses++
}

// ASTTotalLoadTimeInc increments the ASTTotalLoadTime duration metric by the
// d duration specified
func (m *Processor) ASTTotalLoadTimeInc(d time.Duration) {
	m.state.ASTTotalLoadTime += d
}

// ASTTotalDecorationTimeInc increments the ASTTotalDecorationTime duration
// metric by the d duration specified
func (m *Processor) ASTTotalDecorationTimeInc(d time.Duration) {
	m.state.ASTTotalDecorationTime += d
}

// TotalProcessingTimeInc increments the TotalProcessingTime duration metric
// by the d duration specified
func (m *Processor) TotalProcessingTimeInc(d time.Duration) {
	m.state.TotalProcessingTime += d
}

func serializeState(state metricsState) []byte {
	state.ASTTotalLoadTimeStr = state.ASTTotalLoadTime.String()
	state.ASTTotalDecorationTimeStr = state.ASTTotalDecorationTime.String()
	state.TotalProcessingTimeStr = state.TotalProcessingTime.String()

	b, err := json.Marshal(state)
	if err != nil {
		logs.Panic("Could not marshal metrics to JSON", err)
	}

	buf := &bytes.Buffer{}
	err = json.Compact(buf, b)
	if err != nil {
		logs.Panic("Could not compact metrics", err)
	}

	return buf.Bytes()
}
