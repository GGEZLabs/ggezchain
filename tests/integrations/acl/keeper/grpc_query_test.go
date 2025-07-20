package keeper_test

import (
	gocontext "context"
	"fmt"
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/x/acl/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"gotest.tools/v3/assert"
)

func TestGRPCQueryAclAdmin(t *testing.T) {
	f := initFixture(t)
	ctx, queryClient := f.ctx, f.queryClient

	var (
		req      *types.QueryGetAclAdminRequest
		res      *types.QueryGetAclAdminResponse
		aclAdmin types.AclAdmin
	)

	testCases := []struct {
		msg       string
		malleate  func()
		expPass   bool
		expErrMsg string
	}{
		{
			"get acl admin",
			func() {
				f.aclKeeper.SetAclAdmin(ctx, types.AclAdmin{Address: "address"})
				var found bool
				aclAdmin, found = f.aclKeeper.GetAclAdmin(ctx, "address")
				assert.Assert(t, found == true)
				assert.Assert(t, aclAdmin.String() != "")

				req = &types.QueryGetAclAdminRequest{Address: "address"}

				res = &types.QueryGetAclAdminResponse{
					AclAdmin: types.AclAdmin{
						Address: "address",
					},
				}
			},
			true,
			"",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Case %s", testCase.msg), func(t *testing.T) {
			testCase.malleate()

			aclAdmin, err := queryClient.AclAdmin(gocontext.Background(), req)

			if testCase.expPass {
				assert.NilError(t, err)
				assert.Equal(t, res.String(), aclAdmin.String())
			} else {
				assert.ErrorContains(t, err, testCase.expErrMsg)
				assert.Assert(t, aclAdmin == nil)
			}
		})
	}
}

func TestGRPCQueryAllAclAdmin(t *testing.T) {
	f := initFixture(t)
	ctx, queryClient := f.ctx, f.queryClient

	var (
		req         *types.QueryAllAclAdminRequest
		res         *types.QueryAllAclAdminResponse
		aclAdminAll []types.AclAdmin
	)

	testCases := []struct {
		msg       string
		malleate  func()
		expPass   bool
		expErrMsg string
	}{
		{
			"nil request",
			func() {
				req = nil
				res = &types.QueryAllAclAdminResponse{
					AclAdmin:   []types.AclAdmin{},
					Pagination: &query.PageResponse{},
				}
			},
			true,
			"",
		},
		{
			"get all acl admins",
			func() {
				f.aclKeeper.SetAclAdmin(ctx, types.AclAdmin{Address: "address1"})
				f.aclKeeper.SetAclAdmin(ctx, types.AclAdmin{Address: "address2"})
				f.aclKeeper.SetAclAdmin(ctx, types.AclAdmin{Address: "address3"})
				f.aclKeeper.SetAclAdmin(ctx, types.AclAdmin{Address: "address4"})

				aclAdminAll = f.aclKeeper.GetAllAclAdmin(ctx)
				assert.Assert(t, len(aclAdminAll) == 4)

				req = &types.QueryAllAclAdminRequest{}

				res = &types.QueryAllAclAdminResponse{
					AclAdmin: []types.AclAdmin{
						{Address: "address1"},
						{Address: "address2"},
						{Address: "address3"},
						{Address: "address4"},
					},
					Pagination: &query.PageResponse{
						Total: 4,
					},
				}
			},
			true,
			"",
		},
		{
			"get some of acl admins",
			func() {
				aclAdminAll = f.aclKeeper.GetAllAclAdmin(ctx)
				assert.Assert(t, len(aclAdminAll) == 4)

				req = &types.QueryAllAclAdminRequest{
					Pagination: &query.PageRequest{
						Limit: 2,
					},
				}

				res = &types.QueryAllAclAdminResponse{
					AclAdmin: []types.AclAdmin{
						{
							Address: "address1",
						},
						{
							Address: "address2",
						},
					},
				}
			},
			true,
			"",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Case %s", testCase.msg), func(t *testing.T) {
			testCase.malleate()

			aclAdmins, err := queryClient.AclAdminAll(gocontext.Background(), req)

			if testCase.expPass {
				assert.NilError(t, err)
				assert.Equal(t, len(res.AclAdmin), len(aclAdmins.AclAdmin))
			} else {
				assert.ErrorContains(t, err, testCase.expErrMsg)
				assert.Assert(t, aclAdmins == nil)
			}
		})
	}
}

