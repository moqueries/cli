package noexport

type genericByRef[T any] struct{}
type generic[T any] struct{}

func (g *genericByRef[T]) doSomethingPtr() {}

func (g *genericByRef[X]) doSomethingElsePtr() {}

func (g generic[T]) doSomething() {}

func (g generic[X]) doSomethingElse() {}

type genericListByRef[T any, V any] struct{}
type genericList[T any, V any] struct{}

func (g *genericListByRef[T, V]) doSomethingPtr() {}

func (g *genericListByRef[X, Y]) doSomethingElsePtr() {}

func (g genericList[T, V]) doSomething() {}

func (g genericList[X, Y]) doSomethingElse() {}
