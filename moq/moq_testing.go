// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package moq

import (
	"fmt"
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/hash"
)

// MoqT holds the state of a moq of the T type
type MoqT struct {
	Scene  *Scene
	Config Config
	Moq    *MoqT_mock

	ResultsByParams_Errorf []MoqT_Errorf_resultsByParams
	ResultsByParams_Fatalf []MoqT_Fatalf_resultsByParams
	ResultsByParams_Helper []MoqT_Helper_resultsByParams

	Runtime struct {
		ParameterIndexing struct {
			Errorf struct {
				Format ParamIndexing
				Args   ParamIndexing
			}
			Fatalf struct {
				Format ParamIndexing
				Args   ParamIndexing
			}
			Helper struct{}
		}
	}
}

// MoqT_mock isolates the mock interface of the T type
type MoqT_mock struct {
	Moq *MoqT
}

// MoqT_recorder isolates the recorder interface of the T type
type MoqT_recorder struct {
	Moq *MoqT
}

// MoqT_Errorf_params holds the params of the T type
type MoqT_Errorf_params struct {
	Format string
	Args   []interface{}
}

// MoqT_Errorf_paramsKey holds the map key params of the T type
type MoqT_Errorf_paramsKey struct {
	Params struct{ Format string }
	Hashes struct {
		Format hash.Hash
		Args   hash.Hash
	}
}

// MoqT_Errorf_resultsByParams contains the results for a given set of parameters for the T type
type MoqT_Errorf_resultsByParams struct {
	AnyCount  int
	AnyParams uint64
	Results   map[MoqT_Errorf_paramsKey]*MoqT_Errorf_results
}

// MoqT_Errorf_doFn defines the type of function needed when calling AndDo for the T type
type MoqT_Errorf_doFn func(format string, args ...interface{})

// MoqT_Errorf_doReturnFn defines the type of function needed when calling DoReturnResults for the T type
type MoqT_Errorf_doReturnFn func(format string, args ...interface{})

// MoqT_Errorf_results holds the results of the T type
type MoqT_Errorf_results struct {
	Params  MoqT_Errorf_params
	Results []struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqT_Errorf_doFn
		DoReturnFn MoqT_Errorf_doReturnFn
	}
	Index  uint32
	Repeat *RepeatVal
}

// MoqT_Errorf_fnRecorder routes recorded function calls to the MoqT moq
type MoqT_Errorf_fnRecorder struct {
	Params    MoqT_Errorf_params
	AnyParams uint64
	Sequence  bool
	Results   *MoqT_Errorf_results
	Moq       *MoqT
}

// MoqT_Errorf_anyParams isolates the any params functions of the T type
type MoqT_Errorf_anyParams struct {
	Recorder *MoqT_Errorf_fnRecorder
}

// MoqT_Fatalf_params holds the params of the T type
type MoqT_Fatalf_params struct {
	Format string
	Args   []interface{}
}

// MoqT_Fatalf_paramsKey holds the map key params of the T type
type MoqT_Fatalf_paramsKey struct {
	Params struct{ Format string }
	Hashes struct {
		Format hash.Hash
		Args   hash.Hash
	}
}

// MoqT_Fatalf_resultsByParams contains the results for a given set of parameters for the T type
type MoqT_Fatalf_resultsByParams struct {
	AnyCount  int
	AnyParams uint64
	Results   map[MoqT_Fatalf_paramsKey]*MoqT_Fatalf_results
}

// MoqT_Fatalf_doFn defines the type of function needed when calling AndDo for the T type
type MoqT_Fatalf_doFn func(format string, args ...interface{})

// MoqT_Fatalf_doReturnFn defines the type of function needed when calling DoReturnResults for the T type
type MoqT_Fatalf_doReturnFn func(format string, args ...interface{})

// MoqT_Fatalf_results holds the results of the T type
type MoqT_Fatalf_results struct {
	Params  MoqT_Fatalf_params
	Results []struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqT_Fatalf_doFn
		DoReturnFn MoqT_Fatalf_doReturnFn
	}
	Index  uint32
	Repeat *RepeatVal
}

