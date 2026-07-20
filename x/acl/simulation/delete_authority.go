package simulation

import (
	"math/rand"

	"github.com/GGEZLabs/ggezchain/v2/x/acl/keeper"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

func SimulateMsgDeleteAuthority(
	ak types.AuthKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		aclAuthorities, err := k.GetAllAclAuthority(ctx)
		if err != nil {
			return simtypes.OperationMsg{}, nil, err
		}

		if len(aclAuthorities) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, "MsgDeleteAuthority", "no authorities"), nil, nil
		}

		if err := k.AclAdmin.Set(ctx, simAccount.Address.String(), types.AclAdmin{Address: simAccount.Address.String()}); err != nil {
			return simtypes.OperationMsg{}, nil, err
		}
		msg := &types.MsgDeleteAuthority{
			Creator:     simAccount.Address.String(),
			AuthAddress: aclAuthorities[r.Intn(len(aclAuthorities))].Address,
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
