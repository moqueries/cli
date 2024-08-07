// Code generated by Moqueries - https://moqueries.org - DO NOT EDIT.

package metrics

import (
	"fmt"
	"time"

	"moqueries.org/runtime/hash"
	"moqueries.org/runtime/impl"
	"moqueries.org/runtime/moq"
)

// The following type assertion assures that metrics.Metrics is mocked
// completely
var _ Metrics = (*MoqMetrics_mock)(nil)

// MoqMetrics holds the state of a moq of the Metrics type
type MoqMetrics struct {
	Moq *MoqMetrics_mock

	Moq_ASTPkgCacheHitsInc *impl.Moq[
		*MoqMetrics_ASTPkgCacheHitsInc_adaptor,
		MoqMetrics_ASTPkgCacheHitsInc_params,
		MoqMetrics_ASTPkgCacheHitsInc_paramsKey,
		MoqMetrics_ASTPkgCacheHitsInc_results,
	]
	Moq_ASTPkgCacheMissesInc *impl.Moq[
		*MoqMetrics_ASTPkgCacheMissesInc_adaptor,
		MoqMetrics_ASTPkgCacheMissesInc_params,
		MoqMetrics_ASTPkgCacheMissesInc_paramsKey,
		MoqMetrics_ASTPkgCacheMissesInc_results,
	]
	Moq_ASTTotalLoadTimeInc *impl.Moq[
		*MoqMetrics_ASTTotalLoadTimeInc_adaptor,
		MoqMetrics_ASTTotalLoadTimeInc_params,
		MoqMetrics_ASTTotalLoadTimeInc_paramsKey,
		MoqMetrics_ASTTotalLoadTimeInc_results,
	]
	Moq_ASTTotalDecorationTimeInc *impl.Moq[
		*MoqMetrics_ASTTotalDecorationTimeInc_adaptor,
		MoqMetrics_ASTTotalDecorationTimeInc_params,
		MoqMetrics_ASTTotalDecorationTimeInc_paramsKey,
		MoqMetrics_ASTTotalDecorationTimeInc_results,
	]
	Moq_TotalProcessingTimeInc *impl.Moq[
		*MoqMetrics_TotalProcessingTimeInc_adaptor,
		MoqMetrics_TotalProcessingTimeInc_params,
		MoqMetrics_TotalProcessingTimeInc_paramsKey,
		MoqMetrics_TotalProcessingTimeInc_results,
	]
	Moq_Finalize *impl.Moq[
		*MoqMetrics_Finalize_adaptor,
		MoqMetrics_Finalize_params,
		MoqMetrics_Finalize_paramsKey,
		MoqMetrics_Finalize_results,
	]

	Runtime MoqMetrics_runtime
}

// MoqMetrics_mock isolates the mock interface of the Metrics type
type MoqMetrics_mock struct {
	Moq *MoqMetrics
}

// MoqMetrics_recorder isolates the recorder interface of the Metrics type
type MoqMetrics_recorder struct {
	Moq *MoqMetrics
}

// MoqMetrics_runtime holds runtime configuration for the Metrics type
type MoqMetrics_runtime struct {
	ParameterIndexing struct {
		ASTPkgCacheHitsInc        MoqMetrics_ASTPkgCacheHitsInc_paramIndexing
		ASTPkgCacheMissesInc      MoqMetrics_ASTPkgCacheMissesInc_paramIndexing
		ASTTotalLoadTimeInc       MoqMetrics_ASTTotalLoadTimeInc_paramIndexing
		ASTTotalDecorationTimeInc MoqMetrics_ASTTotalDecorationTimeInc_paramIndexing
		TotalProcessingTimeInc    MoqMetrics_TotalProcessingTimeInc_paramIndexing
		Finalize                  MoqMetrics_Finalize_paramIndexing
	}
}

// MoqMetrics_ASTPkgCacheHitsInc_adaptor adapts MoqMetrics as needed by the
// runtime
type MoqMetrics_ASTPkgCacheHitsInc_adaptor struct {
	Moq *MoqMetrics
}

// MoqMetrics_ASTPkgCacheHitsInc_params holds the params of the Metrics type
type MoqMetrics_ASTPkgCacheHitsInc_params struct{}

// MoqMetrics_ASTPkgCacheHitsInc_paramsKey holds the map key params of the
// Metrics type
type MoqMetrics_ASTPkgCacheHitsInc_paramsKey struct {
	Params struct{}
	Hashes struct{}
}

// MoqMetrics_ASTPkgCacheHitsInc_results holds the results of the Metrics type
type MoqMetrics_ASTPkgCacheHitsInc_results struct{}

// MoqMetrics_ASTPkgCacheHitsInc_paramIndexing holds the parameter indexing
// runtime configuration for the Metrics type
type MoqMetrics_ASTPkgCacheHitsInc_paramIndexing struct{}

// MoqMetrics_ASTPkgCacheHitsInc_doFn defines the type of function needed when
// calling AndDo for the Metrics type
type MoqMetrics_ASTPkgCacheHitsInc_doFn func()

// MoqMetrics_ASTPkgCacheHitsInc_doReturnFn defines the type of function needed
// when calling DoReturnResults for the Metrics type
type MoqMetrics_ASTPkgCacheHitsInc_doReturnFn func()

// MoqMetrics_ASTPkgCacheHitsInc_recorder routes recorded function calls to the
// MoqMetrics moq
type MoqMetrics_ASTPkgCacheHitsInc_recorder struct {
	Recorder *impl.Recorder[
		*MoqMetrics_ASTPkgCacheHitsInc_adaptor,
		MoqMetrics_ASTPkgCacheHitsInc_params,
		MoqMetrics_ASTPkgCacheHitsInc_paramsKey,
		MoqMetrics_ASTPkgCacheHitsInc_results,
	]
}

