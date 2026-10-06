package storage

import (
	"context"
	"fmt"
	"url-short/internal/domain"
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
		return fmt.Errorf("Ошибка добавления ссылки: %w", err)
	}

	if row.RowsAffected() == 0 {
		return fmt.Errorf("Было задествовано 0 сторк")
	}

	return nil
}

func (p *Postgres) GetByAlias(ctx context.Context, alias string) (*domain.URL, error) {
	query := `SELECT id, original_url, alias, created_at, updated_at, deleted_at FROM urls WHERE alias=$1`

	d := domain.URL{}

	err := p.pool.QueryRow(ctx, query, alias).Scan(&d.ID, &d.OriginalURL, &d.Alias, &d.CreatedAt, &d.UpdateAt, &d.DeleteAt)

	if err != nil {
		return nil, fmt.Errorf("Ошибка поулчения алиаса: %w", err)
	}

	return &d, nil
}

func (p *Postgres) DeleteByAlias(ctx context.Context, alias string) error {
	query := `DELETE FROM urls WHERE alias=$1`

	row, err := p.pool.Exec(ctx, query, alias)
	if err != nil {
		return fmt.Errorf("Ошибка удаления алиаса: %w", err)
	}

	if row.RowsAffected() == 0 {
		return fmt.Errorf("Было задествовано 0 сторк")
	}

	return nil
}

func (p *Postgres) UpdateURLByAlias(ctx context.Context, alias string, newURL string) error {
	query := `UPDATE urls SET original_url = $1 WHERE alias = $2`

	row, err := p.pool.Exec(ctx, query, newURL, alias)
	if err != nil {
		return fmt.Errorf("Ошибка обновления: %w", err)
	}

	if row.RowsAffected() == 0 {
		return fmt.Errorf("Было задествовано 0 сторк")
	}

	return nil
}
