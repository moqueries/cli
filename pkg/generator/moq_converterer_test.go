// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package generator_test

import (
	"sync/atomic"

	"github.com/dave/dst"
	"github.com/myshkin5/moqueries/pkg/generator"
	"github.com/myshkin5/moqueries/pkg/hash"
	"github.com/myshkin5/moqueries/pkg/testing"
)

// mockConverterer holds the state of a mock of the Converterer type
type mockConverterer struct {
	t                                 testing.MoqT
	resultsByParams_BaseStruct        map[mockConverterer_BaseStruct_params]*mockConverterer_BaseStruct_resultMgr
	params_BaseStruct                 chan mockConverterer_BaseStruct_params
	resultsByParams_IsolationStruct   map[mockConverterer_IsolationStruct_params]*mockConverterer_IsolationStruct_resultMgr
	params_IsolationStruct            chan mockConverterer_IsolationStruct_params
	resultsByParams_MethodStructs     map[mockConverterer_MethodStructs_params]*mockConverterer_MethodStructs_resultMgr
	params_MethodStructs              chan mockConverterer_MethodStructs_params
	resultsByParams_NewFunc           map[mockConverterer_NewFunc_params]*mockConverterer_NewFunc_resultMgr
	params_NewFunc                    chan mockConverterer_NewFunc_params
	resultsByParams_IsolationAccessor map[mockConverterer_IsolationAccessor_params]*mockConverterer_IsolationAccessor_resultMgr
	params_IsolationAccessor          chan mockConverterer_IsolationAccessor_params
	resultsByParams_FuncClosure       map[mockConverterer_FuncClosure_params]*mockConverterer_FuncClosure_resultMgr
	params_FuncClosure                chan mockConverterer_FuncClosure_params
	resultsByParams_MockMethod        map[mockConverterer_MockMethod_params]*mockConverterer_MockMethod_resultMgr
	params_MockMethod                 chan mockConverterer_MockMethod_params
	resultsByParams_RecorderMethods   map[mockConverterer_RecorderMethods_params]*mockConverterer_RecorderMethods_resultMgr
	params_RecorderMethods            chan mockConverterer_RecorderMethods_params
}

// mockConverterer_mock isolates the mock interface of the Converterer type
type mockConverterer_mock struct {
	mock *mockConverterer
}

// mockConverterer_recorder isolates the recorder interface of the Converterer type
type mockConverterer_recorder struct {
	mock *mockConverterer
}

// mockConverterer_BaseStruct_params holds the params of the Converterer type
type mockConverterer_BaseStruct_params struct {
	typeSpec *dst.TypeSpec
	funcs    hash.Hash
}

// mockConverterer_BaseStruct_resultMgr manages multiple results and the state of the Converterer type
type mockConverterer_BaseStruct_resultMgr struct {
	results  []*mockConverterer_BaseStruct_results
	index    uint32
	anyTimes bool
}

// mockConverterer_BaseStruct_results holds the results of the Converterer type
type mockConverterer_BaseStruct_results struct{ structDecl *dst.GenDecl }

// mockConverterer_BaseStruct_fnRecorder routes recorded function calls to the mockConverterer mock
type mockConverterer_BaseStruct_fnRecorder struct {
	params  mockConverterer_BaseStruct_params
	results *mockConverterer_BaseStruct_resultMgr
	mock    *mockConverterer
}

// mockConverterer_IsolationStruct_params holds the params of the Converterer type
type mockConverterer_IsolationStruct_params struct{ typeName, suffix string }

// mockConverterer_IsolationStruct_resultMgr manages multiple results and the state of the Converterer type
type mockConverterer_IsolationStruct_resultMgr struct {
	results  []*mockConverterer_IsolationStruct_results
	index    uint32
	anyTimes bool
}

// mockConverterer_IsolationStruct_results holds the results of the Converterer type
type mockConverterer_IsolationStruct_results struct{ structDecl *dst.GenDecl }

