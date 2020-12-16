package balena

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

const supervisorv2BasePath = "v2"

// SupervisorService handles communication with the supervisor related methods of the
// Balena API.
type SupervisorV2Service struct {
	service
	apiKey     string
	deviceUUID string
	appID      string
	local      bool
}

func (s *SupervisorV2Service) RestartServiceByName(ctx context.Context, name string) error {
	type request struct {
		ServiceName string `json:"serviceName"`
	}
	req, err := s.newRequest(
		ctx,
		http.MethodPost,
		supervisorv2BasePath+"/applications/"+s.appID+"/restart-service",
		&request{ServiceName: name},
	)
	if err != nil {
		return fmt.Errorf("unable to create restart service request: %v", err)
	}
	buf := &bytes.Buffer{}
	err = s.client.Do(req, buf)
	if err != nil {
		return fmt.Errorf("unable to restart service %s: %v", name, err)
	}
	return nil
}

func (s *SupervisorV2Service) StopServiceByName(ctx context.Context, name string) error {
	type request struct {
		ServiceName string `json:"serviceName"`
	}
	req, err := s.newRequest(
		ctx,
		http.MethodPost,
		supervisorv2BasePath+"/applications/"+s.appID+"/stop-service",
		&request{ServiceName: name},
	)
	if err != nil {
		return fmt.Errorf("unable to create stop service request: %v", err)
	}
	buf := &bytes.Buffer{}
	err = s.client.Do(req, buf)
	if err != nil {
		return fmt.Errorf("unable to stop service %s: %v", name, err)
	}
	return nil
}

type ApplicationState struct {
	Services map[string]ServiceState `json:"services,omitempty"`
}

type ServiceState struct {
	Status           string      `json:"status,omitempty"`
	ReleaseID        int64       `json:"releaseId,omitempty"`
	DownloadProgress interface{} `json:"download_progress,omitempty"`
}

type SvAppStateResp struct {
	// Local likely always contains only one single entry where the key is the
	// Application ID
	Local     map[string]ApplicationState `json:"local,omitempty"`
	Dependent interface{}                 `json:"dependent,omitempty"`
	Commit    string                      `json:"commit,omitempty"`
}

func (s *SupervisorV2Service) ApplicationState(ctx context.Context) (*SvAppStateResp, error) {
	req, err := s.newRequest(
		ctx,
		http.MethodGet,
		supervisorv2BasePath+"/applications/"+s.appID+"/state",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create application state request: %v", err)
	}
	resp := &SvAppStateResp{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch application state: %v", err)
	}
	return resp, nil
}

type CloudRequest struct {
	UUID   string      `json:"uuid,omitempty"`
	Method string      `json:"method,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func (s *SupervisorV2Service) newRequest(
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
