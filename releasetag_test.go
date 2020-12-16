package balena

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/einride/balena-go/odata"
	"github.com/stretchr/testify/require"
)

func TestReleaseTagService_List(t *testing.T) {
	// Given
	jsonResp := `{
		"d": [
		{ 
			"id": 92661,
				"release": {
				"__deferred": {
					"uri": "/resin/release(1309764)"
				},
				"__id": 1309764
			},
			"tag_key": "GITCOMMIT",
				"value": "3eca921e92609537d17f8689fa964f88edb13fd7",
				"__metadata": {
				"uri": "/resin/release_tag(@id)?@id=92661"
			}
		}
	]
}`
	releaseID := int64(1309764)
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+releaseTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=release/id+eq+%27" + strconv.FormatInt(releaseID, 10) + "%27"
		if r.URL.RawQuery != expected {
			fmt.Printf("query = %s ; expected = %s\n", r.URL.RawQuery, expected)
			http.Error(w, fmt.Sprintf("query = %s ; expected = %s\n", r.URL.RawQuery, expected), 500)
		}
		fmt.Fprint(w, jsonResp)
	})
	expected := []*ReleaseTagResponse{
		{
			ID: 92661,
			Release: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/release(1309764)"},
				ID:       1309764,
			},
			TagKey: "GITCOMMIT",
			Value:  "3eca921e92609537d17f8689fa964f88edb13fd7",
		},
	}
	// When
	actual, err := client.ReleaseTag.List(context.Background(), releaseID)
	// Then
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestReleaseTagService_ListByCommit(t *testing.T) {
	// Given
	jsonResp := `{
		"d": [
		{ 
			"id": 92661,
				"release": {
				"__deferred": {
					"uri": "/resin/release(1309764)"
				},
				"__id": 1309764
			},
			"tag_key": "GITCOMMIT",
				"value": "3eca921e92609537d17f8689fa964f88edb13fd7",
				"__metadata": {
				"uri": "/resin/release_tag(@id)?@id=92661"
			}
		}
	]
}`
	releaseCommit := "37a9eff78e46f83f591dd34ee6e4b5ce"
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+releaseTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=release/commit+eq+%27" + releaseCommit + "%27"
		if r.URL.RawQuery != expected {
			fmt.Printf("query = %s ; expected = %s\n", r.URL.RawQuery, expected)
			http.Error(w, fmt.Sprintf("query = %s ; expected = %s\n", r.URL.RawQuery, expected), 500)
		}
		fmt.Fprint(w, jsonResp)
	})
	expected := []*ReleaseTagResponse{
		{
			ID: 92661,
			Release: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/release(1309764)"},
				ID:       1309764,
			},
			TagKey: "GITCOMMIT",
			Value:  "3eca921e92609537d17f8689fa964f88edb13fd7",
		},
	}
	// When
	actual, err := client.ReleaseTag.ListByCommit(context.Background(), releaseCommit)
	// Then
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestReleaseTagService_GetWithQuery(t *testing.T) {
	// Given
	jsonResp := `{
		"d": [
		{ 
			"id": 92661,
				"release": {
				"__deferred": {
					"uri": "/resin/release(1309764)"
				},
				"__id": 1309764
			},
			"tag_key": "GITCOMMIT",
				"value": "3eca921e92609537d17f8689fa964f88edb13fd7",
				"__metadata": {
				"uri": "/resin/release_tag(@id)?@id=92661"
			}
		}
	]
}`
	releaseID := int64(1309764)
	client, mux, cleanup := newFixture()
	defer cleanup()
	mux.HandleFunc("/"+releaseTagBasePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		expected := "%24filter=release/commit+eq+%27" + strconv.FormatInt(releaseID, 10) + "%27"
		if r.URL.RawQuery != expected {
			fmt.Printf("query = %s ; expected = %s\n", r.URL.RawQuery, expected)
			http.Error(w, fmt.Sprintf("query = %s ; expected = %s\n", r.URL.RawQuery, expected), 500)
		}
		fmt.Fprint(w, jsonResp)
	})
	expected := []*ReleaseTagResponse{
		{
			ID: 92661,
			Release: odata.Object{
				Deferred: odata.Deferred{URI: "/resin/release(1309764)"},
				ID:       1309764,
			},
			TagKey: "GITCOMMIT",
			Value:  "3eca921e92609537d17f8689fa964f88edb13fd7",
		},
	}
	// When
	query := "%24filter=release/commit+eq+%271309764%27"
	actual, err := client.ReleaseTag.GetWithQuery(context.Background(), query)
	// Then
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
