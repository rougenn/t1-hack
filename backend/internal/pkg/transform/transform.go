package transform

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/unidoc/unioffice/document"
)

// Преобразование файла PDF в текст с использованием pdfcpu
func PDFToText(filePath string) (string, error) {
	// Чтение PDF файла в байты
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Создаем конфигурацию для извлечения текста
	cfg := pdfcpu.NewDefaultConfiguration()

	// Открываем PDF файл
	ctx, err := pdfcpu.NewPDFContext(file)
	if err != nil {
		return "", err
	}

	// Извлекаем текст из всех страниц PDF
	var textBuilder strings.Builder
	for i := 1; i <= ctx.PageCount; i++ {
		text, err := pdfcpu.ExtractText(file, i, nil, cfg)
		if err != nil {
			return "", err
		}
		textBuilder.WriteString(text)
	}

	return textBuilder.String(), nil
}

// Преобразование файла Word (DOCX) в текст
func DOCXToText(filePath string) (string, error) {
	doc, err := document.Open(filePath)
	if err != nil {
		return "", err
	}
	var textContent strings.Builder
	for _, para := range doc.Paragraphs() {
		textContent.WriteString(para.Text())
	}
	return textContent.String(), nil
}

// Универсальная функция для преобразования файлов в текст
func FileToText(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	var text string
	var err error
	switch ext {
	case ".pdf":
		text, err = PDFToText(filePath)
	case ".docx":
		text, err = DOCXToText(filePath)
	// Добавьте обработку других форматов (например, .txt)
	default:
		text, err = ioutil.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("unsupported file format")
		}
	}
	return text, err
}
