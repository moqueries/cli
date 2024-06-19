package export

type Type3 struct{}

type Widget struct{}

func (Widget) Type5() {}

func (Widget) Type6() {}

type WidgetByRef struct{}

func (*WidgetByRef) Type7() {}

func (*WidgetByRef) Type8() {}
