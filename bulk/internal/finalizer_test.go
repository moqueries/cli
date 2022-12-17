package internal_test

import (
	"errors"
	"fmt"
	"go/build"
	"io"
	"testing"

	"moqueries.org/cli/bulk/internal"
	"moqueries.org/cli/generator"
	"moqueries.org/cli/moq"
)

func TestFinalize(t *testing.T) {
	var (
		scene         *moq.Scene
		openFnMoq     *moqOpenFn
		readCloserMoq *moqReadCloser
		generateFnMoq *moqGenerateFn

		appendLine1, appendLine2 string
		req1, req2               generator.GenerateRequest
	)

	beforeEach := func(t *testing.T) {
		t.Helper()
		if scene != nil {
			t.Fatal("afterEach not called")
		}
		scene = moq.NewScene(t)
		config := &moq.Config{Sequence: moq.SeqDefaultOn}
		openFnMoq = newMoqOpenFn(scene, config)
		readCloserMoq = newMoqReadCloser(scene, config)
		generateFnMoq = newMoqGenerateFn(scene, config)

		appendLine1 = "{" +
			"\"types\":[\"Type1\"]," +
			"\"export\":false," +
			"\"destination\":\"\"," +
			"\"destination-dir\":\"\"," +
			"\"package\":\"pkg1\"," +
			"\"import\":\"\"," +
			"\"test-import\":false," +
			"\"working-dir\":\"/my-root-dir\"" +
			"}"
		appendLine2 = "{" +
			"\"types\":[\"Type2\"]," +
			"\"export\":false," +
			"\"destination\":\"\"," +
			"\"destination-dir\":\"\"," +
			"\"package\":\"pkg2\"," +
			"\"import\":\"\"," +
			"\"test-import\":false," +
			"\"working-dir\":\"/my-root-dir\"" +
			"}"
		req1 = generator.GenerateRequest{
			Types:      []string{"Type1"},
			Package:    "pkg1",
			WorkingDir: "/my-root-dir",
		}
		req2 = generator.GenerateRequest{
			Types:      []string{"Type2"},
			Package:    "pkg2",
			WorkingDir: "/my-root-dir",
		}
	}

	afterEach := func() {
		scene.AssertExpectationsMet()
		scene = nil
	}

	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()
		openFnMoq.onCall("some-path-to-state").
			returnResults(readCloserMoq.mock(), nil)
		readCloserMoq.onCall().Read(nil).any().p().
			doReturnResults(readFn(goodInitLine)).
			doReturnResults(readFn(appendLine1)).
			doReturnResults(readFn(appendLine2)).
			doReturnResults(func(p []byte) (int, error) {
				return 0, io.EOF
			})
		generateFnMoq.onCall(req1, req2).returnResults(nil)
		readCloserMoq.onCall().Close().returnResults(errors.New("ignored-error"))

		// ACT
		err := internal.Finalize(
			"some-path-to-state",
			"/my-root-dir",
			openFnMoq.mock(),
			generateFnMoq.mock())
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted no error", err)
		}
	})

	t.Run("verifyState", func(t *testing.T) {
		t.Run("open error", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)
			defer afterEach()
			openFnMoq.onCall("some-path-to-state").
				returnResults(nil, errors.New("open-error"))

			// ACT
			err := internal.Finalize(
				"some-path-to-state",
				"/my-root-dir",
				openFnMoq.mock(),
				generateFnMoq.mock())

			// ASSERT
			if err == nil {
				t.Fatalf("got %#v, wanted an error", err)
			}
			expectedErrMsg := "error opening state file: open-error"
			if err.Error() != expectedErrMsg {
				t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
			}
		})

		t.Run("nothing to read error", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)
			defer afterEach()
			openFnMoq.onCall("some-path-to-state").
				returnResults(readCloserMoq.mock(), nil)
			readCloserMoq.onCall().Read(nil).any().p().
				doReturnResults(func(p []byte) (int, error) {
					return 0, io.EOF
				}).repeat(moq.AnyTimes())
			readCloserMoq.onCall().Close().returnResults(errors.New("ignored-error"))

			// ACT
			err := internal.Finalize(
				"some-path-to-state",
				"/my-root-dir",
				openFnMoq.mock(),
				generateFnMoq.mock())

			// ASSERT
			if err == nil {
				t.Fatalf("got %#v, wanted an error", err)
			}
			if !errors.Is(err, internal.ErrBulkState) {
				t.Errorf("got %#v, want %#v", err, internal.ErrBulkState)
			}
			expectedErrMsg := "bulk state error: state file some-path-to-state not" +
				" initialized properly"
			if err.Error() != expectedErrMsg {
				t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
			}
		})

		t.Run("read error", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)
			defer afterEach()
			openFnMoq.onCall("some-path-to-state").
				returnResults(readCloserMoq.mock(), nil)
			readCloserMoq.onCall().Read(nil).any().p().
				doReturnResults(func(p []byte) (int, error) {
					// Missing proper line ending causes a read error
					copy(p, goodInitLine)
					return len(goodInitLine), nil
				}).
				doReturnResults(func(p []byte) (int, error) {
					return 0, errors.New("bad-file")
				})
			readCloserMoq.onCall().Close().returnResults(errors.New("ignored-error"))

			// ACT
			err := internal.Finalize(
				"some-path-to-state",
				"/my-root-dir",
				openFnMoq.mock(),
				generateFnMoq.mock())

			// ASSERT
			if err == nil {
				t.Fatalf("got %#v, wanted an error", err)
			}
			expectedErrMsg := "error reading state file some-path-to-state: bad-file"
			if err.Error() != expectedErrMsg {
				t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
			}
		})

		t.Run("invalid json", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)
			defer afterEach()
			openFnMoq.onCall("some-path-to-state").
				returnResults(readCloserMoq.mock(), nil)
			readCloserMoq.onCall().Read(nil).any().p().
				doReturnResults(readFn("this is not json")).
				doReturnResults(readFn(goodInitLine)).repeat(moq.AnyTimes())
			readCloserMoq.onCall().Close().noSeq().returnResults(errors.New("ignored-error"))

			// ACT
			err := internal.Finalize(
				"some-path-to-state",
				"/my-root-dir",
				openFnMoq.mock(),
				generateFnMoq.mock())

			// ASSERT
			if err == nil {
				t.Fatalf("got %#v, wanted an error", err)
			}
			expectedErrMsg := "error unmarshalling state file some-path-to-state:" +
				" invalid character 'h' in literal true (expecting 'r')"
			if err.Error() != expectedErrMsg {
				t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
			}
		})

		t.Run("bad GOPATH", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)
			defer afterEach()
			openFnMoq.onCall("some-path-to-state").
				returnResults(readCloserMoq.mock(), nil)
			readCloserMoq.onCall().Read(nil).any().p().doReturnResults(readFn(badGoPathInitLine))
			readCloserMoq.onCall().Close().returnResults(nil)

			// ACT
			err := internal.Finalize(
				"some-path-to-state",
				"/my-root-dir",
				openFnMoq.mock(),
				generateFnMoq.mock())

			// ASSERT
			if err == nil {
				t.Fatalf("got %#v, wanted an error", err)
			}
			if !errors.Is(err, internal.ErrBulkState) {
				t.Errorf("got %#v, want %#v", err, internal.ErrBulkState)
			}
			expectedErrMsg := fmt.Sprintf("bulk state error: current GOPATH"+
				" doesn't match GOPATH from state file some-path-to-state (%s !="+
				" not-the-right-gopath)", build.Default.GOPATH)
			if err.Error() != expectedErrMsg {
				t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
			}
		})

		t.Run("bad base current dir", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)
			defer afterEach()
			openFnMoq.onCall("some-path-to-state").
				returnResults(readCloserMoq.mock(), nil)
			readCloserMoq.onCall().Read(nil).any().p().doReturnResults(readFn(goodInitLine))
			readCloserMoq.onCall().Close().returnResults(nil)

			// ACT
			err := internal.Finalize(
				"some-path-to-state",
				"/some-other-dir",
				openFnMoq.mock(),
				generateFnMoq.mock())

			// ASSERT
			if err == nil {
				t.Fatalf("got %#v, wanted an error", err)
			}
			if !errors.Is(err, internal.ErrBulkState) {
				t.Errorf("got %#v, want %#v", err, internal.ErrBulkState)
			}
			expectedErrMsg := "bulk state error: finalize root directory" +
				" /some-other-dir does not match root directory /my-root-dir from" +
				" state file some-path-to-state"
			if err.Error() != expectedErrMsg {
				t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
			}
		})
	})

	t.Run("invalid req json", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()
		openFnMoq.onCall("some-path-to-state").
			returnResults(readCloserMoq.mock(), nil)
		readCloserMoq.onCall().Read(nil).any().p().
			doReturnResults(readFn(goodInitLine)).
			doReturnResults(readFn("this is not json"))
		readCloserMoq.onCall().Close().returnResults(errors.New("ignored-error"))

		// ACT
		err := internal.Finalize(
			"some-path-to-state",
			"/my-root-dir",
			openFnMoq.mock(),
			generateFnMoq.mock())
		// ASSERT
		if err == nil {
			t.Fatalf("got %#v, wanted an error", err)
		}
		expectedErrMsg := "error unmarshalling request from state file some-path-to-state:" +
			" invalid character 'h' in literal true (expecting 'r')"
		if err.Error() != expectedErrMsg {
			t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
		}
	})

	t.Run("req read error", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()
		openFnMoq.onCall("some-path-to-state").
			returnResults(readCloserMoq.mock(), nil)
		readCloserMoq.onCall().Read(nil).any().p().
			doReturnResults(readFn(goodInitLine)).
			doReturnResults(func(p []byte) (int, error) {
				// Missing proper line ending causes a read error
				copy(p, goodInitLine)
				return len(goodInitLine), nil
			}).
			doReturnResults(func(p []byte) (int, error) {
				return 0, errors.New("bad-file")
			})
		readCloserMoq.onCall().Close().returnResults(errors.New("ignored-error"))

		// ACT
		err := internal.Finalize(
			"some-path-to-state",
			"/my-root-dir",
			openFnMoq.mock(),
			generateFnMoq.mock())
		// ASSERT
		if err == nil {
			t.Fatalf("got %#v, wanted an error", err)
		}
		expectedErrMsg := "error reading request from state file some-path-to-state: bad-file"
		if err.Error() != expectedErrMsg {
			t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
		}
	})

	t.Run("generate error", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()
		openFnMoq.onCall("some-path-to-state").
			returnResults(readCloserMoq.mock(), nil)
		readCloserMoq.onCall().Read(nil).any().p().
			doReturnResults(readFn(goodInitLine)).
			doReturnResults(readFn(appendLine1)).
			doReturnResults(readFn(appendLine2)).
			doReturnResults(func(p []byte) (int, error) {
				return 0, io.EOF
			})
		generateFnMoq.onCall(req1, req2).returnResults(errors.New("generate-error"))
		readCloserMoq.onCall().Close().returnResults(nil)

		// ACT
		err := internal.Finalize(
			"some-path-to-state",
			"/my-root-dir",
			openFnMoq.mock(),
			generateFnMoq.mock())
		// ASSERT
		if err == nil {
			t.Fatalf("got %#v, wanted an error", err)
		}
		expectedErrMsg := "error generating mocks from state file some-path-to-state: generate-error"
		if err.Error() != expectedErrMsg {
			t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
		}
	})
}
