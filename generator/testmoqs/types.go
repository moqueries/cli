package testmoqs

import "io"

// NB: Keep in sync with ../generator_test.go TestGenerating

//nolint:lll // no easy way to break up go:generate lines
//go:generate moqueries --destination moq_testmoqs_test.go UsualFn NoNamesFn NoResultsFn NoParamsFn NothingFn VariadicFn RepeatedIdsFn TimesFn DifficultParamNamesFn DifficultResultNamesFn PassByReferenceFn InterfaceParamFn InterfaceResultFn Usual
//nolint:lll // no easy way to break up go:generate lines
//go:generate moqueries --destination exported/moq_exported_testmoqs.go --export UsualFn NoNamesFn NoResultsFn NoParamsFn NothingFn VariadicFn RepeatedIdsFn TimesFn DifficultParamNamesFn DifficultResultNamesFn PassByReferenceFn InterfaceParamFn InterfaceResultFn Usual

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
type DifficultParamNamesFn func(m, r bool, sequence string, param, params int, result, results float32)

// DifficultResultNamesFn has parameters with names that have been problematic
type DifficultResultNamesFn func() (m, r string, sequence error, param, params int, result, results float32)

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

func (r *InterfaceResultReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

// InterfaceResultFn tests returning interface results
type InterfaceResultFn func(sParam string, bParam bool) (r io.Reader)

// Usual combines all the above function types into an interface
type Usual interface {
	Usual(sParam string, bParam bool) (sResult string, err error)
	NoNames(string, bool) (string, error)
	NoResults(sParam string, bParam bool)
	NoParams() (sResult string, err error)
	Nothing()
	Variadic(other bool, args ...string) (sResult string, err error)
	RepeatedIds(sParam1, sParam2 string, bParam bool) (sResult1, sResult2 string, err error)
	Times(sParam string, times bool) (sResult string, err error)
	DifficultParamNames(m, r bool, sequence string, param, params int, result, results float32)
	DifficultResultNames() (m, r string, sequence error, param, params int, result, results float32)
	PassByReference(p *PassByReferenceParams) (sResult string, err error)
	InterfaceParam(w io.Writer) (sResult string, err error)
	InterfaceResult(sParam string, bParam bool) (r io.Reader)
}
