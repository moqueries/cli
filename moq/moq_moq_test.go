// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package moq_test

import (
	"fmt"
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/moq"
)

// moqMoq holds the state of a moq of the Moq type
type moqMoq struct {
	scene  *moq.Scene
	config moq.Config
	moq    *moqMoq_mock

	resultsByParams_Reset                 []moqMoq_Reset_resultsByParams
	resultsByParams_AssertExpectationsMet []moqMoq_AssertExpectationsMet_resultsByParams

	runtime struct {
		parameterIndexing struct {
			Reset                 struct{}
			AssertExpectationsMet struct{}
		}
	}
}

// moqMoq_mock isolates the mock interface of the Moq type
type moqMoq_mock struct {
	moq *moqMoq
}

// moqMoq_recorder isolates the recorder interface of the Moq type
type moqMoq_recorder struct {
	moq *moqMoq
}

// moqMoq_Reset_params holds the params of the Moq type
type moqMoq_Reset_params struct{}

// moqMoq_Reset_paramsKey holds the map key params of the Moq type
type moqMoq_Reset_paramsKey struct {
	params struct{}
	hashes struct{}
}

// moqMoq_Reset_resultsByParams contains the results for a given set of
// parameters for the Moq type
type moqMoq_Reset_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqMoq_Reset_paramsKey]*moqMoq_Reset_results
}

// moqMoq_Reset_doFn defines the type of function needed when calling andDo for
// the Moq type
type moqMoq_Reset_doFn func()

// moqMoq_Reset_doReturnFn defines the type of function needed when calling
// doReturnResults for the Moq type
type moqMoq_Reset_doReturnFn func()

// moqMoq_Reset_results holds the results of the Moq type
type moqMoq_Reset_results struct {
	params  moqMoq_Reset_params
	results []struct {
		values     *struct{}
		sequence   uint32
		doFn       moqMoq_Reset_doFn
		doReturnFn moqMoq_Reset_doReturnFn
	}
	index  uint32
	repeat *moq.RepeatVal
}

// moqMoq_Reset_fnRecorder routes recorded function calls to the moqMoq moq
type moqMoq_Reset_fnRecorder struct {
	params    moqMoq_Reset_params
	anyParams uint64
	sequence  bool
	results   *moqMoq_Reset_results
	moq       *moqMoq
}

// moqMoq_Reset_anyParams isolates the any params functions of the Moq type
type moqMoq_Reset_anyParams struct {
	recorder *moqMoq_Reset_fnRecorder
}

// moqMoq_AssertExpectationsMet_params holds the params of the Moq type
type moqMoq_AssertExpectationsMet_params struct{}

// moqMoq_AssertExpectationsMet_paramsKey holds the map key params of the Moq
// type
type moqMoq_AssertExpectationsMet_paramsKey struct {
	params struct{}
	hashes struct{}
}

// moqMoq_AssertExpectationsMet_resultsByParams contains the results for a
// given set of parameters for the Moq type
type moqMoq_AssertExpectationsMet_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqMoq_AssertExpectationsMet_paramsKey]*moqMoq_AssertExpectationsMet_results
}

// moqMoq_AssertExpectationsMet_doFn defines the type of function needed when
// calling andDo for the Moq type
type moqMoq_AssertExpectationsMet_doFn func()

// moqMoq_AssertExpectationsMet_doReturnFn defines the type of function needed
// when calling doReturnResults for the Moq type
type moqMoq_AssertExpectationsMet_doReturnFn func()

// moqMoq_AssertExpectationsMet_results holds the results of the Moq type
type moqMoq_AssertExpectationsMet_results struct {
	params  moqMoq_AssertExpectationsMet_params
	results []struct {
		values     *struct{}
		sequence   uint32
		doFn       moqMoq_AssertExpectationsMet_doFn
		doReturnFn moqMoq_AssertExpectationsMet_doReturnFn
	}
	index  uint32
	repeat *moq.RepeatVal
}

// moqMoq_AssertExpectationsMet_fnRecorder routes recorded function calls to
// the moqMoq moq
type moqMoq_AssertExpectationsMet_fnRecorder struct {
	params    moqMoq_AssertExpectationsMet_params
	anyParams uint64
	sequence  bool
	results   *moqMoq_AssertExpectationsMet_results
	moq       *moqMoq
}

