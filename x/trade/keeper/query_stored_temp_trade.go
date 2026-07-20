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

func (q queryServer) ListStoredTempTrade(ctx context.Context, req *types.QueryAllStoredTempTradeRequest) (*types.QueryAllStoredTempTradeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	storedTempTrades, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.StoredTempTrade,
		req.Pagination,
		func(_ uint64, value types.StoredTempTrade) (types.StoredTempTrade, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllStoredTempTradeResponse{StoredTempTrade: storedTempTrades, Pagination: pageRes}, nil
}

func (q queryServer) GetStoredTempTrade(ctx context.Context, req *types.QueryGetStoredTempTradeRequest) (*types.QueryGetStoredTempTradeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.StoredTempTrade.Get(ctx, req.TradeIndex)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetStoredTempTradeResponse{StoredTempTrade: val}, nil
}
