package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"api-students/app/model"
)

type PrestasiRepository interface {
	Create(ctx context.Context, p model.Prestasi) (model.Prestasi, error)
}

type prestasiPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPrestasiRepository(pool *pgxpool.Pool) PrestasiRepository {
	return &prestasiPostgresRepository{pool: pool}
}

func (r *prestasiPostgresRepository) Create(
	ctx context.Context, p model.Prestasi,
) (model.Prestasi, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO prestasi (nama_prestasi, id_student, juara)
		 VALUES ($1, $2, $3)
		 RETURNING id_prestasi, created_at`,
		p.NamaPrestasi, p.IDStudent, p.Juara,
	).Scan(&p.IDPrestasi, &p.CreatedAt)
	if err != nil {

		if isForeignKeyViolation(err) {
			return model.Prestasi{}, ErrNotFound
		}
		return model.Prestasi{}, fmt.Errorf("menyimpan prestasi: %w", err)
	}
	return p, nil
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23503"
	}
	return false
}