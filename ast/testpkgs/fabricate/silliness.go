package fabricate

// Mixed verifies that objects with mixed receivers won't get selected for
// mocking
//nolint: recvcheck // testing mixed receivers is the point
type Mixed int

func (*Mixed) ExportedByRefRecv() {}

func (Mixed) nonExportedByValueRecv() {}

// CantMock verifies that objects with mixed receivers won't get selected for
// mocking
//nolint: recvcheck // testing mixed receivers is the point
type CantMock int

func (CantMock) ByValueRecv() {}

func (*CantMock) ByRefRecv() {}
