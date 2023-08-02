package trade_test

import (
	"testing"

	keepertest "github.com/GGEZLabs/testchain/testutil/keeper"
	"github.com/GGEZLabs/testchain/testutil/nullify"
	"github.com/GGEZLabs/testchain/x/trade"
	"github.com/GGEZLabs/testchain/x/trade/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.TradeKeeper(t)
	trade.InitGenesis(ctx, *k, genesisState)
	got := trade.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
