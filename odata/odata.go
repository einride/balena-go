package odata

import (
	"encoding/json"
	"errors"
)

// EntityURL returns an OData entity ID URL given a base URL and an entity ID.
func EntityURL(base, id string) string {
	return base + "(" + id + ")"
}

type object struct {
	Deferred Deferred `json:"__deferred"`
	ID       int64    `json:"__id"`

	rawData json.RawMessage
}

type Object object

// As unmarshals the raw json data into the struct provided.
// This can be used in UData $expand queries where not the deferred data is
// injected but an array of the actual data itself.
func (o Object) As(d interface{}) error {
	return json.Unmarshal(o.rawData, d)
}

func (o *Object) UnmarshalJSON(data []byte) error {

	d := new(object)
	if err := json.Unmarshal(data, d); err != nil {
		var jsonErr *json.UnmarshalTypeError
		if !errors.As(err, &jsonErr) {
			return err
		}
	}
	d.rawData = data

	*o = *(*Object)(d)

	return nil
}

func New(v interface{}) (*Object, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return &Object{
		rawData: data,
	}, nil
}

type Deferred struct {
	URI string `json:"uri"`
}
