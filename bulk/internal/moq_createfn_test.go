// Code generated by Moqueries - https://moqueries.org - DO NOT EDIT.

package internal_test

import (
	"fmt"
	"io"

	"moqueries.org/cli/bulk/internal"
	"moqueries.org/runtime/hash"
	"moqueries.org/runtime/impl"
	"moqueries.org/runtime/moq"
)

// moqCreateFn holds the state of a moq of the CreateFn type
type moqCreateFn struct {
	moq *impl.Moq[
		*moqCreateFn_adaptor,
		moqCreateFn_params,
		moqCreateFn_paramsKey,
		moqCreateFn_results,
	]

	runtime moqCreateFn_runtime
}

// moqCreateFn_runtime holds runtime configuration for the CreateFn type
type moqCreateFn_runtime struct {
	parameterIndexing moqCreateFn_paramIndexing
}

// moqCreateFn_adaptor adapts moqCreateFn as needed by the runtime
type moqCreateFn_adaptor struct {
	moq *moqCreateFn
}

// moqCreateFn_params holds the params of the CreateFn type
type moqCreateFn_params struct{ name string }

// moqCreateFn_paramsKey holds the map key params of the CreateFn type
type moqCreateFn_paramsKey struct {
	params struct{ name string }
	hashes struct{ name hash.Hash }
}

// moqCreateFn_results holds the results of the CreateFn type
type moqCreateFn_results struct {
	file io.WriteCloser
	err  error
}

// moqCreateFn_paramIndexing holds the parameter indexing runtime configuration
// for the CreateFn type
type moqCreateFn_paramIndexing struct {
	name moq.ParamIndexing
}

// moqCreateFn_doFn defines the type of function needed when calling andDo for
// the CreateFn type
type moqCreateFn_doFn func(name string)

// moqCreateFn_doReturnFn defines the type of function needed when calling
// doReturnResults for the CreateFn type
type moqCreateFn_doReturnFn func(name string) (file io.WriteCloser, err error)

// moqCreateFn_recorder routes recorded function calls to the moqCreateFn moq
type moqCreateFn_recorder struct {
	recorder *impl.Recorder[
		*moqCreateFn_adaptor,
		moqCreateFn_params,
		moqCreateFn_paramsKey,
		moqCreateFn_results,
	]
}

// moqCreateFn_anyParams isolates the any params functions of the CreateFn type
type moqCreateFn_anyParams struct {
	recorder *moqCreateFn_recorder
}

// newMoqCreateFn creates a new moq of the CreateFn type
func newMoqCreateFn(scene *moq.Scene, config *moq.Config) *moqCreateFn {
	adaptor1 := &moqCreateFn_adaptor{}
	m := &moqCreateFn{
		moq: impl.NewMoq[
			*moqCreateFn_adaptor,
			moqCreateFn_params,
			moqCreateFn_paramsKey,
			moqCreateFn_results,
		](scene, adaptor1, config),

		runtime: moqCreateFn_runtime{parameterIndexing: moqCreateFn_paramIndexing{
			name: moq.ParamIndexByValue,
		}},
	}
	adaptor1.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the moq implementation of the CreateFn type
func (m *moqCreateFn) mock() internal.CreateFn {
	return func(name string) (io.WriteCloser, error) {
		m.moq.Scene.T.Helper()
		params := moqCreateFn_params{
			name: name,
		}

		var result1 io.WriteCloser
		var result2 error
		if result := m.moq.Function(params); result != nil {
			result1 = result.file
			result2 = result.err
		}
		return result1, result2
	}
}

func (m *moqCreateFn) onCall(name string) *moqCreateFn_recorder {
	return &moqCreateFn_recorder{
		recorder: m.moq.OnCall(moqCreateFn_params{
			name: name,
		}),
	}
}

func (r *moqCreateFn_recorder) any() *moqCreateFn_anyParams {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.IsAnyPermitted(false) {
		return nil
	}
	return &moqCreateFn_anyParams{recorder: r}
}

func (a *moqCreateFn_anyParams) name() *moqCreateFn_recorder {
	a.recorder.recorder.AnyParam(1)
	return a.recorder
}

func (r *moqCreateFn_recorder) seq() *moqCreateFn_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(true, "seq", false) {
		return nil
	}
	return r
}

func (r *moqCreateFn_recorder) noSeq() *moqCreateFn_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(false, "noSeq", false) {
		return nil
	}
	return r
}

func (r *moqCreateFn_recorder) returnResults(file io.WriteCloser, err error) *moqCreateFn_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.ReturnResults(moqCreateFn_results{
		file: file,
		err:  err,
	})
	return r
}

func (r *moqCreateFn_recorder) andDo(fn moqCreateFn_doFn) *moqCreateFn_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.AndDo(func(params moqCreateFn_params) {
		fn(params.name)
	}, false) {
		return nil
	}
	return r
}

func (r *moqCreateFn_recorder) doReturnResults(fn moqCreateFn_doReturnFn) *moqCreateFn_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.DoReturnResults(func(params moqCreateFn_params) *moqCreateFn_results {
		file, err := fn(params.name)
		return &moqCreateFn_results{
			file: file,
			err:  err,
		}
	})
	return r
}

func (r *moqCreateFn_recorder) repeat(repeaters ...moq.Repeater) *moqCreateFn_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Repeat(repeaters, false) {
		return nil
	}
	return r
}

func (*moqCreateFn_adaptor) PrettyParams(params moqCreateFn_params) string {
	return fmt.Sprintf("CreateFn(%#v)", params.name)
}

func (a *moqCreateFn_adaptor) ParamsKey(params moqCreateFn_params, anyParams uint64) moqCreateFn_paramsKey {
	a.moq.moq.Scene.T.Helper()
	nameUsed, nameUsedHash := impl.ParamKey(
		params.name, 1, a.moq.runtime.parameterIndexing.name, anyParams)
	return moqCreateFn_paramsKey{
		params: struct{ name string }{
			name: nameUsed,
		},
		hashes: struct{ name hash.Hash }{
			name: nameUsedHash,
		},
	}
}

// Reset resets the state of the moq
func (m *moqCreateFn) Reset() {
	m.moq.Reset()
}

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqCreateFn) AssertExpectationsMet() {
	m.moq.Scene.T.Helper()
	m.moq.AssertExpectationsMet()
}
