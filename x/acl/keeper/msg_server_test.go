package keeper_test

import (
	"testing"

	"github.com/GGEZLabs/ggezchain/v2/x/acl/keeper"
	"github.com/stretchr/testify/require"
)

func TestMsgServer(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)
	require.NotNil(t, ms)
	require.NotNil(t, f.ctx)
	require.NotEmpty(t, f.keeper)
}
