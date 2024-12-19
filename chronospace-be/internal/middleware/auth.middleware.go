package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type JWTConfig struct {
	SecretKey string
}

func NewJWTMiddleware(secretKey string) *JWTConfig {
	return &JWTConfig{
		SecretKey: secretKey,
	}
}

func (j *JWTConfig) ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(j.SecretKey), nil
		})

		if err != nil {
			switch {
			case err == jwt.ErrSignatureInvalid:
				c.JSON(401, gin.H{"error": "Invalid token signature"})
			case strings.Contains(err.Error(), "token is expired"):
				c.JSON(401, gin.H{"error": "Token has expired"})
			case strings.Contains(err.Error(), "token is not valid"):
				c.JSON(401, gin.H{"error": "Token structure is invalid"})
			default:
				c.JSON(401, gin.H{"error": "Invalid token: " + err.Error()})
			}
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

// Helper function to get claims from gin context
func GetClaims(c *gin.Context) (jwt.MapClaims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	mapClaims, ok := claims.(jwt.MapClaims)
	return mapClaims, ok
}