// mockConverterer_IsolationStruct_fnRecorder routes recorded function calls to the mockConverterer mock
type mockConverterer_IsolationStruct_fnRecorder struct {
	params  mockConverterer_IsolationStruct_params
	results *mockConverterer_IsolationStruct_resultMgr
	mock    *mockConverterer
}

// mockConverterer_MethodStructs_params holds the params of the Converterer type
type mockConverterer_MethodStructs_params struct {
	typeSpec *dst.TypeSpec
	fn       generator.Func
}

// mockConverterer_MethodStructs_resultMgr manages multiple results and the state of the Converterer type
type mockConverterer_MethodStructs_resultMgr struct {
	results  []*mockConverterer_MethodStructs_results
	index    uint32
	anyTimes bool
}

// mockConverterer_MethodStructs_results holds the results of the Converterer type
type mockConverterer_MethodStructs_results struct{ structDecls []dst.Decl }

// mockConverterer_MethodStructs_fnRecorder routes recorded function calls to the mockConverterer mock
type mockConverterer_MethodStructs_fnRecorder struct {
	params  mockConverterer_MethodStructs_params
	results *mockConverterer_MethodStructs_resultMgr
	mock    *mockConverterer
}

// mockConverterer_NewFunc_params holds the params of the Converterer type
type mockConverterer_NewFunc_params struct {
	typeSpec *dst.TypeSpec
	funcs    hash.Hash
}

// mockConverterer_NewFunc_resultMgr manages multiple results and the state of the Converterer type
type mockConverterer_NewFunc_resultMgr struct {
	results  []*mockConverterer_NewFunc_results
	index    uint32
	anyTimes bool
}

// mockConverterer_NewFunc_results holds the results of the Converterer type
type mockConverterer_NewFunc_results struct{ funcDecl *dst.FuncDecl }

// mockConverterer_NewFunc_fnRecorder routes recorded function calls to the mockConverterer mock
type mockConverterer_NewFunc_fnRecorder struct {
	params  mockConverterer_NewFunc_params
	results *mockConverterer_NewFunc_resultMgr
	mock    *mockConverterer
}

// mockConverterer_IsolationAccessor_params holds the params of the Converterer type
type mockConverterer_IsolationAccessor_params struct{ typeName, suffix, fnName string }

// mockConverterer_IsolationAccessor_resultMgr manages multiple results and the state of the Converterer type
type mockConverterer_IsolationAccessor_resultMgr struct {
	results  []*mockConverterer_IsolationAccessor_results
	index    uint32
	anyTimes bool
}

// mockConverterer_IsolationAccessor_results holds the results of the Converterer type
type mockConverterer_IsolationAccessor_results struct{ funcDecl *dst.FuncDecl }

// mockConverterer_IsolationAccessor_fnRecorder routes recorded function calls to the mockConverterer mock
type mockConverterer_IsolationAccessor_fnRecorder struct {
	params  mockConverterer_IsolationAccessor_params
	results *mockConverterer_IsolationAccessor_resultMgr
	mock    *mockConverterer
}

// mockConverterer_FuncClosure_params holds the params of the Converterer type
type mockConverterer_FuncClosure_params struct {
	typeName, pkgPath string
	fn                generator.Func
}

// mockConverterer_FuncClosure_resultMgr manages multiple results and the state of the Converterer type
type mockConverterer_FuncClosure_resultMgr struct {
	results  []*mockConverterer_FuncClosure_results
	index    uint32
	anyTimes bool
}

// mockConverterer_FuncClosure_results holds the results of the Converterer type
type mockConverterer_FuncClosure_results struct{ funcDecl *dst.FuncDecl }

// mockConverterer_FuncClosure_fnRecorder routes recorded function calls to the mockConverterer mock
type mockConverterer_FuncClosure_fnRecorder struct {
	params  mockConverterer_FuncClosure_params
	results *mockConverterer_FuncClosure_resultMgr
	mock    *mockConverterer
}

// mockConverterer_MockMethod_params holds the params of the Converterer type
type mockConverterer_MockMethod_params struct {
	typeName string
	fn       generator.Func
}

