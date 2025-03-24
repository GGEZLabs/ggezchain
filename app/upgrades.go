package app

import (
	"fmt"

	v1_0_1 "github.com/GGEZLabs/ggezchain/app/upgrade/v1_0_1"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"
)

func (app *App) setupUpgradeHandlers(configurator module.Configurator) {
	app.UpgradeKeeper.SetUpgradeHandler(
		v1_0_1.UpgradeName,
		v1_0_1.CreateUpgradeHandler(app.ModuleManager, configurator),
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
	case v1_0_1.UpgradeName:

		// in case you want to add ,rename or delete module uncomment this:
		// storeUpgrades = &storetypes.StoreUpgrades{
		// 	Added: []string{"testmodule"},
		// }
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}
