package password

import "testing"

func TestPassword(t *testing.T) {
	pw := "testtt"
	hash, err := Hash(pw)
	if err != nil {
		t.Error(err)
	}
	t.Log(hash, len(hash))
	ok := CheckHash(pw, hash)
	t.Log(ok)
}
