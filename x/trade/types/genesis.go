package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		TradeIndex: TradeIndex{
			NextId: DefaultIndex,
		},
		StoredTrades:     []StoredTrade{},
		StoredTempTrades: []StoredTempTrade{},
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

	for _, elem := range gs.StoredTrades {
		if elem.TradeIndex <= 0 {
			return fmt.Errorf("trade_index must be more than 0")
		}

		// Check for duplicated index in storedTrade
		index := string(StoredTradeKey(elem.TradeIndex))
		if _, ok := storedTradeIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for storedTrade")
		}
		storedTradeIndexMap[index] = struct{}{}

		if !elem.TradeType.IsTypeValid() {
			return fmt.Errorf("invalid trade_type, trade_index: %d", elem.TradeIndex)
		}

		if elem.TradeType == TradeTypeBuy ||
			elem.TradeType == TradeTypeSell {
			if !elem.Amount.IsValid() {
				return fmt.Errorf("invalid amount: %s, trade_index: %d", elem.Amount.String(), elem.TradeIndex)
			}

			if elem.Amount.IsZero() {
				return fmt.Errorf("zero amount not allowed: %s, trade_index: %d", elem.Amount.String(), elem.TradeIndex)
			}

			if elem.Amount.Denom != DefaultDenom {
				return fmt.Errorf("invalid denom expected: %s, got: %s, trade_index: %d", DefaultDenom, elem.Amount.Denom, elem.TradeIndex)
			}

			if _, err := sdk.AccAddressFromBech32(elem.ReceiverAddress); err != nil {
				return fmt.Errorf("invalid receiver_address for trade_index %d, address %s, error: %w", elem.TradeIndex, elem.ReceiverAddress, err)
			}
		} else {
			// {"amount":"0","denom":""} Accepted
			// {"amount":"1000","denom":"xxx"} Not Accepted
			if elem.Amount != nil &&
				(elem.Amount.IsValid() &&
					(!elem.Amount.Amount.IsZero() || elem.Amount.Denom != "")) {
				return fmt.Errorf("amount must not be set for trade type: %s, trade_index: %d", elem.TradeType.String(), elem.TradeIndex)
			}

			if elem.ReceiverAddress != "" {
				return fmt.Errorf("receiver_address must not be set for trade type %s, trade_index %d", elem.TradeType.String(), elem.TradeIndex)
			}
		}

		if strings.TrimSpace(elem.CoinMintingPriceUsd) == "" {
			return fmt.Errorf("empty trade price not allowed, trade_index: %d", elem.TradeIndex)
		}

		coinPrice, err := strconv.ParseFloat(elem.CoinMintingPriceUsd, 64)
		if err != nil {
			return fmt.Errorf("invalid trade price err: %s, trade_index: %d", err, elem.TradeIndex)
		}

		if coinPrice <= 0 {
			return fmt.Errorf("price must be more than 0, trade_index: %d", elem.TradeIndex)
		}

		if !elem.Status.IsStatusValid() {
			return fmt.Errorf("invalid status, trade_index: %d", elem.TradeIndex)
		}

		if _, err := sdk.AccAddressFromBech32(elem.Maker); err != nil {
			return fmt.Errorf("invalid maker address for trade_index %d, address %s, error: %w", elem.TradeIndex, elem.Maker, err)
		}

		if elem.Checker != "" && (elem.Status != StatusPending && elem.Status != StatusCanceled) {
			if _, err := sdk.AccAddressFromBech32(elem.Checker); err != nil {
				return fmt.Errorf("invalid checker address for trade_index %d, address %s, error: %w", elem.TradeIndex, elem.Checker, err)
			}
		} else {
			return fmt.Errorf("checker must not be set for trade status %s, trade_index %d", elem.Status.String(), elem.TradeIndex)
		}

		_, err = time.Parse(time.RFC3339, elem.CreateDate)
		if err != nil {
			return fmt.Errorf("invalid create_date format, trade_index: %d", elem.TradeIndex)
		}

		_, err = time.Parse(time.RFC3339, elem.TxDate)
		if err != nil {
			return fmt.Errorf("invalid tx_date format, trade_index: %d", elem.TradeIndex)
		}

		_, err = time.Parse(time.RFC3339, elem.UpdateDate)
		if err != nil {
			return fmt.Errorf("invalid update_date format, trade_index: %d", elem.TradeIndex)
		}

		_, err = time.Parse(time.RFC3339, elem.ProcessDate)
		if err != nil {
			return fmt.Errorf("invalid process_date format, trade_index: %d", elem.TradeIndex)
		}

		_, err = ValidateTradeData(elem.TradeData)
		if err != nil {
			return fmt.Errorf("invalid trade_data, error: %s, trade_index: %d", err, elem.TradeIndex)
		}

		if !json.Valid([]byte(elem.BankingSystemData)) {
			return fmt.Errorf("invalid banking_system_data json format, trade_index: %d", elem.TradeIndex)
		}

		err = ValidateCoinMintingPriceJson(elem.CoinMintingPriceJson)
		if err != nil {
			return fmt.Errorf("invalid coin_minting_price_json, error: %s, trade_index: %d", err, elem.TradeIndex)
		}

		err = ValidateExchangeRateJson(elem.ExchangeRateJson)
		if err != nil {
			return fmt.Errorf("invalid exchange_rate_json, error: %s, trade_index: %d", err, elem.TradeIndex)
		}
	}
	return nil
}

func (gs GenesisState) ValidateStoredTempTrade() error {
	storedTradeIndexMap := make(map[string]struct{}) // tradeIndex

	for _, elem := range gs.StoredTempTrades {
		if elem.TradeIndex <= 0 {
			return fmt.Errorf("trade_index must be more than 0")
		}

		// Check for duplicated index in storedTempTrade
		index := string(StoredTempTradeKey(elem.TradeIndex))
		if _, ok := storedTradeIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for storedTempTrade")
		}
		storedTradeIndexMap[index] = struct{}{}

		_, err := time.Parse(time.RFC3339, elem.TxDate)
		if err != nil {
			return fmt.Errorf("invalid create_date format, trade_index: %d", elem.TradeIndex)
		}
	}
	return nil
}
