package types

import "cosmossdk.io/collections"

// StoredTradeKey is the prefix to retrieve all StoredTrade
var StoredTradeKey = collections.NewPrefix("storedTrade/value/")
