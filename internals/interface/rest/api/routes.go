package routes

import (
	"typing-speed/internals/interface/rest/api/handler"
	"typing-speed/pkg/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

const ACCESS_SECRET string = "access_secret_code"

func SetUpRoutes(handler handler.Handler) *gin.Engine {
	app := gin.New()

	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Replace with your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	pprof.Register(app)

	auth := app.Group("/auth")
	auth.POST("/signup", handler.RegisterUser)
	auth.POST("/signin", handler.LoginUser)

	api := app.Group("/api")
	api.Use(middleware.AuthMiddleware(ACCESS_SECRET))

	api.POST("/typing", handler.TypingDataHandler)

	return app
}
