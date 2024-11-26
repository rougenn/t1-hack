package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your_secret_key")

// GenerateAccessToken генерирует Access Token (действителен 3 минуты)
func GenerateAccessToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(3 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken генерирует Refresh Token (действителен 7 дней)
func GenerateRefreshToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseAccessToken парсит Access Token и возвращает userID
func ParseAccessToken(tokenString string) (int, error) {
	return parseToken(tokenString)
}

// ParseRefreshToken парсит Refresh Token и возвращает userID
func ParseRefreshToken(tokenString string) (int, error) {
	return parseToken(tokenString)
}

// Общая функция парсинга токенов
func parseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if ok && time.Unix(int64(exp), 0).Before(time.Now()) {
		return 0, errors.New("token has expired")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID in token")
	}

	log.Printf("Parsed user_id from token: %d", int(userID))
	return int(userID), nil
}

func getUserIDFromContext(ctx *gin.Context) int {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		ctx.Abort()
		return 0
	}
	return userID.(int)
}

func (r *Server) RefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	userID, err := ParseRefreshToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	newAccessToken, err := GenerateAccessToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	newRefreshToken, err := GenerateRefreshToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			ctx.Abort()
			return
		}

		// Логируем значение заголовка Authorization
		log.Printf("Authorization Header: %s", authHeader)

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format invalid"})
			ctx.Abort()
			return
		}

		tokenString := authHeader[7:]
		userID, err := ParseAccessToken(tokenString)
		if err != nil {
			log.Printf("Error parsing access token: %s", err.Error()) // Логируем ошибку
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
