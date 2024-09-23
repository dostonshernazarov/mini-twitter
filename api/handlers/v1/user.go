package v1

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/etc"
	"github.com/gin-gonic/gin"
)

// CreateUser
// @Security 		BearerAuth
// @Summary 		Create User
// @Description 	this api for creating a new user
// @Tags 			user
// @Accept 			json
// @Produce 		json
// @Param 			data body entity.CreateUserRequest true "Create User Model"
// @Success 		201 {object} entity.CreateUserResponse
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/users [POST]
func (h *HandlerV1) CreateUser(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.CtxTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var request entity.CreateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return
	}

	exists, err := h.storage.User().UniqueEmail(ctx, request.Email)
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

	exists, err = h.storage.User().UniqueUsername(ctx, request.Username)
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
	request.Password = hashed

	createdUser, err := h.storage.User().Create(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}
