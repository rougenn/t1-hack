# Используем официальный Go-образ
FROM golang:1.20 AS build

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем все файлы проекта в контейнер
COPY . .

# Скачиваем зависимости и собираем проект
RUN go mod tidy
RUN go build -o main cmd/app/main.go

# Второй этап: минимизируем размер финального образа
FROM debian:bullseye-slim

# Устанавливаем необходимые зависимости
RUN apt-get update && apt-get install -y libgcc1

# Копируем собранный Go-исполняемый файл из первого этапа
COPY --from=build /app/main /app/main

# Указываем рабочую директорию
WORKDIR /app

# Открываем порт, который будет использовать Go-сервер
EXPOSE 8090

# Запускаем приложение
CMD ["./main"]
