package balena

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"go.einride.tech/balena/odata"
	"gotest.tools/v3/assert"
)

func TestDeviceTagService_List_ID(t *testing.T) {
	// Given
	jsonResp := `{
		"d": [
		{
			"id": 610779,
			"device": {
				"__deferred": {
					"uri": "/resin/device(1701227)"
				},
				"__id": 1701227
			},
			"tag_key": "pinnedTo",
			"value": "rel1",
			"__metadata": {
				"uri": "/resin/device_tag(@id)?@id=610779"
			}
		}
	]
}`
	deviceID := int64(123456)
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=device/id+eq+%27" + strconv.FormatInt(deviceID, 10) + "%27"
		if r.URL.RawQuery != expected {
			t.Logf("query = %s ; expected %s\n", r.URL.RawQuery, expected)
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, jsonResp)
	})
	expected := []*DeviceTagResponse{
		{
			ID:     610779,
			TagKey: "pinnedTo",
			Device: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/device(1701227)"},
				ID:       1701227,
			},
			Value: "rel1",
		},
	}
	// When
	actual, err := client.DeviceTag.List(context.Background(), DeviceID(deviceID))
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestDeviceTagService_List_UUID(t *testing.T) {
	// Given
	jsonResp := `{
		"d": [
		{
			"id": 610779,
			"device": {
				"__deferred": {
					"uri": "/resin/device(1701227)"
				},
				"__id": 1701227
			},
			"tag_key": "pinnedTo",
			"value": "rel1",
			"__metadata": {
				"uri": "/resin/device_tag(@id)?@id=610779"
			}
		}
	]
}`
	uuid := "123456789123456798"
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=device/uuid+eq+%27" + uuid + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, jsonResp)
	})
	expected := []*DeviceTagResponse{
		{
			ID:     610779,
			TagKey: "pinnedTo",
			Device: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/device(1701227)"},
				ID:       1701227,
			},
			Value: "rel1",
		},
	}
	// When
	actual, err := client.DeviceTag.List(context.Background(), DeviceUUID(uuid))
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestDeviceTagService_Create_ID(t *testing.T) {
	// Given
	jsonResp := `{
	"id": 610779,
	"device": {
		"__deferred": {
			"uri": "/resin/device(1701227)"
		},
		"__id": 1701227
	},
	"tag_key": "pinnedTo",
	"value": "rel1",
	"__metadata": {
		"uri": "/resin/device_tag(@id)?@id=610779"
	}
}`
	deviceID := int64(1701227)
	key := "key"
	value := "value"
	client, mux, cleanup := newFixture()
	defer cleanup()
	type request struct {
		DeviceID string `json:"device"`
		Key      string `json:"tag_key"`
		Value    string `json:"value"`
	}
	mux.HandleFunc("/"+deviceTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		b, err := io.ReadAll(r.Body)
		assert.NilError(t, err)
		req := &request{}
		assert.NilError(t, json.Unmarshal(b, req))
		assert.Equal(t, strconv.FormatInt(deviceID, 10), req.DeviceID)
		assert.Equal(t, key, req.Key)
		assert.Equal(t, value, req.Value)
		fmt.Fprint(w, jsonResp)
	})
	expected := &DeviceTagResponse{
		ID:     610779,
		TagKey: "pinnedTo",
		Device: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/device(1701227)"},
			ID:       1701227,
		},
		Value: "rel1",
	}
	// When
	actual, err := client.DeviceTag.Create(context.Background(), DeviceID(deviceID), key, value)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestDeviceTagService_Create_UUID(t *testing.T) {
	// Given
	jsonResp := `{
	"id": 610779,
	"device": {
		"__deferred": {
			"uri": "/resin/device(1701227)"
		},
		"__id": 1701227
	},
	"tag_key": "pinnedTo",
	"value": "rel1",
	"__metadata": {
		"uri": "/resin/device_tag(@id)?@id=610779"
	}
}`
	uuid := "1234567890123456789"
	deviceID := int64(123456)
	key := "key"
	value := "value"
	client, mux, cleanup := newFixture()
	defer cleanup()
	type request struct {
		DeviceID string `json:"device"`
		Key      string `json:"tag_key"`
		Value    string `json:"value"`
	}
	mux.HandleFunc("/"+deviceBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=uuid+eq+%27" + uuid + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		resp := `{"d":[{"id":123456}]}`
		fmt.Fprint(w, resp)
	})
	mux.HandleFunc("/"+deviceTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		b, err := io.ReadAll(r.Body)
		assert.NilError(t, err)
		req := &request{}
		assert.NilError(t, json.Unmarshal(b, req))
		assert.Equal(t, strconv.FormatInt(deviceID, 10), req.DeviceID)
		assert.Equal(t, key, req.Key)
		assert.Equal(t, value, req.Value)
		fmt.Fprint(w, jsonResp)
	})
	expected := &DeviceTagResponse{
		ID:     610779,
		TagKey: "pinnedTo",
		Device: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/device(1701227)"},
			ID:       1701227,
		},
		Value: "rel1",
	}
	// When
	actual, err := client.DeviceTag.Create(context.Background(), DeviceID(deviceID), key, value)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestDeviceTagService_GetWithKey_ID(t *testing.T) {
	// Given
	jsonResp := `{
		"d": [
		{
			"id": 610779,
			"device": {
				"__deferred": {
					"uri": "/resin/device(1701227)"
				},
				"__id": 1701227
			},
			"tag_key": "key",
			"value": "value",
			"__metadata": {
				"uri": "/resin/device_tag(@id)?@id=610779"
			}
		}
	]
}`
	deviceID := int64(1701227)
	key := "key"
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=device/id+eq+%27" + strconv.FormatInt(deviceID, 10) + "%27+and+tag_key+eq+%27" + key + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, jsonResp)
	})
	expected := &DeviceTagResponse{
		ID:     610779,
		TagKey: "key",
		Device: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/device(1701227)"},
			ID:       1701227,
		},
		Value: "value",
	}
	// When
	actual, err := client.DeviceTag.GetWithKey(context.Background(), DeviceID(deviceID), key)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestDeviceTagService_GetWithKey_UUID(t *testing.T) {
	// Given
	jsonResp := `{
		"d": [
		{
			"id": 610779,
			"device": {
				"__deferred": {
					"uri": "/resin/device(1701227)"
				},
				"__id": 1701227
			},
			"tag_key": "key",
			"value": "value",
			"__metadata": {
				"uri": "/resin/device_tag(@id)?@id=610779"
			}
		}
	]
}`
	deviceUUID := "1234567890987654321"
	key := "key"
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=device/uuid+eq+%27" + deviceUUID + "%27+and+tag_key+eq+%27" + key + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, jsonResp)
	})
	expected := &DeviceTagResponse{
		ID:     610779,
		TagKey: "key",
		Device: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/device(1701227)"},
			ID:       1701227,
		},
		Value: "value",
	}
	// When
	actual, err := client.DeviceTag.GetWithKey(context.Background(), DeviceUUID(deviceUUID), key)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestDeviceTagService_GetWithKey_NotFound(t *testing.T) {
	// Given
	jsonResp := `{
		"d": []
}`
	deviceID := int64(1701227)
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, jsonResp)
	})
	// When
	actual, err := client.DeviceTag.GetWithKey(context.Background(), DeviceID(deviceID), "key")
	// Then
	assert.NilError(t, err)
	assert.Assert(t, actual == nil)
}

