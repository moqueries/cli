package testmoqs_test

import (
	"io"
	"reflect"

	"github.com/myshkin5/moqueries/generator/testmoqs"
	"github.com/myshkin5/moqueries/generator/testmoqs/exported"
	"github.com/myshkin5/moqueries/moq"
)

type usualAdaptor struct {
	m *moqUsual
}

func (a *usualAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *usualAdaptor) mock() interface{} { return a.m.mock() }

func (a *usualAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &usualRecorder{r: a.m.onCall().Usual(sParams[0], bParam)}
}

func (a *usualAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock().Usual(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *usualAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqUsual_Usual_params{sParam: sParams[0], bParam: bParam}
}

func (a *usualAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type usualRecorder struct {
	r *moqUsual_Usual_fnRecorder
}

func (r *usualRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.sParam()
	}
}

func (r *usualRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.bParam()
	}
}

func (r *usualRecorder) seq() {
	r.r = r.r.seq()
}

func (r *usualRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *usualRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *usualRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *usualRecorder) doReturnResults(
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

func (r *usualRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *usualRecorder) isNil() bool {
	return r.r == nil
}

type exportedUsualAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedUsualAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedUsualAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedUsualAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedUsualRecorder{r: a.m.OnCall().Usual(sParams[0], bParam)}
}

func (a *exportedUsualAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock().Usual(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedUsualAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqUsual_Usual_params{SParam: sParams[0], BParam: bParam}
}

func (a *exportedUsualAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedUsualRecorder struct {
	r *exported.MoqUsual_Usual_fnRecorder
}

func (r *exportedUsualRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.SParam()
	}
}

func (r *exportedUsualRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.BParam()
	}
}

func (r *exportedUsualRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedUsualRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedUsualRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedUsualRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *exportedUsualRecorder) doReturnResults(
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

func (r *exportedUsualRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedUsualRecorder) isNil() bool {
	return r.r == nil
}

type noNamesAdaptor struct {
	m *moqUsual
}

func (a *noNamesAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *noNamesAdaptor) mock() interface{} { return a.m.mock() }

func (a *noNamesAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &noNamesRecorder{r: a.m.onCall().NoNames(sParams[0], bParam)}
}

func (a *noNamesAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock().NoNames(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *noNamesAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqUsual_NoNames_params{param1: sParams[0], param2: bParam}
}

func (a *noNamesAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type noNamesRecorder struct {
	r *moqUsual_NoNames_fnRecorder
}

func (r *noNamesRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.param1()
	}
}

func (r *noNamesRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.param2()
	}
}

func (r *noNamesRecorder) seq() {
	r.r = r.r.seq()
}

func (r *noNamesRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *noNamesRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *noNamesRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *noNamesRecorder) doReturnResults(
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

func (r *noNamesRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *noNamesRecorder) isNil() bool {
	return r.r == nil
}

type exportedNoNamesAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedNoNamesAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedNoNamesAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedNoNamesAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedNoNamesRecorder{r: a.m.OnCall().NoNames(sParams[0], bParam)}
}

func (a *exportedNoNamesAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock().NoNames(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedNoNamesAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqUsual_NoNames_params{Param1: sParams[0], Param2: bParam}
}

func (a *exportedNoNamesAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedNoNamesRecorder struct {
	r *exported.MoqUsual_NoNames_fnRecorder
}

func (r *exportedNoNamesRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.Param1()
	}
}

func (r *exportedNoNamesRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.Param2()
	}
}

func (r *exportedNoNamesRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedNoNamesRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedNoNamesRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedNoNamesRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *exportedNoNamesRecorder) doReturnResults(
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

func (r *exportedNoNamesRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedNoNamesRecorder) isNil() bool {
	return r.r == nil
}

type noResultsAdaptor struct {
	m *moqUsual
}

func (a *noResultsAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *noResultsAdaptor) mock() interface{} { return a.m.mock() }

func (a *noResultsAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &noResultsRecorder{r: a.m.onCall().NoResults(sParams[0], bParam)}
}

func (a *noResultsAdaptor) invokeMockAndExpectResults(_ moq.T, sParams []string, bParam bool, _ results) {
	a.m.mock().NoResults(sParams[0], bParam)
}

func (a *noResultsAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqUsual_NoResults_params{sParam: sParams[0], bParam: bParam}
}

func (a *noResultsAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type noResultsRecorder struct {
	r *moqUsual_NoResults_fnRecorder
}

func (r *noResultsRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.sParam()
	}
}

func (r *noResultsRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.bParam()
	}
}

func (r *noResultsRecorder) seq() {
	r.r = r.r.seq()
}

func (r *noResultsRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *noResultsRecorder) returnResults([]string, error) {
	r.r = r.r.returnResults()
}

func (r *noResultsRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *noResultsRecorder) doReturnResults(
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

func (r *noResultsRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *noResultsRecorder) isNil() bool {
	return r.r == nil
}

type exportedNoResultsAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedNoResultsAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedNoResultsAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedNoResultsAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedNoResultsRecorder{r: a.m.OnCall().NoResults(sParams[0], bParam)}
}

func (a *exportedNoResultsAdaptor) invokeMockAndExpectResults(_ moq.T, sParams []string, bParam bool, _ results) {
	a.m.Mock().NoResults(sParams[0], bParam)
}

func (a *exportedNoResultsAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqUsual_NoResults_params{SParam: sParams[0], BParam: bParam}
}

func (a *exportedNoResultsAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedNoResultsRecorder struct {
	r *exported.MoqUsual_NoResults_fnRecorder
}

func (r *exportedNoResultsRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.SParam()
	}
}

func (r *exportedNoResultsRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.BParam()
	}
}

