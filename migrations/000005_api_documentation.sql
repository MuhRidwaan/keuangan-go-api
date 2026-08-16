-- ============================================================
-- SCHEMA & SEED: api_documentations
-- Tabel untuk menyimpan dokumentasi endpoint API beserta contoh request & response
-- ============================================================

CREATE TABLE IF NOT EXISTS api_documentations (
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

CREATE INDEX IF NOT EXISTS idx_api_doc_category ON api_documentations (category);
CREATE INDEX IF NOT EXISTS idx_api_doc_method   ON api_documentations (method);

-- Clean up previous seed if re-running
TRUNCATE TABLE api_documentations;

-- ============================================================
-- SEED DATA DOKUMENTASI API
-- ============================================================

-- 1. AUTH
INSERT INTO api_documentations (category, name, method, endpoint, is_protected, description, headers, request_body, response_success, response_error) VALUES
(
    'Auth',
    'Register User Baru',
    'POST',
    '/api/register',
    false,
    'Mendaftarkan pengguna baru ke dalam sistem keuangan',
    '{"Content-Type": "application/json"}'::jsonb,
    '{"name": "Budi Santoso", "email": "budi@gmail.com", "password": "password123"}'::jsonb,
    '{"meta": {"code": 201, "status": "success", "message": "Registrasi berhasil"}, "data": {"id": "d3b07384-d113-4601-a757-5509e564d262", "name": "Budi Santoso", "email": "budi@gmail.com"}}'::jsonb,
    '{"meta": {"code": 409, "status": "error", "message": "email sudah terdaftar"}, "data": null}'::jsonb
),
(
    'Auth',
    'Login User',
    'POST',
    '/api/login',
    false,
    'Otentikasi pengguna untuk mendapatkan JWT Token Bearer',
    '{"Content-Type": "application/json"}'::jsonb,
    '{"email": "budi@gmail.com", "password": "password123"}'::jsonb,
    '{"meta": {"code": 200, "status": "success", "message": "Login berhasil"}, "data": {"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}}'::jsonb,
    '{"meta": {"code": 401, "status": "error", "message": "email atau password salah"}, "data": null}'::jsonb
),
(
    'Auth',
    'Lupa Password (Request OTP)',
    'POST',
    '/api/forgot-password',
    false,
    'Mengirimkan kode OTP reset password ke email terdaftar',
    '{"Content-Type": "application/json"}'::jsonb,
    '{"email": "budi@gmail.com"}'::jsonb,
    '{"meta": {"code": 200, "status": "success", "message": "Jika email terdaftar, link reset password telah dikirim"}, "data": null}'::jsonb,
    '{"meta": {"code": 400, "status": "error", "message": "Email tidak valid"}, "data": null}'::jsonb
),
(
    'Auth',
    'Reset Password dengan OTP',
    'POST',
    '/api/reset-password',
    false,
    'Mereset password pengguna menggunakan OTP dari email',
    '{"Content-Type": "application/json"}'::jsonb,
    '{"email": "budi@gmail.com", "code": "123456", "new_password": "newpassword123"}'::jsonb,
    '{"meta": {"code": 200, "status": "success", "message": "Password berhasil direset, silakan login kembali"}, "data": null}'::jsonb,
    '{"meta": {"code": 400, "status": "error", "message": "Kode OTP tidak valid atau kadaluarsa"}, "data": null}'::jsonb
);

-- 2. CATEGORIES
INSERT INTO api_documentations (category, name, method, endpoint, is_protected, description, headers, request_body, response_success, response_error) VALUES
(
    'Categories',
    'Get List Kategori',
    'GET',
    '/api/categories',
    true,
    'Mengambil daftar kategori transaksi (sistem/global + milik user)',
    '{"Authorization": "Bearer <TOKEN>"}'::jsonb,
    null,
    '{"meta": {"code": 200, "status": "success", "message": "Berhasil mengambil kategori"}, "data": [{"id": "c1f7a091-8888-4c12-9901-5123456789ab", "name": "Gaji", "type": "income", "icon": "briefcase", "user_id": null}, {"id": "d2f7a091-9999-4c12-9901-5123456789cd", "name": "Side Project", "type": "income", "icon": "code", "user_id": "d3b07384-d113-4601-a757-5509e564d262"}]}'::jsonb,
    '{"meta": {"code": 401, "status": "error", "message": "Token tidak ditemukan"}, "data": null}'::jsonb
),
(
    'Categories',
    'Buat Kategori Custom',
    'POST',
    '/api/categories',
    true,
    'Membuat kategori kustom baru untuk user yang sedang login',
    '{"Content-Type": "application/json", "Authorization": "Bearer <TOKEN>"}'::jsonb,
    '{"name": "Investasi Crypto", "type": "income", "icon": "bitcoin"}'::jsonb,
    '{"meta": {"code": 201, "status": "success", "message": "Kategori berhasil dibuat"}, "data": {"id": "e3f7a091-0000-4c12-9901-5123456789ef", "user_id": "d3b07384-d113-4601-a757-5509e564d262", "name": "Investasi Crypto", "type": "income", "icon": "bitcoin"}}'::jsonb,
    '{"meta": {"code": 400, "status": "error", "message": "Nama kategori wajib diisi"}, "data": null}'::jsonb
);

-- 3. TRANSACTIONS
INSERT INTO api_documentations (category, name, method, endpoint, is_protected, description, headers, request_body, response_success, response_error) VALUES
(
    'Transactions',
    'Tambah Transaksi',
    'POST',
    '/api/transactions',
    true,
    'Mencatat pemasukan atau pengeluaran keuangan baru',
    '{"Content-Type": "application/json", "Authorization": "Bearer <TOKEN>"}'::jsonb,
    '{"category_id": "c1f7a091-8888-4c12-9901-5123456789ab", "amount": 5000000, "date": "2026-04-10", "notes": "Gaji Bulan April"}'::jsonb,
    '{"meta": {"code": 201, "status": "success", "message": "Transaksi berhasil disimpan"}, "data": {"id": "f4f7a091-1111-4c12-9901-512345678900", "user_id": "d3b07384-d113-4601-a757-5509e564d262", "category_id": "c1f7a091-8888-4c12-9901-5123456789ab", "amount": 5000000, "date": "2026-04-10T00:00:00Z", "notes": "Gaji Bulan April"}}'::jsonb,
    '{"meta": {"code": 400, "status": "error", "message": "category_id tidak valid"}, "data": null}'::jsonb
),
(
    'Transactions',
    'Get List Transaksi',
    'GET',
    '/api/transactions',
    true,
    'Mengambil seluruh daftar transaksi milik user yang login beserta detail kategorinya',
    '{"Authorization": "Bearer <TOKEN>"}'::jsonb,
    null,
    '{"meta": {"code": 200, "status": "success", "message": "Berhasil mengambil transaksi"}, "data": [{"id": "f4f7a091-1111-4c12-9901-512345678900", "user_id": "d3b07384-d113-4601-a757-5509e564d262", "category_id": "c1f7a091-8888-4c12-9901-5123456789ab", "amount": 5000000, "date": "2026-04-10T00:00:00Z", "notes": "Gaji Bulan April", "category": {"id": "c1f7a091-8888-4c12-9901-5123456789ab", "name": "Gaji", "type": "income", "icon": "briefcase"}}]}'::jsonb,
    '{"meta": {"code": 401, "status": "error", "message": "Unauthorized"}, "data": null}'::jsonb
),
(
    'Transactions',
    'Update Transaksi',
    'PUT',
    '/api/transactions/:id',
    true,
    'Perbarui data transaksi yang sudah ada',
    '{"Content-Type": "application/json", "Authorization": "Bearer <TOKEN>"}'::jsonb,
    '{"category_id": "c1f7a091-8888-4c12-9901-5123456789ab", "amount": 5500000, "date": "2026-04-10", "notes": "Gaji + Bonus April"}'::jsonb,
    '{"meta": {"code": 200, "status": "success", "message": "Transaksi berhasil diperbarui"}, "data": {"id": "f4f7a091-1111-4c12-9901-512345678900", "amount": 5500000, "notes": "Gaji + Bonus April"}}'::jsonb,
    '{"meta": {"code": 404, "status": "error", "message": "Transaksi tidak ditemukan"}, "data": null}'::jsonb
),
(
    'Transactions',
    'Hapus Transaksi',
    'DELETE',
    '/api/transactions/:id',
    true,
    'Menghapus transaksi tertentu berdasarkan ID (soft delete)',
    '{"Authorization": "Bearer <TOKEN>"}'::jsonb,
    null,
    '{"meta": {"code": 200, "status": "success", "message": "Transaksi berhasil dihapus"}, "data": null}'::jsonb,
    '{"meta": {"code": 404, "status": "error", "message": "Transaksi tidak ditemukan"}, "data": null}'::jsonb
);

-- 4. SAVINGS (Tabungan Bersama)
INSERT INTO api_documentations (category, name, method, endpoint, is_protected, description, headers, request_body, response_success, response_error) VALUES
(
    'Savings',
    'Buat Target Tabungan',
    'POST',
    '/api/savings',
    true,
    'Membuat target tabungan bersama / individu baru',
    '{"Content-Type": "application/json", "Authorization": "Bearer <TOKEN>"}'::jsonb,
    '{"title": "Liburan ke Bali", "target_amount": 10000000, "deadline": "2026-12-31"}'::jsonb,
    '{"meta": {"code": 201, "status": "success", "message": "Saving goal berhasil dibuat"}, "data": {"id": "a1b2c3d4-1111-4444-8888-999999999999", "title": "Liburan ke Bali", "target_amount": 10000000, "current_amount": 0, "deadline": "2026-12-31T00:00:00Z"}}'::jsonb,
    '{"meta": {"code": 400, "status": "error", "message": "target_amount minimal 1000"}, "data": null}'::jsonb
),
(
    'Savings',
    'Get List Tabungan Saya',
    'GET',
    '/api/savings',
    true,
    'Mengambil daftar seluruh target tabungan di mana pengguna menjadi owner / member',
    '{"Authorization": "Bearer <TOKEN>"}'::jsonb,
    null,
    '{"meta": {"code": 200, "status": "success", "message": "Berhasil mengambil saving goals"}, "data": [{"id": "a1b2c3d4-1111-4444-8888-999999999999", "title": "Liburan ke Bali", "target_amount": 10000000, "current_amount": 2500000, "deadline": "2026-12-31T00:00:00Z", "my_role": "owner"}]}'::jsonb,
    '{"meta": {"code": 401, "status": "error", "message": "Unauthorized"}, "data": null}'::jsonb
),
(
    'Savings',
    'Tambah Anggota Tabungan',
    'POST',
    '/api/savings/:id/members',
    true,
    'Menambahkan anggota baru ke target tabungan bersama (hanya Owner)',
    '{"Content-Type": "application/json", "Authorization": "Bearer <TOKEN>"}'::jsonb,
    '{"email": "teman@gmail.com"}'::jsonb,
    '{"meta": {"code": 201, "status": "success", "message": "Anggota berhasil ditambahkan"}, "data": {"goal_id": "a1b2c3d4-1111-4444-8888-999999999999", "user_id": "e4c5b6a7-2222-3333-4444-555555555555", "role": "member"}}'::jsonb,
    '{"meta": {"code": 403, "status": "error", "message": "hanya owner yang dapat menambah anggota"}, "data": null}'::jsonb
),
(
    'Savings',
    'Setor Kontribusi Tabungan',
    'POST',
    '/api/savings/:id/contribute',
    true,
    'Menyetorkan dana kontribusi ke target tabungan',
    '{"Content-Type": "application/json", "Authorization": "Bearer <TOKEN>"}'::jsonb,
    '{"amount": 500000, "date": "2026-04-10"}'::jsonb,
    '{"meta": {"code": 201, "status": "success", "message": "Kontribusi berhasil disimpan"}, "data": {"id": "c9b8a7f6-1234-5678-90ab-cdef12345678", "goal_id": "a1b2c3d4-1111-4444-8888-999999999999", "user_id": "d3b07384-d113-4601-a757-5509e564d262", "amount": 500000, "date": "2026-04-10T00:00:00Z"}}'::jsonb,
    '{"meta": {"code": 400, "status": "error", "message": "jumlah kontribusi harus lebih besar dari 0"}, "data": null}'::jsonb
),
(
    'Savings',
    'Tarik Dana Tabungan',
    'POST',
    '/api/savings/:id/withdraw',
    true,
    'Penarikan dana dari target tabungan (mengurangi total terkumpul)',
    '{"Content-Type": "application/json", "Authorization": "Bearer <TOKEN>"}'::jsonb,
    '{"amount": 200000, "date": "2026-04-12"}'::jsonb,
    '{"meta": {"code": 201, "status": "success", "message": "Penarikan berhasil diproses"}, "data": {"id": "w1b2c3d4-0000-1111-2222-333333333333", "goal_id": "a1b2c3d4-1111-4444-8888-999999999999", "user_id": "d3b07384-d113-4601-a757-5509e564d262", "amount": -200000, "date": "2026-04-12T00:00:00Z"}}'::jsonb,
    '{"meta": {"code": 400, "status": "error", "message": "saldo terkumpul tidak mencukupi untuk penarikan"}, "data": null}'::jsonb
),
(
    'Savings',
    'Get Riwayat Kontribusi',
    'GET',
    '/api/savings/:id/contributions',
    true,
    'Mengambil seluruh riwayat mutasi / kontribusi pada tabungan',
    '{"Authorization": "Bearer <TOKEN>"}'::jsonb,
    null,
    '{"meta": {"code": 200, "status": "success", "message": "Berhasil mengambil riwayat kontribusi"}, "data": [{"id": "c9b8a7f6-1234-5678-90ab-cdef12345678", "user_name": "Budi Santoso", "amount": 500000, "date": "2026-04-10T00:00:00Z"}]}'::jsonb,
    '{"meta": {"code": 403, "status": "error", "message": "bukan anggota tabungan"}, "data": null}'::jsonb
);

-- 5. AGENDAS (Jadwal Bersama)
INSERT INTO api_documentations (category, name, method, endpoint, is_protected, description, headers, request_body, response_success, response_error) VALUES
(
    'Agendas',
    'Buat Agenda Baru',
    'POST',
    '/api/agendas',
    true,
    'Membuat jadwal / agenda baru',
    '{"Content-Type": "application/json", "Authorization": "Bearer <TOKEN>"}'::jsonb,
    '{"title": "Meeting Evaluasi Kuartal 2", "description": "Pembahasan target anggaran & belanja", "start_date": "2026-04-15T09:00:00Z", "end_date": "2026-04-15T11:00:00Z"}'::jsonb,
    '{"meta": {"code": 201, "status": "success", "message": "Agenda berhasil dibuat"}, "data": {"id": "b2c3d4e5-9999-8888-7777-666666666666", "title": "Meeting Evaluasi Kuartal 2", "description": "Pembahasan target anggaran & belanja", "start_date": "2026-04-15T09:00:00Z", "end_date": "2026-04-15T11:00:00Z"}}'::jsonb,
    '{"meta": {"code": 400, "status": "error", "message": "end_date harus lebih besar dari start_date"}, "data": null}'::jsonb
),
(
    'Agendas',
    'Get List Agenda Saya',
    'GET',
    '/api/agendas',
    true,
    'Mengambil daftar seluruh agenda di mana pengguna berpartisipasi',
    '{"Authorization": "Bearer <TOKEN>"}'::jsonb,
    null,
    '{"meta": {"code": 200, "status": "success", "message": "Berhasil mengambil agenda"}, "data": [{"id": "b2c3d4e5-9999-8888-7777-666666666666", "title": "Meeting Evaluasi Kuartal 2", "description": "Pembahasan target anggaran & belanja", "start_date": "2026-04-15T09:00:00Z", "end_date": "2026-04-15T11:00:00Z"}]}'::jsonb,
    '{"meta": {"code": 401, "status": "error", "message": "Unauthorized"}, "data": null}'::jsonb
),
(
    'Agendas',
    'Tambah Anggota Agenda',
    'POST',
    '/api/agendas/:id/members',
    true,
    'Menambahkan peserta baru ke dalam agenda (hanya Owner)',
    '{"Content-Type": "application/json", "Authorization": "Bearer <TOKEN>"}'::jsonb,
    '{"email": "rekan@gmail.com"}'::jsonb,
    '{"meta": {"code": 201, "status": "success", "message": "Anggota berhasil ditambahkan"}, "data": {"agenda_id": "b2c3d4e5-9999-8888-7777-666666666666", "user_id": "u5u5u5u5-5555-5555-5555-555555555555", "role": "member"}}'::jsonb,
    '{"meta": {"code": 403, "status": "error", "message": "hanya owner yang dapat menambah anggota agenda"}, "data": null}'::jsonb
);

-- 6. NOTIFICATIONS
INSERT INTO api_documentations (category, name, method, endpoint, is_protected, description, headers, request_body, response_success, response_error) VALUES
(
    'Notifications',
    'Get List Notifikasi',
    'GET',
    '/api/notifications',
    true,
    'Mengambil notifikasi masuk milik pengguna',
    '{"Authorization": "Bearer <TOKEN>"}'::jsonb,
    null,
    '{"meta": {"code": 200, "status": "success", "message": "Berhasil mengambil notifikasi"}, "data": [{"id": "n1n2n3n4-1111-2222-3333-444444444444", "user_id": "d3b07384-d113-4601-a757-5509e564d262", "title": "Undangan Tabungan", "message": "Anda ditambahkan ke tabungan Liburan ke Bali", "is_read": false, "created_at": "2026-04-10T10:00:00Z"}]}'::jsonb,
    '{"meta": {"code": 401, "status": "error", "message": "Unauthorized"}, "data": null}'::jsonb
),
(
    'Notifications',
    'Tandai Read Notifikasi',
    'PUT',
    '/api/notifications/:id/read',
    true,
    'Menandai pesan notifikasi tertentu sebagai sudah dibaca',
    '{"Authorization": "Bearer <TOKEN>"}'::jsonb,
    null,
    '{"meta": {"code": 200, "status": "success", "message": "Notifikasi ditandai dibaca"}, "data": null}'::jsonb,
    '{"meta": {"code": 404, "status": "error", "message": "Notifikasi tidak ditemukan"}, "data": null}'::jsonb
);
