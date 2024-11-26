package server

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"t1/internal/app/models"
	"t1/internal/app/users"
	"t1/internal/pkg/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	ErrIncorrectData = errors.New("password and/or username are incorrect")
	ErrAlreadyExists = errors.New("user with this email or phone already exists")
	ErrNotFound      = errors.New("user not found")
)

type Server struct {
	host string
	DB   *sql.DB
}

func New(host string) *Server {
	database := db.NewDB()

	s := Server{
		host: host,
		DB:   database,
	}
	return &s
}

func (r *Server) Stop() {
	r.DB.Close()
}

func (r *Server) newAPI() *gin.Engine {
	engine := gin.New()

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Разрешить запросы с вашего фронтенда
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	engine.POST("/user/login", r.LogIn)
	engine.POST("/user/signup", r.Register)
	engine.POST("/user/refresh-token", r.RefreshToken)

	protected := engine.Group("/")
	protected.Use(AuthMiddleware())

	return engine
}

func (r *Server) Start() {
	err := r.newAPI().Run(r.host)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (r *Server) LogIn(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := users.SignIn(r.DB, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, users.ErrIncorrectData) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Генерация Access и Refresh токенов
	accessToken, err := GenerateAccessToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := GenerateRefreshToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (r *Server) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := users.Register(r.DB, req)
	if err != nil {
		if errors.Is(err, users.ErrAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}
