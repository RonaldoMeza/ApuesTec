package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/RonaldoMeza/ApuesTec/backend/internal/audit"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/config"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/roles"
	"github.com/RonaldoMeza/ApuesTec/backend/internal/users"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest, userAgent, ipAddress string) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest, userAgent, ipAddress string) (*AuthResponse, error)
	Refresh(ctx context.Context, refreshTokenStr, userAgent, ipAddress string) (*AuthResponse, error)
	Logout(ctx context.Context, refreshTokenStr string) error
	Me(ctx context.Context, userID string) (*UserInfo, error)
	GoogleAuth(ctx context.Context, req GoogleAuthRequest, userAgent, ipAddress string) (*AuthResponse, error)
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error
}

type service struct {
	userRepo  users.Repository
	authRepo  Repository
	roleRepo  roles.Repository
	auditRepo audit.Repository
	cfg       *config.Config
}

func NewService(
	userRepo users.Repository,
	authRepo Repository,
	roleRepo roles.Repository,
	auditRepo audit.Repository,
	cfg *config.Config,
) Service {
	return &service{
		userRepo:  userRepo,
		authRepo:  authRepo,
		roleRepo:  roleRepo,
		auditRepo: auditRepo,
		cfg:       cfg,
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (s *service) Register(ctx context.Context, req RegisterRequest, userAgent, ipAddress string) (*AuthResponse, error) {
	if len(req.Password) < 8 {
		return nil, &ServiceError{Status: http.StatusBadRequest, Code: "WEAK_PASSWORD", Message: ErrWeakPassword.Error()}
	}

	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, &ServiceError{Status: http.StatusConflict, Code: "EMAIL_EXISTS", Message: ErrEmailAlreadyExists.Error()}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.cfg.BcryptCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepo.Create(ctx, req.FullName, req.Email, string(passwordHash))
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	userRole, err := s.roleRepo.GetByName(ctx, "USER")
	if err != nil {
		return nil, fmt.Errorf("get user role: %w", err)
	}

	if err := s.roleRepo.AssignRole(ctx, user.ID, userRole.ID); err != nil {
		return nil, fmt.Errorf("assign role: %w", err)
	}

	userID := user.ID
	_ = s.auditRepo.Log(ctx, &userID, audit.ActionUserRegistered, "user", user.ID, strPtr(ipAddress), strPtr(userAgent))

	return s.generateAuthResponse(ctx, user.ID, user.FullName, user.Email, user.AvatarURL, userAgent, ipAddress)
}

func (s *service) Login(ctx context.Context, req LoginRequest, userAgent, ipAddress string) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &ServiceError{Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS", Message: ErrInvalidCredentials.Error()}
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	if user.Status != "ACTIVE" {
		return nil, &ServiceError{Status: http.StatusForbidden, Code: "USER_DISABLED", Message: "account is disabled"}
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		remaining := time.Until(*user.LockedUntil).Minutes()
		return nil, &ServiceError{
			Status:  http.StatusTooManyRequests,
			Code:    "USER_LOCKED",
			Message: fmt.Sprintf("account locked for %.0f more minutes", remaining),
		}
	}

	if user.LockedUntil != nil && user.LockedUntil.Before(time.Now()) {
		_ = s.userRepo.ResetFailedAttempts(ctx, user.ID)
	}

	if user.PasswordHash == nil {
		return nil, &ServiceError{Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS", Message: ErrInvalidCredentials.Error()}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		attempts, incErr := s.userRepo.IncrementFailedAttempts(ctx, user.ID)
		if incErr != nil {
			return nil, fmt.Errorf("increment attempts: %w", incErr)
		}

		_ = s.auditRepo.Log(ctx, &user.ID, audit.ActionLoginFailed, "user", user.ID, strPtr(ipAddress), strPtr(userAgent))

		if attempts >= s.cfg.LoginMaxAttempts {
			lockUntil := time.Now().Add(time.Duration(s.cfg.LoginLockMinutes) * time.Minute)
			if lockErr := s.userRepo.LockUser(ctx, user.ID, lockUntil); lockErr != nil {
				return nil, fmt.Errorf("lock user: %w", lockErr)
			}
			_ = s.auditRepo.Log(ctx, &user.ID, audit.ActionUserLocked, "user", user.ID, strPtr(ipAddress), strPtr(userAgent))
			return nil, &ServiceError{
				Status:  http.StatusTooManyRequests,
				Code:    "USER_LOCKED",
				Message: fmt.Sprintf("account locked for %d minutes due to too many failed attempts", s.cfg.LoginLockMinutes),
			}
		}

		return nil, &ServiceError{Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS", Message: ErrInvalidCredentials.Error()}
	}

	_ = s.userRepo.ResetFailedAttempts(ctx, user.ID)
	_ = s.auditRepo.Log(ctx, &user.ID, audit.ActionUserLoggedIn, "user", user.ID, strPtr(ipAddress), strPtr(userAgent))

	return s.generateAuthResponse(ctx, user.ID, user.FullName, user.Email, user.AvatarURL, userAgent, ipAddress)
}

func (s *service) Refresh(ctx context.Context, refreshTokenStr, userAgent, ipAddress string) (*AuthResponse, error) {
	tokenHash := hashToken(refreshTokenStr)
	storedToken, err := s.authRepo.FindRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusUnauthorized, Code: "INVALID_REFRESH_TOKEN", Message: ErrInvalidToken.Error()}
	}

	if storedToken.RevokedAt != nil {
		return nil, &ServiceError{Status: http.StatusUnauthorized, Code: "TOKEN_REVOKED", Message: ErrInvalidToken.Error()}
	}

	if storedToken.ExpiresAt.Before(time.Now()) {
		return nil, &ServiceError{Status: http.StatusUnauthorized, Code: "TOKEN_EXPIRED", Message: ErrInvalidToken.Error()}
	}

	if err := s.authRepo.RevokeRefreshToken(ctx, storedToken.ID); err != nil {
		return nil, fmt.Errorf("revoke old refresh token: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	_ = s.auditRepo.Log(ctx, &user.ID, audit.ActionTokenRefreshed, "user", user.ID, strPtr(ipAddress), strPtr(userAgent))

	return s.generateAuthResponse(ctx, user.ID, user.FullName, user.Email, user.AvatarURL, userAgent, ipAddress)
}

func (s *service) Logout(ctx context.Context, refreshTokenStr string) error {
	tokenHash := hashToken(refreshTokenStr)
	storedToken, err := s.authRepo.FindRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil
	}

	_ = s.authRepo.RevokeRefreshToken(ctx, storedToken.ID)
	_ = s.auditRepo.Log(ctx, &storedToken.UserID, audit.ActionUserLoggedOut, "refresh_token", storedToken.ID, nil, nil)
	return nil
}

func (s *service) Me(ctx context.Context, userID string) (*UserInfo, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, &ServiceError{Status: http.StatusNotFound, Code: "USER_NOT_FOUND", Message: ErrUserNotFound.Error()}
	}

	userRoles, err := s.roleRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}

	roles := make([]string, len(userRoles))
	for i, r := range userRoles {
		roles[i] = r.Name
	}

	return &UserInfo{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		Roles:     roles,
	}, nil
}

func (s *service) GoogleAuth(ctx context.Context, req GoogleAuthRequest, userAgent, ipAddress string) (*AuthResponse, error) {
	googleClaims, err := verifyGoogleIDToken(ctx, req.IDToken, s.cfg.GoogleClientID)
	if err != nil {
		_ = s.auditRepo.Log(ctx, nil, audit.ActionGoogleAuthFailed, "auth", "", strPtr(ipAddress), strPtr(userAgent))
		return nil, &ServiceError{Status: http.StatusUnauthorized, Code: "INVALID_ID_TOKEN", Message: "invalid Google ID token"}
	}

	sub, _ := googleClaims["sub"].(string)
	email, _ := googleClaims["email"].(string)
	name, _ := googleClaims["name"].(string)

	existingAccount, err := s.authRepo.FindAuthAccount(ctx, "google", sub)
	if err == nil && existingAccount != nil {
		user, err := s.userRepo.FindByID(ctx, existingAccount.UserID)
		if err != nil {
			return nil, fmt.Errorf("find user for google account: %w", err)
		}
		_ = s.auditRepo.Log(ctx, &user.ID, audit.ActionGoogleAuthSuccess, "user", user.ID, strPtr(ipAddress), strPtr(userAgent))
		return s.generateAuthResponse(ctx, user.ID, user.FullName, user.Email, user.AvatarURL, userAgent, ipAddress)
	}

	if email == "" {
		return nil, &ServiceError{Status: http.StatusBadRequest, Code: "EMAIL_REQUIRED", Message: "Google account must have an email"}
	}

	fullName := name
	if fullName == "" {
		fullName = email
	}

	user, err := s.userRepo.Create(ctx, fullName, email, "")
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			user, err = s.userRepo.FindByEmail(ctx, email)
			if err != nil {
				return nil, fmt.Errorf("find existing user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("create user from google: %w", err)
		}
	}

	if err := s.authRepo.CreateAuthAccount(ctx, user.ID, "google", sub, email); err != nil {
		return nil, fmt.Errorf("create google auth account: %w", err)
	}

	userRole, err := s.roleRepo.GetByName(ctx, "USER")
	if err != nil {
		return nil, fmt.Errorf("get user role: %w", err)
	}
	_ = s.roleRepo.AssignRole(ctx, user.ID, userRole.ID)

	_ = s.auditRepo.Log(ctx, &user.ID, audit.ActionGoogleAuthSuccess, "user", user.ID, strPtr(ipAddress), strPtr(userAgent))

	return s.generateAuthResponse(ctx, user.ID, user.FullName, user.Email, user.AvatarURL, userAgent, ipAddress)
}

func (s *service) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return &ServiceError{Status: http.StatusNotFound, Code: "USER_NOT_FOUND", Message: ErrUserNotFound.Error()}
	}

	if user.PasswordHash == nil || *user.PasswordHash == "" {
		return &ServiceError{Status: http.StatusBadRequest, Code: "NO_PASSWORD_SET", Message: "no password set, use Google login"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(currentPassword)); err != nil {
		return &ServiceError{Status: http.StatusUnauthorized, Code: "WRONG_PASSWORD", Message: "current password is incorrect"}
	}

	if len(newPassword) < 8 {
		return &ServiceError{Status: http.StatusBadRequest, Code: "WEAK_PASSWORD", Message: ErrWeakPassword.Error()}
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), s.cfg.BcryptCost)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, user.ID, string(newHash)); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	if err := s.authRepo.RevokeAllUserRefreshTokens(ctx, user.ID); err != nil {
		return fmt.Errorf("revoke sessions: %w", err)
	}

	_ = s.auditRepo.Log(ctx, &user.ID, audit.ActionPasswordChanged, "user", user.ID, nil, nil)
	return nil
}

