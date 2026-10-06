package domain

import "time"

const (
	AliasAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type URL struct {
	ID          int64
	OriginalURL string
	Alias       string
	CreatedAt   time.Time
	UpdateAt    time.Time
	DeleteAt    *time.Time
}