// moqMoq_AssertExpectationsMet_anyParams isolates the any params functions of
// the Moq type
type moqMoq_AssertExpectationsMet_anyParams struct {
	recorder *moqMoq_AssertExpectationsMet_fnRecorder
}

// newMoqMoq creates a new moq of the Moq type
func newMoqMoq(scene *moq.Scene, config *moq.Config) *moqMoq {
	if config == nil {
		config = &moq.Config{}
	}
	m := &moqMoq{
		scene:  scene,
		config: *config,
		moq:    &moqMoq_mock{},

		runtime: struct {
			parameterIndexing struct {
				Reset                 struct{}
				AssertExpectationsMet struct{}
			}
		}{parameterIndexing: struct {
			Reset                 struct{}
			AssertExpectationsMet struct{}
		}{
			Reset:                 struct{}{},
			AssertExpectationsMet: struct{}{},
		}},
	}
	m.moq.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the mock implementation of the Moq type
func (m *moqMoq) mock() *moqMoq_mock { return m.moq }

func (m *moqMoq_mock) Reset() {
	m.moq.scene.T.Helper()
	params := moqMoq_Reset_params{}
	var results *moqMoq_Reset_results
	for _, resultsByParams := range m.moq.resultsByParams_Reset {
		paramsKey := m.moq.paramsKey_Reset(params, resultsByParams.anyParams)
		var ok bool
		results, ok = resultsByParams.results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.moq.config.Expectation == moq.Strict {
			m.moq.scene.T.Fatalf("Unexpected call to %s", m.moq.prettyParams_Reset(params))
		}
		return
	}

	i := int(atomic.AddUint32(&results.index, 1)) - 1
	if i >= results.repeat.ResultCount {
		if !results.repeat.AnyTimes {
			if m.moq.config.Expectation == moq.Strict {
				m.moq.scene.T.Fatalf("Too many calls to %s", m.moq.prettyParams_Reset(params))
			}
			return
		}
		i = results.repeat.ResultCount - 1
	}

	result := results.results[i]
	if result.sequence != 0 {
		sequence := m.moq.scene.NextMockSequence()
		if (!results.repeat.AnyTimes && result.sequence != sequence) || result.sequence > sequence {
			m.moq.scene.T.Fatalf("Call sequence does not match call to %s", m.moq.prettyParams_Reset(params))
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

func (m *moqMoq_mock) AssertExpectationsMet() {
	m.moq.scene.T.Helper()
	params := moqMoq_AssertExpectationsMet_params{}
	var results *moqMoq_AssertExpectationsMet_results
	for _, resultsByParams := range m.moq.resultsByParams_AssertExpectationsMet {
		paramsKey := m.moq.paramsKey_AssertExpectationsMet(params, resultsByParams.anyParams)
		var ok bool
		results, ok = resultsByParams.results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.moq.config.Expectation == moq.Strict {
			m.moq.scene.T.Fatalf("Unexpected call to %s", m.moq.prettyParams_AssertExpectationsMet(params))
		}
		return
	}

	i := int(atomic.AddUint32(&results.index, 1)) - 1
	if i >= results.repeat.ResultCount {
		if !results.repeat.AnyTimes {
			if m.moq.config.Expectation == moq.Strict {
				m.moq.scene.T.Fatalf("Too many calls to %s", m.moq.prettyParams_AssertExpectationsMet(params))
			}
			return
		}
		i = results.repeat.ResultCount - 1
	}

	result := results.results[i]
	if result.sequence != 0 {
		sequence := m.moq.scene.NextMockSequence()
		if (!results.repeat.AnyTimes && result.sequence != sequence) || result.sequence > sequence {
			m.moq.scene.T.Fatalf("Call sequence does not match call to %s", m.moq.prettyParams_AssertExpectationsMet(params))
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

// onCall returns the recorder implementation of the Moq type
func (m *moqMoq) onCall() *moqMoq_recorder {
	return &moqMoq_recorder{
		moq: m,
	}
}

func (m *moqMoq_recorder) Reset() *moqMoq_Reset_fnRecorder {
	return &moqMoq_Reset_fnRecorder{
		params:   moqMoq_Reset_params{},
		sequence: m.moq.config.Sequence == moq.SeqDefaultOn,
		moq:      m.moq,
	}
}

func (r *moqMoq_Reset_fnRecorder) any() *moqMoq_Reset_anyParams {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_Reset(r.params))
		return nil
	}
	return &moqMoq_Reset_anyParams{recorder: r}
}

func (r *moqMoq_Reset_fnRecorder) seq() *moqMoq_Reset_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_Reset(r.params))
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqMoq_Reset_fnRecorder) noSeq() *moqMoq_Reset_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_Reset(r.params))
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqMoq_Reset_fnRecorder) returnResults() *moqMoq_Reset_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values     *struct{}
		sequence   uint32
		doFn       moqMoq_Reset_doFn
		doReturnFn moqMoq_Reset_doReturnFn
	}{
		values:   &struct{}{},
		sequence: sequence,
	})
	return r
}

