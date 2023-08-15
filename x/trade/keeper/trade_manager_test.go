package keeper_test

import (
	"context"
	"testing"

	//"github.com/GGEZLabs/ggezchain/testutil/sample"
	tradeKeeper "github.com/GGEZLabs/ggezchain/x/trade/keeper"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/stretchr/testify/require"
)

func TestCreateTrade(t *testing.T) {

	tests := []struct {
		name string
		msg  types.StoredTempTrade
		err  error
	}{
		{
			name: "case1",
			msg: types.StoredTempTrade{
				TradeIndex:     1,
				TempTradeIndex: 1,
				CreateDate:     "2023-04-28 12:56",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "case2",
			msg: types.StoredTempTrade{

				TradeIndex:     1,
				TempTradeIndex: 1,
				CreateDate:     "2023-05-1 12:56",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tradeKeeper.Keeper.CancelExpiredPendingTrades(tradeKeeper.Keeper{}, context.Background())
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}

}

// func TestIsAddressWhiteListed(t *testing.T) {
// 	_, k, _ := setupMsgServerCreateTrade(t)

// 	tests := []struct {
// 		name                 string
// 		address              string
// 		msgType              string
// 		isAddressWhitelisted bool
// 		err                  error
// 	}{
// 		{
// 			name:                 "address white listed with type create trade",
// 			address:              "ggez16rmzj02qkqp48v7rhlm99fvqrheenqepfeymsu",
// 			msgType:              types.CreateTrade,
// 			err:                  nil,
// 			isAddressWhitelisted: true,
// 		},
// 		{
// 			name:                 "address is not white listed with type create trade",
// 			address:              "xxxx1ktj48fem0vw7kmg2gnny4dplmmvzvvmdatjjhk",
// 			msgType:              types.CreateTrade,
// 			err:                  nil,
// 			isAddressWhitelisted: false,
// 		}, {
// 			name:                 "address white listed with type process trade",
// 			address:              "ggez16rmzj02qkqp48v7rhlm99fvqrheenqepfeymsu",
// 			msgType:              types.ProcessTrade,
// 			err:                  nil,
// 			isAddressWhitelisted: true,
// 		}, {
// 			name: "address is not white listed with type process trade",

// 			address:              "xxxx1ktj48fem0vw7kmg2gnny4dplmmvzvvmdatjjhk",
// 			msgType:              types.ProcessTrade,
// 			err:                  nil,
// 			isAddressWhitelisted: false,
// 		}, {
// 			name:                 "address white listed with invalid type",
// 			address:              "ggez16rmzj02qkqp48v7rhlm99fvqrheenqepfeymsu",
// 			msgType:              "Invalid Type",
// 			err:                  nil,
// 			isAddressWhitelisted: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			isAddressWhitelisted, err := k.IsAddressWhitelisted(tt.address, tt.msgType)
// 			if tt.err != nil {
// 				require.Equal(t, isAddressWhitelisted, tt.isAddressWhitelisted)
// 				return
// 			}
// 			require.NoError(t, err)
// 		})
// 	}

// }
// func TestIsAddressWhiteListedChainACLNotFount(t *testing.T) {
// 	_, k, _ := setupMsgServerCreateTrade(t)
// 	_, err := k.IsAddressWhitelisted("ggez1ktj48fem0vw7kmg2gnny4dplmmvzvvmdatjjhk", types.CreateTrade)
// 	require.EqualValues(t, err, nil)

// }

// func TestCancelExpiredPendingTrades(t *testing.T) {
// 	_, k, context := setupMsgServerCreateTrade(t)

// 	ctx := sdk.UnwrapSDKContext(context)

// 	storedTradeOne := types.StoredTrade{
// 		TradeIndex:      1,
// 		TradeType:       types.Buy,
// 		Coin:            "uggez",
// 		Price:           "10.59",
// 		Quantity:        "5026505",
// 		ReceiverAddress: sample.AccAddress(),
// 		TradeData:       "TradeData",
// 		Status:          "Pending",
// 		Maker:           sample.AccAddress(),
// 		Checker:         sample.AccAddress(),
// 		CreateDate:      "2023-05-07 12:00",
// 		UpdateDate:      "2023-05-03 12:00",
// 		ProcessDate:     "2023-05-03 12:00",
// 		Result:          "2023-05-03 12:00",
// 	}
// 	storedTradeTwo := types.StoredTrade{
// 		TradeIndex:      2,
// 		TradeType:       types.Sell,
// 		Coin:            "uggez",
// 		Price:           "10.32",
// 		Quantity:        "1000000005",
// 		ReceiverAddress: sample.AccAddress(),
// 		TradeData:       "TradeData",
// 		Status:          "Pending",
// 		Maker:           sample.AccAddress(),
// 		Checker:         sample.AccAddress(),
// 		CreateDate:      "2023-05-06 00:00",
// 		UpdateDate:      "",
// 		ProcessDate:     "",
// 		Result:          "",
// 	}
// 	storedTradeThree := types.StoredTrade{
// 		TradeIndex:      3,
// 		TradeType:       types.Buy,
// 		Coin:            "uggez",
// 		Price:           "15.19",
// 		Quantity:        "5000026000",
// 		ReceiverAddress: sample.AccAddress(),
// 		TradeData:       "TradeData",
// 		Status:          "Pending",
// 		Maker:           sample.AccAddress(),
// 		Checker:         sample.AccAddress(),
// 		CreateDate:      "2023-05-06 2:30",
// 		UpdateDate:      "",
// 		ProcessDate:     "",
// 		Result:          "",
// 	}

// 	k.SetStoredTrade(ctx, storedTradeOne)
// 	k.SetStoredTrade(ctx, storedTradeTwo)
// 	k.SetStoredTrade(ctx, storedTradeThree)

// 	storedTempTradeOne := types.StoredTempTrade{
// 		TradeIndex:     1,
// 		TempTradeIndex: 1,
// 		CreateDate:     "2023-05-07 12:00",
// 	}
// 	storedTempTradeTwo := types.StoredTempTrade{
// 		TradeIndex:     2,
// 		TempTradeIndex: 2,
// 		CreateDate:     "2023-05-06 00:00",
// 	}
// 	storedTempTradeThree := types.StoredTempTrade{
// 		TradeIndex:     3,
// 		TempTradeIndex: 3,
// 		CreateDate:     "2023-05-06 2:30",
// 	}

// 	k.SetStoredTempTrade(ctx, storedTempTradeOne)
// 	k.SetStoredTempTrade(ctx, storedTempTradeTwo)
// 	k.SetStoredTempTrade(ctx, storedTempTradeThree)

// 	tests := []struct {
// 		name           string
// 		expectedStatus string
// 		actualStatus   string
// 		err            error
// 	}{
// 		{
// 			name:           "check default status for stored trade ",
// 			expectedStatus: types.Pending,
// 			actualStatus:   types.Pending,
// 			err:            nil,
// 		}, {
// 			name:           "check canceled status for the expired trade ",
// 			expectedStatus: types.Pending,
// 			actualStatus:   types.Pending,
// 			err:            nil,
// 		},
// 	}
// 	//check if expired trade removed from StoredTempTrade by check length
// 	err := k.CancelExpiredPendingTrades(ctx)
// 	require.NoError(t, err)
// 	require.Equal(t, 1, len(k.GetAllStoredTempTrade(ctx)))

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := k.CancelExpiredPendingTrades(ctx)
// 			if tt.err != nil {
// 				require.Equal(t, tt.expectedStatus, tt.actualStatus)
// 				return
// 			}
// 			require.NoError(t, err)
// 		})
// 	}

// }
