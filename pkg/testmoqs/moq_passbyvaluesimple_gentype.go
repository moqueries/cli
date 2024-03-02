// Code generated by Moqueries - https://moqueries.org - DO NOT EDIT.

package testmoqs

import (
	"fmt"
	"math/bits"
	"sync/atomic"

	"moqueries.org/runtime/hash"
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
	Scene  *moq.Scene
	Config moq.Config
	Moq    *MoqPassByValueSimple_genType_mock

	ResultsByParams_Usual []MoqPassByValueSimple_genType_Usual_resultsByParams

	Runtime struct {
		ParameterIndexing struct {
			Usual struct {
				Param1 moq.ParamIndexing
				Param2 moq.ParamIndexing
			}
		}
	}
	// MoqPassByValueSimple_genType_mock isolates the mock interface of the
}

// PassByValueSimple_genType type
type MoqPassByValueSimple_genType_mock struct {
	Moq *MoqPassByValueSimple_genType
}

// MoqPassByValueSimple_genType_recorder isolates the recorder interface of the
// PassByValueSimple_genType type
type MoqPassByValueSimple_genType_recorder struct {
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

// MoqPassByValueSimple_genType_Usual_resultsByParams contains the results for
// a given set of parameters for the PassByValueSimple_genType type
type MoqPassByValueSimple_genType_Usual_resultsByParams struct {
	AnyCount  int
	AnyParams uint64
	Results   map[MoqPassByValueSimple_genType_Usual_paramsKey]*MoqPassByValueSimple_genType_Usual_results
}

// MoqPassByValueSimple_genType_Usual_doFn defines the type of function needed
// when calling AndDo for the PassByValueSimple_genType type
type MoqPassByValueSimple_genType_Usual_doFn func(string, bool)

// MoqPassByValueSimple_genType_Usual_doReturnFn defines the type of function
// needed when calling DoReturnResults for the PassByValueSimple_genType type
type MoqPassByValueSimple_genType_Usual_doReturnFn func(string, bool) (string, error)

// MoqPassByValueSimple_genType_Usual_results holds the results of the
// PassByValueSimple_genType type
type MoqPassByValueSimple_genType_Usual_results struct {
	Params  MoqPassByValueSimple_genType_Usual_params
	Results []struct {
		Values *struct {
			Result1 string
			Result2 error
		}
		Sequence   uint32
		DoFn       MoqPassByValueSimple_genType_Usual_doFn
		DoReturnFn MoqPassByValueSimple_genType_Usual_doReturnFn
	}
	Index  uint32
	Repeat *moq.RepeatVal
}

// MoqPassByValueSimple_genType_Usual_fnRecorder routes recorded function calls
// to the MoqPassByValueSimple_genType moq
type MoqPassByValueSimple_genType_Usual_fnRecorder struct {
	Params    MoqPassByValueSimple_genType_Usual_params
	AnyParams uint64
	Sequence  bool
	Results   *MoqPassByValueSimple_genType_Usual_results
	Moq       *MoqPassByValueSimple_genType
}

// MoqPassByValueSimple_genType_Usual_anyParams isolates the any params
// functions of the PassByValueSimple_genType type
type MoqPassByValueSimple_genType_Usual_anyParams struct {
	Recorder *MoqPassByValueSimple_genType_Usual_fnRecorder
}

// NewMoqPassByValueSimple_genType creates a new moq of the
// PassByValueSimple_genType type
func NewMoqPassByValueSimple_genType(scene *moq.Scene, config *moq.Config) *MoqPassByValueSimple_genType {
	if config == nil {
		config = &moq.Config{}
	}
	m := &MoqPassByValueSimple_genType{
		Scene:  scene,
		Config: *config,
		Moq:    &MoqPassByValueSimple_genType_mock{},

		Runtime: struct {
			ParameterIndexing struct {
				Usual struct {
					Param1 moq.ParamIndexing
					Param2 moq.ParamIndexing
				}
			}
		}{ParameterIndexing: struct {
			Usual struct {
				Param1 moq.ParamIndexing
				Param2 moq.ParamIndexing
			}
		}{
			Usual: struct {
				Param1 moq.ParamIndexing
				Param2 moq.ParamIndexing
			}{
				Param1: moq.ParamIndexByValue,
				Param2: moq.ParamIndexByValue,
			},
		}},
	}
	m.Moq.Moq = m

	scene.AddMoq(m)
	return m
}

// Mock returns the mock implementation of the PassByValueSimple_genType type
func (m *MoqPassByValueSimple_genType) Mock() *MoqPassByValueSimple_genType_mock { return m.Moq }

func (m *MoqPassByValueSimple_genType_mock) Usual(param1 string, param2 bool) (result1 string, result2 error) {
	m.Moq.Scene.T.Helper()
	params := MoqPassByValueSimple_genType_Usual_params{
		Param1: param1,
		Param2: param2,
	}
	var results *MoqPassByValueSimple_genType_Usual_results
	for _, resultsByParams := range m.Moq.ResultsByParams_Usual {
		paramsKey := m.Moq.ParamsKey_Usual(params, resultsByParams.AnyParams)
		var ok bool
		results, ok = resultsByParams.Results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.Moq.Config.Expectation == moq.Strict {
			m.Moq.Scene.T.Fatalf("Unexpected call to %s", m.Moq.PrettyParams_Usual(params))
		}
		return
	}

	i := int(atomic.AddUint32(&results.Index, 1)) - 1
	if i >= results.Repeat.ResultCount {
		if !results.Repeat.AnyTimes {
			if m.Moq.Config.Expectation == moq.Strict {
				m.Moq.Scene.T.Fatalf("Too many calls to %s", m.Moq.PrettyParams_Usual(params))
			}
			return
		}
		i = results.Repeat.ResultCount - 1
	}

	result := results.Results[i]
	if result.Sequence != 0 {
		sequence := m.Moq.Scene.NextMockSequence()
		if (!results.Repeat.AnyTimes && result.Sequence != sequence) || result.Sequence > sequence {
			m.Moq.Scene.T.Fatalf("Call sequence does not match call to %s", m.Moq.PrettyParams_Usual(params))
		}
	}

	if result.DoFn != nil {
		result.DoFn(param1, param2)
	}

	if result.Values != nil {
		result1 = result.Values.Result1
		result2 = result.Values.Result2
	}
	if result.DoReturnFn != nil {
		result1, result2 = result.DoReturnFn(param1, param2)
	}
	return
}

// OnCall returns the recorder implementation of the PassByValueSimple_genType
// type
func (m *MoqPassByValueSimple_genType) OnCall() *MoqPassByValueSimple_genType_recorder {
	return &MoqPassByValueSimple_genType_recorder{
		Moq: m,
	}
}

func (m *MoqPassByValueSimple_genType_recorder) Usual(param1 string, param2 bool) *MoqPassByValueSimple_genType_Usual_fnRecorder {
	return &MoqPassByValueSimple_genType_Usual_fnRecorder{
		Params: MoqPassByValueSimple_genType_Usual_params{
			Param1: param1,
			Param2: param2,
		},
		Sequence: m.Moq.Config.Sequence == moq.SeqDefaultOn,
		Moq:      m.Moq,
	}
}

func (r *MoqPassByValueSimple_genType_Usual_fnRecorder) Any() *MoqPassByValueSimple_genType_Usual_anyParams {
	r.Moq.Scene.T.Helper()
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Any functions must be called before ReturnResults or DoReturnResults calls, recording %s", r.Moq.PrettyParams_Usual(r.Params))
		return nil
	}
	return &MoqPassByValueSimple_genType_Usual_anyParams{Recorder: r}
}

