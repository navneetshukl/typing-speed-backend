package typing

import (
	"context"
	"time"
)

type TypingData struct {
	Email           string    `json:"email"`
	WPM             int       `json:"wpm"`
	TotalErrors     int       `json:"totalErrors"`
	TotalWords      int       `json:"totalWords"`
	TypedWords      int       `json:"typedWords"`
	TotalTime       int       `json:"totalTime"`       // total time of test in second
	TimeTakenByUser int       `json:"timeTakenByUser"` // total time spend by user
	CreatedAt       time.Time `json:"createdAt"`
}

type TypingService interface {
	AddTestData(ctx context.Context, data *TypingData, email string) error
	RecentTestForProfile(ctx context.Context, email string, month string) ([]*TypingData, error)
	SendTypingSentence(ctx context.Context) string
}
