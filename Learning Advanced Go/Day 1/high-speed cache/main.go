package main

import (
	"fmt"
	"sync"
)

type PriceCache struct {
	mu     sync.RWMutex
	prices map[string]int
}

func (c *PriceCache) getPrice(item string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.prices[item]
}

func (c *PriceCache) updatePrice(item string, price int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prices[item] = price
}

func main() {
	cache := &PriceCache{prices: make(map[string]int)}
	cache.updatePrice("Bitcoin", 45000)

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			_ = cache.getPrice("Bitcoin")
		}()
	}

	wg.Wait()
	fmt.Println("Successfully read prices concurrently!")
}
