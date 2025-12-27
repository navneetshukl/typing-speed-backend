package user

import (
	"context"
	"fmt"
	"log"
	"typing-speed/internals/adapter/external/sendmail"
	"typing-speed/internals/adapter/port"
	"typing-speed/internals/core/user"

	"github.com/golang-jwt/jwt/v5"
)

type UserServiceImpl struct {
	userSvc port.UserRepository
	mailSvc sendmail.MailSender
}

func NewUserService(svc port.UserRepository, mail sendmail.MailSender) user.UserService {
	return &UserServiceImpl{
		userSvc: svc,
		mailSvc: mail,
	}
}

// RegisterUser handles user registration
func (a *UserServiceImpl) RegisterUser(ctx context.Context, userData *user.User) *user.ErrorStruct {
	errorStruct := &user.ErrorStruct{}
	log.Println("UserData is ", userData.Name, " ", userData.Email, " ", userData.Password)

	if userData.Email == "" || userData.Name == "" || userData.Password == "" {
		errorStruct.Error = user.ErrInvalidUserDetail
		errorStruct.ErrorMsg = "user detail is not complete"
		return errorStruct
	}

	data, err := a.userSvc.GetUserByEmail(ctx, userData.Email)
	if err != nil {
		errorStruct.Error = user.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to fetch user by email: " + err.Error()
		return errorStruct
	}
	if data != nil {
		errorStruct.Error = user.ErrUserAlreadyRegistered
		errorStruct.ErrorMsg = "user already registered with this email"
		return errorStruct
	}

	hash, err := user.HashPassword(userData.Password)
	if err != nil {
		errorStruct.Error = user.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to hash password: " + err.Error()
		return errorStruct
	}
	userData.Password = hash

	err = a.userSvc.CreateUser(ctx, userData)
	if err != nil {
		errorStruct.Error = user.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to create user in database: " + err.Error()
		return errorStruct
	}

	//err = a.mailSvc.SendMail("typing@gmail.com", userData.Email, "Register", "User registered successfully")

	return nil
}

// LoginUser userenticates user credentials and generates tokens
func (a *UserServiceImpl) LoginUser(ctx context.Context, userData *user.LoginUser) (*user.LoginResponse, *user.ErrorStruct) {
	errorStruct := &user.ErrorStruct{}

	if userData.Email == "" || userData.Password == "" {
		errorStruct.Error = user.ErrInvalidUserDetail
		errorStruct.ErrorMsg = "email or password cannot be empty"
		return nil, errorStruct
	}

	data, err := a.userSvc.GetUserByEmail(ctx, userData.Email)
	if err != nil {
		errorStruct.Error = user.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to fetch user from DB: " + err.Error()
		return nil, errorStruct
	}

	if data == nil {
		errorStruct.Error = user.ErrInvalidUserDetail
		errorStruct.ErrorMsg = "no user found with this email"
		return nil, errorStruct
	}

	if err = user.ComparePassword(data.Password, userData.Password); err != nil {
		errorStruct.Error = user.ErrInvalidUserDetail
		errorStruct.ErrorMsg = "incorrect password: " + err.Error()
		return nil, errorStruct
	}

	accessToken, err := user.CreateAccessToken(data.Email)
	if err != nil {
		errorStruct.Error = user.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to create access token: " + err.Error()
		return nil, errorStruct
	}

	refreshToken, err := user.CreateRefreshToken(data.Email)
	if err != nil {
		errorStruct.Error = user.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to create refresh token: " + err.Error()
		return nil, errorStruct
	}

	// loginResponse := struct {
	// 	AccessToken  string     `json:"accessToekn"`
	// 	RefreshToken string     `json:"refreshToken"`
	// 	User         *user.User `json:"user"`
	// }{
	// 	AccessToken:  accessToken,
	// 	User:         data,
	// 	RefreshToken: refreshToken,
	// }

	loginResponse := &user.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         data,
	}

	return loginResponse, nil
}

// RefreshToken validates the refresh token and issues new tokens
func (a *UserServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (string, string, *user.ErrorStruct) {
	errorStruct := &user.ErrorStruct{}

	claims := &user.RefreshClaims{}
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, user.ErrUnexpectedSigningMethod
		}
		return user.REFRESH_SECRET, nil
	})

	if err != nil {
		errorStruct.Error = user.ErrInvalidRefreshToken
		errorStruct.ErrorMsg = "invalid or expired refresh token: " + err.Error()
		return "", "", errorStruct
	}

	accessToken, err := user.CreateAccessToken(claims.Email)
	if err != nil {
		errorStruct.Error = user.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to create new access token: " + err.Error()
		return "", "", errorStruct
	}

	newRefreshToken, err := user.CreateRefreshToken(claims.Email)
	if err != nil {
		errorStruct.Error = user.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to create new refresh token: " + err.Error()
		return "", "", errorStruct
	}

	return accessToken, newRefreshToken, nil
}

func (a *UserServiceImpl) UserByEmail(ctx context.Context, email string) (*user.User, *user.ErrorStruct) {
	errorStruct := &user.ErrorStruct{}
	userData, err := a.userSvc.GetUserByEmail(ctx, email)
	if err != nil {
		errorStruct.Error = user.ErrGettingDataFromDB
		errorStruct.ErrorMsg = "failed to get user data " + err.Error()
		return nil, errorStruct
	}
	return userData, nil

}

func (t *UserServiceImpl) TopPerformer(ctx context.Context) ([]*user.TopPerformer, *user.ErrorStruct) {
	errorStruct := &user.ErrorStruct{}
	data, err := t.userSvc.GetTopPerformer(ctx)
	if err != nil {
		errorStruct.Error = user.ErrGettingDataFromDB
		errorStruct.ErrorMsg = fmt.Sprintf("failed to get data from DB %v", err)
		return nil, errorStruct
	}
	return data, nil
}

func (a *UserServiceImpl) GetDataForDashboard(ctx context.Context) (*user.DashboardData, *user.ErrorStruct) {
	errorStruct := &user.ErrorStruct{}
	userData, err := a.userSvc.GetAllUser(ctx)
	if err != nil {
		errorStruct.Error = user.ErrGettingDataFromDB
		errorStruct.ErrorMsg = "failed to get user data " + err.Error()
		return nil, errorStruct
	}

	dashboardData, err := a.userSvc.GetDashboardTopData(ctx)
	if err != nil {
		errorStruct.Error = user.ErrGettingDataFromDB
		errorStruct.ErrorMsg = "failed to get user data " + err.Error()
		return nil, errorStruct
	}

	response := &user.DashboardData{}
	response.User = userData
	response.DashboardTopData = dashboardData

	return response, nil

}
