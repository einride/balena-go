package balena

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/einride/balena-go/odata"
	"gotest.tools/v3/assert"
)

const (
	applicationListResponse = `{
	"d": [
		{
			"id": 1514287,
			"user": {
				"__deferred": {
					"uri": "/resin/user(123437)"
				},
				"__id": 123437
			},
			"depends_on__application": null,
			"actor": 4085719,
			"app_name": "dev-device",
			"slug": "gh_alethenorio/dev-device",
			"commit": "0dfe13d62efb4325385952e1e15361f7",
			"application_type": {
				"__deferred": {
					"uri": "/resin/application_type(5)"
				},
				"__id": 5
			},
			"device_type": "intel-nuc",
			"should_track_latest_release": true,
			"is_accessible_by_support_until__date": null,
			"__metadata": {
				"uri": "/resin/application(@id)?@id=1514287"
			}
		}
	]
}`
)

func TestApplicationService_List(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+applicationBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, applicationListResponse)
	})
	expected := []*ApplicationsResponse{
		{
			ID:                         1514287,
			User:                       odata.Object{Deferred: odata.Deferred{URI: "/resin/user(123437)"}, ID: 123437},
			DependsOnApplication:       nil,
			Actor:                      4085719,
			Name:                       "dev-device",
			Slug:                       "gh_alethenorio/dev-device",
			Commit:                     "0dfe13d62efb4325385952e1e15361f7",
			ApplicationType:            odata.Object{Deferred: odata.Deferred{URI: "/resin/application_type(5)"}, ID: 5},
			DeviceType:                 "intel-nuc",
			TrackLatestRelease:         true,
			IsAccessibleBySupportUntil: nil,
		},
	}
	// When
	actual, err := client.Application.List(context.Background())
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestApplicationService_GetWithQuery(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc(
		"/"+applicationBasePath,
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			expected := "%24filter=app_name%20eq%20%27dev-device%27"
			if r.URL.RawQuery != expected {
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
				return
			}
			fmt.Fprint(w, applicationListResponse)
		},
	)
	expected := []*ApplicationsResponse{
		{
			ID:                         1514287,
			User:                       odata.Object{Deferred: odata.Deferred{URI: "/resin/user(123437)"}, ID: 123437},
			DependsOnApplication:       nil,
			Actor:                      4085719,
			Name:                       "dev-device",
			Slug:                       "gh_alethenorio/dev-device",
			Commit:                     "0dfe13d62efb4325385952e1e15361f7",
			ApplicationType:            odata.Object{Deferred: odata.Deferred{URI: "/resin/application_type(5)"}, ID: 5},
			DeviceType:                 "intel-nuc",
			TrackLatestRelease:         true,
			IsAccessibleBySupportUntil: nil,
		},
	}
	// When
	actual, err := client.Application.GetWithQuery(context.Background(), "%24filter=app_name%20eq%20%27dev-device%27")
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestApplicationService_Get(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityID := int64(1514287)
	mux.HandleFunc(
		"/"+odata.EntityURL(applicationBasePath, strconv.FormatInt(entityID, 10)),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, applicationListResponse)
		},
	)
	expected := &ApplicationsResponse{
		ID:                         1514287,
		User:                       odata.Object{Deferred: odata.Deferred{URI: "/resin/user(123437)"}, ID: 123437},
		DependsOnApplication:       nil,
		Actor:                      4085719,
		Name:                       "dev-device",
		Slug:                       "gh_alethenorio/dev-device",
		Commit:                     "0dfe13d62efb4325385952e1e15361f7",
		ApplicationType:            odata.Object{Deferred: odata.Deferred{URI: "/resin/application_type(5)"}, ID: 5},
		DeviceType:                 "intel-nuc",
		TrackLatestRelease:         true,
		IsAccessibleBySupportUntil: nil,
	}
	// When
	actual, err := client.Application.Get(context.Background(), entityID)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestApplicationService_Get_NotFound(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityID := int64(1514287)
	mux.HandleFunc(
		"/"+odata.EntityURL(applicationBasePath, strconv.FormatInt(entityID, 10)),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, `{"d": []}`)
		},
	)
	// When
	actual, err := client.Application.Get(context.Background(), entityID)
	// Then
	assert.NilError(t, err)
	assert.Assert(t, actual == nil)
}

func TestApplicationService_GetByName(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityName := "dev-device"
	mux.HandleFunc(
		"/"+applicationBasePath,
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			expected := "%24filter=app_name%20eq%20%27dev-device%27"
			if r.URL.RawQuery != expected {
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
				return
			}
			fmt.Fprint(w, applicationListResponse)
		},
	)
	expected := &ApplicationsResponse{
		ID:                         1514287,
		User:                       odata.Object{Deferred: odata.Deferred{URI: "/resin/user(123437)"}, ID: 123437},
		DependsOnApplication:       nil,
		Actor:                      4085719,
		Name:                       "dev-device",
		Slug:                       "gh_alethenorio/dev-device",
		Commit:                     "0dfe13d62efb4325385952e1e15361f7",
		ApplicationType:            odata.Object{Deferred: odata.Deferred{URI: "/resin/application_type(5)"}, ID: 5},
		DeviceType:                 "intel-nuc",
		TrackLatestRelease:         true,
		IsAccessibleBySupportUntil: nil,
	}
	// When
	actual, err := client.Application.GetByName(context.Background(), entityName)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestApplicationService_GetByName_NotFound(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityName := "dev-device"
	mux.HandleFunc(
		"/"+applicationBasePath,
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			expected := "%24filter=app_name%20eq%20%27dev-device%27"
			if r.URL.RawQuery != expected {
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), 500)
				return
			}
			fmt.Fprint(w, `{"d":[]}`)
		},
	)
	// When
	actual, err := client.Application.GetByName(context.Background(), entityName)
	// Then
	assert.NilError(t, err)
	assert.Assert(t, actual == nil)
}

func TestApplicationService_EnableTrackLatestRelease(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityID := int64(1514287)
	mux.HandleFunc(
		"/"+odata.EntityURL(applicationBasePath, strconv.FormatInt(entityID, 10)),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPatch)
			b, err := ioutil.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(t, `{"should_track_latest_release":true}`+"\n", string(b))
			fmt.Fprint(w, "OK")
		},
	)
	// When
	resp, err := client.Application.EnableTrackLatestRelease(context.Background(), entityID)
	// Then
	assert.NilError(t, err)
	assert.Equal(t, "OK", string(resp))
}

func TestApplicationService_DisableTrackLatestRelease(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityID := int64(1514287)
	mux.HandleFunc(
		"/"+odata.EntityURL(applicationBasePath, strconv.FormatInt(entityID, 10)),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPatch)
			b, err := ioutil.ReadAll(r.Body)
			assert.NilError(t, err)
			assert.Equal(t, `{"should_track_latest_release":false}`+"\n", string(b))
			fmt.Fprint(w, "OK")
		},
	)
	// When
	resp, err := client.Application.DisableTrackLatestRelease(context.Background(), entityID)
	// Then
	assert.NilError(t, err)
	assert.Equal(t, "OK", string(resp))
}
