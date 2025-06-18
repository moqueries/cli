package testmoqs_test

import (
	"errors"
	"fmt"
	"io"
	"reflect"

	"moqueries.org/runtime/moq"

	"moqueries.org/cli/generator/testmoqs"
	"moqueries.org/cli/generator/testmoqs/exported"
)

type usualFnAdaptor struct {
	m *moqUsualFn
}

func (a *usualFnAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *usualFnAdaptor) mock() interface{} { return a.m.mock() }

func (a *usualFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &usualFnRecorder{r: a.m.onCall(sParams[0], bParam)}
}

func (a *usualFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *usualFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("UsualFn(%#v, %#v)", sParams[0], bParam)
}

func (a *usualFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type usualFnRecorder struct {
	r *moqUsualFn_recorder
}

func (r *usualFnRecorder) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.sParam()
	}
}

func (r *usualFnRecorder) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.bParam()
	}
}

func (r *usualFnRecorder) seq() {
	r.r = r.r.seq()
}

func (r *usualFnRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *usualFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *usualFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *usualFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.doReturnResults(func(sParam string, bParam bool) (string, error) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return sResults[0], err
	})
}

func (r *usualFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *usualFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedUsualFnAdaptor struct {
	m *exported.MoqUsualFn
}

func (a *exportedUsualFnAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedUsualFnAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedUsualFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedUsualFnRecorder{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedUsualFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedUsualFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("UsualFn(%#v, %#v)", sParams[0], bParam)
}

func (a *exportedUsualFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedUsualFnRecorder struct {
	r *exported.MoqUsualFn_recorder
}

func (r *exportedUsualFnRecorder) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.SParam()
	}
}

func (r *exportedUsualFnRecorder) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.BParam()
	}
}

func (r *exportedUsualFnRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedUsualFnRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedUsualFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedUsualFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.AndDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *exportedUsualFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.DoReturnResults(func(sParam string, bParam bool) (string, error) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return sResults[0], err
	})
}

func (r *exportedUsualFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedUsualFnRecorder) isNil() bool {
	return r.r == nil
}

type noNamesFnAdaptor struct {
	m *moqNoNamesFn
}

func (a *noNamesFnAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *noNamesFnAdaptor) mock() interface{} { return a.m.mock() }

func (a *noNamesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &noNamesFnRecorder{r: a.m.onCall(sParams[0], bParam)}
}

func (a *noNamesFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *noNamesFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("NoNamesFn(%#v, %#v)", sParams[0], bParam)
}

func (a *noNamesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type noNamesFnRecorder struct {
	r *moqNoNamesFn_recorder
}

func (r *noNamesFnRecorder) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.param1()
	}
}

func (r *noNamesFnRecorder) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.param2()
	}
}

func (r *noNamesFnRecorder) seq() {
	r.r = r.r.seq()
}

func (r *noNamesFnRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *noNamesFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *noNamesFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *noNamesFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.doReturnResults(func(sParam string, bParam bool) (string, error) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return sResults[0], err
	})
}

func (r *noNamesFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *noNamesFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedNoNamesFnAdaptor struct {
	m *exported.MoqNoNamesFn
}

func (a *exportedNoNamesFnAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedNoNamesFnAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedNoNamesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedNoNamesFnRecorder{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedNoNamesFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedNoNamesFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("NoNamesFn(%#v, %#v)", sParams[0], bParam)
}

func (a *exportedNoNamesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedNoNamesFnRecorder struct {
	r *exported.MoqNoNamesFn_recorder
}

func (r *exportedNoNamesFnRecorder) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Param1()
	}
}

func (r *exportedNoNamesFnRecorder) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Param2()
	}
}

func (r *exportedNoNamesFnRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedNoNamesFnRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedNoNamesFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedNoNamesFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.AndDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *exportedNoNamesFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.DoReturnResults(func(sParam string, bParam bool) (string, error) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return sResults[0], err
	})
}

func (r *exportedNoNamesFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedNoNamesFnRecorder) isNil() bool {
	return r.r == nil
}

type noResultsFnAdaptor struct {
	m *moqNoResultsFn
}

func (a *noResultsFnAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *noResultsFnAdaptor) mock() interface{} { return a.m.mock() }

func (a *noResultsFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &noResultsFnRecorder{r: a.m.onCall(sParams[0], bParam)}
}

func (a *noResultsFnAdaptor) invokeMockAndExpectResults(_ moq.T, sParams []string, bParam bool, _ results) {
	a.m.mock()(sParams[0], bParam)
}

func (a *noResultsFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("NoResultsFn(%#v, %#v)", sParams[0], bParam)
}

func (a *noResultsFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type noResultsFnRecorder struct {
	r *moqNoResultsFn_recorder
}

func (r *noResultsFnRecorder) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.sParam()
	}
}

func (r *noResultsFnRecorder) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.bParam()
	}
}

func (r *noResultsFnRecorder) seq() {
	r.r = r.r.seq()
}

func (r *noResultsFnRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *noResultsFnRecorder) returnResults([]string, error) {
	r.r = r.r.returnResults()
}

func (r *noResultsFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *noResultsFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, _ []string, _ error,
) {
	r.r = r.r.doReturnResults(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *noResultsFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *noResultsFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedNoResultsFnAdaptor struct {
	m *exported.MoqNoResultsFn
}

func (a *exportedNoResultsFnAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedNoResultsFnAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedNoResultsFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedNoResultsFnRecorder{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedNoResultsFnAdaptor) invokeMockAndExpectResults(_ moq.T, sParams []string, bParam bool, _ results) {
	a.m.Mock()(sParams[0], bParam)
}

func (a *exportedNoResultsFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("NoResultsFn(%#v, %#v)", sParams[0], bParam)
}

func (a *exportedNoResultsFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedNoResultsFnRecorder struct {
	r *exported.MoqNoResultsFn_recorder
}

func (r *exportedNoResultsFnRecorder) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.SParam()
	}
}

func (r *exportedNoResultsFnRecorder) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.BParam()
	}
}

func (r *exportedNoResultsFnRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedNoResultsFnRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedNoResultsFnRecorder) returnResults([]string, error) {
	r.r = r.r.ReturnResults()
}