// MoqMetrics_ASTPkgCacheHitsInc_anyParams isolates the any params functions of
// the Metrics type
type MoqMetrics_ASTPkgCacheHitsInc_anyParams struct {
	Recorder *MoqMetrics_ASTPkgCacheHitsInc_recorder
}

// MoqMetrics_ASTPkgCacheMissesInc_adaptor adapts MoqMetrics as needed by the
// runtime
type MoqMetrics_ASTPkgCacheMissesInc_adaptor struct {
	Moq *MoqMetrics
}

// MoqMetrics_ASTPkgCacheMissesInc_params holds the params of the Metrics type
type MoqMetrics_ASTPkgCacheMissesInc_params struct{}

// MoqMetrics_ASTPkgCacheMissesInc_paramsKey holds the map key params of the
// Metrics type
type MoqMetrics_ASTPkgCacheMissesInc_paramsKey struct {
	Params struct{}
	Hashes struct{}
}

// MoqMetrics_ASTPkgCacheMissesInc_results holds the results of the Metrics
// type
type MoqMetrics_ASTPkgCacheMissesInc_results struct{}

// MoqMetrics_ASTPkgCacheMissesInc_paramIndexing holds the parameter indexing
// runtime configuration for the Metrics type
type MoqMetrics_ASTPkgCacheMissesInc_paramIndexing struct{}

// MoqMetrics_ASTPkgCacheMissesInc_doFn defines the type of function needed
// when calling AndDo for the Metrics type
type MoqMetrics_ASTPkgCacheMissesInc_doFn func()

// MoqMetrics_ASTPkgCacheMissesInc_doReturnFn defines the type of function
// needed when calling DoReturnResults for the Metrics type
type MoqMetrics_ASTPkgCacheMissesInc_doReturnFn func()

// MoqMetrics_ASTPkgCacheMissesInc_recorder routes recorded function calls to
// the MoqMetrics moq
type MoqMetrics_ASTPkgCacheMissesInc_recorder struct {
	Recorder *impl.Recorder[
		*MoqMetrics_ASTPkgCacheMissesInc_adaptor,
		MoqMetrics_ASTPkgCacheMissesInc_params,
		MoqMetrics_ASTPkgCacheMissesInc_paramsKey,
		MoqMetrics_ASTPkgCacheMissesInc_results,
	]
}

// MoqMetrics_ASTPkgCacheMissesInc_anyParams isolates the any params functions
// of the Metrics type
type MoqMetrics_ASTPkgCacheMissesInc_anyParams struct {
	Recorder *MoqMetrics_ASTPkgCacheMissesInc_recorder
}

// MoqMetrics_ASTTotalLoadTimeInc_adaptor adapts MoqMetrics as needed by the
// runtime
type MoqMetrics_ASTTotalLoadTimeInc_adaptor struct {
	Moq *MoqMetrics
}

// MoqMetrics_ASTTotalLoadTimeInc_params holds the params of the Metrics type
type MoqMetrics_ASTTotalLoadTimeInc_params struct{ D time.Duration }

// MoqMetrics_ASTTotalLoadTimeInc_paramsKey holds the map key params of the
// Metrics type
type MoqMetrics_ASTTotalLoadTimeInc_paramsKey struct {
	Params struct{ D time.Duration }
	Hashes struct{ D hash.Hash }
}

// MoqMetrics_ASTTotalLoadTimeInc_results holds the results of the Metrics type
type MoqMetrics_ASTTotalLoadTimeInc_results struct{}

// MoqMetrics_ASTTotalLoadTimeInc_paramIndexing holds the parameter indexing
// runtime configuration for the Metrics type
type MoqMetrics_ASTTotalLoadTimeInc_paramIndexing struct {
	D moq.ParamIndexing
}

// MoqMetrics_ASTTotalLoadTimeInc_doFn defines the type of function needed when
// calling AndDo for the Metrics type
type MoqMetrics_ASTTotalLoadTimeInc_doFn func(d time.Duration)

// MoqMetrics_ASTTotalLoadTimeInc_doReturnFn defines the type of function
// needed when calling DoReturnResults for the Metrics type
type MoqMetrics_ASTTotalLoadTimeInc_doReturnFn func(d time.Duration)

// MoqMetrics_ASTTotalLoadTimeInc_recorder routes recorded function calls to
// the MoqMetrics moq
type MoqMetrics_ASTTotalLoadTimeInc_recorder struct {
	Recorder *impl.Recorder[
		*MoqMetrics_ASTTotalLoadTimeInc_adaptor,
		MoqMetrics_ASTTotalLoadTimeInc_params,
		MoqMetrics_ASTTotalLoadTimeInc_paramsKey,
		MoqMetrics_ASTTotalLoadTimeInc_results,
	]
}

// MoqMetrics_ASTTotalLoadTimeInc_anyParams isolates the any params functions
// of the Metrics type
type MoqMetrics_ASTTotalLoadTimeInc_anyParams struct {
	Recorder *MoqMetrics_ASTTotalLoadTimeInc_recorder
}

// MoqMetrics_ASTTotalDecorationTimeInc_adaptor adapts MoqMetrics as needed by
// the runtime
type MoqMetrics_ASTTotalDecorationTimeInc_adaptor struct {
	Moq *MoqMetrics
}

// MoqMetrics_ASTTotalDecorationTimeInc_params holds the params of the Metrics
// type
type MoqMetrics_ASTTotalDecorationTimeInc_params struct{ D time.Duration }

