package store

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
)

type FileStore struct {
	FilePath string
}

type Record struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (store *FileStore) GetLink(key string) (string, error) {
	file, err := os.ReadFile(store.FilePath)

	if err != nil {
		return "", err
	}

	content := []Record{}

	err = json.Unmarshal(file, &content)

	if err != nil {
		return "", err
	}

	var result string

	for _, rec := range content {
		if rec.ShortURL == key {
			result = rec.OriginalURL
			break
		}
	}

	if result != "" {
		return result, nil
	}

	return result, errors.New("not found")
}

func (store *FileStore) SetLink(key string, link string) {
	file, err := os.OpenFile(store.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return
	}

	defer file.Close()

	newUUID := uuid.NewString()

	record := &Record{
		UUID:        newUUID,
		OriginalURL: link,
		ShortURL:    key,
	}

	result, err := json.Marshal(record)

	if err != nil {
		return
	}

	file.Write(result)

}

func CreateFileStore(filePath string) Store {
	store := &FileStore{FilePath: filePath}

	return store
}