func (r *exportedNoResultsFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.AndDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *exportedNoResultsFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, _ []string, _ error,
) {
	r.r = r.r.DoReturnResults(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *exportedNoResultsFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedNoResultsFnRecorder) isNil() bool {
	return r.r == nil
}

type noParamsFnAdaptor struct {
	m *moqNoParamsFn
}

func (a *noParamsFnAdaptor) config() adaptorConfig {
	return adaptorConfig{noParams: true}
}

func (a *noParamsFnAdaptor) mock() interface{} { return a.m.mock() }

func (a *noParamsFnAdaptor) newRecorder([]string, bool) recorder {
	return &noParamsFnRecorder{r: a.m.onCall()}
}

func (a *noParamsFnAdaptor) invokeMockAndExpectResults(t moq.T, _ []string, _ bool, res results) {
	sResult, err := a.m.mock()()
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *noParamsFnAdaptor) prettyParams([]string, bool) string {
	return "NoParamsFn()"
}

func (a *noParamsFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type noParamsFnRecorder struct {
	r *moqNoParamsFn_recorder
}

func (r *noParamsFnRecorder) anySParam() {}

func (r *noParamsFnRecorder) anyBParam() {}

func (r *noParamsFnRecorder) seq() {
	r.r = r.r.seq()
}

func (r *noParamsFnRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *noParamsFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *noParamsFnRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.andDo(func() {
		fn()
	})
}

func (r *noParamsFnRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, sResults []string, err error,
) {
	r.r = r.r.doReturnResults(func() (string, error) {
		fn()
		return sResults[0], err
	})
}

func (r *noParamsFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *noParamsFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedNoParamsFnAdaptor struct {
	m *exported.MoqNoParamsFn
}

func (a *exportedNoParamsFnAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true, noParams: true}
}

func (a *exportedNoParamsFnAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedNoParamsFnAdaptor) newRecorder([]string, bool) recorder {
	return &exportedNoParamsFnRecorder{r: a.m.OnCall()}
}

func (a *exportedNoParamsFnAdaptor) invokeMockAndExpectResults(t moq.T, _ []string, _ bool, res results) {
	sResult, err := a.m.Mock()()
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedNoParamsFnAdaptor) prettyParams([]string, bool) string {
	return "NoParamsFn()"
}

func (a *exportedNoParamsFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedNoParamsFnRecorder struct {
	r *exported.MoqNoParamsFn_recorder
}

func (r *exportedNoParamsFnRecorder) anySParam() {}

func (r *exportedNoParamsFnRecorder) anyBParam() {}

func (r *exportedNoParamsFnRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedNoParamsFnRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedNoParamsFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedNoParamsFnRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.AndDo(func() {
		fn()
	})
}

func (r *exportedNoParamsFnRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, sResults []string, err error,
) {
	r.r = r.r.DoReturnResults(func() (string, error) {
		fn()
		return sResults[0], err
	})
}

func (r *exportedNoParamsFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedNoParamsFnRecorder) isNil() bool {
	return r.r == nil
}

type nothingFnAdaptor struct {
	m *moqNothingFn
}

func (a *nothingFnAdaptor) config() adaptorConfig {
	return adaptorConfig{noParams: true}
}

func (a *nothingFnAdaptor) mock() interface{} { return a.m.mock() }

func (a *nothingFnAdaptor) newRecorder([]string, bool) recorder {
	return &nothingFnRecorder{r: a.m.onCall()}
}

func (a *nothingFnAdaptor) invokeMockAndExpectResults(moq.T, []string, bool, results) {
	a.m.mock()()
}

func (a *nothingFnAdaptor) prettyParams([]string, bool) string {
	return "NothingFn()"
}

func (a *nothingFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type nothingFnRecorder struct {
	r *moqNothingFn_recorder
}

func (r *nothingFnRecorder) anySParam() {}

func (r *nothingFnRecorder) anyBParam() {}

func (r *nothingFnRecorder) seq() {
	r.r = r.r.seq()
}

func (r *nothingFnRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *nothingFnRecorder) returnResults([]string, error) {
	r.r = r.r.returnResults()
}

func (r *nothingFnRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.andDo(func() {
		fn()
	})
}

func (r *nothingFnRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, _ []string, _ error,
) {
	r.r = r.r.doReturnResults(func() {
		fn()
	})
}

func (r *nothingFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *nothingFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedNothingFnAdaptor struct {
	m *exported.MoqNothingFn
}

func (a *exportedNothingFnAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true, noParams: true}
}

func (a *exportedNothingFnAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedNothingFnAdaptor) newRecorder([]string, bool) recorder {
	return &exportedNothingFnRecorder{r: a.m.OnCall()}
}

func (a *exportedNothingFnAdaptor) invokeMockAndExpectResults(moq.T, []string, bool, results) {
	a.m.Mock()()
}

func (a *exportedNothingFnAdaptor) prettyParams([]string, bool) string {
	return "NothingFn()"
}

func (a *exportedNothingFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedNothingFnRecorder struct {
	r *exported.MoqNothingFn_recorder
}

func (r *exportedNothingFnRecorder) anySParam() {}

func (r *exportedNothingFnRecorder) anyBParam() {}

func (r *exportedNothingFnRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedNothingFnRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedNothingFnRecorder) returnResults([]string, error) {
	r.r = r.r.ReturnResults()
}

func (r *exportedNothingFnRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.AndDo(func() {
		fn()
	})
}

func (r *exportedNothingFnRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, _ []string, _ error,
) {
	r.r = r.r.DoReturnResults(func() {
		fn()
	})
}

func (r *exportedNothingFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedNothingFnRecorder) isNil() bool {
	return r.r == nil
}

type variadicFnAdaptor struct {
	m *moqVariadicFn
}

func (a *variadicFnAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *variadicFnAdaptor) mock() interface{} { return a.m.mock() }

func (a *variadicFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &variadicFnRecorder{r: a.m.onCall(bParam, sParams...)}
}

func (a *variadicFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock()(bParam, sParams...)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *variadicFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("VariadicFn(%#v, %#v)", bParam, sParams)
}

func (a *variadicFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type variadicFnRecorder struct {
	r *moqVariadicFn_recorder
}

func (r *variadicFnRecorder) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.args()
	}
}

func (r *variadicFnRecorder) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.other()
	}
}

func (r *variadicFnRecorder) seq() {
	r.r = r.r.seq()
}

func (r *variadicFnRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *variadicFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *variadicFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(other bool, args ...string) {
		fn()
		if !reflect.DeepEqual(args, expectedSParams) {
			t.Errorf("wanted %#v, got %#v", expectedSParams, args)
		}
		if other != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, other)
		}
	})
}

func (r *variadicFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.doReturnResults(func(other bool, args ...string) (string, error) {
		fn()
		if !reflect.DeepEqual(args, expectedSParams) {
			t.Errorf("wanted %#v, got %#v", expectedSParams, args)
		}
		if other != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, other)
		}
		return sResults[0], err
	})
}

func (r *variadicFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *variadicFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedVariadicFnAdaptor struct {
	m *exported.MoqVariadicFn
}

func (a *exportedVariadicFnAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedVariadicFnAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedVariadicFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedVariadicFnRecorder{r: a.m.OnCall(bParam, sParams...)}
}

func (a *exportedVariadicFnAdaptor) invokeMockAndExpectResults(
	t moq.T, sParams []string, bParam bool, res results,
) {
	sResult, err := a.m.Mock()(bParam, sParams...)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedVariadicFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("VariadicFn(%#v, %#v)", bParam, sParams)
}

func (a *exportedVariadicFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedVariadicFnRecorder struct {
	r *exported.MoqVariadicFn_recorder
}

func (r *exportedVariadicFnRecorder) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Args()
	}
}

func (r *exportedVariadicFnRecorder) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Other()
	}
}

func (r *exportedVariadicFnRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedVariadicFnRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedVariadicFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedVariadicFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.AndDo(func(other bool, args ...string) {
		fn()
		if !reflect.DeepEqual(args, expectedSParams) {
			t.Errorf("wanted %#v, got %#v", expectedSParams, args)
		}
		if other != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, other)
		}
	})
}

func (r *exportedVariadicFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.DoReturnResults(func(other bool, args ...string) (string, error) {
		fn()
		if !reflect.DeepEqual(args, expectedSParams) {
			t.Errorf("wanted %#v, got %#v", expectedSParams, args)
		}
		if other != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, other)
		}
		return sResults[0], err
	})
}

func (r *exportedVariadicFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedVariadicFnRecorder) isNil() bool {
	return r.r == nil
}

type repeatedIdsFnAdaptor struct {
	m *moqRepeatedIdsFn
}

func (a *repeatedIdsFnAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *repeatedIdsFnAdaptor) mock() interface{} { return a.m.mock() }