// MoqMetrics_ASTTotalDecorationTimeInc_paramsKey holds the map key params of
// the Metrics type
type MoqMetrics_ASTTotalDecorationTimeInc_paramsKey struct {
	Params struct{ D time.Duration }
	Hashes struct{ D hash.Hash }
}

// MoqMetrics_ASTTotalDecorationTimeInc_results holds the results of the
// Metrics type
type MoqMetrics_ASTTotalDecorationTimeInc_results struct{}

// MoqMetrics_ASTTotalDecorationTimeInc_paramIndexing holds the parameter
// indexing runtime configuration for the Metrics type
type MoqMetrics_ASTTotalDecorationTimeInc_paramIndexing struct {
	D moq.ParamIndexing
}

// MoqMetrics_ASTTotalDecorationTimeInc_doFn defines the type of function
// needed when calling AndDo for the Metrics type
type MoqMetrics_ASTTotalDecorationTimeInc_doFn func(d time.Duration)

// MoqMetrics_ASTTotalDecorationTimeInc_doReturnFn defines the type of function
// needed when calling DoReturnResults for the Metrics type
type MoqMetrics_ASTTotalDecorationTimeInc_doReturnFn func(d time.Duration)

// MoqMetrics_ASTTotalDecorationTimeInc_recorder routes recorded function calls
// to the MoqMetrics moq
type MoqMetrics_ASTTotalDecorationTimeInc_recorder struct {
	Recorder *impl.Recorder[
		*MoqMetrics_ASTTotalDecorationTimeInc_adaptor,
		MoqMetrics_ASTTotalDecorationTimeInc_params,
		MoqMetrics_ASTTotalDecorationTimeInc_paramsKey,
		MoqMetrics_ASTTotalDecorationTimeInc_results,
	]
}

// MoqMetrics_ASTTotalDecorationTimeInc_anyParams isolates the any params
// functions of the Metrics type
type MoqMetrics_ASTTotalDecorationTimeInc_anyParams struct {
	Recorder *MoqMetrics_ASTTotalDecorationTimeInc_recorder
}

// MoqMetrics_TotalProcessingTimeInc_adaptor adapts MoqMetrics as needed by the
// runtime
type MoqMetrics_TotalProcessingTimeInc_adaptor struct {
	Moq *MoqMetrics
}

// MoqMetrics_TotalProcessingTimeInc_params holds the params of the Metrics
// type
type MoqMetrics_TotalProcessingTimeInc_params struct{ D time.Duration }

// MoqMetrics_TotalProcessingTimeInc_paramsKey holds the map key params of the
// Metrics type
type MoqMetrics_TotalProcessingTimeInc_paramsKey struct {
	Params struct{ D time.Duration }
	Hashes struct{ D hash.Hash }
}

// MoqMetrics_TotalProcessingTimeInc_results holds the results of the Metrics
// type
type MoqMetrics_TotalProcessingTimeInc_results struct{}

// MoqMetrics_TotalProcessingTimeInc_paramIndexing holds the parameter indexing
// runtime configuration for the Metrics type
type MoqMetrics_TotalProcessingTimeInc_paramIndexing struct {
	D moq.ParamIndexing
}

// MoqMetrics_TotalProcessingTimeInc_doFn defines the type of function needed
// when calling AndDo for the Metrics type
type MoqMetrics_TotalProcessingTimeInc_doFn func(d time.Duration)

// MoqMetrics_TotalProcessingTimeInc_doReturnFn defines the type of function
// needed when calling DoReturnResults for the Metrics type
type MoqMetrics_TotalProcessingTimeInc_doReturnFn func(d time.Duration)

// MoqMetrics_TotalProcessingTimeInc_recorder routes recorded function calls to
// the MoqMetrics moq
type MoqMetrics_TotalProcessingTimeInc_recorder struct {
	Recorder *impl.Recorder[
		*MoqMetrics_TotalProcessingTimeInc_adaptor,
		MoqMetrics_TotalProcessingTimeInc_params,
		MoqMetrics_TotalProcessingTimeInc_paramsKey,
		MoqMetrics_TotalProcessingTimeInc_results,
	]
}

// MoqMetrics_TotalProcessingTimeInc_anyParams isolates the any params
// functions of the Metrics type
type MoqMetrics_TotalProcessingTimeInc_anyParams struct {
	Recorder *MoqMetrics_TotalProcessingTimeInc_recorder
}

// MoqMetrics_Finalize_adaptor adapts MoqMetrics as needed by the runtime
type MoqMetrics_Finalize_adaptor struct {
	Moq *MoqMetrics
}

// MoqMetrics_Finalize_params holds the params of the Metrics type
type MoqMetrics_Finalize_params struct{}

// MoqMetrics_Finalize_paramsKey holds the map key params of the Metrics type
type MoqMetrics_Finalize_paramsKey struct {
	Params struct{}
	Hashes struct{}
}

// MoqMetrics_Finalize_results holds the results of the Metrics type
type MoqMetrics_Finalize_results struct{}

// MoqMetrics_Finalize_paramIndexing holds the parameter indexing runtime
// configuration for the Metrics type
type MoqMetrics_Finalize_paramIndexing struct{}

// MoqMetrics_Finalize_doFn defines the type of function needed when calling
// AndDo for the Metrics type
type MoqMetrics_Finalize_doFn func()

// MoqMetrics_Finalize_doReturnFn defines the type of function needed when
// calling DoReturnResults for the Metrics type
type MoqMetrics_Finalize_doReturnFn func()

// MoqMetrics_Finalize_recorder routes recorded function calls to the
// MoqMetrics moq
type MoqMetrics_Finalize_recorder struct {
	Recorder *impl.Recorder[
		*MoqMetrics_Finalize_adaptor,
		MoqMetrics_Finalize_params,
		MoqMetrics_Finalize_paramsKey,
		MoqMetrics_Finalize_results,
	]
}

