package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	ggezchainapp "github.com/GGEZLabs/ggezchain/app"
	"github.com/GGEZLabs/ggezchain/x/trade/keeper"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/cometbft/cometbft/libs/bytes"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"testing"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	app         *ggezchainapp.App
	msgServer   types.MsgServer
	ctx         sdk.Context
	queryClient types.QueryClient
}

var (
	tradeModuleAddress string
)

func TestTradeKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupTestForCreateTrade() {
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	suite.NoError(err)
	// create validator set with single validator
	// validator := tmtypes.NewValidator(pubKey, 1)
	// valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	validator := tmtypes.ValidatorSet{
		Proposer: &tmtypes.Validator{
			Address:          bytes.HexBytes("CF56A1EAC1DB0D0FC12BA4DD1584A6EABB8F907F"),
			PubKey:           pubKey,
			VotingPower:      1,
			ProposerPriority: 0,
		},
		Validators: []*tmtypes.Validator{
			{
				Address:          bytes.HexBytes("CF56A1EAC1DB0D0FC12BA4DD1584A6EABB8F907F"),
				PubKey:           pubKey,
				VotingPower:      1,
				ProposerPriority: 0,
			},
		},
	}

	pvKey := ed25519.GenPrivKeyFromSecret([]byte("5075624B6579456432353531397B313833444644314135303638353141364239373738313933323137444631463141453839324444373143413139443332374646433137383839334132373436427D"))

	// generate genesis account
	acc := authtypes.NewBaseAccount(pvKey.PubKey().Address().Bytes(), pvKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100000000000000))),
	}
	app := ggezchainapp.SetupWithGenesisValSet(suite.T(), &validator, []authtypes.GenesisAccount{acc}, balance)
	ctx := app.BaseApp.NewContext(false)

	// app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	// tradeModuleAddress := app.AccountKeeper.GetModuleAddress(types.ModuleName).String()

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.TradeKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.msgServer = keeper.NewMsgServerImpl(app.TradeKeeper)
	suite.ctx = ctx
	suite.queryClient = queryClient
}

func (suite *IntegrationTestSuite) SetupTestForProcessTrade() {
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	suite.NoError(err)

	// create validator set with single validator
	// validator := tmtypes.NewValidator(pubKey, 1)
	// valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	validator := tmtypes.ValidatorSet{
		Proposer: &tmtypes.Validator{
			Address:          bytes.HexBytes("CF56A1EAC1DB0D0FC12BA4DD1584A6EABB8F907F"),
			PubKey:           pubKey,
			VotingPower:      1,
			ProposerPriority: 0,
		},
		Validators: []*tmtypes.Validator{
			{Address: bytes.HexBytes("CF56A1EAC1DB0D0FC12BA4DD1584A6EABB8F907F"),
				PubKey:           pubKey,
				VotingPower:      1,
				ProposerPriority: 0,
			},
		},
	}

	pvKey := ed25519.GenPrivKeyFromSecret([]byte("5075624B6579456432353531397B313833444644314135303638353141364239373738313933323137444631463141453839324444373143413139443332374646433137383839334132373436427D"))
	delPvKey := ed25519.GenPrivKeyFromSecret([]byte("97B313833444644314135303638353141364239373738313933323137444631463141453839324444373143413139443332374646433137383839334132373436427D"))
	// generate genesis account
	acc := authtypes.NewBaseAccount(pvKey.PubKey().Address().Bytes(), pvKey.PubKey(), 0, 0)
	del := authtypes.NewBaseAccount(delPvKey.PubKey().Address().Bytes(), delPvKey.PubKey(), 0, 0)

	balance := []banktypes.Balance{
		{
			Address: acc.GetAddress().String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(10000000000))),
		},
		{
			Address: del.GetAddress().String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(10000000000))),
		},
	}

	app := ggezchainapp.SetupWithGenesisValSet(suite.T(), &validator, []authtypes.GenesisAccount{acc, del}, balance[0], balance[1])
	ctx := app.BaseApp.NewContext(false)

	// app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	// app.StakingKeeper.SetParams(ctx, stakingtypes.DefaultParams())

	delAdd, err := sdk.AccAddressFromBech32(del.Address)
	if err != nil {
		panic(err)
	}

	val, _ := app.StakingKeeper.GetAllValidators(ctx)
	_, err = app.StakingKeeper.Delegate(ctx, delAdd, sdkmath.NewInt(10000), 1, val[0], true)
	if err != nil {
		panic(err)
	}

	tradeModuleAddress = app.AccountKeeper.GetModuleAddress(types.ModuleName).String()

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.TradeKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.msgServer = keeper.NewMsgServerImpl(app.TradeKeeper)
	suite.ctx = ctx
	suite.queryClient = queryClient
}
