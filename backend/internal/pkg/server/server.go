package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"t1/internal/app/assistants"
	"t1/internal/app/models"
	"t1/internal/app/users"
	"t1/internal/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	// engine.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:5173"}, // Разрешить запросы с вашего фронтенда
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	engine.POST("/api/admin/login", r.LogIn)
	engine.POST("/api/admin/signup", r.Register)
	engine.POST("/api/admin/refresh-token", r.RefreshToken)

	protected := engine.Group("/")
	protected.Use(AuthMiddleware())

	protected.POST("/api/admin/create-chat-assistant", r.CreateChatAssistant) // реализация загрузки файлов
	protected.POST("/api/chats/send/:id", r.SendMessage)                      // запросы и сообщения в чате

	return engine
}

// если отправляем пустой файл питон серв падает !!!!!!!!!
func (r *Server) CreateChatAssistant(ctx *gin.Context) {
	// Получаем userID как uuid.UUID из контекста
	userID := getUserIDFromContext(ctx)
	if userID == uuid.Nil { // Проверка на nil UUID
		r.logger.Error("Failed to get user ID from context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	r.logger.Info("Received request to create assistant", zap.String("user_id", userID.String()))

	var req models.AssistantRequest
	// Логирование перед парсингом запроса
	r.logger.Info("Parsing request body")

	if err := ctx.ShouldBind(&req); err != nil {
		// Логирование ошибки при парсинге запроса
		r.logger.Error("Error parsing request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Логирование успешного парсинга запроса
	r.logger.Info("Request body parsed successfully", zap.Any("assistant_request", req))

	// Создаем ассистента в базе данных
	assistantID, err := assistants.CreateChatAssistant(r.DB, userID, req, ctx)
	if err != nil {
		// Логирование ошибки при создании ассистента
		r.logger.Error("Error creating assistant", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assistant"})
		return
	}

	// Логирование успешного создания ассистента
	r.logger.Info("Successfully created assistant", zap.String("assistant_id", assistantID)) // Логируем как строку UUID

	// Подготовка данных для отправки на Python-сервер
	trainModelRequest := map[string]string{
		"model_name":          fmt.Sprintf("assistant_%s", assistantID), // Используем String() для UUID
		"txt_files_directory": fmt.Sprintf("/tmp/assistants/%s", assistantID),
	}

	// Логирование данных для обучения модели
	r.logger.Info("Preparing data for training model", zap.Any("train_model_request", trainModelRequest))

	trainModelBody, err := json.Marshal(trainModelRequest)
	if err != nil {
		// Логирование ошибки при маршалинге
		r.logger.Error("Error marshalling train model request", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal train model request"})
		return
	}

	// Отправляем запрос на Python-сервер для обучения модели
	r.logger.Info("Sending request to Python server for model training")
	resp, err := http.Post("http://127.0.0.1:5000/train", "application/json", bytes.NewBuffer(trainModelBody))
	if err != nil {
		// Логирование ошибки при отправке запроса
		r.logger.Error("Error sending request to Python server", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to Python server"})
		return
	}
	defer resp.Body.Close()

	// Логирование статуса ответа от Python-сервера
	r.logger.Info("Received response from Python server", zap.Int("status_code", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		// Логирование ошибки, если статус код не OK
		r.logger.Error("Failed to train model", zap.Int("status_code", resp.StatusCode))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to train model"})
		return
	}

	// Успешно, возвращаем ID ассистента
	r.logger.Info("Assistant training completed successfully")
	ctx.JSON(http.StatusOK, gin.H{"assistant_id": assistantID}) // Отправляем как строку UUID
}

// SendMessage обрабатывает сообщения от пользователя и отправляет запрос на Python-сервер
func (r *Server) SendMessage(ctx *gin.Context) {
	// Получаем userID как uuid.UUID из контекста
	userID := getUserIDFromContext(ctx)
	if userID == uuid.Nil { // Проверка на nil UUID
		r.logger.Error("Failed to get user ID from context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Извлекаем assistantID из URL
	assistantIDStr := ctx.Param("id")
	if assistantIDStr == "" {
		r.logger.Error("Assistant ID is required in the URL")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Assistant ID is required"})
		return
	}

	// Преобразуем assistantID в UUID
	assistantID, err := uuid.Parse(assistantIDStr)
	if err != nil {
		r.logger.Error("Invalid Assistant ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Assistant ID"})
		return
	}

	var req models.SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		r.logger.Error("Error parsing send message request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Формируем запрос для отправки на Python-сервер
	questionRequest := map[string]string{
		"assistant_id": assistantID.String(), // Передаем UUID ассистента
		"message":      req.Message,          // Сообщение от пользователя
	}

	questionBody, err := json.Marshal(questionRequest)
	if err != nil {
		r.logger.Error("Error marshalling question request", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal question request"})
		return
	}

	// Отправляем запрос на Python-сервер для получения ответа
	resp, err := http.Post("http://127.0.0.1:5000/ask", "application/json", bytes.NewBuffer(questionBody))
	if err != nil {
		r.logger.Error("Error sending request to Python server", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to Python server"})
		return
	}
	defer resp.Body.Close()

	// Чтение ответа от Python-сервера
	var answerResponse map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&answerResponse); err != nil {
		r.logger.Error("Error decoding response from Python server", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response from Python server"})
		return
	}

	// Отправляем ответ на фронт
	ctx.JSON(http.StatusOK, gin.H{"message": answerResponse["message"]})
}

func (r *Server) GetChat(ctx *gin.Context) {}

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
