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
	releaseResponse = `{
	"d": [
		{
			"id": 1079426,
			"created_at": "2019-09-25T18:51:16.070Z",
			"belongs_to__application": {
				"__deferred": {
					"uri": "/resin/application(1514287)"
				},
				"__id": 1514287
			},
			"is_created_by__user": {
				"__deferred": {
					"uri": "/resin/user(122737)"
				},
				"__id": 122737
			},
			"commit": "0dfe02d62efb4325385952e1e15361f7",
			"composition": {
				"version": "2.1",
				"networks": {},
				"volumes": {
					"resin-data": {},
					"log-data": {},
					"balena-lock": {},
					"deepmap-map": {}
				},
				"services": {
					"nginx": {
						"build": {
							"context": "./nginx"
						},
						"network_mode": "host",
						"volumes": [
							"resin-data:/data",
							"log-data:/data/logs"
						],
						"labels": {
							"io.resin.features.supervisor-api": "1",
							"io.balena.features.balena-api": "1"
						}
					},
					"prometheus": {
						"build": {
							"context": "./prometheus"
						},
						"restart": "always",
						"network_mode": "host",
						"labels": {
							"io.resin.features.supervisor-api": "1",
							"io.balena.features.balena-api": "1"
						}
					},
					"deploymentguard": {
						"build": {
							"context": "./deploymentguard"
						},
						"restart": "always",
						"network_mode": "host",
						"labels": {
							"io.resin.features.supervisor-api": "1",
							"io.balena.features.balena-api": "1"
						}
					}
				}
			},
			"status": "success",
			"source": "local",
			"start_timestamp": "2019-09-25T18:51:15.773Z",
			"end_timestamp": "2019-09-25T18:51:35.835Z",
			"update_timestamp": "2019-09-25T18:51:36.215Z",
			"__metadata": {
				"uri": "/resin/release(@id)?@id=1079426"
			}
		}
	]
}`
)

func TestReleaseService_List(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+releaseBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, releaseResponse)
	})
	expected := []*ReleaseResponse{
		{
			ID:        1079426,
			CreatedAt: "2019-09-25T18:51:16.070Z",
			BelongsToApplication: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/application(1514287)"}, ID: 1514287,
			},
			CreatedByUser: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/user(122737)"}, ID: 122737,
			},
			Commit: "0dfe02d62efb4325385952e1e15361f7",
			Composition: []byte(`{
				"version": "2.1",
				"networks": {},
				"volumes": {
					"resin-data": {},
					"log-data": {},
					"balena-lock": {},
					"deepmap-map": {}
				},
				"services": {
					"nginx": {
						"build": {
							"context": "./nginx"
						},
						"network_mode": "host",
						"volumes": [
							"resin-data:/data",
							"log-data:/data/logs"
						],
						"labels": {
							"io.resin.features.supervisor-api": "1",
							"io.balena.features.balena-api": "1"
						}
					},
					"prometheus": {
						"build": {
							"context": "./prometheus"
						},
						"restart": "always",
						"network_mode": "host",
						"labels": {
							"io.resin.features.supervisor-api": "1",
							"io.balena.features.balena-api": "1"
						}
					},
					"deploymentguard": {
						"build": {
							"context": "./deploymentguard"
						},
						"restart": "always",
						"network_mode": "host",
						"labels": {
							"io.resin.features.supervisor-api": "1",
							"io.balena.features.balena-api": "1"
						}
					}
				}
			}`),
			Status:          "success",
			Source:          "local",
			StartTimestamp:  "2019-09-25T18:51:15.773Z",
			EndTimestamp:    "2019-09-25T18:51:35.835Z",
			UpdateTimestamp: "2019-09-25T18:51:36.215Z",
		},
	}
	// When
	actual, err := client.Release.List(context.Background())
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestReleaseService_Get(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityID := int64(112233)
	mux.HandleFunc(
		"/"+odata.EntityURL(releaseBasePath, strconv.FormatInt(entityID, 10)),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, releaseResponse)
		},
	)
	expected := &ReleaseResponse{
		ID:        1079426,
		CreatedAt: "2019-09-25T18:51:16.070Z",
		BelongsToApplication: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/application(1514287)"}, ID: 1514287,
		},
		CreatedByUser: odata.Object{
			Deferred: odata.Deferred{URI: "/resin/user(122737)"}, ID: 122737,
		},
		Commit: "0dfe02d62efb4325385952e1e15361f7",
		Composition: []byte(`{
				"version": "2.1",
				"networks": {},
				"volumes": {
					"resin-data": {},
					"log-data": {},
					"balena-lock": {},
					"deepmap-map": {}
				},
				"services": {
					"nginx": {
						"build": {
							"context": "./nginx"
						},
						"network_mode": "host",
						"volumes": [
							"resin-data:/data",
							"log-data:/data/logs"
						],
						"labels": {
							"io.resin.features.supervisor-api": "1",
							"io.balena.features.balena-api": "1"
						}
					},
					"prometheus": {
						"build": {
							"context": "./prometheus"
						},
						"restart": "always",
						"network_mode": "host",
						"labels": {
							"io.resin.features.supervisor-api": "1",
							"io.balena.features.balena-api": "1"
						}
					},
					"deploymentguard": {
						"build": {
							"context": "./deploymentguard"
						},
						"restart": "always",
						"network_mode": "host",
						"labels": {
							"io.resin.features.supervisor-api": "1",
							"io.balena.features.balena-api": "1"
						}
					}
				}
			}`),
		Status:          "success",
		Source:          "local",
		StartTimestamp:  "2019-09-25T18:51:15.773Z",
		EndTimestamp:    "2019-09-25T18:51:35.835Z",
		UpdateTimestamp: "2019-09-25T18:51:36.215Z",
	}
	// When
	actual, err := client.Release.Get(context.Background(), entityID)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestReleaseService_Get_NotFound(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityID := int64(112233)
	mux.HandleFunc(
		"/"+odata.EntityURL(releaseBasePath, strconv.FormatInt(entityID, 10)),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, `{"d":[]}`)
		},
	)
	// When
	release, err := client.Release.Get(context.Background(), entityID)
	// Then
	assert.NilError(t, err)
	assert.Assert(t, release == nil)
}

