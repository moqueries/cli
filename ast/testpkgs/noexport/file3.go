package noexport

type type3 struct{}

type widget struct{}

func (widget) method1() {}

func (widget) method2() {}

func (*widget) method3() {}

func (*widget) method4() {}