func (a *MoqPassByValueSimple_genType_Usual_anyParams) Param1() *MoqPassByValueSimple_genType_Usual_fnRecorder {
	a.Recorder.AnyParams |= 1 << 0
	return a.Recorder
}

func (a *MoqPassByValueSimple_genType_Usual_anyParams) Param2() *MoqPassByValueSimple_genType_Usual_fnRecorder {
	a.Recorder.AnyParams |= 1 << 1
	return a.Recorder
}

func (r *MoqPassByValueSimple_genType_Usual_fnRecorder) Seq() *MoqPassByValueSimple_genType_Usual_fnRecorder {
	r.Moq.Scene.T.Helper()
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Seq must be called before ReturnResults or DoReturnResults calls, recording %s", r.Moq.PrettyParams_Usual(r.Params))
		return nil
	}
	r.Sequence = true
	return r
}

func (r *MoqPassByValueSimple_genType_Usual_fnRecorder) NoSeq() *MoqPassByValueSimple_genType_Usual_fnRecorder {
	r.Moq.Scene.T.Helper()
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("NoSeq must be called before ReturnResults or DoReturnResults calls, recording %s", r.Moq.PrettyParams_Usual(r.Params))
		return nil
	}
	r.Sequence = false
	return r
}

func (r *MoqPassByValueSimple_genType_Usual_fnRecorder) ReturnResults(result1 string, result2 error) *MoqPassByValueSimple_genType_Usual_fnRecorder {
	r.Moq.Scene.T.Helper()
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values *struct {
			Result1 string
			Result2 error
		}
		Sequence   uint32
		DoFn       MoqPassByValueSimple_genType_Usual_doFn
		DoReturnFn MoqPassByValueSimple_genType_Usual_doReturnFn
	}{
		Values: &struct {
			Result1 string
			Result2 error
		}{
			Result1: result1,
			Result2: result2,
		},
		Sequence: sequence,
	})
	return r
}

func (r *MoqPassByValueSimple_genType_Usual_fnRecorder) AndDo(fn MoqPassByValueSimple_genType_Usual_doFn) *MoqPassByValueSimple_genType_Usual_fnRecorder {
	r.Moq.Scene.T.Helper()
	if r.Results == nil {
		r.Moq.Scene.T.Fatalf("ReturnResults must be called before calling AndDo")
		return nil
	}
	last := &r.Results.Results[len(r.Results.Results)-1]
	last.DoFn = fn
	return r
}

