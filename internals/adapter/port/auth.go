package port

import (
	"context"
	"typing-speed/internals/core/auth"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*auth.User, error)
	CreateUser(ctx context.Context, user *auth.User) error
	UpdateUser(ctx context.Context, email string, speed, accuracy int) error
	GetTopPerformer(ctx context.Context) ([]*auth.TopPerformer, error)
}
