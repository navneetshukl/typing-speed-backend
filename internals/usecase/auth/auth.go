package auth

import (
	"context"
	"typing-speed/internals/adapter/port"
	"typing-speed/internals/core/auth"
	"github.com/golang-jwt/jwt/v5"
)

type AuthServiceImpl struct {
	authSvc port.AuthRepository
}

func NewAuthService(svc port.AuthRepository) auth.AuthService {
	return &AuthServiceImpl{
		authSvc: svc,
	}
}
func (a *AuthServiceImpl) RegisterUser(ctx context.Context, user *auth.User) error {

	// check if the email is unique

	if user.Email == "" || user.Name == "" || user.Password == "" {
		return auth.ErrInvalidUserDetail
	}

	data, err := a.authSvc.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return auth.ErrSomethingWentWrong
	}
	if data != nil {
		return auth.ErrUserAlreadyRegistered
	}

	hash, err := auth.HashPassword(user.Password)
	if err != nil {
		return auth.ErrSomethingWentWrong
	}
	user.Password = hash

	// register the user
	err = a.authSvc.CreateUser(ctx, user)
	if err != nil {
		return auth.ErrSomethingWentWrong
	}
	return nil

}

func (a *AuthServiceImpl) LoginUser(ctx context.Context, user *auth.LoginUser) (string, string, error) {

	if user.Email == "" || user.Password == "" {
		return "", "", auth.ErrInvalidUserDetail
	}

	data, err := a.authSvc.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", "", auth.ErrSomethingWentWrong
	}
	if data == nil {
		return "", "", auth.ErrSomethingWentWrong
	}

	if err = auth.ComparePassword(data.Password, user.Password); err != nil {
		return "", "", auth.ErrInvalidUserDetail
	}

	// create the jwt tokens
	accessToken, err := auth.CreateAccessToken(data.Email)
	if err != nil {
		return "", "", auth.ErrSomethingWentWrong
	}

	refreshToken, err := auth.CreateRefreshToken(data.Email)
	if err != nil {
		return "", "", auth.ErrSomethingWentWrong
	}

	return accessToken, refreshToken, nil
}

// if refresh token is empty than redirect the user to login page
func (a *AuthServiceImpl) RefreshToken(ctx *context.Context, refreshToken string) (string, string, error) {
	claims := &auth.RefreshClaims{}
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrUnexpectedSigningMethod
		}
		return auth.REFRESH_SECRET, nil
	})
	if err != nil {
		return "", "", auth.ErrInvalidRefreshToken
	}
	// create the jwt tokens
	accessToken, err := auth.CreateAccessToken(claims.Email)
	if err != nil {
		return "", "", auth.ErrSomethingWentWrong
	}

	newRefreshToken, err := auth.CreateRefreshToken(claims.Email)
	if err != nil {
		return "", "", auth.ErrSomethingWentWrong
	}

	return accessToken, newRefreshToken, nil

}
