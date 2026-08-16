# API & Response Documentation — keuangan-api

Dokumentasi lengkap seluruh Endpoint API, Request Payload, Headers, dan Struktur Response JSON untuk proyek `keuangan-api`.

> **Base URL:** `http://localhost:8080/api`  
> **Format Standard Response:**  
> ```json
> {
>   "meta": {
>     "code": 200,
>     "status": "success", // atau "error"
>     "message": "Pesan deskriptif"
>   },
>   "data": { ... } // Objek data / Array / null
> }
> ```

---

## DAFTAR ISI ENDPOINT

1. [Otentikasi & Akun (Auth)](#1-otentikasi--akun-auth)
   - `POST /api/register` — Registrasi User
   - `POST /api/login` — Login User
   - `POST /api/forgot-password` — Lupa Password (Request OTP)
   - `POST /api/reset-password` — Reset Password (Verifikasi OTP)
2. [Kategori Transaksi (Categories)](#2-kategori-transaksi-categories)
   - `GET /api/categories` — List Kategori (Sistem + Custom)
   - `POST /api/categories` — Buat Kategori Custom
3. [Transaksi Keuangan (Transactions)](#3-transaksi-keuangan-transactions)
   - `POST /api/transactions` — Catat Transaksi Baru
   - `GET /api/transactions` — List Transaksi Saya
   - `PUT /api/transactions/:id` — Edit Transaksi
   - `DELETE /api/transactions/:id` — Hapus Transaksi (Soft Delete)
4. [Target Tabungan (Savings / Tabungan Bersama)](#4-target-tabungan-savings--tabungan-bersama)
   - `POST /api/savings` — Buat Target Tabungan Baru
   - `GET /api/savings` — List Target Tabungan Saya
   - `PUT /api/savings/:id` — Update Target Tabungan
   - `DELETE /api/savings/:id` — Hapus Target Tabungan
   - `POST /api/savings/:id/members` — Tambah Anggota Tabungan
   - `POST /api/savings/:id/contribute` — Setor Dana Tabungan
   - `POST /api/savings/:id/withdraw` — Tarik Dana Tabungan
   - `GET /api/savings/:id/contributions` — Riwayat Mutasi / Kontribusi
5. [Agenda & Jadwal Bersama (Agendas)](#5-agenda--jadwal-bersama-agendas)
   - `POST /api/agendas` — Buat Agenda Baru
   - `GET /api/agendas` — List Agenda Saya
   - `PUT /api/agendas/:id` — Update Agenda
   - `DELETE /api/agendas/:id` — Hapus Agenda
   - `POST /api/agendas/:id/members` — Tambah Peserta Agenda
6. [Notifikasi (Notifications)](#6-notifikasi-notifications)
   - `GET /api/notifications` — List Notifikasi Masuk
   - `PUT /api/notifications/:id/read` — Tandai Notifikasi Dibaca

---

## 1. Otentikasi & Akun (Auth)

### 1.1 `POST /register`
Mendaftarkan akun pengguna baru.

- **Headers:** `Content-Type: application/json`
- **Request Body:**
```json
{
  "name": "Budi Santoso",
  "email": "budi@gmail.com",
  "password": "password123"
}
```
- **Response Success (201 Created):**
```json
{
  "meta": {
    "code": 201,
    "status": "success",
    "message": "Registrasi berhasil"
  },
  "data": {
    "id": "d3b07384-d113-4601-a757-5509e564d262",
    "name": "Budi Santoso",
    "email": "budi@gmail.com"
  }
}
```
- **Response Error (409 Conflict):**
```json
{
  "meta": {
    "code": 409,
    "status": "error",
    "message": "email sudah terdaftar"
  },
  "data": null
}
```

---

### 1.2 `POST /login`
Mendapatkan Token JWT Bearer untuk mengakses endpoint yang dilindungi.

- **Headers:** `Content-Type: application/json`
- **Request Body:**
```json
{
  "email": "budi@gmail.com",
  "password": "password123"
}
```
- **Response Success (200 OK):**
```json
{
  "meta": {
    "code": 200,
    "status": "success",
    "message": "Login berhasil"
  },
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZDNiMDczODQtZDExMy00NjAxLWE3NTctNTUwOWU1NjRkMjYyIiwiZW1haWwiOiJidWRpQGdtYWlsLmNvbSIsImV4cCI6MTc3MTg2ODgwMH0..."
  }
}
```
- **Response Error (401 Unauthorized):**
```json
{
  "meta": {
    "code": 401,
    "status": "error",
    "message": "email atau password salah"
  },
  "data": null
}
```

---

### 1.3 `POST /forgot-password`
Mengirimkan OTP reset password ke email terdaftar.

- **Request Body:** `{"email": "budi@gmail.com"}`
- **Response Success (200 OK):**
```json
{
  "meta": {
    "code": 200,
    "status": "success",
    "message": "Jika email terdaftar, link reset password telah dikirim"
  },
  "data": null
}
```

---

### 1.4 `POST /reset-password`
Mereset password dengan memasukkan kode OTP dan password baru.

- **Request Body:**
```json
{
  "email": "budi@gmail.com",
  "code": "123456",
  "new_password": "newpassword123"
}
```
- **Response Success (200 OK):**
```json
{
  "meta": {
    "code": 200,
    "status": "success",
    "message": "Password berhasil direset, silakan login kembali"
  },
  "data": null
}
```

---

## 2. Kategori Transaksi (Categories)

> **Catatan:** Endpoint ini membutuhkan Header: `Authorization: Bearer <TOKEN>`

### 2.1 `GET /categories`
Mengambil semua kategori (kategori global sistem + kategori custom milik user yang login).

- **Headers:** `Authorization: Bearer <TOKEN>`
- **Response Success (200 OK):**
```json
{
  "meta": {
    "code": 200,
    "status": "success",
    "message": "Berhasil mengambil kategori"
  },
  "data": [
    {
      "id": "c1f7a091-8888-4c12-9901-5123456789ab",
      "name": "Gaji",
      "type": "income",
      "icon": "briefcase",
      "user_id": null
    },
    {
      "id": "d2f7a091-9999-4c12-9901-5123456789cd",
      "name": "Makan & Minum",
      "type": "expense",
      "icon": "utensils",
      "user_id": null
    },
    {
      "id": "e3f7a091-0000-4c12-9901-5123456789ef",
      "name": "Side Project",
      "type": "income",
      "icon": "code",
      "user_id": "d3b07384-d113-4601-a757-5509e564d262"
    }
  ]
}
```

---

### 2.2 `POST /categories`
Membuat kategori kustom baru.

- **Headers:** `Content-Type: application/json`, `Authorization: Bearer <TOKEN>`
- **Request Body:**
```json
{
  "name": "Investasi Crypto",
  "type": "income",
  "icon": "bitcoin"
}
```
- **Response Success (201 Created):**
```json
{
  "meta": {
    "code": 201,
    "status": "success",
    "message": "Kategori berhasil dibuat"
  },
  "data": {
    "id": "e3f7a091-0000-4c12-9901-5123456789ef",
    "user_id": "d3b07384-d113-4601-a757-5509e564d262",
    "name": "Investasi Crypto",
    "type": "income",
    "icon": "bitcoin"
  }
}
```

---

## 3. Transaksi Keuangan (Transactions)

### 3.1 `POST /transactions`
Mencatat transaksi keuangan baru.

- **Headers:** `Content-Type: application/json`, `Authorization: Bearer <TOKEN>`
- **Request Body:**
```json
{
  "category_id": "c1f7a091-8888-4c12-9901-5123456789ab",
  "amount": 5000000,
  "date": "2026-04-10",
  "notes": "Gaji Bulan April"
}
```
- **Response Success (201 Created):**
```json
{
  "meta": {
    "code": 201,
    "status": "success",
    "message": "Transaksi berhasil disimpan"
  },
  "data": {
    "id": "f4f7a091-1111-4c12-9901-512345678900",
    "user_id": "d3b07384-d113-4601-a757-5509e564d262",
    "category_id": "c1f7a091-8888-4c12-9901-5123456789ab",
    "amount": 5000000,
    "date": "2026-04-10T00:00:00Z",
    "notes": "Gaji Bulan April"
  }
}
```

---

### 3.2 `GET /transactions`
Mengambil semua transaksi pengguna yang login.

- **Headers:** `Authorization: Bearer <TOKEN>`
- **Response Success (200 OK):**
```json
{
  "meta": {
    "code": 200,
    "status": "success",
    "message": "Berhasil mengambil transaksi"
  },
  "data": [
    {
      "id": "f4f7a091-1111-4c12-9901-512345678900",
      "user_id": "d3b07384-d113-4601-a757-5509e564d262",
      "category_id": "c1f7a091-8888-4c12-9901-5123456789ab",
      "amount": 5000000,
      "date": "2026-04-10T00:00:00Z",
      "notes": "Gaji Bulan April",
      "category": {
        "id": "c1f7a091-8888-4c12-9901-5123456789ab",
        "name": "Gaji",
        "type": "income",
        "icon": "briefcase"
      }
    }
  ]
}
```

---

## 4. Target Tabungan (Savings / Tabungan Bersama)

### 4.1 `POST /savings`
Membuat target tabungan baru.

- **Headers:** `Content-Type: application/json`, `Authorization: Bearer <TOKEN>`
- **Request Body:**
```json
{
  "title": "Liburan ke Bali",
  "target_amount": 10000000,
  "deadline": "2026-12-31"
}
```
- **Response Success (201 Created):**
```json
{
  "meta": {
    "code": 201,
    "status": "success",
    "message": "Saving goal berhasil dibuat"
  },
  "data": {
    "id": "a1b2c3d4-1111-4444-8888-999999999999",
    "title": "Liburan ke Bali",
    "target_amount": 10000000,
    "current_amount": 0,
    "deadline": "2026-12-31T00:00:00Z"
  }
}
```

---

### 4.2 `POST /savings/:id/contribute`
Menyetorkan dana kontribusi ke tabungan bersama / pribadi.

- **Headers:** `Content-Type: application/json`, `Authorization: Bearer <TOKEN>`
- **Request Body:**
```json
{
  "amount": 500000,
  "date": "2026-04-10"
}
```
- **Response Success (201 Created):**
```json
{
  "meta": {
    "code": 201,
    "status": "success",
    "message": "Kontribusi berhasil disimpan"
  },
  "data": {
    "id": "c9b8a7f6-1234-5678-90ab-cdef12345678",
    "goal_id": "a1b2c3d4-1111-4444-8888-999999999999",
    "user_id": "d3b07384-d113-4601-a757-5509e564d262",
    "amount": 500000,
    "date": "2026-04-10T00:00:00Z"
  }
}
```

---

## 5. Agenda & Jadwal Bersama (Agendas)

### 5.1 `POST /agendas`
Membuat agenda kegiatan / acara keuangan baru.

- **Headers:** `Content-Type: application/json`, `Authorization: Bearer <TOKEN>`
- **Request Body:**
```json
{
  "title": "Meeting Evaluasi Anggaran",
  "description": "Pembahasan alokasi dana Q2",
  "start_date": "2026-04-15T09:00:00Z",
  "end_date": "2026-04-15T11:00:00Z"
}
```
- **Response Success (201 Created):**
```json
{
  "meta": {
    "code": 201,
    "status": "success",
    "message": "Agenda berhasil dibuat"
  },
  "data": {
    "id": "b2c3d4e5-9999-8888-7777-666666666666",
    "title": "Meeting Evaluasi Anggaran",
    "description": "Pembahasan alokasi dana Q2",
    "start_date": "2026-04-15T09:00:00Z",
    "end_date": "2026-04-15T11:00:00Z"
  }
}
```

---

## 6. Notifikasi (Notifications)

### 6.1 `GET /notifications`
Mengambil notifikasi milik user.

- **Headers:** `Authorization: Bearer <TOKEN>`
- **Response Success (200 OK):**
```json
{
  "meta": {
    "code": 200,
    "status": "success",
    "message": "Berhasil mengambil notifikasi"
  },
  "data": [
    {
      "id": "n1n2n3n4-1111-2222-3333-444444444444",
      "user_id": "d3b07384-d113-4601-a757-5509e564d262",
      "title": "Undangan Tabungan",
      "message": "Anda ditambahkan ke tabungan Liburan ke Bali",
      "is_read": false,
      "created_at": "2026-04-10T10:00:00Z"
    }
  ]
}
```

---

## TABEL DATABASE DOKUMENTASI API (`api_documentations`)

Untuk menyimpan seluruh dokumentasi ini ke dalam database PostgreSQL, jalankan skrip SQL pada file [api_documentation.sql](file:///d:/Dev/Golang/keuangan-api/migrations/api_documentation.sql).

Struktur tabel di database:
- `id` (UUID Primary Key)
- `category` (VARCHAR, contoh: 'Auth', 'Transactions')
- `name` (VARCHAR, nama endpoint)
- `method` (VARCHAR, contoh: 'POST', 'GET')
- `endpoint` (VARCHAR, contoh: '/api/login')
- `is_protected` (BOOLEAN, butuh token atau tidak)
- `description` (TEXT)
- `headers` (JSONB)
- `request_body` (JSONB)
- `response_success` (JSONB)
- `response_error` (JSONB)
