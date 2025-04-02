package models

import "time"

type Order struct {
	ID         int64     `db:"id"`
	UserID     int64     `db:"user_id"`
	WindowType string    `db:"window_type"`
	Floor      int       `db:"floor"`
	Apartment  string    `db:"apartment"`
	Price      int       `db:"price"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	User       User      `db:"-"`
}
