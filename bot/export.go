package bot

import (
	"bytes"
	"encoding/csv"
	"github.com/eugenepelipets/window-wash-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var moscowLoc *time.Location

func init() {
	// Инициализация при загрузке пакета
	initLocation()
}

func initLocation() {
	var err error
	moscowLoc, err = time.LoadLocation("Europe/Istanbul")
	if err != nil {
		log.Printf("⚠️ Failed to load Istanbul location: %v (falling back to UTC)", err)
		moscowLoc = time.UTC
	}
}

// handleExport обрабатывает команду экспорта
func (b *Bot) handleExport(msg *tgbotapi.Message) {
	// Проверяем, является ли пользователь администратором
	if !b.isAdmin(msg.Chat.ID) {
		b.sendMessage(msg.Chat.ID, "У вас нет прав для выполнения этой команды.")
		return
	}

	// Определяем тип экспорта (по умолчанию - только актуальные)
	onlyCurrent := true
	if strings.Contains(msg.Text, "все") || strings.Contains(msg.Text, "all") {
		onlyCurrent = false
	}

	// Получаем данные из БД
	orders, err := b.db.GetOrdersForExport(onlyCurrent)
	if err != nil {
		log.Printf("⚠️ Ошибка получения данных для экспорта: %v", err)
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при подготовке отчета.")
		return
	}

	// Создаем CSV
	csvData, err := b.createCSV(orders, onlyCurrent)
	if err != nil {
		log.Printf("⚠️ Ошибка создания CSV: %v", err)
		b.sendMessage(msg.Chat.ID, "Произошла ошибка при создании отчета.")
		return
	}

	// Формируем название файла
	fileName := "orders_current_" + time.Now().Format("2006-01-02") + ".csv"
	if !onlyCurrent {
		fileName = "orders_all_" + time.Now().Format("2006-01-02") + ".csv"
	}

	// Отправляем файл
	file := tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: csvData,
	}
	msgConfig := tgbotapi.NewDocument(msg.Chat.ID, file)
	msgConfig.Caption = "Отчет по заказам (" + b.getExportTypeDescription(onlyCurrent) + ")"

	if _, err := b.api.Send(msgConfig); err != nil {
		log.Printf("⚠️ Ошибка отправки файла: %v", err)
		b.sendMessage(msg.Chat.ID, "Не удалось отправить отчет.")
	}
}

// createCSV создает CSV файл из данных заказов
func (b *Bot) createCSV(orders []models.Order, onlyCurrent bool) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Записываем заголовки
	headers := []string{
		"ID", "Дата создания", "Тип окон", "Этаж", "Квартира",
		"Цена", "Статус", "Актуальный", "ID пользователя",
		"Username", "Имя", "Фамилия",
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// Записываем данные
	for _, order := range orders {
		record := []string{
			strconv.FormatInt(order.ID, 10),
			order.CreatedAt.Format("2006-01-02 15:04:05"),
			order.CreatedAt.In(moscowLoc).Format("2006-01-02 15:04:05"), // Локальное
			order.WindowType,
			strconv.Itoa(order.Floor),
			order.Apartment,
			strconv.Itoa(order.Price),
			order.Status,
			strconv.FormatBool(order.IsCurrent),
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

// getExportTypeDescription возвращает описание типа экспорта
func (b *Bot) getExportTypeDescription(onlyCurrent bool) string {
	if onlyCurrent {
		return "только актуальные"
	}
	return "все заказы"
}

// isAdmin проверяет, является ли пользователь администратором
func (b *Bot) isAdmin(chatID int64) bool {
	adminID, _ := strconv.ParseInt(os.Getenv("ADMIN_TELEGRAM_ID"), 10, 64)
	return chatID == adminID
}
