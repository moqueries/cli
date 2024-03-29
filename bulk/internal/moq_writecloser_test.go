// Code generated by Moqueries - https://moqueries.org - DO NOT EDIT.

package internal_test

import (
	"fmt"
	"io"
	"math/bits"
	"sync/atomic"

	"moqueries.org/runtime/hash"
	"moqueries.org/runtime/moq"
)

// The following type assertion assures that io.WriteCloser is mocked
// completely
var _ io.WriteCloser = (*moqWriteCloser_mock)(nil)

// moqWriteCloser holds the state of a moq of the WriteCloser type
type moqWriteCloser struct {
	scene  *moq.Scene
	config moq.Config
	moq    *moqWriteCloser_mock

	resultsByParams_Write []moqWriteCloser_Write_resultsByParams
	resultsByParams_Close []moqWriteCloser_Close_resultsByParams

	runtime struct {
		parameterIndexing struct {
			Write struct {
				p moq.ParamIndexing
			}
			Close struct{}
		}
	}
}

// moqWriteCloser_mock isolates the mock interface of the WriteCloser type
type moqWriteCloser_mock struct {
	moq *moqWriteCloser
}

// moqWriteCloser_recorder isolates the recorder interface of the WriteCloser
// type
type moqWriteCloser_recorder struct {
	moq *moqWriteCloser
}

// moqWriteCloser_Write_params holds the params of the WriteCloser type
type moqWriteCloser_Write_params struct{ p []byte }

// moqWriteCloser_Write_paramsKey holds the map key params of the WriteCloser
// type
type moqWriteCloser_Write_paramsKey struct {
	params struct{}
	hashes struct{ p hash.Hash }
}

// moqWriteCloser_Write_resultsByParams contains the results for a given set of
// parameters for the WriteCloser type
type moqWriteCloser_Write_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqWriteCloser_Write_paramsKey]*moqWriteCloser_Write_results
}

// moqWriteCloser_Write_doFn defines the type of function needed when calling
// andDo for the WriteCloser type
type moqWriteCloser_Write_doFn func(p []byte)

// moqWriteCloser_Write_doReturnFn defines the type of function needed when
// calling doReturnResults for the WriteCloser type
type moqWriteCloser_Write_doReturnFn func(p []byte) (n int, err error)

// moqWriteCloser_Write_results holds the results of the WriteCloser type
type moqWriteCloser_Write_results struct {
	params  moqWriteCloser_Write_params
	results []struct {
		values *struct {
			n   int
			err error
		}
		sequence   uint32
		doFn       moqWriteCloser_Write_doFn
		doReturnFn moqWriteCloser_Write_doReturnFn
	}
	index  uint32
	repeat *moq.RepeatVal
}

// moqWriteCloser_Write_fnRecorder routes recorded function calls to the
// moqWriteCloser moq
type moqWriteCloser_Write_fnRecorder struct {
	params    moqWriteCloser_Write_params
	anyParams uint64
	sequence  bool
	results   *moqWriteCloser_Write_results
	moq       *moqWriteCloser
}

// moqWriteCloser_Write_anyParams isolates the any params functions of the
// WriteCloser type
type moqWriteCloser_Write_anyParams struct {
	recorder *moqWriteCloser_Write_fnRecorder
}

// moqWriteCloser_Close_params holds the params of the WriteCloser type
type moqWriteCloser_Close_params struct{}

// moqWriteCloser_Close_paramsKey holds the map key params of the WriteCloser
// type
type moqWriteCloser_Close_paramsKey struct {
	params struct{}
	hashes struct{}
}

// moqWriteCloser_Close_resultsByParams contains the results for a given set of
// parameters for the WriteCloser type
type moqWriteCloser_Close_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqWriteCloser_Close_paramsKey]*moqWriteCloser_Close_results
}

// moqWriteCloser_Close_doFn defines the type of function needed when calling
// andDo for the WriteCloser type
type moqWriteCloser_Close_doFn func()

// moqWriteCloser_Close_doReturnFn defines the type of function needed when
// calling doReturnResults for the WriteCloser type
type moqWriteCloser_Close_doReturnFn func() error

// moqWriteCloser_Close_results holds the results of the WriteCloser type
type moqWriteCloser_Close_results struct {
	params  moqWriteCloser_Close_params
	results []struct {
		values *struct {
			result1 error
		}
		sequence   uint32
		doFn       moqWriteCloser_Close_doFn
		doReturnFn moqWriteCloser_Close_doReturnFn
	}
	index  uint32
	repeat *moq.RepeatVal
}

// moqWriteCloser_Close_fnRecorder routes recorded function calls to the
// moqWriteCloser moq
type moqWriteCloser_Close_fnRecorder struct {
	params    moqWriteCloser_Close_params
	anyParams uint64
	sequence  bool
	results   *moqWriteCloser_Close_results
	moq       *moqWriteCloser
}

