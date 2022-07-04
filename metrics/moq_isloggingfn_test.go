// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package metrics_test

import (
	"fmt"
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/metrics"
	"github.com/myshkin5/moqueries/moq"
)

// moqIsLoggingFn holds the state of a moq of the IsLoggingFn type
type moqIsLoggingFn struct {
	scene  *moq.Scene
	config moq.Config
	moq    *moqIsLoggingFn_mock

	resultsByParams []moqIsLoggingFn_resultsByParams

	runtime struct {
		parameterIndexing struct{}
	}
}

// moqIsLoggingFn_mock isolates the mock interface of the IsLoggingFn type
type moqIsLoggingFn_mock struct {
	moq *moqIsLoggingFn
}

// moqIsLoggingFn_params holds the params of the IsLoggingFn type
type moqIsLoggingFn_params struct{}

// moqIsLoggingFn_paramsKey holds the map key params of the IsLoggingFn type
type moqIsLoggingFn_paramsKey struct {
	params struct{}
	hashes struct{}
}

// moqIsLoggingFn_resultsByParams contains the results for a given set of parameters for the IsLoggingFn type
type moqIsLoggingFn_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqIsLoggingFn_paramsKey]*moqIsLoggingFn_results
}

// moqIsLoggingFn_doFn defines the type of function needed when calling andDo for the IsLoggingFn type
type moqIsLoggingFn_doFn func()

// moqIsLoggingFn_doReturnFn defines the type of function needed when calling doReturnResults for the IsLoggingFn type
type moqIsLoggingFn_doReturnFn func() bool

// moqIsLoggingFn_results holds the results of the IsLoggingFn type
type moqIsLoggingFn_results struct {
	params  moqIsLoggingFn_params
	results []struct {
		values *struct {
			result1 bool
		}
		sequence   uint32
		doFn       moqIsLoggingFn_doFn
		doReturnFn moqIsLoggingFn_doReturnFn
	}
	index  uint32
	repeat *moq.RepeatVal
}

// moqIsLoggingFn_fnRecorder routes recorded function calls to the moqIsLoggingFn moq
type moqIsLoggingFn_fnRecorder struct {
	params    moqIsLoggingFn_params
	anyParams uint64
	sequence  bool
	results   *moqIsLoggingFn_results
	moq       *moqIsLoggingFn
}

// moqIsLoggingFn_anyParams isolates the any params functions of the IsLoggingFn type
type moqIsLoggingFn_anyParams struct {
	recorder *moqIsLoggingFn_fnRecorder
}

// newMoqIsLoggingFn creates a new moq of the IsLoggingFn type
func newMoqIsLoggingFn(scene *moq.Scene, config *moq.Config) *moqIsLoggingFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &moqIsLoggingFn{
		scene:  scene,
		config: *config,
		moq:    &moqIsLoggingFn_mock{},

		runtime: struct {
			parameterIndexing struct{}
		}{parameterIndexing: struct{}{}},
	}
	m.moq.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the moq implementation of the IsLoggingFn type
func (m *moqIsLoggingFn) mock() metrics.IsLoggingFn {
	return func() bool { moq := &moqIsLoggingFn_mock{moq: m}; return moq.fn() }
}

func (m *moqIsLoggingFn_mock) fn() (result1 bool) {
	m.moq.scene.T.Helper()
	params := moqIsLoggingFn_params{}
	var results *moqIsLoggingFn_results
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
		result.doFn()
	}

	if result.values != nil {
		result1 = result.values.result1
	}
	if result.doReturnFn != nil {
		result1 = result.doReturnFn()
	}
	return
}

func (m *moqIsLoggingFn) onCall() *moqIsLoggingFn_fnRecorder {
	return &moqIsLoggingFn_fnRecorder{
		params:   moqIsLoggingFn_params{},
		sequence: m.config.Sequence == moq.SeqDefaultOn,
		moq:      m,
	}
}

func (r *moqIsLoggingFn_fnRecorder) any() *moqIsLoggingFn_anyParams {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	return &moqIsLoggingFn_anyParams{recorder: r}
}

func (r *moqIsLoggingFn_fnRecorder) seq() *moqIsLoggingFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqIsLoggingFn_fnRecorder) noSeq() *moqIsLoggingFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqIsLoggingFn_fnRecorder) returnResults(result1 bool) *moqIsLoggingFn_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			result1 bool
		}
		sequence   uint32
		doFn       moqIsLoggingFn_doFn
		doReturnFn moqIsLoggingFn_doReturnFn
	}{
		values: &struct {
			result1 bool
		}{
			result1: result1,
		},
		sequence: sequence,
	})
	return r
}

func (r *moqIsLoggingFn_fnRecorder) andDo(fn moqIsLoggingFn_doFn) *moqIsLoggingFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqIsLoggingFn_fnRecorder) doReturnResults(fn moqIsLoggingFn_doReturnFn) *moqIsLoggingFn_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			result1 bool
		}
		sequence   uint32
		doFn       moqIsLoggingFn_doFn
		doReturnFn moqIsLoggingFn_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqIsLoggingFn_fnRecorder) findResults() {
	if r.results != nil {
		r.results.repeat.Increment(r.moq.scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.anyParams)
	insertAt := -1
	var results *moqIsLoggingFn_resultsByParams
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
		results = &moqIsLoggingFn_resultsByParams{
			anyCount:  anyCount,
			anyParams: r.anyParams,
			results:   map[moqIsLoggingFn_paramsKey]*moqIsLoggingFn_results{},
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
		r.results = &moqIsLoggingFn_results{
			params:  r.params,
			results: nil,
			index:   0,
			repeat:  &moq.RepeatVal{},
		}
		results.results[paramsKey] = r.results
	}

	r.results.repeat.Increment(r.moq.scene.T)
}

func (r *moqIsLoggingFn_fnRecorder) repeat(repeaters ...moq.Repeater) *moqIsLoggingFn_fnRecorder {
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
				values *struct {
					result1 bool
				}
				sequence   uint32
				doFn       moqIsLoggingFn_doFn
				doReturnFn moqIsLoggingFn_doReturnFn
			}{
				values:   last.values,
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (m *moqIsLoggingFn) prettyParams(params moqIsLoggingFn_params) string {
	return fmt.Sprintf("IsLoggingFn()")
}

func (m *moqIsLoggingFn) paramsKey(params moqIsLoggingFn_params, anyParams uint64) moqIsLoggingFn_paramsKey {
	return moqIsLoggingFn_paramsKey{
		params: struct{}{},
		hashes: struct{}{},
	}
}

// Reset resets the state of the moq
func (m *moqIsLoggingFn) Reset() { m.resultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqIsLoggingFn) AssertExpectationsMet() {
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