func (a *repeatedIdsFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &repeatedIdsFnRecorder{r: a.m.onCall(sParams[0], sParams[1], bParam)}
}

func (a *repeatedIdsFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult1, sResult2, err := a.m.mock()(sParams[0], sParams[1], bParam)
	if sResult1 != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult1)
	}
	if sResult2 != res.sResults[1] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult2)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *repeatedIdsFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("RepeatedIdsFn(%#v, %#v, %#v)", sParams[0], sParams[1], bParam)
}

func (a *repeatedIdsFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type repeatedIdsFnRecorder struct {
	r *moqRepeatedIdsFn_recorder
}

func (r *repeatedIdsFnRecorder) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.sParam1()
	}
}

func (r *repeatedIdsFnRecorder) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.bParam()
	}
}

func (r *repeatedIdsFnRecorder) seq() {
	r.r = r.r.seq()
}

func (r *repeatedIdsFnRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *repeatedIdsFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], sResults[1], err)
}

func (r *repeatedIdsFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(sParam1, sParam2 string, bParam bool) {
		fn()
		if sParam1 != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam1)
		}
		if sParam2 != expectedSParams[1] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[1], sParam2)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *repeatedIdsFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.doReturnResults(func(sParam1, sParam2 string, bParam bool) (string, string, error) {
		fn()
		if sParam1 != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam1)
		}
		if sParam2 != expectedSParams[1] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[1], sParam2)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return sResults[0], sResults[1], err
	})
}

func (r *repeatedIdsFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *repeatedIdsFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedRepeatedIdsFnAdaptor struct {
	m *exported.MoqRepeatedIdsFn
}

func (a *exportedRepeatedIdsFnAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedRepeatedIdsFnAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedRepeatedIdsFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedRepeatedIdsFnRecorder{r: a.m.OnCall(sParams[0], sParams[1], bParam)}
}

func (a *exportedRepeatedIdsFnAdaptor) invokeMockAndExpectResults(
	t moq.T, sParams []string, bParam bool, res results,
) {
	sResult1, sResult2, err := a.m.Mock()(sParams[0], sParams[1], bParam)
	if sResult1 != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult1)
	}
	if sResult2 != res.sResults[1] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult2)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedRepeatedIdsFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("RepeatedIdsFn(%#v, %#v, %#v)", sParams[0], sParams[1], bParam)
}

func (a *exportedRepeatedIdsFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedRepeatedIdsFnRecorder struct {
	r *exported.MoqRepeatedIdsFn_recorder
}

func (r *exportedRepeatedIdsFnRecorder) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.SParam1()
	}
}

func (r *exportedRepeatedIdsFnRecorder) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.BParam()
	}
}

func (r *exportedRepeatedIdsFnRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedRepeatedIdsFnRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedRepeatedIdsFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], sResults[1], err)
}

func (r *exportedRepeatedIdsFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.AndDo(func(sParam1, sParam2 string, bParam bool) {
		fn()
		if sParam1 != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam1)
		}
		if sParam2 != expectedSParams[1] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[1], sParam2)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *exportedRepeatedIdsFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.DoReturnResults(func(sParam1, sParam2 string, bParam bool) (string, string, error) {
		fn()
		if sParam1 != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam1)
		}
		if sParam2 != expectedSParams[1] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[1], sParam2)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return sResults[0], sResults[1], err
	})
}

func (r *exportedRepeatedIdsFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedRepeatedIdsFnRecorder) isNil() bool {
	return r.r == nil
}

type timesFnAdaptor struct {
	m *moqTimesFn
}

func (a *timesFnAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *timesFnAdaptor) mock() interface{} { return a.m.mock() }

func (a *timesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &timesFnRecorder{r: a.m.onCall(sParams[0], bParam)}
}

func (a *timesFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *timesFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("TimesFn(%#v, %#v)", sParams[0], bParam)
}

func (a *timesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type timesFnRecorder struct {
	r *moqTimesFn_recorder
}

func (r *timesFnRecorder) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.times()
	}
}

func (r *timesFnRecorder) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.bParam()
	}
}

func (r *timesFnRecorder) seq() {
	r.r = r.r.seq()
}

func (r *timesFnRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *timesFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *timesFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *timesFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.doReturnResults(func(sParam string, bParam bool) (string, error) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return sResults[0], err
	})
}

func (r *timesFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *timesFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedTimesFnAdaptor struct {
	m *exported.MoqTimesFn
}

func (a *exportedTimesFnAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedTimesFnAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedTimesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedTimesFnRecorder{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedTimesFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedTimesFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("TimesFn(%#v, %#v)", sParams[0], bParam)
}

func (a *exportedTimesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedTimesFnRecorder struct {
	r *exported.MoqTimesFn_recorder
}

func (r *exportedTimesFnRecorder) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Times()
	}
}

func (r *exportedTimesFnRecorder) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.BParam()
	}
}

func (r *exportedTimesFnRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedTimesFnRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedTimesFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedTimesFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.AndDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *exportedTimesFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.DoReturnResults(func(sParam string, bParam bool) (string, error) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return sResults[0], err
	})
}

func (r *exportedTimesFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedTimesFnRecorder) isNil() bool {
	return r.r == nil
}

type difficultParamNamesFnAdaptor struct {
	m *moqDifficultParamNamesFn
}

func (a *difficultParamNamesFnAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *difficultParamNamesFnAdaptor) mock() interface{} { return a.m.mock() }

func (a *difficultParamNamesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &difficultParamNamesFnRecorder{r: a.m.onCall(bParam, false, sParams[0], 0, 0, 0, 0.0, 0.0, 0.0)}
}

func (a *difficultParamNamesFnAdaptor) invokeMockAndExpectResults(_ moq.T, sParams []string, bParam bool, _ results) {
	a.m.mock()(bParam, false, sParams[0], 0, 0, 0, 0.0, 0.0, 0.0)
}

func (a *difficultParamNamesFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("DifficultParamNamesFn(%#v, false, %#v, 0, 0, 0, 0, 0, 0)", bParam, sParams[0])
}

func (a *difficultParamNamesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type difficultParamNamesFnRecorder struct {
	r *moqDifficultParamNamesFn_recorder
}

func (r *difficultParamNamesFnRecorder) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.param3()
	}
}

func (r *difficultParamNamesFnRecorder) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.param1()
	}
}

func (r *difficultParamNamesFnRecorder) seq() {
	r.r = r.r.seq()
}

func (r *difficultParamNamesFnRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *difficultParamNamesFnRecorder) returnResults([]string, error) {
	r.r = r.r.returnResults()
}

func (r *difficultParamNamesFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(m, _ bool, sequence string, _, _, _ int, _, _, _ float32) {
		fn()
		if sequence != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sequence)
		}
		if m != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, m)
		}
	})
}

func (r *difficultParamNamesFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, _ []string, _ error,
) {
	r.r = r.r.doReturnResults(func(m, _ bool, sequence string, _, _, _ int, _, _, _ float32) {
		fn()
		if sequence != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sequence)
		}
		if m != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, m)
		}
	})
}

