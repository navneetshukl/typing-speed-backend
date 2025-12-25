package port

import (
	"context"
	"typing-speed/internals/core/typing"
)

type TypingRepository interface {
	InsertTestData(ctx context.Context, user *typing.TypingData) error
	GetRecentTestData(ctx context.Context, email string, month int) ([]*typing.TypingData, error)
}
