package geecache

import (
	"geecache/lru"
	"sync"
)

// 实例化 lru，封装 get 和 add 方法，并添加互斥锁 mu。
type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

// add 添加缓存
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//在 add 方法中，判断了 c.lru 是否为 nil，如果等于 nil 再创建实例。
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

// get 获取缓存value值
func (c *cache) get(key string) (value ByteView, ok bool) {
	//获取的时候上锁?
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
