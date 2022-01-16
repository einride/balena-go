package balena

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestSupervisorV2Service_RestartServiceByName_Cloud(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	appID := int64(1514287)
	deviceUUID := "00d859f123685e84772676f09465cc55"
	mux.HandleFunc(
		"/supervisor/v2/applications/1514287/restart-service",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := ioutil.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(
				t,
				`{"uuid":"00d859f123685e84772676f09465cc55","method":"POST","data":{"serviceName":"testsvc"}}`+"\n",
				string(b),
			)
			fmt.Fprint(w, "OK")
		},
	)
	// When
	err := client.SupervisorV2(appID, deviceUUID).RestartServiceByName(context.Background(), "testsvc")
	// Then
	assert.NilError(t, err)
}

func TestSupervisorV2Service_RestartServiceByName_Local(t *testing.T) {
	// Given
	assert.NilError(t, os.Setenv("BALENA_SUPERVISOR_API_KEY", "test"))
	assert.NilError(t, os.Setenv("BALENA_APP_ID", "1122334"))
	assert.NilError(t, os.Setenv("BALENA_DEVICE_UUID", "11223344556677"))
	client, mux := supervisorV2Fixture(t)
	mux.HandleFunc(
		"/v2/applications/1122334/restart-service",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := ioutil.ReadAll(r.Body)
			assert.NilError(t, err)
			expected := "apikey=test"
			if r.URL.RawQuery != expected {
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
				return
			}
			assert.Equal(t, `{"serviceName":"testsvc"}`+"\n", string(b))
			fmt.Fprint(w, "OK")
		},
	)
	// When
	err := client.RestartServiceByName(context.Background(), "testsvc")
	// Then
	assert.NilError(t, err)
}

func TestSupervisorV2Service_StopServiceByName_Cloud(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	appID := int64(1514287)
	deviceUUID := "00d859f123685e84772676f09465cc55"
	mux.HandleFunc(
		"/supervisor/v2/applications/1514287/stop-service",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := ioutil.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(
				t,
				`{"uuid":"00d859f123685e84772676f09465cc55","method":"POST","data":{"serviceName":"testsvc"}}`+"\n",
				string(b),
			)
			fmt.Fprint(w, "OK")
		},
	)
	// When
	err := client.SupervisorV2(appID, deviceUUID).StopServiceByName(context.Background(), "testsvc")
	// Then
	assert.NilError(t, err)
}

func TestSupervisorV2Service_StopServiceByName_Local(t *testing.T) {
	// Given
	assert.NilError(t, os.Setenv("BALENA_SUPERVISOR_API_KEY", "test"))
	assert.NilError(t, os.Setenv("BALENA_APP_ID", "1122334"))
	assert.NilError(t, os.Setenv("BALENA_DEVICE_UUID", "11223344556677"))
	client, mux := supervisorV2Fixture(t)
	mux.HandleFunc(
		"/v2/applications/1122334/stop-service",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := ioutil.ReadAll(r.Body)
			assert.NilError(t, err)
			expected := "apikey=test"
			if r.URL.RawQuery != expected {
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
				return
			}
			assert.Equal(t, `{"serviceName":"testsvc"}`+"\n", string(b))
			fmt.Fprint(w, "OK")
		},
	)
	// When
	err := client.StopServiceByName(context.Background(), "testsvc")
	// Then
	assert.NilError(t, err)
}

func TestSupervisorV2Service_ApplicationState_Cloud(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	appID := int64(1514287)
	deviceUUID := "00d859f123685e84772676f09465cc55"
	services := map[string]ServiceState{
		"2233445": {
			Status:    "Running",
			ReleaseID: 12345,
		},
	}
	expected := &SvAppStateResp{
		Local: map[string]ApplicationState{
			"1122334": {
				Services: services,
			},
		},
		Commit:    "83b49b5eb012bdf0908dac8b3491b6f9",
		Dependent: map[string]interface{}{},
	}
	jsonResp := `{
	"local": {
		"1122334": {
			"services": {
				"2233445": {
					"status": "Running",
					"releaseId": 12345,
					"download_progress": null
				}
			}
		}
	},
	"dependent": {},
	"commit": "83b49b5eb012bdf0908dac8b3491b6f9"
}`
	mux.HandleFunc(
		"/supervisor/v2/applications/1514287/state",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := ioutil.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(t, `{"uuid":"00d859f123685e84772676f09465cc55","method":"GET"}`+"\n", string(b))
			fmt.Fprint(w, jsonResp)
		},
	)
	// When
	actual, err := client.SupervisorV2(appID, deviceUUID).ApplicationState(context.Background())
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestSupervisorV2Service_ApplicationState_Local(t *testing.T) {
	// Given
	assert.NilError(t, os.Setenv("BALENA_SUPERVISOR_API_KEY", "test"))
	assert.NilError(t, os.Setenv("BALENA_APP_ID", "1122334"))
	assert.NilError(t, os.Setenv("BALENA_DEVICE_UUID", "11223344556677"))
	client, mux := supervisorV2Fixture(t)
	services := map[string]ServiceState{
		"2233445": {
			Status:    "Running",
			ReleaseID: 12345,
		},
	}
	expected := &SvAppStateResp{
		Local: map[string]ApplicationState{
			"1122334": {
				Services: services,
			},
		},
		Commit:    "83b49b5eb012bdf0908dac8b3491b6f9",
		Dependent: map[string]interface{}{},
	}
	jsonResp := `{
	"local": {
		"1122334": {
			"services": {
				"2233445": {
					"status": "Running",
					"releaseId": 12345,
					"download_progress": null
				}
			}
		}
	},
	"dependent": {},
	"commit": "83b49b5eb012bdf0908dac8b3491b6f9"
}`
	mux.HandleFunc(
		"/v2/applications/1122334/state",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			b, err := ioutil.ReadAll(r.Body)
			assert.NilError(t, err)
			expected := "apikey=test"
			if r.URL.RawQuery != expected {
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
				return
			}
			assert.Equal(t, "", string(b))
			fmt.Fprint(w, jsonResp)
		},
	)
	// When
	actual, err := client.ApplicationState(context.Background())
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}
