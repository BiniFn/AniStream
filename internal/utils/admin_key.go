package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

const (
	AdminKeyFile   = "admin.key"
	AdminKeyLength = 32
)

func GenerateAdminKey() (string, error) {
	bytes := make([]byte, AdminKeyLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func SaveAdminKey(key string) error {
	if err := os.MkdirAll("tmp", 0755); err != nil {
		return fmt.Errorf("failed to create tmp directory: %w", err)
	}

	filePath := filepath.Join("tmp", AdminKeyFile)
	if err := os.WriteFile(filePath, []byte(key), 0600); err != nil {
		return fmt.Errorf("failed to write admin key file: %w", err)
	}

	return nil
}

func LoadAdminKey() (string, error) {
	filePath := filepath.Join("tmp", AdminKeyFile)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read admin key file: %w", err)
	}
	return string(data), nil
}

func ValidateAdminKey(providedKey string) bool {
	storedKey, err := LoadAdminKey()
	if err != nil {
		return false
	}
	return providedKey == storedKey
}

