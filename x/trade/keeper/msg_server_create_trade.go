package keeper

import (
	"context"
	"strconv"
	"time"

	"github.com/GGEZLabs/ggezchain/x/trade/types"

	errorsmod "cosmossdk.io/errors"

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

	err = types.ValidateTradeData(msg.TradeData)
	if err != nil {
		return nil, err
	}

	tradeIndex, found := k.Keeper.GetTradeIndex(ctx)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrNotFound, "trade with index %d not found", tradeIndex.NextId)
	}

	currentTime := ctx.BlockTime().UTC()
	formattedDate := currentTime.Format(time.RFC3339)

	newIndex := tradeIndex.NextId
	status := types.StatusPending

	storedTrade := types.StoredTrade{
		TradeIndex:           newIndex,
		Status:               status,
		CreateDate:           formattedDate,
		UpdateDate:           formattedDate,
		TradeType:            msg.TradeType,
		Amount:               msg.Amount,
		Price:                msg.Price,
		ReceiverAddress:      msg.ReceiverAddress,
		Maker:                msg.Creator,
		ProcessDate:          formattedDate,
		TradeData:            msg.TradeData,
		BankingSystemData:    msg.BankingSystemData,
		CoinMintingPriceJson: msg.CoinMintingPriceJson,
		ExchangeRateJson:     msg.ExchangeRateJson,
		Result:               types.TradeCreatedSuccessfully,
	}

	storedTempTrade := types.StoredTempTrade{
		TradeIndex:     newIndex,
		TempTradeIndex: newIndex,
		CreateDate:     formattedDate,
	}

	k.Keeper.SetStoredTrade(ctx, storedTrade)
	k.Keeper.SetStoredTempTrade(ctx, storedTempTrade)

	tradeIndex.NextId++

	k.Keeper.SetTradeIndex(ctx, tradeIndex)

	k.Keeper.CancelExpiredPendingTrades(ctx)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateTrade,
			sdk.NewAttribute(types.AttributeKeyTradeIndex, strconv.FormatUint(newIndex, 10)),
			sdk.NewAttribute(types.AttributeKeyStatus, status.String()),
		),
	)

	return &types.MsgCreateTradeResponse{
		TradeIndex: newIndex,
		Status:     status,
	}, nil
}
