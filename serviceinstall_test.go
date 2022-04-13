package balena

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"go.einride.tech/balena/odata"
	"gotest.tools/v3/assert"
)

func TestServiceInstallService_List(t *testing.T) {
	const actualJSON = `{
    "d": [
        {
            "installs__service": [
                {
                    "service_name": "main",
                    "application": {
                        "__id": 1660000,
                        "__deferred": {
                            "uri": "/resin/application(@id)?@id=1660000"
                        }
                    },
                    "created_at": "2020-06-18T20:26:23.470Z",
                    "id": 500000,
                    "__metadata": {
                        "uri": "/resin/service(@id)?@id=500000"
                    }
                }
            ],
            "id": 7180000,
            "created_at": "2021-10-15T08:06:13.353Z",
            "device": {
                "__id": 5000000,
                "__deferred": {
                    "uri": "/resin/device(@id)?@id=5000000"
                }
            },
            "__metadata": {
                "uri": "/resin/service_install(@id)?@id=7180000"
            }
        },
        {
            "installs__service": [
                {
                    "service_name": "other_service",
                    "application": {
                        "__id": 1790000,
                        "__deferred": {
                            "uri": "/resin/application(@id)?@id=1790000"
                        }
                    },
                    "created_at": "2021-01-11T11:46:48.175Z",
                    "id": 860000,
                    "__metadata": {
                        "uri": "/resin/service(@id)?@id=860000"
                    }
                }
            ],
            "id": 8580000,
            "created_at": "2021-11-24T14:03:30.793Z",
            "device": {
                "__id": 5040000,
                "__deferred": {
                    "uri": "/resin/device(@id)?@id=5040000"
                }
            },
            "__metadata": {
                "uri": "/resin/service_install(@id)?@id=8580000"
            }
        }
    ]
}`
	doTest := func(tt *testing.T, expectedQuery string, deviceID IDOrUUID) {
		client, mux, cleanup := newFixture()
		defer cleanup()
		mux.HandleFunc("/"+serviceInstallBasePath, func(w http.ResponseWriter, r *http.Request) {
			testMethod(tt, r, http.MethodGet)
			if r.URL.RawQuery != expectedQuery {
				tt.Errorf("query = %s ; expected %s", r.URL.RawQuery, expectedQuery)
				http.Error(w, "query string mismatch", http.StatusInternalServerError)
				return
			}
			_, _ = fmt.Fprint(w, actualJSON)
		})

		expected := ServiceInstalls{
			{
				InstallsService: []InstallsService{
					{
						ID:          500_000,
						ServiceName: "main",
						Application: odata.Object{
							Deferred: odata.Deferred{
								URI: "/resin/application(@id)?@id=1660000",
							},
							ID: 1660_000,
						},
						CreatedAt: time.Date(2020, 6, 18, 20, 26, 23, 470_000_000, time.UTC),
					},
				},
				ID:        7180_000,
				CreatedAt: time.Date(2021, 10, 15, 8, 6, 13, 353_000_000, time.UTC),
				Device: odata.Object{
					Deferred: odata.Deferred{
						URI: "/resin/device(@id)?@id=5000000",
					},
					ID: 5000000,
				},
			},
			{
				InstallsService: []InstallsService{
					{
						ID:          860_000,
						ServiceName: "other_service",
						Application: odata.Object{
							Deferred: odata.Deferred{
								URI: "/resin/application(@id)?@id=1790000",
							},
							ID: 1790_000,
						},
						CreatedAt: time.Date(2021, 1, 11, 11, 46, 48, 175_000_000, time.UTC),
					},
				},
				ID:        8580_000,
				CreatedAt: time.Date(2021, 11, 24, 14, 3, 30, 793_000_000, time.UTC),
				Device: odata.Object{
					Deferred: odata.Deferred{
						URI: "/resin/device(@id)?@id=5040000",
					},
					ID: 5040_000,
				},
			},
		}

		actual, err := client.ServiceInstall.List(context.Background(), deviceID)
		assert.NilError(tt, err)
		assert.DeepEqual(tt, actual, expected)
	}

	t.Run("device ID", func(t *testing.T) {
		deviceID := DeviceID(5000000)
		//nolint:lll
		expected := "%24filter=device+eq+%275000000%27&%24expand=installs__service(%24select=service_name,application,created_at,id)"
		doTest(t, expected, deviceID)
	})
	t.Run("device UUID", func(t *testing.T) {
		deviceID := DeviceUUID("apabepa123")
		//nolint:lll
		expected := "%24filter=device/uuid+eq+%27apabepa123%27&%24expand=installs__service(%24select=service_name,application,created_at,id)"
		doTest(t, expected, deviceID)
	})
}
