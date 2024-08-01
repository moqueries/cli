// Code generated by Moqueries - https://moqueries.org - DO NOT EDIT.

package internal_test

import (
	"fmt"

	"github.com/dave/dst"
	"moqueries.org/cli/ast"
	"moqueries.org/cli/pkg/internal"
	"moqueries.org/runtime/hash"
	"moqueries.org/runtime/impl"
	"moqueries.org/runtime/moq"
)

// The following type assertion assures that internal.TypeCache is mocked
// completely
var _ internal.TypeCache = (*moqTypeCache_mock)(nil)

// moqTypeCache holds the state of a moq of the TypeCache type
type moqTypeCache struct {
	moq *moqTypeCache_mock

	moq_LoadPackage *impl.Moq[
		*moqTypeCache_LoadPackage_adaptor,
		moqTypeCache_LoadPackage_params,
		moqTypeCache_LoadPackage_paramsKey,
		moqTypeCache_LoadPackage_results,
	]
	moq_MockableTypes *impl.Moq[
		*moqTypeCache_MockableTypes_adaptor,
		moqTypeCache_MockableTypes_params,
		moqTypeCache_MockableTypes_paramsKey,
		moqTypeCache_MockableTypes_results,
	]
	moq_Type *impl.Moq[
		*moqTypeCache_Type_adaptor,
		moqTypeCache_Type_params,
		moqTypeCache_Type_paramsKey,
		moqTypeCache_Type_results,
	]
	moq_IsComparable *impl.Moq[
		*moqTypeCache_IsComparable_adaptor,
		moqTypeCache_IsComparable_params,
		moqTypeCache_IsComparable_paramsKey,
		moqTypeCache_IsComparable_results,
	]
	moq_IsDefaultComparable *impl.Moq[
		*moqTypeCache_IsDefaultComparable_adaptor,
		moqTypeCache_IsDefaultComparable_params,
		moqTypeCache_IsDefaultComparable_paramsKey,
		moqTypeCache_IsDefaultComparable_results,
	]
	moq_FindPackage *impl.Moq[
		*moqTypeCache_FindPackage_adaptor,
		moqTypeCache_FindPackage_params,
		moqTypeCache_FindPackage_paramsKey,
		moqTypeCache_FindPackage_results,
	]

	runtime moqTypeCache_runtime
}

// moqTypeCache_mock isolates the mock interface of the TypeCache type
type moqTypeCache_mock struct {
	moq *moqTypeCache
}

// moqTypeCache_recorder isolates the recorder interface of the TypeCache type
type moqTypeCache_recorder struct {
	moq *moqTypeCache
}

// moqTypeCache_runtime holds runtime configuration for the TypeCache type
type moqTypeCache_runtime struct {
	parameterIndexing struct {
		LoadPackage         moqTypeCache_LoadPackage_paramIndexing
		MockableTypes       moqTypeCache_MockableTypes_paramIndexing
		Type                moqTypeCache_Type_paramIndexing
		IsComparable        moqTypeCache_IsComparable_paramIndexing
		IsDefaultComparable moqTypeCache_IsDefaultComparable_paramIndexing
		FindPackage         moqTypeCache_FindPackage_paramIndexing
	}
}

// moqTypeCache_LoadPackage_adaptor adapts moqTypeCache as needed by the
// runtime
type moqTypeCache_LoadPackage_adaptor struct {
	moq *moqTypeCache
}

// moqTypeCache_LoadPackage_params holds the params of the TypeCache type
type moqTypeCache_LoadPackage_params struct{ pkgPattern string }

// moqTypeCache_LoadPackage_paramsKey holds the map key params of the TypeCache
// type
type moqTypeCache_LoadPackage_paramsKey struct {
	params struct{ pkgPattern string }
	hashes struct{ pkgPattern hash.Hash }
}

// moqTypeCache_LoadPackage_results holds the results of the TypeCache type
type moqTypeCache_LoadPackage_results struct {
	result1 error
}

// moqTypeCache_LoadPackage_paramIndexing holds the parameter indexing runtime
// configuration for the TypeCache type
type moqTypeCache_LoadPackage_paramIndexing struct {
	pkgPattern moq.ParamIndexing
}

// moqTypeCache_LoadPackage_doFn defines the type of function needed when
// calling andDo for the TypeCache type
type moqTypeCache_LoadPackage_doFn func(pkgPattern string)

// moqTypeCache_LoadPackage_doReturnFn defines the type of function needed when
// calling doReturnResults for the TypeCache type
type moqTypeCache_LoadPackage_doReturnFn func(pkgPattern string) error

// moqTypeCache_LoadPackage_recorder routes recorded function calls to the
// moqTypeCache moq
type moqTypeCache_LoadPackage_recorder struct {
	recorder *impl.Recorder[
		*moqTypeCache_LoadPackage_adaptor,
		moqTypeCache_LoadPackage_params,
		moqTypeCache_LoadPackage_paramsKey,
		moqTypeCache_LoadPackage_results,
	]
}

// moqTypeCache_LoadPackage_anyParams isolates the any params functions of the
// TypeCache type
type moqTypeCache_LoadPackage_anyParams struct {
	recorder *moqTypeCache_LoadPackage_recorder
}

// moqTypeCache_MockableTypes_adaptor adapts moqTypeCache as needed by the
// runtime
type moqTypeCache_MockableTypes_adaptor struct {
	moq *moqTypeCache
}

// moqTypeCache_MockableTypes_params holds the params of the TypeCache type
type moqTypeCache_MockableTypes_params struct{ onlyExported bool }

