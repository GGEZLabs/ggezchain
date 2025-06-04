package types

func NewMsgAddAdmin(creator string, admins string) *MsgAddAdmin {
	return &MsgAddAdmin{
		Creator: creator,
		Admins:  admins,
	}
}