// MoqT_Fatalf_fnRecorder routes recorded function calls to the MoqT moq
type MoqT_Fatalf_fnRecorder struct {
	Params    MoqT_Fatalf_params
	AnyParams uint64
	Sequence  bool
	Results   *MoqT_Fatalf_results
	Moq       *MoqT
}

// MoqT_Fatalf_anyParams isolates the any params functions of the T type
type MoqT_Fatalf_anyParams struct {
	Recorder *MoqT_Fatalf_fnRecorder
}

// MoqT_Helper_params holds the params of the T type
type MoqT_Helper_params struct{}

// MoqT_Helper_paramsKey holds the map key params of the T type
type MoqT_Helper_paramsKey struct {
	Params struct{}
	Hashes struct{}
}

// MoqT_Helper_resultsByParams contains the results for a given set of parameters for the T type
type MoqT_Helper_resultsByParams struct {
	AnyCount  int
	AnyParams uint64
	Results   map[MoqT_Helper_paramsKey]*MoqT_Helper_results
}

// MoqT_Helper_doFn defines the type of function needed when calling AndDo for the T type
type MoqT_Helper_doFn func()

// MoqT_Helper_doReturnFn defines the type of function needed when calling DoReturnResults for the T type
type MoqT_Helper_doReturnFn func()

// MoqT_Helper_results holds the results of the T type
type MoqT_Helper_results struct {
	Params  MoqT_Helper_params
	Results []struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqT_Helper_doFn
		DoReturnFn MoqT_Helper_doReturnFn
	}
	Index  uint32
	Repeat *RepeatVal
}

// MoqT_Helper_fnRecorder routes recorded function calls to the MoqT moq
type MoqT_Helper_fnRecorder struct {
	Params    MoqT_Helper_params
	AnyParams uint64
	Sequence  bool
	Results   *MoqT_Helper_results
	Moq       *MoqT
}

// MoqT_Helper_anyParams isolates the any params functions of the T type
type MoqT_Helper_anyParams struct {
	Recorder *MoqT_Helper_fnRecorder
}

// NewMoqT creates a new moq of the T type
func NewMoqT(scene *Scene, config *Config) *MoqT {
	if config == nil {
		config = &Config{}
	}
	m := &MoqT{
		Scene:  scene,
		Config: *config,
		Moq:    &MoqT_mock{},

		Runtime: struct {
			ParameterIndexing struct {
				Errorf struct {
					Format ParamIndexing
					Args   ParamIndexing
				}
				Fatalf struct {
					Format ParamIndexing
					Args   ParamIndexing
				}
				Helper struct{}
			}
		}{ParameterIndexing: struct {
			Errorf struct {
				Format ParamIndexing
				Args   ParamIndexing
			}
			Fatalf struct {
				Format ParamIndexing
				Args   ParamIndexing
			}
			Helper struct{}
		}{
			Errorf: struct {
				Format ParamIndexing
				Args   ParamIndexing
			}{
				Format: ParamIndexByValue,
				Args:   ParamIndexByHash,
			},
			Fatalf: struct {
				Format ParamIndexing
				Args   ParamIndexing
			}{
				Format: ParamIndexByValue,
				Args:   ParamIndexByHash,
			},
			Helper: struct{}{},
		}},
	}
	m.Moq.Moq = m

	scene.AddMoq(m)
	return m
}

// Mock returns the mock implementation of the T type
func (m *MoqT) Mock() *MoqT_mock { return m.Moq }

