package v3_0_0

import (
	"context"
	"encoding/json"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
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

		customEvmGenesis := evmtypes.DefaultGenesisState()
		customEvmGenesis.Params.EvmDenom = "uggez1"
		customEvmGenesis.Params.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{
			ExtendedDenom: "uggez1",
		}

		// Marshal the custom genesis state
		customEvmGenesisJSON := appCodec.MustMarshalJSON(customEvmGenesis)
		_, err := mm.InitGenesis(ctx, appCodec, map[string]json.RawMessage{
			evmtypes.ModuleName: customEvmGenesisJSON,
		})
		if err != nil {
			return nil, err
		}

		if err := evmkeeper.InitEvmCoinInfo(ctx); err != nil {
			return nil, err
		}

		fromVM[evmtypes.ModuleName] = evm.AppModule{}.ConsensusVersion()

		logger.Debug("running module migrations ...")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
