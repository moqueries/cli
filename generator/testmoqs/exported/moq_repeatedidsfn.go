// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package exported

import (
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/generator/testmoqs"
	"github.com/myshkin5/moqueries/hash"
	"github.com/myshkin5/moqueries/moq"
)

// MoqRepeatedIdsFn holds the state of a moq of the RepeatedIdsFn type
type MoqRepeatedIdsFn struct {
	Scene  *moq.Scene
	Config moq.Config
	Moq    *MoqRepeatedIdsFn_mock

	ResultsByParams []MoqRepeatedIdsFn_resultsByParams

	Runtime struct {
		ParameterIndexing struct {
			SParam1 moq.ParamIndexing
			SParam2 moq.ParamIndexing
			BParam  moq.ParamIndexing
		}
	}
}

// MoqRepeatedIdsFn_mock isolates the mock interface of the RepeatedIdsFn type
type MoqRepeatedIdsFn_mock struct {
	Moq *MoqRepeatedIdsFn
}

// MoqRepeatedIdsFn_params holds the params of the RepeatedIdsFn type
type MoqRepeatedIdsFn_params struct {
	SParam1, SParam2 string
	BParam           bool
}

// MoqRepeatedIdsFn_paramsKey holds the map key params of the RepeatedIdsFn type
type MoqRepeatedIdsFn_paramsKey struct {
	Params struct {
		SParam1, SParam2 string
		BParam           bool
	}
	Hashes struct {
		SParam1, SParam2 hash.Hash
		BParam           hash.Hash
	}
}

// MoqRepeatedIdsFn_resultsByParams contains the results for a given set of parameters for the RepeatedIdsFn type
type MoqRepeatedIdsFn_resultsByParams struct {
	AnyCount  int
	AnyParams uint64
	Results   map[MoqRepeatedIdsFn_paramsKey]*MoqRepeatedIdsFn_results
}

// MoqRepeatedIdsFn_doFn defines the type of function needed when calling AndDo for the RepeatedIdsFn type
type MoqRepeatedIdsFn_doFn func(sParam1, sParam2 string, bParam bool)

// MoqRepeatedIdsFn_doReturnFn defines the type of function needed when calling DoReturnResults for the RepeatedIdsFn type
type MoqRepeatedIdsFn_doReturnFn func(sParam1, sParam2 string, bParam bool) (sResult1, sResult2 string, err error)

// MoqRepeatedIdsFn_results holds the results of the RepeatedIdsFn type
type MoqRepeatedIdsFn_results struct {
	Params  MoqRepeatedIdsFn_params
	Results []struct {
		Values *struct {
			SResult1, SResult2 string
			Err                error
		}
		Sequence   uint32
		DoFn       MoqRepeatedIdsFn_doFn
		DoReturnFn MoqRepeatedIdsFn_doReturnFn
	}
	Index  uint32
	Repeat *moq.RepeatVal
}

// MoqRepeatedIdsFn_fnRecorder routes recorded function calls to the MoqRepeatedIdsFn moq
type MoqRepeatedIdsFn_fnRecorder struct {
	Params    MoqRepeatedIdsFn_params
	AnyParams uint64
	Sequence  bool
	Results   *MoqRepeatedIdsFn_results
	Moq       *MoqRepeatedIdsFn
}

// MoqRepeatedIdsFn_anyParams isolates the any params functions of the RepeatedIdsFn type
type MoqRepeatedIdsFn_anyParams struct {
	Recorder *MoqRepeatedIdsFn_fnRecorder
}

// NewMoqRepeatedIdsFn creates a new moq of the RepeatedIdsFn type
func NewMoqRepeatedIdsFn(scene *moq.Scene, config *moq.Config) *MoqRepeatedIdsFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &MoqRepeatedIdsFn{
		Scene:  scene,
		Config: *config,
		Moq:    &MoqRepeatedIdsFn_mock{},

		Runtime: struct {
			ParameterIndexing struct {
				SParam1 moq.ParamIndexing
				SParam2 moq.ParamIndexing
				BParam  moq.ParamIndexing
			}
		}{ParameterIndexing: struct {
			SParam1 moq.ParamIndexing
			SParam2 moq.ParamIndexing
			BParam  moq.ParamIndexing
		}{
			SParam1: moq.ParamIndexByValue,
			SParam2: moq.ParamIndexByValue,
			BParam:  moq.ParamIndexByValue,
		}},
	}
	m.Moq.Moq = m

	scene.AddMoq(m)
	return m
}

