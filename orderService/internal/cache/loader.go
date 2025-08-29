package cache

import (
	"orderService/internal/repository"
)

type LCacheLoader struct {
	repo  repository.Repository
	cache OrderLRuCache
}

func NewLCacheLoader(r repository.Repository, c OrderLRuCache) LCacheLoader {
	return LCacheLoader{
		repo:  r,
		cache: c,
	}
}

func (l *LCacheLoader) LoadCache(c OrderLRuCache, limit int) error {
	recentOrders, err := l.repo.GetRecentOrders(limit)
	if err != nil {
		return err
	}

	for _, order := range recentOrders {
		c.LruCache.Add(order.Uid.String(), order.ToOrderView())
	}

	return nil
}