// moqTypeCache_MockableTypes_paramsKey holds the map key params of the
// TypeCache type
type moqTypeCache_MockableTypes_paramsKey struct {
	params struct{ onlyExported bool }
	hashes struct{ onlyExported hash.Hash }
}

// moqTypeCache_MockableTypes_results holds the results of the TypeCache type
type moqTypeCache_MockableTypes_results struct {
	result1 []dst.Ident
}

// moqTypeCache_MockableTypes_paramIndexing holds the parameter indexing
// runtime configuration for the TypeCache type
type moqTypeCache_MockableTypes_paramIndexing struct {
	onlyExported moq.ParamIndexing
}

// moqTypeCache_MockableTypes_doFn defines the type of function needed when
// calling andDo for the TypeCache type
type moqTypeCache_MockableTypes_doFn func(onlyExported bool)

// moqTypeCache_MockableTypes_doReturnFn defines the type of function needed
// when calling doReturnResults for the TypeCache type
type moqTypeCache_MockableTypes_doReturnFn func(onlyExported bool) []dst.Ident

// moqTypeCache_MockableTypes_recorder routes recorded function calls to the
// moqTypeCache moq
type moqTypeCache_MockableTypes_recorder struct {
	recorder *impl.Recorder[
		*moqTypeCache_MockableTypes_adaptor,
		moqTypeCache_MockableTypes_params,
		moqTypeCache_MockableTypes_paramsKey,
		moqTypeCache_MockableTypes_results,
	]
}

// moqTypeCache_MockableTypes_anyParams isolates the any params functions of
// the TypeCache type
type moqTypeCache_MockableTypes_anyParams struct {
	recorder *moqTypeCache_MockableTypes_recorder
}

// moqTypeCache_Type_adaptor adapts moqTypeCache as needed by the runtime
type moqTypeCache_Type_adaptor struct {
	moq *moqTypeCache
}

// moqTypeCache_Type_params holds the params of the TypeCache type
type moqTypeCache_Type_params struct {
	id         dst.Ident
	contextPkg string
	testImport bool
}

// moqTypeCache_Type_paramsKey holds the map key params of the TypeCache type
type moqTypeCache_Type_paramsKey struct {
	params struct {
		contextPkg string
		testImport bool
	}
	hashes struct {
		id         hash.Hash
		contextPkg hash.Hash
		testImport hash.Hash
	}
}

// moqTypeCache_Type_results holds the results of the TypeCache type
type moqTypeCache_Type_results struct {
	result1 ast.TypeInfo
	result2 error
}

// moqTypeCache_Type_paramIndexing holds the parameter indexing runtime
// configuration for the TypeCache type
type moqTypeCache_Type_paramIndexing struct {
	id         moq.ParamIndexing
	contextPkg moq.ParamIndexing
	testImport moq.ParamIndexing
}

// moqTypeCache_Type_doFn defines the type of function needed when calling
// andDo for the TypeCache type
type moqTypeCache_Type_doFn func(id dst.Ident, contextPkg string, testImport bool)

// moqTypeCache_Type_doReturnFn defines the type of function needed when
// calling doReturnResults for the TypeCache type
type moqTypeCache_Type_doReturnFn func(id dst.Ident, contextPkg string, testImport bool) (ast.TypeInfo, error)

// moqTypeCache_Type_recorder routes recorded function calls to the
// moqTypeCache moq
type moqTypeCache_Type_recorder struct {
	recorder *impl.Recorder[
		*moqTypeCache_Type_adaptor,
		moqTypeCache_Type_params,
		moqTypeCache_Type_paramsKey,
		moqTypeCache_Type_results,
	]
}

// moqTypeCache_Type_anyParams isolates the any params functions of the
// TypeCache type
type moqTypeCache_Type_anyParams struct {
	recorder *moqTypeCache_Type_recorder
}

// moqTypeCache_IsComparable_adaptor adapts moqTypeCache as needed by the
// runtime
type moqTypeCache_IsComparable_adaptor struct {
	moq *moqTypeCache
}

// moqTypeCache_IsComparable_params holds the params of the TypeCache type
type moqTypeCache_IsComparable_params struct {
	expr       dst.Expr
	parentType ast.TypeInfo
}

// moqTypeCache_IsComparable_paramsKey holds the map key params of the
// TypeCache type
type moqTypeCache_IsComparable_paramsKey struct {
	params struct {
		expr       dst.Expr
		parentType ast.TypeInfo
	}
	hashes struct {
		expr       hash.Hash
		parentType hash.Hash
	}
}

// moqTypeCache_IsComparable_results holds the results of the TypeCache type
type moqTypeCache_IsComparable_results struct {
	result1 bool
	result2 error
}

// moqTypeCache_IsComparable_paramIndexing holds the parameter indexing runtime
// configuration for the TypeCache type
type moqTypeCache_IsComparable_paramIndexing struct {
	expr       moq.ParamIndexing
	parentType moq.ParamIndexing
}

// moqTypeCache_IsComparable_doFn defines the type of function needed when
// calling andDo for the TypeCache type
type moqTypeCache_IsComparable_doFn func(expr dst.Expr, parentType ast.TypeInfo)

// moqTypeCache_IsComparable_doReturnFn defines the type of function needed
// when calling doReturnResults for the TypeCache type
type moqTypeCache_IsComparable_doReturnFn func(expr dst.Expr, parentType ast.TypeInfo) (bool, error)

// moqTypeCache_IsComparable_recorder routes recorded function calls to the
// moqTypeCache moq
type moqTypeCache_IsComparable_recorder struct {
	recorder *impl.Recorder[
		*moqTypeCache_IsComparable_adaptor,
		moqTypeCache_IsComparable_params,
		moqTypeCache_IsComparable_paramsKey,
		moqTypeCache_IsComparable_results,
	]
}

