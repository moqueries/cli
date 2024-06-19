package noexport_test

type test_type3 struct{}

type test_widget struct{}

func (test_widget) test_method1() {}

func (test_widget) test_method2() {}

type test_widgetByRef struct{}

func (*test_widgetByRef) test_method3() {}

func (*test_widgetByRef) test_method4() {}
