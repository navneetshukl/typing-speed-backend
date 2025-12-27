package handler

import (
	"context"
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
	start := time.Now()

	logsData := &logs.LogEntry{
		Time:   start,
		Method: c.Request.Method,
		Path:   c.FullPath(),
	}

	defer h.recoverPanic(c, start, logsData)

	email := c.GetString("email")

	var userData typing.TypingData
	if err := c.ShouldBindJSON(&userData); err != nil {
		logsData.RequestData = userData
		h.respondError(c, http.StatusBadRequest, "invalid request body", err, start, logsData)
		return
	}

	err := h.typingUseCase.AddTestData(c.Request.Context(), &userData, email)
	if err != nil {
		h.handleServiceError(c, err, logsData, start)
		return
	}

	h.respondSuccess(c, "user typing data saved successfully", start, logsData, nil)
}

func (h *Handler) RecentTestDashboardHandler(c *gin.Context) {
	start := time.Now()

	logsData := &logs.LogEntry{
		Time:   start,
		Method: c.Request.Method,
		Path:   c.FullPath(),
	}

	defer h.recoverPanic(c, start, logsData)

	email := c.GetString("email")
	month := c.Query("month")

	data, err := h.typingUseCase.RecentTestForProfile(c.Request.Context(), email, month)
	if err != nil {
		h.handleServiceError(c, err, logsData, start)
		return
	}

	if data == nil {
		data = []*typing.TypingData{}
	}

	h.respondSuccess(c, "recent test fetched successfully", start, logsData, data)
}
func (h *Handler) SendWordsToType(c *gin.Context) {
	start := time.Now()

	logsData := &logs.LogEntry{
		Time:   start,
		Method: c.Request.Method,
		Path:   c.FullPath(),
	}

	defer h.recoverPanic(c, start, logsData)

	data := h.typingUseCase.SendTypingSentence(context.Background())

	h.respondSuccess(c, "typing data fetched successfully", start, logsData, data)
}