// moqWriteCloser_Close_anyParams isolates the any params functions of the
// WriteCloser type
type moqWriteCloser_Close_anyParams struct {
	recorder *moqWriteCloser_Close_fnRecorder
}

// newMoqWriteCloser creates a new moq of the WriteCloser type
func newMoqWriteCloser(scene *moq.Scene, config *moq.Config) *moqWriteCloser {
	if config == nil {
		config = &moq.Config{}
	}
	m := &moqWriteCloser{
		scene:  scene,
		config: *config,
		moq:    &moqWriteCloser_mock{},

		runtime: struct {
			parameterIndexing struct {
				Write struct {
					p moq.ParamIndexing
				}
				Close struct{}
			}
		}{parameterIndexing: struct {
			Write struct {
				p moq.ParamIndexing
			}
			Close struct{}
		}{
			Write: struct {
				p moq.ParamIndexing
			}{
				p: moq.ParamIndexByHash,
			},
			Close: struct{}{},
		}},
	}
	m.moq.moq = m

	scene.AddMoq(m)
	return m
}

// mock returns the mock implementation of the WriteCloser type
func (m *moqWriteCloser) mock() *moqWriteCloser_mock { return m.moq }

func (m *moqWriteCloser_mock) Write(p []byte) (n int, err error) {
	m.moq.scene.T.Helper()
	params := moqWriteCloser_Write_params{
		p: p,
	}
	var results *moqWriteCloser_Write_results
	for _, resultsByParams := range m.moq.resultsByParams_Write {
		paramsKey := m.moq.paramsKey_Write(params, resultsByParams.anyParams)
		var ok bool
		results, ok = resultsByParams.results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.moq.config.Expectation == moq.Strict {
			m.moq.scene.T.Fatalf("Unexpected call to %s", m.moq.prettyParams_Write(params))
		}
		return
	}

	i := int(atomic.AddUint32(&results.index, 1)) - 1
	if i >= results.repeat.ResultCount {
		if !results.repeat.AnyTimes {
			if m.moq.config.Expectation == moq.Strict {
				m.moq.scene.T.Fatalf("Too many calls to %s", m.moq.prettyParams_Write(params))
			}
			return
		}
		i = results.repeat.ResultCount - 1
	}

	result := results.results[i]
	if result.sequence != 0 {
		sequence := m.moq.scene.NextMockSequence()
		if (!results.repeat.AnyTimes && result.sequence != sequence) || result.sequence > sequence {
			m.moq.scene.T.Fatalf("Call sequence does not match call to %s", m.moq.prettyParams_Write(params))
		}
	}

	if result.doFn != nil {
		result.doFn(p)
	}

	if result.values != nil {
		n = result.values.n
		err = result.values.err
	}
	if result.doReturnFn != nil {
		n, err = result.doReturnFn(p)
	}
	return
}

func (m *moqWriteCloser_mock) Close() (result1 error) {
	m.moq.scene.T.Helper()
	params := moqWriteCloser_Close_params{}
	var results *moqWriteCloser_Close_results
	for _, resultsByParams := range m.moq.resultsByParams_Close {
		paramsKey := m.moq.paramsKey_Close(params, resultsByParams.anyParams)
		var ok bool
		results, ok = resultsByParams.results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.moq.config.Expectation == moq.Strict {
			m.moq.scene.T.Fatalf("Unexpected call to %s", m.moq.prettyParams_Close(params))
		}
		return
	}

	i := int(atomic.AddUint32(&results.index, 1)) - 1
	if i >= results.repeat.ResultCount {
		if !results.repeat.AnyTimes {
			if m.moq.config.Expectation == moq.Strict {
				m.moq.scene.T.Fatalf("Too many calls to %s", m.moq.prettyParams_Close(params))
			}
			return
		}
		i = results.repeat.ResultCount - 1
	}

	result := results.results[i]
	if result.sequence != 0 {
		sequence := m.moq.scene.NextMockSequence()
		if (!results.repeat.AnyTimes && result.sequence != sequence) || result.sequence > sequence {
			m.moq.scene.T.Fatalf("Call sequence does not match call to %s", m.moq.prettyParams_Close(params))
		}
	}

	if result.doFn != nil {
		result.doFn()
	}

	if result.values != nil {
		result1 = result.values.result1
	}
	if result.doReturnFn != nil {
		result1 = result.doReturnFn()
	}
	return
}

// onCall returns the recorder implementation of the WriteCloser type
func (m *moqWriteCloser) onCall() *moqWriteCloser_recorder {
	return &moqWriteCloser_recorder{
		moq: m,
	}
}

