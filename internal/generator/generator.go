package generator

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString() (string, error) {
	data := make([]byte, 512)

	_, err := rand.Read(data)

	if err != nil {
		return "", err
	}

	h := md5.New()

	h.Write(data)

	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return hash, nil
}

func GenerateUserID() (string, error) {
	data := make([]byte, 512)

	_, err := rand.Read(data)

	if err != nil {
		return "", err
	}

	h := md5.New()

	h.Write(data)

	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return hash, nil
}