// Mock returns the moq implementation of the RepeatedIdsFn type
func (m *MoqRepeatedIdsFn) Mock() testmoqs.RepeatedIdsFn {
	return func(sParam1, sParam2 string, bParam bool) (_, _ string, _ error) {
		moq := &MoqRepeatedIdsFn_mock{Moq: m}
		return moq.Fn(sParam1, sParam2, bParam)
	}
}

func (m *MoqRepeatedIdsFn_mock) Fn(sParam1, sParam2 string, bParam bool) (sResult1, sResult2 string, err error) {
	params := MoqRepeatedIdsFn_params{
		SParam1: sParam1,
		SParam2: sParam2,
		BParam:  bParam,
	}
	var results *MoqRepeatedIdsFn_results
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
		result.DoFn(sParam1, sParam2, bParam)
	}

	if result.Values != nil {
		sResult1 = result.Values.SResult1
		sResult2 = result.Values.SResult2
		err = result.Values.Err
	}
	if result.DoReturnFn != nil {
		sResult1, sResult2, err = result.DoReturnFn(sParam1, sParam2, bParam)
	}
	return
}

func (m *MoqRepeatedIdsFn) OnCall(sParam1, sParam2 string, bParam bool) *MoqRepeatedIdsFn_fnRecorder {
	return &MoqRepeatedIdsFn_fnRecorder{
		Params: MoqRepeatedIdsFn_params{
			SParam1: sParam1,
			SParam2: sParam2,
			BParam:  bParam,
		},
		Sequence: m.Config.Sequence == moq.SeqDefaultOn,
		Moq:      m,
	}
}

func (r *MoqRepeatedIdsFn_fnRecorder) Any() *MoqRepeatedIdsFn_anyParams {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Any functions must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	return &MoqRepeatedIdsFn_anyParams{Recorder: r}
}

func (a *MoqRepeatedIdsFn_anyParams) SParam1() *MoqRepeatedIdsFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 0
	return a.Recorder
}

func (a *MoqRepeatedIdsFn_anyParams) SParam2() *MoqRepeatedIdsFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 1
	return a.Recorder
}

func (a *MoqRepeatedIdsFn_anyParams) BParam() *MoqRepeatedIdsFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 2
	return a.Recorder
}

func (r *MoqRepeatedIdsFn_fnRecorder) Seq() *MoqRepeatedIdsFn_fnRecorder {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Seq must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.Sequence = true
	return r
}

func (r *MoqRepeatedIdsFn_fnRecorder) NoSeq() *MoqRepeatedIdsFn_fnRecorder {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("NoSeq must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.Sequence = false
	return r
}

func (r *MoqRepeatedIdsFn_fnRecorder) ReturnResults(sResult1, sResult2 string, err error) *MoqRepeatedIdsFn_fnRecorder {
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values *struct {
			SResult1, SResult2 string
			Err                error
		}
		Sequence   uint32
		DoFn       MoqRepeatedIdsFn_doFn
		DoReturnFn MoqRepeatedIdsFn_doReturnFn
	}{
		Values: &struct {
			SResult1, SResult2 string
			Err                error
		}{
			SResult1: sResult1,
			SResult2: sResult2,
			Err:      err,
		},
		Sequence: sequence,
	})
	return r
}

func (r *MoqRepeatedIdsFn_fnRecorder) AndDo(fn MoqRepeatedIdsFn_doFn) *MoqRepeatedIdsFn_fnRecorder {
	if r.Results == nil {
		r.Moq.Scene.T.Fatalf("ReturnResults must be called before calling AndDo")
		return nil
	}
	last := &r.Results.Results[len(r.Results.Results)-1]
	last.DoFn = fn
	return r
}

