package keeper_test

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	"github.com/GGEZLabs/ggezchain/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/x/trade/types"

	sdkmath "cosmossdk.io/math"

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
			address:        sample.AccAddress(),
			msgType:        types.TxTypeCreateTrade,
			expectedOutput: false,
			expErr:         true,
			expErrMsg:      "no ACL record found for address",
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
			hasPermission, err := suite.app.TradeKeeper.HasPermission(suite.ctx, tt.address, tt.msgType)
			if tt.expErr {
				suite.Error(err)
				suite.Contains(err.Error(), tt.expErrMsg)
				return
			}
			suite.NoError(err)
			suite.Equal(hasPermission, tt.expectedOutput)
		})
	}
}

func (suite *KeeperTestSuite) TestMintOrBurnCoins() {
	suite.setupTest()
	tests := []struct {
		name           string
		tradeData      types.StoredTrade
		expectedStatus types.TradeStatus
		expectedSupply sdk.Coin
		err            error
	}{
		{
			name: "invalid receiver address",
			tradeData: types.StoredTrade{
				ReceiverAddress: "invalid_address",
				TradeType:       types.TradeTypeSell,
			},
			expectedStatus: types.StatusFailed,
			err:            types.ErrInvalidReceiverAddress,
		},
		{
			name: "invalid trade type",
			tradeData: types.StoredTrade{
				ReceiverAddress: sample.AccAddress(),
				Amount: &sdk.Coin{
					Denom:  types.DefaultCoinDenom,
					Amount: sdkmath.NewInt(10000000),
				},
				TradeType: types.TradeTypeUnspecified,
			},
			expectedStatus: types.StatusFailed,
			err:            types.ErrInvalidTradeType,
		},
		{
			name: "unknown denom",
			tradeData: types.StoredTrade{
				ReceiverAddress: sample.AccAddress(),
				Amount: &sdk.Coin{
					Denom:  "unknown_denom",
					Amount: sdkmath.NewInt(10000000),
				},
				TradeType: types.TradeTypeSell,
			},
			expectedStatus: types.StatusFailed,
			err:            sdkerrors.ErrInsufficientFunds,
		},
		{
			name: "mint max amount coins",
			tradeData: types.StoredTrade{
				Amount: &sdk.Coin{
					Denom:  types.DefaultCoinDenom,
					Amount: sdkmath.NewInt(math.MaxInt64),
				},
				ReceiverAddress: testutil.Alice,
				TradeType:       types.TradeTypeBuy,
			},
			expectedStatus: types.StatusProcessed,
			expectedSupply: sdk.Coin{
				Denom:  types.DefaultCoinDenom,
				Amount: sdkmath.NewInt(math.MaxInt64),
			},
		},
		{
			name: "burn max amount coins",
			tradeData: types.StoredTrade{
				Amount: &sdk.Coin{
					Denom:  types.DefaultCoinDenom,
					Amount: sdkmath.NewInt(math.MaxInt64),
				},
				ReceiverAddress: testutil.Alice,
				TradeType:       types.TradeTypeSell,
			},
			expectedStatus: types.StatusProcessed,
			expectedSupply: sdk.Coin{
				Denom:  types.DefaultCoinDenom,
				Amount: sdkmath.NewInt(0),
			},
		},
		{
			name: "mint coins",
			tradeData: types.StoredTrade{
				Amount: &sdk.Coin{
					Denom:  types.DefaultCoinDenom,
					Amount: sdkmath.NewInt(10000000),
				},
				ReceiverAddress: testutil.Alice,
				TradeType:       types.TradeTypeBuy,
			},
			expectedStatus: types.StatusProcessed,
			expectedSupply: sdk.Coin{
				Denom:  types.DefaultCoinDenom,
				Amount: sdkmath.NewInt(10000000),
			},
		},
		{
			name: "burn fails when burn amount exceeds module account balance",
			tradeData: types.StoredTrade{
				Amount: &sdk.Coin{
					Denom:  types.DefaultCoinDenom,
					Amount: sdkmath.NewInt(100000000000000),
				},
				ReceiverAddress: testutil.Alice,
				TradeType:       types.TradeTypeSell,
			},
			expectedStatus: types.StatusFailed,
			err:            sdkerrors.ErrInsufficientFunds,
		},
		{
			name: "burn coins",
			tradeData: types.StoredTrade{
				Amount: &sdk.Coin{
					Denom:  types.DefaultCoinDenom,
					Amount: sdkmath.NewInt(1000000),
				},
				ReceiverAddress: testutil.Alice,
				TradeType:       types.TradeTypeSell,
			},
			expectedStatus: types.StatusProcessed,
			expectedSupply: sdk.Coin{
				Denom:  types.DefaultCoinDenom,
				Amount: sdkmath.NewInt(9000000),
			},
		},
	}

	for _, tt := range tests {
		suite.Run(fmt.Sprintf("Case %s", tt.name), func() {
			status, err := suite.app.TradeKeeper.MintOrBurnCoins(suite.ctx, tt.tradeData)
			if tt.err != nil {
				suite.Equal(status, tt.expectedStatus)
				suite.ErrorIs(err, tt.err)
				return
			}
			supply := suite.app.BankKeeper.GetSupply(suite.ctx, types.DefaultCoinDenom)
			suite.Equal(supply, tt.expectedSupply)
		})
	}
}

