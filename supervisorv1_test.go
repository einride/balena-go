package balena

import (
	"context"
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
