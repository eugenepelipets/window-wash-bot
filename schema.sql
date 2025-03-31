-- DROP TABLE IF EXISTS users,
--     orders;

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
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (telegram_id)
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);