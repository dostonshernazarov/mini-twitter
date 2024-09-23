package usecase

import (
	"context"
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres/repo"
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

func (u *userService) UniqueUsername(ctx context.Context, username string) (bool, error) {
	return u.repo.UniqueUsername(ctx, username)
}

func (u *userService) UniqueEmail(ctx context.Context, email string) (bool, error) {
	return u.repo.UniqueEmail(ctx, email)
}

func (u *userService) Create(ctx context.Context, user entity.CreateUserRequest) (entity.CreateUserResponse, error) {

	return u.repo.Create(ctx, user)
}

func (u *userService) Update(ctx context.Context, user entity.UpdateUserRequest) (entity.UpdateUserResponse, error) {
	return u.repo.Update(ctx, user)
}

func (u *userService) UpdatePasswd(ctx context.Context, id int, passwd string) error {
	return u.repo.UpdatePasswd(ctx, id, passwd)
}

func (u *userService) UploadImage(ctx context.Context, id int, url string) error {
	return u.repo.UploadImage(ctx, id, url)
}

func (u *userService) Delete(ctx context.Context, id int) error {
	return u.repo.Delete(ctx, id)
}

func (u *userService) Get(ctx context.Context, field map[string]interface{}) (entity.GetUserResponse, error) {
	return u.repo.Get(ctx, field)
}

func (u *userService) List(ctx context.Context, filter entity.Filter) (entity.ListUser, error) {
	return u.repo.List(ctx, filter)
}
