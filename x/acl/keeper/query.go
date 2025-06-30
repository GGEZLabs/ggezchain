package keeper

import (
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
)

var _ types.QueryServer = Keeper{}
