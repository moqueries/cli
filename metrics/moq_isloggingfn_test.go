// Code generated by Moqueries - https://moqueries.org - DO NOT EDIT.

package metrics_test

import (
	"fmt"

	"moqueries.org/cli/metrics"
	"moqueries.org/runtime/impl"
	"moqueries.org/runtime/moq"
)

// moqIsLoggingFn holds the state of a moq of the IsLoggingFn type
type moqIsLoggingFn struct {
	moq *impl.Moq[
		*moqIsLoggingFn_adaptor,
		moqIsLoggingFn_params,
		moqIsLoggingFn_paramsKey,
		moqIsLoggingFn_results,
	]

	runtime moqIsLoggingFn_runtime
}

// moqIsLoggingFn_runtime holds runtime configuration for the IsLoggingFn type
type moqIsLoggingFn_runtime struct {
	parameterIndexing moqIsLoggingFn_paramIndexing
}

// moqIsLoggingFn_adaptor adapts moqIsLoggingFn as needed by the runtime
type moqIsLoggingFn_adaptor struct {
	moq *moqIsLoggingFn
}

// moqIsLoggingFn_params holds the params of the IsLoggingFn type
type moqIsLoggingFn_params struct{}

// moqIsLoggingFn_paramsKey holds the map key params of the IsLoggingFn type
type moqIsLoggingFn_paramsKey struct {
	params struct{}
	hashes struct{}
}

// moqIsLoggingFn_results holds the results of the IsLoggingFn type
type moqIsLoggingFn_results struct {
	result1 bool
}

// moqIsLoggingFn_paramIndexing holds the parameter indexing runtime
// configuration for the IsLoggingFn type
type moqIsLoggingFn_paramIndexing struct{}

// moqIsLoggingFn_doFn defines the type of function needed when calling andDo
// for the IsLoggingFn type
type moqIsLoggingFn_doFn func()

// moqIsLoggingFn_doReturnFn defines the type of function needed when calling
// doReturnResults for the IsLoggingFn type
type moqIsLoggingFn_doReturnFn func() bool

// moqIsLoggingFn_recorder routes recorded function calls to the moqIsLoggingFn
// moq
type moqIsLoggingFn_recorder struct {
	recorder *impl.Recorder[
		*moqIsLoggingFn_adaptor,
		moqIsLoggingFn_params,
		moqIsLoggingFn_paramsKey,
		moqIsLoggingFn_results,
	]
}

// moqIsLoggingFn_anyParams isolates the any params functions of the
// IsLoggingFn type
type moqIsLoggingFn_anyParams struct {
	recorder *moqIsLoggingFn_recorder
}

// newMoqIsLoggingFn creates a new moq of the IsLoggingFn type
func newMoqIsLoggingFn(scene *moq.Scene, config *moq.Config) *moqIsLoggingFn {
	adaptor1 := &moqIsLoggingFn_adaptor{}
	m := &moqIsLoggingFn{
		moq: impl.NewMoq[
			*moqIsLoggingFn_adaptor,
			moqIsLoggingFn_params,
			moqIsLoggingFn_paramsKey,
			moqIsLoggingFn_results,
		](scene, adaptor1, config),

		runtime: moqIsLoggingFn_runtime{parameterIndexing: moqIsLoggingFn_paramIndexing{}},
	}
	adaptor1.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the moq implementation of the IsLoggingFn type
func (m *moqIsLoggingFn) mock() metrics.IsLoggingFn {
	return func() bool {
		m.moq.Scene.T.Helper()
		params := moqIsLoggingFn_params{}

		var result1 bool
		if result := m.moq.Function(params); result != nil {
			result1 = result.result1
		}
		return result1
	}
}

func (m *moqIsLoggingFn) onCall() *moqIsLoggingFn_recorder {
	return &moqIsLoggingFn_recorder{
		recorder: m.moq.OnCall(moqIsLoggingFn_params{}),
	}
}

func (r *moqIsLoggingFn_recorder) any() *moqIsLoggingFn_anyParams {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.IsAnyPermitted(false) {
		return nil
	}
	return &moqIsLoggingFn_anyParams{recorder: r}
}

func (r *moqIsLoggingFn_recorder) seq() *moqIsLoggingFn_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(true, "seq", false) {
		return nil
	}
	return r
}

func (r *moqIsLoggingFn_recorder) noSeq() *moqIsLoggingFn_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(false, "noSeq", false) {
		return nil
	}
	return r
}

func (r *moqIsLoggingFn_recorder) returnResults(result1 bool) *moqIsLoggingFn_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.ReturnResults(moqIsLoggingFn_results{
		result1: result1,
	})
	return r
}

func (r *moqIsLoggingFn_recorder) andDo(fn moqIsLoggingFn_doFn) *moqIsLoggingFn_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.AndDo(func(params moqIsLoggingFn_params) {
		fn()
	}, false) {
		return nil
	}
	return r
}

func (r *moqIsLoggingFn_recorder) doReturnResults(fn moqIsLoggingFn_doReturnFn) *moqIsLoggingFn_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.DoReturnResults(func(params moqIsLoggingFn_params) *moqIsLoggingFn_results {
		result1 := fn()
		return &moqIsLoggingFn_results{
			result1: result1,
		}
	})
	return r
}

func (r *moqIsLoggingFn_recorder) repeat(repeaters ...moq.Repeater) *moqIsLoggingFn_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Repeat(repeaters, false) {
		return nil
	}
	return r
}

func (*moqIsLoggingFn_adaptor) PrettyParams(params moqIsLoggingFn_params) string {
	return fmt.Sprintf("IsLoggingFn()")
}

func (a *moqIsLoggingFn_adaptor) ParamsKey(params moqIsLoggingFn_params, anyParams uint64) moqIsLoggingFn_paramsKey {
	a.moq.moq.Scene.T.Helper()
	return moqIsLoggingFn_paramsKey{
		params: struct{}{},
		hashes: struct{}{},
	}
}

// Reset resets the state of the moq
func (m *moqIsLoggingFn) Reset() {
	m.moq.Reset()
}

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqIsLoggingFn) AssertExpectationsMet() {
	m.moq.Scene.T.Helper()
	m.moq.AssertExpectationsMet()
}
