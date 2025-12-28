package handler

import (
	"errors"
	"net/http"
	"time"
	"typing-speed/internals/core/typing"
	"typing-speed/internals/core/user"
	"typing-speed/pkg/logs"

	"github.com/gin-gonic/gin"
)

// respondError response structure for error
func (h *Handler) respondError(c *gin.Context, status int, msg string, err error, start time.Time, logsData *logs.LogEntry) {
	logsData.Level = LogLevelError
	logsData.Msg = msg
	logsData.Error = err
	logsData.Status = status
	logsData.Latency = logs.Duration(time.Since(start))

	h.logsChan <- *logsData

	c.JSON(status, gin.H{
		"error":  msg,
		"status": status,
	})
}

// respondSuccess response in case of success
func (h *Handler) respondSuccess(c *gin.Context,msg string,start time.Time,logsData *logs.LogEntry,data interface{}) {
	logsData.Level = LogLevelInfo
	logsData.Msg = msg
	logsData.Status = http.StatusOK
	logsData.Latency = logs.Duration(time.Since(start))

	h.logsChan <- *logsData

	c.JSON(http.StatusOK, gin.H{
		"message": msg,
		"status":  http.StatusOK,
		"data":    data,
	})
}

// recoverPanic response in case of panic
func (h *Handler) recoverPanic(c *gin.Context, start time.Time, logsData *logs.LogEntry) {
	if r := recover(); r != nil {
		logsData.Level = LogLevelError
		logsData.Msg = "panic recovered"
		logsData.ExtraData = r
		logsData.Status = http.StatusInternalServerError
		logsData.Latency = logs.Duration(time.Since(start))

		h.logsChan <- *logsData

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "internal server error",
			"status": http.StatusInternalServerError,
		})
	}
}

// handleServiceError handle error from service
func (h *Handler) handleServiceError(c *gin.Context,err error,logsData *logs.LogEntry,start time.Time) {
	//err := errStruct.Error

	status := http.StatusInternalServerError
	message := "something went wrong"

	switch {
	case errors.Is(err, typing.ErrInvalidUser):
		status = http.StatusBadRequest
		message = "user cannot be empty"

	case errors.Is(err, user.ErrUserNotFound):
		status = http.StatusNotFound
		message = "user not found"

	case errors.Is(err, user.ErrInvalidUserDetail):
		status = http.StatusBadRequest
		message = "invalid user details"

	case errors.Is(err, user.ErrInvalidRefreshToken):
		status = http.StatusForbidden
		message = "invalid refresh token"

	case errors.Is(err, user.ErrUserAlreadyRegistered):
		status = http.StatusBadRequest
		message = "user already registered"
	}

	logsData.Status = status
	logsData.Msg = message
	logsData.Latency = logs.Duration(time.Since(start))

	h.logsChan <- *logsData

	c.JSON(status, gin.H{
		"error":  message,
		"status": status,
		"data":   nil,
	})
}

