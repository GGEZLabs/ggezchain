package types

func NewMsgDeleteAuthority(creator string, authAddress string) *MsgDeleteAuthority {
	return &MsgDeleteAuthority{
		Creator:     creator,
		AuthAddress: authAddress,
	}
}
