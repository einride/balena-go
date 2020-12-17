package balena_test

import (
	"context"

	"go.einride.tech/balena"
)

func ExampleDeviceService_Get() {
	token := "mytoken"
	// We supply a nil http client to make use of http.DefaultClient
	client := balena.New(nil, token)

	// We want a device with a given uuid
	deviceUUID := "00a859f211585e9417a676e09434cc55"

	_, _ = client.Device.Get(context.Background(), balena.DeviceUUID(deviceUUID))
}

func ExampleDeviceService_GetWithQuery() {
	token := "mytoken"
	// We supply a nil http client to make use of http.DefaultClient
	client := balena.New(nil, token)

	// We want a device with a given name
	query := "%24filter=device_name+eq+%27mydevice%27"

	_, _ = client.Device.GetWithQuery(context.Background(), query)
}
