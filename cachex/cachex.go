package cachex

import (
	"time"

	"github.com/yu1ec/go-pkg/cachex/driver"
)

func New(driverName string, config any) (Cache, error) {
	d, err := driver.New(driverName, config)
	if err != nil {
		return nil, err
	}

	return &cacheImpl{
		driver: d,
	}, nil

}

type cacheImpl struct {
	driver driver.Driver
}

func (c *cacheImpl) Get(k string) (any, bool) {
	return c.driver.Get(k)
}

func (c *cacheImpl) Put(k string, v any, expireSeconds int64) {
	d := time.Duration(expireSeconds) * time.Second
	c.driver.Set(k, v, d)
}

func (c *cacheImpl) Exists(k string) bool {
	_, exists := c.driver.Get(k)
	return exists
}

func (c *cacheImpl) Remember(k string, expireSeconds int64, create func() (any, error)) (any, error) {
	v, exists := c.Get(k)
	if exists {
		return v, nil
	}

	v, err := create()
	c.Put(k, v, expireSeconds)

	return v, err
}

func (c *cacheImpl) RememberForever(k string, create func() (any, error)) (any, error) {
	return c.Remember(k, -1, create)
}

func (c *cacheImpl) Forget(k string) {
	c.driver.Delete(k)
}

func (c *cacheImpl) Flush() {
	c.driver.Flush()
}
