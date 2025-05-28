package keeper

import (
	"context"
	"fmt"
	"time"

	acltypes "github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HasPermission checks if the given address has permission
// for a specific msgType within this module based on ACL rules.
func (k Keeper) HasPermission(ctx sdk.Context, address string, msgType int32) (bool, error) {
	authority, found := k.aclKeeper.GetAclAuthority(ctx, address)
	if !found {
		return false, acltypes.ErrAuthorityAddressNotExist.Wrapf("unauthorized account %s", address)
	}

	for _, ad := range authority.AccessDefinitions {
		if ad.Module != types.ModuleName {
			continue
		}

		switch msgType {
		case types.TxTypeCreateTrade:
			if ad.IsMaker {
				return true, nil
			}
			return false, nil

		case types.TxTypeProcessTrade:
			if ad.IsChecker {
				return true, nil
			}
			return false, nil

		default:
			return false, types.ErrInvalidMsgType.Wrapf("unrecognized message type: %d", msgType)
		}
	}
	return false, types.ErrModuleNotFound.Wrapf("no permission for module %s", types.ModuleName)
}

// MintOrBurnCoins processes a trade by minting coins for a 'buy' or burning coins for a 'sell',
// handling transfers and rollbacks on failure.
func (k Keeper) MintOrBurnCoins(ctx sdk.Context, tradeData types.StoredTrade) (types.TradeStatus, error) {
	receiverAddress, err := sdk.AccAddressFromBech32(tradeData.ReceiverAddress)
	if err != nil {
		return types.StatusFailed, types.ErrInvalidReceiverAddress.Wrap(err.Error())
	}

	coins := sdk.NewCoins(*tradeData.Amount)

	switch tradeData.TradeType {
	case types.TradeTypeBuy:
		// Mint coins to module account
		if err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
			return types.StatusFailed, err
		}

		// Send coins to user
		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiverAddress, coins); err != nil {
			// Rollback: burn minted coins
			// todo: check update CoinsStuckOnModule to StatusFailed
			if err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
				return types.StatusFailed, err
			}
			return types.StatusFailed, err
		}

		return types.StatusProcessed, nil

	case types.TradeTypeSell:
		// Move coins from user to module
		if err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, receiverAddress, types.ModuleName, coins); err != nil {
			return types.StatusFailed, err
		}

		// Burn coins from module
		if err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
			// Rollback: refund coins to user
			// todo: check update CoinsStuckOnModule to StatusFailed
			if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiverAddress, coins); err != nil {
				return types.StatusFailed, err
			}
			return types.StatusFailed, err
		}

		return types.StatusProcessed, nil

	default:
		return types.StatusFailed, types.ErrInvalidTradeType
	}
}

// CancelExpiredPendingTrades automatically cancels pending trades older than 1 day.
func (k Keeper) CancelExpiredPendingTrades(goCtx context.Context) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	allStoredTempTrade := k.GetAllStoredTempTrade(ctx)
	var canceledIds []uint64

	currentDate := ctx.BlockTime()

	for i := range allStoredTempTrade {
		createDate := allStoredTempTrade[i].CreateDate
		formattedCreateDate, err := time.Parse(time.RFC3339, createDate)
		if err != nil {
			k.logger.Error("an error occurred while canceling expired trades",
				"trade_index", allStoredTempTrade[i].TradeIndex,
				"error", err.Error(),
				"module", types.ModuleName)
			continue
		}
		differenceTime := currentDate.Sub(formattedCreateDate)
		totalDays := int(differenceTime.Hours() / 24)

		if totalDays >= 1 {
			storedTrade, _ := k.GetStoredTrade(ctx, allStoredTempTrade[i].TempTradeIndex)
			storedTrade.Status = types.StatusCanceled
			storedTrade.UpdateDate = currentDate.Format(time.RFC3339)

			k.SetStoredTrade(ctx, storedTrade)
			k.RemoveStoredTempTrade(ctx, allStoredTempTrade[i].TempTradeIndex)

			canceledIds = append(canceledIds, allStoredTempTrade[i].TempTradeIndex)
		}
	}

	if len(canceledIds) > 0 {
		var attributes []sdk.Attribute

		for _, id := range canceledIds {
			attributes = append(attributes, sdk.NewAttribute(types.AttributeKeyTradeIndex, fmt.Sprintf("%d", id)))
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCancelExpiredPendingTrades,
				attributes...,
			),
		)
	}
}
