package balena

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"go.einride.tech/balena/odata"
	"gotest.tools/v3/assert"
)

const (
	deviceTypeResponse = `{
	"d": [
	  {
		"id": 34,
		"slug": "jetson-tx2",
		"name": "Nvidia Jetson TX2",
		"is_private": false,
		"is_of__cpu_architecture": {
		  "__id": 1,
		  "__deferred": {
			"uri": "/resin/cpu_architecture(@id)?@id=1"
		  }
		},
		"belongs_to__device_family": null,
		"__metadata": {
		  "uri": "/resin/device_type(@id)?@id=34"
		}
	  }
	  ]
	}`
)

func TestDeviceTypeService_Get(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityID := int64(34)
	mux.HandleFunc(
		"/"+odata.EntityURL(deviceTypeBasePath, strconv.FormatInt(entityID, 10)),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, deviceTypeResponse)
		},
	)
	expected := &DeviceTypeResponse{
		ID:        34,
		Name:      "Nvidia Jetson TX2",
		Slug:      "jetson-tx2",
		IsPrivate: false,
		IsOfCPUArchitecture: odata.Object{
			ID: 1,
			Deferred: odata.Deferred{
				URI: "/resin/cpu_architecture(@id)?@id=1",
			},
		},
	}
	// When
	actual, err := client.DeviceType.Get(context.Background(), entityID)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestDeviceTypeService_Get_NotFound(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityID := int64(34)
	mux.HandleFunc(
		"/"+odata.EntityURL(deviceTypeBasePath, strconv.FormatInt(entityID, 10)),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, `{"d":[]}`)
		},
	)
	// When
	device, err := client.DeviceType.Get(context.Background(), entityID)
	// Then
	assert.NilError(t, err)
	assert.Assert(t, device == nil)
}

func TestDeviceTypeService_GetWithQuery(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceTypeBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=key+eq+%27value%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, deviceTypeResponse)
	})
	expected := []*DeviceTypeResponse{
		{
			ID:        34,
			Name:      "Nvidia Jetson TX2",
			Slug:      "jetson-tx2",
			IsPrivate: false,
			IsOfCPUArchitecture: odata.Object{
				ID: 1,
				Deferred: odata.Deferred{
					URI: "/resin/cpu_architecture(@id)?@id=1",
				},
			},
		},
	}
	// When
	actual, err := client.DeviceType.GetWithQuery(context.Background(), "%24filter=key+eq+%27value%27")
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}
