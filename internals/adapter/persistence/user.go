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

func (u *UserRepositoryImpl) InsertUserData(ctx context.Context, data *typing.TypingData) error {
	query := `
		INSERT INTO usertypingdata (
		    user_id,
			total_error,
			total_words,
			typed_words,
			total_time,
			total_time_taken_by_user,
			wpm
		) VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	_, err := u.db.ExecContext(
		ctx,
		query,
		data.UserId,
		data.TotalErrors,
		data.TotalWords,
		data.TypedWords,
		data.TotalTime,
		data.TimeTakenByUser,
		data.WPM,
	)

	if err != nil {
		return err
	}
	
	return nil
}
