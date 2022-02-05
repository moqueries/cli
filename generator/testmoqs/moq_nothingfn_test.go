// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package testmoqs_test

import (
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/generator/testmoqs"
	"github.com/myshkin5/moqueries/moq"
)

// moqNothingFn holds the state of a moq of the NothingFn type
type moqNothingFn struct {
	scene  *moq.Scene
	config moq.Config
	moq    *moqNothingFn_mock

	resultsByParams []moqNothingFn_resultsByParams

	runtime struct {
		parameterIndexing struct{}
	}
	// moqNothingFn_mock isolates the mock interface of the NothingFn type
}

type moqNothingFn_mock struct {
	moq *moqNothingFn
}

// moqNothingFn_params holds the params of the NothingFn type
type moqNothingFn_params struct{}

// moqNothingFn_paramsKey holds the map key params of the NothingFn type
type moqNothingFn_paramsKey struct {
	params struct{}
	hashes struct{}
}

// moqNothingFn_resultsByParams contains the results for a given set of parameters for the NothingFn type
type moqNothingFn_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqNothingFn_paramsKey]*moqNothingFn_results
}

// moqNothingFn_doFn defines the type of function needed when calling andDo for the NothingFn type
type moqNothingFn_doFn func()

// moqNothingFn_doReturnFn defines the type of function needed when calling doReturnResults for the NothingFn type
type moqNothingFn_doReturnFn func()

// moqNothingFn_results holds the results of the NothingFn type
type moqNothingFn_results struct {
	params  moqNothingFn_params
	results []struct {
		values     *struct{}
		sequence   uint32
		doFn       moqNothingFn_doFn
		doReturnFn moqNothingFn_doReturnFn
	}
	index  uint32
	repeat *moq.RepeatVal
}

// moqNothingFn_fnRecorder routes recorded function calls to the moqNothingFn moq
type moqNothingFn_fnRecorder struct {
	params    moqNothingFn_params
	anyParams uint64
	sequence  bool
	results   *moqNothingFn_results
	moq       *moqNothingFn
}

// moqNothingFn_anyParams isolates the any params functions of the NothingFn type
type moqNothingFn_anyParams struct {
	recorder *moqNothingFn_fnRecorder
}

// newMoqNothingFn creates a new moq of the NothingFn type
func newMoqNothingFn(scene *moq.Scene, config *moq.Config) *moqNothingFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &moqNothingFn{
		scene:  scene,
		config: *config,
		moq:    &moqNothingFn_mock{},

		runtime: struct {
			parameterIndexing struct{}
		}{parameterIndexing: struct{}{}},
	}
	m.moq.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the moq implementation of the NothingFn type
func (m *moqNothingFn) mock() testmoqs.NothingFn {
	return func() { moq := &moqNothingFn_mock{moq: m}; moq.fn() }
}

func (m *moqNothingFn_mock) fn() {
	params := moqNothingFn_params{}
	var results *moqNothingFn_results
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
			m.moq.scene.T.Fatalf("Unexpected call with parameters %#v", params)
		}
		return
	}

	i := int(atomic.AddUint32(&results.index, 1)) - 1
	if i >= results.repeat.ResultCount {
		if !results.repeat.AnyTimes {
			if m.moq.config.Expectation == moq.Strict {
				m.moq.scene.T.Fatalf("Too many calls to mock with parameters %#v", params)
			}
			return
		}
		i = results.repeat.ResultCount - 1
	}

	result := results.results[i]
	if result.sequence != 0 {
		sequence := m.moq.scene.NextMockSequence()
		if (!results.repeat.AnyTimes && result.sequence != sequence) || result.sequence > sequence {
			m.moq.scene.T.Fatalf("Call sequence does not match %#v", params)
		}
	}

	if result.doFn != nil {
		result.doFn()
	}

	if result.doReturnFn != nil {
		result.doReturnFn()
	}
	return
}

func (m *moqNothingFn) onCall() *moqNothingFn_fnRecorder {
	return &moqNothingFn_fnRecorder{
		params:   moqNothingFn_params{},
		sequence: m.config.Sequence == moq.SeqDefaultOn,
		moq:      m,
	}
}

func (r *moqNothingFn_fnRecorder) any() *moqNothingFn_anyParams {
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	return &moqNothingFn_anyParams{recorder: r}
}

func (r *moqNothingFn_fnRecorder) seq() *moqNothingFn_fnRecorder {
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqNothingFn_fnRecorder) noSeq() *moqNothingFn_fnRecorder {
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqNothingFn_fnRecorder) returnResults() *moqNothingFn_fnRecorder {
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values     *struct{}
		sequence   uint32
		doFn       moqNothingFn_doFn
		doReturnFn moqNothingFn_doReturnFn
	}{
		values:   &struct{}{},
		sequence: sequence,
	})
	return r
}

func (r *moqNothingFn_fnRecorder) andDo(fn moqNothingFn_doFn) *moqNothingFn_fnRecorder {
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqNothingFn_fnRecorder) doReturnResults(fn moqNothingFn_doReturnFn) *moqNothingFn_fnRecorder {
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values     *struct{}
		sequence   uint32
		doFn       moqNothingFn_doFn
		doReturnFn moqNothingFn_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqNothingFn_fnRecorder) findResults() {
	if r.results != nil {
		r.results.repeat.Increment(r.moq.scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.anyParams)
	insertAt := -1
	var results *moqNothingFn_resultsByParams
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
		results = &moqNothingFn_resultsByParams{
			anyCount:  anyCount,
			anyParams: r.anyParams,
			results:   map[moqNothingFn_paramsKey]*moqNothingFn_results{},
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
		r.results = &moqNothingFn_results{
			params:  r.params,
			results: nil,
			index:   0,
			repeat:  &moq.RepeatVal{},
		}
		results.results[paramsKey] = r.results
	}

	r.results.repeat.Increment(r.moq.scene.T)
}

func (r *moqNothingFn_fnRecorder) repeat(repeaters ...moq.Repeater) *moqNothingFn_fnRecorder {
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
				doFn       moqNothingFn_doFn
				doReturnFn moqNothingFn_doReturnFn
			}{
				values:   &struct{}{},
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (m *moqNothingFn) paramsKey(params moqNothingFn_params, anyParams uint64) moqNothingFn_paramsKey {
	return moqNothingFn_paramsKey{
		params: struct{}{},
		hashes: struct{}{},
	}
}

// Reset resets the state of the moq
func (m *moqNothingFn) Reset() { m.resultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqNothingFn) AssertExpectationsMet() {
	for _, res := range m.resultsByParams {
		for _, results := range res.results {
			missing := results.repeat.MinTimes - int(atomic.LoadUint32(&results.index))
			if missing > 0 {
				m.scene.T.Errorf("Expected %d additional call(s) with parameters %#v", missing, results.params)
			}
		}
	}
}
