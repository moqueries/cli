// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package exported

import (
	"io"
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/generator/testmoqs"
	"github.com/myshkin5/moqueries/hash"
	"github.com/myshkin5/moqueries/moq"
)

// MoqInterfaceResultFn holds the state of a moq of the InterfaceResultFn type
type MoqInterfaceResultFn struct {
	Scene  *moq.Scene
	Config moq.Config
	Moq    *MoqInterfaceResultFn_mock

	ResultsByParams []MoqInterfaceResultFn_resultsByParams

	Runtime struct {
		ParameterIndexing struct {
			SParam moq.ParamIndexing
			BParam moq.ParamIndexing
		}
	}
}

// MoqInterfaceResultFn_mock isolates the mock interface of the InterfaceResultFn type
type MoqInterfaceResultFn_mock struct {
	Moq *MoqInterfaceResultFn
}

// MoqInterfaceResultFn_params holds the params of the InterfaceResultFn type
type MoqInterfaceResultFn_params struct {
	SParam string
	BParam bool
}

// MoqInterfaceResultFn_paramsKey holds the map key params of the InterfaceResultFn type
type MoqInterfaceResultFn_paramsKey struct {
	Params struct {
		SParam string
		BParam bool
	}
	Hashes struct {
		SParam hash.Hash
		BParam hash.Hash
	}
}

// MoqInterfaceResultFn_resultsByParams contains the results for a given set of parameters for the InterfaceResultFn type
type MoqInterfaceResultFn_resultsByParams struct {
	AnyCount  int
	AnyParams uint64
	Results   map[MoqInterfaceResultFn_paramsKey]*MoqInterfaceResultFn_results
}

// MoqInterfaceResultFn_doFn defines the type of function needed when calling AndDo for the InterfaceResultFn type
type MoqInterfaceResultFn_doFn func(sParam string, bParam bool)

// MoqInterfaceResultFn_doReturnFn defines the type of function needed when calling DoReturnResults for the InterfaceResultFn type
type MoqInterfaceResultFn_doReturnFn func(sParam string, bParam bool) (r io.Reader)

// MoqInterfaceResultFn_results holds the results of the InterfaceResultFn type
type MoqInterfaceResultFn_results struct {
	Params  MoqInterfaceResultFn_params
	Results []struct {
		Values     *struct{ Result1 io.Reader }
		Sequence   uint32
		DoFn       MoqInterfaceResultFn_doFn
		DoReturnFn MoqInterfaceResultFn_doReturnFn
	}
	Index  uint32
	Repeat *moq.RepeatVal
}

// MoqInterfaceResultFn_fnRecorder routes recorded function calls to the MoqInterfaceResultFn moq
type MoqInterfaceResultFn_fnRecorder struct {
	Params    MoqInterfaceResultFn_params
	AnyParams uint64
	Sequence  bool
	Results   *MoqInterfaceResultFn_results
	Moq       *MoqInterfaceResultFn
}

// MoqInterfaceResultFn_anyParams isolates the any params functions of the InterfaceResultFn type
type MoqInterfaceResultFn_anyParams struct {
	Recorder *MoqInterfaceResultFn_fnRecorder
}

// NewMoqInterfaceResultFn creates a new moq of the InterfaceResultFn type
func NewMoqInterfaceResultFn(scene *moq.Scene, config *moq.Config) *MoqInterfaceResultFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &MoqInterfaceResultFn{
		Scene:  scene,
		Config: *config,
		Moq:    &MoqInterfaceResultFn_mock{},

		Runtime: struct {
			ParameterIndexing struct {
				SParam moq.ParamIndexing
				BParam moq.ParamIndexing
			}
		}{ParameterIndexing: struct {
			SParam moq.ParamIndexing
			BParam moq.ParamIndexing
		}{
			SParam: moq.ParamIndexByValue,
			BParam: moq.ParamIndexByValue,
		}},
	}
	m.Moq.Moq = m

	scene.AddMoq(m)
	return m
}

// Mock returns the moq implementation of the InterfaceResultFn type
func (m *MoqInterfaceResultFn) Mock() testmoqs.InterfaceResultFn {
	return func(sParam string, bParam bool) (_ io.Reader) {
		moq := &MoqInterfaceResultFn_mock{Moq: m}
		return moq.Fn(sParam, bParam)
	}
}

func (m *MoqInterfaceResultFn_mock) Fn(sParam string, bParam bool) (result1 io.Reader) {
	params := MoqInterfaceResultFn_params{
		SParam: sParam,
		BParam: bParam,
	}
	var results *MoqInterfaceResultFn_results
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
		result.DoFn(sParam, bParam)
	}

	if result.Values != nil {
		result1 = result.Values.Result1
	}
	if result.DoReturnFn != nil {
		result1 = result.DoReturnFn(sParam, bParam)
	}
	return
}

func (m *MoqInterfaceResultFn) OnCall(sParam string, bParam bool) *MoqInterfaceResultFn_fnRecorder {
	return &MoqInterfaceResultFn_fnRecorder{
		Params: MoqInterfaceResultFn_params{
			SParam: sParam,
			BParam: bParam,
		},
		Sequence: m.Config.Sequence == moq.SeqDefaultOn,
		Moq:      m,
	}
}

func (r *MoqInterfaceResultFn_fnRecorder) Any() *MoqInterfaceResultFn_anyParams {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Any functions must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	return &MoqInterfaceResultFn_anyParams{Recorder: r}
}

func (a *MoqInterfaceResultFn_anyParams) SParam() *MoqInterfaceResultFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 0
	return a.Recorder
}

