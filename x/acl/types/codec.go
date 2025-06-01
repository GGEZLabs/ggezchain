package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddAuthority{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeleteAuthority{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateAuthority{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgInit{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddAdmin{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeleteAdmin{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateSuperAdmin{},
	)
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
