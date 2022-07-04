// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package generator_test

import (
	"fmt"
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/moq"
)

// moqTestInterface holds the state of a moq of the testInterface type
type moqTestInterface struct {
	scene  *moq.Scene
	config moq.Config
	moq    *moqTestInterface_mock

	resultsByParams_something []moqTestInterface_something_resultsByParams

	runtime struct {
		parameterIndexing struct {
			something struct{}
		}
	}
}

// moqTestInterface_mock isolates the mock interface of the testInterface type
type moqTestInterface_mock struct {
	moq *moqTestInterface
}

// moqTestInterface_recorder isolates the recorder interface of the testInterface type
type moqTestInterface_recorder struct {
	moq *moqTestInterface
}

// moqTestInterface_something_params holds the params of the testInterface type
type moqTestInterface_something_params struct{}

// moqTestInterface_something_paramsKey holds the map key params of the testInterface type
type moqTestInterface_something_paramsKey struct {
	params struct{}
	hashes struct{}
}

// moqTestInterface_something_resultsByParams contains the results for a given set of parameters for the testInterface type
type moqTestInterface_something_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqTestInterface_something_paramsKey]*moqTestInterface_something_results
}

// moqTestInterface_something_doFn defines the type of function needed when calling andDo for the testInterface type
type moqTestInterface_something_doFn func()

// moqTestInterface_something_doReturnFn defines the type of function needed when calling doReturnResults for the testInterface type
type moqTestInterface_something_doReturnFn func()

// moqTestInterface_something_results holds the results of the testInterface type
type moqTestInterface_something_results struct {
	params  moqTestInterface_something_params
	results []struct {
		values     *struct{}
		sequence   uint32
		doFn       moqTestInterface_something_doFn
		doReturnFn moqTestInterface_something_doReturnFn
	}
	index  uint32
	repeat *moq.RepeatVal
}

// moqTestInterface_something_fnRecorder routes recorded function calls to the moqTestInterface moq
type moqTestInterface_something_fnRecorder struct {
	params    moqTestInterface_something_params
	anyParams uint64
	sequence  bool
	results   *moqTestInterface_something_results
	moq       *moqTestInterface
}

// moqTestInterface_something_anyParams isolates the any params functions of the testInterface type
type moqTestInterface_something_anyParams struct {
	recorder *moqTestInterface_something_fnRecorder
}

// newMoqtestInterface creates a new moq of the testInterface type
func newMoqtestInterface(scene *moq.Scene, config *moq.Config) *moqTestInterface {
	if config == nil {
		config = &moq.Config{}
	}
	m := &moqTestInterface{
		scene:  scene,
		config: *config,
		moq:    &moqTestInterface_mock{},

		runtime: struct {
			parameterIndexing struct {
				something struct{}
			}
		}{parameterIndexing: struct {
			something struct{}
		}{
			something: struct{}{},
		}},
	}
	m.moq.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the mock implementation of the testInterface type
func (m *moqTestInterface) mock() *moqTestInterface_mock { return m.moq }

func (m *moqTestInterface_mock) something() {
	m.moq.scene.T.Helper()
	params := moqTestInterface_something_params{}
	var results *moqTestInterface_something_results
	for _, resultsByParams := range m.moq.resultsByParams_something {
		paramsKey := m.moq.paramsKey_something(params, resultsByParams.anyParams)
		var ok bool
		results, ok = resultsByParams.results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.moq.config.Expectation == moq.Strict {
			m.moq.scene.T.Fatalf("Unexpected call to %s", m.moq.prettyParams_something(params))
		}
		return
	}

	i := int(atomic.AddUint32(&results.index, 1)) - 1
	if i >= results.repeat.ResultCount {
		if !results.repeat.AnyTimes {
			if m.moq.config.Expectation == moq.Strict {
				m.moq.scene.T.Fatalf("Too many calls to %s", m.moq.prettyParams_something(params))
			}
			return
		}
		i = results.repeat.ResultCount - 1
	}

	result := results.results[i]
	if result.sequence != 0 {
		sequence := m.moq.scene.NextMockSequence()
		if (!results.repeat.AnyTimes && result.sequence != sequence) || result.sequence > sequence {
			m.moq.scene.T.Fatalf("Call sequence does not match call to %s", m.moq.prettyParams_something(params))
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

// onCall returns the recorder implementation of the testInterface type
func (m *moqTestInterface) onCall() *moqTestInterface_recorder {
	return &moqTestInterface_recorder{
		moq: m,
	}
}

func (m *moqTestInterface_recorder) something() *moqTestInterface_something_fnRecorder {
	return &moqTestInterface_something_fnRecorder{
		params:   moqTestInterface_something_params{},
		sequence: m.moq.config.Sequence == moq.SeqDefaultOn,
		moq:      m.moq,
	}
}

func (r *moqTestInterface_something_fnRecorder) any() *moqTestInterface_something_anyParams {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_something(r.params))
		return nil
	}
	return &moqTestInterface_something_anyParams{recorder: r}
}

func (r *moqTestInterface_something_fnRecorder) seq() *moqTestInterface_something_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_something(r.params))
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqTestInterface_something_fnRecorder) noSeq() *moqTestInterface_something_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_something(r.params))
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqTestInterface_something_fnRecorder) returnResults() *moqTestInterface_something_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values     *struct{}
		sequence   uint32
		doFn       moqTestInterface_something_doFn
		doReturnFn moqTestInterface_something_doReturnFn
	}{
		values:   &struct{}{},
		sequence: sequence,
	})
	return r
}

