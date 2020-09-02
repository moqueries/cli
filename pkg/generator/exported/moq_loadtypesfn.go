// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package exported

import (
	"github.com/dave/dst"
	"github.com/myshkin5/moqueries/pkg/generator"
)

// MockLoadTypesFn holds the state of a mock of the LoadTypesFn type
type MockLoadTypesFn struct {
	ResultsByParams map[MockLoadTypesFn_params]MockLoadTypesFn_results
	Params          chan MockLoadTypesFn_params
}

// MockLoadTypesFn_mock isolates the mock interface of the LoadTypesFn type
type MockLoadTypesFn_mock struct {
	Mock *MockLoadTypesFn
}

// MockLoadTypesFn_recorder isolates the recorder interface of the LoadTypesFn type
type MockLoadTypesFn_recorder struct {
	Mock *MockLoadTypesFn
}

// MockLoadTypesFn_params holds the params of the LoadTypesFn type
type MockLoadTypesFn_params struct {
	Pkg           string
	LoadTestTypes bool
}

// MockLoadTypesFn_results holds the results of the LoadTypesFn type
type MockLoadTypesFn_results struct {
	TypeSpecs []*dst.TypeSpec
	PkgPath   string
	Err       error
}

// MockLoadTypesFn_fnRecorder routes recorded function calls to the MockLoadTypesFn mock
type MockLoadTypesFn_fnRecorder struct {
	Params MockLoadTypesFn_params
	Mock   *MockLoadTypesFn
}

// NewMockLoadTypesFn creates a new mock of the LoadTypesFn type
func NewMockLoadTypesFn() *MockLoadTypesFn {
	return &MockLoadTypesFn{
		ResultsByParams: map[MockLoadTypesFn_params]MockLoadTypesFn_results{},
		Params:          make(chan MockLoadTypesFn_params, 100),
	}
}

// Mock returns the mock implementation of the LoadTypesFn type
func (m *MockLoadTypesFn) Mock() generator.LoadTypesFn {
	return func(pkg string, loadTestTypes bool) (typeSpecs []*dst.TypeSpec, pkgPath string, err error) {
		mock := &MockLoadTypesFn_mock{Mock: m}
		return mock.Fn(pkg, loadTestTypes)
	}
}

func (m *MockLoadTypesFn_mock) Fn(pkg string, loadTestTypes bool) (typeSpecs []*dst.TypeSpec, pkgPath string, err error) {
	params := MockLoadTypesFn_params{
		Pkg:           pkg,
		LoadTestTypes: loadTestTypes,
	}
	m.Mock.Params <- params
	results, ok := m.Mock.ResultsByParams[params]
	if ok {
		typeSpecs = results.TypeSpecs
		pkgPath = results.PkgPath
		err = results.Err
	}
	return typeSpecs, pkgPath, err
}

func (m *MockLoadTypesFn) OnCall(pkg string, loadTestTypes bool) *MockLoadTypesFn_fnRecorder {
	return &MockLoadTypesFn_fnRecorder{
		Params: MockLoadTypesFn_params{
			Pkg:           pkg,
			LoadTestTypes: loadTestTypes,
		},
		Mock: m,
	}
}

func (r *MockLoadTypesFn_fnRecorder) Ret(typeSpecs []*dst.TypeSpec, pkgPath string, err error) {
	r.Mock.ResultsByParams[r.Params] = MockLoadTypesFn_results{
		TypeSpecs: typeSpecs,
		PkgPath:   pkgPath,
		Err:       err,
	}
}
