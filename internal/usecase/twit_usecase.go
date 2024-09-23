package usecase

import (
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres/repo"
	"github.com/google/uuid"
)

type twitService struct {
	ctxTimeout time.Duration
	repo       repo.TweetStorageI
}

func NewTwitService(timeout time.Duration, repository repo.TweetStorageI) Twit {
	return &twitService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}

func (a *twitService) beforeCreate(twit *entity.CreateTwitRequest) {
	twit.ID = uuid.New().String()
	twit.CreatedAt = time.Now().UTC()
	twit.UpdatedAt = time.Now().UTC()
}

func (a *twitService) beforeUpdate(twit *entity.CreateTwitRequest) {
	twit.UpdatedAt = time.Now().UTC()
}
