package cache

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	client := NewClient("dev", "base")
	client.Set("hallo", "world")

	if data, ok := client.Get("hallo"); !ok {
		t.Error("Failed to retrieve data")
	} else {
		if str, ok := data.(string); ok {
			if str != "world" {
				t.Error("Retrieved data does not match")
			}
		} else {
			t.Error("Retrieved data is not a string")
		}
	}
}

var client = NewClient("dev", "baseB")

// BenchmarkSet-4           1000000              1151 ns/op
func BenchmarkSet(b *testing.B) {
	plaintext := make([]byte, 1024) // 1KB 测试数据
	if _, err := rand.Read(plaintext); err != nil {
		b.Fatalf("Failed to generate random plaintext: %v", err)
	}

	for i := 0; i < b.N; i++ {
		client.Set(fmt.Sprintf("k%d", i), plaintext)
	}
}

// BenchmarkGet-4           2347737               535.4 ns/op  = 1,867,985
func BenchmarkGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, ok := client.Get(fmt.Sprintf("k%d", i))
		if !ok {
			// fmt.Println(fmt.Sprintf("fail get key:[k%d]", i))
		}
	}
}
