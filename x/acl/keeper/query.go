package keeper

import (
	"github.com/GGEZLabs/ggezchain/x/acl/types"
)

var _ types.QueryServer = Keeper{}
