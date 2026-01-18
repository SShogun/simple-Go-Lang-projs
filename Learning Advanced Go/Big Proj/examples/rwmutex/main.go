package main

import (
	"fmt"
	"sync"
)

type RWData struct {
	mu    sync.RWMutex
	value map[string]string
}

// NewStore initializes the struct properly
func NewStore() *RWData {
	return &RWData{
		value: make(map[string]string),
	}
}

// Get uses a Read Lock (multiple readers allowed)
func (d *RWData) Get(key string) string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.value[key]
}

// Set uses a Write Lock (exclusive access)
func (d *RWData) Set(key, val string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.value[key] = val
}

func main() {
	store := NewStore()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := fmt.Sprintf("key-%d", i%10)
			v := fmt.Sprintf("value-%d", i)

			store.Set(k, v)
			_ = store.Get(k)
		}(i)
	}

	wg.Wait()
	fmt.Println("Final value for key-0:", store.Get("key-0"))
}
