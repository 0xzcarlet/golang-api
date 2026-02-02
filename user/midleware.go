package user

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware - middleware untuk protect routes
func AuthMiddleware() gin.HandlerFunc {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrResp{Error: "authorization header required"})
			c.Abort()
			return
		}

		// Format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, ErrResp{Error: "invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse & validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, ErrResp{Error: "invalid or expired token"})
			c.Abort()
			return
		}

		// Extract user ID from claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, ErrResp{Error: "invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64) // JWT numbers are float64
		if !ok {
			c.JSON(http.StatusUnauthorized, ErrResp{Error: "invalid user_id in token"})
			c.Abort()
			return
		}

		// Set userID di context untuk digunakan di handler
		c.Set("userID", uint64(userID))
		c.Next()
	}
}
