// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package generator_test

import (
	"fmt"
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/generator"
	"github.com/myshkin5/moqueries/hash"
	"github.com/myshkin5/moqueries/moq"
)

// moqNewConverterFunc holds the state of a moq of the NewConverterFunc type
type moqNewConverterFunc struct {
	scene  *moq.Scene
	config moq.Config
	moq    *moqNewConverterFunc_mock

	resultsByParams []moqNewConverterFunc_resultsByParams

	runtime struct {
		parameterIndexing struct {
			typ    moq.ParamIndexing
			export moq.ParamIndexing
		}
	}
}

// moqNewConverterFunc_mock isolates the mock interface of the NewConverterFunc
// type
type moqNewConverterFunc_mock struct {
	moq *moqNewConverterFunc
}

// moqNewConverterFunc_params holds the params of the NewConverterFunc type
type moqNewConverterFunc_params struct {
	typ    generator.Type
	export bool
}

// moqNewConverterFunc_paramsKey holds the map key params of the
// NewConverterFunc type
type moqNewConverterFunc_paramsKey struct {
	params struct{ export bool }
	hashes struct {
		typ    hash.Hash
		export hash.Hash
	}
}

// moqNewConverterFunc_resultsByParams contains the results for a given set of
// parameters for the NewConverterFunc type
type moqNewConverterFunc_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqNewConverterFunc_paramsKey]*moqNewConverterFunc_results
}

// moqNewConverterFunc_doFn defines the type of function needed when calling
// andDo for the NewConverterFunc type
type moqNewConverterFunc_doFn func(typ generator.Type, export bool)

// moqNewConverterFunc_doReturnFn defines the type of function needed when
// calling doReturnResults for the NewConverterFunc type
type moqNewConverterFunc_doReturnFn func(typ generator.Type, export bool) generator.Converterer

// moqNewConverterFunc_results holds the results of the NewConverterFunc type
type moqNewConverterFunc_results struct {
	params  moqNewConverterFunc_params
	results []struct {
		values *struct {
			result1 generator.Converterer
		}
		sequence   uint32
		doFn       moqNewConverterFunc_doFn
		doReturnFn moqNewConverterFunc_doReturnFn
	}
	index  uint32
	repeat *moq.RepeatVal
}

// moqNewConverterFunc_fnRecorder routes recorded function calls to the
// moqNewConverterFunc moq
type moqNewConverterFunc_fnRecorder struct {
	params    moqNewConverterFunc_params
	anyParams uint64
	sequence  bool
	results   *moqNewConverterFunc_results
	moq       *moqNewConverterFunc
}

// moqNewConverterFunc_anyParams isolates the any params functions of the
// NewConverterFunc type
type moqNewConverterFunc_anyParams struct {
	recorder *moqNewConverterFunc_fnRecorder
}

// newMoqNewConverterFunc creates a new moq of the NewConverterFunc type
func newMoqNewConverterFunc(scene *moq.Scene, config *moq.Config) *moqNewConverterFunc {
	if config == nil {
		config = &moq.Config{}
	}
	m := &moqNewConverterFunc{
		scene:  scene,
		config: *config,
		moq:    &moqNewConverterFunc_mock{},

		runtime: struct {
			parameterIndexing struct {
				typ    moq.ParamIndexing
				export moq.ParamIndexing
			}
		}{parameterIndexing: struct {
			typ    moq.ParamIndexing
			export moq.ParamIndexing
		}{
			typ:    moq.ParamIndexByHash,
			export: moq.ParamIndexByValue,
		}},
	}
	m.moq.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the moq implementation of the NewConverterFunc type
func (m *moqNewConverterFunc) mock() generator.NewConverterFunc {
	return func(typ generator.Type, export bool) generator.Converterer {
		moq := &moqNewConverterFunc_mock{moq: m}
		return moq.fn(typ, export)
	}
}

func (m *moqNewConverterFunc_mock) fn(typ generator.Type, export bool) (result1 generator.Converterer) {
	m.moq.scene.T.Helper()
	params := moqNewConverterFunc_params{
		typ:    typ,
		export: export,
	}
	var results *moqNewConverterFunc_results
	for _, resultsByParams := range m.moq.resultsByParams {
		paramsKey := m.moq.paramsKey(params, resultsByParams.anyParams)
		var ok bool
		results, ok = resultsByParams.results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.moq.config.Expectation == moq.Strict {
			m.moq.scene.T.Fatalf("Unexpected call to %s", m.moq.prettyParams(params))
		}
		return
	}

	i := int(atomic.AddUint32(&results.index, 1)) - 1
	if i >= results.repeat.ResultCount {
		if !results.repeat.AnyTimes {
			if m.moq.config.Expectation == moq.Strict {
				m.moq.scene.T.Fatalf("Too many calls to %s", m.moq.prettyParams(params))
			}
			return
		}
		i = results.repeat.ResultCount - 1
	}

	result := results.results[i]
	if result.sequence != 0 {
		sequence := m.moq.scene.NextMockSequence()
		if (!results.repeat.AnyTimes && result.sequence != sequence) || result.sequence > sequence {
			m.moq.scene.T.Fatalf("Call sequence does not match call to %s", m.moq.prettyParams(params))
		}
	}

	if result.doFn != nil {
		result.doFn(typ, export)
	}

	if result.values != nil {
		result1 = result.values.result1
	}
	if result.doReturnFn != nil {
		result1 = result.doReturnFn(typ, export)
	}
	return
}

func (m *moqNewConverterFunc) onCall(typ generator.Type, export bool) *moqNewConverterFunc_fnRecorder {
	return &moqNewConverterFunc_fnRecorder{
		params: moqNewConverterFunc_params{
			typ:    typ,
			export: export,
		},
		sequence: m.config.Sequence == moq.SeqDefaultOn,
		moq:      m,
	}
}

func (r *moqNewConverterFunc_fnRecorder) any() *moqNewConverterFunc_anyParams {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	return &moqNewConverterFunc_anyParams{recorder: r}
}

func (a *moqNewConverterFunc_anyParams) typ() *moqNewConverterFunc_fnRecorder {
	a.recorder.anyParams |= 1 << 0
	return a.recorder
}

func (a *moqNewConverterFunc_anyParams) export() *moqNewConverterFunc_fnRecorder {
	a.recorder.anyParams |= 1 << 1
	return a.recorder
}

func (r *moqNewConverterFunc_fnRecorder) seq() *moqNewConverterFunc_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqNewConverterFunc_fnRecorder) noSeq() *moqNewConverterFunc_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqNewConverterFunc_fnRecorder) returnResults(result1 generator.Converterer) *moqNewConverterFunc_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			result1 generator.Converterer
		}
		sequence   uint32
		doFn       moqNewConverterFunc_doFn
		doReturnFn moqNewConverterFunc_doReturnFn
	}{
		values: &struct {
			result1 generator.Converterer
		}{
			result1: result1,
		},
		sequence: sequence,
	})
	return r
}

