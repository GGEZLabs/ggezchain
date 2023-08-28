package keeper

import (
	"context"
	"encoding/json"
	"strings"

	//"io/ioutil"
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

func (k Keeper) IsAddressAllowed(stakingKeeper types.StakingKeeper, goCtx context.Context, address string, msgType string) (isAllowed bool, err error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	isAddressWhitelisted, err := k.IsAddressWhitelisted(address, msgType)

	if err != nil {
		return false, err
	}

	if isAddressWhitelisted {
		isAddressLinkedToValidator := k.IsAddressLinkedToValidator(stakingKeeper, ctx, address)
		return isAddressLinkedToValidator, nil
	}

	return false, nil
}

func (k Keeper) IsAddressWhitelisted(address string, msgType string) (isAddressWhitelisted bool, err error) {
	userHomeDir, _ := os.UserHomeDir()
	isWhitelisted := false
	err = nil

	file, err := os.ReadFile(userHomeDir + "/.ggezchain/config/chain_acl.json")
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

/* func (k Keeper) AllStoredTempTrade(goCtx context.Context) (allStoredTempTrade []types.StoredTempTrade,err error){
	ctx := sdk.UnwrapSDKContext(goCtx)
	allStoredTrade:= k.GetAllStoredTempTrade(ctx)
	k.CancelExpiredPendingTrades(ctx,allStoredTrade)
	return allStoredTrade ,nil

} */

func (k Keeper) MintOrBurnCoins(ctx sdk.Context, tradeData types.StoredTrade, coin sdk.Coin) (status string, err error) {
	receiverAddress, errReceiverAddress := sdk.AccAddressFromBech32(tradeData.ReceiverAddress)
	if errReceiverAddress != nil {
		return types.Failed, errReceiverAddress
	}

	if tradeData.TradeType == types.Buy {
		errMintCoins := k.bank.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
		if errMintCoins != nil {
			return types.Failed, errMintCoins
		}

		errSendCoinsFromModuleToAccount := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiverAddress, sdk.NewCoins(coin))
		if errSendCoinsFromModuleToAccount != nil {
			errBurnCoins := k.bank.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
			if errBurnCoins != nil {
				return types.CoinsStuckOnModule, errSendCoinsFromModuleToAccount
			}

			return types.Failed, errSendCoinsFromModuleToAccount
		}
	} else if tradeData.TradeType == types.Sell {
		// Try to send coins from account to module
		errSendCoinsFromAccountToModule := k.bank.SendCoinsFromAccountToModule(ctx, receiverAddress, types.ModuleName, sdk.NewCoins(coin))
		if errSendCoinsFromAccountToModule != nil {
			// returns error
			return types.Failed, errSendCoinsFromAccountToModule
		}

		// Try to burn coins
		errBurnCoins := k.bank.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
		if errBurnCoins != nil {
			// Try to send coins from module to account
			errSendCoinsFromModuleToAccount := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiverAddress, sdk.NewCoins(coin))
			if errSendCoinsFromModuleToAccount != nil {
				// returns error
				return types.CoinsStuckOnModule, errSendCoinsFromModuleToAccount
			}
		}
	}

	return types.Completed, types.ErrTradeProcessedSuccessfully
}

func (k Keeper) IsAddressLinkedToValidator(stakingKeeper types.StakingKeeper, goCtx context.Context, address string) (isAddressLinkedToValidator bool) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddress, _ := sdk.AccAddressFromBech32(address)
	isLinked := false

	chainValidators := stakingKeeper.GetAllValidators(ctx)

	// Loop through the validators
	for _, validator := range chainValidators {
		valAddress, _ := sdk.ValAddressFromBech32(validator.OperatorAddress)

		_, isLinked = stakingKeeper.GetDelegation(ctx, accAddress, valAddress)
		if isLinked {
			break
		}
	}

	return isLinked
}

func (k Keeper) CancelExpiredPendingTrades(goCtx context.Context) (err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	allStoredTrade := k.GetAllStoredTempTrade(ctx)
	status := types.Canceled

	currentDate := time.Now()

	if err != nil {
		return err
	} else {
		for i := 0; i < len(allStoredTrade); i++ {
			createDate := allStoredTrade[i].CreateDate
			formatedCreateDate, err := time.Parse("2006-01-02 15:04", createDate)
			differenceTime := currentDate.Sub(formatedCreateDate)
			totalDays := int(differenceTime.Hours() / 24)

			if err != nil {
				return sdkErrors.Wrapf(types.ErrInvalidDateFormat, types.ErrInvalidDateFormat.Error())
			} else if totalDays >= 1 {

				storedTrade, _ := k.GetStoredTrade(ctx, allStoredTrade[i].TempTradeIndex)
				storedTrade.Status = status
				storedTrade.UpdateDate = currentDate.Format("2006-01-02 15:04")

				k.SetStoredTrade(ctx, storedTrade)
				k.RemoveStoredTempTrade(ctx, allStoredTrade[i].TempTradeIndex)

				ctx.EventManager().EmitEvent(
					sdk.NewEvent(types.CancelExpiredPendingTradesEventType),
				)
			}

		}
	}

	return err
}

