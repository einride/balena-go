package balena

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go.einride.tech/balena/odata"
)

const organizationBasePath = "v6/organization"

// OrganizationService handles communication with the organization related methods of the
// Balena Cloud API.
type OrganizationService service

type OrganizationResponse struct {
	ID                 int64  `json:"id"`
	CreatedAt          string `json:"created_at"`
	Name               string `json:"name"`
	Handle             string `json:"handle"`
	CompanyName        string `json:"company_name"`
	Website            string `json:"website"`
	Industry           string `json:"industry"`
	BillingAccountCode string `json:"billing_account_code"`
}

type OrganizationOData struct {
	D []*OrganizationResponse
	*odata.Object
}

func (d *OrganizationOData) UnmarshalJSON(data []byte) error {
	o := new(odata.Object)
	if err := json.Unmarshal(data, o); err == nil {
		d.Object = o
		return nil
	}

	return json.Unmarshal(data, &d.D)
}

// Get returns information on a single organization given its ID.
// If the organization does not exist, both the response and error are nil.
func (s *OrganizationService) Get(ctx context.Context, id int64) (*OrganizationResponse, error) {
	path := odata.EntityURL(organizationBasePath, strconv.FormatInt(id, 10))

	resp, err := s.getWithQueryAndPath(ctx, path, "")

	if err != nil {
		return nil, fmt.Errorf("unable to get organization: %v", err)
	}
	if len(resp) > 1 {
		return nil, errors.New("received more than 1 organization, expected 0 or 1")
	}
	if len(resp) == 0 {
		return nil, nil
	}
	return resp[0], nil

}

// GetWithQuery allows querying for organizations using a custom open data protocol query.
// The query should be a valid, escaped OData query such as `%24filter=uuid+eq+'12333422'`
//
// Forward slash in filter keys should not be escaped (So `organization/handle` should not be escaped).
func (s *OrganizationService) GetWithQuery(ctx context.Context, query string) ([]*OrganizationResponse, error) {
	return s.getWithQueryAndPath(ctx, organizationBasePath, query)
}

func (s *OrganizationService) getWithQueryAndPath(
	ctx context.Context,
	path string,
	query string,
) ([]*OrganizationResponse, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create get request: %v", err)
	}
	type Response struct {
		D []*OrganizationResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to get application list: %v", err)
	}

	return resp.D, nil
}
