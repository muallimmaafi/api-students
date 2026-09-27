package helper

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	// Custom validation: field tidak boleh kosong atau hanya berisi spasi.
	_ = validate.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})
}

// messageFor menerjemahkan hasil validasi menjadi pesan yang mudah dipahami.
func messageFor(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "wajib diisi"
	case "notblank":
		return "tidak boleh kosong"
	case "min":
		return "tidak boleh kosong"
	case "gte", "lte":
		return "harus di antara 0 dan 100"
	default:
		return "format tidak valid"
	}
}

// ValidateStruct menjalankan validasi berdasarkan tag validate pada struct.
func ValidateStruct(v any) map[string]string {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{
			"request": "format tidak valid",
		}
	}

	errs := map[string]string{}

	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	for _, fieldErr := range validationErrs {
		fieldName := fieldErr.Field()

		if field, ok := t.FieldByName(fieldName); ok {
			jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
			if jsonName != "" && jsonName != "-" {
				fieldName = jsonName
			}
		}

		errs[fieldName] = messageFor(fieldErr)
	}

	return errs
}
