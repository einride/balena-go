package odata

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
)

type testData struct {
	ID   int64  `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

const (
	testMessage = `[
		{
			"id": 5,
			"slug": "raspberrypi3-64",
			"name": "Raspberry Pi 3 (using 64bit OS)"
		}
	]`
	deferredMessage = `{
		"__id": 124474,
		"__deferred": {
			"uri": "/resin/service_instance(@id)?@id=124474"
		}
	}`
)

func TestObject_As(t *testing.T) {
	type fields struct {
		Deferred Deferred
		ID       int64
		rawData  json.RawMessage
	}
	type args struct {
		d interface{}
	}
	tests := map[string]struct {
		fields fields
		args   args
		want   interface{}
	}{
		"devicetype": {
			fields: fields{rawData: json.RawMessage(testMessage)},
			args: args{
				d: new([]*testData),
			},
			want: &[]*testData{{
				ID:   5,
				Slug: "raspberrypi3-64",
				Name: "Raspberry Pi 3 (using 64bit OS)",
			}},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := Object{
				Deferred: tt.fields.Deferred,
				ID:       tt.fields.ID,
				rawData:  tt.fields.rawData,
			}
			err := o.As(tt.args.d)

			assert.NilError(t, err)
			assert.DeepEqual(t, tt.want, tt.args.d, cmp.AllowUnexported(Object{}))
		})
	}
}

func TestObject_UnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		data []byte
		want *Object
	}{
		"Containes deferred data": {
			data: []byte(deferredMessage),
			want: &Object{
				rawData: json.RawMessage(deferredMessage),
				ID:      124474,
				Deferred: Deferred{
					URI: "/resin/service_instance(@id)?@id=124474",
				},
			},
		},
		"Contains real data": {
			data: []byte(testMessage),
			want: &Object{
				rawData: json.RawMessage(testMessage),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			o := new(Object)

			err := o.UnmarshalJSON(tt.data)

			assert.NilError(t, err)
			assert.DeepEqual(t, tt.want, o, cmp.AllowUnexported(Object{}))
		})
	}
}
