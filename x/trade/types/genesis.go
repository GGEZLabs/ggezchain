package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		TradeIndex: TradeIndex{
			NextId: uint64(DefaultIndex),
		},
		StoredTradeList:     []StoredTrade{},
		StoredTempTradeList: []StoredTempTrade{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if gs.TradeIndex.NextId <= 0 {
		return fmt.Errorf("next_id must be more than 0")
	}

	err := gs.ValidateStoredTrade()
	if err != nil {
		return err
	}

	err = gs.ValidateStoredTempTrade()
	if err != nil {
		return err
	}

	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}

func (gs GenesisState) ValidateStoredTrade() error {
	storedTradeIndexMap := make(map[string]struct{})

	for _, elem := range gs.StoredTradeList {

		if elem.TradeIndex <= 0 {
			return fmt.Errorf("trade_index must be more than 0")
		}

		// Check for duplicated index in storedTrade
		index := string(StoredTradeKey(elem.TradeIndex))
		if _, ok := storedTradeIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for storedTrade")
		}
		storedTradeIndexMap[index] = struct{}{}

		if elem.TradeType != TradeTypeBuy &&
			elem.TradeType != TradeTypeSell {
			return fmt.Errorf("trade_type must be buy or sell, trade_index: %d", elem.TradeIndex)
		}

		if elem.Amount.Amount.LTE(math.NewInt(0)) {
			return fmt.Errorf("amount must be more than 0, trade_index: %d", elem.TradeIndex)
		}

		if elem.Amount.Denom != DefaultCoinDenom {
			return fmt.Errorf("invalid denom expected: %s, got: %s, trade_index: %d", DefaultCoinDenom, elem.Amount.Denom, elem.TradeIndex)
		}

		if strings.TrimSpace(elem.Price) == "" {
			return fmt.Errorf("empty trade price not allowed, trade_index: %d", elem.TradeIndex)
		}

		coinPrice, err := strconv.ParseFloat(elem.Price, 64)
		if err != nil {
			return fmt.Errorf("invalid trade price err: %s, trade_index: %d", err, elem.TradeIndex)
		}

		if coinPrice <= 0 {
			return fmt.Errorf("price must be more than 0, trade_index: %d", elem.TradeIndex)
		}

		if _, err := sdk.AccAddressFromBech32(elem.ReceiverAddress); err != nil {
			return fmt.Errorf("invalid receiver_address for trade_index %d, address %s, error: %w", elem.TradeIndex, elem.ReceiverAddress, err)
		}

		if elem.Status != StatusProcessed &&
			elem.Status != StatusRejected &&
			elem.Status != StatusFailed &&
			elem.Status != StatusCanceled &&
			elem.Status != StatusPending {
			return fmt.Errorf("invalid status, trade_index: %d", elem.TradeIndex)
		}

		if _, err := sdk.AccAddressFromBech32(elem.Maker); err != nil {
			return fmt.Errorf("invalid maker address for trade_index %d, address %s, error: %w", elem.TradeIndex, elem.Maker, err)
		}

		if elem.Checker != "" {
			if _, err := sdk.AccAddressFromBech32(elem.Checker); err != nil {
				return fmt.Errorf("invalid checker address for trade_index %d, address %s, error: %w", elem.TradeIndex, elem.Checker, err)
			}
		}

		_, err = time.Parse(time.RFC3339, elem.CreateDate)

		if err != nil {
			return fmt.Errorf("invalid create_date format, trade_index: %d", elem.TradeIndex)
		}

		_, err = time.Parse(time.RFC3339, elem.UpdateDate)

		if err != nil {
			return fmt.Errorf("invalid update_date format, trade_index: %d", elem.TradeIndex)
		}

		_, err = time.Parse(time.RFC3339, elem.ProcessDate)

		if err != nil {
			return fmt.Errorf("invalid process_date format, trade_index: %d", elem.TradeIndex)
		}

		if err := ValidateTradeData(elem.TradeData); err != nil {
			return fmt.Errorf("invalid trade_data, error: %s, trade_index: %d", err, elem.TradeIndex)
		}

		if !json.Valid([]byte(elem.BankingSystemData)) {
			return fmt.Errorf("invalid banking_system_data JSON format, trade_index: %d", elem.TradeIndex)
		}
	}
	return nil
}

func (gs GenesisState) ValidateStoredTempTrade() error {
	storedTempTradeIndexMap := make(map[string]struct{}) // tradeIndex
	tempTradeIndexMap := make(map[int64]struct{})        // tempTradeIndex

	for _, elem := range gs.StoredTempTradeList {
		if elem.TradeIndex <= 0 {
			return fmt.Errorf("trade_index must be more than 0")
		}

		// Check for duplicated index in storedTempTrade
		index := string(StoredTempTradeKey(elem.TradeIndex))
		if _, ok := storedTempTradeIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for storedTempTrade")
		}
		storedTempTradeIndexMap[index] = struct{}{}

		if elem.TempTradeIndex <= 0 {
			return fmt.Errorf("temp_trade_index must be more than 0")
		}

		// Check for duplicated TempTradeIndex
		if _, exists := tempTradeIndexMap[int64(elem.TempTradeIndex)]; exists {
			return fmt.Errorf("duplicated temp_trade_index: %d", elem.TempTradeIndex)
		}
		tempTradeIndexMap[int64(elem.TempTradeIndex)] = struct{}{}

		_, err := time.Parse(time.RFC3339, elem.CreateDate)
		if err != nil {
			return fmt.Errorf("invalid create_date format, trade_index: %d", elem.TradeIndex)
		}

	}
	return nil
}
