package v1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/spf13/cast"

	"github.com/dostonshernazarov/mini-twitter/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadTweetFiles
// @Security 		BearerAuth
// @Summary 		Upload Tweet Files
// @Description 	this api for creating a new tweet
// @Tags			tweet
// @Accept 			multipart/form-data
// @Produce 		json
// @Param 			files formData []file true "Tweet Files"
// @Success 		201 {object} []string
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/tweets/upload [POST]
func (h *HandlerV1) UploadTweetFiles(c *gin.Context) {
	duration, err := time.ParseDuration("15m")
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	_, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: "Error parsing form-data",
		})
		return
	}
	files := form.File["files"]

	var fileURLs []string

	for _, file := range files {
		uuidPath := uuid.NewString() + path.Ext(file.Filename)
		filePath := fmt.Sprintf("./media/tweets/%s", uuidPath)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.UploadingError,
			})
			log.Println(err.Error())
			return
		}

		fileURLs = append(fileURLs, uuidPath)
	}

	c.JSON(http.StatusOK, fileURLs)
}

// CreateTweet
// @Security 		BearerAuth
// @Summary 		Create Tweet
// @Description 	this api for creating a new tweet
// @Tags			tweet
// @Accept 			json
// @Produce 		json
// @Param 			request body entity.TweetRequest true "Create Tweet Model"
// @Success 		201 {object} entity.CreateTweetResponse
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/tweets [POST]
func (h *HandlerV1) CreateTweet(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.TimeOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var request entity.TweetRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return
	}

	claims, err := utils.GetClaimsFromToken(c.Request, h.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	UserId := cast.ToString(claims["sub"])

	// checking: post or repost
	if request.ParentTweetID != nil && request.Content != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		return
	} else if request.ParentTweetID == nil && request.Content == nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		return
	}

	response, err := h.Tweet.CreateTweet(ctx, entity.CreateTweetRequest{
		ID:            uuid.NewString(),
		UserID:        UserId,
		ParentTweetID: request.ParentTweetID,
		Content:       request.Content,
		URLs:          request.URLs,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, response)

}

// UpdateTweet
// @Security 		BearerAuth
// @Summary 		Update Tweet
// @Description 	this api for updating a tweet
// @Tags			tweet
// @Accept 			json
// @Produce 		json
// @Param 			request body entity.UpdateTweetRequest true "Update Tweet Model"
// @Success 		200 {object} entity.UpdateTweetResponse
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/tweets [PUT]
func (h *HandlerV1) UpdateTweet(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.TimeOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var request entity.UpdateTweetRequest

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return
	}

	claims, err := utils.GetClaimsFromToken(c.Request, h.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	userId := cast.ToString(claims["sub"])

	tweet, err := h.Tweet.GetTweet(ctx, request.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, entity.Error{
				Message: entity.NotFoundData,
			})
			log.Println(err.Error())
			return
		} else {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}
	}
	// only owner can update
	if tweet.UserID != userId {
		c.JSON(http.StatusUnauthorized, entity.Error{
			Message: entity.NoAccess,
		})
		return
	}

	// no update reposted tweet
	if tweet.ParentTweetID != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		return
	}

	response, err := h.Tweet.UpdateTweet(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteTweet
// @Security 		BearerAuth
// @Summary 		Delete Tweet
// @Description 	this api for deleting a tweet with id
// @Tags			tweet
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Tweet ID"
// @Success 		200 {object} entity.ResponseWithStatus
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/tweets/{id} [DELETE]
func (h *HandlerV1) DeleteTweet(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.TimeOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	id := c.Param("id")

	claims, err := utils.GetClaimsFromToken(c.Request, h.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	userId := cast.ToString(claims["sub"])

	tweet, err := h.Tweet.GetTweet(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, entity.Error{
				Message: entity.NotFoundData,
			})
			log.Println(err.Error())
			return
		} else {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}
	}

	// only owner can update
	if tweet.UserID != userId {
		c.JSON(http.StatusUnauthorized, entity.Error{
			Message: entity.NoAccess,
		})
		return
	}

	if err := h.Tweet.DeleteTweet(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, entity.Error{
				Message: entity.NotFoundData,
			})
			log.Println(err.Error())
			return
		} else {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, entity.ResponseWithStatus{
		Status: true,
	})
}

// GetTweet
// @Security 		BearerAuth
// @Summary 		Get Tweet
// @Description 	this api for getting a tweet with id
// @Tags			tweet
// @Accept 			json
// @Produce 		json
// @Param 			id path string true "Tweet ID"
// @Success 		200 {object} entity.GetTweetResponse
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/tweets/{id} [GET]
func (h *HandlerV1) GetTweet(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.TimeOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	id := c.Param("id")

	tweet, err := h.Tweet.GetTweet(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, entity.Error{
				Message: entity.NotFoundData,
			})
			log.Println(err.Error())
			return
		} else {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, tweet)
}

// GetTweets
// @Security 		BearerAuth
// @Summary 		List Tweet
// @Description 	this api for getting list of tweet
// @Tags			tweet
// @Accept 			json
// @Produce 		json
// @Param 			page query int false "Page"
// @Param 			limit query int false "Limit"
// @Success 		200 {object} entity.ListTweetsResponse
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/tweets [GET]
func (h *HandlerV1) ListTweets(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.TimeOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	params, errs := utils.ParseQueryParam(c.Request.URL.Query())
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(errs)
		return
	}

	tweets, err := h.Tweet.ListTweets(ctx, entity.Filter{
		Page:  int(params.Page),
		Limit: int(params.Limit),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, entity.Error{
				Message: entity.NotFoundData,
			})
			log.Println(err.Error())
			return
		} else {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, tweets)
}

// UserTweets
// @Security 		BearerAuth
// @Summary 		List User Tweet
// @Description 	this api for getting tweet list of user
// @Tags			tweet
// @Accept 			json
// @Produce 		json
// @Success 		200 {object} entity.ListTweetsResponse
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/tweets/users/{id} [GET]
func (h *HandlerV1) UserTweets(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.TimeOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	claims, err := utils.GetClaimsFromToken(c.Request, h.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	userId := cast.ToString(claims["sub"])

	tweets, err := h.Tweet.UserTweets(ctx, userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, entity.Error{
				Message: entity.NotFoundData,
			})
			log.Println(err.Error())
			return
		} else {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, tweets)
}
