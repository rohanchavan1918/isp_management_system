package models

import (
	"time"
)

type User struct {
	ID uint `gorm:"primaryKey,index"`
	// Person info
	FirstName   string
	LastName    string
	Email       string
	MobileNo    string
	DateOfBirth string
	Gender      string
	Password    string
	// Meta
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsVerified bool
	IsActive   bool
}

type ForgotPassword struct {
	ID         uint `gorm:"primaryKey"`
	User       User
	UserID     int
	OTP        int
	Token      string
	Expired_at time.Time
	Created_at time.Time
}
