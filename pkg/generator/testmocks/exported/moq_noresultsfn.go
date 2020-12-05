// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package exported

import (
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/pkg/generator/testmocks"
	"github.com/myshkin5/moqueries/pkg/moq"
)

// MockNoResultsFn holds the state of a mock of the NoResultsFn type
type MockNoResultsFn struct {
	Scene           *moq.Scene
	Config          moq.MockConfig
	ResultsByParams []MockNoResultsFn_resultsByParams
}

// MockNoResultsFn_mock isolates the mock interface of the NoResultsFn type
type MockNoResultsFn_mock struct {
	Mock *MockNoResultsFn
}

// MockNoResultsFn_params holds the params of the NoResultsFn type
type MockNoResultsFn_params struct {
	SParam string
	BParam bool
}

// MockNoResultsFn_paramsKey holds the map key params of the NoResultsFn type
type MockNoResultsFn_paramsKey struct {
	SParam string
	BParam bool
}

// MockNoResultsFn_resultsByParams contains the results for a given set of parameters for the NoResultsFn type
type MockNoResultsFn_resultsByParams struct {
	AnyCount  int
	AnyParams uint64
	Results   map[MockNoResultsFn_paramsKey]*MockNoResultsFn_results
}

// MockNoResultsFn_doFn defines the type of function needed when calling AndDo for the NoResultsFn type
type MockNoResultsFn_doFn func(sParam string, bParam bool)

// MockNoResultsFn_doReturnFn defines the type of function needed when calling DoReturnResults for the NoResultsFn type
type MockNoResultsFn_doReturnFn func(sParam string, bParam bool)

// MockNoResultsFn_results holds the results of the NoResultsFn type
type MockNoResultsFn_results struct {
	Params  MockNoResultsFn_params
	Results []struct {
		Values *struct {
		}
		Sequence   uint32
		DoFn       MockNoResultsFn_doFn
		DoReturnFn MockNoResultsFn_doReturnFn
	}
	Index    uint32
	AnyTimes bool
}

// MockNoResultsFn_fnRecorder routes recorded function calls to the MockNoResultsFn mock
type MockNoResultsFn_fnRecorder struct {
	Params    MockNoResultsFn_params
	ParamsKey MockNoResultsFn_paramsKey
	AnyParams uint64
	Sequence  bool
	Results   *MockNoResultsFn_results
	Mock      *MockNoResultsFn
}

// NewMockNoResultsFn creates a new mock of the NoResultsFn type
func NewMockNoResultsFn(scene *moq.Scene, config *moq.MockConfig) *MockNoResultsFn {
	if config == nil {
		config = &moq.MockConfig{}
	}
	m := &MockNoResultsFn{
		Scene:  scene,
		Config: *config,
	}
	scene.AddMock(m)
	return m
}

// Mock returns the mock implementation of the NoResultsFn type
func (m *MockNoResultsFn) Mock() testmocks.NoResultsFn {
	return func(sParam string, bParam bool) { mock := &MockNoResultsFn_mock{Mock: m}; mock.Fn(sParam, bParam) }
}

func (m *MockNoResultsFn_mock) Fn(sParam string, bParam bool) {
	params := MockNoResultsFn_params{
		SParam: sParam,
		BParam: bParam,
	}
	var results *MockNoResultsFn_results
	for _, resultsByParams := range m.Mock.ResultsByParams {
		var sParamUsed string
		if resultsByParams.AnyParams&(1<<0) == 0 {
			sParamUsed = sParam
		}
		var bParamUsed bool
		if resultsByParams.AnyParams&(1<<1) == 0 {
			bParamUsed = bParam
		}
		paramsKey := MockNoResultsFn_paramsKey{
			SParam: sParamUsed,
			BParam: bParamUsed,
		}
		var ok bool
		results, ok = resultsByParams.Results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.Mock.Config.Expectation == moq.Strict {
			m.Mock.Scene.MoqT.Fatalf("Unexpected call with parameters %#v", params)
		}
		return
	}

	i := int(atomic.AddUint32(&results.Index, 1)) - 1
	if i >= len(results.Results) {
		if !results.AnyTimes {
			if m.Mock.Config.Expectation == moq.Strict {
				m.Mock.Scene.MoqT.Fatalf("Too many calls to mock with parameters %#v", params)
			}
			return
		}
		i = len(results.Results) - 1
	}

	result := results.Results[i]
	if result.Sequence != 0 {
		sequence := m.Mock.Scene.NextMockSequence()
		if (!results.AnyTimes && result.Sequence != sequence) || result.Sequence > sequence {
			m.Mock.Scene.MoqT.Fatalf("Call sequence does not match %#v", params)
		}
	}

	if result.DoFn != nil {
		result.DoFn(sParam, bParam)
	}

	if result.DoReturnFn != nil {
		result.DoReturnFn(sParam, bParam)
	}
	return
}

func (m *MockNoResultsFn) OnCall(sParam string, bParam bool) *MockNoResultsFn_fnRecorder {
	return &MockNoResultsFn_fnRecorder{
		Params: MockNoResultsFn_params{
			SParam: sParam,
			BParam: bParam,
		},
		ParamsKey: MockNoResultsFn_paramsKey{
			SParam: sParam,
			BParam: bParam,
		},
		Sequence: m.Config.Sequence == moq.SeqDefaultOn,
		Mock:     m,
	}
}

