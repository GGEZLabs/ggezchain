package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	if err := genState.Validate(); err != nil {
		return err
	}

	if err := k.TradeIndex.Set(ctx, genState.TradeIndex); err != nil {
		return err
	}
	for _, elem := range genState.StoredTrades {
		if err := k.StoredTrade.Set(ctx, elem.TradeIndex, elem); err != nil {
			return err
		}
	}
	for _, elem := range genState.StoredTempTrades {
		if err := k.StoredTempTrade.Set(ctx, elem.TradeIndex, elem); err != nil {
			return err
		}
	}

	return k.Params.Set(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	var err error

	genesis := types.DefaultGenesis()
	genesis.Params, err = k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}
	tradeIndex, err := k.TradeIndex.Get(ctx)
	if err != nil {
		if !errors.Is(err, collections.ErrNotFound) {
			return nil, err
		}
	} else {
		genesis.TradeIndex = tradeIndex
	}
	if err := k.StoredTrade.Walk(ctx, nil, func(_ uint64, val types.StoredTrade) (stop bool, err error) {
		genesis.StoredTrades = append(genesis.StoredTrades, val)
		return false, nil
	}); err != nil {
		return nil, err
	}
	if err := k.StoredTempTrade.Walk(ctx, nil, func(_ uint64, val types.StoredTempTrade) (stop bool, err error) {
		genesis.StoredTempTrades = append(genesis.StoredTempTrades, val)
		return false, nil
	}); err != nil {
		return nil, err
	}

	return genesis, nil
}
