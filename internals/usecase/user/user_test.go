package user

import (
	"context"
	"errors"
	"testing"
	"typing-speed/internals/core/user"
)

type FakeUserRepo struct {
	GetByEmailFn func(ctx context.Context, email string) (*user.User, error)
	CreateFn     func(ctx context.Context, u *user.User) error
}

// GetAllUser implements port.UserRepository.
func (f *FakeUserRepo) GetAllUser(ctx context.Context) ([]*user.User, error) {
	panic("unimplemented")
}

// GetDashboardTopData implements port.UserRepository.
func (f *FakeUserRepo) GetDashboardTopData(ctx context.Context) (*user.DashboardTopData, error) {
	panic("unimplemented")
}

// GetTopPerformer implements port.UserRepository.
func (f *FakeUserRepo) GetTopPerformer(ctx context.Context) ([]*user.TopPerformer, error) {
	panic("unimplemented")
}

// UpdateUser implements port.UserRepository.
func (f *FakeUserRepo) UpdateUser(ctx context.Context, email string, speed int, accuracy int, performance int, bestSpeed int) error {
	panic("unimplemented")
}

func (f *FakeUserRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	if f.GetByEmailFn != nil {
		return f.GetByEmailFn(ctx, email)
	}
	return nil, nil
}

func (f *FakeUserRepo) CreateUser(ctx context.Context, u *user.User) error {
	if f.CreateFn != nil {
		return f.CreateFn(ctx, u)
	}
	return nil
}

type FakeMailSender struct {
	SendFn func(from, to, subject, body string) error
}

func (f *FakeMailSender) SendMail(from, to, subject, body string) error {
	if f.SendFn != nil {
		return f.SendFn(from, to, subject, body)
	}
	return nil
}

func TestRegisterUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		input         *user.User
		repo          *FakeUserRepo
		mail          *FakeMailSender
		expectErr     bool
		expectedError error
	}{
		{
			name: "empty user details",
			input: &user.User{
				Name: "Navneet",
				// Email & Password missing
			},
			repo:          &FakeUserRepo{},
			mail:          &FakeMailSender{},
			expectErr:     true,
			expectedError: user.ErrInvalidUserDetail,
		},
		{
			name: "failed to fetch user by email",
			input: &user.User{
				Name:     "Navneet",
				Email:    "navneet@gmail.com",
				Password: "12345",
			},
			repo: &FakeUserRepo{
				GetByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
					return nil, errors.New("db error")

				},
			},
			mail:          &FakeMailSender{},
			expectErr:     true,
			expectedError: user.ErrSomethingWentWrong,
		},
		{
			name: "user already registered",
			input: &user.User{
				Name:     "Navneet",
				Email:    "navneet@gmail.com",
				Password: "12345",
			},
			repo: &FakeUserRepo{
				GetByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
					return &user.User{Email: email}, nil

				},
			},
			mail:          &FakeMailSender{},
			expectErr:     true,
			expectedError: user.ErrUserAlreadyRegistered,
		},
		{
			name: "failing in registration",
			input: &user.User{
				Name:     "Navneet",
				Email:    "navneet@gmail.com",
				Password: "12345",
			},
			repo: &FakeUserRepo{
				GetByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
					return nil, nil
				},
				CreateFn: func(ctx context.Context, u *user.User) error {
					return errors.New("db error")
				},
			},
			mail:          &FakeMailSender{},
			expectErr:     true,
			expectedError: user.ErrSomethingWentWrong,
		},
		{
			name: "successful registration",
			input: &user.User{
				Name:     "Navneet",
				Email:    "navneet@gmail.com",
				Password: "12345",
			},
			repo: &FakeUserRepo{
				GetByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
					return nil, nil
				},
				CreateFn: func(ctx context.Context, u *user.User) error {
					return nil
				},
			},
			mail:      &FakeMailSender{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			service := NewUserService(tt.repo, tt.mail)

			err := service.RegisterUser(ctx, tt.input)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.expectedError != nil && err.Error != tt.expectedError {
					t.Fatalf("expected %v, got %v", tt.expectedError, err.Error)
				}
			} else {
				if err != nil {
					t.Fatalf("expected success, got error")
				}
			}
		})
	}
}
