package memory

import (
    "sync"
    "time"
)

// Item is a cached reference
type Item struct {
    Rates    map[string]float64
    Expiration int64
}

// Expired returns true if the item has expired.
func (item Item) Expired() bool {
    if item.Expiration == 0 {
        return false // never expire the data if Expiration is not set
    }

    return time.Now().Unix() > item.Expiration
}

//Storage mecanism for caching strings in memory
type Storage struct {
    item Item
    mu   *sync.RWMutex
}

//NewStorage creates a new in memory storage
func NewStorage() *Storage {
    return &Storage{
        item: Item{},
        mu:   &sync.RWMutex{},
    }
}

//Get cached rates
func (s Storage) Get() (map[string]float64, int64) {

    s.mu.RLock()
    defer s.mu.RUnlock()

    item := s.item

    if item.Expired() {
        s.item = Item{}
        return nil, 0
    }

    return item.Rates, item.Expiration
}

//Set a cached rates
func (s *Storage) Set(rates map[string]float64, expiration int64) {
    (*s).mu.Lock()
    defer (*s).mu.Unlock()
    (*s).item = Item{rates, expiration}
}
