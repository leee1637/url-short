package domain

import "context"

type URLRepository interface {
	Save(ctx context.Context, url *URL) error
	GetByAlias(ctx context.Context, alias string) (*URL, error)
	DeleteByAlias(ctx context.Context, alias string) error
	UpdateURLByAlias(ctx context.Context, newURL, alias string) error
}
