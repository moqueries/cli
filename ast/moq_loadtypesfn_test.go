// Code generated by Moqueries - https://github.com/myshkin5/moqueries - DO NOT EDIT!

package ast_test

import (
	"math/bits"
	"sync/atomic"

	"github.com/dave/dst"
	"github.com/myshkin5/moqueries/ast"
	"github.com/myshkin5/moqueries/moq"
)

// moqLoadTypesFn holds the state of a moq of the LoadTypesFn type
type moqLoadTypesFn struct {
	scene           *moq.Scene
	config          moq.Config
	resultsByParams []moqLoadTypesFn_resultsByParams
}

// moqLoadTypesFn_mock isolates the mock interface of the LoadTypesFn type
type moqLoadTypesFn_mock struct {
	moq *moqLoadTypesFn
}

// moqLoadTypesFn_params holds the params of the LoadTypesFn type
type moqLoadTypesFn_params struct {
	pkg           string
	loadTestTypes bool
}

// moqLoadTypesFn_paramsKey holds the map key params of the LoadTypesFn type
type moqLoadTypesFn_paramsKey struct {
	pkg           string
	loadTestTypes bool
}

// moqLoadTypesFn_resultsByParams contains the results for a given set of parameters for the LoadTypesFn type
type moqLoadTypesFn_resultsByParams struct {
	anyCount  int
	anyParams uint64
	results   map[moqLoadTypesFn_paramsKey]*moqLoadTypesFn_results
}

// moqLoadTypesFn_doFn defines the type of function needed when calling andDo for the LoadTypesFn type
type moqLoadTypesFn_doFn func(pkg string, loadTestTypes bool)

// moqLoadTypesFn_doReturnFn defines the type of function needed when calling doReturnResults for the LoadTypesFn type
type moqLoadTypesFn_doReturnFn func(pkg string, loadTestTypes bool) (
	typeSpecs []*dst.TypeSpec, pkgPath string, err error)

// moqLoadTypesFn_results holds the results of the LoadTypesFn type
type moqLoadTypesFn_results struct {
	params  moqLoadTypesFn_params
	results []struct {
		values *struct {
			typeSpecs []*dst.TypeSpec
			pkgPath   string
			err       error
		}
		sequence   uint32
		doFn       moqLoadTypesFn_doFn
		doReturnFn moqLoadTypesFn_doReturnFn
	}
	index    uint32
	anyTimes bool
}

// moqLoadTypesFn_fnRecorder routes recorded function calls to the moqLoadTypesFn moq
type moqLoadTypesFn_fnRecorder struct {
	params    moqLoadTypesFn_params
	paramsKey moqLoadTypesFn_paramsKey
	anyParams uint64
	sequence  bool
	results   *moqLoadTypesFn_results
	moq       *moqLoadTypesFn
}

// moqLoadTypesFn_anyParams isolates the any params functions of the LoadTypesFn type
type moqLoadTypesFn_anyParams struct {
	recorder *moqLoadTypesFn_fnRecorder
}

// newMoqLoadTypesFn creates a new moq of the LoadTypesFn type
func newMoqLoadTypesFn(scene *moq.Scene, config *moq.Config) *moqLoadTypesFn {
	if config == nil {
		config = &moq.Config{}
	}
	m := &moqLoadTypesFn{
		scene:  scene,
		config: *config,
	}
	scene.AddMoq(m)
	return m
}

// mock returns the moq implementation of the LoadTypesFn type
func (m *moqLoadTypesFn) mock() ast.LoadTypesFn {
	return func(pkg string, loadTestTypes bool) (
		typeSpecs []*dst.TypeSpec, pkgPath string, err error) {
		moq := &moqLoadTypesFn_mock{moq: m}
		return moq.fn(pkg, loadTestTypes)
	}
}

