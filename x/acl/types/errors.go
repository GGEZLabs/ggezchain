package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/acl module sentinel errors
var (
	ErrInvalidSigner                 = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrEmptyName                     = sdkerrors.Register(ModuleName, 1101, "address name is required and cannot be empty")
	ErrInvalidAccessDefinitionList   = sdkerrors.Register(ModuleName, 1102, "invalid access definition list format")
	ErrInvalidAccessDefinitionObject = sdkerrors.Register(ModuleName, 1103, "invalid access definition object format")
	ErrInvalidModuleName             = sdkerrors.Register(ModuleName, 1104, "invalid module name")
	ErrUnauthorized                  = sdkerrors.Register(ModuleName, 1105, "unauthorized account")
	ErrAuthorityAddressExists        = sdkerrors.Register(ModuleName, 1106, "authority address already exists")
	ErrAuthorityAddressDoesNotExist  = sdkerrors.Register(ModuleName, 1107, "authority address does not exist")
	ErrModuleDoesNotExist            = sdkerrors.Register(ModuleName, 1108, "module name does not exist")
	ErrModuleExists                  = sdkerrors.Register(ModuleName, 1109, "module name already exists")
	ErrEmptyAccessDefinitionList     = sdkerrors.Register(ModuleName, 1110, "access definition list is required and cannot be empty")
	ErrNoUpdateFlags                 = sdkerrors.Register(ModuleName, 1111, "at least one update flag must be provided")
	ErrUpdateAndRemoveModule         = sdkerrors.Register(ModuleName, 1113, "same module(s) cannot be both added/updated and removed in the same request")
	ErrSuperAdminInitialized         = sdkerrors.Register(ModuleName, 1114, "super admin already initialized")
	ErrAdminExist                    = sdkerrors.Register(ModuleName, 1115, "admin(s) already exists")
	ErrAdminNotExist                 = sdkerrors.Register(ModuleName, 1116, "admin(s) does not exist")
)
