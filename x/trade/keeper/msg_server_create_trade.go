package keeper

import (
	"context"
	"encoding/json"
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

	// Validate receiver address if trade type not split or reinvestment
	if td.TradeInfo.TradeType != types.TradeTypeSplit &&
		td.TradeInfo.TradeType != types.TradeTypeReinvestment {
		_, err = sdk.AccAddressFromBech32(msg.ReceiverAddress)
		if err != nil {
			return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid receiver address (%s)", err)
		}
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

	k.Keeper.CancelExpiredPendingTrades(ctx)

	newIndex := tradeIndex.NextId
	tradeType := td.TradeInfo.TradeType
	formattedPrice := types.FormatPrice(td.TradeInfo.Price)

	storedTrade := types.StoredTrade{
		TradeIndex:           newIndex,
		Status:               types.StatusPending,
		CreateDate:           createDateTime,
		UpdateDate:           createDateTime,
		TradeType:            tradeType,
		Price:                formattedPrice,
		Maker:                msg.Creator,
		ProcessDate:          createDateTime,
		BankingSystemData:    msg.BankingSystemData,
		CoinMintingPriceJson: msg.CoinMintingPriceJson,
		ExchangeRateJson:     msg.ExchangeRateJson,
		Result:               types.TradeCreatedSuccessfully,
	}

	switch tradeType {
	case types.TradeTypeSplit, types.TradeTypeReinvestment:
		td.TradeInfo.Quantity = nil
		storedTrade.Amount = nil
		storedTrade.ReceiverAddress = ""

		tdBytes, err := json.Marshal(td)
		if err != nil {
			return nil, errorsmod.Wrapf(sdkerrors.ErrJSONMarshal, "failed to marshal trade data: %s", err)
		}
		storedTrade.TradeData = string(tdBytes)

	default:
		storedTrade.Amount = td.TradeInfo.Quantity
		storedTrade.ReceiverAddress = msg.ReceiverAddress
		storedTrade.TradeData = msg.TradeData
	}

	storedTempTrade := types.StoredTempTrade{
		TradeIndex: newIndex,
		CreateDate: createDateTime,
	}

	k.Keeper.SetStoredTrade(ctx, storedTrade)
	k.Keeper.SetStoredTempTrade(ctx, storedTempTrade)

	tradeIndex.NextId++
	k.Keeper.SetTradeIndex(ctx, tradeIndex)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateTrade,
			sdk.NewAttribute(types.AttributeKeyTradeIndex, fmt.Sprintf("%d", newIndex)),
			sdk.NewAttribute(types.AttributeKeyStatus, types.StatusPending.String()),
		),
	)

	return &types.MsgCreateTradeResponse{
		TradeIndex: newIndex,
		Status:     types.StatusPending,
	}, nil
}
