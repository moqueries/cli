// Package testmoqs contains multiple test mocks and adaptors for use in unit
// testing
package testmoqs

import (
	"io"

	"moqueries.org/cli/generator/testmoqs/other"
)

// NB: Keep in sync with ../generator_test.go TestGenerating

//go:generate moqueries --destination moq_testmoqs_test.go UsualFn NoNamesFn NoResultsFn NoParamsFn NothingFn VariadicFn RepeatedIdsFn TimesFn DifficultParamNamesFn DifficultResultNamesFn PassByArrayFn PassByChanFn PassByEllipsisFn PassByMapFn PassByReferenceFn PassBySliceFn PassByValueFn InterfaceParamFn InterfaceResultFn Usual
//go:generate moqueries --destination exported/moq_exported_testmoqs.go --export UsualFn NoNamesFn NoResultsFn NoParamsFn NothingFn VariadicFn RepeatedIdsFn TimesFn DifficultParamNamesFn DifficultResultNamesFn PassByArrayFn PassByChanFn PassByEllipsisFn PassByMapFn PassByReferenceFn PassBySliceFn PassByValueFn InterfaceParamFn InterfaceResultFn Usual

// UsualFn is a typical function type
type UsualFn func(sParam string, bParam bool) (sResult string, err error)

// NoNamesFn is a typical function type
type NoNamesFn func(string, bool) (string, error)

// NoResultsFn is a function with no return values
type NoResultsFn func(sParam string, bParam bool)

// NoParamsFn is a function with no parameters
type NoParamsFn func() (sResult string, err error)

// NothingFn is a function with no parameters and no return values
type NothingFn func()

// VariadicFn is a function with a variable number of arguments
type VariadicFn func(other bool, args ...string) (sResult string, err error)

// RepeatedIdsFn is a function with multiple arguments of the same type
type RepeatedIdsFn func(sParam1, sParam2 string, bParam bool) (sResult1, sResult2 string, err error)

// TimesFn takes a parameter called times which should generate a valid moq
type TimesFn func(times string, bParam bool) (sResult string, err error)

// DifficultParamNamesFn has parameters with names that have been problematic
type DifficultParamNamesFn func(m, r bool, sequence string, param, params, i int, result, results, _ float32)

// DifficultResultNamesFn has parameters with names that have been problematic
type DifficultResultNamesFn func() (m, r string, sequence error, param, params, i int, result, results, _ float32)

// Params encapsulates the parameters for use in various test types
type Params struct {
	SParam string
	BParam bool
}

// Results encapsulates the results for use in various test types
type Results struct {
	SResult string
	Err     error
}

// PassByArrayFn tests passing parameters and results by array
type PassByArrayFn func(p [3]Params) [3]Results

// PassByChanFn tests passing parameters and results by channel
type PassByChanFn func(p chan Params) chan Results

// PassByEllipsisFn tests passing parameters by ellipsis
type PassByEllipsisFn func(p ...Params) (string, error)

// PassByMapFn tests passing parameters and results by map
type PassByMapFn func(p map[string]Params) map[string]Results

// PassByReferenceFn tests passing parameters and results by reference
type PassByReferenceFn func(p *Params) *Results

// PassBySliceFn tests passing parameters and results by slice
type PassBySliceFn func(p []Params) []Results

// PassByValueFn tests passing parameters and results by reference
type PassByValueFn func(p Params) Results

// InterfaceParamWriter is used for testing functions that take an interface as
// a parameter
type InterfaceParamWriter struct {
	SParam string
	BParam bool
}

func (w *InterfaceParamWriter) Write([]byte) (int, error) {
	return 0, nil
}

// InterfaceParamFn tests passing interface parameters
type InterfaceParamFn func(w io.Writer) (sResult string, err error)

// InterfaceResultReader is used for testing functions that return an interface
type InterfaceResultReader struct {
	SResult string
	Err     error
}

func (r *InterfaceResultReader) Read(p []byte) (int, error) {
	return 0, nil
}

// InterfaceResultFn tests returning interface results
type InterfaceResultFn func(sParam string, bParam bool) (r io.Reader)

// Usual combines all the above function types into an interface
//
//nolint:interfacebloat // Test interface with one of every method type
type Usual interface {
	Usual(sParam string, bParam bool) (sResult string, err error)
	//nolint:inamedparam // Testing interface method with unnamed param for type string
	NoNames(string, bool) (string, error)
	NoResults(sParam string, bParam bool)
	NoParams() (sResult string, err error)
	Nothing()
	Variadic(other bool, args ...string) (sResult string, err error)
	RepeatedIds(sParam1, sParam2 string, bParam bool) (sResult1, sResult2 string, err error)
	Times(sParam string, times bool) (sResult string, err error)
	DifficultParamNames(m, r bool, sequence string, param, params, i int, result, results, _ float32)
	DifficultResultNames() (m, r string, sequence error, param, params, i int, result, results, _ float32)
	PassByArray(p [3]Params) [3]Results
	PassByChan(p chan Params) chan Results
	PassByEllipsis(p ...Params) (string, error)
	PassByMap(p map[string]Params) map[string]Results
	PassByReference(p *Params) *Results
	PassBySlice(p []Params) []Results
	PassByValue(p Params) Results
	InterfaceParam(w io.Writer) (sResult string, err error)
	InterfaceResult(sParam string, bParam bool) (r io.Reader)
	FnParam(fn func())
	other.Other
}
