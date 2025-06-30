package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AclAuthorityAll(ctx context.Context, req *types.QueryAllAclAuthorityRequest) (*types.QueryAllAclAuthorityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var aclAuthoritys []types.AclAuthority

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	aclAuthorityStore := prefix.NewStore(store, types.KeyPrefix(types.AclAuthorityKeyPrefix))

	pageRes, err := query.Paginate(aclAuthorityStore, req.Pagination, func(key []byte, value []byte) error {
		var aclAuthority types.AclAuthority
		if err := k.cdc.Unmarshal(value, &aclAuthority); err != nil {
			return err
		}

		aclAuthoritys = append(aclAuthoritys, aclAuthority)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAclAuthorityResponse{AclAuthority: aclAuthoritys, Pagination: pageRes}, nil
}

func (k Keeper) AclAuthority(ctx context.Context, req *types.QueryGetAclAuthorityRequest) (*types.QueryGetAclAuthorityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetAclAuthority(
		ctx,
		req.Address,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetAclAuthorityResponse{AclAuthority: val}, nil
}