func (m *MoqT_mock) Errorf(format string, args ...interface{}) {
	m.Moq.Scene.T.Helper()
	params := MoqT_Errorf_params{
		Format: format,
		Args:   args,
	}
	var results *MoqT_Errorf_results
	for _, resultsByParams := range m.Moq.ResultsByParams_Errorf {
		paramsKey := m.Moq.ParamsKey_Errorf(params, resultsByParams.AnyParams)
		var ok bool
		results, ok = resultsByParams.Results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.Moq.Config.Expectation == Strict {
			m.Moq.Scene.T.Fatalf("Unexpected call to %s", m.Moq.PrettyParams_Errorf(params))
		}
		return
	}

	i := int(atomic.AddUint32(&results.Index, 1)) - 1
	if i >= results.Repeat.ResultCount {
		if !results.Repeat.AnyTimes {
			if m.Moq.Config.Expectation == Strict {
				m.Moq.Scene.T.Fatalf("Too many calls to %s", m.Moq.PrettyParams_Errorf(params))
			}
			return
		}
		i = results.Repeat.ResultCount - 1
	}

	result := results.Results[i]
	if result.Sequence != 0 {
		sequence := m.Moq.Scene.NextMockSequence()
		if (!results.Repeat.AnyTimes && result.Sequence != sequence) || result.Sequence > sequence {
			m.Moq.Scene.T.Fatalf("Call sequence does not match call to %s", m.Moq.PrettyParams_Errorf(params))
		}
	}

	if result.DoFn != nil {
		result.DoFn(format, args...)
	}

	if result.DoReturnFn != nil {
		result.DoReturnFn(format, args...)
	}
	return
}

func (m *MoqT_mock) Fatalf(format string, args ...interface{}) {
	m.Moq.Scene.T.Helper()
	params := MoqT_Fatalf_params{
		Format: format,
		Args:   args,
	}
	var results *MoqT_Fatalf_results
	for _, resultsByParams := range m.Moq.ResultsByParams_Fatalf {
		paramsKey := m.Moq.ParamsKey_Fatalf(params, resultsByParams.AnyParams)
		var ok bool
		results, ok = resultsByParams.Results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.Moq.Config.Expectation == Strict {
			m.Moq.Scene.T.Fatalf("Unexpected call to %s", m.Moq.PrettyParams_Fatalf(params))
		}
		return
	}

	i := int(atomic.AddUint32(&results.Index, 1)) - 1
	if i >= results.Repeat.ResultCount {
		if !results.Repeat.AnyTimes {
			if m.Moq.Config.Expectation == Strict {
				m.Moq.Scene.T.Fatalf("Too many calls to %s", m.Moq.PrettyParams_Fatalf(params))
			}
			return
		}
		i = results.Repeat.ResultCount - 1
	}

	result := results.Results[i]
	if result.Sequence != 0 {
		sequence := m.Moq.Scene.NextMockSequence()
		if (!results.Repeat.AnyTimes && result.Sequence != sequence) || result.Sequence > sequence {
			m.Moq.Scene.T.Fatalf("Call sequence does not match call to %s", m.Moq.PrettyParams_Fatalf(params))
		}
	}

	if result.DoFn != nil {
		result.DoFn(format, args...)
	}

	if result.DoReturnFn != nil {
		result.DoReturnFn(format, args...)
	}
	return
}

