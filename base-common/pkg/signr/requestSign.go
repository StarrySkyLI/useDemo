package signr

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// GenerateNonce 生成 Nonce
func GenerateNonce() string {
	randomString := generateRandomString(64)
	timestamp := getCurrentTimestamp()
	return randomString + "." + timestamp
}

// 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// 获取当前毫秒级时间戳
func getCurrentTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
}

// 生成 MD5 签名
func generateMD5Hash(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateHeadSign 生成 HeadSign
func GenerateHeadSign(signKey, token, nonce, version string) string {
	concatenatedString := signKey + token + nonce + version
	return strings.ToLower(generateMD5Hash(concatenatedString))
}
