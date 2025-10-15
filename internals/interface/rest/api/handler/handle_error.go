package handler

import (
	"errors"
	"net/http"
	"typing-speed/internals/core/typing"

	"github.com/gin-gonic/gin"
)

func handlerError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, typing.ErrInvalidUser):
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "user cannot be empty",
			"status": http.StatusBadRequest,
			"data":   nil,
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "something went wrong",
			"status": http.StatusInternalServerError,
			"data":   nil,
		})
	}
}