// MoqMetrics_Finalize_anyParams isolates the any params functions of the
// Metrics type
type MoqMetrics_Finalize_anyParams struct {
	Recorder *MoqMetrics_Finalize_recorder
}

// NewMoqMetrics creates a new moq of the Metrics type
func NewMoqMetrics(scene *moq.Scene, config *moq.Config) *MoqMetrics {
	adaptor1 := &MoqMetrics_ASTPkgCacheHitsInc_adaptor{}
	adaptor2 := &MoqMetrics_ASTPkgCacheMissesInc_adaptor{}
	adaptor3 := &MoqMetrics_ASTTotalLoadTimeInc_adaptor{}
	adaptor4 := &MoqMetrics_ASTTotalDecorationTimeInc_adaptor{}
	adaptor5 := &MoqMetrics_TotalProcessingTimeInc_adaptor{}
	adaptor6 := &MoqMetrics_Finalize_adaptor{}
	m := &MoqMetrics{
		Moq: &MoqMetrics_mock{},

		Moq_ASTPkgCacheHitsInc: impl.NewMoq[
			*MoqMetrics_ASTPkgCacheHitsInc_adaptor,
			MoqMetrics_ASTPkgCacheHitsInc_params,
			MoqMetrics_ASTPkgCacheHitsInc_paramsKey,
			MoqMetrics_ASTPkgCacheHitsInc_results,
		](scene, adaptor1, config),
		Moq_ASTPkgCacheMissesInc: impl.NewMoq[
			*MoqMetrics_ASTPkgCacheMissesInc_adaptor,
			MoqMetrics_ASTPkgCacheMissesInc_params,
			MoqMetrics_ASTPkgCacheMissesInc_paramsKey,
			MoqMetrics_ASTPkgCacheMissesInc_results,
		](scene, adaptor2, config),
		Moq_ASTTotalLoadTimeInc: impl.NewMoq[
			*MoqMetrics_ASTTotalLoadTimeInc_adaptor,
			MoqMetrics_ASTTotalLoadTimeInc_params,
			MoqMetrics_ASTTotalLoadTimeInc_paramsKey,
			MoqMetrics_ASTTotalLoadTimeInc_results,
		](scene, adaptor3, config),
		Moq_ASTTotalDecorationTimeInc: impl.NewMoq[
			*MoqMetrics_ASTTotalDecorationTimeInc_adaptor,
			MoqMetrics_ASTTotalDecorationTimeInc_params,
			MoqMetrics_ASTTotalDecorationTimeInc_paramsKey,
			MoqMetrics_ASTTotalDecorationTimeInc_results,
		](scene, adaptor4, config),
		Moq_TotalProcessingTimeInc: impl.NewMoq[
			*MoqMetrics_TotalProcessingTimeInc_adaptor,
			MoqMetrics_TotalProcessingTimeInc_params,
			MoqMetrics_TotalProcessingTimeInc_paramsKey,
			MoqMetrics_TotalProcessingTimeInc_results,
		](scene, adaptor5, config),
		Moq_Finalize: impl.NewMoq[
			*MoqMetrics_Finalize_adaptor,
			MoqMetrics_Finalize_params,
			MoqMetrics_Finalize_paramsKey,
			MoqMetrics_Finalize_results,
		](scene, adaptor6, config),

		Runtime: MoqMetrics_runtime{ParameterIndexing: struct {
			ASTPkgCacheHitsInc        MoqMetrics_ASTPkgCacheHitsInc_paramIndexing
			ASTPkgCacheMissesInc      MoqMetrics_ASTPkgCacheMissesInc_paramIndexing
			ASTTotalLoadTimeInc       MoqMetrics_ASTTotalLoadTimeInc_paramIndexing
			ASTTotalDecorationTimeInc MoqMetrics_ASTTotalDecorationTimeInc_paramIndexing
			TotalProcessingTimeInc    MoqMetrics_TotalProcessingTimeInc_paramIndexing
			Finalize                  MoqMetrics_Finalize_paramIndexing
		}{
			ASTPkgCacheHitsInc:   MoqMetrics_ASTPkgCacheHitsInc_paramIndexing{},
			ASTPkgCacheMissesInc: MoqMetrics_ASTPkgCacheMissesInc_paramIndexing{},
			ASTTotalLoadTimeInc: MoqMetrics_ASTTotalLoadTimeInc_paramIndexing{
				D: moq.ParamIndexByValue,
			},
			ASTTotalDecorationTimeInc: MoqMetrics_ASTTotalDecorationTimeInc_paramIndexing{
				D: moq.ParamIndexByValue,
			},
			TotalProcessingTimeInc: MoqMetrics_TotalProcessingTimeInc_paramIndexing{
				D: moq.ParamIndexByValue,
			},
			Finalize: MoqMetrics_Finalize_paramIndexing{},
		}},
	}
	m.Moq.Moq = m

	adaptor1.Moq = m
	adaptor2.Moq = m
	adaptor3.Moq = m
	adaptor4.Moq = m
	adaptor5.Moq = m
	adaptor6.Moq = m

	scene.AddMoq(m)
	return m
}

// Mock returns the mock implementation of the Metrics type
func (m *MoqMetrics) Mock() *MoqMetrics_mock { return m.Moq }

func (m *MoqMetrics_mock) ASTPkgCacheHitsInc() {
	m.Moq.Moq_ASTPkgCacheHitsInc.Scene.T.Helper()
	params := MoqMetrics_ASTPkgCacheHitsInc_params{}

	m.Moq.Moq_ASTPkgCacheHitsInc.Function(params)
}

