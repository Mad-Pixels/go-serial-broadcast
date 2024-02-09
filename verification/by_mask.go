package verification

import (
	"regexp"
)

type byMask struct {
	serialKey []byte
	pattern   *regexp.Regexp
}

func (b *byMask) Check(msg []byte) bool {
	return b.pattern.Match(msg)
}

func (b *byMask) Key() []byte {
	return b.serialKey
}
