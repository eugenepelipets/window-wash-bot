package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eugenepelipets/window-wash-bot/bot"
	"github.com/eugenepelipets/window-wash-bot/storage"
	"github.com/joho/godotenv"
)

func main() {
	// Настройка логгирования
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
	log.Println("🚀 Запуск бота...")

	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️ Не удалось загрузить .env файл: %v", err)
	}

	// Подключение к БД
	db, err := storage.NewPostgres()
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}
	defer db.Pool.Close()
	log.Println("✅ Подключение к БД установлено")

	// Создание бота
	telegramBot, err := bot.NewBot(db)
	if err != nil {
		log.Fatalf("❌ Ошибка создания бота: %v", err)
	}
	log.Println("✅ Бот инициализирован")

	// Канал для обработки сигналов завершения
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Запуск бота в отдельной горутине
	go func() {
		log.Println("🤖 Бот начал обработку сообщений...")
		telegramBot.Start()
	}()

	// Ожидание сигнала завершения
	<-done
	log.Println("🛑 Получен сигнал завершения, останавливаем бота...")

	// Здесь можно добавить graceful shutdown логику
	log.Println("✅ Бот успешно остановлен")
}
