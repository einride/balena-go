package balena

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestSupervisorV1Service_Reboot_Cloud(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	appID := int64(1514287)
	deviceUUID := "00d859f123685e84772676f09465cc55"
	mux.HandleFunc(
		"/supervisor/v1/reboot",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(
				t,
				`{"uuid":"00d859f123685e84772676f09465cc55","method":"POST","data":{"force":true}}`+"\n",
				string(b),
			)
			fmt.Fprint(w, `{"Data":"OK","Error":""}`)
		},
	)
	// When
	err := client.SupervisorV1(appID, deviceUUID).Reboot(context.Background(), true)
	// Then
	assert.NilError(t, err)
}

func TestSupervisorV1Service_Reboot_CloudError(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	appID := int64(1514287)
	deviceUUID := "00d859f123685e84772676f09465cc55"
	mux.HandleFunc(
		"/supervisor/v1/reboot",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(
				t,
				`{"uuid":"00d859f123685e84772676f09465cc55","method":"POST","data":{"force":true}}`+"\n",
				string(b),
			)
			fmt.Fprint(w, `{"Data":"Not OK","Error":"Something was bad"}`)
		},
	)
	// When
	err := client.SupervisorV1(appID, deviceUUID).Reboot(context.Background(), true)
	// Then
	assert.Assert(t, err != nil)
	assert.ErrorContains(t, err, "Something was bad")
}

func TestSupervisorV1Service_Reboot_Local(t *testing.T) {
	// Given
	assert.NilError(t, os.Setenv("BALENA_SUPERVISOR_API_KEY", "test"))
	assert.NilError(t, os.Setenv("BALENA_APP_ID", "1122334"))
	assert.NilError(t, os.Setenv("BALENA_DEVICE_UUID", "11223344556677"))
	client, mux := supervisorV1Fixture(t)
	mux.HandleFunc(
		"/v1/reboot",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			expected := "apikey=test"
			if r.URL.RawQuery != expected {
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
				return
			}
			assert.Equal(
				t,
				`{"force":true}`+"\n",
				string(b),
			)
			fmt.Fprint(w, `{"Data":"OK","Error":""}`)
		},
	)
	// When
	err := client.Reboot(context.Background(), true)
	// Then
	assert.NilError(t, err)
}

func TestSupervisorV1Service_Blink_Cloud(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	appID := int64(1514287)
	deviceUUID := "00d859f123685e84772676f09465cc55"
	mux.HandleFunc(
		"/supervisor/v1/blink",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(
				t,
				`{"uuid":"00d859f123685e84772676f09465cc55","method":"POST"}`+"\n",
				string(b),
			)
			w.WriteHeader(http.StatusOK)
		},
	)
	// When
	err := client.SupervisorV1(appID, deviceUUID).Blink(context.Background())
	// Then
	assert.NilError(t, err)
}

func TestSupervisorV1Service_Blink_CloudError(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	appID := int64(1514287)
	deviceUUID := "00d859f123685e84772676f09465cc55"
	mux.HandleFunc(
		"/supervisor/v1/blink",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(
				t,
				`{"uuid":"00d859f123685e84772676f09465cc55","method":"POST"}`+"\n",
				string(b),
			)
			w.WriteHeader(http.StatusInternalServerError)
		},
	)
	// When
	err := client.SupervisorV1(appID, deviceUUID).Blink(context.Background())
	// Then
	assert.Assert(t, err != nil)

	var errorResponse *ErrorResponse
	assert.Check(t, errors.As(err, &errorResponse))
	if errorResponse != nil {
		assert.ErrorContains(t, errorResponse, "500")
	}
}

func TestSupervisorV1Service_Blink_Local(t *testing.T) {
	// Given
	assert.NilError(t, os.Setenv("BALENA_SUPERVISOR_API_KEY", "test"))
	assert.NilError(t, os.Setenv("BALENA_APP_ID", "1122334"))
	assert.NilError(t, os.Setenv("BALENA_DEVICE_UUID", "11223344556677"))
	client, mux := supervisorV1Fixture(t)
	mux.HandleFunc(
		"/v1/blink",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			expected := "apikey=test"
			if r.URL.RawQuery != expected {
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
				return
			}
			assert.Equal(t, "", string(b))
			w.WriteHeader(http.StatusOK)
		},
	)
	// When
	err := client.Blink(context.Background())
	// Then
	assert.NilError(t, err)
}

func TestSupervisorV1Service_Update_Cloud(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	appID := int64(1514287)
	deviceUUID := "00d859f123685e84772676f09465cc55"
	mux.HandleFunc(
		"/supervisor/v1/update",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(
				t,
				`{"uuid":"00d859f123685e84772676f09465cc55","method":"POST","data":{"force":true}}`+"\n",
				string(b),
			)
			w.WriteHeader(http.StatusNoContent)
		},
	)
	// When
	err := client.SupervisorV1(appID, deviceUUID).Update(context.Background(), true)
	// Then
	assert.NilError(t, err)
}