// moqTypeCache_IsComparable_anyParams isolates the any params functions of the
// TypeCache type
type moqTypeCache_IsComparable_anyParams struct {
	recorder *moqTypeCache_IsComparable_recorder
}

// moqTypeCache_IsDefaultComparable_adaptor adapts moqTypeCache as needed by
// the runtime
type moqTypeCache_IsDefaultComparable_adaptor struct {
	moq *moqTypeCache
}

// moqTypeCache_IsDefaultComparable_params holds the params of the TypeCache
// type
type moqTypeCache_IsDefaultComparable_params struct {
	expr       dst.Expr
	parentType ast.TypeInfo
}

// moqTypeCache_IsDefaultComparable_paramsKey holds the map key params of the
// TypeCache type
type moqTypeCache_IsDefaultComparable_paramsKey struct {
	params struct {
		expr       dst.Expr
		parentType ast.TypeInfo
	}
	hashes struct {
		expr       hash.Hash
		parentType hash.Hash
	}
}

// moqTypeCache_IsDefaultComparable_results holds the results of the TypeCache
// type
type moqTypeCache_IsDefaultComparable_results struct {
	result1 bool
	result2 error
}

// moqTypeCache_IsDefaultComparable_paramIndexing holds the parameter indexing
// runtime configuration for the TypeCache type
type moqTypeCache_IsDefaultComparable_paramIndexing struct {
	expr       moq.ParamIndexing
	parentType moq.ParamIndexing
}

// moqTypeCache_IsDefaultComparable_doFn defines the type of function needed
// when calling andDo for the TypeCache type
type moqTypeCache_IsDefaultComparable_doFn func(expr dst.Expr, parentType ast.TypeInfo)

// moqTypeCache_IsDefaultComparable_doReturnFn defines the type of function
// needed when calling doReturnResults for the TypeCache type
type moqTypeCache_IsDefaultComparable_doReturnFn func(expr dst.Expr, parentType ast.TypeInfo) (bool, error)

// moqTypeCache_IsDefaultComparable_recorder routes recorded function calls to
// the moqTypeCache moq
type moqTypeCache_IsDefaultComparable_recorder struct {
	recorder *impl.Recorder[
		*moqTypeCache_IsDefaultComparable_adaptor,
		moqTypeCache_IsDefaultComparable_params,
		moqTypeCache_IsDefaultComparable_paramsKey,
		moqTypeCache_IsDefaultComparable_results,
	]
}

// moqTypeCache_IsDefaultComparable_anyParams isolates the any params functions
// of the TypeCache type
type moqTypeCache_IsDefaultComparable_anyParams struct {
	recorder *moqTypeCache_IsDefaultComparable_recorder
}

// moqTypeCache_FindPackage_adaptor adapts moqTypeCache as needed by the
// runtime
type moqTypeCache_FindPackage_adaptor struct {
	moq *moqTypeCache
}

// moqTypeCache_FindPackage_params holds the params of the TypeCache type
type moqTypeCache_FindPackage_params struct{ dir string }

// moqTypeCache_FindPackage_paramsKey holds the map key params of the TypeCache
// type
type moqTypeCache_FindPackage_paramsKey struct {
	params struct{ dir string }
	hashes struct{ dir hash.Hash }
}

// moqTypeCache_FindPackage_results holds the results of the TypeCache type
type moqTypeCache_FindPackage_results struct {
	result1 string
	result2 error
}

// moqTypeCache_FindPackage_paramIndexing holds the parameter indexing runtime
// configuration for the TypeCache type
type moqTypeCache_FindPackage_paramIndexing struct {
	dir moq.ParamIndexing
}

// moqTypeCache_FindPackage_doFn defines the type of function needed when
// calling andDo for the TypeCache type
type moqTypeCache_FindPackage_doFn func(dir string)

// moqTypeCache_FindPackage_doReturnFn defines the type of function needed when
// calling doReturnResults for the TypeCache type
type moqTypeCache_FindPackage_doReturnFn func(dir string) (string, error)

// moqTypeCache_FindPackage_recorder routes recorded function calls to the
// moqTypeCache moq
type moqTypeCache_FindPackage_recorder struct {
	recorder *impl.Recorder[
		*moqTypeCache_FindPackage_adaptor,
		moqTypeCache_FindPackage_params,
		moqTypeCache_FindPackage_paramsKey,
		moqTypeCache_FindPackage_results,
	]
}

// moqTypeCache_FindPackage_anyParams isolates the any params functions of the
// TypeCache type
type moqTypeCache_FindPackage_anyParams struct {
	recorder *moqTypeCache_FindPackage_recorder
}

