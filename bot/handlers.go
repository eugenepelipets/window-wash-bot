package bot

import (
	"fmt"
	"github.com/eugenepelipets/window-wash-bot/models"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	StateDefault                = ""
	StateWaitingForEntrance     = "waiting_for_entrance"
	StateWaitingForFloor        = "waiting_for_floor"
	StateWaitingForApartment    = "waiting_for_apartment"
	StateWindowsSameOrDifferent = "windows_same_or_different"
	StateWindowsSameType        = "windows_same_type"
	StateWindowsSameCount       = "windows_same_count"
	StateWindowsDifferent3      = "windows_diff_3"
	StateWindowsDifferent4      = "windows_diff_4"
	StateWindowsDifferent5      = "windows_diff_5"
	StateWindowsDifferent6_7    = "windows_diff_6_7"
	StateBalconyNeeded          = "balcony_needed"
	StateBalconyType            = "balcony_type"
	StateBalconySash            = "balcony_sash"
	StateTelegramNick           = "telegram_nick"
)

var userState = make(map[int64]string)
var userOrders = make(map[int64]models.Order)

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	switch {
	case data == "new_order":
		b.handleNewOrder(chatID)
	case data == "back":
		b.handleBack(chatID)
	case strings.HasPrefix(data, "entrance_"):
		entrance, _ := strconv.Atoi(data[len("entrance_"):])
		b.handleEntrance(chatID, entrance)
	case data == "windows_same" || data == "windows_different":
		b.handleWindowsSameOrDifferent(chatID, data == "windows_same")
	case strings.HasPrefix(data, "window_"):
		b.handleWindowTypeSelection(chatID, data)
	case strings.HasPrefix(data, "count_"):
		count, _ := strconv.Atoi(data[len("count_"):])
		if userSessions[chatID].Order.WindowsSame {
			b.handleWindowSameCount(chatID, count)
		} else {
			b.handleWindowDifferentCount(chatID, count)
		}
	case strings.HasPrefix(data, "balcony_"):
		log.Printf("Обработка balcony callback: %s", data)
		if count, err := strconv.Atoi(data[len("balcony_"):]); err == nil {
			b.handleBalconyNeeded(chatID, count)
		} else if data == "balcony_standard" || data == "balcony_floor" {
			b.handleBalconyType(chatID, data[len("balcony_"):])
		} else if strings.HasPrefix(data, "balcony_sash_") {
			b.handleBalconySash(chatID, data[len("balcony_sash_"):])
		}
	case data == "skip_nick":
		b.handleTelegramNick(chatID, "")
	case data == "confirm_order":
		b.handleOrderConfirmation(chatID)
	case data == "cancel_order":
		b.handleOrderCancellation(chatID)
	default:
		b.sendMessage(chatID, "Неизвестная команда")
	}

	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := b.api.Request(callbackConfig); err != nil {
		log.Printf("⚠️ Ошибка ответа на callback: %v", err)
	}
}

func (b *Bot) handleNewOrder(chatID int64) {
	// Сбрасываем предыдущий заказ
	userSessions[chatID] = &UserSession{
		Order: models.Order{
			UserID: chatID,
		},
	}
	b.updateState(chatID, StateWaitingForEntrance)
	b.sendMessage(chatID, "Выберите подъезд:", createEntranceKeyboard())
}

func (b *Bot) handleEntrance(chatID int64, entrance int) {
	session := b.getSession(chatID)
	session.Order.Entrance = entrance
	b.updateState(chatID, StateWaitingForFloor)
	b.sendMessage(chatID, "Введите номер этажа (1-24):")
}

func (b *Bot) handleWindowsSameOrDifferent(chatID int64, isSame bool) {
	session := b.getSession(chatID)
	session.Order.WindowsSame = isSame

	if isSame {
		b.updateState(chatID, StateWindowsSameType)
		b.sendMessage(chatID, "Выберите количество створок на окнах:", createWindowTypesKeyboard())
	} else {
		b.updateState(chatID, StateWindowsDifferent3)
		b.sendMessage(chatID, "Сколько 3-створчатых окон? (0-6)", createWindowCountKeyboard())
	}
}

func (b *Bot) handleWindowTypeSelection(chatID int64, sashType string) {
	session := b.getSession(chatID)
	session.Order.Window3Count = 0
	session.Order.Window4Count = 0
	session.Order.Window5Count = 0
	session.Order.Window6_7Count = 0

	switch sashType {
	case "window_3":
		session.Order.Window3Count = 1
	case "window_4":
		session.Order.Window4Count = 1
	case "window_5":
		session.Order.Window5Count = 1
	case "window_6_7":
		session.Order.Window6_7Count = 1
	}

	b.updateState(chatID, StateWindowsSameCount)
	b.sendMessage(chatID, "Сколько всего окон?", createWindowCountKeyboard())
}

