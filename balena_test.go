package balena

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewRequest(t *testing.T) {
	// given
	c := New(nil, "")
	inURL, outURL := "foo", defaultBaseURL+"foo"
	inBody, outBody := &struct {
		Message string
	}{
		Message: "hello world",
	},
		`{"Message":"hello world"}`+"\n"
	// when
	req, _ := c.NewRequest(context.Background(), http.MethodGet, inURL, "", inBody)
	// then
	// test relative URL was expanded
	assert.Equal(t, req.URL.String(), outURL)
	// test body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	assert.Equal(t, string(body), outBody)

	// test default user-agent is attached to the request
	userAgent := req.Header.Get("User-Agent")
	if c.UserAgent != userAgent {
		t.Errorf("NewRequest() User-Agent = %v, expected %v", userAgent, c.UserAgent)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	// given
	c := New(nil, "")
	// when
	_, err := c.NewRequest(context.Background(), http.MethodGet, ":", "", nil)
	// then
	assert.Assert(t, err != nil)
}

func TestNewRequest_Query(t *testing.T) {
	// given
	c := New(nil, "")
	// when
	req, err := c.NewRequest(context.Background(), http.MethodGet, "foo", "key=value", nil)
	// then
	assert.NilError(t, err)
	assert.Equal(t, "key=value", req.URL.RawQuery)
}

func TestDo(t *testing.T) {
	// given
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := New(nil, "")
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url
	defer server.Close()
	type foo struct {
		A string
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})
	req, err := client.NewRequest(context.Background(), http.MethodGet, "/", "", nil)
	assert.NilError(t, err)
	body := &foo{}
	// when
	err = client.Do(req, body)
	// then
	assert.NilError(t, err)
	expected := &foo{"a"}
	assert.DeepEqual(t, expected, body)
}

func newFixture() (*Client, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	client := New(nil, "")
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url

	return client, mux, func() {
		server.Close()
	}
}

// supervisorV2Fixture returns a SupervisorService as if it was used locally inside balena.
// Caller is expected to set the following environmebt variables before
//
// BALENA_SUPERVISOR_API_KEY
// BALENA_APP_ID
// BALENA_DEVICE_UUID.
func supervisorV2Fixture(t *testing.T) (*SupervisorV2Service, *http.ServeMux) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client, err := NewSupervisorV2(nil)
	assert.NilError(t, err)
	url, _ := url.Parse(server.URL + "/")
	client.client.BaseURL = url
	t.Cleanup(func() {
		server.Close()
	})
	return client, mux
}

// supervisorV1Fixture returns a SupervisorService as if it was used locally inside balena.
// Caller is expected to set the following environmebt variables before
//
// BALENA_SUPERVISOR_API_KEY
// BALENA_APP_ID
// BALENA_DEVICE_UUID.
func supervisorV1Fixture(t *testing.T) (*SupervisorV1Service, *http.ServeMux) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client, err := NewSupervisorV1(nil)
	assert.NilError(t, err)
	url, _ := url.Parse(server.URL + "/")
	client.client.BaseURL = url
	t.Cleanup(func() {
		server.Close()
	})
	return client, mux
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}