func (m *MoqMetrics_mock) ASTPkgCacheMissesInc() {
	m.Moq.Moq_ASTPkgCacheMissesInc.Scene.T.Helper()
	params := MoqMetrics_ASTPkgCacheMissesInc_params{}

	m.Moq.Moq_ASTPkgCacheMissesInc.Function(params)
}

func (m *MoqMetrics_mock) ASTTotalLoadTimeInc(d time.Duration) {
	m.Moq.Moq_ASTTotalLoadTimeInc.Scene.T.Helper()
	params := MoqMetrics_ASTTotalLoadTimeInc_params{
		D: d,
	}

	m.Moq.Moq_ASTTotalLoadTimeInc.Function(params)
}

func (m *MoqMetrics_mock) ASTTotalDecorationTimeInc(d time.Duration) {
	m.Moq.Moq_ASTTotalDecorationTimeInc.Scene.T.Helper()
	params := MoqMetrics_ASTTotalDecorationTimeInc_params{
		D: d,
	}

	m.Moq.Moq_ASTTotalDecorationTimeInc.Function(params)
}

func (m *MoqMetrics_mock) TotalProcessingTimeInc(d time.Duration) {
	m.Moq.Moq_TotalProcessingTimeInc.Scene.T.Helper()
	params := MoqMetrics_TotalProcessingTimeInc_params{
		D: d,
	}

	m.Moq.Moq_TotalProcessingTimeInc.Function(params)
}

func (m *MoqMetrics_mock) Finalize() {
	m.Moq.Moq_Finalize.Scene.T.Helper()
	params := MoqMetrics_Finalize_params{}

	m.Moq.Moq_Finalize.Function(params)
}

// OnCall returns the recorder implementation of the Metrics type
func (m *MoqMetrics) OnCall() *MoqMetrics_recorder {
	return &MoqMetrics_recorder{
		Moq: m,
	}
}

func (m *MoqMetrics_recorder) ASTPkgCacheHitsInc() *MoqMetrics_ASTPkgCacheHitsInc_recorder {
	return &MoqMetrics_ASTPkgCacheHitsInc_recorder{
		Recorder: m.Moq.Moq_ASTPkgCacheHitsInc.OnCall(MoqMetrics_ASTPkgCacheHitsInc_params{}),
	}
}

func (r *MoqMetrics_ASTPkgCacheHitsInc_recorder) Any() *MoqMetrics_ASTPkgCacheHitsInc_anyParams {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.IsAnyPermitted(true) {
		return nil
	}
	return &MoqMetrics_ASTPkgCacheHitsInc_anyParams{Recorder: r}
}

func (r *MoqMetrics_ASTPkgCacheHitsInc_recorder) Seq() *MoqMetrics_ASTPkgCacheHitsInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(true, "Seq", true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_ASTPkgCacheHitsInc_recorder) NoSeq() *MoqMetrics_ASTPkgCacheHitsInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(false, "NoSeq", true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_ASTPkgCacheHitsInc_recorder) ReturnResults() *MoqMetrics_ASTPkgCacheHitsInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.ReturnResults(MoqMetrics_ASTPkgCacheHitsInc_results{})
	return r
}

func (r *MoqMetrics_ASTPkgCacheHitsInc_recorder) AndDo(fn MoqMetrics_ASTPkgCacheHitsInc_doFn) *MoqMetrics_ASTPkgCacheHitsInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.AndDo(func(params MoqMetrics_ASTPkgCacheHitsInc_params) {
		fn()
	}, true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_ASTPkgCacheHitsInc_recorder) DoReturnResults(fn MoqMetrics_ASTPkgCacheHitsInc_doReturnFn) *MoqMetrics_ASTPkgCacheHitsInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.DoReturnResults(func(params MoqMetrics_ASTPkgCacheHitsInc_params) *MoqMetrics_ASTPkgCacheHitsInc_results {
		fn()
		return &MoqMetrics_ASTPkgCacheHitsInc_results{}
	})
	return r
}

func (r *MoqMetrics_ASTPkgCacheHitsInc_recorder) Repeat(repeaters ...moq.Repeater) *MoqMetrics_ASTPkgCacheHitsInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Repeat(repeaters, true) {
		return nil
	}
	return r
}

func (*MoqMetrics_ASTPkgCacheHitsInc_adaptor) PrettyParams(params MoqMetrics_ASTPkgCacheHitsInc_params) string {
	return fmt.Sprintf("ASTPkgCacheHitsInc()")
}

func (a *MoqMetrics_ASTPkgCacheHitsInc_adaptor) ParamsKey(params MoqMetrics_ASTPkgCacheHitsInc_params, anyParams uint64) MoqMetrics_ASTPkgCacheHitsInc_paramsKey {
	a.Moq.Moq_ASTPkgCacheHitsInc.Scene.T.Helper()
	return MoqMetrics_ASTPkgCacheHitsInc_paramsKey{
		Params: struct{}{},
		Hashes: struct{}{},
	}
}

func (m *MoqMetrics_recorder) ASTPkgCacheMissesInc() *MoqMetrics_ASTPkgCacheMissesInc_recorder {
	return &MoqMetrics_ASTPkgCacheMissesInc_recorder{
		Recorder: m.Moq.Moq_ASTPkgCacheMissesInc.OnCall(MoqMetrics_ASTPkgCacheMissesInc_params{}),
	}
}

func (r *MoqMetrics_ASTPkgCacheMissesInc_recorder) Any() *MoqMetrics_ASTPkgCacheMissesInc_anyParams {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.IsAnyPermitted(true) {
		return nil
	}
	return &MoqMetrics_ASTPkgCacheMissesInc_anyParams{Recorder: r}
}

