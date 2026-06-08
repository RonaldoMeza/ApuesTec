package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserLocked         = errors.New("user is locked")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrUserNotFound       = errors.New("user not found")
)

type RegisterRequest struct {
	FullName string `json:"fullName" binding:"required,min=2,max=150"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type GoogleAuthRequest struct {
	IDToken string `json:"idToken" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8,max=128"`
}

type AuthResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         UserInfo     `json:"user"`
}

type UserInfo struct {
	ID        string   `json:"id"`
	FullName  string   `json:"fullName"`
	Email     string   `json:"email"`
	AvatarURL *string  `json:"avatarUrl,omitempty"`
	Roles     []string `json:"roles"`
}
