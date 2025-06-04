package types

func NewMsgDeleteAdmin(creator string, admins string) *MsgDeleteAdmin {
	return &MsgDeleteAdmin{
		Creator: creator,
		Admins:  admins,
	}
}