// newMoqTypeCache creates a new moq of the TypeCache type
func newMoqTypeCache(scene *moq.Scene, config *moq.Config) *moqTypeCache {
	adaptor1 := &moqTypeCache_LoadPackage_adaptor{}
	adaptor2 := &moqTypeCache_MockableTypes_adaptor{}
	adaptor3 := &moqTypeCache_Type_adaptor{}
	adaptor4 := &moqTypeCache_IsComparable_adaptor{}
	adaptor5 := &moqTypeCache_IsDefaultComparable_adaptor{}
	adaptor6 := &moqTypeCache_FindPackage_adaptor{}
	m := &moqTypeCache{
		moq: &moqTypeCache_mock{},

		moq_LoadPackage: impl.NewMoq[
			*moqTypeCache_LoadPackage_adaptor,
			moqTypeCache_LoadPackage_params,
			moqTypeCache_LoadPackage_paramsKey,
			moqTypeCache_LoadPackage_results,
		](scene, adaptor1, config),
		moq_MockableTypes: impl.NewMoq[
			*moqTypeCache_MockableTypes_adaptor,
			moqTypeCache_MockableTypes_params,
			moqTypeCache_MockableTypes_paramsKey,
			moqTypeCache_MockableTypes_results,
		](scene, adaptor2, config),
		moq_Type: impl.NewMoq[
			*moqTypeCache_Type_adaptor,
			moqTypeCache_Type_params,
			moqTypeCache_Type_paramsKey,
			moqTypeCache_Type_results,
		](scene, adaptor3, config),
		moq_IsComparable: impl.NewMoq[
			*moqTypeCache_IsComparable_adaptor,
			moqTypeCache_IsComparable_params,
			moqTypeCache_IsComparable_paramsKey,
			moqTypeCache_IsComparable_results,
		](scene, adaptor4, config),
		moq_IsDefaultComparable: impl.NewMoq[
			*moqTypeCache_IsDefaultComparable_adaptor,
			moqTypeCache_IsDefaultComparable_params,
			moqTypeCache_IsDefaultComparable_paramsKey,
			moqTypeCache_IsDefaultComparable_results,
		](scene, adaptor5, config),
		moq_FindPackage: impl.NewMoq[
			*moqTypeCache_FindPackage_adaptor,
			moqTypeCache_FindPackage_params,
			moqTypeCache_FindPackage_paramsKey,
			moqTypeCache_FindPackage_results,
		](scene, adaptor6, config),

		runtime: moqTypeCache_runtime{parameterIndexing: struct {
			LoadPackage         moqTypeCache_LoadPackage_paramIndexing
			MockableTypes       moqTypeCache_MockableTypes_paramIndexing
			Type                moqTypeCache_Type_paramIndexing
			IsComparable        moqTypeCache_IsComparable_paramIndexing
			IsDefaultComparable moqTypeCache_IsDefaultComparable_paramIndexing
			FindPackage         moqTypeCache_FindPackage_paramIndexing
		}{
			LoadPackage: moqTypeCache_LoadPackage_paramIndexing{
				pkgPattern: moq.ParamIndexByValue,
			},
			MockableTypes: moqTypeCache_MockableTypes_paramIndexing{
				onlyExported: moq.ParamIndexByValue,
			},
			Type: moqTypeCache_Type_paramIndexing{
				id:         moq.ParamIndexByHash,
				contextPkg: moq.ParamIndexByValue,
				testImport: moq.ParamIndexByValue,
			},
			IsComparable: moqTypeCache_IsComparable_paramIndexing{
				expr:       moq.ParamIndexByHash,
				parentType: moq.ParamIndexByHash,
			},
			IsDefaultComparable: moqTypeCache_IsDefaultComparable_paramIndexing{
				expr:       moq.ParamIndexByHash,
				parentType: moq.ParamIndexByHash,
			},
			FindPackage: moqTypeCache_FindPackage_paramIndexing{
				dir: moq.ParamIndexByValue,
			},
		}},
	}
	m.moq.moq = m

	adaptor1.moq = m
	adaptor2.moq = m
	adaptor3.moq = m
	adaptor4.moq = m
	adaptor5.moq = m
	adaptor6.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the mock implementation of the TypeCache type
func (m *moqTypeCache) mock() *moqTypeCache_mock { return m.moq }

func (m *moqTypeCache_mock) LoadPackage(pkgPattern string) error {
	m.moq.moq_LoadPackage.Scene.T.Helper()
	params := moqTypeCache_LoadPackage_params{
		pkgPattern: pkgPattern,
	}

	var result1 error
	if result := m.moq.moq_LoadPackage.Function(params); result != nil {
		result1 = result.result1
	}
	return result1
}

func (m *moqTypeCache_mock) MockableTypes(onlyExported bool) []dst.Ident {
	m.moq.moq_MockableTypes.Scene.T.Helper()
	params := moqTypeCache_MockableTypes_params{
		onlyExported: onlyExported,
	}

	var result1 []dst.Ident
	if result := m.moq.moq_MockableTypes.Function(params); result != nil {
		result1 = result.result1
	}
	return result1
}

func (m *moqTypeCache_mock) Type(id dst.Ident, contextPkg string, testImport bool) (ast.TypeInfo, error) {
	m.moq.moq_Type.Scene.T.Helper()
	params := moqTypeCache_Type_params{
		id:         id,
		contextPkg: contextPkg,
		testImport: testImport,
	}

	var result1 ast.TypeInfo
	var result2 error
	if result := m.moq.moq_Type.Function(params); result != nil {
		result1 = result.result1
		result2 = result.result2
	}
	return result1, result2
}

func (m *moqTypeCache_mock) IsComparable(expr dst.Expr, parentType ast.TypeInfo) (bool, error) {
	m.moq.moq_IsComparable.Scene.T.Helper()
	params := moqTypeCache_IsComparable_params{
		expr:       expr,
		parentType: parentType,
	}

	var result1 bool
	var result2 error
	if result := m.moq.moq_IsComparable.Function(params); result != nil {
		result1 = result.result1
		result2 = result.result2
	}
	return result1, result2
}

func (m *moqTypeCache_mock) IsDefaultComparable(expr dst.Expr, parentType ast.TypeInfo) (bool, error) {
	m.moq.moq_IsDefaultComparable.Scene.T.Helper()
	params := moqTypeCache_IsDefaultComparable_params{
		expr:       expr,
		parentType: parentType,
	}

	var result1 bool
	var result2 error
	if result := m.moq.moq_IsDefaultComparable.Function(params); result != nil {
		result1 = result.result1
		result2 = result.result2
	}
	return result1, result2
}

func (m *moqTypeCache_mock) FindPackage(dir string) (string, error) {
	m.moq.moq_FindPackage.Scene.T.Helper()
	params := moqTypeCache_FindPackage_params{
		dir: dir,
	}

	var result1 string
	var result2 error
	if result := m.moq.moq_FindPackage.Function(params); result != nil {
		result1 = result.result1
		result2 = result.result2
	}
	return result1, result2
}

// onCall returns the recorder implementation of the TypeCache type
func (m *moqTypeCache) onCall() *moqTypeCache_recorder {
	return &moqTypeCache_recorder{
		moq: m,
	}
}

func (m *moqTypeCache_recorder) LoadPackage(pkgPattern string) *moqTypeCache_LoadPackage_recorder {
	return &moqTypeCache_LoadPackage_recorder{
		recorder: m.moq.moq_LoadPackage.OnCall(moqTypeCache_LoadPackage_params{
			pkgPattern: pkgPattern,
		}),
	}
}

func (r *moqTypeCache_LoadPackage_recorder) any() *moqTypeCache_LoadPackage_anyParams {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.IsAnyPermitted(false) {
		return nil
	}
	return &moqTypeCache_LoadPackage_anyParams{recorder: r}
}

func (a *moqTypeCache_LoadPackage_anyParams) pkgPattern() *moqTypeCache_LoadPackage_recorder {
	a.recorder.recorder.AnyParam(1)
	return a.recorder
}

func (r *moqTypeCache_LoadPackage_recorder) seq() *moqTypeCache_LoadPackage_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(true, "seq", false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_LoadPackage_recorder) noSeq() *moqTypeCache_LoadPackage_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(false, "noSeq", false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_LoadPackage_recorder) returnResults(result1 error) *moqTypeCache_LoadPackage_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.ReturnResults(moqTypeCache_LoadPackage_results{
		result1: result1,
	})
	return r
}

