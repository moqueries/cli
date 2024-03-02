// Code generated by Moqueries - https://moqueries.org - DO NOT EDIT.

package metrics_test

import (
	"fmt"
	"math/bits"
	"sync/atomic"

	"moqueries.org/cli/metrics"
	"moqueries.org/runtime/hash"
	"moqueries.org/runtime/moq"
)

// moqLoggingfFn holds the state of a moq of the LoggingfFn type
type moqLoggingfFn struct {
	scene  *moq.Scene
	config moq.Config
	moq    *moqLoggingfFn_mock

	resultsByParams []moqLoggingfFn_resultsByParams

	runtime struct {
		parameterIndexing struct {
			format moq.ParamIndexing
			args   moq.ParamIndexing
		}
	}
}

// moqLoggingfFn_mock isolates the mock interface of the LoggingfFn type
type moqLoggingfFn_mock struct {
	moq *moqLoggingfFn
}

// moqLoggingfFn_params holds the params of the LoggingfFn type
type moqLoggingfFn_params struct {
	format string
	args   []interface{}
}

// moqLoggingfFn_paramsKey holds the map key params of the LoggingfFn type
type moqLoggingfFn_paramsKey struct {
	params struct{ format string }
	hashes struct {
		format hash.Hash
		args   hash.Hash
	}
}

// moqLoggingfFn_resultsByParams contains the results for a given set of
// parameters for the LoggingfFn type
type moqLoggingfFn_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqLoggingfFn_paramsKey]*moqLoggingfFn_results
}

// moqLoggingfFn_doFn defines the type of function needed when calling andDo
// for the LoggingfFn type
type moqLoggingfFn_doFn func(format string, args ...interface{})

// moqLoggingfFn_doReturnFn defines the type of function needed when calling
// doReturnResults for the LoggingfFn type
type moqLoggingfFn_doReturnFn func(format string, args ...interface{})

// moqLoggingfFn_results holds the results of the LoggingfFn type
type moqLoggingfFn_results struct {
	params  moqLoggingfFn_params
	results []struct {
		values     *struct{}
		sequence   uint32
		doFn       moqLoggingfFn_doFn
		doReturnFn moqLoggingfFn_doReturnFn
	}
	index  uint32
	repeat *moq.RepeatVal
}

// moqLoggingfFn_fnRecorder routes recorded function calls to the moqLoggingfFn
// moq
type moqLoggingfFn_fnRecorder struct {
	params    moqLoggingfFn_params
	anyParams uint64
	sequence  bool
	results   *moqLoggingfFn_results
	moq       *moqLoggingfFn
}

// moqLoggingfFn_anyParams isolates the any params functions of the LoggingfFn
// type
type moqLoggingfFn_anyParams struct {
	recorder *moqLoggingfFn_fnRecorder
}

// newMoqLoggingfFn creates a new moq of the LoggingfFn type
func newMoqLoggingfFn(scene *moq.Scene, config *moq.Config) *moqLoggingfFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &moqLoggingfFn{
		scene:  scene,
		config: *config,
		moq:    &moqLoggingfFn_mock{},

		runtime: struct {
			parameterIndexing struct {
				format moq.ParamIndexing
				args   moq.ParamIndexing
			}
		}{parameterIndexing: struct {
			format moq.ParamIndexing
			args   moq.ParamIndexing
		}{
			format: moq.ParamIndexByValue,
			args:   moq.ParamIndexByHash,
		}},
	}
	m.moq.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the moq implementation of the LoggingfFn type
func (m *moqLoggingfFn) mock() metrics.LoggingfFn {
	return func(format string, args ...interface{}) {
		m.scene.T.Helper()
		moq := &moqLoggingfFn_mock{moq: m}
		moq.fn(format, args...)
	}
}

func (m *moqLoggingfFn_mock) fn(format string, args ...interface{}) {
	m.moq.scene.T.Helper()
	params := moqLoggingfFn_params{
		format: format,
		args:   args,
	}
	var results *moqLoggingfFn_results
	for _, resultsByParams := range m.moq.resultsByParams {
		paramsKey := m.moq.paramsKey(params, resultsByParams.anyParams)
		var ok bool
		results, ok = resultsByParams.results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.moq.config.Expectation == moq.Strict {
			m.moq.scene.T.Fatalf("Unexpected call to %s", m.moq.prettyParams(params))
		}
		return
	}

	i := int(atomic.AddUint32(&results.index, 1)) - 1
	if i >= results.repeat.ResultCount {
		if !results.repeat.AnyTimes {
			if m.moq.config.Expectation == moq.Strict {
				m.moq.scene.T.Fatalf("Too many calls to %s", m.moq.prettyParams(params))
			}
			return
		}
		i = results.repeat.ResultCount - 1
	}

	result := results.results[i]
	if result.sequence != 0 {
		sequence := m.moq.scene.NextMockSequence()
		if (!results.repeat.AnyTimes && result.sequence != sequence) || result.sequence > sequence {
			m.moq.scene.T.Fatalf("Call sequence does not match call to %s", m.moq.prettyParams(params))
		}
	}

	if result.doFn != nil {
		result.doFn(format, args...)
	}

	if result.doReturnFn != nil {
		result.doReturnFn(format, args...)
	}
	return
}

