package gee_cache

// 缓存值的抽象与封装

// ByteView 包括一个只读的字节切片 b，b 被包装在 ByteView 中是为了防止缓存值被外部程序修改。
type ByteView struct {
	b []byte // b 将会存储真实的缓存值。选择 byte 类型是为了能够支持任意的数据类型的存储，例如字符串、图片等
}

// Len 返回字节切片的长度
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 返回一个拷贝，防止缓存值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String 返回字符串
func (v ByteView) String() string {
	return string(v.b)
}

// cloneBytes 返回一个拷贝，防止缓存值被外部程序修改
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
