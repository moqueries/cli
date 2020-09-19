package testmocks_test

import (
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/pkg/generator/testmocks/exported"
	"github.com/myshkin5/moqueries/pkg/hash"
)

type usualAdaptor struct{ m *mockUsual }

func (a *usualAdaptor) tracksParams() bool { return true }

func (a *usualAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.onCall().Usual(sParams[0], bParam)
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.returnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			rec = rec.times(result.times)
		}
		if result.anyTimes {
			rec.anyTimes()
			rec = nil
		}
	}
	return rec
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

type exportedUsualAdaptor struct{ m *exported.MockUsual }

func (a *exportedUsualAdaptor) tracksParams() bool { return true }

func (a *exportedUsualAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.OnCall().Usual(sParams[0], bParam)
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.ReturnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			rec = rec.Times(result.times)
		}
		if result.anyTimes {
			rec.AnyTimes()
			rec = nil
		}
	}
	return rec
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

type noNamesAdaptor struct{ m *mockUsual }

func (a *noNamesAdaptor) tracksParams() bool { return true }

func (a *noNamesAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.onCall().NoNames(sParams[0], bParam)
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.returnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			rec = rec.times(result.times)
		}
		if result.anyTimes {
			rec.anyTimes()
			rec = nil
		}
	}
	return rec
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

type exportedNoNamesAdaptor struct{ m *exported.MockUsual }

func (a *exportedNoNamesAdaptor) tracksParams() bool { return true }

func (a *exportedNoNamesAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.OnCall().NoNames(sParams[0], bParam)
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.ReturnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			rec = rec.Times(result.times)
		}
		if result.anyTimes {
			rec.AnyTimes()
			rec = nil
		}
	}
	return rec
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

type noResultsAdaptor struct{ m *mockUsual }

func (a *noResultsAdaptor) tracksParams() bool { return true }

func (a *noResultsAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.onCall().NoResults(sParams[0], bParam)
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.returnResults()
		}
		if result.times > 0 {
			rec = rec.times(result.times)
		}
		if result.anyTimes {
			rec.anyTimes()
			rec = nil
		}
	}
	return rec
}

func (a *noResultsAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, _ results) {
	a.m.mock().NoResults(sParams[0], bParam)
}

func (a *noResultsAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return mockUsual_NoResults_params{sParam: sParams[0], bParam: bParam}
}

type exportedNoResultsAdaptor struct{ m *exported.MockUsual }

func (a *exportedNoResultsAdaptor) tracksParams() bool { return true }

func (a *exportedNoResultsAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.OnCall().NoResults(sParams[0], bParam)
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.ReturnResults()
		}
		if result.times > 0 {
			rec = rec.Times(result.times)
		}
		if result.anyTimes {
			rec.AnyTimes()
			rec = nil
		}
	}
	return rec
}

func (a *exportedNoResultsAdaptor) invokeMockAndExpectResults(sParams []string, bParam bool, _ results) {
	a.m.Mock().NoResults(sParams[0], bParam)
}

func (a *exportedNoResultsAdaptor) bundleParams(sParams []string, bParam bool) interface{} {
	return exported.MockUsual_NoResults_params{SParam: sParams[0], BParam: bParam}
}

type noParamsAdaptor struct{ m *mockUsual }

func (a *noParamsAdaptor) tracksParams() bool { return false }

func (a *noParamsAdaptor) expectCall(_ []string, _ bool, results ...results) interface{} {
	rec := a.m.onCall().NoParams()
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.returnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			rec = rec.times(result.times)
		}
		if result.anyTimes {
			rec.anyTimes()
			rec = nil
		}
	}
	return rec
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

type exportedNoParamsAdaptor struct{ m *exported.MockUsual }

func (a *exportedNoParamsAdaptor) tracksParams() bool { return false }

