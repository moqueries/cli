// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package generator_test

import (
	"sync/atomic"

	"github.com/dave/dst"
	"github.com/myshkin5/moqueries/pkg/generator"
	"github.com/myshkin5/moqueries/pkg/testing"
)

// mockLoadTypesFn holds the state of a mock of the LoadTypesFn type
type mockLoadTypesFn struct {
	t               testing.MoqT
	resultsByParams map[mockLoadTypesFn_params]*mockLoadTypesFn_resultMgr
	params          chan mockLoadTypesFn_params
}

// mockLoadTypesFn_mock isolates the mock interface of the LoadTypesFn type
type mockLoadTypesFn_mock struct {
	mock *mockLoadTypesFn
}

// mockLoadTypesFn_recorder isolates the recorder interface of the LoadTypesFn type
type mockLoadTypesFn_recorder struct {
	mock *mockLoadTypesFn
}

// mockLoadTypesFn_params holds the params of the LoadTypesFn type
type mockLoadTypesFn_params struct {
	pkg           string
	loadTestTypes bool
}

// mockLoadTypesFn_resultMgr manages multiple results and the state of the LoadTypesFn type
type mockLoadTypesFn_resultMgr struct {
	results  []*mockLoadTypesFn_results
	index    uint32
	anyTimes bool
}

// mockLoadTypesFn_results holds the results of the LoadTypesFn type
type mockLoadTypesFn_results struct {
	typeSpecs []*dst.TypeSpec
	pkgPath   string
	err       error
}

// mockLoadTypesFn_fnRecorder routes recorded function calls to the mockLoadTypesFn mock
type mockLoadTypesFn_fnRecorder struct {
	params  mockLoadTypesFn_params
	results *mockLoadTypesFn_resultMgr
	mock    *mockLoadTypesFn
}

// newMockLoadTypesFn creates a new mock of the LoadTypesFn type
func newMockLoadTypesFn(t testing.MoqT) *mockLoadTypesFn {
	return &mockLoadTypesFn{
		t:               t,
		resultsByParams: map[mockLoadTypesFn_params]*mockLoadTypesFn_resultMgr{},
		params:          make(chan mockLoadTypesFn_params, 100),
	}
}

// mock returns the mock implementation of the LoadTypesFn type
func (m *mockLoadTypesFn) mock() generator.LoadTypesFn {
	return func(pkg string, loadTestTypes bool) (typeSpecs []*dst.TypeSpec, pkgPath string, err error) {
		mock := &mockLoadTypesFn_mock{mock: m}
		return mock.fn(pkg, loadTestTypes)
	}
}

func (m *mockLoadTypesFn_mock) fn(pkg string, loadTestTypes bool) (typeSpecs []*dst.TypeSpec, pkgPath string, err error) {
	params := mockLoadTypesFn_params{
		pkg:           pkg,
		loadTestTypes: loadTestTypes,
	}
	m.mock.params <- params
	results, ok := m.mock.resultsByParams[params]
	if ok {
		i := int(atomic.AddUint32(&results.index, 1)) - 1
		if i >= len(results.results) {
			if !results.anyTimes {
				m.mock.t.Fatalf("Too many calls to mock with parameters %#v", params)
				return
			}
			i = len(results.results) - 1
		}
		result := results.results[i]
		typeSpecs = result.typeSpecs
		pkgPath = result.pkgPath
		err = result.err
	}
	return typeSpecs, pkgPath, err
}

func (m *mockLoadTypesFn) onCall(pkg string, loadTestTypes bool) *mockLoadTypesFn_fnRecorder {
	return &mockLoadTypesFn_fnRecorder{
		params: mockLoadTypesFn_params{
			pkg:           pkg,
			loadTestTypes: loadTestTypes,
		},
		mock: m,
	}
}

func (r *mockLoadTypesFn_fnRecorder) returnResults(typeSpecs []*dst.TypeSpec, pkgPath string, err error) *mockLoadTypesFn_fnRecorder {
	if r.results == nil {
		if _, ok := r.mock.resultsByParams[r.params]; ok {
			r.mock.t.Fatalf("Expectations already recorded for mock with parameters %#v", r.params)
			return nil
		}

		r.results = &mockLoadTypesFn_resultMgr{results: []*mockLoadTypesFn_results{}, index: 0, anyTimes: false}
		r.mock.resultsByParams[r.params] = r.results
	}
	r.results.results = append(r.results.results, &mockLoadTypesFn_results{
		typeSpecs: typeSpecs,
		pkgPath:   pkgPath,
		err:       err,
	})
	return r
}

func (r *mockLoadTypesFn_fnRecorder) times(count int) *mockLoadTypesFn_fnRecorder {
	if r.results == nil {
		r.mock.t.Fatalf("Return must be called before calling Times")
		return nil
	}
	last := r.results.results[len(r.results.results)-1]
	for n := 0; n < count-1; n++ {
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (r *mockLoadTypesFn_fnRecorder) anyTimes() {
	if r.results == nil {
		r.mock.t.Fatalf("Return must be called before calling AnyTimes")
		return
	}
	r.results.anyTimes = true
}
