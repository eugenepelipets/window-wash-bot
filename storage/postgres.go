package storage

import (
	"context"
	"fmt"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Определяем тип окон для сохранения в БД
	windowType := "different"
	if order.WindowsSame {
		if order.Window3Count > 0 {
			windowType = "3_same"
		} else if order.Window4Count > 0 {
			windowType = "4_same"
		} else if order.Window5Count > 0 {
			windowType = "5_same"
		} else if order.Window6_7Count > 0 {
			windowType = "6_7_same"
		}
	}

	// Проверяем существующие заказы
	var existingOrderID int64
	err = tx.QueryRow(ctx, `
        SELECT id FROM orders 
        WHERE entrance = $1 AND floor = $2 AND apartment = $3 
        AND is_current = true AND status = 'confirmed'
        LIMIT 1`,
		order.Entrance, order.Floor, order.Apartment).Scan(&existingOrderID)

	orderExists := err == nil

	if orderExists {
		order.Status = "needs_clarification"
		order.IsCurrent = true
	} else {
		order.Status = "confirmed"
		order.IsCurrent = true

		_, err = tx.Exec(ctx, `
            UPDATE orders 
            SET is_current = false 
            WHERE entrance = $1 AND floor = $2 AND apartment = $3 AND is_current = true`,
			order.Entrance, order.Floor, order.Apartment)
		if err != nil {
			return fmt.Errorf("ошибка деактивации предыдущих заказов: %v", err)
		}
	}

	// Сохраняем заказ с явным указанием window_type
	_, err = tx.Exec(ctx, `
        INSERT INTO orders (
            user_id, entrance, floor, apartment, windows_same,
            window_3_count, window_4_count, window_5_count, window_6_7_count,
            balcony_count, balcony_type, balcony_sash, telegram_nick,
            price, status, is_current, created_at, window_type
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9,
            $10, $11, $12, $13, $14, $15, $16, NOW(), $17
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
		order.IsCurrent,
		windowType)

	if err != nil {
		return fmt.Errorf("ошибка сохранения заказа: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("ошибка коммита транзакции: %v", err)
	}

	return nil
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

// CheckExistingOrder проверяет наличие активных заказов для указанной квартиры
func (p *Postgres) CheckExistingOrder(entrance int, floor int, apartment string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	err := p.Pool.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM orders 
            WHERE entrance = $1 AND floor = $2 AND apartment = $3 
            AND is_current = true AND status = 'confirmed'
        )`,
		entrance, floor, apartment).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("ошибка проверки заказов: %v", err)
	}

	return exists, nil
}
