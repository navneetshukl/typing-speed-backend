package port

import (
	"context"
	"typing-speed/internals/core/user"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	CreateUser(ctx context.Context, user *user.User) error
	UpdateUser(ctx context.Context, email string, speed, accuracy int,performance int,bestSpeed int) error
	GetTopPerformer(ctx context.Context) ([]*user.TopPerformer, error)
	GetAllUser(ctx context.Context) ([]*user.User, error)
	GetDashboardTopData(ctx context.Context) (*user.DashboardTopData, error)
}
