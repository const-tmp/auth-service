package access

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
