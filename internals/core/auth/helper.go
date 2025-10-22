package auth

import "golang.org/x/crypto/bcrypt"

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