package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/GGEZLabs/ggezchain/v2/testutil/sample"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	td := types.GetSampleTradeDataJson(types.TradeTypeBuy)
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		TradeIndex: types.TradeIndex{
			NextId: 3,
		},
		StoredTrades: []types.StoredTrade{
			{
				TradeIndex:           1,
				TradeType:            types.TradeTypeBuy,
				Amount:               sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
				CoinMintingPriceUsd:  "0.01",
				ReceiverAddress:      sample.AccAddress(),
				Status:               types.StatusProcessed,
				Maker:                sample.AccAddress(),
				Checker:              sample.AccAddress(),
				CreateDate:           "2023-05-11T08:44:00Z",
				TxDate:               "2023-05-11T08:44:00Z",
				UpdateDate:           "2023-05-11T08:44:00Z",
				ProcessDate:          "2023-05-11T08:44:00Z",
				TradeData:            td,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
				ExchangeRateJson:     types.GetSampleExchangeRateJson(),
			},
			{
				TradeIndex:           2,
				TradeType:            types.TradeTypeSell,
				Amount:               sdk.Coin{Denom: types.DefaultDenom, Amount: math.NewInt(100000)},
				CoinMintingPriceUsd:  "0.01",
				ReceiverAddress:      sample.AccAddress(),
				Status:               types.StatusProcessed,
				Maker:                sample.AccAddress(),
				Checker:              sample.AccAddress(),
				CreateDate:           "2023-05-11T08:44:00Z",
				TxDate:               "2023-05-11T08:44:00Z",
				UpdateDate:           "2023-05-11T08:44:00Z",
				ProcessDate:          "2023-05-11T08:44:00Z",
				TradeData:            td,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
				ExchangeRateJson:     types.GetSampleExchangeRateJson(),
			},
			{
				TradeIndex:           3,
				TradeType:            types.TradeTypeReinvestment,
				Amount:               sdk.Coin{Amount: math.NewInt(0)},
				CoinMintingPriceUsd:  "0.01",
				Status:               types.StatusProcessed,
				Maker:                sample.AccAddress(),
				Checker:              sample.AccAddress(),
				CreateDate:           "2023-05-11T08:44:00Z",
				TxDate:               "2023-05-11T08:44:00Z",
				UpdateDate:           "2023-05-11T08:44:00Z",
				ProcessDate:          "2023-05-11T08:44:00Z",
				TradeData:            td,
				BankingSystemData:    "{}",
				CoinMintingPriceJson: types.GetSampleCoinMintingPriceJson(),
				ExchangeRateJson:     types.GetSampleExchangeRateJson(),
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
	}

	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.EqualExportedValues(t, genesisState.TradeIndex, got.TradeIndex)
	require.ElementsMatch(t, genesisState.StoredTrades, got.StoredTrades)
	require.ElementsMatch(t, genesisState.StoredTempTrades, got.StoredTempTrades)
}

// TestInitGenesis_InvalidStateRejected confirms InitGenesis validates its input
// before applying any state, matching the original repo's defensive
// genState.Validate()+panic behavior in x/trade/module/genesis.go. Cosmos SDK's
// module manager does NOT call ValidateGenesis automatically before InitGenesis
// during InitChain, so InitGenesis must guard itself against a malformed
// genesis.json rather than relying solely on the separate `validate-genesis` CLI
// step.
func TestInitGenesis_InvalidStateRejected(t *testing.T) {
	f := initFixture(t)

	// trade_index 0 is invalid per GenesisState.Validate().
	invalidState := types.GenesisState{
		Params:     types.DefaultParams(),
		TradeIndex: types.TradeIndex{NextId: 2},
		StoredTrades: []types.StoredTrade{
			{TradeIndex: 0},
		},
	}

	err := f.keeper.InitGenesis(f.ctx, invalidState)
	require.Error(t, err)
	require.Contains(t, err.Error(), "trade_index must be more than 0")

	// Nothing should have been written to the store.
	_, err = f.keeper.TradeIndex.Get(f.ctx)
	require.Error(t, err)
}

// TestExportGenesis_DefaultTradeIndexWhenUnset confirms ExportGenesis leaves
// TradeIndex at DefaultGenesis's value (NextId: DefaultIndex) when nothing was
// ever set in the store, instead of exporting a zero-valued (NextId: 0)
// TradeIndex, which would fail Validate() on any subsequent import.
//
// Unlike acl's SuperAdmin (legitimately nil in a valid genesis), trade's
// TradeIndex is required by Validate() — InitGenesis would reject a genesis
// state that omits it, so this scenario is exercised directly against a fresh
// keeper (e.g. genesis export tooling run before InitGenesis ever executed),
// not via an InitGenesis round trip.
func TestExportGenesis_DefaultTradeIndexWhenUnset(t *testing.T) {
	f := initFixture(t)

	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got.TradeIndex)
	require.Equal(t, types.DefaultIndex, got.TradeIndex.NextId)

	// The exported state must itself be valid, matching what a real
	// export -> re-import (e.g. during a chain upgrade) would require.
	require.NoError(t, got.Validate())
}