func (b *Bot) handleWindowSameCount(chatID int64, count int) {
	session := b.getSession(chatID)

	if session.Order.Window3Count > 0 {
		session.Order.Window3Count = count
	} else if session.Order.Window4Count > 0 {
		session.Order.Window4Count = count
	} else if session.Order.Window5Count > 0 {
		session.Order.Window5Count = count
	} else if session.Order.Window6_7Count > 0 {
		session.Order.Window6_7Count = count
	}

	b.updateState(chatID, StateBalconyNeeded)
	b.sendMessage(chatID, "Нужно ли мыть окна на лоджии?", createBalconyNeededKeyboard())
}

func (b *Bot) handleWindowDifferentCount(chatID int64, count int) {
	session := b.getSession(chatID)

	switch session.CurrentState {
	case StateWindowsDifferent3:
		session.Order.Window3Count = count
		b.updateState(chatID, StateWindowsDifferent4)
		b.sendMessage(chatID, "Сколько 4-створчатых окон? (0-6)", createWindowCountKeyboard())

	case StateWindowsDifferent4:
		session.Order.Window4Count = count
		b.updateState(chatID, StateWindowsDifferent5)
		b.sendMessage(chatID, "Сколько 5-створчатых окон? (0-6)", createWindowCountKeyboard())

	case StateWindowsDifferent5:
		session.Order.Window5Count = count
		b.updateState(chatID, StateWindowsDifferent6_7)
		b.sendMessage(chatID, "Сколько 6-7-створчатых окон? (0-6)", createWindowCountKeyboard())

	case StateWindowsDifferent6_7:
		session.Order.Window6_7Count = count
		b.updateState(chatID, StateBalconyNeeded)
		b.sendMessage(chatID, "Нужно ли мыть окна на лоджии?", createBalconyNeededKeyboard())
	}
}

func (b *Bot) handleBalconyNeeded(chatID int64, count int) {
	session := b.getSession(chatID)
	session.Order.BalconyCount = count

	if count > 0 {
		// Сбрасываем предыдущие значения
		session.Order.BalconyType = ""
		session.Order.BalconySash = ""

		b.updateState(chatID, StateBalconyType)
		b.sendMessage(chatID, "Окна на лоджии стандартные или до пола?", createBalconyTypeKeyboard())
	} else {
		// Если лоджии не нужны, сразу переходим к нику
		b.updateState(chatID, StateTelegramNick)
		b.sendMessage(chatID, "Введите ваш ник в Telegram (или нажмите 'Пропустить'):", createSkipKeyboard())
	}
}

func (b *Bot) handleBalconyType(chatID int64, balconyType string) {
	session := b.getSession(chatID)
	session.Order.BalconyType = balconyType
	b.updateState(chatID, StateBalconySash)
	b.sendMessage(chatID, "Выберите количество створок на лоджии:", createBalconySashKeyboard())
}

func (b *Bot) handleBalconySash(chatID int64, sashType string) {
	session := b.getSession(chatID)
	session.Order.BalconySash = sashType
	b.updateState(chatID, StateTelegramNick)
	b.sendMessage(chatID, "Введите ваш ник в Telegram (или нажмите 'Пропустить'):", createSkipKeyboard())
}

func (b *Bot) handleTelegramNick(chatID int64, nick string) {
	session := b.getSession(chatID)
	session.Order.TelegramNick = nick

	price, err := CalculatePrice(session.Order)
	if err != nil {
		b.sendMessage(chatID, "Ошибка расчета стоимости. Пожалуйста, начните заново.")
		return
	}
	session.Order.Price = price

	b.showOrderConfirmation(chatID, session.Order)
}

