package generator_test

import (
	"os"
	"testing"

	"github.com/dave/dst"
	"golang.org/x/tools/go/packages"
	"moqueries.org/runtime/logs"

	"moqueries.org/cli/ast"
	"moqueries.org/cli/generator"
	"moqueries.org/cli/metrics"
)

func TestDiamondEmbeddedInterfaces(t *testing.T) {
	// ASSEMBLE
	m := metrics.NewMetrics(logs.IsDebug, logs.Debugf)
	cache := ast.NewCache(packages.Load, os.Stat, os.ReadFile, m)
	gen := generator.New(cache, os.Getwd, func(typ generator.Type, export bool) generator.Converterer {
		return generator.NewConverter(typ, export, cache)
	})

	// ACT
	resp, err := gen.Generate(generator.GenerateRequest{
		Types:       []string{"Diamond"},
		Destination: "moq_diamond_test.go",
		Import:      "moqueries.org/cli/generator/testpkgs/diamond",
	})
	// ASSERT
	if err != nil {
		t.Fatalf("got %#v, wanted no err", err)
	}

	counts := map[string]int{}
	for _, decl := range resp.File.Decls {
		fn, ok := decl.(*dst.FuncDecl)
		if !ok || fn.Recv == nil {
			continue
		}
		if receiverName(t, fn.Recv) == "moqDiamond_mock" {
			counts[fn.Name.Name]++
		}
	}

	if got, want := len(counts), 3; got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	for _, method := range []string{
		"Base",
		"Left",
		"Right",
	} {
		if counts[method] != 1 {
			t.Errorf("got %d declarations of %s, wanted 1", counts[method], method)
		}
	}
}

func receiverName(t *testing.T, recv *dst.FieldList) string {
	t.Helper()

	expr := recv.List[0].Type
	if star, ok := expr.(*dst.StarExpr); ok {
		expr = star.X
	}
	id, ok := expr.(*dst.Ident)
	if !ok {
		t.Fatalf("got %#v, wanted *dst.Ident receiver", expr)
	}

	return id.Name
}