func (r *difficultParamNamesFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *difficultParamNamesFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedDifficultParamNamesFnAdaptor struct {
	m *exported.MoqDifficultParamNamesFn
}

func (a *exportedDifficultParamNamesFnAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedDifficultParamNamesFnAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedDifficultParamNamesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedDifficultParamNamesFnRecorder{
		r: a.m.OnCall(bParam, false, sParams[0], 0, 0, 0, 0.0, 0.0, 0.0),
	}
}

func (a *exportedDifficultParamNamesFnAdaptor) invokeMockAndExpectResults(
	_ moq.T, sParams []string, bParam bool, _ results,
) {
	a.m.Mock()(bParam, false, sParams[0], 0, 0, 0, 0.0, 0.0, 0.0)
}

func (a *exportedDifficultParamNamesFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("DifficultParamNamesFn(%#v, false, %#v, 0, 0, 0, 0, 0, 0)", bParam, sParams[0])
}

func (a *exportedDifficultParamNamesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedDifficultParamNamesFnRecorder struct {
	r *exported.MoqDifficultParamNamesFn_recorder
}

func (r *exportedDifficultParamNamesFnRecorder) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Param3()
	}
}

func (r *exportedDifficultParamNamesFnRecorder) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Param1()
	}
}

func (r *exportedDifficultParamNamesFnRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedDifficultParamNamesFnRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedDifficultParamNamesFnRecorder) returnResults([]string, error) {
	r.r = r.r.ReturnResults()
}

func (r *exportedDifficultParamNamesFnRecorder) andDo(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool,
) {
	r.r = r.r.AndDo(func(m, _ bool, sequence string, _, _, _ int, _, _, _ float32) {
		fn()
		if sequence != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sequence)
		}
		if m != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, m)
		}
	})
}

func (r *exportedDifficultParamNamesFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, _ []string, _ error,
) {
	r.r = r.r.DoReturnResults(func(m, _ bool, sequence string, _, _, _ int, _, _, _ float32) {
		fn()
		if sequence != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sequence)
		}
		if m != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, m)
		}
	})
}

func (r *exportedDifficultParamNamesFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedDifficultParamNamesFnRecorder) isNil() bool {
	return r.r == nil
}

type difficultResultNamesFnAdaptor struct {
	m *moqDifficultResultNamesFn
}

func (a *difficultResultNamesFnAdaptor) config() adaptorConfig {
	return adaptorConfig{noParams: true}
}

func (a *difficultResultNamesFnAdaptor) mock() interface{} { return a.m.mock() }

func (a *difficultResultNamesFnAdaptor) newRecorder([]string, bool) recorder {
	return &difficultResultNamesFnRecorder{r: a.m.onCall()}
}

func (a *difficultResultNamesFnAdaptor) invokeMockAndExpectResults(t moq.T, _ []string, _ bool, res results) {
	m, r, sequence, _, _, _, _, _, _ := a.m.mock()()
	if m != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], m)
	}
	if r != res.sResults[1] {
		t.Errorf("wanted %#v, got %#v", res.sResults[1], m)
	}
	if !errors.Is(sequence, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, sequence)
	}
}

func (a *difficultResultNamesFnAdaptor) prettyParams([]string, bool) string {
	return "DifficultResultNamesFn()"
}

func (a *difficultResultNamesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type difficultResultNamesFnRecorder struct {
	r *moqDifficultResultNamesFn_recorder
}

func (r *difficultResultNamesFnRecorder) anySParam() {}

func (r *difficultResultNamesFnRecorder) anyBParam() {}

func (r *difficultResultNamesFnRecorder) seq() {
	r.r = r.r.seq()
}

func (r *difficultResultNamesFnRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *difficultResultNamesFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], sResults[1], err, 0, 0, 0, 0.0, 0.0, 0.0)
}

func (r *difficultResultNamesFnRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.andDo(func() {
		fn()
	})
}

func (r *difficultResultNamesFnRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, sResults []string, err error,
) {
	//nolint:staticcheck // doReturnFn functions may have error in middle of params
	r.r = r.r.doReturnResults(func() (string, string, error, int, int, int, float32, float32, float32) {
		fn()
		return sResults[0], sResults[1], err, 0, 0, 0, 0.0, 0.0, 0.0
	})
}

func (r *difficultResultNamesFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *difficultResultNamesFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedDifficultResultNamesFnAdaptor struct {
	m *exported.MoqDifficultResultNamesFn
}

func (a *exportedDifficultResultNamesFnAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true, noParams: true}
}

func (a *exportedDifficultResultNamesFnAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedDifficultResultNamesFnAdaptor) newRecorder([]string, bool) recorder {
	return &exportedDifficultResultNamesFnRecorder{r: a.m.OnCall()}
}

func (a *exportedDifficultResultNamesFnAdaptor) invokeMockAndExpectResults(t moq.T, _ []string, _ bool, res results) {
	m, r, sequence, _, _, _, _, _, _ := a.m.Mock()()
	if m != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], m)
	}
	if r != res.sResults[1] {
		t.Errorf("wanted %#v, got %#v", res.sResults[1], m)
	}
	if !errors.Is(sequence, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, sequence)
	}
}

func (a *exportedDifficultResultNamesFnAdaptor) prettyParams([]string, bool) string {
	return "DifficultResultNamesFn()"
}

func (a *exportedDifficultResultNamesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedDifficultResultNamesFnRecorder struct {
	r *exported.MoqDifficultResultNamesFn_recorder
}

func (r *exportedDifficultResultNamesFnRecorder) anySParam() {}

func (r *exportedDifficultResultNamesFnRecorder) anyBParam() {}

func (r *exportedDifficultResultNamesFnRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedDifficultResultNamesFnRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedDifficultResultNamesFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], sResults[1], err, 0, 0, 0, 0.0, 0.0, 0.0)
}

func (r *exportedDifficultResultNamesFnRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.AndDo(func() {
		fn()
	})
}

func (r *exportedDifficultResultNamesFnRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, sResults []string, err error,
) {
	//nolint:staticcheck // doReturnFn functions may have error in middle of params
	r.r = r.r.DoReturnResults(func() (string, string, error, int, int, int, float32, float32, float32) {
		fn()
		return sResults[0], sResults[1], err, 0, 0, 0, 0.0, 0.0, 0.0
	})
}

func (r *exportedDifficultResultNamesFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedDifficultResultNamesFnRecorder) isNil() bool {
	return r.r == nil
}

type passByReferenceFnAdaptor struct {
	m *moqPassByReferenceFn
}

func (a *passByReferenceFnAdaptor) config() adaptorConfig {
	return adaptorConfig{
		opaqueParams: true,
	}
}

func (a *passByReferenceFnAdaptor) mock() interface{} { return a.m.mock() }

func (a *passByReferenceFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &passByReferenceFnRecorder{r: a.m.onCall(&testmoqs.Params{
		SParam: sParams[0],
		BParam: bParam,
	})}
}

func (a *passByReferenceFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	r := a.m.mock()(&testmoqs.Params{
		SParam: sParams[0],
		BParam: bParam,
	})
	if r == nil {
		if res.sResults[0] != "" || res.err != nil {
			t.Fatalf("got nil, want real results")
		}
		return
	}
	if r.SResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], r.SResult)
	}
	if !errors.Is(r.Err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, r.Err)
	}
}

func (a *passByReferenceFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("PassByReferenceFn(%#v)", &testmoqs.Params{
		SParam: sParams[0],
		BParam: bParam,
	})
}

func (a *passByReferenceFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type passByReferenceFnRecorder struct {
	r *moqPassByReferenceFn_recorder
}

func (r *passByReferenceFnRecorder) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.p()
	}
}

func (r *passByReferenceFnRecorder) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.p()
	}
}

