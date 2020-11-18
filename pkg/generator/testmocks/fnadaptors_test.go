package testmocks_test

import (
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/pkg/generator/testmocks/exported"
	"github.com/myshkin5/moqueries/pkg/moq"
)

type usualFnAdaptor struct{ m *mockUsualFn }

func (a *usualFnAdaptor) tracksParams() bool { return true }

func (a *usualFnAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.onCall(sParams[0], bParam)
}

func (a *usualFnAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*mockUsualFn_fnRecorder)
	if sParams {
		cRec = cRec.anySParam()
	}
	if bParam {
		cRec = cRec.anyBParam()
	}
	return cRec
}

func (a *usualFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockUsualFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.seq()
		}
		if result.noSeq {
			cRec = cRec.noSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.returnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			cRec = cRec.times(result.times)
		}
		if result.anyTimes {
			cRec.anyTimes()
			cRec = nil
		}
	}
	return cRec
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

type exportedUsualFnAdaptor struct{ m *exported.MockUsualFn }

func (a *exportedUsualFnAdaptor) tracksParams() bool { return true }

func (a *exportedUsualFnAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.OnCall(sParams[0], bParam)
}

func (a *exportedUsualFnAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*exported.MockUsualFn_fnRecorder)
	if sParams {
		cRec = cRec.AnySParam()
	}
	if bParam {
		cRec = cRec.AnyBParam()
	}
	return cRec
}

func (a *exportedUsualFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockUsualFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.Seq()
		}
		if result.noSeq {
			cRec = cRec.NoSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.ReturnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			cRec = cRec.Times(result.times)
		}
		if result.anyTimes {
			cRec.AnyTimes()
			cRec = nil
		}
	}
	return cRec
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

type noNamesFnAdaptor struct{ m *mockNoNamesFn }

func (a *noNamesFnAdaptor) tracksParams() bool { return true }

func (a *noNamesFnAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.onCall(sParams[0], bParam)
}

func (a *noNamesFnAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*mockNoNamesFn_fnRecorder)
	if sParams {
		cRec = cRec.anyParam1()
	}
	if bParam {
		cRec = cRec.anyParam2()
	}
	return cRec
}

func (a *noNamesFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockNoNamesFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.seq()
		}
		if result.noSeq {
			cRec = cRec.noSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.returnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			cRec = cRec.times(result.times)
		}
		if result.anyTimes {
			cRec.anyTimes()
			cRec = nil
		}
	}
	return cRec
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

type exportedNoNamesFnAdaptor struct{ m *exported.MockNoNamesFn }

func (a *exportedNoNamesFnAdaptor) tracksParams() bool { return true }

func (a *exportedNoNamesFnAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.OnCall(sParams[0], bParam)
}

func (a *exportedNoNamesFnAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*exported.MockNoNamesFn_fnRecorder)
	if sParams {
		cRec = cRec.AnyParam1()
	}
	if bParam {
		cRec = cRec.AnyParam2()
	}
	return cRec
}

func (a *exportedNoNamesFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockNoNamesFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.Seq()
		}
		if result.noSeq {
			cRec = cRec.NoSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.ReturnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			cRec = cRec.Times(result.times)
		}
		if result.anyTimes {
			cRec.AnyTimes()
			cRec = nil
		}
	}
	return cRec
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

type noResultsFnAdaptor struct{ m *mockNoResultsFn }

func (a *noResultsFnAdaptor) tracksParams() bool { return true }

func (a *noResultsFnAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.onCall(sParams[0], bParam)
}

func (a *noResultsFnAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*mockNoResultsFn_fnRecorder)
	if sParams {
		cRec = cRec.anySParam()
	}
	if bParam {
		cRec = cRec.anyBParam()
	}
	return cRec
}

func (a *noResultsFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockNoResultsFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.seq()
		}
		if result.noSeq {
			cRec = cRec.noSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.returnResults()
		}
		if result.times > 0 {
			cRec = cRec.times(result.times)
		}
		if result.anyTimes {
			cRec.anyTimes()
			cRec = nil
		}
	}
	return cRec
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

type exportedNoResultsFnAdaptor struct{ m *exported.MockNoResultsFn }

func (a *exportedNoResultsFnAdaptor) tracksParams() bool { return true }

func (a *exportedNoResultsFnAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.OnCall(sParams[0], bParam)
}

func (a *exportedNoResultsFnAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*exported.MockNoResultsFn_fnRecorder)
	if sParams {
		cRec = cRec.AnySParam()
	}
	if bParam {
		cRec = cRec.AnyBParam()
	}
	return cRec
}

func (a *exportedNoResultsFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockNoResultsFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.Seq()
		}
		if result.noSeq {
			cRec = cRec.NoSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.ReturnResults()
		}
		if result.times > 0 {
			cRec = cRec.Times(result.times)
		}
		if result.anyTimes {
			cRec.AnyTimes()
			cRec = nil
		}
	}
	return cRec
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

type noParamsFnAdaptor struct{ m *mockNoParamsFn }

func (a *noParamsFnAdaptor) tracksParams() bool { return false }

func (a *noParamsFnAdaptor) newRecorder([]string, bool) interface{} {
	return a.m.onCall()
}

func (a *noParamsFnAdaptor) any(rec interface{}, _, _ bool) interface{} {
	return rec
}

func (a *noParamsFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockNoParamsFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.seq()
		}
		if result.noSeq {
			cRec = cRec.noSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.returnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			cRec = cRec.times(result.times)
		}
		if result.anyTimes {
			cRec.anyTimes()
			cRec = nil
		}
	}
	return cRec
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

type exportedNoParamsFnAdaptor struct{ m *exported.MockNoParamsFn }

