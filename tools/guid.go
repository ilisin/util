package tools

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func Guid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}
