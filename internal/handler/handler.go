package handler

import "github.com/alxrusinov/shorturl/internal/model"

type options struct {
	responseAddr string
}

// Handler - structure with information about handler entity
type Handler struct {
	store       Store
	options     *options
	Middlewares *Middlewares
	DeleteChan  chan []model.StoreRecord
}

// Store - interface of store
type Store interface {
	GetLink(arg *model.StoreRecord) (*model.StoreRecord, error)
	SetLink(arg *model.StoreRecord) (*model.StoreRecord, error)
	SetBatchLink(arg []*model.StoreRecord) ([]*model.StoreRecord, error)
	Ping() error
	GetLinks(userID string) ([]model.StoreRecord, error)
	DeleteLinks(shorts [][]model.StoreRecord) error
}

// NewHandler returns new handler instance
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
