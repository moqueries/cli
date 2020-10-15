package testmocks_test

import (
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/pkg/generator/testmocks/exported"
	"github.com/myshkin5/moqueries/pkg/hash"
	"github.com/myshkin5/moqueries/pkg/moq"
)

type usualFnAdaptor struct{ m *mockUsualFn }

func (a *usualFnAdaptor) tracksParams() bool { return true }

func (a *usualFnAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.onCall(sParams[0], bParam)
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

func (a *exportedUsualFnAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.OnCall(sParams[0], bParam)
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

func (a *noNamesFnAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.onCall(sParams[0], bParam)
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
	return mockNoNamesFn_params{sParam: sParams[0], bParam: bParam}
}

func (a *noNamesFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedNoNamesFnAdaptor struct{ m *exported.MockNoNamesFn }

func (a *exportedNoNamesFnAdaptor) tracksParams() bool { return true }

func (a *exportedNoNamesFnAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.OnCall(sParams[0], bParam)
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
	return exported.MockNoNamesFn_params{SParam: sParams[0], BParam: bParam}
}

func (a *exportedNoNamesFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type noResultsFnAdaptor struct{ m *mockNoResultsFn }

func (a *noResultsFnAdaptor) tracksParams() bool { return true }

func (a *noResultsFnAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.onCall(sParams[0], bParam)
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

func (a *exportedNoResultsFnAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.OnCall(sParams[0], bParam)
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

func (a *noParamsFnAdaptor) expectCall(_ []string, _ bool, results ...results) interface{} {
	rec := a.m.onCall()
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

func (a *exportedNoParamsFnAdaptor) expectCall(_ []string, _ bool, results ...results) interface{} {
	rec := a.m.OnCall()
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

func (a *nothingFnAdaptor) expectCall(_ []string, _ bool, results ...results) interface{} {
	rec := a.m.onCall()
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

func (a *exportedNothingFnAdaptor) expectCall(_ []string, _ bool, results ...results) interface{} {
	rec := a.m.OnCall()
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

func (a *variadicFnAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.onCall(bParam, sParams...)
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
	return mockVariadicFn_params{args: hash.DeepHash(sParams), other: bParam}
}

func (a *variadicFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type exportedVariadicFnAdaptor struct{ m *exported.MockVariadicFn }

func (a *exportedVariadicFnAdaptor) tracksParams() bool { return true }

func (a *exportedVariadicFnAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.OnCall(bParam, sParams...)
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
	return exported.MockVariadicFn_params{Args: hash.DeepHash(sParams), Other: bParam}
}

func (a *exportedVariadicFnAdaptor) sceneMock() moq.Mock {
	return a.m
}

type repeatedIdsFnAdaptor struct{ m *mockRepeatedIdsFn }

func (a *repeatedIdsFnAdaptor) tracksParams() bool { return true }

func (a *repeatedIdsFnAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.onCall(sParams[0], sParams[1], bParam)
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

func (a *exportedRepeatedIdsFnAdaptor) expectCall(sParams []string, bParam bool, results ...results) interface{} {
	rec := a.m.OnCall(sParams[0], sParams[1], bParam)
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
