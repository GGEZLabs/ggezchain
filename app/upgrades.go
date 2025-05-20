package app

import (
	"fmt"

	"github.com/GGEZLabs/ggezchain/app/upgrade/v2_0_0"
	acltypes "github.com/GGEZLabs/ggezchain/x/acl/types"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"
)

func (app *App) setupUpgradeHandlers(configurator module.Configurator) {
	app.UpgradeKeeper.SetUpgradeHandler(
		v2_0_0.UpgradeName,
		v2_0_0.CreateUpgradeHandler(app.ModuleManager, configurator),
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
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