func (r *exportedNoResultsRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedNoResultsRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedNoResultsRecorder) returnResults([]string, error) {
	r.r = r.r.ReturnResults()
}

func (r *exportedNoResultsRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *exportedNoResultsRecorder) doReturnResults(
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

func (r *exportedNoResultsRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedNoResultsRecorder) isNil() bool {
	return r.r == nil
}

type noParamsAdaptor struct {
	m *moqUsual
}

func (a *noParamsAdaptor) config() adaptorConfig {
	return adaptorConfig{noParams: true}
}

func (a *noParamsAdaptor) mock() interface{} { return a.m.mock() }

func (a *noParamsAdaptor) newRecorder([]string, bool) recorder {
	return &noParamsRecorder{r: a.m.onCall().NoParams()}
}

func (a *noParamsAdaptor) invokeMockAndExpectResults(t moq.T, _ []string, _ bool, res results) {
	sResult, err := a.m.mock().NoParams()
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *noParamsAdaptor) bundleParams([]string, bool) interface{} {
	return moqUsual_NoParams_params{}
}

func (a *noParamsAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type noParamsRecorder struct {
	r *moqUsual_NoParams_fnRecorder
}

func (r *noParamsRecorder) anySParam() {}

func (r *noParamsRecorder) anyBParam() {}

func (r *noParamsRecorder) seq() {
	r.r = r.r.seq()
}

func (r *noParamsRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *noParamsRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *noParamsRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.andDo(func() {
		fn()
	})
}

func (r *noParamsRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, sResults []string, err error) {
	r.r = r.r.doReturnResults(func() (string, error) {
		fn()
		return sResults[0], err
	})
}

func (r *noParamsRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *noParamsRecorder) isNil() bool {
	return r.r == nil
}

type exportedNoParamsAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedNoParamsAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true, noParams: true}
}

func (a *exportedNoParamsAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedNoParamsAdaptor) newRecorder([]string, bool) recorder {
	return &exportedNoParamsRecorder{r: a.m.OnCall().NoParams()}
}

func (a *exportedNoParamsAdaptor) invokeMockAndExpectResults(t moq.T, _ []string, _ bool, res results) {
	sResult, err := a.m.Mock().NoParams()
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedNoParamsAdaptor) bundleParams([]string, bool) interface{} {
	return exported.MoqUsual_NoParams_params{}
}

