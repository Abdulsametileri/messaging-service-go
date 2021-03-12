package helpers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
)

func Sha256String(metin string) string {
	hash := sha256.New()
	if _, err := io.Copy(hash, strings.NewReader(metin)); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}
