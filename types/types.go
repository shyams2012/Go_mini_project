package types

import (
	"time"

	"github.com/square/go-jose/v3/jwt"
	"gorm.io/gorm"
)

type UserClaims struct {
	jwt.Claims
	Email string `json:"email"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type AuthKey struct {
	Email          string    `json:"email"`
	PrivateKey     string    `json:"privateKey"`
	PublicKey      string    `json:"publicKey"`
	ExpirationTime time.Time `json:"expirationTime"`
}

type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}
