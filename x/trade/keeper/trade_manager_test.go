package keeper_test

import (
	"math"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/GGEZLabs/ggezchain/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (suite *KeeperTestSuite) TestHasPermission() {
	suite.setupTest()
	tests := []struct {
		name           string
		address        string
		msgType        int32
		expectedOutput bool
		expErr         bool
		expErrMsg      string
	}{
		{
			name:           "authority not found",
			address:        testutil.Eve,
			msgType:        types.TxTypeCreateTrade,
			expectedOutput: false,
			expErr:         true,
			expErrMsg:      "unauthorized account",
		},
		{
			name:           "no module match",
			address:        testutil.Carol,
			msgType:        types.TxTypeCreateTrade,
			expectedOutput: false,
			expErr:         true,
			expErrMsg:      "no permission for module",
		},
		{
			name:           "invalid msg type",
			address:        testutil.Alice,
			msgType:        types.TxTypeUnspecified,
			expectedOutput: false,
			expErr:         true,
			expErrMsg:      "unrecognized message type",
		},
		{
			name:           "valid maker permission",
			address:        testutil.Alice,
			msgType:        types.TxTypeCreateTrade,
			expectedOutput: true,
			expErr:         false,
		},
		{
			name:           "valid checker permission",
			address:        testutil.Bob,
			msgType:        types.TxTypeProcessTrade,
			expectedOutput: true,
			expErr:         false,
		},
		{
			name:           "invalid maker permission",
			address:        testutil.Bob,
			msgType:        types.TxTypeCreateTrade,
			expectedOutput: false,
			expErr:         false,
		},
		{
			name:           "invalid checker permission",
			address:        testutil.Alice,
			msgType:        types.TxTypeProcessTrade,
			expectedOutput: false,
			expErr:         false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			hasPermission, err := suite.tradeKeeper.HasPermission(suite.ctx, tt.address, tt.msgType)
			if tt.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tt.expErrMsg)
				return
			}
			suite.Require().NoError(err)
			suite.Require().Equal(hasPermission, tt.expectedOutput)
		})
	}
}

