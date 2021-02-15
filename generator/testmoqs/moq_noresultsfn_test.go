// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package testmoqs_test

import (
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/generator/testmoqs"
	"github.com/myshkin5/moqueries/moq"
)

// moqNoResultsFn holds the state of a moq of the NoResultsFn type
type moqNoResultsFn struct {
	scene           *moq.Scene
	config          moq.Config
	resultsByParams []moqNoResultsFn_resultsByParams
}

// moqNoResultsFn_mock isolates the mock interface of the NoResultsFn type
type moqNoResultsFn_mock struct {
	moq *moqNoResultsFn
}

// moqNoResultsFn_params holds the params of the NoResultsFn type
type moqNoResultsFn_params struct {
	sParam string
	bParam bool
}

// moqNoResultsFn_paramsKey holds the map key params of the NoResultsFn type
type moqNoResultsFn_paramsKey struct {
	sParam string
	bParam bool
}

// moqNoResultsFn_resultsByParams contains the results for a given set of parameters for the NoResultsFn type
type moqNoResultsFn_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqNoResultsFn_paramsKey]*moqNoResultsFn_results
}

// moqNoResultsFn_doFn defines the type of function needed when calling andDo for the NoResultsFn type
type moqNoResultsFn_doFn func(sParam string, bParam bool)

// moqNoResultsFn_doReturnFn defines the type of function needed when calling doReturnResults for the NoResultsFn type
type moqNoResultsFn_doReturnFn func(sParam string, bParam bool)

// moqNoResultsFn_results holds the results of the NoResultsFn type
type moqNoResultsFn_results struct {
	params  moqNoResultsFn_params
	results []struct {
		values *struct {
		}
		sequence   uint32
		doFn       moqNoResultsFn_doFn
		doReturnFn moqNoResultsFn_doReturnFn
	}
	index    uint32
	anyTimes bool
}

// moqNoResultsFn_fnRecorder routes recorded function calls to the moqNoResultsFn moq
type moqNoResultsFn_fnRecorder struct {
	params    moqNoResultsFn_params
	paramsKey moqNoResultsFn_paramsKey
	anyParams uint64
	sequence  bool
	results   *moqNoResultsFn_results
	moq       *moqNoResultsFn
}

// moqNoResultsFn_anyParams isolates the any params functions of the NoResultsFn type
type moqNoResultsFn_anyParams struct {
	recorder *moqNoResultsFn_fnRecorder
}

// newMoqNoResultsFn creates a new moq of the NoResultsFn type
func newMoqNoResultsFn(scene *moq.Scene, config *moq.Config) *moqNoResultsFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &moqNoResultsFn{
		scene:  scene,
		config: *config,
	}
	scene.AddMoq(m)
	return m
}

// mock returns the moq implementation of the NoResultsFn type
func (m *moqNoResultsFn) mock() testmoqs.NoResultsFn {
	return func(sParam string, bParam bool) { moq := &moqNoResultsFn_mock{moq: m}; moq.fn(sParam, bParam) }
}

func (m *moqNoResultsFn_mock) fn(sParam string, bParam bool) {
	params := moqNoResultsFn_params{
		sParam: sParam,
		bParam: bParam,
	}
	var results *moqNoResultsFn_results
	for _, resultsByParams := range m.moq.resultsByParams {
		var sParamUsed string
		if resultsByParams.anyParams&(1<<0) == 0 {
			sParamUsed = sParam
		}
		var bParamUsed bool
		if resultsByParams.anyParams&(1<<1) == 0 {
			bParamUsed = bParam
		}
		paramsKey := moqNoResultsFn_paramsKey{
			sParam: sParamUsed,
			bParam: bParamUsed,
		}
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
	if i >= len(results.results) {
		if !results.anyTimes {
			if m.moq.config.Expectation == moq.Strict {
				m.moq.scene.T.Fatalf("Too many calls to mock with parameters %#v", params)
			}
			return
		}
		i = len(results.results) - 1
	}

	result := results.results[i]
	if result.sequence != 0 {
		sequence := m.moq.scene.NextMockSequence()
		if (!results.anyTimes && result.sequence != sequence) || result.sequence > sequence {
			m.moq.scene.T.Fatalf("Call sequence does not match %#v", params)
		}
	}

	if result.doFn != nil {
		result.doFn(sParam, bParam)
	}

	if result.doReturnFn != nil {
		result.doReturnFn(sParam, bParam)
	}
	return
}

func (m *moqNoResultsFn) onCall(sParam string, bParam bool) *moqNoResultsFn_fnRecorder {
	return &moqNoResultsFn_fnRecorder{
		params: moqNoResultsFn_params{
			sParam: sParam,
			bParam: bParam,
		},
		paramsKey: moqNoResultsFn_paramsKey{
			sParam: sParam,
			bParam: bParam,
		},
		sequence: m.config.Sequence == moq.SeqDefaultOn,
		moq:      m,
	}
}

func (r *moqNoResultsFn_fnRecorder) any() *moqNoResultsFn_anyParams {
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	return &moqNoResultsFn_anyParams{recorder: r}
}

func (a *moqNoResultsFn_anyParams) sParam() *moqNoResultsFn_fnRecorder {
	a.recorder.anyParams |= 1 << 0
	return a.recorder
}

func (a *moqNoResultsFn_anyParams) bParam() *moqNoResultsFn_fnRecorder {
	a.recorder.anyParams |= 1 << 1
	return a.recorder
}

func (r *moqNoResultsFn_fnRecorder) seq() *moqNoResultsFn_fnRecorder {
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqNoResultsFn_fnRecorder) noSeq() *moqNoResultsFn_fnRecorder {
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqNoResultsFn_fnRecorder) returnResults() *moqNoResultsFn_fnRecorder {
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
		}
		sequence   uint32
		doFn       moqNoResultsFn_doFn
		doReturnFn moqNoResultsFn_doReturnFn
	}{
		values: &struct {
		}{},
		sequence: sequence,
	})
	return r
}

