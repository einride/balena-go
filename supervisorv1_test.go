package balena

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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
			b, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t,
				`{"uuid":"00d859f123685e84772676f09465cc55","method":"POST","data":{"force":true}}`+"\n",
				string(b),
			)
			fmt.Fprint(w, `{"Data":"OK","Error":""}`)
		},
	)
	// When
	err := client.SupervisorV1(appID, deviceUUID).Reboot(context.Background(), true)
	// Then
	require.NoError(t, err)
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
			b, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t,
				`{"uuid":"00d859f123685e84772676f09465cc55","method":"POST","data":{"force":true}}`+"\n",
				string(b),
			)
			fmt.Fprint(w, `{"Data":"Not OK","Error":"Something was bad"}`)
		},
	)
	// When
	err := client.SupervisorV1(appID, deviceUUID).Reboot(context.Background(), true)
	// Then
	require.Error(t, err)
	require.Contains(t, err.Error(), "Something was bad")
}

func TestSupervisorV1Service_Reboot_Local(t *testing.T) {
	// Given
	require.NoError(t, os.Setenv("BALENA_SUPERVISOR_API_KEY", "test"))
	require.NoError(t, os.Setenv("BALENA_APP_ID", "1122334"))
	require.NoError(t, os.Setenv("BALENA_DEVICE_UUID", "11223344556677"))
	client, mux := supervisorV1Fixture(t)
	mux.HandleFunc(
		"/v1/reboot",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			b, err := ioutil.ReadAll(r.Body)
			require.NoError(t, err)
			expected := "apikey=test"
			if r.URL.RawQuery != expected {
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
				return
			}
			require.Equal(t,
				`{"force":true}`+"\n",
				string(b),
			)
			fmt.Fprint(w, `{"Data":"OK","Error":""}`)
		},
	)
	// When
	err := client.Reboot(context.Background(), true)
	// Then
	require.NoError(t, err)
}