func (a *exportedNoParamsAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedNoParamsRecorder struct {
	r *exported.MoqUsual_NoParams_fnRecorder
}

func (r *exportedNoParamsRecorder) anySParam() {}

func (r *exportedNoParamsRecorder) anyBParam() {}

func (r *exportedNoParamsRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedNoParamsRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedNoParamsRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedNoParamsRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.AndDo(func() {
		fn()
	})
}

func (r *exportedNoParamsRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, sResults []string, err error) {
	r.r = r.r.DoReturnResults(func() (string, error) {
		fn()
		return sResults[0], err
	})
}

func (r *exportedNoParamsRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedNoParamsRecorder) isNil() bool {
	return r.r == nil
}

type nothingAdaptor struct {
	m *moqUsual
}

func (a *nothingAdaptor) config() adaptorConfig {
	return adaptorConfig{noParams: true}
}

func (a *nothingAdaptor) mock() interface{} { return a.m.mock() }

func (a *nothingAdaptor) newRecorder([]string, bool) recorder {
	return &nothingRecorder{r: a.m.onCall().Nothing()}
}

func (a *nothingAdaptor) invokeMockAndExpectResults(moq.T, []string, bool, results) {
	a.m.mock().Nothing()
}

func (a *nothingAdaptor) bundleParams([]string, bool) interface{} {
	return moqUsual_Nothing_params{}
}

func (a *nothingAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type nothingRecorder struct {
	r *moqUsual_Nothing_fnRecorder
}

func (r *nothingRecorder) anySParam() {}

func (r *nothingRecorder) anyBParam() {}

func (r *nothingRecorder) seq() {
	r.r = r.r.seq()
}

func (r *nothingRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *nothingRecorder) returnResults([]string, error) {
	r.r = r.r.returnResults()
}

func (r *nothingRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.andDo(func() {
		fn()
	})
}

func (r *nothingRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, _ []string, _ error) {
	r.r = r.r.doReturnResults(func() {
		fn()
	})
}

func (r *nothingRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *nothingRecorder) isNil() bool {
	return r.r == nil
}

type exportedNothingAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedNothingAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true, noParams: true}
}

func (a *exportedNothingAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedNothingAdaptor) newRecorder([]string, bool) recorder {
	return &exportedNothingRecorder{r: a.m.OnCall().Nothing()}
}

func (a *exportedNothingAdaptor) invokeMockAndExpectResults(moq.T, []string, bool, results) {
	a.m.Mock().Nothing()
}

func (a *exportedNothingAdaptor) bundleParams([]string, bool) interface{} {
	return exported.MoqUsual_Nothing_params{}
}

func (a *exportedNothingAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedNothingRecorder struct {
	r *exported.MoqUsual_Nothing_fnRecorder
}

func (r *exportedNothingRecorder) anySParam() {}

func (r *exportedNothingRecorder) anyBParam() {}

func (r *exportedNothingRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedNothingRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedNothingRecorder) returnResults([]string, error) {
	r.r = r.r.ReturnResults()
}

func (r *exportedNothingRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.AndDo(func() {
		fn()
	})
}

func (r *exportedNothingRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, _ []string, _ error) {
	r.r = r.r.DoReturnResults(func() {
		fn()
	})
}

func (r *exportedNothingRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedNothingRecorder) isNil() bool {
	return r.r == nil
}

type variadicAdaptor struct {
	m *moqUsual
}

func (a *variadicAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *variadicAdaptor) mock() interface{} { return a.m.mock() }

func (a *variadicAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &variadicRecorder{r: a.m.onCall().Variadic(bParam, sParams...)}
}

func (a *variadicAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock().Variadic(bParam, sParams...)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *variadicAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqUsual_Variadic_params{args: sParams, other: bParam}
}

func (a *variadicAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type variadicRecorder struct {
	r *moqUsual_Variadic_fnRecorder
}

func (r *variadicRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.args()
	}
}

func (r *variadicRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.other()
	}
}

func (r *variadicRecorder) seq() {
	r.r = r.r.seq()
}

func (r *variadicRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *variadicRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *variadicRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *variadicRecorder) doReturnResults(
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

func (r *variadicRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *variadicRecorder) isNil() bool {
	return r.r == nil
}

type exportedVariadicAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedVariadicAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedVariadicAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedVariadicAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedVariadicRecorder{r: a.m.OnCall().Variadic(bParam, sParams...)}
}

func (a *exportedVariadicAdaptor) invokeMockAndExpectResults(
	t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock().Variadic(bParam, sParams...)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedVariadicAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqUsual_Variadic_params{Args: sParams, Other: bParam}
}

func (a *exportedVariadicAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedVariadicRecorder struct {
	r *exported.MoqUsual_Variadic_fnRecorder
}

func (r *exportedVariadicRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.Args()
	}
}

func (r *exportedVariadicRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.Other()
	}
}

