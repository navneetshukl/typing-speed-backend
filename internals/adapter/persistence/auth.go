package db

import (
	"context"
	"database/sql"
	"typing-speed/internals/core/auth"
)

type AuthRepositoryImpl struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepositoryImpl {
	return AuthRepositoryImpl{
		db: db,
	}
}

func(r *AuthRepositoryImpl)GetUserByEmail(ctx context.Context,email string)(*auth.User,error){
	query:=`select * from users where email=$1;`
	rows:=r.db.QueryRowContext(ctx,query,email)
	if rows.Err()!=nil{
		if rows.Err()==sql.ErrNoRows{
			return nil,nil
		}
		return nil,rows.Err()
	}
	user:=&auth.User{}
	err:=rows.Scan(user)
	if err!=nil{
		return nil,err
	}
	return user,nil
}

func(r *AuthRepositoryImpl)CreateUser(ctx context.Context,user *auth.User)(error){
	query:=`insert into users (name,email,password) values($1,$2,$3);`
	_,err:=r.db.ExecContext(ctx,query,user.Name,user.Email,user.Password)
	if err!=nil{
		return err
	}
	return nil
}
