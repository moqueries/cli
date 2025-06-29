// Package bulk is used to generate several Moqueries mocks at once
package bulk

import (
	"io"
	"os"

	"moqueries.org/cli/bulk/internal"
	"moqueries.org/cli/generator"
)

// Initialize initializes bulk processing and creates the bulk processing state
// file
func Initialize(stateFile, rootDir string) error {
	createFn := func(name string) (io.WriteCloser, error) {
		return os.Create(name)
	}

	return internal.Initialize(stateFile, rootDir, createFn)
}

// Append appends a mock generate request to the bulk state
func Append(stateFile string, request generator.GenerateRequest) error {
	openFileFn := func(name string, flag int, perm os.FileMode) (internal.ReadWriteSeekCloser, error) {
		return os.OpenFile(name, flag, perm)
	}

	return internal.Append(stateFile, request, openFileFn)
}

// Finalize complete bulk processing by generating all the requested mocks
func Finalize(stateFile, rootDir string) error {
	openFn := func(name string) (io.ReadCloser, error) {
		return os.Open(name)
	}

	return internal.Finalize(stateFile, rootDir, openFn, generator.Generate)
}
