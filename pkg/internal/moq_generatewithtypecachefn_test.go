// Code generated by Moqueries - https://moqueries.org - DO NOT EDIT!

package internal_test

import (
	"fmt"
	"math/bits"
	"sync/atomic"

	"moqueries.org/cli/generator"
	"moqueries.org/cli/pkg/internal"
	"moqueries.org/runtime/hash"
	"moqueries.org/runtime/moq"
)

// moqGenerateWithTypeCacheFn holds the state of a moq of the
// GenerateWithTypeCacheFn type
type moqGenerateWithTypeCacheFn struct {
	scene  *moq.Scene
	config moq.Config
	moq    *moqGenerateWithTypeCacheFn_mock

	resultsByParams []moqGenerateWithTypeCacheFn_resultsByParams

	runtime struct {
		parameterIndexing struct {
			cache moq.ParamIndexing
			req   moq.ParamIndexing
		}
	}
}

// moqGenerateWithTypeCacheFn_mock isolates the mock interface of the
// GenerateWithTypeCacheFn type
type moqGenerateWithTypeCacheFn_mock struct {
	moq *moqGenerateWithTypeCacheFn
}

// moqGenerateWithTypeCacheFn_params holds the params of the
// GenerateWithTypeCacheFn type
type moqGenerateWithTypeCacheFn_params struct {
	cache generator.TypeCache
	req   generator.GenerateRequest
}

// moqGenerateWithTypeCacheFn_paramsKey holds the map key params of the
// GenerateWithTypeCacheFn type
type moqGenerateWithTypeCacheFn_paramsKey struct {
	params struct{ cache generator.TypeCache }
	hashes struct {
		cache hash.Hash
		req   hash.Hash
	}
}

// moqGenerateWithTypeCacheFn_resultsByParams contains the results for a given
// set of parameters for the GenerateWithTypeCacheFn type
type moqGenerateWithTypeCacheFn_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqGenerateWithTypeCacheFn_paramsKey]*moqGenerateWithTypeCacheFn_results
}

// moqGenerateWithTypeCacheFn_doFn defines the type of function needed when
// calling andDo for the GenerateWithTypeCacheFn type
type moqGenerateWithTypeCacheFn_doFn func(cache generator.TypeCache, req generator.GenerateRequest)

// moqGenerateWithTypeCacheFn_doReturnFn defines the type of function needed
// when calling doReturnResults for the GenerateWithTypeCacheFn type
type moqGenerateWithTypeCacheFn_doReturnFn func(cache generator.TypeCache, req generator.GenerateRequest) error

// moqGenerateWithTypeCacheFn_results holds the results of the
// GenerateWithTypeCacheFn type
type moqGenerateWithTypeCacheFn_results struct {
	params  moqGenerateWithTypeCacheFn_params
	results []struct {
		values *struct {
			result1 error
		}
		sequence   uint32
		doFn       moqGenerateWithTypeCacheFn_doFn
		doReturnFn moqGenerateWithTypeCacheFn_doReturnFn
	}
	index  uint32
	repeat *moq.RepeatVal
}

// moqGenerateWithTypeCacheFn_fnRecorder routes recorded function calls to the
// moqGenerateWithTypeCacheFn moq
type moqGenerateWithTypeCacheFn_fnRecorder struct {
	params    moqGenerateWithTypeCacheFn_params
	anyParams uint64
	sequence  bool
	results   *moqGenerateWithTypeCacheFn_results
	moq       *moqGenerateWithTypeCacheFn
}

// moqGenerateWithTypeCacheFn_anyParams isolates the any params functions of
// the GenerateWithTypeCacheFn type
type moqGenerateWithTypeCacheFn_anyParams struct {
	recorder *moqGenerateWithTypeCacheFn_fnRecorder
}

// newMoqGenerateWithTypeCacheFn creates a new moq of the
// GenerateWithTypeCacheFn type
func newMoqGenerateWithTypeCacheFn(scene *moq.Scene, config *moq.Config) *moqGenerateWithTypeCacheFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &moqGenerateWithTypeCacheFn{
		scene:  scene,
		config: *config,
		moq:    &moqGenerateWithTypeCacheFn_mock{},

		runtime: struct {
			parameterIndexing struct {
				cache moq.ParamIndexing
				req   moq.ParamIndexing
			}
		}{parameterIndexing: struct {
			cache moq.ParamIndexing
			req   moq.ParamIndexing
		}{
			cache: moq.ParamIndexByHash,
			req:   moq.ParamIndexByHash,
		}},
	}
	m.moq.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the moq implementation of the GenerateWithTypeCacheFn type
func (m *moqGenerateWithTypeCacheFn) mock() internal.GenerateWithTypeCacheFn {
	return func(cache generator.TypeCache, req generator.GenerateRequest) error {
		m.scene.T.Helper()
		moq := &moqGenerateWithTypeCacheFn_mock{moq: m}
		return moq.fn(cache, req)
	}
}

func (m *moqGenerateWithTypeCacheFn_mock) fn(cache generator.TypeCache, req generator.GenerateRequest) (result1 error) {
	m.moq.scene.T.Helper()
	params := moqGenerateWithTypeCacheFn_params{
		cache: cache,
		req:   req,
	}
	var results *moqGenerateWithTypeCacheFn_results
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
		result.doFn(cache, req)
	}

	if result.values != nil {
		result1 = result.values.result1
	}
	if result.doReturnFn != nil {
		result1 = result.doReturnFn(cache, req)
	}
	return
}

