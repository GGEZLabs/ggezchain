package app

import (
	"fmt"
	V2 "github.com/GGEZLabs/ggezchain/app/upgrade/v2"
	V3 "github.com/GGEZLabs/ggezchain/app/upgrade/v3"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func (app *App) setupUpgradeHandlers(configurator module.Configurator) {
	app.UpgradeKeeper.SetUpgradeHandler(
		V3.UpgradeName,
		V3.CreateUpgradeHandler(app.mm, configurator),
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()

	if err != nil {
		panic(fmt.Errorf("failed to read upgrade info from disk: %w", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	switch upgradeInfo.Name {
	case V2.UpgradeName:
	case V3.UpgradeName:
	}
}
