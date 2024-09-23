package v1

import (
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/pkg/config"
	tokens "github.com/dostonshernazarov/mini-twitter/internal/pkg/token"
	"github.com/dostonshernazarov/mini-twitter/internal/usecase"
	"go.uber.org/zap"
)

type HandlerV1 struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	User           usecase.User
	Tweet          usecase.Twit
	Follow         usecase.Follow
	Search         usecase.Search
	Like           usecase.Like
}

type HandlerV1Config struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	User           usecase.User
	Tweet          usecase.Twit
	Follow         usecase.Follow
	Search         usecase.Search
	Like           usecase.Like
}

func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		Config:         c.Config,
		Logger:         c.Logger,
		ContextTimeout: c.ContextTimeout,
		JwtHandler:     c.JwtHandler,
		User:           c.User,
		Tweet:          c.Tweet,
		Follow:         c.Follow,
		Search:         c.Search,
		Like:           c.Like,
	}
}
