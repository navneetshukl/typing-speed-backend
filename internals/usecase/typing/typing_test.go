package typing

import (
	"context"
	"errors"
	"strings"
	"testing"
	"typing-speed/internals/core/typing"
	"typing-speed/internals/core/user"
)

type FakeTypingRepo struct {
	InsertFn    func(ctx context.Context, data *typing.TypingData) error
	GetRecentFn func(ctx context.Context, email string, month int) ([]*typing.TypingData, error)
}
type FakeUserRepo struct {
	GetByEmailFn          func(ctx context.Context, email string) (*user.User, error)
	CreateFn              func(ctx context.Context, u *user.User) error
	GetTopPerformerFn     func(ctx context.Context) ([]*user.TopPerformer, error)
	GetAllUserFn          func(ctx context.Context) ([]*user.User, error)
	GetDashboardTopDataFn func(ctx context.Context) (*user.DashboardTopData, error)
	UpdateUserFn          func(ctx context.Context, email string, speed, acc, perf, best int) error
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
	if f.UpdateUserFn != nil {
		return f.UpdateUserFn(ctx, email, speed, accuracy, performance, bestSpeed)
	}
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

func (f *FakeTypingRepo) InsertTestData(ctx context.Context, data *typing.TypingData) error {
	if f.InsertFn != nil {
		return f.InsertFn(ctx, data)
	}
	return nil
}

func (f *FakeTypingRepo) GetRecentTestData(ctx context.Context, email string, month int) ([]*typing.TypingData, error) {
	if f.GetRecentFn != nil {
		return f.GetRecentFn(ctx, email, month)
	}
	return nil, nil
}

func TestAddTestData(t *testing.T) {
	//ctx := context.Background()

	tests := []struct {
		name          string
		data          *typing.TypingData
		userRepo      *FakeUserRepo
		testRepo      *FakeTypingRepo
		expectErr     bool
		expectedError error
	}{
		{
			name: "insert test data failure",
			data: &typing.TypingData{TotalWords: 10, TypedWords: 8},
			testRepo: &FakeTypingRepo{
				InsertFn: func(ctx context.Context, data *typing.TypingData) error {
					return errors.New("db error")
				},
			},
			expectErr:     true,
			expectedError: typing.ErrInsertingData,
		},
		{
			name: "user fetch failure",
			data: &typing.TypingData{TotalWords: 10, TypedWords: 8},
			testRepo: &FakeTypingRepo{
				InsertFn: func(ctx context.Context, data *typing.TypingData) error {
					return nil
				},
			},
			userRepo: &FakeUserRepo{
				GetByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
					return nil, errors.New("db error")
				},
			},
			expectErr:     true,
			expectedError: typing.ErrGettingDataFromDB,
		},
		{
			name: "update failure",
			data: &typing.TypingData{
				TypedWords: 80,
				TotalWords: 100,
				WPM:        70,
			},
			testRepo: &FakeTypingRepo{
				InsertFn: func(ctx context.Context, data *typing.TypingData) error {
					return nil
				},
			},
			userRepo: &FakeUserRepo{
				GetByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
					return &user.User{
						TotalTest:   1,
						AvgAccuracy: 80,
						AvgSpeed:    60,
						BestSpeed:   70,
					}, nil
				},
				UpdateUserFn: func(ctx context.Context, email string, speed, acc, perf, best int) error {
					return errors.New("db error")
				},
			},
			expectErr: true,
			expectedError: typing.ErrUpdatingTotalTest,
		},
		{
			name: "successful insert",
			data: &typing.TypingData{
				TypedWords: 80,
				TotalWords: 100,
				WPM:        70,
			},
			testRepo: &FakeTypingRepo{
				InsertFn: func(ctx context.Context, data *typing.TypingData) error {
					return nil
				},
			},
			userRepo: &FakeUserRepo{
				GetByEmailFn: func(ctx context.Context, email string) (*user.User, error) {
					return &user.User{
						TotalTest:   1,
						AvgAccuracy: 80,
						AvgSpeed:    60,
						BestSpeed:   70,
					}, nil
				},
				UpdateUserFn: func(ctx context.Context, email string, speed, acc, perf, best int) error {
					return nil
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &TypingServiceImpl{
				userSvc: tt.userRepo,
				testSvc: tt.testRepo,
			}

			err := service.AddTestData(context.Background(), tt.data, "test@mail.com")

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if err != tt.expectedError {
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

func TestRecentTestForProfile(t *testing.T) {
	//ctx := context.Background()

	tests := []struct {
		name          string
		month         string
		repo          *FakeTypingRepo
		expectErr     bool
		expectedError error
	}{
		{
			name:          "invalid month format",
			month:         "abc",
			repo:          &FakeTypingRepo{},
			expectErr:     true,
			expectedError: typing.ErrSomethingWentWrong,
		},
		{
			name:  "db error",
			month: "2",
			repo: &FakeTypingRepo{
				GetRecentFn: func(ctx context.Context, email string, month int) ([]*typing.TypingData, error) {
					return nil, errors.New("db error")
				},
			},
			expectErr:     true,
			expectedError: typing.ErrGettingDataFromDB,
		},
		{
			name:  "success",
			month: "1",
			repo: &FakeTypingRepo{
				GetRecentFn: func(ctx context.Context, email string, month int) ([]*typing.TypingData, error) {
					return []*typing.TypingData{
						{TypedWords: 100, TotalWords: 120},
					}, nil
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &TypingServiceImpl{
				testSvc: tt.repo,
			}

			data, err := service.RecentTestForProfile(context.Background(), "test@example.com", tt.month)

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

func TestSendTypingSentence(t *testing.T) {
	service := &TypingServiceImpl{}

	result := service.SendTypingSentence(context.Background())

	if result == "" {
		t.Fatalf("expected non-empty string, got empty")
	}
	if len(result) != 150 {
		t.Fatalf("expected length 150, got %d", len(result))
	}
	allowed := "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ@#$&"
	for _, ch := range result {
		if !strings.ContainsRune(allowed, ch) {
			t.Fatalf("invalid character found: %c", ch)
		}
	}
}
