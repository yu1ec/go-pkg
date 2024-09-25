package driver

import (
	"fmt"
	"sync"
	"time"
)

// BaseDriver 是所有缓存驱动程序的基本接口
type BaseDriver interface {
	// Add 仅当给定键的项目尚不存在或现有项目已过期时，才将项目添加到缓存。否则返回错误。
	Add(k string, v any, d time.Duration) error

	// Delete 从缓存中删除一个项目。如果密钥不在缓存中，则不执行任何操作。
	Delete(k string)

	// DeleteExpired 删除过期的缓存
	DeleteExpired()

	// Flush 清空缓存
	Flush()

	// Get 从缓存中获取一个项目。返回该项或 nil，以及一个指示是否找到该键的布尔值。
	Get(k string) (any, bool)

	// GetWithExpiration 从缓存中返回一个项目及其过期时间。它返回该项目或 nil、过期时间（如果已设置）
	// (如果该项目永不过期，则返回时间的零值。Time 返回) 以及指示是否找到该键的 bool。
	GetWithExpiration(k string) (any, time.Time, bool)

	// Replace 替换缓存,如果缓存不存在,则返回错误
	Replace(k string, x any, d time.Duration) error

	// Set 添加/替换现有的缓存设置,包括过期时间,如果过期时间是0,则使用默认过期时间,如果为-1则表示永不过期
	Set(k string, x any, d time.Duration)

	// SetDefault 添加/替换现有的缓存设置,使用默认过期时间
	SetDefault(k string, x any)
}

// NumericOperations 是所有数值类型驱动程序的基本接口
type NumericOperations interface {
	IncrementInt(k string, n int) (int, error)
	DecrementInt(k string, n int) (int, error)

	IncrementInt64(k string, n int64) (int64, error)
	DecrementInt64(k string, n int64) (int64, error)

	IncrementUint(k string, n uint) (uint, error)
	DecrementUint(k string, n uint) (uint, error)

	IncrementUint64(k string, n uint64) (uint64, error)
	DecrementUint64(k string, n uint64) (uint64, error)
}

type Driver interface {
	BaseDriver
	NumericOperations
}

type DriverFactory func(config any) (Driver, error)

var (
	driverMu sync.RWMutex
	drivers  = make(map[string]DriverFactory)
)

func Register(name string, factory DriverFactory) {
	driverMu.Lock()
	defer driverMu.Unlock()
	if factory == nil {
		panic("缓存: 注册的驱动是 nil")
	}
	if _, dup := drivers[name]; dup {
		panic("缓存：注册重复的驱动名称 " + name)
	}
	drivers[name] = factory
}

// New 创建一个新的驱动实例
func New(driverName string, config any) (Driver, error) {
	driverMu.RLock()
	driverFactory, ok := drivers[driverName]
	driverMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("缓存：未注册的驱动名称 %s", driverName)
	}
	return driverFactory(config)
}
