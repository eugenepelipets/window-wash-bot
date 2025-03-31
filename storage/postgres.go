package storage

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/eugenepelipets/window-wash-bot/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

// Подключение к БД
func NewPostgres() (*Postgres, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ Переменная DATABASE_URL не задана! Проверь .env")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	// Проверяем соединение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("✅ Подключение к БД установлено")
	return &Postgres{Pool: pool}, nil
}

// Сохранение пользователя
func (p *Postgres) SaveUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (telegram_id) DO UPDATE
		SET username = EXCLUDED.username, first_name = EXCLUDED.first_name, last_name = EXCLUDED.last_name;
	`

	_, err := p.Pool.Exec(ctx, query, user.TelegramID, user.UserName, user.FirstName, user.LastName, time.Now())
	if err != nil {
		log.Printf("⚠️ Ошибка при сохранении пользователя: %v", err)
		return err
	}

	return nil
}

// Сохранение заказа
func (p *Postgres) SaveOrder(order models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверяем обязательные поля
	if order.UserID == 0 || order.WindowType == "" || order.Apartment == "" {
		return errors.New("не заполнены обязательные поля заказа")
	}

	query := `
        INSERT INTO orders (user_id, window_type, floor, apartment, price, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := p.Pool.Exec(ctx, query,
		order.UserID,
		order.WindowType,
		order.Floor,
		order.Apartment,
		order.Price,
		order.Status,
		time.Now(),
	)

	if err != nil {
		log.Printf("⚠️ Ошибка при сохранении заказа: %v", err)
		return err
	}

	log.Printf("✅ Заказ сохранен для пользователя %d", order.UserID)
	return nil
}
