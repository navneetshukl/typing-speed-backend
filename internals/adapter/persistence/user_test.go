package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
	"typing-speed/internals/core/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserByEmail_UserFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "created_at",
		"avg_speed", "avg_accuracy", "total_test", "level",
		"last_test_time", "streak", "best_speed", "avg_performance",
	}).AddRow(
		1, "Navneet", "test@test.com", "hashed",
		time.Now(), 50, 95, 10, 1,
		time.Now(), 5, 70, 80,
	)
	mock.ExpectQuery("SELECT (.+) FROM users WHERE email =").
		WithArgs("test@test.com").
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail(context.Background(), "test@test.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Navneet", user.Name)
	assert.Equal(t, "test@test.com", user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestGetUserByEmail_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE email =").
		WithArgs("test@test.com").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByEmail(context.Background(), "test@test.com")

	require.NoError(t, err)
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	u := &user.User{
		Name:     "Navneet",
		Email:    "test@test.com",
		Password: "hashed",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(u.Name, u.Email, u.Password).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected

	err = repo.CreateUser(context.Background(), u)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
func TestCreateUser_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	u := &user.User{
		Name:     "Navneet",
		Email:    "test@test.com",
		Password: "hashed",
	}

	dbErr := errors.New("duplicate key value violates unique constraint")

	mock.ExpectExec("INSERT INTO users").
		WithArgs(u.Name, u.Email, u.Password).
		WillReturnError(dbErr)

	err = repo.CreateUser(context.Background(), u)

	require.Error(t, err)
	assert.Equal(t, dbErr, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
func TestCreateUser_ContextCancelled(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	u := &user.User{
		Name:     "Navneet",
		Email:    "test@test.com",
		Password: "hashed",
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel BEFORE call

	err = repo.CreateUser(ctx, u)

	require.Error(t, err)
	assert.Equal(t, context.Canceled, err)

	// IMPORTANT: no ExpectExec, so no DB call expected
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTopPerformer_Success(t *testing.T) {
	//	query := `SELECT name, avg_performance FROM users ORDER BY avg_performance DESC LIMIT 10`

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	rows := sqlmock.NewRows([]string{
		"name", "avg_performance",
	}).AddRow("navneet", 123)

	mock.ExpectQuery("SELECT name, avg_performance FROM users ORDER BY avg_performance DESC LIMIT 10").WithArgs().WillReturnRows(rows)
	data, err := repo.GetTopPerformer(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "navneet", data[0].Name)
	assert.Equal(t, 123, data[0].Performance)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestGetAllUsers_UserFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)
	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "created_at",
		"avg_speed", "avg_accuracy", "total_test", "level",
		"last_test_time", "streak", "best_speed", "avg_performance",
	}).AddRow(
		1, "Navneet", "test@test.com", "hashed",
		time.Now(), 50, 95, 10, 1,
		time.Now(), 5, 70, 80,
	)
	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs().
		WillReturnRows(rows)

	user, err := repo.GetAllUser(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Navneet", user[0].Name)
	assert.Equal(t, "test@test.com", user[0].Email)

	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestGetAllUser_NoUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "password", "created_at",
		"avg_speed", "avg_accuracy", "total_test", "level",
		"last_test_time", "streak", "best_speed", "avg_performance",
	})
	// no AddRow → empty result set

	mock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnRows(rows)

	users, err := repo.GetAllUser(context.Background())

	require.NoError(t, err)
	assert.Empty(t, users)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDashboardTopData_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{
		"avg_speed", "avg_accuracy", "total_test",
	}).AddRow(
		55.5, 92.3, int64(120),
	)

	mock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnRows(rows)

	data, err := repo.GetDashboardTopData(context.Background())

	require.NoError(t, err)
	require.NotNil(t, data)
	assert.Equal(t, int64(120), data.TotalTest)
	assert.Equal(t, 55, data.AverageSpeed)
	assert.Equal(t, 92, data.AverageAccuracy)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDashboardTopData_NoUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{
		"avg_speed", "avg_accuracy", "total_test",
	}).AddRow(
		0.0, 0.0, int64(0),
	)

	mock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnRows(rows)

	data, err := repo.GetDashboardTopData(context.Background())

	require.NoError(t, err)
	assert.Equal(t, int64(0), data.TotalTest)
	assert.Equal(t, 0, data.AverageSpeed)
	assert.Equal(t, 0, data.AverageAccuracy)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDashboardTopData_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnError(errors.New("db down"))

	data, err := repo.GetDashboardTopData(context.Background())

	require.Error(t, err)
	assert.Nil(t, data)

	require.NoError(t, mock.ExpectationsWereMet())
}
func TestUpdateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectExec("UPDATE users").
		WithArgs(
			"test@test.com",
			60,
			95,
			85,
			sqlmock.AnyArg(), // time.Now()
			70,
		).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err = repo.UpdateUser(
		context.Background(),
		"test@test.com",
		60,
		95,
		85,
		70,
	)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
func TestUpdateUser_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectExec("UPDATE users").
		WithArgs(
			"missing@test.com",
			60,
			95,
			85,
			sqlmock.AnyArg(),
			70,
		).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	err = repo.UpdateUser(
		context.Background(),
		"missing@test.com",
		60,
		95,
		85,
		70,
	)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
func TestUpdateUser_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	dbErr := errors.New("db down")

	mock.ExpectExec("UPDATE users").
		WithArgs(
			"test@test.com",
			60,
			95,
			85,
			sqlmock.AnyArg(),
			70,
		).
		WillReturnError(dbErr)

	err = repo.UpdateUser(
		context.Background(),
		"test@test.com",
		60,
		95,
		85,
		70,
	)

	require.Error(t, err)
	assert.Equal(t, dbErr, err)

	require.NoError(t, mock.ExpectationsWereMet())
}
func TestUpdateUser_ContextCancelled(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = repo.UpdateUser(
		ctx,
		"test@test.com",
		60,
		95,
		85,
		70,
	)

	require.Error(t, err)
	assert.Equal(t, context.Canceled, err)

	// No DB call expected
	require.NoError(t, mock.ExpectationsWereMet())
}
