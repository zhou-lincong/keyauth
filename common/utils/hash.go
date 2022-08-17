package utils

import (
	"crypto/sha1"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Hash(data string) string {
	if data == "" {
		return ""
	}

	h := sha1.New()
	_, err := h.Write([]byte(data))
	if err != nil {
		panic(err)
	}

	hd := h.Sum(nil)
	return fmt.Sprintf("%x", hd)
}
