package keeper

import (
	"github.com/GGEZLabs/ggezchain/x/trade/types"
)

var _ types.QueryServer = Keeper{}
