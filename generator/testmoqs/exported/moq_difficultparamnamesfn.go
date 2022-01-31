// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package exported

import (
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/generator/testmoqs"
	"github.com/myshkin5/moqueries/hash"
	"github.com/myshkin5/moqueries/moq"
)

// MoqDifficultParamNamesFn holds the state of a moq of the DifficultParamNamesFn type
type MoqDifficultParamNamesFn struct {
	Scene  *moq.Scene
	Config moq.Config
	Moq    *MoqDifficultParamNamesFn_mock

	ResultsByParams []MoqDifficultParamNamesFn_resultsByParams

	Runtime struct {
		ParameterIndexing struct {
			Param1 moq.ParamIndexing
			Param2 moq.ParamIndexing
			Param3 moq.ParamIndexing
			Param  moq.ParamIndexing
			Param5 moq.ParamIndexing
			Param6 moq.ParamIndexing
			Param7 moq.ParamIndexing
		}
	}
}

// MoqDifficultParamNamesFn_mock isolates the mock interface of the DifficultParamNamesFn type
type MoqDifficultParamNamesFn_mock struct {
	Moq *MoqDifficultParamNamesFn
}

// MoqDifficultParamNamesFn_params holds the params of the DifficultParamNamesFn type
type MoqDifficultParamNamesFn_params struct {
	Param1, Param2 bool
	Param3         string
	Param, Param5  int
	Param6, Param7 float32
}

// MoqDifficultParamNamesFn_paramsKey holds the map key params of the DifficultParamNamesFn type
type MoqDifficultParamNamesFn_paramsKey struct {
	Params struct {
		Param1, Param2 bool
		Param3         string
		Param, Param5  int
		Param6, Param7 float32
	}
	Hashes struct {
		Param1, Param2 hash.Hash
		Param3         hash.Hash
		Param, Param5  hash.Hash
		Param6, Param7 hash.Hash
	}
}

// MoqDifficultParamNamesFn_resultsByParams contains the results for a given set of parameters for the DifficultParamNamesFn type
type MoqDifficultParamNamesFn_resultsByParams struct {
	AnyCount  int
	AnyParams uint64
	Results   map[MoqDifficultParamNamesFn_paramsKey]*MoqDifficultParamNamesFn_results
}

// MoqDifficultParamNamesFn_doFn defines the type of function needed when calling AndDo for the DifficultParamNamesFn type
type MoqDifficultParamNamesFn_doFn func(m, r bool, sequence string, param, params int, result, results float32)

// MoqDifficultParamNamesFn_doReturnFn defines the type of function needed when calling DoReturnResults for the DifficultParamNamesFn type
type MoqDifficultParamNamesFn_doReturnFn func(m, r bool, sequence string, param, params int, result, results float32)

// MoqDifficultParamNamesFn_results holds the results of the DifficultParamNamesFn type
type MoqDifficultParamNamesFn_results struct {
	Params  MoqDifficultParamNamesFn_params
	Results []struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqDifficultParamNamesFn_doFn
		DoReturnFn MoqDifficultParamNamesFn_doReturnFn
	}
	Index  uint32
	Repeat *moq.RepeatVal
}

// MoqDifficultParamNamesFn_fnRecorder routes recorded function calls to the MoqDifficultParamNamesFn moq
type MoqDifficultParamNamesFn_fnRecorder struct {
	Params    MoqDifficultParamNamesFn_params
	AnyParams uint64
	Sequence  bool
	Results   *MoqDifficultParamNamesFn_results
	Moq       *MoqDifficultParamNamesFn
}

// MoqDifficultParamNamesFn_anyParams isolates the any params functions of the DifficultParamNamesFn type
type MoqDifficultParamNamesFn_anyParams struct {
	Recorder *MoqDifficultParamNamesFn_fnRecorder
}

