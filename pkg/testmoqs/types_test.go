package testmoqs_test

import (
	"os"
	"strings"
	"testing"

	"moqueries.org/cli/pkg"
)

func TestPackageGeneration(t *testing.T) {
	// ASSEMBLE
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("got %#v, want no error", err)
	}
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "moq_") || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}

		err = os.Remove(e.Name())
		if err != nil {
			t.Errorf("got %#v, want no error", err)
		}
	}
	expectedFiles := []string{
		"moq_dosomethingwithparam.go",
		"moq_passbyrefsimple.go",
		"moq_passbyvaluesimple.go",
		"moq_reduced.go",
		"moq_standalonefunc.go",
	}

	// ACT
	err = pkg.Generate(pkg.PackageGenerateRequest{
		DestinationDir: "pkgout",
		SkipPkgDirs:    4,
		PkgPatterns:    []string{"moqueries.org/cli/pkg/testmoqs"},
	})
	// ASSERT
	if err != nil {
		t.Fatalf("got %#v, want no error", err)
	}

	entries, err = os.ReadDir("pkgout")
	if err != nil {
		t.Fatalf("got %#v, want no error", err)
	}
	moqs := map[string]struct{}{}
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), "moq_") || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}

		moqs[e.Name()] = struct{}{}
	}
	if len(moqs) != len(expectedFiles) {
		t.Errorf("got %#v mocks, want length %d", moqs, len(expectedFiles))
	}
	for _, f := range expectedFiles {
		if _, ok := moqs[f]; !ok {
			t.Errorf("got %#v, want to contain %s", moqs, f)
		}
	}
	// Minimal testing here just to make sure the right types were found (full
	// mock testing in the generator package)
}
