package keeper

import (
	"context"

	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ProcessTrade(goCtx context.Context, msg *types.MsgProcessTrade) (*types.MsgProcessTradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed, _ := k.IsAddressAllowed(ctx, msg.Creator, types.ProcessTrade)
	if !isAllowed {
		return nil, types.ErrInvalidCheckerPermission
	}

	currentTime := ctx.BlockTime().UTC()
	formattedDate := currentTime.Format("2006-01-02 03:04")
	tradeData, found := k.Keeper.GetStoredTrade(ctx, msg.TradeIndex)
	if !found {
		panic("Trade Index not found")
	}

	err := msg.ValidateProcess(tradeData.Status, tradeData.Maker, msg.Creator)
	if err != nil {
		return nil, err
	}

	errResult := *new(error)
	result := "No Result"

	status := tradeData.Status
	if msg.ProcessType == types.Reject {
		status = types.Rejected
		errResult = types.ErrTradeProcessedSuccessfully
		result = errResult.Error()

		storedTrade := types.StoredTrade{
			TradeIndex:      msg.TradeIndex,
			TradeType:       tradeData.TradeType,
			Coin:            tradeData.Coin,
			Price:           tradeData.Price,
			Quantity:        tradeData.Quantity,
			ReceiverAddress: tradeData.ReceiverAddress,
			Status:          status,
			Maker:           tradeData.Maker,
			Checker:         msg.Creator,
			UpdateDate:      formattedDate,
			CreateDate:      tradeData.CreateDate,
			ProcessDate:     formattedDate,
			TradeData:       tradeData.TradeData,
			Result:          result,
		}

		k.Keeper.SetStoredTrade(ctx, storedTrade)
		k.RemoveStoredTempTrade(ctx, msg.TradeIndex)
	} else if msg.ProcessType == types.Confirm {

		coin, err := msg.GetPrepareCoin(tradeData)
		if err != nil {
			return nil, err
		}
		status, errResult = k.MintOrBurnCoins(ctx, tradeData, coin)

		if errResult != nil {
			result = errResult.Error()
		}

		storedTrade := types.StoredTrade{
			TradeIndex:      msg.TradeIndex,
			TradeType:       tradeData.TradeType,
			Coin:            tradeData.Coin,
			Price:           tradeData.Price,
			Quantity:        tradeData.Quantity,
			ReceiverAddress: tradeData.ReceiverAddress,
			Status:          status,
			Maker:           tradeData.Maker,
			Checker:         msg.Creator,
			CreateDate:      tradeData.CreateDate,
			UpdateDate:      formattedDate,
			ProcessDate:     formattedDate,
			TradeData:       tradeData.TradeData,
			Result:          result,
		}

		k.Keeper.SetStoredTrade(ctx, storedTrade)
		k.RemoveStoredTempTrade(ctx, msg.TradeIndex)
	}

	return &types.MsgProcessTradeResponse{
		TradeIndex: msg.TradeIndex,
		Status:     status,
	}, err
}
