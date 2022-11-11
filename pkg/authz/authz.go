package authz

import (
	"auth/pkg/access"
	"fmt"
)

const (
	NMax = 64
)

type Authorizer struct {
	names  map[string]access.Access
	values map[access.Access]string
}

func New(permissions ...string) Authorizer {
	am := Authorizer{
		names:  make(map[string]access.Access),
		values: make(map[access.Access]string),
	}

	for i, permission := range permissions {
		if i >= NMax {
			panic("uint64 overflow")
		}
		a := access.Access(1) << i
		am.names[permission] = a
		am.values[a] = permission
	}

	return am
}

func (a Authorizer) Access(permissions ...string) access.Access {
	ac := access.Access(0)
	for _, permission := range permissions {
		tmp, ok := a.names[permission]
		if !ok {
			panic(fmt.Sprintf("unknown permission: %s", permission))
		}
		ac = ac | tmp
	}
	return ac
}

func (a Authorizer) ByName() map[string]access.Access {
	return a.names
}

func (a Authorizer) ByAccess() map[access.Access]string {
	return a.values
}
