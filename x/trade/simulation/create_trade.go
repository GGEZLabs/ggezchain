package simulation

import (
	"math/rand"
	"strconv"
	"time"

	"cosmossdk.io/math"
	acltypes "github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/GGEZLabs/ggezchain/x/trade/keeper"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

func SimulateMsgCreateTrade(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	aclk types.AclKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		i := r.Int()

		// Set authority before create trades
		aclk.SetAclAuthority(ctx, acltypes.AclAuthority{
			Address: simAccount.Address.String(),
			Name:    strconv.Itoa(i),
			AccessDefinitions: []*acltypes.AccessDefinition{
				{
					Module:    types.ModuleName,
					IsMaker:   true,
					IsChecker: false,
				},
			},
		})

		randomDate := RandomDate(r)
		if IsFutureDate(randomDate) {
			return simtypes.NoOpMsg(types.ModuleName, "MsgCreateTrade", "create date is future date"), nil, nil
		}

		msg := &types.MsgCreateTrade{
			Creator:   simAccount.Address.String(),
			TradeType: randomTradeType(r),
			Amount: &sdk.Coin{
				Denom:  types.DefaultDenom,
				Amount: math.NewInt(int64(r.Uint64())).Abs(),
			},
			Price:                strconv.Itoa(i),
			ReceiverAddress:      simAccount.Address.String(),
			TradeData:            types.GetSampleTradeData(),
			BankingSystemData:    `{}`,
			CoinMintingPriceJson: `{}`,
			ExchangeRateJson:     `{}`,
			CreateDate:           randomDate,
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			Context:         ctx,
			SimAccount:      simAccount,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: sdk.NewCoins(),
			AccountKeeper:   ak,
			Bankkeeper:      bk,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// randomTradeType Pick a random trade type
func randomTradeType(r *rand.Rand) types.TradeType {
	switch r.Intn(2) {
	case 0:
		return types.TradeTypeBuy
	case 1:
		return types.TradeTypeSell
	default:
		panic("invalid trade type")
	}
}

// RandomDate generates a random RFC3339-formatted date.
func RandomDate(r *rand.Rand) string {
	start := time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC)
	year := time.Now().Year() + 5
	end := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	duration := end.Sub(start)
	randomOffset := time.Duration(r.Int63n(int64(duration)))
	randomTime := start.Add(randomOffset)

	return randomTime.Format(time.RFC3339)
}

func IsFutureDate(dateStr string) bool {
	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		panic(err.Error())
	}
	return parsedDate.After(time.Now())
}
