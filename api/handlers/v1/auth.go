package v1

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dostonshernazarov/mini-twitter/internal/entity"
	cache "github.com/dostonshernazarov/mini-twitter/internal/infrastructure/repository/redis"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/etc"
	tokens "github.com/dostonshernazarov/mini-twitter/internal/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

// SignUp
// @Summary 		Sign Up
// @Description		this api for sign-up a new users
// @Tags 			auth
// @Accept 			json
// @Produce 		json
// @Param 			request body entity.SignUpRequest true "Sign Up Model"
// @Success 		200 {object} entity.ResponseWithMessage
// @Failure 		400 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/auth/sign-up [POST]
func (h *HandlerV1) SignUp(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.TimeOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var request entity.SignUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return
	}

	emailExists, err := h.User.UniqueEmail(ctx, request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	if emailExists {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.EmailUsed,
		})
		log.Println(request.Email, entity.EmailUsed)
		return
	}

	usernameExists, err := h.User.UniqueUsername(ctx, request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	if usernameExists {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.UsernameTaken,
		})
		log.Println(request.Username, entity.UsernameTaken)
		return
	}

	otp := etc.GenerateCode(6)

	err = etc.SendMessage([]string{request.Email}, entity.SMTPCode{Code: otp}, "./internal/pkg/etc/otp.html", *h.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	hashedPasswd, err := etc.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	data := entity.RedisData{
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPasswd,
		OTP:      otp,
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	if err := cache.Set(ctx, request.Email, string(bytes), time.Minute*5); err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.ResponseWithMessage{
		Message: entity.SendOTP,
	})
}

// VerifySignUp
// @Summary 		Verify Sign Up
// @Description		this api for verify sign-up request
// @Tags 			auth
// @Accept 			json
// @Produce 		json
// @Param 			request body entity.VerifySignUpRequest true "Verify Sign Up Model"
// @Success 		200 {object} entity.AuthResponse
// @Failure 		400 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/auth/verify [POST]
func (h *HandlerV1) VerifySignUp(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.TimeOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var request entity.VerifySignUpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return
	}

	value, err := cache.Get(ctx, request.Email)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.JSON(http.StatusBadRequest, entity.Error{
				Message: entity.OTPExpired,
			})
			log.Println(request.Email, entity.OTPExpired)
			return
		} else {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}
	}

	var data entity.RedisData
	if err := json.Unmarshal([]byte(cast.ToString(value)), &data); err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	user, err := h.User.Create(ctx, entity.CreateUserRequest{
		ID:       uuid.NewString(),
		Name:     data.Name,
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
		Role:     entity.RoleUser,
	})

	if err != nil {
		println("error")
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	jwtHandler := tokens.JwtHandler{
		Sub:       user.ID,
		Role:      user.Role,
		SigninKey: h.Config.SigningKey,
	}

	access, refresh, err := jwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		UserID:       user.ID,
	})
}

// LogIn
// @Summary 		Log In
// @Description		this api for login to site
// @Tags 			auth
// @Accept 			json
// @Produce 		json
// @Param 			request body entity.LoginRequest true "Login Model"
// @Success 		200 {object} entity.AuthResponse
// @Failure 		400 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/auth/login [POST]
func (h *HandlerV1) LogIn(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.TimeOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var request entity.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return
	}

	user, err := h.User.Get(ctx, map[string]interface{}{
		"username": request.Username,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, entity.Error{
				Message: entity.WrongLoginOrPasswd,
			})
			log.Println(err.Error())
			return
		} else {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}
	}

	if !etc.CheckPasswordHash(request.Password, user.Password) {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.WrongLoginOrPasswd,
		})
		return
	}

	jwtHandler := tokens.JwtHandler{
		Sub:  user.ID,
		Role: user.Role,
	}

	access, refresh, err := jwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		UserID:       user.ID,
	})
}

