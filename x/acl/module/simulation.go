package acl

import (
	"math/rand"

	"github.com/GGEZLabs/ggezchain/testutil/sample"
	aclsimulation "github.com/GGEZLabs/ggezchain/x/acl/simulation"
	"github.com/GGEZLabs/ggezchain/x/acl/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = aclsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgAddAuthority = "op_weight_msg_add_authority"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddAuthority int = 100

	opWeightMsgDeleteAuthority = "op_weight_msg_delete_authority"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteAuthority int = 100

	opWeightMsgUpdateAuthority = "op_weight_msg_update_authority"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateAuthority int = 100

	opWeightMsgInitAclAdmin = "op_weight_msg_init_acl_admin"
	// TODO: Determine the simulation weight value
	defaultWeightMsgInitAclAdmin int = 100

	opWeightMsgAddAclAdmin = "op_weight_msg_add_acl_admin"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddAclAdmin int = 100

	opWeightMsgDeleteAclAdmin = "op_weight_msg_delete_acl_admin"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteAclAdmin int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	aclGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&aclGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgAddAuthority int
	simState.AppParams.GetOrGenerate(opWeightMsgAddAuthority, &weightMsgAddAuthority, nil,
		func(_ *rand.Rand) {
			weightMsgAddAuthority = defaultWeightMsgAddAuthority
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddAuthority,
		aclsimulation.SimulateMsgAddAuthority(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteAuthority int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteAuthority, &weightMsgDeleteAuthority, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteAuthority = defaultWeightMsgDeleteAuthority
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteAuthority,
		aclsimulation.SimulateMsgDeleteAuthority(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateAuthority int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateAuthority, &weightMsgUpdateAuthority, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateAuthority = defaultWeightMsgUpdateAuthority
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateAuthority,
		aclsimulation.SimulateMsgUpdateAuthority(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgInitAclAdmin int
	simState.AppParams.GetOrGenerate(opWeightMsgInitAclAdmin, &weightMsgInitAclAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgInitAclAdmin = defaultWeightMsgInitAclAdmin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgInitAclAdmin,
		aclsimulation.SimulateMsgInitAclAdmin(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgAddAclAdmin int
	simState.AppParams.GetOrGenerate(opWeightMsgAddAclAdmin, &weightMsgAddAclAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgAddAclAdmin = defaultWeightMsgAddAclAdmin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddAclAdmin,
		aclsimulation.SimulateMsgAddAclAdmin(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteAclAdmin int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteAclAdmin, &weightMsgDeleteAclAdmin, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteAclAdmin = defaultWeightMsgDeleteAclAdmin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteAclAdmin,
		aclsimulation.SimulateMsgDeleteAclAdmin(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddAuthority,
			defaultWeightMsgAddAuthority,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgAddAuthority(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteAuthority,
			defaultWeightMsgDeleteAuthority,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgDeleteAuthority(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUpdateAuthority,
			defaultWeightMsgUpdateAuthority,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgUpdateAuthority(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgInitAclAdmin,
			defaultWeightMsgInitAclAdmin,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgInitAclAdmin(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddAclAdmin,
			defaultWeightMsgAddAclAdmin,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgAddAclAdmin(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgDeleteAclAdmin,
			defaultWeightMsgDeleteAclAdmin,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				aclsimulation.SimulateMsgDeleteAclAdmin(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