func TestGRPCQueryAclAuthority(t *testing.T) {
	f := initFixture(t)
	ctx, queryClient := f.ctx, f.queryClient

	var (
		req      *types.QueryGetAclAuthorityRequest
		res      *types.QueryGetAclAuthorityResponse
		aclAdmin types.AclAuthority
	)

	testCases := []struct {
		msg       string
		malleate  func()
		expPass   bool
		expErrMsg string
	}{
		{
			"get acl authority",
			func() {
				f.aclKeeper.SetAclAuthority(ctx, types.AclAuthority{
					Address: "address",
					Name:    "Alice",
					AccessDefinitions: []*types.AccessDefinition{
						{Module: "module", IsMaker: true, IsChecker: true},
					},
				})
				var found bool
				aclAdmin, found = f.aclKeeper.GetAclAuthority(ctx, "address")
				assert.Assert(t, found == true)
				assert.Assert(t, aclAdmin.String() != "")

				req = &types.QueryGetAclAuthorityRequest{Address: "address"}

				res = &types.QueryGetAclAuthorityResponse{
					AclAuthority: types.AclAuthority{
						Address: "address",
						Name:    "Alice",
						AccessDefinitions: []*types.AccessDefinition{
							{Module: "module", IsMaker: true, IsChecker: true},
						},
					},
				}
			},
			true,
			"",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Case %s", testCase.msg), func(t *testing.T) {
			testCase.malleate()

			aclAdmin, err := queryClient.AclAuthority(gocontext.Background(), req)

			if testCase.expPass {
				assert.NilError(t, err)
				assert.Equal(t, res.String(), aclAdmin.String())
			} else {
				assert.ErrorContains(t, err, testCase.expErrMsg)
				assert.Assert(t, aclAdmin == nil)
			}
		})
	}
}

func TestGRPCQueryAllAclAuthority(t *testing.T) {
	f := initFixture(t)
	ctx, queryClient := f.ctx, f.queryClient

	var (
		req             *types.QueryAllAclAuthorityRequest
		res             *types.QueryAllAclAuthorityResponse
		aclAuthorityAll []types.AclAuthority
	)

	testCases := []struct {
		msg       string
		malleate  func()
		expPass   bool
		expErrMsg string
	}{
		{
			"nil request",
			func() {
				req = nil
				res = &types.QueryAllAclAuthorityResponse{
					AclAuthority: []types.AclAuthority{},
					Pagination:   &query.PageResponse{},
				}
			},
			true,
			"",
		},
		{
			"get all acl authorities",
			func() {
				f.aclKeeper.SetAclAuthority(ctx, types.AclAuthority{
					Address: "address1",
					Name:    "Name1",
					AccessDefinitions: []*types.AccessDefinition{
						{Module: "module1", IsMaker: true, IsChecker: false},
					},
				})
				f.aclKeeper.SetAclAuthority(ctx, types.AclAuthority{
					Address: "address2",
					Name:    "Name2",
					AccessDefinitions: []*types.AccessDefinition{
						{Module: "module2", IsMaker: false, IsChecker: true},
					},
				})
				f.aclKeeper.SetAclAuthority(ctx, types.AclAuthority{
					Address: "address3",
					Name:    "Name3",
					AccessDefinitions: []*types.AccessDefinition{
						{Module: "module3", IsMaker: true, IsChecker: false},
					},
				})
				f.aclKeeper.SetAclAuthority(ctx, types.AclAuthority{
					Address: "address4",
					Name:    "Name4",
					AccessDefinitions: []*types.AccessDefinition{
						{Module: "module4", IsMaker: false, IsChecker: true},
					},
				})

				aclAuthorityAll = f.aclKeeper.GetAllAclAuthority(ctx)
				assert.Assert(t, len(aclAuthorityAll) == 4)

				req = &types.QueryAllAclAuthorityRequest{}

				res = &types.QueryAllAclAuthorityResponse{
					AclAuthority: []types.AclAuthority{
						{Address: "address1"},
						{Address: "address2"},
						{Address: "address3"},
						{Address: "address4"},
					},
					Pagination: &query.PageResponse{
						Total: 4,
					},
				}
			},
			true,
			"",
		},
		{
			"get some of acl authorities",
			func() {
				aclAuthorityAll = f.aclKeeper.GetAllAclAuthority(ctx)
				assert.Assert(t, len(aclAuthorityAll) == 4)

				req = &types.QueryAllAclAuthorityRequest{
					Pagination: &query.PageRequest{
						Limit: 2,
					},
				}

				res = &types.QueryAllAclAuthorityResponse{
					AclAuthority: []types.AclAuthority{
						{
							Address: "address1",
						},
						{
							Address: "address2",
						},
					},
				}
			},
			true,
			"",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Case %s", testCase.msg), func(t *testing.T) {
			testCase.malleate()

			aclAuthorities, err := queryClient.AclAuthorityAll(gocontext.Background(), req)

			if testCase.expPass {
				assert.NilError(t, err)
				assert.Equal(t, len(res.AclAuthority), len(aclAuthorities.AclAuthority))
			} else {
				assert.ErrorContains(t, err, testCase.expErrMsg)
				assert.Assert(t, aclAuthorities == nil)
			}
		})
	}
}
