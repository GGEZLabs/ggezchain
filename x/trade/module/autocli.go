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
					Use:            "show-stored-trade [trade-index]",
					Short:          "Shows a storedTrade",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "trade_index"}},
				},
				{
					RpcMethod: "StoredTempTradeAll",
					Use:       "list-stored-temp-trade",
					Short:     "List all storedTempTrade",
				},
				{
					RpcMethod:      "StoredTempTrade",
					Use:            "show-stored-temp-trade [trade-index]",
					Short:          "Shows a storedTempTrade",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "trade_index"}},
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
					Use:            "create-trade [trade-type] [amount] [price] [receiver-address] [trade-data] [banking-system-data] [coin-minting-price-json] [exchange-rate-json]",
					Short:          "Send a createTrade tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "trade_type"}, {ProtoField: "amount"}, {ProtoField: "price"}, {ProtoField: "receiver_address"}, {ProtoField: "trade_data"}, {ProtoField: "banking_system_data"}, {ProtoField: "coin_minting_price_json"}, {ProtoField: "exchange_rate_json"}},
				},
				{
					RpcMethod:      "ProcessTrade",
					Use:            "process-trade [trade-index] [process-type]",
					Short:          "Send a processTrade tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "trade_index"}, {ProtoField: "process_type"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