func (r *exportedVariadicRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedVariadicRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedVariadicRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedVariadicRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *exportedVariadicRecorder) doReturnResults(
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

func (r *exportedVariadicRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedVariadicRecorder) isNil() bool {
	return r.r == nil
}

type repeatedIdsAdaptor struct {
	m *moqUsual
}

func (a *repeatedIdsAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *repeatedIdsAdaptor) mock() interface{} { return a.m.mock() }

func (a *repeatedIdsAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &repeatedIdsRecorder{r: a.m.onCall().RepeatedIds(sParams[0], sParams[1], bParam)}
}

func (a *repeatedIdsAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult1, sResult2, err := a.m.mock().RepeatedIds(sParams[0], sParams[1], bParam)
	if sResult1 != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult1)
	}
	if sResult2 != res.sResults[1] {
		t.Errorf("wanted %#v, got %#v", res.sResults[1], sResult2)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *repeatedIdsAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqUsual_RepeatedIds_params{sParam1: sParams[0], sParam2: sParams[1], bParam: bParam}
}

func (a *repeatedIdsAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type repeatedIdsRecorder struct {
	r *moqUsual_RepeatedIds_fnRecorder
}

func (r *repeatedIdsRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.sParam1()
	}
}

func (r *repeatedIdsRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.bParam()
	}
}

func (r *repeatedIdsRecorder) seq() {
	r.r = r.r.seq()
}

func (r *repeatedIdsRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *repeatedIdsRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], sResults[1], err)
}

func (r *repeatedIdsRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *repeatedIdsRecorder) doReturnResults(
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

func (r *repeatedIdsRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *repeatedIdsRecorder) isNil() bool {
	return r.r == nil
}

type exportedRepeatedIdsAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedRepeatedIdsAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedRepeatedIdsAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedRepeatedIdsAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedRepeatedIdsRecorder{r: a.m.OnCall().RepeatedIds(sParams[0], sParams[1], bParam)}
}

func (a *exportedRepeatedIdsAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult1, sResult2, err := a.m.Mock().RepeatedIds(sParams[0], sParams[1], bParam)
	if sResult1 != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult1)
	}
	if sResult2 != res.sResults[1] {
		t.Errorf("wanted %#v, got %#v", res.sResults[1], sResult2)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedRepeatedIdsAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqUsual_RepeatedIds_params{SParam1: sParams[0], SParam2: sParams[1], BParam: bParam}
}

func (a *exportedRepeatedIdsAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedRepeatedIdsRecorder struct {
	r *exported.MoqUsual_RepeatedIds_fnRecorder
}

func (r *exportedRepeatedIdsRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.SParam1()
	}
}

func (r *exportedRepeatedIdsRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.BParam()
	}
}

func (r *exportedRepeatedIdsRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedRepeatedIdsRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedRepeatedIdsRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], sResults[1], err)
}

func (r *exportedRepeatedIdsRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *exportedRepeatedIdsRecorder) doReturnResults(
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

func (r *exportedRepeatedIdsRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedRepeatedIdsRecorder) isNil() bool {
	return r.r == nil
}

type timesAdaptor struct {
	m *moqUsual
}

func (a *timesAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *timesAdaptor) mock() interface{} { return a.m.mock() }

func (a *timesAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &timesRecorder{r: a.m.onCall().Times(sParams[0], bParam)}
}

func (a *timesAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock().Times(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *timesAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqUsual_Times_params{sParam: sParams[0], times: bParam}
}

func (a *timesAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type timesRecorder struct {
	r *moqUsual_Times_fnRecorder
}

func (r *timesRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.sParam()
	}
}

func (r *timesRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.times()
	}
}

func (r *timesRecorder) seq() {
	r.r = r.r.seq()
}

func (r *timesRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *timesRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *timesRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *timesRecorder) doReturnResults(
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

func (r *timesRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *timesRecorder) isNil() bool {
	return r.r == nil
}

type exportedTimesAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedTimesAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedTimesAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedTimesAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedTimesRecorder{r: a.m.OnCall().Times(sParams[0], bParam)}
}

func (a *exportedTimesAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock().Times(sParams[0], bParam)
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedTimesAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqUsual_Times_params{SParam: sParams[0], Times: bParam}
}

func (a *exportedTimesAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedTimesRecorder struct {
	r *exported.MoqUsual_Times_fnRecorder
}

func (r *exportedTimesRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.SParam()
	}
}

func (r *exportedTimesRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.Times()
	}
}

