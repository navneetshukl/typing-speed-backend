package typing

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"typing-speed/internals/adapter/external/sendmail"
	"typing-speed/internals/adapter/port"
	"typing-speed/internals/core/typing"
	"typing-speed/internals/core/user"
)

type TypingServiceImpl struct {
	userSvc port.UserRepository
	mailSvc sendmail.MailSender
	testSvc port.TestRepository
}

func NewTypingService(svc port.UserRepository, mail sendmail.MailSender, test port.TestRepository) typing.TypingService {
	return &TypingServiceImpl{
		userSvc: svc,
		mailSvc: mail,
		testSvc: test,
	}
}

func (t *TypingServiceImpl) AddTestData(ctx context.Context, data *typing.TypingData, email string) *user.ErrorStruct {
	errorStruct := &user.ErrorStruct{}

	// insert data into db
	data.Email = email
	err := t.testSvc.InsertTestData(ctx, data)
	if err != nil {
		errorStruct.Error = typing.ErrInsertingData
		errorStruct.ErrorMsg = fmt.Sprintf("failed to insert typing data: %v", err)
		return errorStruct
	}

	// update the total test of user to +1

	userData, err := t.userSvc.GetUserByEmail(ctx, email)
	if err != nil {
		errorStruct.Error = typing.ErrGettingDataFromDB
		errorStruct.ErrorMsg = fmt.Sprintf("failed to get user from db: %v", err)
		return errorStruct
	}
	var currentAccuracy int
	if data.TotalWords == 0 {
		currentAccuracy = 0
	} else {
		currentAccuracy = ((data.TypedWords - data.TotalErrors) * 100) / (data.TotalWords)
	}
	updatedAccuracy := (userData.AvgAccuracy*userData.TotalTest + currentAccuracy) / (userData.TotalTest + 1)
	bestSpeed := data.WPM
	if userData.BestSpeed > (bestSpeed) {
		bestSpeed = int(userData.BestSpeed)
	}
	updatedSpeed := (data.WPM + (userData.AvgSpeed * userData.TotalTest)) / (userData.TotalTest + 1)
	currentPerformance := (data.WPM * currentAccuracy)

	updatedPerformance := (userData.AvgPerformance*(userData.TotalTest) + currentPerformance) /
		(userData.TotalTest + 1)

	err = t.userSvc.UpdateUser(ctx, email, updatedSpeed, updatedAccuracy, updatedPerformance, bestSpeed)
	if err != nil {
		errorStruct.Error = typing.ErrUpdatingTotalTest
		errorStruct.ErrorMsg = fmt.Sprintf("failed to updating test count: %v", err)
		return errorStruct
	}

	return nil
}

func (t *TypingServiceImpl) RecentTestForProfile(ctx context.Context, email string, month string) ([]*typing.TypingData, *user.ErrorStruct) {
	errorStruct := &user.ErrorStruct{}
	var m int
	var err error
	if month != "" {
		m, err = strconv.Atoi(month)
		if err != nil {
			errorStruct.Error = typing.ErrSomethingWentWrong
			errorStruct.ErrorMsg = fmt.Sprintf("error converting string to int: %v", err)
			return nil, errorStruct
		}
	}
	m = -1
	data, err := t.testSvc.GetRecentTestData(ctx, email, m)
	if err != nil {
		errorStruct.Error = typing.ErrGettingDataFromDB
		errorStruct.ErrorMsg = fmt.Sprintf("failed to insert typing data: %v", err)
		return nil, errorStruct
	}
	return data, nil
}

func (t *TypingServiceImpl) SendTypingSentence(ctx context.Context) string {
	words := "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ@#$&"
	str := strings.Builder{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	len := len(words)

	for i := 1; i <= 150; i++ {
		idx := r.Intn(len)
		str.WriteByte(words[idx])

	}
	return str.String()

}
