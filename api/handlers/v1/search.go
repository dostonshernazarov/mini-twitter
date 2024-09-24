package v1

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"

	"github.com/gin-gonic/gin"
)

// Search
// @Security 		BearerAuth
// @Summary 		Search
// @Tags 			search
// @Accept			json
// @Produce 		json
// @Param 			data path string true "Search Content"
// @Success 		200 {object} entity.ListUser
// @Failure 		400 {object} entity.Error
// @Failure 		401 {object} entity.Error
// @Failure 		403 {object} entity.Error
// @Failure 		404 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/search/{data} [GET]
func (h *HandlerV1) SearchTweet(c *gin.Context) {
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

	data := c.Param("data")

	response, err := h.Search.Search(ctx, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