func (r *exportedTimesRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedTimesRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedTimesRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedTimesRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *exportedTimesRecorder) doReturnResults(
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

func (r *exportedTimesRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedTimesRecorder) isNil() bool {
	return r.r == nil
}

type difficultParamNamesAdaptor struct {
	m *moqUsual
}

func (a *difficultParamNamesAdaptor) config() adaptorConfig { return adaptorConfig{} }

func (a *difficultParamNamesAdaptor) mock() interface{} { return a.m.mock() }

func (a *difficultParamNamesAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &difficultParamNamesRecorder{r: a.m.onCall().DifficultParamNames(bParam, false, sParams[0], 0, 0, 0.0, 0.0)}
}

func (a *difficultParamNamesAdaptor) invokeMockAndExpectResults(_ moq.T, sParams []string, bParam bool, _ results) {
	a.m.mock().DifficultParamNames(bParam, false, sParams[0], 0, 0, 0.0, 0.0)
}

func (a *difficultParamNamesAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqDifficultParamNamesFn_params{param1: bParam, param2: false, param3: sParams[0]}
}

func (a *difficultParamNamesAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type difficultParamNamesRecorder struct {
	r *moqUsual_DifficultParamNames_fnRecorder
}

func (r *difficultParamNamesRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.param3()
	}
}

func (r *difficultParamNamesRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.param1()
	}
}

func (r *difficultParamNamesRecorder) seq() {
	r.r = r.r.seq()
}

func (r *difficultParamNamesRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *difficultParamNamesRecorder) returnResults([]string, error) {
	r.r = r.r.returnResults()
}

func (r *difficultParamNamesRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *difficultParamNamesRecorder) doReturnResults(
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

func (r *difficultParamNamesRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *difficultParamNamesRecorder) isNil() bool {
	return r.r == nil
}

type exportedDifficultParamNamesAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedDifficultParamNamesAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true}
}

func (a *exportedDifficultParamNamesAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedDifficultParamNamesAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedDifficultParamNamesRecorder{r: a.m.OnCall().DifficultParamNames(
		bParam, false, sParams[0], 0, 0, 0.0, 0.0)}
}

func (a *exportedDifficultParamNamesAdaptor) invokeMockAndExpectResults(_ moq.T, sParams []string, bParam bool, _ results) {
	a.m.Mock().DifficultParamNames(bParam, false, sParams[0], 0, 0, 0.0, 0.0)
}

func (a *exportedDifficultParamNamesAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MoqDifficultParamNamesFn_params{Param1: bParam, Param2: false, Param3: sParams[0]}
}

func (a *exportedDifficultParamNamesAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedDifficultParamNamesRecorder struct {
	r *exported.MoqUsual_DifficultParamNames_fnRecorder
}

func (r *exportedDifficultParamNamesRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.Param3()
	}
}

func (r *exportedDifficultParamNamesRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.Param1()
	}
}

func (r *exportedDifficultParamNamesRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedDifficultParamNamesRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedDifficultParamNamesRecorder) returnResults([]string, error) {
	r.r = r.r.ReturnResults()
}

func (r *exportedDifficultParamNamesRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *exportedDifficultParamNamesRecorder) doReturnResults(
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

func (r *exportedDifficultParamNamesRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedDifficultParamNamesRecorder) isNil() bool {
	return r.r == nil
}

type difficultResultNamesAdaptor struct {
	m *moqUsual
}

func (a *difficultResultNamesAdaptor) config() adaptorConfig {
	return adaptorConfig{noParams: true}
}

func (a *difficultResultNamesAdaptor) mock() interface{} { return a.m.mock() }

func (a *difficultResultNamesAdaptor) newRecorder([]string, bool) recorder {
	return &difficultResultNamesRecorder{r: a.m.onCall().DifficultResultNames()}
}

func (a *difficultResultNamesAdaptor) invokeMockAndExpectResults(t moq.T, _ []string, _ bool, res results) {
	m, r, sequence, _, _, _, _ := a.m.mock().DifficultResultNames()
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

func (a *difficultResultNamesAdaptor) bundleParams([]string, bool) interface{} {
	return moqDifficultResultNamesFn_params{}
}

func (a *difficultResultNamesAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type difficultResultNamesRecorder struct {
	r *moqUsual_DifficultResultNames_fnRecorder
}

func (r *difficultResultNamesRecorder) anySParam() {}

func (r *difficultResultNamesRecorder) anyBParam() {}

func (r *difficultResultNamesRecorder) seq() {
	r.r = r.r.seq()
}

func (r *difficultResultNamesRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *difficultResultNamesRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], sResults[1], err, 0, 0, 0.0, 0.0)
}

func (r *difficultResultNamesRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.andDo(func() {
		fn()
	})
}

func (r *difficultResultNamesRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, sResults []string, err error) {
	r.r = r.r.doReturnResults(func() (m, r string, sequence error, _, _ int, _, _ float32) {
		fn()
		return sResults[0], sResults[1], err, 0, 0, 0.0, 0.0
	})
}

func (r *difficultResultNamesRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *difficultResultNamesRecorder) isNil() bool {
	return r.r == nil
}

type exportedDifficultResultNamesAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedDifficultResultNamesAdaptor) config() adaptorConfig {
	return adaptorConfig{exported: true, noParams: true}
}

func (a *exportedDifficultResultNamesAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedDifficultResultNamesAdaptor) newRecorder([]string, bool) recorder {
	return &exportedDifficultResultNamesRecorder{r: a.m.OnCall().DifficultResultNames()}
}

func (a *exportedDifficultResultNamesAdaptor) invokeMockAndExpectResults(t moq.T, _ []string, _ bool, res results) {
	m, r, sequence, _, _, _, _ := a.m.Mock().DifficultResultNames()
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

func (a *exportedDifficultResultNamesAdaptor) bundleParams([]string, bool) interface{} {
	return exported.MoqDifficultResultNamesFn_params{}
}

func (a *exportedDifficultResultNamesAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedDifficultResultNamesRecorder struct {
	r *exported.MoqUsual_DifficultResultNames_fnRecorder
}

func (r *exportedDifficultResultNamesRecorder) anySParam() {}

func (r *exportedDifficultResultNamesRecorder) anyBParam() {}

func (r *exportedDifficultResultNamesRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedDifficultResultNamesRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedDifficultResultNamesRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], sResults[1], err, 0, 0, 0.0, 0.0)
}

func (r *exportedDifficultResultNamesRecorder) andDo(_ moq.T, fn func(), _ []string, _ bool) {
	r.r = r.r.AndDo(func() {
		fn()
	})
}

func (r *exportedDifficultResultNamesRecorder) doReturnResults(
	_ moq.T, fn func(), _ []string, _ bool, sResults []string, err error) {
	r.r = r.r.DoReturnResults(func() (m, r string, sequence error, _, _ int, _, _ float32) {
		fn()
		return sResults[0], sResults[1], err, 0, 0, 0.0, 0.0
	})
}

func (r *exportedDifficultResultNamesRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedDifficultResultNamesRecorder) isNil() bool {
	return r.r == nil
}

type passByReferenceAdaptor struct {
	m *moqUsual
}

func (a *passByReferenceAdaptor) config() adaptorConfig {
	return adaptorConfig{
		opaqueParams: true,
	}
}

func (a *passByReferenceAdaptor) mock() interface{} { return a.m.mock() }

func (a *passByReferenceAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &passByReferenceRecorder{r: a.m.onCall().PassByReference(&testmoqs.PassByReferenceParams{
		SParam: sParams[0],
		BParam: bParam,
	})}
}

func (a *passByReferenceAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock().PassByReference(&testmoqs.PassByReferenceParams{
		SParam: sParams[0],
		BParam: bParam,
	})
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *passByReferenceAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqUsual_PassByReference_params{p: &testmoqs.PassByReferenceParams{
		SParam: sParams[0],
		BParam: bParam,
	}}
}

func (a *passByReferenceAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type passByReferenceRecorder struct {
	r *moqUsual_PassByReference_fnRecorder
}

func (r *passByReferenceRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.p()
	}
}

func (r *passByReferenceRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.p()
	}
}

func (r *passByReferenceRecorder) seq() {
	r.r = r.r.seq()
}

func (r *passByReferenceRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *passByReferenceRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *passByReferenceRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.andDo(func(p *testmoqs.PassByReferenceParams) {
		fn()
		if p.SParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], p.SParam)
		}
		if p.BParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, p.BParam)
		}
	})
}

func (r *passByReferenceRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
	r.r = r.r.doReturnResults(func(p *testmoqs.PassByReferenceParams) (string, error) {
		fn()
		if p.SParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], p.SParam)
		}
		if p.BParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, p.BParam)
		}
		return sResults[0], err
	})
}

func (r *passByReferenceRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *passByReferenceRecorder) isNil() bool {
	return r.r == nil
}

type exportedPassByReferenceAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedPassByReferenceAdaptor) config() adaptorConfig {
	return adaptorConfig{
		exported:     true,
		opaqueParams: true,
	}
}

