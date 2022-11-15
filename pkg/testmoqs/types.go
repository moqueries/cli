package testmoqs

// StandaloneFunc is used to test that standalone functions can be mocked
func StandaloneFunc(_ string, bParam bool) (string, error) {
	return "", nil
}

type PassByValueSimple struct{}

func (PassByValueSimple) Usual(string, bool) (string, error) {
	return "", nil
}

type PassByRefSimple struct{}

func (*PassByRefSimple) Usual(string, bool) (string, error) {
	return "", nil
}
