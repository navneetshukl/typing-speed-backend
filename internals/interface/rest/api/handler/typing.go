package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"typing-speed/internals/core/typing"
	"typing-speed/internals/core/user"
	"typing-speed/pkg/logs"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	typingUseCase typing.TypingService
	userUseCase   user.UserService
	logsChan      chan logs.LogEntry
}

func NewHandler(ty typing.TypingService, auth user.UserService, ch chan logs.LogEntry) Handler {
	return Handler{
		typingUseCase: ty,
		logsChan:      ch,
		userUseCase:   auth,
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
	defer func() {
		if r := recover(); r != nil {

			logsData.Method = c.Request.Method
			logsData.Path = c.FullPath()
			logsData.ExtraData = r
			h.logsChan <- logsData
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "something went wrong",
				"status": http.StatusInternalServerError,
				"data":   nil,
			})
		}

	}()

	email := c.GetString("email")

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

	fmt.Println("Email in handler is ", email)

	er := h.typingUseCase.AddTestData(context.Background(), &userData, email)
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

func (h *Handler) RecentTestDashboardHandler(c *gin.Context) {
	logsData := logs.LogEntry{}
	logsData.Method = c.Request.Method
	logsData.Path = c.FullPath()
	defer func() {
		if r := recover(); r != nil {

			logsData.Method = c.Request.Method
			logsData.Path = c.FullPath()
			logsData.ExtraData = r
			h.logsChan <- logsData
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "something went wrong",
				"status": http.StatusInternalServerError,
				"data":   nil,
			})
		}

	}()
	start := time.Now()

	email := c.GetString("email")
	month:=c.Query("month")

	fmt.Println("Month is ",month)
	
	data, err := h.typingUseCase.RecentTestForProfile(context.Background(), email,month)
	if err != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		h.handlerError(c, err, &logsData)
		return
	}

	logsData.Latency = logs.Duration(time.Since(start))
	logsData.Level = LogLevelInfo
	logsData.Msg = "recent test fetched successfully"
	logsData.Status = http.StatusOK
	logsData.ResponseData=data
	h.logsChan <- logsData
	c.JSON(http.StatusOK, gin.H{
		"message": "recent test fetched successfully",
		"status":  http.StatusOK,
		"data":    data,
	})
}


