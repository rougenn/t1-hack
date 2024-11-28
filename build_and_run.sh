#!/bin/bash

# Сборка фронтенда
echo "Сборка фронтенда..."
cd frontend
npm install
npm run build
cp -r dist/* ../backend/static/
cd ..

# Сборка Go-бинарника
echo "Сборка Go приложения..."
cd backend
go build -o app cmd/app/main.go
cd ..

# Запуск Go-сервера
echo "Запуск Go сервера..."
./backend/app &

# Запуск Python-сервера
echo "Запуск Python сервера..."
cd ml
python3 server.py &
cd ..

echo "Процесс завершён. Серверы запущены."