func (r *MoqMetrics_ASTPkgCacheMissesInc_recorder) Seq() *MoqMetrics_ASTPkgCacheMissesInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(true, "Seq", true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_ASTPkgCacheMissesInc_recorder) NoSeq() *MoqMetrics_ASTPkgCacheMissesInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(false, "NoSeq", true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_ASTPkgCacheMissesInc_recorder) ReturnResults() *MoqMetrics_ASTPkgCacheMissesInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.ReturnResults(MoqMetrics_ASTPkgCacheMissesInc_results{})
	return r
}

func (r *MoqMetrics_ASTPkgCacheMissesInc_recorder) AndDo(fn MoqMetrics_ASTPkgCacheMissesInc_doFn) *MoqMetrics_ASTPkgCacheMissesInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.AndDo(func(params MoqMetrics_ASTPkgCacheMissesInc_params) {
		fn()
	}, true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_ASTPkgCacheMissesInc_recorder) DoReturnResults(fn MoqMetrics_ASTPkgCacheMissesInc_doReturnFn) *MoqMetrics_ASTPkgCacheMissesInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.DoReturnResults(func(params MoqMetrics_ASTPkgCacheMissesInc_params) *MoqMetrics_ASTPkgCacheMissesInc_results {
		fn()
		return &MoqMetrics_ASTPkgCacheMissesInc_results{}
	})
	return r
}

func (r *MoqMetrics_ASTPkgCacheMissesInc_recorder) Repeat(repeaters ...moq.Repeater) *MoqMetrics_ASTPkgCacheMissesInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Repeat(repeaters, true) {
		return nil
	}
	return r
}

func (*MoqMetrics_ASTPkgCacheMissesInc_adaptor) PrettyParams(params MoqMetrics_ASTPkgCacheMissesInc_params) string {
	return fmt.Sprintf("ASTPkgCacheMissesInc()")
}

func (a *MoqMetrics_ASTPkgCacheMissesInc_adaptor) ParamsKey(params MoqMetrics_ASTPkgCacheMissesInc_params, anyParams uint64) MoqMetrics_ASTPkgCacheMissesInc_paramsKey {
	a.Moq.Moq_ASTPkgCacheMissesInc.Scene.T.Helper()
	return MoqMetrics_ASTPkgCacheMissesInc_paramsKey{
		Params: struct{}{},
		Hashes: struct{}{},
	}
}

func (m *MoqMetrics_recorder) ASTTotalLoadTimeInc(d time.Duration) *MoqMetrics_ASTTotalLoadTimeInc_recorder {
	return &MoqMetrics_ASTTotalLoadTimeInc_recorder{
		Recorder: m.Moq.Moq_ASTTotalLoadTimeInc.OnCall(MoqMetrics_ASTTotalLoadTimeInc_params{
			D: d,
		}),
	}
}

func (r *MoqMetrics_ASTTotalLoadTimeInc_recorder) Any() *MoqMetrics_ASTTotalLoadTimeInc_anyParams {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.IsAnyPermitted(true) {
		return nil
	}
	return &MoqMetrics_ASTTotalLoadTimeInc_anyParams{Recorder: r}
}

func (a *MoqMetrics_ASTTotalLoadTimeInc_anyParams) D() *MoqMetrics_ASTTotalLoadTimeInc_recorder {
	a.Recorder.Recorder.AnyParam(1)
	return a.Recorder
}

func (r *MoqMetrics_ASTTotalLoadTimeInc_recorder) Seq() *MoqMetrics_ASTTotalLoadTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(true, "Seq", true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_ASTTotalLoadTimeInc_recorder) NoSeq() *MoqMetrics_ASTTotalLoadTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(false, "NoSeq", true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_ASTTotalLoadTimeInc_recorder) ReturnResults() *MoqMetrics_ASTTotalLoadTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.ReturnResults(MoqMetrics_ASTTotalLoadTimeInc_results{})
	return r
}

func (r *MoqMetrics_ASTTotalLoadTimeInc_recorder) AndDo(fn MoqMetrics_ASTTotalLoadTimeInc_doFn) *MoqMetrics_ASTTotalLoadTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.AndDo(func(params MoqMetrics_ASTTotalLoadTimeInc_params) {
		fn(params.D)
	}, true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_ASTTotalLoadTimeInc_recorder) DoReturnResults(fn MoqMetrics_ASTTotalLoadTimeInc_doReturnFn) *MoqMetrics_ASTTotalLoadTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.DoReturnResults(func(params MoqMetrics_ASTTotalLoadTimeInc_params) *MoqMetrics_ASTTotalLoadTimeInc_results {
		fn(params.D)
		return &MoqMetrics_ASTTotalLoadTimeInc_results{}
	})
	return r
}

func (r *MoqMetrics_ASTTotalLoadTimeInc_recorder) Repeat(repeaters ...moq.Repeater) *MoqMetrics_ASTTotalLoadTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Repeat(repeaters, true) {
		return nil
	}
	return r
}

func (*MoqMetrics_ASTTotalLoadTimeInc_adaptor) PrettyParams(params MoqMetrics_ASTTotalLoadTimeInc_params) string {
	return fmt.Sprintf("ASTTotalLoadTimeInc(%#v)", params.D)
}

