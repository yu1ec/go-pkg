package gocache

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/yu1ec/go-pkg/cachex/driver"
)

func init() {
	driver.Register("memory_gocache", New)
}

type GoCacheDriver struct {
	cache *cache.Cache
}

type GoCacheConfig struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

func New(config any) (driver.Driver, error) {
	cfg, ok := config.(*GoCacheConfig)
	if !ok {
		cfg = &GoCacheConfig{
			DefaultExpiration: 5 * time.Minute,
			CleanupInterval:   10 * time.Minute,
		}
	}

	c := cache.New(cfg.DefaultExpiration, cfg.CleanupInterval)

	return &GoCacheDriver{cache: c}, nil
}

// Add 仅当给定键的项目尚不存在或现有项目已过期时，才将项目添加到缓存。否则返回错误。
func (g *GoCacheDriver) Add(k string, v any, d time.Duration) error {
	return g.cache.Add(k, v, d)
}

func (g *GoCacheDriver) IncrementInt(k string, n int) (int, error) {
	return g.cache.IncrementInt(k, n)
}

func (g *GoCacheDriver) DecrementInt(k string, n int) (int, error) {
	return g.cache.DecrementInt(k, n)
}

func (g *GoCacheDriver) IncrementInt64(k string, n int64) (int64, error) {
	return g.cache.IncrementInt64(k, n)
}

func (g *GoCacheDriver) DecrementInt64(k string, n int64) (int64, error) {
	return g.cache.DecrementInt64(k, n)
}

func (g *GoCacheDriver) IncrementUint(k string, n uint) (uint, error) {
	return g.cache.IncrementUint(k, n)
}
func (g *GoCacheDriver) DecrementUint(k string, n uint) (uint, error) {
	return g.cache.DecrementUint(k, n)
}

func (g *GoCacheDriver) IncrementUint64(k string, n uint64) (uint64, error) {
	return g.cache.IncrementUint64(k, n)
}

func (g *GoCacheDriver) DecrementUint64(k string, n uint64) (uint64, error) {
	return g.cache.DecrementUint64(k, n)
}

// Delete 从缓存中删除一个项目。如果密钥不在缓存中，则不执行任何操作。
func (g *GoCacheDriver) Delete(k string) {
	g.cache.Delete(k)
}

// DeleteExpired 删除过期的缓存
func (g *GoCacheDriver) DeleteExpired() {
	g.cache.DeleteExpired()
}

// Flush 清空缓存
func (g *GoCacheDriver) Flush() {
	g.cache.Flush()
}

// Get 从缓存中获取一个项目。返回该项或 nil，以及一个指示是否找到该键的布尔值。
func (g *GoCacheDriver) Get(k string) (any, bool) {
	return g.cache.Get(k)
}

// GetWithExpiration 从缓存中返回一个项目及其过期时间。它返回该项目或 nil、过期时间（如果已设置）
// (如果该项目永不过期，则返回时间的零值。Time 返回) 以及指示是否找到该键的 bool。
func (g *GoCacheDriver) GetWithExpiration(k string) (any, time.Time, bool) {
	return g.cache.GetWithExpiration(k)
}

// Replace 替换缓存,如果缓存不存在,则返回错误
func (g *GoCacheDriver) Replace(k string, x any, d time.Duration) error {
	return g.cache.Replace(k, x, d)
}

// Set 添加/替换现有的缓存设置,包括过期时间,如果过期时间是0,则使用默认过期时间,如果为-1则表示永不过期
func (g *GoCacheDriver) Set(k string, x any, d time.Duration) {
	g.cache.Set(k, x, d)
}

// SetDefault 添加/替换现有的缓存设置,使用默认过期时间
func (g *GoCacheDriver) SetDefault(k string, x any) {
	g.cache.SetDefault(k, x)
}
