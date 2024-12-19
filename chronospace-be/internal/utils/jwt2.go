package utils

import (
	"chronospace-be/internal/models"
	"context"
	"errors"
	"fmt"
	"time"

	"encoding/hex"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func GenerateTokens(ctx context.Context, userID pgtype.UUID, secretKey string) (models.Tokens, error) {
	// Set expiry times
	accessExpiry := time.Now().Add(15 * time.Minute)
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)

	// Create JWT claims for access token
	accessClaims := jwt.MapClaims{
		"user_id": hex.EncodeToString(userID.Bytes[:]),
		"exp":     accessExpiry.Unix(),
		"type":    "access",
	}

	// Create JWT claims for refresh token
	refreshClaims := jwt.MapClaims{
		"user_id": hex.EncodeToString(userID.Bytes[:]),
		"exp":     refreshExpiry.Unix(),
		"type":    "refresh",
	}

	// Sign tokens with secret key from config
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(secretKey))
	if err != nil {
		return models.Tokens{}, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretKey))
	if err != nil {
		return models.Tokens{}, fmt.Errorf("error generating refresh token: %w", err)
	}

	return models.Tokens{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
	}, nil
}

type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil // Replace with your actual secret key
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
