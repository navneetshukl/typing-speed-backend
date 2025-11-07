package auth

import (
	"context"
	"typing-speed/internals/adapter/external/sendmail"
	"typing-speed/internals/adapter/port"
	"typing-speed/internals/core/auth"

	"github.com/golang-jwt/jwt/v5"
)

type AuthServiceImpl struct {
	authSvc port.AuthRepository
	mailSvc sendmail.MailSender
}

func NewAuthService(svc port.AuthRepository,mail sendmail.MailSender) auth.AuthService {
	return &AuthServiceImpl{
		authSvc: svc,
		mailSvc: mail,
	}
}

// RegisterUser handles user registration
func (a *AuthServiceImpl) RegisterUser(ctx context.Context, user *auth.User) *auth.ErrorStruct {
	errorStruct := &auth.ErrorStruct{}

	if user.Email == "" || user.Name == "" || user.Password == "" {
		errorStruct.Error = auth.ErrInvalidUserDetail
		errorStruct.ErrorMsg = "user detail is not complete"
		return errorStruct
	}

	data, err := a.authSvc.GetUserByEmail(ctx, user.Email)
	if err != nil {
		errorStruct.Error = auth.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to fetch user by email: " + err.Error()
		return errorStruct
	}
	if data != nil {
		errorStruct.Error = auth.ErrUserAlreadyRegistered
		errorStruct.ErrorMsg = "user already registered with this email"
		return errorStruct
	}

	hash, err := auth.HashPassword(user.Password)
	if err != nil {
		errorStruct.Error = auth.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to hash password: " + err.Error()
		return errorStruct
	}
	user.Password = hash

	err = a.authSvc.CreateUser(ctx, user)
	if err != nil {
		errorStruct.Error = auth.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to create user in database: " + err.Error()
		return errorStruct
	}

	err=a.mailSvc.SendMail("typing@gmail.com",user.Email,"Register","User registered successfully")

	return nil
}

// LoginUser authenticates user credentials and generates tokens
func (a *AuthServiceImpl) LoginUser(ctx context.Context, user *auth.LoginUser) (string, string, *auth.ErrorStruct) {
	errorStruct := &auth.ErrorStruct{}

	if user.Email == "" || user.Password == "" {
		errorStruct.Error = auth.ErrInvalidUserDetail
		errorStruct.ErrorMsg = "email or password cannot be empty"
		return "", "", errorStruct
	}

	data, err := a.authSvc.GetUserByEmail(ctx, user.Email)
	if err != nil {
		errorStruct.Error = auth.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to fetch user from DB: " + err.Error()
		return "", "", errorStruct
	}

	if data == nil {
		errorStruct.Error = auth.ErrInvalidUserDetail
		errorStruct.ErrorMsg = "no user found with this email"
		return "", "", errorStruct
	}

	if err = auth.ComparePassword(data.Password, user.Password); err != nil {
		errorStruct.Error = auth.ErrInvalidUserDetail
		errorStruct.ErrorMsg = "incorrect password: " + err.Error()
		return "", "", errorStruct
	}

	accessToken, err := auth.CreateAccessToken(data.Email)
	if err != nil {
		errorStruct.Error = auth.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to create access token: " + err.Error()
		return "", "", errorStruct
	}

	refreshToken, err := auth.CreateRefreshToken(data.Email)
	if err != nil {
		errorStruct.Error = auth.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to create refresh token: " + err.Error()
		return "", "", errorStruct
	}

	return accessToken, refreshToken, nil
}

// RefreshToken validates the refresh token and issues new tokens
func (a *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (string, string, *auth.ErrorStruct) {
	errorStruct := &auth.ErrorStruct{}

	claims := &auth.RefreshClaims{}
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrUnexpectedSigningMethod
		}
		return auth.REFRESH_SECRET, nil
	})

	if err != nil {
		errorStruct.Error = auth.ErrInvalidRefreshToken
		errorStruct.ErrorMsg = "invalid or expired refresh token: " + err.Error()
		return "", "", errorStruct
	}

	accessToken, err := auth.CreateAccessToken(claims.Email)
	if err != nil {
		errorStruct.Error = auth.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to create new access token: " + err.Error()
		return "", "", errorStruct
	}

	newRefreshToken, err := auth.CreateRefreshToken(claims.Email)
	if err != nil {
		errorStruct.Error = auth.ErrSomethingWentWrong
		errorStruct.ErrorMsg = "failed to create new refresh token: " + err.Error()
		return "", "", errorStruct
	}

	return accessToken, newRefreshToken, nil
}