func (r *moqMoq_Reset_fnRecorder) andDo(fn moqMoq_Reset_doFn) *moqMoq_Reset_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqMoq_Reset_fnRecorder) doReturnResults(fn moqMoq_Reset_doReturnFn) *moqMoq_Reset_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values     *struct{}
		sequence   uint32
		doFn       moqMoq_Reset_doFn
		doReturnFn moqMoq_Reset_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqMoq_Reset_fnRecorder) findResults() {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.results.repeat.Increment(r.moq.scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.anyParams)
	insertAt := -1
	var results *moqMoq_Reset_resultsByParams
	for n, res := range r.moq.resultsByParams_Reset {
		if res.anyParams == r.anyParams {
			results = &res
			break
		}
		if res.anyCount > anyCount {
			insertAt = n
		}
	}
	if results == nil {
		results = &moqMoq_Reset_resultsByParams{
			anyCount:  anyCount,
			anyParams: r.anyParams,
			results:   map[moqMoq_Reset_paramsKey]*moqMoq_Reset_results{},
		}
		r.moq.resultsByParams_Reset = append(r.moq.resultsByParams_Reset, *results)
		if insertAt != -1 && insertAt+1 < len(r.moq.resultsByParams_Reset) {
			copy(r.moq.resultsByParams_Reset[insertAt+1:], r.moq.resultsByParams_Reset[insertAt:0])
			r.moq.resultsByParams_Reset[insertAt] = *results
		}
	}

	paramsKey := r.moq.paramsKey_Reset(r.params, r.anyParams)

	var ok bool
	r.results, ok = results.results[paramsKey]
	if !ok {
		r.results = &moqMoq_Reset_results{
			params:  r.params,
			results: nil,
			index:   0,
			repeat:  &moq.RepeatVal{},
		}
		results.results[paramsKey] = r.results
	}

	r.results.repeat.Increment(r.moq.scene.T)
}

func (r *moqMoq_Reset_fnRecorder) repeat(repeaters ...moq.Repeater) *moqMoq_Reset_fnRecorder {
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
				doFn       moqMoq_Reset_doFn
				doReturnFn moqMoq_Reset_doReturnFn
			}{
				values:   last.values,
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (m *moqMoq) prettyParams_Reset(params moqMoq_Reset_params) string { return fmt.Sprintf("Reset()") }

func (m *moqMoq) paramsKey_Reset(params moqMoq_Reset_params, anyParams uint64) moqMoq_Reset_paramsKey {
	m.scene.T.Helper()
	return moqMoq_Reset_paramsKey{
		params: struct{}{},
		hashes: struct{}{},
	}
}

func (m *moqMoq_recorder) AssertExpectationsMet() *moqMoq_AssertExpectationsMet_fnRecorder {
	return &moqMoq_AssertExpectationsMet_fnRecorder{
		params:   moqMoq_AssertExpectationsMet_params{},
		sequence: m.moq.config.Sequence == moq.SeqDefaultOn,
		moq:      m.moq,
	}
}

func (r *moqMoq_AssertExpectationsMet_fnRecorder) any() *moqMoq_AssertExpectationsMet_anyParams {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_AssertExpectationsMet(r.params))
		return nil
	}
	return &moqMoq_AssertExpectationsMet_anyParams{recorder: r}
}

func (r *moqMoq_AssertExpectationsMet_fnRecorder) seq() *moqMoq_AssertExpectationsMet_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_AssertExpectationsMet(r.params))
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqMoq_AssertExpectationsMet_fnRecorder) noSeq() *moqMoq_AssertExpectationsMet_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_AssertExpectationsMet(r.params))
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqMoq_AssertExpectationsMet_fnRecorder) returnResults() *moqMoq_AssertExpectationsMet_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values     *struct{}
		sequence   uint32
		doFn       moqMoq_AssertExpectationsMet_doFn
		doReturnFn moqMoq_AssertExpectationsMet_doReturnFn
	}{
		values:   &struct{}{},
		sequence: sequence,
	})
	return r
}

