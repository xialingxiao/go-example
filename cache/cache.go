package cache

//Storage mecanism for caching strings
type Storage interface {
    Get() map[string]float64
    Set(content map[string]float64, expiration int64)
}
