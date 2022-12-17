package internal

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"

	"moqueries.org/cli/generator"
	"moqueries.org/cli/logs"
)

var (
	// ErrBadAppendRequest is returned when a caller passes bad parameters to
	// Append
	ErrBadAppendRequest = errors.New("bad request")
	// ErrBulkState is returned when the bulk state is invalid
	ErrBulkState = errors.New("bulk state error")
)

//go:generate moqueries OpenFileFn

// OpenFileFn is the function type of os.OpenFile
type OpenFileFn func(name string, flag int, perm os.FileMode) (ReadWriteSeekCloser, error)

//go:generate moqueries ReadWriteSeekCloser

// ReadWriteSeekCloser is the interface that groups the basic Read, Write,
// Seek and Close methods.
type ReadWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}

// Append appends a mock generate request to the bulk state
func Append(stateFile string, req generator.GenerateRequest, openFileFn OpenFileFn) error {
	if !filepath.IsAbs(req.WorkingDir) {
		return fmt.Errorf("%w: the request working directory must be absolute: %s",
			ErrBadAppendRequest, req.WorkingDir)
	}

	f, err := openFileFn(stateFile, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		return fmt.Errorf("error opening state file: %w", err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			logs.Error("error closing state file", err)
		}
	}()

	_, err = verifyState(f, stateFile, req.WorkingDir, false)
	if err != nil {
		return err
	}

	err = appendRequest(f, stateFile, req)
	if err != nil {
		return err
	}

	return nil
}

func verifyState(f io.ReadCloser, stateFile, workingDir string, rootDirOnly bool) (*bufio.Scanner, error) {
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return nil, fmt.Errorf("%w: state file %s not initialized properly",
			ErrBulkState, stateFile)
	}

	txt := scanner.Text()
	err := scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("error reading state file %s: %w", stateFile, err)
	}

	var state initialState
	err = json.Unmarshal([]byte(txt), &state)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling state file %s: %w", stateFile, err)
	}

	if state.GoPath != build.Default.GOPATH {
		return nil, fmt.Errorf("%w: current GOPATH doesn't match GOPATH from state file %s (%s != %s)",
			ErrBulkState, stateFile, build.Default.GOPATH, state.GoPath)
	}

	if rootDirOnly {
		if state.RootDir != workingDir {
			return nil, fmt.Errorf("%w: finalize root directory %s does"+
				" not match root directory %s from state file %s",
				ErrBulkState, workingDir, state.RootDir, stateFile)
		}
	} else {
		rel, err := filepath.Rel(state.RootDir, workingDir)
		if err != nil {
			logs.Panicf("error getting relative path %s from %s: %#v",
				state.RootDir, workingDir, err)
		}

		if strings.HasPrefix(rel, "..") {
			return nil, fmt.Errorf("%w: working directory %s is not a"+
				" child of root directory %s from state file %s",
				ErrBulkState, workingDir, state.RootDir, stateFile)
		}
	}

	return scanner, nil
}

func appendRequest(f ReadWriteSeekCloser, stateFile string, req generator.GenerateRequest) error {
	_, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("error seeking end of state file %s: %w", stateFile, err)
	}
	_, err = f.Write(compact(req))
	if err != nil {
		return fmt.Errorf("error writing state file %s: %w", stateFile, err)
	}
	_, err = f.Write([]byte("\n"))
	if err != nil {
		return fmt.Errorf("error finishing writing of state file %s: %w", stateFile, err)
	}
	return nil
}
