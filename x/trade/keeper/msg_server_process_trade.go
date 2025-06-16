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

	st, found := k.Keeper.GetStoredTrade(ctx, msg.TradeIndex)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrNotFound, "trade with index %d not found", msg.TradeIndex)
	}

	err = msg.Validate(st.Status, st.Maker)
	if err != nil {
		return nil, err
	}

	currentTime := ctx.BlockTime()
	formattedDate := currentTime.Format(time.RFC3339)

	result := types.TradeProcessedSuccessfully
	status := st.Status
	var storedTrade types.StoredTrade

	if msg.ProcessType == types.ProcessTypeReject {
		status = types.StatusRejected

		storedTrade = types.StoredTrade{
			TradeIndex:           msg.TradeIndex,
			TradeType:            st.TradeType,
			Amount:               st.Amount,
			Price:                st.Price,
			ReceiverAddress:      st.ReceiverAddress,
			Status:               status,
			Maker:                st.Maker,
			Checker:              msg.Creator,
			UpdateDate:           formattedDate,
			CreateDate:           st.CreateDate,
			ProcessDate:          formattedDate,
			TradeData:            st.TradeData,
			BankingSystemData:    st.BankingSystemData,
			CoinMintingPriceJson: st.CoinMintingPriceJson,
			ExchangeRateJson:     st.ExchangeRateJson,
			Result:               result,
		}
	} else if msg.ProcessType == types.ProcessTypeConfirm {
		status, err = k.MintOrBurnCoins(ctx, st)
		if err != nil {
			result = err.Error()
		}

		storedTrade = types.StoredTrade{
			TradeIndex:           msg.TradeIndex,
			TradeType:            st.TradeType,
			Amount:               st.Amount,
			Price:                st.Price,
			ReceiverAddress:      st.ReceiverAddress,
			Status:               status,
			Maker:                st.Maker,
			Checker:              msg.Creator,
			CreateDate:           st.CreateDate,
			UpdateDate:           formattedDate,
			ProcessDate:          formattedDate,
			TradeData:            st.TradeData,
			BankingSystemData:    st.BankingSystemData,
			CoinMintingPriceJson: st.CoinMintingPriceJson,
			ExchangeRateJson:     st.ExchangeRateJson,
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
			sdk.NewAttribute(types.AttributeKeyMaker, st.Maker),
			sdk.NewAttribute(types.AttributeKeyTradeData, st.TradeData),
			sdk.NewAttribute(types.AttributeKeyCreateDate, st.CreateDate),
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
