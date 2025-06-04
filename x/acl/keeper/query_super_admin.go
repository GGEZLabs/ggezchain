package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"

	"github.com/GGEZLabs/ramichain/x/acl/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) GetSuperAdmin(ctx context.Context, req *types.QueryGetSuperAdminRequest) (*types.QueryGetSuperAdminResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	val, err := q.k.SuperAdmin.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetSuperAdminResponse{SuperAdmin: val}, nil
}
