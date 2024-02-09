package verification

import (
	"regexp"
)

type Interface interface {
	Check([]byte) bool
	Key() []byte
}

func NewByMask(serialKey []byte, pattern string) (Interface, error) {
	p, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &byMask{
		serialKey: serialKey,
		pattern:   p,
	}, nil
}
