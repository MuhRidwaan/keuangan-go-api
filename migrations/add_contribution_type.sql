-- Migration: tambah kolom type ke saving_contributions
-- Jalankan ini jika tabel saving_contributions sudah ada sebelumnya.
-- Jika fresh install, jalankan init.sql yang sudah diupdate.

ALTER TABLE saving_contributions
    ADD COLUMN IF NOT EXISTS type VARCHAR(3) NOT NULL DEFAULT 'in'
        CHECK (type IN ('in', 'out'));
