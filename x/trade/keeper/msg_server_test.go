package keeper_test

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/store/types"
	acltypes "github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/keeper"
	module "github.com/GGEZLabs/ggezchain/v2/x/trade/module"
	tradetestutil "github.com/GGEZLabs/ggezchain/v2/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
)

// KeeperTestSuite exercises msgServer/trade_manager behavior that needs real
// (mocked) bankKeeper/aclKeeper collaborators, unlike the plain fixture in
// keeper_test.go (which wires the keeper with nil bankKeeper/aclKeeper and is
// only used by the collections-backed query/genesis/params tests).
type KeeperTestSuite struct {
	suite.Suite

	tradeKeeper keeper.Keeper
	bankKeeper  *tradetestutil.MockBankKeeper
	aclKeeper   *tradetestutil.MockAclKeeper
	msgServer   types.MsgServer
	ctx         sdk.Context
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func TestMsgServer(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)
	require.NotNil(t, ms)
	require.NotNil(t, f.ctx)
	require.NotEmpty(t, f.keeper)
}

// setupTest (re)creates a trade keeper wired with fresh gomock bank/acl
// keepers, and grants Alice/Bob/Carol/Eve/Trent their well-known permissions.
func (suite *KeeperTestSuite) setupTest() {
	sdk.GetConfig().SetBech32PrefixForAccount("ggez", "ggez")

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	// testutil.DefaultContextWithDB seeds the block header with time.Now();
	// zero it out here so tests match the old repo's sdk.NewContext(...,
	// cmtproto.Header{}, ...) default (zero time.Time, formatted as
	// "0001-01-01T00:00:00Z") unless a test explicitly overrides BlockTime.
	ctx := testutil.DefaultContextWithDB(suite.T(), storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx.WithBlockTime(time.Time{})

	authority := authtypes.NewModuleAddress(types.GovModuleName)

	ctrl := gomock.NewController(suite.T())
	bankKeeper := tradetestutil.NewMockBankKeeper(ctrl)
	aclKeeper := tradetestutil.NewMockAclKeeper(ctrl)

	tradeKeeper := keeper.NewKeeper(
		storeService,
		encCfg.Codec,
		addressCodec,
		authority,
		bankKeeper,
		aclKeeper,
	)

	suite.Require().NoError(tradeKeeper.Params.Set(ctx, types.DefaultParams()))
	suite.Require().NoError(tradeKeeper.TradeIndex.Set(ctx, types.TradeIndex{NextId: 1}))

	suite.ctx = ctx
	suite.tradeKeeper = tradeKeeper
	suite.bankKeeper = bankKeeper
	suite.aclKeeper = aclKeeper
	suite.msgServer = keeper.NewMsgServerImpl(tradeKeeper)
	suite.setAclAuthority()
}

// setAclAuthority sets up the well-known ACL permissions used across the
// trade module's keeper tests. Unlike the old repo, AclKeeper.GetAclAuthority
// returns (AclAuthority, error) instead of (AclAuthority, bool); a missing
// authority is now signaled via collections.ErrNotFound.
func (suite *KeeperTestSuite) setAclAuthority() {
	// Alice has maker permission in trade module
	suite.aclKeeper.EXPECT().GetAclAuthority(suite.ctx, tradetestutil.Alice).Return(acltypes.AclAuthority{
		Address: tradetestutil.Alice,
		Name:    "Alice",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    types.ModuleName,
				IsMaker:   true,
				IsChecker: false,
			},
		},
	}, nil).AnyTimes()

	// Bob has checker permission in trade module
	suite.aclKeeper.EXPECT().GetAclAuthority(suite.ctx, tradetestutil.Bob).Return(acltypes.AclAuthority{
		Address: tradetestutil.Bob,
		Name:    "Bob",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    types.ModuleName,
				IsMaker:   false,
				IsChecker: true,
			},
		},
	}, nil).AnyTimes()

	// Carol has maker permission in acl module (no permissions for trade module)
	suite.aclKeeper.EXPECT().GetAclAuthority(suite.ctx, tradetestutil.Carol).Return(acltypes.AclAuthority{
		Address: tradetestutil.Carol,
		Name:    "Carol",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    acltypes.ModuleName,
				IsMaker:   true,
				IsChecker: false,
			},
		},
	}, nil).AnyTimes()

	// Eve does not exist in AclAuthority
	suite.aclKeeper.EXPECT().GetAclAuthority(suite.ctx, tradetestutil.Eve).Return(acltypes.AclAuthority{}, collections.ErrNotFound).AnyTimes()

	// Trent has maker and checker permission in trade module
	suite.aclKeeper.EXPECT().GetAclAuthority(suite.ctx, tradetestutil.Trent).Return(acltypes.AclAuthority{
		Address: tradetestutil.Trent,
		Name:    "Trent",
		AccessDefinitions: []*acltypes.AccessDefinition{
			{
				Module:    types.ModuleName,
				IsMaker:   true,
				IsChecker: true,
			},
		},
	}, nil).AnyTimes()
}

