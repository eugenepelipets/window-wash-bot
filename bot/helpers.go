package bot

import (
	"errors"
	"strconv"
)

// Расчет цены на основе типа окон и этажа
func CalculatePrice(windowType string, floor int) (int, error) {
	basePrice := 0
	floorMultiplier := 1.0

	if floor <= 0 {
		return 0, errors.New("floor must be positive")
	}

	switch windowType {
	case "regular":
		basePrice = 1000
	case "panoramic":
		basePrice = 2000
	case "shop":
		basePrice = 3000
	default:
		basePrice = 1500
	}

	// Увеличиваем цену на 5% за каждый этаж выше 5
	if floor > 5 {
		floorMultiplier = 1.0 + float64(floor-5)*0.05
	}

	totalPrice := float64(basePrice) * floorMultiplier
	return int(totalPrice), nil
}

// Проверка, что строка содержит только цифры
func IsDigitsOnly(s string) bool {
	if s == "" {
		return false
	}

	_, err := strconv.Atoi(s)
	return err == nil
}
