// Code generated by Moqueries - https://moqueries.org - DO NOT EDIT.

package pkgout

import (
	"fmt"

	"moqueries.org/runtime/hash"
	"moqueries.org/runtime/impl"
	"moqueries.org/runtime/moq"
)

// The following type assertion assures that testmoqs.PassByValueSimple_genType
// is mocked completely
var _ PassByValueSimple_genType = (*MoqPassByValueSimple_genType_mock)(nil)

// PassByValueSimple_genType is the fabricated implementation type of this mock
// (emitted when mocking a collections of methods directly and not from an
// interface type)
type PassByValueSimple_genType interface {
	Usual(string, bool) (string, error)
}

// MoqPassByValueSimple_genType holds the state of a moq of the
// PassByValueSimple_genType type
type MoqPassByValueSimple_genType struct {
	Moq *MoqPassByValueSimple_genType_mock

	Moq_Usual *impl.Moq[
		*MoqPassByValueSimple_genType_Usual_adaptor,
		MoqPassByValueSimple_genType_Usual_params,
		MoqPassByValueSimple_genType_Usual_paramsKey,
		MoqPassByValueSimple_genType_Usual_results,
	]

	Runtime MoqPassByValueSimple_genType_runtime
}

// MoqPassByValueSimple_genType_mock isolates the mock interface of the
// PassByValueSimple_genType type
type MoqPassByValueSimple_genType_mock struct {
	Moq *MoqPassByValueSimple_genType
}

// MoqPassByValueSimple_genType_recorder isolates the recorder interface of the
// PassByValueSimple_genType type
type MoqPassByValueSimple_genType_recorder struct {
	Moq *MoqPassByValueSimple_genType
}

// MoqPassByValueSimple_genType_runtime holds runtime configuration for the
// PassByValueSimple_genType type
type MoqPassByValueSimple_genType_runtime struct {
	ParameterIndexing struct {
		Usual MoqPassByValueSimple_genType_Usual_paramIndexing
	}
}

// MoqPassByValueSimple_genType_Usual_adaptor adapts
// MoqPassByValueSimple_genType as needed by the runtime
type MoqPassByValueSimple_genType_Usual_adaptor struct {
	Moq *MoqPassByValueSimple_genType
}

// MoqPassByValueSimple_genType_Usual_params holds the params of the
// PassByValueSimple_genType type
type MoqPassByValueSimple_genType_Usual_params struct {
	Param1 string
	Param2 bool
}

// MoqPassByValueSimple_genType_Usual_paramsKey holds the map key params of the
// PassByValueSimple_genType type
type MoqPassByValueSimple_genType_Usual_paramsKey struct {
	Params struct {
		Param1 string
		Param2 bool
	}
	Hashes struct {
		Param1 hash.Hash
		Param2 hash.Hash
	}
}

// MoqPassByValueSimple_genType_Usual_results holds the results of the
// PassByValueSimple_genType type
type MoqPassByValueSimple_genType_Usual_results struct {
	Result1 string
	Result2 error
}

// MoqPassByValueSimple_genType_Usual_paramIndexing holds the parameter
// indexing runtime configuration for the PassByValueSimple_genType type
type MoqPassByValueSimple_genType_Usual_paramIndexing struct {
	Param1 moq.ParamIndexing
	Param2 moq.ParamIndexing
}

// MoqPassByValueSimple_genType_Usual_doFn defines the type of function needed
// when calling AndDo for the PassByValueSimple_genType type
type MoqPassByValueSimple_genType_Usual_doFn func(string, bool)

// MoqPassByValueSimple_genType_Usual_doReturnFn defines the type of function
// needed when calling DoReturnResults for the PassByValueSimple_genType type
type MoqPassByValueSimple_genType_Usual_doReturnFn func(string, bool) (string, error)

// MoqPassByValueSimple_genType_Usual_recorder routes recorded function calls
// to the MoqPassByValueSimple_genType moq
type MoqPassByValueSimple_genType_Usual_recorder struct {
	Recorder *impl.Recorder[
		*MoqPassByValueSimple_genType_Usual_adaptor,
		MoqPassByValueSimple_genType_Usual_params,
		MoqPassByValueSimple_genType_Usual_paramsKey,
		MoqPassByValueSimple_genType_Usual_results,
	]
}

// MoqPassByValueSimple_genType_Usual_anyParams isolates the any params
// functions of the PassByValueSimple_genType type
type MoqPassByValueSimple_genType_Usual_anyParams struct {
	Recorder *MoqPassByValueSimple_genType_Usual_recorder
}

// NewMoqPassByValueSimple_genType creates a new moq of the
// PassByValueSimple_genType type
func NewMoqPassByValueSimple_genType(scene *moq.Scene, config *moq.Config) *MoqPassByValueSimple_genType {
	adaptor1 := &MoqPassByValueSimple_genType_Usual_adaptor{}
	m := &MoqPassByValueSimple_genType{
		Moq: &MoqPassByValueSimple_genType_mock{},

		Moq_Usual: impl.NewMoq[
			*MoqPassByValueSimple_genType_Usual_adaptor,
			MoqPassByValueSimple_genType_Usual_params,
			MoqPassByValueSimple_genType_Usual_paramsKey,
			MoqPassByValueSimple_genType_Usual_results,
		](scene, adaptor1, config),

		Runtime: MoqPassByValueSimple_genType_runtime{ParameterIndexing: struct {
			Usual MoqPassByValueSimple_genType_Usual_paramIndexing
		}{
			Usual: MoqPassByValueSimple_genType_Usual_paramIndexing{
				Param1: moq.ParamIndexByValue,
				Param2: moq.ParamIndexByValue,
			},
		}},
	}
	m.Moq.Moq = m

	adaptor1.Moq = m

	scene.AddMoq(m)
	return m
}

