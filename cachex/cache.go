package cachex

type Cache interface {
	// Get 从缓存中获取一个项目。返回该项或 nil，以及一个指示是否找到该键的布尔值。
	Get(k string) (any, bool)
	// Put 添加/替换现有的缓存设置,包括过期时间,如果过期时间是0,则使用默认过期时间,如果为-1则表示永不过期 单位/秒
	Put(k string, value any, expireSeconds int64)
	// Exists 检查给定的键是否存在于缓存中。
	Exists(k string) bool
	// Remember 如果缓存中不存在该键，则从 create 函数创建一个新值，并将其添加到缓存中。单位/秒
	Remember(k string, expireSeconds int64, create func() (any, error)) (any, error)
	// RememberForever 如果缓存中不存在该键，则从 create 函数创建一个新值，并将其添加到缓存中。
	RememberForever(key string, create func() (any, error)) (any, error)
	// Forget 删除给定的键。
	Forget(key string)
	// Flush 清空缓存
	Flush()
}
