package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListStoredTrade(ctx context.Context, req *types.QueryAllStoredTradeRequest) (*types.QueryAllStoredTradeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	storedTrades, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.StoredTrade,
		req.Pagination,
		func(_ uint64, value types.StoredTrade) (types.StoredTrade, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllStoredTradeResponse{StoredTrade: storedTrades, Pagination: pageRes}, nil
}

func (q queryServer) GetStoredTrade(ctx context.Context, req *types.QueryGetStoredTradeRequest) (*types.QueryGetStoredTradeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.StoredTrade.Get(ctx, req.TradeIndex)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetStoredTradeResponse{StoredTrade: val}, nil
}
