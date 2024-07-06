package trade

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/GGEZLabs/ggezchain/api/ggezchain/trade"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "TradeIndex",
					Use:       "show-trade-index",
					Short:     "show tradeIndex",
				},
				{
					RpcMethod: "StoredTradeAll",
					Use:       "list-stored-trade",
					Short:     "List all storedTrade",
				},
				{
					RpcMethod:      "StoredTrade",
					Use:            "show-stored-trade [id]",
					Short:          "Shows a storedTrade",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "tradeIndex"}},
				},
				{
					RpcMethod: "StoredTempTradeAll",
					Use:       "list-stored-temp-trade",
					Short:     "List all storedTempTrade",
				},
				{
					RpcMethod:      "StoredTempTrade",
					Use:            "show-stored-temp-trade [id]",
					Short:          "Shows a storedTempTrade",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "tradeIndex"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateTrade",
					Use:            "create-trade [trade-type] [coin] [price] [quantity] [receiver-address] [trade-data]",
					Short:          "Send a createTrade tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "tradeType"}, {ProtoField: "coin"}, {ProtoField: "price"}, {ProtoField: "quantity"}, {ProtoField: "receiverAddress"}, {ProtoField: "tradeData"}},
				},
				{
					RpcMethod:      "ProcessTrade",
					Use:            "process-trade [process-type] [trade-index]",
					Short:          "Send a processTrade tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "processType"}, {ProtoField: "tradeIndex"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
