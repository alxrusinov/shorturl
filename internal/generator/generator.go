package generator

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomString generate random string
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

// GenerateUserID return uaer ID as string
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
