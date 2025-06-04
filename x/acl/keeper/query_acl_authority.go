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

func (q queryServer) ListAclAuthority(ctx context.Context, req *types.QueryAllAclAuthorityRequest) (*types.QueryAllAclAuthorityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	aclAuthoritys, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.AclAuthority,
		req.Pagination,
		func(_ string, value types.AclAuthority) (types.AclAuthority, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAclAuthorityResponse{AclAuthority: aclAuthoritys, Pagination: pageRes}, nil
}

func (q queryServer) GetAclAuthority(ctx context.Context, req *types.QueryGetAclAuthorityRequest) (*types.QueryGetAclAuthorityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.AclAuthority.Get(ctx, req.Address)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetAclAuthorityResponse{AclAuthority: val}, nil
}
