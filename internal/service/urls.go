package service

import (
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
