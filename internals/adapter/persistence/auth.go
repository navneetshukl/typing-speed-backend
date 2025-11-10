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

func (r *AuthRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	query := `SELECT id, name, email, password, created_at FROM users WHERE email = $1;`

	user := &auth.User{}

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // no user found
		}
		return nil, err
	}

	return user, nil
}

func (r *AuthRepositoryImpl) CreateUser(ctx context.Context, user *auth.User) error {
	query := `insert into users (name,email,password) values($1,$2,$3);`
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepositoryImpl) UpdateTotalTest(ctx context.Context, email string) error {
	query := `
        UPDATE users
        SET total_test = total_test + 1
        WHERE email = $1;
    `

	_, err := r.db.ExecContext(ctx, query, email)
	if err != nil {
		return err
	}

	return nil
}
