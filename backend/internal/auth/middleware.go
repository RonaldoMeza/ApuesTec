package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/RonaldoMeza/ApuesTec/backend/internal/response"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "authorization header required")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "INVALID_TOKEN", "invalid or expired access token")
			c.Abort()
			return
		}

		c.Set("userID", claims.Subject)
		c.Set("userEmail", claims.Email)
		c.Set("userRoles", claims.Roles)
		c.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRolesRaw, exists := c.Get("userRoles")
		if !exists {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "access denied")
			c.Abort()
			return
		}

		userRoles, ok := userRolesRaw.([]string)
		if !ok {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "access denied")
			c.Abort()
			return
		}

		for _, userRole := range userRoles {
			for _, allowed := range allowedRoles {
				if userRole == allowed {
					c.Next()
					return
				}
			}
		}

		response.Error(c, http.StatusForbidden, "FORBIDDEN", "insufficient permissions")
		c.Abort()
	}
}