func (r *moqMoq_AssertExpectationsMet_fnRecorder) andDo(fn moqMoq_AssertExpectationsMet_doFn) *moqMoq_AssertExpectationsMet_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqMoq_AssertExpectationsMet_fnRecorder) doReturnResults(fn moqMoq_AssertExpectationsMet_doReturnFn) *moqMoq_AssertExpectationsMet_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values     *struct{}
		sequence   uint32
		doFn       moqMoq_AssertExpectationsMet_doFn
		doReturnFn moqMoq_AssertExpectationsMet_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqMoq_AssertExpectationsMet_fnRecorder) findResults() {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.results.repeat.Increment(r.moq.scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.anyParams)
	insertAt := -1
	var results *moqMoq_AssertExpectationsMet_resultsByParams
	for n, res := range r.moq.resultsByParams_AssertExpectationsMet {
		if res.anyParams == r.anyParams {
			results = &res
			break
		}
		if res.anyCount > anyCount {
			insertAt = n
		}
	}
	if results == nil {
		results = &moqMoq_AssertExpectationsMet_resultsByParams{
			anyCount:  anyCount,
			anyParams: r.anyParams,
			results:   map[moqMoq_AssertExpectationsMet_paramsKey]*moqMoq_AssertExpectationsMet_results{},
		}
		r.moq.resultsByParams_AssertExpectationsMet = append(r.moq.resultsByParams_AssertExpectationsMet, *results)
		if insertAt != -1 && insertAt+1 < len(r.moq.resultsByParams_AssertExpectationsMet) {
			copy(r.moq.resultsByParams_AssertExpectationsMet[insertAt+1:], r.moq.resultsByParams_AssertExpectationsMet[insertAt:0])
			r.moq.resultsByParams_AssertExpectationsMet[insertAt] = *results
		}
	}

	paramsKey := r.moq.paramsKey_AssertExpectationsMet(r.params, r.anyParams)

	var ok bool
	r.results, ok = results.results[paramsKey]
	if !ok {
		r.results = &moqMoq_AssertExpectationsMet_results{
			params:  r.params,
			results: nil,
			index:   0,
			repeat:  &moq.RepeatVal{},
		}
		results.results[paramsKey] = r.results
	}

	r.results.repeat.Increment(r.moq.scene.T)
}

func (r *moqMoq_AssertExpectationsMet_fnRecorder) repeat(repeaters ...moq.Repeater) *moqMoq_AssertExpectationsMet_fnRecorder {
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
				doFn       moqMoq_AssertExpectationsMet_doFn
				doReturnFn moqMoq_AssertExpectationsMet_doReturnFn
			}{
				values:   last.values,
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (m *moqMoq) prettyParams_AssertExpectationsMet(params moqMoq_AssertExpectationsMet_params) string {
	return fmt.Sprintf("AssertExpectationsMet()")
}

func (m *moqMoq) paramsKey_AssertExpectationsMet(params moqMoq_AssertExpectationsMet_params, anyParams uint64) moqMoq_AssertExpectationsMet_paramsKey {
	m.scene.T.Helper()
	return moqMoq_AssertExpectationsMet_paramsKey{
		params: struct{}{},
		hashes: struct{}{},
	}
}

// Reset resets the state of the moq
func (m *moqMoq) Reset() {
	m.resultsByParams_Reset = nil
	m.resultsByParams_AssertExpectationsMet = nil
}

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqMoq) AssertExpectationsMet() {
	m.scene.T.Helper()
	for _, res := range m.resultsByParams_Reset {
		for _, results := range res.results {
			missing := results.repeat.MinTimes - int(atomic.LoadUint32(&results.index))
			if missing > 0 {
				m.scene.T.Errorf("Expected %d additional call(s) to %s", missing, m.prettyParams_Reset(results.params))
			}
		}
	}
	for _, res := range m.resultsByParams_AssertExpectationsMet {
		for _, results := range res.results {
			missing := results.repeat.MinTimes - int(atomic.LoadUint32(&results.index))
			if missing > 0 {
				m.scene.T.Errorf("Expected %d additional call(s) to %s", missing, m.prettyParams_AssertExpectationsMet(results.params))
			}
		}
	}
}
