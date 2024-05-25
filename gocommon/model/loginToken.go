package model

import (
	"time"
)

type LoginToken struct {
	ID                  string    `json:"id"`
	Email               string    `json:"email"`
	Phone               string    `json:"phone"`
	Timestamp           time.Time `json:"timestamp" redis:"timestamp"`
	PasswordTimestamp   time.Time `json:"password_timestamp" redis:"timestamp"`
	Expire              int64     `json:"expire" redis:"expire"`
	Token               string    `json:"token" redis:"token"`
	SingpassID          string    `json:"singpass_id"`
	SingpassAccessToken string    `json:"singpass_access_token"`
}
