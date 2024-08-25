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
	UUID          string `json:"uuid"`
	ShortURL      string `json:"short_url"`
	OriginalURL   string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

func (store *FileStore) GetLink(arg *StoreRecord) (*StoreRecord, error) {
	file, err := os.OpenFile(store.filePath, os.O_RDONLY, 0666)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		record := &Record{}
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err == nil && record.ShortURL == arg.ShortLink {
			arg.OriginalLink = record.OriginalURL
			return arg, nil
		}
	}

	if scanner.Err() != nil {
		return nil, err
	}

	return nil, errors.New("not found")
}

func (store *FileStore) SetLink(arg *StoreRecord) (*StoreRecord, error) {
	file, err := os.OpenFile(store.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		record := &Record{}
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err == nil && record.OriginalURL == arg.OriginalLink {
			arg.OriginalLink = record.OriginalURL
			return arg, &DuplicateValueError{Err: errors.New("record already exists")}
		}
	}

	if scanner.Err() != nil {
		return nil, err
	}

	newUUID := uuid.NewString()

	record := &Record{
		UUID:          newUUID,
		CorrelationID: arg.CorrelationID,
		OriginalURL:   arg.OriginalLink,
		ShortURL:      arg.ShortLink,
	}

	result, err := json.Marshal(record)

	if err != nil {
		return nil, err
	}

	_, err = file.Write(append(result, []byte("\n")...))

	if err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return arg, nil

}

func (store *FileStore) Ping() error {
	file, err := os.OpenFile(store.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	return file.Close()
}

func (store *FileStore) SetBatchLink(arg []*StoreRecord) ([]*StoreRecord, error) {
	file, err := os.OpenFile(store.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return nil, err
	}

	newUUID := uuid.NewString()

	for _, val := range arg {

		record := &Record{
			UUID:          newUUID,
			CorrelationID: val.CorrelationID,
			OriginalURL:   val.OriginalLink,
			ShortURL:      val.ShortLink,
		}

		result, err := json.Marshal(record)

		if err != nil {
			return nil, err
		}

		_, err = file.Write(append(result, []byte("\n")...))

		if err != nil {
			return nil, err
		}

	}

	err = file.Close()

	if err != nil {
		return nil, err
	}

	return arg, nil

}

func CreateFileStore(filePath string) Store {
	store := &FileStore{filePath: filePath}

	return store
}
