package testmoqs_test

import (
	"sync"
	"testing"

	"moqueries.org/runtime/moq"
)

const (
	callsPerRoutine = 100
	routines        = 10
)

func TestAtomicSequences(t *testing.T) {
	t.Skip("This test is used to show how sequences and expected " +
		"call order are not atomic. This test is rarely successful so " +
		"should remain skipped.")

	// ASSEMBLE
	scene := moq.NewScene(t)
	usualMoq := newMoqUsualFn(scene, &moq.Config{
		Sequence: moq.SeqDefaultOn,
	})

	usualMoq.onCall("Hi", false).
		returnResults("Bye", nil).
		// We say we're expecting 10 times more results just in case one
		// routine iterates more quickly. We don't actually expect all calls
		// to be made.
		repeat(moq.Times(callsPerRoutine * routines * 10))

	start := make(chan struct{})
	done := sync.WaitGroup{}
	mockFn := usualMoq.mock()

	for n := 0; n < routines; n++ {
		done.Add(1)
		go func() {
			defer done.Done()

			<-start

			for m := 0; m < callsPerRoutine; m++ {
				res, err := mockFn("Hi", false)
				if err != nil {
					t.Errorf("wanted no err, got %#v", err)
				}
				if res != "Bye" {
					t.Errorf("wanted Bye, got %s", res)
				}
			}
		}()
	}

	// ACT
	close(start)
	done.Wait()

	// ASSERT
}
