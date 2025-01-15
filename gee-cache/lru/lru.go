package lru

import (
	"container/list"
)

// Cache是一个LRU缓存，并发访问是不安全的
type Cache struct {
	maxBytes  int64      // 允许使用的最大内存
	nbytes    int64      // 当前已使用的内存
	ll        *list.List // 双向链表
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value) // 可选，在清除条目时执行
}

type entry struct {
	key   string
	value Value
}

// 值使用Len来计算它需要多少字节
type Value interface {
	Len() int
}

// New是Cache的构造函数
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 增加
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { // 如果键存在，则更新对应节点的值，并将该节点移到队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { // 不存在则是新增场景，首先队尾添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// 如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// 查询
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 如果键对应的链表节点存在，则将对应节点移动到对位，并返回查找到的值
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 删除
func (c *Cache) RemoveOldest() {
	// 实际上是缓存淘汰，即移除最近最少访问的节点（队首）
	ele := c.ll.Back() // 取到队首节点，从链表中删除
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                // 从字典中c.cache删除该节点的映射关系
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) // 更新当前所用内存
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
