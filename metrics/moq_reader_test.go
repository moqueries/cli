// Code generated by Moqueries - https://moqueries.org - DO NOT EDIT.

package metrics_test

import (
	"fmt"
	"io"

	"moqueries.org/runtime/hash"
	"moqueries.org/runtime/impl"
	"moqueries.org/runtime/moq"
)

// The following type assertion assures that io.Reader is mocked completely
var _ io.Reader = (*moqReader_mock)(nil)

// moqReader holds the state of a moq of the Reader type
type moqReader struct {
	moq *moqReader_mock

	moq_Read *impl.Moq[
		*moqReader_Read_adaptor,
		moqReader_Read_params,
		moqReader_Read_paramsKey,
		moqReader_Read_results,
	]

	runtime moqReader_runtime
}

// moqReader_mock isolates the mock interface of the Reader type
type moqReader_mock struct {
	moq *moqReader
}

// moqReader_recorder isolates the recorder interface of the Reader type
type moqReader_recorder struct {
	moq *moqReader
}

// moqReader_runtime holds runtime configuration for the Reader type
type moqReader_runtime struct {
	parameterIndexing struct {
		Read moqReader_Read_paramIndexing
	}
}

// moqReader_Read_adaptor adapts moqReader as needed by the runtime
type moqReader_Read_adaptor struct {
	moq *moqReader
}

// moqReader_Read_params holds the params of the Reader type
type moqReader_Read_params struct{ p []byte }

// moqReader_Read_paramsKey holds the map key params of the Reader type
type moqReader_Read_paramsKey struct {
	params struct{}
	hashes struct{ p hash.Hash }
}

// moqReader_Read_results holds the results of the Reader type
type moqReader_Read_results struct {
	n   int
	err error
}

// moqReader_Read_paramIndexing holds the parameter indexing runtime
// configuration for the Reader type
type moqReader_Read_paramIndexing struct {
	p moq.ParamIndexing
}

// moqReader_Read_doFn defines the type of function needed when calling andDo
// for the Reader type
type moqReader_Read_doFn func(p []byte)

// moqReader_Read_doReturnFn defines the type of function needed when calling
// doReturnResults for the Reader type
type moqReader_Read_doReturnFn func(p []byte) (n int, err error)

// moqReader_Read_recorder routes recorded function calls to the moqReader moq
type moqReader_Read_recorder struct {
	recorder *impl.Recorder[
		*moqReader_Read_adaptor,
		moqReader_Read_params,
		moqReader_Read_paramsKey,
		moqReader_Read_results,
	]
}

// moqReader_Read_anyParams isolates the any params functions of the Reader
// type
type moqReader_Read_anyParams struct {
	recorder *moqReader_Read_recorder
}

// newMoqReader creates a new moq of the Reader type
func newMoqReader(scene *moq.Scene, config *moq.Config) *moqReader {
	adaptor1 := &moqReader_Read_adaptor{}
	m := &moqReader{
		moq: &moqReader_mock{},

		moq_Read: impl.NewMoq[
			*moqReader_Read_adaptor,
			moqReader_Read_params,
			moqReader_Read_paramsKey,
			moqReader_Read_results,
		](scene, adaptor1, config),

		runtime: moqReader_runtime{parameterIndexing: struct {
			Read moqReader_Read_paramIndexing
		}{
			Read: moqReader_Read_paramIndexing{
				p: moq.ParamIndexByHash,
			},
		}},
	}
	m.moq.moq = m

	adaptor1.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the mock implementation of the Reader type
func (m *moqReader) mock() *moqReader_mock { return m.moq }

func (m *moqReader_mock) Read(p []byte) (int, error) {
	m.moq.moq_Read.Scene.T.Helper()
	params := moqReader_Read_params{
		p: p,
	}

	var result1 int
	var result2 error
	if result := m.moq.moq_Read.Function(params); result != nil {
		result1 = result.n
		result2 = result.err
	}
	return result1, result2
}

// onCall returns the recorder implementation of the Reader type
func (m *moqReader) onCall() *moqReader_recorder {
	return &moqReader_recorder{
		moq: m,
	}
}

func (m *moqReader_recorder) Read(p []byte) *moqReader_Read_recorder {
	return &moqReader_Read_recorder{
		recorder: m.moq.moq_Read.OnCall(moqReader_Read_params{
			p: p,
		}),
	}
}

func (r *moqReader_Read_recorder) any() *moqReader_Read_anyParams {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.IsAnyPermitted(false) {
		return nil
	}
	return &moqReader_Read_anyParams{recorder: r}
}

func (a *moqReader_Read_anyParams) p() *moqReader_Read_recorder {
	a.recorder.recorder.AnyParam(1)
	return a.recorder
}

func (r *moqReader_Read_recorder) seq() *moqReader_Read_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(true, "seq", false) {
		return nil
	}
	return r
}

func (r *moqReader_Read_recorder) noSeq() *moqReader_Read_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Seq(false, "noSeq", false) {
		return nil
	}
	return r
}

func (r *moqReader_Read_recorder) returnResults(n int, err error) *moqReader_Read_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.ReturnResults(moqReader_Read_results{
		n:   n,
		err: err,
	})
	return r
}

func (r *moqReader_Read_recorder) andDo(fn moqReader_Read_doFn) *moqReader_Read_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.AndDo(func(params moqReader_Read_params) {
		fn(params.p)
	}, false) {
		return nil
	}
	return r
}

func (r *moqReader_Read_recorder) doReturnResults(fn moqReader_Read_doReturnFn) *moqReader_Read_recorder {
	r.recorder.Moq.Scene.T.Helper()
	r.recorder.DoReturnResults(func(params moqReader_Read_params) *moqReader_Read_results {
		n, err := fn(params.p)
		return &moqReader_Read_results{
			n:   n,
			err: err,
		}
	})
	return r
}

func (r *moqReader_Read_recorder) repeat(repeaters ...moq.Repeater) *moqReader_Read_recorder {
	r.recorder.Moq.Scene.T.Helper()
	if !r.recorder.Repeat(repeaters, false) {
		return nil
	}
	return r
}

func (*moqReader_Read_adaptor) PrettyParams(params moqReader_Read_params) string {
	return fmt.Sprintf("Read(%#v)", params.p)
}

func (a *moqReader_Read_adaptor) ParamsKey(params moqReader_Read_params, anyParams uint64) moqReader_Read_paramsKey {
	a.moq.moq_Read.Scene.T.Helper()
	pUsedHash := impl.HashOnlyParamKey(a.moq.moq_Read.Scene.T,
		params.p, "p", 1, a.moq.runtime.parameterIndexing.Read.p, anyParams)
	return moqReader_Read_paramsKey{
		params: struct{}{},
		hashes: struct{ p hash.Hash }{
			p: pUsedHash,
		},
	}
}

// Reset resets the state of the moq
func (m *moqReader) Reset() {
	m.moq_Read.Reset()
}

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqReader) AssertExpectationsMet() {
	m.moq_Read.Scene.T.Helper()
	m.moq_Read.AssertExpectationsMet()
}
