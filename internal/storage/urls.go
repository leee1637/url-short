package storage

import (
	"context"
	"errors"
	"fmt"
	"url-short/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// type URLRepository interface {
// 	Save(ctx context.Context, url *URL) error
// 	GetByAlias(ctx context.Context, alias string) (*URL, error)
// 	DeleteByAlias(ctx context.Context, alias string) error
// 	UpdateByAlias(ctx context.Context, oldAlias, newAlias string) error
// }

func (p *Postgres) Save(ctx context.Context, url *domain.URL) error {
	query := `INSERT INTO urls (original_url, alias) VALUES ($1, $2)`

	row, err := p.pool.Exec(ctx, query, url.OriginalURL, url.Alias)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrConflict
		}
		return fmt.Errorf("db: %w", err)
	}

	if row.RowsAffected() == 0 {
		return fmt.Errorf("save url: no rows inserted")
	}

	return nil
}

func (p *Postgres) GetByAlias(ctx context.Context, alias string) (*domain.URL, error) {
	query := `SELECT id, original_url, alias, created_at, update_at, delete_at FROM urls WHERE alias=$1`

	d := domain.URL{}

	err := p.pool.QueryRow(ctx, query, alias).Scan(&d.ID, &d.OriginalURL, &d.Alias, &d.CreatedAt, &d.UpdateAt, &d.DeleteAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("db: %w", err)
	}
	return &d, nil
}

func (p *Postgres) DeleteByAlias(ctx context.Context, alias string) error {
	query := `DELETE FROM urls WHERE alias=$1`

	row, err := p.pool.Exec(ctx, query, alias)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}
	if row.RowsAffected() == 0 {
		return fmt.Errorf("delete: %w", domain.ErrNotFound)
	}

	return nil
}

func (p *Postgres) UpdateURLByAlias(ctx context.Context, newURL string, alias string) error {
	query := `UPDATE urls SET original_url = $1, update_at = NOW() WHERE alias = $2`

	row, err := p.pool.Exec(ctx, query, newURL, alias)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}

	if row.RowsAffected() == 0 {
		return fmt.Errorf("update: %w", domain.ErrNotFound)
	}

	return nil
}
