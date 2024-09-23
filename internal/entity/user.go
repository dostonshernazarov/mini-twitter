package entity

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CreateUserRequest struct {
	ID        string
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserResponse struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	Bio            *string `json:"bio"`
	Role           string  `json:"role"`
	ProfilePicture *string `json:"profile_picture"`
}

type UpdateUserRequest struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
}

type UpdateUserResponse struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	Bio            *string `json:"bio"`
	Role           string  `json:"role"`
	ProfilePicture *string `json:"profile_picture"`
}

type UpdateUserColumnsRequest struct {
	ID     int                    `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

type GetUserResponse struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	Bio            *string `json:"bio"`
	Role           string  `json:"role"`
	Password       string  `json:"-"`
	ProfilePicture *string `json:"profile_picture"`
	FollowingCount int     `json:"following_count"`
	FollowersCount int     `json:"followers_count"`
}

type Filter struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type ListUser struct {
	Users []GetUserResponse `json:"users"`
	Count int               `json:"count"`
}

func (u *CreateUserRequest) Verify() error {
	return validation.ValidateStruct(
		u,
		validation.Field(
			&u.Role,
			validation.Required,
			validation.In(RoleUser),
		),
		validation.Field(
			&u.Email,
			validation.Required,
			is.Email,
		),
		validation.Field(
			&u.Password,
			validation.Required,
			validation.Length(8, 30),
			validation.Match(regexp.MustCompile("^[a-zA-Z0-9]+$")),
		),
	)
}
