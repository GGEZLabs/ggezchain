package trade

import (
	"github.com/GGEZLabs/ggezchain/v2/x/trade/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	err := genState.Validate()
	if err != nil {
		panic(err)
	}

	// Set if defined
	k.SetTradeIndex(ctx, genState.TradeIndex)
	// Set all the storedTrade
	for _, elem := range genState.StoredTrades {
		k.SetStoredTrade(ctx, elem)
	}
	// Set all the storedTempTrade
	for _, elem := range genState.StoredTempTrades {
		k.SetStoredTempTrade(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Get all tradeIndex
	tradeIndex, found := k.GetTradeIndex(ctx)
	if found {
		genesis.TradeIndex = tradeIndex
	}
	genesis.StoredTrades = k.GetAllStoredTrade(ctx)
	genesis.StoredTempTrades = k.GetAllStoredTempTrade(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
