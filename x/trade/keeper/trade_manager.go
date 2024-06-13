package keeper

import (
	"context"
	"encoding/json"
	"strings"

	"os"
	"time"

	errors "cosmossdk.io/errors"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type TradeDataObject struct {
	TradeData struct {
		TradeRequestID int     `json:"tradeRequestID"`
		AssetHolderID  int     `json:"assetHolderID"`
		AssetID        int     `json:"assetID"`
		TradeType      string  `json:"tradeType"`
		TradeValue     float64 `json:"tradeValue"`
		Currency       string  `json:"currency"`
		Exchange       string  `json:"exchange"`
		FundName       string  `json:"fundName"`
		Issuer         string  `json:"issuer"`
		NoShares       string  `json:"noShares"`
		Price          string  `json:"price"`
		Quantity       string  `json:"quantity"`
		Segment        string  `json:"segment"`
		SharePrice     string  `json:"sharePrice"`
		Ticker         string  `json:"ticker"`
		TradeFee       string  `json:"tradeFee"`
		TradeNetPrice  string  `json:"tradeNetPrice"`
		TradeNetValue  string  `json:"tradeNetValue"`
	} `json:"TradeData"`
	Brokerage struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Country string `json:"country"`
	} `json:"Brokerage"`
}

type ChainACL struct {
	Trade struct {
		AllowedMakers []struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"allowed_makers"`
		AllowedCheckers []struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"allowed_checkers"`
	} `json:"trade"`
}

func (k Keeper) IsAddressAllowed(goCtx context.Context, address string, msgType string) (isAllowed bool, err error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	isAddressWhitelisted, err := k.IsAddressWhitelisted(address, msgType, types.ACLFilePath)

	if err != nil {
		return false, err
	}

	if isAddressWhitelisted {
		isAddressLinkedToValidator, err := k.IsAddressLinkedToValidator(ctx, address)
		if err != nil {
			return false, err
		}
		return isAddressLinkedToValidator, nil
	}

	return false, nil
}

func (k Keeper) IsAddressWhitelisted(address string, msgType string, ACLFilePath string) (isAddressWhitelisted bool, err error) {
	userHomeDir, _ := os.UserHomeDir()
	isWhitelisted := false
	err = nil

	file, err := os.ReadFile(userHomeDir + ACLFilePath)
	if err != nil {
		return isWhitelisted, types.ErrInvalidPath
	}

	// Unmarshal the chain ACL into a struct
	var chainACL ChainACL
	if err := json.Unmarshal(file, &chainACL); err != nil {
		return isWhitelisted, err
	}

	if msgType == types.CreateTrade {
		// Loop through the validators for Allowed Makers
		for _, maker := range chainACL.Trade.AllowedMakers {
			if address == maker.Address {
				isWhitelisted = true
				break
			}
		}
	} else if msgType == types.ProcessTrade {
		// Loop through the validators for Allowed Checkers
		for _, checker := range chainACL.Trade.AllowedCheckers {
			if address == checker.Address {
				isWhitelisted = true
				break
			}
		}
	} else {
		isWhitelisted = false
	}

	return isWhitelisted, err
}

func (k Keeper) MintOrBurnCoins(ctx sdk.Context, tradeData types.StoredTrade, coin sdk.Coin) (status string, err error) {
	receiverAddress, err := sdk.AccAddressFromBech32(tradeData.ReceiverAddress)
	if err != nil {
		return types.Failed, types.ErrInvalidReceiverAddress
	}

	if tradeData.TradeType == types.Buy {
		err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
		if err != nil {
			return types.Failed, err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiverAddress, sdk.NewCoins(coin))
		if err != nil {
			errBurnCoins := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
			if errBurnCoins != nil {
				return types.CoinsStuckOnModule, err
			}

			return types.Failed, err
		}
	} else if tradeData.TradeType == types.Sell {
		// Try to send coins from account to module
		err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, receiverAddress, types.ModuleName, sdk.NewCoins(coin))
		if err != nil {
			// returns error
			return types.Failed, err
		}

		// Try to burn coins
		err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
		if err != nil {
			// Try to send coins from module to account
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiverAddress, sdk.NewCoins(coin))
			if err != nil {
				// returns error
				return types.CoinsStuckOnModule, err
			}
		}
	}

	return types.Completed, types.ErrTradeProcessedSuccessfully
}

func (k Keeper) IsAddressLinkedToValidator(goCtx context.Context, address string) (isAddressLinkedToValidator bool, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return false, err
	}
	isLinked := false

	chainValidators := k.stakingKeeper.GetAllValidators(ctx)

	// Loop through the validators
	for _, validator := range chainValidators {
		valAddress, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		if err != nil {
			return false, err
		}
		_, isLinked = k.stakingKeeper.GetDelegation(ctx, accAddress, valAddress)
		if isLinked {
			break
		}
	}

	return isLinked, nil
}

