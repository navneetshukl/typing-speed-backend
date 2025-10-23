package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const(
	ACCESS_SECRET string="access_secret_code"
	REFRESH_SECRET string="refresh_secret_code"

)

func HashPassword(password string)(string,error){
	pass,err:=bcrypt.GenerateFromPassword([]byte(password),10)
	if err!=nil{
		return "",err
	}
	return string(pass),nil
}

func ComparePassword(hashed,password string)error{
	err:=bcrypt.CompareHashAndPassword([]byte(hashed),[]byte(password))
	if err!=nil{
		return err
	}
	return nil
}

func CreateAccessToken(email string)(string,error){
	now:=time.Now()
	claims:=AccessClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: uuid.NewString(),
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24*time.Hour)),
			Issuer: "typing-app",
			Subject: email,

		},
	}
	t:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return t.SignedString(ACCESS_SECRET)
}

func CreateRefreshToken(email string)(string,error){
	now:=time.Now()
	claims:=RefreshClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: uuid.NewString(),
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24*7*time.Hour)),
			Issuer: "typing-app",
			Subject: email,

		},
	}
	t:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return t.SignedString(REFRESH_SECRET)
}