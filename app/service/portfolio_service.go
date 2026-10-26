package service

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
)

type PortfolioService struct {
	portfolioRepo *repository.PortfolioRepository
	studentRepo   repository.StudentRepository
}

func NewPortfolioService(
	portfolioRepo *repository.PortfolioRepository,
	studentRepo repository.StudentRepository,
) *PortfolioService {
	return &PortfolioService{
		portfolioRepo: portfolioRepo,
		studentRepo:   studentRepo,
	}
}

// Get mengambil portofolio milik mahasiswa yang terdaftar di PostgreSQL.
func (s *PortfolioService) Get(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	studentID, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id harus berupa angka positif",
		)
	}

	// Pastikan mahasiswa terdaftar di PostgreSQL.
	_, err := s.studentRepo.FindByID(ctx, studentID)
	if err != nil {
		return translatePortfolioError(c, err, "gagal memeriksa mahasiswa")
	}

	portfolio, err := s.portfolioRepo.FindByStudentID(ctx, studentID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(
				c,
				fiber.StatusNotFound,
				"portofolio belum tersedia",
			)
		}

		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal mengambil portofolio",
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"portofolio berhasil ditemukan",
		portfolio,
	)
}

func (s *PortfolioService) Put(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	studentID, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest,
			"id harus berupa angka positif")
	}

	// Pastikan mahasiswa terdaftar di PostgreSQL.
	_, err := s.studentRepo.FindByID(ctx, studentID)
	if err != nil {
		return translatePortfolioError(c, err,
			"gagal memeriksa mahasiswa")
	}

	var req model.Portfolio
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest,
			"body harus berupa JSON yang valid")
	}

	req.Headline = strings.TrimSpace(req.Headline)
	req.Bio = strings.TrimSpace(req.Bio)

	// ID mahasiswa selalu mengikuti parameter URL,
	// bukan nilai student_id dari body request.
	req.StudentID = studentID

	if req.Headline == "" {
		return helper.Fail(c, fiber.StatusBadRequest,
			"headline wajib diisi")
	}

	if req.Bio == "" {
		return helper.Fail(c, fiber.StatusBadRequest,
			"bio wajib diisi")
	}

	if len(req.Projects) > 20 {
		return helper.Fail(c, fiber.StatusBadRequest,
			"portofolio maksimal memiliki 20 proyek")
	}

	for i := range req.Projects {
		req.Projects[i].Title =
			strings.TrimSpace(req.Projects[i].Title)

		if req.Projects[i].Title == "" {
			return helper.Fail(c, fiber.StatusBadRequest,
				"judul setiap proyek wajib diisi")
		}

		for j := range req.Projects[i].Technologies {
			req.Projects[i].Technologies[j] =
				strings.TrimSpace(req.Projects[i].Technologies[j])
		}
	}

	if err := s.portfolioRepo.Replace(ctx, req); err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError,
			"gagal menyimpan portofolio")
	}

	return helper.Success(c, fiber.StatusOK,
		"portofolio berhasil disimpan", req)
}

func (s *PortfolioService) AddProject(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	studentID, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(c, fiber.StatusBadRequest,
			"id harus berupa angka positif")
	}

	// Pastikan mahasiswa terdaftar di PostgreSQL.
	_, err := s.studentRepo.FindByID(ctx, studentID)
	if err != nil {
		return translatePortfolioError(c, err,
			"gagal memeriksa mahasiswa")
	}

	var project model.Project
	if err := c.BodyParser(&project); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest,
			"body harus berupa JSON yang valid")
	}

	project.Title = strings.TrimSpace(project.Title)
	project.Link = strings.TrimSpace(project.Link)

	if project.Title == "" {
		return helper.Fail(c, fiber.StatusBadRequest,
			"judul proyek wajib diisi")
	}

	if len(project.Technologies) == 0 {
		return helper.Fail(c, fiber.StatusBadRequest,
			"minimal satu teknologi harus diisi")
	}

	for i := range project.Technologies {
		project.Technologies[i] =
			strings.TrimSpace(project.Technologies[i])

		if project.Technologies[i] == "" {
			return helper.Fail(c, fiber.StatusBadRequest,
				"nama teknologi tidak boleh kosong")
		}
	}

	if project.Year < 1900 || project.Year > 2100 {
		return helper.Fail(c, fiber.StatusBadRequest,
			"tahun proyek tidak valid")
	}

	// Pastikan portofolio sudah tersedia.
	_, err = s.portfolioRepo.FindByStudentID(ctx, studentID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound,
				"portofolio belum tersedia")
		}

		return helper.Fail(c, fiber.StatusInternalServerError,
			"gagal memeriksa portofolio")
	}

	// Simpan proyek dan dapatkan ID yang dibuat repository.
	projectID, err := s.portfolioRepo.AddProject(
		ctx,
		studentID,
		project,
	)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrPortfolioNotFound):
			return helper.Fail(c, fiber.StatusNotFound,
				"portofolio tidak ditemukan")

		case errors.Is(err, repository.ErrProjectLimit):
			return helper.Fail(c, fiber.StatusConflict,
				"batas 20 proyek tercapai")

		default:
			return helper.Fail(c, fiber.StatusInternalServerError,
				"gagal menambahkan proyek")
		}
	}

	// Sertakan ID proyek agar bisa digunakan untuk endpoint DELETE.
	project.ID = projectID

	return helper.Success(c, fiber.StatusCreated,
		"proyek berhasil ditambahkan",
		map[string]any{
			"student_id": studentID,
			"project":    project,
		})
}

func (s *PortfolioService) DeleteProject(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	studentID, valid := helper.ParamID(c)
	if !valid {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id harus berupa angka positif",
		)
	}

	projectID := strings.TrimSpace(c.Params("projectId"))
	if projectID == "" {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"projectId wajib diisi",
		)
	}

	_, err := s.studentRepo.FindByID(ctx, studentID)
	if err != nil {
		return translatePortfolioError(
			c,
			err,
			"gagal memeriksa mahasiswa",
		)
	}

	err = s.portfolioRepo.DeleteProject(ctx, studentID, projectID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound),
			errors.Is(err, repository.ErrPortfolioNotFound):
			return helper.Fail(
				c,
				fiber.StatusNotFound,
				"portofolio atau proyek tidak ditemukan",
			)
		default:
			return helper.Fail(
				c,
				fiber.StatusInternalServerError,
				"gagal menghapus proyek",
			)
		}
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"proyek berhasil dihapus",
		nil,
	)
}

func (s *PortfolioService) TechnologyStats(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	stats, err := s.portfolioRepo.TechnologyStats(ctx)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError,
			"gagal mengambil statistik teknologi")
	}

	return helper.Success(c, fiber.StatusOK,
		"statistik teknologi berhasil diambil", stats)
}

func translatePortfolioError(
	c *fiber.Ctx,
	err error,
	generalMessage string,
) error {
	if errors.Is(err, repository.ErrNotFound) {
		return helper.Fail(
			c,
			fiber.StatusNotFound,
			"mahasiswa tidak ditemukan",
		)
	}

	return helper.Fail(
		c,
		fiber.StatusInternalServerError,
		generalMessage,
	)
}

// Placeholder agar model tetap menjadi bagian dari service
// ketika kita melanjutkan implementasi endpoint berikutnya.
var _ model.Portfolio
