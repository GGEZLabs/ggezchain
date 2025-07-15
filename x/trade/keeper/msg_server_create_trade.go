package keeper

import (
	"context"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
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

	err = types.ValidateCoinMintingPriceJson(msg.CoinMintingPriceJson)
	if err != nil {
		return nil, err
	}

	err = types.ValidateExchangeRateJson(msg.ExchangeRateJson)
	if err != nil {
		return nil, err
	}

	// Validate receiver address if trade type not split or reinvestment
	if td.TradeInfo.TradeType != types.TradeTypeSplit &&
		td.TradeInfo.TradeType != types.TradeTypeReverseSplit &&
		td.TradeInfo.TradeType != types.TradeTypeReinvestment &&
		td.TradeInfo.TradeType != types.TradeTypeDividends {
		_, err = sdk.AccAddressFromBech32(msg.ReceiverAddress)
		if err != nil {
			return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid receiver address (%s)", err)
		}
	} else if msg.ReceiverAddress != "" {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("receiver address must not be set for trade type %s", td.TradeInfo.TradeType.String())
	}

	tradeIndex, found := k.Keeper.GetTradeIndex(ctx)
	if !found {
		return nil, errorsmod.Wrapf(sdkerrors.ErrNotFound, "trade with index %d not found", tradeIndex.NextId)
	}

	currentDateTime := ctx.BlockTime()
	formattedDateTime := currentDateTime.Format(time.RFC3339)
	createDateTime := currentDateTime.Format(time.RFC3339)

	if msg.CreateDate != "" {
		if err = types.ValidateDate(currentDateTime, msg.CreateDate); err != nil {
			return nil, err
		}
		createDateTime = msg.CreateDate
	}

	newIndex := tradeIndex.NextId
	tradeType := td.TradeInfo.TradeType
	formattedPrice := types.FormatPrice(td.TradeInfo.CoinMintingPriceUsd)

	storedTrade := types.StoredTrade{
		TradeIndex:           newIndex,
		Status:               types.StatusPending,
		TxDate:               formattedDateTime,
		CreateDate:           createDateTime,
		UpdateDate:           formattedDateTime,
		TradeType:            tradeType,
		Amount:               td.TradeInfo.Quantity,
		TradeData:            msg.TradeData,
		ReceiverAddress:      msg.ReceiverAddress,
		CoinMintingPriceUsd:  formattedPrice,
		Maker:                msg.Creator,
		ProcessDate:          formattedDateTime,
		BankingSystemData:    msg.BankingSystemData,
		CoinMintingPriceJson: msg.CoinMintingPriceJson,
		ExchangeRateJson:     msg.ExchangeRateJson,
		Result:               types.TradeCreatedSuccessfully,
	}

	storedTempTrade := types.StoredTempTrade{
		TradeIndex: newIndex,
		TxDate:     formattedDateTime,
	}

	k.Keeper.SetStoredTrade(ctx, storedTrade)
	k.Keeper.SetStoredTempTrade(ctx, storedTempTrade)

	tradeIndex.NextId++
	k.Keeper.SetTradeIndex(ctx, tradeIndex)

	// Cancel expired trades
	k.Keeper.CancelExpiredPendingTrades(ctx)

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
