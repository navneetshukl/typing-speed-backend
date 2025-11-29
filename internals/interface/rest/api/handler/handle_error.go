package handler

import (
	"errors"
	"net/http"
	"typing-speed/internals/core/typing"
	"typing-speed/internals/core/user"
	"typing-speed/pkg/logs"

	"github.com/gin-gonic/gin"
)

func (h *Handler) handlerError(c *gin.Context, er *user.ErrorStruct, logs *logs.LogEntry) {
	logs.Msg = er.ErrorMsg
	err:=er.Error

	// {403->login, 401->refresh token, 404->register user}
	switch {
	case errors.Is(err, typing.ErrInvalidUser):
		logs.Status = http.StatusBadRequest
		h.logsChan <- *logs
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "user cannot be empty",
			"status": http.StatusBadRequest,
			"data":   nil,
		})

	case errors.Is(err,user.ErrUserNotFound):
		logs.Status = http.StatusNotFound
		h.logsChan <- *logs
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "user not found.please register.",
			"status": http.StatusNotFound,
			"data":   nil,
		})

		case errors.Is(err,user.ErrInvalidUserDetail):
		logs.Status = http.StatusBadRequest
		h.logsChan <- *logs
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "user detail is not valid",
			"status": http.StatusBadRequest,
			"data":   nil,
		})

		case errors.Is(err,user.ErrInvalidRefreshToken):
		logs.Status = http.StatusForbidden
		h.logsChan <- *logs
		c.JSON(http.StatusForbidden, gin.H{
			"error":  "user signin again",
			"status": http.StatusForbidden,
			"data":   nil,
		})

		case errors.Is(err,user.ErrUserAlreadyRegistered):
		logs.Status = http.StatusBadRequest
		h.logsChan <- *logs
		c.JSON(http.StatusForbidden, gin.H{
			"error":  "user already registered",
			"status": http.StatusForbidden,
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
