package filestore

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"

	"github.com/alxrusinov/shorturl/internal/customerrors"
	"github.com/alxrusinov/shorturl/internal/model"
)

type FileStore struct {
	filePath string
}

func (store *FileStore) GetLink(arg *model.StoreRecord) (*model.StoreRecord, error) {
	file, err := os.OpenFile(store.filePath, os.O_RDONLY, 0666)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		record := &model.StoreRecord{}
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

func (store *FileStore) SetLink(arg *model.StoreRecord) (*model.StoreRecord, error) {
	file, err := os.OpenFile(store.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		record := &model.StoreRecord{}
		err := json.Unmarshal(scanner.Bytes(), &record)
		if err == nil && record.OriginalLink == arg.OriginalLink {
			arg.OriginalLink = record.OriginalLink
			return arg, &customerrors.DuplicateValueError{Err: errors.New("record already exists")}
		}
	}

	if scanner.Err() != nil {
		return nil, err
	}

	newUUID := uuid.NewString()

	record := &model.StoreRecord{
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

func (store *FileStore) SetBatchLink(arg []*model.StoreRecord) ([]*model.StoreRecord, error) {
	file, err := os.OpenFile(store.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return nil, err
	}

	newUUID := uuid.NewString()

	for _, val := range arg {

		record := &model.StoreRecord{
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

func (store *FileStore) GetLinks(userID string) ([]model.StoreRecord, error) {
	file, err := os.OpenFile(store.filePath, os.O_RDONLY, 0666)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var result []model.StoreRecord

	for scanner.Scan() {
		record := &model.StoreRecord{}
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

func (store *FileStore) DeleteLinks(shorts [][]model.StoreRecord) error {
	file, err := os.OpenFile(store.filePath, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	var preparedShorts []model.StoreRecord

	for _, val := range shorts {
		preparedShorts = append(preparedShorts, val...)
	}

	scanner := bufio.NewScanner(file)

	var result []model.StoreRecord

	for scanner.Scan() {
		record := &model.StoreRecord{}
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

func NewFileStore(filePath string) *FileStore {
	store := &FileStore{filePath: filePath}

	return store
}
