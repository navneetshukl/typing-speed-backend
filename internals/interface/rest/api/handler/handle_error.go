package handler

import (
	"errors"
	"net/http"
	"typing-speed/internals/core/typing"
	"typing-speed/pkg/logs"

	"github.com/gin-gonic/gin"
)

func (h *TypingHandler) handlerError(c *gin.Context, err error, logs *logs.LogEntry) {
	logs.Msg = err.Error()

	switch {
	case errors.Is(err, typing.ErrInvalidUser):
		logs.Status = http.StatusBadRequest
		h.logsChan <- *logs
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "user cannot be empty",
			"status": http.StatusBadRequest,
			"data":   nil,
		})

	default:
		logs.Status = http.StatusInternalServerError
		h.logsChan <- *logs
		c.JSON(http.StatusInternalServerError, gin.H{

			"error":  "something went wrong",
			"status": http.StatusInternalServerError,
			"data":   nil,
		})
	}
}
