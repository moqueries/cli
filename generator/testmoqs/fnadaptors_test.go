package testmoqs_test

import (
	"reflect"

	"github.com/myshkin5/moqueries/generator/testmoqs/exported"
	"github.com/myshkin5/moqueries/moq"
)

type usualFnAdaptor struct {
	m *moqUsualFn
}

func (a *usualFnAdaptor) exported() bool { return false }

func (a *usualFnAdaptor) tracksParams() bool { return true }

func (a *usualFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &usualFnRecorder{r: a.m.onCall(sParams[0], bParam)}
}

func (a *usualFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *usualFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqUsualFn_params{sParam: sParams[0], bParam: bParam}
}

func (a *usualFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type usualFnRecorder struct {
	r *moqUsualFn_fnRecorder
}

func (r *usualFnRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.sParam()
	}
}

func (r *usualFnRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (a *exportedUsualFnAdaptor) exported() bool { return true }

func (a *exportedUsualFnAdaptor) tracksParams() bool { return true }

func (a *exportedUsualFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedUsualFnRecorder{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedUsualFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedUsualFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqUsualFn_params{SParam: sParams[0], BParam: bParam}
}

func (a *exportedUsualFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedUsualFnRecorder struct {
	r *exported.MoqUsualFn_fnRecorder
}

func (r *exportedUsualFnRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.SParam()
	}
}

func (r *exportedUsualFnRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (a *noNamesFnAdaptor) exported() bool { return false }

func (a *noNamesFnAdaptor) tracksParams() bool { return true }

func (a *noNamesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &noNamesFnRecorder{r: a.m.onCall(sParams[0], bParam)}
}

func (a *noNamesFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *noNamesFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqNoNamesFn_params{param1: sParams[0], param2: bParam}
}

func (a *noNamesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type noNamesFnRecorder struct {
	r *moqNoNamesFn_fnRecorder
}

func (r *noNamesFnRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.param1()
	}
}

func (r *noNamesFnRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (a *exportedNoNamesFnAdaptor) exported() bool { return true }

func (a *exportedNoNamesFnAdaptor) tracksParams() bool { return true }

func (a *exportedNoNamesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedNoNamesFnRecorder{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedNoNamesFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedNoNamesFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqNoNamesFn_params{Param1: sParams[0], Param2: bParam}
}

func (a *exportedNoNamesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedNoNamesFnRecorder struct {
	r *exported.MoqNoNamesFn_fnRecorder
}

func (r *exportedNoNamesFnRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.Param1()
	}
}

func (r *exportedNoNamesFnRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (a *noResultsFnAdaptor) exported() bool { return false }

func (a *noResultsFnAdaptor) tracksParams() bool { return true }

func (a *noResultsFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &noResultsFnRecorder{r: a.m.onCall(sParams[0], bParam)}
}

func (a *noResultsFnAdaptor) invokeMockAndExpectResults(_ moq.T, sParams []string, bParam bool, _ results) {
	a.m.mock()(sParams[0], bParam)
}

func (a *noResultsFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqNoResultsFn_params{sParam: sParams[0], bParam: bParam}
}

func (a *noResultsFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type noResultsFnRecorder struct {
	r *moqNoResultsFn_fnRecorder
}

func (r *noResultsFnRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.sParam()
	}
}

func (r *noResultsFnRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, _ []string, _ error) {
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

func (a *exportedNoResultsFnAdaptor) exported() bool { return true }

func (a *exportedNoResultsFnAdaptor) tracksParams() bool { return true }

func (a *exportedNoResultsFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedNoResultsFnRecorder{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedNoResultsFnAdaptor) invokeMockAndExpectResults(_ moq.T, sParams []string, bParam bool, _ results) {
	a.m.Mock()(sParams[0], bParam)
}

func (a *exportedNoResultsFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqNoResultsFn_params{SParam: sParams[0], BParam: bParam}
}

func (a *exportedNoResultsFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedNoResultsFnRecorder struct {
	r *exported.MoqNoResultsFn_fnRecorder
}

func (r *exportedNoResultsFnRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.SParam()
	}
}

func (r *exportedNoResultsFnRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, _ []string, _ error) {
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

func (a *noParamsFnAdaptor) exported() bool { return false }

func (a *noParamsFnAdaptor) tracksParams() bool { return false }

func (a *noParamsFnAdaptor) newRecorder([]string, bool) recorder {
	return &noParamsFnRecorder{r: a.m.onCall()}
}

func (a *noParamsFnAdaptor) invokeMockAndExpectResults(t moq.T, _ []string, _ bool, res results) {
	sResult, err := a.m.mock()()
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *noParamsFnAdaptor) bundleParams([]string, bool) interface{} {
	return moqNoParamsFn_params{}
}

func (a *noParamsFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type noParamsFnRecorder struct {
	r *moqNoParamsFn_fnRecorder
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
	_ moq.T, fn func(), _ []string, _ bool, sResults []string, err error) {
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

func (a *exportedNoParamsFnAdaptor) exported() bool { return true }

func (a *exportedNoParamsFnAdaptor) tracksParams() bool { return false }

func (a *exportedNoParamsFnAdaptor) newRecorder([]string, bool) recorder {
	return &exportedNoParamsFnRecorder{r: a.m.OnCall()}
}

func (a *exportedNoParamsFnAdaptor) invokeMockAndExpectResults(t moq.T, _ []string, _ bool, res results) {
	sResult, err := a.m.Mock()()
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedNoParamsFnAdaptor) bundleParams([]string, bool) interface{} {
	return exported.MoqNoParamsFn_params{}
}

func (a *exportedNoParamsFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedNoParamsFnRecorder struct {
	r *exported.MoqNoParamsFn_fnRecorder
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
	_ moq.T, fn func(), _ []string, _ bool, sResults []string, err error) {
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

func (a *nothingFnAdaptor) exported() bool { return false }

func (a *nothingFnAdaptor) tracksParams() bool { return false }

func (a *nothingFnAdaptor) newRecorder([]string, bool) recorder {
	return &nothingFnRecorder{r: a.m.onCall()}
}

func (a *nothingFnAdaptor) invokeMockAndExpectResults(moq.T, []string, bool, results) {
	a.m.mock()()
}

func (a *nothingFnAdaptor) bundleParams([]string, bool) interface{} {
	return moqNothingFn_params{}
}

func (a *nothingFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type nothingFnRecorder struct {
	r *moqNothingFn_fnRecorder
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
	_ moq.T, fn func(), _ []string, _ bool, _ []string, _ error) {
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

func (a *exportedNothingFnAdaptor) exported() bool { return true }

func (a *exportedNothingFnAdaptor) tracksParams() bool { return false }

func (a *exportedNothingFnAdaptor) newRecorder([]string, bool) recorder {
	return &exportedNothingFnRecorder{r: a.m.OnCall()}
}

func (a *exportedNothingFnAdaptor) invokeMockAndExpectResults(moq.T, []string, bool, results) {
	a.m.Mock()()
}

func (a *exportedNothingFnAdaptor) bundleParams([]string, bool) interface{} {
	return exported.MoqNothingFn_params{}
}

func (a *exportedNothingFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedNothingFnRecorder struct {
	r *exported.MoqNothingFn_fnRecorder
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
	_ moq.T, fn func(), _ []string, _ bool, _ []string, _ error) {
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

func (a *variadicFnAdaptor) exported() bool { return false }

func (a *variadicFnAdaptor) tracksParams() bool { return true }

func (a *variadicFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &variadicFnRecorder{r: a.m.onCall(bParam, sParams...)}
}

func (a *variadicFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock()(bParam, sParams...)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *variadicFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqVariadicFn_params{args: sParams, other: bParam}
}

func (a *variadicFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type variadicFnRecorder struct {
	r *moqVariadicFn_fnRecorder
}

func (r *variadicFnRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.args()
	}
}

func (r *variadicFnRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (a *exportedVariadicFnAdaptor) exported() bool { return true }

func (a *exportedVariadicFnAdaptor) tracksParams() bool { return true }

func (a *exportedVariadicFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedVariadicFnRecorder{r: a.m.OnCall(bParam, sParams...)}
}

func (a *exportedVariadicFnAdaptor) invokeMockAndExpectResults(
	t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock()(bParam, sParams...)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedVariadicFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqVariadicFn_params{Args: sParams, Other: bParam}
}

func (a *exportedVariadicFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedVariadicFnRecorder struct {
	r *exported.MoqVariadicFn_fnRecorder
}

func (r *exportedVariadicFnRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.Args()
	}
}

func (r *exportedVariadicFnRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (a *repeatedIdsFnAdaptor) exported() bool { return false }

func (a *repeatedIdsFnAdaptor) tracksParams() bool { return true }

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
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *repeatedIdsFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqRepeatedIdsFn_params{sParam1: sParams[0], sParam2: sParams[1], bParam: bParam}
}

func (a *repeatedIdsFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type repeatedIdsFnRecorder struct {
	r *moqRepeatedIdsFn_fnRecorder
}

func (r *repeatedIdsFnRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.sParam1()
	}
}

func (r *repeatedIdsFnRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (a *exportedRepeatedIdsFnAdaptor) exported() bool { return true }

func (a *exportedRepeatedIdsFnAdaptor) tracksParams() bool { return true }

func (a *exportedRepeatedIdsFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedRepeatedIdsFnRecorder{r: a.m.OnCall(sParams[0], sParams[1], bParam)}
}

func (a *exportedRepeatedIdsFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult1, sResult2, err := a.m.Mock()(sParams[0], sParams[1], bParam)
	if sResult1 != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult1)
	}
	if sResult2 != res.sResults[1] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult2)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedRepeatedIdsFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqRepeatedIdsFn_params{SParam1: sParams[0], SParam2: sParams[1], BParam: bParam}
}

func (a *exportedRepeatedIdsFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedRepeatedIdsFnRecorder struct {
	r *exported.MoqRepeatedIdsFn_fnRecorder
}

func (r *exportedRepeatedIdsFnRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.SParam1()
	}
}

func (r *exportedRepeatedIdsFnRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (a *timesFnAdaptor) exported() bool { return false }

func (a *timesFnAdaptor) tracksParams() bool { return true }

func (a *timesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &timesFnRecorder{r: a.m.onCall(sParams[0], bParam)}
}

func (a *timesFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *timesFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqTimesFn_params{times: sParams[0], bParam: bParam}
}

func (a *timesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type timesFnRecorder struct {
	r *moqTimesFn_fnRecorder
}

func (r *timesFnRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.times()
	}
}

func (r *timesFnRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (a *exportedTimesFnAdaptor) exported() bool { return true }

func (a *exportedTimesFnAdaptor) tracksParams() bool { return true }

func (a *exportedTimesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedTimesFnRecorder{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedTimesFnAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock()(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedTimesFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqTimesFn_params{Times: sParams[0], BParam: bParam}
}

func (a *exportedTimesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedTimesFnRecorder struct {
	r *exported.MoqTimesFn_fnRecorder
}

func (r *exportedTimesFnRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.Times()
	}
}

func (r *exportedTimesFnRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (a *difficultParamNamesFnAdaptor) exported() bool { return false }

func (a *difficultParamNamesFnAdaptor) tracksParams() bool { return true }

func (a *difficultParamNamesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &difficultParamNamesFnRecorder{r: a.m.onCall(bParam, false, sParams[0], 0, 0, 0.0, 0.0)}
}

func (a *difficultParamNamesFnAdaptor) invokeMockAndExpectResults(_ moq.T, sParams []string, bParam bool, _ results) {
	a.m.mock()(bParam, false, sParams[0], 0, 0, 0.0, 0.0)
}

func (a *difficultParamNamesFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqDifficultParamNamesFn_params{param1: bParam, param2: false, param3: sParams[0]}
}

func (a *difficultParamNamesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type difficultParamNamesFnRecorder struct {
	r *moqDifficultParamNamesFn_fnRecorder
}

func (r *difficultParamNamesFnRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.param3()
	}
}

func (r *difficultParamNamesFnRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
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
	r.r = r.r.andDo(func(m, _ bool, sequence string, _, _ int, _, _ float32) {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, _ []string, _ error) {
	r.r = r.r.doReturnResults(func(m, _ bool, sequence string, _, _ int, _, _ float32) {
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

func (a *exportedDifficultParamNamesFnAdaptor) exported() bool { return true }

func (a *exportedDifficultParamNamesFnAdaptor) tracksParams() bool { return true }

func (a *exportedDifficultParamNamesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedDifficultParamNamesFnRecorder{r: a.m.OnCall(bParam, false, sParams[0], 0, 0, 0.0, 0.0)}
}

func (a *exportedDifficultParamNamesFnAdaptor) invokeMockAndExpectResults(_ moq.T, sParams []string, bParam bool, _ results) {
	a.m.Mock()(bParam, false, sParams[0], 0, 0, 0.0, 0.0)
}

func (a *exportedDifficultParamNamesFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqDifficultParamNamesFn_params{Param1: bParam, Param2: false, Param3: sParams[0]}
}

func (a *exportedDifficultParamNamesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedDifficultParamNamesFnRecorder struct {
	r *exported.MoqDifficultParamNamesFn_fnRecorder
}

func (r *exportedDifficultParamNamesFnRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.Param3()
	}
}

func (r *exportedDifficultParamNamesFnRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
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

func (r *exportedDifficultParamNamesFnRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.AndDo(func(m, _ bool, sequence string, _, _ int, _, _ float32) {
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
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, _ []string, _ error) {
	r.r = r.r.DoReturnResults(func(m, _ bool, sequence string, _, _ int, _, _ float32) {
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

func (a *difficultResultNamesFnAdaptor) exported() bool { return false }

func (a *difficultResultNamesFnAdaptor) tracksParams() bool { return false }

func (a *difficultResultNamesFnAdaptor) newRecorder([]string, bool) recorder {
	return &difficultResultNamesFnRecorder{r: a.m.onCall()}
}

func (a *difficultResultNamesFnAdaptor) invokeMockAndExpectResults(t moq.T, _ []string, _ bool, res results) {
	m, r, sequence, _, _, _, _ := a.m.mock()()
	if m != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], m)
	}
	if r != res.sResults[1] {
		t.Errorf("wanted %#v, got %#v", res.sResults[1], m)
	}
	if sequence != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, sequence)
	}
}

func (a *difficultResultNamesFnAdaptor) bundleParams([]string, bool) interface{} {
	return moqDifficultResultNamesFn_params{}
}

func (a *difficultResultNamesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type difficultResultNamesFnRecorder struct {
	r *moqDifficultResultNamesFn_fnRecorder
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
	r.r = r.r.returnResults(sResults[0], sResults[1], err, 0, 0, 0.0, 0.0)
}

func (r *difficultResultNamesFnRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.andDo(func() {
		fn()
	})
}

func (r *difficultResultNamesFnRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, sResults []string, err error) {
	r.r = r.r.doReturnResults(func() (m, r string, sequence error, _, _ int, _, _ float32) {
		fn()
		return sResults[0], sResults[1], err, 0, 0, 0.0, 0.0
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

func (a *exportedDifficultResultNamesFnAdaptor) exported() bool { return true }

func (a *exportedDifficultResultNamesFnAdaptor) tracksParams() bool { return false }

func (a *exportedDifficultResultNamesFnAdaptor) newRecorder([]string, bool) recorder {
	return &exportedDifficultResultNamesFnRecorder{r: a.m.OnCall()}
}

func (a *exportedDifficultResultNamesFnAdaptor) invokeMockAndExpectResults(t moq.T, _ []string, _ bool, res results) {
	m, r, sequence, _, _, _, _ := a.m.Mock()()
	if m != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], m)
	}
	if r != res.sResults[1] {
		t.Errorf("wanted %#v, got %#v", res.sResults[1], m)
	}
	if sequence != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, sequence)
	}
}

func (a *exportedDifficultResultNamesFnAdaptor) bundleParams([]string, bool) interface{} {
	return exported.MoqDifficultResultNamesFn_params{}
}

func (a *exportedDifficultResultNamesFnAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedDifficultResultNamesFnRecorder struct {
	r *exported.MoqDifficultResultNamesFn_fnRecorder
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
	r.r = r.r.ReturnResults(sResults[0], sResults[1], err, 0, 0, 0.0, 0.0)
}

func (r *exportedDifficultResultNamesFnRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.AndDo(func() {
		fn()
	})
}

func (r *exportedDifficultResultNamesFnRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, sResults []string, err error) {
	r.r = r.r.DoReturnResults(func() (m, r string, sequence error, _, _ int, _, _ float32) {
		fn()
		return sResults[0], sResults[1], err, 0, 0, 0.0, 0.0
	})
}

func (r *exportedDifficultResultNamesFnRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedDifficultResultNamesFnRecorder) isNil() bool {
	return r.r == nil
}
