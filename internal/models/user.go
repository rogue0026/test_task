package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
)

var (
	ErrInvalidLogin = errors.New("invalid user login")
	ErrInvalidName  = errors.New("invalid user name")
	ErrBadPassword  = errors.New("password length must be eight or greater symbols")
)

type User struct {
	ID                      string `json:"id,omitempty"`
	Login                   string `json:"login"`
	Name                    string `json:"name"`
	Email                   string `json:"email"`
	Password                string `json:"password"`
	IsVerified              bool   `json:"is_verified,omitempty"`
	VerificationCode        string `json:"verification_code,omitempty"`
	VerificationCodeExpires string `json:"verification_code_expires"`
}

func (u *User) HashPassword() error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}
	u.Password = string(passwordHash)
	return nil
}

func (u *User) ValidateData() error {
	if len(u.Login) == 0 {
		return ErrInvalidLogin
	}
	if len(u.Name) == 0 {
		return ErrInvalidName
	}
	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		return err
	}
	if len([]rune(u.Password)) < 8 {
		return ErrBadPassword
	}
	return nil
}
