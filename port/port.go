package port

type Interface interface {
	Read(b []byte) (int, error)
}

func NewPort() (Interface, error) {
	return NewSerial()
}
