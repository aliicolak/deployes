package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword: kullanıcı şifresini bcrypt ile hashler
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash: plain şifre ile hashed şifreyi karşılaştırır
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
