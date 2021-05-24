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
				"id": 1798244,
    	  		"created_at": "2021-05-13T21:56:07.112Z",
    	  		"belongs_to__application": {
    	  		  "__id": 1234567,
    	  		  "__deferred": {
    	  		    "uri": "/resin/application(@id)?@id=1234567"
    	  		  }
    	  		},
    	  		"is_created_by__user": {
    	  		  "__id": 7654321,
    	  		  "__deferred": {
    	  		    "uri": "/resin/user(@id)?@id=7654321"
    	  		  }
    	  		},
    	  		"commit": "9c4e6991df722d2d13693e7ee5ad6039",
    	  		"composition": {
    	  		  "version": "2.1",
    	  		  "volumes": {
    	  		    "settings": {}
    	  		  },
    	  		  "services": {
    	  		    "stellarium": {
    	  		      "build": {
    	  		        "context": "./stellarium"
    	  		      },
    	  		      "restart": "always",
    	  		      "network_mode": "host",
    	  		      "privileged": true,
    	  		      "volumes": [
    	  		        "settings:/data/stellarium"
    	  		      ],
    	  		      "shm_size": "2gb"
    	  		    }
    	  		  }
    	  		},
    	  		"status": "success",
    	  		"source": "cloud",
    	  		"build_log": "[Info] ...",
    	  		"start_timestamp": "2021-05-13T21:56:07.084Z",
    	  		"end_timestamp": "2021-05-13T21:59:40.919Z",
    	  		"update_timestamp": "2021-05-13T21:59:46.976Z",
    	  		"is_invalidated": false,
    	  		"release_version": null,
    	  		"release_type": "final",
    	  		"is_passing_tests": true,
    	  		"contract": "{\"name\":\"Stellarium\",\"description\":\"An educational star, constellation, planet\"}"
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
			ID:        1798244,
			CreatedAt: "2021-05-13T21:56:07.112Z",
			BelongsToApplication: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/application(@id)?@id=1234567"}, ID: 1234567,
			},
			CreatedByUser: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/user(@id)?@id=7654321"}, ID: 7654321,
			},
			Commit: "9c4e6991df722d2d13693e7ee5ad6039",
			Composition: []byte(`{
    	  		  "version": "2.1",
    	  		  "volumes": {
    	  		    "settings": {}
    	  		  },
    	  		  "services": {
    	  		    "stellarium": {
    	  		      "build": {
    	  		        "context": "./stellarium"
    	  		      },
    	  		      "restart": "always",
    	  		      "network_mode": "host",
    	  		      "privileged": true,
    	  		      "volumes": [
    	  		        "settings:/data/stellarium"
    	  		      ],
    	  		      "shm_size": "2gb"
    	  		    }
    	  		  }
    	  		}`),
			Status:          "success",
			Source:          "cloud",
			StartTimestamp:  "2021-05-13T21:56:07.084Z",
			EndTimestamp:    "2021-05-13T21:59:40.919Z",
			UpdateTimestamp: "2021-05-13T21:59:46.976Z",
			BuildLog:        "[Info] ...",
			IsInvalidated:   false,
			ReleaseVersion:  nil,
			ReleaseType:     "final",
			IsPassingTests:  true,
			Contract:        "{\"name\":\"Stellarium\",\"description\":\"An educational star, constellation, planet\"}",
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
		ID:        1798244,
		CreatedAt: "2021-05-13T21:56:07.112Z",
		BelongsToApplication: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/application(@id)?@id=1234567"}, ID: 1234567,
		},
		CreatedByUser: &odata.Object{
			Deferred: odata.Deferred{URI: "/resin/user(@id)?@id=7654321"}, ID: 7654321,
		},
		Commit: "9c4e6991df722d2d13693e7ee5ad6039",
		Composition: []byte(`{
    	  		  "version": "2.1",
    	  		  "volumes": {
    	  		    "settings": {}
    	  		  },
    	  		  "services": {
    	  		    "stellarium": {
    	  		      "build": {
    	  		        "context": "./stellarium"
    	  		      },
    	  		      "restart": "always",
    	  		      "network_mode": "host",
    	  		      "privileged": true,
    	  		      "volumes": [
    	  		        "settings:/data/stellarium"
    	  		      ],
    	  		      "shm_size": "2gb"
    	  		    }
    	  		  }
    	  		}`),
		Status:          "success",
		Source:          "cloud",
		StartTimestamp:  "2021-05-13T21:56:07.084Z",
		EndTimestamp:    "2021-05-13T21:59:40.919Z",
		UpdateTimestamp: "2021-05-13T21:59:46.976Z",
		BuildLog:        "[Info] ...",
		IsInvalidated:   false,
		ReleaseVersion:  nil,
		ReleaseType:     "final",
		IsPassingTests:  true,
		Contract:        "{\"name\":\"Stellarium\",\"description\":\"An educational star, constellation, planet\"}",
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
			ID:        1798244,
			CreatedAt: "2021-05-13T21:56:07.112Z",
			BelongsToApplication: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/application(@id)?@id=1234567"}, ID: 1234567,
			},
			CreatedByUser: &odata.Object{
				Deferred: odata.Deferred{URI: "/resin/user(@id)?@id=7654321"}, ID: 7654321,
			},
			Commit: "9c4e6991df722d2d13693e7ee5ad6039",
			Composition: []byte(`{
    	  		  "version": "2.1",
    	  		  "volumes": {
    	  		    "settings": {}
    	  		  },
    	  		  "services": {
    	  		    "stellarium": {
    	  		      "build": {
    	  		        "context": "./stellarium"
    	  		      },
    	  		      "restart": "always",
    	  		      "network_mode": "host",
    	  		      "privileged": true,
    	  		      "volumes": [
    	  		        "settings:/data/stellarium"
    	  		      ],
    	  		      "shm_size": "2gb"
    	  		    }
    	  		  }
    	  		}`),
			Status:          "success",
			Source:          "cloud",
			StartTimestamp:  "2021-05-13T21:56:07.084Z",
			EndTimestamp:    "2021-05-13T21:59:40.919Z",
			UpdateTimestamp: "2021-05-13T21:59:46.976Z",
			BuildLog:        "[Info] ...",
			IsInvalidated:   false,
			ReleaseVersion:  nil,
			ReleaseType:     "final",
			IsPassingTests:  true,
			Contract:        "{\"name\":\"Stellarium\",\"description\":\"An educational star, constellation, planet\"}",
		},
	}
	// When
	actual, err := client.Release.GetWithQuery(context.Background(), "%24filter=key+eq+%27value%27")
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}
