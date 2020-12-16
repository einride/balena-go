package balena

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/einride/balena-go/odata"
	"github.com/stretchr/testify/require"
)

const (
	deviceResponse = `{
	"d": [
		{
			"id": 1702297,
			"belongs_to__application": {
				"__deferred": {
					"uri": "/resin/application(1234587)"
				},
				"__id": 1234587
			},
			"belongs_to__user": {
				"__deferred": {
					"uri": "/resin/user(126547)"
				},
				"__id": 126547
			},
			"is_managed_by__device": null,
			"actor": 4081727,
			"should_be_running__release": null,
			"device_name": "old-moon",
			"device_type": "qemux86-64",
			"uuid": "00d859f123685e84772676f09465cc55",
			"is_on__commit": "0dfe02d12efb2325385952e1e15361f7",
			"note": null,
			"local_id": null,
			"status": "Idle",
			"is_online": true,
			"last_connectivity_event": "2019-09-26T06:50:31.238Z",
			"is_connected_to_vpn": true,
			"last_vpn_event": "2019-09-26T06:50:31.238Z",
			"ip_address": "10.0.2.15",
			"vpn_address": "10.240.63.4",
			"public_address": "83.241.129.195",
			"os_version": "balenaOS 2.38.0+rev1",
			"os_variant": "dev",
			"supervisor_version": "9.15.7",
			"should_be_managed_by__supervisor_release": null,
			"is_managed_by__service_instance": {
				"__deferred": {
					"uri": "/resin/service_instance(123987)"
				},
				"__id": 123987
			},
			"provisioning_progress": null,
			"provisioning_state": "",
			"download_progress": null,
			"is_web_accessible": false,
			"longitude": "11.9672",
			"latitude": "57.7066",
			"location": "Gothenburg, Västra Götaland, Sweden",
			"custom_longitude": "",
			"custom_latitude": "",
			"logs_channel": null,
			"is_locked_until__date": null,
			"is_accessible_by_support_until__date": null,
			"created_at": "2019-09-23T11:55:58.654Z",
			"is_active": true,
			"__metadata": {
				"uri": "/resin/device(@id)?@id=1702297"
			}
		}
	]
}`
)

