package balena

import "go.einride.tech/balena/odata"

type DeviceType struct {
	ID           uint64       `json:"id"`
	Name         string       `json:"name"`
	Slug         string       `json:"slug"`
	UUID         string       `json:"uuid"`
	IsPrivate    bool         `json:"is_private"`
	Architecture odata.Object `json:"is_of__cpu_architecture"`
	BelongsToFamily     *odata.Object `json:"belongs_to__device_family,omitempty"`
}
