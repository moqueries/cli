// Package internal implements the internals for bulk operations
package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"io"

	"moqueries.org/runtime/logs"
)

//go:generate moqueries CreateFn

// CreateFn is the function type of os.Create
type CreateFn func(name string) (file io.WriteCloser, err error)

//go:generate moqueries --import io WriteCloser

type initialState struct {
	RootDir string `json:"root-dir"`
	GoPath  string `json:"go-path"`
}

// Initialize initializes bulk processing and creates the bulk processing state
// file
func Initialize(stateFile, rootDir string, createFn CreateFn) error {
	f, err := createFn(stateFile)
	if err != nil {
		return fmt.Errorf("error creating state file: %w", err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			logs.Error("error closing state file", err)
		}
	}()

	state := initialState{
		RootDir: rootDir,
		GoPath:  build.Default.GOPATH,
	}
	_, err = f.Write(compact(state))
	if err != nil {
		return fmt.Errorf("error writing state file %s: %w", stateFile, err)
	}
	_, err = f.Write([]byte("\n"))
	if err != nil {
		return fmt.Errorf("error finishing writing of state file %s: %w", stateFile, err)
	}

	return nil
}

func compact(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		logs.Panic("error marshalling state info", err)
	}
	buf := bytes.NewBuffer([]byte{})
	err = json.Compact(buf, b)
	if err != nil {
		logs.Panic("error compacting state info", err)
	}
	return buf.Bytes()
}
