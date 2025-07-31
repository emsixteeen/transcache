package transcache

import (
	"bytes"
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
