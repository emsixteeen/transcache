package transcache

import (
	"bytes"
	"context"
	"io"
	"sync"
)

type MemoryCache struct {
	mu   sync.Mutex
	data map[string]*bytes.Buffer
}

func (m *MemoryCache) Get(key string) io.Reader {
	m.mu.Lock()
	defer m.mu.Unlock()

	if r, ok := m.data[key]; ok {
		return bytes.NewReader(r.Bytes())
	}

	return nil
}

func (m *MemoryCache) Set(key string) io.Writer {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = new(bytes.Buffer)
	return m.data[key]
}

func (m *MemoryCache) SetCtx(ctx context.Context, key string) io.WriteCloser {
	return &contextWriter{
		cache: m,
		buf:   new(bytes.Buffer),
		key:   key,
		ctx:   ctx,
	}
}

type contextWriter struct {
	cache *MemoryCache
	buf   *bytes.Buffer
	key   string
	ctx   context.Context
}

func (c *contextWriter) Write(p []byte) (int, error) {
	return c.buf.Write(p)
}

func (c *contextWriter) Close() error {
	select {
	case <-c.ctx.Done():
		// Discard, an error
		return c.ctx.Err()
	default:
		c.cache.mu.Lock()
		defer c.cache.mu.Unlock()
		c.cache.data[c.key] = c.buf
	}

	return nil
}
