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
		return fmt.Errorf("unable to create restart service request: %v", err)
	}
	var resp struct {
		Data  string `json:"Data"`
		Error string `json:"Error"`
	}
	err = s.client.Do(req, &resp)
	if err != nil {
		return fmt.Errorf("unable to reboot device: %v", err)
	}
	if resp.Data != "OK" {
		return fmt.Errorf("unable to reboot device: %v", resp.Error)
	}
	return nil
}

func (s *SupervisorV1Service) newRequest(
	ctx context.Context,
	method, urlStr string,
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
