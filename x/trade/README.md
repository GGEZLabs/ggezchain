# x/trade

## Abstract

The `x/trade` module enables seamless investing and trading of tokenized sustainability assets.

Next is a list of features the module currently supports:
* Creation: A user defines the trade assets, and target recipient.
* Processing: Authorized accounts (as defined in the `acl` module) review and execute
  trades using ProcessTrade messages.

## Contents

- [Abstract](#abstract)
- [State](#state)
  - [TradeIndex](#tradeindex)
  - [StoredTrade](#storedtrade)
  - [StoredTempTrade](#storedtemptrade)
- [Messages](#messages)
  - [MsgCreateTrade](#msgcreatetrade)
  - [MsgProcessTrade](#msgprocesstrade)
- [Events](#events)
  - [Message Events](#message-events)
  - [Keeper Events](#keeper-events)
- [Client](#client)
  - [CLI](#cli)
  - [Query](#query)
  - [Transactions](#transactions)
---

## State

### TradeIndex

Each trade in the `x/trade` module is uniquely identified by a `TradeIndex`.

### StoredTrade

The `StoredTrade` represents the complete data structure of a trade within the `x/trade` module.

### StoredTempTrade

The `StoredTempTrade` represents a trade that is currently in a pending state.

---

## Messages

In this section we describe the processing of the `trade` messages and the corresponding updates to the state. All created/modified state objects specified by each message are defined within the [state](#state) section.

### MsgCreateTrade

The `MsgCreateTrade` message creates both a `StoredTrade` and `StoredTempTrade`.

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/trade/tx.proto#L21
```

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/trade/tx.proto#L42-L51
```

This message is expected to fail if:
* signer does not have maker permission.
* invalid trade data format.
* trade index does not found.
* invalid create date format.

### MsgProcessTrade

The `MsgProcessTrade` message creates both a `StoredTrade` and `StoredTempTrade`.

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/trade/tx.proto#L22
```

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/trade/tx.proto#L58-L63
```

This message is expected to fail if:

* signer does not have checker permission.
* `StoredTrade` does not found.
* the maker and checker are the same address.
* the `StoredTrade` is not in a pending state.

---

## Events

The `trade` module emits the following events:

### Message Events

### MsgCreateTrade

| Type         | Attribute Key | Attribute Value |
| ------------ | ------------- | --------------- |
| create_trade | trade_index   | {TradeIndex}    |
| create_trade | status        | {status}        |

### MsgProcessTrade

| Type          | Attribute Key | Attribute Value |
| ------------- | ------------- | --------------- |
| process_trade | trade_index   | {TradeIndex}    |
| process_trade | status        | {status}        |
| process_trade | checker       | {checker}       |
| process_trade | maker         | {maker}         |
| process_trade | trade_data    | {tradeData}     |
| process_trade | create_date   | {createDate}    |
| process_trade | update_date   | {updateDate}    |
| process_trade | process_da    | {processDate}   |
| process_trade | result        | {result}        |


### Keeper Events

### CancelExpiredPendingTrades

```json
{
  "type": "canceled_trades",
  "attributes": [
    {
      "key": "trade_index",
      "value": "{{trade_index}}",
      "index": true
    },
  ]
}
```

---

## Client

### CLI

A user can query and interact with the `trade` module using the CLI.

#### Query

The `query` commands allow users to query `trade` state.

```shell
ggezchaind query trade --help
```

##### trade-index

The `trade-index` command allows users to query trade-index.

```shell
ggezchaind query trade trade-index [flags]
```

Example:

```shell
ggezchaind query trade trade-index
```

Example Output:

```yml
trade_index:
  next_id: "1"
```

##### stored-trades

The `stored-trades` command allows users to query all stored-trades.

```shell
ggezchaind query trade stored-trades [flags]
```

Example:

```shell
ggezchaind query trade stored-trades
```

Example Output:

```yml
pagination:
  total: "1"
stored_trade:
- amount:
    amount: "162075000000000"
    denom: uggz
  banking_system_data: '{}'
  coin_minting_price_json: '{}'
  create_date: "2025-06-17T07:21:33Z"
  exchange_rate_json: '{}'
  maker: ggez1..
  coin_minting_price_usd: "1.2e-11"
  process_date: "2025-06-17T07:21:33Z"
  receiver_address: ggez1..
  result: trade created successfully
  status: TRADE_STATUS_PENDING
  trade_data: '{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low
    Carbon Target ETF","issuer":"Blackrock","no_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity:
    Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive
    Brokers LLC","type":"Brokerage Firm","country":"US"}}'
  trade_index: "1"
  trade_type: TRADE_TYPE_BUY
  update_date: "2025-06-17T07:21:33Z"
```

##### stored-trade

The `stored-trade` command allows users to query all stored-trade.

```shell
ggezchaind query trade stored-trade [trade-index] [flags]
```

Example:

```shell
ggezchaind query trade stored-trade 1
```

Example Output:

```yml
stored_trade:
  amount:
    amount: "162075000000000"
    denom: uggz
  banking_system_data: '{}'
  coin_minting_price_json: '{}'
  create_date: "2025-06-17T07:21:33Z"
  exchange_rate_json: '{}'
  maker: ggez1..
  coin_minting_price_usd: "1.2e-11"
  process_date: "2025-06-17T07:21:33Z"
  receiver_address: ggez1..
  result: trade created successfully
  status: TRADE_STATUS_PENDING
  trade_data: '{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low
    Carbon Target ETF","issuer":"Blackrock","no_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity:
    Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive
    Brokers LLC","type":"Brokerage Firm","country":"US"}}'
  trade_index: "1"
  trade_type: TRADE_TYPE_BUY
  update_date: "2025-06-17T07:21:33Z"

```

##### stored-temp-trades

The `stored-temp-trades` command allows users to query all stored-temp-trades.

```shell
ggezchaind query trade stored-temp-trades [flags]
```

Example:

```shell
ggezchaind query trade stored-temp-trades
```

Example Output:

```yml
pagination:
  total: "1"
stored_temp_trade:
- create_date: "2025-06-17T07:21:33Z"
  trade_index: "1"
```

##### stored-temp-trade

The `stored-temp-trade` command allows users to query all stored-temp-trade.

```shell
ggezchaind query trade stored-temp-trade [trade-index] [flags]
```

Example:

```shell
ggezchaind query trade stored-temp-trade 1
```

Example Output:

```yml
stored_temp_trade:
  create_date: "2025-06-17T07:21:33Z"
  trade_index: "1"
```

---

#### Transactions

The `tx` commands allow users to interact with the `trade` module.

```shell
ggezchaind tx trade --help
```

##### create-trade

The `create-trade` command create the `StoredTrade`. Must have authority to do so.

```shell
ggezchaind tx trade create-trade [trade-data] [banking-system-data] [coin-minting-price-json] [exchange-rate-json] [receiver-address] [flags]
```

Example:

```shell
ggezchaind tx trade create-trade '{"trade_info":{"asset_holder_id":1,"asset_id":1,"trade_type":1,"trade_value":1944.9,"base_currency":"USD","settlement_currency":"USD","exchange_rate":1,"exchange":"US","fund_name":"Low Carbon Target ETF","issuer":"Blackrock","no_shares":10,"coin_minting_price_usd":0.000000000012,"quantity":{"amount":"162075000000000","denom":"uggz"},"segment":"Equity: Global Low Carbon","share_price":194.49,"ticker":"CRBN","trade_fee":0,"trade_net_price":194.49,"trade_net_value":1944.9},"brokerage":{"name":"Interactive Brokers LLC","type":"Brokerage Firm","country":"US"}}' '{}' '{}' '{}' ggez1..
```

##### process-trade

The `process-trade` command process the `StoredTrade`. Must have authority to do so.

```shell
ggezchaind tx trade process-trade [trade-index] [process-type] [flags]
```

Example:

```shell
ggezchaind tx trade process-trade 1 confirm
```