func (r *moqTypeCache_LoadPackage_recorder) andDo(fn moqTypeCache_LoadPackage_doFn) *moqTypeCache_LoadPackage_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.AndDo(func(params moqTypeCache_LoadPackage_params) {
		fn(params.pkgPattern)
	}, false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_LoadPackage_recorder) doReturnResults(fn moqTypeCache_LoadPackage_doReturnFn) *moqTypeCache_LoadPackage_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.DoReturnResults(func(params moqTypeCache_LoadPackage_params) *moqTypeCache_LoadPackage_results {
		result1 := fn(params.pkgPattern)
		return &moqTypeCache_LoadPackage_results{
			result1: result1,
		}
	})
	return r
}

func (r *moqTypeCache_LoadPackage_recorder) repeat(repeaters ...moq.Repeater) *moqTypeCache_LoadPackage_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Repeat(repeaters, false) {
		return nil
	}
	return r
}

func (*moqTypeCache_LoadPackage_adaptor) PrettyParams(params moqTypeCache_LoadPackage_params) string {
	return fmt.Sprintf("LoadPackage(%#v)", params.pkgPattern)
}

func (a *moqTypeCache_LoadPackage_adaptor) ParamsKey(params moqTypeCache_LoadPackage_params, anyParams uint64) moqTypeCache_LoadPackage_paramsKey {
	a.moq.moq_LoadPackage.Scene.T.Helper()
	pkgPatternUsed, pkgPatternUsedHash := impl.ParamKey(
		params.pkgPattern, 1, a.moq.runtime.parameterIndexing.LoadPackage.pkgPattern, anyParams)
	return moqTypeCache_LoadPackage_paramsKey{
		params: struct{ pkgPattern string }{
			pkgPattern: pkgPatternUsed,
		},
		hashes: struct{ pkgPattern hash.Hash }{
			pkgPattern: pkgPatternUsedHash,
		},
	}
}

func (m *moqTypeCache_recorder) MockableTypes(onlyExported bool) *moqTypeCache_MockableTypes_recorder {
	return &moqTypeCache_MockableTypes_recorder{
		recorder: m.moq.moq_MockableTypes.OnCall(moqTypeCache_MockableTypes_params{
			onlyExported: onlyExported,
		}),
	}
}

func (r *moqTypeCache_MockableTypes_recorder) any() *moqTypeCache_MockableTypes_anyParams {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.IsAnyPermitted(false) {
		return nil
	}
	return &moqTypeCache_MockableTypes_anyParams{recorder: r}
}

func (a *moqTypeCache_MockableTypes_anyParams) onlyExported() *moqTypeCache_MockableTypes_recorder {
	a.recorder.recorder.AnyParam(1)
	return a.recorder
}

func (r *moqTypeCache_MockableTypes_recorder) seq() *moqTypeCache_MockableTypes_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(true, "seq", false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_MockableTypes_recorder) noSeq() *moqTypeCache_MockableTypes_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(false, "noSeq", false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_MockableTypes_recorder) returnResults(result1 []dst.Ident) *moqTypeCache_MockableTypes_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.ReturnResults(moqTypeCache_MockableTypes_results{
		result1: result1,
	})
	return r
}

func (r *moqTypeCache_MockableTypes_recorder) andDo(fn moqTypeCache_MockableTypes_doFn) *moqTypeCache_MockableTypes_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.AndDo(func(params moqTypeCache_MockableTypes_params) {
		fn(params.onlyExported)
	}, false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_MockableTypes_recorder) doReturnResults(fn moqTypeCache_MockableTypes_doReturnFn) *moqTypeCache_MockableTypes_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.DoReturnResults(func(params moqTypeCache_MockableTypes_params) *moqTypeCache_MockableTypes_results {
		result1 := fn(params.onlyExported)
		return &moqTypeCache_MockableTypes_results{
			result1: result1,
		}
	})
	return r
}

