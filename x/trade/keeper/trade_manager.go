package keeper

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

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

	file, err := ioutil.ReadFile(userHomeDir + "/.ggezchain/config/chain_acl.json")
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
				return types.Failed, errSendCoinsFromModuleToAccount
			}

			return types.Failed, errSendCoinsFromModuleToAccount
		}
	} else if tradeData.TradeType == types.Sell {

		errSendCoinsFromAccountToModule := k.bank.SendCoinsFromAccountToModule(ctx, receiverAddress, types.ModuleName, sdk.NewCoins(coin))
		if errSendCoinsFromAccountToModule != nil {
			errMintCoins := k.bank.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
			if errMintCoins != nil {
				return types.Failed, errSendCoinsFromAccountToModule
			}

			return types.Failed, errSendCoinsFromAccountToModule
		}
		errBurnCoins := k.bank.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
		if errBurnCoins != nil {
			return types.Failed, errBurnCoins
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
