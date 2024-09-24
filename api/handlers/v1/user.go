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
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/etc"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	access, refresh, err := h.JwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	createdUser, err := h.User.Create(ctx, entity.CreateUserRequest{
		ID:       uuid.NewString(),
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
		Password: hashed,
		Role:     "user",
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
		Username: *&createdUser.Username,
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
// @Param 			data body entity.UpdateUserRequest true "Update User Model"
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

	var request entity.UpdateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return
	}

	user, err := h.User.Get(ctx, map[string]interface{}{
		"id": request.ID,
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

	err = h.User.Update(ctx, request)
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
// @Success 		201 {object} entity.GetUserResponse
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

	userID := c.PostForm("id")

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return
	}

	uuidPath := uuid.NewString() + path.Ext(file.Filename)

	if err := c.SaveUploadedFile(file, fmt.Sprintf("./media/users/%s", uuidPath)); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	if err := h.User.UploadImage(ctx, userID, uuidPath); err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	user, err := h.User.Get(ctx, map[string]interface{}{
		"id": userID,
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

	c.JSON(http.StatusOK, user)
}
