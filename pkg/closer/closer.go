package closer

import (
	"context"
	"log/slog"
	"sync"

	"go.uber.org/multierr"
)

type closeFn func(ctx context.Context) error

type item struct {
	name string
	fn   closeFn
}

type Closer struct {
	log *slog.Logger
	mu sync.Mutex
	items []item
}

func New(log *slog.Logger) *Closer {
	return &Closer{log: log}
}

func (c *Closer) Add(name string, fn closeFn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = append(c.items, item{name: name, fn: fn})
}

func (c *Closer) AddFunc(name string, fn func()) {
	c.Add(name, func(ctx context.Context) error {
		fn()

		return nil
	})
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	items := append([]item(nil), c.items...)
	c.mu.Unlock()

	var result error
	for _, i := range items {
		if err := i.fn(ctx); err != nil {
			result = multierr.Append(result, err)
			c.log.Error("shutdown hook faileld", "name", i.name, "error", err)
		}
	}
	c.log.Info("Shutdown hook finished", "result", result)
	
	return result
}