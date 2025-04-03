package storage

import (
	"context"
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

	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Помечаем предыдущие заказы как неактуальные
	_, err = tx.Exec(ctx, `
        UPDATE orders 
        SET is_current = false 
        WHERE apartment = $1 AND is_current = true`,
		order.Apartment)
	if err != nil {
		return err
	}

	// Сохраняем новый заказ с явным указанием window_type
	_, err = tx.Exec(ctx, `
        INSERT INTO orders (
            user_id, entrance, floor, apartment, windows_same,
            window_3_count, window_4_count, window_5_count, window_6_7_count,
            balcony_count, balcony_type, balcony_sash, telegram_nick,
            price, status, is_current, created_at, window_type
        ) VALUES (
            $1, $2, $3, $4, $5,
            $6, $7, $8, $9,
            $10, $11, $12, $13,
            $14, $15, true, NOW(), $16
        )`,
		order.UserID,
		order.Entrance,
		order.Floor,
		order.Apartment,
		order.WindowsSame,
		order.Window3Count,
		order.Window4Count,
		order.Window5Count,
		order.Window6_7Count,
		order.BalconyCount,
		order.BalconyType,
		order.BalconySash,
		order.TelegramNick,
		order.Price,
		order.Status,
		"custom_type", // Здесь нужно указать актуальное значение
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// storage/postgres.go
// Добавляем в конец файла

// GetOrdersForExport получает заказы для экспорта (с фильтром по актуальности)
func (p *Postgres) GetOrdersForExport(onlyCurrent bool) ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
        SELECT 
            o.id, o.entrance, o.floor, o.apartment, o.windows_same,
            o.window_3_count, o.window_4_count, o.window_5_count, o.window_6_7_count,
            o.balcony_count, o.balcony_type, o.balcony_sash, o.telegram_nick,
            o.price, o.status, o.is_current, o.created_at,
            u.telegram_id, u.username, u.first_name, u.last_name
        FROM orders o
        JOIN users u ON o.user_id = u.telegram_id
        WHERE $1 = false OR o.is_current = true
        ORDER BY o.created_at DESC
    `

	rows, err := p.Pool.Query(ctx, query, onlyCurrent)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var user models.User
		err := rows.Scan(
			&order.ID,
			&order.Entrance,
			&order.Floor,
			&order.Apartment,
			&order.WindowsSame,
			&order.Window3Count,
			&order.Window4Count,
			&order.Window5Count,
			&order.Window6_7Count,
			&order.BalconyCount,
			&order.BalconyType,
			&order.BalconySash,
			&order.TelegramNick,
			&order.Price,
			&order.Status,
			&order.IsCurrent,
			&order.CreatedAt,
			&user.TelegramID,
			&user.UserName,
			&user.FirstName,
			&user.LastName,
		)
		if err != nil {
			return nil, err
		}
		order.User = user
		orders = append(orders, order)
	}

	return orders, nil
}