func (r *passByReferenceFnRecorder) seq() {
	r.r = r.r.seq()
}

func (r *passByReferenceFnRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *passByReferenceFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(&testmoqs.Results{SResult: sResults[0], Err: err})
}

func (r *passByReferenceFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(p *testmoqs.Params) {
		fn()
		if p.SParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], p.SParam)
		}
		if p.BParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, p.BParam)
		}
	})
}

func (r *passByReferenceFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.doReturnResults(func(p *testmoqs.Params) *testmoqs.Results {
		fn()
		if p.SParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], p.SParam)
		}
		if p.BParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, p.BParam)
		}
		return &testmoqs.Results{SResult: sResults[0], Err: err}
	})
}

func (r *passByReferenceFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *passByReferenceFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedPassByReferenceFnAdaptor struct {
	m *exported.MoqPassByReferenceFn
}

func (a *exportedPassByReferenceFnAdaptor) config() adaptorConfig {
	return adaptorConfig{
		exported:     true,
		opaqueParams: true,
	}
}

func (a *exportedPassByReferenceFnAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedPassByReferenceFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedPassByReferenceFnRecorder{r: a.m.OnCall(&testmoqs.Params{
		SParam: sParams[0],
		BParam: bParam,
	})}
}

func (a *exportedPassByReferenceFnAdaptor) invokeMockAndExpectResults(
	t moq.T, sParams []string, bParam bool, res results,
) {
	r := a.m.Mock()(&testmoqs.Params{
		SParam: sParams[0],
		BParam: bParam,
	})
	if r == nil {
		if res.sResults[0] != "" || res.err != nil {
			t.Fatalf("got nil, want real results")
		}
		return
	}
	if r.SResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], r.SResult)
	}
	if !errors.Is(r.Err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, r.Err)
	}
}

func (a *exportedPassByReferenceFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("PassByReferenceFn(%#v)", &testmoqs.Params{
		SParam: sParams[0],
		BParam: bParam,
	})
}

func (a *exportedPassByReferenceFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedPassByReferenceFnRecorder struct {
	r *exported.MoqPassByReferenceFn_recorder
}

func (r *exportedPassByReferenceFnRecorder) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.P()
	}
}

func (r *exportedPassByReferenceFnRecorder) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.P()
	}
}

func (r *exportedPassByReferenceFnRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedPassByReferenceFnRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedPassByReferenceFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(&testmoqs.Results{SResult: sResults[0], Err: err})
}

func (r *exportedPassByReferenceFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.AndDo(func(p *testmoqs.Params) {
		fn()
		if p.SParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], p.SParam)
		}
		if p.BParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, p.BParam)
		}
	})
}

func (r *exportedPassByReferenceFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.DoReturnResults(func(p *testmoqs.Params) *testmoqs.Results {
		fn()
		if p.SParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], p.SParam)
		}
		if p.BParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, p.BParam)
		}
		return &testmoqs.Results{SResult: sResults[0], Err: err}
	})
}

func (r *exportedPassByReferenceFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedPassByReferenceFnRecorder) isNil() bool {
	return r.r == nil
}

type interfaceParamFnAdaptor struct {
	m *moqInterfaceParamFn
}

func (a *interfaceParamFnAdaptor) config() adaptorConfig {
	return adaptorConfig{
		opaqueParams: true,
	}
}

func (a *interfaceParamFnAdaptor) mock() interface{} { return a.m.mock() }

func (a *interfaceParamFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &interfaceParamFnRecorder{r: a.m.onCall(&testmoqs.InterfaceParamWriter{
		SParam: sParams[0],
		BParam: bParam,
	})}
}

func (a *interfaceParamFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock()(&testmoqs.InterfaceParamWriter{
		SParam: sParams[0],
		BParam: bParam,
	})
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *interfaceParamFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("InterfaceParamFn(%#v)", &testmoqs.InterfaceParamWriter{
		SParam: sParams[0],
		BParam: bParam,
	})
}

func (a *interfaceParamFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type interfaceParamFnRecorder struct {
	r *moqInterfaceParamFn_recorder
}

func (r *interfaceParamFnRecorder) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.w()
	}
}

func (r *interfaceParamFnRecorder) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.w()
	}
}

func (r *interfaceParamFnRecorder) seq() {
	r.r = r.r.seq()
}

func (r *interfaceParamFnRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *interfaceParamFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *interfaceParamFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(w io.Writer) {
		fn()
		ipw, ok := w.(*testmoqs.InterfaceParamWriter)
		if !ok && w != nil {
			t.Fatalf("wanted a *testmoqs.InterfaceParamWriter, got %#v", w)
		}
		if ipw.SParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], ipw.SParam)
		}
		if ipw.BParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, ipw.BParam)
		}
	})
}

func (r *interfaceParamFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.doReturnResults(func(w io.Writer) (string, error) {
		fn()
		ipw, ok := w.(*testmoqs.InterfaceParamWriter)
		if !ok && w != nil {
			t.Fatalf("wanted a *testmoqs.InterfaceParamWriter, got %#v", w)
		}
		if ipw.SParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], ipw.SParam)
		}
		if ipw.BParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, ipw.BParam)
		}
		return sResults[0], err
	})
}

func (r *interfaceParamFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *interfaceParamFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedInterfaceParamFnAdaptor struct {
	m *exported.MoqInterfaceParamFn
}

func (a *exportedInterfaceParamFnAdaptor) config() adaptorConfig {
	return adaptorConfig{
		exported:     true,
		opaqueParams: true,
	}
}

func (a *exportedInterfaceParamFnAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedInterfaceParamFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedInterfaceParamFnRecorder{r: a.m.OnCall(&testmoqs.InterfaceParamWriter{
		SParam: sParams[0],
		BParam: bParam,
	})}
}

func (a *exportedInterfaceParamFnAdaptor) invokeMockAndExpectResults(
	t moq.T, sParams []string, bParam bool, res results,
) {
	sResult, err := a.m.Mock()(&testmoqs.InterfaceParamWriter{
		SParam: sParams[0],
		BParam: bParam,
	})
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedInterfaceParamFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("InterfaceParamFn(%#v)", &testmoqs.InterfaceParamWriter{
		SParam: sParams[0],
		BParam: bParam,
	})
}

func (a *exportedInterfaceParamFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedInterfaceParamFnRecorder struct {
	r *exported.MoqInterfaceParamFn_recorder
}

func (r *exportedInterfaceParamFnRecorder) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.W()
	}
}

func (r *exportedInterfaceParamFnRecorder) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.W()
	}
}

func (r *exportedInterfaceParamFnRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedInterfaceParamFnRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedInterfaceParamFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedInterfaceParamFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.AndDo(func(w io.Writer) {
		fn()
		ipw, ok := w.(*testmoqs.InterfaceParamWriter)
		if !ok && w != nil {
			t.Fatalf("wanted a *testmoqs.InterfaceParamWriter, got %#v", w)
		}
		if ipw.SParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], ipw.SParam)
		}
		if ipw.BParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, ipw.BParam)
		}
	})
}

func (r *exportedInterfaceParamFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.DoReturnResults(func(w io.Writer) (string, error) {
		fn()
		ipw, ok := w.(*testmoqs.InterfaceParamWriter)
		if !ok && w != nil {
			t.Fatalf("wanted a *testmoqs.InterfaceParamWriter, got %#v", w)
		}
		if ipw.SParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], ipw.SParam)
		}
		if ipw.BParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, ipw.BParam)
		}
		return sResults[0], err
	})
}

