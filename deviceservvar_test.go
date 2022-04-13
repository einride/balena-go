package balena

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"go.einride.tech/balena/odata"
	"gotest.tools/v3/assert"
)

func TestDeviceServVarService_List_ID(t *testing.T) {
	// Given
	jsonResp := `{
	"d": [
		{
            "service_install": [
                {
                    "installs__service": [
                        {
                            "id": 1101713,
							"application": {
								"__id": 1660000,
								"__deferred": {
									"uri": "/resin/application(@id)?@id=1660000"
								}
							},
							"created_at": "2020-06-18T20:26:23.470Z",
                            "service_name": "vehiclelogger",
                            "__metadata": {
                                "uri": "/resin/service(@id)?@id=1101713"
                            }
                        }
                    ],
                    "id": 7185462,
            		"created_at": "2021-10-15T08:06:13.353Z",
                    "device": {
                        "__id": 1702297,
                        "__deferred": {
                            "uri": "/resin/device(@id)?@id=1702297"
                        }
                    },
                    "__metadata": {
                        "uri": "/resin/service_install(@id)?@id=7185462"
                    }
                }
            ],
			"id": 183330,
			"created_at": "2019-09-27T17:54:21.559Z",
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
		//nolint:lll
		expected := "%24filter=service_install/device+eq+%27123456%27&$expand=service_install($select=id,device,created_at;$expand=installs__service($select=id,service_name,created_at,application))"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		_, _ = fmt.Fprint(w, jsonResp)
	})
	expected := []*DeviceServVarResponse{
		{
			ID:        183330,
			CreatedAt: "2019-09-27T17:54:21.559Z",
			Name:      "key",
			Value:     "test",
			ServiceInstall: []*ServiceInstallResponse{
				{
					InstallsService: []InstallsService{
						{
							ID:          1101713,
							ServiceName: "vehiclelogger",
							Application: odata.Object{
								Deferred: odata.Deferred{
									URI: "/resin/application(@id)?@id=1660000",
								},
								ID: 1660000,
							},
							CreatedAt: time.Date(2020, 6, 18, 20, 26, 23, 470_000_000, time.UTC),
						},
					},
					ID:        7185462,
					CreatedAt: time.Date(2021, 10, 15, 8, 6, 13, 353_000_000, time.UTC),
					Device: odata.Object{
						Deferred: odata.Deferred{URI: "/resin/device(@id)?@id=1702297"},
						ID:       1702297,
					},
				},
			},
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
            "service_install": [
                {
                    "installs__service": [
                        {
                            "id": 1101713,
							"application": {
								"__id": 1660000,
								"__deferred": {
									"uri": "/resin/application(@id)?@id=1660000"
								}
							},
							"created_at": "2020-06-18T20:26:23.470Z",
                            "service_name": "vehiclelogger",
                            "__metadata": {
                                "uri": "/resin/service(@id)?@id=1101713"
                            }
                        }
                    ],
                    "id": 7185462,
            		"created_at": "2021-10-15T08:06:13.353Z",
                    "device": {
                        "__id": 1702297,
                        "__deferred": {
                            "uri": "/resin/device(@id)?@id=1702297"
                        }
                    },
                    "__metadata": {
                        "uri": "/resin/service_install(@id)?@id=7185462"
                    }
                }
            ],
			"id": 183330,
			"created_at": "2019-09-27T17:54:21.559Z",
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
		//nolint:lll
		expected := "%24filter=service_install/device/uuid+eq+%27" + uuid + "%27&$expand=service_install($select=id,device,created_at;$expand=installs__service($select=id,service_name,created_at,application))"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		_, _ = fmt.Fprint(w, jsonResp)
	})
	expected := []*DeviceServVarResponse{
		{
			ID:        183330,
			CreatedAt: "2019-09-27T17:54:21.559Z",
			Name:      "key",
			Value:     "test",
			ServiceInstall: []*ServiceInstallResponse{
				{
					InstallsService: []InstallsService{
						{
							ID: 1101713,
							Application: odata.Object{
								Deferred: odata.Deferred{
									URI: "/resin/application(@id)?@id=1660000",
								},
								ID: 1660000,
							},
							CreatedAt:   time.Date(2020, 6, 18, 20, 26, 23, 470_000_000, time.UTC),
							ServiceName: "vehiclelogger",
						},
					},
					ID:        7185462,
					CreatedAt: time.Date(2021, 10, 15, 8, 6, 13, 353_000_000, time.UTC),
					Device: odata.Object{
						Deferred: odata.Deferred{URI: "/resin/device(@id)?@id=1702297"},
						ID:       1702297,
					},
				},
			},
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
    "service_install": {
        "__id": 7185462,
        "__deferred": {
            "uri": "/resin/service_install(@id)?@id=7185462"
        }
    },
	"name": "key",
	"value": "test",
	"__metadata": {
		"uri": "/resin/device_service_environment_variable(@id)?@id=183330"
	}
}`
	serviceInstallID := int64(7185462)
	key := "key"
	value := "value"
	client, mux, cleanup := newFixture()
	defer cleanup()
	type request struct {
		ServiceInstallID int64  `json:"service_install"`
		Name             string `json:"name"`
		Value            string `json:"value"`
	}
	mux.HandleFunc("/v6/service_install", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"D": [{"id": %d}]}`, serviceInstallID)
	})
	mux.HandleFunc("/"+deviceServVarBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		b, err := ioutil.ReadAll(r.Body)
		assert.NilError(t, err)
		req := &request{}
		assert.NilError(t, json.Unmarshal(b, req))
		assert.Equal(t, serviceInstallID, req.ServiceInstallID)
		assert.Equal(t, key, req.Name)
		assert.Equal(t, value, req.Value)
		_, _ = fmt.Fprint(w, jsonResp)
	})
	expected := &DeviceServVarCreateResponse{
		ID:        183330,
		CreatedAt: "2019-09-27T17:54:21.559Z",
		ServiceInstall: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/service_install(@id)?@id=7185462"},
			ID:       7185462,
		},
		Name:  "key",
		Value: "test",
	}
	// When
	actual, err := client.DeviceServVar.Create(context.Background(), 7185462, key, value)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestDeviceServVarService_Create_UUID(t *testing.T) {
	// Given
	jsonResp := `{
	"id": 183330,
	"created_at": "2019-09-27T17:54:21.559Z",
    "service_install": {
        "__id": 7185462,
        "__deferred": {
            "uri": "/resin/service_install(@id)?@id=7185462"
        }
    },
	"name": "key",
	"value": "test",
	"__metadata": {
		"uri": "/resin/device_service_environment_variable(@id)?@id=183330"
	}
}`
	serviceInstallID := int64(7185462)
	key := "key"
	value := "value"
	client, mux, cleanup := newFixture()
	defer cleanup()
	type request struct {
		ServiceInstallID int64  `json:"service_install"`
		Name             string `json:"name"`
		Value            string `json:"value"`
	}
	mux.HandleFunc("/v6/service_install", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"D": [{"id": %d}]}`, serviceInstallID)
	})
	mux.HandleFunc("/"+deviceServVarBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		b, err := ioutil.ReadAll(r.Body)
		assert.NilError(t, err)
		req := &request{}
		assert.NilError(t, json.Unmarshal(b, req))
		assert.Equal(t, serviceInstallID, req.ServiceInstallID)
		assert.Equal(t, key, req.Name)
		assert.Equal(t, value, req.Value)
		_, _ = fmt.Fprint(w, jsonResp)
	})
	expected := &DeviceServVarCreateResponse{
		ID:        183330,
		CreatedAt: "2019-09-27T17:54:21.559Z",
		ServiceInstall: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/service_install(@id)?@id=7185462"},
			ID:       7185462,
		},
		Name:  "key",
		Value: "test",
	}
	// When
	actual, err := client.DeviceServVar.Create(context.Background(), 7185462, key, value)
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
		_, _ = fmt.Fprint(w, `{
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

func TestDeviceServVarService_Update(t *testing.T) {
	// Given
	uuid := "12345678901234567890"
	key := "key"
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceServVarBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		expected := "%24filter=service_install/device/uuid+eq+%27" + uuid + "%27+and+name+eq+%27" + key + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		_, _ = fmt.Fprint(w, "OK")
	})
	// When
	err := client.DeviceServVar.Update(context.Background(), DeviceUUID(uuid), key, "newVal")
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
		_, _ = fmt.Fprint(w, `{
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
		_, _ = fmt.Fprint(w, `{
	"d": []
}`)
	})
	// When
	err := client.DeviceServVar.DeleteWithName(context.Background(), DeviceID(deviceID), key)
	// Then
	assert.NilError(t, err)
}
