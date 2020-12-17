package balena

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/einride/balena-go/odata"
)

const releaseTagBasePath = "v5/release_tag"

// ReleaseTag handles communication with the release tag related methods of the
// Balena Cloud API.
type ReleaseTagService service

type ReleaseTagResponse struct {
	ID      int64        `json:"id,omitempty"`
	Release odata.Object `json:"release,omitempty"`
	TagKey  string       `json:"tag_key,omitempty"`
	Value   string       `json:"value,omitempty"`
}

// List lists all release tags for a given release ID.
func (s *ReleaseTagService) List(ctx context.Context, releaseID int64) ([]*ReleaseTagResponse, error) {
	query := "%24filter=release/id+eq+%27" + strconv.FormatInt(releaseID, 10) + "%27"
	return s.GetWithQuery(ctx, query)
}

// ListByCommit lists all release tags for a given release commit.
func (s *ReleaseTagService) ListByCommit(ctx context.Context, commit string) ([]*ReleaseTagResponse, error) {
	query := "%24filter=release/commit+eq+%27" + commit + "%27"
	return s.GetWithQuery(ctx, query)
}

// GetWithQuery allows querying for release tags using a custom Open Data Protocol query.
// The query should be a valid, escaped OData query such as `%24filter=uuid+eq+%2712333422%27`.
//
// Forward slash in filter keys should not be escaped (So `device/uuid` should not be escaped).
func (s *ReleaseTagService) GetWithQuery(ctx context.Context, query string) ([]*ReleaseTagResponse, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, releaseTagBasePath, query, nil)
	if err != nil {
		return nil, fmt.Errorf("get release tag with query NewRequest: %v", err)
	}
	type Response struct {
		D []*ReleaseTagResponse `json:"d,omitempty"`
	}
	resp := &Response{}
	err = s.client.Do(req, resp)
	if err != nil {
		return nil, fmt.Errorf("get release tag with query: %v", err)
	}
	return resp.D, nil
}
