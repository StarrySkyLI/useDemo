package configServer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logc"
	"gitlab.coolgame.world/go-template/base-common/pkg/cache"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var once sync.Once

type Client struct {
	ctx         context.Context
	conn        *clientv3.Client
	cacheClient *cache.Client

	conf Config
}

type Config struct {
	Env     string `json:"env"`
	AppName string `json:"appName"`

	// etcd
	Endpoints   []string      `json:"endpoints"`
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	DialTimeout time.Duration `json:"dial-timeout"`
}

func NewClient(ctx context.Context, conf Config) *Client {
	if conf.DialTimeout == 0 {
		conf.DialTimeout = 5 * time.Second
	}

	return &Client{
		ctx:  ctx,
		conf: conf,
	}
}

func (c *Client) MustStart() {
	// 配置 etcd 客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.conf.Endpoints, // etcd 服务器地址
		DialTimeout: c.conf.DialTimeout,
		Username:    c.conf.Username,
		Password:    c.conf.Password,
	})
	if err != nil {
		logc.Must(err)
	}

	c.conn = cli
	c.cacheClient = cache.NewClient(c.conf.Env, c.conf.AppName)

	var initCan = make(chan struct{}, 1)
	defer close(initCan)
	once.Do(func() {
		if err = c.initEtcdKey(initCan); err != nil {
			logc.Must(err)
		}
	})

	// watch
	go c.WatchEtcd(initCan)
}

func (c *Client) Get(key string) (interface{}, bool) {
	return c.cacheClient.Get(c.cacheClient.BuildKey(key))
}

// 初始化 key 到本地内存
func (c *Client) initEtcdKey(initChan chan struct{}) error {
	var key = c.cacheClient.BuildKey("")
	res, err := c.conn.Get(c.ctx, key, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	c.batchSetKey(res.Kvs)
	initChan <- struct{}{}

	return nil
}

func (c *Client) batchSetKey(data []*mvccpb.KeyValue) {
	if len(data) > 100 {
		var waite sync.WaitGroup
		for i := 0; i < len(data); i += 100 {
			waite.Add(1)
			go func(start, end int) {
				defer waite.Done()
				for j := start; j < end; j++ {
					c.cacheClient.Set(string(data[j].Key), string(data[j].Value))
				}
			}(i, i+100)
		}

		waite.Wait()
	} else {
		for i := 0; i < len(data); i++ {
			c.cacheClient.Set(string(data[i].Key), string(data[i].Value))
		}
	}
}

func (c *Client) WatchEtcd(initCan chan struct{}) {
	// 创建一个 Watcher
	rch := c.conn.Watch(c.ctx, c.cacheClient.BuildKey(""), clientv3.WithPrefix())

	<-initCan

	for {
		select {
		case <-c.ctx.Done():
			_ = c.conn.Close()
			fmt.Println("Config of etcd watch down. ")
			return
		case wresp := <-rch:
			for _, ev := range wresp.Events {
				switch ev.Type {
				case mvccpb.PUT:
					c.cacheClient.Set(string(ev.Kv.Key), string(ev.Kv.Value))
				case mvccpb.DELETE:
					c.cacheClient.Delete(string(ev.Kv.Key))
				}
			}
		}
	}
}
