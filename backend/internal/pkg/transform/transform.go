package transform

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/document/convert"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// Инициализация лицензии UniDoc
var licenseInitialized = false

func InitLicense() {
	if licenseInitialized {
		return
	}

	err := license.SetMeteredKey("a6e0d818aba9dc2726be5a434cdaccfd2312692f72346561a37cca63e750f2e9")
	if err != nil {
		log.Fatalf("Failed to set UniPDF license key: %v", err)
	}

	log.Println("UniPDF license key set successfully.")
	licenseInitialized = true
}

// Функция для извлечения текста из PDF файла
func PDftotext(filePath string) (string, error) {
	// Open the PDF file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	// Create a PDF reader
	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		return "", fmt.Errorf("Error creating PDF reader: %v", err)
	}

	// Get the number of pages in the PDF
	numOfPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", fmt.Errorf("Error getting number of pages: %v", err)
	}

	var textBuilder strings.Builder

	// Loop through each page and extract text
	for i := 0; i < numOfPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return "", fmt.Errorf("Error getting page %d: %v", pageNum, err)
		}

		textExtractor, err := extractor.New(page)
		if err != nil {
			return "", fmt.Errorf("Error creating text extractor for page %d: %v", pageNum, err)
		}

		pageText, err := textExtractor.ExtractText()
		if err != nil {
			return "", fmt.Errorf("Error extracting text from page %d: %v", pageNum, err)
		}

		textBuilder.WriteString(pageText)
		textBuilder.WriteString("\n")
	}

	return textBuilder.String(), nil
}

// Функция для конвертации DOCX в PDF
func DOCXToPDF(inputFilePath string, outputFilePath string) error {
	doc, err := document.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("error opening document: %s", err)
	}
	defer doc.Close()

	c := convert.ConvertToPdf(doc)
	err = c.WriteToFile(outputFilePath)
	if err != nil {
		return fmt.Errorf("error converting document: %s", err)
	}

	return nil
}

// Главная функция, которая обрабатывает файлы в зависимости от их расширения
func FileToText(filePath string) (string, error) {
	InitLicense()

	// Проверяем, существует ли файл
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("Файл не найден: %s", filePath)
	}

	if info.IsDir() {
		return "", fmt.Errorf("Указанный путь является директорией, а не файлом: %s", filePath)
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".txt":
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("Ошибка при чтении файла: %v", err)
		}
		return string(content), nil

	case ".pdf":
		text, err := PDftotext(filePath)
		if err != nil {
			return "", fmt.Errorf("Ошибка при извлечении текста из PDF: %v", err)
		}
		return text, nil

	case ".docx":
		tempDir := os.TempDir()
		tempPDFPath := filepath.Join(tempDir, strings.TrimSuffix(filepath.Base(filePath), ".docx")+".pdf")

		err := DOCXToPDF(filePath, tempPDFPath)
		if err != nil {
			return "", fmt.Errorf("Ошибка при конвертации DOCX в PDF: %v", err)
		}
		defer os.Remove(tempPDFPath) // Удаляем временный файл после использования

		text, err := PDftotext(tempPDFPath)
		if err != nil {
			return "", fmt.Errorf("Ошибка при извлечении текста из PDF: %v", err)
		}
		return text, nil

	default:
		return "", fmt.Errorf("Неподдерживаемый тип файла: %s", ext)
	}
}
