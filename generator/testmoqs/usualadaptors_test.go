package testmoqs_test

import (
	"reflect"

	"github.com/myshkin5/moqueries/generator/testmoqs/exported"
	"github.com/myshkin5/moqueries/moq"
)

type usualAdaptor struct{ m *moqUsual }

func (a *usualAdaptor) exported() bool { return false }

func (a *usualAdaptor) tracksParams() bool { return true }

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

type usualRecorder struct{ r *moqUsual_Usual_fnRecorder }

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

func (r *usualRecorder) doReturnResults(t moq.T,
	fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

type exportedUsualAdaptor struct{ m *exported.MoqUsual }

func (a *exportedUsualAdaptor) exported() bool { return true }

func (a *exportedUsualAdaptor) tracksParams() bool { return true }

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

func (r *exportedUsualRecorder) doReturnResults(t moq.T,
	fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

type noNamesAdaptor struct{ m *moqUsual }

func (a *noNamesAdaptor) exported() bool { return false }

func (a *noNamesAdaptor) tracksParams() bool { return true }

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

type noNamesRecorder struct{ r *moqUsual_NoNames_fnRecorder }

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

func (r *noNamesRecorder) doReturnResults(t moq.T,
	fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

type exportedNoNamesAdaptor struct{ m *exported.MoqUsual }

func (a *exportedNoNamesAdaptor) exported() bool { return true }

func (a *exportedNoNamesAdaptor) tracksParams() bool { return true }

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

func (r *exportedNoNamesRecorder) doReturnResults(t moq.T,
	fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

type noResultsAdaptor struct{ m *moqUsual }

func (a *noResultsAdaptor) exported() bool { return false }

func (a *noResultsAdaptor) tracksParams() bool { return true }

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

func (r *noResultsRecorder) doReturnResults(t moq.T,
	fn func(), expectedSParams []string, expectedBParam bool, _ []string, _ error) {
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

type exportedNoResultsAdaptor struct{ m *exported.MoqUsual }

func (a *exportedNoResultsAdaptor) exported() bool { return true }

func (a *exportedNoResultsAdaptor) tracksParams() bool { return true }

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

func (r *exportedNoResultsRecorder) doReturnResults(t moq.T,
	fn func(), expectedSParams []string, expectedBParam bool, _ []string, _ error) {
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

type noParamsAdaptor struct{ m *moqUsual }

func (a *noParamsAdaptor) exported() bool { return false }

func (a *noParamsAdaptor) tracksParams() bool { return false }

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

func (r *noParamsRecorder) doReturnResults(_ moq.T,
	fn func(), _ []string, _ bool, sResults []string, err error) {
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

type exportedNoParamsAdaptor struct{ m *exported.MoqUsual }

func (a *exportedNoParamsAdaptor) exported() bool { return true }

func (a *exportedNoParamsAdaptor) tracksParams() bool { return false }

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

func (r *exportedNoParamsRecorder) doReturnResults(_ moq.T,
	fn func(), _ []string, _ bool, sResults []string, err error) {
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

type nothingAdaptor struct{ m *moqUsual }

func (a *nothingAdaptor) exported() bool { return false }

func (a *nothingAdaptor) tracksParams() bool { return false }

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

type nothingRecorder struct{ r *moqUsual_Nothing_fnRecorder }

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

func (r *nothingRecorder) doReturnResults(_ moq.T,
	fn func(), _ []string, _ bool, _ []string, _ error) {
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

type exportedNothingAdaptor struct{ m *exported.MoqUsual }

func (a *exportedNothingAdaptor) exported() bool { return true }

func (a *exportedNothingAdaptor) tracksParams() bool { return false }

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

func (r *exportedNothingRecorder) doReturnResults(_ moq.T,
	fn func(), _ []string, _ bool, _ []string, _ error) {
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

type variadicAdaptor struct{ m *moqUsual }

func (a *variadicAdaptor) exported() bool { return false }

func (a *variadicAdaptor) tracksParams() bool { return true }

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

func (r *variadicRecorder) doReturnResults(t moq.T,
	fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

type exportedVariadicAdaptor struct{ m *exported.MoqUsual }

func (a *exportedVariadicAdaptor) exported() bool { return true }

func (a *exportedVariadicAdaptor) tracksParams() bool { return true }

func (a *exportedVariadicAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedVariadicRecorder{r: a.m.OnCall().Variadic(bParam, sParams...)}
}

func (a *exportedVariadicAdaptor) invokeMockAndExpectResults(t moq.T,
	sParams []string, bParam bool, res results) {
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

func (r *exportedVariadicRecorder) doReturnResults(t moq.T,
	fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

type repeatedIdsAdaptor struct{ m *moqUsual }

func (a *repeatedIdsAdaptor) exported() bool { return false }

func (a *repeatedIdsAdaptor) tracksParams() bool { return true }

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

func (r *repeatedIdsRecorder) doReturnResults(t moq.T,
	fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

type exportedRepeatedIdsAdaptor struct{ m *exported.MoqUsual }

func (a *exportedRepeatedIdsAdaptor) exported() bool { return true }

func (a *exportedRepeatedIdsAdaptor) tracksParams() bool { return true }

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

func (r *exportedRepeatedIdsRecorder) doReturnResults(t moq.T,
	fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

type timesAdaptor struct{ m *moqUsual }

func (a *timesAdaptor) exported() bool { return false }

func (a *timesAdaptor) tracksParams() bool { return true }

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

type timesRecorder struct{ r *moqUsual_Times_fnRecorder }

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

func (r *timesRecorder) doReturnResults(t moq.T,
	fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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

type exportedTimesAdaptor struct{ m *exported.MoqUsual }

func (a *exportedTimesAdaptor) exported() bool { return true }

func (a *exportedTimesAdaptor) tracksParams() bool { return true }

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

func (r *exportedTimesRecorder) doReturnResults(t moq.T,
	fn func(), expectedSParams []string, expectedBParam bool, sResults []string, err error) {
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
