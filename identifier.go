package balena

import "strconv"

// IDOrUUID represents an ID which can be an Entity ID or an UUID.
type IDOrUUID struct {
	id     string
	isUUID bool
}

func DeviceUUID(uuid string) IDOrUUID {
	return IDOrUUID{
		id:     uuid,
		isUUID: true,
	}
}

func DeviceID(id int64) IDOrUUID {
	return IDOrUUID{
		id:     strconv.FormatInt(id, 10),
		isUUID: false,
	}
}