func TestSupervisorV1Service_Update_CloudError(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	appID := int64(1514287)
	deviceUUID := "00d859f123685e84772676f09465cc55"
	mux.HandleFunc(
		"/supervisor/v1/update",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(
				t,
				`{"uuid":"00d859f123685e84772676f09465cc55","method":"POST","data":{"force":true}}`+"\n",
				string(b),
			)
			w.WriteHeader(http.StatusInternalServerError)
		},
	)
	// When
	err := client.SupervisorV1(appID, deviceUUID).Update(context.Background(), true)
	// Then
	assert.Assert(t, err != nil)

	var errorResponse *ErrorResponse
	assert.Check(t, errors.As(err, &errorResponse))
	if errorResponse != nil {
		assert.ErrorContains(t, errorResponse, "500")
	}
}

func TestSupervisorV1Service_Update_Local(t *testing.T) {
	// Given
	assert.NilError(t, os.Setenv("BALENA_SUPERVISOR_API_KEY", "test"))
	assert.NilError(t, os.Setenv("BALENA_APP_ID", "1122334"))
	assert.NilError(t, os.Setenv("BALENA_DEVICE_UUID", "11223344556677"))
	client, mux := supervisorV1Fixture(t)
	mux.HandleFunc(
		"/v1/update",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			expected := "apikey=test"
			if r.URL.RawQuery != expected {
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
				return
			}
			assert.Equal(
				t,
				`{"force":true}`+"\n",
				string(b),
			)
			w.WriteHeader(http.StatusNoContent)
		},
	)
	// When
	err := client.Update(context.Background(), true)
	// Then
	assert.NilError(t, err)
}

const (
	supervisorV1DeviceResponse = `{
		"api_port":48484,
		"ip_address":"192.168.0.114 10.42.0.3",
		"commit":"414e65cd378a69a96f403b75f14b40b55856f860",
		"status":"Downloading",
		"download_progress":84,
		"os_version":"Resin OS 1.0.4 (fido)",
		"supervisor_version":"1.6.0",
		"update_pending":true,
		"update_downloaded":false,
		"update_failed":false
	}`
)

func TestSupervisorV1Service_Device_Cloud(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	appID := int64(1514287)
	deviceUUID := "00d859f123685e84772676f09465cc55"
	mux.HandleFunc(
		"/supervisor/v1/device",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(
				t,
				`{"uuid":"00d859f123685e84772676f09465cc55","method":"GET"}`+"\n",
				string(b),
			)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, supervisorV1DeviceResponse)
		},
	)
	// When
	response, err := client.SupervisorV1(appID, deviceUUID).Device(context.Background())
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, response, &SupervisorV1DeviceResponse{
		APIPort:           48484,
		Commit:            "414e65cd378a69a96f403b75f14b40b55856f860",
		IPAddress:         "192.168.0.114 10.42.0.3",
		Status:            "Downloading",
		DownloadProgress:  84,
		OSVersion:         "Resin OS 1.0.4 (fido)",
		SupervisorVersion: "1.6.0",
		UpdatePending:     true,
		UpdateDownloaded:  false,
		UpdateFailed:      false,
	})
}

func TestSupervisorV1Service_Device_CloudError(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	appID := int64(1514287)
	deviceUUID := "00d859f123685e84772676f09465cc55"
	mux.HandleFunc(
		"/supervisor/v1/device",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(
				t,
				`{"uuid":"00d859f123685e84772676f09465cc55","method":"GET"}`+"\n",
				string(b),
			)
			w.WriteHeader(http.StatusInternalServerError)
		},
	)
	// When
	response, err := client.SupervisorV1(appID, deviceUUID).Device(context.Background())
	// Then
	assert.Assert(t, err != nil)

	var errorResponse *ErrorResponse
	assert.Check(t, errors.As(err, &errorResponse))
	if errorResponse != nil {
		assert.ErrorContains(t, errorResponse, "500")
	}

	assert.Assert(t, response == nil)
}

func TestSupervisorV1Service_Device_Local(t *testing.T) {
	// Given
	assert.NilError(t, os.Setenv("BALENA_SUPERVISOR_API_KEY", "test"))
	assert.NilError(t, os.Setenv("BALENA_APP_ID", "1122334"))
	assert.NilError(t, os.Setenv("BALENA_DEVICE_UUID", "11223344556677"))
	client, mux := supervisorV1Fixture(t)
	mux.HandleFunc(
		"/v1/device",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			b, err := io.ReadAll(r.Body)
			assert.NilError(t, err)
			expected := "apikey=test"
			if r.URL.RawQuery != expected {
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
				return
			}
			assert.Equal(t, "", string(b))
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, supervisorV1DeviceResponse)
		},
	)
	// When
	response, err := client.Device(context.Background())
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, response, &SupervisorV1DeviceResponse{
		APIPort:           48484,
		Commit:            "414e65cd378a69a96f403b75f14b40b55856f860",
		IPAddress:         "192.168.0.114 10.42.0.3",
		Status:            "Downloading",
		DownloadProgress:  84,
		OSVersion:         "Resin OS 1.0.4 (fido)",
		SupervisorVersion: "1.6.0",
		UpdatePending:     true,
		UpdateDownloaded:  false,
		UpdateFailed:      false,
	})
}
