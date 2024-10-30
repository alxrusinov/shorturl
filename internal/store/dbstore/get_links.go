package dbstore

import (
	"context"

	"github.com/alxrusinov/shorturl/internal/model"
)

// GetLinks returns information about all users links from data base
func (store *DBStore) GetLinks(userID string) ([]model.StoreRecord, error) {
	rows, err := store.db.QueryContext(context.Background(), "SELECT user_id, short, original, correlation_id, is_deleted FROM links WHERE user_id = $1", userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []model.StoreRecord

	for rows.Next() {
		var row model.StoreRecord

		if err := rows.Scan(&row.UUID, &row.ShortLink, &row.OriginalLink, &row.CorrelationID, &row.Deleted); err != nil {
			return nil, err
		}

		result = append(result, row)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil

}
