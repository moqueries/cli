package noexport

type type3 struct{}

type widget struct{}

func (widget) method1() {}

func (widget) method2() {}

type widgetByRef struct{}

func (*widgetByRef) method3() {}

func (*widgetByRef) method4() {}
