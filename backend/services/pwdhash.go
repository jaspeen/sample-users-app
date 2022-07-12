package services

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcrypt_cost = 14
)

func HashPasword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt_cost)
	return string(hash), err
}

func CheckPassword(passHash string, password string) bool {
	log.Debug().Msgf("CheckPass: %v %v", passHash, password)
	err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password))
	return err == nil
}
