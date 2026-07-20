package keeper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	acltypes "github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HasPermission checks if the given address has permission
// for a specific msgType within this module based on ACL rules.
func (k Keeper) HasPermission(ctx context.Context, address string, msgType int32) (bool, error) {
	authority, err := k.aclKeeper.GetAclAuthority(ctx, address)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return false, acltypes.ErrAuthorityAddressDoesNotExist.Wrapf("unauthorized account %s", address)
		}
		return false, err
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
func (k Keeper) MintOrBurnCoins(ctx context.Context, storedTrade types.StoredTrade) (types.TradeStatus, error) {
	receiverAddress, err := sdk.AccAddressFromBech32(storedTrade.ReceiverAddress)
	if err != nil {
		return types.StatusFailed, types.ErrInvalidReceiverAddress.Wrap(err.Error())
	}

	coins := sdk.NewCoins(storedTrade.Amount)

	switch storedTrade.TradeType {
	case types.TradeTypeBuy:
		// Mint coins to module account
		if err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
			return types.StatusFailed, err
		}

		// Send coins to user
		if err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiverAddress, coins); err != nil {
			// Rollback: burn minted coins
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
	logger := ctx.Logger()

	// Snapshot all pending temp trades first; mutating StoredTrade/StoredTempTrade
	// while iterating them directly would mutate the collection out from under
	// the walk, so collect the indexes and dates up front instead.
	type tempTrade struct {
		tradeIndex uint64
		txDate     string
	}
	var allStoredTempTrade []tempTrade
	err := k.StoredTempTrade.Walk(ctx, nil, func(tradeIndex uint64, storedTempTrade types.StoredTempTrade) (stop bool, err error) {
		allStoredTempTrade = append(allStoredTempTrade, tempTrade{tradeIndex: tradeIndex, txDate: storedTempTrade.TxDate})
		return false, nil
	})
	if err != nil {
		logger.Error("an error occurred while canceling expired trades",
			"error", err.Error(),
			"module", types.ModuleName)
		return
	}

	var canceledIds []uint64
	currentDate := ctx.BlockTime()

	for _, tt := range allStoredTempTrade {
		formattedTxDate, err := time.Parse(time.RFC3339, tt.txDate)
		if err != nil {
			logger.Error("an error occurred while canceling expired trades",
				"trade_index", tt.tradeIndex,
				"error", err.Error(),
				"module", types.ModuleName)
			continue
		}
		differenceTime := currentDate.Sub(formattedTxDate)
		totalDays := int(differenceTime.Hours() / 24)

		if totalDays >= 1 {
			storedTrade, err := k.StoredTrade.Get(ctx, tt.tradeIndex)
			if err != nil {
				logger.Error("an error occurred while canceling expired trades",
					"trade_index", tt.tradeIndex,
					"error", err.Error(),
					"module", types.ModuleName)
				continue
			}
			storedTrade.Status = types.StatusCanceled
			storedTrade.UpdateDate = currentDate.Format(time.RFC3339)
			storedTrade.Result = types.TradeIsCanceled

			if err := k.StoredTrade.Set(ctx, tt.tradeIndex, storedTrade); err != nil {
				logger.Error("an error occurred while canceling expired trades",
					"trade_index", tt.tradeIndex,
					"error", err.Error(),
					"module", types.ModuleName)
				continue
			}
			if err := k.StoredTempTrade.Remove(ctx, tt.tradeIndex); err != nil {
				logger.Error("an error occurred while canceling expired trades",
					"trade_index", tt.tradeIndex,
					"error", err.Error(),
					"module", types.ModuleName)
				continue
			}

			canceledIds = append(canceledIds, tt.tradeIndex)
		}
	}

	if len(canceledIds) > 0 {
		var attributes []sdk.Attribute

		for _, id := range canceledIds {
			attributes = append(attributes, sdk.NewAttribute(types.AttributeKeyTradeIndex, fmt.Sprintf("%d", id)))
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCanceledTrades,
				attributes...,
			),
		)
	}
}
