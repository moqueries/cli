// Package testmoqs contains multiple test mocks and adaptors for use in unit
// testing
package testmoqs

import "io"

// TODO: Need to have a full copy of this file in a test package?

// NB: Keep in sync with ../generator_test.go TestGenerating

//go:generate moqueries --destination moq_testmoqs_test.go UsualFn NoNamesFn NoResultsFn NoParamsFn NothingFn VariadicFn RepeatedIdsFn TimesFn DifficultParamNamesFn DifficultResultNamesFn PassByReferenceFn InterfaceParamFn InterfaceResultFn GenericParamsFn PartialGenericParamsFn GenericResultsFn PartialGenericResultsFn GenericInterfaceParamFn GenericInterfaceResultFn Usual GenericParams PartialGenericParams GenericResults PartialGenericResults GenericInterfaceParam GenericInterfaceResult
//go:generate moqueries --destination exported/moq_exported_testmoqs.go --export UsualFn NoNamesFn NoResultsFn NoParamsFn NothingFn VariadicFn RepeatedIdsFn TimesFn DifficultParamNamesFn DifficultResultNamesFn PassByReferenceFn InterfaceParamFn InterfaceResultFn GenericParamsFn PartialGenericParamsFn GenericResultsFn PartialGenericResultsFn GenericInterfaceParamFn GenericInterfaceResultFn Usual GenericParams PartialGenericParams GenericResults PartialGenericResults GenericInterfaceParam GenericInterfaceResult

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

// PassByReferenceParams encapsulates the parameters for passing by reference
// tests
type PassByReferenceParams struct {
	SParam string
	BParam bool
}

// PassByReferenceFn tests passing parameters by reference
type PassByReferenceFn func(p *PassByReferenceParams) (sResult string, err error)

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

func (r *InterfaceResultReader) Read([]byte) (int, error) {
	return 0, nil
}

// InterfaceResultFn tests returning interface results
type InterfaceResultFn func(sParam string, bParam bool) (r io.Reader)

// GenericParamsFn has all generic parameters
type GenericParamsFn[S, B any] func(S, B) (string, error)

// PartialGenericParamsFn has some generic parameters
type PartialGenericParamsFn[S any] func(S, bool) (string, error)

// GenericResultsFn has all generic results
type GenericResultsFn[S ~string, E error] func(string, bool) (S, E)

// PartialGenericResultsFn has some generic results
type PartialGenericResultsFn[S ~string] func(string, bool) (S, error)

type MyWriter io.Writer

// GenericInterfaceParamFn tests passing generic interface parameters
type GenericInterfaceParamFn[W MyWriter] func(w W) (sResult string, err error)

type MyReader io.Reader

// GenericInterfaceResultFn tests returning generic interface results
type GenericInterfaceResultFn[R MyReader] func(sParam string, bParam bool) (r R)

// Usual combines all the above function types into an interface
//
//nolint:interfacebloat // Test interface with one of every method type
type Usual interface {
	Usual(sParam string, bParam bool) (sResult string, err error)
	NoNames(string, bool) (string, error)
	NoResults(sParam string, bParam bool)
	NoParams() (sResult string, err error)
	Nothing()
	Variadic(other bool, args ...string) (sResult string, err error)
	RepeatedIds(sParam1, sParam2 string, bParam bool) (sResult1, sResult2 string, err error)
	Times(sParam string, times bool) (sResult string, err error)
	DifficultParamNames(m, r bool, sequence string, param, params, i int, result, results, _ float32)
	DifficultResultNames() (m, r string, sequence error, param, params, i int, result, results, _ float32)
	PassByReference(p *PassByReferenceParams) (sResult string, err error)
	InterfaceParam(w io.Writer) (sResult string, err error)
	InterfaceResult(sParam string, bParam bool) (r io.Reader)
	FnParam(fn func())
}

type GenericParams[S, B any] interface {
	Usual(S, B) (string, error)
}

type PartialGenericParams[S any] interface {
	Usual(S, bool) (string, error)
}

type GenericResults[S ~string, E error] interface {
	Usual(string, bool) (S, E)
}

type PartialGenericResults[S ~string] interface {
	Usual(string, bool) (S, error)
}

type GenericInterfaceParam[W MyWriter] interface {
	Usual(w W) (sResult string, err error)
}

type GenericInterfaceResult[R MyReader] interface {
	Usual(sParam string, bParam bool) (r R)
}
