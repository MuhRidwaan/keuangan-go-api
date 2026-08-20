# 🚀 Proposal Cetak Biru Fitur Web Application (Next.js) — Finance & Agenda Platform

Dokumen ini berisi rancangan fitur lengkap untuk pengembangan **Aplikasi Web Keuangan & Agenda** berbasis **Next.js (App Router)**. Rancangan ini mencakup 100% fitur dari aplikasi mobile (Flutter) ditambah **fitur-fitur eksklusif versi Web** untuk memberikan pengalaman pengolahan data keuangan & agenda yang lebih kaya, interaktif, dan komprehensif.

---

## 🛠️ Teknologi & Desain Sistem (Tech Stack)

| Komponen | Teknologi yang Direkomondesikan |
| :--- | :--- |
| **Framework** | **Next.js 14+ (App Router)** |
| **Bahasa** | **TypeScript** |
| **State Management** | **Zustand** atau **TanStack React Query** |
| **Styling & Theme** | **TailwindCSS** + **Shadcn UI** (Modern Glassmorphism, Dark/Light Mode) |
| **Visualisasi Data** | **Recharts** / **Chart.js** (Grafik Pemasukan, Pengeluaran & Trends) |
| **Kalender Interaktif** | **FullCalendar** / **React Big Calendar** |
| **Data Grid & Table** | **TanStack Table (React Table)** (Sorting, Filtering, Pagination, Multi-select) |
| **Export Data** | **jsPDF / pdfmake** (PDF Report) & **XLSX** (Excel Export) |

---

## 📋 Matriks Komparasi Fitur (Apps vs Next.js Web)

| Modul Fitur | Mobile Apps (Flutter) | Next.js Web Application (Enhanced) |
| :--- | :---: | :---: |
| **Dashboard** | Mobile Card Summary | **Interactive Dashboard Analytics (Charts & Heatmap)** |
| **Manajemen Transaksi** | ListView Scroll + Date Filter | **Advanced Data Grid Table (Multi-select, Filter, Export PDF/Excel/CSV)** |
| **Financial Reports** | Chart Progress Sederhana | **Interactive Chart Analytics (Line, Pie, Bar, MoM Growth)** |
| **Budget Bulanan** | Progress Bar Limit | **Budget Visual Manager (Alert Thresholds, Category Breakdown)** |
| **Agenda & Jadwal** | List View Agenda Cards | **Full-Screen Interactive Calendar (Drag & Drop, Month/Week/Day View)** |
| **Tabungan Bersama** | Card Goal + Setor/Tarik | **Goal Tracker Board + Log Kontribusi Detail** |
| **Riwayat Kontak** | Quick Suggest Chips | **Contact Management & Quick Invite Center** |
| **Theme & UI** | Mobile Responsive | **Responsive Sidebar, Dark/Light Mode, Glassmorphism Dashboard** |

---

## 💡 Rincian Fitur Lengkap Versi Web (Next.js)

### 1. 📊 Interactive Analytics Dashboard (Beranda Utama)
*Semua fitur mobile + Analisis Grafik Visual:*
- **Hero Balance & Financial Summary**: Menampilkan Total Saldo, Pemasukan Bulan Ini, Pengeluaran Bulan Ini, dan **Pengeluaran Hari Ini**.
- **Monthly Income vs Expense Chart**: Grafik batang/garis perbandingan arus kas pemasukan dan pengeluaran per bulan.
- **Spending Category Distribution**: Grafik Pie / Donut interaktif distribusi pengeluaran per kategori.
- **Financial Activity Heatmap**: Visualisasi intensitas pengeluaran harian dalam bentuk kalender hijau/merah (mirip GitHub activity graph).
- **Recent Transactions & Upcoming Agendas Widget**: Quick-access widget transaksi terakhir dan 3 agenda terdekat yang akan datang.

---

### 2. 📑 Advanced Transactions Data Table (Pusat Transaksi)
*Membuat pengelolaan ribuan transaksi menjadi sangat cepat & fleksibel:*
- **Data Table Interaktif (TanStack Table)**:
  - Pencarian instan (*Instant Search*) berdasarkan deskripsi/catatan.
  - Sorting kolom (Tanggal, Kategori, Jumlah).
  - Filter multi-kategori & filter rentang tanggal (*Date Range Picker*).
  - Pagination dinamis (10, 25, 50, 100 baris per halaman).
- **Batch Operations**: Multi-select transaksi untuk hapus masal sekaligus.
- **Multi-Format Export**:
  - Export **CSV** (seperti di mobile).
  - Export **Excel (.xlsx)** dengan format sel angka & tanggal otomatis.
  - Export **PDF Report** dengan kop laporan resmi, grafik summary, dan tabel transaksi.
- **Quick CRUD Modal**: Tambah & Edit transaksi tanpa pindah halaman (*Modal Sheet / Dialog*).

---

