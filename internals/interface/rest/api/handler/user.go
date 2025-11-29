package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"typing-speed/internals/core/user"
	"typing-speed/pkg/logs"
	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUser(c *gin.Context) {
	logsData := logs.LogEntry{}
	defer func() {
		if r := recover(); r != nil {
			h.logsChan <- logsData
			logsData.Method = c.Request.Method
			logsData.Path = c.FullPath()
			logsData.ExtraData = r
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "something went wrong",
				"status": http.StatusInternalServerError,
				"data":   nil,
			})
		}

	}()

	start := time.Now()
	var userData user.User
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

	er := h.userUseCase.RegisterUser(context.Background(), &userData)
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

func (h *Handler) LoginUser(c *gin.Context) {
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
	var userData user.LoginUser
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

	accToken, refToken, er := h.userUseCase.LoginUser(context.Background(), &userData)
	if er != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		h.handlerError(c, er, &logsData)
		return
	}
	logsData.Latency = logs.Duration(time.Since(start))
	logsData.Level = LogLevelInfo
	logsData.Msg = "user login successfully"
	logsData.Status = http.StatusOK
	h.logsChan <- logsData
	c.SetCookie("refresh_token", refToken, int(24*7*time.Hour), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message":      "user login successfully",
		"status":       http.StatusOK,
		"data":         nil,
		"access_token": accToken,
	})

}

func (h *Handler) RefreshHandler(c *gin.Context) {
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
	cookie, err := c.Cookie("refresh_token")
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

	accToken, refToken, er := h.userUseCase.RefreshToken(context.Background(), cookie)
	if er != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		h.handlerError(c, er, &logsData)
		return
	}
	logsData.Latency = logs.Duration(time.Since(start))
	logsData.Level = LogLevelInfo
	logsData.Msg = "refresh token generated successfully"
	logsData.Status = http.StatusOK
	h.logsChan <- logsData
	c.SetCookie("refresh_token", refToken, int(24*7*time.Hour), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message":      "refresh token generated successfully",
		"status":       http.StatusOK,
		"data":         nil,
		"access_token": accToken,
	})

}

func (h *Handler) UserByEmailHandler(c *gin.Context) {
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
	email = "a@a.com"
	fmt.Println("Email is ", email)

	data, err := h.userUseCase.UserByEmail(context.Background(), email)
	if err != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		h.handlerError(c, err, &logsData)
		return
	}

	logsData.Latency = logs.Duration(time.Since(start))
	logsData.Level = LogLevelInfo
	logsData.Msg = "user data fetched successfully"
	logsData.Status = http.StatusOK
	logsData.ResponseData = data
	h.logsChan <- logsData
	c.JSON(http.StatusOK, gin.H{
		"message": "user data fetched successfully",
		"status":  http.StatusOK,
		"data":    data,
	})
}

func (h *Handler) TopPerformerHandler(c *gin.Context) {
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

	data, err := h.userUseCase.TopPerformer(context.Background())
	if err != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		h.handlerError(c, err, &logsData)
		return
	}

	fmt.Println("Top performer is ", data)

	logsData.Latency = logs.Duration(time.Since(start))
	logsData.Level = LogLevelInfo
	logsData.Msg = "top performer fetched successfully"
	logsData.Status = http.StatusOK
	logsData.ResponseData = data
	h.logsChan <- logsData
	c.JSON(http.StatusOK, gin.H{
		"message": "top performer fetched successfully",
		"status":  http.StatusOK,
		"data":    data,
	})
}


func (h *Handler) AllUserHandler(c *gin.Context) {
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

	// email := c.GetString("email")
	// email = "a@a.com"
	// fmt.Println("Email is ", email)

	data, err := h.userUseCase.GetAllUser(context.Background())
	if err != nil {
		logsData.Latency = logs.Duration(time.Since(start))
		logsData.Level = LogLevelError
		h.handlerError(c, err, &logsData)
		return
	}

	logsData.Latency = logs.Duration(time.Since(start))
	logsData.Level = LogLevelInfo
	logsData.Msg = "user data fetched successfully"
	logsData.Status = http.StatusOK
	logsData.ResponseData = data
	h.logsChan <- logsData
	c.JSON(http.StatusOK, gin.H{
		"message": "user data fetched successfully",
		"status":  http.StatusOK,
		"data":    data,
	})
}
