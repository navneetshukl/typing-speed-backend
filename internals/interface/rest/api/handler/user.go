package handler

import (
	"context"
	"net/http"
	"time"
	"typing-speed/internals/core/user"
	"typing-speed/pkg/logs"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUser(c *gin.Context) {
	start := time.Now()

	logsData := &logs.LogEntry{
		Time:   start,
		Method: c.Request.Method,
		Path:   c.FullPath(),
	}

	defer h.recoverPanic(c, start, logsData)

	var userData user.User
	if err := c.ShouldBindJSON(&userData); err != nil {
		logsData.RequestData = userData
		h.respondError(c, http.StatusBadRequest, "invalid request body", err, start, logsData)
		return
	}

	logsData.RequestData = userData

	if err := h.userUseCase.RegisterUser(c.Request.Context(), &userData); err != nil {
		h.handleServiceError(c, err, logsData, start)
		return
	}

	h.respondSuccess(c, "user registered successfully", start, logsData, nil)
}

func (h *Handler) LoginUser(c *gin.Context) {
	start := time.Now()

	logsData := &logs.LogEntry{
		Time:   start,
		Method: c.Request.Method,
		Path:   c.FullPath(),
	}

	defer h.recoverPanic(c, start, logsData)

	var userData user.LoginUser
	if err := c.ShouldBindJSON(&userData); err != nil {
		logsData.RequestData = userData
		h.respondError(c, http.StatusBadRequest, "invalid request body", err, start, logsData)
		return
	}

	logsData.RequestData = userData

	loginData, err := h.userUseCase.LoginUser(c.Request.Context(), &userData)
	if err != nil {
		h.handleServiceError(c, err, logsData, start)
		return
	}

	// ✅ Use common success handler
	h.respondSuccess(c, "user login successful", start, logsData, loginData)

	// Set refresh token cookie AFTER response log is prepared
	c.SetCookie(
		"refresh_token",
		loginData.RefreshToken,
		int((24 * time.Hour).Seconds()),
		"/",
		"",
		false,
		true,
	)
}

func (h *Handler) RefreshHandlerV1(c *gin.Context) {
	start := time.Now()

	logsData := &logs.LogEntry{
		Time:   start,
		Method: c.Request.Method,
		Path:   c.FullPath(),
	}

	defer h.recoverPanic(c, start, logsData)
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		h.respondError(c, http.StatusUnauthorized, "refresh token not present", err, start, logsData)
		return
	}

	accToken, refToken, er := h.userUseCase.RefreshToken(context.Background(), cookie)
	if er != nil {
		h.handleServiceError(c, er, logsData, start)
		return

	}
	c.SetCookie("refresh_token", refToken, int(24*7*time.Hour), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message":      "refresh token generated successfully",
		"status":       http.StatusOK,
		"data":         nil,
		"access_token": accToken,
	})

}

func (h *Handler) UserByEmailHandler(c *gin.Context) {
	start := time.Now()

	logsData := &logs.LogEntry{
		Time:   start,
		Method: c.Request.Method,
		Path:   c.FullPath(),
	}

	defer h.recoverPanic(c, start, logsData)

	email := c.GetString("email")

	data, err := h.userUseCase.UserByEmail(c.Request.Context(), email)
	if err != nil {
		h.handleServiceError(c, err, logsData, start)
		return
	}

	if data == nil {
		data = &user.User{}
	}

	h.respondSuccess(c, "user data fetched successfully", start, logsData, data)
}
func (h *Handler) TopPerformerHandler(c *gin.Context) {
	start := time.Now()

	logsData := &logs.LogEntry{
		Time:   start,
		Method: c.Request.Method,
		Path:   c.FullPath(),
	}

	defer h.recoverPanic(c, start, logsData)

	data, err := h.userUseCase.TopPerformer(c.Request.Context())
	if err != nil {
		h.handleServiceError(c, err, logsData, start)
		return
	}

	if data == nil {
		data = []*user.TopPerformer{}
	}

	h.respondSuccess(c, "top performer fetched successfully", start, logsData, data)
}

func (h *Handler) DataForDashboardHandler(c *gin.Context) {
	start := time.Now()

	logsData := &logs.LogEntry{
		Time:   start,
		Method: c.Request.Method,
		Path:   c.FullPath(),
	}

	defer h.recoverPanic(c, start, logsData)

	data, err := h.userUseCase.GetDataForDashboard(c.Request.Context())
	if err != nil {
		h.handleServiceError(c, err, logsData, start)
		return
	}

	if data == nil {
		data = &user.DashboardData{}
	}

	h.respondSuccess(c, "user data fetched successfully", start, logsData, data)
}
