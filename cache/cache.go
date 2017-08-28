package cache

//Storage mecanism for caching
type Storage interface {
    Get() (map[string]float64, int64)
    Set(content map[string]float64, expiration int64)
}
