package safemap

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
	m map[string]string
}

// 分段锁
func NewShardMap() *ShardMap {
	sm := &ShardMap{}
	for i := 0; i < shardCount; i++ {
		sm.shards[i] = &shard{m: make(map[string]string)}
	}
	return sm
}

func (sm *ShardMap) getShard(key string) *shard {
	hash := fnv32(key)
	return sm.shards[hash%shardCount]
}

func (sm *ShardMap) Store(key string, value string) {
	s := sm.getShard(key)
	s.Lock()
	defer s.Unlock()
	s.m[key] = value
}

func (sm *ShardMap) Load(key string) (string, bool) {
	s := sm.getShard(key)
	s.RLock()
	defer s.RUnlock()
	val, ok := s.m[key]
	return val, ok
}

func fnv32(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}
