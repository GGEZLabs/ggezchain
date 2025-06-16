package keeper

import (
	"context"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateTrade(goCtx context.Context, msg *types.MsgCreateTrade) (*types.MsgCreateTradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	hasPermission, err := k.HasPermission(ctx, msg.Creator, types.TxTypeCreateTrade)
	if err != nil {
		return nil, err
	}

	if !hasPermission {
		return nil, types.ErrInvalidMakerPermission
	}

	td, err := types.ValidateTradeData(msg.TradeData)
	if err != nil {
		return nil, err
	}

	tradeIndex, found := k.Keeper.GetTradeIndex(ctx)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrNotFound, "trade with index %d not found", tradeIndex.NextId)
	}

	currentDateTime := ctx.BlockTime()
	createDateTime := currentDateTime.Format(time.RFC3339)

	if msg.CreateDate != "" {
		if err = types.ValidateDate(currentDateTime, msg.CreateDate); err != nil {
			return nil, err
		}
		createDateTime = msg.CreateDate
	}

	newIndex := tradeIndex.NextId
	status := types.StatusPending

	storedTrade := types.StoredTrade{
		TradeIndex:           newIndex,
		Status:               status,
		CreateDate:           createDateTime,
		UpdateDate:           createDateTime,
		TradeType:            td.TradeInfo.TradeType,
		Amount:               td.TradeInfo.Quantity,
		Price:                fmt.Sprint(td.TradeInfo.Price),
		ReceiverAddress:      msg.ReceiverAddress,
		Maker:                msg.Creator,
		ProcessDate:          createDateTime,
		TradeData:            msg.TradeData,
		BankingSystemData:    msg.BankingSystemData,
		CoinMintingPriceJson: msg.CoinMintingPriceJson,
		ExchangeRateJson:     msg.ExchangeRateJson,
		Result:               types.TradeCreatedSuccessfully,
	}

	storedTempTrade := types.StoredTempTrade{
		TradeIndex: newIndex,
		CreateDate: createDateTime,
	}

	k.Keeper.SetStoredTrade(ctx, storedTrade)
	k.Keeper.SetStoredTempTrade(ctx, storedTempTrade)

	tradeIndex.NextId++

	k.Keeper.SetTradeIndex(ctx, tradeIndex)

	k.Keeper.CancelExpiredPendingTrades(ctx)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateTrade,
			sdk.NewAttribute(types.AttributeKeyTradeIndex, fmt.Sprintf("%d", newIndex)),
			sdk.NewAttribute(types.AttributeKeyStatus, status.String()),
		),
	)

	return &types.MsgCreateTradeResponse{
		TradeIndex: newIndex,
		Status:     status,
	}, nil
}
