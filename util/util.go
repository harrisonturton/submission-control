package util

import (
	"fmt"
	"os"
)

func Uuid() string {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	return fmt.Sprintf("%x", b[0:4])
}