func (suite *KeeperTestSuite) TestMintOrBurnCoins() {
	suite.setupTest()

	suite.Run("invalid receiver address", func() {
		tradeData := types.StoredTrade{
			ReceiverAddress: "invalid_address",
			TradeType:       types.TradeTypeSell,
		}

		status, err := suite.tradeKeeper.MintOrBurnCoins(suite.ctx, tradeData)
		suite.Require().Equal(status, types.StatusFailed)
		suite.Require().ErrorIs(err, types.ErrInvalidReceiverAddress)
	})

	suite.Run("invalid trade type", func() {
		tradeData := types.StoredTrade{
			ReceiverAddress: testutil.Alice,
			Amount: &sdk.Coin{
				Denom:  types.DefaultDenom,
				Amount: sdkmath.NewInt(10000000),
			},
			TradeType: types.TradeTypeUnspecified,
		}

		status, err := suite.tradeKeeper.MintOrBurnCoins(suite.ctx, tradeData)
		suite.Require().Equal(status, types.StatusFailed)
		suite.Require().ErrorIs(err, types.ErrInvalidTradeType)
	})

	suite.Run("unknown denom", func() {
		tradeData := types.StoredTrade{
			ReceiverAddress: testutil.Alice,
			Amount: &sdk.Coin{
				Denom:  "unknown_denom",
				Amount: sdkmath.NewInt(10000000),
			},
			TradeType: types.TradeTypeSell,
		}

		receiverAddress, err := sdk.AccAddressFromBech32(testutil.Alice)
		suite.Require().NoErrorf(err, "invalid receiver address: %s", err)

		suite.bankKeeper.EXPECT().SendCoinsFromAccountToModule(suite.ctx, receiverAddress, types.ModuleName, sdk.Coins{
			{
				Denom:  "unknown_denom",
				Amount: sdkmath.NewInt(10000000),
			},
		}).Return(sdkerrors.ErrInsufficientFunds).Times(1)

		status, err := suite.tradeKeeper.MintOrBurnCoins(suite.ctx, tradeData)
		suite.Require().Equal(status, types.StatusFailed)
		suite.Require().ErrorIs(err, sdkerrors.ErrInsufficientFunds)
	})

	suite.Run("mint max amount coins", func() {
		tradeData := types.StoredTrade{
			Amount: &sdk.Coin{
				Denom:  types.DefaultDenom,
				Amount: sdkmath.NewInt(math.MaxInt64),
			},
			ReceiverAddress: testutil.Alice,
			TradeType:       types.TradeTypeBuy,
		}

		suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName,
			sdk.Coins{
				{
					Denom:  types.DefaultDenom,
					Amount: sdkmath.NewInt(math.MaxInt64),
				},
			}).Return(nil).Times(1)

		receiverAddress, err := sdk.AccAddressFromBech32(testutil.Alice)
		suite.Require().NoErrorf(err, "invalid receiver address: %s", err)

		suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName,
			receiverAddress,
			sdk.Coins{
				{
					Denom:  types.DefaultDenom,
					Amount: sdkmath.NewInt(math.MaxInt64),
				},
			}).Return(nil).Times(1)

		suite.bankKeeper.EXPECT().GetSupply(suite.ctx, types.DefaultDenom).Return(sdk.Coin{
			Denom:  types.DefaultDenom,
			Amount: sdkmath.NewInt(math.MaxInt64),
		}).Times(1)

		status, err := suite.tradeKeeper.MintOrBurnCoins(suite.ctx, tradeData)
		suite.Require().Equal(status, types.StatusProcessed)
		suite.Require().NoError(err)
		supply := suite.bankKeeper.GetSupply(suite.ctx, types.DefaultDenom)
		suite.Require().Equal(supply, sdk.Coin{
			Denom:  types.DefaultDenom,
			Amount: sdkmath.NewInt(math.MaxInt64),
		})
	})

	suite.Run("burn max amount coins", func() {
		tradeData := types.StoredTrade{
			Amount: &sdk.Coin{
				Denom:  types.DefaultDenom,
				Amount: sdkmath.NewInt(math.MaxInt64),
			},
			ReceiverAddress: testutil.Alice,
			TradeType:       types.TradeTypeSell,
		}

		receiverAddress, err := sdk.AccAddressFromBech32(testutil.Alice)
		suite.Require().NoErrorf(err, "invalid receiver address: %s", err)

		suite.bankKeeper.EXPECT().SendCoinsFromAccountToModule(suite.ctx, receiverAddress,
			types.ModuleName,
			sdk.Coins{
				{
					Denom:  types.DefaultDenom,
					Amount: sdkmath.NewInt(math.MaxInt64),
				},
			}).Return(nil).Times(1)

		suite.bankKeeper.EXPECT().BurnCoins(suite.ctx, types.ModuleName,
			sdk.Coins{
				{
					Denom:  types.DefaultDenom,
					Amount: sdkmath.NewInt(math.MaxInt64),
				},
			}).Return(nil).Times(1)

		suite.bankKeeper.EXPECT().GetSupply(suite.ctx, types.DefaultDenom).Return(sdk.Coin{
			Denom:  types.DefaultDenom,
			Amount: sdkmath.NewInt(0),
		}).Times(1)

		status, err := suite.tradeKeeper.MintOrBurnCoins(suite.ctx, tradeData)
		suite.Require().Equal(status, types.StatusProcessed)
		suite.Require().NoError(err)
		supply := suite.bankKeeper.GetSupply(suite.ctx, types.DefaultDenom)
		suite.Require().Equal(supply, sdk.Coin{
			Denom:  types.DefaultDenom,
			Amount: sdkmath.NewInt(0),
		})
	})

	suite.Run("mint coins", func() {
		tradeData := types.StoredTrade{
			Amount: &sdk.Coin{
				Denom:  types.DefaultDenom,
				Amount: sdkmath.NewInt(10000000),
			},
			ReceiverAddress: testutil.Alice,
			TradeType:       types.TradeTypeBuy,
		}

		suite.bankKeeper.EXPECT().MintCoins(suite.ctx, types.ModuleName,
			sdk.Coins{
				{
					Denom:  types.DefaultDenom,
					Amount: sdkmath.NewInt(10000000),
				},
			}).Return(nil).Times(1)

		receiverAddress, err := sdk.AccAddressFromBech32(testutil.Alice)
		suite.Require().NoErrorf(err, "invalid receiver address: %s", err)

		suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName,
			receiverAddress,
			sdk.Coins{
				{
					Denom:  types.DefaultDenom,
					Amount: sdkmath.NewInt(10000000),
				},
			}).Return(nil).Times(1)

		suite.bankKeeper.EXPECT().GetSupply(suite.ctx, types.DefaultDenom).Return(sdk.Coin{
			Denom:  types.DefaultDenom,
			Amount: sdkmath.NewInt(10000000),
		}).Times(1)

		status, err := suite.tradeKeeper.MintOrBurnCoins(suite.ctx, tradeData)
		suite.Require().Equal(status, types.StatusProcessed)
		suite.Require().NoError(err)
		supply := suite.bankKeeper.GetSupply(suite.ctx, types.DefaultDenom)
		suite.Require().Equal(supply, sdk.Coin{
			Denom:  types.DefaultDenom,
			Amount: sdkmath.NewInt(10000000),
		})
	})

	suite.Run("send coins from account to module fails when sent amount exceeds account balance", func() {
		tradeData := types.StoredTrade{
			Amount: &sdk.Coin{
				Denom:  types.DefaultDenom,
				Amount: sdkmath.NewInt(100000000000000),
			},
			ReceiverAddress: testutil.Alice,
			TradeType:       types.TradeTypeSell,
		}

		receiverAddress, err := sdk.AccAddressFromBech32(testutil.Alice)
		suite.Require().NoErrorf(err, "invalid receiver address: %s", err)

		suite.bankKeeper.EXPECT().SendCoinsFromAccountToModule(suite.ctx, receiverAddress,
			types.ModuleName,
			sdk.Coins{
				{
					Denom:  types.DefaultDenom,
					Amount: sdkmath.NewInt(100000000000000),
				},
			}).Return(sdkerrors.ErrInsufficientFunds).Times(1)

		status, err := suite.tradeKeeper.MintOrBurnCoins(suite.ctx, tradeData)
		suite.Require().Equal(status, types.StatusFailed)
		suite.Require().ErrorIs(err, sdkerrors.ErrInsufficientFunds)
	})

	suite.Run("burn coins", func() {
		tradeData := types.StoredTrade{
			Amount: &sdk.Coin{
				Denom:  types.DefaultDenom,
				Amount: sdkmath.NewInt(1000000),
			},
			ReceiverAddress: testutil.Alice,
			TradeType:       types.TradeTypeSell,
		}

		receiverAddress, err := sdk.AccAddressFromBech32(testutil.Alice)
		suite.Require().NoErrorf(err, "invalid receiver address: %s", err)

		suite.bankKeeper.EXPECT().SendCoinsFromAccountToModule(suite.ctx, receiverAddress,
			types.ModuleName,
			sdk.Coins{
				{
					Denom:  types.DefaultDenom,
					Amount: sdkmath.NewInt(1000000),
				},
			}).Return(nil).Times(1)

		suite.bankKeeper.EXPECT().BurnCoins(suite.ctx, types.ModuleName,
			sdk.Coins{
				{
					Denom:  types.DefaultDenom,
					Amount: sdkmath.NewInt(1000000),
				},
			}).Return(nil).Times(1)

		suite.bankKeeper.EXPECT().GetSupply(suite.ctx, types.DefaultDenom).Return(sdk.Coin{
			Denom:  types.DefaultDenom,
			Amount: sdkmath.NewInt(9000000),
		}).Times(1)

		status, err := suite.tradeKeeper.MintOrBurnCoins(suite.ctx, tradeData)
		suite.Require().Equal(status, types.StatusProcessed)
		suite.Require().NoError(err)
		supply := suite.bankKeeper.GetSupply(suite.ctx, types.DefaultDenom)
		suite.Require().Equal(supply, sdk.Coin{
			Denom:  types.DefaultDenom,
			Amount: sdkmath.NewInt(9000000),
		})
	})
}