// createNTrades resets the suite via setupTest and creates numberOfTrades
// sample trades (all as Alice, who has maker permission), returning their indexes.
func (suite *KeeperTestSuite) createNTrades(numberOfTrades uint64) (tradeIndex []uint64) {
	suite.setupTest()
	var indexes []uint64

	for i := uint64(0); i < numberOfTrades; i++ {
		createResponse, err := suite.msgServer.CreateTrade(suite.ctx, types.GetSampleMsgCreateTrade())
		suite.Require().NoError(err)
		suite.Require().Equal(types.MsgCreateTradeResponse{
			TradeIndex: i + 1,
			Status:     types.StatusPending,
		}, *createResponse)

		indexes = append(indexes, i+1)
	}
	return indexes
}

// --- collections-backed helpers replacing the old repo's raw KVStore wrapper
// methods (k.GetStoredTrade/k.SetStoredTrade/k.GetAllStoredTrade/...), which
// don't exist anymore on the new Keeper (replaced by direct collections
// fields: k.StoredTrade, k.StoredTempTrade, k.TradeIndex). ---

func getStoredTrade(ctx context.Context, k keeper.Keeper, idx uint64) (types.StoredTrade, bool) {
	v, err := k.StoredTrade.Get(ctx, idx)
	if err != nil {
		return types.StoredTrade{}, false
	}
	return v, true
}

func setStoredTrade(ctx context.Context, k keeper.Keeper, st types.StoredTrade) {
	if err := k.StoredTrade.Set(ctx, st.TradeIndex, st); err != nil {
		panic(err)
	}
}

func getAllStoredTrade(ctx context.Context, k keeper.Keeper) []types.StoredTrade {
	var out []types.StoredTrade
	_ = k.StoredTrade.Walk(ctx, nil, func(_ uint64, v types.StoredTrade) (bool, error) {
		out = append(out, v)
		return false, nil
	})
	return out
}

func getStoredTempTrade(ctx context.Context, k keeper.Keeper, idx uint64) (types.StoredTempTrade, bool) {
	v, err := k.StoredTempTrade.Get(ctx, idx)
	if err != nil {
		return types.StoredTempTrade{}, false
	}
	return v, true
}

func setStoredTempTrade(ctx context.Context, k keeper.Keeper, tt types.StoredTempTrade) {
	if err := k.StoredTempTrade.Set(ctx, tt.TradeIndex, tt); err != nil {
		panic(err)
	}
}

func getAllStoredTempTrade(ctx context.Context, k keeper.Keeper) []types.StoredTempTrade {
	var out []types.StoredTempTrade
	_ = k.StoredTempTrade.Walk(ctx, nil, func(_ uint64, v types.StoredTempTrade) (bool, error) {
		out = append(out, v)
		return false, nil
	})
	return out
}

func getTradeIndex(ctx context.Context, k keeper.Keeper) (types.TradeIndex, bool) {
	v, err := k.TradeIndex.Get(ctx)
	if err != nil {
		return types.TradeIndex{}, false
	}
	return v, true
}
