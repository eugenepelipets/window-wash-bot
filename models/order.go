package models

import "time"

type Order struct {
	ID             int64     `db:"id"`
	UserID         int64     `db:"user_id"`
	Entrance       int       `db:"entrance"`
	Floor          int       `db:"floor"`
	Apartment      string    `db:"apartment"`
	WindowsSame    bool      `db:"windows_same"`
	Window3Count   int       `db:"window_3_count"`
	Window4Count   int       `db:"window_4_count"`
	Window5Count   int       `db:"window_5_count"`
	Window6_7Count int       `db:"window_6_7_count"`
	BalconyCount   int       `db:"balcony_count"`
	BalconyType    string    `db:"balcony_type"` // "standard" или "floor_to_ceiling"
	BalconySash    string    `db:"balcony_sash"` // "3", "4", "5", "6_7"
	TelegramNick   string    `db:"telegram_nick"`
	Price          int       `db:"price"`
	Status         string    `db:"status"`
	CreatedAt      time.Time `db:"created_at"`
	IsCurrent      bool      `db:"is_current"`
	User           User      `db:"-"`
}
