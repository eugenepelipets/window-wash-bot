package bot

import (
	"log"
	"os"

	"github.com/eugenepelipets/window-wash-bot/models"
	"github.com/eugenepelipets/window-wash-bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbotapi.BotAPI
	db  *storage.Postgres
}

// Создаем бота
func NewBot(db *storage.Postgres) (*Bot, error) {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("❌ Переменная TELEGRAM_TOKEN не задана! Проверь .env")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = true // Показываем запросы к API в логах

	log.Printf("✅ Бот авторизован как %s", bot.Self.UserName)

	return &Bot{api: bot, db: db}, nil
}

// Запуск бота
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			switch userState[update.Message.Chat.ID] {
			case "waiting_for_floor":
				b.validateFloor(update.Message)
			case "waiting_for_apartment":
				b.validateApartment(update.Message)
			default:
				b.handleMessage(update.Message)
			}
		} else if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
		}
	}
}

// Обработка сообщений
func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	switch msg.Text {
	case "/start":
		b.handleStart(msg)
	default:
		b.sendMessage(msg.Chat.ID, "Я не понимаю эту команду 🤔")
	}
}

// Обработка команды /start
func (b *Bot) handleStart(msg *tgbotapi.Message) {
	user := models.User{
		TelegramID: msg.Chat.ID,
		UserName:   msg.From.UserName,
		FirstName:  msg.From.FirstName,
		LastName:   msg.From.LastName,
	}

	err := b.db.SaveUser(user)
	if err != nil {
		log.Printf("⚠️ Ошибка сохранения пользователя: %v", err)
	}

	// Отправляем приветственное сообщение с кнопкой
	replyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Новый заказ", "new_order"),
		),
	)
	b.sendMessage(msg.Chat.ID, "Привет! Я помогу тебе записаться на мытье окон. Нажми кнопку ниже 👇", replyMarkup)
}

func (b *Bot) sendMessage(chatID int64, text string, replyMarkup ...tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	if len(replyMarkup) > 0 {
		msg.ReplyMarkup = replyMarkup[0]
	}
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("⚠️ Ошибка отправки сообщения: %v", err)
	}
}