func (m *moqLoadTypesFn_mock) fn(pkg string, loadTestTypes bool) (
	typeSpecs []*dst.TypeSpec, pkgPath string, err error) {
	params := moqLoadTypesFn_params{
		pkg:           pkg,
		loadTestTypes: loadTestTypes,
	}
	var results *moqLoadTypesFn_results
	for _, resultsByParams := range m.moq.resultsByParams {
		var pkgUsed string
		if resultsByParams.anyParams&(1<<0) == 0 {
			pkgUsed = pkg
		}
		var loadTestTypesUsed bool
		if resultsByParams.anyParams&(1<<1) == 0 {
			loadTestTypesUsed = loadTestTypes
		}
		paramsKey := moqLoadTypesFn_paramsKey{
			pkg:           pkgUsed,
			loadTestTypes: loadTestTypesUsed,
		}
		var ok bool
		results, ok = resultsByParams.results[paramsKey]
		if ok {
			break
		}
	}
	if results == nil {
		if m.moq.config.Expectation == moq.Strict {
			m.moq.scene.T.Fatalf("Unexpected call with parameters %#v", params)
		}
		return
	}

	i := int(atomic.AddUint32(&results.index, 1)) - 1
	if i >= len(results.results) {
		if !results.anyTimes {
			if m.moq.config.Expectation == moq.Strict {
				m.moq.scene.T.Fatalf("Too many calls to mock with parameters %#v", params)
			}
			return
		}
		i = len(results.results) - 1
	}

	result := results.results[i]
	if result.sequence != 0 {
		sequence := m.moq.scene.NextMockSequence()
		if (!results.anyTimes && result.sequence != sequence) || result.sequence > sequence {
			m.moq.scene.T.Fatalf("Call sequence does not match %#v", params)
		}
	}

	if result.doFn != nil {
		result.doFn(pkg, loadTestTypes)
	}

	if result.values != nil {
		typeSpecs = result.values.typeSpecs
		pkgPath = result.values.pkgPath
		err = result.values.err
	}
	if result.doReturnFn != nil {
		typeSpecs, pkgPath, err = result.doReturnFn(pkg, loadTestTypes)
	}
	return
}

func (m *moqLoadTypesFn) onCall(pkg string, loadTestTypes bool) *moqLoadTypesFn_fnRecorder {
	return &moqLoadTypesFn_fnRecorder{
		params: moqLoadTypesFn_params{
			pkg:           pkg,
			loadTestTypes: loadTestTypes,
		},
		paramsKey: moqLoadTypesFn_paramsKey{
			pkg:           pkg,
			loadTestTypes: loadTestTypes,
		},
		sequence: m.config.Sequence == moq.SeqDefaultOn,
		moq:      m,
	}
}

func (r *moqLoadTypesFn_fnRecorder) any() *moqLoadTypesFn_anyParams {
	if r.results != nil {
		r.moq.scene.T.Fatalf("Any functions must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	return &moqLoadTypesFn_anyParams{recorder: r}
}

func (a *moqLoadTypesFn_anyParams) pkg() *moqLoadTypesFn_fnRecorder {
	a.recorder.anyParams |= 1 << 0
	return a.recorder
}

func (a *moqLoadTypesFn_anyParams) loadTestTypes() *moqLoadTypesFn_fnRecorder {
	a.recorder.anyParams |= 1 << 1
	return a.recorder
}

func (r *moqLoadTypesFn_fnRecorder) seq() *moqLoadTypesFn_fnRecorder {
	if r.results != nil {
		r.moq.scene.T.Fatalf("seq must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	r.sequence = true
	return r
}

func (r *moqLoadTypesFn_fnRecorder) noSeq() *moqLoadTypesFn_fnRecorder {
	if r.results != nil {
		r.moq.scene.T.Fatalf("noSeq must be called before returnResults or doReturnResults calls, parameters: %#v", r.params)
		return nil
	}
	r.sequence = false
	return r
}

func (r *moqLoadTypesFn_fnRecorder) returnResults(
	typeSpecs []*dst.TypeSpec, pkgPath string, err error) *moqLoadTypesFn_fnRecorder {
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			typeSpecs []*dst.TypeSpec
			pkgPath   string
			err       error
		}
		sequence   uint32
		doFn       moqLoadTypesFn_doFn
		doReturnFn moqLoadTypesFn_doReturnFn
	}{
		values: &struct {
			typeSpecs []*dst.TypeSpec
			pkgPath   string
			err       error
		}{
			typeSpecs: typeSpecs,
			pkgPath:   pkgPath,
			err:       err,
		},
		sequence: sequence,
	})
	return r
}

func (r *moqLoadTypesFn_fnRecorder) andDo(fn moqLoadTypesFn_doFn) *moqLoadTypesFn_fnRecorder {
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults must be called before calling andDo")
		return nil
	}
	last := &r.results.results[len(r.results.results)-1]
	last.doFn = fn
	return r
}

