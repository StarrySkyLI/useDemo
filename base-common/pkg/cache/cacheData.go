package cache

import (
	"fmt"

	"github.com/patrickmn/go-cache"
)

type Client struct {
	AppName string
	Prefix  string

	cache *cache.Cache
}

func NewClient(env, appName string) *Client {
	return &Client{
		AppName: appName,
		Prefix:  env,
		cache:   cache.New(0, 0),
	}
}

func (c *Client) BuildKey(key string) string {
	return fmt.Sprintf("/%s/%s/%s", c.Prefix, c.AppName, key)
}

func (c *Client) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *Client) Set(key string, value interface{}) {
	c.cache.SetDefault(key, value)
}

func (c *Client) Delete(key string) {
	c.cache.Delete(key)
}
