package uuid

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/satori/go.uuid"
	"math/rand"
	"time"
)

func GenUUID() uuid.UUID {
	v1 := uuid.NewV1()
	return v1
}

// 随机字符串
func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	randomString := make([]byte, length)
	for i := 0; i < length; i++ {
		randomString[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomString)
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
