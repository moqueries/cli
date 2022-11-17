package noexport_test

type test_type3 struct{}

type test_widget struct{}

func (test_widget) test_method1() {}

func (test_widget) test_method2() {}

func (*test_widget) test_method3() {}

func (*test_widget) test_method4() {}