func TestDeviceTagService_UpdateWithKey_UUID(t *testing.T) {
	// Given
	uuid := "1234567890987654321"
	key := "key"
	value := "newValue"
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		expected := "%24filter=device/uuid+eq+%27" + uuid + "%27+and+tag_key+eq+%27" + key + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, "OK")
	})
	// When
	err := client.DeviceTag.UpdateWithKey(context.Background(), DeviceUUID(uuid), key, value)
	// Then
	assert.NilError(t, err)
}

func TestDeviceTagService_UpdateWithKey_ID(t *testing.T) {
	// Given
	id := int64(123)
	key := "key"
	value := "newValue"
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		expected := "%24filter=device/id+eq+%27" + strconv.FormatInt(id, 10) + "%27+and+tag_key+eq+%27" + key + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, "OK")
	})
	// When
	err := client.DeviceTag.UpdateWithKey(context.Background(), DeviceID(id), key, value)
	// Then
	assert.NilError(t, err)
}

func TestDeviceTagService_DeleteWithKey_ID(t *testing.T) {
	// Given
	jsonResp := `{
		"d": [
		{
			"id": 610779,
			"device": {
				"__deferred": {
					"uri": "/resin/device(1701227)"
				},
				"__id": 1701227
			},
			"tag_key": "key",
			"value": "value",
			"__metadata": {
				"uri": "/resin/device_tag(@id)?@id=610779"
			}
		}
	]
}`
	deviceID := int64(1701227)
	key := "key"
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		expected := "%24filter=device/id+eq+%27" + strconv.FormatInt(deviceID, 10) + "%27+and+tag_key+eq+%27" + key + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, jsonResp)
	})
	// When
	err := client.DeviceTag.DeleteWithKey(context.Background(), DeviceID(deviceID), key)
	// Then
	assert.NilError(t, err)
}

func TestDeviceTagService_DeleteWithKey_UUID(t *testing.T) {
	// Given
	jsonResp := `{
		"d": [
		{
			"id": 610779,
			"device": {
				"__deferred": {
					"uri": "/resin/device(1701227)"
				},
				"__id": 1701227
			},
			"tag_key": "key",
			"value": "value",
			"__metadata": {
				"uri": "/resin/device_tag(@id)?@id=610779"
			}
		}
	]
}`
	uuid := "1234567890987654321"
	key := "key"
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		expected := "%24filter=device/uuid+eq+%27" + uuid + "%27+and+tag_key+eq+%27" + key + "%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, jsonResp)
	})
	// When
	err := client.DeviceTag.DeleteWithKey(context.Background(), DeviceUUID(uuid), key)
	// Then
	assert.NilError(t, err)
}

func TestDeviceTagService_GetWithQuery(t *testing.T) {
	// Given
	jsonResp := `{
		"d": [
		{
			"id": 610779,
			"device": {
				"__deferred": {
					"uri": "/resin/device(1701227)"
				},
				"__id": 1701227
			},
			"tag_key": "key",
			"value": "value",
			"__metadata": {
				"uri": "/resin/device_tag(@id)?@id=610779"
			}
		}
	]
}`
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+deviceTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=key+eq+%27value%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, jsonResp)
	})
	expected := []*DeviceTagResponse{
		{
			ID:     610779,
			TagKey: "key",
			Device: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/device(1701227)"},
				ID:       1701227,
			},
			Value: "value",
		},
	}
	// When
	actual, err := client.DeviceTag.GetWithQuery(context.Background(), "%24filter=key+eq+%27value%27")
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}
