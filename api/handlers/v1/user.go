package v1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	awss3 "github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/awsS3"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/etc"
	tokens "github.com/dostonshernazarov/mini-twitter/internal/pkg/token"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

// CreateUser
// @Security 		BearerAuth
// @Summary 		Create User
// @Description 	this api for creating a new user
// @Tags 			user
// @Accept 			json
// @Produce 		json
// @Param 			data body entity.User true "Create User Model"
// @Success 		201 {object} entity.UserResponse
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/users [POST]
func (h *HandlerV1) CreateUser(c *gin.Context) {
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

	var request entity.User

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return
	}

	exists, err := h.User.UniqueEmail(ctx, request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
	}

	if exists {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.EmailUsed,
		})
		return
	}

	exists, err = h.User.UniqueUsername(ctx, request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.UsernameTaken,
		})
		return
	}

	hashed, err := etc.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	userID := uuid.NewString()

	h.JwtHandler = tokens.JwtHandler{
		Sub:       userID,
		Role:      entity.RoleUser,
		SigninKey: h.Config.SigningKey,
	}

	access, refresh, err := h.JwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	createdUser, err := h.User.Create(ctx, entity.CreateUserRequest{
		ID:       userID,
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
		Password: hashed,
		Role:     entity.RoleUser,
		Refresh:  refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, entity.UserResponse{
		ID:       createdUser.ID,
		Name:     createdUser.Name,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Role:     createdUser.Role,
		Access:   access,
	})
}

// UpdateUser
// @Security 		BearerAuth
// @Summary 		Update User
// @Description 	this api for updating a user
// @Tags 			user
// @Accept 			json
// @Produce 		json
// @Param 			data body entity.UpdateUserRequestSwag true "Update User Model"
// @Success 		201 {object} entity.ResponseWithStatus
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/users [PUT]
func (h *HandlerV1) UpdateUser(c *gin.Context) {
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

	var request entity.UpdateUserRequestSwag

	if err := c.ShouldBindJSON(&request); err != nil {
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

	id := cast.ToString(claims["sub"])

	user, err := h.User.Get(ctx, map[string]interface{}{
		"id": id,
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

	if request.Username != user.Username {
		exists, err := h.User.UniqueUsername(ctx, request.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}

		if exists {
			c.JSON(http.StatusBadRequest, entity.Error{
				Message: entity.UsernameTaken,
			})
			return
		}
	}

	err = h.User.Update(ctx, entity.UpdateUserRequest{
		ID:       id,
		Name:     request.Name,
		Username: request.Username,
		Bio:      request.Bio,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.ResponseWithStatus{
		Status: true,
	})
}

// DeleteUser
// @Security 		BearerAuth
// @Summary 		Delete User
// @Description 	this api for deleting a user
// @Tags 			user
// @Accept 			json
// @Produce 		json
// @Param 			id path int true "User ID"
// @Success 		200 {object} entity.ResponseWithStatus
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/users/{id} [DELETE]
func (h *HandlerV1) DeleteUser(c *gin.Context) {
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

	userID := c.Param("id")

	err = h.User.Delete(ctx, userID)
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

	c.JSON(http.StatusOK, entity.ResponseWithStatus{
		Status: true,
	})
}

// GetUser
// @Security 		BearerAuth
// @Summary 		Get User
// @Description 	this api for getting a user
// @Tags 			user
// @Accept 			json
// @Produce 		json
// @Param 			key query string true "Key"
// @Param 			value query string true "Value"
// @Success 		200 {object} entity.GetUserResponse
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/users [GET]
func (h *HandlerV1) GetUser(c *gin.Context) {
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

	key := c.Query("key")
	value := c.Query("value")

	if key != "id" && key != "email" && key != "username" {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(key)
		return
	}

	user, err := h.User.Get(ctx, map[string]interface{}{
		key: value,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{
			Message: entity.NotFoundData,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListUsers
// @Security 		BearerAuth
// @Summary 		List User
// @Description 	this api for getting list of user
// @Tags 			user
// @Accept 			json
// @Produce 		json
// @Param 			page query int false "Page"
// @Param 			limit query int false "Limit"
// @Success 		201 {object} entity.ListUser
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/users/list [GET]
func (h *HandlerV1) ListUsers(c *gin.Context) {
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

	users, err := h.User.List(ctx, entity.Filter{
		Page:  int(params.Page),
		Limit: int(params.Limit),
	})

	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{
			Message: entity.NotFoundData,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

// UploadProfilePhoto
// @Security 		BearerAuth
// @Summary 		Upload User Profile
// @Description 	this api for uploading user profile
// @Tags 			user
// @Accept 			json
// @Produce 		json
// @Param 			id formData int true "User ID"
// @Param 			avatar formData file true "User Profile Photo"
// @Success 		201 {object} string
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/users/upload-photo [POST]
func (h *HandlerV1) UploadProfilePhoto(c *gin.Context) {
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

	id := cast.ToString(claims["sub"])

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return

	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	defer src.Close()

	url, uniqueFile, err := awss3.UploadFileToS3(h.Config, src, file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	if err := c.SaveUploadedFile(file, fmt.Sprintf("./media/users/%s", uniqueFile)); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	if err := h.User.UploadImage(ctx, id, url); err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": url,
	})
}

// GetUser
// @Security 		BearerAuth
// @Summary 		Get User profile
// @Description 	this api for getting a user data by access token
// @Tags 			user
// @Accept 			json
// @Produce 		json
// @Success 		200 {object} entity.GetUserResponse
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/users/profile [GET]
func (h *HandlerV1) GetUserProfile(c *gin.Context) {
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

	id := cast.ToString(claims["sub"])

	user, err := h.User.Get(ctx, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		c.JSON(http.StatusNotFound, entity.Error{
			Message: entity.NotFoundData,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}