// mockConverterer_MockMethod_resultMgr manages multiple results and the state of the Converterer type
type mockConverterer_MockMethod_resultMgr struct {
	results  []*mockConverterer_MockMethod_results
	index    uint32
	anyTimes bool
}

// mockConverterer_MockMethod_results holds the results of the Converterer type
type mockConverterer_MockMethod_results struct{ funcDecl *dst.FuncDecl }

// mockConverterer_MockMethod_fnRecorder routes recorded function calls to the mockConverterer mock
type mockConverterer_MockMethod_fnRecorder struct {
	params  mockConverterer_MockMethod_params
	results *mockConverterer_MockMethod_resultMgr
	mock    *mockConverterer
}

// mockConverterer_RecorderMethods_params holds the params of the Converterer type
type mockConverterer_RecorderMethods_params struct {
	typeName string
	fn       generator.Func
}

// mockConverterer_RecorderMethods_resultMgr manages multiple results and the state of the Converterer type
type mockConverterer_RecorderMethods_resultMgr struct {
	results  []*mockConverterer_RecorderMethods_results
	index    uint32
	anyTimes bool
}

// mockConverterer_RecorderMethods_results holds the results of the Converterer type
type mockConverterer_RecorderMethods_results struct{ funcDecls []dst.Decl }

// mockConverterer_RecorderMethods_fnRecorder routes recorded function calls to the mockConverterer mock
type mockConverterer_RecorderMethods_fnRecorder struct {
	params  mockConverterer_RecorderMethods_params
	results *mockConverterer_RecorderMethods_resultMgr
	mock    *mockConverterer
}

// newMockConverterer creates a new mock of the Converterer type
func newMockConverterer(t testing.MoqT) *mockConverterer {
	return &mockConverterer{
		t:                                 t,
		resultsByParams_BaseStruct:        map[mockConverterer_BaseStruct_params]*mockConverterer_BaseStruct_resultMgr{},
		params_BaseStruct:                 make(chan mockConverterer_BaseStruct_params, 100),
		resultsByParams_IsolationStruct:   map[mockConverterer_IsolationStruct_params]*mockConverterer_IsolationStruct_resultMgr{},
		params_IsolationStruct:            make(chan mockConverterer_IsolationStruct_params, 100),
		resultsByParams_MethodStructs:     map[mockConverterer_MethodStructs_params]*mockConverterer_MethodStructs_resultMgr{},
		params_MethodStructs:              make(chan mockConverterer_MethodStructs_params, 100),
		resultsByParams_NewFunc:           map[mockConverterer_NewFunc_params]*mockConverterer_NewFunc_resultMgr{},
		params_NewFunc:                    make(chan mockConverterer_NewFunc_params, 100),
		resultsByParams_IsolationAccessor: map[mockConverterer_IsolationAccessor_params]*mockConverterer_IsolationAccessor_resultMgr{},
		params_IsolationAccessor:          make(chan mockConverterer_IsolationAccessor_params, 100),
		resultsByParams_FuncClosure:       map[mockConverterer_FuncClosure_params]*mockConverterer_FuncClosure_resultMgr{},
		params_FuncClosure:                make(chan mockConverterer_FuncClosure_params, 100),
		resultsByParams_MockMethod:        map[mockConverterer_MockMethod_params]*mockConverterer_MockMethod_resultMgr{},
		params_MockMethod:                 make(chan mockConverterer_MockMethod_params, 100),
		resultsByParams_RecorderMethods:   map[mockConverterer_RecorderMethods_params]*mockConverterer_RecorderMethods_resultMgr{},
		params_RecorderMethods:            make(chan mockConverterer_RecorderMethods_params, 100),
	}
}

// mock returns the mock implementation of the Converterer type
func (m *mockConverterer) mock() *mockConverterer_mock {
	return &mockConverterer_mock{
		mock: m,
	}
}

