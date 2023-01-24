package balena

import "go.einride.tech/balena/odata"

type Image struct {
	CreatedAt               string       `json:"created_at,omitempty"`
	ModifiedAt              string       `json:"modified_at,omitempty"`
	ID                      int64        `json:"id,omitempty"`
	StartTimestamp          string       `json:"start_timestamp,omitempty"`
	EndTimestamp            string       `json:"end_timestamp,omitempty"`
	Dockerfile              string       `json:"dockerfile,omitempty"`
	IsABuildOfService       odata.Object `json:"is_a_build_of__service,omitempty"`
	ImageSize               int64        `json:"image_size,omitempty"`
	IsStoredAtImageLocation string       `json:"is_stored_at__image_location,omitempty"`
	ProjectType             string       `json:"project_type,omitempty"`
	ErrorMessage            string       `json:"error_message,omitempty"`
	BuildLog                string       `json:"build_log,omitempty"`
	PushTimestamp           string       `json:"push_timestamp,omitempty"`
	Status                  string       `json:"status,omitempty"`
	ContentHash             string       `json:"content_hash,omitempty"`
	Contract                string       `json:"contract,omitempty"`
}

type ImageResponse struct {
	CreatedAt       string       `json:"created_at,omitempty"`
	ID              int64        `json:"id,omitempty"`
	IsPartOfRelease odata.Object `json:"is_part_of__release,omitempty"`
	Image           []*Image     `json:"image,omitempty"`
}
