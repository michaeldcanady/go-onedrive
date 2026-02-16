package edit

import (
	"crypto/sha256"
	"io"
	"log"
	"path/filepath"
	"strings"
)

func Name(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func hashSha256(reader io.Reader) []byte {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		log.Fatal(err)
	}

	return hash.Sum(nil)
}
