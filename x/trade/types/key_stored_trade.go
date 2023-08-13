package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// StoredTradeKeyPrefix is the prefix to retrieve all StoredTrade
	StoredTradeKeyPrefix = "StoredTrade/value/"
)

// StoredTradeKey returns the store key to retrieve a StoredTrade from the index fields
func StoredTradeKey(
	tradeIndex string,
) []byte {
	var key []byte

	tradeIndexBytes := []byte(tradeIndex)
	key = append(key, tradeIndexBytes...)
	key = append(key, []byte("/")...)

	return key
}
