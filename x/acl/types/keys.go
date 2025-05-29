package types

const (
	// ModuleName defines the module name
	ModuleName = "acl"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_acl"
)

var ParamsKey = []byte("p_acl")

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	SuperAdminKey = "SuperAdmin/value/"
)