func (a *MoqInterfaceResultFn_anyParams) BParam() *MoqInterfaceResultFn_fnRecorder {
	a.Recorder.AnyParams |= 1 << 1
	return a.Recorder
}

func (r *MoqInterfaceResultFn_fnRecorder) Seq() *MoqInterfaceResultFn_fnRecorder {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Seq must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.Sequence = true
	return r
}

func (r *MoqInterfaceResultFn_fnRecorder) NoSeq() *MoqInterfaceResultFn_fnRecorder {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("NoSeq must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.Sequence = false
	return r
}

func (r *MoqInterfaceResultFn_fnRecorder) ReturnResults(result1 io.Reader) *MoqInterfaceResultFn_fnRecorder {
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values     *struct{ Result1 io.Reader }
		Sequence   uint32
		DoFn       MoqInterfaceResultFn_doFn
		DoReturnFn MoqInterfaceResultFn_doReturnFn
	}{
		Values: &struct{ Result1 io.Reader }{
			Result1: result1,
		},
		Sequence: sequence,
	})
	return r
}

func (r *MoqInterfaceResultFn_fnRecorder) AndDo(fn MoqInterfaceResultFn_doFn) *MoqInterfaceResultFn_fnRecorder {
	if r.Results == nil {
		r.Moq.Scene.T.Fatalf("ReturnResults must be called before calling AndDo")
		return nil
	}
	last := &r.Results.Results[len(r.Results.Results)-1]
	last.DoFn = fn
	return r
}

func (r *MoqInterfaceResultFn_fnRecorder) DoReturnResults(fn MoqInterfaceResultFn_doReturnFn) *MoqInterfaceResultFn_fnRecorder {
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values     *struct{ Result1 io.Reader }
		Sequence   uint32
		DoFn       MoqInterfaceResultFn_doFn
		DoReturnFn MoqInterfaceResultFn_doReturnFn
	}{Sequence: sequence, DoReturnFn: fn})
	return r
}

func (r *MoqInterfaceResultFn_fnRecorder) FindResults() {
	if r.Results != nil {
		r.Results.Repeat.Increment(r.Moq.Scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.AnyParams)
	insertAt := -1
	var results *MoqInterfaceResultFn_resultsByParams
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
		results = &MoqInterfaceResultFn_resultsByParams{
			AnyCount:  anyCount,
			AnyParams: r.AnyParams,
			Results:   map[MoqInterfaceResultFn_paramsKey]*MoqInterfaceResultFn_results{},
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
		r.Results = &MoqInterfaceResultFn_results{
			Params:  r.Params,
			Results: nil,
			Index:   0,
			Repeat:  &moq.RepeatVal{},
		}
		results.Results[paramsKey] = r.Results
	}

	r.Results.Repeat.Increment(r.Moq.Scene.T)
}

func (r *MoqInterfaceResultFn_fnRecorder) Repeat(repeaters ...moq.Repeater) *MoqInterfaceResultFn_fnRecorder {
	if r.Results == nil {
		r.Moq.Scene.T.Fatalf("ReturnResults or DoReturnResults must be called before calling Repeat")
		return nil
	}
	r.Results.Repeat.Repeat(r.Moq.Scene.T, repeaters)
	last := r.Results.Results[len(r.Results.Results)-1]
	for n := 0; n < r.Results.Repeat.ResultCount-1; n++ {
		if r.Sequence {
			last = struct {
				Values     *struct{ Result1 io.Reader }
				Sequence   uint32
				DoFn       MoqInterfaceResultFn_doFn
				DoReturnFn MoqInterfaceResultFn_doReturnFn
			}{
				Values: &struct{ Result1 io.Reader }{
					Result1: last.Values.Result1,
				},
				Sequence: r.Moq.Scene.NextRecorderSequence(),
			}
		}
		r.Results.Results = append(r.Results.Results, last)
	}
	return r
}

func (m *MoqInterfaceResultFn) ParamsKey(params MoqInterfaceResultFn_params, anyParams uint64) MoqInterfaceResultFn_paramsKey {
	var sParamUsed string
	var sParamUsedHash hash.Hash
	if anyParams&(1<<0) == 0 {
		if m.Runtime.ParameterIndexing.SParam == moq.ParamIndexByValue {
			sParamUsed = params.SParam
		} else {
			sParamUsedHash = hash.DeepHash(params.SParam)
		}
	}
	var bParamUsed bool
	var bParamUsedHash hash.Hash
	if anyParams&(1<<1) == 0 {
		if m.Runtime.ParameterIndexing.BParam == moq.ParamIndexByValue {
			bParamUsed = params.BParam
		} else {
			bParamUsedHash = hash.DeepHash(params.BParam)
		}
	}
	return MoqInterfaceResultFn_paramsKey{
		Params: struct {
			SParam string
			BParam bool
		}{
			SParam: sParamUsed,
			BParam: bParamUsed,
		},
		Hashes: struct {
			SParam hash.Hash
			BParam hash.Hash
		}{
			SParam: sParamUsedHash,
			BParam: bParamUsedHash,
		}}
}

// Reset resets the state of the moq
func (m *MoqInterfaceResultFn) Reset() { m.ResultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *MoqInterfaceResultFn) AssertExpectationsMet() {
	for _, res := range m.ResultsByParams {
		for _, results := range res.Results {
			missing := results.Repeat.MinTimes - int(atomic.LoadUint32(&results.Index))
			if missing > 0 {
				m.Scene.T.Errorf("Expected %d additional call(s) with parameters %#v", missing, results.Params)
			}
		}
	}
}