func (a *exportedPassByReferenceAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedPassByReferenceAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedPassByReferenceRecorder{r: a.m.OnCall().PassByReference(&testmoqs.PassByReferenceParams{
		SParam: sParams[0],
		BParam: bParam,
	})}
}

func (a *exportedPassByReferenceAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock().PassByReference(&testmoqs.PassByReferenceParams{
		SParam: sParams[0],
		BParam: bParam,
	})
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedPassByReferenceAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqUsual_PassByReference_params{p: &testmoqs.PassByReferenceParams{
		SParam: sParams[0],
		BParam: bParam,
	}}
}

func (a *exportedPassByReferenceAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedPassByReferenceRecorder struct {
	r *exported.MoqUsual_PassByReference_fnRecorder
}

func (r *exportedPassByReferenceRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.P()
	}
}

func (r *exportedPassByReferenceRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.P()
	}
}

func (r *exportedPassByReferenceRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedPassByReferenceRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedPassByReferenceRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedPassByReferenceRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
	r.r = r.r.AndDo(func(p *testmoqs.PassByReferenceParams) {
		fn()
		if p.SParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], p.SParam)
		}
		if p.BParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, p.BParam)
		}
	})
}

func (r *exportedPassByReferenceRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
	r.r = r.r.DoReturnResults(func(p *testmoqs.PassByReferenceParams) (string, error) {
		fn()
		if p.SParam != expectedSParams[0] {
			t.Errorf("wanted %#v, got %#v", expectedSParams[0], p.SParam)
		}
		if p.BParam != expectedBParam {
			t.Errorf("wanted %t, got %#v", expectedBParam, p.BParam)
		}
		return sResults[0], err
	})
}

func (r *exportedPassByReferenceRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedPassByReferenceRecorder) isNil() bool {
	return r.r == nil
}

type interfaceParamAdaptor struct {
	m *moqUsual
}

func (a *interfaceParamAdaptor) config() adaptorConfig {
	return adaptorConfig{
		opaqueParams: true,
	}
}

func (a *interfaceParamAdaptor) mock() interface{} { return a.m.mock() }

func (a *interfaceParamAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &interfaceParamRecorder{r: a.m.onCall().InterfaceParam(&testmoqs.InterfaceParamWriter{
		SParam: sParams[0],
		BParam: bParam,
	})}
}

func (a *interfaceParamAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock().InterfaceParam(&testmoqs.InterfaceParamWriter{
		SParam: sParams[0],
		BParam: bParam,
	})
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *interfaceParamAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqUsual_InterfaceParam_params{w: &testmoqs.InterfaceParamWriter{
		SParam: sParams[0],
		BParam: bParam,
	}}
}

func (a *interfaceParamAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type interfaceParamRecorder struct {
	r *moqUsual_InterfaceParam_fnRecorder
}

func (r *interfaceParamRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.w()
	}
}

func (r *interfaceParamRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.w()
	}
}

func (r *interfaceParamRecorder) seq() {
	r.r = r.r.seq()
}

func (r *interfaceParamRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *interfaceParamRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(sResults[0], err)
}

func (r *interfaceParamRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *interfaceParamRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (r *interfaceParamRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *interfaceParamRecorder) isNil() bool {
	return r.r == nil
}

type exportedInterfaceParamAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedInterfaceParamAdaptor) config() adaptorConfig {
	return adaptorConfig{
		exported:     true,
		opaqueParams: true,
	}
}

func (a *exportedInterfaceParamAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedInterfaceParamAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedInterfaceParamRecorder{r: a.m.OnCall().InterfaceParam(&testmoqs.InterfaceParamWriter{
		SParam: sParams[0],
		BParam: bParam,
	})}
}

func (a *exportedInterfaceParamAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock().InterfaceParam(&testmoqs.InterfaceParamWriter{
		SParam: sParams[0],
		BParam: bParam,
	})
	if sResult != res.sResults[0] {
		t.Errorf("wanted %#v, got %#v", res.sResults[0], sResult)
	}
	if err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, err)
	}
}

func (a *exportedInterfaceParamAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqUsual_InterfaceParam_params{w: &testmoqs.InterfaceParamWriter{
		SParam: sParams[0],
		BParam: bParam,
	}}
}

func (a *exportedInterfaceParamAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedInterfaceParamRecorder struct {
	r *exported.MoqUsual_InterfaceParam_fnRecorder
}

func (r *exportedInterfaceParamRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.W()
	}
}