func (k Keeper) CancelExpiredPendingTrades(goCtx context.Context) (err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	allStoredTempTrade := k.GetAllStoredTempTrade(ctx)
	status := types.Canceled

	currentDate := time.Now()

	for i := 0; i < len(allStoredTempTrade); i++ {
		createDate := allStoredTempTrade[i].CreateDate
		formattedCreateDate, err := time.Parse("2006-01-02 15:04", createDate)
		if err != nil {
			return sdkErrors.Wrapf(types.ErrInvalidDateFormat, types.ErrInvalidDateFormat.Error())
		}
		differenceTime := currentDate.Sub(formattedCreateDate)
		totalDays := int(differenceTime.Hours() / 24)

		if totalDays >= 1 {

			storedTrade, _ := k.GetStoredTrade(ctx, allStoredTempTrade[i].TempTradeIndex)
			storedTrade.Status = status
			storedTrade.UpdateDate = currentDate.Format("2006-01-02 15:04")

			k.SetStoredTrade(ctx, storedTrade)
			k.RemoveStoredTempTrade(ctx, allStoredTempTrade[i].TempTradeIndex)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(types.CancelExpiredPendingTradesEventType),
			)
		}
	}
	return err
}

func (k Keeper) ValidateTradeData(tradeData string) (err error) {

	isJson := k.IsJSON(tradeData)

	if isJson {
		var data TradeDataObject
		err = json.Unmarshal([]byte(tradeData), &data)
		if err != nil {
			return errors.Wrap(types.ErrInvalidTradeDataObject, "Invalid Trade Data Object")
		}
		tradeData := data.TradeData
		if tradeData.AssetHolderID <= 0 {
			return errors.Wrap(types.ErrTradeDataAssetHolderID, "Invalid Trade Data Object")
		}
		if tradeData.AssetID <= 0 {
			return errors.Wrap(types.ErrTradeDataAssetID, "Invalid Trade Data Object")
		}
		if tradeData.TradeRequestID <= 0 {
			return errors.Wrap(types.ErrTradeDataRequestID, "Invalid Trade Data Object")
		}
		if tradeData.TradeValue == 0 {
			return errors.Wrap(types.ErrTradeDataValue, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Currency) == "" {
			return errors.Wrap(types.ErrTradeDataCurrency, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Exchange) == "" {
			return errors.Wrap(types.ErrTradeDataExchange, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.FundName) == "" {
			return errors.Wrap(types.ErrTradeDataFundName, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Issuer) == "" {
			return errors.Wrap(types.ErrTradeDataIssuer, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.NoShares) == "" {
			return errors.Wrap(types.ErrTradeDataNoShares, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Price) == "" {
			return errors.Wrap(types.ErrTradeDataPrice, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Quantity) == "" {
			return errors.Wrap(types.ErrTradeDataQuantity, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Segment) == "" {
			return errors.Wrap(types.ErrTradeDataSegment, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.SharePrice) == "" {
			return errors.Wrap(types.ErrTradeDataSharePrice, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Ticker) == "" {
			return errors.Wrap(types.ErrTradeDataTicker, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.TradeFee) == "" {
			return errors.Wrap(types.ErrTradeDataFee, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.TradeNetPrice) == "" {
			return errors.Wrap(types.ErrTradeDataNetPrice, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.TradeNetValue) == "" {
			return errors.Wrap(types.ErrTradeDataNetValue, "Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.TradeType) == "" {
			return errors.Wrap(types.ErrInvalidTradeType, "Invalid Trade Data Object")
		}

		brokerage := data.Brokerage
		if strings.TrimSpace(brokerage.Country) == "" {
			return errors.Wrap(types.ErrBrokerageCountry, "Invalid Brokerage Country Object")
		}
		if strings.TrimSpace(brokerage.Type) == "" {
			return errors.Wrap(types.ErrBrokerageType, "Invalid Brokerage Type Object")
		}
		if strings.TrimSpace(brokerage.Name) == "" {
			return errors.Wrap(types.ErrBrokerageName, "Invalid Brokerage Name Object")
		}
		return nil
	}
	return errors.Wrap(types.ErrInvalidTradeDataJSON, "Invalid Trade Data JSON")

}

func (k Keeper) IsJSON(str string) bool {
	var jsonFormat json.RawMessage
	return json.Unmarshal([]byte(str), &jsonFormat) == nil
}
