# GGEZChain Documentation

## 1. Local Setup & Running the Node

### 1.1 Install Go

GGEZChain requires **Go 1.25.0** or later.

- Download and install Go from the official site: https://go.dev/dl/
- Verify the installation:

  ```bash
  go version
  # should print: go version go1.25.0 (or later)
  ```

- Make sure `$GOPATH/bin` is on your `PATH` (needed later to run the installed binary):

  ```bash
  export GOROOT=/usr/local/go
  export GOPATH=$HOME/go
  export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
  ```

### 1.2 Clone the Repository

```bash
git clone https://github.com/GGEZLabs/ggezchain.git
cd ggezchain
```

### 1.3 Build / Install

Two options, depending on whether you want the binary in a local `build/` folder or installed to your `GOPATH`:

```bash
make install   # installs ggezchaind to $GOPATH/bin
make build     # builds binary to ./build/ggezchaind
```

Verify the binary is available:

```bash
ggezchaind version
# or, if using the local build:
./build/ggezchaind version
```

### 1.4 Run the Node

```bash
./scripts/init-chain.sh
ggezchaind start
```

This starts the node using the configuration and genesis found in the default home directory (`~/.ggezchain` unless overridden with `--home`).

### 1.5 Resetting the Node

Use this when you want to wipe local chain state and re-sync from scratch (e.g. after a corrupted database, or to switch networks/genesis on the same machine).

```bash
ggezchaind comet reset-state
# or, to also remove address book and config overrides:
ggezchaind comet unsafe-reset-all
```

- `comet reset-state` — removes the blockchain database (`data/`) but leaves `config/` untouched, including `addrbook.json`.
- `comet unsafe-reset-all` — removes `data/` **and** resets `config/addrbook.json` and other local state; as the name warns, this is a broader wipe.