func (r *moqNoResultsFn_fnRecorder) andDo(fn moqNoResultsFn_doFn) *moqNoResultsFn_fnRecorder {
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqNoResultsFn_fnRecorder) doReturnResults(fn moqNoResultsFn_doReturnFn) *moqNoResultsFn_fnRecorder {
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
		}
		sequence   uint32
		doFn       moqNoResultsFn_doFn
		doReturnFn moqNoResultsFn_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqNoResultsFn_fnRecorder) findResults() {
	if r.results == nil {
		anyCount := bits.OnesCount64(r.anyParams)
		insertAt := -1
		var results *moqNoResultsFn_resultsByParams
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
			results = &moqNoResultsFn_resultsByParams{
				anyCount:  anyCount,
				anyParams: r.anyParams,
				results:   map[moqNoResultsFn_paramsKey]*moqNoResultsFn_results{},
			}
			r.moq.resultsByParams = append(r.moq.resultsByParams, *results)
			if insertAt != -1 && insertAt+1 < len(r.moq.resultsByParams) {
				copy(r.moq.resultsByParams[insertAt+1:], r.moq.resultsByParams[insertAt:0])
				r.moq.resultsByParams[insertAt] = *results
			}
		}

		var sParamUsed string
		if r.anyParams&(1<<0) == 0 {
			sParamUsed = r.paramsKey.sParam
		}
		var bParamUsed bool
		if r.anyParams&(1<<1) == 0 {
			bParamUsed = r.paramsKey.bParam
		}
		paramsKey := moqNoResultsFn_paramsKey{
			sParam: sParamUsed,
			bParam: bParamUsed,
		}

		var ok bool
		r.results, ok = results.results[paramsKey]
		if !ok {
			r.results = &moqNoResultsFn_results{
				params:   r.params,
				results:  nil,
				index:    0,
				anyTimes: false,
			}
			results.results[paramsKey] = r.results
		}
	}
}

func (r *moqNoResultsFn_fnRecorder) repeat(repeaters ...moq.Repeater) *moqNoResultsFn_fnRecorder {
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults or doReturnResults must be called before calling repeat")
		return nil
	}
	repeat := moq.Repeat(r.moq.scene.T, repeaters)
	last := r.results.results[len(r.results.results)-1]
	for n := 0; n < repeat.MaxTimes-1; n++ {
		if last.sequence != 0 {
			last = struct {
				values *struct {
				}
				sequence   uint32
				doFn       moqNoResultsFn_doFn
				doReturnFn moqNoResultsFn_doReturnFn
			}{
				values: &struct {
				}{},
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	r.results.anyTimes = repeat.AnyTimes
	return r
}

// Reset resets the state of the moq
func (m *moqNoResultsFn) Reset() { m.resultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqNoResultsFn) AssertExpectationsMet() {
	for _, res := range m.resultsByParams {
		for _, results := range res.results {
			missing := len(results.results) - int(atomic.LoadUint32(&results.index))
			if missing == 1 && results.anyTimes == true {
				continue
			}
			if missing > 0 {
				m.scene.T.Errorf("Expected %d additional call(s) with parameters %#v", missing, results.params)
			}
		}
	}
}