func (r *MockNoResultsFn_fnRecorder) AnySParam() *MockNoResultsFn_fnRecorder {
	if r.Results != nil {
		r.Mock.Scene.MoqT.Fatalf("Any functions must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.AnyParams |= 1 << 0
	return r
}

func (r *MockNoResultsFn_fnRecorder) AnyBParam() *MockNoResultsFn_fnRecorder {
	if r.Results != nil {
		r.Mock.Scene.MoqT.Fatalf("Any functions must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.AnyParams |= 1 << 1
	return r
}

func (r *MockNoResultsFn_fnRecorder) Seq() *MockNoResultsFn_fnRecorder {
	if r.Results != nil {
		r.Mock.Scene.MoqT.Fatalf("Seq must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.Sequence = true
	return r
}

func (r *MockNoResultsFn_fnRecorder) NoSeq() *MockNoResultsFn_fnRecorder {
	if r.Results != nil {
		r.Mock.Scene.MoqT.Fatalf("NoSeq must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.Sequence = false
	return r
}

func (r *MockNoResultsFn_fnRecorder) ReturnResults() *MockNoResultsFn_fnRecorder {
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Mock.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values *struct {
		}
		Sequence   uint32
		DoFn       MockNoResultsFn_doFn
		DoReturnFn MockNoResultsFn_doReturnFn
	}{
		Values: &struct {
		}{},
		Sequence: sequence,
	})
	return r
}

func (r *MockNoResultsFn_fnRecorder) AndDo(fn MockNoResultsFn_doFn) *MockNoResultsFn_fnRecorder {
	if r.Results == nil {
		r.Mock.Scene.MoqT.Fatalf("ReturnResults must be called before calling AndDo")
		return nil
	}
	last := &r.Results.Results[len(r.Results.Results)-1]
	last.DoFn = fn
	return r
}

func (r *MockNoResultsFn_fnRecorder) DoReturnResults(fn MockNoResultsFn_doReturnFn) *MockNoResultsFn_fnRecorder {
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Mock.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values *struct {
		}
		Sequence   uint32
		DoFn       MockNoResultsFn_doFn
		DoReturnFn MockNoResultsFn_doReturnFn
	}{Sequence: sequence, DoReturnFn: fn})
	return r
}

func (r *MockNoResultsFn_fnRecorder) FindResults() {
	if r.Results == nil {
		anyCount := bits.OnesCount64(r.AnyParams)
		insertAt := -1
		var results *MockNoResultsFn_resultsByParams
		for n, res := range r.Mock.ResultsByParams {
			if res.AnyParams == r.AnyParams {
				results = &res
				break
			}
			if res.AnyCount > anyCount {
				insertAt = n
			}
		}
		if results == nil {
			results = &MockNoResultsFn_resultsByParams{
				AnyCount:  anyCount,
				AnyParams: r.AnyParams,
				Results:   map[MockNoResultsFn_paramsKey]*MockNoResultsFn_results{},
			}
			r.Mock.ResultsByParams = append(r.Mock.ResultsByParams, *results)
			if insertAt != -1 && insertAt+1 < len(r.Mock.ResultsByParams) {
				copy(r.Mock.ResultsByParams[insertAt+1:], r.Mock.ResultsByParams[insertAt:0])
				r.Mock.ResultsByParams[insertAt] = *results
			}
		}

		var sParamUsed string
		if r.AnyParams&(1<<0) == 0 {
			sParamUsed = r.ParamsKey.SParam
		}
		var bParamUsed bool
		if r.AnyParams&(1<<1) == 0 {
			bParamUsed = r.ParamsKey.BParam
		}
		paramsKey := MockNoResultsFn_paramsKey{
			SParam: sParamUsed,
			BParam: bParamUsed,
		}

		var ok bool
		r.Results, ok = results.Results[paramsKey]
		if !ok {
			r.Results = &MockNoResultsFn_results{
				Params:   r.Params,
				Results:  nil,
				Index:    0,
				AnyTimes: false,
			}
			results.Results[paramsKey] = r.Results
		}
	}
}

func (r *MockNoResultsFn_fnRecorder) Times(count int) *MockNoResultsFn_fnRecorder {
	if r.Results == nil {
		r.Mock.Scene.MoqT.Fatalf("ReturnResults or DoReturnResults must be called before calling Times")
		return nil
	}
	last := r.Results.Results[len(r.Results.Results)-1]
	for n := 0; n < count-1; n++ {
		if last.Sequence != 0 {
			last = struct {
				Values *struct {
				}
				Sequence   uint32
				DoFn       MockNoResultsFn_doFn
				DoReturnFn MockNoResultsFn_doReturnFn
			}{
				Values: &struct {
				}{},
				Sequence: r.Mock.Scene.NextRecorderSequence(),
			}
		}
		r.Results.Results = append(r.Results.Results, last)
	}
	return r
}

func (r *MockNoResultsFn_fnRecorder) AnyTimes() {
	if r.Results == nil {
		r.Mock.Scene.MoqT.Fatalf("ReturnResults or DoReturnResults must be called before calling AnyTimes")
		return
	}
	r.Results.AnyTimes = true
}

// Reset resets the state of the mock
func (m *MockNoResultsFn) Reset() { m.ResultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *MockNoResultsFn) AssertExpectationsMet() {
	for _, res := range m.ResultsByParams {
		for _, results := range res.Results {
			missing := len(results.Results) - int(atomic.LoadUint32(&results.Index))
			if missing == 1 && results.AnyTimes == true {
				continue
			}
			if missing > 0 {
				m.Scene.MoqT.Errorf("Expected %d additional call(s) with parameters %#v", missing, results.Params)
			}
		}
	}
}