func (r *exportedInterfaceParamRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.W()
	}
}

func (r *exportedInterfaceParamRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedInterfaceParamRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedInterfaceParamRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(sResults[0], err)
}

func (r *exportedInterfaceParamRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *exportedInterfaceParamRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (r *exportedInterfaceParamRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedInterfaceParamRecorder) isNil() bool {
	return r.r == nil
}

type interfaceResultAdaptor struct {
	m *moqUsual
}

func (a *interfaceResultAdaptor) config() adaptorConfig {
	return adaptorConfig{
		opaqueParams: true,
	}
}

func (a *interfaceResultAdaptor) mock() interface{} { return a.m.mock() }

func (a *interfaceResultAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &interfaceResultRecorder{r: a.m.onCall().InterfaceResult(sParams[0], bParam)}
}

func (a *interfaceResultAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	r := a.m.mock().InterfaceResult(sParams[0], bParam)
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
	if irr.Err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, irr)
	}
}

func (a *interfaceResultAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqUsual_InterfaceResult_params{
		sParam: sParams[0],
		bParam: bParam,
	}
}

func (a *interfaceResultAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type interfaceResultRecorder struct {
	r *moqUsual_InterfaceResult_fnRecorder
}

func (r *interfaceResultRecorder) anySParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.sParam()
	}
}

func (r *interfaceResultRecorder) anyBParam() {
	a := r.r.any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.bParam()
	}
}

func (r *interfaceResultRecorder) seq() {
	r.r = r.r.seq()
}

func (r *interfaceResultRecorder) noSeq() {
	r.r = r.r.noSeq()
}

func (r *interfaceResultRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.returnResults(&testmoqs.InterfaceResultReader{
		SResult: sResults[0],
		Err:     err,
	})
}

func (r *interfaceResultRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *interfaceResultRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (r *interfaceResultRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.repeat(repeaters...)
}

func (r *interfaceResultRecorder) isNil() bool {
	return r.r == nil
}

type exportedInterfaceResultAdaptor struct {
	m *exported.MoqUsual
}

func (a *exportedInterfaceResultAdaptor) config() adaptorConfig {
	return adaptorConfig{
		exported:     true,
		opaqueParams: true,
	}
}

func (a *exportedInterfaceResultAdaptor) mock() interface{} { return a.m.Mock() }

func (a *exportedInterfaceResultAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedInterfaceResultRecorder{r: a.m.OnCall().InterfaceResult(sParams[0], bParam)}
}

func (a *exportedInterfaceResultAdaptor) invokeMockAndExpectResults(t moq.T, sParams []string, bParam bool, res results) {
	r := a.m.Mock().InterfaceResult(sParams[0], bParam)
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
	if irr.Err != res.err {
		t.Errorf("wanted %#v, got %#v", res.err, irr)
	}
}

func (a *exportedInterfaceResultAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return moqUsual_InterfaceResult_params{
		sParam: sParams[0],
		bParam: bParam,
	}
}

func (a *exportedInterfaceResultAdaptor) sceneMoq() moq.Moq {
	return a.m
}

type exportedInterfaceResultRecorder struct {
	r *exported.MoqUsual_InterfaceResult_fnRecorder
}

func (r *exportedInterfaceResultRecorder) anySParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.SParam()
	}
}

func (r *exportedInterfaceResultRecorder) anyBParam() {
	a := r.r.Any()
	if a == nil {
		r.r = nil
	} else {
		r.r = a.BParam()
	}
}

func (r *exportedInterfaceResultRecorder) seq() {
	r.r = r.r.Seq()
}

func (r *exportedInterfaceResultRecorder) noSeq() {
	r.r = r.r.NoSeq()
}

func (r *exportedInterfaceResultRecorder) returnResults(sResults []string, err error) {
	r.r = r.r.ReturnResults(&testmoqs.InterfaceResultReader{
		SResult: sResults[0],
		Err:     err,
	})
}

func (r *exportedInterfaceResultRecorder) andDo(t moq.T, fn func(), expectedSParams []string, expectedBParam bool) {
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

func (r *exportedInterfaceResultRecorder) doReturnResults(
	t moq.T, fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

func (r *exportedInterfaceResultRecorder) repeat(repeaters ...moq.Repeater) {
	r.r = r.r.Repeat(repeaters...)
}

func (r *exportedInterfaceResultRecorder) isNil() bool {
	return r.r == nil
}
