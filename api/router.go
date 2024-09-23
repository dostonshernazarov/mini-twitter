package api

import (
	"time"

	v1 "github.com/dostonshernazarov/mini-twitter/api/handlers/v1"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/config"
	"github.com/dostonshernazarov/mini-twitter/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	tokens "github.com/dostonshernazarov/mini-twitter/internal/pkg/token"

	"go.uber.org/zap"
)

type RouteOption struct {
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

// NewRoute
// @title Welcome To Farmish API
// @Description API for Farmer
func NewRoute(option RouteOption) *gin.Engine {

	gin.SetMode(option.Config.GinMode)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	HandlerV1 := v1.New(&v1.HandlerV1Config{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		JwtHandler:     option.JwtHandler,
		User:           option.User,
		Tweet:          option.Tweet,
		Follow:         option.Follow,
		Search:         option.Search,
		Like:           option.Like,
	})

	api := router.Group("/v1")

	{
		api.POST("/users", HandlerV1.CreateUser)
	}

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router
}
