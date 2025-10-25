package handler

import (
	"context"
	"net/http"
	"time"
	"typing-speed/internals/core/auth"
	"typing-speed/internals/core/typing"
	"typing-speed/pkg/logs"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	typingUseCase typing.TypingService
	authUseCase   auth.AuthService
	logsChan      chan logs.LogEntry
}

func NewHandler(ty typing.TypingService,auth auth.AuthService, ch chan logs.LogEntry) Handler {
	return Handler{
		typingUseCase: ty,
		logsChan:      ch,
		authUseCase: auth,
	}
}

const (
	LogLevelInfo  = "INFO"
	LogLevelError = "ERROR"
)

func (h *Handler) TypingDataHandler(c *gin.Context) {
	logsData := logs.LogEntry{}
	logsData.Method = c.Request.Method
	logsData.Path = c.FullPath()

	start := time.Now()
	var userData typing.TypingData
	err := c.ShouldBindJSON(&userData)
	logsData.RequestData = userData
	if err != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		logsData.Msg = err.Error()
		logsData.Status = http.StatusInternalServerError
		h.logsChan <- logsData
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "something went wrong",
			"status": http.StatusInternalServerError,
			"data":   nil,
		})
		return
	}

	er := h.typingUseCase.AddUserData(context.Background(), &userData)
	if er != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		h.handlerError(c, er, &logsData)
		return
	}
	logsData.Latency = logs.Duration(time.Since(start))
	logsData.Level = LogLevelInfo
	logsData.Msg = "user registered successfully"
	logsData.Status = http.StatusOK
	h.logsChan <- logsData
	c.JSON(http.StatusOK, gin.H{
		"message": "user registered successfully",
		"status":  http.StatusOK,
		"data":    nil,
	})

}
