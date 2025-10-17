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

func (h *TypingHandler) TypingDataHandler(c *gin.Context) {
	logsData := logs.LogEntry{}
	logsData.Method = c.Request.Method
	logsData.Path = c.FullPath()
	start := time.Now()
	fmt.Println("Typing Handler is hit")
	var userData *typing.TypingData
	err := c.ShouldBindJSON(&userData)
	if err != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = "ERROR"
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

	err = h.typingUseCase.AddUserData(context.Background(), userData)
	if err != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = "ERROR"
		h.handlerError(c, err,&logsData)
		return
	}
	logsData.Latency = logs.Duration(time.Since(start))
	logsData.Level = "SUCCESS"
	logsData.Msg = "user registered successfully"
	logsData.Status=http.StatusOK
	h.logsChan <- logsData
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": "user registered successfully",
		"status":  http.StatusOK,
		"data":    nil,
	})
	return

}
