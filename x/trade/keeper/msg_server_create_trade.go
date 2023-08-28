package keeper

import (
	"context"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func (k msgServer) CreateTrade(goCtx context.Context, msg *types.MsgCreateTrade) (*types.MsgCreateTradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isValidTradeDataObject , validateTradeDataErr := k.ValidateTradeData(msg.TradeData)
	if !isValidTradeDataObject {
		return nil, validateTradeDataErr
	}
	
	isAllowed, _ := k.IsAddressAllowed(k.stakingKeeper, ctx, msg.Creator, types.CreateTrade)
	if !isAllowed {
		//panic("you don't have permission to perform this action")
		return nil, types.ErrInvalidMakerPermission
	}
	currentTime := time.Now()
	formattedDate := currentTime.Format("2006-01-02 03:04")

	err := msg.Validate()
	if err != nil {
		return nil, err
	}

	tradeIndex, found := k.Keeper.GetTradeIndex(ctx)
	if !found {
		panic("Trade Index not found")
	}

	newIndex := tradeIndex.NextId
	status := types.Pending

	storedTrade := types.StoredTrade{
		TradeIndex:      newIndex,
		Status:          status,
		CreateDate:      formattedDate,
		TradeType:       msg.TradeType,
		Coin:            msg.Coin,
		Price:           msg.Price,
		Quantity:        msg.Quantity,
		ReceiverAddress: msg.ReceiverAddress,
		Maker:           msg.Creator,
		Checker:         "",
		ProcessDate:     formattedDate,
		TradeData:       msg.TradeData,
		Result:          types.ErrTradeCreatedSuccessfully.Error(),
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
		sdk.NewEvent(types.CancelExpiredPendingTradesEventType),
	)
	return &types.MsgCreateTradeResponse{
		TradeIndex: newIndex,
		Status:     status,
	}, nil

}
