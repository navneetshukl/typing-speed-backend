package routes

import (
	"net/http"
	"time"
	"typing-speed/internals/interface/rest/api/handler"
	"typing-speed/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(handler handler.Handler) *gin.Engine {
	app := gin.New()

	// ✅ Custom CORS middleware (credentials-safe)
	app.Use(cors.New(cors.Config{
		// Add your Vercel URL and localhost for testing
		AllowOrigins:     []string{"https://typing-speed-frontend.vercel.app", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Debug profiling
	pprof.Register(app)

	app.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "Health is Good",
		})
	})

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
	api.GET("/typingWord", handler.SendWordsToType)

	dashboard := protected.Group("/dashboard")
	dashboard.GET("/recentTest", handler.RecentTestDashboardHandler)

	return app
}
