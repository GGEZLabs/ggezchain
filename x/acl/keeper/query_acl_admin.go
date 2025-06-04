package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"github.com/GGEZLabs/ramichain/x/acl/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListAclAdmin(ctx context.Context, req *types.QueryAllAclAdminRequest) (*types.QueryAllAclAdminResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	aclAdmins, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.AclAdmin,
		req.Pagination,
		func(_ string, value types.AclAdmin) (types.AclAdmin, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAclAdminResponse{AclAdmin: aclAdmins, Pagination: pageRes}, nil
}

func (q queryServer) GetAclAdmin(ctx context.Context, req *types.QueryGetAclAdminRequest) (*types.QueryGetAclAdminResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.AclAdmin.Get(ctx, req.Address)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetAclAdminResponse{AclAdmin: val}, nil
}
