package internal

import (
	"encoding/json"
	"fmt"
	"io"

	"moqueries.org/runtime/logs"

	"moqueries.org/cli/generator"
)

//go:generate moqueries OpenFn

// OpenFn is the function type of os.Open
type OpenFn func(name string) (file io.ReadCloser, err error)

//go:generate moqueries --import io ReadCloser

//go:generate moqueries GenerateFn

// GenerateFn is the function type of generator.Generate
type GenerateFn func(reqs ...generator.GenerateRequest) error

// Finalize complete bulk processing by generating all the requested mocks
func Finalize(stateFile, rootDir string, openFn OpenFn, generateFn GenerateFn) error {
	f, err := openFn(stateFile)
	if err != nil {
		return fmt.Errorf("error opening state file: %w", err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			logs.Error("error closing state file", err)
		}
	}()
	scanner, err := verifyState(f, stateFile, rootDir, true)
	if err != nil {
		return err
	}

	var reqs []generator.GenerateRequest
	for scanner.Scan() {
		txt := scanner.Text()

		var req generator.GenerateRequest
		err = json.Unmarshal([]byte(txt), &req)
		if err != nil {
			return fmt.Errorf("error unmarshalling request from state file %s: %w", stateFile, err)
		}
		reqs = append(reqs, req)
	}

	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("error reading request from state file %s: %w", stateFile, err)
	}

	err = generateFn(reqs...)
	if err != nil {
		return fmt.Errorf("error generating mocks from state file %s: %w", stateFile, err)
	}

	return nil
}
