package helper

import (
	"encoding/csv"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
)

func NegotiateStudents(c *fiber.Ctx, students []model.Student, meta *model.CursorMeta) error {
	accept := c.Get("Accept")

	switch accept {
	case "", "*/*", "application/json":
		return SuccessList(
			c,
			"daftar mahasiswa berhasil diambil",
			students,
			meta,
		)

	case "text/csv":
		return WriteStudentsCSV(c, students)

	default:
		return Fail(
			c,
			fiber.StatusNotAcceptable,
			"format response tidak didukung",
		)
	}
}

func WriteStudentsCSV(c *fiber.Ctx, students []model.Student) error {
	c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename=students.csv")

	w := csv.NewWriter(c.Response().BodyWriter())

	if err := w.Write([]string{
		"id",
		"nim",
		"name",
		"grade",
		"is_active",
		"created_at",
		"owner_id",
	}); err != nil {
		return err
	}

	for _, s := range students {
		ownerID := ""
		if s.OwnerID != nil {
			ownerID = fmt.Sprintf("%d", *s.OwnerID)
		}

		if err := w.Write([]string{
			fmt.Sprintf("%d", s.ID),
			s.NIM,
			s.Name,
			fmt.Sprintf("%.2f", s.Grade),
			fmt.Sprintf("%t", s.IsActive),
			s.CreatedAt.Format("2006-01-02T15:04:05.999999Z07:00"),
			ownerID,
		}); err != nil {
			return err
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}

	return nil
}