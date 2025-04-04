package bot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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
			// Обрабатываем текстовые сообщения
			if update.Message.IsCommand() {
				b.handleMessage(update.Message)
			} else {
				// Проверяем текущее состояние пользователя
				session := b.getSession(update.Message.Chat.ID)
				switch session.CurrentState {
				case StateWaitingForFloor, StateWaitingForApartment, StateTelegramNick:
					b.handleTextMessage(update.Message)
				default:
					b.sendMessage(update.Message.Chat.ID, "Пожалуйста, используйте кнопки для продолжения.")
				}
			}
		} else if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
		}
	}
}

// Обработка сообщений
func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	switch {
	case msg.Text == "/start":
		b.handleStart(msg)
	case strings.HasPrefix(msg.Text, "/export"):
		b.handleExport(msg)
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
	b.sendMessage(msg.Chat.ID, "Привет! Чтобы записаться на мойку окон нажми на кнопку 👇", replyMarkup)
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

func (b *Bot) sendMainMenu(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Главное меню:")
	msg.ReplyMarkup = createMainMenuKeyboard()
	b.api.Send(msg)
}

func (b *Bot) sendEntranceKeyboard(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите подъезд:")
	msg.ReplyMarkup = createEntranceKeyboard()
	b.api.Send(msg)
}

func (b *Bot) notifyAdminAboutDuplicate(chatID int64, order models.Order) {
	adminID, _ := strconv.ParseInt(os.Getenv("ADMIN_TELEGRAM_ID"), 10, 64)
	if adminID == 0 {
		return
	}

	msgText := fmt.Sprintf(
		"⚠️ Обнаружен дублирующий заказ!\n\n"+
			"Подъезд: %d\nЭтаж: %d\nКвартира: %s\n"+
			"Пользователь: @%s (%d)\n\n"+
			"Первый заказ будет подтвержден, этот - на уточнении.",
		order.Entrance, order.Floor, order.Apartment,
		order.User.UserName, order.User.TelegramID)

	msg := tgbotapi.NewMessage(adminID, msgText)
	b.api.Send(msg)
}