func (s *service) generateAuthResponse(ctx context.Context, userID, fullName, email string, avatarURL *string, userAgent, ipAddress string) (*AuthResponse, error) {
	userRoles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}

	roles := make([]string, len(userRoles))
	for i, r := range userRoles {
		roles[i] = r.Name
	}

	accessToken, err := s.generateAccessToken(userID, email, roles)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshTokenRaw, refreshTokenHash, err := generateRefreshTokenPair()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	expiresAt := time.Now().Add(s.cfg.JWTRefreshTTL)
	if err := s.authRepo.CreateRefreshToken(ctx, userID, refreshTokenHash, expiresAt, strPtr(userAgent), strPtr(ipAddress)); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenRaw,
		User: UserInfo{
			ID:        userID,
			FullName:  fullName,
			Email:     email,
			AvatarURL: avatarURL,
			Roles:     roles,
		},
	}, nil
}

type Claims struct {
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func (s *service) generateAccessToken(userID, email string, roles []string) (string, error) {
	now := time.Now()
	claims := Claims{
		Email: email,
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.JWTAccessTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTAccessSecret))
}

func generateRefreshTokenPair() (raw string, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate random bytes: %w", err)
	}
	raw = hex.EncodeToString(b)
	hash = hashToken(raw)
	return raw, hash, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

type ServiceError struct {
	Status  int
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

func verifyGoogleIDToken(ctx context.Context, idToken, clientID string) (map[string]interface{}, error) {
	if clientID == "" {
		return nil, fmt.Errorf("google client ID not configured")
	}

	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("parse unverified google token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid google token claims")
	}

	if aud, ok := claims["aud"].(string); !ok || aud != clientID {
		return nil, fmt.Errorf("token audience mismatch")
	}

	if iss, ok := claims["iss"].(string); ok {
		if iss != "accounts.google.com" && iss != "https://accounts.google.com" {
			return nil, fmt.Errorf("invalid issuer: %s", iss)
		}
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token expired")
		}
	}

	return claims, nil
}