// NewMoqDifficultParamNamesFn creates a new moq of the DifficultParamNamesFn type
func NewMoqDifficultParamNamesFn(scene *moq.Scene, config *moq.Config) *MoqDifficultParamNamesFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &MoqDifficultParamNamesFn{
		Scene:  scene,
		Config: *config,
		Moq:    &MoqDifficultParamNamesFn_mock{},

		Runtime: struct {
			ParameterIndexing struct {
				Param1 moq.ParamIndexing
				Param2 moq.ParamIndexing
				Param3 moq.ParamIndexing
				Param  moq.ParamIndexing
				Param5 moq.ParamIndexing
				Param6 moq.ParamIndexing
				Param7 moq.ParamIndexing
			}
		}{ParameterIndexing: struct {
			Param1 moq.ParamIndexing
			Param2 moq.ParamIndexing
			Param3 moq.ParamIndexing
			Param  moq.ParamIndexing
			Param5 moq.ParamIndexing
			Param6 moq.ParamIndexing
			Param7 moq.ParamIndexing
		}{
			Param1: moq.ParamIndexByValue,
			Param2: moq.ParamIndexByValue,
			Param3: moq.ParamIndexByValue,
			Param:  moq.ParamIndexByValue,
			Param5: moq.ParamIndexByValue,
			Param6: moq.ParamIndexByValue,
			Param7: moq.ParamIndexByValue,
		}},
	}
	m.Moq.Moq = m

	scene.AddMoq(m)
	return m
}

// Mock returns the moq implementation of the DifficultParamNamesFn type
func (m *MoqDifficultParamNamesFn) Mock() testmoqs.DifficultParamNamesFn {
	return func(param1, param2 bool, param3 string, param, param5 int, param6, param7 float32) {
		moq := &MoqDifficultParamNamesFn_mock{Moq: m}
		moq.Fn(param1, param2, param3, param, param5, param6, param7)
	}
}

func (m *MoqDifficultParamNamesFn_mock) Fn(param1, param2 bool, param3 string, param, param5 int, param6, param7 float32) {
	params := MoqDifficultParamNamesFn_params{
		Param1: param1,
		Param2: param2,
		Param3: param3,
		Param:  param,
		Param5: param5,
		Param6: param6,
		Param7: param7,
	}
	var results *MoqDifficultParamNamesFn_results
	for _, resultsByParams := range m.Moq.ResultsByParams {
		paramsKey := m.Moq.ParamsKey(params, resultsByParams.AnyParams)
		var ok bool
		results, ok = resultsByParams.Results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.Moq.Config.Expectation == moq.Strict {
			m.Moq.Scene.T.Fatalf("Unexpected call with parameters %#v", params)
		}
		return
	}

	i := int(atomic.AddUint32(&results.Index, 1)) - 1
	if i >= results.Repeat.ResultCount {
		if !results.Repeat.AnyTimes {
			if m.Moq.Config.Expectation == moq.Strict {
				m.Moq.Scene.T.Fatalf("Too many calls to mock with parameters %#v", params)
			}
			return
		}
		i = results.Repeat.ResultCount - 1
	}

	result := results.Results[i]
	if result.Sequence != 0 {
		sequence := m.Moq.Scene.NextMockSequence()
		if (!results.Repeat.AnyTimes && result.Sequence != sequence) || result.Sequence > sequence {
			m.Moq.Scene.T.Fatalf("Call sequence does not match %#v", params)
		}
	}

	if result.DoFn != nil {
		result.DoFn(param1, param2, param3, param, param5, param6, param7)
	}

	if result.DoReturnFn != nil {
		result.DoReturnFn(param1, param2, param3, param, param5, param6, param7)
	}
	return
}

func (m *MoqDifficultParamNamesFn) OnCall(param1, param2 bool, param3 string, param, param5 int, param6, param7 float32) *MoqDifficultParamNamesFn_fnRecorder {
	return &MoqDifficultParamNamesFn_fnRecorder{
		Params: MoqDifficultParamNamesFn_params{
			Param1: param1,
			Param2: param2,
			Param3: param3,
			Param:  param,
			Param5: param5,
			Param6: param6,
			Param7: param7,
		},
		Sequence: m.Config.Sequence == moq.SeqDefaultOn,
		Moq:      m,
	}
}

func (r *MoqDifficultParamNamesFn_fnRecorder) Any() *MoqDifficultParamNamesFn_anyParams {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Any functions must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	return &MoqDifficultParamNamesFn_anyParams{Recorder: r}
}

func (a *MoqDifficultParamNamesFn_anyParams) Param1() *MoqDifficultParamNamesFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 0
	return a.Recorder
}

func (a *MoqDifficultParamNamesFn_anyParams) Param2() *MoqDifficultParamNamesFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 1
	return a.Recorder
}

func (a *MoqDifficultParamNamesFn_anyParams) Param3() *MoqDifficultParamNamesFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 2
	return a.Recorder
}

