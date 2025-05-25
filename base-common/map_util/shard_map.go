package map_util

import (
	"hash/fnv"
	"sync"
)

const shardCount = 32

type ShardMap struct {
	shards [shardCount]*shard
}

type shard struct {
	sync.RWMutex
	m map[string]any
}

// 分段锁
func NewShardMap() *ShardMap {
	sm := &ShardMap{}
	for i := 0; i < shardCount; i++ {
		sm.shards[i] = &shard{m: make(map[string]any)}
	}
	return sm
}

func (sm *ShardMap) getShard(key string) *shard {
	hash := fnv32(key)
	return sm.shards[hash%shardCount]
}

func (sm *ShardMap) Store(key string, value any) {
	s := sm.getShard(key)
	s.Lock()
	defer s.Unlock()
	s.m[key] = value
}

func (sm *ShardMap) Load(key string) (any, bool) {
	s := sm.getShard(key)
	s.RLock()
	defer s.RUnlock()
	val, ok := s.m[key]
	return val, ok
}
func (sm *ShardMap) Delete(key string) bool {
	s := sm.getShard(key)
	s.Lock()
	defer s.Unlock()

	if _, ok := s.m[key]; ok {
		delete(s.m, key)
		return true
	}
	return false
}
func fnv32(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

// 安全遍历所有键值对（仿 sync.Map 的 Range 设计）
func (sm *ShardMap) Range(f func(key string, value any) error) error {
	for _, shard := range sm.shards {
		shard.RLock()
		for k, v := range shard.m {
			if err := f(k, v); err != nil {
				shard.RUnlock()
				return err // 遇到错误提前终止
			}
		}
		shard.RUnlock()
	}
	return nil
}

// 获取所有键（线程安全）
func (sm *ShardMap) Keys() []string {
	keys := make([]string, 0)
	for _, shard := range sm.shards {
		shard.RLock()
		for k := range shard.m {
			keys = append(keys, k)
		}
		shard.RUnlock()
	}
	return keys
}

// 获取所有值（线程安全）
func (sm *ShardMap) Values() []any {
	values := make([]any, 0)
	for _, shard := range sm.shards {
		shard.RLock()
		for _, v := range shard.m {
			values = append(values, v)
		}
		shard.RUnlock()
	}
	return values
}
