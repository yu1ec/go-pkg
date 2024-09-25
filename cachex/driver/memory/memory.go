package memory

import (
	"fmt"

	"github.com/yu1ec/go-pkg/cachex/driver"
	"github.com/yu1ec/go-pkg/cachex/driver/memory/gocache"
)

func init() {
	driver.Register("memory", NewMemoryCache)
}

// NewMemoryCache 创建一个内存缓存
func NewMemoryCache(config any) (driver.Driver, error) {
	cfg, ok := config.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid config for memory cache")
	}

	implementation, ok := cfg["implementation"].(string)
	if !ok {
		implementation = "gocache" // 默认使用 gocache
	}
	switch implementation {
	case "gocache":
		return gocache.New(config)
	default:
		return nil, fmt.Errorf("unsupported implementation: %s", implementation)
	}
}
