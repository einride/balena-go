package balena

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"go.einride.tech/balena/odata"
	"gotest.tools/v3/assert"
)

const (
	applicationListResponse = `{
	"d": [
		{
    	  "id": 1234567,
    	  "organization": {
    	    "__id": 122333,
    	    "__deferred": {
    	      "uri": "/resin/organization(@id)?@id=122333"
    	    }
    	  },
    	  "depends_on__application": null,
    	  "actor": 700068,
    	  "app_name": "Stellarium",
    	  "slug": "david_tischler1/stellarium",
    	  "should_be_running__release": {
    	    "__id": 1798244,
    	    "__deferred": {
    	      "uri": "/resin/release(@id)?@id=1798244"
    	    }
    	  },
    	  "application_type": {
    	    "__id": 4,
    	    "__deferred": {
    	      "uri": "/resin/application_type(@id)?@id=4"
    	    }
    	  },
    	  "is_for__device_type": {
    	    "__id": 60,
    	    "__deferred": {
    	      "uri": "/resin/device_type(@id)?@id=60"
    	    }
    	  },
    	  "is_of__class": "fleet",
    	  "should_track_latest_release": true,
    	  "is_accessible_by_support_until__date": null,
    	  "is_public": true,
    	  "is_host": false,
    	  "is_archived": false,
    	  "is_discoverable": true,
    	  "is_stored_at__repository_url": "https://github.com/balenalabs-incubator/stellarium",
    	  "created_at": "2020-12-08T21:02:03.327Z",
    	  "uuid": "fc02cb0c1f174d10811303446cde8aae",
    	  "is_of__class": "fleet"
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
			ID:   1234567,
			UUID: "fc02cb0c1f174d10811303446cde8aae",
			Organization: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/organization(@id)?@id=122333"}, ID: 122333,
			},
			DependsOnApplication: nil,
			Actor:                700068,
			AppName:              "Stellarium",
			Slug:                 "david_tischler1/stellarium",
			IsOfClass:            "fleet",
			ShouldBeRunningRelease: &ReleaseOData{
				Object: &odata.Object{
					Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1798244"}, ID: 1798244,
				},
			},
			IsForDeviceType: &DeviceTypeOData{
				Object: &odata.Object{
					Deferred: odata.Deferred{URI: "/resin/device_type(@id)?@id=60"}, ID: 60,
				},
			},
			ShouldTrackLatestRelease:       true,
			IsAccessibleBySupportUntilDate: nil,
			IsPublic:                       true,
			IsHost:                         false,
			IsArchived:                     false,
			IsDiscoverable:                 true,
			IsStoredAtRepositoryURL:        "https://github.com/balenalabs-incubator/stellarium",
			CreatedAt:                      "2020-12-08T21:02:03.327Z",
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
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, applicationListResponse)
		},
	)
	expected := []*ApplicationsResponse{
		{
			ID:   1234567,
			UUID: "fc02cb0c1f174d10811303446cde8aae",
			Organization: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/organization(@id)?@id=122333"}, ID: 122333,
			},
			DependsOnApplication: nil,
			Actor:                700068,
			AppName:              "Stellarium",
			Slug:                 "david_tischler1/stellarium",
			IsOfClass:            "fleet",
			ShouldBeRunningRelease: &ReleaseOData{
				Object: &odata.Object{
					Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1798244"}, ID: 1798244,
				},
			},
			IsForDeviceType: &DeviceTypeOData{
				Object: &odata.Object{
					Deferred: odata.Deferred{URI: "/resin/device_type(@id)?@id=60"}, ID: 60,
				},
			},
			ShouldTrackLatestRelease:       true,
			IsAccessibleBySupportUntilDate: nil,
			IsPublic:                       true,
			IsHost:                         false,
			IsArchived:                     false,
			IsDiscoverable:                 true,
			IsStoredAtRepositoryURL:        "https://github.com/balenalabs-incubator/stellarium",
			CreatedAt:                      "2020-12-08T21:02:03.327Z",
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
		ID:   1234567,
		UUID: "fc02cb0c1f174d10811303446cde8aae",
		Organization: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/organization(@id)?@id=122333"}, ID: 122333,
		},
		DependsOnApplication: nil,
		Actor:                700068,
		AppName:              "Stellarium",
		Slug:                 "david_tischler1/stellarium",
		IsOfClass:            "fleet",
		ShouldBeRunningRelease: &ReleaseOData{
			Object: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1798244"}, ID: 1798244,
			},
		},
		IsForDeviceType: &DeviceTypeOData{
			Object: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/device_type(@id)?@id=60"}, ID: 60,
			},
		},
		ShouldTrackLatestRelease:       true,
		IsAccessibleBySupportUntilDate: nil,
		IsPublic:                       true,
		IsHost:                         false,
		IsArchived:                     false,
		IsDiscoverable:                 true,
		IsStoredAtRepositoryURL:        "https://github.com/balenalabs-incubator/stellarium",
		CreatedAt:                      "2020-12-08T21:02:03.327Z",
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
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, applicationListResponse)
		},
	)
	expected := &ApplicationsResponse{
		ID:   1234567,
		UUID: "fc02cb0c1f174d10811303446cde8aae",
		Organization: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/organization(@id)?@id=122333"}, ID: 122333,
		},
		DependsOnApplication: nil,
		Actor:                700068,
		AppName:              "Stellarium",
		Slug:                 "david_tischler1/stellarium",
		IsOfClass:            "fleet",
		ShouldBeRunningRelease: &ReleaseOData{
			Object: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/release(@id)?@id=1798244"}, ID: 1798244,
			},
		},
		IsForDeviceType: &DeviceTypeOData{
			Object: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/device_type(@id)?@id=60"}, ID: 60,
			},
		},
		ShouldTrackLatestRelease:       true,
		IsAccessibleBySupportUntilDate: nil,
		IsPublic:                       true,
		IsHost:                         false,
		IsArchived:                     false,
		IsDiscoverable:                 true,
		IsStoredAtRepositoryURL:        "https://github.com/balenalabs-incubator/stellarium",
		CreatedAt:                      "2020-12-08T21:02:03.327Z",
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
				http.Error(w, fmt.Sprintf("query = %s ; expected %s", r.URL.RawQuery, expected), http.StatusInternalServerError)
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
			b, err := io.ReadAll(r.Body)
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
			b, err := io.ReadAll(r.Body)
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
