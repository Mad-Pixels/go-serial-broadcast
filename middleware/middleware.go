package middleware

import (
	"regexp"
)

type Middleware interface {
	Check([]byte) bool
	Key() []byte
}

func NewVerifyByMask(serialKey []byte, pattern string) (Middleware, error) {
	p, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &byMask{
		serialKey: serialKey,
		pattern:   p,
	}, nil
}
