package entity

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"regexp"
	"strings"
)

type SMTPCode struct {
	Code string `json:"code"`
}

type RedisData struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	OTP      string `json:"otp"`
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifySignUpRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type AuthResponse struct {
	UserID       int    `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type VerifyForgotPasswordRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

func (s *SignUpRequest) Validate() error {
	s.Email = strings.ToLower(s.Email)
	s.Email = strings.TrimSpace(s.Email)

	return validation.ValidateStruct(
		s,
		validation.Field(
			&s.Email,
			validation.Required,
			is.Email,
		),
		validation.Field(
			&s.Password,
			validation.Required,
			validation.Length(8, 30),
			validation.Match(regexp.MustCompile("^[a-zA-Z0-9]")),
		),
	)
}

func (s *VerifySignUpRequest) Validate() error {
	s.Email = strings.ToLower(s.Email)
	s.Email = strings.TrimSpace(s.Email)

	return validation.ValidateStruct(
		s,
		validation.Field(
			&s.Email,
			validation.Required,
			is.Email,
		),
	)
}

func (s *ResetPasswordRequest) Validate() error {
	s.Email = strings.ToLower(s.Email)
	s.Email = strings.TrimSpace(s.Email)

	return validation.ValidateStruct(
		s,
		validation.Field(
			&s.Email,
			validation.Required,
			is.Email,
		),
		validation.Field(
			&s.NewPassword,
			validation.Required,
			validation.Length(8, 30),
			validation.Match(regexp.MustCompile("^[a-zA-Z0-9]+$")),
		),
	)
}
