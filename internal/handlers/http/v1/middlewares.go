package v1

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)
		c.Next()
	}
}

func SetLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := slog.Default().With("request_id", c.GetString("RequestID"))
		c.Set("Logger", logger)
		c.Next()
	}
}

func AuthMiddleware(JWTSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Authorization header is required")
			return
		}
		prefix := "Bearer "
		if len(accessToken) < len(prefix) || accessToken[:len(prefix)] != prefix {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Authorization header must start with 'Bearer '")
			return
		}
		accessToken = accessToken[len(prefix):]

		token, err := jwt.ParseWithClaims(accessToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSecret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid token")
			return
		}
		var userID uuid.UUID
		var login string
		if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
			userID, err = uuid.Parse((*claims)["sub"].(string))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid token")
				return
			}
			login = (*claims)["login"].(string)
		}
		c.Set("UserID", userID)
		c.Set("Login", login)
		c.Next()
	}
}

func SetUserInfoMiddleware(JWTSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		prefix := "Bearer "
		if len(accessToken) < len(prefix) || accessToken[:len(prefix)] != prefix {
			c.Next()
			return
		}
		accessToken = accessToken[len(prefix):]

		token, err := jwt.ParseWithClaims(accessToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSecret), nil
		})
		if err != nil {
			c.Next()
			return
		}
		var userID uuid.UUID
		var login string
		if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
			userID, err = uuid.Parse((*claims)["sub"].(string))
			if err != nil {
				c.Next()
				return
			}
			login = (*claims)["login"].(string)
		}
		c.Set("UserID", userID)
		c.Set("Login", login)
		c.Next()
	}
}
