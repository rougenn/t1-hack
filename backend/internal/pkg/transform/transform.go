package transform

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Преобразование файла PDF в текст с использованием pdfcpu
func PDFToText(filePath string) (string, error) {
	// Реализуйте преобразование PDF в текст
	// Например, с помощью pdfcpu
	return "", nil
}

// Преобразование файла Word (DOCX) в текст
func DOCXToText(filePath string) (string, error) {
	// Реализуйте преобразование DOCX в текст
	// Например, с помощью unioffice
	return "", nil
}

// Универсальная функция для преобразования файлов в текст
func FileToText(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	var text string
	var err error

	switch ext {
	case ".txt":
		// Если это .txt, просто читаем и сохраняем файл
		text, err = saveTxtFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to save .txt file: %v", err)
		}
	case ".pdf":
		text, err = PDFToText(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to convert PDF to text: %v", err)
		}
	case ".docx":
		text, err = DOCXToText(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to convert DOCX to text: %v", err)
		}
	default:
		return "", fmt.Errorf("unsupported file format: %s", ext)
	}
	return text, err
}

// Функция для сохранения .txt файла в новую директорию
func saveTxtFile(filePath string) (string, error) {
	// Читаем содержимое .txt файла
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Определим, куда сохранять файл. Например, в папку "output"
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		// Создаем директорию, если она не существует
		err := os.Mkdir(outputDir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	// Определяем имя нового файла
	newFileName := filepath.Join(outputDir, filepath.Base(filePath))

	// Сохраняем файл в новую директорию
	err = ioutil.WriteFile(newFileName, content, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write .txt file: %v", err)
	}

	// Возвращаем путь сохраненного файла
	return newFileName, nil
}
