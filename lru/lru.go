package lru

import "container/list"

// Cache 是一个LRU 缓存。并发不安全。
type Cache struct {
	maxBytes int64                    // 允许使用的最大内存
	nbytes   int64                    // 当前已使用的内存
	ll       *list.List               // 双向链表
	cache    map[string]*list.Element // 键是字符串，值是双向链表中对应节点的指针
	// 可选，在某条记录被移除时的回调函数
	OnEvicted func(key string, value Value)
}

// entry 是双向链表节点的数据类型，在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射。
type entry struct {
	key   string
	value Value
}

// Value 使用 Len 来返回其在内存中的大小
type Value interface {
	Len() int
}

// New 创建一个新的 Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get 查找一个 key
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 移除最久未使用的记录
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() // 取到队首节点，从链表中删除。
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// 从字典中 c.cache 删除该节点的映射关系。
		delete(c.cache, kv.key)
		// 更新当前所用的内存 c.nbytes。
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())

		// 如果回调函数 OnEvicted 不为 nil，则调用回调函数。
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 向缓存添加一个值
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { // 如果键存在，则更新对应节点的值，并将该节点移到队尾。
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 更新值
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { // 不存在则是新增场景，首先队尾添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关系。
		// 添加新元素
		ele = c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// 更新 c.nbytes，如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点。
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len 返回当前缓存的元素个数
func (c *Cache) Len() int {
	return c.ll.Len()
}
