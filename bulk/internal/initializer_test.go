package internal_test

import (
	"errors"
	"go/build"
	"testing"

	"moqueries.org/runtime/moq"

	"moqueries.org/cli/bulk/internal"
)

var goodInitLine = "{\"root-dir\":\"/my-root-dir\",\"go-path\":\"" + build.Default.GOPATH + "\"}"

func TestInitialize(t *testing.T) {
	var (
		scene          *moq.Scene
		createFnMoq    *moqCreateFn
		writeCloserMoq *moqWriteCloser
	)

	beforeEach := func(t *testing.T) {
		t.Helper()
		if scene != nil {
			t.Fatal("afterEach not called")
		}
		scene = moq.NewScene(t)
		config := &moq.Config{Sequence: moq.SeqDefaultOn}
		createFnMoq = newMoqCreateFn(scene, config)
		writeCloserMoq = newMoqWriteCloser(scene, config)
	}

	afterEach := func() {
		scene.AssertExpectationsMet()
		scene = nil
	}

	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()
		createFnMoq.onCall("some-path-to-state").
			returnResults(writeCloserMoq.mock(), nil)
		writeCloserMoq.onCall().Write([]byte(goodInitLine)).returnResults(0, nil)
		writeCloserMoq.onCall().Write([]byte("\n")).returnResults(0, nil)
		writeCloserMoq.onCall().Close().returnResults(errors.New("ignored-error"))

		// ACT
		err := internal.Initialize(
			"some-path-to-state", "/my-root-dir", createFnMoq.mock())
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted no error", err)
		}
	})

	t.Run("create error", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()
		createFnMoq.onCall("some-path-to-state").
			returnResults(nil, errors.New("create-error"))

		// ACT
		err := internal.Initialize(
			"some-path-to-state", "root-dir", createFnMoq.mock())

		// ASSERT
		if err == nil {
			t.Fatalf("got %#v, wanted an error", err)
		}
		expectedErrMsg := "error creating state file: create-error"
		if err.Error() != expectedErrMsg {
			t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
		}
	})

	t.Run("write error", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()
		createFnMoq.onCall("some-path-to-state").
			returnResults(writeCloserMoq.mock(), nil)
		writeCloserMoq.onCall().Write([]byte(goodInitLine)).
			returnResults(0, errors.New("write-error"))
		writeCloserMoq.onCall().Close().returnResults(nil)

		// ACT
		err := internal.Initialize(
			"some-path-to-state", "/my-root-dir", createFnMoq.mock())

		// ASSERT
		if err == nil {
			t.Fatalf("got %#v, wanted an error", err)
		}
		expectedErrMsg := "error writing state file some-path-to-state: write-error"
		if err.Error() != expectedErrMsg {
			t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
		}
	})

	t.Run("final write error", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()
		createFnMoq.onCall("some-path-to-state").
			returnResults(writeCloserMoq.mock(), nil)
		writeCloserMoq.onCall().Write([]byte(goodInitLine)).
			returnResults(0, nil)
		writeCloserMoq.onCall().Write([]byte("\n")).
			returnResults(0, errors.New("write-error"))
		writeCloserMoq.onCall().Close().returnResults(nil)

		// ACT
		err := internal.Initialize(
			"some-path-to-state", "/my-root-dir", createFnMoq.mock())

		// ASSERT
		if err == nil {
			t.Fatalf("got %#v, wanted an error", err)
		}
		expectedErrMsg := "error finishing writing of state file some-path-to-state: write-error"
		if err.Error() != expectedErrMsg {
			t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
		}
	})
}