func (m *MoqT_mock) Helper() {
	m.Moq.Scene.T.Helper()
	params := MoqT_Helper_params{}
	var results *MoqT_Helper_results
	for _, resultsByParams := range m.Moq.ResultsByParams_Helper {
		paramsKey := m.Moq.ParamsKey_Helper(params, resultsByParams.AnyParams)
		var ok bool
		results, ok = resultsByParams.Results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.Moq.Config.Expectation == Strict {
			m.Moq.Scene.T.Fatalf("Unexpected call to %s", m.Moq.PrettyParams_Helper(params))
		}
		return
	}

	i := int(atomic.AddUint32(&results.Index, 1)) - 1
	if i >= results.Repeat.ResultCount {
		if !results.Repeat.AnyTimes {
			if m.Moq.Config.Expectation == Strict {
				m.Moq.Scene.T.Fatalf("Too many calls to %s", m.Moq.PrettyParams_Helper(params))
			}
			return
		}
		i = results.Repeat.ResultCount - 1
	}

	result := results.Results[i]
	if result.Sequence != 0 {
		sequence := m.Moq.Scene.NextMockSequence()
		if (!results.Repeat.AnyTimes && result.Sequence != sequence) || result.Sequence > sequence {
			m.Moq.Scene.T.Fatalf("Call sequence does not match call to %s", m.Moq.PrettyParams_Helper(params))
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

// OnCall returns the recorder implementation of the T type
func (m *MoqT) OnCall() *MoqT_recorder {
	return &MoqT_recorder{
		Moq: m,
	}
}

func (m *MoqT_recorder) Errorf(format string, args ...interface{}) *MoqT_Errorf_fnRecorder {
	return &MoqT_Errorf_fnRecorder{
		Params: MoqT_Errorf_params{
			Format: format,
			Args:   args,
		},
		Sequence: m.Moq.Config.Sequence == SeqDefaultOn,
		Moq:      m.Moq,
	}
}

func (r *MoqT_Errorf_fnRecorder) Any() *MoqT_Errorf_anyParams {
	r.Moq.Scene.T.Helper()
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Any functions must be called before ReturnResults or DoReturnResults calls, recording %s", r.Moq.PrettyParams_Errorf(r.Params))
		return nil
	}
	return &MoqT_Errorf_anyParams{Recorder: r}
}

func (a *MoqT_Errorf_anyParams) Format() *MoqT_Errorf_fnRecorder {
	a.Recorder.AnyParams |= 1 << 0
	return a.Recorder
}

func (a *MoqT_Errorf_anyParams) Args() *MoqT_Errorf_fnRecorder {
	a.Recorder.AnyParams |= 1 << 1
	return a.Recorder
}

func (r *MoqT_Errorf_fnRecorder) Seq() *MoqT_Errorf_fnRecorder {
	r.Moq.Scene.T.Helper()
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Seq must be called before ReturnResults or DoReturnResults calls, recording %s", r.Moq.PrettyParams_Errorf(r.Params))
		return nil
	}
	r.Sequence = true
	return r
}

func (r *MoqT_Errorf_fnRecorder) NoSeq() *MoqT_Errorf_fnRecorder {
	r.Moq.Scene.T.Helper()
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("NoSeq must be called before ReturnResults or DoReturnResults calls, recording %s", r.Moq.PrettyParams_Errorf(r.Params))
		return nil
	}
	r.Sequence = false
	return r
}

func (r *MoqT_Errorf_fnRecorder) ReturnResults() *MoqT_Errorf_fnRecorder {
	r.Moq.Scene.T.Helper()
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqT_Errorf_doFn
		DoReturnFn MoqT_Errorf_doReturnFn
	}{
		Values:   &struct{}{},
		Sequence: sequence,
	})
	return r
}

func (r *MoqT_Errorf_fnRecorder) AndDo(fn MoqT_Errorf_doFn) *MoqT_Errorf_fnRecorder {
	r.Moq.Scene.T.Helper()
	if r.Results == nil {
		r.Moq.Scene.T.Fatalf("ReturnResults must be called before calling AndDo")
		return nil
	}
	last := &r.Results.Results[len(r.Results.Results)-1]
	last.DoFn = fn
	return r
}

func (r *MoqT_Errorf_fnRecorder) DoReturnResults(fn MoqT_Errorf_doReturnFn) *MoqT_Errorf_fnRecorder {
	r.Moq.Scene.T.Helper()
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqT_Errorf_doFn
		DoReturnFn MoqT_Errorf_doReturnFn
	}{Sequence: sequence, DoReturnFn: fn})
	return r
}