func TestDeviceService_List_ID(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, deviceResponse)
	})
	expected := []*DeviceResponse{
		{
			ID: 1702297,
			BelongsToApplication: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/application(1234587)"}, ID: 1234587,
			},
			BelongsToUser: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/user(126547)"}, ID: 126547,
			},
			IsManagedByDevice:          nil,
			Actor:                      4081727,
			ShouldBeRunningRelease:     nil,
			Name:                       "old-moon",
			DeviceType:                 "qemux86-64",
			UUID:                       "00d859f123685e84772676f09465cc55",
			OnCommit:                   "0dfe02d12efb2325385952e1e15361f7",
			Note:                       nil,
			LocalID:                    nil,
			Status:                     "Idle",
			IsOnline:                   true,
			LastConnectivityEvent:      "2019-09-26T06:50:31.238Z",
			IsConnectedToVPN:           true,
			LastVPNEvent:               "2019-09-26T06:50:31.238Z",
			IPAddress:                  "10.0.2.15",
			VPNAddress:                 "10.240.63.4",
			PublicAddress:              "83.241.129.195",
			OSVersion:                  "balenaOS 2.38.0+rev1",
			OSVariant:                  "dev",
			SupervisorVersion:          "9.15.7",
			ManagedBySupervisorRelease: nil,
			IsManagedByServiceInstance: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/service_instance(123987)"}, ID: 123987,
			},
			ProvisioningProgress:       nil,
			ProvisioningState:          "",
			DownloadProgress:           nil,
			IsWebAccessible:            false,
			Longitude:                  "11.9672",
			Latitude:                   "57.7066",
			Location:                   "Gothenburg, Västra Götaland, Sweden",
			CustomLongitude:            "",
			CustomLatitude:             "",
			LogsChannel:                nil,
			IsLockedUntil:              nil,
			IsAccessibleBySupportUntil: nil,
			CreatedAt:                  "2019-09-23T11:55:58.654Z",
			IsActive:                   true,
		},
	}
	// When
	actual, err := client.Device.List(context.Background())
	// Then
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestDeviceService_ListByApplication(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	applicationID := int64(1234587)
	mux.HandleFunc("/"+deviceBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=belongs_to__application%20eq%20%271234587%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		fmt.Fprint(w, deviceResponse)
	})
	expected := []*DeviceResponse{
		{
			ID: 1702297,
			BelongsToApplication: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/application(1234587)"}, ID: 1234587,
			},
			BelongsToUser: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/user(126547)"}, ID: 126547,
			},
			IsManagedByDevice:          nil,
			Actor:                      4081727,
			ShouldBeRunningRelease:     nil,
			Name:                       "old-moon",
			DeviceType:                 "qemux86-64",
			UUID:                       "00d859f123685e84772676f09465cc55",
			OnCommit:                   "0dfe02d12efb2325385952e1e15361f7",
			Note:                       nil,
			LocalID:                    nil,
			Status:                     "Idle",
			IsOnline:                   true,
			LastConnectivityEvent:      "2019-09-26T06:50:31.238Z",
			IsConnectedToVPN:           true,
			LastVPNEvent:               "2019-09-26T06:50:31.238Z",
			IPAddress:                  "10.0.2.15",
			VPNAddress:                 "10.240.63.4",
			PublicAddress:              "83.241.129.195",
			OSVersion:                  "balenaOS 2.38.0+rev1",
			OSVariant:                  "dev",
			SupervisorVersion:          "9.15.7",
			ManagedBySupervisorRelease: nil,
			IsManagedByServiceInstance: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/service_instance(123987)"}, ID: 123987,
			},
			ProvisioningProgress:       nil,
			ProvisioningState:          "",
			DownloadProgress:           nil,
			IsWebAccessible:            false,
			Longitude:                  "11.9672",
			Latitude:                   "57.7066",
			Location:                   "Gothenburg, Västra Götaland, Sweden",
			CustomLongitude:            "",
			CustomLatitude:             "",
			LogsChannel:                nil,
			IsLockedUntil:              nil,
			IsAccessibleBySupportUntil: nil,
			CreatedAt:                  "2019-09-23T11:55:58.654Z",
			IsActive:                   true,
		},
	}
	// When
	actual, err := client.Device.ListByApplication(context.Background(), applicationID)
	// Then
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestDeviceService_Get_ID(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityID := int64(112233)
	mux.HandleFunc(
		"/"+odata.EntityURL(deviceBasePath, strconv.FormatInt(entityID, 10)),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, deviceResponse)
		},
	)
	expected := &DeviceResponse{
		ID: 1702297,
		BelongsToApplication: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/application(1234587)"}, ID: 1234587,
		},
		BelongsToUser: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/user(126547)"}, ID: 126547,
		},
		IsManagedByDevice:          nil,
		Actor:                      4081727,
		ShouldBeRunningRelease:     nil,
		Name:                       "old-moon",
		DeviceType:                 "qemux86-64",
		UUID:                       "00d859f123685e84772676f09465cc55",
		OnCommit:                   "0dfe02d12efb2325385952e1e15361f7",
		Note:                       nil,
		LocalID:                    nil,
		Status:                     "Idle",
		IsOnline:                   true,
		LastConnectivityEvent:      "2019-09-26T06:50:31.238Z",
		IsConnectedToVPN:           true,
		LastVPNEvent:               "2019-09-26T06:50:31.238Z",
		IPAddress:                  "10.0.2.15",
		VPNAddress:                 "10.240.63.4",
		PublicAddress:              "83.241.129.195",
		OSVersion:                  "balenaOS 2.38.0+rev1",
		OSVariant:                  "dev",
		SupervisorVersion:          "9.15.7",
		ManagedBySupervisorRelease: nil,
		IsManagedByServiceInstance: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/service_instance(123987)"}, ID: 123987,
		},
		ProvisioningProgress:       nil,
		ProvisioningState:          "",
		DownloadProgress:           nil,
		IsWebAccessible:            false,
		Longitude:                  "11.9672",
		Latitude:                   "57.7066",
		Location:                   "Gothenburg, Västra Götaland, Sweden",
		CustomLongitude:            "",
		CustomLatitude:             "",
		LogsChannel:                nil,
		IsLockedUntil:              nil,
		IsAccessibleBySupportUntil: nil,
		CreatedAt:                  "2019-09-23T11:55:58.654Z",
		IsActive:                   true,
	}
	// When
	actual, err := client.Device.Get(context.Background(), DeviceID(entityID))
	// Then
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestDeviceService_Get_UUID(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	uuid := "123456789123456789"
	mux.HandleFunc("/"+deviceBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=uuid+eq+%27123456789123456789%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		fmt.Fprint(w, deviceResponse)
	})
	expected := &DeviceResponse{
		ID: 1702297,
		BelongsToApplication: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/application(1234587)"}, ID: 1234587,
		},
		BelongsToUser: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/user(126547)"}, ID: 126547,
		},
		IsManagedByDevice:          nil,
		Actor:                      4081727,
		ShouldBeRunningRelease:     nil,
		Name:                       "old-moon",
		DeviceType:                 "qemux86-64",
		UUID:                       "00d859f123685e84772676f09465cc55",
		OnCommit:                   "0dfe02d12efb2325385952e1e15361f7",
		Note:                       nil,
		LocalID:                    nil,
		Status:                     "Idle",
		IsOnline:                   true,
		LastConnectivityEvent:      "2019-09-26T06:50:31.238Z",
		IsConnectedToVPN:           true,
		LastVPNEvent:               "2019-09-26T06:50:31.238Z",
		IPAddress:                  "10.0.2.15",
		VPNAddress:                 "10.240.63.4",
		PublicAddress:              "83.241.129.195",
		OSVersion:                  "balenaOS 2.38.0+rev1",
		OSVariant:                  "dev",
		SupervisorVersion:          "9.15.7",
		ManagedBySupervisorRelease: nil,
		IsManagedByServiceInstance: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/service_instance(123987)"}, ID: 123987,
		},
		ProvisioningProgress:       nil,
		ProvisioningState:          "",
		DownloadProgress:           nil,
		IsWebAccessible:            false,
		Longitude:                  "11.9672",
		Latitude:                   "57.7066",
		Location:                   "Gothenburg, Västra Götaland, Sweden",
		CustomLongitude:            "",
		CustomLatitude:             "",
		LogsChannel:                nil,
		IsLockedUntil:              nil,
		IsAccessibleBySupportUntil: nil,
		CreatedAt:                  "2019-09-23T11:55:58.654Z",
		IsActive:                   true,
	}
	// When
	actual, err := client.Device.Get(context.Background(), DeviceUUID(uuid))
	// Then
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestDeviceService_Get_NotFound(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	uuid := "123456789123456789"
	mux.HandleFunc("/"+deviceBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=uuid+eq+%27123456789123456789%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		fmt.Fprint(w, `{"d":[]}`)
	})
	// When
	device, err := client.Device.Get(context.Background(), DeviceUUID(uuid))
	// Then
	require.NoError(t, err)
	require.Nil(t, device)
}

