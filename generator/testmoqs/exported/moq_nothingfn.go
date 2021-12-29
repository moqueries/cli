// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package exported

import (
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/generator/testmoqs"
	"github.com/myshkin5/moqueries/moq"
)

// MoqNothingFn holds the state of a moq of the NothingFn type
type MoqNothingFn struct {
	Scene           *moq.Scene
	Config          moq.Config
	ResultsByParams []MoqNothingFn_resultsByParams
}

// MoqNothingFn_mock isolates the mock interface of the NothingFn type
type MoqNothingFn_mock struct {
	Moq *MoqNothingFn
}

// MoqNothingFn_params holds the params of the NothingFn type
type MoqNothingFn_params struct{}

// MoqNothingFn_paramsKey holds the map key params of the NothingFn type
type MoqNothingFn_paramsKey struct{}

// MoqNothingFn_resultsByParams contains the results for a given set of parameters for the NothingFn type
type MoqNothingFn_resultsByParams struct {
	AnyCount  int
	AnyParams uint64
	Results   map[MoqNothingFn_paramsKey]*MoqNothingFn_results
}

// MoqNothingFn_doFn defines the type of function needed when calling AndDo for the NothingFn type
type MoqNothingFn_doFn func()

// MoqNothingFn_doReturnFn defines the type of function needed when calling DoReturnResults for the NothingFn type
type MoqNothingFn_doReturnFn func()

// MoqNothingFn_results holds the results of the NothingFn type
type MoqNothingFn_results struct {
	Params  MoqNothingFn_params
	Results []struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqNothingFn_doFn
		DoReturnFn MoqNothingFn_doReturnFn
	}
	Index  uint32
	Repeat *moq.RepeatVal
}

// MoqNothingFn_fnRecorder routes recorded function calls to the MoqNothingFn moq
type MoqNothingFn_fnRecorder struct {
	Params    MoqNothingFn_params
	ParamsKey MoqNothingFn_paramsKey
	AnyParams uint64
	Sequence  bool
	Results   *MoqNothingFn_results
	Moq       *MoqNothingFn
}

// MoqNothingFn_anyParams isolates the any params functions of the NothingFn type
type MoqNothingFn_anyParams struct {
	Recorder *MoqNothingFn_fnRecorder
}

// NewMoqNothingFn creates a new moq of the NothingFn type
func NewMoqNothingFn(scene *moq.Scene, config *moq.Config) *MoqNothingFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &MoqNothingFn{
		Scene:  scene,
		Config: *config,
	}
	scene.AddMoq(m)
	return m
}

// Mock returns the moq implementation of the NothingFn type
func (m *MoqNothingFn) Mock() testmoqs.NothingFn {
	return func() { moq := &MoqNothingFn_mock{Moq: m}; moq.Fn() }
}

func (m *MoqNothingFn_mock) Fn() {
	params := MoqNothingFn_params{}
	var results *MoqNothingFn_results
	for _, resultsByParams := range m.Moq.ResultsByParams {
		paramsKey := MoqNothingFn_paramsKey{}
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
		result.DoFn()
	}

	if result.DoReturnFn != nil {
		result.DoReturnFn()
	}
	return
}

func (m *MoqNothingFn) OnCall() *MoqNothingFn_fnRecorder {
	return &MoqNothingFn_fnRecorder{
		Params:    MoqNothingFn_params{},
		ParamsKey: MoqNothingFn_paramsKey{},
		Sequence:  m.Config.Sequence == moq.SeqDefaultOn,
		Moq:       m,
	}
}

func (r *MoqNothingFn_fnRecorder) Any() *MoqNothingFn_anyParams {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Any functions must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	return &MoqNothingFn_anyParams{Recorder: r}
}

func (r *MoqNothingFn_fnRecorder) Seq() *MoqNothingFn_fnRecorder {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Seq must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.Sequence = true
	return r
}

func (r *MoqNothingFn_fnRecorder) NoSeq() *MoqNothingFn_fnRecorder {
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("NoSeq must be called before ReturnResults or DoReturnResults calls, parameters: %#v", r.Params)
		return nil
	}
	r.Sequence = false
	return r
}

func (r *MoqNothingFn_fnRecorder) ReturnResults() *MoqNothingFn_fnRecorder {
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqNothingFn_doFn
		DoReturnFn MoqNothingFn_doReturnFn
	}{
		Values:   &struct{}{},
		Sequence: sequence,
	})
	return r
}

func (r *MoqNothingFn_fnRecorder) AndDo(fn MoqNothingFn_doFn) *MoqNothingFn_fnRecorder {
	if r.Results == nil {
		r.Moq.Scene.T.Fatalf("ReturnResults must be called before calling AndDo")
		return nil
	}
	last := &r.Results.Results[len(r.Results.Results)-1]
	last.DoFn = fn
	return r
}

func (r *MoqNothingFn_fnRecorder) DoReturnResults(fn MoqNothingFn_doReturnFn) *MoqNothingFn_fnRecorder {
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqNothingFn_doFn
		DoReturnFn MoqNothingFn_doReturnFn
	}{Sequence: sequence, DoReturnFn: fn})
	return r
}

func (r *MoqNothingFn_fnRecorder) FindResults() {
	if r.Results == nil {
		anyCount := bits.OnesCount64(r.AnyParams)
		insertAt := -1
		var results *MoqNothingFn_resultsByParams
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
			results = &MoqNothingFn_resultsByParams{
				AnyCount:  anyCount,
				AnyParams: r.AnyParams,
				Results:   map[MoqNothingFn_paramsKey]*MoqNothingFn_results{},
			}
			r.Moq.ResultsByParams = append(r.Moq.ResultsByParams, *results)
			if insertAt != -1 && insertAt+1 < len(r.Moq.ResultsByParams) {
				copy(r.Moq.ResultsByParams[insertAt+1:], r.Moq.ResultsByParams[insertAt:0])
				r.Moq.ResultsByParams[insertAt] = *results
			}
		}

		paramsKey := MoqNothingFn_paramsKey{}

		var ok bool
		r.Results, ok = results.Results[paramsKey]
		if !ok {
			r.Results = &MoqNothingFn_results{
				Params:  r.Params,
				Results: nil,
				Index:   0,
				Repeat:  &moq.RepeatVal{},
			}
			results.Results[paramsKey] = r.Results
		}
	}
	r.Results.Repeat.Increment(r.Moq.Scene.T)
}

func (r *MoqNothingFn_fnRecorder) Repeat(repeaters ...moq.Repeater) *MoqNothingFn_fnRecorder {
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
				DoFn       MoqNothingFn_doFn
				DoReturnFn MoqNothingFn_doReturnFn
			}{
				Values:   &struct{}{},
				Sequence: r.Moq.Scene.NextRecorderSequence(),
			}
		}
		r.Results.Results = append(r.Results.Results, last)
	}
	return r
}

// Reset resets the state of the moq
func (m *MoqNothingFn) Reset() { m.ResultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *MoqNothingFn) AssertExpectationsMet() {
	for _, res := range m.ResultsByParams {
		for _, results := range res.Results {
			missing := results.Repeat.MinTimes - int(atomic.LoadUint32(&results.Index))
			if missing > 0 {
				m.Scene.T.Errorf("Expected %d additional call(s) with parameters %#v", missing, results.Params)
			}
		}
	}
}