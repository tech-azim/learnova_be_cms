package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type ClaimStruct struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Fungsi helper untuk get JWT Secret
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable not set")
	}
	return []byte(secret)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)

		jwtSecret := getJWTSecret()

		token, err := jwt.ParseWithClaims(tokenString, &ClaimStruct{}, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			fmt.Println("Signing method:", t.Method.Alg())
			return jwtSecret, nil
		})

		if err != nil {
			fmt.Println("Parse error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  "Invalid or expired token",
				"detail": err.Error(),
			})
			c.Abort()
			return
		}

		if !token.Valid {
			fmt.Println("Token not valid")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*ClaimStruct)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Error invalid token claims",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}