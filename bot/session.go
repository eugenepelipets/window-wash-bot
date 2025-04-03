package bot

import (
	"github.com/eugenepelipets/window-wash-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type UserSession struct {
	CurrentState   string
	PreviousStates []string // История состояний для реализации "Назад"
	Order          models.Order
	TempData       map[string]interface{} // Для временных данных
}

var userSessions = make(map[int64]*UserSession)

func (b *Bot) getSession(chatID int64) *UserSession {
	if _, ok := userSessions[chatID]; !ok {
		userSessions[chatID] = &UserSession{
			PreviousStates: make([]string, 0),
			TempData:       make(map[string]interface{}),
		}
	}
	return userSessions[chatID]
}

func (b *Bot) updateState(chatID int64, newState string) {
	session := b.getSession(chatID)
	if session.CurrentState != "" {
		session.PreviousStates = append(session.PreviousStates, session.CurrentState)
	}
	session.CurrentState = newState
}

func (b *Bot) handleBack(chatID int64) {
	session := b.getSession(chatID)
	if len(session.PreviousStates) == 0 {
		b.sendMainMenu(chatID)
		return
	}

	// Извлекаем предыдущее состояние
	prevState := session.PreviousStates[len(session.PreviousStates)-1]
	session.PreviousStates = session.PreviousStates[:len(session.PreviousStates)-1]
	session.CurrentState = prevState

	// Восстанавливаем предыдущий шаг
	b.restorePreviousStep(chatID, prevState)
}

func (b *Bot) restorePreviousStep(chatID int64, state string) {
	switch state {
	case StateWaitingForEntrance:
		b.sendMessage(chatID, "Выберите подъезд:", createEntranceKeyboard())
	case StateWaitingForFloor:
		b.sendMessage(chatID, "Введите этаж (1-24):")
	case StateWaitingForApartment:
		b.sendMessage(chatID, "Введите номер квартиры:")
	case StateWindowsSameOrDifferent:
		b.sendMessage(chatID, "Количество створок на окнах одинаковое или разное?",
			createWindowsSameOrDifferentKeyboard())
	case StateWindowsSameType:
		b.sendMessage(chatID, "Выберите количество створок на окнах:",
			createWindowTypesKeyboard())
	case StateWindowsSameCount:
		b.sendMessage(chatID, "Сколько всего окон?",
			createWindowCountKeyboard())
	case StateWindowsDifferent3:
		b.sendMessage(chatID, "Сколько 3-створчатых окон? (0-6)",
			createWindowCountKeyboard())
	case StateWindowsDifferent4:
		b.sendMessage(chatID, "Сколько 4-створчатых окон? (0-6)",
			createWindowCountKeyboard())
	case StateWindowsDifferent5:
		b.sendMessage(chatID, "Сколько 5-створчатых окон? (0-6)",
			createWindowCountKeyboard())
	case StateWindowsDifferent6_7:
		b.sendMessage(chatID, "Сколько 6-7-створчатых окон? (0-6)",
			createWindowCountKeyboard())
	case StateBalconyNeeded:
		b.sendMessage(chatID, "Нужно ли мыть окна на лоджии?",
			createBalconyNeededKeyboard())
		//todo остальные состояния дописать
	default:
		b.sendMainMenu(chatID)
	}
}

func (b *Bot) handleTextMessage(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	text := msg.Text

	session := b.getSession(chatID)

	switch session.CurrentState {
	case StateWaitingForFloor:
		floor, err := strconv.Atoi(text)
		if err != nil || floor < 1 || floor > 24 {
			b.sendMessage(chatID, "Некорректный этаж. Введите цифру от 1 до 24:")
			return
		}
		session.Order.Floor = floor
		b.updateState(chatID, StateWaitingForApartment)
		b.sendMessage(chatID, "Введите номер квартиры (1-1500):")

	case StateWaitingForApartment:
		if !IsDigitsOnly(text) {
			b.sendMessage(chatID, "Некорректный номер квартиры. Введите только цифры:")
			return
		}
		apartment, err := strconv.Atoi(text)
		if err != nil || apartment < 1 || apartment > 1500 {
			b.sendMessage(chatID, "Некорректный номер квартиры. Введите цифру от 1 до 1500:")
			return
		}
		session.Order.Apartment = text
		b.updateState(chatID, StateWindowsSameOrDifferent)
		b.sendMessage(chatID, "Количество створок на окнах одинаковое или разное?",
			createWindowsSameOrDifferentKeyboard())

	case StateTelegramNick:
		b.handleTelegramNick(chatID, text)

	default:
		b.sendMessage(chatID, "Пожалуйста, используйте кнопки для продолжения.")
	}
}
