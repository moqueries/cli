// Package bulk is used to generate several Moqueries mocks at once
package bulk

import (
	"io"
	"os"

	"github.com/myshkin5/moqueries/bulk/internal"
	"github.com/myshkin5/moqueries/generator"
)

// Initialize initializes bulk processing and creates the bulk processing state
// file
func Initialize(stateFile, rootDir string) error {
	createFn := func(name string) (io.WriteCloser, error) {
		//nolint:gosec // Users can use any file for bulk operations
		return os.Create(name)
	}

	return internal.Initialize(stateFile, rootDir, createFn)
}

// Append appends a mock generate request to the bulk state
func Append(stateFile string, request generator.GenerateRequest) error {
	openFileFn := func(name string, flag int, perm os.FileMode) (internal.ReadWriteSeekCloser, error) {
		//nolint:gosec // Users can use any file for bulk operations
		return os.OpenFile(name, flag, perm)
	}

	return internal.Append(stateFile, request, openFileFn)
}

// Finalize complete bulk processing by generating all the requested mocks
func Finalize(stateFile, rootDir string) error {
	openFn := func(name string) (io.ReadCloser, error) {
		//nolint:gosec // Users can use any file for bulk operations
		return os.Open(name)
	}

	return internal.Finalize(stateFile, rootDir, openFn, generator.Generate)
}
