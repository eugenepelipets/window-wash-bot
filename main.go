package main

import (
	"log"

	"github.com/eugenepelipets/window-wash-bot/bot"
	"github.com/eugenepelipets/window-wash-bot/storage"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️ Warning: %v", err)
	}

	db, err := storage.NewPostgres()
	if err != nil {
		log.Fatalf("❌ Database connection error: %v", err)
	}
	defer db.Pool.Close() // Упрощаем, так как Close() не возвращает ошибку

	telegramBot, err := bot.NewBot(db)
	if err != nil {
		log.Fatalf("❌ Bot creation error: %v", err)
	}

	log.Println("🚀 Bot started successfully")
	telegramBot.Start()
}
