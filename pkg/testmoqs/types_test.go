package testmoqs_test

import (
	"os"
	"strings"
	"testing"

	"github.com/myshkin5/moqueries/pkg"
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
	err = pkg.Generate(".", 5, "github.com/myshkin5/moqueries/pkg/testmoqs")

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
	if len(moqs) != 4 {
		t.Errorf("got %#v mocks, want 3", moqs)
	}
	if _, ok := moqs["moq_passbyrefsimple_stargentype.go"]; !ok {
		t.Errorf("got %#v, want to contain moq_passbyrefsimple_stargentype.go", moqs)
	}
	if _, ok := moqs["moq_passbyvaluesimple_gentype.go"]; !ok {
		t.Errorf("got %#v, want to contain moq_passbyvaluesimple_gentype.go", moqs)
	}
	if _, ok := moqs["moq_standalonefunc_gentype.go"]; !ok {
		t.Errorf("got %#v, want to contain moq_standalonefunc_gentype.go", moqs)
	}
	if _, ok := moqs["moq_reduced.go"]; !ok {
		t.Errorf("got %#v, want to contain moq_standalonefunc_gentype.go", moqs)
	}
	// Minimal testing here just to make sure the right types were found (full
	// mock testing in the generator package)
}
