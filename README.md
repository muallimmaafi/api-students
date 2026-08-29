# API Students

REST API sederhana untuk mengelola data mahasiswa, dibangun menggunakan Go, Fiber v2, dan PostgreSQL (melalui pola Repository). Data tersimpan permanen di database, sehingga tidak hilang saat server dihentikan.

## Struktur Proyek
api-students/
├── app/
│ ├── model/
│ │ └── student.go # struct entitas, request, dan respons
│ └── repository/
│ └── student_repository.go # kontrak dan implementasi akses data
├── config/
│ └── env.go # memuat variabel environment
├── database/
│ └── postgres.go # koneksi dan connection pool
├── migrations/
│ ├── 001_create_students.sql
│ └── 002_fix_students_grade_type.sql
├── .env.example
├── main.go
├── handler.go
└── helper.go


## Skema Tabel

```sql
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    nim VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    grade DOUBLE PRECISION NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX students_nim_key ON students (nim);
CREATE INDEX students_name_lower_idx ON students (LOWER(name));
```

| Kolom | Tipe | Keterangan |
|---|---|---|
| id | SERIAL | Primary key, dibuat otomatis |
| nim | VARCHAR(20) | Wajib unik (dijaga UNIQUE INDEX) |
| name | VARCHAR(100) | Nama mahasiswa |
| grade | DOUBLE PRECISION | Nilai/IPK mahasiswa |
| is_active | BOOLEAN | Status aktif, default true |
| created_at | TIMESTAMPTZ | Waktu dibuat, default waktu sekarang |

## Cara Setup dari Nol

1. **Clone repositori dan masuk ke foldernya**
```bash
   git clone https://github.com/muallimmaafi/api-students.git
   cd api-students
```

2. **Buat database PostgreSQL**
```bash
   psql -U postgres -c "CREATE DATABASE db_api_students;"
```

3. **Jalankan migrasi secara berurutan**
```bash
   psql -U postgres -d db_api_students -f migrations/001_create_students.sql
   psql -U postgres -d db_api_students -f migrations/002_fix_students_grade_type.sql
```

4. **Salin `.env.example` menjadi `.env`, lalu isi sesuai konfigurasi lokal**
```bash
   cp .env.example .env
```

5. **Install dependency dan jalankan**
```bash
   go mod tidy
   go run .
```
   Server berjalan di `http://localhost:3000` (atau sesuai `APP_PORT` di `.env`).

## Variabel Environment

| Variabel | Keterangan | Contoh |
|---|---|---|
| APP_PORT | Port server HTTP | 3000 |
| DB_HOST | Host PostgreSQL | localhost |
| DB_PORT | Port PostgreSQL | 5432 |
| DB_USER | Username database | postgres |
| DB_PASSWORD | Password database | (isi sendiri) |
| DB_NAME | Nama database | db_api_students |
| DB_SSLMODE | Mode SSL koneksi | disable |
| DB_MAX_CONNS | Maksimum koneksi pada connection pool | 10 |

## Format Respons

Seluruh endpoint menggunakan amplop respons yang konsisten:

```json
{ "success": true, "message": "...", "data": { ... } }
```

Untuk daftar data, ditambahkan `meta`:
```json
{ "success": true, "message": "...", "data": [ ... ], "meta": { "page": 1, "limit": 10, "total": 3, "total_pages": 1 } }
```

Untuk kegagalan validasi, ditambahkan `errors` per field:
```json
{ "success": false, "message": "validasi gagal", "errors": { "name": "wajib diisi" } }
```

## Kontrak API

| Metode | Endpoint | Parameter | Contoh Body | Status | Contoh Respons |
|---|---|---|---|---|---|
| GET | `/api/v1/health` | – | – | 200, 503 | `{"success":true,"message":"server dan database berjalan","data":{"timestamp":"..."}}` |
| GET | `/api/v1/students` | Query: `page`, `limit` (maks 100), `search`, `sort` (`id`\|`nim`\|`name`\|`grade`\|`created_at`), `order` (`asc`\|`desc`), `is_active` (`true`\|`false`) | – | 200 | `{"success":true,"data":[...],"meta":{"page":1,"limit":10,"total":3,"total_pages":1}}` |
| GET | `/api/v1/students/:id` | Path: `id` (angka) | – | 200, 400, 404 | `{"success":true,"data":{"id":1,"nim":"2024001","name":"Ali",...}}` |
| POST | `/api/v1/students` | – | `{"nim":"2024001","name":"Ali","grade":85.5}` | 201, 400, 409, 415, 422 | `{"success":true,"message":"mahasiswa berhasil dibuat","data":{...}}` (header `Location` disertakan) |
| PUT | `/api/v1/students/:id` | Path: `id`. Semua field wajib | `{"nim":"2024001","name":"Ali Ramadhan","grade":88,"is_active":false}` | 200, 400, 404, 409, 415, 422 | `{"success":true,"message":"mahasiswa berhasil diganti seluruhnya","data":{...}}` |
| PATCH | `/api/v1/students/:id` | Path: `id`. Field opsional, hanya yang dikirim yang diubah | `{"is_active":true}` | 200, 400, 404, 409, 415, 422 | `{"success":true,"message":"mahasiswa berhasil diperbarui sebagian","data":{...}}` |
| DELETE | `/api/v1/students/:id` | Path: `id` (angka) | – | 204, 400, 404 | *(tanpa body)* |

## Status HTTP yang Digunakan

| Status | Arti |
|---|---|
| 200 | Berhasil (GET, PUT, PATCH) |
| 201 | Data baru berhasil dibuat |
| 204 | Berhasil dihapus, tanpa body |
| 400 | Permintaan salah bentuk (id bukan angka, body bukan JSON valid) |
| 404 | Data tidak ditemukan |
| 409 | Bertentangan dengan data yang ada (NIM duplikat) |
| 415 | Content-Type bukan `application/json` |
| 422 | Validasi isi gagal |
| 500 | Kesalahan tak terduga pada server/database |
| 503 | Database tidak dapat dihubungi (khusus `/health`) |

## Sumber Bantuan

Dibantu Claude (Anthropic) untuk: struktur query migrasi SQL, pola connection pool (`pgxpool`), implementasi pola Repository (interface + implementasi Postgres), pemetaan error database ke status HTTP, dan penyusunan dokumentasi ini.