func (m *moqLoggingfFn) onCall(format string, args ...interface{}) *moqLoggingfFn_fnRecorder {
	return &moqLoggingfFn_fnRecorder{
		params: moqLoggingfFn_params{
			format: format,
			args:   args,
		},
		sequence: m.config.Sequence == moq.SeqDefaultOn,
		moq:      m,
	}
}

func (r *moqLoggingfFn_fnRecorder) any() *moqLoggingfFn_anyParams {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	return &moqLoggingfFn_anyParams{recorder: r}
}

func (a *moqLoggingfFn_anyParams) format() *moqLoggingfFn_fnRecorder {
	a.recorder.anyParams |= 1 << 0
	return a.recorder
}

func (a *moqLoggingfFn_anyParams) args() *moqLoggingfFn_fnRecorder {
	a.recorder.anyParams |= 1 << 1
	return a.recorder
}

func (r *moqLoggingfFn_fnRecorder) seq() *moqLoggingfFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqLoggingfFn_fnRecorder) noSeq() *moqLoggingfFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqLoggingfFn_fnRecorder) returnResults() *moqLoggingfFn_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values     *struct{}
		sequence   uint32
		doFn       moqLoggingfFn_doFn
		doReturnFn moqLoggingfFn_doReturnFn
	}{
		values:   &struct{}{},
		sequence: sequence,
	})
	return r
}

func (r *moqLoggingfFn_fnRecorder) andDo(fn moqLoggingfFn_doFn) *moqLoggingfFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqLoggingfFn_fnRecorder) doReturnResults(fn moqLoggingfFn_doReturnFn) *moqLoggingfFn_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values     *struct{}
		sequence   uint32
		doFn       moqLoggingfFn_doFn
		doReturnFn moqLoggingfFn_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqLoggingfFn_fnRecorder) findResults() {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.results.repeat.Increment(r.moq.scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.anyParams)
	insertAt := -1
	var results *moqLoggingfFn_resultsByParams
	for n, res := range r.moq.resultsByParams {
		if res.anyParams == r.anyParams {
			results = &res
			break
		}
		if res.anyCount > anyCount {
			insertAt = n
		}
	}
	if results == nil {
		results = &moqLoggingfFn_resultsByParams{
			anyCount:  anyCount,
			anyParams: r.anyParams,
			results:   map[moqLoggingfFn_paramsKey]*moqLoggingfFn_results{},
		}
		r.moq.resultsByParams = append(r.moq.resultsByParams, *results)
		if insertAt != -1 && insertAt+1 < len(r.moq.resultsByParams) {
			copy(r.moq.resultsByParams[insertAt+1:], r.moq.resultsByParams[insertAt:0])
			r.moq.resultsByParams[insertAt] = *results
		}
	}

	paramsKey := r.moq.paramsKey(r.params, r.anyParams)

	var ok bool
	r.results, ok = results.results[paramsKey]
	if !ok {
		r.results = &moqLoggingfFn_results{
			params:  r.params,
			results: nil,
			index:   0,
			repeat:  &moq.RepeatVal{},
		}
		results.results[paramsKey] = r.results
	}

	r.results.repeat.Increment(r.moq.scene.T)
}

func (r *moqLoggingfFn_fnRecorder) repeat(repeaters ...moq.Repeater) *moqLoggingfFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults or doReturnResults must be called before calling repeat")
		return nil
	}
	r.results.repeat.Repeat(r.moq.scene.T, repeaters)
	last := r.results.results[len(r.results.results)-1]
	for n := 0; n < r.results.repeat.ResultCount-1; n++ {
		if r.sequence {
			last = struct {
				values     *struct{}
				sequence   uint32
				doFn       moqLoggingfFn_doFn
				doReturnFn moqLoggingfFn_doReturnFn
			}{
				values:   last.values,
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (m *moqLoggingfFn) prettyParams(params moqLoggingfFn_params) string {
	return fmt.Sprintf("LoggingfFn(%#v, %#v)", params.format, params.args)
}

func (m *moqLoggingfFn) paramsKey(params moqLoggingfFn_params, anyParams uint64) moqLoggingfFn_paramsKey {
	m.scene.T.Helper()
	var formatUsed string
	var formatUsedHash hash.Hash
	if anyParams&(1<<0) == 0 {
		if m.runtime.parameterIndexing.format == moq.ParamIndexByValue {
			formatUsed = params.format
		} else {
			formatUsedHash = hash.DeepHash(params.format)
		}
	}
	var argsUsedHash hash.Hash
	if anyParams&(1<<1) == 0 {
		if m.runtime.parameterIndexing.args == moq.ParamIndexByValue {
			m.scene.T.Fatalf("The args parameter can't be indexed by value")
		}
		argsUsedHash = hash.DeepHash(params.args)
	}
	return moqLoggingfFn_paramsKey{
		params: struct{ format string }{
			format: formatUsed,
		},
		hashes: struct {
			format hash.Hash
			args   hash.Hash
		}{
			format: formatUsedHash,
			args:   argsUsedHash,
		},
	}
}

// Reset resets the state of the moq
func (m *moqLoggingfFn) Reset() { m.resultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqLoggingfFn) AssertExpectationsMet() {
	m.scene.T.Helper()
	for _, res := range m.resultsByParams {
		for _, results := range res.results {
			missing := results.repeat.MinTimes - int(atomic.LoadUint32(&results.index))
			if missing > 0 {
				m.scene.T.Errorf("Expected %d additional call(s) to %s", missing, m.prettyParams(results.params))
			}
		}
	}
}
