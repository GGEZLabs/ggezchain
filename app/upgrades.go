package app

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/GGEZLabs/ggezchain/v2/app/upgrade/v2_0_0"
	"github.com/GGEZLabs/ggezchain/v2/app/upgrade/v2_1_0"
	acltypes "github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	epochstypes "github.com/cosmos/cosmos-sdk/x/epochs/types"
	protocolpooltypes "github.com/cosmos/cosmos-sdk/x/protocolpool/types"
)

func (app *App) setupUpgradeHandlers(configurator module.Configurator) {
	app.UpgradeKeeper.SetUpgradeHandler(
		v2_0_0.UpgradeName,
		v2_0_0.CreateUpgradeHandler(app.ModuleManager, configurator),
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		v2_1_0.UpgradeName,
		v2_1_0.CreateUpgradeHandler(app.ModuleManager, configurator),
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Errorf("failed to read upgrade info from disk: %w", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	var storeUpgrades *storetypes.StoreUpgrades

	switch upgradeInfo.Name {
	case v2_0_0.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{acltypes.ModuleName},
		}
	case v2_1_0.UpgradeName:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{
				epochstypes.ModuleName,
				protocolpooltypes.ModuleName,
			},
		}
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
