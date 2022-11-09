package password

import "golang.org/x/crypto/bcrypt"

func Hash(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return bytes, err
}

func CheckHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
