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

// Like
// @Security 		BearerAuth
// @Summary 		Like-Unlike
// @Tags 			like
// @Accept			json
// @Produce 		json
// @Param 			request body entity.LikeAction true "Like Model"
// @Success 		200 {object} entity.ResponseWithStatus
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/likes [POST]
func (h *HandlerV1) LikeTweet(c *gin.Context) {
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

	var request entity.LikeAction

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

	status, err := h.Like.Like(ctx, request)
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
