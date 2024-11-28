package assistants

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"t1/internal/app/models"
	"t1/internal/pkg/transform"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid" // Подключаем для работы с UUID
)

// CreateChatAssistant создает нового чат-ассистента и добавляет его в таблицу assistant_links
func CreateChatAssistant(DB *sql.DB, adminID uuid.UUID, req models.AssistantRequest, ctx *gin.Context) (string, error) {
	// Генерация уникального ID для ассистента
	assistantID := uuid.New() // Новый UUID для ассистента
	log.Printf("Creating assistant with ID: %s", assistantID.String())

	// 2. Добавляем ссылку для ассистента в таблицу assistant_links
	linkID, err := AddAssistantLink(DB, adminID)
	if err != nil {
		log.Printf("Error adding link for assistant %s: %v", assistantID, err)
		return "", fmt.Errorf("failed to create assistant link: %w", err)
	}
	log.Printf("Successfully added link for assistant %s with ID %s", assistantID, linkID)

	// 3. Получаем файлы из формы запроса (multipart.Form)
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Printf("Error getting multipart form: %v", err)
		return "", fmt.Errorf("error getting multipart form: %w", err)
	}
	log.Printf("Received multipart form with %d files", len(form.File["files"]))

	// Получаем файлы из формы
	files := form.File["files"]

	// Создаем директорию для ассистента, где будут храниться только .txt файлы
	tempDir := "/tmp/assistants/" + assistantID.String() // Преобразуем UUID в строку для пути
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		log.Printf("Error creating temporary directory: %v", err)
		return "", fmt.Errorf("error creating temporary directory: %w", err)
	}
	log.Printf("Created temporary directory: %s", tempDir)

	// Преобразуем и сохраняем только .txt файлы
	var textContent []string
	for _, fileHeader := range files {
		// Проверяем расширение файла
		if filepath.Ext(fileHeader.Filename) != ".txt" {
			log.Printf("Skipping file %s (not a .txt file)", fileHeader.Filename)
			continue // Пропускаем файлы, которые не .txt
		}

		// Сохраняем файл на диск
		filePath := filepath.Join(tempDir, fileHeader.Filename)
		if err := ctx.SaveUploadedFile(fileHeader, filePath); err != nil {
			log.Printf("Error saving file %s: %v", fileHeader.Filename, err)
			return "", fmt.Errorf("error saving file %s: %w", fileHeader.Filename, err)
		}
		log.Printf("File %s saved to path %s", fileHeader.Filename, filePath)

		// Преобразуем файл в текстовый формат
		text, err := transform.FileToText(filePath)
		if err != nil {
			log.Printf("Error converting file %s to text: %v", fileHeader.Filename, err)
			return "", fmt.Errorf("error converting file %s to text: %w", fileHeader.Filename, err)
		}
		textContent = append(textContent, text)
		log.Printf("Successfully converted file %s to text", fileHeader.Filename)

		// Удаляем временный файл после обработки
		// if err := os.Remove(filePath); err != nil {
		// 	log.Printf("Error removing temporary file %s: %v", filePath, err)
		// } else {
		// 	log.Printf("Successfully removed temporary file %s", filePath)
		// }
	}

	// 4. Возвращаем URL ассистента
	return assistantID.String(), nil
}
