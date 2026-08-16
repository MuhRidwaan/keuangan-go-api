-- ============================================================
-- MIGRATION 000001: Initial Core Schema (Safe & Non-Destructive)
-- ============================================================

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 1. users
CREATE TABLE IF NOT EXISTS users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by    VARCHAR(255),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by    VARCHAR(255),
    deleted_at    TIMESTAMPTZ,
    deleted_by    VARCHAR(255),

    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE INDEX IF NOT EXISTS idx_users_email      ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);

-- 2. categories
CREATE TABLE IF NOT EXISTS categories (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(255),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(255),

    user_id    UUID         REFERENCES users (id) ON DELETE SET NULL,
    name       VARCHAR(100) NOT NULL,
    type       VARCHAR(10)  NOT NULL CHECK (type IN ('income', 'expense')),
    icon       VARCHAR(100)
);

CREATE INDEX IF NOT EXISTS idx_categories_user_id    ON categories (user_id);
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories (deleted_at);

-- 3. transactions
CREATE TABLE IF NOT EXISTS transactions (
    id          UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    created_by  VARCHAR(255),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMPTZ,
    deleted_by  VARCHAR(255),

    user_id     UUID           NOT NULL REFERENCES users      (id) ON DELETE CASCADE,
    category_id UUID           NOT NULL REFERENCES categories (id) ON DELETE RESTRICT,
    amount      NUMERIC(15, 2) NOT NULL,
    date        DATE           NOT NULL,
    notes       TEXT
);

CREATE INDEX IF NOT EXISTS idx_transactions_user_id     ON transactions (user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions (category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_deleted_at  ON transactions (deleted_at);

-- 4. saving_goals
CREATE TABLE IF NOT EXISTS saving_goals (
    id             UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at     TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    created_by     VARCHAR(255),
    updated_at     TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_by     VARCHAR(255),
    deleted_at     TIMESTAMPTZ,
    deleted_by     VARCHAR(255),

    title          VARCHAR(255)   NOT NULL,
    target_amount  NUMERIC(15, 2) NOT NULL,
    current_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    deadline       DATE           NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_saving_goals_deleted_at ON saving_goals (deleted_at);

-- 5. saving_members
CREATE TABLE IF NOT EXISTS saving_members (
    goal_id    UUID        NOT NULL REFERENCES saving_goals (id) ON DELETE CASCADE,
    user_id    UUID        NOT NULL REFERENCES users        (id) ON DELETE CASCADE,
    role       VARCHAR(10) NOT NULL DEFAULT 'member' CHECK (role IN ('owner', 'member')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),

    PRIMARY KEY (goal_id, user_id)
);

-- 6. saving_contributions
CREATE TABLE IF NOT EXISTS saving_contributions (
    id         UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(255),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(255),

    goal_id    UUID           NOT NULL REFERENCES saving_goals (id) ON DELETE CASCADE,
    user_id    UUID           NOT NULL REFERENCES users        (id) ON DELETE CASCADE,
    amount     NUMERIC(15, 2) NOT NULL,
    date       DATE           NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_saving_contributions_goal_id    ON saving_contributions (goal_id);
CREATE INDEX IF NOT EXISTS idx_saving_contributions_user_id    ON saving_contributions (user_id);
CREATE INDEX IF NOT EXISTS idx_saving_contributions_deleted_at ON saving_contributions (deleted_at);

-- 7. agendas
CREATE TABLE IF NOT EXISTS agendas (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by  VARCHAR(255),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by  VARCHAR(255),
    deleted_at  TIMESTAMPTZ,
    deleted_by  VARCHAR(255),

    title       VARCHAR(255) NOT NULL,
    description TEXT,
    start_date  TIMESTAMPTZ  NOT NULL,
    end_date    TIMESTAMPTZ  NOT NULL,

    CONSTRAINT agendas_end_after_start CHECK (end_date > start_date)
);

CREATE INDEX IF NOT EXISTS idx_agendas_deleted_at ON agendas (deleted_at);

-- 8. agenda_members
CREATE TABLE IF NOT EXISTS agenda_members (
    agenda_id  UUID        NOT NULL REFERENCES agendas (id) ON DELETE CASCADE,
    user_id    UUID        NOT NULL REFERENCES users   (id) ON DELETE CASCADE,
    role       VARCHAR(10) NOT NULL DEFAULT 'member' CHECK (role IN ('owner', 'member')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),

    PRIMARY KEY (agenda_id, user_id)
);

-- SEED: Kategori default sistem
INSERT INTO categories (name, type, icon) VALUES
    ('Gaji',          'income',  'briefcase'),
    ('Freelance',     'income',  'laptop'),
    ('Investasi',     'income',  'trending-up'),
    ('Bonus',         'income',  'gift'),
    ('Makan & Minum', 'expense', 'utensils'),
    ('Transportasi',  'expense', 'car'),
    ('Belanja',       'expense', 'shopping-bag'),
    ('Tagihan',       'expense', 'file-text'),
    ('Hiburan',       'expense', 'film'),
    ('Kesehatan',     'expense', 'heart')
ON CONFLICT DO NOTHING;
