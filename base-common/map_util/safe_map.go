package map_util

import "sync"

type SafeMap struct {
	sync.RWMutex // 读写锁，保护下面的map
	m            map[string]string
}

func NewSafeMap(cap int) *SafeMap {
	return &SafeMap{
		m: make(map[string]string, cap),
	}
}
func (sm *SafeMap) Set(key, value string) { // 设置键值对
	sm.Lock()
	defer sm.Unlock()
	sm.m[key] = value
}
func (sm *SafeMap) Get(key string) string { // 获取键值对
	sm.RLock()
	defer sm.RUnlock()
	return sm.m[key]
}
func (sm *SafeMap) Delete(key string) { // 删除键值对
	sm.Lock()
	defer sm.Unlock()
	delete(sm.m, key)
}
func (sm *SafeMap) Len() int { // 获取map长度
	sm.RLock()
	defer sm.RUnlock()
	return len(sm.m)
}
func (sm *SafeMap) Each(f func(key, value string) bool) { // 遍历map
	sm.RLock()
	defer sm.RUnlock()
	for k, v := range sm.m {
		if !f(k, v) {
			break
		}
	}
}
