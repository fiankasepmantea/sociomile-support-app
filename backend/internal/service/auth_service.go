package service

import (
	"context"
	"time"

	"backend/internal/appcontext"
	"backend/internal/apperrors"
	"backend/internal/cache"
	"backend/internal/repository"
	"backend/internal/repository/pg"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *pg.Repository
	redis *cache.Redis
	jwtSecret string
	jwtExpiry int
}

func NewAuthService(repo *pg.Repository, jwtSecret string, jwtExpiryHours int) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiryHours,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	ExpiresAt time.Time   `json:"expires_at"`
	User      interface{} `json:"user"`
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Get tenant from context (injected by tenant middleware)
	appCtx, err := appcontext.FromContext(ctx)
	if err != nil {
		return nil, apperrors.ErrInternal("tenant context missing")
	}

	// Find user by email in this tenant
	user, err := s.repo.User.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, apperrors.ErrUnauthorized("invalid email or password")
		}
		return nil, apperrors.ErrInternal("database error")
	}

	// Verify user belongs to current tenant
	if user.TenantID != appCtx.TenantID {
		return nil, apperrors.ErrUnauthorized("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrUnauthorized("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, apperrors.ErrForbidden("user account is inactive")
	}

	// Generate JWT token
	expiresAt := time.Now().Add(time.Hour * time.Duration(s.jwtExpiry))
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"role":      user.Role,
		"exp":       expiresAt.Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, apperrors.ErrInternal("failed to generate token")
	}

	return &LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
		User:      user,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) {
	if s.redis == nil {
		return
	}

	key := "jwt:blacklist:" + token
	s.redis.Client.Set(ctx, key, "1", 24*time.Hour)
}
