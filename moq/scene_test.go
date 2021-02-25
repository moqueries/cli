package moq_test

import (
	"testing"

	"github.com/myshkin5/moqueries/moq"
)

func TestScene(t *testing.T) {
	var (
		scene *moq.Scene
		moq1  *moqMoq
		moq2  *moqMoq
		tMoq  *moq.MoqT
		moqT  moq.T

		testScene *moq.Scene
	)

	beforeEach := func(t *testing.T) {
		if scene != nil {
			t.Fatal("afterEach not called")
		}
		scene = moq.NewScene(t)
		moq1 = newMoqMoq(scene, nil)
		moq2 = newMoqMoq(scene, nil)
		tMoq = moq.NewMoqT(scene, nil)
		moqT = tMoq.Mock()

		testScene = moq.NewScene(moqT)
		testScene.AddMoq(moq1.mock())
		testScene.AddMoq(moq2.mock())
	}

	afterEach := func() {
		scene.AssertExpectationsMet()
		scene = nil
	}

	t.Run("resets all of its moqs", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		moq1.onCall().Reset().returnResults()
		moq2.onCall().Reset().returnResults()

		// ACT
		testScene.Reset()

		// ASSERT
		afterEach()
	})

	t.Run("asserts all of its moqs meet expectations", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		moq1.onCall().AssertExpectationsMet().returnResults()
		moq2.onCall().AssertExpectationsMet().returnResults()

		// ACT
		testScene.AssertExpectationsMet()

		// ASSERT
		afterEach()
	})

	t.Run("returns the same MoqT it is given", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)

		// ACT
		actualMoqT := testScene.T

		// ASSERT
		if actualMoqT != moqT {
			t.Errorf("got %#v, wanted %#v", actualMoqT, moqT)
		}
		afterEach()
	})
}
