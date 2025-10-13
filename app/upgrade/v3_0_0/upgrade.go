package v3_0_0

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	evmkeeper *evmkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(context context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(context)

		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Debug("running module migrations ...")

		fromVM, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return fromVM, err
		}

		evmParams := evmkeeper.GetParams(ctx)
		evmParams.EvmDenom = BaseDenom
		evmParams.AllowUnprotectedTxs = true // TODO:
		if err := evmkeeper.SetParams(ctx, evmParams); err != nil {
			return nil, err
		}
		
		logger.Info("Upgrade v3 complete")
		return fromVM, nil
	}
}
