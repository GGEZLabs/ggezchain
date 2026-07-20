#!/usr/bin/env bash

set -euo pipefail

command -v ggezchaind >/dev/null 2>&1 || {
    echo "ggezchaind not found in PATH"
    exit 1
}

echo "Initializing ggezchaind with chain ID 'ggezchain' and default denomination 'uggez1'..."
echo "later rookie jazz alter minute group share scan random try brain brain task afraid roast fuel ring autumn awake diamond length sadness please round" | ggezchaind keys add alice --recover --keyring-backend test
echo "float elder spice lamp blue cause office surge reopen brass stone garbage pistol noodle breeze fortune jewel flavor asthma dirt rubber lyrics vicious picnic" | ggezchaind keys add bob --recover --keyring-backend test

ggezchaind init ggezchain --default-denom uggez1 --chain-id ggezchain

ALICE_ADDR=$(ggezchaind keys show alice -a --keyring-backend test)
BOB_ADDR=$(ggezchaind keys show bob -a --keyring-backend test)

echo "Adding genesis account "$ALICE_ADDR" with balance '1000000000000uggez1'..."
ggezchaind genesis add-genesis-account "$ALICE_ADDR" 1000000000000uggez1
echo "Adding genesis account "$BOB_ADDR" with balance '1000000000000uggez1'..."
ggezchaind genesis add-genesis-account "$BOB_ADDR" 1000000000000uggez1

echo "Generating genesis transaction for 'alice' with amount '12222222uggez1'..."
ggezchaind genesis gentx alice 12222222uggez1 --keyring-backend test --chain-id ggezchain

echo "Collecting genesis transactions..."
ggezchaind genesis collect-gentxs 

echo "Setting minimum gas prices to '0uggez1'..."
ggezchaind config set app minimum-gas-prices 0uggez1
ggezchaind config set app api.enable true
ggezchaind config set app api.swagger true
ggezchaind config set app api.address tcp://0.0.0.0:1317
ggezchaind config set client chain-id ggezchain

echo "ggezchaind initialization complete."