func (r *moqTypeCache_MockableTypes_recorder) repeat(repeaters ...moq.Repeater) *moqTypeCache_MockableTypes_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Repeat(repeaters, false) {
		return nil
	}
	return r
}

func (*moqTypeCache_MockableTypes_adaptor) PrettyParams(params moqTypeCache_MockableTypes_params) string {
	return fmt.Sprintf("MockableTypes(%#v)", params.onlyExported)
}

func (a *moqTypeCache_MockableTypes_adaptor) ParamsKey(params moqTypeCache_MockableTypes_params, anyParams uint64) moqTypeCache_MockableTypes_paramsKey {
	a.moq.moq_MockableTypes.Scene.T.Helper()
	onlyExportedUsed, onlyExportedUsedHash := impl.ParamKey(
		params.onlyExported, 1, a.moq.runtime.parameterIndexing.MockableTypes.onlyExported, anyParams)
	return moqTypeCache_MockableTypes_paramsKey{
		params: struct{ onlyExported bool }{
			onlyExported: onlyExportedUsed,
		},
		hashes: struct{ onlyExported hash.Hash }{
			onlyExported: onlyExportedUsedHash,
		},
	}
}

func (m *moqTypeCache_recorder) Type(id dst.Ident, contextPkg string, testImport bool) *moqTypeCache_Type_recorder {
	return &moqTypeCache_Type_recorder{
		recorder: m.moq.moq_Type.OnCall(moqTypeCache_Type_params{
			id:         id,
			contextPkg: contextPkg,
			testImport: testImport,
		}),
	}
}

func (r *moqTypeCache_Type_recorder) any() *moqTypeCache_Type_anyParams {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.IsAnyPermitted(false) {
		return nil
	}
	return &moqTypeCache_Type_anyParams{recorder: r}
}

func (a *moqTypeCache_Type_anyParams) id() *moqTypeCache_Type_recorder {
	a.recorder.recorder.AnyParam(1)
	return a.recorder
}

func (a *moqTypeCache_Type_anyParams) contextPkg() *moqTypeCache_Type_recorder {
	a.recorder.recorder.AnyParam(2)
	return a.recorder
}

func (a *moqTypeCache_Type_anyParams) testImport() *moqTypeCache_Type_recorder {
	a.recorder.recorder.AnyParam(3)
	return a.recorder
}

func (r *moqTypeCache_Type_recorder) seq() *moqTypeCache_Type_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(true, "seq", false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_Type_recorder) noSeq() *moqTypeCache_Type_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(false, "noSeq", false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_Type_recorder) returnResults(result1 ast.TypeInfo, result2 error) *moqTypeCache_Type_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.ReturnResults(moqTypeCache_Type_results{
		result1: result1,
		result2: result2,
	})
	return r
}

func (r *moqTypeCache_Type_recorder) andDo(fn moqTypeCache_Type_doFn) *moqTypeCache_Type_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.AndDo(func(params moqTypeCache_Type_params) {
		fn(params.id, params.contextPkg, params.testImport)
	}, false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_Type_recorder) doReturnResults(fn moqTypeCache_Type_doReturnFn) *moqTypeCache_Type_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.DoReturnResults(func(params moqTypeCache_Type_params) *moqTypeCache_Type_results {
		result1, result2 := fn(params.id, params.contextPkg, params.testImport)
		return &moqTypeCache_Type_results{
			result1: result1,
			result2: result2,
		}
	})
	return r
}

func (r *moqTypeCache_Type_recorder) repeat(repeaters ...moq.Repeater) *moqTypeCache_Type_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Repeat(repeaters, false) {
		return nil
	}
	return r
}

func (*moqTypeCache_Type_adaptor) PrettyParams(params moqTypeCache_Type_params) string {
	return fmt.Sprintf("Type(%#v, %#v, %#v)", params.id, params.contextPkg, params.testImport)
}

func (a *moqTypeCache_Type_adaptor) ParamsKey(params moqTypeCache_Type_params, anyParams uint64) moqTypeCache_Type_paramsKey {
	a.moq.moq_Type.Scene.T.Helper()
	idUsedHash := impl.HashOnlyParamKey(a.moq.moq_Type.Scene.T,
		params.id, "id", 1, a.moq.runtime.parameterIndexing.Type.id, anyParams)
	contextPkgUsed, contextPkgUsedHash := impl.ParamKey(
		params.contextPkg, 2, a.moq.runtime.parameterIndexing.Type.contextPkg, anyParams)
	testImportUsed, testImportUsedHash := impl.ParamKey(
		params.testImport, 3, a.moq.runtime.parameterIndexing.Type.testImport, anyParams)
	return moqTypeCache_Type_paramsKey{
		params: struct {
			contextPkg string
			testImport bool
		}{
			contextPkg: contextPkgUsed,
			testImport: testImportUsed,
		},
		hashes: struct {
			id         hash.Hash
			contextPkg hash.Hash
			testImport hash.Hash
		}{
			id:         idUsedHash,
			contextPkg: contextPkgUsedHash,
			testImport: testImportUsedHash,
		},
	}
}

func (m *moqTypeCache_recorder) IsComparable(expr dst.Expr, parentType ast.TypeInfo) *moqTypeCache_IsComparable_recorder {
	return &moqTypeCache_IsComparable_recorder{
		recorder: m.moq.moq_IsComparable.OnCall(moqTypeCache_IsComparable_params{
			expr:       expr,
			parentType: parentType,
		}),
	}
}

