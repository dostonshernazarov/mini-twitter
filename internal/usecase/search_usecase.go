package usecase

import (
	"context"
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres/repo"
)

type searchService struct {
	ctxTimeout time.Duration
	repo       repo.SearchStorageI
}

func NewSearchService(timeout time.Duration, repository repo.SearchStorageI) Search {
	return &searchService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}

func (s *searchService) Search(ctx context.Context, data string) (entity.SearchResponse, error) {
	return s.repo.Search(ctx, data)
}
