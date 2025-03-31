package main

import (
	"log"

	"github.com/eugenepelipets/window-wash-bot/bot"
	"github.com/eugenepelipets/window-wash-bot/storage"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Файл .env не найден, используем переменные окружения")
	}

	// Подключение к БД
	db, err := storage.NewPostgres()
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}
	defer db.Pool.Close()

	// Создание и запуск бота
	telegramBot, err := bot.NewBot(db)
	if err != nil {
		log.Fatalf("❌ Ошибка создания бота: %v", err)
	}

	log.Println("🚀 Бот успешно запущен!")
	telegramBot.Start()
}
