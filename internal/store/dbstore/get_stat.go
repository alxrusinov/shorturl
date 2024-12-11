package dbstore

import (
	"github.com/alxrusinov/shorturl/internal/model"
)

// GetStat - gets dtatistics of urls and users
func (store *DBStore) GetStat() (*model.StatResponse, error) {
	result := new(model.StatResponse)

	err := store.db.QueryRow("SELECT count(*) AS users, count(DISTINCT user_id) AS links FROM links").Scan(&result.URLS, &result.Users)

	if err != nil {
		return nil, err
	}

	return result, nil
}