func (r *moqTypeCache_IsComparable_recorder) any() *moqTypeCache_IsComparable_anyParams {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.IsAnyPermitted(false) {
		return nil
	}
	return &moqTypeCache_IsComparable_anyParams{recorder: r}
}

func (a *moqTypeCache_IsComparable_anyParams) expr() *moqTypeCache_IsComparable_recorder {
	a.recorder.recorder.AnyParam(1)
	return a.recorder
}

func (a *moqTypeCache_IsComparable_anyParams) parentType() *moqTypeCache_IsComparable_recorder {
	a.recorder.recorder.AnyParam(2)
	return a.recorder
}

func (r *moqTypeCache_IsComparable_recorder) seq() *moqTypeCache_IsComparable_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(true, "seq", false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_IsComparable_recorder) noSeq() *moqTypeCache_IsComparable_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(false, "noSeq", false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_IsComparable_recorder) returnResults(result1 bool, result2 error) *moqTypeCache_IsComparable_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.ReturnResults(moqTypeCache_IsComparable_results{
		result1: result1,
		result2: result2,
	})
	return r
}

func (r *moqTypeCache_IsComparable_recorder) andDo(fn moqTypeCache_IsComparable_doFn) *moqTypeCache_IsComparable_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.AndDo(func(params moqTypeCache_IsComparable_params) {
		fn(params.expr, params.parentType)
	}, false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_IsComparable_recorder) doReturnResults(fn moqTypeCache_IsComparable_doReturnFn) *moqTypeCache_IsComparable_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.DoReturnResults(func(params moqTypeCache_IsComparable_params) *moqTypeCache_IsComparable_results {
		result1, result2 := fn(params.expr, params.parentType)
		return &moqTypeCache_IsComparable_results{
			result1: result1,
			result2: result2,
		}
	})
	return r
}

func (r *moqTypeCache_IsComparable_recorder) repeat(repeaters ...moq.Repeater) *moqTypeCache_IsComparable_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Repeat(repeaters, false) {
		return nil
	}
	return r
}

func (*moqTypeCache_IsComparable_adaptor) PrettyParams(params moqTypeCache_IsComparable_params) string {
	return fmt.Sprintf("IsComparable(%#v, %#v)", params.expr, params.parentType)
}

func (a *moqTypeCache_IsComparable_adaptor) ParamsKey(params moqTypeCache_IsComparable_params, anyParams uint64) moqTypeCache_IsComparable_paramsKey {
	a.moq.moq_IsComparable.Scene.T.Helper()
	exprUsed, exprUsedHash := impl.ParamKey(
		params.expr, 1, a.moq.runtime.parameterIndexing.IsComparable.expr, anyParams)
	parentTypeUsed, parentTypeUsedHash := impl.ParamKey(
		params.parentType, 2, a.moq.runtime.parameterIndexing.IsComparable.parentType, anyParams)
	return moqTypeCache_IsComparable_paramsKey{
		params: struct {
			expr       dst.Expr
			parentType ast.TypeInfo
		}{
			expr:       exprUsed,
			parentType: parentTypeUsed,
		},
		hashes: struct {
			expr       hash.Hash
			parentType hash.Hash
		}{
			expr:       exprUsedHash,
			parentType: parentTypeUsedHash,
		},
	}
}

func (m *moqTypeCache_recorder) IsDefaultComparable(expr dst.Expr, parentType ast.TypeInfo) *moqTypeCache_IsDefaultComparable_recorder {
	return &moqTypeCache_IsDefaultComparable_recorder{
		recorder: m.moq.moq_IsDefaultComparable.OnCall(moqTypeCache_IsDefaultComparable_params{
			expr:       expr,
			parentType: parentType,
		}),
	}
}

func (r *moqTypeCache_IsDefaultComparable_recorder) any() *moqTypeCache_IsDefaultComparable_anyParams {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.IsAnyPermitted(false) {
		return nil
	}
	return &moqTypeCache_IsDefaultComparable_anyParams{recorder: r}
}

func (a *moqTypeCache_IsDefaultComparable_anyParams) expr() *moqTypeCache_IsDefaultComparable_recorder {
	a.recorder.recorder.AnyParam(1)
	return a.recorder
}

func (a *moqTypeCache_IsDefaultComparable_anyParams) parentType() *moqTypeCache_IsDefaultComparable_recorder {
	a.recorder.recorder.AnyParam(2)
	return a.recorder
}

func (r *moqTypeCache_IsDefaultComparable_recorder) seq() *moqTypeCache_IsDefaultComparable_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(true, "seq", false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_IsDefaultComparable_recorder) noSeq() *moqTypeCache_IsDefaultComparable_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(false, "noSeq", false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_IsDefaultComparable_recorder) returnResults(result1 bool, result2 error) *moqTypeCache_IsDefaultComparable_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.ReturnResults(moqTypeCache_IsDefaultComparable_results{
		result1: result1,
		result2: result2,
	})
	return r
}

func (r *moqTypeCache_IsDefaultComparable_recorder) andDo(fn moqTypeCache_IsDefaultComparable_doFn) *moqTypeCache_IsDefaultComparable_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.AndDo(func(params moqTypeCache_IsDefaultComparable_params) {
		fn(params.expr, params.parentType)
	}, false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_IsDefaultComparable_recorder) doReturnResults(fn moqTypeCache_IsDefaultComparable_doReturnFn) *moqTypeCache_IsDefaultComparable_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.DoReturnResults(func(params moqTypeCache_IsDefaultComparable_params) *moqTypeCache_IsDefaultComparable_results {
		result1, result2 := fn(params.expr, params.parentType)
		return &moqTypeCache_IsDefaultComparable_results{
			result1: result1,
			result2: result2,
		}
	})
	return r
}

