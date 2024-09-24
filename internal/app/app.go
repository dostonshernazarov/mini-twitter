package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	"github.com/dostonshernazarov/mini-twitter/api"
	"github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/postgres"
	cache "github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/redis"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/config"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/logger"
	postgresdb "github.com/dostonshernazarov/mini-twitter/internal/pkg/postgres"

	"github.com/dostonshernazarov/mini-twitter/internal/usecase"
	"go.uber.org/zap"
)

type App struct {
	Config   *config.Config
	Logger   *zap.Logger
	DB       *postgresdb.PostgresDB
	server   *http.Server
	Enforcer *casbin.Enforcer
	User     usecase.User
	Tweet    usecase.Twit
	Follow   usecase.Follow
	Search   usecase.Search
	Like     usecase.Like
}

func NewApp(cfg config.Config) (*App, error) {
	// logger init
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	// postgres init
	db, err := postgresdb.New(&cfg)
	if err != nil {
		return nil, err
	}

	redisClient, err := cache.NewRedisStorage(&cfg)
	if err != nil {
		log.Fatalf("failed to create redis storage: %v", err)
	}

	cache.Init(redisClient)

	// context timeout initialization
	contextTimeout, err := time.ParseDuration(cfg.Context.TimeOut)
	if err != nil {
		return nil, err
	}

	// initialization enforcer
	enforcer, err := casbin.NewEnforcer("./internal/pkg/config/auth.conf", "./internal/pkg/config/auth.csv")
	if err != nil {
		return nil, err
	}

	// Storage init
	userRepo := postgres.NewUserRepo(db)
	tweetRepo := postgres.NewTweetRepo(db)
	followRepo := postgres.NewFollowRepo(db)
	likeRepo := postgres.NewLikeRepo(db)
	searchRepo := postgres.NewSearchRepo(db)

	//Usecase init
	usecase.NewUserService(contextTimeout, userRepo)
	usecase.NewTweetService(contextTimeout, tweetRepo)
	usecase.NewFollowService(contextTimeout, followRepo)
	usecase.NewLikeService(contextTimeout, likeRepo)
	usecase.NewSearchService(contextTimeout, searchRepo)

	return &App{
		Config:   &cfg,
		Logger:   logger,
		DB:       db,
		Enforcer: enforcer,
		User:     userRepo,
		Tweet:    tweetRepo,
		Follow:   followRepo,
		Search:   searchRepo,
		Like:     likeRepo,
	}, nil
}

func (a *App) Run() error {
	contextTimeout, err := time.ParseDuration(a.Config.Context.TimeOut)
	if err != nil {
		return fmt.Errorf("error while parsing context timeout test: %v", err)
	}

	// api init
	handler := api.NewRoute(api.RouteOption{
		Config:         a.Config,
		Logger:         a.Logger,
		ContextTimeout: contextTimeout,
		Enforcer:       a.Enforcer,
		User:           a.User,
		Tweet:          a.Tweet,
		Follow:         a.Follow,
		Search:         a.Search,
		Like:           a.Like,
	})

	err = a.Enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	if roleManager, ok := a.Enforcer.GetRoleManager().(*defaultrolemanager.RoleManager); ok {
		// Use the roleManager as expected
		roleManager.AddMatchingFunc("keyMatch", util.KeyMatch)
		roleManager.AddMatchingFunc("keyMatch3", util.KeyMatch3)
	} else {
		return fmt.Errorf("unexpected RoleManager type: %T", a.Enforcer.GetRoleManager())
	}

	// server init
	a.server, err = api.NewServer(a.Config, handler)
	if err != nil {
		return fmt.Errorf("error while initializing server: %v", err)
	}

	return a.server.ListenAndServe()
}

func (a *App) Stop() {

	// close database
	a.DB.Close()

	// shutdown server http
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.Logger.Error("shutdown server http ", zap.Error(err))
	}

	// zap logger sync
	a.Logger.Sync()
}