/* func (k Keeper) CalculateTradeValueLast30Days(goCtx context.Context) (err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	allStoredTrade := k.GetAllStoredTrade(ctx)
	totalMint := 0
	totalBurn := 0
	GGEZMintTotal := 0
	GGEZBurnTotal := 0
	currentDate := time.Now()

	if err != nil {
		return err
	} else {
		for i := 0; i < len(allStoredTrade); i++ {
			createDate := allStoredTrade[i].CreateDate
			formatedCreateDate, err := time.Parse("2006-01-02 15:04", createDate)
			differenceTime := currentDate.Sub(formatedCreateDate)
			totalDays := int(differenceTime.Hours() / 24)

			if err != nil {

				return sdkErrors.Wrapf(types.ErrInvalidDateFormat, types.ErrInvalidDateFormat.Error())

			} else if allStoredTrade[i].TradeType == "buy" && allStoredTrade[i].Status == "Completed" && totalDays >= 30 {
				intTotalMint, err := strconv.Atoi(allStoredTrade[i].Price)

				if err != nil {
					return err
				}

				totalMint = totalMint + intTotalMint

			}

		}
	}

	return err
} */

func (k Keeper) GetStakingTest(staking types.StakingKeeper) (stakingKeeper types.StakingKeeper) {
	return staking
}

func (k Keeper) ValidateTradeData(tradeData string) (valid bool,err error) {

	isJson := IsJSON(tradeData)

	if isJson == true {
		var data TradeDataObject
		err = json.Unmarshal([]byte(tradeData), &data)
		if err != nil {
			return false,errors.Wrap(types.ErrInvalidTradeDataObject,"Invalid Trade Data Object")
		}
		tradeData := data.TradeData
		if tradeData.AssetHolderID <= 0 {
			return false,errors.Wrap(types.ErrTradeDataAssetHolderID,"Invalid Trade Data Object")
		}
		if tradeData.AssetID <= 0 {
			return false,errors.Wrap(types.ErrTradeDataAssetID,"Invalid Trade Data Object")
		}
		if tradeData.TradeRequestID <= 0 {
			return false,errors.Wrap(types.ErrTradeDataRequestID,"Invalid Trade Data Object")
		}
		if tradeData.TradeValue <= 0 {
			return false,errors.Wrap(types.ErrTradeDataValue,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Currency) == "" {
			return false,errors.Wrap(types.ErrTradeDataCurrency,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Exchange) == "" {
			return false,errors.Wrap(types.ErrTradeDataExchange,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.FundName) == "" {
			return false,errors.Wrap(types.ErrTradeDataFundName,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Issuer) == "" {
			return false,errors.Wrap(types.ErrTradeDataIssuer,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.NoShares) == "" {
			return false,errors.Wrap(types.ErrTradeDataNoShares,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Price) == "" {
			return false,errors.Wrap(types.ErrTradeDataPrice,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Quantity) == "" {
			return false,errors.Wrap(types.ErrTradeDataQuantity,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Segment) == "" {
			return false,errors.Wrap(types.ErrTradeDataSegment,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.SharePrice) == "" {
			return false,errors.Wrap(types.ErrTradeDataSharePrice,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.Ticker) == "" {
			return false,errors.Wrap(types.ErrTradeDataTicker,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.TradeFee) == "" {
			return false,errors.Wrap(types.ErrTradeDataFee,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.TradeNetPrice) == "" {
			return false,errors.Wrap(types.ErrTradeDataNetPrice,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.TradeNetValue) == "" {
			return false,errors.Wrap(types.ErrTradeDataNetValue,"Invalid Trade Data Object")
		}
		if strings.TrimSpace(tradeData.TradeType) == "" {
			return false,errors.Wrap(types.ErrInvalidTradeType,"Invalid Trade Data Object")
		}


		brokerage := data.Brokerage
		if strings.TrimSpace(brokerage.Country) == "" {
			return false,errors.Wrap(types.ErrBrokerageCountry,"Invalid Brokerage Country Object")
		}
		if strings.TrimSpace(brokerage.Type) == "" {
			return false,errors.Wrap(types.ErrBrokerageType,"Invalid Brokerage Type Object")
		}
		if strings.TrimSpace(brokerage.Name) == "" {
			return false,errors.Wrap(types.ErrBrokerageName,"Invalid Brokerage Name Object")
		}
		return true,err
	}
	return false, errors.Wrap(types.ErrInvalidTradeDataJSON,"Invalid Trade Data JSON")
	
}

func IsJSON(str string) bool {
	var jsonFormat json.RawMessage
    return json.Unmarshal([]byte(str), &jsonFormat) == nil
}
