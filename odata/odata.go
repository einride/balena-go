package odata

// EntityURL returns an OData entity ID URL given a base URL and an entity ID.
func EntityURL(base string, id string) string {
	return base + "(" + id + ")"
}

type Object struct {
	Deferred Deferred `json:"__deferred"`
	ID       int64    `json:"__id"`
}

type Deferred struct {
	URI string `json:"uri"`
}
