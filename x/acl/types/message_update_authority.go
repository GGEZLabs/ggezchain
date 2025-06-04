package types

func NewMsgUpdateAuthority(creator string, authAddress string, newName string, overwriteAccessDefinitions string, addAccessDefinitions string, updateAccessDefinition string, deleteAccessDefinitions []string, clearAllAccessDefinitions bool) *MsgUpdateAuthority {
	return &MsgUpdateAuthority{
		Creator:                    creator,
		AuthAddress:                authAddress,
		NewName:                    newName,
		OverwriteAccessDefinitions: overwriteAccessDefinitions,
		AddAccessDefinitions:       addAccessDefinitions,
		UpdateAccessDefinition:     updateAccessDefinition,
		DeleteAccessDefinitions:    deleteAccessDefinitions,
		ClearAllAccessDefinitions:  clearAllAccessDefinitions,
	}
}
