package api

import (
	"time"

	"github.com/casbin/casbin/v2"
	_ "github.com/dostonshernazarov/mini-twitter/api/docs"
	v1 "github.com/dostonshernazarov/mini-twitter/api/handlers/v1"
	"github.com/dostonshernazarov/mini-twitter/api/middleware"
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
	Enforcer       *casbin.Enforcer
	User           usecase.User
	Tweet          usecase.Twit
	Follow         usecase.Follow
	Search         usecase.Search
	Like           usecase.Like
}

// NewRoute
// @title Welcome To Mini Twitter API
// @Description API for Mini Twitter
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRoute(option RouteOption) *gin.Engine {

	gin.SetMode(option.Config.GinMode)
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CheckCasbinPermission(option.Enforcer, *option.Config))

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
		api.POST("/auth/sign-up", HandlerV1.SignUp)
		api.POST("/auth/verify", HandlerV1.VerifySignUp)
		api.POST("/auth/login", HandlerV1.LogIn)
		api.POST("/auth/forgot-password/:email", HandlerV1.ForgotPassword)
		api.POST("/auth/verify-forgot-password", HandlerV1.VerifyForgotPassword)
		api.PUT("/auth/reset-password", HandlerV1.ResetPassword)
		api.GET("/auth/refresh/:refresh", HandlerV1.GetNewToken)

		api.POST("/users", HandlerV1.CreateUser)
		api.PUT("/users", HandlerV1.UpdateUser)
		api.DELETE("/users/:id", HandlerV1.DeleteUser)
		api.GET("/users", HandlerV1.GetUser)
		api.GET("/users/list", HandlerV1.ListUsers)
		api.POST("/users/upload-photo", HandlerV1.UploadProfilePhoto)

		api.POST("/tweets/upload", HandlerV1.UploadTweetFiles)
		api.POST("/tweets", HandlerV1.CreateTweet)
		api.PUT("/tweets", HandlerV1.UpdateTweet)
		api.DELETE("/tweets/:id", HandlerV1.DeleteTweet)
		api.GET("/tweets/:id", HandlerV1.GetTweet)
		api.GET("/tweets", HandlerV1.ListTweets)
		api.GET("/tweets/users/:id", HandlerV1.UserTweets)

		api.GET("/search/:data", HandlerV1.SearchTweet)
		api.POST("/likes", HandlerV1.LikeTweet)

		api.POST("/follows", HandlerV1.FollowUnfollow)
		api.GET("/followings/:id", HandlerV1.Followings)
		api.GET("/followers/:id", HandlerV1.Followers)

	}

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
