-- ============================================================
-- FULL DATABASE SCHEMA & SEED DATA (SUPABASE / POSTGRESQL)
-- Jalankan file ini di Supabase SQL Editor
-- ============================================================

-- 1. Enable Extension pgcrypto untuk gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 2. Bersihkan tabel lama jika ada
DROP TABLE IF EXISTS notifications            CASCADE;
DROP TABLE IF EXISTS password_reset_tokens     CASCADE;
DROP TABLE IF EXISTS api_documentations        CASCADE;
DROP TABLE IF EXISTS agenda_members            CASCADE;
DROP TABLE IF EXISTS agendas                   CASCADE;
DROP TABLE IF EXISTS saving_contributions      CASCADE;
DROP TABLE IF EXISTS saving_members            CASCADE;
DROP TABLE IF EXISTS saving_goals              CASCADE;
DROP TABLE IF EXISTS transactions              CASCADE;
DROP TABLE IF EXISTS categories                CASCADE;
DROP TABLE IF EXISTS users                     CASCADE;

-- ============================================================
-- 1. TABEL users
-- ============================================================
CREATE TABLE users (
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

CREATE INDEX idx_users_email      ON users (email);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);

-- ============================================================
-- 2. TABEL categories
-- ============================================================
CREATE TABLE categories (
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

CREATE INDEX idx_categories_user_id    ON categories (user_id);
CREATE INDEX idx_categories_deleted_at ON categories (deleted_at);

-- SEED: Kategori Sistem Default (user_id NULL = global/system)
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
    ('Kesehatan',     'expense', 'heart');

-- ============================================================
-- 3. TABEL transactions
-- ============================================================
CREATE TABLE transactions (
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

CREATE INDEX idx_transactions_user_id     ON transactions (user_id);
CREATE INDEX idx_transactions_category_id ON transactions (category_id);
CREATE INDEX idx_transactions_deleted_at  ON transactions (deleted_at);

-- ============================================================
-- 4. TABEL saving_goals
-- ============================================================
CREATE TABLE saving_goals (
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

CREATE INDEX idx_saving_goals_deleted_at ON saving_goals (deleted_at);

-- ============================================================
-- 5. TABEL saving_members
-- ============================================================
CREATE TABLE saving_members (
    goal_id    UUID        NOT NULL REFERENCES saving_goals (id) ON DELETE CASCADE,
    user_id    UUID        NOT NULL REFERENCES users        (id) ON DELETE CASCADE,
    role       VARCHAR(10) NOT NULL DEFAULT 'member' CHECK (role IN ('owner', 'member')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),

    PRIMARY KEY (goal_id, user_id)
);

-- ============================================================
-- 6. TABEL saving_contributions
-- ============================================================
CREATE TABLE saving_contributions (
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

CREATE INDEX idx_saving_contributions_goal_id    ON saving_contributions (goal_id);
CREATE INDEX idx_saving_contributions_user_id    ON saving_contributions (user_id);
CREATE INDEX idx_saving_contributions_deleted_at ON saving_contributions (deleted_at);

-- ============================================================
-- 7. TABEL agendas
-- ============================================================
CREATE TABLE agendas (
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

CREATE INDEX idx_agendas_deleted_at ON agendas (deleted_at);

-- ============================================================
-- 8. TABEL agenda_members
-- ============================================================
CREATE TABLE agenda_members (
    agenda_id  UUID        NOT NULL REFERENCES agendas (id) ON DELETE CASCADE,
    user_id    UUID        NOT NULL REFERENCES users   (id) ON DELETE CASCADE,
    role       VARCHAR(10) NOT NULL DEFAULT 'member' CHECK (role IN ('owner', 'member')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255),

    PRIMARY KEY (agenda_id, user_id)
);

-- ============================================================
-- 9. TABEL notifications
-- ============================================================
CREATE TABLE notifications (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    message    TEXT        NOT NULL,
    is_read    BOOLEAN     NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications (user_id);

-- ============================================================
-- 10. TABEL password_reset_tokens
-- ============================================================
CREATE TABLE password_reset_tokens (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at    TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_password_reset_user_id ON password_reset_tokens (user_id);
CREATE INDEX idx_password_reset_token   ON password_reset_tokens (token);

-- ============================================================
-- 11. TABEL api_documentations & SEED DATA
-- ============================================================
CREATE TABLE api_documentations (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category         VARCHAR(50)  NOT NULL,
    name             VARCHAR(100) NOT NULL,
    method           VARCHAR(10)  NOT NULL,
    endpoint         VARCHAR(255) NOT NULL,
    is_protected     BOOLEAN      NOT NULL DEFAULT false,
    description      TEXT,
    headers          JSONB,
    request_body     JSONB,
    response_success JSONB,
    response_error   JSONB,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_api_doc_category ON api_documentations (category);
CREATE INDEX idx_api_doc_method   ON api_documentations (method);

-- SEED DOKUMENTASI API
INSERT INTO api_documentations (category, name, method, endpoint, is_protected, description, headers, request_body, response_success, response_error) VALUES
('Auth', 'Register User Baru', 'POST', '/api/register', false, 'Mendaftarkan pengguna baru', '{"Content-Type": "application/json"}'::jsonb, '{"name": "Budi Santoso", "email": "budi@gmail.com", "password": "password123"}'::jsonb, '{"meta": {"code": 201, "status": "success", "message": "Registrasi berhasil"}, "data": {"id": "uuid", "name": "Budi Santoso", "email": "budi@gmail.com"}}'::jsonb, '{"meta": {"code": 409, "status": "error", "message": "email sudah terdaftar"}, "data": null}'::jsonb),
('Auth', 'Login User', 'POST', '/api/login', false, 'Otentikasi pengguna untuk mendapatkan JWT Token', '{"Content-Type": "application/json"}'::jsonb, '{"email": "budi@gmail.com", "password": "password123"}'::jsonb, '{"meta": {"code": 200, "status": "success", "message": "Login berhasil"}, "data": {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}}'::jsonb, '{"meta": {"code": 401, "status": "error", "message": "email atau password salah"}, "data": null}'::jsonb),
('Auth', 'Lupa Password (Request OTP)', 'POST', '/api/forgot-password', false, 'Mengirimkan token reset password ke email terdaftar', '{"Content-Type": "application/json"}'::jsonb, '{"email": "budi@gmail.com"}'::jsonb, '{"meta": {"code": 200, "status": "success", "message": "Jika email terdaftar, link reset password telah dikirim"}, "data": null}'::jsonb, '{"meta": {"code": 400, "status": "error", "message": "Email tidak valid"}, "data": null}'::jsonb),
('Auth', 'Reset Password dengan OTP', 'POST', '/api/reset-password', false, 'Mereset password pengguna menggunakan token', '{"Content-Type": "application/json"}'::jsonb, '{"email": "budi@gmail.com", "token": "123456", "new_password": "newpassword123"}'::jsonb, '{"meta": {"code": 200, "status": "success", "message": "Password berhasil direset, silakan login kembali"}, "data": null}'::jsonb, '{"meta": {"code": 400, "status": "error", "message": "Kode Token tidak valid"}, "data": null}'::jsonb),
('Categories', 'Get List Kategori', 'GET', '/api/categories', true, 'Mengambil daftar kategori transaksi', '{"Authorization": "Bearer <TOKEN>"}'::jsonb, null, '{"meta": {"code": 200, "status": "success", "message": "Berhasil mengambil kategori"}, "data": [{"id": "uuid", "name": "Gaji", "type": "income", "icon": "briefcase", "user_id": null}]}'::jsonb, '{"meta": {"code": 401, "status": "error", "message": "Token tidak ditemukan"}, "data": null}'::jsonb),
('Categories', 'Buat Kategori Custom', 'POST', '/api/categories', true, 'Membuat kategori kustom baru', '{"Content-Type": "application/json", "Authorization": "Bearer <TOKEN>"}'::jsonb, '{"name": "Investasi Crypto", "type": "income", "icon": "bitcoin"}'::jsonb, '{"meta": {"code": 201, "status": "success", "message": "Kategori berhasil dibuat"}, "data": {"id": "uuid", "name": "Investasi Crypto"}}'::jsonb, '{"meta": {"code": 400, "status": "error", "message": "Nama kategori wajib diisi"}, "data": null}'::jsonb),
('Transactions', 'Tambah Transaksi', 'POST', '/api/transactions', true, 'Mencatat transaksi keuangan baru', '{"Content-Type": "application/json", "Authorization": "Bearer <TOKEN>"}'::jsonb, '{"category_id": "uuid-cat", "amount": 5000000, "date": "2026-04-10", "notes": "Gaji Bulan April"}'::jsonb, '{"meta": {"code": 201, "status": "success", "message": "Transaksi berhasil disimpan"}, "data": {"id": "uuid", "amount": 5000000}}'::jsonb, '{"meta": {"code": 400, "status": "error", "message": "category_id tidak valid"}, "data": null}'::jsonb),
('Transactions', 'Get List Transaksi', 'GET', '/api/transactions', true, 'Mengambil seluruh daftar transaksi', '{"Authorization": "Bearer <TOKEN>"}'::jsonb, null, '{"meta": {"code": 200, "status": "success", "message": "Berhasil mengambil transaksi"}, "data": []}'::jsonb, '{"meta": {"code": 401, "status": "error", "message": "Unauthorized"}, "data": null}'::jsonb),
('Savings', 'Buat Target Tabungan', 'POST', '/api/savings', true, 'Membuat target tabungan baru', '{"Content-Type": "application/json", "Authorization": "Bearer <TOKEN>"}'::jsonb, '{"title": "Liburan ke Bali", "target_amount": 10000000, "deadline": "2026-12-31"}'::jsonb, '{"meta": {"code": 201, "status": "success", "message": "Saving goal berhasil dibuat"}, "data": {"id": "uuid", "title": "Liburan ke Bali"}}'::jsonb, '{"meta": {"code": 400, "status": "error", "message": "target_amount minimal 1000"}, "data": null}'::jsonb),
('Savings', 'Get List Tabungan Saya', 'GET', '/api/savings', true, 'Mengambil daftar seluruh target tabungan', '{"Authorization": "Bearer <TOKEN>"}'::jsonb, null, '{"meta": {"code": 200, "status": "success", "message": "Berhasil mengambil saving goals"}, "data": []}'::jsonb, '{"meta": {"code": 401, "status": "error", "message": "Unauthorized"}, "data": null}'::jsonb),
('Agendas', 'Buat Agenda Baru', 'POST', '/api/agendas', true, 'Membuat jadwal / agenda baru', '{"Content-Type": "application/json", "Authorization": "Bearer <TOKEN>"}'::jsonb, '{"title": "Meeting Rapat", "description": "Pembahasan anggaran", "start_date": "2026-04-15T09:00:00Z", "end_date": "2026-04-15T11:00:00Z"}'::jsonb, '{"meta": {"code": 201, "status": "success", "message": "Agenda berhasil dibuat"}, "data": {"id": "uuid"}}'::jsonb, '{"meta": {"code": 400, "status": "error", "message": "end_date harus lebih besar"}, "data": null}'::jsonb),
('Agendas', 'Get List Agenda Saya', 'GET', '/api/agendas', true, 'Mengambil daftar agenda', '{"Authorization": "Bearer <TOKEN>"}'::jsonb, null, '{"meta": {"code": 200, "status": "success", "message": "Berhasil mengambil agenda"}, "data": []}'::jsonb, '{"meta": {"code": 401, "status": "error", "message": "Unauthorized"}, "data": null}'::jsonb),
('Notifications', 'Get List Notifikasi', 'GET', '/api/notifications', true, 'Mengambil notifikasi masuk', '{"Authorization": "Bearer <TOKEN>"}'::jsonb, null, '{"meta": {"code": 200, "status": "success", "message": "Berhasil mengambil notifikasi"}, "data": []}'::jsonb, '{"meta": {"code": 401, "status": "error", "message": "Unauthorized"}, "data": null}'::jsonb),
('Notifications', 'Tandai Read Notifikasi', 'PUT', '/api/notifications/:id/read', true, 'Menandai notifikasi dibaca', '{"Authorization": "Bearer <TOKEN>"}'::jsonb, null, '{"meta": {"code": 200, "status": "success", "message": "Notifikasi ditandai dibaca"}, "data": null}'::jsonb, '{"meta": {"code": 404, "status": "error", "message": "Notifikasi tidak ditemukan"}, "data": null}'::jsonb);
