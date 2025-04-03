package bot

import (
	"github.com/eugenepelipets/window-wash-bot/models"
)

// Расчет цены на основе типа окон и этажа
func CalculatePrice(order models.Order) (int, error) {
	total := 0

	// Окна
	total += order.Window3Count * 1000
	total += order.Window4Count * 1500
	total += order.Window5Count * 2000
	total += order.Window6_7Count * 2500

	// Лоджии
	if order.BalconyCount > 0 {
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
		total += balconyPrice * order.BalconyCount
	}

	return total, nil
}

// Проверка, что строка содержит только цифры
func IsDigitsOnly(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
