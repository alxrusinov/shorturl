package dbstore

import (
	"context"

	"github.com/alxrusinov/shorturl/internal/model"
)

// GetLink returns original link from data base by shorten
func (store *DBStore) GetLink(arg *model.StoreRecord) (*model.StoreRecord, error) {
	var s string
	err := store.db.QueryRowContext(context.Background(), "SELECT original FROM links WHERE short = $1 and is_deleted = FALSE", arg.ShortLink).Scan(&s)

	if err != nil {
		return nil, err
	}

	arg.OriginalLink = s

	return arg, nil
}
