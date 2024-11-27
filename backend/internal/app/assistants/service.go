package assistants

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"t1/internal/app/models"
	"t1/internal/pkg/transform"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Функция для создания нового чат-ассистента
func CreateChatAssistant(DB *sql.DB, userID int, req models.AssistantRequest, ctx *gin.Context) (string, error) {
	// Генерация уникального ID для ассистента
	assistantID := uuid.New().String()

	// Получаем файлы из формы запроса (multipart.Form)
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Printf("Error getting multipart form: %v", err)
		return "", err
	}

	// Получаем файлы из формы
	files := form.File["files"]

	// Сохраняем файлы в папке /tmp
	tempDir := "/tmp/assistants/" + assistantID
	os.MkdirAll(tempDir, 0755)

	var textContent []string
	for _, fileHeader := range files {
		// Сохраняем файл на диск
		filePath := filepath.Join(tempDir, fileHeader.Filename)
		if err := ctx.SaveUploadedFile(fileHeader, filePath); err != nil {
			log.Printf("Error saving file: %v", err)
			return "", err
		}

		// Преобразуем файл в текстовый формат
		text, err := transform.FileToText(filePath)
		if err != nil {
			log.Printf("Error converting file: %v", err)
			return "", err
		}
		textContent = append(textContent, text)

		// Удаляем временный файл после обработки
		os.Remove(filePath)
	}

	// Добавляем ссылку для чат-ассистента в базу данных
	linkID, err := AddAssistantLink(DB, userID, assistantID, req.URL)
	if err != nil {
		log.Printf("Error adding link: %v", err)
		return "", err
	}

	// Записываем статистику запросов и ответов
	for _, text := range textContent {
		err = AddAssistantStatistics(DB, linkID, "User request", text)
		if err != nil {
			log.Printf("Error saving statistics: %v", err)
			return "", err
		}
	}

	return assistantID, nil
}
