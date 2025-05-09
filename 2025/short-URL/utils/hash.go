package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

// GenerateShortURL 生成短链URL
func GenerateShortURL(longURL string) string {
	hash := sha256.Sum256([]byte(longURL))
	return base64.URLEncoding.EncodeToString(hash[:])[:8]
}
