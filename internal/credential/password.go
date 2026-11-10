package credential

import (
	"unicode"

	"github.com/alexedwards/argon2id"
)

var hasUpper, hasSpecialCharacter, hasNum = false, false, false

func Hash(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, err
}

func Verify(password string, hashedPassword string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hashedPassword)
	if err != nil {
		return false, err
	}
	return match, err
}

func FormatCheck(password string) bool {
	// match, err := regexp.MatchString(password, passwordRegex)
	// if err != nil {
	// 	return !match
	// }
	// return match
	if len(password) < 8 {
		return false
	}

	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsNumber(r) {
			hasNum = true
		}
		if unicode.IsSymbol(r) {
			hasSpecialCharacter = true
		}
	}

	return hasUpper && hasNum && hasSpecialCharacter
}