### 3. 📈 Financial Reports & Insights (Laporan Keuangan)
- **Period Filter**: Filter laporan bulanan, tahunan, atau kustom rentang tanggal.
- **Net Cash Flow Analysis**: Kalkulasi Saldo Bersih (*Net Income*), Rata-rata Pengeluaran Harian, dan Estimasi Sisa Hari Hemat.
- **Category Deep Dive**: Detail breakdown tiap kategori dilengkapi dengan progress bar persentase terhadap total pengeluaran.
- **Comparison Comparison (Month-over-Month)**: Perbandingan persentase kenaikan/penurunan pengeluaran dibanding bulan sebelumnya.

---

### 4. 🎯 Budget Management & Limits (Kelola Budget Bulanan)
- **Monthly Budget Setup**: Buat budget bulanan per kategori atau budget global.
- **Visual Limit Bar**: Progress bar berwarna dinamis:
  - 🟢 **Hijau**: Aman (< 80% limit).
  - 🟡 **Kuning**: Peringatan (80% - 99% limit).
  - 🔴 **Merah**: Melebihi limit (≥ 100% budget exceeded).
- **Smart Alert System**: Banner peringatan otomatis di top dashboard saat budget mendekati atau melewati batas.

---

### 5. 📅 Full Calendar Agenda & Task Management (Agenda Interaktif)
- **Multiple Calendar Views**: Tampilan **Bulan (Month)**, **Minggu (Week)**, **Hari (Day)**, dan **List View**.
- **Drag & Drop Re-scheduling**: Geser agenda untuk memperbarui tanggal/waktu mulai & selesai secara instan.
- **Status Toggle**: Ubah status agenda (**Pending / Terlaksana**) dengan satu klik checkmark.
- **Collaborative Agenda Details**:
  - Tampilan daftar seluruh anggota agenda beserta role (**Owner** / **Anggota**).
  - Undang anggota baru via email dengan *auto-complete kontak*.

---

### 6. 💰 Savings & Financial Goals (Tabungan Bersama / Target)
- **Goal Progress Cards**: Tampilan kartu visual progres pencapaian tabungan (Persentase target tercapai & sisa kekurangannya).
- **Contribute & Withdraw Operations**: Form Setor dan Tarik tabungan dengan pencatatan riwayat langsung.
- **Member Contribution Log**: Tabel riwayat kontribusi anggota pada tabungan bersama (Siapa menyetor berapa dan kapan).

---

### 7. 👥 Contact Management & Invites (Pusat Kontak & Kolaborasi)
- **Recent Contacts List**: Daftar kontak teman/anggota keluarga yang pernah diundang atau diundang bersama.
- **One-Click Invite**: Mengundang anggota ke Agenda / Tabungan tanpa mengetik ulang email (cukup klik nama dari daftar kontak).

---

### 8. 🌓 UI/UX & System Preferences
- **Dark Mode & Light Mode**: Dukungan tema gelap dan terang dengan peralihan animasi yang halus.
- **Responsive Layout**: Sidebar navigasi yang dapat di-collapse untuk pengalaman terbaik di layar Laptop, Tablet, maupun HP.
- **JWT Token Management**: Penyimpanan token JWT aman di Cookies/LocalStorage + auto redirect jika token expired.

---

## 🏗️ Struktur Direktori Proyek Next.js yang Diusulkan

```text
keuangan-web/
├── app/
│   ├── (auth)/
│   │   ├── login/
│   │   └── register/
│   ├── (dashboard)/
│   │   ├── dashboard/
│   │   ├── transactions/
│   │   ├── budget/
│   │   ├── reports/
│   │   ├── agendas/
│   │   ├── savings/
│   │   └── settings/
│   ├── layout.tsx
│   └── page.tsx
├── components/
│   ├── ui/               # Component Shadcn (Button, Dialog, Table, Input)
│   ├── dashboard/        # Widget Dashboard (Charts, Cards, Heatmap)
│   ├── transactions/     # Component Transaksi (DataTable, FilterBar)
│   ├── agendas/          # Component Agenda (FullCalendar Integration)
│   └── layout/           # Sidebar, Navbar, UserDropdown
├── lib/
│   ├── api/              # Dio / Axios Client untuk terhubung ke keuangan-api
│   ├── utils.ts          # Formatters (Rupiah, Format Tanggal)
│   └── constants.ts
└── types/                # TypeScript Interfaces (Transaction, Budget, Agenda)
```

---

## 🚀 Rekomendasi Langkah Eksekusi

1. **Inisialisasi Proyek Next.js**: Menggunakan `npx create-next-app@latest ./ --typescript --tailwind --app`.
2. **Integrasi Design System & UI Components**: Memasang `shadcn-ui`, `lucide-react`, dan setup tema dark/light.
3. **Setup API Client**: Membuat `axios` instance dengan interceptor `Authorization: Bearer <token>` terhubung ke backend Golang `https://keuangan-go-api.vercel.app` (atau local `http://localhost:8080`).
4. **Implementasi Halaman per Halaman**: Dimulai dari Auth (Login/Register) ➔ Dashboard ➔ Transaksi ➔ Budget ➔ Laporan ➔ Agenda ➔ Tabungan.
