package users

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPwd hashes a password.
func HashPwd(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// HintPwd hint a password.
func HintPwd(password string) string {
	str := password
	if len(str) >= 7 {
		star := ""
		for i := 1; i <= len(str)-4; i++ {
			star += "*"
		}
		return str[0:2] + star + str[len(str)-2:]
	} else {
		return "INVALID_PASSWORD"
	}
	return str
}

// CheckPwd checks if a password is correct.
func CheckPwd(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
