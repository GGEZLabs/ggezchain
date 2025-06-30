package keeper

import (
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
)

var _ types.QueryServer = Keeper{}
