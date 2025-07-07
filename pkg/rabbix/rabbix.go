package rabbix

type TestCase struct {
	Name     string         `json:"name"`
	RouteKey string         `json:"route_key"`
	JSONPool map[string]any `json:"json_pool"`
}
