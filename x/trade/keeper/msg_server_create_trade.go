package keeper

import (
	"context"
	"time"

	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateTrade(goCtx context.Context, msg *types.MsgCreateTrade) (*types.MsgCreateTradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validateTradeDataErr := k.ValidateTradeData(msg.TradeData)
	if validateTradeDataErr != nil {
		return nil, validateTradeDataErr
	}

	isAllowed, _ := k.IsAddressAllowed(ctx, msg.Creator, types.CreateTrade)
	if !isAllowed {
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

	return &types.MsgCreateTradeResponse{
		TradeIndex: newIndex,
		Status:     status,
	}, nil}
