package user

import (
	"context"
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
func (a *UserServiceImpl) RegisterUser(ctx context.Context, userData *user.User) error {
	log.Println("UserData is ", userData.Name, " ", userData.Email, " ", userData.Password)

	if userData.Email == "" || userData.Name == "" || userData.Password == "" {
		return user.ErrInvalidUserDetail
	}

	data, err := a.userSvc.GetUserByEmail(ctx, userData.Email)
	if err != nil {
		return user.ErrSomethingWentWrong
	}
	if data != nil {
		return user.ErrUserAlreadyRegistered
	}

	hash, err := user.HashPassword(userData.Password)
	if err != nil {
		return user.ErrSomethingWentWrong
	}
	userData.Password = hash

	err = a.userSvc.CreateUser(ctx, userData)
	if err != nil {
		return user.ErrSomethingWentWrong
	}

	//err = a.mailSvc.SendMail("typing@gmail.com", userData.Email, "Register", "User registered successfully")

	return nil
}

// LoginUser userenticates user credentials and generates tokens
func (a *UserServiceImpl) LoginUser(ctx context.Context, userData *user.LoginUser) (*user.LoginResponse, error) {

	if userData.Email == "" || userData.Password == "" {
		return nil, user.ErrInvalidUserDetail
	}

	data, err := a.userSvc.GetUserByEmail(ctx, userData.Email)
	if err != nil {
		return nil, user.ErrSomethingWentWrong
	}

	if data == nil {
		return nil, user.ErrInvalidUserDetail
	}

	if err = user.ComparePassword(data.Password, userData.Password); err != nil {
		return nil, user.ErrInvalidUserDetail
	}

	accessToken, err := user.CreateAccessToken(data.Email)
	if err != nil {
		return nil, user.ErrSomethingWentWrong
	}

	refreshToken, err := user.CreateRefreshToken(data.Email)
	if err != nil {
		return nil, user.ErrSomethingWentWrong
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
func (a *UserServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {

	claims := &user.RefreshClaims{}
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, user.ErrUnexpectedSigningMethod
		}
		return user.REFRESH_SECRET, nil
	})

	if err != nil {
		return "", "", user.ErrInvalidRefreshToken
	}

	accessToken, err := user.CreateAccessToken(claims.Email)
	if err != nil {
		return "", "", user.ErrSomethingWentWrong
	}

	newRefreshToken, err := user.CreateRefreshToken(claims.Email)
	if err != nil {
		return "", "", user.ErrSomethingWentWrong
	}

	return accessToken, newRefreshToken, nil
}

func (a *UserServiceImpl) UserByEmail(ctx context.Context, email string) (*user.User, error) {
	userData, err := a.userSvc.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, user.ErrGettingDataFromDB
	}
	return userData, nil

}

func (t *UserServiceImpl) TopPerformer(ctx context.Context) ([]*user.TopPerformer, error) {
	data, err := t.userSvc.GetTopPerformer(ctx)
	if err != nil {

		return nil, user.ErrGettingDataFromDB
	}
	return data, nil
}

func (a *UserServiceImpl) GetDataForDashboard(ctx context.Context) (*user.DashboardData, error) {
	userData, err := a.userSvc.GetAllUser(ctx)
	if err != nil {
		return nil, user.ErrGettingDataFromDB
	}

	dashboardData, err := a.userSvc.GetDashboardTopData(ctx)
	if err != nil {
		return nil, user.ErrGettingDataFromDB
	}

	response := &user.DashboardData{}
	response.User = userData
	response.DashboardTopData = dashboardData

	return response, nil

}
