package utils

import "github.com/google/uuid"

func GenerateUUID() string {
	uuidObj := uuid.New()

	// Convert UUID to string
	uuidStr := uuidObj.String()

	return uuidStr
}
