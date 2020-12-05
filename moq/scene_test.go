package moq_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/moq"
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
	})

	It("asserts all of its mocks meet expectations", func() {
		// ASSEMBLE
		mock1.onCall().AssertExpectationsMet().returnResults()
		mock2.onCall().AssertExpectationsMet().returnResults()

		// ACT
		testScene.AssertExpectationsMet()

		// ASSERT
	})

	It("returns the same MoqT it is given", func() {
		// ASSEMBLE

		// ACT
		actualMoqT := testScene.MoqT

		// ASSERT
		Expect(actualMoqT).To(BeIdenticalTo(moqT))
	})
})