func (a *MoqMetrics_ASTTotalLoadTimeInc_adaptor) ParamsKey(params MoqMetrics_ASTTotalLoadTimeInc_params, anyParams uint64) MoqMetrics_ASTTotalLoadTimeInc_paramsKey {
	a.Moq.Moq_ASTTotalLoadTimeInc.Scene.T.Helper()
	dUsed, dUsedHash := impl.ParamKey(
		params.D, 1, a.Moq.Runtime.ParameterIndexing.ASTTotalLoadTimeInc.D, anyParams)
	return MoqMetrics_ASTTotalLoadTimeInc_paramsKey{
		Params: struct{ D time.Duration }{
			D: dUsed,
		},
		Hashes: struct{ D hash.Hash }{
			D: dUsedHash,
		},
	}
}

func (m *MoqMetrics_recorder) ASTTotalDecorationTimeInc(d time.Duration) *MoqMetrics_ASTTotalDecorationTimeInc_recorder {
	return &MoqMetrics_ASTTotalDecorationTimeInc_recorder{
		Recorder: m.Moq.Moq_ASTTotalDecorationTimeInc.OnCall(MoqMetrics_ASTTotalDecorationTimeInc_params{
			D: d,
		}),
	}
}

func (r *MoqMetrics_ASTTotalDecorationTimeInc_recorder) Any() *MoqMetrics_ASTTotalDecorationTimeInc_anyParams {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.IsAnyPermitted(true) {
		return nil
	}
	return &MoqMetrics_ASTTotalDecorationTimeInc_anyParams{Recorder: r}
}

func (a *MoqMetrics_ASTTotalDecorationTimeInc_anyParams) D() *MoqMetrics_ASTTotalDecorationTimeInc_recorder {
	a.Recorder.Recorder.AnyParam(1)
	return a.Recorder
}

func (r *MoqMetrics_ASTTotalDecorationTimeInc_recorder) Seq() *MoqMetrics_ASTTotalDecorationTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(true, "Seq", true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_ASTTotalDecorationTimeInc_recorder) NoSeq() *MoqMetrics_ASTTotalDecorationTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(false, "NoSeq", true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_ASTTotalDecorationTimeInc_recorder) ReturnResults() *MoqMetrics_ASTTotalDecorationTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.ReturnResults(MoqMetrics_ASTTotalDecorationTimeInc_results{})
	return r
}

func (r *MoqMetrics_ASTTotalDecorationTimeInc_recorder) AndDo(fn MoqMetrics_ASTTotalDecorationTimeInc_doFn) *MoqMetrics_ASTTotalDecorationTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.AndDo(func(params MoqMetrics_ASTTotalDecorationTimeInc_params) {
		fn(params.D)
	}, true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_ASTTotalDecorationTimeInc_recorder) DoReturnResults(fn MoqMetrics_ASTTotalDecorationTimeInc_doReturnFn) *MoqMetrics_ASTTotalDecorationTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.DoReturnResults(func(params MoqMetrics_ASTTotalDecorationTimeInc_params) *MoqMetrics_ASTTotalDecorationTimeInc_results {
		fn(params.D)
		return &MoqMetrics_ASTTotalDecorationTimeInc_results{}
	})
	return r
}

func (r *MoqMetrics_ASTTotalDecorationTimeInc_recorder) Repeat(repeaters ...moq.Repeater) *MoqMetrics_ASTTotalDecorationTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Repeat(repeaters, true) {
		return nil
	}
	return r
}

func (*MoqMetrics_ASTTotalDecorationTimeInc_adaptor) PrettyParams(params MoqMetrics_ASTTotalDecorationTimeInc_params) string {
	return fmt.Sprintf("ASTTotalDecorationTimeInc(%#v)", params.D)
}

func (a *MoqMetrics_ASTTotalDecorationTimeInc_adaptor) ParamsKey(params MoqMetrics_ASTTotalDecorationTimeInc_params, anyParams uint64) MoqMetrics_ASTTotalDecorationTimeInc_paramsKey {
	a.Moq.Moq_ASTTotalDecorationTimeInc.Scene.T.Helper()
	dUsed, dUsedHash := impl.ParamKey(
		params.D, 1, a.Moq.Runtime.ParameterIndexing.ASTTotalDecorationTimeInc.D, anyParams)
	return MoqMetrics_ASTTotalDecorationTimeInc_paramsKey{
		Params: struct{ D time.Duration }{
			D: dUsed,
		},
		Hashes: struct{ D hash.Hash }{
			D: dUsedHash,
		},
	}
}

func (m *MoqMetrics_recorder) TotalProcessingTimeInc(d time.Duration) *MoqMetrics_TotalProcessingTimeInc_recorder {
	return &MoqMetrics_TotalProcessingTimeInc_recorder{
		Recorder: m.Moq.Moq_TotalProcessingTimeInc.OnCall(MoqMetrics_TotalProcessingTimeInc_params{
			D: d,
		}),
	}
}

func (r *MoqMetrics_TotalProcessingTimeInc_recorder) Any() *MoqMetrics_TotalProcessingTimeInc_anyParams {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.IsAnyPermitted(true) {
		return nil
	}
	return &MoqMetrics_TotalProcessingTimeInc_anyParams{Recorder: r}
}

func (a *MoqMetrics_TotalProcessingTimeInc_anyParams) D() *MoqMetrics_TotalProcessingTimeInc_recorder {
	a.Recorder.Recorder.AnyParam(1)
	return a.Recorder
}

func (r *MoqMetrics_TotalProcessingTimeInc_recorder) Seq() *MoqMetrics_TotalProcessingTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(true, "Seq", true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_TotalProcessingTimeInc_recorder) NoSeq() *MoqMetrics_TotalProcessingTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(false, "NoSeq", true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_TotalProcessingTimeInc_recorder) ReturnResults() *MoqMetrics_TotalProcessingTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.ReturnResults(MoqMetrics_TotalProcessingTimeInc_results{})
	return r
}

