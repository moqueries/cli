// Package testmoqs contains multiple test mocks and adaptors for use in unit
// testing
package testmoqs

// StandaloneFunc is used to test that standalone functions can be mocked
func StandaloneFunc(_ string, bParam bool) (string, error) {
	return "", nil
}

type PassByValueSimple struct{}

func (PassByValueSimple) Usual(string, bool) (string, error) {
	return "", nil
}

type PassByRefSimple struct{}

func (*PassByRefSimple) Usual(string, bool) (string, error) {
	return "", nil
}

// Reduced creates a mock with an embedded reduced interface with only the
// exported methods mocked when using the ExcludeNonExported flag
type Reduced interface {
	Usual(sParam string, bParam bool) (sResult string, err error)
	notSoUsual()
	ReallyUnusualParams(struct{ a string })
	ReallyUnusualResults() struct{ a string }
}

type Generic[T any, V any] struct{}

func (g *Generic[T, V]) DoSomethingPtr() {}

func (g *Generic[X, Y]) DoSomethingElsePtr() {}

func (g Generic[T, V]) DoSomething() {}

func (g Generic[X, Y]) DoSomethingElse() {}

type Param struct{}

func DoSomethingWithParam(Param) {}