func (a *exportedNoParamsAdaptor) expectCall(_ []string, _ bool, results ...results) interface{} {
	rec := a.m.OnCall().NoParams()
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.ReturnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			rec = rec.Times(result.times)
		}
		if result.anyTimes {
			rec.AnyTimes()
			rec = nil
		}
	}
	return rec
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

type nothingAdaptor struct{ m *mockUsual }

func (a *nothingAdaptor) tracksParams() bool { return false }

func (a *nothingAdaptor) expectCall(_ []string, _ bool, results ...results) interface{} {
	rec := a.m.onCall().Nothing()
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.returnResults()
		}
		if result.times > 0 {
			rec = rec.times(result.times)
		}
		if result.anyTimes {
			rec.anyTimes()
			rec = nil
		}
	}
	return rec
}

func (a *nothingAdaptor) invokeMockAndExpectResults([]string, bool, results) {
	a.m.mock().Nothing()
}

func (a *nothingAdaptor) bundleParams([]string, bool) interface{} {
	return mockUsual_Nothing_params{}
}

type exportedNothingAdaptor struct{ m *exported.MockUsual }

func (a *exportedNothingAdaptor) tracksParams() bool { return false }

func (a *exportedNothingAdaptor) expectCall(_ []string, _ bool, results ...results) interface{} {
	rec := a.m.OnCall().Nothing()
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.ReturnResults()
		}
		if result.times > 0 {
			rec = rec.Times(result.times)
		}
		if result.anyTimes {
			rec.AnyTimes()
			rec = nil
		}
	}
	return rec
}

func (a *exportedNothingAdaptor) invokeMockAndExpectResults([]string, bool, results) {
	a.m.Mock().Nothing()
}

func (a *exportedNothingAdaptor) bundleParams([]string, bool) interface{} {
	return exported.MockUsual_Nothing_params{}
}

type variadicAdaptor struct{ m *mockUsual }

func (a *variadicAdaptor) tracksParams() bool { return true }

func (a *variadicAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.onCall().Variadic(bParam, sParams...)
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.returnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			rec = rec.times(result.times)
		}
		if result.anyTimes {
			rec.anyTimes()
			rec = nil
		}
	}
	return rec
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
	return mockUsual_Variadic_params{args: hash.DeepHash(sParams), other: bParam}
}

type exportedVariadicAdaptor struct{ m *exported.MockUsual }

func (a *exportedVariadicAdaptor) tracksParams() bool { return true }

func (a *exportedVariadicAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.OnCall().Variadic(bParam, sParams...)
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.ReturnResults(result.sResults[0], result.err)
		}
		if result.times > 0 {
			rec = rec.Times(result.times)
		}
		if result.anyTimes {
			rec.AnyTimes()
			rec = nil
		}
	}
	return rec
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
	return exported.MockUsual_Variadic_params{Args: hash.DeepHash(sParams), Other: bParam}
}

type repeatedIdsAdaptor struct{ m *mockUsual }

func (a *repeatedIdsAdaptor) tracksParams() bool { return true }

func (a *repeatedIdsAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.onCall().RepeatedIds(sParams[0], sParams[1], bParam)
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.returnResults(result.sResults[0], result.sResults[1], result.err)
		}
		if result.times > 0 {
			rec = rec.times(result.times)
		}
		if result.anyTimes {
			rec.anyTimes()
			rec = nil
		}
	}
	return rec
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

type exportedRepeatedIdsAdaptor struct{ m *exported.MockUsual }

func (a *exportedRepeatedIdsAdaptor) tracksParams() bool { return true }

func (a *exportedRepeatedIdsAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.OnCall().RepeatedIds(sParams[0], sParams[1], bParam)
	for _, result := range results {
		if !result.noReturnResults {
			rec = rec.ReturnResults(result.sResults[0], result.sResults[1], result.err)
		}
		if result.times > 0 {
			rec = rec.Times(result.times)
		}
		if result.anyTimes {
			rec.AnyTimes()
			rec = nil
		}
	}
	return rec
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
