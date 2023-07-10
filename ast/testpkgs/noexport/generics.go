package noexport

type generic[T any] struct{}

func (g *generic[T]) doSomethingPtr() {}

func (g *generic[X]) doSomethingElsePtr() {}

func (g generic[T]) doSomething() {}

func (g generic[X]) doSomethingElse() {}

type genericList[T any, V any] struct{}

func (g *genericList[T, V]) doSomethingPtr() {}

func (g *genericList[X, Y]) doSomethingElsePtr() {}

func (g genericList[T, V]) doSomething() {}

func (g genericList[X, Y]) doSomethingElse() {}
