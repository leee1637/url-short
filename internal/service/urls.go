package service

import (
	"context"
	"log/slog"
	"url-short/internal/domain"
)

// type URLRepository interface {
// 	Save(ctx context.Context, url *URL) error
// 	GetByAlias(ctx context.Context, alias string) (*URL, error)
// 	DeleteByAlias(ctx context.Context, alias string) error
// 	UpdateByAlias(ctx context.Context, oldAlias, newAlias string) error
// }

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

func (s *Service) SaveUrl(ctx context.Context, originalURL, alias string) error {

	err := AliasValidate(alias)
	if err != nil {
		s.logger.WarnContext(ctx, "Ошибка валидации алиаса",
			slog.Any("error", err),
			slog.String("alias", alias))

		return err
	}

	d := &domain.URL{
		OriginalURL: originalURL,
		Alias:       alias,
	}

	err = s.repo.Save(ctx, d)
	if err != nil {
		s.logger.ErrorContext(ctx, "Ошибка запроса к бд",
			slog.Any("error", err))
		return err
	}

	s.logger.Debug("Запрос на создания ссылки успешно обработан")
	return nil
}

func (s *Service) GetUrlByAlias(ctx context.Context, a string) (*domain.URL, error) {

	err := AliasValidate(a)
	if err != nil {
		s.logger.WarnContext(ctx, "Ошибка валидации алиаса",
			slog.Any("error", err),
			slog.String("alias", a))

		return nil, err
	}

	d, err := s.repo.GetByAlias(ctx, a)
	if err != nil {
		s.logger.ErrorContext(ctx, "Ошибка запроса к бд",
			slog.Any("error", err))
		return nil, err
	}

	s.logger.Debug("Запрос на получения данных ссылки успешно обработан")
	return d, nil
}

func (s *Service) DeleteUrlByAlias(ctx context.Context, a string) error {

	err := AliasValidate(a)
	if err != nil {
		s.logger.WarnContext(ctx, "Ошибка валидации алиаса",
			slog.Any("error", err),
			slog.String("alias", a))

		return err
	}

	err = s.repo.DeleteByAlias(ctx, a)
	if err != nil {
		s.logger.ErrorContext(ctx, "Ошибка запроса к бд",
			slog.Any("error", err))
		return err
	}

	s.logger.Debug("Запрос на удаление ссылки успешно обработан")
	return nil
}

func (s *Service) UpdateUrlByAlias(ctx context.Context, oldAlias, newURL string) error {

	err := AliasValidate(oldAlias)
	if err != nil {
		s.logger.WarnContext(ctx, "Ошибка валидации алиаса",
			slog.Any("error", err),
			slog.String("alias", oldAlias))

		return err
	}

	err = s.repo.UpdateURLByAlias(ctx, newURL, oldAlias)
	if err != nil {
		s.logger.ErrorContext(ctx, "Ошибка запроса к бд",
			slog.Any("error", err))
		return err
	}

	s.logger.Debug("Запрос на обновление ссылки успешно обработан")
	return nil
}
