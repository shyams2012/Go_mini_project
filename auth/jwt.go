package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/shyams2012/Go_mini_project/types"
	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
)

const (
	JWTIssuer    = "ShyamSah" // The issuer of JWT
	JWTAddExpiry = time.Minute * 3000
)

var (
	ErrInvalidJWT error = fmt.Errorf("invalid jwt")
)

func NewJWT(email string, priKey *rsa.PrivateKey, expireAt time.Time) (string, error) {
	claims := types.UserClaims{
		Email: email,
		Claims: jwt.Claims{
			Issuer: JWTIssuer,
			Expiry: jwt.NewNumericDate(expireAt),
		},
	}
	opts := jose.SignerOptions{}
	opts.WithType("JWT")

	signKey := jose.SigningKey{
		Algorithm: jose.RS256,
		Key:       priKey,
	}

	signer, err := jose.NewSigner(signKey, &opts)
	if err != nil {
		return "", err
	}

	return jwt.Signed(signer).
		Claims(claims).
		CompactSerialize()
}

func ParseJWT(signedJWT string, pubKey *rsa.PublicKey) (*types.UserClaims, error) {
	token, err := jwt.ParseSigned(signedJWT)
	if err != nil {
		return nil, fmt.Errorf("invalid jwt")
	}

	claims := new(types.UserClaims)
	if err := token.Claims(pubKey, claims); err != nil {
		return nil, fmt.Errorf("invalid jwt")
	}

	err = claims.Validate(jwt.Expected{
		Issuer: JWTIssuer,
		Time:   time.Now(),
	})
	if err != nil {
		if err == jwt.ErrExpired {
			return nil, ErrInvalidJWT
		}

		return nil, ErrInvalidJWT
	}

	return claims, nil
}

func GetUnsafeClaims(signedJWT string) (*types.UserClaims, error) {
	token, err := jwt.ParseSigned(signedJWT)
	if err != nil {
		return nil, fmt.Errorf("invalid jwt")
	}

	claims := new(types.UserClaims)

	if err := token.UnsafeClaimsWithoutVerification(claims); err != nil {
		return nil, fmt.Errorf("could not get unsafe claims")
	}

	return claims, nil
}
