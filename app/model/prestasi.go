package model

import "time"

type Prestasi struct {
	IDPrestasi   int       `json:"id_prestasi"`
	NamaPrestasi string    `json:"nama_prestasi"`
	IDStudent    int       `json:"id_student"`
	Juara        string    `json:"juara"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreatePrestasiRequest struct {
	NamaPrestasi string `json:"nama_prestasi"`
	IDStudent    int    `json:"id_student"`
	Juara        string `json:"juara"`
}