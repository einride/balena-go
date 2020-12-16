package balena

import (
	"bytes"
	"context"
	"net/http"
	"strconv"

	"github.com/einride/balena-go/odata"
	"golang.org/x/xerrors"
)

const deviceTagBasePath = "v4/device_tag"

// DeviceTagService handles communication with the device tag related methods of the
// Balena Cloud API.
type DeviceTagService service

type DeviceTagResponse struct {
	ID     int64        `json:"id,omitempty"`
	Device odata.Object `json:"device,omitempty"`
	TagKey string       `json:"tag_key,omitempty"`
	Value  string       `json:"value,omitempty"`
}

// List lists all device tags for a given device ID/UUID.
func (s *DeviceTagService) List(ctx context.Context, deviceID IDOrUUID) ([]*DeviceTagResponse, error) {
	query := "%24filter=device/id+eq+%27" + deviceID.id + "%27"
	if deviceID.isUUID {
		query = "%24filter=device/uuid+eq+%27" + deviceID.id + "%27"
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, deviceTagBasePath, query, nil)
	if err != nil {
		return nil, xerrors.Errorf("list device tag NewRequest: %v", err)
	}
	type Response struct {
		D []*DeviceTagResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, xerrors.Errorf("list device tag: %v", err)
	}
	return resp.D, nil
}

// Create creates a device tag with key=value given a device ID/UUID.
// An error is returned if the key already exists.
func (s *DeviceTagService) Create(
	ctx context.Context,
	deviceID IDOrUUID,
	key string,
	value string,
) (*DeviceTagResponse, error) {
	// If UUID, retrieve device ID
	id := deviceID.id
	if deviceID.isUUID {
		resp, err := s.client.Device.Get(ctx, deviceID)
		if err != nil {
			return nil, err
		}
		id = strconv.FormatInt(resp.ID, 10)
	}
	type request struct {
		DeviceID string `json:"device"`
		Key      string `json:"tag_key"`
		Value    string `json:"value"`
	}
	req, err := s.client.NewRequest(ctx, http.MethodPost, deviceTagBasePath, "", &request{
		DeviceID: id,
		Key:      key,
		Value:    value,
	})
	if err != nil {
		return nil, xerrors.Errorf("create device tag NewRequest: %v", err)
	}
	resp := &DeviceTagResponse{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, xerrors.Errorf("create device tag: %v", err)
	}
	return resp, nil
}

// GetWithKey retrieves a tag with the given key from the given device ID/UUID.
// If no key is found both the response and error returned are nil.
func (s *DeviceTagService) GetWithKey(ctx context.Context, deviceID IDOrUUID, key string) (*DeviceTagResponse, error) {
	// Get the variable ID
	query := "%24filter=device/id+eq+%27" + deviceID.id + "%27"
	if deviceID.isUUID {
		query = "%24filter=device/uuid+eq+%27" + deviceID.id + "%27"
	}
	query = query + "+and+tag_key+eq+%27" + key + "%27"
	req, err := s.client.NewRequest(ctx, http.MethodGet, deviceTagBasePath, query, nil)
	if err != nil {
		return nil, xerrors.Errorf("get device tag with key NewRequest: %v", err)
	}
	type Response struct {
		D []*DeviceTagResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, xerrors.Errorf("get device tag with key: %v", err)
	}
	if len(resp.D) > 1 {
		return nil, xerrors.Errorf("expected 1 tag but got %d", len(resp.D))
	}

	if len(resp.D) == 0 {
		return nil, nil
	}
	return resp.D[0], nil
}

// UpdateWithKey updates the value of a device tag matching the given key and device ID/UUID.
// No error is returned if the key or device does not exist.
func (s *DeviceTagService) UpdateWithKey(ctx context.Context, deviceID IDOrUUID, key string, value string) error {
	// Get the variable ID
	query := "%24filter=device/id+eq+%27" + deviceID.id + "%27"
	if deviceID.isUUID {
		query = "%24filter=device/uuid+eq+%27" + deviceID.id + "%27"
	}
	query = query + "+and+tag_key+eq+%27" + key + "%27"
	type request struct {
		Value string `json:"value"`
	}
	req, err := s.client.NewRequest(ctx, http.MethodPatch, deviceTagBasePath, query, &request{
		Value: value,
	})
	if err != nil {
		return xerrors.Errorf("update device tag with key NewRequest: %v", err)
	}
	buf := &bytes.Buffer{}
	err = s.client.Do(req, buf)
	if err != nil {
		return xerrors.Errorf("update device tag with key: %v", err)
	}
	return nil
}

// DeleteWithKey deletes a device tag from a given device ID/UUID and key.
// No error is returned if the tag does not exist.
func (s *DeviceTagService) DeleteWithKey(ctx context.Context, deviceID IDOrUUID, key string) error {
	// Get the variable ID
	query := "%24filter=device/id+eq+%27" + deviceID.id + "%27"
	if deviceID.isUUID {
		query = "%24filter=device/uuid+eq+%27" + deviceID.id + "%27"
	}
	query = query + "+and+tag_key+eq+%27" + key + "%27"
	req, err := s.client.NewRequest(ctx, http.MethodDelete, deviceTagBasePath, query, nil)
	if err != nil {
		return xerrors.Errorf("delete device tag with key NewRequest: %v", err)
	}
	buf := &bytes.Buffer{}
	err = s.client.Do(req, buf)
	if err != nil {
		return xerrors.Errorf("delete device tag with key: %v", err)
	}
	return nil
}

// GetWithQuery allows querying for device tags using a custom Open Data Protocol query.
// The query should be a valid, escaped OData query such as `%24filter=uuid+eq+%2712333422%27`.
//
// Forward slash in filter keys should not be escaped (So `device/uuid` should not be escaped).
func (s *DeviceTagService) GetWithQuery(ctx context.Context, query string) ([]*DeviceTagResponse, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, deviceTagBasePath, query, nil)
	if err != nil {
		return nil, xerrors.Errorf("get device tag with query NewRequest: %v", err)
	}
	type Response struct {
		D []*DeviceTagResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, xerrors.Errorf("get device tag with query: %v", err)
	}
	return resp.D, nil
}
