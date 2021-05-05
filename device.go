package balena

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go.einride.tech/balena/odata"
)

const deviceBasePath = "v4/device"

// DeviceService handles communication with the device related methods of the
// Balena Cloud API.
type DeviceService service

type DeviceResponse struct {
	ID         int64  `json:"id,omitempty"`
	Actor      int64  `json:"actor,omitempty"`
	Name       string `json:"device_name,omitempty"`
	DeviceType string `json:"device_type,omitempty"`
	UUID       string `json:"uuid,omitempty"`
	OnCommit   string `json:"is_on__commit,omitempty"`
	// TODO: Change to time.Time maybe?
	LastConnectivityEvent string `json:"last_connectivity_event,omitempty"`
	Status                string `json:"status,omitempty"`
	LastVPNEvent          string `json:"last_vpn_event,omitempty"`
	// TODO: Should we change to net.IP maybe?
	IPAddress string `json:"ip_address,omitempty"`
	// TODO: Should we change to net.IP maybe?
	VPNAddress string `json:"vpn_address,omitempty"`
	// TODO: Should we change to net.IP maybe?
	PublicAddress              string        `json:"public_address,omitempty"`
	OSVersion                  string        `json:"os_version,omitempty"`
	OSVariant                  string        `json:"os_variant,omitempty"`
	SupervisorVersion          string        `json:"supervisor_version,omitempty"`
	ProvisioningState          string        `json:"provisioning_state,omitempty"`
	Longitude                  string        `json:"longitude,omitempty"`
	Latitude                   string        `json:"latitude,omitempty"`
	Location                   string        `json:"location,omitempty"`
	CustomLongitude            string        `json:"custom_longitude,omitempty"`
	CustomLatitude             string        `json:"custom_latitude,omitempty"`
	CreatedAt                  string        `json:"created_at,omitempty"`
	IsOnline                   bool          `json:"is_online,omitempty"`
	IsConnectedToVPN           bool          `json:"is_connected_to_vpn,omitempty"`
	IsWebAccessible            bool          `json:"is_web_accessible,omitempty"`
	IsActive                   bool          `json:"is_active,omitempty"`
	BelongsToApplication       odata.Object  `json:"belongs_to__application,omitempty"`
	BelongsToUser              odata.Object  `json:"belongs_to__user,omitempty"`
	IsManagedByServiceInstance odata.Object  `json:"is_managed_by__service_instance,omitempty"`
	IsManagedByDevice          interface{}   `json:"is_managed_by__device,omitempty"`
	ShouldBeRunningRelease     *odata.Object `json:"should_be_running__release,omitempty"`
	Note                       interface{}   `json:"note,omitempty"`
	LocalID                    interface{}   `json:"local_id,omitempty"`
	ManagedBySupervisorRelease interface{}   `json:"should_be_managed_by__supervisor_release,omitempty"`
	ProvisioningProgress       interface{}   `json:"provisioning_progress,omitempty"`
	DownloadProgress           interface{}   `json:"download_progress,omitempty"`
	LogsChannel                interface{}   `json:"logs_channel,omitempty"`
	IsLockedUntil              interface{}   `json:"is_locked_until__date,omitempty"`
	IsAccessibleBySupportUntil interface{}   `json:"is_accessible_by_support_until__date,omitempty"`
	MACAddress                 interface{}   `json:"mac_address,omitempty"`
}

// List returns a list of all devices.
func (s *DeviceService) List(ctx context.Context) ([]*DeviceResponse, error) {
	return s.GetWithQuery(ctx, "")
}

// ListByApplication returns a list of all devices owned by the single application given its ID.
func (s *DeviceService) ListByApplication(ctx context.Context, applicationID int64) ([]*DeviceResponse, error) {
	query := "%24filter=belongs_to__application%20eq%20%27" + strconv.FormatInt(applicationID, 10) + "%27"
	return s.GetWithQuery(ctx, query)
}

// Get returns information on a single device given its ID or UUID.
// If the device does not exist, both the response and error are nil.
func (s *DeviceService) Get(ctx context.Context, deviceID IDOrUUID) (*DeviceResponse, error) {
	var query string
	path := odata.EntityURL(deviceBasePath, deviceID.id)
	if deviceID.isUUID {
		query = "%24filter=uuid+eq+%27" + deviceID.id + "%27"
		path = deviceBasePath
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create get request: %v", err)
	}
	type Response struct {
		D []DeviceResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to get device: %v", err)
	}
	if len(resp.D) > 1 {
		return nil, errors.New("received more than 1 device, expected 0 or 1")
	}
	if len(resp.D) == 0 {
		return nil, nil
	}
	return &resp.D[0], nil
}

// GetWithQuery allows querying for devices using a custom open data protocol query.
// The query should be a valid, escaped OData query such as `%24filter=uuid+eq+'12333422'`
//
// Forward slash in filter keys should not be escaped (So `device/uuid` should not be escaped).
func (s *DeviceService) GetWithQuery(ctx context.Context, query string) ([]*DeviceResponse, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, deviceBasePath, query, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create device request: %v", err)
	}
	type Response struct {
		D []*DeviceResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to query device: %v", err)
	}
	return resp.D, nil
}

// PinRelease pins a device to a specific release.
func (s *DeviceService) PinRelease(ctx context.Context, deviceID IDOrUUID, releaseID int64) ([]byte, error) {
	type request struct {
		ShouldRunRelease string `json:"should_be_running__release"`
	}
	release := strconv.FormatInt(releaseID, 10)
	var query string
	path := odata.EntityURL(deviceBasePath, deviceID.id)
	if deviceID.isUUID {
		query = "%24filter=uuid+eq+%27" + deviceID.id + "%27"
		path = deviceBasePath
	}
	req, err := s.client.NewRequest(
		ctx,
		http.MethodPatch,
		path,
		query,
		&request{ShouldRunRelease: release},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create setRelease request: %v", err)
	}
	buf := &bytes.Buffer{}
	err = s.client.Do(req, buf)
	if err != nil {
		return nil, fmt.Errorf("unable to patch device: %v", err)
	}
	return buf.Bytes(), nil
}

// TrackLatestRelease sets a device to track the latest available release.
func (s *DeviceService) TrackLatestRelease(ctx context.Context, deviceID IDOrUUID) ([]byte, error) {
	type request struct {
		ShouldRunRelease interface{} `json:"should_be_running__release"`
	}
	var query string
	path := odata.EntityURL(deviceBasePath, deviceID.id)
	if deviceID.isUUID {
		query = "%24filter=uuid+eq+%27" + deviceID.id + "%27"
		path = deviceBasePath
	}
	req, err := s.client.NewRequest(
		ctx,
		http.MethodPatch,
		path,
		query,
		&request{ShouldRunRelease: nil},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create setRelease request: %v", err)
	}
	buf := &bytes.Buffer{}
	err = s.client.Do(req, buf)
	if err != nil {
		return nil, fmt.Errorf("unable to patch device: %v", err)
	}
	return buf.Bytes(), nil
}
