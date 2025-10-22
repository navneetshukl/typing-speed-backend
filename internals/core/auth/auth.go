package auth

import "time"

type User struct{
	ID int `json:"id,omitempty"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginUser struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type AuthService interface{}