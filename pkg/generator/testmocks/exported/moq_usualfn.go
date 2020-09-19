// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package exported

import (
	"sync/atomic"

	"github.com/myshkin5/moqueries/pkg/generator/testmocks"
	"github.com/myshkin5/moqueries/pkg/testing"
)

// MockUsualFn holds the state of a mock of the UsualFn type
type MockUsualFn struct {
	T               testing.MoqT
	ResultsByParams map[MockUsualFn_params]*MockUsualFn_resultMgr
	Params          chan MockUsualFn_params
}

// MockUsualFn_mock isolates the mock interface of the UsualFn type
type MockUsualFn_mock struct {
	Mock *MockUsualFn
}

// MockUsualFn_recorder isolates the recorder interface of the UsualFn type
type MockUsualFn_recorder struct {
	Mock *MockUsualFn
}

// MockUsualFn_params holds the params of the UsualFn type
type MockUsualFn_params struct {
	SParam string
	BParam bool
}

// MockUsualFn_resultMgr manages multiple results and the state of the UsualFn type
type MockUsualFn_resultMgr struct {
	Results  []*MockUsualFn_results
	Index    uint32
	AnyTimes bool
}

// MockUsualFn_results holds the results of the UsualFn type
type MockUsualFn_results struct {
	SResult string
	Err     error
}

// MockUsualFn_fnRecorder routes recorded function calls to the MockUsualFn mock
type MockUsualFn_fnRecorder struct {
	Params  MockUsualFn_params
	Results *MockUsualFn_resultMgr
	Mock    *MockUsualFn
}

// NewMockUsualFn creates a new mock of the UsualFn type
func NewMockUsualFn(t testing.MoqT) *MockUsualFn {
	return &MockUsualFn{
		T:               t,
		ResultsByParams: map[MockUsualFn_params]*MockUsualFn_resultMgr{},
		Params:          make(chan MockUsualFn_params, 100),
	}
}

// Mock returns the mock implementation of the UsualFn type
func (m *MockUsualFn) Mock() testmocks.UsualFn {
	return func(sParam string, bParam bool) (sResult string, err error) {
		mock := &MockUsualFn_mock{Mock: m}
		return mock.Fn(sParam, bParam)
	}
}

func (m *MockUsualFn_mock) Fn(sParam string, bParam bool) (sResult string, err error) {
	params := MockUsualFn_params{
		SParam: sParam,
		BParam: bParam,
	}
	m.Mock.Params <- params
	results, ok := m.Mock.ResultsByParams[params]
	if ok {
		i := int(atomic.AddUint32(&results.Index, 1)) - 1
		if i >= len(results.Results) {
			if !results.AnyTimes {
				m.Mock.T.Fatalf("Too many calls to mock with parameters %#v", params)
				return
			}
			i = len(results.Results) - 1
		}
		result := results.Results[i]
		sResult = result.SResult
		err = result.Err
	}
	return sResult, err
}

func (m *MockUsualFn) OnCall(sParam string, bParam bool) *MockUsualFn_fnRecorder {
	return &MockUsualFn_fnRecorder{
		Params: MockUsualFn_params{
			SParam: sParam,
			BParam: bParam,
		},
		Mock: m,
	}
}

func (r *MockUsualFn_fnRecorder) ReturnResults(sResult string, err error) *MockUsualFn_fnRecorder {
	if r.Results == nil {
		if _, ok := r.Mock.ResultsByParams[r.Params]; ok {
			r.Mock.T.Fatalf("Expectations already recorded for mock with parameters %#v", r.Params)
			return nil
		}

		r.Results = &MockUsualFn_resultMgr{Results: []*MockUsualFn_results{}, Index: 0, AnyTimes: false}
		r.Mock.ResultsByParams[r.Params] = r.Results
	}
	r.Results.Results = append(r.Results.Results, &MockUsualFn_results{
		SResult: sResult,
		Err:     err,
	})
	return r
}

func (r *MockUsualFn_fnRecorder) Times(count int) *MockUsualFn_fnRecorder {
	if r.Results == nil {
		r.Mock.T.Fatalf("Return must be called before calling Times")
		return nil
	}
	last := r.Results.Results[len(r.Results.Results)-1]
	for n := 0; n < count-1; n++ {
		r.Results.Results = append(r.Results.Results, last)
	}
	return r
}

func (r *MockUsualFn_fnRecorder) AnyTimes() {
	if r.Results == nil {
		r.Mock.T.Fatalf("Return must be called before calling AnyTimes")
		return
	}
	r.Results.AnyTimes = true
}
