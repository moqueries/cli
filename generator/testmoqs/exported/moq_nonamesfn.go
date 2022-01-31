// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package exported

import (
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/generator/testmoqs"
	"github.com/myshkin5/moqueries/hash"
	"github.com/myshkin5/moqueries/moq"
)

// MoqNoNamesFn holds the state of a moq of the NoNamesFn type
type MoqNoNamesFn struct {
	Scene  *moq.Scene
	Config moq.Config
	Moq    *MoqNoNamesFn_mock

	ResultsByParams []MoqNoNamesFn_resultsByParams

	Runtime struct {
		ParameterIndexing struct {
			Param1 moq.ParamIndexing
			Param2 moq.ParamIndexing
		}
	}
}

// MoqNoNamesFn_mock isolates the mock interface of the NoNamesFn type
type MoqNoNamesFn_mock struct {
	Moq *MoqNoNamesFn
}

// MoqNoNamesFn_params holds the params of the NoNamesFn type
type MoqNoNamesFn_params struct {
	Param1 string
	Param2 bool
}

// MoqNoNamesFn_paramsKey holds the map key params of the NoNamesFn type
type MoqNoNamesFn_paramsKey struct {
	Params struct {
		Param1 string
		Param2 bool
	}
	Hashes struct {
		Param1 hash.Hash
		Param2 hash.Hash
	}
}

// MoqNoNamesFn_resultsByParams contains the results for a given set of parameters for the NoNamesFn type
type MoqNoNamesFn_resultsByParams struct {
	AnyCount  int
	AnyParams uint64
	Results   map[MoqNoNamesFn_paramsKey]*MoqNoNamesFn_results
}

// MoqNoNamesFn_doFn defines the type of function needed when calling AndDo for the NoNamesFn type
type MoqNoNamesFn_doFn func(string, bool)

// MoqNoNamesFn_doReturnFn defines the type of function needed when calling DoReturnResults for the NoNamesFn type
type MoqNoNamesFn_doReturnFn func(string, bool) (string, error)

// MoqNoNamesFn_results holds the results of the NoNamesFn type
type MoqNoNamesFn_results struct {
	Params  MoqNoNamesFn_params
	Results []struct {
		Values *struct {
			Result1 string
			Result2 error
		}
		Sequence   uint32
		DoFn       MoqNoNamesFn_doFn
		DoReturnFn MoqNoNamesFn_doReturnFn
	}
	Index  uint32
	Repeat *moq.RepeatVal
}

// MoqNoNamesFn_fnRecorder routes recorded function calls to the MoqNoNamesFn moq
type MoqNoNamesFn_fnRecorder struct {
	Params    MoqNoNamesFn_params
	AnyParams uint64
	Sequence  bool
	Results   *MoqNoNamesFn_results
	Moq       *MoqNoNamesFn
}

// MoqNoNamesFn_anyParams isolates the any params functions of the NoNamesFn type
type MoqNoNamesFn_anyParams struct {
	Recorder *MoqNoNamesFn_fnRecorder
}

// NewMoqNoNamesFn creates a new moq of the NoNamesFn type
func NewMoqNoNamesFn(scene *moq.Scene, config *moq.Config) *MoqNoNamesFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &MoqNoNamesFn{
		Scene:  scene,
		Config: *config,
		Moq:    &MoqNoNamesFn_mock{},

		Runtime: struct {
			ParameterIndexing struct {
				Param1 moq.ParamIndexing
				Param2 moq.ParamIndexing
			}
		}{ParameterIndexing: struct {
			Param1 moq.ParamIndexing
			Param2 moq.ParamIndexing
		}{
			Param1: moq.ParamIndexByValue,
			Param2: moq.ParamIndexByValue,
		}},
	}
	m.Moq.Moq = m

	scene.AddMoq(m)
	return m
}

// Mock returns the moq implementation of the NoNamesFn type
func (m *MoqNoNamesFn) Mock() testmoqs.NoNamesFn {
	return func(param1 string, param2 bool) (string, error) {
		moq := &MoqNoNamesFn_mock{Moq: m}
		return moq.Fn(param1, param2)
	}
}

