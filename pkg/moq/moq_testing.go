// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package moq

import (
	"sync/atomic"

	"github.com/myshkin5/moqueries/pkg/hash"
)

// MockMoqT holds the state of a mock of the MoqT type
type MockMoqT struct {
	Scene                  *Scene
	Config                 MockConfig
	ResultsByParams_Errorf map[MockMoqT_Errorf_params]*MockMoqT_Errorf_resultMgr
	ResultsByParams_Fatalf map[MockMoqT_Fatalf_params]*MockMoqT_Fatalf_resultMgr
}

// MockMoqT_mock isolates the mock interface of the MoqT type
type MockMoqT_mock struct {
	Mock *MockMoqT
}

// MockMoqT_recorder isolates the recorder interface of the MoqT type
type MockMoqT_recorder struct {
	Mock *MockMoqT
}

// MockMoqT_Errorf_params holds the params of the MoqT type
type MockMoqT_Errorf_params struct {
	Format string
	Args   hash.Hash
}

// MockMoqT_Errorf_resultMgr manages multiple results and the state of the MoqT type
type MockMoqT_Errorf_resultMgr struct {
	Results  []*MockMoqT_Errorf_results
	Index    uint32
	AnyTimes bool
}

// MockMoqT_Errorf_results holds the results of the MoqT type
type MockMoqT_Errorf_results struct {
}

// MockMoqT_Errorf_fnRecorder routes recorded function calls to the MockMoqT mock
type MockMoqT_Errorf_fnRecorder struct {
	Params  MockMoqT_Errorf_params
	Results *MockMoqT_Errorf_resultMgr
	Mock    *MockMoqT
}

// MockMoqT_Fatalf_params holds the params of the MoqT type
type MockMoqT_Fatalf_params struct {
	Format string
	Args   hash.Hash
}

// MockMoqT_Fatalf_resultMgr manages multiple results and the state of the MoqT type
type MockMoqT_Fatalf_resultMgr struct {
	Results  []*MockMoqT_Fatalf_results
	Index    uint32
	AnyTimes bool
}

// MockMoqT_Fatalf_results holds the results of the MoqT type
type MockMoqT_Fatalf_results struct {
}

// MockMoqT_Fatalf_fnRecorder routes recorded function calls to the MockMoqT mock
type MockMoqT_Fatalf_fnRecorder struct {
	Params  MockMoqT_Fatalf_params
	Results *MockMoqT_Fatalf_resultMgr
	Mock    *MockMoqT
}

// NewMockMoqT creates a new mock of the MoqT type
func NewMockMoqT(scene *Scene, config *MockConfig) *MockMoqT {
	if config == nil {
		config = &MockConfig{}
	}
	m := &MockMoqT{
		Scene:  scene,
		Config: *config,
	}
	m.Reset()
	scene.AddMock(m)
	return m
}

// Mock returns the mock implementation of the MoqT type
func (m *MockMoqT) Mock() *MockMoqT_mock {
	return &MockMoqT_mock{
		Mock: m,
	}
}

func (m *MockMoqT_mock) Errorf(format string, args ...interface{}) {
	params := MockMoqT_Errorf_params{
		Format: format,
		Args:   hash.DeepHash(args),
	}
	results, ok := m.Mock.ResultsByParams_Errorf[params]
	if !ok {
		if m.Mock.Config.Expectation == Strict {
			m.Mock.Scene.MoqT.Fatalf("Unexpected call with parameters %#v", params)
		}
		return
	}

	i := int(atomic.AddUint32(&results.Index, 1)) - 1
	if i >= len(results.Results) {
		if !results.AnyTimes {
			if m.Mock.Config.Expectation == Strict {
				m.Mock.Scene.MoqT.Fatalf("Too many calls to mock with parameters %#v", params)
			}
			return
		}
		i = len(results.Results) - 1
	}
	return
}

func (m *MockMoqT_mock) Fatalf(format string, args ...interface{}) {
	params := MockMoqT_Fatalf_params{
		Format: format,
		Args:   hash.DeepHash(args),
	}
	results, ok := m.Mock.ResultsByParams_Fatalf[params]
	if !ok {
		if m.Mock.Config.Expectation == Strict {
			m.Mock.Scene.MoqT.Fatalf("Unexpected call with parameters %#v", params)
		}
		return
	}

	i := int(atomic.AddUint32(&results.Index, 1)) - 1
	if i >= len(results.Results) {
		if !results.AnyTimes {
			if m.Mock.Config.Expectation == Strict {
				m.Mock.Scene.MoqT.Fatalf("Too many calls to mock with parameters %#v", params)
			}
			return
		}
		i = len(results.Results) - 1
	}
	return
}

// OnCall returns the recorder implementation of the MoqT type
func (m *MockMoqT) OnCall() *MockMoqT_recorder {
	return &MockMoqT_recorder{
		Mock: m,
	}
}

