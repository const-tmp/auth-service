package access

import "fmt"

//go:generate stringer -type Access

type Access uint64

func (a Access) Check(a2 Access) bool {
	return a&a2 == a
}

const (
	Active Access = 1 << iota
	SelfLedger
	ReadLedger
	WriteLedger
)

const (
	NMax = 64
)

type Helper struct {
	names  map[string]Access
	values map[Access]string
}

func NewHelperFromPermissions(permissions ...string) Helper {
	am := Helper{
		names:  make(map[string]Access),
		values: make(map[Access]string),
	}

	for i, permission := range permissions {
		if i >= NMax {
			panic("uint64 overflow")
		}
		a := Access(1) << i
		am.names[permission] = a
		am.values[a] = permission
	}

	return am
}

func (h Helper) Access(permissions ...string) Access {
	ac := Access(0)
	for _, permission := range permissions {
		tmp, ok := h.names[permission]
		if !ok {
			panic(fmt.Sprintf("unknown permission: %s", permission))
		}
		ac = ac | tmp
	}
	return ac
}

func (h Helper) Permissions(a Access) (res []string) {
	for name, access := range h.names {
		if access.Check(a) {
			res = append(res, name)
		}
	}
	return
}

func (h Helper) AllPermissions() (res []string) {
	for name, _ := range h.names {
		res = append(res, name)
	}
	return
}

func (h Helper) ByName() map[string]Access {
	return h.names
}

func (h Helper) ByAccess() map[Access]string {
	return h.values
}