func (a *MoqDifficultParamNamesFn_anyParams) Param() *MoqDifficultParamNamesFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 3
	return a.Recorder
}

func (a *MoqDifficultParamNamesFn_anyParams) Param5() *MoqDifficultParamNamesFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 4
	return a.Recorder
}

func (a *MoqDifficultParamNamesFn_anyParams) Param6() *MoqDifficultParamNamesFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 5
	return a.Recorder
}

func (a *MoqDifficultParamNamesFn_anyParams) Param7() *MoqDifficultParamNamesFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 6
	return a.Recorder
}

func (r *MoqDifficultParamNamesFn_fnRecorder) Seq() *MoqDifficultParamNamesFn_fnRecorder {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Seq must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.Sequence = true
	return r
}

func (r *MoqDifficultParamNamesFn_fnRecorder) NoSeq() *MoqDifficultParamNamesFn_fnRecorder {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("NoSeq must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.Sequence = false
	return r
}

func (r *MoqDifficultParamNamesFn_fnRecorder) ReturnResults() *MoqDifficultParamNamesFn_fnRecorder {
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqDifficultParamNamesFn_doFn
		DoReturnFn MoqDifficultParamNamesFn_doReturnFn
	}{
		Values:   &struct{}{},
		Sequence: sequence,
	})
	return r
}

func (r *MoqDifficultParamNamesFn_fnRecorder) AndDo(fn MoqDifficultParamNamesFn_doFn) *MoqDifficultParamNamesFn_fnRecorder {
	if r.Results == nil {
		r.Moq.Scene.T.Fatalf("ReturnResults must be called before calling AndDo")
		return nil
	}
	last := &r.Results.Results[len(r.Results.Results)-1]
	last.DoFn = fn
	return r
}

func (r *MoqDifficultParamNamesFn_fnRecorder) DoReturnResults(fn MoqDifficultParamNamesFn_doReturnFn) *MoqDifficultParamNamesFn_fnRecorder {
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqDifficultParamNamesFn_doFn
		DoReturnFn MoqDifficultParamNamesFn_doReturnFn
	}{Sequence: sequence, DoReturnFn: fn})
	return r
}

