package helpers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
)

func Sha256String(text string) string {
	hash := sha256.New()
	if _, err := io.Copy(hash, strings.NewReader(text)); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func LowerTrimString(text string) string {
	return strings.ToLower(strings.TrimSpace(text))
}
