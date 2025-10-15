package routes

import (
	"typing-speed/internals/interface/rest/api/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(typingHandler handler.TypingHandler) *gin.Engine {
	app := gin.New()

	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Replace with your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Example routes
	api := app.Group("/api")
	api.POST("/typing",typingHandler.TypingDataHandler)

	return app
}