func (suite *KeeperTestSuite) TestCancelExpiredPendingTrades() {
	suite.setupTest()
	blockHeight := int64(1)
	blockTime := time.Now().UTC()
	ctx := suite.ctx.WithBlockHeight(blockHeight).WithBlockTime(blockTime)

	suite.Run("no temp trades", func() {
		suite.tradeKeeper.CancelExpiredPendingTrades(ctx)

		trades := suite.tradeKeeper.GetAllStoredTrade(ctx)
		suite.Require().Equal(0, len(trades))

		tempTrades := suite.tradeKeeper.GetAllStoredTempTrade(ctx)
		suite.Require().Equal(0, len(tempTrades))
	})

	suite.Run("invalid date format", func() {
		storedTrade := types.StoredTrade{
			TradeIndex:  1,
			Status:      types.StatusPending,
			TxDate:      "2023-05-11T08:44:00Z",
			ProcessDate: "2023-05-11T08:44:00Z",
		}

		storedTempTrade := types.StoredTempTrade{
			TradeIndex: 1,
			TxDate:     "2023-05-06",
		}

		suite.tradeKeeper.SetStoredTrade(ctx, storedTrade)
		suite.tradeKeeper.SetStoredTempTrade(ctx, storedTempTrade)

		suite.tradeKeeper.CancelExpiredPendingTrades(ctx)

		trades := suite.tradeKeeper.GetAllStoredTrade(ctx)
		suite.Require().Equal(1, len(trades))
		suite.Require().Equal(types.StatusPending, trades[0].Status)

		tempTrades := suite.tradeKeeper.GetAllStoredTempTrade(ctx)
		suite.Require().Equal(1, len(tempTrades))
	})

	suite.Run("cancel trade", func() {
		storedTrade := types.StoredTrade{
			TradeIndex:  1,
			Status:      types.StatusPending,
			CreateDate:  "2023-05-11T08:44:00Z",
			ProcessDate: "2023-05-11T08:44:00Z",
		}

		storedTempTrade := types.StoredTempTrade{
			TradeIndex: 1,
			TxDate:     "2023-05-11T08:44:00Z",
		}

		suite.tradeKeeper.SetStoredTrade(ctx, storedTrade)
		suite.tradeKeeper.SetStoredTempTrade(ctx, storedTempTrade)

		suite.tradeKeeper.CancelExpiredPendingTrades(ctx)

		trades := suite.tradeKeeper.GetAllStoredTrade(ctx)
		suite.Require().Equal(1, len(trades))
		suite.Require().Equal(types.StatusCanceled, trades[0].Status)
		suite.Require().Equal(types.TradeIsCanceled, trades[0].Result)

		tempTrades := suite.tradeKeeper.GetAllStoredTempTrade(ctx)
		suite.Require().Equal(0, len(tempTrades))
	})

	suite.Run("cancel multiple trades", func() {
		storedTrades, storedTempTrades := generateStoredTrades(10, 15)

		// Set stored trades
		for _, storedTrade := range storedTrades {
			suite.tradeKeeper.SetStoredTrade(ctx, storedTrade)
		}

		// Set stored temp trades
		for _, storedTempTrade := range storedTempTrades {
			suite.tradeKeeper.SetStoredTempTrade(ctx, storedTempTrade)
		}

		suite.tradeKeeper.CancelExpiredPendingTrades(ctx)

		// Total length
		trades := suite.tradeKeeper.GetAllStoredTrade(ctx)
		suite.Require().Equal(25, len(trades))

		// Cancelled trade length
		cancelledTrades := getCancelledTrades(trades)
		suite.Require().Equal(10, len(cancelledTrades))

		// Pending trade length
		tempTradeLength := len(suite.tradeKeeper.GetAllStoredTempTrade(ctx))
		suite.Require().Equal(15, tempTradeLength)
	})
}

