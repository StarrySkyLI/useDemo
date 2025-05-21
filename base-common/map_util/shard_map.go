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