func (suite *KeeperTestSuite) TestCancelExpiredPendingTrades() {
	suite.setupTest()
	blockHeight := int64(1)
	blockTime := time.Now().UTC()
	ctx := suite.ctx.WithBlockHeight(blockHeight).WithBlockTime(blockTime)

	suite.Run("no temp trades", func() {
		suite.app.TradeKeeper.CancelExpiredPendingTrades(ctx)

		trades := suite.app.TradeKeeper.GetAllStoredTrade(ctx)
		suite.Equal(0, len(trades))

		tempTrades := suite.app.TradeKeeper.GetAllStoredTrade(ctx)
		suite.Equal(0, len(tempTrades))
	})

	suite.Run("invalid date format", func() {
		storedTrade := types.StoredTrade{
			TradeIndex:  1,
			Status:      types.StatusPending,
			CreateDate:  "2023-05-11T08:44:00Z",
			UpdateDate:  "",
			ProcessDate: "2023-05-11T08:44:00Z",
		}

		storedTempTrade := types.StoredTempTrade{
			TradeIndex:     1,
			TempTradeIndex: 1,
			CreateDate:     "2023-05-06",
		}

		suite.app.TradeKeeper.SetStoredTrade(ctx, storedTrade)
		suite.app.TradeKeeper.SetStoredTempTrade(ctx, storedTempTrade)

		suite.app.TradeKeeper.CancelExpiredPendingTrades(ctx)

		trades := suite.app.TradeKeeper.GetAllStoredTrade(ctx)
		suite.Equal(1, len(trades))
		suite.Equal(types.StatusPending, trades[0].Status)

		tempTrades := suite.app.TradeKeeper.GetAllStoredTrade(ctx)
		suite.Equal(1, len(tempTrades))
	})

	suite.Run("cancel multiple trades", func() {
		storedTrades, storedTempTrades := generateStoredTrades(10, 15)

		// set stored trades
		for _, storedTrade := range storedTrades {
			suite.app.TradeKeeper.SetStoredTrade(ctx, storedTrade)
		}

		// set stored temp trades
		for _, storedTempTrade := range storedTempTrades {
			suite.app.TradeKeeper.SetStoredTempTrade(ctx, storedTempTrade)
		}

		suite.app.TradeKeeper.CancelExpiredPendingTrades(ctx)

		// total length
		trades := suite.app.TradeKeeper.GetAllStoredTrade(ctx)
		suite.Equal(25, len(trades))

		// cancelled trade length
		cancelledTrades := getCancelledTrades(trades)
		suite.Equal(10, len(cancelledTrades))

		// pending trade length
		tempTradeLength := len(suite.app.TradeKeeper.GetAllStoredTempTrade(ctx))
		suite.Equal(15, tempTradeLength)
	})
}

// generateStoredTrades generates lists of StoredTrade and StoredTempTrade
func generateStoredTrades(expiredCount, notExpiredCount int) ([]types.StoredTrade, []types.StoredTempTrade) {
	var storedTrades []types.StoredTrade
	var storedTempTrades []types.StoredTempTrade
	tradeIndex := uint64(1)

	createStoredTrade := func(isExpired bool) (types.StoredTrade, types.StoredTempTrade) {
		var createDate, processDate string
		if isExpired {
			createDate = generateRandomPastDate()
			processDate = createDate
		} else {
			createDate = time.Now().UTC().Format(time.RFC3339)
			processDate = createDate
		}

		st := types.StoredTrade{
			TradeIndex:  tradeIndex,
			Status:      types.StatusPending,
			CreateDate:  createDate,
			UpdateDate:  "",
			ProcessDate: processDate,
		}
		temp := types.StoredTempTrade{
			TradeIndex:     tradeIndex,
			TempTradeIndex: tradeIndex,
			CreateDate:     createDate,
		}
		tradeIndex++
		return st, temp
	}

	// generate expired trades
	for range expiredCount {
		st, temp := createStoredTrade(true)
		storedTrades = append(storedTrades, st)
		storedTempTrades = append(storedTempTrades, temp)
	}

	// generate not expired trades
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
	// generate a random duration between 24 and 48 hours ago
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
