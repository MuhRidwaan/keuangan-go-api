# API Testing Guide — keuangan-api

Base URL: `http://localhost:8080`

Gunakan **Postman**, **Thunder Client** (VS Code extension), atau **curl**.

---

## ALUR TESTING YANG BENAR

```
1. Register         → dapat user baru
2. Login            → dapat TOKEN
3. Pakai TOKEN      → untuk semua request selanjutnya
```

> Semua endpoint selain Register & Login wajib pakai header:
> `Authorization: Bearer <TOKEN_DARI_LOGIN>`

---

## 1. AUTH

### POST /api/register
```
POST http://localhost:8080/api/register
Content-Type: application/json

{
  "name": "Budi Santoso",
  "email": "budi@gmail.com",
  "password": "password123"
}
```

**Response sukses (201):**
```json
{
  "meta": { "code": 201, "status": "success", "message": "Registrasi berhasil" },
  "data": {
    "id": "uuid-user",
    "name": "Budi Santoso",
    "email": "budi@gmail.com"
  }
}
```

---

### POST /api/login
```
POST http://localhost:8080/api/login
Content-Type: application/json

{
  "email": "budi@gmail.com",
  "password": "password123"
}
```

**Response sukses (200):**
```json
{
  "meta": { "code": 200, "status": "success", "message": "Login berhasil" },
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

> ⚠️ **Simpan token ini!** Pakai di semua request berikutnya sebagai:
> `Authorization: Bearer eyJhbGci...`

---

## 2. CATEGORIES

### GET /api/categories — Ambil semua kategori
```
GET http://localhost:8080/api/categories
Authorization: Bearer <TOKEN>
```

**Response sukses (200):**
```json
{
  "meta": { "code": 200, "status": "success", "message": "Berhasil mengambil kategori" },
  "data": [
    { "id": "uuid", "name": "Gaji", "type": "income", "icon": "briefcase", "user_id": null },
    { "id": "uuid", "name": "Makan & Minum", "type": "expense", "icon": "utensils", "user_id": null }
  ]
}
```

---

### POST /api/categories — Buat kategori custom
```
POST http://localhost:8080/api/categories
Authorization: Bearer <TOKEN>
Content-Type: application/json

{
  "name": "Side Project",
  "type": "income",
  "icon": "code"
}
```

**Response sukses (201):**
```json
{
  "meta": { "code": 201, "status": "success", "message": "Kategori berhasil dibuat" },
  "data": {
    "id": "uuid-kategori",
    "user_id": "uuid-user-kamu",
    "name": "Side Project",
    "type": "income",
    "icon": "code"
  }
}
```

---

## 3. TRANSACTIONS

### POST /api/transactions — Catat transaksi baru
> Butuh `category_id` dari response GET /api/categories

```
POST http://localhost:8080/api/transactions
Authorization: Bearer <TOKEN>
Content-Type: application/json

{
  "category_id": "uuid-kategori-dari-GET-categories",
  "amount": 5000000,
  "date": "2026-04-10",
  "notes": "Gaji bulan April"
}
```

**Response sukses (201):**
```json
{
  "meta": { "code": 201, "status": "success", "message": "Transaksi berhasil disimpan" },
  "data": {
    "id": "uuid-transaksi",
    "user_id": "uuid-user",
    "category_id": "uuid-kategori",
    "amount": 5000000,
    "date": "2026-04-10T00:00:00Z",
    "notes": "Gaji bulan April"
  }
}
```

---

### GET /api/transactions — Ambil semua transaksi milik user
```
GET http://localhost:8080/api/transactions
Authorization: Bearer <TOKEN>
```

---

## 4. SAVINGS (Tabungan Bersama)

### POST /api/savings — Buat saving goal baru
```
POST http://localhost:8080/api/savings
Authorization: Bearer <TOKEN>
Content-Type: application/json

{
  "title": "Liburan ke Jepang",
  "target_amount": 30000000,
  "deadline": "2026-12-31"
}
```

**Response sukses (201):**
```json
{
  "meta": { "code": 201, "status": "success", "message": "Saving goal berhasil dibuat" },
  "data": {
    "id": "uuid-goal",
    "title": "Liburan ke Jepang",
    "target_amount": 30000000,
    "current_amount": 0,
    "deadline": "2026-12-31T00:00:00Z"
  }
}
```

---

### GET /api/savings — Ambil semua goal milik user
```
GET http://localhost:8080/api/savings
Authorization: Bearer <TOKEN>
```

---

### POST /api/savings/:id/members — Tambah anggota (hanya owner)
> Daftarkan dulu user kedua via `/api/register`, lalu pakai ID-nya

```
POST http://localhost:8080/api/savings/uuid-goal/members
Authorization: Bearer <TOKEN_OWNER>
Content-Type: application/json

