package typing

import (
	"context"
	"fmt"
	"typing-speed/internals/adapter/external/sendmail"
	"typing-speed/internals/adapter/port"
	"typing-speed/internals/core/auth"
	"typing-speed/internals/core/typing"
)

type TypingServiceImpl struct {
	userSvc port.UserRepository
	mailSvc sendmail.MailSender
	authSvc port.AuthRepository
}

func NewTypingService(svc port.UserRepository, mail sendmail.MailSender, auth port.AuthRepository) typing.TypingService {
	return &TypingServiceImpl{
		userSvc: svc,
		mailSvc: mail,
		authSvc: auth,
	}
}

func (t *TypingServiceImpl) AddUserData(ctx context.Context, data *typing.TypingData, email string) *auth.ErrorStruct {
	errorStruct := &auth.ErrorStruct{}

	// insert data into db
	data.Email = email
	err := t.userSvc.InsertUserData(ctx, data)
	if err != nil {
		errorStruct.Error = typing.ErrInsertingData
		errorStruct.ErrorMsg = fmt.Sprintf("failed to insert typing data: %v", err)
		return errorStruct
	}

	// update the total test of user to +1

	err = t.authSvc.UpdateTotalTest(ctx, email)
	if err != nil {
		errorStruct.Error = typing.ErrUpdatingTotalTest
		errorStruct.ErrorMsg = fmt.Sprintf("failed to updating test count: %v", err)
		return errorStruct
	}

	return nil
}
