package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/smoreg/freezino/backend/internal/config"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// TokenType represents the type of JWT token
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims represents the JWT claims
type Claims struct {
	UserID uint      `json:"user_id"`
	Email  string    `json:"email"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations
type JWTManager struct {
	config *config.Config
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(cfg *config.Config) *JWTManager {
	return &JWTManager{config: cfg}
}

// GenerateAccessToken generates an access token for a user
func (jm *JWTManager) GenerateAccessToken(userID uint, email string) (string, error) {
	duration, err := time.ParseDuration(jm.config.JWTAccessExpiration)
	if err != nil {
		duration = 15 * time.Minute // default 15 minutes
	}

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Type:   AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.config.JWTSecret))
}

// GenerateRefreshToken generates a refresh token for a user
func (jm *JWTManager) GenerateRefreshToken(userID uint, email string) (string, error) {
	duration, err := time.ParseDuration(jm.config.JWTRefreshExpiration)
	if err != nil {
		duration = 7 * 24 * time.Hour // default 7 days
	}

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Type:   RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.config.JWTSecret))
}

// ValidateToken validates a JWT token and returns the claims
func (jm *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(jm.config.JWTSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GenerateTokenPair generates both access and refresh tokens
func (jm *JWTManager) GenerateTokenPair(userID uint, email string) (accessToken, refreshToken string, err error) {
	accessToken, err = jm.GenerateAccessToken(userID, email)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = jm.GenerateRefreshToken(userID, email)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
