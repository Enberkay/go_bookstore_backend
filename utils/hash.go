package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword แปลง password เป็น hash โดยใช้ bcrypt
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}
	return string(bytes), nil
}

// CheckPassword ตรวจสอบ password กับ hashed password
func CheckPassword(hashedPassword, password string) error {
	if hashedPassword == "" || password == "" {
		return errors.New("password and hash cannot be empty")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}
