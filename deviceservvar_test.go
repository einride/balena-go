package balena

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"go.einride.tech/balena/odata"
	"gotest.tools/v3/assert"
)

func TestDeviceServVarService_List_ID(t *testing.T) {
	// Given
	jsonResp := `{
	"d": [
		{
			"id": 183330,
			"created_at": "2019-09-27T17:54:21.559Z",
			"device": {
				"__deferred": {
					"uri": "/resin/device(1702297)"
				},
				"__id": 1702297
			},
			"name": "key",
			"value": "test",
			"__metadata": {
				"uri": "/resin/device_service_environment_variable(@id)?@id=183330"
			}
		}
	]
}`
	deviceID := int64(123456)
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceServVarBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=service_install/device+eq+%27123456%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		fmt.Fprint(w, jsonResp)
	})
	expected := []*DeviceServVarResponse{
		{
			ID:        183330,
			CreatedAt: "2019-09-27T17:54:21.559Z",
			Device: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/device(1702297)"},
				ID:       1702297,
			},
			Name:  "key",
			Value: "test",
		},
	}
	// When
	actual, err := client.DeviceServVar.List(context.Background(), DeviceID(deviceID))
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestDeviceServVarService_List_UUID(t *testing.T) {
	// Given
	jsonResp := `{
	"d": [
		{
			"id": 183330,
			"created_at": "2019-09-27T17:54:21.559Z",
			"device": {
				"__deferred": {
					"uri": "/resin/device(1702297)"
				},
				"__id": 1702297
			},
			"name": "key",
			"value": "test",
			"__metadata": {
				"uri": "/resin/device_service_environment_variable(@id)?@id=183330"
			}
		}
	]
}`
	uuid := "123456789123456798"
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceServVarBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=service_install/device/uuid+eq+%27" + uuid + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		fmt.Fprint(w, jsonResp)
	})
	expected := []*DeviceServVarResponse{
		{
			ID:        183330,
			CreatedAt: "2019-09-27T17:54:21.559Z",
			Device: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/device(1702297)"},
				ID:       1702297,
			},
			Name:  "key",
			Value: "test",
		},
	}
	// When
	actual, err := client.DeviceServVar.List(context.Background(), DeviceUUID(uuid))
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestDeviceServVarService_Create_ID(t *testing.T) {
	// Given
	jsonResp := `{
	"id": 183330,
	"created_at": "2019-09-27T17:54:21.559Z",
	"device": {
		"__deferred": {
			"uri": "/resin/device(1702297)"
		},
		"__id": 1702297
	},
	"name": "key",
	"value": "test",
	"__metadata": {
		"uri": "/resin/device_service_environment_variable(@id)?@id=183330"
	}
}`
	deviceID := int64(123456)
	key := "key"
	value := "value"
	client, mux, cleanup := newFixture()
	defer cleanup()
	type request struct {
		Device string `json:"device"`
		Name   string `json:"name"`
		Value  string `json:"value"`
	}
	mux.HandleFunc("/"+deviceServVarBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		b, err := ioutil.ReadAll(r.Body)
		assert.NilError(t, err)
		req := &request{}
		assert.NilError(t, json.Unmarshal(b, req))
		assert.Equal(t, strconv.FormatInt(deviceID, 10), req.Device)
		assert.Equal(t, key, req.Name)
		assert.Equal(t, value, req.Value)
		fmt.Fprint(w, jsonResp)
	})
	expected := &DeviceServVarResponse{
		ID:        183330,
		CreatedAt: "2019-09-27T17:54:21.559Z",
		Device: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/device(1702297)"},
			ID:       1702297,
		},
		Name:  "key",
		Value: "test",
	}
	// When
	actual, err := client.DeviceServVar.Create(context.Background(), DeviceID(deviceID), key, value)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestDeviceServVarService_Create_UUID(t *testing.T) {
	// Given
	jsonResp := `{
	"id": 183330,
	"created_at": "2019-09-27T17:54:21.559Z",
	"device": {
		"__deferred": {
			"uri": "/resin/device(1702297)"
		},
		"__id": 1702297
	},
	"name": "key",
	"value": "test",
	"__metadata": {
		"uri": "/resin/device_service_environment_variable(@id)?@id=183330"
	}
}`
	uuid := "1234567890123456789"
	deviceID := int64(123456)
	key := "key"
	value := "value"
	client, mux, cleanup := newFixture()
	defer cleanup()
	type request struct {
		Device string `json:"device"`
		Name   string `json:"name"`
		Value  string `json:"value"`
	}
	mux.HandleFunc("/"+deviceBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=uuid+eq+%27" + uuid + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		resp := `{"d":[{"id":123456}]}`
		fmt.Fprint(w, resp)
	})
	mux.HandleFunc("/"+deviceServVarBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		b, err := ioutil.ReadAll(r.Body)
		assert.NilError(t, err)
		req := &request{}
		assert.NilError(t, json.Unmarshal(b, req))
		assert.Equal(t, strconv.FormatInt(deviceID, 10), req.Device)
		assert.Equal(t, key, req.Name)
		assert.Equal(t, value, req.Value)
		fmt.Fprint(w, jsonResp)
	})
	expected := &DeviceServVarResponse{
		ID:        183330,
		CreatedAt: "2019-09-27T17:54:21.559Z",
		Device: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/device(1702297)"},
			ID:       1702297,
		},
		Name:  "key",
		Value: "test",
	}
	// When
	actual, err := client.DeviceServVar.Create(context.Background(), DeviceUUID(uuid), key, value)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestDeviceServVarService_DeleteWithName_ID_OK(t *testing.T) {
	// Given
	deviceID := int64(123456)
	key := "key"
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceServVarBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		expected := "%24filter=service_install/device+eq+%27" +
			strconv.FormatInt(deviceID, 10) + "%27+and+name+eq+%27" + key + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		fmt.Fprint(w, `{
	"d": [
		{
			"id": 183330,
			"created_at": "2019-09-27T17:54:21.559Z",
			"device": {
				"__deferred": {
					"uri": "/resin/device(1702297)"
				},
				"__id": 1702297
			},
			"name": "key",
			"value": "test",
			"__metadata": {
				"uri": "/resin/device_service_environment_variable(@id)?@id=183330"
			}
		}
	]
}`)
	})
	// When
	err := client.DeviceServVar.DeleteWithName(context.Background(), DeviceID(deviceID), key)
	// Then
	assert.NilError(t, err)
}

