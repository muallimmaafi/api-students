# API Students

REST API sederhana untuk mengelola data mahasiswa, dibangun menggunakan Go dan Fiber v2. Data disimpan di memori (in-memory), sehingga akan hilang setiap kali server dihentikan.

## Menjalankan Proyek

```bash
go mod tidy
go run .
```

Server akan berjalan di `http://localhost:3000`.

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