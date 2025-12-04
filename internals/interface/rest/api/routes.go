package routes

import (
	"typing-speed/internals/interface/rest/api/handler"
	"typing-speed/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(handler handler.Handler) *gin.Engine {
	app := gin.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Debug profiling
	pprof.Register(app)

	auth := app.Group("/auth")
	auth.POST("/signup", handler.RegisterUser)
	auth.POST("/signin", handler.LoginUser)

	protected := app.Group("/")
	protected.Use(middleware.AuthMiddleware())

	api := protected.Group("/api")
	api.POST("/typing", handler.TypingDataHandler)
	api.GET("/userData", handler.UserByEmailHandler)
	api.GET("/topPerformer", handler.TopPerformerHandler)
	api.GET("/allUser", handler.DataForDashboardHandler)

	dashboard := protected.Group("/dashboard")
	dashboard.GET("/recentTest", handler.RecentTestDashboardHandler)

	return app
}
