package moq_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/moqueries/moq"
)

var _ = Describe("Scene", func() {
	var (
		scene *moq.Scene
		moq1  *moqMoq
		moq2  *moqMoq
		tMoq  *moq.MoqT
		moqT  moq.T

		testScene *moq.Scene
	)

	BeforeEach(func() {
		scene = moq.NewScene(GinkgoT())
		moq1 = newMoqMoq(scene, nil)
		moq2 = newMoqMoq(scene, nil)
		tMoq = moq.NewMoqT(scene, nil)
		moqT = tMoq.Mock()

		testScene = moq.NewScene(moqT)
		testScene.AddMoq(moq1.mock())
		testScene.AddMoq(moq2.mock())
	})

	AfterEach(func() {
		scene.AssertExpectationsMet()
	})

	It("resets all of its moqs", func() {
		// ASSEMBLE
		moq1.onCall().Reset().returnResults()
		moq2.onCall().Reset().returnResults()

		// ACT
		testScene.Reset()

		// ASSERT
	})

	It("asserts all of its moqs meet expectations", func() {
		// ASSEMBLE
		moq1.onCall().AssertExpectationsMet().returnResults()
		moq2.onCall().AssertExpectationsMet().returnResults()

		// ACT
		testScene.AssertExpectationsMet()

		// ASSERT
	})

	It("returns the same MoqT it is given", func() {
		// ASSEMBLE

		// ACT
		actualMoqT := testScene.T

		// ASSERT
		Expect(actualMoqT).To(BeIdenticalTo(moqT))
	})
})
