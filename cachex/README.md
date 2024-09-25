# Cachex

缓存处理库，支持多种缓存实现，包括内存缓存，Redis 缓存等

## 使用

```go
import (
	"github.com/yu1ec/go-pkg/cachex"
	_ "github.com/yu1ec/go-pkg/cachex/driver/memory"  // 导入 memory 驱动
	// _ "github.com/yu1ec/go-pkg/cachex/driver/redis"  // 导入 redis 驱动（需要实现）
	// _ "github.com/yu1ec/go-pkg/cachex/driver/mysql"  // 导入 mysql 驱动（需要实现）
)

func main() {
	// 使用 memory 驱动的 gocache 实现
	memoryCache, err := cachex.New("memory", map[string]any{
		"implementation": "gocache",
		"CleanupInterval": 10 * time.Minute,
	})
	if err != nil {
		// 处理错误
	}

	// 使用 memoryCache...

	// 使用 redis 驱动（假设已实现）
	redisCache, err := cachex.New("redis", map[string]any{
		"address": "localhost:6379",
		"password": "",
		"db": 0,
	})
	if err != nil {
		// 处理错误
	}

	// 使用 redisCache...
}
```