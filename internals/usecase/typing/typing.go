package typing

import (
	"context"
	"fmt"
	"typing-speed/internals/adapter/port"
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

func (t *TypingServiceImpl) AddUserData(ctx context.Context, data *typing.TypingData) error {
	// err := typing.TypingDataValid(data)
	// if err != nil {
	// 	return err
	// }

	fmt.Println("Request is ",data)

	// insert this data to db
	err:=t.userSvc.InsertUserData(ctx,data)
	if err!=nil{
		return typing.ErrInsertingData
	}
	return nil
}
