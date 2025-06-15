#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=x/trade/types/expected_keepers.go -package testutil -destination x/trade/testutil/expected_keepers_mocks.go
