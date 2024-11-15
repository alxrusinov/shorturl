package inmemorystore

import (
	"testing"

	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryStore_GetLink(t *testing.T) {

	tests := []struct {
		name  string
		store *InMemoryStore
		arg   *model.StoreRecord
		want  *model.StoreRecord
		err   bool
	}{
		{
			name:  "1# success",
			store: NewInMemoryStore(),
			arg: &model.StoreRecord{
				ShortLink:    "123",
				OriginalLink: "http://example.com",
			},
		},
		{
			name:  "2# success",
			store: NewInMemoryStore(),
			arg: &model.StoreRecord{
				ShortLink:    "123",
				OriginalLink: "http://example.com",
			},
			err: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case tests[0].name:
				tt.store.data[tt.arg.ShortLink] = tt.arg
				got, err := tt.store.GetLink(tt.arg)

				assert.Equal(t, tt.arg.OriginalLink, got.OriginalLink)
				assert.Nil(t, err)

			case tests[1].name:
				got, err := tt.store.GetLink(tt.arg)

				assert.Nil(t, got)
				assert.Error(t, err)
			}

		})
	}
}

func TestInMemoryStore_SetLink(t *testing.T) {
	tests := []struct {
		name  string
		store *InMemoryStore
		arg   *model.StoreRecord
		want  *model.StoreRecord
		err   bool
	}{
		{
			name:  "1# success",
			store: NewInMemoryStore(),
			arg: &model.StoreRecord{
				ShortLink:    "123",
				OriginalLink: "http://example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.store.SetLink(tt.arg)

			assert.Equal(t, tt.store.data[tt.arg.ShortLink], got)
			assert.Nil(t, err)
		})
	}
}

func TestInMemoryStore_Ping(t *testing.T) {

	t.Run("1# success", func(t *testing.T) {
		store := NewInMemoryStore()

		err := store.Ping()

		assert.Nil(t, err)
	})

}

func TestInMemoryStore_SetBatchLink(t *testing.T) {

	t.Run("1# success", func(t *testing.T) {
		store := NewInMemoryStore()

		batch := make([]*model.StoreRecord, 0)

		batch = append(batch, &model.StoreRecord{
			ShortLink: "123",
		})

		got, err := store.SetBatchLink(batch)

		assert.Equal(t, batch, got)
		assert.Nil(t, err)
	})
}

func TestInMemoryStore_GetLinks(t *testing.T) {
	t.Run("1# success", func(t *testing.T) {
		store := NewInMemoryStore()

		userID := "1"
		anotherUserID := "2"

		batch := make([]*model.StoreRecord, 0)

		withUserID := &model.StoreRecord{
			ShortLink: "123",
			UUID:      userID,
		}

		withAnotherUserID := &model.StoreRecord{
			ShortLink: "321",
			UUID:      anotherUserID,
		}

		batch = append(batch, withUserID, withAnotherUserID)

		store.SetBatchLink(batch)

		got, err := store.GetLinks(userID)

		assert.Len(t, got, 1)
		assert.Nil(t, err)

		for _, val := range got {
			assert.Equal(t, val.UUID, userID)
		}
	})
}

func TestInMemoryStore_DeleteLinks(t *testing.T) {
	tests := []struct {
		name string
		err  bool
	}{
		{
			name: "1# success",
			err:  false,
		},
		{
			name: "2# error",
			err:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewInMemoryStore()

			userID := "1"
			anotherUserID := "2"

			shorts := make([][]model.StoreRecord, 0)
			batchOne := make([]*model.StoreRecord, 0)
			batchTwo := make([]*model.StoreRecord, 0)

			wrongRecord := model.StoreRecord{
				ShortLink: "777",
				UUID:      userID,
			}

			batchOne = append(batchOne, &model.StoreRecord{
				ShortLink: "111",
				UUID:      userID,
			}, &model.StoreRecord{
				ShortLink: "222",
				UUID:      anotherUserID,
			})

			batchTwo = append(batchTwo, &model.StoreRecord{
				ShortLink: "999",
				UUID:      userID,
			}, &model.StoreRecord{
				ShortLink: "888",
				UUID:      anotherUserID,
			})

			store.SetBatchLink(batchOne)
			store.SetBatchLink(batchTwo)

			switch tt.name {
			case tests[0].name:
				shorts = append(shorts, []model.StoreRecord{*batchOne[0]}, []model.StoreRecord{*batchTwo[1]})

				err := store.DeleteLinks(shorts)

				assert.Nil(t, err)

				for _, short := range shorts {
					for _, val := range short {
						assert.Equal(t, true, store.data[val.ShortLink].Deleted)
					}
				}
			case tests[1].name:
				shorts = append(shorts, []model.StoreRecord{*batchOne[0]}, []model.StoreRecord{*batchTwo[1], wrongRecord})

				err := store.DeleteLinks(shorts)

				assert.NotNil(t, err)

			}

		})
	}
}
