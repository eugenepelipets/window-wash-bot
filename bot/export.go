package bot

import (
	"bytes"
	"encoding/csv"
	"github.com/eugenepelipets/window-wash-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"time"
)

// handleExport обрабатывает команду экспорта
func (b *Bot) handleExport(chatID int64) {
	// Проверяем, является ли пользователь администратором
	if !b.isAdmin(chatID) {
		b.sendMessage(chatID, "У вас нет прав для выполнения этой команды.")
		return
	}

	// Получаем данные из БД
	orders, err := b.db.GetOrdersForExport()
	if err != nil {
		log.Printf("⚠️ Ошибка получения данных для экспорта: %v", err)
		b.sendMessage(chatID, "Произошла ошибка при подготовке отчета.")
		return
	}

	// Создаем CSV
	csvData, err := b.createCSV(orders)
	if err != nil {
		log.Printf("⚠️ Ошибка создания CSV: %v", err)
		b.sendMessage(chatID, "Произошла ошибка при создании отчета.")
		return
	}

	// Отправляем файл
	file := tgbotapi.FileBytes{
		Name:  "orders_" + time.Now().Format("2006-01-02") + ".csv",
		Bytes: csvData,
	}
	msg := tgbotapi.NewDocument(chatID, file)
	msg.Caption = "Отчет по заказам"

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("⚠️ Ошибка отправки файла: %v", err)
		b.sendMessage(chatID, "Не удалось отправить отчет.")
	}
}

// createCSV создает CSV файл из данных заказов
func (b *Bot) createCSV(orders []models.Order) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Записываем заголовки
	headers := []string{
		"ID", "Дата создания", "Тип окон", "Этаж", "Квартира",
		"Цена", "Статус", "ID пользователя", "Username",
		"Имя", "Фамилия",
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// Записываем данные
	for _, order := range orders {
		record := []string{
			strconv.FormatInt(order.ID, 10),
			order.CreatedAt.Format("2006-01-02 15:04:05"),
			order.WindowType,
			strconv.Itoa(order.Floor),
			order.Apartment,
			strconv.Itoa(order.Price),
			order.Status,
			strconv.FormatInt(order.User.TelegramID, 10),
			order.User.UserName,
			order.User.FirstName,
			order.User.LastName,
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// isAdmin проверяет, является ли пользователь администратором
func (b *Bot) isAdmin(chatID int64) bool {
	adminID, _ := strconv.ParseInt(os.Getenv("ADMIN_TELEGRAM_ID"), 10, 64)
	return chatID == adminID
}
