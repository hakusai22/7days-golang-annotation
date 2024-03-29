package lru

import "container/list"

// Cache 并发不安全的缓存
type Cache struct {
	maxBytes  int64                         //允许使用的最大内存
	nBytes    int64                         //当前已使用的内存
	ll        *list.List                    //Go 语言标准库实现的双向链表list.List
	cache     map[string]*list.Element      //键是字符串，值是双向链表中对应节点的指针。
	OnEvicted func(key string, value Value) //某条记录被移除时的回调函数，可以为 nil
}

// 键值对 entry 是双向链表节点的数据类型，在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射。
type entry struct {
	key   string
	value Value
}

// Value 为了通用性，我们允许值是实现了 Value 接口的任意类型，该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小。
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Add adds a value to the cache. 缓存的方法
/**
1. 如果键存在，则更新对应节点的值，并将该节点移到队尾。
2. 不存在则是新增场景，首先队尾添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关系。
3. 更新 c.nbytes，如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点。
*/
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// Get 查找主要有 2 个步骤，第一步是从字典中找到对应的双向链表的节点，第二步，将该节点移动到队尾。
/**
1. 如果键对应的链表节点存在，则将对应节点移动到队尾，并返回查找到的值。
2. c.ll.MoveToFront(ele)，即将链表中的节点 ele 移动到队尾（双向链表作为队列，队首队尾是相对的，在这里约定 front 为队尾）
*/
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 这里的删除，实际上是缓存淘汰。即移除最近最少访问的节点（队首）
/**
1. c.ll.Back() 取到队首节点，从链表中删除。
2. delete(c.cache, kv.key)，从字典中 c.cache 删除该节点的映射关系。
3. 更新当前所用的内存 c.nbytes。
4. 如果回调函数 OnEvicted 不为 nil，则调用回调函数。
*/
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
