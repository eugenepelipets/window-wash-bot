package bot

import (
	"github.com/eugenepelipets/window-wash-bot/models"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	StateDefault             = ""
	StateWaitingForFloor     = "waiting_for_floor"
	StateWaitingForApartment = "waiting_for_apartment"
)

var userState = make(map[int64]string)
var userOrders = make(map[int64]models.Order)

// Обработка callback-запросов от кнопок
func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	switch data {
	case "new_order":
		b.handleNewOrder(chatID)
	case "confirm_order":
		b.handleConfirmOrder(chatID)
	default:
		if len(data) > 7 && data[:7] == "window_" {
			b.handleWindowSelection(chatID, data[7:])
		}
	}

	// Ответ на callback, чтобы убрать "часики" у кнопки
	// Заменяем проблемный участок в handlers.go
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := b.api.Request(callbackConfig); err != nil {
		log.Printf("⚠️ Ошибка ответа на callback: %v", err)
	}
}

// Обработка нового заказа
func (b *Bot) handleNewOrder(chatID int64) {
	keyboard := createWindowTypesKeyboard()
	b.sendMessage(chatID, "Выберите тип окон:", keyboard)
}

// Обработка выбора типа окон
func (b *Bot) handleWindowSelection(chatID int64, windowType string) {
	userState[chatID] = "waiting_for_floor"
	// Сохраняем тип окон во временную структуру
	userOrders[chatID] = models.Order{
		UserID:     chatID,
		WindowType: windowType,
	}
	b.sendMessage(chatID, "Введите этаж (только цифры):")
}

// Подтверждение заказа
func (b *Bot) handleConfirmOrder(chatID int64) {
	// Получаем сохраненные данные
	order := userOrders[chatID]

	// Рассчитываем цену
	price, _ := CalculatePrice(order.WindowType, order.Floor)
	order.Price = price
	order.Status = "confirmed"

	// Сохраняем в БД
	err := b.db.SaveOrder(order)
	if err != nil {
		log.Printf("⚠️ Ошибка сохранения заказа: %v", err)
		b.sendMessage(chatID, "Произошла ошибка при сохранении заказа. Попробуйте позже.")
		return
	}

	// Удаляем временные данные
	delete(userOrders, chatID)
	delete(userState, chatID)

	b.sendMessage(chatID, "Ваш заказ подтвержден! Ожидайте мастера.")
}

// Валидация этажа
func (b *Bot) validateFloor(msg *tgbotapi.Message) {
	floor, err := strconv.Atoi(msg.Text)
	if err != nil || floor < 1 || floor > 100 {
		b.sendMessage(msg.Chat.ID, "Некорректный этаж. Введите цифру от 1 до 100:")
		return
	}

	order := userOrders[msg.Chat.ID]
	order.Floor = floor
	userOrders[msg.Chat.ID] = order

	userState[msg.Chat.ID] = "waiting_for_apartment"
	b.sendMessage(msg.Chat.ID, "Введите номер квартиры:")
}

// Валидация номера квартиры
func (b *Bot) validateApartment(msg *tgbotapi.Message) {
	if len(msg.Text) < 1 || len(msg.Text) > 10 {
		b.sendMessage(msg.Chat.ID, "Некорректный номер квартиры. Введите снова:")
		return
	}

	// Обновляем заказ
	order := userOrders[msg.Chat.ID]
	order.Apartment = msg.Text
	userOrders[msg.Chat.ID] = order

	userState[msg.Chat.ID] = ""
	b.sendConfirmation(msg.Chat.ID)
}

// Отправка подтверждения заказа
func (b *Bot) sendConfirmation(chatID int64) {
	keyboard := createConfirmationKeyboard()
	b.sendMessage(chatID, "Подтвердите заказ:", keyboard)
}

// Создание клавиатуры с типами окон
func createWindowTypesKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Обычные окна", "window_regular"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Панорамные окна", "window_panoramic"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Витрины", "window_shop"),
		),
	)
}

// Создание клавиатуры подтверждения
func createConfirmationKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвердить", "confirm_order"),
			tgbotapi.NewInlineKeyboardButtonData("Отменить", "cancel_order"),
		),
	)
}