func TestDeviceService_GetWithQuery(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=key+eq+%27value%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		fmt.Fprint(w, deviceResponse)
	})
	expected := []*DeviceResponse{
		{
			ID: 1702297,
			BelongsToApplication: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/application(1234587)"}, ID: 1234587,
			},
			BelongsToUser: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/user(126547)"}, ID: 126547,
			},
			IsManagedByDevice:          nil,
			Actor:                      4081727,
			ShouldBeRunningRelease:     nil,
			Name:                       "old-moon",
			DeviceType:                 "qemux86-64",
			UUID:                       "00d859f123685e84772676f09465cc55",
			OnCommit:                   "0dfe02d12efb2325385952e1e15361f7",
			Note:                       nil,
			LocalID:                    nil,
			Status:                     "Idle",
			IsOnline:                   true,
			LastConnectivityEvent:      "2019-09-26T06:50:31.238Z",
			IsConnectedToVPN:           true,
			LastVPNEvent:               "2019-09-26T06:50:31.238Z",
			IPAddress:                  "10.0.2.15",
			VPNAddress:                 "10.240.63.4",
			PublicAddress:              "83.241.129.195",
			OSVersion:                  "balenaOS 2.38.0+rev1",
			OSVariant:                  "dev",
			SupervisorVersion:          "9.15.7",
			ManagedBySupervisorRelease: nil,
			IsManagedByServiceInstance: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/service_instance(123987)"}, ID: 123987,
			},
			ProvisioningProgress:       nil,
			ProvisioningState:          "",
			DownloadProgress:           nil,
			IsWebAccessible:            false,
			Longitude:                  "11.9672",
			Latitude:                   "57.7066",
			Location:                   "Gothenburg, Västra Götaland, Sweden",
			CustomLongitude:            "",
			CustomLatitude:             "",
			LogsChannel:                nil,
			IsLockedUntil:              nil,
			IsAccessibleBySupportUntil: nil,
			CreatedAt:                  "2019-09-23T11:55:58.654Z",
			IsActive:                   true,
		},
	}
	// When
	actual, err := client.Device.GetWithQuery(context.Background(), "%24filter=key+eq+%27value%27")
	// Then
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestDeviceService_PinRelease_ID(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityID := int64(112233)
	releaseID := int64(14332)
	mux.HandleFunc(
		"/"+odata.EntityURL(deviceBasePath, strconv.FormatInt(entityID, 10)),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPatch)
			b, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, `{"should_be_running__release":"14332"}`+"\n", string(b))
			fmt.Fprint(w, "OK")
		},
	)
	// When
	resp, err := client.Device.PinRelease(context.Background(), DeviceID(entityID), releaseID)
	// Then
	require.NoError(t, err)
	require.Equal(t, "OK", string(resp))
}

