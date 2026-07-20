package types

import "cosmossdk.io/collections"

// StoredTempTradeKey is the prefix to retrieve all StoredTempTrade
var StoredTempTradeKey = collections.NewPrefix("storedTempTrade/value/")