func (r *MoqT_Errorf_fnRecorder) FindResults() {
	if r.Results != nil {
		r.Results.Repeat.Increment(r.Moq.Scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.AnyParams)
	insertAt := -1
	var results *MoqT_Errorf_resultsByParams
	for n, res := range r.Moq.ResultsByParams_Errorf {
		if res.AnyParams == r.AnyParams {
			results = &res
			break
		}
		if res.AnyCount > anyCount {
			insertAt = n
		}
	}
	if results == nil {
		results = &MoqT_Errorf_resultsByParams{
			AnyCount:  anyCount,
			AnyParams: r.AnyParams,
			Results:   map[MoqT_Errorf_paramsKey]*MoqT_Errorf_results{},
		}
		r.Moq.ResultsByParams_Errorf = append(r.Moq.ResultsByParams_Errorf, *results)
		if insertAt != -1 && insertAt+1 < len(r.Moq.ResultsByParams_Errorf) {
			copy(r.Moq.ResultsByParams_Errorf[insertAt+1:], r.Moq.ResultsByParams_Errorf[insertAt:0])
			r.Moq.ResultsByParams_Errorf[insertAt] = *results
		}
	}

	paramsKey := r.Moq.ParamsKey_Errorf(r.Params, r.AnyParams)

	var ok bool
	r.Results, ok = results.Results[paramsKey]
	if !ok {
		r.Results = &MoqT_Errorf_results{
			Params:  r.Params,
			Results: nil,
			Index:   0,
			Repeat:  &RepeatVal{},
		}
		results.Results[paramsKey] = r.Results
	}

	r.Results.Repeat.Increment(r.Moq.Scene.T)
}

func (r *MoqT_Errorf_fnRecorder) Repeat(repeaters ...Repeater) *MoqT_Errorf_fnRecorder {
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
				Values     *struct{}
				Sequence   uint32
				DoFn       MoqT_Errorf_doFn
				DoReturnFn MoqT_Errorf_doReturnFn
			}{
				Values:   last.Values,
				Sequence: r.Moq.Scene.NextRecorderSequence(),
			}
		}
		r.Results.Results = append(r.Results.Results, last)
	}
	return r
}

func (m *MoqT) PrettyParams_Errorf(params MoqT_Errorf_params) string {
	return fmt.Sprintf("Errorf(%#v, %#v)", params.Format, params.Args)
}

func (m *MoqT) ParamsKey_Errorf(params MoqT_Errorf_params, anyParams uint64) MoqT_Errorf_paramsKey {
	var formatUsed string
	var formatUsedHash hash.Hash
	if anyParams&(1<<0) == 0 {
		if m.Runtime.ParameterIndexing.Errorf.Format == ParamIndexByValue {
			formatUsed = params.Format
		} else {
			formatUsedHash = hash.DeepHash(params.Format)
		}
	}
	var argsUsedHash hash.Hash
	if anyParams&(1<<1) == 0 {
		if m.Runtime.ParameterIndexing.Errorf.Args == ParamIndexByValue {
			m.Scene.T.Fatalf("The args parameter of the Errorf function can't be indexed by value")
		}
		argsUsedHash = hash.DeepHash(params.Args)
	}
	return MoqT_Errorf_paramsKey{
		Params: struct{ Format string }{
			Format: formatUsed,
		},
		Hashes: struct {
			Format hash.Hash
			Args   hash.Hash
		}{
			Format: formatUsedHash,
			Args:   argsUsedHash,
		},
	}
}

func (m *MoqT_recorder) Fatalf(format string, args ...interface{}) *MoqT_Fatalf_fnRecorder {
	return &MoqT_Fatalf_fnRecorder{
		Params: MoqT_Fatalf_params{
			Format: format,
			Args:   args,
		},
		Sequence: m.Moq.Config.Sequence == SeqDefaultOn,
		Moq:      m.Moq,
	}
}

func (r *MoqT_Fatalf_fnRecorder) Any() *MoqT_Fatalf_anyParams {
	r.Moq.Scene.T.Helper()
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Any functions must be called before ReturnResults or DoReturnResults calls, recording %s", r.Moq.PrettyParams_Fatalf(r.Params))
		return nil
	}
	return &MoqT_Fatalf_anyParams{Recorder: r}
}

func (a *MoqT_Fatalf_anyParams) Format() *MoqT_Fatalf_fnRecorder {
	a.Recorder.AnyParams |= 1 << 0
	return a.Recorder
}

func (a *MoqT_Fatalf_anyParams) Args() *MoqT_Fatalf_fnRecorder {
	a.Recorder.AnyParams |= 1 << 1
	return a.Recorder
}

