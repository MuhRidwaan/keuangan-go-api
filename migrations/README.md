# 🗄️ System Migration Incremental — Keuangan API

Seluruh file migrasi disusun secara **Forward-Only Incremental & Non-Destructive** (tanpa file rollback `.down.sql` dan tanpa menghapus data/tabel yang sudah ada).

## 📜 Daftar Versi Migrasi

| Nomor | File SQL Migrasi | Penjelasan Perubahan |
| :--- | :--- | :--- |
| **000001** | [`000001_init_schema.sql`](./000001_init_schema.sql) | Inisialisasi tabel dasar (`users`, `categories`, `transactions`, `saving_goals`, `saving_members`, `saving_contributions`, `agendas`, `agenda_members`) |
| **000002** | [`000002_add_notifications.sql`](./000002_add_notifications.sql) | Penambahan tabel notifikasi (`notifications`) |
| **000003** | [`000003_add_password_reset.sql`](./000003_add_password_reset.sql) | Penambahan tabel reset password (`password_reset_tokens`) |
| **000004** | [`000004_add_contribution_type.sql`](./000004_add_contribution_type.sql) | Menambahkan kolom `type` (`in`/`out`) ke tabel `saving_contributions` |
| **000005** | [`000005_api_documentation.sql`](./000005_api_documentation.sql) | Penambahan tabel & seed data dokumentasi API (`api_documentations`) |

---

## 🔒 Keamanan Data User
* File `.sql` **TIDAK PERNAH** menggunakan `DROP TABLE`.
* Menggunakan `CREATE TABLE IF NOT EXISTS` untuk tabel baru.
* Menggunakan `ALTER TABLE ... ADD COLUMN IF NOT EXISTS` untuk penambahan kolom baru.
* Data yang sudah ada di database local kamu **100% aman dan tidak akan terhapus**.

---

## 🚀 Cara Menjalankan Migrasi Manual
```bash
psql -U postgres -d db_keuangan -f migrations/000001_init_schema.sql
psql -U postgres -d db_keuangan -f migrations/000002_add_notifications.sql
psql -U postgres -d db_keuangan -f migrations/000003_add_password_reset.sql
psql -U postgres -d db_keuangan -f migrations/000004_add_contribution_type.sql
psql -U postgres -d db_keuangan -f migrations/000005_api_documentation.sql
```
