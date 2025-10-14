package db

import (
	"context"
	"database/sql"
	"typing-speed/internals/core/typing"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepositoryImpl {
	return UserRepositoryImpl{
		db: db,
	}
}
func (u *UserRepositoryImpl) InsertUserData(ctx context.Context, data *typing.TypingData) error{
	query:=`INSERT into users (user_id,wpm,total_error,total_words,typed_words,createdAt,total_time,time_taken_by_user) VALUES 
	($1,$2,$3,$4,$5,$6,$7,$8);`

	_,err:=u.db.ExecContext(ctx,query,data.UserId,data.WPM,data.TotalErrors,data.TotalErrors,data.TypedWords,data.CreatedAt,
	data.TotalTime,data.TimeTakenByUser)
	if err!=nil{
		return err
	}
	return nil
}
