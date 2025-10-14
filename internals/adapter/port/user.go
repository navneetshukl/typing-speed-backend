package port

import (
	"context"
	"typing-speed/internals/core/typing"
)

type UserRepository interface {
	InsertUserData(ctx context.Context, user *typing.TypingData) error
}