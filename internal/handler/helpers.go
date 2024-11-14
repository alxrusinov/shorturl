package handler

import (
	"errors"
	"fmt"
)

type errReader struct{}

func (er *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func createShortLink(host string, shorten string) string {
	return fmt.Sprintf("%s/%s", host, shorten)
}

func newErrReader() *errReader {
	return &errReader{}
}
