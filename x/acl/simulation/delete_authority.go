package simulation

import (
	"math/rand"

	"github.com/GGEZLabs/ggezchain/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

func SimulateMsgDeleteAuthority(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		aclAuthorities := k.GetAllAclAuthority(ctx)

		if len(aclAuthorities) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, "MsgDeleteAuthority", "no authorities"), nil, nil
		}

		k.SetAclAdmin(ctx, types.AclAdmin{Address: simAccount.Address.String()})
		msg := &types.MsgDeleteAuthority{
			Creator:     simAccount.Address.String(),
			AuthAddress: aclAuthorities[r.Intn(len(aclAuthorities))].Address,
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
