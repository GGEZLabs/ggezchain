package app_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/app"
	acltypes "github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	tradetypes "github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"
)

func TestRunMigrations(t *testing.T) {
	cases := []struct {
		name         string
		tradeFromV   uint64
		wantTradeToV uint64
		aclFromV     uint64
		wantAclToV   uint64
	}{
		{name: "matches live state", tradeFromV: 2, wantTradeToV: 3, aclFromV: 1, wantAclToV: 2},
		{name: "trade still at v1", tradeFromV: 1, wantTradeToV: 3, aclFromV: 1, wantAclToV: 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := app.Setup(t, false)
			ctx := a.NewContextLegacy(true, cmtproto.Header{Height: a.LastBlockHeight()})

			fromVM := a.ModuleManager.GetVersionMap()
			fromVM[tradetypes.ModuleName] = tc.tradeFromV
			fromVM[acltypes.ModuleName] = tc.aclFromV

			toVM, err := a.ModuleManager.RunMigrations(ctx, a.Configurator(), fromVM)
			require.NoError(t, err)
			require.Equal(t, tc.wantTradeToV, toVM[tradetypes.ModuleName])
			require.Equal(t, tc.wantAclToV, toVM[acltypes.ModuleName])
		})
	}
}