func (a *exportedNoParamsFnAdaptor) tracksParams() bool { return false }

func (a *exportedNoParamsFnAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.OnCall()
}

func (a *exportedNoParamsFnAdaptor) any(rec interface{}, _, _ bool) interface{} {
	return rec
}

func (a *exportedNoParamsFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockNoParamsFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.Seq()
		}
		if result.noSeq {
			cRec = cRec.NoSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.ReturnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			cRec = cRec.Times(result.times)
		}
		if result.anyTimes {
			cRec.AnyTimes()
			cRec = nil
		}
	}
	return cRec
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

type nothingFnAdaptor struct{ m *mockNothingFn }

func (a *nothingFnAdaptor) tracksParams() bool { return false }

func (a *nothingFnAdaptor) newRecorder([]string, bool) interface{} {
	return a.m.onCall()
}

func (a *nothingFnAdaptor) any(rec interface{}, _, _ bool) interface{} {
	return rec
}

func (a *nothingFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockNothingFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.seq()
		}
		if result.noSeq {
			cRec = cRec.noSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.returnResults()
		}
		if result.times > 0 {
			cRec = cRec.times(result.times)
		}
		if result.anyTimes {
			cRec.anyTimes()
			cRec = nil
		}
	}
	return cRec
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

type exportedNothingFnAdaptor struct{ m *exported.MockNothingFn }

func (a *exportedNothingFnAdaptor) tracksParams() bool { return false }

func (a *exportedNothingFnAdaptor) newRecorder([]string, bool) interface{} {
	return a.m.OnCall()
}

func (a *exportedNothingFnAdaptor) any(rec interface{}, _, _ bool) interface{} {
	return rec
}

func (a *exportedNothingFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockNothingFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.Seq()
		}
		if result.noSeq {
			cRec = cRec.NoSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.ReturnResults()
		}
		if result.times > 0 {
			cRec = cRec.Times(result.times)
		}
		if result.anyTimes {
			cRec.AnyTimes()
			cRec = nil
		}
	}
	return cRec
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

type variadicFnAdaptor struct{ m *mockVariadicFn }

func (a *variadicFnAdaptor) tracksParams() bool { return true }

func (a *variadicFnAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.onCall(bParam, sParams...)
}

func (a *variadicFnAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*mockVariadicFn_fnRecorder)
	if sParams {
		cRec = cRec.anyArgs()
	}
	if bParam {
		cRec = cRec.anyOther()
	}
	return cRec
}

func (a *variadicFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockVariadicFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.seq()
		}
		if result.noSeq {
			cRec = cRec.noSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.returnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			cRec = cRec.times(result.times)
		}
		if result.anyTimes {
			cRec.anyTimes()
			cRec = nil
		}
	}
	return cRec
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

type exportedVariadicFnAdaptor struct{ m *exported.MockVariadicFn }

func (a *exportedVariadicFnAdaptor) tracksParams() bool { return true }

func (a *exportedVariadicFnAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.OnCall(bParam, sParams...)
}

func (a *exportedVariadicFnAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*exported.MockVariadicFn_fnRecorder)
	if sParams {
		cRec = cRec.AnyArgs()
	}
	if bParam {
		cRec = cRec.AnyOther()
	}
	return cRec
}

func (a *exportedVariadicFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockVariadicFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.Seq()
		}
		if result.noSeq {
			cRec = cRec.NoSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.ReturnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			cRec = cRec.Times(result.times)
		}
		if result.anyTimes {
			cRec.AnyTimes()
			cRec = nil
		}
	}
	return cRec
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

type repeatedIdsFnAdaptor struct{ m *mockRepeatedIdsFn }

func (a *repeatedIdsFnAdaptor) tracksParams() bool { return true }

func (a *repeatedIdsFnAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.onCall(sParams[0], sParams[1], bParam)
}

func (a *repeatedIdsFnAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*mockRepeatedIdsFn_fnRecorder)
	if sParams {
		cRec = cRec.anySParam1()
	}
	if bParam {
		cRec = cRec.anyBParam()
	}
	return cRec
}

func (a *repeatedIdsFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockRepeatedIdsFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.seq()
		}
		if result.noSeq {
			cRec = cRec.noSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.returnResults(result.sResults[0], result.sResults[1], result.err)
		}
		if result.times > 0 {
			cRec = cRec.times(result.times)
		}
		if result.anyTimes {
			cRec.anyTimes()
			cRec = nil
		}
	}
	return cRec
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

type exportedRepeatedIdsFnAdaptor struct{ m *exported.MockRepeatedIdsFn }

func (a *exportedRepeatedIdsFnAdaptor) tracksParams() bool { return true }

func (a *exportedRepeatedIdsFnAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.OnCall(sParams[0], sParams[1], bParam)
}

func (a *exportedRepeatedIdsFnAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*exported.MockRepeatedIdsFn_fnRecorder)
	if sParams {
		cRec = cRec.AnySParam1()
	}
	if bParam {
		cRec = cRec.AnyBParam()
	}
	return cRec
}

func (a *exportedRepeatedIdsFnAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockRepeatedIdsFn_fnRecorder)
	for _, result := range results {
		if result.seq {
			cRec = cRec.Seq()
		}
		if result.noSeq {
			cRec = cRec.NoSeq()
		}
		if !result.noReturnResults {
			cRec = cRec.ReturnResults(result.sResults[0], result.sResults[1], result.err)
		}
		if result.times > 0 {
			cRec = cRec.Times(result.times)
		}
		if result.anyTimes {
			cRec.AnyTimes()
			cRec = nil
		}
	}
	return cRec
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
