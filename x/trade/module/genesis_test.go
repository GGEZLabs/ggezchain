package trade_test

import (
	"testing"

	"cosmossdk.io/math"
	keepertest "github.com/GGEZLabs/ggezchain/v2/testutil/keeper"
	"github.com/GGEZLabs/ggezchain/v2/testutil/nullify"
	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	trade "github.com/GGEZLabs/ggezchain/v2/x/trade/module"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	td := types.GetSampleTradeData(types.TradeTypeBuy)
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		TradeIndex: types.TradeIndex{
			NextId: 3,
		},
		StoredTrades: []types.StoredTrade{
			{
				TradeIndex:        1,
				TradeType:         types.TradeTypeBuy,
				Amount:            &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
				Price:             "0.01",
				ReceiverAddress:   sample.AccAddress(),
				Status:            types.StatusPending,
				Maker:             sample.AccAddress(),
				Checker:           sample.AccAddress(),
				CreateDate:        "2023-05-11T08:44:00Z",
				UpdateDate:        "2023-05-11T08:44:00Z",
				ProcessDate:       "2023-05-11T08:44:00Z",
				TradeData:         td,
				BankingSystemData: "{}",
			},
			{
				TradeIndex:        2,
				TradeType:         types.TradeTypeSell,
				Amount:            &sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
				Price:             "0.01",
				ReceiverAddress:   sample.AccAddress(),
				Status:            types.StatusPending,
				Maker:             sample.AccAddress(),
				Checker:           sample.AccAddress(),
				CreateDate:        "2023-05-11T08:44:00Z",
				UpdateDate:        "2023-05-11T08:44:00Z",
				ProcessDate:       "2023-05-11T08:44:00Z",
				TradeData:         td,
				BankingSystemData: "{}",
			},
		},
		StoredTempTrades: []types.StoredTempTrade{
			{
				TradeIndex: 1,
				TxDate:     "2023-05-11T08:44:00Z",
			},
			{
				TradeIndex: 2,
				TxDate:     "2023-05-11T08:44:00Z",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.TradeKeeper(t)
	trade.InitGenesis(ctx, k, genesisState)
	got := trade.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.TradeIndex, got.TradeIndex)
	require.ElementsMatch(t, genesisState.StoredTrades, got.StoredTrades)
	require.ElementsMatch(t, genesisState.StoredTempTrades, got.StoredTempTrades)
	// this line is used by starport scaffolding # genesis/test/assert
}
