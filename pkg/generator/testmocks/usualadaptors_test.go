package testmocks_test

import (
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/pkg/generator/testmocks/exported"
	"github.com/myshkin5/moqueries/pkg/moq"
)

type usualAdaptor struct{ m *mockUsual }

func (a *usualAdaptor) tracksParams() bool { return true }

func (a *usualAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.onCall().Usual(sParams[0], bParam)
}

func (a *usualAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*mockUsual_Usual_fnRecorder)
	if sParams {
		cRec = cRec.anySParam()
	}
	if bParam {
		cRec = cRec.anyBParam()
	}
	return cRec
}

func (a *usualAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockUsual_Usual_fnRecorder)
	for _, result := range results {
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

func (a *usualAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock().Usual(sParams[0], bParam)
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *usualAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return mockUsual_Usual_params{sParam: sParams[0], bParam: bParam}
}

func (a *usualAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedUsualAdaptor struct{ m *exported.MockUsual }

func (a *exportedUsualAdaptor) tracksParams() bool { return true }

func (a *exportedUsualAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.OnCall().Usual(sParams[0], bParam)
}

func (a *exportedUsualAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*exported.MockUsual_Usual_fnRecorder)
	if sParams {
		cRec = cRec.AnySParam()
	}
	if bParam {
		cRec = cRec.AnyBParam()
	}
	return cRec
}

func (a *exportedUsualAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockUsual_Usual_fnRecorder)
	for _, result := range results {
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

func (a *exportedUsualAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock().Usual(sParams[0], bParam)
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *exportedUsualAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MockUsual_Usual_params{SParam: sParams[0], BParam: bParam}
}

func (a *exportedUsualAdaptor) sceneMock() moq.Mock {
	return a.m
}

type noNamesAdaptor struct{ m *mockUsual }

func (a *noNamesAdaptor) tracksParams() bool { return true }

func (a *noNamesAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.onCall().NoNames(sParams[0], bParam)
}

func (a *noNamesAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*mockUsual_NoNames_fnRecorder)
	if sParams {
		cRec = cRec.anyParam1()
	}
	if bParam {
		cRec = cRec.anyParam2()
	}
	return cRec
}

func (a *noNamesAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockUsual_NoNames_fnRecorder)
	for _, result := range results {
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

func (a *noNamesAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock().NoNames(sParams[0], bParam)
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *noNamesAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return mockUsual_NoNames_params{param1: sParams[0], param2: bParam}
}

func (a *noNamesAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedNoNamesAdaptor struct{ m *exported.MockUsual }

func (a *exportedNoNamesAdaptor) tracksParams() bool { return true }

func (a *exportedNoNamesAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.OnCall().NoNames(sParams[0], bParam)
}

func (a *exportedNoNamesAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*exported.MockUsual_NoNames_fnRecorder)
	if sParams {
		cRec = cRec.AnyParam1()
	}
	if bParam {
		cRec = cRec.AnyParam2()
	}
	return cRec
}

func (a *exportedNoNamesAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockUsual_NoNames_fnRecorder)
	for _, result := range results {
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

func (a *exportedNoNamesAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock().NoNames(sParams[0], bParam)
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *exportedNoNamesAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MockUsual_NoNames_params{Param1: sParams[0], Param2: bParam}
}

func (a *exportedNoNamesAdaptor) sceneMock() moq.Mock {
	return a.m
}

type noResultsAdaptor struct{ m *mockUsual }

func (a *noResultsAdaptor) tracksParams() bool { return true }

func (a *noResultsAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.onCall().NoResults(sParams[0], bParam)
}

func (a *noResultsAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*mockUsual_NoResults_fnRecorder)
	if sParams {
		cRec = cRec.anySParam()
	}
	if bParam {
		cRec = cRec.anyBParam()
	}
	return cRec
}

func (a *noResultsAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockUsual_NoResults_fnRecorder)
	for _, result := range results {
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

func (a *noResultsAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, _ results) {
	a.m.mock().NoResults(sParams[0], bParam)
}

func (a *noResultsAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return mockUsual_NoResults_params{sParam: sParams[0], bParam: bParam}
}

func (a *noResultsAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedNoResultsAdaptor struct{ m *exported.MockUsual }

func (a *exportedNoResultsAdaptor) tracksParams() bool { return true }

func (a *exportedNoResultsAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.OnCall().NoResults(sParams[0], bParam)
}

func (a *exportedNoResultsAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*exported.MockUsual_NoResults_fnRecorder)
	if sParams {
		cRec = cRec.AnySParam()
	}
	if bParam {
		cRec = cRec.AnyBParam()
	}
	return cRec
}

func (a *exportedNoResultsAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockUsual_NoResults_fnRecorder)
	for _, result := range results {
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

func (a *exportedNoResultsAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, _ results) {
	a.m.Mock().NoResults(sParams[0], bParam)
}

func (a *exportedNoResultsAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MockUsual_NoResults_params{SParam: sParams[0], BParam: bParam}
}

func (a *exportedNoResultsAdaptor) sceneMock() moq.Mock {
	return a.m
}

type noParamsAdaptor struct{ m *mockUsual }

func (a *noParamsAdaptor) tracksParams() bool { return false }

func (a *noParamsAdaptor) newRecorder([]string, bool) interface{} {
	return a.m.onCall().NoParams()
}

func (a *noParamsAdaptor) any(rec interface{}, _, _ bool) interface{} {
	return rec
}

func (a *noParamsAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockUsual_NoParams_fnRecorder)
	for _, result := range results {
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

func (a *noParamsAdaptor) invokeMockAndExpectResults(_ []string, _ bool, res results) {
	sResult, err := a.m.mock().NoParams()
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *noParamsAdaptor) bundleParams([]string, bool) interface{} {
	return mockUsual_NoParams_params{}
}

func (a *noParamsAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedNoParamsAdaptor struct{ m *exported.MockUsual }

func (a *exportedNoParamsAdaptor) tracksParams() bool { return false }

func (a *exportedNoParamsAdaptor) newRecorder([]string, bool) interface{} {
	return a.m.OnCall().NoParams()
}

func (a *exportedNoParamsAdaptor) any(rec interface{}, _, _ bool) interface{} {
	return rec
}

func (a *exportedNoParamsAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockUsual_NoParams_fnRecorder)
	for _, result := range results {
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

func (a *exportedNoParamsAdaptor) invokeMockAndExpectResults(_ []string, _ bool, res results) {
	sResult, err := a.m.Mock().NoParams()
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *exportedNoParamsAdaptor) bundleParams([]string, bool) interface{} {
	return exported.MockUsual_NoParams_params{}
}

func (a *exportedNoParamsAdaptor) sceneMock() moq.Mock {
	return a.m
}

type nothingAdaptor struct{ m *mockUsual }

func (a *nothingAdaptor) tracksParams() bool { return false }

func (a *nothingAdaptor) newRecorder([]string, bool) interface{} {
	return a.m.onCall().Nothing()
}

func (a *nothingAdaptor) any(rec interface{}, _, _ bool) interface{} {
	return rec
}

func (a *nothingAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockUsual_Nothing_fnRecorder)
	for _, result := range results {
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

func (a *nothingAdaptor) invokeMockAndExpectResults([]string, bool, results) {
	a.m.mock().Nothing()
}

func (a *nothingAdaptor) bundleParams([]string, bool) interface{} {
	return mockUsual_Nothing_params{}
}

func (a *nothingAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedNothingAdaptor struct{ m *exported.MockUsual }

func (a *exportedNothingAdaptor) tracksParams() bool { return false }

func (a *exportedNothingAdaptor) newRecorder([]string, bool) interface{} {
	return a.m.OnCall().Nothing()
}

func (a *exportedNothingAdaptor) any(rec interface{}, _, _ bool) interface{} {
	return rec
}

func (a *exportedNothingAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockUsual_Nothing_fnRecorder)
	for _, result := range results {
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

func (a *exportedNothingAdaptor) invokeMockAndExpectResults([]string, bool, results) {
	a.m.Mock().Nothing()
}

func (a *exportedNothingAdaptor) bundleParams([]string, bool) interface{} {
	return exported.MockUsual_Nothing_params{}
}

func (a *exportedNothingAdaptor) sceneMock() moq.Mock {
	return a.m
}

type variadicAdaptor struct{ m *mockUsual }

func (a *variadicAdaptor) tracksParams() bool { return true }

func (a *variadicAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.onCall().Variadic(bParam, sParams...)
}

func (a *variadicAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*mockUsual_Variadic_fnRecorder)
	if sParams {
		cRec = cRec.anyArgs()
	}
	if bParam {
		cRec = cRec.anyOther()
	}
	return cRec
}

func (a *variadicAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockUsual_Variadic_fnRecorder)
	for _, result := range results {
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

func (a *variadicAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult, err := a.m.mock().Variadic(bParam, sParams...)
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *variadicAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return mockUsual_Variadic_params{args: sParams, other: bParam}
}

func (a *variadicAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedVariadicAdaptor struct{ m *exported.MockUsual }

func (a *exportedVariadicAdaptor) tracksParams() bool { return true }

func (a *exportedVariadicAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.OnCall().Variadic(bParam, sParams...)
}

func (a *exportedVariadicAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*exported.MockUsual_Variadic_fnRecorder)
	if sParams {
		cRec = cRec.AnyArgs()
	}
	if bParam {
		cRec = cRec.AnyOther()
	}
	return cRec
}

func (a *exportedVariadicAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockUsual_Variadic_fnRecorder)
	for _, result := range results {
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

func (a *exportedVariadicAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult, err := a.m.Mock().Variadic(bParam, sParams...)
	Expect(sResult).To(Equal(res.sResults[0]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *exportedVariadicAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MockUsual_Variadic_params{Args: sParams, Other: bParam}
}

func (a *exportedVariadicAdaptor) sceneMock() moq.Mock {
	return a.m
}

type repeatedIdsAdaptor struct{ m *mockUsual }

func (a *repeatedIdsAdaptor) tracksParams() bool { return true }

func (a *repeatedIdsAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.onCall().RepeatedIds(sParams[0], sParams[1], bParam)
}

func (a *repeatedIdsAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*mockUsual_RepeatedIds_fnRecorder)
	if sParams {
		cRec = cRec.anySParam1()
	}
	if bParam {
		cRec = cRec.anyBParam()
	}
	return cRec
}

func (a *repeatedIdsAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*mockUsual_RepeatedIds_fnRecorder)
	for _, result := range results {
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

func (a *repeatedIdsAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult1, sResult2, err := a.m.mock().RepeatedIds(sParams[0], sParams[1], bParam)
	Expect(sResult1).To(Equal(res.sResults[0]))
	Expect(sResult2).To(Equal(res.sResults[1]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *repeatedIdsAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return mockUsual_RepeatedIds_params{sParam1: sParams[0], sParam2: sParams[1], bParam: bParam}
}

func (a *repeatedIdsAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedRepeatedIdsAdaptor struct{ m *exported.MockUsual }

func (a *exportedRepeatedIdsAdaptor) tracksParams() bool { return true }

func (a *exportedRepeatedIdsAdaptor) newRecorder(sParams []string, bParam bool) interface{} {
	return a.m.OnCall().RepeatedIds(sParams[0], sParams[1], bParam)
}

func (a *exportedRepeatedIdsAdaptor) any(rec interface{}, sParams, bParam bool) interface{} {
	cRec := rec.(*exported.MockUsual_RepeatedIds_fnRecorder)
	if sParams {
		cRec = cRec.AnySParam1()
	}
	if bParam {
		cRec = cRec.AnyBParam()
	}
	return cRec
}

func (a *exportedRepeatedIdsAdaptor) results(rec interface{}, results ...results) interface{} {
	cRec := rec.(*exported.MockUsual_RepeatedIds_fnRecorder)
	for _, result := range results {
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

func (a *exportedRepeatedIdsAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, res results) {
	sResult1, sResult2, err := a.m.Mock().RepeatedIds(sParams[0], sParams[1], bParam)
	Expect(sResult1).To(Equal(res.sResults[0]))
	Expect(sResult2).To(Equal(res.sResults[1]))
	if res.err == nil {
		Expect(err).To(BeNil())
	} else {
		Expect(err).To(Equal(res.err))
	}
}

func (a *exportedRepeatedIdsAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MockUsual_RepeatedIds_params{SParam1: sParams[0], SParam2: sParams[1], BParam: bParam}
}

func (a *exportedRepeatedIdsAdaptor) sceneMock() moq.Mock {
	return a.m
}
