DROP TABLE IF EXISTS users,
    orders;

CREATE TABLE IF NOT EXISTS users
(
    id          SERIAL PRIMARY KEY,
    telegram_id BIGINT       NOT NULL UNIQUE,
    username    VARCHAR(100) NOT NULL,
    first_name  VARCHAR(100) NOT NULL,
    last_name   VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders
(
    id          SERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL,
    window_type VARCHAR(20) NOT NULL,
    floor       INTEGER     NOT NULL,
    apartment   VARCHAR(10) NOT NULL,
    price       INTEGER     NOT NULL,
    status      VARCHAR(20)              DEFAULT 'pending',
    is_current  BOOLEAN                  DEFAULT TRUE,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    entrance INTEGER NOT NULL DEFAULT 1,
    windows_same BOOLEAN NOT NULL DEFAULT TRUE,
    window_3_count INTEGER DEFAULT 0,
    window_4_count INTEGER DEFAULT 0,
    window_5_count INTEGER DEFAULT 0,
    window_6_7_count INTEGER DEFAULT 0,
    balcony_count INTEGER DEFAULT 0,
    balcony_type VARCHAR(20),
    balcony_sash VARCHAR(10),
    telegram_nick VARCHAR(100),
    FOREIGN KEY (user_id) REFERENCES users (telegram_id)
);

CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);
CREATE INDEX idx_orders_current ON orders (user_id, apartment, is_current);