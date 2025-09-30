#!/bin/bash

set -ex

# initialize Hermes relayer configuration
mkdir -p /root/.hermes/
touch /root/.hermes/config.toml

echo $GGEZ_B_E2E_RLY_MNEMONIC >/root/.hermes/GGEZ_B_E2E_RLY_MNEMONIC.txt
echo $GGEZ_A_E2E_RLY_MNEMONIC >/root/.hermes/GGEZ_A_E2E_RLY_MNEMONIC.txt

# setup Hermes relayer configuration with non-zero gas_price
tee /root/.hermes/config.toml <<EOF
[global]
log_level = 'info'

[mode]

[mode.clients]
enabled = true
refresh = true
misbehaviour = true

[mode.connections]
enabled = false

[mode.channels]
enabled = true

[mode.packets]
enabled = true
clear_interval = 100
clear_on_start = true
tx_confirmation = true

[rest]
enabled = true
host = '0.0.0.0'
port = 3031

[telemetry]
enabled = true
host = '127.0.0.1'
port = 3001

[[chains]]
id = '$GGEZ_A_E2E_CHAIN_ID'
rpc_addr = 'http://$GGEZ_A_E2E_VAL_HOST:26657'
grpc_addr = 'http://$GGEZ_A_E2E_VAL_HOST:9090'
event_source = { mode = 'pull', interval = '1s', max_retries = 4 }
rpc_timeout = '10s'
account_prefix = 'ggez'
key_name = 'rly01-ggez-a'
store_prefix = 'ibc'
max_gas = 6000000
gas_price = { price = 2000000000, denom = 'uggez1' }
gas_multiplier = 1.5
clock_drift = '1m' # to accommodate docker containers
trusting_period = '14days'
trust_threshold = { numerator = '1', denominator = '3' }


[[chains]]
id = '$GGEZ_B_E2E_CHAIN_ID'
rpc_addr = 'http://$GGEZ_B_E2E_VAL_HOST:26657'
grpc_addr = 'http://$GGEZ_B_E2E_VAL_HOST:9090'
event_source = { mode = 'pull', interval = '1s', max_retries = 4 }
rpc_timeout = '10s'
account_prefix = 'ggez'
key_name = 'rly01-ggez-b'
store_prefix = 'ibc'
max_gas = 6000000
gas_price = { price = 2000000000, denom = 'uggez1' }
gas_multiplier = 1.5
clock_drift = '1m' # to accommodate docker containers
trusting_period = '14days'
trust_threshold = { numerator = '1', denominator = '3' }
EOF

# import keys
hermes keys add --key-name rly01-ggez-b --chain $GGEZ_B_E2E_CHAIN_ID --mnemonic-file /root/.hermes/GGEZ_B_E2E_RLY_MNEMONIC.txt
sleep 5
hermes keys add --key-name rly01-ggez-a --chain $GGEZ_A_E2E_CHAIN_ID --mnemonic-file /root/.hermes/GGEZ_A_E2E_RLY_MNEMONIC.txt
