package aesGCM

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

var IsOpenAesGcm = false
var EncryptKey = make([]byte, 0)

// AES-GCM 加密
func Encrypt(key, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 使用 GCM 模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	res := base64.URLEncoding.EncodeToString(ciphertext)
	res = strings.ReplaceAll(res, "=", "")

	return res, nil
}

func EncryptSame(key, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 使用 GCM 模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	for i := 0; i < aesGCM.NonceSize(); i++ {
		nonce[i] = 1
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	res := base64.URLEncoding.EncodeToString(ciphertext)
	res = strings.ReplaceAll(res, "=", "")

	return res, nil
}

// AES-GCM 解密
func Decrypt(key []byte, encryptedText string) ([]byte, error) {
	if l := len(encryptedText) % 4; l != 0 {
		encryptedText += strings.Repeat("=", 4-l)
	}

	ciphertext, err := base64.URLEncoding.DecodeString(encryptedText)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
