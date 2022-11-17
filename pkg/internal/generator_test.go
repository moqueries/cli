package internal_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dave/dst"

	"github.com/myshkin5/moqueries/generator"
	"github.com/myshkin5/moqueries/metrics"
	"github.com/myshkin5/moqueries/moq"
	"github.com/myshkin5/moqueries/pkg/internal"
)

func TestGenerate(t *testing.T) {
	var (
		scene      *moq.Scene
		cacheMoq   *moqTypeCache
		metricsMoq *metrics.MoqMetrics
		genFnMoq   *moqGenerateWithTypeCacheFn
	)

	beforeEach := func(t *testing.T) {
		t.Helper()

		if scene != nil {
			t.Fatal("afterEach not called")
		}
		scene = moq.NewScene(t)
		cacheMoq = newMoqTypeCache(scene, nil)
		metricsMoq = metrics.NewMoqMetrics(scene, nil)
		genFnMoq = newMoqGenerateWithTypeCacheFn(scene, nil)
		genFnMoq.runtime.parameterIndexing.cache = moq.ParamIndexByValue
	}

	afterEach := func(t *testing.T) {
		t.Helper()
		scene.AssertExpectationsMet()
		scene = nil
	}

	t.Run("happy path", func(t *testing.T) {
		testCases := map[string]struct {
			skipPkgDirs int
			destDir1    string
			destDir2    string
		}{
			"no skips": {
				skipPkgDirs: 0,
				destDir1:    "that-dir/there/pkg1",
				destDir2:    "that-dir/there/pkg2",
			},
			"skips": {
				skipPkgDirs: 3,
				destDir1:    ".",
				destDir2:    ".",
			},
		}
		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				beforeEach(t)
				defer afterEach(t)

				cacheMoq.onCall().LoadPackage("pkg1").returnResults(nil)
				cacheMoq.onCall().LoadPackage("pkg2").returnResults(nil)
				id1 := dst.Ident{Name: "Typ1", Path: "pkg1"}
				id2 := dst.Ident{Name: "Typ2", Path: "pkg2"}
				cacheMoq.onCall().MockableTypes(true).returnResults([]dst.Ident{id1, id2})
				req1 := generator.GenerateRequest{
					Types:              []string{"Typ1"},
					Export:             true,
					DestinationDir:     tc.destDir1,
					Import:             "pkg1",
					ExcludeNonExported: true,
				}
				genFnMoq.onCall(cacheMoq.mock(), req1).returnResults(nil)
				req2 := generator.GenerateRequest{
					Types:              []string{"Typ2"},
					Export:             true,
					DestinationDir:     tc.destDir2,
					Import:             "pkg2",
					ExcludeNonExported: true,
				}
				genFnMoq.onCall(cacheMoq.mock(), req2).returnResults(nil)
				metricsMoq.OnCall().TotalProcessingTimeInc(0).Any().D().ReturnResults()
				metricsMoq.OnCall().Finalize().ReturnResults()

				// ACT
				err := internal.Generate(
					cacheMoq.mock(),
					metricsMoq.Mock(),
					genFnMoq.mock(),
					"./that-dir/there",
					tc.skipPkgDirs,
					[]string{"pkg1", "pkg2"})
				// ASSERT
				if err != nil {
					t.Fatalf("got %#v, want no error", err)
				}
			})
		}
	})

	t.Run("skip too many package dirs", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach(t)

		cacheMoq.onCall().LoadPackage("pkg1").returnResults(nil)
		cacheMoq.onCall().LoadPackage("pkg2").returnResults(nil)
		id1 := dst.Ident{Name: "Typ1", Path: "pkg1"}
		id2 := dst.Ident{Name: "Typ2", Path: "pkg2"}
		cacheMoq.onCall().MockableTypes(true).returnResults([]dst.Ident{id1, id2})

		// ACT
		err := internal.Generate(
			cacheMoq.mock(),
			metricsMoq.Mock(),
			genFnMoq.mock(),
			"./that-dir/there",
			4,
			[]string{"pkg1", "pkg2"})
		// ASSERT
		if err == nil {
			t.Fatalf("got no error, want error")
		}

		if !errors.Is(err, internal.ErrSkipTooManyPackageDirs) {
			t.Errorf("got %#v, want an internal.ErrSkipTooManyPackageDirs", err)
		}

		expectedMsg := "skipping too many package dirs: skipping 4 directories on that-dir/there/pkg1 path"
		if err.Error() != expectedMsg {
			t.Errorf("got %s, want %s", err.Error(), expectedMsg)
		}
	})

	t.Run("load error", func(t *testing.T) {
		// ASSEMBLE
		beforeEach(t)
		defer afterEach(t)

		cacheMoq.onCall().LoadPackage("pkg1").returnResults(nil)
		expectedErr := errors.New("load-error")
		cacheMoq.onCall().LoadPackage("pkg2").returnResults(expectedErr)

		// ACT
		err := internal.Generate(
			cacheMoq.mock(),
			metricsMoq.Mock(),
			genFnMoq.mock(),
			"./that-dir/there",
			0,
			[]string{"pkg1", "pkg2"})
		// ASSERT
		if err != expectedErr {
			t.Fatalf("got no error, want %#v", expectedErr)
		}
	})

	t.Run("generate error", func(t *testing.T) {
		tc := map[string]error{
			"no first err":           nil,
			"non-exported first err": fmt.Errorf("%w: wha-wha", generator.ErrNonExported),
		}

		for name, firstErr := range tc {
			t.Run(name, func(t *testing.T) {
				// ASSEMBLE
				beforeEach(t)
				defer afterEach(t)

				cacheMoq.onCall().LoadPackage("pkg1").returnResults(nil)
				cacheMoq.onCall().LoadPackage("pkg2").returnResults(nil)
				id1 := dst.Ident{Name: "Typ1", Path: "pkg1"}
				id2 := dst.Ident{Name: "Typ2", Path: "pkg2"}
				cacheMoq.onCall().MockableTypes(true).returnResults([]dst.Ident{id1, id2})
				req1 := generator.GenerateRequest{
					Types:              []string{"Typ1"},
					Export:             true,
					DestinationDir:     "that-dir/there/pkg1",
					Import:             "pkg1",
					ExcludeNonExported: true,
				}
				genFnMoq.onCall(cacheMoq.mock(), req1).returnResults(firstErr)
				req2 := generator.GenerateRequest{
					Types:              []string{"Typ2"},
					Export:             true,
					DestinationDir:     "that-dir/there/pkg2",
					Import:             "pkg2",
					ExcludeNonExported: true,
				}
				expectedErr := errors.New("generate-error")
				genFnMoq.onCall(cacheMoq.mock(), req2).returnResults(expectedErr)

				// ACT
				err := internal.Generate(
					cacheMoq.mock(),
					metricsMoq.Mock(),
					genFnMoq.mock(),
					"./that-dir/there",
					0,
					[]string{"pkg1", "pkg2"})
				// ASSERT
				if err != expectedErr {
					t.Fatalf("got %#v, want %#v", err, expectedErr)
				}
			})
		}
	})
}
