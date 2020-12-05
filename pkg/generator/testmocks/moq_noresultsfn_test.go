// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package testmocks_test

import (
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/pkg/generator/testmocks"
	"github.com/myshkin5/moqueries/pkg/moq"
)

// mockNoResultsFn holds the state of a mock of the NoResultsFn type
type mockNoResultsFn struct {
	scene           *moq.Scene
	config          moq.MockConfig
	resultsByParams []mockNoResultsFn_resultsByParams
}

// mockNoResultsFn_mock isolates the mock interface of the NoResultsFn type
type mockNoResultsFn_mock struct {
	mock *mockNoResultsFn
}

// mockNoResultsFn_params holds the params of the NoResultsFn type
type mockNoResultsFn_params struct {
	sParam string
	bParam bool
}

// mockNoResultsFn_paramsKey holds the map key params of the NoResultsFn type
type mockNoResultsFn_paramsKey struct {
	sParam string
	bParam bool
}

// mockNoResultsFn_resultsByParams contains the results for a given set of parameters for the NoResultsFn type
type mockNoResultsFn_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[mockNoResultsFn_paramsKey]*mockNoResultsFn_results
}

// mockNoResultsFn_doFn defines the type of function needed when calling andDo for the NoResultsFn type
type mockNoResultsFn_doFn func(sParam string, bParam bool)

// mockNoResultsFn_doReturnFn defines the type of function needed when calling doReturnResults for the NoResultsFn type
type mockNoResultsFn_doReturnFn func(sParam string, bParam bool)

// mockNoResultsFn_results holds the results of the NoResultsFn type
type mockNoResultsFn_results struct {
	params  mockNoResultsFn_params
	results []struct {
		values *struct {
		}
		sequence   uint32
		doFn       mockNoResultsFn_doFn
		doReturnFn mockNoResultsFn_doReturnFn
	}
	index    uint32
	anyTimes bool
}

// mockNoResultsFn_fnRecorder routes recorded function calls to the mockNoResultsFn mock
type mockNoResultsFn_fnRecorder struct {
	params    mockNoResultsFn_params
	paramsKey mockNoResultsFn_paramsKey
	anyParams uint64
	sequence  bool
	results   *mockNoResultsFn_results
	mock      *mockNoResultsFn
}

// newMockNoResultsFn creates a new mock of the NoResultsFn type
func newMockNoResultsFn(scene *moq.Scene, config *moq.MockConfig) *mockNoResultsFn {
	if config == nil {
		config = &moq.MockConfig{}
	}
	m := &mockNoResultsFn{
		scene:  scene,
		config: *config,
	}
	scene.AddMock(m)
	return m
}

// mock returns the mock implementation of the NoResultsFn type
func (m *mockNoResultsFn) mock() testmocks.NoResultsFn {
	return func(sParam string, bParam bool) { mock := &mockNoResultsFn_mock{mock: m}; mock.fn(sParam, bParam) }
}

