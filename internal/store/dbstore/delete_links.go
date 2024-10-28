package dbstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/alxrusinov/shorturl/internal/model"
)

func (store *DBStore) DeleteLinks(shorts [][]model.StoreRecord) error {
	tx, err := store.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Commit()

	preparedShorts := []string{}
	preparedIDs := []string{}

	for _, val := range shorts {
		userID := val[0].UUID
		preparedIDs = append(preparedIDs, fmt.Sprint("'"+userID+"'"))

		for _, shortLink := range val {
			preparedShorts = append(preparedShorts, fmt.Sprint("'"+shortLink.ShortLink+"'"))
		}
	}

	shortsPlaceholders := strings.Join(preparedShorts, ", ")
	userIDsPlaceholders := strings.Join(preparedIDs, ", ")

	rows, err := tx.QueryContext(context.Background(), `UPDATE links SET is_deleted = TRUE WHERE user_id = ANY(ARRAY[`+userIDsPlaceholders+`]) and short = ANY(ARRAY[`+shortsPlaceholders+`]) RETURNING id;`)

	if err != nil {
		return err
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil

}
