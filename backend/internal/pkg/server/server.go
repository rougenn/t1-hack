package server

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"t1/internal/app/models" 
	"t1/internal/app/users"
	"github.com/google/uuid"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	_ "time"
	"t1/internal/pkg/db"
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
		AllowOrigins:     []string{"http://localhost:5173"}, 
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
<<<<<<< Updated upstream
	engine.POST("/api/admin/create-manager", r.CreateManager) // надо реализовать загрузку файла\файлов в запросе
	// + генерацию ссылки + генерацию уникального айди(его можно в ссылке использовать)

	engine.POST("/api/get-chat/:id", r.GetChat) // просто должны подгружать страничку

	engine.POST("/api/chats/:id/send", r.SendMessage) // тут тебе просто кидают запрос с фронта уже в чате

=======
	engine.POST("/api/user/create-manager", r.CreateManager)
	engine.POST("/api/admin/create-manager", r.CreateManager) 
	// надо реализовать загрузку файла\файлов в запросе
	// + генерацию ссылки + генерацию уникального айди(его можно в ссылке использовать)
	engine.GET("/api/get-chat/:id", r.GetChat) // исправлено на GET
	engine.POST("/api/chats/:id/send", r.SendMessage) // тут тебе просто кидают запрос с фронта уже в чате
	engine.POST("/api/chats", r.CreateChat)
>>>>>>> Stashed changes
	protected := engine.Group("/api/")
	protected.Use(AuthMiddleware())

	return engine
}

func (r *Server) CreateChat(ctx *gin.Context) {
    var req models.ChatRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    chatID, err := db.CreateChatInDB(r.DB, req.UserID, req.Title)
    if err != nil {
        log.Println("Error in CreateChatInDB:", err) // Логируем ошибку
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{
        "message": "Chat created successfully",
        "chat_id": chatID,
    })
}



func (r *Server) GetChat(ctx *gin.Context) {
	// 1. Получаем ID чата из параметра пути
	chatID := ctx.Param("id")
	if chatID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Chat ID is required"})
		return
	}

	// 2. Заглушка для получения данных чата
	// Здесь в будущем можно добавить запрос в базу данных. Пока используем фиктивные данные.
	chatData := map[string]interface{}{
		"chat_id": chatID,
		"messages": []map[string]interface{}{
			{
				"sender":   "user",
				"message":  "Hello! How can I help you?",
				"sent_at":  "2024-11-27T10:00:00Z",
			},
			{
				"sender":   "assistant",
				"message":  "I need help with my account.",
				"sent_at":  "2024-11-27T10:01:00Z",
			},
		},
	}

	// 3. Возвращаем данные чата
	ctx.JSON(http.StatusOK, gin.H{
		"chat": chatData,
	})
}





type SendMessageRequest struct {
	Sender  string `json:"sender" binding:"required"`  // Кто отправил сообщение: "user" или "assistant"
	Message string `json:"message" binding:"required"` // Текст сообщения
}

func (r *Server) SendMessage(ctx *gin.Context) {
	// 1. Получаем ID чата из параметра пути
	chatID := ctx.Param("id")

	// 2. Получаем сообщение из тела запроса
	var req struct {
		Message string `json:"message"`
	}

	// 3. Проверяем, что тело запроса правильно
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 4. Проверяем, что ID чата и сообщение не пустые
	if chatID == "" || req.Message == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Chat ID and message are required"})
		return
	}

	// Тут можно добавить логику сохранения сообщения в базе данных

	// 5. Возвращаем успешный ответ
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Message sent successfully",
		"chat_id": chatID,
		"message_content": req.Message,
	})
}



func (r *Server) CreateManager(ctx *gin.Context) {
	// 1. Получаем файл из запроса
	file, header, err := ctx.Request.FormFile("file") // Ключ "file" должен совпадать с Postman
	if err != nil {
		log.Printf("Failed to get file from request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}
	defer file.Close()

	// 2. Генерация уникального ID
	uniqueID := uuid.New().String()

	// 3. Формирование пути для сохранения файла
	filename := uniqueID + filepath.Ext(header.Filename)
	savePath := "/Users/chrizantona/t1-hack/backend/uploads/" + filename

	// Логируем путь, чтобы проверить, куда сохраняется файл
	log.Println("Saving file to:", savePath)

	// 4. Сохранение файла
	out, err := os.Create(savePath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	defer out.Close()

	// Копируем содержимое файла в новый файл
	if _, err := io.Copy(out, file); err != nil {
		log.Printf("Failed to write file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
		return
	}

	// 5. Генерация ссылки на сохранённый файл
	fileURL := "http://localhost:8090/uploads/" + filename

	// 6. Ответ клиенту
	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Manager created successfully",
		"file_url": fileURL,
		"file_id":  uniqueID,
	})
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
