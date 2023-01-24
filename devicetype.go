package balena

import (
	"encoding/json"

	"go.einride.tech/balena/odata"
)

type DeviceType struct {
	ID           uint64       `json:"id,omitempty"`
	Name         string       `json:"name,omitempty"`
	Slug         string       `json:"slug,omitempty"`
	UUID         string       `json:"uuid,omitempty"`
	IsPrivate    bool         `json:"is_private,omitempty"`
	Architecture odata.Object `json:"is_of__cpu_architecture,omitempty"`
	BelongsToFamily     *odata.Object `json:"belongs_to__device_family,omitempty"`
}

type DeviceTypeOData struct {
	D []*DeviceTypeResponse
	*odata.Object
}

func (d *DeviceTypeOData) UnmarshalJSON(data []byte) error {
	o := new(odata.Object)
	if err := json.Unmarshal(data, o); err == nil {
		d.Object = o
		return nil
	}

	return json.Unmarshal(data, &d.D)
}
