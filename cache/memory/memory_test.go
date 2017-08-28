package memory_test

import (
    "testing"
    "time"
    "reflect"
    "github.com/xialingxiao/go-example/cache/memory"
)

func TestGetEmpty(t *testing.T) {
    storage := memory.NewStorage()
    content := storage.Get()
    if len(content) != 0 {
        t.Errorf("Expected an empty map but got '%s'", content)
    }
}

func TestGetValue(t *testing.T) {
    storage := memory.NewStorage()
    storage.Set(map[string]float64{"SGD": 1.355804}, time.Now().Unix()+10)
    time.Sleep(time.Duration(1)*time.Second)
    content := storage.Get()
    if !reflect.DeepEqual(content, map[string]float64{"SGD": 1.355804}) {
        t.Errorf(`Expected to get map[SGD:1.355804] but got '%s'`, content)
    }
}

func TestGetExpiredValue(t *testing.T) {
    storage := memory.NewStorage()
    storage.Set(map[string]float64{"SGD": 1.355804}, time.Now().Unix())
    time.Sleep(time.Duration(1)*time.Second)
    content := storage.Get()
    if content != nil {
        t.Errorf(`Expected nil but got '%s'`, content)
    }
}