func (m *moqGenerateWithTypeCacheFn) onCall(cache generator.TypeCache, req generator.GenerateRequest) *moqGenerateWithTypeCacheFn_fnRecorder {
	return &moqGenerateWithTypeCacheFn_fnRecorder{
		params: moqGenerateWithTypeCacheFn_params{
			cache: cache,
			req:   req,
		},
		sequence: m.config.Sequence == moq.SeqDefaultOn,
		moq:      m,
	}
}

func (r *moqGenerateWithTypeCacheFn_fnRecorder) any() *moqGenerateWithTypeCacheFn_anyParams {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	return &moqGenerateWithTypeCacheFn_anyParams{recorder: r}
}

func (a *moqGenerateWithTypeCacheFn_anyParams) cache() *moqGenerateWithTypeCacheFn_fnRecorder {
	a.recorder.anyParams |= 1 << 0
	return a.recorder
}

func (a *moqGenerateWithTypeCacheFn_anyParams) req() *moqGenerateWithTypeCacheFn_fnRecorder {
	a.recorder.anyParams |= 1 << 1
	return a.recorder
}

func (r *moqGenerateWithTypeCacheFn_fnRecorder) seq() *moqGenerateWithTypeCacheFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqGenerateWithTypeCacheFn_fnRecorder) noSeq() *moqGenerateWithTypeCacheFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams(r.params))
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqGenerateWithTypeCacheFn_fnRecorder) returnResults(result1 error) *moqGenerateWithTypeCacheFn_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			result1 error
		}
		sequence   uint32
		doFn       moqGenerateWithTypeCacheFn_doFn
		doReturnFn moqGenerateWithTypeCacheFn_doReturnFn
	}{
		values: &struct {
			result1 error
		}{
			result1: result1,
		},
		sequence: sequence,
	})
	return r
}

func (r *moqGenerateWithTypeCacheFn_fnRecorder) andDo(fn moqGenerateWithTypeCacheFn_doFn) *moqGenerateWithTypeCacheFn_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqGenerateWithTypeCacheFn_fnRecorder) doReturnResults(fn moqGenerateWithTypeCacheFn_doReturnFn) *moqGenerateWithTypeCacheFn_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			result1 error
		}
		sequence   uint32
		doFn       moqGenerateWithTypeCacheFn_doFn
		doReturnFn moqGenerateWithTypeCacheFn_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqGenerateWithTypeCacheFn_fnRecorder) findResults() {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.results.repeat.Increment(r.moq.scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.anyParams)
	insertAt := -1
	var results *moqGenerateWithTypeCacheFn_resultsByParams
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
		results = &moqGenerateWithTypeCacheFn_resultsByParams{
			anyCount:  anyCount,
			anyParams: r.anyParams,
			results:   map[moqGenerateWithTypeCacheFn_paramsKey]*moqGenerateWithTypeCacheFn_results{},
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
		r.results = &moqGenerateWithTypeCacheFn_results{
			params:  r.params,
			results: nil,
			index:   0,
			repeat:  &moq.RepeatVal{},
		}
		results.results[paramsKey] = r.results
	}

	r.results.repeat.Increment(r.moq.scene.T)
}

func (r *moqGenerateWithTypeCacheFn_fnRecorder) repeat(repeaters ...moq.Repeater) *moqGenerateWithTypeCacheFn_fnRecorder {
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
					result1 error
				}
				sequence   uint32
				doFn       moqGenerateWithTypeCacheFn_doFn
				doReturnFn moqGenerateWithTypeCacheFn_doReturnFn
			}{
				values:   last.values,
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (m *moqGenerateWithTypeCacheFn) prettyParams(params moqGenerateWithTypeCacheFn_params) string {
	return fmt.Sprintf("GenerateWithTypeCacheFn(%#v, %#v)", params.cache, params.req)
}

func (m *moqGenerateWithTypeCacheFn) paramsKey(params moqGenerateWithTypeCacheFn_params, anyParams uint64) moqGenerateWithTypeCacheFn_paramsKey {
	m.scene.T.Helper()
	var cacheUsed generator.TypeCache
	var cacheUsedHash hash.Hash
	if anyParams&(1<<0) == 0 {
		if m.runtime.parameterIndexing.cache == moq.ParamIndexByValue {
			cacheUsed = params.cache
		} else {
			cacheUsedHash = hash.DeepHash(params.cache)
		}
	}
	var reqUsedHash hash.Hash
	if anyParams&(1<<1) == 0 {
		if m.runtime.parameterIndexing.req == moq.ParamIndexByValue {
			m.scene.T.Fatalf("The req parameter can't be indexed by value")
		}
		reqUsedHash = hash.DeepHash(params.req)
	}
	return moqGenerateWithTypeCacheFn_paramsKey{
		params: struct{ cache generator.TypeCache }{
			cache: cacheUsed,
		},
		hashes: struct {
			cache hash.Hash
			req   hash.Hash
		}{
			cache: cacheUsedHash,
			req:   reqUsedHash,
		},
	}
}

// Reset resets the state of the moq
func (m *moqGenerateWithTypeCacheFn) Reset() { m.resultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqGenerateWithTypeCacheFn) AssertExpectationsMet() {
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
