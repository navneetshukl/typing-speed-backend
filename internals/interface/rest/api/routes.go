package routes

import (
	"typing-speed/internals/interface/rest/api/handler"
	"typing-speed/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(handler handler.Handler) *gin.Engine {
	app := gin.New()

	// ✅ Custom CORS middleware (credentials-safe)
	app.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

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
	api.GET("/typingWord", handler.SendWordsToType)

	dashboard := protected.Group("/dashboard")
	dashboard.GET("/recentTest", handler.RecentTestDashboardHandler)

	return app
}
