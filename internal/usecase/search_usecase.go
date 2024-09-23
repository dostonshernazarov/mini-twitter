package usecase

import (
	"time"

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
