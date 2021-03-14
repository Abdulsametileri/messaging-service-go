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

func HashPassword(text string) string {
	return Sha256String(text)
}

func StripBearer(tok string) (string, error) {
	if len(tok) > 6 && strings.ToLower(tok[0:7]) == "bearer " {
		return tok[7:], nil
	}
	return tok, nil
}
