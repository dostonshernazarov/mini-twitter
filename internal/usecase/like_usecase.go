package usecase

import (
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres/repo"
)

type likeService struct {
	ctxTimeout time.Duration
	repo       repo.LikeStorageI
}

func NewLikeService(timeout time.Duration, repository repo.LikeStorageI) Like {
	return &likeService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}
