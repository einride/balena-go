package balena

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go.einride.tech/balena/odata"
)

const deviceTypeBasePath = "v6/device_type"

// DeviceTypeService handles communication with the device type related methods of the
// Balena Cloud API.
type DeviceTypeService service

type DeviceTypeResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	IsPrivate bool   `json:"is_private"`
	// Logo will only be populated when explicitly selected through a query parameter: `$select=logo`.
	Logo                string       `json:"logo"`
	IsOfCPUArchitecture odata.Object `json:"is_of__cpu_architecture"`
}

// Get returns information on a single device type given its ID.
// If the device type does not exist, both the response and error are nil.
func (s *DeviceTypeService) Get(ctx context.Context, id int64) (*DeviceTypeResponse, error) {
	path := odata.EntityURL(deviceTypeBasePath, strconv.FormatInt(id, 10))
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, "", nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create get request: %w", err)
	}
	type Response struct {
		D []DeviceTypeResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to get device type: %w", err)
	}
	if len(resp.D) > 1 {
		return nil, errors.New("received more than 1 device type, expected 0 or 1")
	}
	if len(resp.D) == 0 {
		return nil, nil
	}
	return &resp.D[0], nil
}

// GetWithQuery allows querying for device types using a custom open data protocol query.
// The query should be a valid, escaped OData query such as `%24filter=slug+eq+'jetson-tx2'`
//
// Forward slash in filter keys should not be escaped (So `device_type/slug` should not be escaped).
func (s *DeviceTypeService) GetWithQuery(ctx context.Context, query string) ([]*DeviceTypeResponse, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, deviceTypeBasePath, query, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create device type request: %w", err)
	}
	type Response struct {
		D []*DeviceTypeResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to query device type: %w", err)
	}
	return resp.D, nil
}
