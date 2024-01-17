DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'money_movement_type') THEN
        CREATE TYPE money_movement_type AS ENUM ('income', 'expense');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS category (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL CHECK (name <> ''),
    money_movement_type money_movement_type NOT NULL
);

CREATE TABLE IF NOT EXISTS "user" (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(100) UNIQUE NOT NULL CHECK (user_id <> '')
);

CREATE TABLE IF NOT EXISTS user_money_movement (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    category_id BIGINT,
    amount NUMERIC(30, 2),
    date DATE,
    description TEXT,
    CONSTRAINT fk_category_id FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_id_date_category_id ON user_money_movement (user_id, date, category_id);

CREATE TABLE IF NOT EXISTS family (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL CHECK (name <> '')
);

CREATE TABLE IF NOT EXISTS user_family (
    user_id BIGINT,
    family_id BIGINT,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    CONSTRAINT fk_family_id FOREIGN KEY (family_id) REFERENCES family (id) ON DELETE CASCADE
);

