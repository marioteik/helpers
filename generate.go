package helpers

import (
	"crypto/rand"
	"errors"
	"io"
)

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func CompareCode(hashedCode, plainText string) error {
	if hashedCode != plainText {
		return errors.New("invalid code")
	}

	return nil
}
