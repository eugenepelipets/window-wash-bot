package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func createMainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Новый заказ", "new_order"),
		),
	)
}

func createEntranceKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подъезд 1", "entrance_1"),
			tgbotapi.NewInlineKeyboardButtonData("Подъезд 2", "entrance_2"),
			tgbotapi.NewInlineKeyboardButtonData("Подъезд 3", "entrance_3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подъезд 4", "entrance_4"),
			tgbotapi.NewInlineKeyboardButtonData("Подъезд 5", "entrance_5"),
			tgbotapi.NewInlineKeyboardButtonData("Подъезд 6", "entrance_6"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}

func createWindowsSameOrDifferentKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Одинаковые", "windows_same"),
			tgbotapi.NewInlineKeyboardButtonData("Разные", "windows_different"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}

func createWindowTypesKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3-створчатые", "window_3"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("4-створчатые", "window_4"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("5-створчатые", "window_5"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("6-7-створчатые", "window_6_7"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}

func createWindowCountKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0", "count_0"),
			tgbotapi.NewInlineKeyboardButtonData("1", "count_1"),
			tgbotapi.NewInlineKeyboardButtonData("2", "count_2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3", "count_3"),
			tgbotapi.NewInlineKeyboardButtonData("4", "count_4"),
			tgbotapi.NewInlineKeyboardButtonData("5", "count_5"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("6", "count_6"),
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
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

func createBalconyNeededKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("1 лоджия", "balcony_1"),
			tgbotapi.NewInlineKeyboardButtonData("2 лоджии", "balcony_2"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3 лоджии", "balcony_3"),
			tgbotapi.NewInlineKeyboardButtonData("Не нужно", "balcony_0"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}

func createBalconyTypeKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Стандартные", "balcony_standard"),
			tgbotapi.NewInlineKeyboardButtonData("До пола", "balcony_floor"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}

func createSkipKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пропустить", "skip_nick"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад", "back"),
		),
	)
}
