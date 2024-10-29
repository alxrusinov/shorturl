package dbstore

import (
	"context"
	"errors"
	"io"

	"github.com/alxrusinov/shorturl/internal/model"
)

// SetBatchLink adds links to data base by batch
func (store *DBStore) SetBatchLink(arg []*model.StoreRecord) ([]*model.StoreRecord, error) {
	tx, err := store.db.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	stmt := tx.Stmt(store.insertQuery)

	defer stmt.Close()

	response := make([]*model.StoreRecord, 0)

	for _, val := range arg {
		res := &model.StoreRecord{}
		err := stmt.QueryRowContext(context.Background(), val.ShortLink, val.OriginalLink, val.CorrelationID, val.UUID).Scan(&res.ShortLink, &res.OriginalLink, &res.CorrelationID, &res.UUID)

		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		response = append(response, res)

	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return response, nil
}