func (r *MoqT_Fatalf_fnRecorder) Seq() *MoqT_Fatalf_fnRecorder {
	r.Moq.Scene.T.Helper()
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Seq must be called before ReturnResults or DoReturnResults calls, recording %s", r.Moq.PrettyParams_Fatalf(r.Params))
		return nil
	}
	r.Sequence = true
	return r
}

func (r *MoqT_Fatalf_fnRecorder) NoSeq() *MoqT_Fatalf_fnRecorder {
	r.Moq.Scene.T.Helper()
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("NoSeq must be called before ReturnResults or DoReturnResults calls, recording %s", r.Moq.PrettyParams_Fatalf(r.Params))
		return nil
	}
	r.Sequence = false
	return r
}

func (r *MoqT_Fatalf_fnRecorder) ReturnResults() *MoqT_Fatalf_fnRecorder {
	r.Moq.Scene.T.Helper()
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqT_Fatalf_doFn
		DoReturnFn MoqT_Fatalf_doReturnFn
	}{
		Values:   &struct{}{},
		Sequence: sequence,
	})
	return r
}

func (r *MoqT_Fatalf_fnRecorder) AndDo(fn MoqT_Fatalf_doFn) *MoqT_Fatalf_fnRecorder {
	r.Moq.Scene.T.Helper()
	if r.Results == nil {
		r.Moq.Scene.T.Fatalf("ReturnResults must be called before calling AndDo")
		return nil
	}
	last := &r.Results.Results[len(r.Results.Results)-1]
	last.DoFn = fn
	return r
}

func (r *MoqT_Fatalf_fnRecorder) DoReturnResults(fn MoqT_Fatalf_doReturnFn) *MoqT_Fatalf_fnRecorder {
	r.Moq.Scene.T.Helper()
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqT_Fatalf_doFn
		DoReturnFn MoqT_Fatalf_doReturnFn
	}{Sequence: sequence, DoReturnFn: fn})
	return r
}

func (r *MoqT_Fatalf_fnRecorder) FindResults() {
	if r.Results != nil {
		r.Results.Repeat.Increment(r.Moq.Scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.AnyParams)
	insertAt := -1
	var results *MoqT_Fatalf_resultsByParams
	for n, res := range r.Moq.ResultsByParams_Fatalf {
		if res.AnyParams == r.AnyParams {
			results = &res
			break
		}
		if res.AnyCount > anyCount {
			insertAt = n
		}
	}
	if results == nil {
		results = &MoqT_Fatalf_resultsByParams{
			AnyCount:  anyCount,
			AnyParams: r.AnyParams,
			Results:   map[MoqT_Fatalf_paramsKey]*MoqT_Fatalf_results{},
		}
		r.Moq.ResultsByParams_Fatalf = append(r.Moq.ResultsByParams_Fatalf, *results)
		if insertAt != -1 && insertAt+1 < len(r.Moq.ResultsByParams_Fatalf) {
			copy(r.Moq.ResultsByParams_Fatalf[insertAt+1:], r.Moq.ResultsByParams_Fatalf[insertAt:0])
			r.Moq.ResultsByParams_Fatalf[insertAt] = *results
		}
	}

	paramsKey := r.Moq.ParamsKey_Fatalf(r.Params, r.AnyParams)

	var ok bool
	r.Results, ok = results.Results[paramsKey]
	if !ok {
		r.Results = &MoqT_Fatalf_results{
			Params:  r.Params,
			Results: nil,
			Index:   0,
			Repeat:  &RepeatVal{},
		}
		results.Results[paramsKey] = r.Results
	}

	r.Results.Repeat.Increment(r.Moq.Scene.T)
}

func (r *MoqT_Fatalf_fnRecorder) Repeat(repeaters ...Repeater) *MoqT_Fatalf_fnRecorder {
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
				Values     *struct{}
				Sequence   uint32
				DoFn       MoqT_Fatalf_doFn
				DoReturnFn MoqT_Fatalf_doReturnFn
			}{
				Values:   last.Values,
				Sequence: r.Moq.Scene.NextRecorderSequence(),
			}
		}
		r.Results.Results = append(r.Results.Results, last)
	}
	return r
}

