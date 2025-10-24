package handler

import (
	"context"
	"net/http"
	"time"
	"typing-speed/internals/core/auth"
	"typing-speed/pkg/logs"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUser(c *gin.Context) {
	logsData := logs.LogEntry{}
	logsData.Method = c.Request.Method
	logsData.Path = c.FullPath()

	start := time.Now()
	var userData auth.User
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

	err = h.authUseCase.RegisterUser(context.Background(), &userData)
	if err != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		h.handlerError(c, err, &logsData)
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
	return
}

func (h *Handler) LoginUser(c *gin.Context) {
	logsData := logs.LogEntry{}
	logsData.Method = c.Request.Method
	logsData.Path = c.FullPath()

	start := time.Now()
	var userData auth.LoginUser
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

	accToken, refToken, err := h.authUseCase.LoginUser(context.Background(), &userData)
	if err != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		h.handlerError(c, err, &logsData)
		return
	}
	logsData.Latency = logs.Duration(time.Since(start))
	logsData.Level = LogLevelInfo
	logsData.Msg = "user registered successfully"
	logsData.Status = http.StatusOK
	h.logsChan <- logsData
	c.SetCookie("refresh_token",refToken,int(24*7*time.Hour),"/","",false,true)
	c.JSON(http.StatusOK, gin.H{
		"message":      "user registered successfully",
		"status":       http.StatusOK,
		"data":         nil,
		"access_token": accToken,
	})
	return
}

func(h *Handler)RefreshHandler(c *gin.Context){
	logsData := logs.LogEntry{}
	logsData.Method = c.Request.Method
	logsData.Path = c.FullPath()

	start := time.Now()
	cookie,err:=c.Cookie("refresh_token")
	if err != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		logsData.Msg = "Refresh Token Not Present"
		logsData.Status = http.StatusUnauthorized
		h.logsChan <- logsData
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":  "refres token not present",
			"status": http.StatusUnauthorized,
			"data":   nil,
		})
		return
	}

	accToken, refToken, err := h.authUseCase.RefreshToken(context.Background(), cookie)
	if err != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		h.handlerError(c, err, &logsData)
		return
	}
	logsData.Latency = logs.Duration(time.Since(start))
	logsData.Level = LogLevelInfo
	logsData.Msg = "user registered successfully"
	logsData.Status = http.StatusOK
	h.logsChan <- logsData
	c.SetCookie("refresh_token",refToken,int(24*7*time.Hour),"/","",false,true)
	c.JSON(http.StatusOK, gin.H{
		"message":      "user registered successfully",
		"status":       http.StatusOK,
		"data":         nil,
		"access_token": accToken,
	})
	return
}
