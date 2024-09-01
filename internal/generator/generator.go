package generator

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	newRand "math/rand"
	"time"
)

func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := newRand.NewSource(time.Now().UnixNano())
	random := newRand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[random.Intn(len(charset))]
	}
	return string(result), nil
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