func (r *exportedInterfaceParamFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedInterfaceParamFnRecorder) isNil() bool {
	return r.r == nil
}

type interfaceResultFnAdaptor struct {
	m *moqInterfaceResultFn
}

func (a *interfaceResultFnAdaptor) config() adaptorConfig {
	return adaptorConfig{
		opaqueParams: true,
	}
}

func (a *interfaceResultFnAdaptor) mock() interface{} { return a.m.mock() }

func (a *interfaceResultFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &interfaceResultFnRecorder{r: a.m.onCall(sParams[0], bParam)}
}

func (a *interfaceResultFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	r := a.m.mock()(sParams[0], bParam)
	irr, ok := r.(*testmoqs.InterfaceResultReader)
	if !ok {
		if r != nil {
			t.Fatalf("wanted a *testmoqs.InterfaceResultReader, got %#v", r)
		}
		if res.sResults[0] != "" || res.err != nil {
			t.Fatalf("wanted real results, got %#v", r)
		}
		return
	}
	if irr.SResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], irr)
	}
	if !errors.Is(irr.Err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, irr)
	}
}

func (a *interfaceResultFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("InterfaceResultFn(%#v, %#v)", sParams[0], bParam)
}

func (a *interfaceResultFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type interfaceResultFnRecorder struct {
	r *moqInterfaceResultFn_recorder
}

func (r *interfaceResultFnRecorder) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.sParam()
	}
}

func (r *interfaceResultFnRecorder) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.bParam()
	}
}

func (r *interfaceResultFnRecorder) seq() {
	r.r = r.r.seq()
}

func (r *interfaceResultFnRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *interfaceResultFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(&testmoqs.InterfaceResultReader{
		SResult: sResults[0],
		Err:     err,
	})
}

func (r *interfaceResultFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *interfaceResultFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.doReturnResults(func(sParam string, bParam bool) io.Reader {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return &testmoqs.InterfaceResultReader{
			SResult: sResults[0],
			Err:     err,
		}
	})
}

func (r *interfaceResultFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *interfaceResultFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedInterfaceResultFnAdaptor struct {
	m *exported.MoqInterfaceResultFn
}

func (a *exportedInterfaceResultFnAdaptor) config() adaptorConfig {
	return adaptorConfig{
		exported:     true,
		opaqueParams: true,
	}
}

func (a *exportedInterfaceResultFnAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedInterfaceResultFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedInterfaceResultFnRecorder{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedInterfaceResultFnAdaptor) invokeMockAndExpectResults(
	t moq.T, sParams []string, bParam bool, res results,
) {
	r := a.m.Mock()(sParams[0], bParam)
	irr, ok := r.(*testmoqs.InterfaceResultReader)
	if !ok {
		if r != nil {
			t.Fatalf("wanted a *testmoqs.InterfaceResultReader, got %#v", r)
		}
		if res.sResults[0] != "" || res.err != nil {
			t.Fatalf("wanted real results, got %#v", r)
		}
		return
	}
	if irr.SResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], irr)
	}
	if !errors.Is(irr.Err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, irr)
	}
}

func (a *exportedInterfaceResultFnAdaptor) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("InterfaceResultFn(%#v, %#v)", sParams[0], bParam)
}

func (a *exportedInterfaceResultFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedInterfaceResultFnRecorder struct {
	r *exported.MoqInterfaceResultFn_recorder
}

func (r *exportedInterfaceResultFnRecorder) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.SParam()
	}
}

func (r *exportedInterfaceResultFnRecorder) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.BParam()
	}
}

func (r *exportedInterfaceResultFnRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedInterfaceResultFnRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedInterfaceResultFnRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(&testmoqs.InterfaceResultReader{
		SResult: sResults[0],
		Err:     err,
	})
}

func (r *exportedInterfaceResultFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.AndDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *exportedInterfaceResultFnRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.DoReturnResults(func(sParam string, bParam bool) io.Reader {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return &testmoqs.InterfaceResultReader{
			SResult: sResults[0],
			Err:     err,
		}
	})
}

func (r *exportedInterfaceResultFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedInterfaceResultFnRecorder) isNil() bool {
	return r.r == nil
}

type genericParamsFnAdaptor[S, B any] struct {
	m *moqGenericParamsFn[S, B]
}

func (a *genericParamsFnAdaptor[S, B]) config() adaptorConfig { return adaptorConfig{} }

func (a *genericParamsFnAdaptor[S, B]) mock() interface{} { return a.m.mock() }

func (a *genericParamsFnAdaptor[S, B]) newRecorder(sParams []S, bParam B) recorder {
	return &genericParamsFnRecorder[S, B]{r: a.m.onCall(sParams[0], bParam)}
}

func (a *genericParamsFnAdaptor[S, B]) invokeMockAndExpectResults(t moq.T, sParams []S, bParam B, res results) {
	sResult, err := a.m.mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *genericParamsFnAdaptor[S, B]) prettyParams(sParams []S, bParam B) string {
	return fmt.Sprintf("GenericParamsFn(%#v, %#v)", sParams[0], bParam)
}

func (a *genericParamsFnAdaptor[S, B]) sceneMoq() moq.Moq {
	return a.m
}

type genericParamsFnRecorder[S, B any] struct {
	r *moqGenericParamsFn_recorder[S, B]
}

func (r *genericParamsFnRecorder[S, B]) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.param1()
	}
}

func (r *genericParamsFnRecorder[S, B]) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.param2()
	}
}

func (r *genericParamsFnRecorder[S, B]) seq() {
	r.r = r.r.seq()
}

func (r *genericParamsFnRecorder[S, B]) noSeq() {
	r.r = r.r.noSeq()
}

func (r *genericParamsFnRecorder[S, B]) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *genericParamsFnRecorder[S, B]) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(sParam S, bParam B) {
		fn()
		if !reflect.DeepEqual(sParam, expectedSParams[0]) {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if !reflect.DeepEqual(bParam, expectedBParam) {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *genericParamsFnRecorder[S, B]) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.doReturnResults(func(sParam S, bParam B) (string, error) {
		fn()
		if !reflect.DeepEqual(sParam, expectedSParams[0]) {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if !reflect.DeepEqual(bParam, expectedBParam) {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return sResults[0], err
	})
}

func (r *genericParamsFnRecorder[S, B]) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *genericParamsFnRecorder[S, B]) isNil() bool {
	return r.r == nil
}

type exportedGenericParamsFnAdaptor[S, B any] struct {
	m *exported.MoqGenericParamsFn[S, B]
}

func (a *exportedGenericParamsFnAdaptor[S, B]) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedGenericParamsFnAdaptor[S, B]) mock() interface{} { return a.m.Mock() }

func (a *exportedGenericParamsFnAdaptor[S, B]) newRecorder(sParams []S, bParam B) recorder {
	return &exportedGenericParamsFnRecorder[S, B]{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedGenericParamsFnAdaptor[S, B]) invokeMockAndExpectResults(
	t moq.T, sParams []S, bParam B, res results,
) {
	sResult, err := a.m.Mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedGenericParamsFnAdaptor[S, B]) prettyParams(sParams []S, bParam B) string {
	return fmt.Sprintf("GenericParamsFn(%#v, %#v)", sParams[0], bParam)
}

func (a *exportedGenericParamsFnAdaptor[S, B]) sceneMoq() moq.Moq {
	return a.m
}

type exportedGenericParamsFnRecorder[S, B any] struct {
	r *exported.MoqGenericParamsFn_recorder[S, B]
}

func (r *exportedGenericParamsFnRecorder[S, B]) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Param1()
	}
}

