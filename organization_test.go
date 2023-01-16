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
	organizationGetResponse = `{
  "d": [
    {
      "id": 241261,
      "created_at": "2022-01-27T15:25:20.412Z",
      "name": "SGG",
      "handle": "sgg",
      "company_name": null,
      "website": null,
      "industry": null,
      "billing_account_code": "4c8a1ef1-986c-448a-8829-55395952e804",
      "has_past_due_invoice_since__date": null,
      "__metadata": {
        "uri": "/resin/organization(@id)?@id=241261"
      }
    }
  ]
}`
)

func TestOrganizationService_Get(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	entityID := int64(241261)

	mux.HandleFunc(
		"/"+odata.EntityURL(organizationBasePath, strconv.FormatInt(entityID, 10)),
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			fmt.Fprint(w, organizationGetResponse)
		},
	)
	expected := &OrganizationResponse{
		ID:                 entityID,
		CreatedAt:          "2022-01-27T15:25:20.412Z",
		Name:               "SGG",
		Handle:             "sgg",
		CompanyName:        "",
		BillingAccountCode: "4c8a1ef1-986c-448a-8829-55395952e804",
	}
	// When
	actual, err := client.Organization.Get(context.Background(), entityID)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)

}

func TestOrganizationService_GetWithQuery(t *testing.T) {
	// Given
	client, mux, cleanup := newFixture()
	defer cleanup()
	query := "%24filter=handle%20eq%20%27sgg%27"

	mux.HandleFunc(
		"/"+organizationBasePath,
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodGet)
			if r.URL.RawQuery != query {
				http.Error(w, fmt.Sprintf("query = %s, expected %s", r.URL.RawQuery, query), 500)
				return
			}
			fmt.Fprint(w, organizationGetResponse)
		},
	)
	expected := []*OrganizationResponse{
		{
			ID:                 241261,
			CreatedAt:          "2022-01-27T15:25:20.412Z",
			Name:               "SGG",
			Handle:             "sgg",
			CompanyName:        "",
			BillingAccountCode: "4c8a1ef1-986c-448a-8829-55395952e804",
		},
	}
	// When
	actual, err := client.Organization.GetWithQuery(context.Background(), query)
	// Then
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)

}
