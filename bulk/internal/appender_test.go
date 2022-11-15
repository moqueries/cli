package internal_test

import (
	"errors"
	"fmt"
	"go/build"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/myshkin5/moqueries/bulk/internal"
	"github.com/myshkin5/moqueries/generator"
	"github.com/myshkin5/moqueries/moq"
)

var badGoPathInitLine = "{\"root-dir\":\"/my-root-dir\",\"go-path\":\"not-the-right-gopath\"}"

func readFn(line string) func(p []byte) (int, error) {
	return func(p []byte) (int, error) {
		line += "\n"
		copy(p, line)
		return len(line), nil
	}
}

func TestAppend(t *testing.T) {
	var (
		appendLine1 = "{" +
			"\"types\":null," +
			"\"export\":false," +
			"\"destination\":\"\"," +
			"\"destination-dir\":\"\"," +
			"\"package\":\"\"," +
			"\"import\":\"\"," +
			"\"test-import\":false," +
			"\"working-dir\":\"/my-root-dir\"," +
			"\"exclude-non-exported\":false" +
			"}"

		scene         *moq.Scene
		openFileFnMoq *moqOpenFileFn
		fileMoq       *moqReadWriteSeekCloser
	)

	beforeEach := func(t *testing.T) {
		t.Helper()
		if scene != nil {
			t.Fatal("afterEach not called")
		}
		scene = moq.NewScene(t)
		config := &moq.Config{Sequence: moq.SeqDefaultOn}
		openFileFnMoq = newMoqOpenFileFn(scene, config)
		fileMoq = newMoqReadWriteSeekCloser(scene, config)
	}

	afterEach := func() {
		scene.AssertExpectationsMet()
		scene = nil
	}

	t.Run("simple", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()
		openFileFnMoq.onCall("some-path-to-state", os.O_RDWR|os.O_APPEND, 0).
			returnResults(fileMoq.mock(), nil)
		fileMoq.onCall().Read(nil).any().p().doReturnResults(readFn(goodInitLine))
		fileMoq.onCall().Seek(0, io.SeekEnd).returnResults(0, nil)
		fileMoq.onCall().Write([]byte(appendLine1)).returnResults(0, nil)
		fileMoq.onCall().Write([]byte("\n")).returnResults(0, nil)
		fileMoq.onCall().Close().returnResults(errors.New("ignored-error"))

		req := generator.GenerateRequest{WorkingDir: "/my-root-dir"}

		// ACT
		err := internal.Append("some-path-to-state", req, openFileFnMoq.mock())
		// ASSERT
		if err != nil {
			t.Errorf("got %#v, wanted no error", err)
		}
	})

	t.Run("simple subdir", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()
		openFileFnMoq.onCall("some-path-to-state", os.O_RDWR|os.O_APPEND, 0).
			returnResults(fileMoq.mock(), nil)
		fileMoq.onCall().Read(nil).any().p().doReturnResults(readFn(goodInitLine))
		fileMoq.onCall().Seek(0, io.SeekEnd).returnResults(0, nil)
		l := strings.ReplaceAll(appendLine1, "/my-root-dir", "/my-root-dir/subdir")
		fileMoq.onCall().Write([]byte(l)).returnResults(0, nil)
		fileMoq.onCall().Write([]byte("\n")).returnResults(0, nil)
		fileMoq.onCall().Close().returnResults(errors.New("ignored-error"))

		req := generator.GenerateRequest{WorkingDir: "/my-root-dir/subdir"}

		// ACT
		err := internal.Append("some-path-to-state", req, openFileFnMoq.mock())
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
			openFileFnMoq.onCall("some-path-to-state", os.O_RDWR|os.O_APPEND, 0).
				returnResults(nil, errors.New("open-error"))

			req := generator.GenerateRequest{WorkingDir: "/my-root-dir/subdir"}

			// ACT
			err := internal.Append("some-path-to-state", req, openFileFnMoq.mock())

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
			openFileFnMoq.onCall("some-path-to-state", os.O_RDWR|os.O_APPEND, 0).
				returnResults(fileMoq.mock(), nil)
			fileMoq.onCall().Read(nil).any().p().
				doReturnResults(func(p []byte) (int, error) {
					return 0, io.EOF
				}).repeat(moq.AnyTimes())
			fileMoq.onCall().Close().returnResults(errors.New("ignored-error"))

			req := generator.GenerateRequest{WorkingDir: "/my-root-dir/subdir"}

			// ACT
			err := internal.Append("some-path-to-state", req, openFileFnMoq.mock())

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
			openFileFnMoq.onCall("some-path-to-state", os.O_RDWR|os.O_APPEND, 0).
				returnResults(fileMoq.mock(), nil)
			fileMoq.onCall().Read(nil).any().p().
				doReturnResults(func(p []byte) (int, error) {
					// Missing proper line ending causes a read error
					copy(p, goodInitLine)
					return len(goodInitLine), nil
				}).
				doReturnResults(func(p []byte) (int, error) {
					return 0, errors.New("bad-file")
				})
			fileMoq.onCall().Close().returnResults(errors.New("ignored-error"))

			req := generator.GenerateRequest{WorkingDir: "/my-root-dir/subdir"}

			// ACT
			err := internal.Append("some-path-to-state", req, openFileFnMoq.mock())

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
			openFileFnMoq.onCall("some-path-to-state", os.O_RDWR|os.O_APPEND, 0).
				returnResults(fileMoq.mock(), nil)
			fileMoq.onCall().Read(nil).any().p().
				doReturnResults(readFn("this is not json")).
				doReturnResults(readFn(goodInitLine)).repeat(moq.AnyTimes())
			fileMoq.onCall().Close().noSeq().returnResults(errors.New("ignored-error"))

			req := generator.GenerateRequest{WorkingDir: "/my-root-dir/subdir"}

			// ACT
			err := internal.Append("some-path-to-state", req, openFileFnMoq.mock())

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
			openFileFnMoq.onCall("some-path-to-state", os.O_RDWR|os.O_APPEND, 0).
				returnResults(fileMoq.mock(), nil)
			fileMoq.onCall().Read(nil).any().p().doReturnResults(readFn(badGoPathInitLine))
			fileMoq.onCall().Close().returnResults(nil)

			req := generator.GenerateRequest{WorkingDir: "/my-root-dir/subdir"}

			// ACT
			err := internal.Append("some-path-to-state", req, openFileFnMoq.mock())

			// ASSERT
			if err == nil {
				t.Fatalf("got %#v, wanted an error", err)
			}
			if !errors.Is(err, internal.ErrBulkState) {
				t.Errorf("got %#v, want %#v", err, internal.ErrBulkState)
			}
			expectedErrMsg := fmt.Sprintf("bulk state error: current GOPATH"+
				" doesn't match GOPATH from state file some-path-to-state"+
				" (%s != not-the-right-gopath)", build.Default.GOPATH)
			if err.Error() != expectedErrMsg {
				t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
			}
		})

		t.Run("bad base current dir", func(t *testing.T) {
			// ASSEMBLE
			beforeEach(t)
			defer afterEach()
			openFileFnMoq.onCall("some-path-to-state", os.O_RDWR|os.O_APPEND, 0).
				returnResults(fileMoq.mock(), nil)
			fileMoq.onCall().Read(nil).any().p().doReturnResults(readFn(goodInitLine))
			fileMoq.onCall().Close().returnResults(nil)

			req := generator.GenerateRequest{WorkingDir: "/some-other-dir"}

			// ACT
			err := internal.Append("some-path-to-state", req, openFileFnMoq.mock())

			// ASSERT
			if err == nil {
				t.Fatalf("got %#v, wanted an error", err)
			}
			if !errors.Is(err, internal.ErrBulkState) {
				t.Errorf("got %#v, want %#v", err, internal.ErrBulkState)
			}
			expectedErrMsg := "bulk state error: working directory /some-other-dir is" +
				" not a child of root directory /my-root-dir from state file" +
				" some-path-to-state"
			if err.Error() != expectedErrMsg {
				t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
			}
		})
	})

	t.Run("bad relative current dir", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()

		req := generator.GenerateRequest{WorkingDir: "./subdir"}

		// ACT
		err := internal.Append("some-path-to-state", req, openFileFnMoq.mock())

		// ASSERT
		if err == nil {
			t.Fatalf("got %#v, wanted an error", err)
		}
		if !errors.Is(err, internal.ErrBadAppendRequest) {
			t.Errorf("got %#v, want %#v", err, internal.ErrBadAppendRequest)
		}
		expectedErrMsg := "bad request: the request working directory must be absolute: ./subdir"
		if err.Error() != expectedErrMsg {
			t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
		}
	})

	t.Run("seek error", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()
		openFileFnMoq.onCall("some-path-to-state", os.O_RDWR|os.O_APPEND, 0).
			returnResults(fileMoq.mock(), nil)
		fileMoq.onCall().Read(nil).any().p().doReturnResults(readFn(goodInitLine))
		fileMoq.onCall().Seek(0, io.SeekEnd).returnResults(0, errors.New("seek-error"))
		fileMoq.onCall().Close().returnResults(nil)

		req := generator.GenerateRequest{WorkingDir: "/my-root-dir"}

		// ACT
		err := internal.Append("some-path-to-state", req, openFileFnMoq.mock())

		// ASSERT
		if err == nil {
			t.Fatalf("got %#v, wanted an error", err)
		}
		expectedErrMsg := "error seeking end of state file some-path-to-state: seek-error"
		if err.Error() != expectedErrMsg {
			t.Errorf("got %s, wanted %s", err.Error(), expectedErrMsg)
		}
	})

	t.Run("write error", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach()
		openFileFnMoq.onCall("some-path-to-state", os.O_RDWR|os.O_APPEND, 0).
			returnResults(fileMoq.mock(), nil)
		fileMoq.onCall().Read(nil).any().p().doReturnResults(readFn(goodInitLine))
		fileMoq.onCall().Seek(0, io.SeekEnd).returnResults(0, nil)
		fileMoq.onCall().Write([]byte(appendLine1)).returnResults(0, errors.New("write-error"))
		fileMoq.onCall().Close().returnResults(nil)

		req := generator.GenerateRequest{WorkingDir: "/my-root-dir"}

		// ACT
		err := internal.Append("some-path-to-state", req, openFileFnMoq.mock())

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
		openFileFnMoq.onCall("some-path-to-state", os.O_RDWR|os.O_APPEND, 0).
			returnResults(fileMoq.mock(), nil)
		fileMoq.onCall().Read(nil).any().p().doReturnResults(readFn(goodInitLine))
		fileMoq.onCall().Seek(0, io.SeekEnd).returnResults(0, nil)
		fileMoq.onCall().Write([]byte(appendLine1)).returnResults(0, nil)
		fileMoq.onCall().Write([]byte("\n")).returnResults(0, errors.New("write-error"))
		fileMoq.onCall().Close().returnResults(nil)

		req := generator.GenerateRequest{WorkingDir: "/my-root-dir"}

		// ACT
		err := internal.Append("some-path-to-state", req, openFileFnMoq.mock())

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
