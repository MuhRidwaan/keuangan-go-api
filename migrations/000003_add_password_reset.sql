-- ============================================================
-- MIGRATION 000003: Add Password Reset Tokens Table
-- ============================================================

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ  NOT NULL,
    used_at    TIMESTAMPTZ  DEFAULT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT password_reset_tokens_token_unique UNIQUE (token)
);

CREATE INDEX IF NOT EXISTS idx_prt_user_id ON password_reset_tokens (user_id);
CREATE INDEX IF NOT EXISTS idx_prt_token   ON password_reset_tokens (token);