func (r *moqTypeCache_IsDefaultComparable_recorder) repeat(repeaters ...moq.Repeater) *moqTypeCache_IsDefaultComparable_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Repeat(repeaters, false) {
		return nil
	}
	return r
}

func (*moqTypeCache_IsDefaultComparable_adaptor) PrettyParams(params moqTypeCache_IsDefaultComparable_params) string {
	return fmt.Sprintf("IsDefaultComparable(%#v, %#v)", params.expr, params.parentType)
}

func (a *moqTypeCache_IsDefaultComparable_adaptor) ParamsKey(params moqTypeCache_IsDefaultComparable_params, anyParams uint64) moqTypeCache_IsDefaultComparable_paramsKey {
	a.moq.moq_IsDefaultComparable.Scene.T.Helper()
	exprUsed, exprUsedHash := impl.ParamKey(
		params.expr, 1, a.moq.runtime.parameterIndexing.IsDefaultComparable.expr, anyParams)
	parentTypeUsed, parentTypeUsedHash := impl.ParamKey(
		params.parentType, 2, a.moq.runtime.parameterIndexing.IsDefaultComparable.parentType, anyParams)
	return moqTypeCache_IsDefaultComparable_paramsKey{
		params: struct {
			expr       dst.Expr
			parentType ast.TypeInfo
		}{
			expr:       exprUsed,
			parentType: parentTypeUsed,
		},
		hashes: struct {
			expr       hash.Hash
			parentType hash.Hash
		}{
			expr:       exprUsedHash,
			parentType: parentTypeUsedHash,
		},
	}
}

func (m *moqTypeCache_recorder) FindPackage(dir string) *moqTypeCache_FindPackage_recorder {
	return &moqTypeCache_FindPackage_recorder{
		recorder: m.moq.moq_FindPackage.OnCall(moqTypeCache_FindPackage_params{
			dir: dir,
		}),
	}
}

func (r *moqTypeCache_FindPackage_recorder) any() *moqTypeCache_FindPackage_anyParams {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.IsAnyPermitted(false) {
		return nil
	}
	return &moqTypeCache_FindPackage_anyParams{recorder: r}
}

func (a *moqTypeCache_FindPackage_anyParams) dir() *moqTypeCache_FindPackage_recorder {
	a.recorder.recorder.AnyParam(1)
	return a.recorder
}

func (r *moqTypeCache_FindPackage_recorder) seq() *moqTypeCache_FindPackage_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(true, "seq", false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_FindPackage_recorder) noSeq() *moqTypeCache_FindPackage_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(false, "noSeq", false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_FindPackage_recorder) returnResults(result1 string, result2 error) *moqTypeCache_FindPackage_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.ReturnResults(moqTypeCache_FindPackage_results{
		result1: result1,
		result2: result2,
	})
	return r
}

func (r *moqTypeCache_FindPackage_recorder) andDo(fn moqTypeCache_FindPackage_doFn) *moqTypeCache_FindPackage_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.AndDo(func(params moqTypeCache_FindPackage_params) {
		fn(params.dir)
	}, false) {
		return nil
	}
	return r
}

func (r *moqTypeCache_FindPackage_recorder) doReturnResults(fn moqTypeCache_FindPackage_doReturnFn) *moqTypeCache_FindPackage_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.DoReturnResults(func(params moqTypeCache_FindPackage_params) *moqTypeCache_FindPackage_results {
		result1, result2 := fn(params.dir)
		return &moqTypeCache_FindPackage_results{
			result1: result1,
			result2: result2,
		}
	})
	return r
}

func (r *moqTypeCache_FindPackage_recorder) repeat(repeaters ...moq.Repeater) *moqTypeCache_FindPackage_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Repeat(repeaters, false) {
		return nil
	}
	return r
}

func (*moqTypeCache_FindPackage_adaptor) PrettyParams(params moqTypeCache_FindPackage_params) string {
	return fmt.Sprintf("FindPackage(%#v)", params.dir)
}

func (a *moqTypeCache_FindPackage_adaptor) ParamsKey(params moqTypeCache_FindPackage_params, anyParams uint64) moqTypeCache_FindPackage_paramsKey {
	a.moq.moq_FindPackage.Scene.T.Helper()
	dirUsed, dirUsedHash := impl.ParamKey(
		params.dir, 1, a.moq.runtime.parameterIndexing.FindPackage.dir, anyParams)
	return moqTypeCache_FindPackage_paramsKey{
		params: struct{ dir string }{
			dir: dirUsed,
		},
		hashes: struct{ dir hash.Hash }{
			dir: dirUsedHash,
		},
	}
}

// Reset resets the state of the moq
func (m *moqTypeCache) Reset() {
	m.moq_LoadPackage.Reset()
	m.moq_MockableTypes.Reset()
	m.moq_Type.Reset()
	m.moq_IsComparable.Reset()
	m.moq_IsDefaultComparable.Reset()
	m.moq_FindPackage.Reset()
}

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqTypeCache) AssertExpectationsMet() {
	m.moq_LoadPackage.Scene.T.Helper()
	m.moq_LoadPackage.AssertExpectationsMet()
	m.moq_MockableTypes.AssertExpectationsMet()
	m.moq_Type.AssertExpectationsMet()
	m.moq_IsComparable.AssertExpectationsMet()
	m.moq_IsDefaultComparable.AssertExpectationsMet()
	m.moq_FindPackage.AssertExpectationsMet()
}
