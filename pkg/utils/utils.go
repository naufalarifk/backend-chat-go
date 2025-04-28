package utils

import (
	"github.com/google/uuid"
)

func GenerateRandomID() string {
	return uuid.New().String()
}
