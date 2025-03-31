package models

import "time"

type User struct {
	ID         int64     `db:"id"`
	TelegramID int64     `db:"telegram_id"`
	UserName   string    `db:"username"`
	FirstName  string    `db:"first_name"`
	LastName   string    `db:"last_name"`
	CreatedAt  time.Time `db:"created_at"`
}