func (r *moqNewConverterFunc_fnRecorder) andDo(fn moqNewConverterFunc_doFn) *moqNewConverterFunc_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqNewConverterFunc_fnRecorder) doReturnResults(fn moqNewConverterFunc_doReturnFn) *moqNewConverterFunc_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			result1 generator.Converterer
		}
		sequence   uint32
		doFn       moqNewConverterFunc_doFn
		doReturnFn moqNewConverterFunc_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqNewConverterFunc_fnRecorder) findResults() {
	if r.results != nil {
		r.results.repeat.Increment(r.moq.scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.anyParams)
	insertAt := -1
	var results *moqNewConverterFunc_resultsByParams
	for n, res := range r.moq.resultsByParams {
		if res.anyParams == r.anyParams {
			results = &res
			break
		}
		if res.anyCount > anyCount {
			insertAt = n
		}
	}
	if results == nil {
		results = &moqNewConverterFunc_resultsByParams{
			anyCount:  anyCount,
			anyParams: r.anyParams,
			results:   map[moqNewConverterFunc_paramsKey]*moqNewConverterFunc_results{},
		}
		r.moq.resultsByParams = append(r.moq.resultsByParams, *results)
		if insertAt != -1 && insertAt+1 < len(r.moq.resultsByParams) {
			copy(r.moq.resultsByParams[insertAt+1:], r.moq.resultsByParams[insertAt:0])
			r.moq.resultsByParams[insertAt] = *results
		}
	}

	paramsKey := r.moq.paramsKey(r.params, r.anyParams)

	var ok bool
	r.results, ok = results.results[paramsKey]
	if !ok {
		r.results = &moqNewConverterFunc_results{
			params:  r.params,
			results: nil,
			index:   0,
			repeat:  &moq.RepeatVal{},
		}
		results.results[paramsKey] = r.results
	}

	r.results.repeat.Increment(r.moq.scene.T)
}

func (r *moqNewConverterFunc_fnRecorder) repeat(repeaters ...moq.Repeater) *moqNewConverterFunc_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults or doReturnResults must be called before calling repeat")
		return nil
	}
	r.results.repeat.Repeat(r.moq.scene.T, repeaters)
	last := r.results.results[len(r.results.results)-1]
	for n := 0; n < r.results.repeat.ResultCount-1; n++ {
		if r.sequence {
			last = struct {
				values *struct {
					result1 generator.Converterer
				}
				sequence   uint32
				doFn       moqNewConverterFunc_doFn
				doReturnFn moqNewConverterFunc_doReturnFn
			}{
				values:   last.values,
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (m *moqNewConverterFunc) prettyParams(params moqNewConverterFunc_params) string {
	return fmt.Sprintf("NewConverterFunc(%#v, %#v)", params.typ, params.export)
}

func (m *moqNewConverterFunc) paramsKey(params moqNewConverterFunc_params, anyParams uint64) moqNewConverterFunc_paramsKey {
	var typUsedHash hash.Hash
	if anyParams&(1<<0) == 0 {
		if m.runtime.parameterIndexing.typ == moq.ParamIndexByValue {
			m.scene.T.Fatalf("The typ parameter can't be indexed by value")
		}
		typUsedHash = hash.DeepHash(params.typ)
	}
	var exportUsed bool
	var exportUsedHash hash.Hash
	if anyParams&(1<<1) == 0 {
		if m.runtime.parameterIndexing.export == moq.ParamIndexByValue {
			exportUsed = params.export
		} else {
			exportUsedHash = hash.DeepHash(params.export)
		}
	}
	return moqNewConverterFunc_paramsKey{
		params: struct{ export bool }{
			export: exportUsed,
		},
		hashes: struct {
			typ    hash.Hash
			export hash.Hash
		}{
			typ:    typUsedHash,
			export: exportUsedHash,
		},
	}
}

// Reset resets the state of the moq
func (m *moqNewConverterFunc) Reset() { m.resultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqNewConverterFunc) AssertExpectationsMet() {
	m.scene.T.Helper()
	for _, res := range m.resultsByParams {
		for _, results := range res.results {
			missing := results.repeat.MinTimes - int(atomic.LoadUint32(&results.index))
			if missing > 0 {
				m.scene.T.Errorf("Expected %d additional call(s) to %s", missing, m.prettyParams(results.params))
			}
		}
	}
}