func (m *MoqT) PrettyParams_Fatalf(params MoqT_Fatalf_params) string {
	return fmt.Sprintf("Fatalf(%#v, %#v)", params.Format, params.Args)
}

func (m *MoqT) ParamsKey_Fatalf(params MoqT_Fatalf_params, anyParams uint64) MoqT_Fatalf_paramsKey {
	var formatUsed string
	var formatUsedHash hash.Hash
	if anyParams&(1<<0) == 0 {
		if m.Runtime.ParameterIndexing.Fatalf.Format == ParamIndexByValue {
			formatUsed = params.Format
		} else {
			formatUsedHash = hash.DeepHash(params.Format)
		}
	}
	var argsUsedHash hash.Hash
	if anyParams&(1<<1) == 0 {
		if m.Runtime.ParameterIndexing.Fatalf.Args == ParamIndexByValue {
			m.Scene.T.Fatalf("The args parameter of the Fatalf function can't be indexed by value")
		}
		argsUsedHash = hash.DeepHash(params.Args)
	}
	return MoqT_Fatalf_paramsKey{
		Params: struct{ Format string }{
			Format: formatUsed,
		},
		Hashes: struct {
			Format hash.Hash
			Args   hash.Hash
		}{
			Format: formatUsedHash,
			Args:   argsUsedHash,
		},
	}
}

func (m *MoqT_recorder) Helper() *MoqT_Helper_fnRecorder {
	return &MoqT_Helper_fnRecorder{
		Params:   MoqT_Helper_params{},
		Sequence: m.Moq.Config.Sequence == SeqDefaultOn,
		Moq:      m.Moq,
	}
}

func (r *MoqT_Helper_fnRecorder) Any() *MoqT_Helper_anyParams {
	r.Moq.Scene.T.Helper()
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Any functions must be called before ReturnResults or DoReturnResults calls, recording %s", r.Moq.PrettyParams_Helper(r.Params))
		return nil
	}
	return &MoqT_Helper_anyParams{Recorder: r}
}

func (r *MoqT_Helper_fnRecorder) Seq() *MoqT_Helper_fnRecorder {
	r.Moq.Scene.T.Helper()
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("Seq must be called before ReturnResults or DoReturnResults calls, recording %s", r.Moq.PrettyParams_Helper(r.Params))
		return nil
	}
	r.Sequence = true
	return r
}

func (r *MoqT_Helper_fnRecorder) NoSeq() *MoqT_Helper_fnRecorder {
	r.Moq.Scene.T.Helper()
	if r.Results != nil {
		r.Moq.Scene.T.Fatalf("NoSeq must be called before ReturnResults or DoReturnResults calls, recording %s", r.Moq.PrettyParams_Helper(r.Params))
		return nil
	}
	r.Sequence = false
	return r
}

func (r *MoqT_Helper_fnRecorder) ReturnResults() *MoqT_Helper_fnRecorder {
	r.Moq.Scene.T.Helper()
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqT_Helper_doFn
		DoReturnFn MoqT_Helper_doReturnFn
	}{
		Values:   &struct{}{},
		Sequence: sequence,
	})
	return r
}

func (r *MoqT_Helper_fnRecorder) AndDo(fn MoqT_Helper_doFn) *MoqT_Helper_fnRecorder {
	r.Moq.Scene.T.Helper()
	if r.Results == nil {
		r.Moq.Scene.T.Fatalf("ReturnResults must be called before calling AndDo")
		return nil
	}
	last := &r.Results.Results[len(r.Results.Results)-1]
	last.DoFn = fn
	return r
}

func (r *MoqT_Helper_fnRecorder) DoReturnResults(fn MoqT_Helper_doReturnFn) *MoqT_Helper_fnRecorder {
	r.Moq.Scene.T.Helper()
	r.FindResults()

	var sequence uint32
	if r.Sequence {
		sequence = r.Moq.Scene.NextRecorderSequence()
	}

	r.Results.Results = append(r.Results.Results, struct {
		Values     *struct{}
		Sequence   uint32
		DoFn       MoqT_Helper_doFn
		DoReturnFn MoqT_Helper_doReturnFn
	}{Sequence: sequence, DoReturnFn: fn})
	return r
}