{
  "user_id": "uuid-user-yang-mau-ditambah"
}
```

**Response sukses (201):**
```json
{
  "meta": { "code": 201, "status": "success", "message": "Anggota berhasil ditambahkan" },
  "data": {
    "goal_id": "uuid-goal",
    "user_id": "uuid-user-baru",
    "role": "member"
  }
}
```

**Response gagal — bukan owner (403):**
```json
{
  "meta": { "code": 403, "status": "error", "message": "hanya owner yang dapat menambah anggota" }
}
```

---

### POST /api/savings/:id/contribute — Setor kontribusi
> Bisa dilakukan oleh semua member (owner maupun member biasa)

```
POST http://localhost:8080/api/savings/uuid-goal/contribute
Authorization: Bearer <TOKEN>
Content-Type: application/json

{
  "amount": 1000000,
  "date": "2026-04-10"
}
```

**Response sukses (201):**
```json
{
  "meta": { "code": 201, "status": "success", "message": "Kontribusi berhasil disimpan" },
  "data": {
    "id": "uuid-kontribusi",
    "goal_id": "uuid-goal",
    "user_id": "uuid-user",
    "amount": 1000000,
    "date": "2026-04-10T00:00:00Z"
  }
}
```

> ✅ Setelah ini, `current_amount` di saving_goals otomatis bertambah 1.000.000

---

## 5. AGENDAS (Jadwal Bersama)

### POST /api/agendas — Buat agenda baru
> Format tanggal: RFC3339 → `2026-04-15T09:00:00Z`

```
POST http://localhost:8080/api/agendas
Authorization: Bearer <TOKEN>
Content-Type: application/json

{
  "title": "Meeting Bulanan Tim",
  "description": "Review progress dan planning sprint berikutnya",
  "start_date": "2026-04-15T09:00:00Z",
  "end_date": "2026-04-15T11:00:00Z"
}
```

**Response sukses (201):**
```json
{
  "meta": { "code": 201, "status": "success", "message": "Agenda berhasil dibuat" },
  "data": {
    "id": "uuid-agenda",
    "title": "Meeting Bulanan Tim",
    "description": "Review progress dan planning sprint berikutnya",
    "start_date": "2026-04-15T09:00:00Z",
    "end_date": "2026-04-15T11:00:00Z"
  }
}
```

---

### GET /api/agendas — Ambil semua agenda milik user
```
GET http://localhost:8080/api/agendas
Authorization: Bearer <TOKEN>
```

---

### POST /api/agendas/:id/members — Tambah anggota agenda (hanya owner)
```
POST http://localhost:8080/api/agendas/uuid-agenda/members
Authorization: Bearer <TOKEN_OWNER>
Content-Type: application/json

{
  "user_id": "uuid-user-yang-mau-ditambah"
}
```

---

## CONTOH ERROR RESPONSES

### 400 — Bad Request (validasi gagal)
```json
{
  "meta": { "code": 400, "status": "error", "message": "Key: 'RegisterInput.Email' Error:Field validation for 'Email' failed on the 'email' tag" },
  "data": null
}
```

### 401 — Unauthorized (token tidak ada / expired)
```json
{
  "meta": { "code": 401, "status": "error", "message": "Token tidak ditemukan" },
  "data": null
}
```

### 403 — Forbidden (bukan owner)
```json
{
  "meta": { "code": 403, "status": "error", "message": "hanya owner yang dapat menambah anggota" },
  "data": null
}
```

### 409 — Conflict (email sudah terdaftar / user sudah jadi member)
```json
{
  "meta": { "code": 409, "status": "error", "message": "email sudah terdaftar" },
  "data": null
}
```

---

## URUTAN TESTING LENGKAP (Happy Path)

```
1.  POST /api/register          → daftar user A
2.  POST /api/register          → daftar user B (untuk test multi-user)
3.  POST /api/login             → login user A, simpan TOKEN_A
4.  POST /api/login             → login user B, simpan TOKEN_B
5.  GET  /api/categories        → lihat kategori sistem (pakai TOKEN_A)
6.  POST /api/categories        → buat kategori custom (pakai TOKEN_A)
7.  POST /api/transactions      → catat transaksi (pakai TOKEN_A + category_id dari step 5)
8.  GET  /api/transactions      → lihat transaksi user A
9.  POST /api/savings           → buat saving goal (pakai TOKEN_A, jadi owner)
10. GET  /api/savings           → lihat goals user A
11. POST /api/savings/:id/members   → tambah user B sebagai member (pakai TOKEN_A)
12. POST /api/savings/:id/contribute → setor dari user A (pakai TOKEN_A)
13. POST /api/savings/:id/contribute → setor dari user B (pakai TOKEN_B)
14. POST /api/agendas           → buat agenda (pakai TOKEN_A)
15. POST /api/agendas/:id/members   → tambah user B ke agenda (pakai TOKEN_A)
16. GET  /api/agendas           → user B lihat agendanya (pakai TOKEN_B)
```
