package user

import (
	"context"
	"errors"
	"testing"
	"time"
	"typing-speed/internals/core/user"

	"github.com/golang-jwt/jwt/v5"
)

type FakeUserRepo struct {
	GetByEmailFn          func(ctx context.Context, email string) (*user.User, error)
	CreateFn              func(ctx context.Context, u *user.User) error
	GetTopPerformerFn     func(ctx context.Context) ([]*user.TopPerformer, error)
	GetAllUserFn          func(ctx context.Context) ([]*user.User, error)
	GetDashboardTopDataFn func(ctx context.Context) (*user.DashboardTopData, error)
}

// GetAllUser implements port.UserRepository.
func (f *FakeUserRepo) GetAllUser(ctx context.Context) ([]*user.User, error) {
	if f.GetAllUserFn != nil {
		return f.GetAllUserFn(ctx)
	}
	return nil, nil
}

// GetDashboardTopData implements port.UserRepository.
func (f *FakeUserRepo) GetDashboardTopData(ctx context.Context) (*user.DashboardTopData, error) {
	if f.GetDashboardTopDataFn != nil {
		return f.GetDashboardTopDataFn(ctx)
	}
	return nil, nil
}

// GetTopPerformer implements port.UserRepository.
func (f *FakeUserRepo) GetTopPerformer(ctx context.Context) ([]*user.TopPerformer, error) {
	if f.GetTopPerformerFn != nil {
		return f.GetTopPerformerFn(ctx)
	}
	return nil, nil
}

// UpdateUser implements port.UserRepository.
func (f *FakeUserRepo) UpdateUser(ctx context.Context, email string, speed int, accuracy int, performance int, bestSpeed int) error {
	return nil
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
				if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
					t.Fatalf("expected %v, got %v", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected success, got error")
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		input         *user.LoginUser
		repo          *FakeUserRepo
		mail          *FakeMailSender
		expectErr     bool
		expectedError error
	}{
		{
			name: "empty user details",
			input: &user.LoginUser{
				Email: "Navneet",
				// Email & Password missing
			},
			repo:          &FakeUserRepo{},
			mail:          &FakeMailSender{},
			expectErr:     true,
			expectedError: user.ErrInvalidUserDetail,
		},
		{
			name: "failed to fetch user by email",
			input: &user.LoginUser{
				//Name:     "Navneet",
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
			name: "no user found",
			input: &user.LoginUser{
				//Name:     "Navneet",
				Email:    "navneet@gmail.com",
				Password: "12345",
			},
			repo: &FakeUserRepo{
				GetByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
					return nil, nil

				},
			},
			mail:          &FakeMailSender{},
			expectErr:     true,
			expectedError: user.ErrInvalidUserDetail,
		},
		{
			name: "successful login",
			input: &user.LoginUser{
				//Name:     "Navneet",
				Email:    "navneet@gmail.com",
				Password: "12345",
			},
			repo: &FakeUserRepo{
				GetByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
					hash, _ := user.HashPassword("12345")
					return &user.User{Email: email, Password: hash}, nil
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

			_, err := service.LoginUser(ctx, tt.input)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.expectedError != nil && err != tt.expectedError {
					t.Fatalf("expected %v, got %v", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected success, got error")
				}
			}
		})
	}
}

func TestUserByEmail(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		input         *user.User
		repo          *FakeUserRepo
		mail          *FakeMailSender
		expectErr     bool
		expectedError error
	}{{
		name: "error from db",
		input: &user.User{
			Email:    "navneet",
			Password: "123",
		},
		repo: &FakeUserRepo{
			GetByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
				return nil, errors.New("db error")
			},
		},
		mail:          &FakeMailSender{},
		expectErr:     true,
		expectedError: user.ErrGettingDataFromDB,
	}, {
		name:  "successfully getting user",
		input: &user.User{Email: "navneet"},
		repo: &FakeUserRepo{
			GetByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
				return &user.User{}, nil
			},
		},
		mail:      &FakeMailSender{},
		expectErr: false,
	}}

	for _, tt := range tests {
		service := NewUserService(tt.repo, tt.mail)
		_, err := service.UserByEmail(ctx, tt.input.Email)
		if tt.expectErr {
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if tt.expectedError != nil && err != tt.expectedError {
				t.Fatalf("expected %v, got %v", tt.expectedError, err.Error())
			}
		} else {
			if err != nil {
				t.Fatalf("expected success, got error")
			}
		}

	}
}

