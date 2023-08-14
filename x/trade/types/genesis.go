package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		TradeIndex:          TradeIndex{NextId: uint64(DefaultIndex)},
		StoredTradeList:     []StoredTrade{},
		StoredTempTradeList: []StoredTempTrade{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in storedTrade
	storedTradeIndexMap := make(map[string]struct{})

	for _, elem := range gs.StoredTradeList {
		index := string(StoredTradeKey(elem.TradeIndex))
		if _, ok := storedTradeIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for storedTrade")
		}
		storedTradeIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in storedTempTrade
	storedTempTradeIndexMap := make(map[string]struct{})

	for _, elem := range gs.StoredTempTradeList {
		index := string(StoredTempTradeKey(elem.TradeIndex))
		if _, ok := storedTempTradeIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for storedTempTrade")
		}
		storedTempTradeIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
