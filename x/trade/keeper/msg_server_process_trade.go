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

	st.Checker = msg.Creator
	st.UpdateDate = formattedDate
	st.ProcessDate = formattedDate

	defaultResult := types.TradeProcessedSuccessfully
	var finalStatus types.TradeStatus
	var finalResult string

	switch msg.ProcessType {
	case types.ProcessTypeReject:
		finalStatus = types.StatusRejected
		finalResult = defaultResult

	case types.ProcessTypeConfirm:
		if st.TradeType == types.TradeTypeSplit || st.TradeType == types.TradeTypeReinvestment {
			finalStatus = types.StatusProcessed
			finalResult = defaultResult
		} else {
			status, err := k.MintOrBurnCoins(ctx, st)
			if err != nil {
				finalResult = err.Error()
			} else {
				finalResult = defaultResult
			}
			finalStatus = status
		}

	default:
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "unsupported process type: %v", msg.ProcessType)
	}

	st.Status = finalStatus
	st.Result = finalResult

	k.Keeper.SetStoredTrade(ctx, st)
	k.RemoveStoredTempTrade(ctx, msg.TradeIndex)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProcessTrade,
			sdk.NewAttribute(types.AttributeKeyTradeIndex, fmt.Sprintf("%d", msg.TradeIndex)),
			sdk.NewAttribute(types.AttributeKeyStatus, st.Status.String()),
			sdk.NewAttribute(types.AttributeKeyChecker, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyMaker, st.Maker),
			sdk.NewAttribute(types.AttributeKeyTradeData, st.TradeData),
			sdk.NewAttribute(types.AttributeKeyCreateDate, st.CreateDate),
			sdk.NewAttribute(types.AttributeKeyUpdateDate, formattedDate),
			sdk.NewAttribute(types.AttributeKeyProcessDate, formattedDate),
			sdk.NewAttribute(types.AttributeKeyResult, st.Result),
		),
	)

	return &types.MsgProcessTradeResponse{
		TradeIndex: msg.TradeIndex,
		Status:     st.Status,
	}, nil
}