func TestTopPerformer(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		repo          *FakeUserRepo
		expectErr     bool
		expectedError error
	}{
		{
			name: "error while fetching top performer",
			repo: &FakeUserRepo{
				GetTopPerformerFn: func(ctx context.Context) ([]*user.TopPerformer, error) {
					return nil, errors.New("db error")
				},
			},
			expectErr:     true,
			expectedError: user.ErrGettingDataFromDB,
		},
		{
			name: "successfully fetched top performer",
			repo: &FakeUserRepo{
				GetTopPerformerFn: func(ctx context.Context) ([]*user.TopPerformer, error) {
					return []*user.TopPerformer{
						{Name: "1", Performance: 100},
					}, nil
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewUserService(tt.repo, nil)

			data, err := service.TopPerformer(ctx)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.expectedError != nil && err != tt.expectedError {
					t.Fatalf("expected %v, got %v", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected success, got error")
				}
				if data == nil {
					t.Fatalf("expected data, got nil")
				}
			}
		})
	}
}

func TestGetDataForDashboard(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		repo          *FakeUserRepo
		expectErr     bool
		expectedError error
	}{
		{
			name: "error while fetching users",
			repo: &FakeUserRepo{
				GetAllUserFn: func(ctx context.Context) ([]*user.User, error) {
					return nil, errors.New("db error")
				},
			},
			expectErr:     true,
			expectedError: user.ErrGettingDataFromDB,
		},
		{
			name: "error while fetching dashboard data",
			repo: &FakeUserRepo{
				GetAllUserFn: func(ctx context.Context) ([]*user.User, error) {
					return []*user.User{}, nil
				},
				GetDashboardTopDataFn: func(ctx context.Context) (*user.DashboardTopData, error) {
					return nil, errors.New("db error")
				},
			},
			expectErr:     true,
			expectedError: user.ErrGettingDataFromDB,
		},
		{
			name: "successful dashboard fetch",
			repo: &FakeUserRepo{
				GetAllUserFn: func(ctx context.Context) ([]*user.User, error) {
					return []*user.User{
						{Email: "test@test.com"},
					}, nil
				},
				GetDashboardTopDataFn: func(ctx context.Context) (*user.DashboardTopData, error) {
					return &user.DashboardTopData{
						TotalTest:       1,
						AverageSpeed:    1,
						AverageAccuracy: 1,
					}, nil
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewUserService(tt.repo, nil)

			result, err := service.GetDataForDashboard(ctx)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.expectedError != nil && err != tt.expectedError {
					t.Fatalf("expected %v, got %v", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected success, got error")
				}
				if result == nil {
					t.Fatalf("expected dashboard data, got nil")
				}
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		refreshToken  string
		expectErr     bool
		expectedError error
	}{
		{
			name:          "invalid refresh token format",
			refreshToken:  "invalid.token.value",
			expectErr:     true,
			expectedError: user.ErrInvalidRefreshToken,
		},
		{
			name: "expired refresh token",
			refreshToken: func() string {
				// create an expired token
				claims := &user.RefreshClaims{
					Email: "test@example.com",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
					},
				}
				token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
					SignedString(user.REFRESH_SECRET)
				return token
			}(),
			expectErr:     true,
			expectedError: user.ErrInvalidRefreshToken,
		},
		{
			name: "successful token refresh",
			refreshToken: func() string {
				token, _ := user.CreateRefreshToken("test@example.com")
				return token
			}(),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewUserService(nil, nil)

			accessToken, refreshToken, err := service.RefreshToken(ctx, tt.refreshToken)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.expectedError != nil && err != tt.expectedError {
					t.Fatalf("expected %v, got %v", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected success, got error")
				}
				if accessToken == "" || refreshToken == "" {
					t.Fatalf("expected tokens, got empty values")
				}
			}
		})
	}
}
