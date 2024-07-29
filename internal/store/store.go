package store

type Store interface {
	GetLink(key string) (string, error)
	SetLink(key string, link string)
}