func (r *MoqPassByValueSimple_genType_Usual_fnRecorder) DoReturnResults(fn MoqPassByValueSimple_genType_Usual_doReturnFn) *MoqPassByValueSimple_genType_Usual_fnRecorder {
	r.Moq.Scene.T.Helper()
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values *struct {
			Result1 string
			Result2 error
		}
		Sequence   uint32
		DoFn       MoqPassByValueSimple_genType_Usual_doFn
		DoReturnFn MoqPassByValueSimple_genType_Usual_doReturnFn
	}{Sequence: sequence, DoReturnFn: fn})
	return r
}

func (r *MoqPassByValueSimple_genType_Usual_fnRecorder) FindResults() {
	r.Moq.Scene.T.Helper()
	if r.Results != nil {
		r.Results.Repeat.Increment(r.Moq.Scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.AnyParams)
	insertAt := -1
	var results *MoqPassByValueSimple_genType_Usual_resultsByParams
	for n, res := range r.Moq.ResultsByParams_Usual {
		if res.AnyParams == r.AnyParams {
			results = &res
			break
		}
		if res.AnyCount > anyCount {
			insertAt = n
		}
	}
	if results == nil {
		results = &MoqPassByValueSimple_genType_Usual_resultsByParams{
			AnyCount:  anyCount,
			AnyParams: r.AnyParams,
			Results:   map[MoqPassByValueSimple_genType_Usual_paramsKey]*MoqPassByValueSimple_genType_Usual_results{},
		}
		r.Moq.ResultsByParams_Usual = append(r.Moq.ResultsByParams_Usual, *results)
		if insertAt != -1 && insertAt+1 < len(r.Moq.ResultsByParams_Usual) {
			copy(r.Moq.ResultsByParams_Usual[insertAt+1:], r.Moq.ResultsByParams_Usual[insertAt:0])
			r.Moq.ResultsByParams_Usual[insertAt] = *results
		}
	}

	paramsKey := r.Moq.ParamsKey_Usual(r.Params, r.AnyParams)

	var ok bool
	r.Results, ok = results.Results[paramsKey]
	if !ok {
		r.Results = &MoqPassByValueSimple_genType_Usual_results{
			Params:  r.Params,
			Results: nil,
			Index:   0,
			Repeat:  &moq.RepeatVal{},
		}
		results.Results[paramsKey] = r.Results
	}

	r.Results.Repeat.Increment(r.Moq.Scene.T)
}

func (r *MoqPassByValueSimple_genType_Usual_fnRecorder) Repeat(repeaters ...moq.Repeater) *MoqPassByValueSimple_genType_Usual_fnRecorder {
	r.Moq.Scene.T.Helper()
	if r.Results == nil {
		r.Moq.Scene.T.Fatalf("ReturnResults or DoReturnResults must be called before calling Repeat")
		return nil
	}
	r.Results.Repeat.Repeat(r.Moq.Scene.T, repeaters)
	last := r.Results.Results[len(r.Results.Results)-1]
	for n := 0; n < r.Results.Repeat.ResultCount-1; n++ {
		if r.Sequence {
			last = struct {
				Values *struct {
					Result1 string
					Result2 error
				}
				Sequence   uint32
				DoFn       MoqPassByValueSimple_genType_Usual_doFn
				DoReturnFn MoqPassByValueSimple_genType_Usual_doReturnFn
			}{
				Values:   last.Values,
				Sequence: r.Moq.Scene.NextRecorderSequence(),
			}
		}
		r.Results.Results = append(r.Results.Results, last)
	}
	return r
}

func (m *MoqPassByValueSimple_genType) PrettyParams_Usual(params MoqPassByValueSimple_genType_Usual_params) string {
	return fmt.Sprintf("Usual(%#v, %#v)", params.Param1, params.Param2)
}

func (m *MoqPassByValueSimple_genType) ParamsKey_Usual(params MoqPassByValueSimple_genType_Usual_params, anyParams uint64) MoqPassByValueSimple_genType_Usual_paramsKey {
	m.Scene.T.Helper()
	var param1Used string
	var param1UsedHash hash.Hash
	if anyParams&(1<<0) == 0 {
		if m.Runtime.ParameterIndexing.Usual.Param1 == moq.ParamIndexByValue {
			param1Used = params.Param1
		} else {
			param1UsedHash = hash.DeepHash(params.Param1)
		}
	}
	var param2Used bool
	var param2UsedHash hash.Hash
	if anyParams&(1<<1) == 0 {
		if m.Runtime.ParameterIndexing.Usual.Param2 == moq.ParamIndexByValue {
			param2Used = params.Param2
		} else {
			param2UsedHash = hash.DeepHash(params.Param2)
		}
	}
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
func (m *MoqPassByValueSimple_genType) Reset() { m.ResultsByParams_Usual = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *MoqPassByValueSimple_genType) AssertExpectationsMet() {
	m.Scene.T.Helper()
	for _, res := range m.ResultsByParams_Usual {
		for _, results := range res.Results {
			missing := results.Repeat.MinTimes - int(atomic.LoadUint32(&results.Index))
			if missing > 0 {
				m.Scene.T.Errorf("Expected %d additional call(s) to %s", missing, m.PrettyParams_Usual(results.Params))
			}
		}
	}
}