// Mock returns the mock implementation of the PassByValueSimple_genType type
func (m *MoqPassByValueSimple_genType) Mock() *MoqPassByValueSimple_genType_mock { return m.Moq }

func (m *MoqPassByValueSimple_genType_mock) Usual(param1 string, param2 bool) (string, error) {
	m.Moq.Moq_Usual.Scene.T.Helper()
	params := MoqPassByValueSimple_genType_Usual_params{
		Param1: param1,
		Param2: param2,
	}

	var result1 string
	var result2 error
	if result := m.Moq.Moq_Usual.Function(params); result != nil {
		result1 = result.Result1
		result2 = result.Result2
	}
	return result1, result2
}

// OnCall returns the recorder implementation of the PassByValueSimple_genType
// type
func (m *MoqPassByValueSimple_genType) OnCall() *MoqPassByValueSimple_genType_recorder {
	return &MoqPassByValueSimple_genType_recorder{
		Moq: m,
	}
}

func (m *MoqPassByValueSimple_genType_recorder) Usual(param1 string, param2 bool) *MoqPassByValueSimple_genType_Usual_recorder {
	return &MoqPassByValueSimple_genType_Usual_recorder{
		Recorder: m.Moq.Moq_Usual.OnCall(MoqPassByValueSimple_genType_Usual_params{
			Param1: param1,
			Param2: param2,
		}),
	}
}

func (r *MoqPassByValueSimple_genType_Usual_recorder) Any() *MoqPassByValueSimple_genType_Usual_anyParams {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.IsAnyPermitted(true) {
		return nil
	}
	return &MoqPassByValueSimple_genType_Usual_anyParams{Recorder: r}
}

func (a *MoqPassByValueSimple_genType_Usual_anyParams) Param1() *MoqPassByValueSimple_genType_Usual_recorder {
	a.Recorder.Recorder.AnyParam(1)
	return a.Recorder
}

func (a *MoqPassByValueSimple_genType_Usual_anyParams) Param2() *MoqPassByValueSimple_genType_Usual_recorder {
	a.Recorder.Recorder.AnyParam(2)
	return a.Recorder
}

func (r *MoqPassByValueSimple_genType_Usual_recorder) Seq() *MoqPassByValueSimple_genType_Usual_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(true, "Seq", true) {
		return nil
	}
	return r
}

func (r *MoqPassByValueSimple_genType_Usual_recorder) NoSeq() *MoqPassByValueSimple_genType_Usual_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(false, "NoSeq", true) {
		return nil
	}
	return r
}

func (r *MoqPassByValueSimple_genType_Usual_recorder) ReturnResults(result1 string, result2 error) *MoqPassByValueSimple_genType_Usual_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.ReturnResults(MoqPassByValueSimple_genType_Usual_results{
		Result1: result1,
		Result2: result2,
	})
	return r
}

func (r *MoqPassByValueSimple_genType_Usual_recorder) AndDo(fn MoqPassByValueSimple_genType_Usual_doFn) *MoqPassByValueSimple_genType_Usual_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.AndDo(func(params MoqPassByValueSimple_genType_Usual_params) {
		fn(params.Param1, params.Param2)
	}, true) {
		return nil
	}
	return r
}

func (r *MoqPassByValueSimple_genType_Usual_recorder) DoReturnResults(fn MoqPassByValueSimple_genType_Usual_doReturnFn) *MoqPassByValueSimple_genType_Usual_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.DoReturnResults(func(params MoqPassByValueSimple_genType_Usual_params) *MoqPassByValueSimple_genType_Usual_results {
		result1, result2 := fn(params.Param1, params.Param2)
		return &MoqPassByValueSimple_genType_Usual_results{
			Result1: result1,
			Result2: result2,
		}
	})
	return r
}

func (r *MoqPassByValueSimple_genType_Usual_recorder) Repeat(repeaters ...moq.Repeater) *MoqPassByValueSimple_genType_Usual_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Repeat(repeaters, true) {
		return nil
	}
	return r
}

func (*MoqPassByValueSimple_genType_Usual_adaptor) PrettyParams(params MoqPassByValueSimple_genType_Usual_params) string {
	return fmt.Sprintf("Usual(%#v, %#v)", params.Param1, params.Param2)
}

func (a *MoqPassByValueSimple_genType_Usual_adaptor) ParamsKey(params MoqPassByValueSimple_genType_Usual_params, anyParams uint64) MoqPassByValueSimple_genType_Usual_paramsKey {
	a.Moq.Moq_Usual.Scene.T.Helper()
	param1Used, param1UsedHash := impl.ParamKey(
		params.Param1, 1, a.Moq.Runtime.ParameterIndexing.Usual.Param1, anyParams)
	param2Used, param2UsedHash := impl.ParamKey(
		params.Param2, 2, a.Moq.Runtime.ParameterIndexing.Usual.Param2, anyParams)
	return MoqPassByValueSimple_genType_Usual_paramsKey{
		Params: struct {
			Param1 string
			Param2 bool
		}{
			Param1: param1Used,
			Param2: param2Used,
		},
		Hashes: struct {
			Param1 hash.Hash
			Param2 hash.Hash
		}{
			Param1: param1UsedHash,
			Param2: param2UsedHash,
		},
	}
}

// Reset resets the state of the moq
func (m *MoqPassByValueSimple_genType) Reset() {
	m.Moq_Usual.Reset()
}

// AssertExpectationsMet asserts that all expectations have been met
func (m *MoqPassByValueSimple_genType) AssertExpectationsMet() {
	m.Moq_Usual.Scene.T.Helper()
	m.Moq_Usual.AssertExpectationsMet()
}
