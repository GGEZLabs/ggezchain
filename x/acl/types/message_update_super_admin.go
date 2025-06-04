package types

func NewMsgUpdateSuperAdmin(creator string, newSuperAdmin string) *MsgUpdateSuperAdmin {
	return &MsgUpdateSuperAdmin{
		Creator:       creator,
		NewSuperAdmin: newSuperAdmin,
	}
}
