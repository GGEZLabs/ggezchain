package v3_0_0

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	precisebank "github.com/cosmos/evm/x/precisebank"
	precisebanktypes "github.com/cosmos/evm/x/precisebank/types"
	evm "github.com/cosmos/evm/x/vm"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// TODO: Check upgrade
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	evmkeeper *evmkeeper.Keeper,
	appCodec codec.Codec,
) upgradetypes.UpgradeHandler {
	return func(context context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(context)

		logger := ctx.Logger().With("upgrade", UpgradeName)

		evmParams := evmtypes.DefaultParams()
		evmParams.EvmDenom = BaseDenom
		evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: BaseDenom}

		if err := evmkeeper.SetParams(ctx, evmParams); err != nil {
			return nil, err
		}

		if err := evmkeeper.InitEvmCoinInfo(ctx); err != nil {
			return nil, err
		}

		fromVM[evmtypes.ModuleName] = evm.AppModule{}.ConsensusVersion()
		fromVM[precisebanktypes.ModuleName] = precisebank.AppModule{}.ConsensusVersion()

		logger.Debug("running module migrations ...")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
