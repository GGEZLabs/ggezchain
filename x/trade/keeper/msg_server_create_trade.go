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

func (k msgServer) CreateTrade(ctx context.Context, msg *types.MsgCreateTrade) (*types.MsgCreateTradeResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

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
	if td.TradeInfo.TradeType == types.TradeTypeBuy ||
		td.TradeInfo.TradeType == types.TradeTypeSell {
		_, err = sdk.AccAddressFromBech32(msg.ReceiverAddress)
		if err != nil {
			return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid receiver address (%s)", err)
		}
	} else if msg.ReceiverAddress != "" {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("receiver address must not be set for trade type %s", td.TradeInfo.TradeType.String())
	}

	tradeIndex, err := k.TradeIndex.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrNotFound, "trade index not found")
	}

	currentDateTime := sdkCtx.BlockTime()
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

	// StoredTrade.Amount is a non-nullable sdk.Coin (unlike TradeInfo.Quantity,
	// which is a nullable *sdk.Coin left unset for trade types that carry no
	// quantity, e.g. split/dividends) so fall back to the zero-value Coin.
	var amount sdk.Coin
	if td.TradeInfo.Quantity != nil {
		amount = *td.TradeInfo.Quantity
	}

	storedTrade := types.StoredTrade{
		TradeIndex:           newIndex,
		Status:               types.StatusPending,
		TxDate:               formattedDateTime,
		CreateDate:           createDateTime,
		UpdateDate:           formattedDateTime,
		TradeType:            tradeType,
		Amount:               amount,
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

	if err := k.StoredTrade.Set(ctx, newIndex, storedTrade); err != nil {
		return nil, err
	}
	if err := k.StoredTempTrade.Set(ctx, newIndex, storedTempTrade); err != nil {
		return nil, err
	}

	tradeIndex.NextId++
	if err := k.TradeIndex.Set(ctx, tradeIndex); err != nil {
		return nil, err
	}

	// Cancel expired trades
	k.CancelExpiredPendingTrades(ctx)

	sdkCtx.EventManager().EmitEvent(
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
