package types

func NewMsgInit(creator string, superAdmin string) *MsgInit {
	return &MsgInit{
		Creator:    creator,
		SuperAdmin: superAdmin,
	}
}
