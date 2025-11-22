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
		INSERT INTO user_typing_data (
		    email,
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
		data.Email,
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

func (u *UserRepositoryImpl) GetRecentTestForProfile(ctx context.Context, email string) ([]*typing.TypingData, error) {
	query := `SELECT total_error,total_words,typed_words,total_time,total_time_taken_by_user,wpm,created_at
			  from user_typing_data WHERE email=$1`

	rows, err := u.db.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	records := []*typing.TypingData{}
	for rows.Next() {
		record := &typing.TypingData{}
		err := rows.Scan(
			&record.TotalErrors,
			&record.TotalWords,
			&record.TypedWords,
			&record.TotalTime,
			&record.TimeTakenByUser,
			&record.WPM,
			&record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		records = append(records, record)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}