func TestDeviceService_PinRelease_UUID(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	uuid := "123456789123456789"
	releaseID := int64(14332)
	mux.HandleFunc(
		"/"+deviceBasePath,
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPatch)
			expected := "%24filter=uuid+eq+%27" + uuid + "%27"
			if r.URL.RawQuery != expected {
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
				return
			}
			b, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, `{"should_be_running__release":"14332"}`+"\n", string(b))
			fmt.Fprint(w, "OK")
		},
	)
	// When
	resp, err := client.Device.PinRelease(context.Background(), DeviceUUID(uuid), releaseID)
	// Then
	require.NoError(t, err)
	require.Equal(t, "OK", string(resp))
}

func TestDeviceService_TrackLatestRelease_ID(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityID := int64(112233)
	mux.HandleFunc(
		"/"+odata.EntityURL(deviceBasePath, strconv.FormatInt(entityID, 10)),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPatch)
			b, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, `{"should_be_running__release":null}`+"\n", string(b))
			fmt.Fprint(w, "OK")
		},
	)
	// When
	resp, err := client.Device.TrackLatestRelease(context.Background(), DeviceID(entityID))
	// Then
	require.NoError(t, err)
	require.Equal(t, "OK", string(resp))
}

func TestDeviceService_TrackLatestRelease_UUID(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	uuid := "123456789123456789"
	mux.HandleFunc(
		"/"+deviceBasePath,
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPatch)
			expected := "%24filter=uuid+eq+%27" + uuid + "%27"
			if r.URL.RawQuery != expected {
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
				return
			}
			b, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, `{"should_be_running__release":null}`+"\n", string(b))
			fmt.Fprint(w, "OK")
		},
	)
	// When
	resp, err := client.Device.TrackLatestRelease(context.Background(), DeviceUUID(uuid))
	// Then
	require.NoError(t, err)
	require.Equal(t, "OK", string(resp))
}
