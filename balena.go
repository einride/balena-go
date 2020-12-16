package balena

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	defaultBaseURL = "https://api.balena-cloud.com/"
	userAgent      = "einride/go-balena"
)

type Client struct {
	client *http.Client // HTTP client used to communicate with the API.

	// Base URL for API requests. Defaults to the public Balena Cloud API
	// BaseURL should always be specified with a trailing slash.
	BaseURL   *url.URL
	UserAgent string
	authToken string

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the Balena API
	Application  *ApplicationService
	Device       *DeviceService
	Release      *ReleaseService
	ReleaseTag   *ReleaseTagService
	DeviceEnvVar *DeviceEnvVarService
	DeviceTag    *DeviceTagService
}

type service struct {
	client *Client
}

// An ErrorResponse reports the error caused by an API request.
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode)
}

// New returns a new Balena API client.
func New(httpClient *http.Client, authToken string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent, authToken: authToken}
	c.common.client = c
	c.Application = (*ApplicationService)(&c.common)
	c.Device = (*DeviceService)(&c.common)
	c.Release = (*ReleaseService)(&c.common)
	c.DeviceEnvVar = (*DeviceEnvVarService)(&c.common)
	c.DeviceTag = (*DeviceTagService)(&c.common)
	c.ReleaseTag = (*ReleaseTagService)(&c.common)
	return c
}

// SupervisorV2 returns a SupervisorV2Service to be used with balena cloud.Supervisor
// For local supervisor communication use balena.NewSupervisorV2() instead.
func (c *Client) SupervisorV2(applicationID int64, deviceUUID string) *SupervisorV2Service {
	return &SupervisorV2Service{
		service:    c.common,
		deviceUUID: deviceUUID,
		appID:      strconv.FormatInt(applicationID, 10),
	}
}

// NewSupervisorV2 returns a new supervisor v2 API client.
// This is only meant to be used within balena where the application can talk to
// the supervisor directly. For talking to the supervisor through balene cloud
// use the balena API client through the client.SupervisorV2() method.
func NewSupervisorV2(httpClient *http.Client) (*SupervisorV2Service, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	addr := os.Getenv("BALENA_SUPERVISOR_ADDRESS")
	baseURL, err := url.Parse(addr + "/")
	if err != nil {
		return nil, fmt.Errorf("unable to find supervisor address: %v", err)
	}
	key := os.Getenv("BALENA_SUPERVISOR_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("unable to retrieve supervisor API key. BALENA_SUPERVISOR_API_KEY not set?")
	}
	uuid := os.Getenv("BALENA_DEVICE_UUID")
	if uuid == "" {
		return nil, fmt.Errorf("unable to retrieve supervisor API key. BALENA_DEVICE_UUID not set?")
	}
	appID := os.Getenv("BALENA_APP_ID")
	if appID == "" {
		return nil, fmt.Errorf("unable to retrieve balena application ID. BALENA_APP_ID not set?")
	}
	c := Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	svc := service{
		client: &c,
	}

	return &SupervisorV2Service{
		service:    svc,
		deviceUUID: uuid,
		apiKey:     key,
		appID:      appID,
		local:      true,
	}, nil
}

// SupervisorV1 returns a SupervisorV1Service to be used with balena cloud.Supervisor
// For local supervisor communication use balena.NewSupervisorV1() instead.
func (c *Client) SupervisorV1(applicationID int64, deviceUUID string) *SupervisorV1Service {
	return &SupervisorV1Service{
		service:    c.common,
		deviceUUID: deviceUUID,
		appID:      strconv.FormatInt(applicationID, 10),
	}
}

// NewSupervisorV1 returns a new supervisor v1 API client.
// This is only meant to be used within balena where the application can talk to
// the supervisor directly. For talking to the supervisor through balene cloud
// use the balena API client through the client.SupervisorV1() method.
func NewSupervisorV1(httpClient *http.Client) (*SupervisorV1Service, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	addr := os.Getenv("BALENA_SUPERVISOR_ADDRESS")
	baseURL, err := url.Parse(addr + "/")
	if err != nil {
		return nil, fmt.Errorf("unable to find supervisor address: %v", err)
	}
	key := os.Getenv("BALENA_SUPERVISOR_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("unable to retrieve supervisor API key. BALENA_SUPERVISOR_API_KEY not set?")
	}
	uuid := os.Getenv("BALENA_DEVICE_UUID")
	if uuid == "" {
		return nil, fmt.Errorf("unable to retrieve supervisor API key. BALENA_DEVICE_UUID not set?")
	}
	appID := os.Getenv("BALENA_APP_ID")
	if appID == "" {
		return nil, fmt.Errorf("unable to retrieve balena application ID. BALENA_APP_ID not set?")
	}
	c := Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	svc := service{
		client: &c,
	}

	return &SupervisorV1Service{
		service:    svc,
		deviceUUID: uuid,
		apiKey:     key,
		appID:      appID,
		local:      true,
	}, nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
// A raw query string can be specified by rawQuery.
func (c *Client) NewRequest(
	ctx context.Context,
	method, urlStr, rawQuery string,
	body interface{},
) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if len(rawQuery) > 0 {
		u.RawQuery = rawQuery
	}
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// req.Header.Set("Accept", mediaTypeV3)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	if c.authToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.authToken)
	}
	return req, nil
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(req *http.Request, v interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()
	err = CheckResponse(resp)
	if err != nil {
		return err
	}
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return err
			}
		}
	}
	return err
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	return errorResponse
}
