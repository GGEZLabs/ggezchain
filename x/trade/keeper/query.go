package keeper

import (
	"github.com/GGEZLabs/testchain/x/trade/types"
)

var _ types.QueryServer = Keeper{}
