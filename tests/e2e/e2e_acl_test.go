package e2e

import (
	"fmt"
	"time"
)

func (s *IntegrationTestSuite) testAcl() {
	chainEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))

	superAdmin, err := s.chainA.genesisAccounts[0].keyInfo.GetAddress()
	s.Require().NoError(err)

	admin1, err := s.chainA.genesisAccounts[0].keyInfo.GetAddress()
	s.Require().NoError(err)

	admin2, err := s.chainA.genesisAccounts[1].keyInfo.GetAddress()
	s.Require().NoError(err)

	admin3, err := s.chainA.genesisAccounts[2].keyInfo.GetAddress()
	s.Require().NoError(err)

	// Init ACL module
	s.execInitAcl(s.chainA, 0, superAdmin.String(), superAdmin.String(), ggezHomePath, standardFees.String())
	s.Require().Eventually(
		func() bool {
			res, err := querySuperAdmin(chainEndpoint)
			s.Require().NoError(err)

			return res.GetSuperAdmin().Admin == superAdmin.String()
		},
		20*time.Second,
		5*time.Second,
	)

	// Add admin by super admin
	s.execAddAclAdmin(s.chainA, 0, admin1.String(), superAdmin.String(), ggezHomePath, standardFees.String())

	s.Require().Eventually(
		func() bool {
			res, err := queryAclAdmin(chainEndpoint, admin1.String())
			s.Require().NoError(err)

			return res.GetAclAdmin().Address == admin1.String()
		},
		20*time.Second,
		5*time.Second,
	)

	// Add another admin by super admin
	s.execAddAclAdmin(s.chainA, 0, admin2.String(), superAdmin.String(), ggezHomePath, standardFees.String())

	s.Require().Eventually(
		func() bool {
			aclAdmin, err := queryAclAdmin(chainEndpoint, admin2.String())
			s.Require().NoError(err)

			aclAdminAll, err := queryAllAclAdmin(chainEndpoint)
			s.Require().NoError(err)

			return aclAdmin.GetAclAdmin().Address == admin2.String() && len(aclAdminAll.AclAdmin) == 2
		},
		20*time.Second,
		5*time.Second,
	)

	// Add admin3 by super admin
	s.execAddAclAdmin(s.chainA, 0, admin3.String(), superAdmin.String(), ggezHomePath, standardFees.String())

	s.Require().Eventually(
		func() bool {
			aclAdmin, err := queryAclAdmin(chainEndpoint, admin3.String())
			s.Require().NoError(err)

			aclAdminAll, err := queryAllAclAdmin(chainEndpoint)
			s.Require().NoError(err)

			return aclAdmin.GetAclAdmin().Address == admin3.String() && len(aclAdminAll.AclAdmin) == 3
		},
		20*time.Second,
		5*time.Second,
	)

	// Delete admin3
	s.execDeleteAclAdmin(s.chainA, 0, admin3.String(), superAdmin.String(), ggezHomePath, standardFees.String())

	s.Require().Eventually(
		func() bool {
			aclAdminAll, err := queryAllAclAdmin(chainEndpoint)
			s.Require().NoError(err)

			return len(aclAdminAll.AclAdmin) == 2
		},
		20*time.Second,
		5*time.Second,
	)

	// Add authority for super admin
	accessDefinitions := `[{"module":"trade","is_maker":true,"is_checker":true},{"module":"acl","is_maker":true,"is_checker":true}]`
	s.execAddAclAuthority(s.chainA, 0, superAdmin.String(), "Alice", accessDefinitions, admin1.String(), ggezHomePath, standardFees.String())

	s.Require().Eventually(
		func() bool {
			aclAuthority, err := queryAclAuthority(chainEndpoint, superAdmin.String())
			s.Require().NoError(err)

			aclAuthorityAll, err := queryAllAclAuthority(chainEndpoint)
			s.Require().NoError(err)

			return len(aclAuthorityAll.AclAuthority) == 1 && aclAuthority.GetAclAuthority().Address == superAdmin.String()
		},
		20*time.Second,
		5*time.Second,
	)

	// Add authority for admin2
	accessDefinitions = `[{"module":"trade","is_maker":false,"is_checker":true}]`
	s.execAddAclAuthority(s.chainA, 0, admin2.String(), "Bob", accessDefinitions, admin1.String(), ggezHomePath, standardFees.String())

	s.Require().Eventually(
		func() bool {
			aclAuthority, err := queryAclAuthority(chainEndpoint, admin2.String())
			s.Require().NoError(err)

			aclAuthorityAll, err := queryAllAclAuthority(chainEndpoint)
			s.Require().NoError(err)

			return len(aclAuthorityAll.AclAuthority) == 2 && aclAuthority.GetAclAuthority().Address == admin2.String()
		},
		20*time.Second,
		5*time.Second,
	)

	// Add authority for admin3
	accessDefinitions = `[{"module":"trade","is_maker":true,"is_checker":true}]`
	s.execAddAclAuthority(s.chainA, 0, admin3.String(), "Carol", accessDefinitions, admin1.String(), ggezHomePath, standardFees.String())

	s.Require().Eventually(
		func() bool {
			aclAuthority, err := queryAclAuthority(chainEndpoint, admin3.String())
			s.Require().NoError(err)

			aclAuthorityAll, err := queryAllAclAuthority(chainEndpoint)
			s.Require().NoError(err)

			return len(aclAuthorityAll.AclAuthority) == 3 && aclAuthority.GetAclAuthority().Address == admin3.String()
		},
		20*time.Second,
		5*time.Second,
	)

	// Delete admin3 authority
	s.execDeleteAclAuthority(s.chainA, 0, admin3.String(), admin1.String(), ggezHomePath, standardFees.String())

	s.Require().Eventually(
		func() bool {
			aclAuthorityAll, err := queryAllAclAuthority(chainEndpoint)
			s.Require().NoError(err)

			return len(aclAuthorityAll.AclAuthority) == 2
		},
		20*time.Second,
		5*time.Second,
	)

	// update authority
	updateAccessDefinitions := `{"module":"trade","is_maker":true,"is_checker":false}`
	s.execUpdateAclAuthority(s.chainA, 0, superAdmin.String(), "acl", updateAccessDefinitions, admin1.String(), ggezHomePath, standardFees.String())

	s.Require().Eventually(
		func() bool {
			res, err := queryAclAuthority(chainEndpoint, superAdmin.String())
			s.Require().NoError(err)
			s.Require().Equal("trade", res.AclAuthority.AccessDefinitions[0].Module)
			s.Require().Equal(true, res.AclAuthority.AccessDefinitions[0].IsMaker)
			s.Require().Equal(false, res.AclAuthority.AccessDefinitions[0].IsChecker)
			return len(res.AclAuthority.AccessDefinitions) == 1
		},
		20*time.Second,
		5*time.Second,
	)
}