func (b *Bot) showOrderConfirmation(chatID int64, order models.Order) {
	// Рассчитываем стоимость для каждого типа
	var window3Sum, window4Sum, window5Sum, window6_7Sum, balconySum int
	var details strings.Builder

	// Добавляем основную информацию
	details.WriteString(fmt.Sprintf("Подтвердите заказ:\n\nПодъезд: %d\nЭтаж: %d\nКвартира: %s\n\nОкна:\n",
		order.Entrance, order.Floor, order.Apartment))

	// Расчёт стоимости окон
	if order.Window3Count > 0 {
		window3Sum = order.Window3Count * 1000
		details.WriteString(fmt.Sprintf("- 3-створчатые: %d * 1000 = %d руб.\n", order.Window3Count, window3Sum))
	}
	if order.Window4Count > 0 {
		window4Sum = order.Window4Count * 1500
		details.WriteString(fmt.Sprintf("- 4-створчатые: %d * 1500 = %d руб.\n", order.Window4Count, window4Sum))
	}
	if order.Window5Count > 0 {
		window5Sum = order.Window5Count * 2000
		details.WriteString(fmt.Sprintf("- 5-створчатые: %d * 2000 = %d руб.\n", order.Window5Count, window5Sum))
	}
	if order.Window6_7Count > 0 {
		window6_7Sum = order.Window6_7Count * 2500
		details.WriteString(fmt.Sprintf("- 6-7-створчатые: %d * 2500 = %d руб.\n", order.Window6_7Count, window6_7Sum))
	}

	// Расчёт стоимости лоджий
	if order.BalconyCount > 0 {
		details.WriteString("\nЛоджии:\n")
		var balconyPrice int
		switch order.BalconySash {
		case "3":
			if order.BalconyType == "standard" {
				balconyPrice = 1000
			} else {
				balconyPrice = 1500
			}
		case "4":
			if order.BalconyType == "standard" {
				balconyPrice = 1500
			} else {
				balconyPrice = 2000
			}
		case "5":
			if order.BalconyType == "standard" {
				balconyPrice = 2000
			} else {
				balconyPrice = 2500
			}
		case "6_7":
			if order.BalconyType == "standard" {
				balconyPrice = 2500
			} else {
				balconyPrice = 3000
			}
		}

		balconySum = order.BalconyCount * balconyPrice
		balconyType := "стандартные"
		if order.BalconyType == "floor" {
			balconyType = "до пола"
		}

		details.WriteString(fmt.Sprintf("- %d лоджии (%s, %s створки): %d * %d = %d руб.\n",
			order.BalconyCount, balconyType, order.BalconySash,
			order.BalconyCount, balconyPrice, balconySum))
	}

	// Итоговая стоимость
	total := window3Sum + window4Sum + window5Sum + window6_7Sum + balconySum
	details.WriteString(fmt.Sprintf("\nИтого стоимость: %d руб.", total))

	b.updateState(chatID, "waiting_confirmation")
	b.sendMessage(chatID, details.String(), createConfirmationKeyboard())
}

func (b *Bot) handleOrderConfirmation(chatID int64) {
	session := b.getSession(chatID)
	order := session.Order
	order.Status = "confirmed"

	if err := b.db.SaveOrder(order); err != nil {
		b.sendMessage(chatID, "Ошибка сохранения заказа. Пожалуйста, попробуйте позже.")
		return
	}

	delete(userSessions, chatID)
	delete(userState, chatID)

	b.sendMessage(chatID, "Ваш заказ подтвержден! Ожидайте мастера.", createMainMenuKeyboard())
}

func (b *Bot) handleOrderCancellation(chatID int64) {
	delete(userSessions, chatID)
	delete(userState, chatID)
	b.sendMessage(chatID, "Заказ отменен.", createMainMenuKeyboard())
}

func (b *Bot) validateFloor(msg *tgbotapi.Message) {
	floor, err := strconv.Atoi(msg.Text)
	if err != nil || floor < 1 || floor > 24 {
		b.sendMessage(msg.Chat.ID, "Некорректный этаж. Введите цифру от 1 до 24:")
		return
	}

	session := b.getSession(msg.Chat.ID)
	session.Order.Floor = floor
	b.updateState(msg.Chat.ID, StateWaitingForApartment)
	b.sendMessage(msg.Chat.ID, "Введите номер квартиры (1-1500):")
}

func (b *Bot) validateApartment(msg *tgbotapi.Message) {
	if !IsDigitsOnly(msg.Text) {
		b.sendMessage(msg.Chat.ID, "Некорректный номер квартиры. Введите только цифры:")
		return
	}

	apartment, _ := strconv.Atoi(msg.Text)
	if apartment < 1 || apartment > 1500 {
		b.sendMessage(msg.Chat.ID, "Некорректный номер квартиры. Введите цифру от 1 до 1500:")
		return
	}

	session := b.getSession(msg.Chat.ID)
	session.Order.Apartment = msg.Text
	b.updateState(msg.Chat.ID, StateWindowsSameOrDifferent)
	b.sendMessage(msg.Chat.ID, "Количество створок на окнах одинаковое или разное?", createWindowsSameOrDifferentKeyboard())
}
