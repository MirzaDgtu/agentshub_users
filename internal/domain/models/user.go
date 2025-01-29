package models

import "time"

type User struct {
	ID              uint      `json:"id"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	FirstName       string    `json:"first_name"`
	MiddleName      string    `json:"middle_name"`
	LastName        string    `json:"last_name"`
	ProfileImageURL string    `json:"profile_imageURL"`
	SignIn          bool      `json:"sign_in"`
	IsBlocked       bool      `json:"is_blocked"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
