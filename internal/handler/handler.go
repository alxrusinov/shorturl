package handler

import "github.com/alxrusinov/shorturl/internal/model"

type options struct {
	responseAddr string
}

type Handler struct {
	store       Store
	options     *options
	Middlewares *Middlewares
	DeleteChan  chan []model.StoreRecord
}

type Store interface {
	GetLink(arg *model.StoreRecord) (*model.StoreRecord, error)
	SetLink(arg *model.StoreRecord) (*model.StoreRecord, error)
	SetBatchLink(arg []*model.StoreRecord) ([]*model.StoreRecord, error)
	Ping() error
	GetLinks(userID string) ([]model.StoreRecord, error)
	DeleteLinks(shorts [][]model.StoreRecord) error
}

func NewHandler(sStore Store, responseAddr string) *Handler {
	handler := &Handler{
		store: sStore,
		options: &options{
			responseAddr: responseAddr,
		},
		DeleteChan: make(chan []model.StoreRecord),
	}

	return handler
}