func (r *MoqRepeatedIdsFn_fnRecorder) DoReturnResults(fn MoqRepeatedIdsFn_doReturnFn) *MoqRepeatedIdsFn_fnRecorder {
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values *struct {
			SResult1, SResult2 string
			Err                error
		}
		Sequence   uint32
		DoFn       MoqRepeatedIdsFn_doFn
		DoReturnFn MoqRepeatedIdsFn_doReturnFn
	}{Sequence: sequence, DoReturnFn: fn})
	return r
}

func (r *MoqRepeatedIdsFn_fnRecorder) FindResults() {
	if r.Results != nil {
		r.Results.Repeat.Increment(r.Moq.Scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.AnyParams)
	insertAt := -1
	var results *MoqRepeatedIdsFn_resultsByParams
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
		results = &MoqRepeatedIdsFn_resultsByParams{
			AnyCount:  anyCount,
			AnyParams: r.AnyParams,
			Results:   map[MoqRepeatedIdsFn_paramsKey]*MoqRepeatedIdsFn_results{},
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
		r.Results = &MoqRepeatedIdsFn_results{
			Params:  r.Params,
			Results: nil,
			Index:   0,
			Repeat:  &moq.RepeatVal{},
		}
		results.Results[paramsKey] = r.Results
	}

	r.Results.Repeat.Increment(r.Moq.Scene.T)
}

func (r *MoqRepeatedIdsFn_fnRecorder) Repeat(repeaters ...moq.Repeater) *MoqRepeatedIdsFn_fnRecorder {
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
					SResult1, SResult2 string
					Err                error
				}
				Sequence   uint32
				DoFn       MoqRepeatedIdsFn_doFn
				DoReturnFn MoqRepeatedIdsFn_doReturnFn
			}{
				Values: &struct {
					SResult1, SResult2 string
					Err                error
				}{
					SResult1: last.Values.SResult1,
					SResult2: last.Values.SResult2,
					Err:      last.Values.Err,
				},
				Sequence: r.Moq.Scene.NextRecorderSequence(),
			}
		}
		r.Results.Results = append(r.Results.Results, last)
	}
	return r
}

func (m *MoqRepeatedIdsFn) ParamsKey(params MoqRepeatedIdsFn_params, anyParams uint64) MoqRepeatedIdsFn_paramsKey {
	var sParam1Used string
	var sParam1UsedHash hash.Hash
	if anyParams&(1<<0) == 0 {
		if m.Runtime.ParameterIndexing.SParam1 == moq.ParamIndexByValue {
			sParam1Used = params.SParam1
		} else {
			sParam1UsedHash = hash.DeepHash(params.SParam1)
		}
	}
	var sParam2Used string
	var sParam2UsedHash hash.Hash
	if anyParams&(1<<1) == 0 {
		if m.Runtime.ParameterIndexing.SParam2 == moq.ParamIndexByValue {
			sParam2Used = params.SParam2
		} else {
			sParam2UsedHash = hash.DeepHash(params.SParam2)
		}
	}
	var bParamUsed bool
	var bParamUsedHash hash.Hash
	if anyParams&(1<<2) == 0 {
		if m.Runtime.ParameterIndexing.BParam == moq.ParamIndexByValue {
			bParamUsed = params.BParam
		} else {
			bParamUsedHash = hash.DeepHash(params.BParam)
		}
	}
	return MoqRepeatedIdsFn_paramsKey{
		Params: struct {
			SParam1, SParam2 string
			BParam           bool
		}{
			SParam1: sParam1Used,
			SParam2: sParam2Used,
			BParam:  bParamUsed,
		},
		Hashes: struct {
			SParam1, SParam2 hash.Hash
			BParam           hash.Hash
		}{
			SParam1: sParam1UsedHash,
			SParam2: sParam2UsedHash,
			BParam:  bParamUsedHash,
		}}
}

// Reset resets the state of the moq
func (m *MoqRepeatedIdsFn) Reset() { m.ResultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *MoqRepeatedIdsFn) AssertExpectationsMet() {
	for _, res := range m.ResultsByParams {
		for _, results := range res.Results {
			missing := results.Repeat.MinTimes - int(atomic.LoadUint32(&results.Index))
			if missing > 0 {
				m.Scene.T.Errorf("Expected %d additional call(s) with parameters %#v", missing, results.Params)
			}
		}
	}
}
