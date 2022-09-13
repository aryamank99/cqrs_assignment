package models

type LoginCredentials struct {
	ID           string `json:"_id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Salt         string `json:"salt"`
}