func (r *moqTestInterface_something_fnRecorder) andDo(fn moqTestInterface_something_doFn) *moqTestInterface_something_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqTestInterface_something_fnRecorder) doReturnResults(fn moqTestInterface_something_doReturnFn) *moqTestInterface_something_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values     *struct{}
		sequence   uint32
		doFn       moqTestInterface_something_doFn
		doReturnFn moqTestInterface_something_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqTestInterface_something_fnRecorder) findResults() {
	if r.results != nil {
		r.results.repeat.Increment(r.moq.scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.anyParams)
	insertAt := -1
	var results *moqTestInterface_something_resultsByParams
	for n, res := range r.moq.resultsByParams_something {
		if res.anyParams == r.anyParams {
			results = &res
			break
		}
		if res.anyCount > anyCount {
			insertAt = n
		}
	}
	if results == nil {
		results = &moqTestInterface_something_resultsByParams{
			anyCount:  anyCount,
			anyParams: r.anyParams,
			results:   map[moqTestInterface_something_paramsKey]*moqTestInterface_something_results{},
		}
		r.moq.resultsByParams_something = append(r.moq.resultsByParams_something, *results)
		if insertAt != -1 && insertAt+1 < len(r.moq.resultsByParams_something) {
			copy(r.moq.resultsByParams_something[insertAt+1:], r.moq.resultsByParams_something[insertAt:0])
			r.moq.resultsByParams_something[insertAt] = *results
		}
	}

	paramsKey := r.moq.paramsKey_something(r.params, r.anyParams)

	var ok bool
	r.results, ok = results.results[paramsKey]
	if !ok {
		r.results = &moqTestInterface_something_results{
			params:  r.params,
			results: nil,
			index:   0,
			repeat:  &moq.RepeatVal{},
		}
		results.results[paramsKey] = r.results
	}

	r.results.repeat.Increment(r.moq.scene.T)
}

func (r *moqTestInterface_something_fnRecorder) repeat(repeaters ...moq.Repeater) *moqTestInterface_something_fnRecorder {
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
				doFn       moqTestInterface_something_doFn
				doReturnFn moqTestInterface_something_doReturnFn
			}{
				values:   last.values,
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (m *moqTestInterface) prettyParams_something(params moqTestInterface_something_params) string {
	return fmt.Sprintf("something()")
}

func (m *moqTestInterface) paramsKey_something(params moqTestInterface_something_params, anyParams uint64) moqTestInterface_something_paramsKey {
	return moqTestInterface_something_paramsKey{
		params: struct{}{},
		hashes: struct{}{},
	}
}

// Reset resets the state of the moq
func (m *moqTestInterface) Reset() { m.resultsByParams_something = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqTestInterface) AssertExpectationsMet() {
	m.scene.T.Helper()
	for _, res := range m.resultsByParams_something {
		for _, results := range res.results {
			missing := results.repeat.MinTimes - int(atomic.LoadUint32(&results.index))
			if missing > 0 {
				m.scene.T.Errorf("Expected %d additional call(s) to %s", missing, m.prettyParams_something(results.params))
			}
		}
	}
}
