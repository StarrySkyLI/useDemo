package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

func NewCache() {
	c := cache.New(3600*time.Second, 10*time.Second)
	c.Set("key", "value", cache.DefaultExpiration)
}
