package cache

import "time"

type CacheEntry struct {
	Name      string    `json:"name"`
	RouteKey  string    `json:"route_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CacheStr struct {
	Tests   []CacheEntry `json:"tests"`
	Version string       `json:"version"`
}

type TestCase struct {
	Name     string         `json:"name"`
	RouteKey string         `json:"route_key"`
	JSONPool map[string]any `json:"json_pool"`
}
