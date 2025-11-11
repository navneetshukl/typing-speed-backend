package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID           int64      `db:"id" json:"id"`
	Name         string     `db:"name" json:"name"`
	Email        string     `db:"email" json:"email"`
	Password     string     `db:"password" json:"-"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	AvgSpeed     int        `db:"avg_speed" json:"avg_speed"`
	AvgAccuracy  int        `db:"avg_accuracy" json:"avg_accuracy"`
	TotalTest    int        `db:"total_test" json:"total_test"`
	Level        int        `db:"level" json:"level"`
	LastTestTime *time.Time `db:"last_test_time" json:"last_test_time,omitempty"`
	Streak       int        `db:"streak" json:"streak"`
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
	RegisterUser(ctx context.Context, user *User) *ErrorStruct
	LoginUser(ctx context.Context, user *LoginUser) (string, string, *ErrorStruct)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, *ErrorStruct)
}
