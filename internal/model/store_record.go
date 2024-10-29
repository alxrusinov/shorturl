package model

// StoreRecord has information about pair of original and shorten links and theri state
type StoreRecord struct {
	UUID          string `json:"user_id" db:"user_id"`
	ShortLink     string `json:"short_url,omitempty" db:"short"`
	OriginalLink  string `json:"original_url,omitempty" db:"original"`
	CorrelationID string `json:"correlation_id" db:"correlation_id"`
	Deleted       bool   `json:"is_deleted" db:"is_deleted"`
}
