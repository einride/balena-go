package balena

type CPUArchitectureResponse struct {
	ID   int64  `json:"id,omitempty"`
	Slug string `json:"slug,omitempty"`
	Name string `json:"name,omitempty"`
}
