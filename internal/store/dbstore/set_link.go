package dbstore

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/alxrusinov/shorturl/internal/customerrors"
	"github.com/alxrusinov/shorturl/internal/model"
)

func (store *DBStore) SetLink(arg *model.StoreRecord) (*model.StoreRecord, error) {
	var err error
	dbQuery := `INSERT INTO links (short, original, correlation_id, user_id)
				VALUES ($1, $2, $3, $4);
				`

	selectQuery := `SELECT short FROM links WHERE original = $1 `

	_, err = store.db.ExecContext(context.Background(), dbQuery, arg.ShortLink, arg.OriginalLink, arg.CorrelationID, arg.UUID)

	if err != nil {
		if dbErr, ok := err.(*pgconn.PgError); ok {
			if dbErr.Code == pgerrcode.UniqueViolation {

				err := store.db.QueryRowContext(context.Background(), selectQuery, arg.OriginalLink).Scan(&arg.ShortLink)

				if err != nil {
					return nil, err
				}

				return arg, &customerrors.DuplicateValueError{Err: dbErr}
			}
		}
		return nil, err
	}

	return arg, nil
}
