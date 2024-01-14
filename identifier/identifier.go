package identifier

// Identifier object.
type Identifier interface {
	Ports() ([]string, error)
	Check([]byte) bool
}

// New initialize Identifier object.
func New(config Config) (Identifier, error) {
	return newRecognition(config)
}