func (r *MoqDifficultParamNamesFn_fnRecorder) FindResults() {
	if r.Results != nil {
		r.Results.Repeat.Increment(r.Moq.Scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.AnyParams)
	insertAt := -1
	var results *MoqDifficultParamNamesFn_resultsByParams
	for n, res := range r.Moq.ResultsByParams {
		if res.AnyParams == r.AnyParams {
			results = &res
			break
		}
		if res.AnyCount > anyCount {
			insertAt = n
		}
	}
	if results == nil {
		results = &MoqDifficultParamNamesFn_resultsByParams{
			AnyCount:  anyCount,
			AnyParams: r.AnyParams,
			Results:   map[MoqDifficultParamNamesFn_paramsKey]*MoqDifficultParamNamesFn_results{},
		}
		r.Moq.ResultsByParams = append(r.Moq.ResultsByParams, *results)
		if insertAt != -1 && insertAt+1 < len(r.Moq.ResultsByParams) {
			copy(r.Moq.ResultsByParams[insertAt+1:], r.Moq.ResultsByParams[insertAt:0])
			r.Moq.ResultsByParams[insertAt] = *results
		}
	}

	paramsKey := r.Moq.ParamsKey(r.Params, r.AnyParams)

	var ok bool
	r.Results, ok = results.Results[paramsKey]
	if !ok {
		r.Results = &MoqDifficultParamNamesFn_results{
			Params:  r.Params,
			Results: nil,
			Index:   0,
			Repeat:  &moq.RepeatVal{},
		}
		results.Results[paramsKey] = r.Results
	}

	r.Results.Repeat.Increment(r.Moq.Scene.T)
}

func (r *MoqDifficultParamNamesFn_fnRecorder) Repeat(repeaters ...moq.Repeater) *MoqDifficultParamNamesFn_fnRecorder {
	if r.Results == nil {
		r.Moq.Scene.T.Fatalf("ReturnResults or DoReturnResults must be called before calling Repeat")
		return nil
	}
	r.Results.Repeat.Repeat(r.Moq.Scene.T, repeaters)
	last := r.Results.Results[len(r.Results.Results)-1]
	for n := 0; n < r.Results.Repeat.ResultCount-1; n++ {
		if r.Sequence {
			last = struct {
				Values     *struct{}
				Sequence   uint32
				DoFn       MoqDifficultParamNamesFn_doFn
				DoReturnFn MoqDifficultParamNamesFn_doReturnFn
			}{
				Values:   &struct{}{},
				Sequence: r.Moq.Scene.NextRecorderSequence(),
			}
		}
		r.Results.Results = append(r.Results.Results, last)
	}
	return r
}

func (m *MoqDifficultParamNamesFn) ParamsKey(params MoqDifficultParamNamesFn_params, anyParams uint64) MoqDifficultParamNamesFn_paramsKey {
	var param1Used bool
	var param1UsedHash hash.Hash
	if anyParams&(1<<0) == 0 {
		if m.Runtime.ParameterIndexing.Param1 == moq.ParamIndexByValue {
			param1Used = params.Param1
		} else {
			param1UsedHash = hash.DeepHash(params.Param1)
		}
	}
	var param2Used bool
	var param2UsedHash hash.Hash
	if anyParams&(1<<1) == 0 {
		if m.Runtime.ParameterIndexing.Param2 == moq.ParamIndexByValue {
			param2Used = params.Param2
		} else {
			param2UsedHash = hash.DeepHash(params.Param2)
		}
	}
	var param3Used string
	var param3UsedHash hash.Hash
	if anyParams&(1<<2) == 0 {
		if m.Runtime.ParameterIndexing.Param3 == moq.ParamIndexByValue {
			param3Used = params.Param3
		} else {
			param3UsedHash = hash.DeepHash(params.Param3)
		}
	}
	var paramUsed int
	var paramUsedHash hash.Hash
	if anyParams&(1<<3) == 0 {
		if m.Runtime.ParameterIndexing.Param == moq.ParamIndexByValue {
			paramUsed = params.Param
		} else {
			paramUsedHash = hash.DeepHash(params.Param)
		}
	}
	var param5Used int
	var param5UsedHash hash.Hash
	if anyParams&(1<<4) == 0 {
		if m.Runtime.ParameterIndexing.Param5 == moq.ParamIndexByValue {
			param5Used = params.Param5
		} else {
			param5UsedHash = hash.DeepHash(params.Param5)
		}
	}
	var param6Used float32
	var param6UsedHash hash.Hash
	if anyParams&(1<<5) == 0 {
		if m.Runtime.ParameterIndexing.Param6 == moq.ParamIndexByValue {
			param6Used = params.Param6
		} else {
			param6UsedHash = hash.DeepHash(params.Param6)
		}
	}
	var param7Used float32
	var param7UsedHash hash.Hash
	if anyParams&(1<<6) == 0 {
		if m.Runtime.ParameterIndexing.Param7 == moq.ParamIndexByValue {
			param7Used = params.Param7
		} else {
			param7UsedHash = hash.DeepHash(params.Param7)
		}
	}
	return MoqDifficultParamNamesFn_paramsKey{
		Params: struct {
			Param1, Param2 bool
			Param3         string
			Param, Param5  int
			Param6, Param7 float32
		}{
			Param1: param1Used,
			Param2: param2Used,
			Param3: param3Used,
			Param:  paramUsed,
			Param5: param5Used,
			Param6: param6Used,
			Param7: param7Used,
		},
		Hashes: struct {
			Param1, Param2 hash.Hash
			Param3         hash.Hash
			Param, Param5  hash.Hash
			Param6, Param7 hash.Hash
		}{
			Param1: param1UsedHash,
			Param2: param2UsedHash,
			Param3: param3UsedHash,
			Param:  paramUsedHash,
			Param5: param5UsedHash,
			Param6: param6UsedHash,
			Param7: param7UsedHash,
		}}
}

// Reset resets the state of the moq
func (m *MoqDifficultParamNamesFn) Reset() { m.ResultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *MoqDifficultParamNamesFn) AssertExpectationsMet() {
	for _, res := range m.ResultsByParams {
		for _, results := range res.Results {
			missing := results.Repeat.MinTimes - int(atomic.LoadUint32(&results.Index))
			if missing > 0 {
				m.Scene.T.Errorf("Expected %d additional call(s) with parameters %#v", missing, results.Params)
			}
		}
	}
}
