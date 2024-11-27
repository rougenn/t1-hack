package server

import (
	"database/sql"
	"errors"
	"net/http"
	"t1/internal/app/assistants"
	"t1/internal/app/models"
	"t1/internal/app/users"
	"t1/internal/pkg/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	ErrIncorrectData = errors.New("password and/or username are incorrect")
	ErrAlreadyExists = errors.New("user with this email or phone already exists")
	ErrNotFound      = errors.New("user not found")
)

type Server struct {
	host   string
	DB     *sql.DB
	logger *zap.Logger // Логгер теперь является частью структуры Server
}

// Новый конструктор сервера, который создает логгер
func New(host string) *Server {
	database := db.NewDB()

	// Создаем новый логгер
	logger, err := zap.NewProduction() // Можно также использовать zap.NewDevelopment() для менее формального логирования
	if err != nil {
		panic("failed to initialize zap logger")
	}

	s := Server{
		host:   host,
		DB:     database,
		logger: logger, // Передаем логгер в структуру
	}
	return &s
}

func (r *Server) Stop() {
	r.DB.Close()
	// Закрытие логгера
	if err := r.logger.Sync(); err != nil {
		r.logger.Error("Error syncing logger", zap.Error(err))
	}
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

	engine.POST("/api/admin/login", r.LogIn)
	engine.POST("/api/admin/signup", r.Register)
	engine.POST("/api/admin/refresh-token", r.RefreshToken)

	engine.POST("/api/admin/create-chat-assistant", r.CreateChatAssistant) // реализация загрузки файлов
	engine.POST("/api/get-chat/:id", r.GetChat)                            // подгрузка страницы чата
	engine.POST("/api/chats/:id/send", r.SendMessage)                      // запросы и сообщения в чате

	protected := engine.Group("/api/")
	protected.Use(AuthMiddleware())

	return engine
}

func (r *Server) CreateChatAssistant(ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)

	var req models.AssistantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// Логируем ошибку с помощью zap
		r.logger.Error("Error parsing request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Передаем контекст ctx как аргумент в функцию для создания ассистента
	assistantID, err := assistants.CreateChatAssistant(r.DB, userID, req, ctx)
	if err != nil {
		// Логируем ошибку с помощью zap
		r.logger.Error("Error creating assistant", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assistant"})
		return
	}

	// Временно закомментировали отправку в Kafka
	// Отправляем запрос в Kafka
	/*
		err = kafka.SendToKafka("chat-assistants", assistantID, nil)
		if err != nil {
			// Логируем ошибку с помощью zap
			r.logger.Error("Error sending to Kafka", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send to Kafka"})
			return
		}
	*/

	// Временно отправляем ответ напрямую, без отправки в Kafka
	ctx.JSON(http.StatusOK, gin.H{"assistant_id": assistantID})
}

func (r *Server) GetChat(ctx *gin.Context) {}

func (r *Server) SendMessage(ctx *gin.Context) {}

func (r *Server) Start() {
	err := r.newAPI().Run(r.host)
	if err != nil {
		r.logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func (r *Server) LogIn(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.logger.Error("Error binding login request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := users.SignIn(r.DB, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, users.ErrIncorrectData) {
			r.logger.Warn("Incorrect login attempt", zap.String("email", req.Email))
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		r.logger.Error("Error signing in", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Генерация Access и Refresh токенов
	accessToken, err := GenerateAccessToken(user.ID)
	if err != nil {
		r.logger.Error("Failed to generate access token", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := GenerateRefreshToken(user.ID)
	if err != nil {
		r.logger.Error("Failed to generate refresh token", zap.Error(err))
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
		r.logger.Error("Error binding registration request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := users.Register(r.DB, req)
	if err != nil {
		if errors.Is(err, users.ErrAlreadyExists) {
			r.logger.Warn("User already exists", zap.String("email", req.Email))
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		r.logger.Error("Error registering user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}
