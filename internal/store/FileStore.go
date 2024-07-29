package store

import (
	"bufio"
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
	file, err := os.OpenFile(store.FilePath, os.O_RDONLY, 0666)

	if err != nil {
		return "", err
	}

	defer file.Close()

	var result string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		record := &Record{}
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err == nil && record.ShortURL == key {
			result = record.OriginalURL
			break
		}
	}

	if scanner.Err() != nil {
		return "", err
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

	file.Write(append(result, []byte("\n")...))

}

func CreateFileStore(filePath string) Store {
	store := &FileStore{FilePath: filePath}

	return store
}
