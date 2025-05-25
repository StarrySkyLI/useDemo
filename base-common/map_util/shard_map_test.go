package map_util

import (
	"errors"
	"strconv"
	"sync"
	"testing"
)

func TestNewShardMap_Initialization(t *testing.T) {
	sm := NewShardMap()
	if len(sm.shards) != shardCount {
		t.Fatalf("Expected %d shards, got %d", shardCount, len(sm.shards))
	}

	for i, s := range sm.shards {
		if s == nil {
			t.Fatalf("Shard %d is nil", i)
		}
		if len(s.m) != 0 {
			t.Fatalf("Shard %d map should be empty", i)
		}
	}
}

func TestShardMap_BasicOperations(t *testing.T) {
	sm := NewShardMap()

	t.Run("Store and Load existing key", func(t *testing.T) {
		key := "exist"
		value := "value"
		sm.Store(key, value)

		if val, ok := sm.Load(key); !ok || val != value {
			t.Errorf("Expected (%v, true), got (%v, %v)", value, val, ok)
		}
	})

	t.Run("Load non-existing key", func(t *testing.T) {
		if val, ok := sm.Load("ghost"); ok || val != nil {
			t.Errorf("Expected (nil, false), got (%v, %v)", val, ok)
		}
	})

	t.Run("Overwrite existing key", func(t *testing.T) {
		key := "overwrite"
		sm.Store(key, "old")
		sm.Store(key, "new")

		if val, _ := sm.Load(key); val != "new" {
			t.Errorf("Expected 'new', got %v", val)
		}
	})

	t.Run("Empty key handling", func(t *testing.T) {
		sm.Store("", "empty")
		if val, ok := sm.Load(""); !ok || val != "empty" {
			t.Errorf("Empty key storage failed")
		}
	})

	t.Run("Nil value storage", func(t *testing.T) {
		sm.Store("nil", nil)
		if val, ok := sm.Load("nil"); !ok || val != nil {
			t.Errorf("Nil value storage failed")
		}
	})
}

func TestShardMap_ShardSelection(t *testing.T) {
	sm := NewShardMap()
	key := "consistent"

	// 验证相同key总是返回同一分片
	shard1 := sm.getShard(key)
	shard2 := sm.getShard(key)
	if shard1 != shard2 {
		t.Fatal("Same key returns different shards")
	}

	// 验证不同key可能分配到不同分片（概率性测试）
	diffKeys := []string{"a", "b", "c", "d", "e"}
	shards := make(map[*shard]bool)
	for _, k := range diffKeys {
		shards[sm.getShard(k)] = true
	}
	if len(shards) == 1 {
		t.Log("Warning: All test keys hashed to same shard. Consider better key selection")
	}
}

func TestShardMap_Concurrency(t *testing.T) {
	sm := NewShardMap()
	const workers = 100
	var wg sync.WaitGroup

	t.Run("Parallel writes to different keys", func(t *testing.T) {
		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func(i int) {
				defer wg.Done()
				key := string(rune(i)) // 生成不同key
				sm.Store(key, i)
				if val, ok := sm.Load(key); !ok || val != i {
					t.Errorf("Concurrent write failed for key %s", key)
				}
			}(i)
		}
		wg.Wait()
	})

	t.Run("Contention on same shard", func(t *testing.T) {
		// 所有操作使用相同分片
		key := "contention"
		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func(v int) {
				defer wg.Done()
				sm.Store(key, v)
				sm.Load(key)
			}(i)
		}
		wg.Wait()

		// 最终值应该是最后一个写入的值
		if val, _ := sm.Load(key); val != workers-1 {
			t.Log("Expected last-write-wins behavior, value:", val)
		}
	})
}

func TestFnv32_Consistency(t *testing.T) {
	tests := []struct {
		input  string
		output uint32 // 预计算值
	}{
		{"test", 0xbc2c0be9},
		{"", 0x811c9dc5},
		{"hello", 0x4f9f2cab},
	}

	for _, tt := range tests {
		if got := fnv32(tt.input); got != tt.output {
			t.Errorf("fnv32(%q) = 0x%x, want 0x%x", tt.input, got, tt.output)
		}
	}
}

func TestRange(t *testing.T) {
	sm := NewShardMap()

	// 填充测试数据
	for i := 0; i < 1000; i++ {
		sm.Store(strconv.Itoa(i), i)
	}

	t.Run("正常遍历", func(t *testing.T) {
		count := 0
		err := sm.Range(func(k string, v any) error {
			count++
			return nil
		})
		if err != nil || count != 1000 {
			t.Fatalf("遍历失败，数量不符: %d, 错误: %v", count, err)
		}
	})

	t.Run("中途终止", func(t *testing.T) {
		stopErr := errors.New("stop")
		count := 0
		err := sm.Range(func(k string, v any) error {
			count++
			if count == 500 {
				return stopErr
			}
			return nil
		})
		if err != stopErr || count != 500 {
			t.Fatal("未正确终止遍历")
		}
	})

	t.Run("并发安全", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				sm.Store("key"+strconv.Itoa(n), n)
				sm.Range(func(k string, v any) error { return nil })
			}(i)
		}
		wg.Wait()
	})
}

func TestKeysValues(t *testing.T) {
	sm := NewShardMap()

	// 空map测试
	if len(sm.Keys()) != 0 || len(sm.Values()) != 0 {
		t.Fatal("空map返回非空结果")
	}

	// 填充测试数据
	expected := make(map[string]int)
	for i := 0; i < 1000; i++ {
		key := strconv.Itoa(i)
		sm.Store(key, i)
		expected[key] = i
	}

	t.Run("Keys验证", func(t *testing.T) {
		keys := sm.Keys()
		if len(keys) != 1000 {
			t.Fatal("键数量不符")
		}
		for _, k := range keys {
			if _, ok := expected[k]; !ok {
				t.Fatal("存在无效键")
			}
		}
	})

	t.Run("Values验证", func(t *testing.T) {
		values := sm.Values()
		if len(values) != 1000 {
			t.Fatal("值数量不符")
		}
		valueMap := make(map[int]bool)
		for _, v := range values {
			valueMap[v.(int)] = true
		}
		for i := 0; i < 1000; i++ {
			if !valueMap[i] {
				t.Fatal("值缺失")
			}
		}
	})
}

func BenchmarkRange(b *testing.B) {
	sm := NewShardMap()
	for i := 0; i < 100000; i++ {
		sm.Store(strconv.Itoa(i), i)
	}

	b.Run("纯遍历", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			sm.Range(func(k string, v any) error { return nil })
		}
	})

	b.Run("并发读写", func(b *testing.B) {
		var wg sync.WaitGroup
		b.RunParallel(func(pb *testing.PB) {
			wg.Add(1)
			defer wg.Done()
			counter := 0
			for pb.Next() {
				// 混合读写操作
				key := strconv.Itoa(counter)
				sm.Store(key, counter)
				sm.Delete(key)
				counter++
				sm.Range(func(k string, v any) error { return nil })
			}
		})
		wg.Wait()
	})
}

func BenchmarkKeys(b *testing.B) {
	sm := NewShardMap()
	for i := 0; i < 100000; i++ {
		sm.Store(strconv.Itoa(i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sm.Keys()
	}
}

func BenchmarkValues(b *testing.B) {
	sm := NewShardMap()
	for i := 0; i < 100000; i++ {
		sm.Store(strconv.Itoa(i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sm.Values()
	}
}
