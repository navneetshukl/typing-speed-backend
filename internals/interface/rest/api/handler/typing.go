package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"typing-speed/internals/core/typing"
	"typing-speed/pkg/logs"

	"github.com/gin-gonic/gin"
)

type TypingHandler struct {
	typingUseCase typing.TypingService
	logsChan      chan logs.LogEntry
}

func NewTypingHandler(ty typing.TypingService, ch chan logs.LogEntry) TypingHandler {
	return TypingHandler{
		typingUseCase: ty,
		logsChan:      ch,
	}
}

const (
	LogLevelInfo  = "INFO"
	LogLevelError = "ERROR"
)

func (h *TypingHandler) TypingDataHandler(c *gin.Context) {
	logsData := logs.LogEntry{}
	logsData.Method = c.Request.Method
	logsData.Path = c.FullPath()
	
	start := time.Now()
	fmt.Println("Typing Handler is hit")
	var userData typing.TypingData
	err := c.ShouldBindJSON(&userData)
	logsData.RequestData=userData
	if err != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		logsData.Msg = err.Error()
		logsData.Status=http.StatusInternalServerError
		h.logsChan <- logsData
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "something went wrong",
			"status": http.StatusInternalServerError,
			"data":   nil,
		})
		return
	}

	err = h.typingUseCase.AddUserData(context.Background(), &userData)
	if err != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		h.handlerError(c, err,&logsData)
		return
	}
	logsData.Latency = logs.Duration(time.Since(start))
	logsData.Level = LogLevelInfo
	logsData.Msg = "user registered successfully"
	logsData.Status=http.StatusOK
	h.logsChan <- logsData
	c.JSON(http.StatusOK, gin.H{
		"message": "user registered successfully",
		"status":  http.StatusOK,
		"data":    nil,
	})
	return

}
