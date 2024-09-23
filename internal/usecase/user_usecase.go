package usecase

import (
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres/repo"
	"github.com/google/uuid"
)

type userService struct {
	ctxTimeout time.Duration
	repo       repo.UserStorageI
}

func NewUserService(timeout time.Duration, repository repo.UserStorageI) User {
	return &userService{
		ctxTimeout: timeout,
		repo:       repository,
	}
}

func (a *userService) beforeCreate(user *entity.CreateUserRequest) {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
}

func (a *userService) beforeUpdate(user *entity.CreateUserRequest) {
	user.UpdatedAt = time.Now().UTC()
}
