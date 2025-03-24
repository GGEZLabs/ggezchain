package keeper

import (
	"context"

	"github.com/GGEZLabs/ggezchain/x/trade/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (k Keeper) StoredTempTradeAll(ctx context.Context, req *types.QueryAllStoredTempTradeRequest) (*types.QueryAllStoredTempTradeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var storedTempTrades []types.StoredTempTrade

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	storedTempTradeStore := prefix.NewStore(store, types.KeyPrefix(types.StoredTempTradeKeyPrefix))

	pageRes, err := query.Paginate(storedTempTradeStore, req.Pagination, func(key []byte, value []byte) error {
		var storedTempTrade types.StoredTempTrade
		if err := k.cdc.Unmarshal(value, &storedTempTrade); err != nil {
			return err
		}

		storedTempTrades = append(storedTempTrades, storedTempTrade)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllStoredTempTradeResponse{StoredTempTrade: storedTempTrades, Pagination: pageRes}, nil
}

func (k Keeper) StoredTempTrade(ctx context.Context, req *types.QueryGetStoredTempTradeRequest) (*types.QueryGetStoredTempTradeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetStoredTempTrade(
		ctx,
		req.TradeIndex,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetStoredTempTradeResponse{StoredTempTrade: val}, nil
}