// ForgotPassword
// @Summary 		Forgot Password
// @Description		this api for sending request about forgot password
// @Tags 			auth
// @Accept 			json
// @Produce 		json
// @Param 			email path string true "Email"
// @Success 		200 {object} entity.ResponseWithMessage
// @Failure 		400 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/auth/forgot-password/{email} [POST]
func (h *HandlerV1) ForgotPassword(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.TimeOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	email := c.Param("email")
	fmt.Println(email)

	user, err := h.User.Get(ctx, map[string]interface{}{
		"email": email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, entity.Error{
				Message: entity.NotFoundData,
			})
			log.Println(err.Error())
			return
		} else {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}
	}

	otp := etc.GenerateCode(6)

	bytes, err := json.Marshal(&otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	if err := cache.Set(ctx, user.Email, string(bytes), time.Minute*3); err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	if err := etc.SendMessage([]string{user.Email}, entity.SMTPCode{Code: otp}, "./internal/pkg/etc/otp.html", *h.Config); err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.ResponseWithMessage{
		Message: entity.SendOTP,
	})
}

// VerifyForgotPassword
// @Summary 		Verify Forgot Password
// @Description		this api for verify forgot password
// @Tags 			auth
// @Accept 			json
// @Produce 		json
// @Param 			request body entity.VerifyForgotPasswordRequest true "Verify Forgot Password Model"
// @Success 		200 {object} entity.ResponseWithStatus
// @Failure 		400 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/auth/verify-forgot-password [POST]
func (h *HandlerV1) VerifyForgotPassword(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.TimeOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var request entity.VerifyForgotPasswordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return
	}

	value, err := cache.Get(ctx, request.Email)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.JSON(http.StatusBadRequest, entity.Error{
				Message: entity.OTPExpired,
			})
			log.Println(request.Email, entity.OTPExpired)
			return
		} else {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}
	}

	var checkOTP string

	if err := json.Unmarshal([]byte(cast.ToString(value)), &checkOTP); err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	if checkOTP != request.Code {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(checkOTP, request.Code)
		return
	}

	c.JSON(http.StatusOK, entity.ResponseWithStatus{
		Status: true,
	})
}

// ResetPassword
// @Summary 		Reset Password
// @Description		this api for reset password
// @Tags 			auth
// @Accept 			json
// @Produce 		json
// @Param 			request body entity.ResetPasswordRequest true "Reset Password Model"
// @Success 		200 {object} entity.ResponseWithStatus
// @Failure 		400 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/auth/reset-password [PUT]
func (h *HandlerV1) ResetPassword(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.TimeOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var request entity.ResetPasswordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, entity.Error{
			Message: entity.IncorrectData,
		})
		log.Println(err.Error())
		return
	}

	user, err := h.User.Get(ctx, map[string]interface{}{
		"email": request.Email,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, entity.Error{
				Message: entity.NotFoundData,
			})
			log.Println(err.Error())
			return
		} else {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}
	}

	hashed, err := etc.HashPassword(request.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	if err := h.User.UpdatePasswd(ctx, user.ID, hashed); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, entity.Error{
				Message: entity.NotFoundData,
			})
			log.Println(err.Error())
			return
		} else {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, entity.ResponseWithStatus{
		Status: true,
	})
}

// GetNewToken
// @Summary 		Get New Access
// @Description		this api for getting a new access and refresh token
// @Tags 			auth
// @Accept 			json
// @Produce 		json
// @Param 			refresh path string true "Refresh token"
// @Success 		200 {object} entity.AuthResponse
// @Failure 		400 {object} entity.Error
// @Failure 		500 {object} entity.Error
// @Router 			/v1/auth/refresh/{refresh} [GET]
func (h *HandlerV1) GetNewToken(c *gin.Context) {
	duration, err := time.ParseDuration(h.Config.Context.TimeOut)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	refreshToken := c.Param("refresh")

	claims, err := tokens.ExtractClaim(refreshToken, []byte(h.Config.SigningKey))
	if err != nil {
		c.JSON(http.StatusUnauthorized, entity.Error{
			Message: entity.TokenExpired,
		})
		log.Println(err.Error())
		return
	}

	userID := cast.ToInt(claims["sub"])

	user, err := h.User.Get(ctx, map[string]interface{}{
		"id": userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, entity.Error{
				Message: entity.NotFoundData,
			})
			log.Println(err.Error())
			return
		} else {
			c.JSON(http.StatusInternalServerError, entity.Error{
				Message: entity.ServerError,
			})
			log.Println(err.Error())
			return
		}
	}

	jwtHandler := tokens.JwtHandler{
		Sub:       user.ID,
		Role:      user.Role,
		SigninKey: h.Config.SigningKey,
	}

	access, refresh, err := jwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.Error{
			Message: entity.ServerError,
		})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.AuthResponse{
		UserID:       user.ID,
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
