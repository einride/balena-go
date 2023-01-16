package balena

import (
	"encoding/json"

	"go.einride.tech/balena/odata"
)

type CPUArchitectureResponse struct {
	ID   int64  `json:"id,omitempty"`
	Slug string `json:"slug,omitempty"`
	Name string `json:"name,omitempty"`
}

type CPUArchitectureOData struct {
	D []*CPUArchitectureResponse
	*odata.Object
}

func (d *CPUArchitectureOData) UnmarshalJSON(data []byte) error {
	o := new(odata.Object)
	if err := json.Unmarshal(data, o); err == nil {
		d.Object = o
		return nil
	}

	return json.Unmarshal(data, &d.D)
}