func (m *moqWriteCloser_recorder) Write(p []byte) *moqWriteCloser_Write_fnRecorder {
	return &moqWriteCloser_Write_fnRecorder{
		params: moqWriteCloser_Write_params{
			p: p,
		},
		sequence: m.moq.config.Sequence == moq.SeqDefaultOn,
		moq:      m.moq,
	}
}

func (r *moqWriteCloser_Write_fnRecorder) any() *moqWriteCloser_Write_anyParams {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_Write(r.params))
		return nil
	}
	return &moqWriteCloser_Write_anyParams{recorder: r}
}

func (a *moqWriteCloser_Write_anyParams) p() *moqWriteCloser_Write_fnRecorder {
	a.recorder.anyParams |= 1 << 0
	return a.recorder
}

func (r *moqWriteCloser_Write_fnRecorder) seq() *moqWriteCloser_Write_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_Write(r.params))
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqWriteCloser_Write_fnRecorder) noSeq() *moqWriteCloser_Write_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_Write(r.params))
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqWriteCloser_Write_fnRecorder) returnResults(n int, err error) *moqWriteCloser_Write_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			n   int
			err error
		}
		sequence   uint32
		doFn       moqWriteCloser_Write_doFn
		doReturnFn moqWriteCloser_Write_doReturnFn
	}{
		values: &struct {
			n   int
			err error
		}{
			n:   n,
			err: err,
		},
		sequence: sequence,
	})
	return r
}

func (r *moqWriteCloser_Write_fnRecorder) andDo(fn moqWriteCloser_Write_doFn) *moqWriteCloser_Write_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqWriteCloser_Write_fnRecorder) doReturnResults(fn moqWriteCloser_Write_doReturnFn) *moqWriteCloser_Write_fnRecorder {
	r.moq.scene.T.Helper()
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			n   int
			err error
		}
		sequence   uint32
		doFn       moqWriteCloser_Write_doFn
		doReturnFn moqWriteCloser_Write_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqWriteCloser_Write_fnRecorder) findResults() {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.results.repeat.Increment(r.moq.scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.anyParams)
	insertAt := -1
	var results *moqWriteCloser_Write_resultsByParams
	for n, res := range r.moq.resultsByParams_Write {
		if res.anyParams == r.anyParams {
			results = &res
			break
		}
		if res.anyCount > anyCount {
			insertAt = n
		}
	}
	if results == nil {
		results = &moqWriteCloser_Write_resultsByParams{
			anyCount:  anyCount,
			anyParams: r.anyParams,
			results:   map[moqWriteCloser_Write_paramsKey]*moqWriteCloser_Write_results{},
		}
		r.moq.resultsByParams_Write = append(r.moq.resultsByParams_Write, *results)
		if insertAt != -1 && insertAt+1 < len(r.moq.resultsByParams_Write) {
			copy(r.moq.resultsByParams_Write[insertAt+1:], r.moq.resultsByParams_Write[insertAt:0])
			r.moq.resultsByParams_Write[insertAt] = *results
		}
	}

	paramsKey := r.moq.paramsKey_Write(r.params, r.anyParams)

	var ok bool
	r.results, ok = results.results[paramsKey]
	if !ok {
		r.results = &moqWriteCloser_Write_results{
			params:  r.params,
			results: nil,
			index:   0,
			repeat:  &moq.RepeatVal{},
		}
		results.results[paramsKey] = r.results
	}

	r.results.repeat.Increment(r.moq.scene.T)
}

func (r *moqWriteCloser_Write_fnRecorder) repeat(repeaters ...moq.Repeater) *moqWriteCloser_Write_fnRecorder {
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
					n   int
					err error
				}
				sequence   uint32
				doFn       moqWriteCloser_Write_doFn
				doReturnFn moqWriteCloser_Write_doReturnFn
			}{
				values:   last.values,
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (m *moqWriteCloser) prettyParams_Write(params moqWriteCloser_Write_params) string {
	return fmt.Sprintf("Write(%#v)", params.p)
}

func (m *moqWriteCloser) paramsKey_Write(params moqWriteCloser_Write_params, anyParams uint64) moqWriteCloser_Write_paramsKey {
	m.scene.T.Helper()
	var pUsedHash hash.Hash
	if anyParams&(1<<0) == 0 {
		if m.runtime.parameterIndexing.Write.p == moq.ParamIndexByValue {
			m.scene.T.Fatalf("The p parameter of the Write function can't be indexed by value")
		}
		pUsedHash = hash.DeepHash(params.p)
	}
	return moqWriteCloser_Write_paramsKey{
		params: struct{}{},
		hashes: struct{ p hash.Hash }{
			p: pUsedHash,
		},
	}
}

func (m *moqWriteCloser_recorder) Close() *moqWriteCloser_Close_fnRecorder {
	return &moqWriteCloser_Close_fnRecorder{
		params:   moqWriteCloser_Close_params{},
		sequence: m.moq.config.Sequence == moq.SeqDefaultOn,
		moq:      m.moq,
	}
}

func (r *moqWriteCloser_Close_fnRecorder) any() *moqWriteCloser_Close_anyParams {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_Close(r.params))
		return nil
	}
	return &moqWriteCloser_Close_anyParams{recorder: r}
}