func TestDeviceServVarService_DeleteWithName_UUID_OK(t *testing.T) {
	// Given
	uuid := "12345678901234567890"
	key := "key"
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceServVarBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		expected := "%24filter=service_install/device/uuid+eq+%27" + uuid + "%27+and+name+eq+%27" + key + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		fmt.Fprint(w, `{
	"d": [
		{
			"id": 183330,
			"created_at": "2019-09-27T17:54:21.559Z",
			"device": {
				"__deferred": {
					"uri": "/resin/device(1702297)"
				},
				"__id": 1702297
			},
			"name": "key",
			"value": "test",
			"__metadata": {
				"uri": "/resin/device_service_environment_variable(@id)?@id=183330"
			}
		}
	]
}`)
	})
	// When
	err := client.DeviceServVar.DeleteWithName(context.Background(), DeviceUUID(uuid), key)
	// Then
	assert.NilError(t, err)
}

func TestDeviceServVarService_DeleteWithName_NotFound(t *testing.T) {
	// Given
	deviceID := int64(123456)
	key := "key"
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceServVarBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		expected := "%24filter=service_install/device+eq+%27" +
			strconv.FormatInt(deviceID, 10) + "%27+and+name+eq+%27" + key + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		fmt.Fprint(w, `{
	"d": []
}`)
	})
	// When
	err := client.DeviceServVar.DeleteWithName(context.Background(), DeviceID(deviceID), key)
	// Then
	assert.NilError(t, err)
}
