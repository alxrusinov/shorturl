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

func (store *FileStore) GetLink(arg *StoreArgs) (*StoreArgs, error) {
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

func (store *FileStore) SetLink(arg *StoreArgs) (*StoreArgs, error) {
	file, err := os.OpenFile(store.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return nil, err
	}

	var rows []*Record

	err = json.Unmarshal([]byte(file), &rows)

	if err != nil {
		return nil, err
	}

	for _, val := range rows {
		if val.OriginalURL == arg.OriginalLink {
			arg.ShortLink = val.ShortURL
			return arg, nil
		}
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

func (store *FileStore) SetBatchLink(arg []*StoreArgs) ([]*StoreArgs, error) {
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
