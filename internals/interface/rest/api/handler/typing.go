package handler

import (
	"context"
	"fmt"
	"net/http"
	"typing-speed/internals/core/typing"

	"github.com/gin-gonic/gin"
)

type TypingHandler struct {
	typingUseCase typing.TypingService
}

func NewTypingHandler(ty typing.TypingService) TypingHandler {
	return TypingHandler{
		typingUseCase: ty,
	}
}

func (h *TypingHandler) TypingDataHandler(c *gin.Context) {
	fmt.Println("Typing Handler is hit")
	var userData *typing.TypingData
	err := c.ShouldBindJSON(&userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "something went wrong",
			"status": http.StatusInternalServerError,
			"data":   nil,
		})
		return
	}

	err = h.typingUseCase.AddUserData(context.Background(), userData)
	if err != nil {
		handlerError(c, err)
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": "user registered successfully",
		"status":  http.StatusOK,
		"data":    nil,
	})
	return

}
