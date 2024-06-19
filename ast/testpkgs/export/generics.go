package export

type GenericByRef[T any] struct{}
type Generic[T any] struct{}

func (g *GenericByRef[T]) DoSomethingPtr() {}

func (g *GenericByRef[X]) DoSomethingElsePtr() {}

func (g Generic[T]) DoSomething() {}

func (g Generic[X]) DoSomethingElse() {}

type GenericListByRef[T any, V any] struct{}
type GenericList[T any, V any] struct{}

func (g *GenericListByRef[T, V]) DoSomethingPtr() {}

func (g *GenericListByRef[X, Y]) DoSomethingElsePtr() {}

func (g GenericList[T, V]) DoSomething() {}

func (g GenericList[X, Y]) DoSomethingElse() {}
