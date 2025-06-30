package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/GGEZLabs/ggezchain/v2/x/trade/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) StoredTradeAll(ctx context.Context, req *types.QueryAllStoredTradeRequest) (*types.QueryAllStoredTradeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var storedTrades []types.StoredTrade

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	storedTradeStore := prefix.NewStore(store, types.KeyPrefix(types.StoredTradeKeyPrefix))

	pageRes, err := query.Paginate(storedTradeStore, req.Pagination, func(key []byte, value []byte) error {
		var storedTrade types.StoredTrade
		if err := k.cdc.Unmarshal(value, &storedTrade); err != nil {
			return err
		}

		storedTrades = append(storedTrades, storedTrade)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllStoredTradeResponse{StoredTrade: storedTrades, Pagination: pageRes}, nil
}

func (k Keeper) StoredTrade(ctx context.Context, req *types.QueryGetStoredTradeRequest) (*types.QueryGetStoredTradeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetStoredTrade(
		ctx,
		req.TradeIndex,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetStoredTradeResponse{StoredTrade: val}, nil
}