func (m *mockConverterer_mock) BaseStruct(typeSpec *dst.TypeSpec, funcs []generator.Func) (structDecl *dst.GenDecl) {
	params := mockConverterer_BaseStruct_params{
		typeSpec: typeSpec,
		funcs:    hash.DeepHash(funcs),
	}
	m.mock.params_BaseStruct <- params
	results, ok := m.mock.resultsByParams_BaseStruct[params]
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
		structDecl = result.structDecl
	}
	return structDecl
}

func (m *mockConverterer_mock) IsolationStruct(typeName, suffix string) (structDecl *dst.GenDecl) {
	params := mockConverterer_IsolationStruct_params{
		typeName: typeName,
		suffix:   suffix,
	}
	m.mock.params_IsolationStruct <- params
	results, ok := m.mock.resultsByParams_IsolationStruct[params]
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
		structDecl = result.structDecl
	}
	return structDecl
}

func (m *mockConverterer_mock) MethodStructs(typeSpec *dst.TypeSpec, fn generator.Func) (structDecls []dst.Decl) {
	params := mockConverterer_MethodStructs_params{
		typeSpec: typeSpec,
		fn:       fn,
	}
	m.mock.params_MethodStructs <- params
	results, ok := m.mock.resultsByParams_MethodStructs[params]
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
		structDecls = result.structDecls
	}
	return structDecls
}

func (m *mockConverterer_mock) NewFunc(typeSpec *dst.TypeSpec, funcs []generator.Func) (funcDecl *dst.FuncDecl) {
	params := mockConverterer_NewFunc_params{
		typeSpec: typeSpec,
		funcs:    hash.DeepHash(funcs),
	}
	m.mock.params_NewFunc <- params
	results, ok := m.mock.resultsByParams_NewFunc[params]
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
		funcDecl = result.funcDecl
	}
	return funcDecl
}

func (m *mockConverterer_mock) IsolationAccessor(typeName, suffix, fnName string) (funcDecl *dst.FuncDecl) {
	params := mockConverterer_IsolationAccessor_params{
		typeName: typeName,
		suffix:   suffix,
		fnName:   fnName,
	}
	m.mock.params_IsolationAccessor <- params
	results, ok := m.mock.resultsByParams_IsolationAccessor[params]
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
		funcDecl = result.funcDecl
	}
	return funcDecl
}

func (m *mockConverterer_mock) FuncClosure(typeName, pkgPath string, fn generator.Func) (funcDecl *dst.FuncDecl) {
	params := mockConverterer_FuncClosure_params{
		typeName: typeName,
		pkgPath:  pkgPath,
		fn:       fn,
	}
	m.mock.params_FuncClosure <- params
	results, ok := m.mock.resultsByParams_FuncClosure[params]
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
		funcDecl = result.funcDecl
	}
	return funcDecl
}

func (m *mockConverterer_mock) MockMethod(typeName string, fn generator.Func) (funcDecl *dst.FuncDecl) {
	params := mockConverterer_MockMethod_params{
		typeName: typeName,
		fn:       fn,
	}
	m.mock.params_MockMethod <- params
	results, ok := m.mock.resultsByParams_MockMethod[params]
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
		funcDecl = result.funcDecl
	}
	return funcDecl
}

func (m *mockConverterer_mock) RecorderMethods(typeName string, fn generator.Func) (funcDecls []dst.Decl) {
	params := mockConverterer_RecorderMethods_params{
		typeName: typeName,
		fn:       fn,
	}
	m.mock.params_RecorderMethods <- params
	results, ok := m.mock.resultsByParams_RecorderMethods[params]
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
		funcDecls = result.funcDecls
	}
	return funcDecls
}

// onCall returns the recorder implementation of the Converterer type
func (m *mockConverterer) onCall() *mockConverterer_recorder {
	return &mockConverterer_recorder{
		mock: m,
	}
}

func (m *mockConverterer_recorder) BaseStruct(typeSpec *dst.TypeSpec, funcs []generator.Func) *mockConverterer_BaseStruct_fnRecorder {
	return &mockConverterer_BaseStruct_fnRecorder{
		params: mockConverterer_BaseStruct_params{
			typeSpec: typeSpec,
			funcs:    hash.DeepHash(funcs),
		},
		mock: m.mock,
	}
}

