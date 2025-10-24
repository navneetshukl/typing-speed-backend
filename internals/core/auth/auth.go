package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccessClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type AuthService interface {
	RegisterUser(ctx context.Context, user *User) error
	LoginUser(ctx context.Context, user *LoginUser) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}
