package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func Descrypt(password string, hashPassword string) (bool, error){
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

	if err != nil {
		return false, err
	}

	return true, err
}