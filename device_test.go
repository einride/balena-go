package balena

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"go.einride.tech/balena/odata"
	"gotest.tools/v3/assert"
)

const (
	deviceResponse = `{
	"d": [
		{
   		   "id": 4218895,
   		   "belongs_to__application": {
   		     "__id": 1827427,
   		     "__deferred": {
   		       "uri": "/resin/application(@id)?@id=1827427"
   		     }
   		   },
   		   "belongs_to__user": null,
   		   "is_managed_by__device": null,
   		   "actor": 7288314,
   		   "should_be_running__release": {
   		     "__id": 1796078,
   		     "__deferred": {
   		       "uri": "/resin/release(@id)?@id=1796078"
   		     }
   		   },
   		   "device_name": "log-station-office",
		   "device_type": {
   		     "__id": 58,
   		     "__deferred": {
   		       "uri": "/resin/device_type(@id)?@id=58"
   		     }
   		   },
   		   "uuid": "6fe2836d9bbebc5b399f5fc28b840e8e",
   		   "is_running__release": {
   		     "__id": 1796078,
   		     "__deferred": {
   		       "uri": "/resin/release(@id)?@id=1796078"
   		     }
   		   },
   		   "note": null,
   		   "local_id": null,
   		   "status": "idle",
		   "overall_status": "idle",
   		   "is_online": true,
   		   "last_connectivity_event": "2021-05-23T04:13:21.629Z",
   		   "is_connected_to_vpn": true,
   		   "last_vpn_event": "2021-05-23T04:13:21.629Z",
   		   "ip_address": "10.1.20.198",
   		   "mac_address": "b8:27:eb:72:f9:5e b8:40:eb:27:ac:0b",
   		   "vpn_address": "12.345.95.246",
   		   "public_address": "12.345.41.74",
		   "os_version": "balenaOS 2.75.0+rev1",
		   "os_variant": "prod",
		   "supervisor_version": "12.5.10",
		   "should_be_managed_by__supervisor_release": {
			   "__id": 1764685,
			   "__deferred": {
				   "uri": "/resin/supervisor_release(@id)?@id=1764685"
			   }
		   },
		   "should_be_operated_by__release": {
			   "__id": 1781572,
			   "__deferred": {
				   "uri": "/resin/release(@id)?@id=1781572"
			   }
		   },
		   "is_managed_by__service_instance": {
			   "__id": 124474,
			   "__deferred": {
				   "uri": "/resin/service_instance(@id)?@id=124474"
			   }
		   },
		   "provisioning_progress": null,
		   "provisioning_state": "",
		   "download_progress": null,
		   "is_web_accessible": false,
		   "longitude": "12.2103",
		   "latitude": "57.6828",
		   "location": "Landvetter, Västra Götaland County, Sweden",
		   "custom_longitude": "",
		   "custom_latitude": "",
		   "logs_channel": null,
		   "is_locked_until__date": null,
		   "is_accessible_by_support_until__date": null,
		   "created_at": "2021-05-11T08:05:16.634Z",
		   "is_active": true,
		   "api_heartbeat_state": "online",
		   "memory_usage": 321,
		   "memory_total": 973,
		   "storage_block_device": "/dev/mmcblk0p6",
		   "storage_usage": 191,
		   "storage_total": 14138,
		   "cpu_temp": 63,
		   "cpu_usage": 34,
		   "cpu_id": "000000008e72f95e",
		   "is_undervolted": false,
		   "__metadata": {
			   "uri": "/resin/device(@id)?@id=4218895"
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
			ID:         4218895,
			Actor:      7288314,
			DeviceName: "log-station-office",
			DeviceType: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/device_type(@id)?@id=58"}, ID: 58,
			},
			UUID:                  "6fe2836d9bbebc5b399f5fc28b840e8e",
			LastConnectivityEvent: "2021-05-23T04:13:21.629Z",
			Status:                "idle",
			OverallStatus:         "idle",
			LastVPNEvent:          "2021-05-23T04:13:21.629Z",
			IPAddress:             "10.1.20.198",
			VPNAddress:            "12.345.95.246",
			PublicAddress:         "12.345.41.74",
			OSVersion:             "balenaOS 2.75.0+rev1",
			OSVariant:             "prod",
			SupervisorVersion:     "12.5.10",
			ProvisioningState:     "",
			Longitude:             "12.2103",
			Latitude:              "57.6828",
			Location:              "Landvetter, Västra Götaland County, Sweden",
			CustomLongitude:       "",
			CustomLatitude:        "",
			CreatedAt:             "2021-05-11T08:05:16.634Z",
			IsOnline:              true,
			IsConnectedToVPN:      true,
			IsWebAccessible:       false,
			IsActive:              true,
			BelongsToApplication: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/application(@id)?@id=1827427"}, ID: 1827427,
			},
			BelongsToUser: nil,
			IsManagedByServiceInstance: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/service_instance(@id)?@id=124474"}, ID: 124474,
			},
			IsManagedByDevice: nil,
			IsRunningRelease: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1796078"}, ID: 1796078,
			},
			ShouldBeRunningRelease: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1796078"}, ID: 1796078,
			},
			Note:    nil,
			LocalID: nil,
			ShouldBeManagedBySupervisorRelease: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/supervisor_release(@id)?@id=1764685"}, ID: 1764685,
			},
			ProvisioningProgress:       nil,
			DownloadProgress:           nil,
			LogsChannel:                nil,
			IsLockedUntil:              nil,
			IsAccessibleBySupportUntil: nil,
			MACAddress:                 "b8:27:eb:72:f9:5e b8:40:eb:27:ac:0b",
			APIHeartbeatState:          "online",
			MemoryUsage:                321,
			MemoryTotal:                973,
			StorageBlockDevice:         "/dev/mmcblk0p6",
			StorageUsage:               191,
			StorageTotal:               14138,
			CPUTemp:                    63,
			CPUUsage:                   34,
			CPUID:                      "000000008e72f95e",
			IsUndervolted:              false,
		},
	}
	// When
	actual, err := client.Device.List(context.Background())
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
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
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, deviceResponse)
	})
	expected := []*DeviceResponse{
		{
			ID:         4218895,
			Actor:      7288314,
			DeviceName: "log-station-office",
			DeviceType: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/device_type(@id)?@id=58"}, ID: 58,
			},
			UUID:                  "6fe2836d9bbebc5b399f5fc28b840e8e",
			LastConnectivityEvent: "2021-05-23T04:13:21.629Z",
			Status:                "idle",
			OverallStatus:         "idle",
			LastVPNEvent:          "2021-05-23T04:13:21.629Z",
			IPAddress:             "10.1.20.198",
			VPNAddress:            "12.345.95.246",
			PublicAddress:         "12.345.41.74",
			OSVersion:             "balenaOS 2.75.0+rev1",
			OSVariant:             "prod",
			SupervisorVersion:     "12.5.10",
			ProvisioningState:     "",
			Longitude:             "12.2103",
			Latitude:              "57.6828",
			Location:              "Landvetter, Västra Götaland County, Sweden",
			CustomLongitude:       "",
			CustomLatitude:        "",
			CreatedAt:             "2021-05-11T08:05:16.634Z",
			IsOnline:              true,
			IsConnectedToVPN:      true,
			IsWebAccessible:       false,
			IsActive:              true,
			BelongsToApplication: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/application(@id)?@id=1827427"}, ID: 1827427,
			},
			BelongsToUser: nil,
			IsManagedByServiceInstance: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/service_instance(@id)?@id=124474"}, ID: 124474,
			},
			IsManagedByDevice: nil,
			IsRunningRelease: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1796078"}, ID: 1796078,
			},
			ShouldBeRunningRelease: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1796078"}, ID: 1796078,
			},
			Note:    nil,
			LocalID: nil,
			ShouldBeManagedBySupervisorRelease: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/supervisor_release(@id)?@id=1764685"}, ID: 1764685,
			},
			ProvisioningProgress:       nil,
			DownloadProgress:           nil,
			LogsChannel:                nil,
			IsLockedUntil:              nil,
			IsAccessibleBySupportUntil: nil,
			MACAddress:                 "b8:27:eb:72:f9:5e b8:40:eb:27:ac:0b",
			APIHeartbeatState:          "online",
			MemoryUsage:                321,
			MemoryTotal:                973,
			StorageBlockDevice:         "/dev/mmcblk0p6",
			StorageUsage:               191,
			StorageTotal:               14138,
			CPUTemp:                    63,
			CPUUsage:                   34,
			CPUID:                      "000000008e72f95e",
			IsUndervolted:              false,
		},
	}
	// When
	actual, err := client.Device.ListByApplication(context.Background(), applicationID)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
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
		ID:         4218895,
		Actor:      7288314,
		DeviceName: "log-station-office",
		DeviceType: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/device_type(@id)?@id=58"}, ID: 58,
		},
		UUID:                  "6fe2836d9bbebc5b399f5fc28b840e8e",
		LastConnectivityEvent: "2021-05-23T04:13:21.629Z",
		Status:                "idle",
		OverallStatus:         "idle",
		LastVPNEvent:          "2021-05-23T04:13:21.629Z",
		IPAddress:             "10.1.20.198",
		VPNAddress:            "12.345.95.246",
		PublicAddress:         "12.345.41.74",
		OSVersion:             "balenaOS 2.75.0+rev1",
		OSVariant:             "prod",
		SupervisorVersion:     "12.5.10",
		ProvisioningState:     "",
		Longitude:             "12.2103",
		Latitude:              "57.6828",
		Location:              "Landvetter, Västra Götaland County, Sweden",
		CustomLongitude:       "",
		CustomLatitude:        "",
		CreatedAt:             "2021-05-11T08:05:16.634Z",
		IsOnline:              true,
		IsConnectedToVPN:      true,
		IsWebAccessible:       false,
		IsActive:              true,
		BelongsToApplication: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/application(@id)?@id=1827427"}, ID: 1827427,
		},
		BelongsToUser: nil,
		IsManagedByServiceInstance: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/service_instance(@id)?@id=124474"}, ID: 124474,
		},
		IsManagedByDevice: nil,
		IsRunningRelease: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1796078"}, ID: 1796078,
		},
		ShouldBeRunningRelease: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1796078"}, ID: 1796078,
		},
		Note:    nil,
		LocalID: nil,
		ShouldBeManagedBySupervisorRelease: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/supervisor_release(@id)?@id=1764685"}, ID: 1764685,
		},
		ProvisioningProgress:       nil,
		DownloadProgress:           nil,
		LogsChannel:                nil,
		IsLockedUntil:              nil,
		IsAccessibleBySupportUntil: nil,
		MACAddress:                 "b8:27:eb:72:f9:5e b8:40:eb:27:ac:0b",
		APIHeartbeatState:          "online",
		MemoryUsage:                321,
		MemoryTotal:                973,
		StorageBlockDevice:         "/dev/mmcblk0p6",
		StorageUsage:               191,
		StorageTotal:               14138,
		CPUTemp:                    63,
		CPUUsage:                   34,
		CPUID:                      "000000008e72f95e",
		IsUndervolted:              false,
	}
	// When
	actual, err := client.Device.Get(context.Background(), DeviceID(entityID))
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
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
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, deviceResponse)
	})
	expected := &DeviceResponse{
		ID:         4218895,
		Actor:      7288314,
		DeviceName: "log-station-office",
		DeviceType: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/device_type(@id)?@id=58"}, ID: 58,
		},
		UUID:                  "6fe2836d9bbebc5b399f5fc28b840e8e",
		LastConnectivityEvent: "2021-05-23T04:13:21.629Z",
		Status:                "idle",
		OverallStatus:         "idle",
		LastVPNEvent:          "2021-05-23T04:13:21.629Z",
		IPAddress:             "10.1.20.198",
		VPNAddress:            "12.345.95.246",
		PublicAddress:         "12.345.41.74",
		OSVersion:             "balenaOS 2.75.0+rev1",
		OSVariant:             "prod",
		SupervisorVersion:     "12.5.10",
		ProvisioningState:     "",
		Longitude:             "12.2103",
		Latitude:              "57.6828",
		Location:              "Landvetter, Västra Götaland County, Sweden",
		CustomLongitude:       "",
		CustomLatitude:        "",
		CreatedAt:             "2021-05-11T08:05:16.634Z",
		IsOnline:              true,
		IsConnectedToVPN:      true,
		IsWebAccessible:       false,
		IsActive:              true,
		BelongsToApplication: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/application(@id)?@id=1827427"}, ID: 1827427,
		},
		BelongsToUser: nil,
		IsManagedByServiceInstance: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/service_instance(@id)?@id=124474"}, ID: 124474,
		},
		IsManagedByDevice: nil,
		IsRunningRelease: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1796078"}, ID: 1796078,
		},
		ShouldBeRunningRelease: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1796078"}, ID: 1796078,
		},
		Note:    nil,
		LocalID: nil,
		ShouldBeManagedBySupervisorRelease: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/supervisor_release(@id)?@id=1764685"}, ID: 1764685,
		},
		ProvisioningProgress:       nil,
		DownloadProgress:           nil,
		LogsChannel:                nil,
		IsLockedUntil:              nil,
		IsAccessibleBySupportUntil: nil,
		MACAddress:                 "b8:27:eb:72:f9:5e b8:40:eb:27:ac:0b",
		APIHeartbeatState:          "online",
		MemoryUsage:                321,
		MemoryTotal:                973,
		StorageBlockDevice:         "/dev/mmcblk0p6",
		StorageUsage:               191,
		StorageTotal:               14138,
		CPUTemp:                    63,
		CPUUsage:                   34,
		CPUID:                      "000000008e72f95e",
		IsUndervolted:              false,
	}
	// When
	actual, err := client.Device.Get(context.Background(), DeviceUUID(uuid))
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
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
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, `{"d":[]}`)
	})
	// When
	device, err := client.Device.Get(context.Background(), DeviceUUID(uuid))
	// Then
	assert.NilError(t, err)
	assert.Assert(t, device == nil)
}

func TestDeviceService_GetWithQuery(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=key+eq+%27value%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, deviceResponse)
	})
	expected := []*DeviceResponse{
		{
			ID:         4218895,
			Actor:      7288314,
			DeviceName: "log-station-office",
			DeviceType: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/device_type(@id)?@id=58"}, ID: 58,
			},
			UUID:                  "6fe2836d9bbebc5b399f5fc28b840e8e",
			LastConnectivityEvent: "2021-05-23T04:13:21.629Z",
			Status:                "idle",
			OverallStatus:         "idle",
			LastVPNEvent:          "2021-05-23T04:13:21.629Z",
			IPAddress:             "10.1.20.198",
			VPNAddress:            "12.345.95.246",
			PublicAddress:         "12.345.41.74",
			OSVersion:             "balenaOS 2.75.0+rev1",
			OSVariant:             "prod",
			SupervisorVersion:     "12.5.10",
			ProvisioningState:     "",
			Longitude:             "12.2103",
			Latitude:              "57.6828",
			Location:              "Landvetter, Västra Götaland County, Sweden",
			CustomLongitude:       "",
			CustomLatitude:        "",
			CreatedAt:             "2021-05-11T08:05:16.634Z",
			IsOnline:              true,
			IsConnectedToVPN:      true,
			IsWebAccessible:       false,
			IsActive:              true,
			BelongsToApplication: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/application(@id)?@id=1827427"}, ID: 1827427,
			},
			BelongsToUser: nil,
			IsManagedByServiceInstance: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/service_instance(@id)?@id=124474"}, ID: 124474,
			},
			IsManagedByDevice: nil,
			IsRunningRelease: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1796078"}, ID: 1796078,
			},
			ShouldBeRunningRelease: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1796078"}, ID: 1796078,
			},
			Note:    nil,
			LocalID: nil,
			ShouldBeManagedBySupervisorRelease: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/supervisor_release(@id)?@id=1764685"}, ID: 1764685,
			},
			ProvisioningProgress:       nil,
			DownloadProgress:           nil,
			LogsChannel:                nil,
			IsLockedUntil:              nil,
			IsAccessibleBySupportUntil: nil,
			MACAddress:                 "b8:27:eb:72:f9:5e b8:40:eb:27:ac:0b",
			APIHeartbeatState:          "online",
			MemoryUsage:                321,
			MemoryTotal:                973,
			StorageBlockDevice:         "/dev/mmcblk0p6",
			StorageUsage:               191,
			StorageTotal:               14138,
			CPUTemp:                    63,
			CPUUsage:                   34,
			CPUID:                      "000000008e72f95e",
			IsUndervolted:              false,
		},
	}
	// When
	actual, err := client.Device.GetWithQuery(context.Background(), "%24filter=key+eq+%27value%27")
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
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
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(t, `{"should_be_running__release":"14332"}`+"\n", string(b))
			fmt.Fprint(w, "OK")
		},
	)
	// When
	resp, err := client.Device.PinRelease(context.Background(), DeviceID(entityID), releaseID)
	// Then
	assert.NilError(t, err)
	assert.Equal(t, "OK", string(resp))
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
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
				return
			}
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(t, `{"should_be_running__release":"14332"}`+"\n", string(b))
			fmt.Fprint(w, "OK")
		},
	)
	// When
	resp, err := client.Device.PinRelease(context.Background(), DeviceUUID(uuid), releaseID)
	// Then
	assert.NilError(t, err)
	assert.Equal(t, "OK", string(resp))
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
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(t, `{"should_be_running__release":null}`+"\n", string(b))
			fmt.Fprint(w, "OK")
		},
	)
	// When
	resp, err := client.Device.TrackLatestRelease(context.Background(), DeviceID(entityID))
	// Then
	assert.NilError(t, err)
	assert.Equal(t, "OK", string(resp))
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
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
				return
			}
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(t, `{"should_be_running__release":null}`+"\n", string(b))
			fmt.Fprint(w, "OK")
		},
	)
	// When
	resp, err := client.Device.TrackLatestRelease(context.Background(), DeviceUUID(uuid))
	// Then
	assert.NilError(t, err)
	assert.Equal(t, "OK", string(resp))
}
