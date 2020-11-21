package testmocks_test

import (
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/pkg/generator/testmocks/exported"
	"github.com/myshkin5/moqueries/pkg/moq"
)

type usualFnAdaptor struct{ m *mockUsualFn }

func (a *usualFnAdaptor) tracksParams() bool { return true }

func (a *usualFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &usualFnRecorder{r: a.m.onCall(sParams[0], bParam)}
}

func (a *usualFnAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock()(sParams[0], bParam)
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *usualFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return mockUsualFn_params{sParam: sParams[0], bParam: bParam}
}

func (a *usualFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type usualFnRecorder struct{ r *mockUsualFn_fnRecorder }

func (r *usualFnRecorder) anySParam() {
	r.r = r.r.anySParam()
}

func (r *usualFnRecorder) anyBParam() {
	r.r = r.r.anyBParam()
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

func (r *usualFnRecorder) times(count int) {
	r.r = r.r.times(count)
}

func (r *usualFnRecorder) anyTimes() {
	r.r.anyTimes()
	r.r = nil
}

func (r *usualFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedUsualFnAdaptor struct{ m *exported.MockUsualFn }

func (a *exportedUsualFnAdaptor) tracksParams() bool { return true }

func (a *exportedUsualFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedUsualFnRecorder{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedUsualFnAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock()(sParams[0], bParam)
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *exportedUsualFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MockUsualFn_params{SParam: sParams[0], BParam: bParam}
}

func (a *exportedUsualFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedUsualFnRecorder struct {
	r *exported.MockUsualFn_fnRecorder
}

func (r *exportedUsualFnRecorder) anySParam() {
	r.r = r.r.AnySParam()
}

func (r *exportedUsualFnRecorder) anyBParam() {
	r.r = r.r.AnyBParam()
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

func (r *exportedUsualFnRecorder) times(count int) {
	r.r = r.r.Times(count)
}

func (r *exportedUsualFnRecorder) anyTimes() {
	r.r.AnyTimes()
	r.r = nil
}

func (r *exportedUsualFnRecorder) isNil() bool {
	return r.r == nil
}

type noNamesFnAdaptor struct{ m *mockNoNamesFn }

func (a *noNamesFnAdaptor) tracksParams() bool { return true }

func (a *noNamesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &noNamesFnRecorder{r: a.m.onCall(sParams[0], bParam)}
}

func (a *noNamesFnAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock()(sParams[0], bParam)
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *noNamesFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return mockNoNamesFn_params{param1: sParams[0], param2: bParam}
}

func (a *noNamesFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type noNamesFnRecorder struct{ r *mockNoNamesFn_fnRecorder }

func (r *noNamesFnRecorder) anySParam() {
	r.r = r.r.anyParam1()
}

func (r *noNamesFnRecorder) anyBParam() {
	r.r = r.r.anyParam2()
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

func (r *noNamesFnRecorder) times(count int) {
	r.r = r.r.times(count)
}

func (r *noNamesFnRecorder) anyTimes() {
	r.r.anyTimes()
	r.r = nil
}

func (r *noNamesFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedNoNamesFnAdaptor struct{ m *exported.MockNoNamesFn }

func (a *exportedNoNamesFnAdaptor) tracksParams() bool { return true }

func (a *exportedNoNamesFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedNoNamesFnRecorder{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedNoNamesFnAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock()(sParams[0], bParam)
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *exportedNoNamesFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MockNoNamesFn_params{Param1: sParams[0], Param2: bParam}
}

func (a *exportedNoNamesFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedNoNamesFnRecorder struct {
	r *exported.MockNoNamesFn_fnRecorder
}

func (r *exportedNoNamesFnRecorder) anySParam() {
	r.r = r.r.AnyParam1()
}

func (r *exportedNoNamesFnRecorder) anyBParam() {
	r.r = r.r.AnyParam2()
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

func (r *exportedNoNamesFnRecorder) times(count int) {
	r.r = r.r.Times(count)
}

func (r *exportedNoNamesFnRecorder) anyTimes() {
	r.r.AnyTimes()
	r.r = nil
}

func (r *exportedNoNamesFnRecorder) isNil() bool {
	return r.r == nil
}

type noResultsFnAdaptor struct{ m *mockNoResultsFn }

func (a *noResultsFnAdaptor) tracksParams() bool { return true }

func (a *noResultsFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &noResultsFnRecorder{r: a.m.onCall(sParams[0], bParam)}
}

func (a *noResultsFnAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, _ results) {
	a.m.mock()(sParams[0], bParam)
}

func (a *noResultsFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return mockNoResultsFn_params{sParam: sParams[0], bParam: bParam}
}

func (a *noResultsFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type noResultsFnRecorder struct{ r *mockNoResultsFn_fnRecorder }

func (r *noResultsFnRecorder) anySParam() {
	r.r = r.r.anySParam()
}

func (r *noResultsFnRecorder) anyBParam() {
	r.r = r.r.anyBParam()
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

func (r *noResultsFnRecorder) times(count int) {
	r.r = r.r.times(count)
}

func (r *noResultsFnRecorder) anyTimes() {
	r.r.anyTimes()
	r.r = nil
}

func (r *noResultsFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedNoResultsFnAdaptor struct{ m *exported.MockNoResultsFn }

func (a *exportedNoResultsFnAdaptor) tracksParams() bool { return true }

func (a *exportedNoResultsFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedNoResultsFnRecorder{r: a.m.OnCall(sParams[0], bParam)}
}

func (a *exportedNoResultsFnAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, _ results) {
	a.m.Mock()(sParams[0], bParam)
}

func (a *exportedNoResultsFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MockNoResultsFn_params{SParam: sParams[0], BParam: bParam}
}

func (a *exportedNoResultsFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedNoResultsFnRecorder struct {
	r *exported.MockNoResultsFn_fnRecorder
}

func (r *exportedNoResultsFnRecorder) anySParam() {
	r.r = r.r.AnySParam()
}

func (r *exportedNoResultsFnRecorder) anyBParam() {
	r.r = r.r.AnyBParam()
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

func (r *exportedNoResultsFnRecorder) times(count int) {
	r.r = r.r.Times(count)
}

func (r *exportedNoResultsFnRecorder) anyTimes() {
	r.r.AnyTimes()
	r.r = nil
}

func (r *exportedNoResultsFnRecorder) isNil() bool {
	return r.r == nil
}

type noParamsFnAdaptor struct{ m *mockNoParamsFn }

func (a *noParamsFnAdaptor) tracksParams() bool { return false }

func (a *noParamsFnAdaptor) newRecorder([]string, bool) recorder {
	return &noParamsFnRecorder{r: a.m.onCall()}
}

func (a *noParamsFnAdaptor) invokeMockAndExpectResults(_ []string, _ bool, res results) {
	sResult, err := a.m.mock()()
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *noParamsFnAdaptor) bundleParams([]string, bool) interface{} {
	return mockNoParamsFn_params{}
}

func (a *noParamsFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type noParamsFnRecorder struct{ r *mockNoParamsFn_fnRecorder }

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

func (r *noParamsFnRecorder) times(count int) {
	r.r = r.r.times(count)
}

func (r *noParamsFnRecorder) anyTimes() {
	r.r.anyTimes()
	r.r = nil
}

func (r *noParamsFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedNoParamsFnAdaptor struct{ m *exported.MockNoParamsFn }

func (a *exportedNoParamsFnAdaptor) tracksParams() bool { return false }

func (a *exportedNoParamsFnAdaptor) newRecorder([]string, bool) recorder {
	return &exportedNoParamsFnRecorder{r: a.m.OnCall()}
}

func (a *exportedNoParamsFnAdaptor) invokeMockAndExpectResults(_ []string, _ bool, res results) {
	sResult, err := a.m.Mock()()
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *exportedNoParamsFnAdaptor) bundleParams([]string, bool) interface{} {
	return exported.MockNoParamsFn_params{}
}

func (a *exportedNoParamsFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedNoParamsFnRecorder struct {
	r *exported.MockNoParamsFn_fnRecorder
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

func (r *exportedNoParamsFnRecorder) times(count int) {
	r.r = r.r.Times(count)
}

func (r *exportedNoParamsFnRecorder) anyTimes() {
	r.r.AnyTimes()
	r.r = nil
}

func (r *exportedNoParamsFnRecorder) isNil() bool {
	return r.r == nil
}

type nothingFnAdaptor struct{ m *mockNothingFn }

func (a *nothingFnAdaptor) tracksParams() bool { return false }

func (a *nothingFnAdaptor) newRecorder([]string, bool) recorder {
	return &nothingFnRecorder{r: a.m.onCall()}
}

func (a *nothingFnAdaptor) invokeMockAndExpectResults([]string, bool, results) {
	a.m.mock()()
}

func (a *nothingFnAdaptor) bundleParams([]string, bool) interface{} {
	return mockNothingFn_params{}
}

func (a *nothingFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type nothingFnRecorder struct{ r *mockNothingFn_fnRecorder }

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

func (r *nothingFnRecorder) times(count int) {
	r.r = r.r.times(count)
}

func (r *nothingFnRecorder) anyTimes() {
	r.r.anyTimes()
	r.r = nil
}

func (r *nothingFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedNothingFnAdaptor struct{ m *exported.MockNothingFn }

func (a *exportedNothingFnAdaptor) tracksParams() bool { return false }

func (a *exportedNothingFnAdaptor) newRecorder([]string, bool) recorder {
	return &exportedNothingFnRecorder{r: a.m.OnCall()}
}

func (a *exportedNothingFnAdaptor) invokeMockAndExpectResults([]string, bool, results) {
	a.m.Mock()()
}

func (a *exportedNothingFnAdaptor) bundleParams([]string, bool) interface{} {
	return exported.MockNothingFn_params{}
}

func (a *exportedNothingFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedNothingFnRecorder struct {
	r *exported.MockNothingFn_fnRecorder
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

func (r *exportedNothingFnRecorder) times(count int) {
	r.r = r.r.Times(count)
}

func (r *exportedNothingFnRecorder) anyTimes() {
	r.r.AnyTimes()
	r.r = nil
}

func (r *exportedNothingFnRecorder) isNil() bool {
	return r.r == nil
}

type variadicFnAdaptor struct{ m *mockVariadicFn }

func (a *variadicFnAdaptor) tracksParams() bool { return true }

func (a *variadicFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &variadicFnRecorder{r: a.m.onCall(bParam, sParams...)}
}

func (a *variadicFnAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock()(bParam, sParams...)
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *variadicFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return mockVariadicFn_params{args: sParams, other: bParam}
}

func (a *variadicFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type variadicFnRecorder struct{ r *mockVariadicFn_fnRecorder }

func (r *variadicFnRecorder) anySParam() {
	r.r = r.r.anyArgs()
}

func (r *variadicFnRecorder) anyBParam() {
	r.r = r.r.anyOther()
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

func (r *variadicFnRecorder) times(count int) {
	r.r = r.r.times(count)
}

func (r *variadicFnRecorder) anyTimes() {
	r.r.anyTimes()
	r.r = nil
}

func (r *variadicFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedVariadicFnAdaptor struct{ m *exported.MockVariadicFn }

func (a *exportedVariadicFnAdaptor) tracksParams() bool { return true }

func (a *exportedVariadicFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedVariadicFnRecorder{r: a.m.OnCall(bParam, sParams...)}
}

func (a *exportedVariadicFnAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock()(bParam, sParams...)
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *exportedVariadicFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MockVariadicFn_params{Args: sParams, Other: bParam}
}

func (a *exportedVariadicFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedVariadicFnRecorder struct {
	r *exported.MockVariadicFn_fnRecorder
}

func (r *exportedVariadicFnRecorder) anySParam() {
	r.r = r.r.AnyArgs()
}

func (r *exportedVariadicFnRecorder) anyBParam() {
	r.r = r.r.AnyOther()
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

func (r *exportedVariadicFnRecorder) times(count int) {
	r.r = r.r.Times(count)
}

func (r *exportedVariadicFnRecorder) anyTimes() {
	r.r.AnyTimes()
	r.r = nil
}

func (r *exportedVariadicFnRecorder) isNil() bool {
	return r.r == nil
}

type repeatedIdsFnAdaptor struct{ m *mockRepeatedIdsFn }

func (a *repeatedIdsFnAdaptor) tracksParams() bool { return true }

func (a *repeatedIdsFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &repeatedIdsFnRecorder{r: a.m.onCall(sParams[0], sParams[1], bParam)}
}

func (a *repeatedIdsFnAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult1, sResult2, err := a.m.mock()(sParams[0], sParams[1], bParam)
	Expect(sResult1).To(Equal(res.sResults[0]))
	Expect(sResult2).To(Equal(res.sResults[1]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *repeatedIdsFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return mockRepeatedIdsFn_params{sParam1: sParams[0], sParam2: sParams[1], bParam: bParam}
}

func (a *repeatedIdsFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type repeatedIdsFnRecorder struct{ r *mockRepeatedIdsFn_fnRecorder }

func (r *repeatedIdsFnRecorder) anySParam() {
	r.r = r.r.anySParam1()
}

func (r *repeatedIdsFnRecorder) anyBParam() {
	r.r = r.r.anyBParam()
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

func (r *repeatedIdsFnRecorder) times(count int) {
	r.r = r.r.times(count)
}

func (r *repeatedIdsFnRecorder) anyTimes() {
	r.r.anyTimes()
	r.r = nil
}

func (r *repeatedIdsFnRecorder) isNil() bool {
	return r.r == nil
}

type exportedRepeatedIdsFnAdaptor struct{ m *exported.MockRepeatedIdsFn }

func (a *exportedRepeatedIdsFnAdaptor) tracksParams() bool { return true }

func (a *exportedRepeatedIdsFnAdaptor) newRecorder(sParams []string, bParam bool) recorder {
	return &exportedRepeatedIdsFnRecorder{r: a.m.OnCall(sParams[0], sParams[1], bParam)}
}

func (a *exportedRepeatedIdsFnAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult1, sResult2, err := a.m.Mock()(sParams[0], sParams[1], bParam)
	Expect(sResult1).To(Equal(res.sResults[0]))
	Expect(sResult2).To(Equal(res.sResults[1]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *exportedRepeatedIdsFnAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MockRepeatedIdsFn_params{SParam1: sParams[0], SParam2: sParams[1], BParam: bParam}
}

func (a *exportedRepeatedIdsFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedRepeatedIdsFnRecorder struct {
	r *exported.MockRepeatedIdsFn_fnRecorder
}

func (r *exportedRepeatedIdsFnRecorder) anySParam() {
	r.r = r.r.AnySParam1()
}

func (r *exportedRepeatedIdsFnRecorder) anyBParam() {
	r.r = r.r.AnyBParam()
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

func (r *exportedRepeatedIdsFnRecorder) times(count int) {
	r.r = r.r.Times(count)
}

func (r *exportedRepeatedIdsFnRecorder) anyTimes() {
	r.r.AnyTimes()
	r.r = nil
}

func (r *exportedRepeatedIdsFnRecorder) isNil() bool {
	return r.r == nil
}
