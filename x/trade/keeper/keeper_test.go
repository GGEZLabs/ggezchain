package keeper_test

import (
	"testing"

	ggezchainapp "github.com/GGEZLabs/ggezchain/app"
	acltypes "github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/GGEZLabs/ggezchain/x/trade/keeper"
	"github.com/GGEZLabs/ggezchain/x/trade/testutil"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *ggezchainapp.App
	msgServer   types.MsgServer
	ctx         sdk.Context
	queryClient types.QueryClient
}

func TestTradeKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) setupTest() {
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	suite.NoError(err)
	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100000000000000))),
	}

	app := ggezchainapp.SetupWithGenesisValSet(suite.T(), valSet, []authtypes.GenesisAccount{acc}, balance)
	ctx := app.BaseApp.NewContext(false)
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.msgServer = keeper.NewMsgServerImpl(app.TradeKeeper)
	suite.ctx = ctx
	suite.queryClient = queryClient
	suite.setAclAuthority(suite.app)
}

func (suite *KeeperTestSuite) setAclAuthority(app *ggezchainapp.App) {
	authorities := []acltypes.AclAuthority{
		{
			Address: testutil.Alice,
			Name:    "Alice",
			AccessDefinitions: []*acltypes.AccessDefinition{
				{
					Module:    types.ModuleName,
					IsMaker:   true,
					IsChecker: false,
				},
			},
		},
		{
			Address: testutil.Bob,
			Name:    "Bob",
			AccessDefinitions: []*acltypes.AccessDefinition{
				{
					Module:    types.ModuleName,
					IsMaker:   false,
					IsChecker: true,
				},
			},
		},
		{
			Address: testutil.Carol,
			Name:    "Carol",
			AccessDefinitions: []*acltypes.AccessDefinition{
				{
					Module:    acltypes.ModuleName,
					IsMaker:   true,
					IsChecker: false,
				},
			},
		},
		{
			Address: testutil.Trent,
			Name:    "Trent",
			AccessDefinitions: []*acltypes.AccessDefinition{
				{
					Module:    types.ModuleName,
					IsMaker:   true,
					IsChecker: true,
				},
			},
		},
	}

	for _, auth := range authorities {
		app.AclKeeper.SetAclAuthority(suite.ctx, acltypes.AclAuthority{
			Address:           auth.Address,
			Name:              auth.Name,
			AccessDefinitions: auth.AccessDefinitions,
		})
	}
}

func (suite *KeeperTestSuite) createTrade(numberOfTrades uint64) (tradeIndex []uint64) {
	suite.setupTest()
	var indexes []uint64

	for i := uint64(0); i < numberOfTrades; i++ {
		createResponse, err := suite.msgServer.CreateTrade(suite.ctx, types.GetSampleMsgCreateTrade())
		suite.Nil(err)
		suite.EqualValues(types.MsgCreateTradeResponse{
			TradeIndex: i + 1,
			Status:     types.StatusPending,
		}, *createResponse)

		indexes = append(indexes, i+1)
	}
	return indexes
}
