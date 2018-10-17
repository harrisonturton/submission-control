package util

import (
	"crypto/sha256"
	"encoding/hex"
)

type Digest string

func FromBytes(bytes []byte) Digest {
	h := sha256.New()
	h.Write(bytes)
	return hex.EncodeToString(h.Sum(nil))
}

func FromString(str string) Digest {
	return FromBytes([]byte(str))
}
