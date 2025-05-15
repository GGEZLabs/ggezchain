package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// AclAdminKeyPrefix is the prefix to retrieve all AclAdmin
	AclAdminKeyPrefix = "AclAdmin/value/"
)

// AclAdminKey returns the store key to retrieve a AclAdmin from the index fields
func AclAdminKey(
	address string,
) []byte {
	var key []byte

	addressBytes := []byte(address)
	key = append(key, addressBytes...)
	key = append(key, []byte("/")...)

	return key
}
