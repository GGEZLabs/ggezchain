# GGEZChain — Project Context

> Working notes for AI-assisted development on this repository.
> Location: `AIDocs/CLAUDE.md`. This is a living document — sections will be added over time.

---

## 1. What this project is

**ggezchain** is a sovereign Cosmos SDK application-chain built by
[GGEZ Labs](https://github.com/GGEZLabs) for **ReFi (Regenerative Finance)** and the
**tokenization of Real-World Assets (RWAs)**.

The thesis: real-world, income-generating assets (sustainable and socially responsible
investments) are the most credible form of collateral for on-chain value. ggezchain turns
those assets into digital tokens that are tradable and composable across DeFi, while keeping
the issuance path auditable and permissioned.

The practical consequence for the codebase: **token supply is not open-mint**. Coins enter
and leave circulation only through a reviewed, two-party trade workflow that mirrors an
off-chain banking/custody process. That workflow is the heart of the chain and lives in
`x/trade`, gated by `x/acl`.

| | |
|---|---|
| Repo | https://github.com/GGEZLabs/ggezchain |
| Module path | `github.com/GGEZLabs/ggezchain/v2` |
| Binary | `ggezchaind` |
| Chain ID (local) | `ggezchain` |
| Base denom | `uggez1` (micro — 1 GGEZ = 1,000,000 uggez1) |
| Address prefix | `ggez` |
| Current version | v2.1.3 |
| Docs | https://docs.ggez.one |

---

## 2. Tech stack

| Layer | Choice | Version |
|---|---|---|
| Language | Go | 1.25.10+ (repo builds on 1.26.x) |
| Framework | Cosmos SDK | v0.53.6 |
| Consensus | CometBFT (BFT PoS) | v0.38.21 |
| Interoperability | IBC-Go | v10.4.0 |
| Smart contracts | CosmWasm (`wasmd`) | v0.60.1 |
| Wiring | depinject / app_config (`app/app_config.go`) | — |
| Protobuf | `buf` (+ Docker for `make proto-gen`) | — |

Notable `go.mod` `replace` directives exist for `sonic`, `gin`, `goleveldb`, and
`nhooyr.io/websocket` — check them before bumping those deps.

---

## 3. Repository layout

```
app/            App wiring: app.go, app_config.go, ante.go, ibc.go, wasm.go,
                genesis, exports, and upgrade handlers (upgrades.go + upgrade/)
cmd/ggezchaind/ CLI entrypoint
x/trade/        CUSTOM — RWA trade lifecycle, mint/burn of uggez1
x/acl/          CUSTOM — access control: super admin, admins, authorities
proto/ggezchain/{trade,acl}/v1/   Protobuf source of truth for custom modules
scripts/        init-chain.sh (local devnet), mockgen.sh, makefiles/
tests/e2e/      End-to-end tests (require Docker)
testutil/       Shared test helpers
docs/           Swagger / OpenAPI assets
AIDocs/         AI + documentation workspace (this folder)
```

---

## 4. The two custom modules

### `x/trade` — RWA trade lifecycle

Implements a **maker/checker (four-eyes) workflow**. A trade is proposed by one party and
must be independently confirmed or rejected by another before it settles. Only on
confirmation does the module mint or burn `uggez1` — the `trade` module account holds
`Minter`, `Burner`, and `Staking` permissions.

**Messages** (`proto/ggezchain/trade/v1/tx.proto`):

| Msg | Role | Effect |
|---|---|---|
| `MsgCreateTrade` | maker | Creates a trade in `TRADE_STATUS_PENDING` |
| `MsgProcessTrade` | checker | `PROCESS_TYPE_CONFIRM` or `PROCESS_TYPE_REJECT` |
| `MsgUpdateParams` | governance | Module params |

**Trade statuses:** `PENDING` → `PROCESSED` \| `REJECTED` \| `CANCELED` \| `FAILED`.

**Trade types:** `BUY`, `SELL`, `SPLIT`, `REVERSE_SPLIT`, `REINVESTMENT` (buy shares with
dividends), `DIVIDENDS`, `DIVIDEND_DEDUCTION` (deduct dividends when reinvested) — i.e. the
vocabulary of a securities register, not a generic DEX.

**State:** `StoredTrade` (settled record) and `StoredTempTrade` (pending), with a global
`TradeIndex.next_id` counter. A `StoredTrade` carries a full audit trail: `maker`, `checker`,
`tx_date`, `create_date`, `update_date`, `process_date`, plus off-chain context as JSON
strings — `trade_data`, `coin_minting_price_json`, `exchange_rate_json`,
`banking_system_data`, and `result`.

### `x/acl` — access control

Defines **who may act as maker and who may act as checker, per module**. Hierarchy:
`SuperAdmin` → `AclAdmin`s → `AclAuthority`s, where an authority holds a list of
`AccessDefinition{ module, is_maker, is_checker }`.

**Messages:** `Init`, `UpdateSuperAdmin`, `AddAdmin`, `DeleteAdmin`, `AddAuthority`,
`UpdateAuthority`, `DeleteAuthority`, `UpdateParams`.

`Init` bootstraps the ACL state — the chain's permissioning must be initialized before
trades can be processed.

---

## 5. Standard SDK modules in use

`auth`, `bank`, `staking`, `distribution`, `protocolpool`, `slashing`, `evidence`, `mint`,
`gov`, `group`, `authz`, `feegrant`, `nft`, `epochs`, `consensus`, `upgrade`, `circuit`,
plus `ibc` / `transfer` and `wasm`.

---

## 6. Build, run, test

```sh
make install        # build + install ggezchaind into $(go env GOPATH)/bin
make build          # build into ./build/ggezchaind
make test-unit      # unit tests
make test-race      # unit tests with -race
make test-cover     # coverage profile
make test-e2e       # e2e tests — REQUIRES Docker (builds ggezlabs/ggezchain image)
make mocks          # regenerate mocks (mockgen v0.6.0)
make proto-gen      # regenerate protobufs — REQUIRES Docker
make format         # gofumpt + misspell + gci
```

Builds carry the `ledger` build tag by default, so **gcc must be present**
(or set `LEDGER_ENABLED=false`).

### Local devnet

```sh
export PATH="$HOME/go/bin:$PATH"
bash scripts/init-chain.sh          # creates alice + bob (test keyring), genesis, gentx
ggezchaind start
```

`scripts/init-chain.sh` sets chain-id `ggezchain`, denom `uggez1`, funds `alice` and `bob`
with 1,000,000 GGEZ each, sets `minimum-gas-prices` to `0uggez1`, and enables the REST API
and Swagger.

**Default endpoints:** RPC `:26657` · REST/Swagger `:1317` · gRPC `:9090` · chain home `~/.ggezchain`

> The `alice`/`bob` mnemonics in `init-chain.sh` are **public test keys**. Never use them
> outside a local devnet.

---

## 7. Upgrades

On-chain upgrade handlers live in `app/upgrades.go` with per-version packages under
`app/upgrade/`: `v2_0_0`, `v2_1_0`, `v2_1_1`, `v2_1_3`. Any state-breaking change needs a new
versioned handler registered there. `app/upgrade_migrations_test.go` covers migrations.

---

## 8. Branches

`main` (default) · `dev` · `local` · `v53` · `ggezchain-cosmos-v0.53` ·
`ggezchain-evm` (EVM track) · `ggezchain-manual-wiring` (non-depinject wiring track)

---

## 9. Conventions & gotchas

- **Proto is the source of truth.** Change `proto/ggezchain/**` first, then `make proto-gen`
  (needs Docker) — never hand-edit generated `*.pb.go`.
- **Respect the maker/checker split.** The same address must not be able to both create and
  confirm a trade; that separation is the chain's core security property, and any change
  touching `x/trade` or `x/acl` should be evaluated against it.
- **Mint/burn is privileged.** Only `x/trade` should move total supply, and only via a
  confirmed trade.
- **Amounts are integers in `uggez1`.** Convert at the presentation layer, never in state.
- Lint config: `.golangci.yml`; import ordering enforced by `gci` via `make format`.
- Spelling is CI-checked (`.codespellrc`); markdown by `.markdownlint.yml`.

---

## 10. TODO / to expand

- [ ] Trade validation rules and `trade_data` JSON schema
- [ ] ACL bootstrap runbook (`MsgInit` → admins → authorities)
- [ ] End-to-end worked example: create → confirm → mint
- [ ] Testnet/mainnet endpoints and chain IDs
- [ ] CosmWasm contract deployment notes
- [ ] IBC channel/relayer setup
- [ ] Release & upgrade-proposal checklist
