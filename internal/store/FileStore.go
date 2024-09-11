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

func (store *FileStore) GetLink(arg *StoreRecord) (*StoreRecord, error) {
	file, err := os.OpenFile(store.filePath, os.O_RDONLY, 0666)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		record := &StoreRecord{}
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err == nil && record.ShortLink == arg.ShortLink {
			arg.OriginalLink = record.OriginalLink
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
		record := &StoreRecord{}
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err == nil && record.OriginalLink == arg.OriginalLink {
			arg.OriginalLink = record.OriginalLink
			return arg, &DuplicateValueError{Err: errors.New("record already exists")}
		}
	}

	if scanner.Err() != nil {
		return nil, err
	}

	newUUID := uuid.NewString()

	record := &StoreRecord{
		UUID:          newUUID,
		CorrelationID: arg.CorrelationID,
		OriginalLink:  arg.OriginalLink,
		ShortLink:     arg.ShortLink,
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

		record := &StoreRecord{
			UUID:          newUUID,
			CorrelationID: val.CorrelationID,
			OriginalLink:  val.OriginalLink,
			ShortLink:     val.ShortLink,
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

func (store *FileStore) GetLinks(userID string) ([]StoreRecord, error) {
	file, err := os.OpenFile(store.filePath, os.O_RDONLY, 0666)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var result []StoreRecord

	for scanner.Scan() {
		record := &StoreRecord{}
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err == nil && userID == record.UUID {
			result = append(result, *record)
		}
	}

	if scanner.Err() != nil {
		return nil, err
	}

	return result, nil

}

func (store *FileStore) DeleteLinks(shorts [][]StoreRecord) error {
	file, err := os.OpenFile(store.filePath, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	var preparedShorts []StoreRecord

	for _, val := range shorts {
		for _, short := range val {
			preparedShorts = append(preparedShorts, short)
		}
	}

	scanner := bufio.NewScanner(file)

	var result []StoreRecord

	for scanner.Scan() {
		record := &StoreRecord{}
		err := json.Unmarshal(scanner.Bytes(), &record)

		if err == nil {
			for _, rec := range preparedShorts {
				if rec.UUID == record.UUID && rec.ShortLink == record.ShortLink {
					record.Deleted = true
					break
				}
			}

			result = append(result, *record)
		}
	}

	if scanner.Err() != nil {
		return err
	}

	content, err := json.Marshal(&result)

	if err != nil {
		return err
	}

	_, err = file.Write(content)

	if err != nil {
		return err
	}

	return nil

}

func CreateFileStore(filePath string) Store {
	store := &FileStore{filePath: filePath}

	return store
}
