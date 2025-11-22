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

	userData, err := t.authSvc.GetUserByEmail(ctx, email)
	if err != nil {

	}
	calculatedAccuracy := (data.TotalErrors * 100) / (data.TotalWords)
	accuracy := userData.AvgAccuracy * userData.TotalTest
	speed := userData.AvgSpeed * userData.TotalTest
	updatedSpeed := (speed + data.WPM) / (userData.TotalTest + 1)
	updatedAccuracy := (calculatedAccuracy + accuracy) / (userData.TotalTest + 1)

	err = t.authSvc.UpdateUser(ctx, email, updatedSpeed, updatedAccuracy)
	if err != nil {
		errorStruct.Error = typing.ErrUpdatingTotalTest
		errorStruct.ErrorMsg = fmt.Sprintf("failed to updating test count: %v", err)
		return errorStruct
	}

	return nil
}

func (t *TypingServiceImpl) RecentTestForProfile(ctx context.Context, email string) ([]*typing.TypingData, *auth.ErrorStruct) {
	errorStruct := &auth.ErrorStruct{}
	data, err := t.userSvc.GetRecentTestForProfile(ctx, email)
	if err != nil {
		errorStruct.Error = typing.ErrGettingDataFromDB
		errorStruct.ErrorMsg = fmt.Sprintf("failed to insert typing data: %v", err)
		return nil, errorStruct
	}
	return data, nil
}
