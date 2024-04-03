package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/kulinsky/gotestexample/internal/common"
)

type InMemoryRepository struct {
	store map[string]string
	mu    sync.RWMutex
}

func New() *InMemoryRepository {
	return &InMemoryRepository{
		store: make(map[string]string),
	}
}

func (r *InMemoryRepository) set(key, val string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[key] = val
}

func (r *InMemoryRepository) get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.store[key]

	return v, ok
}

func (r *InMemoryRepository) Save(ctx context.Context, id, fullURL string) error {
	r.set(id, fullURL)

	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (string, error) {
	v, ok := r.get(id)
	if !ok {
		return "", fmt.Errorf("%w: store has no key", common.ErrNotFound)
	}

	return v, nil
}
