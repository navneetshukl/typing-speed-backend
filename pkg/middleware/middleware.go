package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || len(authHeader) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format", "status": http.StatusUnauthorized, "data": nil})
			c.Abort()
			return
		}

		tokenStr := authHeader

		claims := &AccessClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token", "status": http.StatusUnauthorized, "data": nil})
			c.Abort()
			return
		}

		fmt.Println("Email is ",claims.Email)

		c.Set("email", claims.Email)
		c.Next()
	}
}
