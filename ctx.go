package gootelinstrument

import (
	"context"
	"sync"

	"github.com/pijng/gls"
)

var (
	cacheOnce sync.Once
	cache     ctxCache
)

type ctxCache struct {
	mu      sync.RWMutex
	entries map[int64]context.Context
}

func (c *ctxCache) get(id int64) (context.Context, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ctx, ok := c.entries[id]
	return ctx, ok
}

func (c *ctxCache) set(id int64, ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[id] = ctx
}

func init() {
	cacheOnce.Do(func() {
		cache = ctxCache{
			entries: make(map[int64]context.Context),
		}
	})
}

func setContext(ctx context.Context) {
	id := gls.ID()
	cache.set(id, ctx)
}

func getContext() context.Context {
	id := gls.ID()

	ctx, ok := cache.get(id)
	if !ok {
		ctx = context.Background()
		setContext(ctx)
	}

	return ctx
}

func getParentContext() context.Context {
	ctx, ok := cache.get(gls.ParentID())
	if !ok {
		return nil
	}

	return ctx
}
