package models

import "time"

type Order struct {
	ID             int64     `db:"id"`
	UserID         int64     `db:"user_id"`
	Entrance       int       `db:"entrance"`
	Floor          int       `db:"floor"`
	Apartment      string    `db:"apartment"`
	WindowType     string    `db:"window_type"` // "3_same", "4_same", "5_same", "6_7_same", "different"
	WindowsSame    bool      `db:"windows_same"`
	Window3Count   int       `db:"window_3_count"`
	Window4Count   int       `db:"window_4_count"`
	Window5Count   int       `db:"window_5_count"`
	Window6_7Count int       `db:"window_6_7_count"`
	BalconyCount   int       `db:"balcony_count"`
	BalconyType    string    `db:"balcony_type"`
	BalconySash    string    `db:"balcony_sash"`
	TelegramNick   string    `db:"telegram_nick"`
	Price          int       `db:"price"`
	Status         string    `db:"status"` // "confirmed", "needs_clarification", "canceled"
	IsCurrent      bool      `db:"is_current"`
	CreatedAt      time.Time `db:"created_at"`
	User           User      `db:"-"`
}
