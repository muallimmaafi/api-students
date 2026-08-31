package service

import (
	"testing"

	"api-students/app/model"
)

// Perhatikan: pengujian ini tidak menyalakan server, tidak menyentuh
// database, dan tidak membuat fiber.Ctx.

func TestCountTotalPages(t *testing.T) {
	cases := []struct{ total, limit, want int }{
		{0, 10, 0},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{137, 20, 7},
	}
	for _, tc := range cases {
		if got := CountTotalPages(tc.total, tc.limit); got != tc.want {
			t.Errorf("total=%d limit=%d: harap %d, dapat %d",
				tc.total, tc.limit, tc.want, got)
		}
	}
}

func TestValidateCreate(t *testing.T) {
	errs := ValidateCreate(model.CreateStudentRequest{
		NIM: "", Name: "", Grade: 150,
	})
	if len(errs) != 3 {
		t.Errorf("harap 3 error, dapat %d: %v", len(errs), errs)
	}

	errsValid := ValidateCreate(model.CreateStudentRequest{
		NIM: "12345", Name: "Budi", Grade: 85,
	})
	if len(errsValid) != 0 {
		t.Errorf("tidak seharusnya ada error: %v", errsValid)
	}
}

func TestApplyPatch(t *testing.T) {
	initial := model.Student{ID: 1, NIM: "12345", Name: "Budi", Grade: 80, IsActive: true}
	inactive := false

	result, errs := ApplyPatch(initial, model.PatchStudentRequest{IsActive: &inactive})
	if len(errs) != 0 {
		t.Fatalf("tidak seharusnya ada error: %v", errs)
	}
	if result.IsActive {
		t.Error("is_active seharusnya berubah menjadi false")
	}
	if result.Name != "Budi" {
		t.Error("field yang tidak dikirim seharusnya tidak berubah")
	}
}