> **Warning:** neither command touches `config/priv_validator_key.json` or `config/node_key.json`, so your validator/node identity is preserved. But both **do** reset `data/priv_validator_state.json` back to height 0. If this node is an active validator on a live chain, resetting that file and then restarting can cause **double-signing** (the node may re-sign blocks/votes it already signed at heights it now thinks are new) — never run these commands against a live validator without understanding that risk. They're intended for local/test networks or a node that's fully caught up in a fresh restore (see [Section 5.6](#56-restore-the-validator-key)).
>
> After resetting, you'll need a fresh `genesis.json` (see [Section 1.4](#14-run-the-node) or [Section 4.5](#45-sync-genesis-with-the-network)) before starting the node again.

### 1.6 Recovering a Key from a Mnemonic

If you already have a key's 24-word mnemonic (e.g. from a backup) and need to re-add it to the keyring, use `--recover` instead of generating a new key:

```bash
ggezchaind keys add KEY_NAME \
  --keyring-backend [os|file|kwallet|pass|test|memory] \
  --keyring-dir /path/to/keys/file \
  --recover
```

This prompts you to type/paste in the existing mnemonic instead of creating a new one, and re-derives the same key/address from it.

- `KEY_NAME` — local name to give the recovered key in the keyring (can be anything; it doesn't need to match the original name used when the key was first created).
- `--keyring-backend` — must match the backend used elsewhere for this node/user (`file` is used throughout this doc, e.g. [Section 4.6](#46-add-the-validator-key)).
- `--keyring-dir` — path to the keyring directory; omit to use the default (`~/.ggezchain/keys` when `--keyring-backend file` and no `--home` override).

> **Warning:** treat the mnemonic as a secret on par with the private key itself — anyone with it can recover the same key and control its funds/validator identity.

## 2. Configuring Cosmovisor

[Cosmovisor](https://docs.cosmos.network/main/build/tooling/cosmovisor) manages the `ggezchaind` binary and automates upgrades. It wraps the node process and swaps in new binaries when an upgrade is triggered on-chain.

### 2.1 Install Cosmovisor

```bash
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
# or pin a specific version:
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.6.0
```

Verify:

```bash
cosmovisor version
```

### 2.2 Required Environment Variables

Cosmovisor needs these set before `init` and `run` will work:

```bash
export DAEMON_NAME=ggezchain
export DAEMON_HOME=$HOME/.ggezchain
export DAEMON_RESTART_AFTER_UPGRADE=true
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
```

### 2.3 Initialize Cosmovisor

```bash
cosmovisor init /path/to/ggezchaind
```

This sets up the Cosmovisor directory layout under `$DAEMON_HOME/cosmovisor`, copying the given binary in as the genesis (`genesis/bin`) version:

```
$DAEMON_HOME/cosmovisor/
├── genesis/bin/ggezchaind
├── upgrades/<upgrade-name>/bin/ggezchaind
└── current -> genesis (or upgrades/<upgrade-name>)
```

### 2.4 Run the Node via Cosmovisor

```bash
cosmovisor run start
```

Cosmovisor runs `ggezchaind start` under the hood, monitors the chain for scheduled upgrades, and automatically switches to the upgrade binary in `$DAEMON_HOME/cosmovisor/upgrades/<upgrade-name>/bin` once it's staged and the halt height is reached.

## 3. Debugging

The repo includes a `.vscode/launch.json` with two Go debug configurations:

### 3.1 Attach to a Running Process

Use this to debug a node that's already running (e.g. started via `ggezchaind start` or `cosmovisor run start`). VS Code prompts you to pick the process to attach to.

```json
{
  "name": "Attach to Running Process",
  "type": "go",
  "request": "attach",
  "mode": "local",
  "processId": "${command:pickProcess}"
}
```

To use it: start the node, open the "Run and Debug" panel in VS Code, select **Attach to Running Process**, and pick the `ggezchaind` process from the list.

### 3.2 Debug a Sample Command

Use this to step through a specific CLI command, e.g. `ggezchaind export`. Change `args` to debug a different subcommand.

```json
{
  "name": "Debug Command (export)",
  "type": "go",
  "request": "launch",
  "mode": "debug",
  "program": "${workspaceFolder}/cmd/ggezchaind",
  "args": [
    "export"
  ]
}
```

## 4. Adding a New Validator from Scratch

This walks through setting up a brand-new validator node and joining a running network.

### 4.1 Install Go

```bash
sudo apt update
wget https://go.dev/dl/go1.25.10.linux-amd64.tar.gz
sha256sum go1.25.10.linux-amd64.tar.gz   # compare against the checksum published on go.dev/dl
sudo tar -C /usr/local -xzf go1.25.10.linux-amd64.tar.gz
```

Add Go to `.profile`:

```bash
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

### 4.2 Install Cosmovisor

```bash
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.7.1
```

### 4.3 Install the Network Binary

```bash
wget https://github.com/GGEZLabs/ggezchain/releases/download/v1.0.0/ggezchaind-v1.0.0-linux-amd64 -O go/bin/ggezchaind
chmod +x /home/ggez/go/bin/ggezchaind
ggezchaind version
```

### 4.4 Initialize the Node

```bash
ggezchaind init ggez --chain-id ggezchain-testnet --default-denom uggez1
```

### 4.5 Sync Genesis with the Network

Export the genesis file from an existing node (after stopping it):

```bash
ggezchaind export --for-zero-height --output-document ./genesis.json
```

Copy that `genesis.json` over to replace the one generated by `init` in this new node's `~/.ggezchain/config/`.

### 4.6 Add the Validator Key

```bash
ggezchaind keys add ggez --keyring-backend file --keyring-dir /home/ggez/.ggezchain/keys
```

### 4.7 Set Cosmovisor Environment Variables

Add to `.profile`:

```bash
export DAEMON_NAME=ggezchain
export DAEMON_HOME=$HOME/.ggezchain
export DAEMON_RESTART_AFTER_UPGRADE=true
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
```

```bash
source .profile
```

### 4.8 Initialize Cosmovisor

```bash
cosmovisor init go/bin/ggezchaind
```

### 4.9 Enable State Sync

Confirm at least one peer on the network has snapshots enabled, then set the following in `~/.ggezchain/config/config.toml` under `[statesync]`:

```toml
enable = true
rpc_servers = "http://172.23.100.40:26657,http://172.23.100.50:26657"
trust_height = 942316
trust_hash = "FF28A050FE9E72BA849C044CD8B79A00F11FB09B9E1CA7E198CD9ED4FA4877FD"
trust_period = "168h0m0s"
```

`rpc_servers` should point to two live RPC endpoints on the network. `trust_height`/`trust_hash` should come from a recent block (within `trust_period`) on those nodes.

> **Important:** the nodes listed in `rpc_servers` must have snapshotting enabled in their own `~/.ggezchain/config/app.toml` under `[state-sync]` — if they don't, they have no snapshots to serve and your node will fail to state-sync:
>
> ```toml
> [state-sync]
>
> # snapshot-interval specifies the block interval at which local state sync snapshots are
> # taken (0 to disable).
> snapshot-interval = 1000
>
> # snapshot-keep-recent specifies the number of recent snapshots to keep and serve (0 to keep all).
> snapshot-keep-recent = 2
> ```

### 4.10 Create the systemd Service

```bash
echo "[Unit]
Description=GGEZ Chain
After=network-online.target

[Service]
User=ggez
ExecStart=/bin/bash -l -c 'cosmovisor run start'
Environment=DAEMON_NAME=ggezchain DAEMON_HOME=/home/ggez/.ggezchain DAEMON_RESTART_AFTER_UPGRADE=true DAEMON_ALLOW_DOWNLOAD_BINARIES=true
Restart=always
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target" | sudo tee /etc/systemd/system/ggezchain.service > /dev/null
```

> Note: `systemd`'s `Environment=` directive doesn't expand `$HOME`, so the absolute path is used here instead.

### 4.11 Enable & Start the Service

```bash
sudo systemctl enable ggezchain.service
sudo systemctl start ggezchain.service
```

Check status/logs:

```bash
sudo systemctl status ggezchain.service
journalctl -u ggezchain.service -f
```

### 4.12 Fund the Validator Address

Send tokens to the new validator's address from an already-funded account so it has enough balance to self-delegate and cover fees.

### 4.13 Get the Validator's Public Key

```bash
ggezchaind comet show-validator
```

This prints the node's consensus `pubkey` — use it in `validator.json` below.

### 4.14 Create `validator.json`

```json
{
    "pubkey": {
        "@type": "/cosmos.crypto.ed25519.PubKey",
        "key": "rtJsf2QK+f7CzYVxilSGMMpaPNOwUY2mhuHmjG7wgCE="
    },
    "amount": "200000000000uggez",
    "moniker": "GGEZ Admin",
    "website": "https://ggez.com",
    "commission-rate": "0.1",
    "commission-max-rate": "0.2",
    "commission-max-change-rate": "0.01",
    "min-self-delegation": "1"
}
```

- `pubkey` — the consensus pubkey from step 4.13.
- `amount` — self-delegation amount, in `uggez`.
- `commission-rate` / `commission-max-rate` / `commission-max-change-rate` — commission parameters (max change per day).
- `min-self-delegation` — minimum self-delegated stake the validator must maintain.

### 4.15 Submit the Create-Validator Transaction

```bash
ggezchaind tx staking create-validator ./validator.json \
  --from ggez1dgw3amkw9egu7wdhsqxguw30sj2hpwftdsm8c7 \
  --keyring-dir /home/ggez/.ggezchain/keys \
  --keyring-backend file \
  --fees 125000uggez \
  --gas 250000 \
  -y
```

`--from` can be the key name (e.g. `ggez`, from step 4.6) or its address, as shown here.

### 4.16 Verify the Validator

```bash
ggezchaind query staking validator ggezvaloper1dgw3amkw9egu7wdhsqxguw30sj2hpwftl62lcf
```

A successful result shows the validator with `status: BOND_STATUS_BONDED` (once it's in the active set) and the delegation/commission details from `validator.json`.

## 5. Restoring a Validator

Used when moving an existing, already-registered validator to a new machine (e.g. hardware migration or disaster recovery) rather than creating a new one.

> **Note:** This only works if you have a backup of the validator's `~/.ggezchain/config/priv_validator_key.json`. Without it, the validator's on-chain identity can't be recovered — you'd have to create a brand-new validator instead (see [Section 4](#4-adding-a-new-validator-from-scratch)).

### 5.1 Install Go

```bash
sudo apt update
wget https://go.dev/dl/go1.23.7.linux-amd64.tar.gz
sha256sum go1.23.7.linux-amd64.tar.gz   # compare against the checksum published on go.dev/dl
sudo tar -C /usr/local -xzf go1.23.7.linux-amd64.tar.gz
```

### 5.2 Configure `.profile`

```bash
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

export DAEMON_NAME=ggezchain
export DAEMON_HOME=$HOME/.ggezchain
export DAEMON_RESTART_AFTER_UPGRADE=true
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
```

```bash
source .profile
```

### 5.3 Install Cosmovisor

```bash
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.6.0
```

### 5.4 Download the Network Binary

```bash
wget https://github.com/GGEZLabs/ggezchain/releases/download/v2.0.0/ggezchaind-v2.0.0-linux-amd64 -O go/bin/ggezchaind
chmod +x go/bin/ggezchaind
ggezchaind version
```

### 5.5 Initialize the Node

```bash
ggezchaind init ggezchain-val-4-13 --chain-id ggezchain --default-denom uggez1
```

### 5.6 Restore the Validator Key

Copy the backed-up key files into the new node's config directory, overwriting the ones generated by `init`:

```bash
cp /path/to/backup/priv_validator_key.json ~/.ggezchain/config/priv_validator_key.json
cp /path/to/backup/node_key.json ~/.ggezchain/config/node_key.json   # optional, preserves the same node/peer identity
```

This is the step that makes it a *restore* rather than a new validator — the consensus pubkey in `priv_validator_key.json` is what ties this node back to the validator's existing on-chain registration.

### 5.7 Sync Genesis with the Network

```bash
ggezchaind export --for-zero-height --output-document ./genesis.json
```

Copy the exported `genesis.json` over the one in `~/.ggezchain/config/`, replacing the one from `init`.

### 5.8 Review the Configuration Files

Check `~/.ggezchain/config/config.toml` and `app.toml` against the network's recommended settings (persistent peers/seeds, pruning, minimum gas price, ports) before starting the node.

### 5.9 Initialize Cosmovisor

```bash
cosmovisor init go/bin/ggezchaind
```

### 5.10 Enable State Sync

Confirm at least one peer on the network has snapshots enabled, then set the following in `config.toml` under `[statesync]`:

```toml
enable = true
rpc_servers = "https://bc.ggez.one:8443,http://bc-2-2.ggez.one:8443"
trust_height = 7685135
trust_hash = "02CC1115CFDE1C1320260A8C07F9FD6934065A0CE941882C12BC148A88DAC459"
trust_period = "168h0m0s"
```

> **Important:** the nodes listed in `rpc_servers` must have snapshotting enabled in their own `~/.ggezchain/config/app.toml` under `[state-sync]` — if they don't, they have no snapshots to serve and your node will fail to state-sync:
>
> ```toml
> [state-sync]
>
> # snapshot-interval specifies the block interval at which local state sync snapshots are
> # taken (0 to disable).
> snapshot-interval = 1000
>
> # snapshot-keep-recent specifies the number of recent snapshots to keep and serve (0 to keep all).
> snapshot-keep-recent = 2
> ```

### 5.11 Create the systemd Service

```bash
echo "[Unit]
Description=GGEZ Chain
After=network-online.target

[Service]
User=mutaz_ggez_one
ExecStart=/bin/bash -l -c 'cosmovisor run start'
Environment=DAEMON_NAME=ggezchain DAEMON_HOME=/home/mutaz_ggez_one/.ggezchain DAEMON_RESTART_AFTER_UPGRADE=true DAEMON_ALLOW_DOWNLOAD_BINARIES=true
Restart=always
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target" | sudo tee /etc/systemd/system/ggezchain.service > /dev/null
```

### 5.12 Enable & Start the Service

```bash
sudo systemctl enable ggezchain.service
sudo systemctl start ggezchain.service
```

Check status/logs:

```bash
sudo systemctl status ggezchain.service
journalctl -u ggezchain.service -f
```

## 6. Testing a Software Upgrade Locally

This walks through simulating a governance-gated software upgrade (e.g. `v2.1.0`) on a local single-node network via Cosmovisor.

### 6.1 Initialize the Current Version

Init the node as usual for the current binary/version (see [Section 1](#1-local-setup--running-the-node)), then init Cosmovisor with that binary:

```bash
cosmovisor init /home/ggez/go/bin/ggezchaind
```

### 6.2 Speed Up Governance Params

To avoid waiting on real voting/deposit periods, edit `app_state.gov.params` in `~/.ggezchain/config/genesis.json` before starting the node:

```json
"max_deposit_period": "20s",
"voting_period": "20s",
"expedited_voting_period": "10s",
```

### 6.3 Build the New Binary

```bash
make build
```

This produces the upgrade target at `build/ggezchaind`.

### 6.4 Run the Upgrade Flow

Open a terminal to run the node, and use a second terminal for the commands below — keep the first terminal open so you can watch the upgrade happen in its logs.

**Terminal 1 — start the network:**

```bash
cosmovisor run start
```

**Terminal 2 — stage the upgrade binary:**

```bash
cosmovisor add-upgrade v2.1.0 build/ggezchaind
```

**Terminal 2 — submit the upgrade proposal:**

```bash
ggezchaind tx upgrade software-upgrade v2.1.0 \
  --upgrade-height 10 \
  --title upgrade \
  --summary upgrade \
  --no-validate \
  --deposit 10000000uggez1 \
  --from ggez1q3sfaepes35ly4sa5ppguf6gs49un4uz858e33 \
  -y
```

**Terminal 2 — vote on the proposal:**

```bash
ggezchaind tx gov vote 1 yes --from ggez1q3sfaepes35ly4sa5ppguf6gs49un4uz858e33 -y
```

Replace `1` with the actual proposal ID returned when the proposal was submitted.

### 6.5 Verify the Upgrade

```bash
ggezchaind query gov proposal 1
```

Then watch Terminal 1's logs: once the chain reaches block height `10`, Cosmovisor should halt the old binary and automatically swap in `build/ggezchaind` (staged as `v2.1.0` in step 6.4), then resume block production on the new version.

## 7. Creating a New GitHub Release

Releases are built with [GoReleaser](https://goreleaser.com/) (config in `.goreleaser.yaml`), which cross-compiles Linux, ARM, and macOS binaries in one pass — macOS builds require an osxcross toolchain.

### 7.1 Prerequisite: Set Up osxcross

Cross-compiling the macOS/Darwin binaries requires [tpoechtrager/osxcross](https://github.com/tpoechtrager/osxcross) installed on the Linux (or *BSD) machine doing the release.

### 7.2 Confirm `main` Is Up to Date

Make sure the latest intended code is merged into `main` before tagging a release.

### 7.3 Draft the Release on GitHub

Go to **Releases → Draft new release** on the [GGEZLabs/ggezchain](https://github.com/GGEZLabs/ggezchain/releases) repo and create a new tag, e.g. `v2.1.1`.

### 7.4 Check Out `main` Locally

```bash
git checkout main
git pull
```

### 7.5 Set the GitHub Token

```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

Needs `repo` scope so GoReleaser can push the release assets.

### 7.6 Run GoReleaser

```bash
sudo -E PATH=$PATH COSMWASM_VERSION=v2.2.4 TM_VERSION=v0.14.1 REPO_OWNER=GGEZLabs REPO_NAME=ggezchain goreleaser release --clean
```

- `sudo -E` preserves the current environment (including `GITHUB_TOKEN` and `PATH`) while running as root, which the cross-compile toolchain needs.
- `COSMWASM_VERSION` — pins the `libwasmvm` static libs downloaded for each target in `.goreleaser.yaml`.
- `TM_VERSION` — CometBFT version to embed.
- `--clean` — wipes any previous `dist/` output before building.

If the token is valid and the working tree has no pending/uncommitted changes, GoReleaser builds all targets and publishes the release directly to GitHub — see [Release ggezchain v2.1.1](https://github.com/GGEZLabs/ggezchain/releases/tag/v2.1.1) for an example of a successful run.

## 8. Creating a Software Upgrade on Live

This is the mainnet/testnet counterpart to [Section 6](#6-testing-a-software-upgrade-locally): submitting a real `MsgSoftwareUpgrade` governance proposal against a live network, pointing at binaries published via [Section 7](#7-creating-a-new-github-release).

### 8.1 Connect to a Node

Go to any full node with access to the network you're upgrading (mainnet or testnet).

### 8.2 Estimate the Upgrade Height

`scripts/estimate-upgrade-height.sh` queries the node to work out a safe `plan.height`: it reads the current height, computes the average block time from recent blocks, fetches the voting period from gov params, and adds a buffer on top of the estimated vote-end height.

```bash
./scripts/estimate-upgrade-height.sh \
  --rpc http://localhost:26657 \
  --rest http://localhost:1317
```

Optional flags: `--sample <n>` (blocks used to compute average block time, default 100), `--buffer <n>` (extra blocks added after the vote-end estimate, default 500), `--expedited` (use `expedited_voting_period` instead of `voting_period`). It prints a recommended `plan.height` — use that value in step 8.5.

### 8.3 Draft the Proposal

```bash
ggezchaind tx gov draft-proposal --skip-metadata
```

This launches an interactive prompt — choose **software-upgrade** as the proposal type. It writes out a `draft_proposal.json` for you to fill in and review.

### 8.4 Generate the Binaries Info

`./scripts/create_binaries.py` fetches the `sha256sum.txt` published with a GitHub release (see [Section 7](#7-creating-a-new-github-release)) and builds the escaped `plan.info` JSON string for `linux/amd64`, `linux/arm64`, `darwin/amd64`, and `darwin/arm64`:

```bash
python3 scripts/create_binaries.py v2.1.1
```

Paste its output as the value of `plan.info` in the proposal file below.

### 8.5 Fill In the Proposal File

Edit the generated file into something like this, using the `plan.height` from step 8.2 and the `plan.info` from step 8.4:

```json
{
 "messages": [
  {
   "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
   "authority": "ggez10d07y265gmmuvt4z0w9aw880jnsr700jm4garg",
   "plan": {
    "name": "v2.1.0",
    "time": "0001-01-01T00:00:00Z",
    "height": "10",
    "info": "{\"binaries\": {\"linux/amd64\": \"https://github.com/GGEZLabs/ggezchain/releases/download/v2.1.1/ggezchaind-v2.1.1-linux-amd64.tar.gz?checksum=sha256:e537f2dc2fe16f9a3ed1cf7dd4e97a7607b5d70e8fbed9278f783e7947fcc12c\", \"linux/arm64\": \"https://github.com/GGEZLabs/ggezchain/releases/download/v2.1.1/ggezchaind-v2.1.1-linux-arm64.tar.gz?checksum=sha256:fbbed2b1e34b23f5c0714619c5b0faea84db91c2fd2f5db3d2bcc89e8aca7e10\", \"darwin/amd64\": \"https://github.com/GGEZLabs/ggezchain/releases/download/v2.1.1/ggezchaind-v2.1.1-darwin-amd64.tar.gz?checksum=sha256:575d41b0ef7c198ce40cd00b05b6b841d5368bf5a84285973b12b7846e567411\", \"darwin/arm64\": \"https://github.com/GGEZLabs/ggezchain/releases/download/v2.1.1/ggezchaind-v2.1.1-darwin-arm64.tar.gz?checksum=sha256:0585bb934b71ca38f4150c410ed05c67af5860c99f1e4b2b9396dbc10434274f\"}}",
    "upgraded_client_state": null
   }
  }
 ],
 "metadata": "ipfs://CID",
 "deposit": "10000000uggez1",
 "title": "Software Upgrade v2.1.1",
 "summary": "**GGEZ Chain – Software Upgrade Proposal (v2.1.1)**\n\nThis proposal requests approval to upgrade the GGEZChain mainnet from v2.0.0 to v2.1.1.\n\nKey changes included in this release:\n\nState machine / module changes:\n- Upgraded Cosmos SDK to v0.53, bringing updated app wiring, ante handler, IBC module setup, genesis export, and simulation tooling in line with the new SDK version\n- Added the `epochs` and `protocolpool` modules to the chain\n- Disabled `EnableUnorderedTransactions` in the auth module configuration\n- Registered the v2.1.0 and v2.1.1 upgrade handlers, which run standard module migrations (no manual data migrations required)\n\nFixes:\n- Fixed genesis export to honor the specified height, ensuring `export` produces a consistent state snapshot for state-sync and chain-recovery procedures",
 "expedited": false
}
```

Key fields to double-check before submitting:

- `authority` — must be the `x/gov` module account address, not a regular key.
- `plan.name` — must exactly match the upgrade handler name registered in `app/upgrades.go` for the target binary.
- `plan.height` — the block height the upgrade halts at; set far enough ahead to allow the voting period to complete.
- `plan.info.binaries` — one URL per platform, each with a `?checksum=sha256:<hash>` query param; Cosmovisor uses this to auto-download and verify the binary if `DAEMON_ALLOW_DOWNLOAD_BINARIES=true`.
- `deposit` — must meet the chain's minimum deposit to enter voting.

### 8.6 Submit the Proposal

```bash
ggezchaind tx gov submit-proposal proposal.json --from <your-key-name>
```

### 8.7 Check Proposal Status

```bash
ggezchaind query gov proposal <proposal-id>
```

### 8.8 Vote on the Proposal

```bash
ggezchaind tx gov vote <proposal-id> yes --from <your-key-name>
```

### 8.9 Wait for the Upgrade

Once the proposal passes, the upgrade is scheduled at `plan.height`. Validators running Cosmovisor (see [Section 2](#2-configuring-cosmovisor)) will automatically halt and swap to the new binary at that height, provided it's staged (via `cosmovisor add-upgrade`) or downloadable from `plan.info.binaries`.

## 9. Ignite CLI Version

If you're using [Ignite CLI](https://docs.ignite.com/) to scaffold a new message, module, or otherwise modify the chain (`ignite scaffold message`, `ignite scaffold module`, etc.), use this exact version:

```
v29.10.1-dev
```

Scaffolding with a different Ignite version can generate code that doesn't match this repo's SDK version (`v0.53.3`) and existing module layout, causing drift or build errors. Install/switch to it with:

```bash
go install github.com/ignite/cli/v29/ignite/cmd/ignite@v29.10.1-dev
```

Verify with:

```bash
ignite version
```
