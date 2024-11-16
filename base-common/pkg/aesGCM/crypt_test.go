package aesGCM

import (
	"bytes"
	"crypto/rand"
	"strings"
	"testing"
)

func TestEncrypt(t *testing.T) {
	var (
		key = []byte("key111kjjfkdsflahdf9829uihfu32hu")
		str = `{
	"name":"test001",
	"mobile":"11111111510",
	"password":"123456",
	"verification_code":"808213"
	}`
		data = []byte(str)
	)

	enStr, err := Encrypt(key, data)
	if err != nil {
		t.Error("Encrypt err:", err.Error())
		return
	}

	dStr, err := Decrypt(key, enStr)
	if err != nil {
		t.Error("Decrypt err:", err.Error())
		return
	}

	if 0 != bytes.Compare(data, dStr) {
		t.Error("Decrypt not equal data. ")
		return
	}
}

// BenchmarkEncrypt-4        295299              3632 ns/op
func BenchmarkEncrypt(b *testing.B) {
	key := []byte("examplekey1234567890123456111111") // 32 字节密钥（AES-256）
	plaintext := make([]byte, 1024)                   // 1KB 测试数据
	if _, err := rand.Read(plaintext); err != nil {
		b.Fatalf("Failed to generate random plaintext: %v", err)
	}

	for i := 0; i < b.N; i++ {
		_, err := Encrypt(key, plaintext)
		if err != nil {
			b.Fatalf("Encrypt failed: %v", err)
		}
	}
}

// BenchmarkDecrypt-4        445113              2606 ns/op
func BenchmarkDecrypt(b *testing.B) {
	key := []byte("examplekey1234567890123456111111") // 32 字节密钥（AES-256）
	plaintext := make([]byte, 1024)                   // 1KB 测试数据
	if _, err := rand.Read(plaintext); err != nil {
		b.Fatalf("Failed to generate random plaintext: %v", err)
	}

	encrypted, err := Encrypt(key, plaintext)
	if err != nil {
		b.Fatalf("Encrypt failed: %v", err)
	}

	for i := 0; i < b.N; i++ {
		_, err := Decrypt(key, encrypted)
		if err != nil {
			b.Fatalf("Decrypt failed: %v", err)
		}
	}
}

func TestEncryptSame(t *testing.T) {
	var (
		key = []byte("key111kjjfkdsflahdf9829uihfu32hu")
		str = `{
	"name":"test001",
	"mobile":"11111111510",
	"password":"123456",
	"verification_code":"808213"
	}`
		data = []byte(str)
	)

	enStr1, err := EncryptSame(key, data)
	if err != nil {
		t.Error("Encrypt err:", err.Error())
		return
	}
	enStr2, err := EncryptSame(key, data)
	if err != nil {
		t.Error("Encrypt err:", err.Error())
		return
	}
	if 0 != strings.Compare(enStr1, enStr2) {
		t.Error("Encrypt Same err:")
		return
	}

	dStr, err := Decrypt(key, enStr1)
	if err != nil {
		t.Error("Decrypt err:", err.Error())
		return
	}

	if 0 != bytes.Compare(data, dStr) {
		t.Error("Decrypt not equal data. ")
		return
	}
}
