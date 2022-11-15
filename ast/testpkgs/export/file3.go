package export

type Type3 struct{}

type Widget struct{}

func (Widget) Type5() {}

func (Widget) Type6() {}

func (*Widget) Type7() {}

func (*Widget) Type8() {}
