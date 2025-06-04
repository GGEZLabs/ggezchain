package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registrar codectypes.InterfaceRegistry) {

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateAuthority{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeleteAuthority{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddAuthority{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateSuperAdmin{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeleteAdmin{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddAdmin{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgInit{},
	)

	registrar.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registrar, &_Msg_serviceDesc)
}
