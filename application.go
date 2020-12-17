package balena

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/einride/balena-go/odata"
)

const applicationBasePath = "v4/application"

// ApplicationService handles communication with the application related methods of the
// Balena Cloud API.
type ApplicationService service

type ApplicationsResponse struct {
	ID                         int64        `json:"id,omitempty"`
	User                       odata.Object `json:"user,omitempty"`
	DependsOnApplication       interface{}  `json:"depends_on__application,omitempty"`
	Actor                      int64        `json:"actor,omitempty"`
	Name                       string       `json:"app_name,omitempty"`
	Slug                       string       `json:"slug,omitempty"`
	Commit                     string       `json:"commit,omitempty"`
	ApplicationType            odata.Object `json:"application_type,omitempty"`
	DeviceType                 string       `json:"device_type,omitempty"`
	TrackLatestRelease         bool         `json:"should_track_latest_release,omitempty"`
	IsAccessibleBySupportUntil interface{}  `json:"is_accessible_by_support_until__date,omitempty"`
}

func (s *ApplicationService) List(ctx context.Context) ([]*ApplicationsResponse, error) {
	return s.GetWithQuery(ctx, "")
}

// GetWithQuery allows querying for devices using a custom open data protocol query.
// The query should be a valid, escaped OData query such as `%24filter=uuid+eq+'12333422'`.
//
// Forward slash in filter keys should not be escaped (So `device/uuid` should not be escaped).
func (s *ApplicationService) GetWithQuery(ctx context.Context, query string) ([]*ApplicationsResponse, error) {
	return s.getWithQueryAndPath(ctx, applicationBasePath, query)
}

// Get returns information on a single application given its ID.
// If the application does not exist, both the response and error are nil.
func (s *ApplicationService) Get(ctx context.Context, applicationID int64) (*ApplicationsResponse, error) {
	path := odata.EntityURL(applicationBasePath, strconv.FormatInt(applicationID, 10))
	resp, err := s.getWithQueryAndPath(ctx, path, "")
	if len(resp) > 1 {
		return nil, errors.New("received more than 1 application, expected 0 or 1")
	}
	if len(resp) == 0 {
		return nil, nil
	}
	return resp[0], err
}

// GetByName returns information on a single application given its Name
// If the application does not exist, both the response and error are nil.
func (s *ApplicationService) GetByName(ctx context.Context, applicationName string) (*ApplicationsResponse, error) {
	query := "%24filter=app_name%20eq%20%27" + applicationName + "%27"
	resp, err := s.getWithQueryAndPath(ctx, applicationBasePath, query)
	if len(resp) > 1 {
		return nil, errors.New("received more than 1 application, expected 0 or 1")
	}
	if len(resp) == 0 {
		return nil, nil
	}
	return resp[0], err
}

func (s *ApplicationService) getWithQueryAndPath(
	ctx context.Context,
	path string,
	query string,
) ([]*ApplicationsResponse, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create get request: %v", err)
	}
	type Response struct {
		D []ApplicationsResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to get application list: %v", err)
	}
	apps := make([]*ApplicationsResponse, 0, len(resp.D))
	for _, app := range resp.D {
		app := app
		apps = append(apps, &app)
	}
	return apps, nil
}

// EnableTrackLatestRelease sets all devices owned by the application to track the latest available release.
func (s *ApplicationService) EnableTrackLatestRelease(ctx context.Context, applicationID int64) ([]byte, error) {
	type request struct {
		ShouldTrackLatestRelease bool `json:"should_track_latest_release"`
	}
	var query string
	path := odata.EntityURL(applicationBasePath, strconv.FormatInt(applicationID, 10))
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, query, &request{ShouldTrackLatestRelease: true})
	if err != nil {
		return nil, fmt.Errorf("unable to create setTrackLatestRelease request: %v", err)
	}
	var buf bytes.Buffer
	err = s.client.Do(req, &buf)
	if err != nil {
		return nil, fmt.Errorf("unable to path application: %v", err)
	}
	return buf.Bytes(), nil
}

// DisableTrackLatestRelease sets all devices owned by the application to NOT track the latest available release.
func (s *ApplicationService) DisableTrackLatestRelease(ctx context.Context, applicationID int64) ([]byte, error) {
	type request struct {
		ShouldTrackLatestRelease bool `json:"should_track_latest_release"`
	}
	var query string
	path := odata.EntityURL(applicationBasePath, strconv.FormatInt(applicationID, 10))
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, query, &request{ShouldTrackLatestRelease: false})
	if err != nil {
		return nil, fmt.Errorf("unable to create setTrackLatestRelease request: %v", err)
	}
	var buf bytes.Buffer
	err = s.client.Do(req, &buf)
	if err != nil {
		return nil, fmt.Errorf("unable to path application: %v", err)
	}
	return buf.Bytes(), nil
}
