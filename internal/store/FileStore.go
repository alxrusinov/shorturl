package store

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
)

type FileStore struct {
	filePath string
}

type Record struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (store *FileStore) GetLink(key string) (string, error) {
	file, err := os.OpenFile(store.filePath, os.O_RDONLY, 0666)

	if err != nil {
		return "", err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		record := &Record{}
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err == nil && record.ShortURL == key {
			return record.OriginalURL, nil
		}
	}

	if scanner.Err() != nil {
		return "", err
	}

	return "", errors.New("not found")
}

func (store *FileStore) SetLink(key string, link string) error {
	file, err := os.OpenFile(store.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	newUUID := uuid.NewString()

	record := &Record{
		UUID:        newUUID,
		OriginalURL: link,
		ShortURL:    key,
	}

	result, err := json.Marshal(record)

	if err != nil {
		return err
	}

	_, err = file.Write(append(result, []byte("\n")...))

	if err != nil {
		return err
	}

	return file.Close()

}

funct(store *FileStore) Ping() error {
file, err := os.OpenFile(store.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	return file.Close()
}

func CreateFileStore(filePath string) Store {
	store := &FileStore{filePath: filePath}

	return store
}
