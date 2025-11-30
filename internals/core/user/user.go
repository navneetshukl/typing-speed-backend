package user

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID             int64      `db:"id" json:"-"`
	Name           string     `db:"name" json:"name"`
	Email          string     `db:"email" json:"email"`
	Password       string     `db:"password" json:"password"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	AvgSpeed       int        `db:"avg_speed" json:"avgSpeed"`
	AvgAccuracy    int        `db:"avg_accuracy" json:"avgAccuracy"`
	TotalTest      int        `db:"total_test" json:"totalTest"`
	Level          int        `db:"level" json:"level"`
	LastTestTime   *time.Time `db:"last_test_time" json:"lastTestTime,omitempty"`
	Streak         int        `db:"streak" json:"streak"`
	BestSpeed      float64    `db:"best_speed" json:"bestSpeed"`
	AvgPerformance float64    `db:"avg_performance" json:"avgPerformance"`
}

type TopPerformer struct {
	Name        string `json:"name"`
	Performance int    `json:"performance"`
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

type DashboardTopData struct {
	TotalTest       int64 `json:"totalTest"`
	AverageSpeed    int   `json:"avgSpeed"`
	AverageAccuracy int   `json:"avgAccuracy"`
}

type DashboardData struct{
	User []*User  `json:"user"`
	DashboardTopData *DashboardTopData `json:"dashboardTopData"`
}

type UserService interface {
	RegisterUser(ctx context.Context, user *User) *ErrorStruct
	LoginUser(ctx context.Context, user *LoginUser) (string, string, *ErrorStruct)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, *ErrorStruct)
	UserByEmail(ctx context.Context, email string) (*User, *ErrorStruct)
	TopPerformer(ctx context.Context) ([]*TopPerformer, *ErrorStruct)
	GetDataForDashboard(ctx context.Context) (*DashboardData, *ErrorStruct)
}
