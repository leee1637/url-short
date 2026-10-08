package service

import (
	"context"
	"log/slog"
	"url-short/internal/domain"
)

type Service struct {
	repo   domain.URLRepository
	logger *slog.Logger
}

func New(repo domain.URLRepository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// logRepoErr штатные ошибки (404/409/400) — Warn, остальное — Error.
func (s *Service) logRepoErr(ctx context.Context, op string, err error) {
	if isExpected(err) {
		s.logger.WarnContext(ctx, op, slog.Any("error", err))
		return
	}
	s.logger.ErrorContext(ctx, op, slog.Any("error", err))
}

func (s *Service) SaveUrl(ctx context.Context, originalURL, alias string) error {
	if err := URLValidate(originalURL); err != nil {
		s.logger.WarnContext(ctx, "валидация url",
			slog.Any("error", err),
			slog.String("url", originalURL))
		return err
	}

	if alias == "" {
		alias = GenerateAlias()
	}

	if err := AliasValidate(alias); err != nil {
		s.logger.WarnContext(ctx, "валидация алиаса",
			slog.Any("error", err),
			slog.String("alias", alias))
		return err
	}

	d := &domain.URL{
		OriginalURL: originalURL,
		Alias:       alias,
	}

	if err := s.repo.Save(ctx, d); err != nil {
		s.logRepoErr(ctx, "ошибка сохранения ссылки", err)
		return err
	}

	s.logger.InfoContext(ctx, "ссылка создана", slog.String("alias", alias))
	return nil
}

func (s *Service) GetUrlByAlias(ctx context.Context, a string) (*domain.URL, error) {
	if err := AliasValidate(a); err != nil {
		s.logger.WarnContext(ctx, "валидация алиаса",
			slog.Any("error", err),
			slog.String("alias", a))
		return nil, err
	}

	d, err := s.repo.GetByAlias(ctx, a)
	if err != nil {
		s.logRepoErr(ctx, "ошибка получения ссылки", err)
		return nil, err
	}

	s.logger.DebugContext(ctx, "ссылка получена", slog.String("alias", a))
	return d, nil
}

func (s *Service) DeleteUrlByAlias(ctx context.Context, a string) error {
	if err := AliasValidate(a); err != nil {
		s.logger.WarnContext(ctx, "валидация алиаса",
			slog.Any("error", err),
			slog.String("alias", a))
		return err
	}

	if err := s.repo.DeleteByAlias(ctx, a); err != nil {
		s.logRepoErr(ctx, "ошибка удаления ссылки", err)
		return err
	}

	s.logger.InfoContext(ctx, "ссылка удалена", slog.String("alias", a))
	return nil
}

func (s *Service) UpdateUrlByAlias(ctx context.Context, alias, newURL string) error {
	if err := AliasValidate(alias); err != nil {
		s.logger.WarnContext(ctx, "валидация алиаса",
			slog.Any("error", err),
			slog.String("alias", alias))
		return err
	}

	if err := URLValidate(newURL); err != nil {
		s.logger.WarnContext(ctx, "валидация url",
			slog.Any("error", err),
			slog.String("url", newURL))
		return err
	}

	if err := s.repo.UpdateURLByAlias(ctx, newURL, alias); err != nil {
		s.logRepoErr(ctx, "ошибка обновления ссылки", err)
		return err
	}

	s.logger.InfoContext(ctx, "ссылка обновлена", slog.String("alias", alias))
	return nil
}
