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

func (k msgServer) ProcessTrade(goCtx context.Context, msg *types.MsgProcessTrade) (*types.MsgProcessTradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	hasPermission, err := k.HasPermission(ctx, msg.Creator, types.TxTypeProcessTrade)
	if err != nil {
		return nil, err
	}

	if !hasPermission {
		return nil, types.ErrInvalidCheckerPermission
	}

	tradeData, found := k.Keeper.GetStoredTrade(ctx, msg.TradeIndex)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrNotFound, "trade with index %d not found", msg.TradeIndex)
	}

	err = msg.Validate(tradeData.Status, tradeData.Maker)
	if err != nil {
		return nil, err
	}

	currentTime := ctx.BlockTime()
	formattedDate := currentTime.Format(time.RFC3339)

	result := types.TradeProcessedSuccessfully
	status := tradeData.Status
	var storedTrade types.StoredTrade

	if msg.ProcessType == types.ProcessTypeReject {
		status = types.StatusRejected

		storedTrade = types.StoredTrade{
			TradeIndex:           msg.TradeIndex,
			TradeType:            tradeData.TradeType,
			Amount:               tradeData.Amount,
			Price:                tradeData.Price,
			ReceiverAddress:      tradeData.ReceiverAddress,
			Status:               status,
			Maker:                tradeData.Maker,
			Checker:              msg.Creator,
			UpdateDate:           formattedDate,
			CreateDate:           tradeData.CreateDate,
			ProcessDate:          formattedDate,
			TradeData:            tradeData.TradeData,
			BankingSystemData:    tradeData.BankingSystemData,
			CoinMintingPriceJson: tradeData.CoinMintingPriceJson,
			ExchangeRateJson:     tradeData.ExchangeRateJson,
			Result:               result,
		}
	} else if msg.ProcessType == types.ProcessTypeConfirm {
		status, err = k.MintOrBurnCoins(ctx, tradeData)
		if err != nil {
			result = err.Error()
		}

		storedTrade = types.StoredTrade{
			TradeIndex:           msg.TradeIndex,
			TradeType:            tradeData.TradeType,
			Amount:               tradeData.Amount,
			Price:                tradeData.Price,
			ReceiverAddress:      tradeData.ReceiverAddress,
			Status:               status,
			Maker:                tradeData.Maker,
			Checker:              msg.Creator,
			CreateDate:           tradeData.CreateDate,
			UpdateDate:           formattedDate,
			ProcessDate:          formattedDate,
			TradeData:            tradeData.TradeData,
			BankingSystemData:    tradeData.BankingSystemData,
			CoinMintingPriceJson: tradeData.CoinMintingPriceJson,
			ExchangeRateJson:     tradeData.ExchangeRateJson,
			Result:               result,
		}
	}

	k.Keeper.SetStoredTrade(ctx, storedTrade)
	k.RemoveStoredTempTrade(ctx, msg.TradeIndex)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProcessTrade,
			sdk.NewAttribute(types.AttributeKeyTradeIndex, fmt.Sprintf("%d", msg.TradeIndex)),
			sdk.NewAttribute(types.AttributeKeyStatus, status.String()),
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
	}, nil
}
