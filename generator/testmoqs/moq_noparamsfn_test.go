// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package testmoqs_test

import (
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/generator/testmoqs"
	"github.com/myshkin5/moqueries/moq"
)

// moqNoParamsFn holds the state of a moq of the NoParamsFn type
type moqNoParamsFn struct {
	scene  *moq.Scene
	config moq.Config
	moq    *moqNoParamsFn_mock

	resultsByParams []moqNoParamsFn_resultsByParams

	runtime struct {
		parameterIndexing struct{}
	}
	// moqNoParamsFn_mock isolates the mock interface of the NoParamsFn type
}

type moqNoParamsFn_mock struct {
	moq *moqNoParamsFn
}

// moqNoParamsFn_params holds the params of the NoParamsFn type
type moqNoParamsFn_params struct{}

// moqNoParamsFn_paramsKey holds the map key params of the NoParamsFn type
type moqNoParamsFn_paramsKey struct {
	params struct{}
	hashes struct{}
}

// moqNoParamsFn_resultsByParams contains the results for a given set of parameters for the NoParamsFn type
type moqNoParamsFn_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqNoParamsFn_paramsKey]*moqNoParamsFn_results
}

// moqNoParamsFn_doFn defines the type of function needed when calling andDo for the NoParamsFn type
type moqNoParamsFn_doFn func()

// moqNoParamsFn_doReturnFn defines the type of function needed when calling doReturnResults for the NoParamsFn type
type moqNoParamsFn_doReturnFn func() (sResult string, err error)

// moqNoParamsFn_results holds the results of the NoParamsFn type
type moqNoParamsFn_results struct {
	params  moqNoParamsFn_params
	results []struct {
		values *struct {
			sResult string
			err     error
		}
		sequence   uint32
		doFn       moqNoParamsFn_doFn
		doReturnFn moqNoParamsFn_doReturnFn
	}
	index  uint32
	repeat *moq.RepeatVal
}

// moqNoParamsFn_fnRecorder routes recorded function calls to the moqNoParamsFn moq
type moqNoParamsFn_fnRecorder struct {
	params    moqNoParamsFn_params
	anyParams uint64
	sequence  bool
	results   *moqNoParamsFn_results
	moq       *moqNoParamsFn
}

// moqNoParamsFn_anyParams isolates the any params functions of the NoParamsFn type
type moqNoParamsFn_anyParams struct {
	recorder *moqNoParamsFn_fnRecorder
}

// newMoqNoParamsFn creates a new moq of the NoParamsFn type
func newMoqNoParamsFn(scene *moq.Scene, config *moq.Config) *moqNoParamsFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &moqNoParamsFn{
		scene:  scene,
		config: *config,
		moq:    &moqNoParamsFn_mock{},

		runtime: struct {
			parameterIndexing struct{}
		}{parameterIndexing: struct{}{}},
	}
	m.moq.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the moq implementation of the NoParamsFn type
func (m *moqNoParamsFn) mock() testmoqs.NoParamsFn {
	return func() (_ string, _ error) { moq := &moqNoParamsFn_mock{moq: m}; return moq.fn() }
}

func (m *moqNoParamsFn_mock) fn() (sResult string, err error) {
	params := moqNoParamsFn_params{}
	var results *moqNoParamsFn_results
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

	if result.values != nil {
		sResult = result.values.sResult
		err = result.values.err
	}
	if result.doReturnFn != nil {
		sResult, err = result.doReturnFn()
	}
	return
}

func (m *moqNoParamsFn) onCall() *moqNoParamsFn_fnRecorder {
	return &moqNoParamsFn_fnRecorder{
		params:   moqNoParamsFn_params{},
		sequence: m.config.Sequence == moq.SeqDefaultOn,
		moq:      m,
	}
}

func (r *moqNoParamsFn_fnRecorder) any() *moqNoParamsFn_anyParams {
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	return &moqNoParamsFn_anyParams{recorder: r}
}

func (r *moqNoParamsFn_fnRecorder) seq() *moqNoParamsFn_fnRecorder {
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqNoParamsFn_fnRecorder) noSeq() *moqNoParamsFn_fnRecorder {
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqNoParamsFn_fnRecorder) returnResults(sResult string, err error) *moqNoParamsFn_fnRecorder {
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			sResult string
			err     error
		}
		sequence   uint32
		doFn       moqNoParamsFn_doFn
		doReturnFn moqNoParamsFn_doReturnFn
	}{
		values: &struct {
			sResult string
			err     error
		}{
			sResult: sResult,
			err:     err,
		},
		sequence: sequence,
	})
	return r
}

func (r *moqNoParamsFn_fnRecorder) andDo(fn moqNoParamsFn_doFn) *moqNoParamsFn_fnRecorder {
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqNoParamsFn_fnRecorder) doReturnResults(fn moqNoParamsFn_doReturnFn) *moqNoParamsFn_fnRecorder {
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			sResult string
			err     error
		}
		sequence   uint32
		doFn       moqNoParamsFn_doFn
		doReturnFn moqNoParamsFn_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqNoParamsFn_fnRecorder) findResults() {
	if r.results != nil {
		r.results.repeat.Increment(r.moq.scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.anyParams)
	insertAt := -1
	var results *moqNoParamsFn_resultsByParams
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
		results = &moqNoParamsFn_resultsByParams{
			anyCount:  anyCount,
			anyParams: r.anyParams,
			results:   map[moqNoParamsFn_paramsKey]*moqNoParamsFn_results{},
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
		r.results = &moqNoParamsFn_results{
			params:  r.params,
			results: nil,
			index:   0,
			repeat:  &moq.RepeatVal{},
		}
		results.results[paramsKey] = r.results
	}

	r.results.repeat.Increment(r.moq.scene.T)
}

func (r *moqNoParamsFn_fnRecorder) repeat(repeaters ...moq.Repeater) *moqNoParamsFn_fnRecorder {
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
					sResult string
					err     error
				}
				sequence   uint32
				doFn       moqNoParamsFn_doFn
				doReturnFn moqNoParamsFn_doReturnFn
			}{
				values: &struct {
					sResult string
					err     error
				}{
					sResult: last.values.sResult,
					err:     last.values.err,
				},
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (m *moqNoParamsFn) paramsKey(params moqNoParamsFn_params, anyParams uint64) moqNoParamsFn_paramsKey {
	return moqNoParamsFn_paramsKey{
		params: struct{}{},
		hashes: struct{}{},
	}
}

// Reset resets the state of the moq
func (m *moqNoParamsFn) Reset() { m.resultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqNoParamsFn) AssertExpectationsMet() {
	for _, res := range m.resultsByParams {
		for _, results := range res.results {
			missing := results.repeat.MinTimes - int(atomic.LoadUint32(&results.index))
			if missing > 0 {
				m.scene.T.Errorf("Expected %d additional call(s) with parameters %#v", missing, results.params)
			}
		}
	}
}
