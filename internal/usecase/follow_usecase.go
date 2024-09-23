package usecase

import (
	"context"
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
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

func (f *followService) Follow(ctx context.Context, follow entity.FollowAction) (bool, error) {
	return f.repo.Follow(ctx, follow)
}

func (f *followService) GetFollowings(ctx context.Context, id int) (entity.ListUser, error) {
	return f.repo.GetFollowings(ctx, id)
}

func (f *followService) GetFollowers(ctx context.Context, id int) (entity.ListUser, error) {
	return f.repo.GetFollowers(ctx, id)
}
