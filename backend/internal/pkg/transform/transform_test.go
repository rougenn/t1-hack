package transform

import (
    "testing"
    "fmt"
    "path/filepath"
)

func TestFileToText(t *testing.T) {
    // Гарантируем настройку лицензии
    InitLicense()

    files := []string{
        "sample.txt",               // Текстовый файл
        "styled_paragraph-2.pdf",   // PDF-файл
        "example.docx",             // DOCX-файл
        "unsupported.xlsx",         // Для проверки неподдерживаемого формата
    }

    for _, file := range files {
        t.Run(fmt.Sprintf("Testing file: %s", file), func(t *testing.T) {
            text, err := FileToText(file)
            if err != nil {
                if filepath.Ext(file) == ".xlsx" {
                    // Ожидаемая ошибка для неподдерживаемого формата
                    t.Logf("Ожидаемая ошибка для файла %s: %v", file, err)
                } else {
                    t.Errorf("Ошибка при обработке файла %s: %v", file, err)
                }
            } else {
                t.Logf("Извлеченный текст из %s:\n%s", file, text)
            }
        })
    }
}
