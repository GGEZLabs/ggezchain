package simulation

import (
	"math/rand"
	"strconv"

	acltypes "github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/GGEZLabs/ggezchain/x/trade/keeper"
	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

func SimulateMsgProcessTrade(
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
					IsMaker:   false,
					IsChecker: true,
				},
			},
		})
		var indexes []uint64
		allStoredTrade := k.GetAllStoredTrade(ctx)
		for _, storedTrade := range allStoredTrade {
			if storedTrade.Status == types.StatusPending {
				indexes = append(indexes, storedTrade.TradeIndex)
			}
		}

		if len(indexes) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, "MsgProcessTrade", "no pending trades available"), nil, nil
		}

		tradeIndex := indexes[r.Intn(len(indexes))]
		trade, _ := k.GetStoredTrade(ctx, tradeIndex)

		if trade.Maker == simAccount.Address.String() {
			return simtypes.NoOpMsg(types.ModuleName, "MsgProcessTrade", "checker must be different than maker"), nil, nil
		}

		msg := &types.MsgProcessTrade{
			Creator:     simAccount.Address.String(),
			ProcessType: randomProcessType(r),
			TradeIndex:  tradeIndex,
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

// randomProcessType Pick a random  process type
func randomProcessType(r *rand.Rand) types.ProcessType {
	switch r.Intn(2) {
	case 0:
		return types.ProcessTypeConfirm
	case 1:
		return types.ProcessTypeReject
	default:
		panic("invalid process type")
	}
}