func (m *MoqNoNamesFn_mock) Fn(param1 string, param2 bool) (result1 string, result2 error) {
	params := MoqNoNamesFn_params{
		Param1: param1,
		Param2: param2,
	}
	var results *MoqNoNamesFn_results
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

func (m *MoqNoNamesFn) OnCall(param1 string, param2 bool) *MoqNoNamesFn_fnRecorder {
	return &MoqNoNamesFn_fnRecorder{
		Params: MoqNoNamesFn_params{
			Param1: param1,
			Param2: param2,
		},
		Sequence: m.Config.Sequence == moq.SeqDefaultOn,
		Moq:      m,
	}
}

func (r *MoqNoNamesFn_fnRecorder) Any() *MoqNoNamesFn_anyParams {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Any functions must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	return &MoqNoNamesFn_anyParams{Recorder: r}
}

func (a *MoqNoNamesFn_anyParams) Param1() *MoqNoNamesFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 0
	return a.Recorder
}

func (a *MoqNoNamesFn_anyParams) Param2() *MoqNoNamesFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 1
	return a.Recorder
}

func (r *MoqNoNamesFn_fnRecorder) Seq() *MoqNoNamesFn_fnRecorder {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Seq must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.Sequence = true
	return r
}

func (r *MoqNoNamesFn_fnRecorder) NoSeq() *MoqNoNamesFn_fnRecorder {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("NoSeq must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.Sequence = false
	return r
}

func (r *MoqNoNamesFn_fnRecorder) ReturnResults(result1 string, result2 error) *MoqNoNamesFn_fnRecorder {
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
		DoFn       MoqNoNamesFn_doFn
		DoReturnFn MoqNoNamesFn_doReturnFn
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

func (r *MoqNoNamesFn_fnRecorder) AndDo(fn MoqNoNamesFn_doFn) *MoqNoNamesFn_fnRecorder {
	if r.Results == nil {
		r.Moq.Scene.T.Fatalf("ReturnResults must be called before calling AndDo")
		return nil
	}
	last := &r.Results.Results[len(r.Results.Results)-1]
	last.DoFn = fn
	return r
}

func (r *MoqNoNamesFn_fnRecorder) DoReturnResults(fn MoqNoNamesFn_doReturnFn) *MoqNoNamesFn_fnRecorder {
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
		DoFn       MoqNoNamesFn_doFn
		DoReturnFn MoqNoNamesFn_doReturnFn
	}{Sequence: sequence, DoReturnFn: fn})
	return r
}

func (r *MoqNoNamesFn_fnRecorder) FindResults() {
	if r.Results != nil {
		r.Results.Repeat.Increment(r.Moq.Scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.AnyParams)
	insertAt := -1
	var results *MoqNoNamesFn_resultsByParams
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
		results = &MoqNoNamesFn_resultsByParams{
			AnyCount:  anyCount,
			AnyParams: r.AnyParams,
			Results:   map[MoqNoNamesFn_paramsKey]*MoqNoNamesFn_results{},
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
		r.Results = &MoqNoNamesFn_results{
			Params:  r.Params,
			Results: nil,
			Index:   0,
			Repeat:  &moq.RepeatVal{},
		}
		results.Results[paramsKey] = r.Results
	}

	r.Results.Repeat.Increment(r.Moq.Scene.T)
}

func (r *MoqNoNamesFn_fnRecorder) Repeat(repeaters ...moq.Repeater) *MoqNoNamesFn_fnRecorder {
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
				DoFn       MoqNoNamesFn_doFn
				DoReturnFn MoqNoNamesFn_doReturnFn
			}{
				Values: &struct {
					Result1 string
					Result2 error
				}{
					Result1: last.Values.Result1,
					Result2: last.Values.Result2,
				},
				Sequence: r.Moq.Scene.NextRecorderSequence(),
			}
		}
		r.Results.Results = append(r.Results.Results, last)
	}
	return r
}

func (m *MoqNoNamesFn) ParamsKey(params MoqNoNamesFn_params, anyParams uint64) MoqNoNamesFn_paramsKey {
	var param1Used string
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
	return MoqNoNamesFn_paramsKey{
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
		}}
}

// Reset resets the state of the moq
func (m *MoqNoNamesFn) Reset() { m.ResultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *MoqNoNamesFn) AssertExpectationsMet() {
	for _, res := range m.ResultsByParams {
		for _, results := range res.Results {
			missing := results.Repeat.MinTimes - int(atomic.LoadUint32(&results.Index))
			if missing > 0 {
				m.Scene.T.Errorf("Expected %d additional call(s) with parameters %#v", missing, results.Params)
			}
		}
	}
}
