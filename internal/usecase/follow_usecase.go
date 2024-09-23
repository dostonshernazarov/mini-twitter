package usecase

import (
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres/repo"
)

type followService struct {
	ctxTimeout time.Duration
	repo       repo.FollowStorageI
}

func NewFollowService(timeout time.Duration, repository repo.FollowStorageI) Follow {
	return &followService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}
