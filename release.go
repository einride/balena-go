package balena

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/einride/balena-go/odata"
)

const releaseBasePath = "v4/release"

// ReleaseService handles communication with the release related methods of the
// Balena Cloud API.
type ReleaseService service

type ReleaseResponse struct {
	ID                   int64           `json:"id,omitempty"`
	CreatedAt            string          `json:"created_at,omitempty"`
	BelongsToApplication odata.Object    `json:"belongs_to__application,omitempty"`
	CreatedByUser        odata.Object    `json:"is_created_by__user,omitempty"`
	Composition          json.RawMessage `json:"composition,omitempty"`
	Commit               string          `json:"commit,omitempty"`
	Status               string          `json:"status,omitempty"`
	Source               string          `json:"source,omitempty"`
	StartTimestamp       string          `json:"start_timestamp,omitempty"`
	EndTimestamp         string          `json:"end_timestamp,omitempty"`
	UpdateTimestamp      string          `json:"update_timestamp,omitempty"`
}

// List lists all releases.
func (s *ReleaseService) List(ctx context.Context) ([]*ReleaseResponse, error) {
	return s.GetWithQuery(ctx, "")
}

// Get returns a release given a release ID.
// If no such release exists, both the response and error returned are nil.
func (s *ReleaseService) Get(ctx context.Context, id int64) (*ReleaseResponse, error) {
	req, err := s.client.NewRequest(
		ctx,
		http.MethodGet,
		odata.EntityURL(releaseBasePath, strconv.FormatInt(id, 10)),
		"",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create get request: %v", err)
	}
	type Response struct {
		D []ReleaseResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to get release: %v", err)
	}
	if len(resp.D) > 1 {
		return nil, errors.New("received more than 1 release, expected 0 or 1")
	}
	if len(resp.D) == 0 {
		return nil, nil
	}
	return &resp.D[0], nil
}

// GetWithQuery allows querying releases using a custom open data protocol query.
// The query should be a valid, escaped OData query such as `%24filter=uuid+eq+'12333422'`.
//
// Forward slash in filter keys should not be escaped (So `device/uuid` should not be escaped).
func (s *ReleaseService) GetWithQuery(ctx context.Context, query string) ([]*ReleaseResponse, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, releaseBasePath, query, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create release request: %v", err)
	}
	type Response struct {
		D []*ReleaseResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to query release: %v", err)
	}
	return resp.D, nil
}
