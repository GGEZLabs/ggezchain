package keeper

import (
	"context"

	"github.com/GGEZLabs/testchain/x/trade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateTrade(goCtx context.Context, msg *types.MsgCreateTrade) (*types.MsgCreateTradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgCreateTradeResponse{}, nil
}
