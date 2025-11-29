package db

import (
	"context"
	"database/sql"
	"errors"
	"typing-speed/internals/core/user"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepositoryImpl {
	return UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, name, email, password, created_at, avg_speed, avg_accuracy, total_test, level, last_test_time, streak,
        best_speed,avg_performance
		FROM users
		WHERE email = $1;
	`

	user := &user.User{}

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
func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *user.User) error {
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

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, email string, speed, accuracy int) error {
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

func (u *UserRepositoryImpl) GetTopPerformer(ctx context.Context) ([]*user.TopPerformer, error) {
	query := `SELECT name, avg_performance FROM users ORDER BY avg_performance DESC LIMIT 10`

	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	performers := []*user.TopPerformer{}

	for rows.Next() {
		p := &user.TopPerformer{}
		if err := rows.Scan(&p.Name, &p.Performance); err != nil {
			return nil, err
		}
		performers = append(performers, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return performers, nil
}

func (u *UserRepositoryImpl) GetAllUser(ctx context.Context) ([]*user.User, error) {
	query := `
        SELECT 
            id, name, email, password, created_at, 
            avg_speed, avg_accuracy, total_test, level, 
            last_test_time, streak, best_speed, avg_performance
        FROM users;
    `

	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*user.User

	for rows.Next() {
		u := &user.User{}

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
			&u.AvgSpeed,
			&u.AvgAccuracy,
			&u.TotalTest,
			&u.Level,
			&u.LastTestTime,
			&u.Streak,
			&u.BestSpeed,
			&u.AvgPerformance,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	// Check if iteration had an error
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
