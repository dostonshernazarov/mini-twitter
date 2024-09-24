package v1

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// Follow
// @Security 		BearerAuth
// @Summary 		Follow-Unfollow
// @Tags 			follow
// @Accept			json
// @Produce 		json
// @Param 			request body entity.FollowAction true "Follow Model"
// @Success 		200 {object} entity.ResponseWithStatus
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/follows [POST]
func (h *HandlerV1) FollowUnfollow(c *gin.Context) {
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

	var request entity.FollowAction

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

	request.UserID = cast.ToString(claims["sub"])

	status, err := h.Follow.Follow(ctx, request)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, entity.Error{
				Message: entity.IncorrectData,
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
		Status: status,
	})
}

// Followings
// @Security 		BearerAuth
// @Summary 		User Followings
// @Tags 			follow
// @Accept			json
// @Produce 		json
// @Success 		200 {object} entity.ListUser
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/followings [GET]
func (h *HandlerV1) Followings(c *gin.Context) {
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

	followings, err := h.Follow.GetFollowings(ctx, id)
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

	c.JSON(http.StatusOK, followings)
}

// Followers
// @Security 		BearerAuth
// @Summary 		User Followers
// @Tags 			follow
// @Accept			json
// @Produce 		json
// @Success 		200 {object} entity.ListUser
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/followers [GET]
func (h *HandlerV1) Followers(c *gin.Context) {
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

	followers, err := h.Follow.GetFollowers(ctx, id)
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

	c.JSON(http.StatusOK, followers)
}
