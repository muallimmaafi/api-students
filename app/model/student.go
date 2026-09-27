package model

import "time"

type Student struct {
	ID        int       `json:"id"`
	NIM       string    `json:"nim"`
	Name      string    `json:"name"`
	Grade     float64   `json:"grade"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	OwnerID   *int      `json:"owner_id"`
}

// POST — semua field wajib
type CreateStudentRequest struct {
	NIM   string  `json:"nim" validate:"required,notblank"`
	Name  string  `json:"name" validate:"required,notblank"`
	Grade float64 `json:"grade" validate:"gte=0,lte=100"`
}

// PUT — ganti seluruh isi, semua field wajib
type ReplaceStudentRequest struct {
	NIM      string  `json:"nim" validate:"required,notblank"`
	Name     string  `json:"name" validate:"required,notblank"`
	Grade    float64 `json:"grade" validate:"gte=0,lte=100"`
	IsActive bool    `json:"is_active"`
}

// PATCH — field pointer agar bisa membedakan
// "tidak dikirim" (nil) dan "dikirim bernilai kosong"
type PatchStudentRequest struct {
	NIM      *string  `json:"nim,omitempty" validate:"omitnil,notblank"`
	Name     *string  `json:"name,omitempty" validate:"omitnil,notblank"`
	Grade    *float64 `json:"grade,omitempty" validate:"omitnil,gte=0,lte=100"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// Amplop baku untuk semua respons
type WebResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ListQuery menampung parameter query untuk endpoint daftar:
// pencarian, penyaringan, pengurutan, dan paginasi.
type ListQuery struct {
	Page     int
	Limit    int
	Search   string
	Sort     string
	Order    string
	IsActive *bool
}

// Offset menghitung berapa baris yang dilewati untuk halaman ini.
// Dipakai nanti oleh repository untuk LIMIT/OFFSET di SQL.
func (q ListQuery) Offset() int {
	return (q.Page - 1) * q.Limit
}

// Cursor menyimpan posisi terakhir pada pagination.
// Kombinasi created_at + id digunakan agar urutan selalu unik.
type Cursor struct {
	CreatedAt time.Time
	ID        int
}

// CursorQuery menampung parameter untuk pagination berbasis cursor.
type CursorQuery struct {
	Limit    int
	Search   string
	IsActive *bool
	After    *Cursor
}

// CursorMeta adalah metadata pagination berbasis cursor.
type CursorMeta struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}
