package fabricate

type Mixed int

func (Mixed) nonExportedByValueRecv() {}

func (*Mixed) ExportedByRefRecv() {}

type CantMock int

func (CantMock) ByValueRecv() {}

func (*CantMock) ByRefRecv() {}