func TestReleaseService_GetWithQuery(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+releaseBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=key+eq+%27value%27"
		if r.URL.RawQuery != expected {
			http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
			return
		}
		fmt.Fprint(w, releaseResponse)
	})
	expected := []*ReleaseResponse{
		{
			ID:        1079426,
			CreatedAt: "2019-09-25T18:51:16.070Z",
			BelongsToApplication: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/application(1514287)"}, ID: 1514287,
			},
			CreatedByUser: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/user(122737)"}, ID: 122737,
			},
			Commit: "0dfe02d62efb4325385952e1e15361f7",
			Composition: []byte(`{
				"version": "2.1",
				"networks": {},
				"volumes": {
					"resin-data": {},
					"log-data": {},
					"balena-lock": {},
					"deepmap-map": {}
				},
				"services": {
					"nginx": {
						"build": {
							"context": "./nginx"
						},
						"network_mode": "host",
						"volumes": [
							"resin-data:/data",
							"log-data:/data/logs"
						],
						"labels": {
							"io.resin.features.supervisor-api": "1",
							"io.balena.features.balena-api": "1"
						}
					},
					"prometheus": {
						"build": {
							"context": "./prometheus"
						},
						"restart": "always",
						"network_mode": "host",
						"labels": {
							"io.resin.features.supervisor-api": "1",
							"io.balena.features.balena-api": "1"
						}
					},
					"deploymentguard": {
						"build": {
							"context": "./deploymentguard"
						},
						"restart": "always",
						"network_mode": "host",
						"labels": {
							"io.resin.features.supervisor-api": "1",
							"io.balena.features.balena-api": "1"
						}
					}
				}
			}`),
			Status:          "success",
			Source:          "local",
			StartTimestamp:  "2019-09-25T18:51:15.773Z",
			EndTimestamp:    "2019-09-25T18:51:35.835Z",
			UpdateTimestamp: "2019-09-25T18:51:36.215Z",
		},
	}
	// When
	actual, err := client.Release.GetWithQuery(context.Background(), "%24filter=key+eq+%27value%27")
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}
