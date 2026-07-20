package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/acl module sentinel errors
var (
	ErrInvalidSigner                 = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrEmptyName                     = errors.Register(ModuleName, 1101, "address name is required and cannot be empty")
	ErrInvalidAccessDefinitionList   = errors.Register(ModuleName, 1102, "invalid access definition list format")
	ErrInvalidAccessDefinitionObject = errors.Register(ModuleName, 1103, "invalid access definition object format")
	ErrInvalidModuleName             = errors.Register(ModuleName, 1104, "invalid module name")
	ErrUnauthorized                  = errors.Register(ModuleName, 1105, "unauthorized account")
	ErrAuthorityAddressExists        = errors.Register(ModuleName, 1106, "authority address already exists")
	ErrAuthorityAddressDoesNotExist  = errors.Register(ModuleName, 1107, "authority address does not exist")
	ErrModuleDoesNotExist            = errors.Register(ModuleName, 1108, "module name does not exist")
	ErrModuleExists                  = errors.Register(ModuleName, 1109, "module name already exists")
	ErrEmptyAccessDefinitionList     = errors.Register(ModuleName, 1110, "access definition list is required and cannot be empty")
	ErrNoUpdateFlags                 = errors.Register(ModuleName, 1111, "at least one update flag must be provided")
	ErrUpdateAndRemoveModule         = errors.Register(ModuleName, 1113, "same module(s) cannot be both added/updated and removed in the same request")
	ErrSuperAdminInitialized         = errors.Register(ModuleName, 1114, "super admin already initialized")
	ErrAdminExist                    = errors.Register(ModuleName, 1115, "admin(s) already exists")
	ErrAdminNotExist                 = errors.Register(ModuleName, 1116, "admin(s) does not exist")
)
