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

func SimulateMsgDeleteAdmin(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		aclAdmins := k.GetAllAclAdmin(ctx)

		if len(aclAdmins) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, "MsgDeleteAdmin", "no admins"), nil, nil
		}

		if len(aclAdmins) == 1 {
			return simtypes.NoOpMsg(types.ModuleName, "MsgDeleteAdmin", "at least one admin must remain"), nil, nil
		}

		k.SetSuperAdmin(ctx, types.SuperAdmin{Admin: simAccount.Address.String()})
		msg := &types.MsgDeleteAdmin{
			Creator: simAccount.Address.String(),
			Admins:  []string{aclAdmins[r.Intn(len(aclAdmins))].Address},
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