func (m *MockMoqT_recorder) Errorf(format string, args ...interface{}) *MockMoqT_Errorf_fnRecorder {
	return &MockMoqT_Errorf_fnRecorder{
		Params: MockMoqT_Errorf_params{
			Format: format,
			Args:   hash.DeepHash(args),
		},
		Mock: m.Mock,
	}
}

func (r *MockMoqT_Errorf_fnRecorder) ReturnResults() *MockMoqT_Errorf_fnRecorder {
	if r.Results == nil {
		if _, ok := r.Mock.ResultsByParams_Errorf[r.Params]; ok {
			r.Mock.Scene.MoqT.Fatalf("Expectations already recorded for mock with parameters %#v", r.Params)
			return nil
		}

		r.Results = &MockMoqT_Errorf_resultMgr{Results: []*MockMoqT_Errorf_results{}, Index: 0, AnyTimes: false}
		r.Mock.ResultsByParams_Errorf[r.Params] = r.Results
	}
	r.Results.Results = append(r.Results.Results, &MockMoqT_Errorf_results{})
	return r
}

func (r *MockMoqT_Errorf_fnRecorder) Times(count int) *MockMoqT_Errorf_fnRecorder {
	if r.Results == nil {
		r.Mock.Scene.MoqT.Fatalf("Return must be called before calling Times")
		return nil
	}
	last := r.Results.Results[len(r.Results.Results)-1]
	for n := 0; n < count-1; n++ {
		r.Results.Results = append(r.Results.Results, last)
	}
	return r
}

func (r *MockMoqT_Errorf_fnRecorder) AnyTimes() {
	if r.Results == nil {
		r.Mock.Scene.MoqT.Fatalf("Return must be called before calling AnyTimes")
		return
	}
	r.Results.AnyTimes = true
}

func (m *MockMoqT_recorder) Fatalf(format string, args ...interface{}) *MockMoqT_Fatalf_fnRecorder {
	return &MockMoqT_Fatalf_fnRecorder{
		Params: MockMoqT_Fatalf_params{
			Format: format,
			Args:   hash.DeepHash(args),
		},
		Mock: m.Mock,
	}
}

func (r *MockMoqT_Fatalf_fnRecorder) ReturnResults() *MockMoqT_Fatalf_fnRecorder {
	if r.Results == nil {
		if _, ok := r.Mock.ResultsByParams_Fatalf[r.Params]; ok {
			r.Mock.Scene.MoqT.Fatalf("Expectations already recorded for mock with parameters %#v", r.Params)
			return nil
		}

		r.Results = &MockMoqT_Fatalf_resultMgr{Results: []*MockMoqT_Fatalf_results{}, Index: 0, AnyTimes: false}
		r.Mock.ResultsByParams_Fatalf[r.Params] = r.Results
	}
	r.Results.Results = append(r.Results.Results, &MockMoqT_Fatalf_results{})
	return r
}

func (r *MockMoqT_Fatalf_fnRecorder) Times(count int) *MockMoqT_Fatalf_fnRecorder {
	if r.Results == nil {
		r.Mock.Scene.MoqT.Fatalf("Return must be called before calling Times")
		return nil
	}
	last := r.Results.Results[len(r.Results.Results)-1]
	for n := 0; n < count-1; n++ {
		r.Results.Results = append(r.Results.Results, last)
	}
	return r
}

func (r *MockMoqT_Fatalf_fnRecorder) AnyTimes() {
	if r.Results == nil {
		r.Mock.Scene.MoqT.Fatalf("Return must be called before calling AnyTimes")
		return
	}
	r.Results.AnyTimes = true
}

// Reset resets the state of the mock
func (m *MockMoqT) Reset() {
	m.ResultsByParams_Errorf = map[MockMoqT_Errorf_params]*MockMoqT_Errorf_resultMgr{}
	m.ResultsByParams_Fatalf = map[MockMoqT_Fatalf_params]*MockMoqT_Fatalf_resultMgr{}
}

// AssertExpectationsMet asserts that all expectations have been met
func (m *MockMoqT) AssertExpectationsMet() {
	for params, results := range m.ResultsByParams_Errorf {
		missing := len(results.Results) - int(atomic.LoadUint32(&results.Index))
		if missing == 1 && results.AnyTimes == true {
			continue
		}
		if missing > 0 {
			m.Scene.MoqT.Errorf("Expected %d additional call(s) with parameters %#v", missing, params)
		}
	}
	for params, results := range m.ResultsByParams_Fatalf {
		missing := len(results.Results) - int(atomic.LoadUint32(&results.Index))
		if missing == 1 && results.AnyTimes == true {
			continue
		}
		if missing > 0 {
			m.Scene.MoqT.Errorf("Expected %d additional call(s) with parameters %#v", missing, params)
		}
	}
}