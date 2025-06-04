package types

func NewMsgAddAuthority(creator string, authAddress string, name string, accessDefinitions string) *MsgAddAuthority {
	return &MsgAddAuthority{
		Creator:           creator,
		AuthAddress:       authAddress,
		Name:              name,
		AccessDefinitions: accessDefinitions,
	}
}