// generateStoredTrades generates lists of StoredTrade and StoredTempTrade
func generateStoredTrades(expiredCount, notExpiredCount int) ([]types.StoredTrade, []types.StoredTempTrade) {
	var storedTrades []types.StoredTrade
	var storedTempTrades []types.StoredTempTrade
	tradeIndex := uint64(1)

	createStoredTrade := func(isExpired bool) (types.StoredTrade, types.StoredTempTrade) {
		var txDate, processDate string
		if isExpired {
			txDate = generateRandomPastDate()
			processDate = txDate
		} else {
			txDate = time.Now().UTC().Format(time.RFC3339)
			processDate = txDate
		}

		st := types.StoredTrade{
			TradeIndex:  tradeIndex,
			Status:      types.StatusPending,
			UpdateDate:  "",
			ProcessDate: processDate,
			TxDate:      txDate,
		}
		temp := types.StoredTempTrade{
			TradeIndex: tradeIndex,
			TxDate:     txDate,
		}
		tradeIndex++
		return st, temp
	}

	// Generate expired trades
	for range expiredCount {
		st, temp := createStoredTrade(true)
		storedTrades = append(storedTrades, st)
		storedTempTrades = append(storedTempTrades, temp)
	}

	// Generate not expired trades
	for range notExpiredCount {
		st, temp := createStoredTrade(false)
		storedTrades = append(storedTrades, st)
		storedTempTrades = append(storedTempTrades, temp)
	}

	return storedTrades, storedTempTrades
}

// generateRandomPastDate returns a string date in RFC3339 format, before 24 hours.
func generateRandomPastDate() string {
	now := time.Now().UTC()
	// Generate a random duration between 24 and 48 hours ago
	min := int64(24 * time.Hour)
	max := int64(48 * time.Hour)
	randomDuration := time.Duration(rand.Int63n(max-min) + min)
	return now.Add(-randomDuration).Format(time.RFC3339)
}

// getCancelledTrades extract cancelled trades
func getCancelledTrades(trades []types.StoredTrade) []types.StoredTrade {
	if len(trades) == 0 {
		return nil
	}

	var cancelledTrades []types.StoredTrade
	for _, trade := range trades {
		if trade.Status == types.StatusCanceled {
			cancelledTrades = append(cancelledTrades, trade)
		}
	}
	return cancelledTrades
}
