package db

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"
	"typing-speed/internals/core/typing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertTestData_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	data := &typing.TypingData{
		Email:           "test@test.com",
		TotalErrors:     1,
		TotalWords:      10,
		TypedWords:      20,
		TotalTime:       20,
		TimeTakenByUser: 15,
		WPM:             10,
	}
	mock.ExpectExec("INSERT INTO user_typing_data").
		WithArgs(data.Email, data.TotalErrors, data.TotalWords,
			data.TypedWords, data.TotalTime, data.TimeTakenByUser, data.WPM).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewTestRepository(db)
	err = repo.InsertTestData(context.Background(), data)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())

}

func TestInsertTestData_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	data := &typing.TypingData{
		Email:           "test@test.com",
		TotalErrors:     1,
		TotalWords:      10,
		TypedWords:      20,
		TotalTime:       20,
		TimeTakenByUser: 15,
		WPM:             10,
	}

	mock.ExpectExec("INSERT INTO user_typing_data").
		WithArgs(
			data.Email,
			data.TotalErrors,
			data.TotalWords,
			data.TypedWords,
			data.TotalTime,
			data.TimeTakenByUser,
			data.WPM,
		).
		WillReturnError(errors.New("insert failed"))

	repo := NewTestRepository(db)

	err = repo.InsertTestData(context.Background(), data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "insert failed")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRecentTest_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTestRepository(db)

	now := time.Now()

	rows := mock.NewRows([]string{
		"total_error",
		"total_words",
		"typed_words",
		"total_time",
		"total_time_taken_by_user",
		"wpm",
		"created_at",
	}).AddRow(1, 2, 3, 4, 5, 6, now)

	query := `
		SELECT total_error, total_words, typed_words, total_time,
		       total_time_taken_by_user, wpm, created_at
		FROM user_typing_data
		WHERE email = $1
		ORDER BY created_at DESC
	`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("test@test.com").
		WillReturnRows(rows)

	data, err := repo.GetRecentTestData(
		context.Background(),
		"test@test.com",
		-1,
	)

	require.NoError(t, err)
	require.Len(t, data, 1)

	assert.Equal(t, 1, data[0].TotalErrors)
	assert.Equal(t, 2, data[0].TotalWords)
	assert.Equal(t, 3, data[0].TypedWords)
	assert.Equal(t, 6, data[0].WPM)
	assert.WithinDuration(t, now, data[0].CreatedAt, time.Second)

	require.NoError(t, mock.ExpectationsWereMet())
}
func TestGetRecentTest_WithMonthFilter_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewTestRepository(db)

	now := time.Now()

	rows := mock.NewRows([]string{
		"total_error",
		"total_words",
		"typed_words",
		"total_time",
		"total_time_taken_by_user",
		"wpm",
		"created_at",
	}).AddRow(2, 10, 8, 60, 55, 80, now)

	// month = 1 → 30 days
	days := 30

	query := fmt.Sprintf(`
		SELECT total_error, total_words, typed_words, total_time,
		       total_time_taken_by_user, wpm, created_at
		FROM user_typing_data
		WHERE email = $1
		  AND created_at >= NOW() - INTERVAL '%d days'
		ORDER BY created_at DESC
	`, days)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("test@test.com").
		WillReturnRows(rows)

	data, err := repo.GetRecentTestData(
		context.Background(),
		"test@test.com",
		1,
	)

	require.NoError(t, err)
	require.Len(t, data, 1)

	assert.Equal(t, 2, data[0].TotalErrors)
	assert.Equal(t, 10, data[0].TotalWords)
	assert.Equal(t, 8, data[0].TypedWords)
	assert.Equal(t, 80, data[0].WPM)
	assert.WithinDuration(t, now, data[0].CreatedAt, time.Second)

	require.NoError(t, mock.ExpectationsWereMet())
}
