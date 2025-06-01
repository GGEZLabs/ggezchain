package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/acl module sentinel errors
var (
	ErrInvalidSigner                 = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrEmptyName                     = sdkerrors.Register(ModuleName, 1101, "empty name not allowed")
	ErrInvalidAccessDefinitionList   = sdkerrors.Register(ModuleName, 1102, "invalid access definition list format")
	ErrInvalidAccessDefinitionObject = sdkerrors.Register(ModuleName, 1103, "invalid access definition object format")
	ErrInvalidModuleName             = sdkerrors.Register(ModuleName, 1104, "invalid module name")
	ErrUnauthorized                  = sdkerrors.Register(ModuleName, 1105, "unauthorized account")
	ErrAuthorityAddressExist         = sdkerrors.Register(ModuleName, 1106, "authority address already exist")
	ErrAuthorityAddressNotExist      = sdkerrors.Register(ModuleName, 1107, "authority address not exist")
	ErrModuleNotExist                = sdkerrors.Register(ModuleName, 1108, "module not exist")
	ErrModuleExist                   = sdkerrors.Register(ModuleName, 1109, "module already exist")
	ErrEmptyAccessDefinitionList     = sdkerrors.Register(ModuleName, 1110, "access definition list is empty")
	ErrNoUpdateFlags                 = sdkerrors.Register(ModuleName, 1111, "at least one update flag must be provided")
	ErrUpdateAndRemoveModule         = sdkerrors.Register(ModuleName, 1113, "module(s) cannot be both added/updated and removed in the same request")
	ErrRequireMakerOrChecker         = sdkerrors.Register(ModuleName, 1114, "at least one of is_maker or is_checker must be true")
	ErrSuperAdminInitialized         = sdkerrors.Register(ModuleName, 1115, "super admin already initialized")
	ErrAdminExist                    = sdkerrors.Register(ModuleName, 1116, "admin(s) already exist")
	ErrAllAdminDeletion              = sdkerrors.Register(ModuleName, 1117, "cannot delete all admins, at least one admin must remain")
	ErrAdminNotExist                 = sdkerrors.Register(ModuleName, 1118, "admin(s) not exist")
)
