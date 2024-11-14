package generator

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
)

// Generator has methods for creating random strings
type Generator struct{}

// GenerateRandomString generate random string
func (g *Generator) GenerateRandomString() (string, error) {
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
func (g *Generator) GenerateUserID() (string, error) {
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

// NewGenerator generate generator
func NewGenerator() *Generator {
	return &Generator{}
}
