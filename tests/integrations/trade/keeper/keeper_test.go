package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	aclkeeper "github.com/GGEZLabs/ggezchain/v2/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/module"
	acltypes "github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/module"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"gotest.tools/v3/assert"
)

type fixture struct {
	ctx         sdk.Context
	queryClient types.QueryClient
	aclKeeper   aclkeeper.Keeper
	bankKeeper  bankkeeper.Keeper
	tradeKeeper *keeper.Keeper
}

func initFixture(tb testing.TB) *fixture {
	tb.Helper()
	sdk.GetConfig().SetBech32PrefixForAccount("ggez", "ggez")

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, acltypes.StoreKey, types.StoreKey,
	)
	cdc := moduletestutil.MakeTestEncodingConfig(auth.AppModuleBasic{}, bank.AppModuleBasic{}, trade.AppModule{}, acl.AppModule{}).Codec

	logger := log.NewTestLogger(tb)
	cms := integration.CreateMultiStore(keys, logger)

	newCtx := sdk.NewContext(cms, cmtproto.Header{}, true, logger)

	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	maccPerms := map[string][]string{
		types.ModuleName: {authtypes.Minter, authtypes.Burner},
	}

	addressCodec := addresscodec.NewBech32Codec("ggez")

	accountKeeper := authkeeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		addressCodec,
		sdk.Bech32MainPrefix,
		authority.String(),
	)

	blockedAddresses := map[string]bool{
		accountKeeper.GetAuthority(): false,
	}
	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		accountKeeper,
		blockedAddresses,
		authority.String(),
		log.NewNopLogger(),
	)

	aclKeeper := aclkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[acltypes.StoreKey]),
		cdc,
		addressCodec,
		authority,
	)

	router := baseapp.NewMsgServiceRouter()
	router.SetInterfaceRegistry(cdc.InterfaceRegistry())

	tradeKeeper := keeper.NewKeeper(
		runtime.NewKVStoreService(keys[types.StoreKey]),
		cdc,
		addressCodec,
		authority,
		bankKeeper,
		aclKeeper,
	)

	err := tradeKeeper.TradeIndex.Set(newCtx, types.TradeIndex{NextId: 1})
	assert.NilError(tb, err)
	err = tradeKeeper.Params.Set(newCtx, types.DefaultParams())
	assert.NilError(tb, err)

	authModule := auth.NewAppModule(cdc, accountKeeper, authsims.RandomGenesisAccounts, nil)
	bankModule := bank.NewAppModule(cdc, bankKeeper, accountKeeper, nil)
	aclModule := acl.NewAppModule(cdc, aclKeeper, accountKeeper, bankKeeper)
	tradeModule := trade.NewAppModule(cdc, tradeKeeper, accountKeeper, bankKeeper, aclKeeper)

	integrationApp := integration.NewIntegrationApp(newCtx, logger, keys, cdc, map[string]appmodule.AppModule{
		authtypes.ModuleName: authModule,
		banktypes.ModuleName: bankModule,
		acltypes.ModuleName:  aclModule,
		types.ModuleName:     tradeModule,
	})

	sdkCtx := sdk.UnwrapSDKContext(integrationApp.Context())

	msgSrvr := keeper.NewMsgServerImpl(tradeKeeper)

	// Register MsgServer and QueryServer
	types.RegisterMsgServer(router, msgSrvr)
	types.RegisterQueryServer(integrationApp.QueryHelper(), keeper.NewQueryServerImpl(tradeKeeper))

	queryClient := types.NewQueryClient(integrationApp.QueryHelper())

	return &fixture{
		ctx:         sdkCtx,
		queryClient: queryClient,
		bankKeeper:  bankKeeper,
		aclKeeper:   aclKeeper,
		tradeKeeper: &tradeKeeper,
	}
}

func setAclAuthority(tb testing.TB, ctx sdk.Context, aclKeeper aclkeeper.Keeper) {
	tb.Helper()

	// Alice has maker permission in trade module
	err := aclKeeper.SetAclAuthority(ctx, acltypes.AclAuthority{
		Address: testutil.Alice,
		Name:    "Alice",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    types.ModuleName,
				IsMaker:   true,
				IsChecker: false,
			},
		},
	})
	assert.NilError(tb, err)

	// Bob has checker permission in trade module
	err = aclKeeper.SetAclAuthority(ctx, acltypes.AclAuthority{
		Address: testutil.Bob,
		Name:    "Bob",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    types.ModuleName,
				IsMaker:   false,
				IsChecker: true,
			},
		},
	})
	assert.NilError(tb, err)

	// Carol has maker permission in acl module (no permissions for trade module)
	err = aclKeeper.SetAclAuthority(ctx, acltypes.AclAuthority{
		Address: testutil.Carol,
		Name:    "Carol",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    acltypes.ModuleName,
				IsMaker:   true,
				IsChecker: false,
			},
		},
	})
	assert.NilError(tb, err)

	// Trent has maker and checker permission in trade module
	err = aclKeeper.SetAclAuthority(ctx, acltypes.AclAuthority{
		Address: testutil.Trent,
		Name:    "Trent",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    types.ModuleName,
				IsMaker:   true,
				IsChecker: true,
			},
		},
	})
	assert.NilError(tb, err)
}

// getAllStoredTrade walks the StoredTrade collection, mirroring the idiom used
// by x/trade/keeper/genesis.go's ExportGenesis. The new collections-based
// keeper does not expose a GetAllStoredTrade convenience method, so tests that
// need every stored trade replicate the walk here.
func getAllStoredTrade(tb testing.TB, ctx context.Context, k keeper.Keeper) []types.StoredTrade {
	tb.Helper()

	var list []types.StoredTrade
	err := k.StoredTrade.Walk(ctx, nil, func(_ uint64, val types.StoredTrade) (stop bool, err error) {
		list = append(list, val)
		return false, nil
	})
	assert.NilError(tb, err)

	return list
}

// getAllStoredTempTrade walks the StoredTempTrade collection; see
// getAllStoredTrade above for why this is inlined rather than a keeper method.
func getAllStoredTempTrade(tb testing.TB, ctx context.Context, k keeper.Keeper) []types.StoredTempTrade {
	tb.Helper()

	var list []types.StoredTempTrade
	err := k.StoredTempTrade.Walk(ctx, nil, func(_ uint64, val types.StoredTempTrade) (stop bool, err error) {
		list = append(list, val)
		return false, nil
	})
	assert.NilError(tb, err)

	return list
}
