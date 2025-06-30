package keeper_test

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	acltypes "github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

type KeeperTestSuite struct {
	suite.Suite

	tradeKeeper *keeper.Keeper
	bankKeeper  *testutil.MockBankKeeper
	aclKeeper   *testutil.MockAclKeeper
	msgServer   types.MsgServer
	ctx         sdk.Context
	queryClient types.QueryClient
}

func TestTradeKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) setupTest() {
	sdk.GetConfig().SetBech32PrefixForAccount("ggez", "ggez")
	tradeKeeper, bankKeeper, aclKeeper, encCfg, ctx := setupTradeKeeper(suite.T())
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	queryClient := types.NewQueryClient(queryHelper)
	types.RegisterQueryServer(queryHelper, tradeKeeper)

	suite.ctx = ctx
	suite.tradeKeeper = tradeKeeper
	suite.bankKeeper = bankKeeper
	suite.aclKeeper = aclKeeper
	suite.msgServer = keeper.NewMsgServerImpl(*suite.tradeKeeper)
	suite.queryClient = queryClient
	suite.setAclAuthority()
}

// setupTradeKeeper creates a tradeKeeper as well as all its dependencies.
func setupTradeKeeper(t *testing.T) (
	*keeper.Keeper,
	*testutil.MockBankKeeper,
	*testutil.MockAclKeeper,
	moduletestutil.TestEncodingConfig,
	sdk.Context,
) {
	t.Helper()
	keys := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(keys, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	encCfg := moduletestutil.MakeTestEncodingConfig()
	types.RegisterInterfaces(encCfg.InterfaceRegistry)
	acltypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	banktypes.RegisterInterfaces(encCfg.InterfaceRegistry)
	authtypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	// gomock initializations
	ctrl := gomock.NewController(t)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	aclKeeper := testutil.NewMockAclKeeper(ctrl)

	tradeKeeper := keeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keys),
		log.NewNopLogger(),
		authority.String(),
		bankKeeper,
		aclKeeper,
	)

	// Initialize params
	if err := tradeKeeper.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	tradeKeeper.SetTradeIndex(ctx, types.TradeIndex{NextId: 1})

	return &tradeKeeper, bankKeeper, aclKeeper, encCfg, ctx
}

func (suite *KeeperTestSuite) setAclAuthority() {
	// Alice has maker permission in trade module
	suite.aclKeeper.EXPECT().GetAclAuthority(suite.ctx, testutil.Alice).Return(acltypes.AclAuthority{
		Address: testutil.Alice,
		Name:    "Alice",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    types.ModuleName,
				IsMaker:   true,
				IsChecker: false,
			},
		},
	}, true).AnyTimes()

	// Bob has checker permission in trade module
	suite.aclKeeper.EXPECT().GetAclAuthority(suite.ctx, testutil.Bob).Return(acltypes.AclAuthority{
		Address: testutil.Bob,
		Name:    "Bob",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    types.ModuleName,
				IsMaker:   false,
				IsChecker: true,
			},
		},
	}, true).AnyTimes()

	// Carol has maker permission in acl module (no permissions for trade module)
	suite.aclKeeper.EXPECT().GetAclAuthority(suite.ctx, testutil.Carol).Return(acltypes.AclAuthority{
		Address: testutil.Carol,
		Name:    "Carol",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    acltypes.ModuleName,
				IsMaker:   true,
				IsChecker: false,
			},
		},
	}, true).AnyTimes()

	// Eve not exist in AclAuthority
	suite.aclKeeper.EXPECT().GetAclAuthority(suite.ctx, testutil.Eve).Return(acltypes.AclAuthority{}, false).AnyTimes()

	// Trent has maker and checker permission in trade module
	suite.aclKeeper.EXPECT().GetAclAuthority(suite.ctx, testutil.Trent).Return(acltypes.AclAuthority{
		Address: testutil.Trent,
		Name:    "Trent",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    types.ModuleName,
				IsMaker:   true,
				IsChecker: true,
			},
		},
	}, true).AnyTimes()
}

func (suite *KeeperTestSuite) createNTrades(numberOfTrades uint64) (tradeIndex []uint64) {
	suite.setupTest()
	var indexes []uint64

	for i := uint64(0); i < numberOfTrades; i++ {
		createResponse, err := suite.msgServer.CreateTrade(suite.ctx, types.GetSampleMsgCreateTrade())
		suite.Require().NoError(err)
		suite.Require().EqualValues(types.MsgCreateTradeResponse{
			TradeIndex: i + 1,
			Status:     types.StatusPending,
		}, *createResponse)

		indexes = append(indexes, i+1)
	}
	return indexes
}