func (r *MoqMetrics_TotalProcessingTimeInc_recorder) AndDo(fn MoqMetrics_TotalProcessingTimeInc_doFn) *MoqMetrics_TotalProcessingTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.AndDo(func(params MoqMetrics_TotalProcessingTimeInc_params) {
		fn(params.D)
	}, true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_TotalProcessingTimeInc_recorder) DoReturnResults(fn MoqMetrics_TotalProcessingTimeInc_doReturnFn) *MoqMetrics_TotalProcessingTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.DoReturnResults(func(params MoqMetrics_TotalProcessingTimeInc_params) *MoqMetrics_TotalProcessingTimeInc_results {
		fn(params.D)
		return &MoqMetrics_TotalProcessingTimeInc_results{}
	})
	return r
}

func (r *MoqMetrics_TotalProcessingTimeInc_recorder) Repeat(repeaters ...moq.Repeater) *MoqMetrics_TotalProcessingTimeInc_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Repeat(repeaters, true) {
		return nil
	}
	return r
}

func (*MoqMetrics_TotalProcessingTimeInc_adaptor) PrettyParams(params MoqMetrics_TotalProcessingTimeInc_params) string {
	return fmt.Sprintf("TotalProcessingTimeInc(%#v)", params.D)
}

func (a *MoqMetrics_TotalProcessingTimeInc_adaptor) ParamsKey(params MoqMetrics_TotalProcessingTimeInc_params, anyParams uint64) MoqMetrics_TotalProcessingTimeInc_paramsKey {
	a.Moq.Moq_TotalProcessingTimeInc.Scene.T.Helper()
	dUsed, dUsedHash := impl.ParamKey(
		params.D, 1, a.Moq.Runtime.ParameterIndexing.TotalProcessingTimeInc.D, anyParams)
	return MoqMetrics_TotalProcessingTimeInc_paramsKey{
		Params: struct{ D time.Duration }{
			D: dUsed,
		},
		Hashes: struct{ D hash.Hash }{
			D: dUsedHash,
		},
	}
}

func (m *MoqMetrics_recorder) Finalize() *MoqMetrics_Finalize_recorder {
	return &MoqMetrics_Finalize_recorder{
		Recorder: m.Moq.Moq_Finalize.OnCall(MoqMetrics_Finalize_params{}),
	}
}

func (r *MoqMetrics_Finalize_recorder) Any() *MoqMetrics_Finalize_anyParams {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.IsAnyPermitted(true) {
		return nil
	}
	return &MoqMetrics_Finalize_anyParams{Recorder: r}
}

func (r *MoqMetrics_Finalize_recorder) Seq() *MoqMetrics_Finalize_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(true, "Seq", true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_Finalize_recorder) NoSeq() *MoqMetrics_Finalize_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Seq(false, "NoSeq", true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_Finalize_recorder) ReturnResults() *MoqMetrics_Finalize_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.ReturnResults(MoqMetrics_Finalize_results{})
	return r
}

func (r *MoqMetrics_Finalize_recorder) AndDo(fn MoqMetrics_Finalize_doFn) *MoqMetrics_Finalize_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.AndDo(func(params MoqMetrics_Finalize_params) {
		fn()
	}, true) {
		return nil
	}
	return r
}

func (r *MoqMetrics_Finalize_recorder) DoReturnResults(fn MoqMetrics_Finalize_doReturnFn) *MoqMetrics_Finalize_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	r.Recorder.DoReturnResults(func(params MoqMetrics_Finalize_params) *MoqMetrics_Finalize_results {
		fn()
		return &MoqMetrics_Finalize_results{}
	})
	return r
}

func (r *MoqMetrics_Finalize_recorder) Repeat(repeaters ...moq.Repeater) *MoqMetrics_Finalize_recorder {
	r.Recorder.Moq.Scene.T.Helper()
	if !r.Recorder.Repeat(repeaters, true) {
		return nil
	}
	return r
}

func (*MoqMetrics_Finalize_adaptor) PrettyParams(params MoqMetrics_Finalize_params) string {
	return fmt.Sprintf("Finalize()")
}

func (a *MoqMetrics_Finalize_adaptor) ParamsKey(params MoqMetrics_Finalize_params, anyParams uint64) MoqMetrics_Finalize_paramsKey {
	a.Moq.Moq_Finalize.Scene.T.Helper()
	return MoqMetrics_Finalize_paramsKey{
		Params: struct{}{},
		Hashes: struct{}{},
	}
}

// Reset resets the state of the moq
func (m *MoqMetrics) Reset() {
	m.Moq_ASTPkgCacheHitsInc.Reset()
	m.Moq_ASTPkgCacheMissesInc.Reset()
	m.Moq_ASTTotalLoadTimeInc.Reset()
	m.Moq_ASTTotalDecorationTimeInc.Reset()
	m.Moq_TotalProcessingTimeInc.Reset()
	m.Moq_Finalize.Reset()
}

// AssertExpectationsMet asserts that all expectations have been met
func (m *MoqMetrics) AssertExpectationsMet() {
	m.Moq_ASTPkgCacheHitsInc.Scene.T.Helper()
	m.Moq_ASTPkgCacheHitsInc.AssertExpectationsMet()
	m.Moq_ASTPkgCacheMissesInc.AssertExpectationsMet()
	m.Moq_ASTTotalLoadTimeInc.AssertExpectationsMet()
	m.Moq_ASTTotalDecorationTimeInc.AssertExpectationsMet()
	m.Moq_TotalProcessingTimeInc.AssertExpectationsMet()
	m.Moq_Finalize.AssertExpectationsMet()
}