func (r *moqLoadTypesFn_fnRecorder) doReturnResults(fn moqLoadTypesFn_doReturnFn) *moqLoadTypesFn_fnRecorder {
	r.findResults()

	var sequence uint32
	if r.sequence {
		sequence = r.moq.scene.NextRecorderSequence()
	}

	r.results.results = append(r.results.results, struct {
		values *struct {
			typeSpecs []*dst.TypeSpec
			pkgPath   string
			err       error
		}
		sequence   uint32
		doFn       moqLoadTypesFn_doFn
		doReturnFn moqLoadTypesFn_doReturnFn
	}{sequence: sequence, doReturnFn: fn})
	return r
}

func (r *moqLoadTypesFn_fnRecorder) findResults() {
	if r.results == nil {
		anyCount := bits.OnesCount64(r.anyParams)
		insertAt := -1
		var results *moqLoadTypesFn_resultsByParams
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
			results = &moqLoadTypesFn_resultsByParams{
				anyCount:  anyCount,
				anyParams: r.anyParams,
				results:   map[moqLoadTypesFn_paramsKey]*moqLoadTypesFn_results{},
			}
			r.moq.resultsByParams = append(r.moq.resultsByParams, *results)
			if insertAt != -1 && insertAt+1 < len(r.moq.resultsByParams) {
				copy(r.moq.resultsByParams[insertAt+1:], r.moq.resultsByParams[insertAt:0])
				r.moq.resultsByParams[insertAt] = *results
			}
		}

		var pkgUsed string
		if r.anyParams&(1<<0) == 0 {
			pkgUsed = r.paramsKey.pkg
		}
		var loadTestTypesUsed bool
		if r.anyParams&(1<<1) == 0 {
			loadTestTypesUsed = r.paramsKey.loadTestTypes
		}
		paramsKey := moqLoadTypesFn_paramsKey{
			pkg:           pkgUsed,
			loadTestTypes: loadTestTypesUsed,
		}

		var ok bool
		r.results, ok = results.results[paramsKey]
		if !ok {
			r.results = &moqLoadTypesFn_results{
				params:   r.params,
				results:  nil,
				index:    0,
				anyTimes: false,
			}
			results.results[paramsKey] = r.results
		}
	}
}

func (r *moqLoadTypesFn_fnRecorder) repeat(repeaters ...moq.Repeater) *moqLoadTypesFn_fnRecorder {
	if r.results == nil {
		r.moq.scene.T.Fatalf("returnResults or doReturnResults must be called before calling repeat")
		return nil
	}
	repeat := moq.Repeat(r.moq.scene.T, repeaters)
	last := r.results.results[len(r.results.results)-1]
	for n := 0; n < repeat.MaxTimes-1; n++ {
		if last.sequence != 0 {
			last = struct {
				values *struct {
					typeSpecs []*dst.TypeSpec
					pkgPath   string
					err       error
				}
				sequence   uint32
				doFn       moqLoadTypesFn_doFn
				doReturnFn moqLoadTypesFn_doReturnFn
			}{
				values: &struct {
					typeSpecs []*dst.TypeSpec
					pkgPath   string
					err       error
				}{
					typeSpecs: last.values.typeSpecs,
					pkgPath:   last.values.pkgPath,
					err:       last.values.err,
				},
				sequence: r.moq.scene.NextRecorderSequence(),
			}
		}
		r.results.results = append(r.results.results, last)
	}
	r.results.anyTimes = repeat.AnyTimes
	return r
}

// Reset resets the state of the moq
func (m *moqLoadTypesFn) Reset() { m.resultsByParams = nil }

// AssertExpectationsMet asserts that all expectations have been met
func (m *moqLoadTypesFn) AssertExpectationsMet() {
	for _, res := range m.resultsByParams {
		for _, results := range res.results {
			missing := len(results.results) - int(atomic.LoadUint32(&results.index))
			if missing == 1 && results.anyTimes == true {
				continue
			}
			if missing > 0 {
				m.scene.T.Errorf("Expected %d additional call(s) with parameters %#v", missing, results.params)
			}
		}
	}
}
