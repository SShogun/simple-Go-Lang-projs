package main

import (
	"sync"
)

type RWData struct {
	value map[string]string
	mu    sync.RWMutex
}

func Get(data *RWData, key string) string {
	data.mu.RLock()
	defer data.mu.RUnlock()
	return data.value[key]
}

func Set(data *RWData, key, val string) {
	data.mu.Lock()
	defer data.mu.Unlock()
	data.value[key] = val
}

func main() {
	data := &RWData{
		value: make(map[string]string),
	}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + string(i%10+'0')
			Set(data, key, "value"+string(i+'0'))
			_ = Get(data, key)
		}(i)
	}
	wg.Wait()

}
