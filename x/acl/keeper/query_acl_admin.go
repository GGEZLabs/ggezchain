package keeper

import (
	"context"

	"github.com/GGEZLabs/ggezchain/x/acl/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (k Keeper) AclAdminAll(ctx context.Context, req *types.QueryAllAclAdminRequest) (*types.QueryAllAclAdminResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var aclAdmins []types.AclAdmin

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	aclAdminStore := prefix.NewStore(store, types.KeyPrefix(types.AclAdminKeyPrefix))

	pageRes, err := query.Paginate(aclAdminStore, req.Pagination, func(key []byte, value []byte) error {
		var aclAdmin types.AclAdmin
		if err := k.cdc.Unmarshal(value, &aclAdmin); err != nil {
			return err
		}

		aclAdmins = append(aclAdmins, aclAdmin)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAclAdminResponse{AclAdmin: aclAdmins, Pagination: pageRes}, nil
}

func (k Keeper) AclAdmin(ctx context.Context, req *types.QueryGetAclAdminRequest) (*types.QueryGetAclAdminResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, found := k.GetAclAdmin(
		ctx,
		req.Address,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetAclAdminResponse{AclAdmin: val}, nil
}
