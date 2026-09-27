package helper

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"api-students/app/model"
)

// EncodeCursor mengubah posisi cursor menjadi string yang aman
// untuk dikirim melalui query parameter.
func EncodeCursor(cursor model.Cursor) string {
	raw := cursor.CreatedAt.UTC().Format(time.RFC3339Nano) + "|" +
		strconv.Itoa(cursor.ID)

	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

// DecodeCursor membaca cursor dari query parameter.
func DecodeCursor(value string) (model.Cursor, error) {
	var cursor model.Cursor

	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return cursor, fmt.Errorf("cursor tidak valid")
	}

	parts := strings.Split(string(decoded), "|")
	if len(parts) != 2 {
		return cursor, fmt.Errorf("cursor tidak valid")
	}

	createdAt, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return cursor, fmt.Errorf("cursor tidak valid")
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil || id < 1 {
		return cursor, fmt.Errorf("cursor tidak valid")
	}

	cursor.CreatedAt = createdAt
	cursor.ID = id

	return cursor, nil
}
