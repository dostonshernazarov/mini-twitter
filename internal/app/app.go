package app

import (
	"log"
	"net/http"
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/redis"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/config"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/logger"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/postgres"
	"go.uber.org/zap"
)

type App struct {
	Config *config.Config
	Logger *zap.Logger
	DB     *postgres.PostgresDB
	server *http.Server
}

func NewApp(cfg config.Config) (*App, error) {
	// logger init
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	// postgres init
	db, err := postgres.New(&cfg)
	if err != nil {
		return nil, err
	}

	redisClient, err := cache.NewRedisStorage(&cfg)
	if err != nil {
		log.Fatalf("failed to create redis storage: %v", err)
	}

	cache.Init(redisClient)

	// context timeout initialization
	contextTimeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return nil, err
	}

	postgres.NewFollowRepo(db)

	return &App{
		Config: &config.Config{},
		Logger: logger,
		DB:     db,
		server: &http.Server{},
	}, nil
}
