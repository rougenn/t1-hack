package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// в будущем брать его откуда то (конфиг думаю)
var jwtSecret = []byte("kljasdf;j;lasjdfhjkjk")

// GenerateAccessToken генерирует Access Token (действителен 3 минуты)
func GenerateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(), // Сохраняем UUID как строку
		"exp":     time.Now().Add(2 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken генерирует Refresh Token (действителен 7 дней)
func GenerateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(), // Сохраняем UUID как строку
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseAccessToken парсит Access Token и возвращает userID
func ParseAccessToken(tokenString string) (uuid.UUID, error) {
	return parseToken(tokenString)
}

// ParseRefreshToken парсит Refresh Token и возвращает userID
func ParseRefreshToken(tokenString string) (uuid.UUID, error) {
	return parseToken(tokenString)
}

// Общая функция парсинга токенов
func parseToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if ok && time.Unix(int64(exp), 0).Before(time.Now()) {
		return uuid.Nil, errors.New("token has expired")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("invalid user ID in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID format")
	}

	log.Printf("Parsed user_id from token: %s", userID.String())
	return userID, nil
}

func getUserIDFromContext(ctx *gin.Context) uuid.UUID {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		ctx.Abort()
		return uuid.Nil
	}
	return userID.(uuid.UUID)
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

		ctx.Set("user_id", userID) // Сохраняем userID в контексте как uuid.UUID
		ctx.Next()
	}
}
