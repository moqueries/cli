// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package internal_test

import (
	"fmt"
	"io"
	"math/bits"
	"sync/atomic"

	"github.com/myshkin5/moqueries/bulk/internal"
	"github.com/myshkin5/moqueries/hash"
	"github.com/myshkin5/moqueries/moq"
)

// moqOpenFn holds the state of a moq of the OpenFn type
type moqOpenFn struct {
	scene  *moq.Scene
	config moq.Config
	moq    *moqOpenFn_mock

	resultsByParams []moqOpenFn_resultsByParams

	runtime struct {
		parameterIndexing struct {
			name moq.ParamIndexing
		}
	}
}

// moqOpenFn_mock isolates the mock interface of the OpenFn type
type moqOpenFn_mock struct {
	moq *moqOpenFn
}

// moqOpenFn_params holds the params of the OpenFn type
type moqOpenFn_params struct{ name string }

// moqOpenFn_paramsKey holds the map key params of the OpenFn type
type moqOpenFn_paramsKey struct {
	params struct{ name string }
	hashes struct{ name hash.Hash }
}

// moqOpenFn_resultsByParams contains the results for a given set of parameters
// for the OpenFn type
type moqOpenFn_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqOpenFn_paramsKey]*moqOpenFn_results
}

// moqOpenFn_doFn defines the type of function needed when calling andDo for
// the OpenFn type
type moqOpenFn_doFn func(name string)

// moqOpenFn_doReturnFn defines the type of function needed when calling
// doReturnResults for the OpenFn type
type moqOpenFn_doReturnFn func(name string) (file io.ReadCloser, err error)

// moqOpenFn_results holds the results of the OpenFn type
type moqOpenFn_results struct {
	params  moqOpenFn_params
	results []struct {
		values *struct {
			file io.ReadCloser
			err  error
		}
		sequence   uint32
		doFn       moqOpenFn_doFn
		doReturnFn moqOpenFn_doReturnFn
	}
	index  uint32
	repeat *moq.RepeatVal
}

// moqOpenFn_fnRecorder routes recorded function calls to the moqOpenFn moq
type moqOpenFn_fnRecorder struct {
	params    moqOpenFn_params
	anyParams uint64
	sequence  bool
	results   *moqOpenFn_results
	moq       *moqOpenFn
}

// moqOpenFn_anyParams isolates the any params functions of the OpenFn type
type moqOpenFn_anyParams struct {
	recorder *moqOpenFn_fnRecorder
}

// newMoqOpenFn creates a new moq of the OpenFn type
func newMoqOpenFn(scene *moq.Scene, config *moq.Config) *moqOpenFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &moqOpenFn{
		scene:  scene,
		config: *config,
		moq:    &moqOpenFn_mock{},

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

// mock returns the moq implementation of the OpenFn type
func (m *moqOpenFn) mock() internal.OpenFn {
	return func(name string) (_ io.ReadCloser, _ error) { moq := &moqOpenFn_mock{moq: m}; return moq.fn(name) }
}

func (m *moqOpenFn_mock) fn(name string) (file io.ReadCloser, err error) {
	m.moq.scene.T.Helper()
	params := moqOpenFn_params{
		name: name,
	}
	var results *moqOpenFn_results
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
		file = result.values.file
		err = result.values.err
	}
	if result.doReturnFn != nil {
		file, err = result.doReturnFn(name)
	}
	return
}

func (m *moqOpenFn) onCall(name string) *moqOpenFn_fnRecorder {
	return &moqOpenFn_fnRecorder{
		params: moqOpenFn_params{
			name: name,
		},
		sequence: m.config.Sequence == moq.SeqDefaultOn,
		moq:      m,
	}
}

func (r *moqOpenFn_fnRecorder) any() *moqOpenFn_anyParams {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	return &moqOpenFn_anyParams{recorder: r}
}

func (a *moqOpenFn_anyParams) name() *moqOpenFn_fnRecorder {
	a.recorder.anyParams |= 1 << 0
	return a.recorder
}

func (r *moqOpenFn_fnRecorder) seq() *moqOpenFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqOpenFn_fnRecorder) noSeq() *moqOpenFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqOpenFn_fnRecorder) returnResults(file io.ReadCloser, err error) *moqOpenFn_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			file io.ReadCloser
			err  error
		}
		sequence   uint32
		doFn       moqOpenFn_doFn
		doReturnFn moqOpenFn_doReturnFn
	}{
		values: &struct {
			file io.ReadCloser
			err  error
		}{
			file: file,
			err:  err,
		},
		sequence: sequence,
	})
	return r
}

func (r *moqOpenFn_fnRecorder) andDo(fn moqOpenFn_doFn) *moqOpenFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqOpenFn_fnRecorder) doReturnResults(fn moqOpenFn_doReturnFn) *moqOpenFn_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			file io.ReadCloser
			err  error
		}
		sequence   uint32
		doFn       moqOpenFn_doFn
		doReturnFn moqOpenFn_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqOpenFn_fnRecorder) findResults() {
	if r.results != nil {
		r.results.repeat.Increment(r.moq.scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.anyParams)
	insertAt := -1
	var results *moqOpenFn_resultsByParams
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
		results = &moqOpenFn_resultsByParams{
			anyCount:  anyCount,
			anyParams: r.anyParams,
			results:   map[moqOpenFn_paramsKey]*moqOpenFn_results{},
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
		r.results = &moqOpenFn_results{
			params:  r.params,
			results: nil,
			index:   0,
			repeat:  &moq.RepeatVal{},
		}
		results.results[paramsKey] = r.results
	}

	r.results.repeat.Increment(r.moq.scene.T)
}

func (r *moqOpenFn_fnRecorder) repeat(repeaters ...moq.Repeater) *moqOpenFn_fnRecorder {
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
					file io.ReadCloser
					err  error
				}
				sequence   uint32
				doFn       moqOpenFn_doFn
				doReturnFn moqOpenFn_doReturnFn
			}{
				values:   last.values,
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (m *moqOpenFn) prettyParams(params moqOpenFn_params) string {
	return fmt.Sprintf("OpenFn(%#v)", params.name)
}

func (m *moqOpenFn) paramsKey(params moqOpenFn_params, anyParams uint64) moqOpenFn_paramsKey {
	var nameUsed string
	var nameUsedHash hash.Hash
	if anyParams&(1<<0) == 0 {
		if m.runtime.parameterIndexing.name == moq.ParamIndexByValue {
			nameUsed = params.name
		} else {
			nameUsedHash = hash.DeepHash(params.name)
		}
	}
	return moqOpenFn_paramsKey{
		params: struct{ name string }{
			name: nameUsed,
		},
		hashes: struct{ name hash.Hash }{
			name: nameUsedHash,
		},
	}
}

// Reset resets the state of the moq
func (m *moqOpenFn) Reset() { m.resultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqOpenFn) AssertExpectationsMet() {
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
