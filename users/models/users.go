package model

import "time"

type Users struct {
	Id                       int64     `json:"id"`
	UserUuid                 string    `json:"user_uuid"`
	FirstName                string    `json:"first_name"`
	LastName                 string    `json:"last_name"`
	Email                    string    `json:"email"`
	Phone                    string    `json:"phone"`
	Photo                    string    `json:"photo"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
	AccountId                int64     `json:"account_id"`
	Password                 string    `json:"password"`
	Token                    string    `json:"token"`
	Verified                 bool      `json:"verified"`
	VerificationToken        string    `json:"verification_token"`
	PreferredCommunicationId int       `json:"preferred_communication_id"`
}

// UserLoginRequest :
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type Paymentdatails struct {
	Name   string  `form:"name" binding:"required"`
	Email  string  `form:"email" binding:"required"`
	Amount float64 `form:"amount" binding:"required"`
}
