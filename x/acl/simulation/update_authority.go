package simulation

import (
	"math/rand"
	"strconv"

	"github.com/GGEZLabs/ggezchain/v2/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

func SimulateMsgUpdateAuthority(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		i := r.Int()

		aclAuthorities := k.GetAllAclAuthority(ctx)

		if len(aclAuthorities) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, "MsgUpdateAuthority", "no authorities"), nil, nil
		}

		k.SetAclAdmin(ctx, types.AclAdmin{Address: simAccount.Address.String()})
		msg := &types.MsgUpdateAuthority{
			Creator:                    simAccount.Address.String(),
			AuthAddress:                aclAuthorities[r.Intn(len(aclAuthorities))].Address,
			NewName:                    strconv.Itoa(i),
			OverwriteAccessDefinitions: `[{"module":"trade","is_maker":true,"is_checker":true}]`,
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           txGen,
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
