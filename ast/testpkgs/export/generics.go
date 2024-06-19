package export

type GenericByRef[T any] struct{}

func (g *GenericByRef[T]) DoSomethingPtr() {}

func (g *GenericByRef[X]) DoSomethingElsePtr() {}

type Generic[T any] struct{}

func (g Generic[T]) DoSomething() {}

func (g Generic[X]) DoSomethingElse() {}

type GenericListByRef[T any, V any] struct{}

func (g *GenericListByRef[T, V]) DoSomethingPtr() {}

func (g *GenericListByRef[X, Y]) DoSomethingElsePtr() {}

type GenericList[T any, V any] struct{}

func (g GenericList[T, V]) DoSomething() {}

func (g GenericList[X, Y]) DoSomethingElse() {}
