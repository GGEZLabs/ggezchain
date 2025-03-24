package keeper

import (
	"context"
	"strconv"
	"time"

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
	formattedDate := currentTime.Format(time.RFC3339)
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
			TradeIndex:           msg.TradeIndex,
			TradeType:            tradeData.TradeType,
			Coin:                 tradeData.Coin,
			Price:                tradeData.Price,
			Quantity:             tradeData.Quantity,
			ReceiverAddress:      tradeData.ReceiverAddress,
			Status:               status,
			Maker:                tradeData.Maker,
			Checker:              msg.Creator,
			UpdateDate:           formattedDate,
			CreateDate:           tradeData.CreateDate,
			ProcessDate:          formattedDate,
			TradeData:            tradeData.TradeData,
			BankingSystemData:    tradeData.BankingSystemData,
			CoinMintingPriceJSON: tradeData.CoinMintingPriceJSON,
			ExchangeRateJSON:     tradeData.ExchangeRateJSON,
			Result:               result,
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
			TradeIndex:           msg.TradeIndex,
			TradeType:            tradeData.TradeType,
			Coin:                 tradeData.Coin,
			Price:                tradeData.Price,
			Quantity:             tradeData.Quantity,
			ReceiverAddress:      tradeData.ReceiverAddress,
			Status:               status,
			Maker:                tradeData.Maker,
			Checker:              msg.Creator,
			CreateDate:           tradeData.CreateDate,
			UpdateDate:           formattedDate,
			ProcessDate:          formattedDate,
			TradeData:            tradeData.TradeData,
			BankingSystemData:    tradeData.BankingSystemData,
			CoinMintingPriceJSON: tradeData.CoinMintingPriceJSON,
			ExchangeRateJSON:     tradeData.ExchangeRateJSON,
			Result:               result,
		}

		k.Keeper.SetStoredTrade(ctx, storedTrade)
		k.RemoveStoredTempTrade(ctx, msg.TradeIndex)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProcessTrade,
			sdk.NewAttribute(types.AttributeKeyTradeIndex, strconv.FormatUint(msg.TradeIndex, 10)),
			sdk.NewAttribute(types.AttributeKeyStatus, status),
			sdk.NewAttribute(types.AttributeKeyChecker, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyMaker, tradeData.Maker),
			sdk.NewAttribute(types.AttributeKeyTradeData, tradeData.TradeData),
			sdk.NewAttribute(types.AttributeKeyCreateDate, tradeData.CreateDate),
			sdk.NewAttribute(types.AttributeKeyUpdateDate, formattedDate),
			sdk.NewAttribute(types.AttributeKeyProcessDate, formattedDate),
			sdk.NewAttribute(types.AttributeKeyResult, result),
		),
	)

	return &types.MsgProcessTradeResponse{
		TradeIndex: msg.TradeIndex,
		Status:     status,
	}, err
}
