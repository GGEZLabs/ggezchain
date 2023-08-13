package keeper

import (
	"context"

	"github.com/GGEZLabs/testchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ProcessTrade(goCtx context.Context, msg *types.MsgProcessTrade) (*types.MsgProcessTradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgProcessTradeResponse{}, nil
}
