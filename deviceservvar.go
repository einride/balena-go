package balena

import (
	"context"
	"fmt"
	"net/http"

	"go.einride.tech/balena/odata"
)

const deviceServVarBasePath = "v6/device_service_environment_variable"

type DeviceServVarService service

type DeviceServVarResponse struct {
	ID             int64           `json:"id,omitempty"`
	CreatedAt      string          `json:"created_at,omitempty"`
	Name           string          `json:"name,omitempty"`
	Value          string          `json:"value,omitempty"`
	ServiceInstall ServiceInstalls `json:"service_install,omitempty"`
}

type DeviceServVarCreateResponse struct {
	ID             int64        `json:"id,omitempty"`
	CreatedAt      string       `json:"created_at,omitempty"`
	Name           string       `json:"name,omitempty"`
	Value          string       `json:"value,omitempty"`
	ServiceInstall odata.Object `json:"service_install,omitempty"`
}

// List lists all environment variables given a specific device ID/UUID.
func (s *DeviceServVarService) List(ctx context.Context, deviceID IDOrUUID) ([]*DeviceServVarResponse, error) {
	query := "%24filter=service_install/device+eq+%27" + deviceID.id + "%27"
	if deviceID.isUUID {
		query = "%24filter=service_install/device/uuid+eq+%27" + deviceID.id + "%27"
	}
	//nolint:lll
	query += "&$expand=service_install($select=id,device,created_at;$expand=installs__service($select=id,service_name,created_at,application))"
	req, err := s.client.NewRequest(ctx, http.MethodGet, deviceServVarBasePath, query, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}
	type Response struct {
		D []*DeviceServVarResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request: %v", err)
	}
	return resp.D, nil
}

// Create creates an environment variable with name=value given a ServiceInstall ID (see ServiceInstallService).
func (s *DeviceServVarService) Create(
	ctx context.Context,
	serviceInstallID int64,
	name string,
	value string,
) (*DeviceServVarCreateResponse, error) {
	type request struct {
		ServiceInstallID int64  `json:"service_install"`
		Name             string `json:"name"`
		Value            string `json:"value"`
	}
	req, err := s.client.NewRequest(ctx, http.MethodPost, deviceServVarBasePath, "", &request{
		ServiceInstallID: serviceInstallID,
		Name:             name,
		Value:            value,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}
	resp := &DeviceServVarCreateResponse{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request: %v", err)
	}
	return resp, nil
}

// Update a variable with the given name from the device with given ID/UUID to a new value.
// No error is returned if no variable with such name exists.
func (s *DeviceServVarService) Update(ctx context.Context, deviceID IDOrUUID, name, newValue string) error {
	query := "%24filter=service_install/device+eq+%27" + deviceID.id + "%27"
	if deviceID.isUUID {
		query = "%24filter=service_install/device/uuid+eq+%27" + deviceID.id + "%27"
	}
	query = query + "+and+name+eq+%27" + name + "%27"
	body := struct {
		Value string `json:"value"`
	}{Value: newValue}
	req, err := s.client.NewRequest(ctx, http.MethodPatch, deviceServVarBasePath, query, body)
	if err != nil {
		return fmt.Errorf("unable to create request: %v", err)
	}
	err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("unable to perform request: %v", err)
	}
	return nil
}

// DeleteWithName deletes a variable with the given name from the device with given ID/UUID.
// No error is returned if no variable with such name exists.
func (s *DeviceServVarService) DeleteWithName(ctx context.Context, deviceID IDOrUUID, name string) error {
	// Get the variable ID
	query := "%24filter=service_install/device+eq+%27" + deviceID.id + "%27"
	if deviceID.isUUID {
		query = "%24filter=service_install/device/uuid+eq+%27" + deviceID.id + "%27"
	}
	query = query + "+and+name+eq+%27" + name + "%27"
	req, err := s.client.NewRequest(ctx, http.MethodDelete, deviceServVarBasePath, query, nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %v", err)
	}
	err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("unable to perform request: %v", err)
	}
	return nil
}