func (r *moqWriteCloser_Close_fnRecorder) seq() *moqWriteCloser_Close_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_Close(r.params))
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqWriteCloser_Close_fnRecorder) noSeq() *moqWriteCloser_Close_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, recording %s", r.moq.prettyParams_Close(r.params))
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqWriteCloser_Close_fnRecorder) returnResults(result1 error) *moqWriteCloser_Close_fnRecorder {
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
		doFn       moqWriteCloser_Close_doFn
		doReturnFn moqWriteCloser_Close_doReturnFn
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

func (r *moqWriteCloser_Close_fnRecorder) andDo(fn moqWriteCloser_Close_doFn) *moqWriteCloser_Close_fnRecorder {
	r.moq.scene.T.Helper()
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqWriteCloser_Close_fnRecorder) doReturnResults(fn moqWriteCloser_Close_doReturnFn) *moqWriteCloser_Close_fnRecorder {
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
		doFn       moqWriteCloser_Close_doFn
		doReturnFn moqWriteCloser_Close_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqWriteCloser_Close_fnRecorder) findResults() {
	r.moq.scene.T.Helper()
	if r.results != nil {
		r.results.repeat.Increment(r.moq.scene.T)
		return
	}

	anyCount := bits.OnesCount64(r.anyParams)
	insertAt := -1
	var results *moqWriteCloser_Close_resultsByParams
	for n, res := range r.moq.resultsByParams_Close {
		if res.anyParams == r.anyParams {
			results = &res
			break
		}
		if res.anyCount > anyCount {
			insertAt = n
		}
	}
	if results == nil {
		results = &moqWriteCloser_Close_resultsByParams{
			anyCount:  anyCount,
			anyParams: r.anyParams,
			results:   map[moqWriteCloser_Close_paramsKey]*moqWriteCloser_Close_results{},
		}
		r.moq.resultsByParams_Close = append(r.moq.resultsByParams_Close, *results)
		if insertAt != -1 && insertAt+1 < len(r.moq.resultsByParams_Close) {
			copy(r.moq.resultsByParams_Close[insertAt+1:], r.moq.resultsByParams_Close[insertAt:0])
			r.moq.resultsByParams_Close[insertAt] = *results
		}
	}

	paramsKey := r.moq.paramsKey_Close(r.params, r.anyParams)

	var ok bool
	r.results, ok = results.results[paramsKey]
	if !ok {
		r.results = &moqWriteCloser_Close_results{
			params:  r.params,
			results: nil,
			index:   0,
			repeat:  &moq.RepeatVal{},
		}
		results.results[paramsKey] = r.results
	}

	r.results.repeat.Increment(r.moq.scene.T)
}

func (r *moqWriteCloser_Close_fnRecorder) repeat(repeaters ...moq.Repeater) *moqWriteCloser_Close_fnRecorder {
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
				doFn       moqWriteCloser_Close_doFn
				doReturnFn moqWriteCloser_Close_doReturnFn
			}{
				values:   last.values,
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	return r
}

func (m *moqWriteCloser) prettyParams_Close(params moqWriteCloser_Close_params) string {
	return fmt.Sprintf("Close()")
}

func (m *moqWriteCloser) paramsKey_Close(params moqWriteCloser_Close_params, anyParams uint64) moqWriteCloser_Close_paramsKey {
	m.scene.T.Helper()
	return moqWriteCloser_Close_paramsKey{
		params: struct{}{},
		hashes: struct{}{},
	}
}

// Reset resets the state of the moq
func (m *moqWriteCloser) Reset() { m.resultsByParams_Write = nil; m.resultsByParams_Close = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqWriteCloser) AssertExpectationsMet() {
	m.scene.T.Helper()
	for _, res := range m.resultsByParams_Write {
		for _, results := range res.results {
			missing := results.repeat.MinTimes - int(atomic.LoadUint32(&results.index))
			if missing > 0 {
				m.scene.T.Errorf("Expected %d additional call(s) to %s", missing, m.prettyParams_Write(results.params))
			}
		}
	}
	for _, res := range m.resultsByParams_Close {
		for _, results := range res.results {
			missing := results.repeat.MinTimes - int(atomic.LoadUint32(&results.index))
			if missing > 0 {
				m.scene.T.Errorf("Expected %d additional call(s) to %s", missing, m.prettyParams_Close(results.params))
			}
		}
	}
}
