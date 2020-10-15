package moq_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/pkg/moq"
)

var _ = Describe("Scene", func() {
	var (
		scene    *moq.Scene
		mock1    *mockMock
		mock2    *mockMock
		mockMoqT *moq.MockMoqT
		moqT     moq.MoqT

		testScene *moq.Scene
	)

	BeforeEach(func() {
		scene = moq.NewScene(GinkgoT())
		mock1 = newMockMock(scene, nil)
		mock2 = newMockMock(scene, nil)
		mockMoqT = moq.NewMockMoqT(scene, nil)
		moqT = mockMoqT.Mock()

		testScene = moq.NewScene(moqT)
		testScene.AddMock(mock1.mock())
		testScene.AddMock(mock2.mock())
	})

	AfterEach(func() {
		scene.AssertExpectationsMet()
	})

	It("resets all of its mocks", func() {
		// ASSEMBLE
		mock1.onCall().Reset().returnResults()
		mock2.onCall().Reset().returnResults()

		// ACT
		testScene.Reset()

		// ASSERT
		Expect(mock1.resultsByParams_Reset[mockMock_Reset_params{}].index).
			To(Equal(uint32(1)))
		Expect(mock2.resultsByParams_Reset[mockMock_Reset_params{}].index).
			To(Equal(uint32(1)))
	})

	It("asserts all of its mocks meet expectations", func() {
		// ASSEMBLE
		mock1.onCall().AssertExpectationsMet().returnResults()
		mock2.onCall().AssertExpectationsMet().returnResults()

		// ACT
		testScene.AssertExpectationsMet()

		// ASSERT
		Expect(mock1.resultsByParams_AssertExpectationsMet[mockMock_AssertExpectationsMet_params{}].index).
			To(Equal(uint32(1)))
		Expect(mock2.resultsByParams_AssertExpectationsMet[mockMock_AssertExpectationsMet_params{}].index).
			To(Equal(uint32(1)))
	})

	It("returns the same MoqT it is given", func() {
		// ASSEMBLE

		// ACT
		actualMoqT := testScene.MoqT

		// ASSERT
		Expect(actualMoqT).To(BeIdenticalTo(moqT))
	})
})