func (r *MoqT_Helper_fnRecorder) FindResults() {
	if r.Results != nil {
		r.Results.Repeat.Increment(r.Moq.Scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.AnyParams)
	insertAt := -1
	var results *MoqT_Helper_resultsByParams
	for n, res := range r.Moq.ResultsByParams_Helper {
		if res.AnyParams == r.AnyParams {
			results = &res
			break
		}
		if res.AnyCount > anyCount {
			insertAt = n
		}
	}
	if results == nil {
		results = &MoqT_Helper_resultsByParams{
			AnyCount:  anyCount,
			AnyParams: r.AnyParams,
			Results:   map[MoqT_Helper_paramsKey]*MoqT_Helper_results{},
		}
		r.Moq.ResultsByParams_Helper = append(r.Moq.ResultsByParams_Helper, *results)
		if insertAt != -1 && insertAt+1 < len(r.Moq.ResultsByParams_Helper) {
			copy(r.Moq.ResultsByParams_Helper[insertAt+1:], r.Moq.ResultsByParams_Helper[insertAt:0])
			r.Moq.ResultsByParams_Helper[insertAt] = *results
		}
	}

	paramsKey := r.Moq.ParamsKey_Helper(r.Params, r.AnyParams)

	var ok bool
	r.Results, ok = results.Results[paramsKey]
	if !ok {
		r.Results = &MoqT_Helper_results{
			Params:  r.Params,
			Results: nil,
			Index:   0,
			Repeat:  &RepeatVal{},
		}
		results.Results[paramsKey] = r.Results
	}

	r.Results.Repeat.Increment(r.Moq.Scene.T)
}

func (r *MoqT_Helper_fnRecorder) Repeat(repeaters ...Repeater) *MoqT_Helper_fnRecorder {
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
				Values     *struct{}
				Sequence   uint32
				DoFn       MoqT_Helper_doFn
				DoReturnFn MoqT_Helper_doReturnFn
			}{
				Values:   last.Values,
				Sequence: r.Moq.Scene.NextRecorderSequence(),
			}
		}
		r.Results.Results = append(r.Results.Results, last)
	}
	return r
}

func (m *MoqT) PrettyParams_Helper(params MoqT_Helper_params) string { return fmt.Sprintf("Helper()") }

func (m *MoqT) ParamsKey_Helper(params MoqT_Helper_params, anyParams uint64) MoqT_Helper_paramsKey {
	return MoqT_Helper_paramsKey{
		Params: struct{}{},
		Hashes: struct{}{},
	}
}

// Reset resets the state of the moq
func (m *MoqT) Reset() {
	m.ResultsByParams_Errorf = nil
	m.ResultsByParams_Fatalf = nil
	m.ResultsByParams_Helper = nil
}

// AssertExpectationsMet asserts that all expectations have been met
func (m *MoqT) AssertExpectationsMet() {
	m.Scene.T.Helper()
	for _, res := range m.ResultsByParams_Errorf {
		for _, results := range res.Results {
			missing := results.Repeat.MinTimes - int(atomic.LoadUint32(&results.Index))
			if missing > 0 {
				m.Scene.T.Errorf("Expected %d additional call(s) to %s", missing, m.PrettyParams_Errorf(results.Params))
			}
		}
	}
	for _, res := range m.ResultsByParams_Fatalf {
		for _, results := range res.Results {
			missing := results.Repeat.MinTimes - int(atomic.LoadUint32(&results.Index))
			if missing > 0 {
				m.Scene.T.Errorf("Expected %d additional call(s) to %s", missing, m.PrettyParams_Fatalf(results.Params))
			}
		}
	}
	for _, res := range m.ResultsByParams_Helper {
		for _, results := range res.Results {
			missing := results.Repeat.MinTimes - int(atomic.LoadUint32(&results.Index))
			if missing > 0 {
				m.Scene.T.Errorf("Expected %d additional call(s) to %s", missing, m.PrettyParams_Helper(results.Params))
			}
		}
	}
}
