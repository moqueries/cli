// Code generated by Moqueries - https://moqueries.org - DO NOT EDIT!

package ast_test

import (
	"fmt"
	"math/bits"
	"os"
	"sync/atomic"

	"moqueries.org/cli/ast"
	"moqueries.org/runtime/hash"
	"moqueries.org/runtime/moq"
)

// moqStatFn holds the state of a moq of the StatFn type
type moqStatFn struct {
	scene  *moq.Scene
	config moq.Config
	moq    *moqStatFn_mock

	resultsByParams []moqStatFn_resultsByParams

	runtime struct {
		parameterIndexing struct {
			name moq.ParamIndexing
		}
	}
}

// moqStatFn_mock isolates the mock interface of the StatFn type
type moqStatFn_mock struct {
	moq *moqStatFn
}

// moqStatFn_params holds the params of the StatFn type
type moqStatFn_params struct{ name string }

// moqStatFn_paramsKey holds the map key params of the StatFn type
type moqStatFn_paramsKey struct {
	params struct{ name string }
	hashes struct{ name hash.Hash }
}

// moqStatFn_resultsByParams contains the results for a given set of parameters
// for the StatFn type
type moqStatFn_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqStatFn_paramsKey]*moqStatFn_results
}

// moqStatFn_doFn defines the type of function needed when calling andDo for
// the StatFn type
type moqStatFn_doFn func(name string)

// moqStatFn_doReturnFn defines the type of function needed when calling
// doReturnResults for the StatFn type
type moqStatFn_doReturnFn func(name string) (os.FileInfo, error)

// moqStatFn_results holds the results of the StatFn type
type moqStatFn_results struct {
	params  moqStatFn_params
	results []struct {
		values *struct {
			result1 os.FileInfo
			result2 error
		}
		sequence   uint32
		doFn       moqStatFn_doFn
		doReturnFn moqStatFn_doReturnFn
	}
	index  uint32
	repeat *moq.RepeatVal
}

// moqStatFn_fnRecorder routes recorded function calls to the moqStatFn moq
type moqStatFn_fnRecorder struct {
	params    moqStatFn_params
	anyParams uint64
	sequence  bool
	results   *moqStatFn_results
	moq       *moqStatFn
}

// moqStatFn_anyParams isolates the any params functions of the StatFn type
type moqStatFn_anyParams struct {
	recorder *moqStatFn_fnRecorder
}

// newMoqStatFn creates a new moq of the StatFn type
func newMoqStatFn(scene *moq.Scene, config *moq.Config) *moqStatFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &moqStatFn{
		scene:  scene,
		config: *config,
		moq:    &moqStatFn_mock{},

		runtime: struct {
			parameterIndexing struct {
				name moq.ParamIndexing
			}
		}{parameterIndexing: struct {
			name moq.ParamIndexing
		}{
			name: moq.ParamIndexByValue,
		}},
	}
	m.moq.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the moq implementation of the StatFn type
func (m *moqStatFn) mock() ast.StatFn {
	return func(name string) (os.FileInfo, error) {
		m.scene.T.Helper()
		moq := &moqStatFn_mock{moq: m}
		return moq.fn(name)
	}
}

func (m *moqStatFn_mock) fn(name string) (result1 os.FileInfo, result2 error) {
	m.moq.scene.T.Helper()
	params := moqStatFn_params{
		name: name,
	}
	var results *moqStatFn_results
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
		result.doFn(name)
	}

	if result.values != nil {
		result1 = result.values.result1
		result2 = result.values.result2
	}
	if result.doReturnFn != nil {
		result1, result2 = result.doReturnFn(name)
	}
	return
}

func (m *moqStatFn) onCall(name string) *moqStatFn_fnRecorder {
	return &moqStatFn_fnRecorder{
		params: moqStatFn_params{
			name: name,
		},
		sequence: m.config.Sequence == moq.SeqDefaultOn,
		moq:      m,
	}
}

func (r *moqStatFn_fnRecorder) any() *moqStatFn_anyParams {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	return &moqStatFn_anyParams{recorder: r}
}

func (a *moqStatFn_anyParams) name() *moqStatFn_fnRecorder {
	a.recorder.anyParams |= 1 << 0
	return a.recorder
}

func (r *moqStatFn_fnRecorder) seq() *moqStatFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqStatFn_fnRecorder) noSeq() *moqStatFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqStatFn_fnRecorder) returnResults(result1 os.FileInfo, result2 error) *moqStatFn_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			result1 os.FileInfo
			result2 error
		}
		sequence   uint32
		doFn       moqStatFn_doFn
		doReturnFn moqStatFn_doReturnFn
	}{
		values: &struct {
			result1 os.FileInfo
			result2 error
		}{
			result1: result1,
			result2: result2,
		},
		sequence: sequence,
	})
	return r
}

func (r *moqStatFn_fnRecorder) andDo(fn moqStatFn_doFn) *moqStatFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqStatFn_fnRecorder) doReturnResults(fn moqStatFn_doReturnFn) *moqStatFn_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			result1 os.FileInfo
			result2 error
		}
		sequence   uint32
		doFn       moqStatFn_doFn
		doReturnFn moqStatFn_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqStatFn_fnRecorder) findResults() {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.results.repeat.Increment(r.moq.scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.anyParams)
	insertAt := -1
	var results *moqStatFn_resultsByParams
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
		results = &moqStatFn_resultsByParams{
			anyCount:  anyCount,
			anyParams: r.anyParams,
			results:   map[moqStatFn_paramsKey]*moqStatFn_results{},
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
		r.results = &moqStatFn_results{
			params:  r.params,
			results: nil,
			index:   0,
			repeat:  &moq.RepeatVal{},
		}
		results.results[paramsKey] = r.results
	}

	r.results.repeat.Increment(r.moq.scene.T)
}

func (r *moqStatFn_fnRecorder) repeat(repeaters ...moq.Repeater) *moqStatFn_fnRecorder {
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
					result1 os.FileInfo
					result2 error
				}
				sequence   uint32
				doFn       moqStatFn_doFn
				doReturnFn moqStatFn_doReturnFn
			}{
				values:   last.values,
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (m *moqStatFn) prettyParams(params moqStatFn_params) string {
	return fmt.Sprintf("StatFn(%#v)", params.name)
}

func (m *moqStatFn) paramsKey(params moqStatFn_params, anyParams uint64) moqStatFn_paramsKey {
	m.scene.T.Helper()
	var nameUsed string
	var nameUsedHash hash.Hash
	if anyParams&(1<<0) == 0 {
		if m.runtime.parameterIndexing.name == moq.ParamIndexByValue {
			nameUsed = params.name
		} else {
			nameUsedHash = hash.DeepHash(params.name)
		}
	}
	return moqStatFn_paramsKey{
		params: struct{ name string }{
			name: nameUsed,
		},
		hashes: struct{ name hash.Hash }{
			name: nameUsedHash,
		},
	}
}

// Reset resets the state of the moq
func (m *moqStatFn) Reset() { m.resultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqStatFn) AssertExpectationsMet() {
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
