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

	// ACT
	err = pkg.Generate(pkg.PackageGenerateRequest{
		DestinationDir: ".",
		SkipPkgDirs:    4,
		PkgPatterns:    []string{"moqueries.org/cli/pkg/testmoqs"},
	})
	// ASSERT
	if err != nil {
		t.Fatalf("got %#v, want no error", err)
	}

	entries, err = os.ReadDir(".")
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
	if count := 4; len(moqs) != count {
		t.Errorf("got %#v mocks, want %d", moqs, count)
	}
	if _, ok := moqs["moq_passbyrefsimple.go"]; !ok {
		t.Errorf("got %#v, want to contain moq_passbyrefsimple_stargentype.go", moqs)
	}
	if _, ok := moqs["moq_passbyvaluesimple.go"]; !ok {
		t.Errorf("got %#v, want to contain moq_passbyvaluesimple_gentype.go", moqs)
	}
	if _, ok := moqs["moq_standalonefunc.go"]; !ok {
		t.Errorf("got %#v, want to contain moq_standalonefunc_gentype.go", moqs)
	}
	if _, ok := moqs["moq_reduced.go"]; !ok {
		t.Errorf("got %#v, want to contain moq_standalonefunc_gentype.go", moqs)
	}
	// Minimal testing here just to make sure the right types were found (full
	// mock testing in the generator package)
}