func (r *mockConverterer_BaseStruct_fnRecorder) returnResults(structDecl *dst.GenDecl) *mockConverterer_BaseStruct_fnRecorder {
	if r.results == nil {
		if _, ok := r.mock.resultsByParams_BaseStruct[r.params]; ok {
			r.mock.t.Fatalf("Expectations already recorded for mock with parameters %#v", r.params)
			return nil
		}

		r.results = &mockConverterer_BaseStruct_resultMgr{results: []*mockConverterer_BaseStruct_results{}, index: 0, anyTimes: false}
		r.mock.resultsByParams_BaseStruct[r.params] = r.results
	}
	r.results.results = append(r.results.results, &mockConverterer_BaseStruct_results{
		structDecl: structDecl,
	})
	return r
}

func (r *mockConverterer_BaseStruct_fnRecorder) times(count int) *mockConverterer_BaseStruct_fnRecorder {
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

func (r *mockConverterer_BaseStruct_fnRecorder) anyTimes() {
	if r.results == nil {
		r.mock.t.Fatalf("Return must be called before calling AnyTimes")
		return
	}
	r.results.anyTimes = true
}

func (m *mockConverterer_recorder) IsolationStruct(typeName, suffix string) *mockConverterer_IsolationStruct_fnRecorder {
	return &mockConverterer_IsolationStruct_fnRecorder{
		params: mockConverterer_IsolationStruct_params{
			typeName: typeName,
			suffix:   suffix,
		},
		mock: m.mock,
	}
}

func (r *mockConverterer_IsolationStruct_fnRecorder) returnResults(structDecl *dst.GenDecl) *mockConverterer_IsolationStruct_fnRecorder {
	if r.results == nil {
		if _, ok := r.mock.resultsByParams_IsolationStruct[r.params]; ok {
			r.mock.t.Fatalf("Expectations already recorded for mock with parameters %#v", r.params)
			return nil
		}

		r.results = &mockConverterer_IsolationStruct_resultMgr{results: []*mockConverterer_IsolationStruct_results{}, index: 0, anyTimes: false}
		r.mock.resultsByParams_IsolationStruct[r.params] = r.results
	}
	r.results.results = append(r.results.results, &mockConverterer_IsolationStruct_results{
		structDecl: structDecl,
	})
	return r
}

func (r *mockConverterer_IsolationStruct_fnRecorder) times(count int) *mockConverterer_IsolationStruct_fnRecorder {
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

func (r *mockConverterer_IsolationStruct_fnRecorder) anyTimes() {
	if r.results == nil {
		r.mock.t.Fatalf("Return must be called before calling AnyTimes")
		return
	}
	r.results.anyTimes = true
}

func (m *mockConverterer_recorder) MethodStructs(typeSpec *dst.TypeSpec, fn generator.Func) *mockConverterer_MethodStructs_fnRecorder {
	return &mockConverterer_MethodStructs_fnRecorder{
		params: mockConverterer_MethodStructs_params{
			typeSpec: typeSpec,
			fn:       fn,
		},
		mock: m.mock,
	}
}

func (r *mockConverterer_MethodStructs_fnRecorder) returnResults(structDecls []dst.Decl) *mockConverterer_MethodStructs_fnRecorder {
	if r.results == nil {
		if _, ok := r.mock.resultsByParams_MethodStructs[r.params]; ok {
			r.mock.t.Fatalf("Expectations already recorded for mock with parameters %#v", r.params)
			return nil
		}

		r.results = &mockConverterer_MethodStructs_resultMgr{results: []*mockConverterer_MethodStructs_results{}, index: 0, anyTimes: false}
		r.mock.resultsByParams_MethodStructs[r.params] = r.results
	}
	r.results.results = append(r.results.results, &mockConverterer_MethodStructs_results{
		structDecls: structDecls,
	})
	return r
}

func (r *mockConverterer_MethodStructs_fnRecorder) times(count int) *mockConverterer_MethodStructs_fnRecorder {
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

func (r *mockConverterer_MethodStructs_fnRecorder) anyTimes() {
	if r.results == nil {
		r.mock.t.Fatalf("Return must be called before calling AnyTimes")
		return
	}
	r.results.anyTimes = true
}

func (m *mockConverterer_recorder) NewFunc(typeSpec *dst.TypeSpec, funcs []generator.Func) *mockConverterer_NewFunc_fnRecorder {
	return &mockConverterer_NewFunc_fnRecorder{
		params: mockConverterer_NewFunc_params{
			typeSpec: typeSpec,
			funcs:    hash.DeepHash(funcs),
		},
		mock: m.mock,
	}
}

func (r *mockConverterer_NewFunc_fnRecorder) returnResults(funcDecl *dst.FuncDecl) *mockConverterer_NewFunc_fnRecorder {
	if r.results == nil {
		if _, ok := r.mock.resultsByParams_NewFunc[r.params]; ok {
			r.mock.t.Fatalf("Expectations already recorded for mock with parameters %#v", r.params)
			return nil
		}

		r.results = &mockConverterer_NewFunc_resultMgr{results: []*mockConverterer_NewFunc_results{}, index: 0, anyTimes: false}
		r.mock.resultsByParams_NewFunc[r.params] = r.results
	}
	r.results.results = append(r.results.results, &mockConverterer_NewFunc_results{
		funcDecl: funcDecl,
	})
	return r
}

func (r *mockConverterer_NewFunc_fnRecorder) times(count int) *mockConverterer_NewFunc_fnRecorder {
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

func (r *mockConverterer_NewFunc_fnRecorder) anyTimes() {
	if r.results == nil {
		r.mock.t.Fatalf("Return must be called before calling AnyTimes")
		return
	}
	r.results.anyTimes = true
}

func (m *mockConverterer_recorder) IsolationAccessor(typeName, suffix, fnName string) *mockConverterer_IsolationAccessor_fnRecorder {
	return &mockConverterer_IsolationAccessor_fnRecorder{
		params: mockConverterer_IsolationAccessor_params{
			typeName: typeName,
			suffix:   suffix,
			fnName:   fnName,
		},
		mock: m.mock,
	}
}

func (r *mockConverterer_IsolationAccessor_fnRecorder) returnResults(funcDecl *dst.FuncDecl) *mockConverterer_IsolationAccessor_fnRecorder {
	if r.results == nil {
		if _, ok := r.mock.resultsByParams_IsolationAccessor[r.params]; ok {
			r.mock.t.Fatalf("Expectations already recorded for mock with parameters %#v", r.params)
			return nil
		}

		r.results = &mockConverterer_IsolationAccessor_resultMgr{results: []*mockConverterer_IsolationAccessor_results{}, index: 0, anyTimes: false}
		r.mock.resultsByParams_IsolationAccessor[r.params] = r.results
	}
	r.results.results = append(r.results.results, &mockConverterer_IsolationAccessor_results{
		funcDecl: funcDecl,
	})
	return r
}

func (r *mockConverterer_IsolationAccessor_fnRecorder) times(count int) *mockConverterer_IsolationAccessor_fnRecorder {
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

func (r *mockConverterer_IsolationAccessor_fnRecorder) anyTimes() {
	if r.results == nil {
		r.mock.t.Fatalf("Return must be called before calling AnyTimes")
		return
	}
	r.results.anyTimes = true
}

func (m *mockConverterer_recorder) FuncClosure(typeName, pkgPath string, fn generator.Func) *mockConverterer_FuncClosure_fnRecorder {
	return &mockConverterer_FuncClosure_fnRecorder{
		params: mockConverterer_FuncClosure_params{
			typeName: typeName,
			pkgPath:  pkgPath,
			fn:       fn,
		},
		mock: m.mock,
	}
}

func (r *mockConverterer_FuncClosure_fnRecorder) returnResults(funcDecl *dst.FuncDecl) *mockConverterer_FuncClosure_fnRecorder {
	if r.results == nil {
		if _, ok := r.mock.resultsByParams_FuncClosure[r.params]; ok {
			r.mock.t.Fatalf("Expectations already recorded for mock with parameters %#v", r.params)
			return nil
		}

		r.results = &mockConverterer_FuncClosure_resultMgr{results: []*mockConverterer_FuncClosure_results{}, index: 0, anyTimes: false}
		r.mock.resultsByParams_FuncClosure[r.params] = r.results
	}
	r.results.results = append(r.results.results, &mockConverterer_FuncClosure_results{
		funcDecl: funcDecl,
	})
	return r
}

func (r *mockConverterer_FuncClosure_fnRecorder) times(count int) *mockConverterer_FuncClosure_fnRecorder {
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

func (r *mockConverterer_FuncClosure_fnRecorder) anyTimes() {
	if r.results == nil {
		r.mock.t.Fatalf("Return must be called before calling AnyTimes")
		return
	}
	r.results.anyTimes = true
}

func (m *mockConverterer_recorder) MockMethod(typeName string, fn generator.Func) *mockConverterer_MockMethod_fnRecorder {
	return &mockConverterer_MockMethod_fnRecorder{
		params: mockConverterer_MockMethod_params{
			typeName: typeName,
			fn:       fn,
		},
		mock: m.mock,
	}
}

func (r *mockConverterer_MockMethod_fnRecorder) returnResults(funcDecl *dst.FuncDecl) *mockConverterer_MockMethod_fnRecorder {
	if r.results == nil {
		if _, ok := r.mock.resultsByParams_MockMethod[r.params]; ok {
			r.mock.t.Fatalf("Expectations already recorded for mock with parameters %#v", r.params)
			return nil
		}

		r.results = &mockConverterer_MockMethod_resultMgr{results: []*mockConverterer_MockMethod_results{}, index: 0, anyTimes: false}
		r.mock.resultsByParams_MockMethod[r.params] = r.results
	}
	r.results.results = append(r.results.results, &mockConverterer_MockMethod_results{
		funcDecl: funcDecl,
	})
	return r
}

func (r *mockConverterer_MockMethod_fnRecorder) times(count int) *mockConverterer_MockMethod_fnRecorder {
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

func (r *mockConverterer_MockMethod_fnRecorder) anyTimes() {
	if r.results == nil {
		r.mock.t.Fatalf("Return must be called before calling AnyTimes")
		return
	}
	r.results.anyTimes = true
}

func (m *mockConverterer_recorder) RecorderMethods(typeName string, fn generator.Func) *mockConverterer_RecorderMethods_fnRecorder {
	return &mockConverterer_RecorderMethods_fnRecorder{
		params: mockConverterer_RecorderMethods_params{
			typeName: typeName,
			fn:       fn,
		},
		mock: m.mock,
	}
}

func (r *mockConverterer_RecorderMethods_fnRecorder) returnResults(funcDecls []dst.Decl) *mockConverterer_RecorderMethods_fnRecorder {
	if r.results == nil {
		if _, ok := r.mock.resultsByParams_RecorderMethods[r.params]; ok {
			r.mock.t.Fatalf("Expectations already recorded for mock with parameters %#v", r.params)
			return nil
		}

		r.results = &mockConverterer_RecorderMethods_resultMgr{results: []*mockConverterer_RecorderMethods_results{}, index: 0, anyTimes: false}
		r.mock.resultsByParams_RecorderMethods[r.params] = r.results
	}
	r.results.results = append(r.results.results, &mockConverterer_RecorderMethods_results{
		funcDecls: funcDecls,
	})
	return r
}

func (r *mockConverterer_RecorderMethods_fnRecorder) times(count int) *mockConverterer_RecorderMethods_fnRecorder {
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

func (r *mockConverterer_RecorderMethods_fnRecorder) anyTimes() {
	if r.results == nil {
		r.mock.t.Fatalf("Return must be called before calling AnyTimes")
		return
	}
	r.results.anyTimes = true
}
