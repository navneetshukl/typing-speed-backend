package db

import (
	"context"
	"database/sql"
	"errors"
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

func (r *AuthRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	query := `
		SELECT id, name, email, password, created_at, avg_speed, avg_accuracy, total_test, level, last_test_time, streak,
        best_speed,avg_performance
		FROM users
		WHERE email = $1;
	`

	user := &auth.User{}

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.AvgSpeed,
		&user.AvgAccuracy,
		&user.TotalTest,
		&user.Level,
		&user.LastTestTime,
		&user.Streak,
		&user.BestSpeed,
		&user.AvgPerformance,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // user not found
		}
		return nil, err
	}

	return user, nil
}

// CreateUser inserts a new user into the database (no return)
func (r *AuthRepositoryImpl) CreateUser(ctx context.Context, user *auth.User) error {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3);
	`

	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepositoryImpl) UpdateUser(ctx context.Context, email string, speed, accuracy int) error {
	query := `
        UPDATE users
        SET 
            total_test = total_test + 1,
            avg_speed = $2,
            avg_accuracy = $3
        WHERE email = $1;
    `

	_, err := r.db.ExecContext(ctx, query, email, speed, accuracy)
	if err != nil {
		return err
	}

	return nil
}

