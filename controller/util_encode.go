package controller

import (
	"crypto/md5"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GetMD5(raw string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
}

func GetBcrypt(raw string) string {
	password := []byte(raw)
	hashedPwd, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPwd)
}

func MatchBcrypt(raw string, bcryptStr string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(bcryptStr), []byte(raw))

	return err == nil
}
