package utils

import (
	"math/rand"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateId (leng int) string {
	var sb strings.Builder
	sb.Grow(leng)
	for range leng {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}