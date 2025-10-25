package typing

import (
	"context"
	"fmt"
	"typing-speed/internals/adapter/port"
	"typing-speed/internals/core/auth"
	"typing-speed/internals/core/typing"
)

type TypingServiceImpl struct {
	userSvc port.UserRepository
}

func NewTypingService(svc port.UserRepository) typing.TypingService {
	return &TypingServiceImpl{
		userSvc: svc,
	}
}

func (t *TypingServiceImpl) AddUserData(ctx context.Context, data *typing.TypingData) *auth.ErrorStruct {
	errorStruct := &auth.ErrorStruct{}

	// insert data into db
	err := t.userSvc.InsertUserData(ctx, data)
	if err != nil {
		errorStruct.Error = typing.ErrInsertingData
		errorStruct.ErrorMsg = fmt.Sprintf("failed to insert typing data: %v", err)
		return errorStruct
	}

	return nil
}

