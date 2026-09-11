package service

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
)

type PrestasiService struct {
	repo repository.PrestasiRepository
}

func NewPrestasiService(repo repository.PrestasiRepository) *PrestasiService {
	return &PrestasiService{repo: repo}
}

func (s *PrestasiService) Create(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.CreatePrestasiRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	req.NamaPrestasi = strings.TrimSpace(req.NamaPrestasi)
	req.Juara = strings.TrimSpace(req.Juara)

	errs := map[string]string{}
	if req.NamaPrestasi == "" {
		errs["nama_prestasi"] = "wajib diisi"
	}
	if req.IDStudent <= 0 {
		errs["id_student"] = "wajib diisi dengan id mahasiswa yang valid"
	}
	if req.Juara == "" {
		errs["juara"] = "wajib diisi"
	}
	if len(errs) > 0 {
		return helper.FailValidation(c, errs)
	}

	prestasi, err := s.repo.Create(ctx, model.Prestasi{
		NamaPrestasi: req.NamaPrestasi,
		IDStudent:    req.IDStudent,
		Juara:        req.Juara,
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound,
				"mahasiswa dengan id_student tersebut tidak ditemukan")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "gagal menyimpan prestasi")
	}

	return helper.Success(c, fiber.StatusCreated, "prestasi berhasil ditambahkan", prestasi)
}