func (r *exportedGenericParamsFnRecorder[S, B]) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Param2()
	}
}

func (r *exportedGenericParamsFnRecorder[S, B]) seq() {
	r.r = r.r.Seq()
}

func (r *exportedGenericParamsFnRecorder[S, B]) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedGenericParamsFnRecorder[S, B]) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedGenericParamsFnRecorder[S, B]) andDo(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool,
) {
	r.r = r.r.AndDo(func(sParam S, bParam B) {
		fn()
		if !reflect.DeepEqual(sParam, expectedSParams[0]) {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if !reflect.DeepEqual(bParam, expectedBParam) {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *exportedGenericParamsFnRecorder[S, B]) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.DoReturnResults(func(sParam S, bParam B) (string, error) {
		fn()
		if !reflect.DeepEqual(sParam, expectedSParams[0]) {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if !reflect.DeepEqual(bParam, expectedBParam) {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return sResults[0], err
	})
}

func (r *exportedGenericParamsFnRecorder[S, B]) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedGenericParamsFnRecorder[S, B]) isNil() bool {
	return r.r == nil
}

type partialGenericParamsFnAdaptor[S any] struct {
	m *moqPartialGenericParamsFn[S]
}

func (a *partialGenericParamsFnAdaptor[S]) config() adaptorConfig { return adaptorConfig{} }

func (a *partialGenericParamsFnAdaptor[S]) mock() interface{} { return a.m.mock() }

func (a *partialGenericParamsFnAdaptor[S]) newRecorder(sParams []S, bParam bool) recorder {
	return &partialGenericParamsFnRecorder[S]{r: a.m.onCall(sParams[0], bParam)}
}

func (a *partialGenericParamsFnAdaptor[S]) invokeMockAndExpectResults(
	t moq.T, sParams []S, bParam bool, res results,
) {
	sResult, err := a.m.mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *partialGenericParamsFnAdaptor[S]) prettyParams(sParams []S, bParam bool) string {
	return fmt.Sprintf("PartialGenericParamsFn(%#v, %#v)", sParams[0], bParam)
}

func (a *partialGenericParamsFnAdaptor[S]) sceneMoq() moq.Moq {
	return a.m
}

type partialGenericParamsFnRecorder[S any] struct {
	r *moqPartialGenericParamsFn_recorder[S]
}

func (r *partialGenericParamsFnRecorder[S]) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.param1()
	}
}

func (r *partialGenericParamsFnRecorder[S]) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.param2()
	}
}

func (r *partialGenericParamsFnRecorder[S]) seq() {
	r.r = r.r.seq()
}

func (r *partialGenericParamsFnRecorder[S]) noSeq() {
	r.r = r.r.noSeq()
}

func (r *partialGenericParamsFnRecorder[S]) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *partialGenericParamsFnRecorder[S]) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(sParam S, bParam bool) {
		fn()
		if !reflect.DeepEqual(sParam, expectedSParams[0]) {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *partialGenericParamsFnRecorder[S]) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.doReturnResults(func(sParam S, bParam bool) (string, error) {
		fn()
		if !reflect.DeepEqual(sParam, expectedSParams[0]) {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return sResults[0], err
	})
}

func (r *partialGenericParamsFnRecorder[S]) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *partialGenericParamsFnRecorder[S]) isNil() bool {
	return r.r == nil
}

type exportedPartialGenericParamsFnAdaptor[S any] struct {
	m *exported.MoqPartialGenericParamsFn[S]
}

func (a *exportedPartialGenericParamsFnAdaptor[S]) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedPartialGenericParamsFnAdaptor[S]) mock() interface{} { return a.m.Mock() }

func (a *exportedPartialGenericParamsFnAdaptor[S]) newRecorder(sParams []S, bParam bool) recorder {
	return &exportedPartialGenericParamsFnRecorder[S]{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedPartialGenericParamsFnAdaptor[S]) invokeMockAndExpectResults(
	t moq.T, sParams []S, bParam bool, res results,
) {
	sResult, err := a.m.Mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !errors.Is(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedPartialGenericParamsFnAdaptor[S]) prettyParams(sParams []S, bParam bool) string {
	return fmt.Sprintf("PartialGenericParamsFn(%#v, %#v)", sParams[0], bParam)
}

func (a *exportedPartialGenericParamsFnAdaptor[S]) sceneMoq() moq.Moq {
	return a.m
}

type exportedPartialGenericParamsFnRecorder[S any] struct {
	r *exported.MoqPartialGenericParamsFn_recorder[S]
}

func (r *exportedPartialGenericParamsFnRecorder[S]) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Param1()
	}
}

func (r *exportedPartialGenericParamsFnRecorder[S]) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Param2()
	}
}

func (r *exportedPartialGenericParamsFnRecorder[S]) seq() {
	r.r = r.r.Seq()
}

func (r *exportedPartialGenericParamsFnRecorder[S]) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedPartialGenericParamsFnRecorder[S]) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedPartialGenericParamsFnRecorder[S]) andDo(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool,
) {
	r.r = r.r.AndDo(func(sParam S, bParam bool) {
		fn()
		if !reflect.DeepEqual(sParam, expectedSParams[0]) {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *exportedPartialGenericParamsFnRecorder[S]) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.DoReturnResults(func(sParam S, bParam bool) (string, error) {
		fn()
		if !reflect.DeepEqual(sParam, expectedSParams[0]) {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return sResults[0], err
	})
}

func (r *exportedPartialGenericParamsFnRecorder[S]) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedPartialGenericParamsFnRecorder[S]) isNil() bool {
	return r.r == nil
}

type genericResultsFnAdaptor[S ~string, E error] struct {
	m *moqGenericResultsFn[S, E]
}

func (a *genericResultsFnAdaptor[S, E]) config() adaptorConfig { return adaptorConfig{} }

func (a *genericResultsFnAdaptor[S, E]) mock() interface{} { return a.m.mock() }

func (a *genericResultsFnAdaptor[S, E]) newRecorder(sParams []string, bParam bool) recorder {
	return &genericResultsFnRecorder[S, E]{r: a.m.onCall(sParams[0], bParam)}
}

func (a *genericResultsFnAdaptor[S, E]) invokeMockAndExpectResults(
	t moq.T, sParams []string, bParam bool, res results,
) {
	sResult, err := a.m.mock()(sParams[0], bParam)
	if !reflect.DeepEqual(sResult, res.sResults[0]) {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !reflect.DeepEqual(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *genericResultsFnAdaptor[S, E]) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("GenericResultsFn(%#v, %#v)", sParams[0], bParam)
}

func (a *genericResultsFnAdaptor[S, E]) sceneMoq() moq.Moq {
	return a.m
}

type genericResultsFnRecorder[S ~string, E error] struct {
	r *moqGenericResultsFn_recorder[S, E]
}

func (r *genericResultsFnRecorder[S, E]) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.param1()
	}
}

func (r *genericResultsFnRecorder[S, E]) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.param2()
	}
}

func (r *genericResultsFnRecorder[S, E]) seq() {
	r.r = r.r.seq()
}

func (r *genericResultsFnRecorder[S, E]) noSeq() {
	r.r = r.r.noSeq()
}

func (r *genericResultsFnRecorder[S, E]) returnResults(sResults []string, err error) {
	var e E
	if err != nil {
		e = err.(E)
	}
	r.r = r.r.returnResults(S(sResults[0]), e)
}

