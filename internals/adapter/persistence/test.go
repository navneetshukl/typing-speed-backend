package db

import (
	"context"
	"database/sql"
	"fmt"
	"typing-speed/internals/core/typing"
)

type TestRepositoryImpl struct {
	db *sql.DB
}

func NewTestRepository(db *sql.DB) TestRepositoryImpl {
	return TestRepositoryImpl{
		db: db,
	}
}

func (u *TestRepositoryImpl) InsertTestData(ctx context.Context, data *typing.TypingData) error {
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

func (u *TestRepositoryImpl) GetRecentTestData(ctx context.Context, email string, month int) ([]*typing.TypingData, error) {
	days := 30 * month

	var query string

	if month == -1 {
		query = `
			SELECT total_error, total_words, typed_words, total_time,
			       total_time_taken_by_user, wpm, created_at
			FROM user_typing_data
			WHERE email = $1
			ORDER BY created_at DESC
		`
	} else {
		query = fmt.Sprintf(`
			SELECT total_error, total_words, typed_words, total_time,
			       total_time_taken_by_user, wpm, created_at
			FROM user_typing_data
			WHERE email = $1
			  AND created_at >= NOW() - INTERVAL '%d days'
			ORDER BY created_at DESC
		`, days)
	}

	rows, err := u.db.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*typing.TypingData

	for rows.Next() {
		record := &typing.TypingData{}
		if err := rows.Scan(
			&record.TotalErrors,
			&record.TotalWords,
			&record.TypedWords,
			&record.TotalTime,
			&record.TimeTakenByUser,
			&record.WPM,
			&record.CreatedAt,
		); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}


