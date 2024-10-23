package redis_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/yu1ec/go-pkg/cachex/driver"
	"github.com/yu1ec/go-pkg/cachex/driver/redis"
)

func setupRedis(t *testing.T) (*miniredis.Miniredis, driver.Driver) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	config := &redis.RedisConfig{
		Addr: mr.Addr(),
	}

	d, err := redis.New(config)
	if err != nil {
		t.Fatalf("failed to create Redis driver: %v", err)
	}

	return mr, d
}

func TestRedisDriver(t *testing.T) {
	mr, d := setupRedis(t)
	defer mr.Close()

	t.Run("Set and Get", func(t *testing.T) {
		d.Set("key1", "value1", 1*time.Minute)
		val, exists := d.Get("key1")
		assert.True(t, exists)
		assert.Equal(t, "value1", val)
	})

	t.Run("Add", func(t *testing.T) {
		err := d.Add("key2", "value2", 1*time.Minute)
		assert.NoError(t, err)
		err = d.Add("key2", "value3", 1*time.Minute)
		assert.Error(t, err)
	})

	t.Run("Replace", func(t *testing.T) {
		d.Set("key3", "value3", 1*time.Minute)
		err := d.Replace("key3", "new_value3", 1*time.Minute)
		assert.NoError(t, err)
		val, exists := d.Get("key3")
		assert.True(t, exists)
		assert.Equal(t, "new_value3", val)
	})

	t.Run("Delete", func(t *testing.T) {
		d.Set("key4", "value4", 1*time.Minute)
		d.Delete("key4")
		_, exists := d.Get("key4")
		assert.False(t, exists)
	})

	t.Run("Flush", func(t *testing.T) {
		d.Set("key5", "value5", 1*time.Minute)
		d.Flush()
		_, exists := d.Get("key5")
		assert.False(t, exists)
	})

	t.Run("GetWithExpiration", func(t *testing.T) {
		d.Set("key6", "value6", 1*time.Minute)
		val, expiration, exists := d.GetWithExpiration("key6")
		assert.True(t, exists)
		assert.Equal(t, "value6", val)
		assert.True(t, expiration.After(time.Now()))
	})

	t.Run("Increment and Decrement", func(t *testing.T) {
		d.Set("counter", "10", 1*time.Minute)

		val, err := d.IncrementInt("counter", 5)
		assert.NoError(t, err)
		assert.Equal(t, 15, val)

		val, err = d.DecrementInt("counter", 3)
		assert.NoError(t, err)
		assert.Equal(t, 12, val)
	})
}
