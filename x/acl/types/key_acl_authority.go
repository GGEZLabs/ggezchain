package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// AclAuthorityKeyPrefix is the prefix to retrieve all AclAuthority
	AclAuthorityKeyPrefix = "AclAuthority/value/"
)

// AclAuthorityKey returns the store key to retrieve a AclAuthority from the index fields
func AclAuthorityKey(
	address string,
) []byte {
	var key []byte

	addressBytes := []byte(address)
	key = append(key, addressBytes...)
	key = append(key, []byte("/")...)

	return key
}
