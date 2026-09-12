package service

import (
	"testing"

	"api-students/app/model"
)

func TestCheckPasswordStrength(t *testing.T) {
	cases := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"terlalu pendek", "abc123", true},
		{"tanpa angka", "abcdefgh", true},
		{"tanpa huruf", "12345678", true},
		{"password umum", "password123", true},
		{"password kuat", "rahasia123", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := checkPasswordStrength(tc.password)
			gotErr := msg != ""
			if gotErr != tc.wantErr {
				t.Errorf("password=%q: harap error=%v, dapat pesan=%q",
					tc.password, tc.wantErr, msg)
			}
		})
	}
}

func TestValidateRegister(t *testing.T) {
	// Kasus semua field salah
	errs := ValidateRegister(model.RegisterRequest{
		Username: "ab", Email: "bukan-email", Password: "lemah",
	})
	if len(errs) != 3 {
		t.Errorf("harap 3 error, dapat %d: %v", len(errs), errs)
	}

	// Kasus semua field valid
	errsValid := ValidateRegister(model.RegisterRequest{
		Username: "budi_santoso", Email: "budi@example.com", Password: "rahasia123",
	})
	if len(errsValid) != 0 {
		t.Errorf("tidak seharusnya ada error: %v", errsValid)
	}
}

func TestIsValidUsername(t *testing.T) {
	valid := []string{"budi", "budi123", "budi.santoso", "budi_santoso"}
	for _, u := range valid {
		if !isValidUsername(u) {
			t.Errorf("%q seharusnya valid", u)
		}
	}

	invalid := []string{"budi santoso", "budi@santoso", "budi!"}
	for _, u := range invalid {
		if isValidUsername(u) {
			t.Errorf("%q seharusnya TIDAK valid", u)
		}
	}
}