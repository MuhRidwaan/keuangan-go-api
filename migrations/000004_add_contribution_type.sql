-- ============================================================
-- MIGRATION 000004: Add Column 'type' to Table saving_contributions
-- ============================================================

ALTER TABLE saving_contributions
    ADD COLUMN IF NOT EXISTS type VARCHAR(3) NOT NULL DEFAULT 'in'
        CHECK (type IN ('in', 'out'));