func (m *mockNoResultsFn_mock) fn(sParam string, bParam bool) {
	params := mockNoResultsFn_params{
		sParam: sParam,
		bParam: bParam,
	}
	var results *mockNoResultsFn_results
	for _, resultsByParams := range m.mock.resultsByParams {
		var sParamUsed string
		if resultsByParams.anyParams&(1<<0) == 0 {
			sParamUsed = sParam
		}
		var bParamUsed bool
		if resultsByParams.anyParams&(1<<1) == 0 {
			bParamUsed = bParam
		}
		paramsKey := mockNoResultsFn_paramsKey{
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
		if m.mock.config.Expectation == moq.Strict {
			m.mock.scene.MoqT.Fatalf("Unexpected call with parameters %#v", params)
		}
		return
	}

	i := int(atomic.AddUint32(&results.index, 1)) - 1
	if i >= len(results.results) {
		if !results.anyTimes {
			if m.mock.config.Expectation == moq.Strict {
				m.mock.scene.MoqT.Fatalf("Too many calls to mock with parameters %#v", params)
			}
			return
		}
		i = len(results.results) - 1
	}

	result := results.results[i]
	if result.sequence != 0 {
		sequence := m.mock.scene.NextMockSequence()
		if (!results.anyTimes && result.sequence != sequence) || result.sequence > sequence {
			m.mock.scene.MoqT.Fatalf("Call sequence does not match %#v", params)
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

func (m *mockNoResultsFn) onCall(sParam string, bParam bool) *mockNoResultsFn_fnRecorder {
	return &mockNoResultsFn_fnRecorder{
		params: mockNoResultsFn_params{
			sParam: sParam,
			bParam: bParam,
		},
		paramsKey: mockNoResultsFn_paramsKey{
			sParam: sParam,
			bParam: bParam,
		},
		sequence: m.config.Sequence == moq.SeqDefaultOn,
		mock:     m,
	}
}

func (r *mockNoResultsFn_fnRecorder) anySParam() *mockNoResultsFn_fnRecorder {
	if r.results != nil {
		r.mock.scene.MoqT.Fatalf("Any functions must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	r.anyParams |= 1 << 0
	return r
}

func (r *mockNoResultsFn_fnRecorder) anyBParam() *mockNoResultsFn_fnRecorder {
	if r.results != nil {
		r.mock.scene.MoqT.Fatalf("Any functions must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	r.anyParams |= 1 << 1
	return r
}

func (r *mockNoResultsFn_fnRecorder) seq() *mockNoResultsFn_fnRecorder {
	if r.results != nil {
		r.mock.scene.MoqT.Fatalf("seq must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	r.sequence = true
	return r
}

func (r *mockNoResultsFn_fnRecorder) noSeq() *mockNoResultsFn_fnRecorder {
	if r.results != nil {
		r.mock.scene.MoqT.Fatalf("noSeq must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	r.sequence = false
	return r
}

func (r *mockNoResultsFn_fnRecorder) returnResults() *mockNoResultsFn_fnRecorder {
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.mock.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
		}
		sequence   uint32
		doFn       mockNoResultsFn_doFn
		doReturnFn mockNoResultsFn_doReturnFn
	}{
		values: &struct {
		}{},
		sequence: sequence,
	})
	return r
}

func (r *mockNoResultsFn_fnRecorder) andDo(fn mockNoResultsFn_doFn) *mockNoResultsFn_fnRecorder {
	if r.results == nil {
		r.mock.scene.MoqT.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *mockNoResultsFn_fnRecorder) doReturnResults(fn mockNoResultsFn_doReturnFn) *mockNoResultsFn_fnRecorder {
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.mock.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
		}
		sequence   uint32
		doFn       mockNoResultsFn_doFn
		doReturnFn mockNoResultsFn_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *mockNoResultsFn_fnRecorder) findResults() {
	if r.results == nil {
		anyCount := bits.OnesCount64(r.anyParams)
		insertAt := -1
		var results *mockNoResultsFn_resultsByParams
		for n, res := range r.mock.resultsByParams {
			if res.anyParams == r.anyParams {
				results = &res
				break
			}
			if res.anyCount > anyCount {
				insertAt = n
			}
		}
		if results == nil {
			results = &mockNoResultsFn_resultsByParams{
				anyCount:  anyCount,
				anyParams: r.anyParams,
				results:   map[mockNoResultsFn_paramsKey]*mockNoResultsFn_results{},
			}
			r.mock.resultsByParams = append(r.mock.resultsByParams, *results)
			if insertAt != -1 && insertAt+1 < len(r.mock.resultsByParams) {
				copy(r.mock.resultsByParams[insertAt+1:], r.mock.resultsByParams[insertAt:0])
				r.mock.resultsByParams[insertAt] = *results
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
		paramsKey := mockNoResultsFn_paramsKey{
			sParam: sParamUsed,
			bParam: bParamUsed,
		}

		var ok bool
		r.results, ok = results.results[paramsKey]
		if !ok {
			r.results = &mockNoResultsFn_results{
				params:   r.params,
				results:  nil,
				index:    0,
				anyTimes: false,
			}
			results.results[paramsKey] = r.results
		}
	}
}

func (r *mockNoResultsFn_fnRecorder) times(count int) *mockNoResultsFn_fnRecorder {
	if r.results == nil {
		r.mock.scene.MoqT.Fatalf("returnResults or doReturnResults must be called before calling times")
		return nil
	}
	last := r.results.results[len(r.results.results)-1]
	for n := 0; n < count-1; n++ {
		if last.sequence != 0 {
			last = struct {
				values *struct {
				}
				sequence   uint32
				doFn       mockNoResultsFn_doFn
				doReturnFn mockNoResultsFn_doReturnFn
			}{
				values: &struct {
				}{},
				sequence: r.mock.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (r *mockNoResultsFn_fnRecorder) anyTimes() {
	if r.results == nil {
		r.mock.scene.MoqT.Fatalf("returnResults or doReturnResults must be called before calling anyTimes")
		return
	}
	r.results.anyTimes = true
}

// Reset resets the state of the mock
func (m *mockNoResultsFn) Reset() { m.resultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *mockNoResultsFn) AssertExpectationsMet() {
	for _, res := range m.resultsByParams {
		for _, results := range res.results {
			missing := len(results.results) - int(atomic.LoadUint32(&results.index))
			if missing == 1 && results.anyTimes == true {
				continue
			}
			if missing > 0 {
				m.scene.MoqT.Errorf("Expected %d additional call(s) with parameters %#v", missing, results.params)
			}
		}
	}
}
