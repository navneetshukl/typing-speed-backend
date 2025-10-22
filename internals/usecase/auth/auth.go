package auth

import (
	"context"
	"typing-speed/internals/adapter/port"
	"typing-speed/internals/core/auth"
)

type AuthServiceImpl struct {
	authSvc port.AuthRepository
}

func NewAuthService(svc port.AuthRepository) auth.AuthService {
	return &AuthServiceImpl{
		authSvc: svc,
	}
}
func(a *AuthServiceImpl)RegisterUser(ctx context.Context,user *auth.User)error{

	// check if the email is unique

	if user.Email=="" || user.Name=="" || user.Password==""{
		return auth.ErrInvalidUserDetail
	}

	data,err:=a.authSvc.GetUserByEmail(ctx,user.Email)
	if err!=nil{
		return auth.ErrSomethingWentWrong
	}
	if data!=nil{
		return auth.ErrUserAlreadyRegistered
	}

	hash,err:=auth.HashPassword(user.Password)
	if err!=nil{
		return auth.ErrSomethingWentWrong
	}
	user.Password=hash

	// register the user
	err=a.authSvc.CreateUser(ctx,user)
	if err!=nil{
		return auth.ErrSomethingWentWrong
	}
	return nil

}

func(a *AuthServiceImpl)LoginUser(ctx context.Context,user *auth.LoginUser)error{
	if user.Email=="" || user.Password==""{
		return auth.ErrInvalidUserDetail
	}

	data,err:=a.authSvc.GetUserByEmail(ctx,user.Email)
	if err!=nil{
		return auth.ErrSomethingWentWrong
	}
	if data==nil{
		return auth.ErrUserNotFound
	}

	err=auth.ComparePassword(data.Password,user.Password)
	if err!=nil{
		return auth.ErrInvalidUserDetail
	}
	return nil
}