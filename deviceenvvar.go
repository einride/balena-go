package balena

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"go.einride.tech/balena/odata"
)

const deviceEnvVarBasePath = "v4/device_environment_variable"

type DeviceEnvVarService service

type DeviceEnvVarResponse struct {
	ID        int64        `json:"id,omitempty"`
	CreatedAt string       `json:"created_at,omitempty"`
	Device    odata.Object `json:"device,omitempty"`
	Name      string       `json:"name,omitempty"`
	Value     string       `json:"value,omitempty"`
}

// List lists all environment variables given a specific device ID/UUID.
func (s *DeviceEnvVarService) List(ctx context.Context, deviceID IDOrUUID) ([]*DeviceEnvVarResponse, error) {
	query := "%24filter=device+eq+%27" + deviceID.id + "%27"
	if deviceID.isUUID {
		query = "%24filter=device/uuid+eq+%27" + deviceID.id + "%27"
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, deviceEnvVarBasePath, query, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}
	type Response struct {
		D []*DeviceEnvVarResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request: %v", err)
	}
	return resp.D, nil
}

// Create creates an environment variable with name=value given a device ID/UUID.
func (s *DeviceEnvVarService) Create(
	ctx context.Context,
	deviceID IDOrUUID,
	name string,
	value string,
) (*DeviceEnvVarResponse, error) {
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
		Name     string `json:"name"`
		Value    string `json:"value"`
	}
	req, err := s.client.NewRequest(ctx, http.MethodPost, deviceEnvVarBasePath, "", &request{
		DeviceID: id,
		Name:     name,
		Value:    value,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}
	resp := &DeviceEnvVarResponse{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request: %v", err)
	}
	return resp, nil
}

// DeleteWithName deletes a variable with the given name from the device with given ID/UUID.
// No error is returned if no variable with such name exists.
func (s *DeviceEnvVarService) DeleteWithName(ctx context.Context, deviceID IDOrUUID, name string) error {
	// Get the variable ID
	query := "%24filter=device+eq+%27" + deviceID.id + "%27"
	if deviceID.isUUID {
		query = "%24filter=device/uuid+eq+%27" + deviceID.id + "%27"
	}
	query = query + "+and+name+eq+%27" + name + "%27"
	req, err := s.client.NewRequest(ctx, http.MethodDelete, deviceEnvVarBasePath, query, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %v", err)
	}
	err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("unable to perform request: %v", err)
	}
	return nil
}
