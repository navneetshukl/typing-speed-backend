package typing

import (
	"context"
	"time"
	"typing-speed/internals/core/auth"
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
	AddUserData(ctx context.Context, data *TypingData, email string) *auth.ErrorStruct
	RecentTestForProfile(ctx context.Context, email string) ([]*TypingData, *auth.ErrorStruct) 
}
