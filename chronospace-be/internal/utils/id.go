package utils

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetUserIDFromContext(ctx *gin.Context) (pgtype.UUID, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return pgtype.UUID{}, errors.New("no authorization header")
	}

	if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
		return pgtype.UUID{}, errors.New("invalid authorization header format")
	}
	tokenString := authHeader[7:]

	// Parse the token and extract user ID
	claims, err := ParseToken(tokenString)
	if err != nil {
		return pgtype.UUID{}, err
	}

	// Assuming the ID is stored in claims as string
	var userID pgtype.UUID
	if err := userID.Scan(claims.ID); err != nil {
		return pgtype.UUID{}, errors.New("failed to parse user ID from token")
	}

	return userID, nil
}

func ParseUUID(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return pgtype.UUID{}, err
	}
	return uuid, nil
}
