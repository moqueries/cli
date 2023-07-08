package export

type Generic[T any] struct{}

func (g *Generic[T]) DoSomethingPtr() {}

func (g *Generic[X]) DoSomethingElsePtr() {}

func (g Generic[T]) DoSomething() {}

func (g Generic[X]) DoSomethingElse() {}

type GenericList[T any, V any] struct{}

func (g *GenericList[T, V]) DoSomethingPtr() {}

func (g *GenericList[X, Y]) DoSomethingElsePtr() {}

func (g GenericList[T, V]) DoSomething() {}

func (g GenericList[X, Y]) DoSomethingElse() {}
