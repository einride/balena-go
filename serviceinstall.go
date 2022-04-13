package balena

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.einride.tech/balena/odata"
)

const serviceInstallBasePath = "v6/service_install"

type ServiceInstallService service

type InstallsService struct {
	ID          int64        `json:"id"`
	ServiceName string       `json:"service_name"`
	Application odata.Object `json:"application"`
	CreatedAt   time.Time    `json:"created_at"`
}

type ServiceInstallResponse struct {
	InstallsService []InstallsService `json:"installs__service"`
	ID              int64             `json:"id"`
	CreatedAt       time.Time         `json:"created_at"`
	Device          odata.Object      `json:"device"`
}

type ServiceInstalls []*ServiceInstallResponse

// ServiceNames returns a slice of all service names contained in the ServiceInstallResponse(s).
func (i ServiceInstalls) ServiceNames() []string {
	names := make([]string, 0, len(i))
	for _, install := range i {
		for _, installsService := range install.InstallsService {
			names = append(names, installsService.ServiceName)
		}
	}
	return names
}

// FindByServiceName returns the InstallService that contains the given ServiceName.
// Returns true if found, else false.
func (i ServiceInstalls) FindByServiceName(serviceName string) (InstallsService, bool) {
	for _, install := range i {
		for _, installsService := range install.InstallsService {
			if installsService.ServiceName == serviceName {
				return installsService, true
			}
		}
	}
	return InstallsService{}, false
}

// List all service installs for a particular device.
func (s *ServiceInstallService) List(ctx context.Context, deviceID IDOrUUID) (ServiceInstalls, error) {
	query := "%24filter=device+eq+%27" + deviceID.id + "%27"
	if deviceID.isUUID {
		query = "%24filter=device/uuid+eq+%27" + deviceID.id + "%27"
	}
	query += "&%24expand=installs__service(%24select=service_name,application,created_at,id)"

	type response struct {
		D ServiceInstalls `json:"d"`
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, serviceInstallBasePath, query, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}
	resp := &response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request: %v", err)
	}
	return resp.D, nil
}