func (r *genericResultsFnRecorder[S, E]) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *genericResultsFnRecorder[S, E]) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.doReturnResults(func(sParam string, bParam bool) (S, E) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		var e E
		if err != nil {
			e = err.(E)
		}
		return S(sResults[0]), e
	})
}

func (r *genericResultsFnRecorder[S, E]) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *genericResultsFnRecorder[S, E]) isNil() bool {
	return r.r == nil
}

type exportedGenericResultsFnAdaptor[S ~string, E error] struct {
	m *exported.MoqGenericResultsFn[S, E]
}

func (a *exportedGenericResultsFnAdaptor[S, E]) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedGenericResultsFnAdaptor[S, E]) mock() interface{} { return a.m.Mock() }

func (a *exportedGenericResultsFnAdaptor[S, E]) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedGenericResultsFnRecorder[S, E]{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedGenericResultsFnAdaptor[S, E]) invokeMockAndExpectResults(
	t moq.T, sParams []string, bParam bool, res results,
) {
	sResult, err := a.m.Mock()(sParams[0], bParam)
	if !reflect.DeepEqual(sResult, res.sResults[0]) {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !reflect.DeepEqual(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedGenericResultsFnAdaptor[S, E]) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("GenericResultsFn(%#v, %#v)", sParams[0], bParam)
}

func (a *exportedGenericResultsFnAdaptor[S, E]) sceneMoq() moq.Moq {
	return a.m
}

type exportedGenericResultsFnRecorder[S ~string, E error] struct {
	r *exported.MoqGenericResultsFn_recorder[S, E]
}

func (r *exportedGenericResultsFnRecorder[S, E]) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Param1()
	}
}

func (r *exportedGenericResultsFnRecorder[S, E]) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Param2()
	}
}

func (r *exportedGenericResultsFnRecorder[S, E]) seq() {
	r.r = r.r.Seq()
}

func (r *exportedGenericResultsFnRecorder[S, E]) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedGenericResultsFnRecorder[S, E]) returnResults(sResults []string, err error) {
	var e E
	if err != nil {
		e = err.(E)
	}
	r.r = r.r.ReturnResults(S(sResults[0]), e)
}

func (r *exportedGenericResultsFnRecorder[S, E]) andDo(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool,
) {
	r.r = r.r.AndDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *exportedGenericResultsFnRecorder[S, E]) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.DoReturnResults(func(sParam string, bParam bool) (S, E) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		var e E
		if err != nil {
			e = err.(E)
		}
		return S(sResults[0]), e
	})
}

func (r *exportedGenericResultsFnRecorder[S, E]) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedGenericResultsFnRecorder[S, E]) isNil() bool {
	return r.r == nil
}

type partialGenericResultsFnAdaptor[S ~string] struct {
	m *moqPartialGenericResultsFn[S]
}

func (a *partialGenericResultsFnAdaptor[S]) config() adaptorConfig { return adaptorConfig{} }

func (a *partialGenericResultsFnAdaptor[S]) mock() interface{} { return a.m.mock() }

func (a *partialGenericResultsFnAdaptor[S]) newRecorder(sParams []string, bParam bool) recorder {
	return &partialGenericResultsFnRecorder[S]{r: a.m.onCall(sParams[0], bParam)}
}

func (a *partialGenericResultsFnAdaptor[S]) invokeMockAndExpectResults(
	t moq.T, sParams []string, bParam bool, res results,
) {
	sResult, err := a.m.mock()(sParams[0], bParam)
	if !reflect.DeepEqual(sResult, res.sResults[0]) {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !reflect.DeepEqual(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *partialGenericResultsFnAdaptor[S]) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("PartialGenericResultsFn(%#v, %#v)", sParams[0], bParam)
}

func (a *partialGenericResultsFnAdaptor[S]) sceneMoq() moq.Moq {
	return a.m
}

type partialGenericResultsFnRecorder[S ~string] struct {
	r *moqPartialGenericResultsFn_recorder[S]
}

func (r *partialGenericResultsFnRecorder[S]) anySParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.param1()
	}
}

func (r *partialGenericResultsFnRecorder[S]) anyBParam() {
	if a := r.r.any(); a == nil {
		r.r = nil
	} else {
		r.r = a.param2()
	}
}

func (r *partialGenericResultsFnRecorder[S]) seq() {
	r.r = r.r.seq()
}

func (r *partialGenericResultsFnRecorder[S]) noSeq() {
	r.r = r.r.noSeq()
}

func (r *partialGenericResultsFnRecorder[S]) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(S(sResults[0]), err)
}

func (r *partialGenericResultsFnRecorder[S]) andDo(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool,
) {
	r.r = r.r.andDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *partialGenericResultsFnRecorder[S]) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.doReturnResults(func(sParam string, bParam bool) (S, error) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return S(sResults[0]), err
	})
}

func (r *partialGenericResultsFnRecorder[S]) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *partialGenericResultsFnRecorder[S]) isNil() bool {
	return r.r == nil
}

type exportedPartialGenericResultsFnAdaptor[S ~string] struct {
	m *exported.MoqPartialGenericResultsFn[S]
}

func (a *exportedPartialGenericResultsFnAdaptor[S]) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedPartialGenericResultsFnAdaptor[S]) mock() interface{} { return a.m.Mock() }

func (a *exportedPartialGenericResultsFnAdaptor[S]) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedPartialGenericResultsFnRecorder[S]{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedPartialGenericResultsFnAdaptor[S]) invokeMockAndExpectResults(
	t moq.T, sParams []string, bParam bool, res results,
) {
	sResult, err := a.m.Mock()(sParams[0], bParam)
	if !reflect.DeepEqual(sResult, res.sResults[0]) {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if !reflect.DeepEqual(err, res.err) {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedPartialGenericResultsFnAdaptor[S]) prettyParams(sParams []string, bParam bool) string {
	return fmt.Sprintf("PartialGenericResultsFn(%#v, %#v)", sParams[0], bParam)
}

func (a *exportedPartialGenericResultsFnAdaptor[S]) sceneMoq() moq.Moq {
	return a.m
}

type exportedPartialGenericResultsFnRecorder[S ~string] struct {
	r *exported.MoqPartialGenericResultsFn_recorder[S]
}

func (r *exportedPartialGenericResultsFnRecorder[S]) anySParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Param1()
	}
}

func (r *exportedPartialGenericResultsFnRecorder[S]) anyBParam() {
	if a := r.r.Any(); a == nil {
		r.r = nil
	} else {
		r.r = a.Param2()
	}
}

func (r *exportedPartialGenericResultsFnRecorder[S]) seq() {
	r.r = r.r.Seq()
}

func (r *exportedPartialGenericResultsFnRecorder[S]) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedPartialGenericResultsFnRecorder[S]) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(S(sResults[0]), err)
}

func (r *exportedPartialGenericResultsFnRecorder[S]) andDo(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool,
) {
	r.r = r.r.AndDo(func(sParam string, bParam bool) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
	})
}

func (r *exportedPartialGenericResultsFnRecorder[S]) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error,
) {
	r.r = r.r.DoReturnResults(func(sParam string, bParam bool) (S, error) {
		fn()
		if sParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], sParam)
		}
		if bParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, bParam)
		}
		return S(sResults[0]), err
	})
}

func (r *exportedPartialGenericResultsFnRecorder[S]) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedPartialGenericResultsFnRecorder[S]) isNil() bool {
	return r.r == nil
}
