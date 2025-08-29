package cache

import (
	"github.com/hashicorp/golang-lru/v2/expirable"
	"orderService/internal/models"
	"time"
)

//go:generate mockery --name=ILruCache --output=mocks --outpkg=mocks --case=snake --with-expecter
type ILruCache interface {
	Get(key string) (models.OrderView, bool)
	Add(key string, value models.OrderView) bool
}

type OrderLRuCache struct {
	LruCache *expirable.LRU[string, models.OrderView]
}

func NewCache(size, ttl int) OrderLRuCache {
	cache := expirable.NewLRU[string, models.OrderView](size, nil, time.Duration(ttl)*time.Second)
	return OrderLRuCache{cache}
}

func (o OrderLRuCache) Get(key string) (models.OrderView, bool) {
	return o.LruCache.Get(key)
}

func (o OrderLRuCache) Add(key string, value models.OrderView) bool {
	return o.LruCache.Add(key, value)
}
