package postgres_test

import (
	"context"
	"testing"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	postgresql "github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/config"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	// Setup Config
	cfg := config.Load()

	// Setup Database
	db, err := postgres.New(cfg)
	assert.NoError(t, err)
	defer db.Close()

	// Setup Repository
	repo := postgresql.NewUserRepo(db)
	ctx := context.Background()
	defaultUserID := uuid.New().String()

	// Create
	createdUserModel := &entity.CreateUserRequest{
		ID:       defaultUserID,
		Name:     "James Bond",
		Username: "username_test",
		Email:    "jamesbond@gmail.com",
		Password: "12345678",
		Role:     "user",
	}

	createdUser, err := repo.Create(ctx, entity.CreateUserRequest{
		ID:       defaultUserID,
		Name:     "James Bond",
		Username: "username_test",
		Email:    "jamesbond@gmail.com",
		Password: "12345678",
		Role:     "user",
	})
	assert.NoError(t, err)
	assert.Equal(t, createdUserModel.ID, createdUser.ID)
	assert.Equal(t, createdUserModel.Name, createdUser.Name)
	assert.Equal(t, createdUserModel.Username, createdUser.Username)
	assert.Equal(t, createdUserModel.Email, createdUser.Email)
	assert.Equal(t, createdUser.Role, createdUser.Role)

	// Update
	updatedUserModel := &entity.UpdateUserRequest{
		ID:       defaultUserID,
		Name:     "New Test User",
		Username: "username_test",
		Bio:      new(string),
	}
	err = repo.Update(ctx, *updatedUserModel)
	assert.NoError(t, err)

	// Get
	getUser, err := repo.Get(ctx, map[string]interface{}{"id": defaultUserID})
	assert.NoError(t, err)
	assert.Equal(t, getUser.Name, updatedUserModel.Name)
	assert.Equal(t, getUser.ID, updatedUserModel.ID)
	assert.Equal(t, getUser.Username, updatedUserModel.Username)

	// Delete
	err = repo.Delete(ctx, defaultUserID)
	assert.NoError(t, err)
	notUser, err := repo.Get(ctx, map[string]interface{}{"id": defaultUserID})
	assert.Error(t, err)
	assert.Nil(t, notUser)
	assert.Equal(t, err, pgx.ErrNoRows)
}
