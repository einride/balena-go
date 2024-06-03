package balena

import (
	"context"
	"fmt"
	"net/http"
)

const supervisorv1BasePath = "v1"

// SupervisorV1Service handles communication with the supervisor v1 related methods of the
// Balena API.
type SupervisorV1Service struct {
	service
	apiKey     string
	deviceUUID string
	appID      string
	local      bool
}

func (s *SupervisorV1Service) Reboot(ctx context.Context, force bool) error {
	type request struct {
		Force bool `json:"force"`
	}
	req, err := s.newRequest(
		ctx,
		http.MethodPost,
		supervisorv1BasePath+"/reboot",
		&request{Force: force},
	)
	if err != nil {
		return fmt.Errorf("unable to create restart service request: %w", err)
	}
	var resp struct {
		Data  string `json:"Data"`
		Error string `json:"Error"`
	}
	err = s.client.Do(req, &resp)
	if err != nil {
		return fmt.Errorf("unable to reboot device: %w", err)
	}
	if resp.Data != "OK" {
		return fmt.Errorf("unable to reboot device: %v", resp.Error)
	}
	return nil
}

// Blink starts a blink pattern on a LED for 15 seconds, if your device has one.
// https://docs.balena.io/reference/supervisor/supervisor-api/#post-v1blink
func (s *SupervisorV1Service) Blink(ctx context.Context) error {
	req, err := s.newRequest(
		ctx,
		http.MethodPost,
		supervisorv1BasePath+"/blink",
		nil,
	)
	if err != nil {
		return fmt.Errorf("unable to create blink request: %w", err)
	}

	err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("unable to start blink pattern: %w", err)
	}

	return nil
}

// Update triggers a check for the target state of configurations and app services.
// https://docs.balena.io/reference/supervisor/supervisor-api/#post-v1update
func (s *SupervisorV1Service) Update(ctx context.Context, force bool) error {
	type request struct {
		Force bool `json:"force"`
	}
	req, err := s.newRequest(
		ctx,
		http.MethodPost,
		supervisorv1BasePath+"/update",
		&request{Force: force},
	)
	if err != nil {
		return fmt.Errorf("unable to create update request: %w", err)
	}

	err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("unable to trigger update: %w", err)
	}

	return nil
}

type SupervisorV1DeviceResponse struct {
	APIPort           uint   `json:"api_port,omitempty"`
	Commit            string `json:"commit,omitempty"`
	IPAddress         string `json:"ip_address,omitempty"`
	MACAddress        string `json:"mac_address,omitempty"`
	Status            string `json:"status,omitempty"`
	DownloadProgress  uint   `json:"download_progress,omitempty"`
	OSVersion         string `json:"os_version,omitempty"`
	SupervisorVersion string `json:"supervisor_version,omitempty"`
	UpdatePending     bool   `json:"update_pending,omitempty"`
	UpdateDownloaded  bool   `json:"update_downloaded,omitempty"`
	UpdateFailed      bool   `json:"update_failed,omitempty"`
}

// Device returns the current device state.
// https://docs.balena.io/reference/supervisor/supervisor-api/#get-v1device
func (s *SupervisorV1Service) Device(ctx context.Context) (*SupervisorV1DeviceResponse, error) {
	req, err := s.newRequest(
		ctx,
		http.MethodGet,
		supervisorv1BasePath+"/device",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create device request: %w", err)
	}

	var response SupervisorV1DeviceResponse

	err = s.client.Do(req, &response)
	if err != nil {
		return nil, fmt.Errorf("unable to get device state: %w", err)
	}

	return &response, nil
}

func (s *SupervisorV1Service) newRequest(
	ctx context.Context,
	method string,
	urlStr string,
	body interface{},
) (*http.Request, error) {
	u := urlStr
	m := method
	q := fmt.Sprintf("apikey=%s", s.apiKey)
	b := body
	if !s.local {
		u = "supervisor/" + u
		m = http.MethodPost
		q = ""
		b = &CloudRequest{
			UUID:   s.deviceUUID,
			Method: method,
			Data:   body,
		}
	}
	return s.client.NewRequest(ctx, m, u, q, b)
}
